package logger

import (
	"context"
	"os"
	"time"

	"github.com/mjmorales/doan/pkg/agent"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

func SetLogFile(agentConfig agent.AgentConfig) (*os.File, error) {
	logfile := agentConfig.LogFile
	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		zerolog.Ctx(context.TODO()).Fatal().Err(err).Msg("failed to open log file")
	}

	log.Logger = log.With().Caller().Logger().Output(f)
	log.Debug().Msgf("logging to file %s", logfile)
	return f, nil
}
