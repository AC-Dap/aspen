package logging

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var outputs []io.Writer

func DisableLogger() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

func InitializeLogger(level zerolog.Level) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(level)
}

func AddConsoleOutput(prettyPrint bool) {
	if prettyPrint {
		outputs = append(outputs, zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		outputs = append(outputs, os.Stderr)
	}
	updateOutputs()
}

func AddFileOutput(filename string) error {
	f, err := os.OpenFile(
		filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		return fmt.Errorf("unable to open log file: %w", err)
	}
	outputs = append(outputs, f)
	updateOutputs()
	return nil
}

func updateOutputs() {
	var writer io.Writer
	if len(outputs) > 1 {
		writer = zerolog.MultiLevelWriter(outputs...)
	} else {
		writer = outputs[0]
	}
	log.Logger = zerolog.New(writer).With().Timestamp().Logger()
}
