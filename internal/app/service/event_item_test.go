package service

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/util"
)

func TestService_GetEventItemsCount(t *testing.T) {
	t.Parallel()

	t.Run("count with results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		eventIDs := []int{3, 4}

		mock.ExpectQuery("^SELECT (.+) FROM event_").
			WithArgs(eventIDs).
			WillReturnRows(pgxmock.NewRows(
				[]string{"event_id", "count"}).
				AddRow(eventIDs[0], 2).
				AddRow(eventIDs[1], 3),
			)

		result, err := svc.GetEventItemsCount(ctx, eventIDs)
		assert.NilError(t, err)
		assert.Equal(t, len(result), 2)
		assert.Equal(t, result[0].EventID, eventIDs[0])
		assert.Equal(t, result[0].Count, 2)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("count with empty results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		eventIDs := []int{3, 4}

		mock.ExpectQuery("^SELECT (.+) FROM event_").
			WithArgs(eventIDs).
			WillReturnRows(pgxmock.NewRows(
				[]string{"event_id", "count"}),
			)

		result, err := svc.GetEventItemsCount(ctx, eventIDs)
		assert.NilError(t, err)
		assert.Equal(t, len(result), 0)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("count with empty input should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		result, err := svc.GetEventItemsCount(ctx, []int{})
		assert.NilError(t, err)
		assert.Equal(t, len(result), 0)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetEventItemsByEvent(t *testing.T) {
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
		ID:           3,
		RefID:        util.Must(model.NewEventItemRefID()),
		EventID:      event.ID,
		Description:  "event-item",
		Created:      ts,
		LastModified: ts,
	}

	t.Run("get should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name,
					event.Description, event.Archived,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_id", "description",
				}).
				AddRow(
					eventItem.ID, eventItem.RefID,
					eventItem.EventID, eventItem.Description,
				),
			)

		result, err := svc.GetEventItemsByEvent(ctx, event.RefID)
		assert.NilError(t, err)
		assert.Equal(t, len(result), 1)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get no results should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name,
					event.Description, event.Archived,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_id", "description",
				}),
			)

		result, err := svc.GetEventItemsByEvent(ctx, event.RefID)
		assert.NilError(t, err)
		assert.Equal(t, len(result), 0)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get event not found should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.RefID).
			WillReturnError(pgx.ErrNoRows)

		_, err := svc.GetEventItemsByEvent(ctx, event.RefID)
		errs.AssertError(t, err, errs.NotFound, "event not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetEventItemsByEventID(t *testing.T) {
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
		ID:           3,
		RefID:        util.Must(model.NewEventItemRefID()),
		EventID:      event.ID,
		Description:  "event-item",
		Created:      ts,
		LastModified: ts,
	}

	t.Run("get should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_id", "description",
				}).
				AddRow(
					eventItem.ID, eventItem.RefID,
					eventItem.EventID, eventItem.Description,
				),
			)

		result, err := svc.GetEventItemsByEventID(ctx, event.ID)
		assert.NilError(t, err)
		assert.Equal(t, len(result), 1)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get no results should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_id", "description",
				}),
			)

		result, err := svc.GetEventItemsByEventID(ctx, event.ID)
		assert.NilError(t, err)
		assert.Equal(t, len(result), 0)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetEventItemsByIDs(t *testing.T) {
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
		ID:           3,
		RefID:        util.Must(model.NewEventItemRefID()),
		EventID:      event.ID,
		Description:  "event-item",
		Created:      ts,
		LastModified: ts,
	}
	eventItem2 := &model.EventItem{
		ID:           4,
		RefID:        util.Must(model.NewEventItemRefID()),
		EventID:      event.ID,
		Description:  "event-item2",
		Created:      ts,
		LastModified: ts,
	}

	t.Run("get should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		eventItemIDs := []int{eventItem.ID, eventItem2.ID}

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(eventItemIDs).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_id", "description",
				}).
				AddRow(
					eventItem.ID, eventItem.RefID,
					eventItem.EventID, eventItem.Description,
				).
				AddRow(
					eventItem2.ID, eventItem2.RefID,
					eventItem2.EventID, eventItem2.Description,
				),
			)

		result, err := svc.GetEventItemsByIDs(ctx, eventItemIDs)
		assert.NilError(t, err)
		assert.Equal(t, len(result), 2)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get no results should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		eventItemIDs := []int{eventItem.ID, eventItem2.ID}

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(eventItemIDs).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_id", "description",
				}),
			)

		result, err := svc.GetEventItemsByIDs(ctx, eventItemIDs)
		assert.NilError(t, err)
		assert.Equal(t, len(result), 0)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get empty input should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		result, err := svc.GetEventItemsByIDs(ctx, []int{})
		assert.NilError(t, err)
		assert.Equal(t, len(result), 0)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetEventItem(t *testing.T) {
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
		ID:           3,
		RefID:        util.Must(model.NewEventItemRefID()),
		EventID:      event.ID,
		Description:  "event-item",
		Created:      ts,
		LastModified: ts,
	}

	t.Run("get should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(eventItem.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_id", "description",
				}).
				AddRow(
					eventItem.ID, eventItem.RefID,
					eventItem.EventID, eventItem.Description,
				),
			)

		result, err := svc.GetEventItem(ctx, eventItem.RefID)
		assert.NilError(t, err)
		assert.Equal(t, result.RefID, eventItem.RefID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get not found should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(eventItem.RefID).
			WillReturnError(pgx.ErrNoRows)

		_, err := svc.GetEventItem(ctx, eventItem.RefID)
		errs.AssertError(t, err, errs.NotFound, "event-item not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetEventItemByID(t *testing.T) {
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
		ID:           3,
		RefID:        util.Must(model.NewEventItemRefID()),
		EventID:      event.ID,
		Description:  "event-item",
		Created:      ts,
		LastModified: ts,
	}

	t.Run("get should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(eventItem.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_id", "description",
				}).
				AddRow(
					eventItem.ID, eventItem.RefID,
					eventItem.EventID, eventItem.Description,
				),
			)

		result, err := svc.GetEventItemByID(ctx, eventItem.ID)
		assert.NilError(t, err)
		assert.Equal(t, result.RefID, eventItem.RefID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get not found should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(eventItem.ID).
			WillReturnError(pgx.ErrNoRows)

		_, err := svc.GetEventItemByID(ctx, eventItem.ID)
		errs.AssertError(t, err, errs.NotFound, "event-item not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_RemoveEventItem(t *testing.T) {
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
		ID:           3,
		RefID:        util.Must(model.NewEventItemRefID()),
		EventID:      event.ID,
		Description:  "event-item",
		Created:      ts,
		LastModified: ts,
	}

	t.Run("remove should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(eventItem.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_id", "description",
				}).
				AddRow(
					eventItem.ID, eventItem.RefID,
					eventItem.EventID, eventItem.Description,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name,
					event.Description, event.Archived,
				),
			)
		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM event_item_ ").
			WithArgs(eventItem.ID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		err := svc.RemoveEventItem(ctx, user.ID, eventItem.RefID, nil)
		assert.NilError(t, err)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete with archived event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(eventItem.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_id", "description",
				}).
				AddRow(
					eventItem.ID, eventItem.RefID,
					eventItem.EventID, eventItem.Description,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name,
					event.Description, true,
				),
			)

		err := svc.RemoveEventItem(ctx, user.ID, eventItem.RefID, nil)
		errs.AssertError(t, err, errs.PermissionDenied, "event is archived")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete with other owner should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(eventItem.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_id", "description",
				}).
				AddRow(
					eventItem.ID, eventItem.RefID,
					eventItem.EventID, eventItem.Description,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID+1, event.Name,
					event.Description, event.Archived,
				),
			)

		err := svc.RemoveEventItem(ctx, user.ID, eventItem.RefID, nil)
		errs.AssertError(t, err, errs.PermissionDenied, "not event owner")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete with missing event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(eventItem.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_id", "description",
				}).
				AddRow(
					eventItem.ID, eventItem.RefID,
					eventItem.EventID, eventItem.Description,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.ID).
			WillReturnError(pgx.ErrNoRows)

		err := svc.RemoveEventItem(ctx, user.ID, eventItem.RefID, nil)
		errs.AssertError(t, err, errs.NotFound, "event not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete with failing failIfCheck should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(eventItem.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_id", "description",
				}).
				AddRow(
					eventItem.ID, eventItem.RefID,
					eventItem.EventID, eventItem.Description,
				),
			)

		err := svc.RemoveEventItem(ctx, user.ID, eventItem.RefID,
			func(ei *model.EventItem) bool { return true },
		)
		errs.AssertError(t, err, errs.FailedPrecondition, "extra checks failed")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete with missing event-item should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(eventItem.RefID).
			WillReturnError(pgx.ErrNoRows)

		err := svc.RemoveEventItem(ctx, user.ID, eventItem.RefID, nil)
		errs.AssertError(t, err, errs.NotFound, "event-item not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_AddEventItem(t *testing.T) {
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
		ID:           3,
		RefID:        util.Must(model.NewEventItemRefID()),
		EventID:      event.ID,
		Description:  "event-item",
		Created:      ts,
		LastModified: ts,
	}

	t.Run("add item should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name,
					event.Description, event.Archived,
				),
			)
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO event_item_").
			WithArgs(pgx.NamedArgs{
				"refID":       EventItemRefIDMatcher,
				"eventID":     eventItem.EventID,
				"description": eventItem.Description,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_id", "description",
				}).
				AddRow(
					eventItem.ID, eventItem.RefID,
					eventItem.EventID, eventItem.Description,
				),
			)
		mock.ExpectCommit()
		mock.ExpectRollback()

		result, err := svc.AddEventItem(
			ctx, user.ID, event.RefID, eventItem.Description,
		)
		assert.NilError(t, err)
		assert.Equal(t, result.RefID, eventItem.RefID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("add item with archived event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name,
					event.Description, true,
				),
			)

		_, err := svc.AddEventItem(
			ctx, user.ID, event.RefID, eventItem.Description,
		)
		errs.AssertError(t, err, errs.PermissionDenied, "event is archived")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("add item not owner should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID+1, event.Name,
					event.Description, event.Archived,
				),
			)

		_, err := svc.AddEventItem(
			ctx, user.ID, event.RefID, eventItem.Description,
		)
		errs.AssertError(t, err, errs.PermissionDenied, "not event owner")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("add item missing event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.RefID).
			WillReturnError(pgx.ErrNoRows)

		_, err := svc.AddEventItem(
			ctx, user.ID, event.RefID, eventItem.Description,
		)
		errs.AssertError(t, err, errs.NotFound, "event not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("add item with bad value should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		_, err := svc.AddEventItem(
			ctx, user.ID, event.RefID, "",
		)
		errs.AssertError(t, err, errs.InvalidArgument, "description bad value")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_UpdateEventItem(t *testing.T) {
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
		ID:           3,
		RefID:        util.Must(model.NewEventItemRefID()),
		EventID:      event.ID,
		Description:  "event-item",
		Created:      ts,
		LastModified: ts,
	}
	earmark := &model.Earmark{
		ID:           4,
		RefID:        util.Must(model.NewEarmarkRefID()),
		EventItemID:  eventItem.ID,
		UserID:       user.ID,
		Note:         "earmark",
		Created:      ts,
		LastModified: ts,
	}

	t.Run("update not earmarked should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		description := "hodor"

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(eventItem.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_id", "description",
				}).
				AddRow(
					eventItem.ID, eventItem.RefID,
					eventItem.EventID, eventItem.Description,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name,
					event.Description, event.Archived,
				),
			)
		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(eventItem.ID).
			WillReturnError(pgx.ErrNoRows)
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE event_item_ ").
			WithArgs(pgx.NamedArgs{
				"description": description,
				"eventItemID": eventItem.ID,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		result, err := svc.UpdateEventItem(ctx, user.ID, eventItem.RefID, description, nil)
		assert.NilError(t, err)
		assert.Equal(t, result.RefID, eventItem.RefID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update earmarked by self should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		description := "hodor"

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(eventItem.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_id", "description",
				}).
				AddRow(
					eventItem.ID, eventItem.RefID,
					eventItem.EventID, eventItem.Description,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name,
					event.Description, event.Archived,
				),
			)
		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(eventItem.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_item_id", "user_id", "note",
				}).
				AddRow(
					earmark.ID, earmark.RefID, earmark.EventItemID,
					earmark.UserID, earmark.Note,
				),
			)
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE event_item_ ").
			WithArgs(pgx.NamedArgs{
				"description": description,
				"eventItemID": eventItem.ID,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		result, err := svc.UpdateEventItem(ctx, user.ID, eventItem.RefID, description, nil)
		assert.NilError(t, err)
		assert.Equal(t, result.RefID, eventItem.RefID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update earmarked by other user should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		description := "hodor"

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(eventItem.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_id", "description",
				}).
				AddRow(
					eventItem.ID, eventItem.RefID,
					eventItem.EventID, eventItem.Description,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name,
					event.Description, event.Archived,
				),
			)
		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(eventItem.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_item_id", "user_id", "note",
				}).
				AddRow(
					earmark.ID, earmark.RefID, earmark.EventItemID,
					earmark.UserID+1, earmark.Note,
				),
			)

		_, err := svc.UpdateEventItem(ctx, user.ID, eventItem.RefID, description, nil)
		errs.AssertError(t, err, errs.PermissionDenied, "earmarked by other user")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update archived event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		description := "hodor"

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(eventItem.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_id", "description",
				}).
				AddRow(
					eventItem.ID, eventItem.RefID,
					eventItem.EventID, eventItem.Description,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name,
					event.Description, true,
				),
			)

		_, err := svc.UpdateEventItem(ctx, user.ID, eventItem.RefID, description, nil)
		errs.AssertError(t, err, errs.PermissionDenied, "event is archived")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update other owner should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		description := "hodor"

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(eventItem.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_id", "description",
				}).
				AddRow(
					eventItem.ID, eventItem.RefID,
					eventItem.EventID, eventItem.Description,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID+1, event.Name,
					event.Description, event.Archived,
				),
			)

		_, err := svc.UpdateEventItem(ctx, user.ID, eventItem.RefID, description, nil)
		errs.AssertError(t, err, errs.PermissionDenied, "not event owner")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update no event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		description := "hodor"

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(eventItem.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_id", "description",
				}).
				AddRow(
					eventItem.ID, eventItem.RefID,
					eventItem.EventID, eventItem.Description,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.ID).
			WillReturnError(pgx.ErrNoRows)

		_, err := svc.UpdateEventItem(ctx, user.ID, eventItem.RefID, description, nil)
		errs.AssertError(t, err, errs.NotFound, "event not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update with failing failIfCheck should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		description := "hodor"

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(eventItem.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_id", "description",
				}).
				AddRow(
					eventItem.ID, eventItem.RefID,
					eventItem.EventID, eventItem.Description,
				),
			)

		_, err := svc.UpdateEventItem(
			ctx, user.ID, eventItem.RefID, description,
			func(ei *model.EventItem) bool { return true },
		)
		errs.AssertError(t, err, errs.FailedPrecondition, "extra checks failed")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update no event-item should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		description := "hodor"

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(eventItem.RefID).
			WillReturnError(pgx.ErrNoRows)

		_, err := svc.UpdateEventItem(ctx, user.ID, eventItem.RefID, description, nil)
		errs.AssertError(t, err, errs.NotFound, "event-item not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update bad value should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		_, err := svc.UpdateEventItem(ctx, user.ID, eventItem.RefID, "", nil)
		errs.AssertError(t, err, errs.InvalidArgument, "description bad value")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
