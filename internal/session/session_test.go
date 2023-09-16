package session

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/dropwhile/icbt/internal/util"
	"github.com/rs/zerolog/log"
)

var logBuffer = &bytes.Buffer{}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()
	log.Logger = util.NewTestLogger(logBuffer)
	os.Exit(m.Run())
}
