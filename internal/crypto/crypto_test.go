package crypto

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/logger"
)

var logBuffer = &bytes.Buffer{}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()
	log.Logger = logger.NewTestLogger(logBuffer)
	os.Exit(m.Run())
}
