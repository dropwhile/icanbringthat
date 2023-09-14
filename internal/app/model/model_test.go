package model

import (
	"bytes"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/cactus/mlog"
)

var tstTs time.Time

func init() {
	tstTs, _ = time.Parse(time.RFC3339, "2023-01-01T03:04:05Z")
}

var logBuffer = &bytes.Buffer{}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()

	debug := os.Getenv("DEBUG")
	verbose := os.Getenv("TESTS_LOG_VERBOSE")

	// now configure a standard logger
	mlog.SetFlags(mlog.Lstd)

	if debug != "" {
		// don't log timestamps for test logs,
		// to enable log matching if desired.
		mlog.SetFlags(mlog.Llevel | mlog.Lsort | mlog.Ldebug)
	}

	// log to bytes buffer
	if verbose != "" {
		mlog.DefaultLogger = mlog.New(logBuffer, mlog.Lstd)
	}
	mlog.Debug("debug logging enabled")

	os.Exit(m.Run())
}
