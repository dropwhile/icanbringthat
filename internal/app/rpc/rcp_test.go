package rpc

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v3"
	pg_query "github.com/pganalyze/pg_query_go/v4"
	"github.com/twitchtv/twirp"
	"gotest.tools/v3/assert"

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

func SetupDBMock(t *testing.T, ctx context.Context) pgxmock.PgxConnIface {
	t.Helper()

	var queryMatcher pgxmock.QueryMatcher = pgxmock.QueryMatcherFunc(func(expectedSQL, actualSQL string) error {
		err := pgxmock.QueryMatcherRegexp.Match(expectedSQL, actualSQL)
		if err != nil {
			return err
		}
		_, err = pg_query.Parse(actualSQL)
		if err != nil {
			return fmt.Errorf("error parsing sql '%s': %w", actualSQL, err)
		}

		return nil
	})

	mock, err := pgxmock.NewConn(
		pgxmock.QueryMatcherOption(queryMatcher),
	)
	assert.NilError(t, err)
	return mock
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
