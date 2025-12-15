// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package session

import (
	"bytes"
	"flag"
	"testing"

	"github.com/dropwhile/icanbringthat/internal/logger"
)

var logBuffer = &bytes.Buffer{}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()
	logger.SetupLogging(logger.NewTestLogger,
		&logger.Options{
			Sink: logBuffer,
		},
	)
	m.Run()
}
