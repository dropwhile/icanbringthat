package util

import (
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewTestLogger(w io.Writer) zerolog.Logger {
	// don't log timestamps for test logs,
	// to enable log matching if desired.
	logger := log.Output(
		zerolog.ConsoleWriter{Out: w},
	)

	switch strings.ToLower(os.Getenv("log_level")) {
	case "debug":
		logger = logger.Level(zerolog.DebugLevel)
		log.Debug().Msg("setting log level to debug")
	case "trace":
		logger = logger.Level(zerolog.TraceLevel)
		log.Trace().Msg("setting log level to trace")
	default:
		logger = logger.Level(zerolog.InfoLevel)
		log.Info().Msg("setting log level to info")
	}
	return logger

}
