package logging

import (
	"fmt"
	"io"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const TAG_FIELD_NAME = "aspen.tagger_logger.tag"

/*
 * TaggedLogger wraps around zerolog's logger and marks every message with a tag.
 * This tag is interpreted by TaggedLoggerWriter and prepends every log like so:
 *		2025-12-24 INF [TAG] ...
 */
type TaggedLogger struct {
	tag string
}

func NewTaggedLogger(tag string) TaggedLogger {
	return TaggedLogger{
		tag,
	}
}

func (t *TaggedLogger) Panic() *zerolog.Event {
	return log.Panic().Str(TAG_FIELD_NAME, t.tag)
}

func (t *TaggedLogger) Fatal() *zerolog.Event {
	return log.Fatal().Str(TAG_FIELD_NAME, t.tag)
}

func (t *TaggedLogger) Error() *zerolog.Event {
	return log.Error().Str(TAG_FIELD_NAME, t.tag)
}

func (t *TaggedLogger) Warn() *zerolog.Event {
	return log.Warn().Str(TAG_FIELD_NAME, t.tag)
}

func (t *TaggedLogger) Info() *zerolog.Event {
	return log.Info().Str(TAG_FIELD_NAME, t.tag)
}

func (t *TaggedLogger) Debug() *zerolog.Event {
	return log.Debug().Str(TAG_FIELD_NAME, t.tag)
}

func (t *TaggedLogger) Trace() *zerolog.Event {
	return log.Trace().Str(TAG_FIELD_NAME, t.tag)
}

func NewTaggedLoggerWriter(out io.Writer, noColor bool) io.Writer {
	// Output log lines like
	// 		2025-12-24 13:04:05 INF [TAG] ...
	writer := zerolog.ConsoleWriter{
		Out: out, NoColor: noColor,
		TimeFormat:    time.DateTime,
		PartsOrder:    []string{"time", "level", TAG_FIELD_NAME, "message"},
		FieldsExclude: []string{TAG_FIELD_NAME},
	}
	writer.FormatPartValueByName = func(i any, _ string) string {
		return fmt.Sprintf("[%s]", i)
	}

	return writer
}
