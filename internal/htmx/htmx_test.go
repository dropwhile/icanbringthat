package htmx

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/dropwhile/icbt/internal/logger"
)

var logBuffer = &bytes.Buffer{}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()
	logger.SetupLogging(logger.NewTestLogger,
		&logger.Options{
			Sink: logBuffer,
		},
	)
	os.Exit(m.Run())
}
