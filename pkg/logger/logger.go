package logger

import (
	"context"
	"os"
	"time"

	"github.com/mjmorales/doan/pkg/agent"
	"github.com/rs/zerolog"
)

// SetGlobalLogConfig sets the global log configuration for zerolog.
func SetGlobalLogConfig(agentConfig agent.AgentConfig) {
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

	logfile := agentConfig.LogFile
	if logfile != "" {
		f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			zerolog.Ctx(context.TODO()).Fatal().Err(err).Msg("failed to open log file")
		}

		zerolog.Ctx(context.TODO()).Info().Str("logfile", logfile).Msg("logging to file")
		zerolog.Ctx(context.TODO()).Output(f)
	}
}
