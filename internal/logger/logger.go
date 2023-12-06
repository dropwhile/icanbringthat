package logger

import (
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func NewTestLogger(w io.Writer) zerolog.Logger {
	// don't log timestamps for test logs,
	// to enable log matching if desired.
	logger := log.Output(
		zerolog.ConsoleWriter{
			Out:          w,
			PartsExclude: []string{zerolog.TimestampFieldName},
		},
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

func NewLogger(w io.Writer) zerolog.Logger {
	// don't log timestamps for test logs,
	// to enable log matching if desired.
	logger := log.Output(
		zerolog.ConsoleWriter{
			Out:        w,
			TimeFormat: time.RFC3339,
		},
	).With().Caller().Logger()
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		counter := 0
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				if counter > 0 {
					short = file[i+1:]
					break
				}
				counter += 1
			}
		}

		// prune from after @ to next /
		atIdx := strings.Index(short, "@")
		if atIdx >= 0 && atIdx+7 <= len(short) {
			for i := atIdx; i < len(short); i++ {
				if short[i] == '/' {
					if i-atIdx > 7 {
						short = short[:atIdx+7] + "..." + short[i:]
					}
					break
				}
			}
		}

		file = short
		return file + ":" + strconv.Itoa(line)
	}
	return logger
}
