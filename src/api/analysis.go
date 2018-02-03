package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/de0gee/de0gee-data/src/database"
	"github.com/de0gee/de0gee-data/src/models"
	cache "github.com/robfig/go-cache"
)

// AIPort designates the port for the AI processing
var AIPort = "8002"

var (
	httpClient *http.Client
	routeCache *cache.Cache
)

const (
	MaxIdleConnections int = 20
	RequestTimeout     int = 5
)

// init HTTPClient
func init() {
	httpClient = createHTTPClient()
	routeCache = cache.New(5*time.Minute, 10*time.Minute)
}

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
		Timeout: time.Duration(RequestTimeout) * time.Second,
	}

	return client
}

type AnalysisResponse struct {
	Data    models.LocationAnalysis `json:"analysis"`
	Message string                  `json:"message"`
	Success bool                    `json:"success"`
}

func AnalyzeSensorData(s models.SensorData) (aidata models.LocationAnalysis, err error) {
	d, err := database.Open(s.Family)
	if err != nil {
		return
	}
	defer d.Close()

	// check if its already been classified
	// aidata, err = d.GetPrediction(s.Timestamp)
	// if err == nil {
	// 	return
	// }

	// inquire the AI
	var target AnalysisResponse
	type ClassifyPayload struct {
		Sensor     models.SensorData `json:"sensor_data"`
		DataFolder string            `json:"data_folder"`
	}
	var p2 ClassifyPayload
	p2.Sensor = s
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	p2.DataFolder = dir
	url := "http://localhost:" + AIPort + "/classify"
	bPayload, err := json.Marshal(p2)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bPayload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&target)
	if err != nil {
		return
	}
	if !target.Success {
		err = errors.New("unable to analyze: " + target.Message)
		return
	}
	if len(target.Data.Predictions) == 0 {
		err = errors.New("problem analyzing: no predictions")
		return
	}

	aidata = target.Data
	// add prediction to the database
	err = d.AddPrediction(s.Timestamp, aidata)
	if err != nil {
		return
	}

	return
}
