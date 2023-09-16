package model

import (
	"bytes"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/dropwhile/icbt/internal/util"
	"github.com/rs/zerolog/log"
)

var tstTs time.Time

func init() {
	tstTs, _ = time.Parse(time.RFC3339, "2023-01-01T03:04:05Z")
}

var logBuffer = &bytes.Buffer{}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()
	log.Logger = util.NewTestLogger(logBuffer)
	os.Exit(m.Run())
}
