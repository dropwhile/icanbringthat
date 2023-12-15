package rpc

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/logger"
)

var tstTs time.Time = MustParseTime(time.RFC3339, "2030-01-01T03:04:05Z")

func MustParseTime(layout, value string) time.Time {
	ts, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return ts
}

func assertTwirpError(t *testing.T, err error, code twirp.ErrorCode, msg string) {
	t.Helper()
	twerr, ok := err.(twirp.Error)
	if !ok {
		t.Errorf("not a twirp error type")
		return
	}
	if twerr.Code() != code {
		t.Errorf("wrong code. have=%q, want=%q", twerr.Code(), code)
	}
	if twerr.Msg() != msg {
		t.Errorf("wrong msg. have=%q, want=%q", twerr.Msg(), msg)
	}
}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()
	logger.SetupLogging(logger.NewTestLogger, nil)
	os.Exit(m.Run())
}
