package logging

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func DisableLogger() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

func InitializeLogger(level zerolog.Level) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(level)
}

func SetOutputToConsole() {
	writer := zerolog.ConsoleWriter{Out: os.Stderr}
	log.Logger = zerolog.New(writer).With().Timestamp().Logger()
}

func SetOutputToFile(filepath string) error {
	f, err := os.OpenFile(
		filepath,
		os.O_TRUNC|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		return fmt.Errorf("unable to open log file: %w", err)
	}

	writer := zerolog.ConsoleWriter{Out: f, NoColor: true}
	log.Logger = zerolog.New(writer).With().Timestamp().Logger()
	return nil
}
