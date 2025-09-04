// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package service

import (
	"context"
	"testing"

	"github.com/dropwhile/assert"
	"github.com/pashagolub/pgxmock/v4"
)

func TestService_ArchiveOldEvents(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mock := SetupDBMock(t, ctx)
	svc := New(Options{Db: mock})

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE event_").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	mock.ExpectRollback()

	err := svc.ArchiveOldEvents(ctx)
	assert.Nil(t, err)
	// we make sure that all expectations were met
	assert.Nil(t, mock.ExpectationsWereMet(),
		"there were unfulfilled expectations")
}
