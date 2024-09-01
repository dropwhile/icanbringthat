// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package service

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/util"
)

func TestService_GetEarmarksByEventID(t *testing.T) {
	t.Parallel()

	t.Run("get with results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		eventID := 21

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(eventID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified",
				}).
				AddRow(
					1, util.Must(model.NewEarmarkRefID()), 1, 1,
					"some note 1", tstTs, tstTs,
				).
				AddRow(
					2, util.Must(model.NewEarmarkRefID()), 2, 1,
					"some note 2", tstTs, tstTs,
				),
			)

		results, err := svc.GetEarmarksByEventID(ctx, eventID)
		assert.NilError(t, err)
		assert.Equal(t, len(results), 2)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get with no/empty results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		eventID := 21

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(eventID).
			WillReturnError(pgx.ErrNoRows)

		results, err := svc.GetEarmarksByEventID(ctx, eventID)
		assert.NilError(t, err)
		assert.Equal(t, len(results), 0)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetEarmarkByEventItemID(t *testing.T) {
	t.Parallel()

	t.Run("get with result should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		eventItemID := 21
		earmarkID := 2

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(eventItemID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified",
				}).
				AddRow(
					earmarkID, util.Must(model.NewEarmarkRefID()), eventItemID, 1,
					"some note 2", tstTs, tstTs,
				),
			)

		result, err := svc.GetEarmarkByEventItemID(ctx, eventItemID)
		assert.NilError(t, err)
		assert.Equal(t, result.ID, earmarkID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get with no/empty results should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		eventItemID := 21

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(eventItemID).
			WillReturnError(pgx.ErrNoRows)

		_, err := svc.GetEarmarkByEventItemID(ctx, eventItemID)
		errs.AssertError(t, err, errs.NotFound, "earmark not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetEarmarksCount(t *testing.T) {
	t.Parallel()

	t.Run("count with results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		userID := 4
		currentCount := 4
		archivedCount := 2

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"current", "archived"}).
				AddRow(currentCount, archivedCount),
			)

		result, err := svc.GetEarmarksCount(ctx, userID)
		assert.NilError(t, err)
		assert.Equal(t, result.Current, currentCount)
		assert.Equal(t, result.Archived, archivedCount)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetEarmarksPaginated(t *testing.T) {
	t.Parallel()

	t.Run("get with results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		userID := 4
		limit := 5
		offset := 0
		archived := false
		currentCount := 2
		archivedCount := 3

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"current", "archived"}).
				AddRow(currentCount, archivedCount),
			)
		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(pgx.NamedArgs{
				"userID":   userID,
				"limit":    limit,
				"offset":   offset,
				"archived": archived,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified",
				}).
				AddRow(
					1, util.Must(model.NewEarmarkRefID()), 1, 1,
					"some note 1", tstTs, tstTs,
				).
				AddRow(
					2, util.Must(model.NewEarmarkRefID()), 2, 1,
					"some note 2", tstTs, tstTs,
				),
			)

		earmarks, pagination, err := svc.GetEarmarksPaginated(ctx, userID, limit, offset, archived)
		assert.NilError(t, err)
		assert.Equal(t, len(earmarks), currentCount)
		assert.Equal(t, pagination.Limit, limit)
		assert.Equal(t, pagination.Offset, offset)
		assert.Equal(t, pagination.Count, currentCount)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get with archived results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		userID := 4
		limit := 5
		offset := 0
		archived := true
		currentCount := 2
		archivedCount := 3

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"current", "archived"}).
				AddRow(currentCount, archivedCount),
			)
		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(pgx.NamedArgs{
				"userID":   userID,
				"limit":    limit,
				"offset":   offset,
				"archived": archived,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified",
				}).
				AddRow(
					1, util.Must(model.NewEarmarkRefID()), 1, 1,
					"some note 1", tstTs, tstTs,
				).
				AddRow(
					2, util.Must(model.NewEarmarkRefID()), 2, 1,
					"some note 2", tstTs, tstTs,
				).
				AddRow(
					3, util.Must(model.NewEarmarkRefID()), 3, 1,
					"some note 2", tstTs, tstTs,
				),
			)

		earmarks, pagination, err := svc.GetEarmarksPaginated(ctx, userID, limit, offset, archived)
		assert.NilError(t, err)
		assert.Equal(t, len(earmarks), archivedCount)
		assert.Equal(t, pagination.Limit, limit)
		assert.Equal(t, pagination.Offset, offset)
		assert.Equal(t, pagination.Count, archivedCount)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get with no results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		userID := 4
		limit := 5
		offset := 0
		archived := false
		currentCount := 0
		archivedCount := 2

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"current", "archived"}).
				AddRow(currentCount, archivedCount),
			)

		earmarks, pagination, err := svc.GetEarmarksPaginated(ctx, userID, limit, offset, archived)
		assert.NilError(t, err)
		assert.Equal(t, len(earmarks), currentCount)
		assert.Equal(t, pagination.Limit, limit)
		assert.Equal(t, pagination.Offset, offset)
		assert.Equal(t, pagination.Count, currentCount)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetEarmarks(t *testing.T) {
	t.Parallel()

	t.Run("get with results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		userID := 4
		archived := false

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(pgx.NamedArgs{
				"userID":   userID,
				"archived": archived,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified",
				}).
				AddRow(
					1, util.Must(model.NewEarmarkRefID()), 1, 1,
					"some note 1", tstTs, tstTs,
				).
				AddRow(
					2, util.Must(model.NewEarmarkRefID()), 2, 1,
					"some note 2", tstTs, tstTs,
				),
			)

		earmarks, err := svc.GetEarmarks(ctx, userID, archived)
		assert.NilError(t, err)
		assert.Equal(t, len(earmarks), 2)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get with no results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		userID := 4
		archived := false

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(pgx.NamedArgs{
				"userID":   userID,
				"archived": archived,
			}).
			WillReturnError(pgx.ErrNoRows)

		earmarks, err := svc.GetEarmarks(ctx, userID, archived)
		assert.NilError(t, err)
		assert.Equal(t, len(earmarks), 0)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_NewEarmark(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:           1,
		RefID:        util.Must(model.NewEventRefID()),
		UserID:       user.ID,
		Name:         "event",
		Description:  "description",
		Archived:     false,
		StartTime:    ts,
		StartTimeTz:  util.Must(ParseTimeZone("Etc/UTC")),
		Created:      ts,
		LastModified: ts,
	}
	eventItem := &model.EventItem{
		ID:           2,
		RefID:        util.Must(model.NewEventItemRefID()),
		EventID:      event.ID,
		Description:  "eventitem",
		Created:      ts,
		LastModified: ts,
	}
	earmark := &model.Earmark{
		ID:           3,
		RefID:        util.Must(model.NewEarmarkRefID()),
		EventItemID:  eventItem.ID,
		UserID:       user.ID,
		Note:         "nothing",
		Created:      ts,
		LastModified: ts,
	}

	t.Run("create earmark", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(eventItem.ID).
			WillReturnError(pgx.ErrNoRows)
		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(eventItem.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"start_time", "start_time_tz", "created", "last_modified",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					event.StartTime, event.StartTimeTz, ts, ts,
				),
			)
		mock.ExpectBegin()
		mock.ExpectQuery("^INSERT INTO earmark_").
			WithArgs(pgx.NamedArgs{
				"refID":       EarmarkRefIDMatcher,
				"eventItemID": earmark.EventItemID,
				"userID":      earmark.UserID,
				"note":        "some note",
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified",
				}).
				AddRow(
					earmark.ID, earmark.RefID, earmark.EventItemID, earmark.UserID,
					earmark.Note, ts, ts,
				),
			)
		mock.ExpectCommit()
		mock.ExpectRollback()

		_, err := svc.NewEarmark(ctx, user, earmark.EventItemID, "some note")
		assert.NilError(t, err)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create earmark missing event", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(eventItem.ID).
			WillReturnError(pgx.ErrNoRows)
		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(eventItem.ID).
			WillReturnError(pgx.ErrNoRows)

		_, err := svc.NewEarmark(ctx, user, eventItem.ID, "some note")
		errs.AssertError(t, err, errs.NotFound, "event not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create earmark user not verified", func(t *testing.T) {
		t.Parallel()

		user := &model.User{
			ID:           33,
			RefID:        util.Must(model.NewUserRefID()),
			Email:        "user@example.com",
			Name:         "user",
			PWHash:       []byte("00x00"),
			Verified:     false,
			Created:      ts,
			LastModified: ts,
		}

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("SELECT (.+) FROM earmark_").
			WithArgs(eventItem.ID).
			WillReturnError(pgx.ErrNoRows)
		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(eventItem.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived", "created", "last_modified",
				}).
				AddRow(
					event.ID, util.Must(model.NewEventRefID()), 34,
					"event name", "event desc",
					false, tstTs, tstTs,
				),
			)

		_, err := svc.NewEarmark(ctx, user, eventItem.ID, "some note")
		errs.AssertError(t, err, errs.PermissionDenied, "Account must be verified before earmarking is allowed.")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetEarmark(t *testing.T) {
	t.Parallel()

	t.Run("get with result should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		earmarkRefID := util.Must(model.NewEarmarkRefID())

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(earmarkRefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified",
				}).
				AddRow(
					1, util.Must(model.NewEarmarkRefID()), 1, 1,
					"some note 1", tstTs, tstTs,
				),
			)

		earmark, err := svc.GetEarmark(ctx, earmarkRefID)
		assert.NilError(t, err)
		assert.Equal(t, earmark.ID, 1)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get without result should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		earmarkRefID := util.Must(model.NewEarmarkRefID())

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(earmarkRefID).
			WillReturnError(pgx.ErrNoRows)

		_, err := svc.GetEarmark(ctx, earmarkRefID)
		errs.AssertError(t, err, errs.NotFound, "earmark not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_DeleteEarmark(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:           1,
		RefID:        util.Must(model.NewEventRefID()),
		UserID:       user.ID,
		Name:         "event",
		Description:  "description",
		Archived:     false,
		StartTime:    ts,
		StartTimeTz:  util.Must(ParseTimeZone("Etc/UTC")),
		Created:      ts,
		LastModified: ts,
	}
	eventItem := &model.EventItem{
		ID:           2,
		RefID:        util.Must(model.NewEventItemRefID()),
		EventID:      event.ID,
		Description:  "eventitem",
		Created:      ts,
		LastModified: ts,
	}
	earmark := &model.Earmark{
		ID:           3,
		RefID:        util.Must(model.NewEarmarkRefID()),
		EventItemID:  eventItem.ID,
		UserID:       user.ID,
		Note:         "nothing",
		Created:      ts,
		LastModified: ts,
	}

	t.Run("delete should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(eventItem.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description", "archived",
					"start_time", "start_time_tz", "created", "last_modified",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				),
			)
		mock.ExpectBegin()
		mock.ExpectExec("^DELETE FROM earmark_").
			WithArgs(earmark.ID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		err := svc.DeleteEarmark(ctx, user.ID, earmark)
		assert.NilError(t, err)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete with different user owner should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		err := svc.DeleteEarmark(ctx, user.ID+1, earmark)
		errs.AssertError(t, err, errs.PermissionDenied, "permission denied")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete with missing event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(eventItem.ID).
			WillReturnError(pgx.ErrNoRows)

		err := svc.DeleteEarmark(ctx, user.ID, earmark)
		errs.AssertError(t, err, errs.NotFound, "event not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete with archived event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(eventItem.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description", "archived",
					"start_time", "start_time_tz", "created", "last_modified",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					true, event.StartTime, event.StartTimeTz, ts, ts,
				),
			)

		err := svc.DeleteEarmark(ctx, user.ID, earmark)
		errs.AssertError(t, err, errs.PermissionDenied, "event is archived")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_DeleteEarmarkByRefID(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:           1,
		RefID:        util.Must(model.NewEventRefID()),
		UserID:       user.ID,
		Name:         "event",
		Description:  "description",
		Archived:     false,
		StartTime:    ts,
		StartTimeTz:  util.Must(ParseTimeZone("Etc/UTC")),
		Created:      ts,
		LastModified: ts,
	}
	eventItem := &model.EventItem{
		ID:           2,
		RefID:        util.Must(model.NewEventItemRefID()),
		EventID:      event.ID,
		Description:  "eventitem",
		Created:      ts,
		LastModified: ts,
	}
	earmark := &model.Earmark{
		ID:           3,
		RefID:        util.Must(model.NewEarmarkRefID()),
		EventItemID:  eventItem.ID,
		UserID:       user.ID,
		Note:         "nothing",
		Created:      ts,
		LastModified: ts,
	}

	t.Run("delete should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(earmark.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified",
				}).
				AddRow(
					earmark.ID, earmark.RefID, earmark.EventItemID,
					earmark.UserID, earmark.Note,
					earmark.Created, earmark.LastModified,
				),
			)
		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(eventItem.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description", "archived",
					"start_time", "start_time_tz", "created", "last_modified",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				),
			)
		mock.ExpectBegin()
		mock.ExpectExec("^DELETE FROM earmark_").
			WithArgs(earmark.ID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		err := svc.DeleteEarmarkByRefID(ctx, user.ID, earmark.RefID)
		assert.NilError(t, err)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete with different user owner should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(earmark.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified",
				}).
				AddRow(
					earmark.ID, earmark.RefID, earmark.EventItemID,
					earmark.UserID, earmark.Note,
					earmark.Created, earmark.LastModified,
				),
			)
		err := svc.DeleteEarmarkByRefID(ctx, user.ID+1, earmark.RefID)
		errs.AssertError(t, err, errs.PermissionDenied, "permission denied")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete with missing event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(earmark.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified",
				}).
				AddRow(
					earmark.ID, earmark.RefID, earmark.EventItemID,
					earmark.UserID, earmark.Note,
					earmark.Created, earmark.LastModified,
				),
			)
		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(eventItem.ID).
			WillReturnError(pgx.ErrNoRows)

		err := svc.DeleteEarmarkByRefID(ctx, user.ID, earmark.RefID)
		errs.AssertError(t, err, errs.NotFound, "event not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete with archived event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(earmark.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified",
				}).
				AddRow(
					earmark.ID, earmark.RefID, earmark.EventItemID,
					earmark.UserID, earmark.Note,
					earmark.Created, earmark.LastModified,
				),
			)
		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(eventItem.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description", "archived",
					"start_time", "start_time_tz", "created", "last_modified",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					true, event.StartTime, event.StartTimeTz, ts, ts,
				),
			)

		err := svc.DeleteEarmarkByRefID(ctx, user.ID, earmark.RefID)
		errs.AssertError(t, err, errs.PermissionDenied, "event is archived")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete with missing earmark should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(earmark.RefID).
			WillReturnError(pgx.ErrNoRows)

		err := svc.DeleteEarmarkByRefID(ctx, user.ID, earmark.RefID)
		errs.AssertError(t, err, errs.NotFound, "earmark not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
