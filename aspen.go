package main

import (
	"aspen/config"
	"aspen/logging"
	"aspen/middleware"
	"aspen/resources"
	"aspen/router"
	"aspen/router/service"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	// Add handler for ctrl-c shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Start server
	log.Info().Int("port", *serverPort).Msg("Starting server")
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *serverPort),
		Handler: &router.GlobalRouter,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	// Wait for signal
	sig := <-quit
	log.Info().Str("signal", sig.String()).Msg("Received shutdown signal")

	// Shutdown server and services
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error shutting down server")
	}

	err = router.GlobalRouter.Shutdown()
	if err != nil {
		log.Error().Err(err).Msg("Error stopping services")
	}

	log.Info().Msg("Server shutdown complete")
}
