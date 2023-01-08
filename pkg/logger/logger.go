package logger

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// SetGlobalLogConfig sets the global log configuration for zerolog.
func SetGlobalLogConfig() {
	// set the time format to RFC3339
	zerolog.TimeFieldFormat = time.RFC3339

	// sets the log level to debug if the DEBUG envvar is passed
	if os.Getenv("DEBUG") != "" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// notify the user if the DEBUG envvar is set
	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		zerolog.Ctx(context.TODO()).Info().Msg("debug logging enabled")
	}
}
