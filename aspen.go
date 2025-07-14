package main

import (
	"aspen/config"
	"aspen/logging"
	"aspen/middleware"
	"aspen/resources"
	"aspen/router"
	"aspen/router/service"
	"flag"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var serverPort = flag.Int("port", 8080, "the port to open this server on")
var serviceFolder = flag.String("services", "./services", "the folder to place service files in")

func main() {
	// Init
	logging.InitializeLogger(zerolog.DebugLevel)
	logging.AddConsoleOutput(true)
	middleware.RegisterMiddleware()
	resources.RegisterResources()
	flag.Parse()

	// Set service folder
	service.SetGlobalFolder(*serviceFolder)

	// Parse config path
	if len(flag.Args()) == 0 {
		log.Fatal().Msg("Error: Configuration file path is required. Usage: go run ./aspen.go [flags] <config-file>")
	}
	configPath := flag.Args()[0]
	err := config.SetGlobalConfigFile(configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Error setting global config file")
	}

	// Load config
	instance, err := config.ParseGlobalConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading config")
	}

	// Start instance
	err = instance.BuildAndStartServices()
	if err != nil {
		log.Fatal().Err(err).Msg("Error starting services")
	}

	// Init router
	router.UpdateRouter(instance)

	// Start server
	log.Info().Int("port", *serverPort).Msg("Starting server")
	err = http.ListenAndServe(fmt.Sprintf(":%d", *serverPort), &router.GlobalRouter)
	log.Fatal().Err(err).Msg("Server closed")
}
