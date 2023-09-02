package model

import (
	"bytes"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cactus/mlog"
)

var tstTs time.Time

func init() {
	tstTs, _ = time.Parse(time.RFC3339, "2023-01-01T03:04:05Z")
}

func setupDBMock(t *testing.T) (*DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	model, err := SetupFromDb(db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return model, mock
}

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
