package rpc

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/logger"
	"github.com/dropwhile/icbt/internal/util"
)

var tstTs time.Time = util.MustParseTime(time.RFC3339, "2030-01-01T03:04:05Z")

func assertTwirpError(t *testing.T, err error, code twirp.ErrorCode, msg string, meta ...map[string]string) {
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
	for _, m := range meta {
		for k, v := range m {
			x := twerr.Meta(k)
			if x != v {
				t.Errorf("meta value %q mismatch. have=%q, want=%q",
					k, x, v)
			}
		}
	}
}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()
	logger.SetupLogging(logger.NewTestLogger, nil)
	os.Exit(m.Run())
}
