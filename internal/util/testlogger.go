package util

import (
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewTestLogger(w io.Writer) zerolog.Logger {
	// don't log timestamps for test logs,
	// to enable log matching if desired.
	logger := log.Output(
		zerolog.ConsoleWriter{Out: w},
	).With().Caller().Logger()
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}

	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		logger = logger.Level(zerolog.DebugLevel)
		logger.Debug().Msg("setting log level to debug")
	case "trace":
		logger = logger.Level(zerolog.TraceLevel)
		logger.Trace().Msg("setting log level to trace")
	default:
		logger = logger.Level(zerolog.InfoLevel)
		logger.Info().Msgf("unexpected LOG_LEVEL env var set to '%s'",
			strings.ToLower(os.Getenv("LOG_LEVEL")))
		logger.Info().Msg("setting log level to info")
	}
	return logger
}
