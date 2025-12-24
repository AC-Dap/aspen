package main

import (
	"aspen/auth"
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
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var serverPort = flag.Int("port", 8080, "The port to open this server on.")
var serviceFolder = flag.String("services", "./services", "The folder to place service files in.")
var authFile = flag.String("auth", "auth.sqlite", "The file to use for authentication storage.")

// Logging flags
var logLevel = flag.String("log-level", "DEBUG", `The logging level to use. From highest to lowest, the levels are:
	PANIC, FATAL, ERROR, WARN, INFO, DEBUG, TRACE, DISABLED`)
var loggingOutput = flag.String("out", "", "The output file to log to. Logs to stderr if not provided.")

func parseLoggingLevelString() {
	switch strings.ToUpper(*logLevel) {
	case "PANIC":
		logging.InitializeLogger(zerolog.PanicLevel)
	case "FATAL":
		logging.InitializeLogger(zerolog.FatalLevel)
	case "ERROR":
		logging.InitializeLogger(zerolog.ErrorLevel)
	case "WARN":
		logging.InitializeLogger(zerolog.WarnLevel)
	case "INFO":
		logging.InitializeLogger(zerolog.InfoLevel)
	case "DEBUG":
		logging.InitializeLogger(zerolog.DebugLevel)
	case "TRACE":
		logging.InitializeLogger(zerolog.TraceLevel)
	case "DISABLED":
		logging.DisableLogger()
	default:
		log.Fatal().Str("log-level", *logLevel).Msg("Invalid log level.")
	}
}

func main() {
	flag.Parse()

	// Set up logger before anything else
	if *loggingOutput == "" {
		logging.SetOutputToConsole()
	} else {
		err := logging.SetOutputToFile(*loggingOutput)
		if err != nil {
			log.Fatal().Err(err).Msg("Error opening log file")
		}
	}
	parseLoggingLevelString()

	middleware.RegisterMiddleware()
	resources.RegisterResources()

	err := auth.Initialize(*authFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing authentication")
	} else {
		log.Info().Msg("Authentication initialized successfully")
	}

	service.SetGlobalFolder(*serviceFolder)

	// Parse config path
	if len(flag.Args()) == 0 {
		log.Fatal().Msg("Error: Configuration file path is required. Usage: go run ./aspen.go [flags] <config-file>")
	}
	configPath := flag.Args()[0]
	err = config.SetGlobalConfigFile(configPath)
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
