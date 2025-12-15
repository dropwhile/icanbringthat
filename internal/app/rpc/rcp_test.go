// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rpc

import (
	"errors"
	"flag"
	"testing"
	"time"

	"connectrpc.com/connect"

	"github.com/dropwhile/icanbringthat/internal/logger"
	"github.com/dropwhile/icanbringthat/internal/util"
)

var tstTs time.Time = util.MustParseTime(time.RFC3339, "2030-01-01T03:04:05Z")

// CodeOf returns the error's status code if it is or wraps an [*Error] and
// [CodeUnknown] otherwise.
func AsConnectError(t *testing.T, err error) *connect.Error {
	t.Helper()
	if err == nil {
		t.Fatal("got a nil error when we expected a connectRpc error")
	}
	if connectErr := new(connect.Error); errors.As(err, &connectErr) {
		return connectErr
	}
	t.Fatalf("got a non-connectRpc error: %s", err)
	return nil
}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()
	logger.SetupLogging(logger.NewTestLogger, nil)
	m.Run()
}
