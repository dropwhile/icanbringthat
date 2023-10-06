package model

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/util"
)

var tstTs time.Time

func init() {
	tstTs, _ = time.Parse(time.RFC3339, "2023-01-01T03:04:05Z")
}

// var logBuffer = &bytes.Buffer{}
func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()
	// log.Logger = util.NewTestLogger(logBuffer)
	log.Logger = util.NewTestLogger(os.Stderr)
	os.Exit(m.Run())
}
