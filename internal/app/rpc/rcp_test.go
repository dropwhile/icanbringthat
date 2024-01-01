package rpc

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/dropwhile/icbt/internal/logger"
	"github.com/dropwhile/icbt/internal/util"
)

var tstTs time.Time = util.MustParseTime(time.RFC3339, "2030-01-01T03:04:05Z")

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()
	logger.SetupLogging(logger.NewTestLogger, nil)
	os.Exit(m.Run())
}
