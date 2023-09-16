package util

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/rs/zerolog/log"
)

var logBuffer = &bytes.Buffer{}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()
	log.Logger = NewTestLogger(logBuffer)
	os.Exit(m.Run())
}
