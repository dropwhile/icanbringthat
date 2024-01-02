package service

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	pg_query "github.com/pganalyze/pg_query_go/v4"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/logger"
	"github.com/dropwhile/icbt/internal/util"
)

var tstTs time.Time = util.MustParseTime(time.RFC3339, "2030-01-01T03:04:05Z")

func SetupDBMock(t *testing.T, ctx context.Context) pgxmock.PgxConnIface {
	t.Helper()

	var queryMatcher pgxmock.QueryMatcher = pgxmock.QueryMatcherFunc(func(expectedSQL, actualSQL string) error {
		err := pgxmock.QueryMatcherRegexp.Match(expectedSQL, actualSQL)
		if err != nil {
			return err
		}
		matchSQL := actualSQL
		// rewrite query like a query rewriter would, if @ is in the query (named args)
		if strings.Contains(actualSQL, "@") {

			newSQL, _, err := (pgx.NamedArgs{}).RewriteQuery(context.Background(), nil, actualSQL, nil)
			if err != nil {
				return fmt.Errorf("error rewriting sql '%s': %w", actualSQL, err)
			}
			matchSQL = newSQL
		}
		_, err = pg_query.Parse(matchSQL)
		if err != nil {
			return fmt.Errorf("error parsing sql '%s': %w", actualSQL, err)
		}

		return nil
	})

	mock, err := pgxmock.NewConn(
		pgxmock.QueryMatcherOption(queryMatcher),
	)
	t.Cleanup(func() { _ = mock.Close(ctx) })
	assert.NilError(t, err)
	return mock
}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()
	logger.SetupLogging(logger.NewTestLogger, nil)
	os.Exit(m.Run())
}