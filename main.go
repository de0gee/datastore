package main

import (
	"flag"

	"github.com/de0gee/de0gee-data/src/database"
	"github.com/de0gee/de0gee-data/src/server"
)

func main() {
	port := flag.String("port", "8003", "port for the server")
	debug := flag.Bool("debug", false, "turn on debug mode")
	flag.Parse()
	if *debug {
		database.Debug(true)
		server.Debug(true)
	}
	server.Run(*port)
}
