package main

import (
	"aspen/config"
	"aspen/logging"
	"aspen/middleware"
	"aspen/resources"
	"aspen/router"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var serverPort = flag.Int("port", 8080, "the port to open this server on")

func main() {
	// Init
	logging.InitializeLogger(zerolog.InfoLevel)
	logging.AddConsoleOutput(true)
	resources.RegisterResources()
	flag.Parse()

	// Parse config path
	if len(flag.Args()) == 0 {
		log.Fatal().Msg("Error: Configuration file path is required. Usage: go run ./aspen.go [flags] <config-file>")
	}
	configPath := flag.Args()[0]

	// Load config
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading file")
	}

	config, err := config.ParseJSON(data)
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing JSON")
	}

	resource_routes, err := config.GetResourceRoutes()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading routes")
	}

	// Init router
	router.UpdateRouter(router.NewRouterInstance(
		[]router.Middleware{middleware.Logger{}},
		resource_routes,
	))

	// Start server
	log.Info().Int("port", *serverPort).Msg("Starting server")
	err = http.ListenAndServe(fmt.Sprintf(":%d", *serverPort), &router.GlobalRouter)
	log.Fatal().Err(err).Msg("Server closed")
}
