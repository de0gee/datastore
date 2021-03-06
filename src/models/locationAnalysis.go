package models

type LocationAnalysis struct {
	LocationNames map[string]string `json:"location_names"`
	Predictions   []struct {
		Locations     []string  `json:"locations"`
		Name          string    `json:"name"`
		Probabilities []float64 `json:"probabilities"`
	} `json:"predictions"`
	BestGuess LocationPrediction `json:"best_guess",omitempty`
}

type LocationPrediction struct {
	Location    string  `json:"location",omitempty`
	Name        string  `json:"name",omitempty`
	Probability float64 `json:"probability",omitempty`
}
