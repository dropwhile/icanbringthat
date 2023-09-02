package session

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/cactus/mlog"
)

var logBuffer = &bytes.Buffer{}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()

	debug := os.Getenv("DEBUG")
	// now configure a standard logger
	mlog.SetFlags(mlog.Lstd)

	if debug != "" {
		// don't log timestamps for test logs,
		// to enable log matching if desired.
		mlog.SetFlags(mlog.Llevel | mlog.Lsort | mlog.Ldebug)
	}

	// log to bytes buffer
	mlog.DefaultLogger = mlog.New(logBuffer, mlog.Lstd)
	mlog.Debug("debug logging enabled")

	os.Exit(m.Run())
}
