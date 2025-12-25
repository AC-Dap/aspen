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
)

var lg = logging.NewTaggedLogger("Main")

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
		logging.SetLoggingLevel(zerolog.PanicLevel)
	case "FATAL":
		logging.SetLoggingLevel(zerolog.FatalLevel)
	case "ERROR":
		logging.SetLoggingLevel(zerolog.ErrorLevel)
	case "WARN":
		logging.SetLoggingLevel(zerolog.WarnLevel)
	case "INFO":
		logging.SetLoggingLevel(zerolog.InfoLevel)
	case "DEBUG":
		logging.SetLoggingLevel(zerolog.DebugLevel)
	case "TRACE":
		logging.SetLoggingLevel(zerolog.TraceLevel)
	case "DISABLED":
		logging.SetLoggingLevel(zerolog.Disabled)
	default:
		lg.Fatal().Str("log-level", *logLevel).Msg("Invalid log level.")
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
			lg.Fatal().Err(err).Msg("Error opening log file")
		}
	}
	parseLoggingLevelString()

	middleware.RegisterMiddleware()
	resources.RegisterResources()

	err := auth.Initialize(*authFile)
	if err != nil {
		lg.Fatal().Err(err).Msg("Error initializing authentication")
	} else {
		lg.Info().Msg("Authentication initialized successfully")
	}

	service.SetGlobalFolder(*serviceFolder)

	// Parse config path
	if len(flag.Args()) == 0 {
		lg.Fatal().Msg("Error: Configuration file path is required. Usage: go run ./aspen.go [flags] <config-file>")
	}
	configPath := flag.Args()[0]
	err = config.SetGlobalConfigFile(configPath)
	if err != nil {
		lg.Fatal().Err(err).Msg("Error setting global config file")
	}

	// Load config
	instance, err := config.ParseGlobalConfig()
	if err != nil {
		lg.Fatal().Err(err).Msg("Error loading config")
	}

	// Start instance
	err = instance.BuildAndStartServices()
	if err != nil {
		lg.Fatal().Err(err).Msg("Error starting services")
	}

	// Init router
	router.UpdateRouter(instance)

	// Add handler for ctrl-c shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Start server
	lg.Info().Int("port", *serverPort).Msg("Starting server")
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *serverPort),
		Handler: &router.GlobalRouter,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			lg.Fatal().Err(err).Msg("Server failed")
		}
	}()

	// Wait for signal
	sig := <-quit
	lg.Info().Str("signal", sig.String()).Msg("Received shutdown signal")

	// Shutdown server and services
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		lg.Error().Err(err).Msg("Error shutting down server")
	}

	err = router.GlobalRouter.Shutdown()
	if err != nil {
		lg.Error().Err(err).Msg("Error stopping services")
	}

	lg.Info().Msg("Server shutdown complete")
}
