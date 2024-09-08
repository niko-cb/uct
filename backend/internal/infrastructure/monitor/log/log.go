package log

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var log *zerolog.Logger

// Initialize initializes zerolog with the necessary settings
func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano           // Set timestamp format
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack // Enable stack trace marshaling

	l := zerolog.New(os.Stdout).With().Timestamp().Logger()
	log = &l
}

func Debug(ctx context.Context, msg string) {
	log.Debug().
		Str("severity", "DEBUG").
		Msg(msg)
}

func Info(ctx context.Context, msg string) {
	log.Info().
		Str("severity", "INFO").
		Msg(msg)
}

func Warning(ctx context.Context, err error) {
	log.Warn().
		Stack().
		Err(err).
		Str("severity", "WARNING").
		Msg(err.Error())
}

func Error(ctx context.Context, err error) {
	log.Error().
		Stack().
		Err(err).
		Str("severity", "ERROR").
		Msg(err.Error())
}

func Fatal(ctx context.Context, err error) {
	log.Fatal().
		Stack().
		Err(err).
		Str("severity", "ALERT").
		Msg(err.Error())
}
