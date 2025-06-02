package main

import (
	"aspen/config"
	"aspen/resources"
	"aspen/router"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var serverPort = flag.Int("port", 8080, "the port to open this server on")

func main() {
	// Init
	resources.RegisterResources()
	flag.Parse()

	// Parse config path
	if len(flag.Args()) == 0 {
		log.Fatal("Error: Configuration file path is required. Usage: go run ./aspen.go [flags] <config-file>")
	}
	configPath := flag.Args()[0]

	// Load config
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal("Error reading file: ", err)
	}

	config, err := config.ParseJSON(data)
	if err != nil {
		log.Fatal("Error parsing JSON: ", err)
	}

	resource_routes, err := config.ToResourceRoutes()
	if err != nil {
		log.Fatal("Error loading routes: ", err)
	}

	// Init router
	router.UpdateRouter(resource_routes)

	// Start server
	log.Printf("Starting server on port %d", *serverPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *serverPort), &router.GlobalRouter))
}
