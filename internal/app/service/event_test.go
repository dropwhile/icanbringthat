package service

import (
	"context"
	"testing"
	"time"

	"github.com/dropwhile/refid/v2"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/samber/mo"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/util"
)

func TestService_GetEvent(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:           1,
		RefID:        refid.Must(model.NewEventRefID()),
		UserID:       user.ID,
		Name:         "event",
		Description:  "description",
		Archived:     false,
		StartTime:    ts,
		StartTimeTz:  util.Must(ParseTimeZone("Etc/UTC")),
		Created:      ts,
		LastModified: ts,
	}

	t.Run("get with results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM event_").
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

		result, err := svc.GetEvent(ctx, event.RefID)
		assert.NilError(t, err)
		assert.Equal(t, result.RefID, event.RefID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get with not results should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM event_").
			WithArgs(event.RefID).
			WillReturnError(pgx.ErrNoRows)

		_, err := svc.GetEvent(ctx, event.RefID)
		errs.AssertError(t, err, errs.NotFound, "event not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetEventByID(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:           1,
		RefID:        refid.Must(model.NewEventRefID()),
		UserID:       user.ID,
		Name:         "event",
		Description:  "description",
		Archived:     false,
		StartTime:    ts,
		StartTimeTz:  util.Must(ParseTimeZone("Etc/UTC")),
		Created:      ts,
		LastModified: ts,
	}

	t.Run("get with results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM event_").
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

		result, err := svc.GetEventByID(ctx, event.ID)
		assert.NilError(t, err)
		assert.Equal(t, result.RefID, event.RefID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get with not results should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM event_").
			WithArgs(event.ID).
			WillReturnError(pgx.ErrNoRows)

		_, err := svc.GetEventByID(ctx, event.ID)
		errs.AssertError(t, err, errs.NotFound, "event not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetEventsByIDs(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:           1,
		RefID:        refid.Must(model.NewEventRefID()),
		UserID:       user.ID,
		Name:         "event",
		Description:  "description",
		Archived:     false,
		StartTime:    ts,
		StartTimeTz:  util.Must(ParseTimeZone("Etc/UTC")),
		Created:      ts,
		LastModified: ts,
	}

	t.Run("get with results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM event_").
			WithArgs([]int{event.ID}).
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

		result, err := svc.GetEventsByIDs(ctx, []int{event.ID})
		assert.NilError(t, err)
		assert.Equal(t, len(result), 1)
		assert.Equal(t, result[0].RefID, event.RefID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get with no results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM event_").
			WithArgs([]int{event.ID}).
			WillReturnError(pgx.ErrNoRows)

		result, err := svc.GetEventsByIDs(ctx, []int{event.ID})
		assert.NilError(t, err)
		assert.Equal(t, len(result), 0)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get with empty input should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		result, err := svc.GetEventsByIDs(ctx, []int{})
		assert.NilError(t, err)
		assert.Equal(t, len(result), 0)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_DeleteEvent(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:           1,
		RefID:        refid.Must(model.NewEventRefID()),
		UserID:       user.ID,
		Name:         "event",
		Description:  "description",
		Archived:     false,
		StartTime:    ts,
		StartTimeTz:  util.Must(ParseTimeZone("Etc/UTC")),
		Created:      ts,
		LastModified: ts,
	}

	t.Run("delete should succeed", func(t *testing.T) {
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
		mock.ExpectExec("DELETE FROM event_ ").
			WithArgs(event.ID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		err := svc.DeleteEvent(ctx, user.ID, event.RefID)
		assert.NilError(t, err)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete not owner should fail", func(t *testing.T) {
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

		err := svc.DeleteEvent(ctx, user.ID, event.RefID)
		errs.AssertError(t, err, errs.PermissionDenied, "permission denied")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete not exist should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.RefID).
			WillReturnError(pgx.ErrNoRows)

		err := svc.DeleteEvent(ctx, user.ID, event.RefID)
		errs.AssertError(t, err, errs.NotFound, "event not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_UpdateEvent(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:           1,
		RefID:        refid.Must(model.NewEventRefID()),
		UserID:       user.ID,
		Name:         "event",
		Description:  "description",
		Archived:     false,
		StartTime:    ts,
		StartTimeTz:  util.Must(ParseTimeZone("Etc/UTC")),
		Created:      ts,
		LastModified: ts,
	}

	t.Run("update name should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		euvs := &EventUpdateValues{
			Name:          mo.Some("new-name"),
			Description:   mo.None[string](),
			ItemSortOrder: mo.None[[]int](),
			StartTime:     mo.None[time.Time](),
			Tz:            mo.None[string](),
		}

		startTimeTz, terr := util.OptionMapConvert(
			euvs.Tz, func(input string) (*model.TimeZone, error) {
				return ParseTimeZone(input)
			},
		)
		assert.NilError(t, terr)

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
		mock.ExpectExec("UPDATE event_ ").
			WithArgs(pgx.NamedArgs{
				"name":          euvs.Name,
				"description":   euvs.Description,
				"startTime":     euvs.StartTime,
				"startTimeTz":   startTimeTz,
				"itemSortOrder": euvs.ItemSortOrder,
				"eventID":       event.ID,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		err := svc.UpdateEvent(ctx, user.ID, event.RefID, euvs)
		assert.NilError(t, err)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update description should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		euvs := &EventUpdateValues{
			Name:          mo.None[string](),
			Description:   mo.Some("some desc"),
			ItemSortOrder: mo.None[[]int](),
			StartTime:     mo.None[time.Time](),
			Tz:            mo.None[string](),
		}

		startTimeTz, terr := util.OptionMapConvert(
			euvs.Tz, func(input string) (*model.TimeZone, error) {
				return ParseTimeZone(input)
			},
		)
		assert.NilError(t, terr)

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
		mock.ExpectExec("UPDATE event_ ").
			WithArgs(pgx.NamedArgs{
				"name":          euvs.Name,
				"description":   euvs.Description,
				"startTime":     euvs.StartTime,
				"startTimeTz":   startTimeTz,
				"itemSortOrder": euvs.ItemSortOrder,
				"eventID":       event.ID,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		err := svc.UpdateEvent(ctx, user.ID, event.RefID, euvs)
		assert.NilError(t, err)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update tz should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		euvs := &EventUpdateValues{
			Name:          mo.None[string](),
			Description:   mo.None[string](),
			ItemSortOrder: mo.None[[]int](),
			StartTime:     mo.None[time.Time](),
			Tz:            mo.Some("US/Pacific"),
		}

		startTimeTz, terr := util.OptionMapConvert(
			euvs.Tz, func(input string) (*model.TimeZone, error) {
				return ParseTimeZone(input)
			},
		)
		assert.NilError(t, terr)

		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived", "start_time", "start_time_tz",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name,
					event.Description, event.Archived,
					event.StartTime, event.StartTimeTz,
				),
			)
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE event_ ").
			WithArgs(pgx.NamedArgs{
				"name":          euvs.Name,
				"description":   euvs.Description,
				"startTime":     euvs.StartTime,
				"startTimeTz":   startTimeTz,
				"itemSortOrder": euvs.ItemSortOrder,
				"eventID":       event.ID,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		err := svc.UpdateEvent(ctx, user.ID, event.RefID, euvs)
		assert.NilError(t, err)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update starttime should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		euvs := &EventUpdateValues{
			Name:          mo.None[string](),
			Description:   mo.None[string](),
			ItemSortOrder: mo.None[[]int](),
			StartTime:     mo.Some(tstTs),
			Tz:            mo.None[string](),
		}

		startTimeTz, terr := util.OptionMapConvert(
			euvs.Tz, func(input string) (*model.TimeZone, error) {
				return ParseTimeZone(input)
			},
		)
		assert.NilError(t, terr)

		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived", "start_time", "start_time_tz",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name,
					event.Description, event.Archived,
					event.StartTime, event.StartTimeTz,
				),
			)
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE event_ ").
			WithArgs(pgx.NamedArgs{
				"name":          euvs.Name,
				"description":   euvs.Description,
				"startTime":     euvs.StartTime,
				"startTimeTz":   startTimeTz,
				"itemSortOrder": euvs.ItemSortOrder,
				"eventID":       event.ID,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		err := svc.UpdateEvent(ctx, user.ID, event.RefID, euvs)
		assert.NilError(t, err)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update item-sort-order should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		euvs := &EventUpdateValues{
			Name:          mo.None[string](),
			Description:   mo.None[string](),
			ItemSortOrder: mo.Some([]int{1, 2, 3}),
			StartTime:     mo.None[time.Time](),
			Tz:            mo.None[string](),
		}

		startTimeTz, terr := util.OptionMapConvert(
			euvs.Tz, func(input string) (*model.TimeZone, error) {
				return ParseTimeZone(input)
			},
		)
		assert.NilError(t, terr)

		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived", "start_time", "start_time_tz",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name,
					event.Description, event.Archived,
					event.StartTime, event.StartTimeTz,
				),
			)
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE event_ ").
			WithArgs(pgx.NamedArgs{
				"name":          euvs.Name,
				"description":   euvs.Description,
				"startTime":     euvs.StartTime,
				"startTimeTz":   startTimeTz,
				"itemSortOrder": euvs.ItemSortOrder,
				"eventID":       event.ID,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		err := svc.UpdateEvent(ctx, user.ID, event.RefID, euvs)
		assert.NilError(t, err)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update no change should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		euvs := &EventUpdateValues{
			Name:          mo.None[string](),
			Description:   mo.None[string](),
			ItemSortOrder: mo.None[[]int](),
			StartTime:     mo.None[time.Time](),
			Tz:            mo.None[string](),
		}

		err := svc.UpdateEvent(ctx, user.ID, event.RefID, euvs)
		errs.AssertError(t, err, errs.InvalidArgument, "missing fields")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update archived event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		euvs := &EventUpdateValues{
			Name:          mo.Some("new-name"),
			Description:   mo.None[string](),
			ItemSortOrder: mo.None[[]int](),
			StartTime:     mo.None[time.Time](),
			Tz:            mo.None[string](),
		}

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

		err := svc.UpdateEvent(ctx, user.ID, event.RefID, euvs)
		errs.AssertError(t, err, errs.PermissionDenied, "event is archived")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update not owner should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		euvs := &EventUpdateValues{
			Name:          mo.Some("new-name"),
			Description:   mo.None[string](),
			ItemSortOrder: mo.None[[]int](),
			StartTime:     mo.None[time.Time](),
			Tz:            mo.None[string](),
		}

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

		err := svc.UpdateEvent(ctx, user.ID, event.RefID, euvs)
		errs.AssertError(t, err, errs.PermissionDenied, "permission denied")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update no event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		euvs := &EventUpdateValues{
			Name:          mo.Some("new-name"),
			Description:   mo.None[string](),
			ItemSortOrder: mo.None[[]int](),
			StartTime:     mo.None[time.Time](),
			Tz:            mo.None[string](),
		}

		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.RefID).
			WillReturnError(pgx.ErrNoRows)

		err := svc.UpdateEvent(ctx, user.ID, event.RefID, euvs)
		errs.AssertError(t, err, errs.NotFound, "event not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update bad timezone should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		euvs := &EventUpdateValues{
			Name:          mo.None[string](),
			Description:   mo.None[string](),
			ItemSortOrder: mo.None[[]int](),
			StartTime:     mo.None[time.Time](),
			Tz:            mo.Some("hodor"),
		}

		err := svc.UpdateEvent(ctx, user.ID, event.RefID, euvs)
		errs.AssertError(t, err, errs.InvalidArgument, "Tz bad value")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update empty timezone should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		euvs := &EventUpdateValues{
			Name:          mo.None[string](),
			Description:   mo.None[string](),
			ItemSortOrder: mo.None[[]int](),
			StartTime:     mo.None[time.Time](),
			Tz:            mo.Some(""),
		}

		err := svc.UpdateEvent(ctx, user.ID, event.RefID, euvs)
		errs.AssertError(t, err, errs.InvalidArgument, "Tz bad value")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update bad name should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		euvs := &EventUpdateValues{
			Name:          mo.Some(""),
			Description:   mo.None[string](),
			ItemSortOrder: mo.None[[]int](),
			StartTime:     mo.None[time.Time](),
			Tz:            mo.None[string](),
		}

		err := svc.UpdateEvent(ctx, user.ID, event.RefID, euvs)
		errs.AssertError(t, err, errs.InvalidArgument, "Name bad value")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update bad description should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		euvs := &EventUpdateValues{
			Name:          mo.None[string](),
			Description:   mo.Some(""),
			ItemSortOrder: mo.None[[]int](),
			StartTime:     mo.None[time.Time](),
			Tz:            mo.None[string](),
		}

		err := svc.UpdateEvent(ctx, user.ID, event.RefID, euvs)
		errs.AssertError(t, err, errs.InvalidArgument, "Description bad value")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update bad starttime should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		euvs := &EventUpdateValues{
			Name:          mo.None[string](),
			Description:   mo.None[string](),
			ItemSortOrder: mo.None[[]int](),
			StartTime:     mo.Some(time.Time{}),
			Tz:            mo.None[string](),
		}

		err := svc.UpdateEvent(ctx, user.ID, event.RefID, euvs)
		errs.AssertError(t, err, errs.InvalidArgument, "start_time bad value")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update bad item sort order should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		euvs := &EventUpdateValues{
			Name:          mo.None[string](),
			Description:   mo.None[string](),
			ItemSortOrder: mo.Some([]int{}),
			StartTime:     mo.None[time.Time](),
			Tz:            mo.None[string](),
		}

		err := svc.UpdateEvent(ctx, user.ID, event.RefID, euvs)
		errs.AssertError(t, err, errs.InvalidArgument, "ItemSortOrder bad value")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_UpdateEventItemSorting(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:           1,
		RefID:        refid.Must(model.NewEventRefID()),
		UserID:       user.ID,
		Name:         "event",
		Description:  "description",
		Archived:     false,
		StartTime:    ts,
		StartTimeTz:  util.Must(ParseTimeZone("Etc/UTC")),
		Created:      ts,
		LastModified: ts,
	}

	t.Run("update should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		itemSortOrder := []int{5, 4, 3}

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
		mock.ExpectExec("UPDATE event_ ").
			WithArgs(pgx.NamedArgs{
				"name":          mo.None[string](),
				"description":   mo.None[string](),
				"startTime":     mo.None[time.Time](),
				"startTimeTz":   mo.None[*model.TimeZone](),
				"itemSortOrder": mo.Some(itemSortOrder),
				"eventID":       event.ID,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		result, err := svc.UpdateEventItemSorting(ctx, user.ID, event.RefID, itemSortOrder)
		assert.NilError(t, err)
		assert.DeepEqual(t, result.ItemSortOrder, itemSortOrder)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update no change should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		itemSortOrder := []int{5, 4, 3}

		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived", "item_sort_order",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name,
					event.Description, event.Archived,
					itemSortOrder,
				),
			)

		_, err := svc.UpdateEventItemSorting(ctx, user.ID, event.RefID, itemSortOrder)
		errs.AssertError(t, err, errs.FailedPrecondition, "no changes")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update archived event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		itemSortOrder := []int{5, 4, 3}

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

		_, err := svc.UpdateEventItemSorting(ctx, user.ID, event.RefID, itemSortOrder)
		errs.AssertError(t, err, errs.PermissionDenied, "event is archived")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update not owner should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		itemSortOrder := []int{5, 4, 3}

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

		_, err := svc.UpdateEventItemSorting(ctx, user.ID, event.RefID, itemSortOrder)
		errs.AssertError(t, err, errs.PermissionDenied, "permission denied")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update no event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		itemSortOrder := []int{5, 4, 3}

		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(event.RefID).
			WillReturnError(pgx.ErrNoRows)

		_, err := svc.UpdateEventItemSorting(ctx, user.ID, event.RefID, itemSortOrder)
		errs.AssertError(t, err, errs.NotFound, "event not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_CreateEvent(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:           1,
		RefID:        refid.Must(model.NewEventRefID()),
		UserID:       user.ID,
		Name:         "event",
		Description:  "description",
		Archived:     false,
		StartTime:    ts,
		StartTimeTz:  util.Must(ParseTimeZone("Etc/UTC")),
		Created:      ts,
		LastModified: ts,
	}

	t.Run("create should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO event_ ").
			WithArgs(pgx.NamedArgs{
				"refID":       EventRefIDMatcher,
				"userID":      event.UserID,
				"name":        event.Name,
				"description": event.Description,
				"startTime":   event.StartTime,
				"startTimeTz": event.StartTimeTz,
			}).
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
		mock.ExpectCommit()
		mock.ExpectRollback()

		result, err := svc.CreateEvent(
			ctx, user, event.Name, event.Description, event.StartTime,
			event.StartTimeTz.String(),
		)
		assert.NilError(t, err)
		assert.Equal(t, result.RefID, event.RefID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create with bad tz should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		_, err := svc.CreateEvent(
			ctx, user, event.Name, event.Description, event.StartTime,
			"hodor",
		)
		errs.AssertError(t, err, errs.InvalidArgument, "tz bad value")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create with bad description should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		_, err := svc.CreateEvent(
			ctx, user, event.Name, "", event.StartTime,
			event.StartTimeTz.String(),
		)
		errs.AssertError(t, err, errs.InvalidArgument, "description bad value")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create with bad name should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		_, err := svc.CreateEvent(
			ctx, user, "", event.Description, event.StartTime,
			event.StartTimeTz.String(),
		)
		errs.AssertError(t, err, errs.InvalidArgument, "name bad value")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create with starttime (zero time) should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		_, err := svc.CreateEvent(
			ctx, user, event.Name, event.Description, time.Time{},
			event.StartTimeTz.String(),
		)
		errs.AssertError(t, err, errs.InvalidArgument, "start_time bad value")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create with starttime (before unix epoch) should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		_, err := svc.CreateEvent(
			ctx, user, event.Name, event.Description, time.Unix(0, 0).UTC(),
			event.StartTimeTz.String(),
		)
		errs.AssertError(t, err, errs.InvalidArgument, "start_time bad value")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetEventsPaginated(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:           1,
		RefID:        refid.Must(model.NewEventRefID()),
		UserID:       user.ID,
		Name:         "event",
		Description:  "description",
		Archived:     false,
		StartTime:    ts,
		StartTimeTz:  util.Must(ParseTimeZone("Etc/UTC")),
		Created:      ts,
		LastModified: ts,
	}

	t.Run("get with results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		userID := 4
		limit := 5
		offset := 0
		archived := false
		currentCount := 1
		archivedCount := 3

		mock.ExpectQuery("^SELECT count(.+) FROM event_").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"current", "archived"}).
				AddRow(currentCount, archivedCount),
			)
		mock.ExpectQuery("^SELECT (.+) FROM event_").
			WithArgs(pgx.NamedArgs{
				"userID":   userID,
				"limit":    limit,
				"offset":   offset,
				"archived": archived,
			}).
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

		results, pagination, err := svc.GetEventsPaginated(ctx, userID, limit, offset, archived)
		assert.NilError(t, err)
		assert.Equal(t, len(results), currentCount)
		assert.Equal(t, pagination.Limit, uint32(limit))
		assert.Equal(t, pagination.Offset, uint32(offset))
		assert.Equal(t, pagination.Count, uint32(currentCount))
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
		currentCount := 3
		archivedCount := 1

		mock.ExpectQuery("^SELECT count(.+) FROM event_").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"current", "archived"}).
				AddRow(currentCount, archivedCount),
			)
		mock.ExpectQuery("^SELECT (.+) FROM event_").
			WithArgs(pgx.NamedArgs{
				"userID":   userID,
				"limit":    limit,
				"offset":   offset,
				"archived": archived,
			}).
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

		results, pagination, err := svc.GetEventsPaginated(ctx, userID, limit, offset, archived)
		assert.NilError(t, err)
		assert.Equal(t, len(results), archivedCount)
		assert.Equal(t, pagination.Limit, uint32(limit))
		assert.Equal(t, pagination.Offset, uint32(offset))
		assert.Equal(t, pagination.Count, uint32(archivedCount))
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

		mock.ExpectQuery("^SELECT count(.+) FROM event_").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"current", "archived"}).
				AddRow(0, 0),
			)

		results, pagination, err := svc.GetEventsPaginated(ctx, userID, limit, offset, archived)
		assert.NilError(t, err)
		assert.Equal(t, len(results), 0)
		assert.Equal(t, pagination.Limit, uint32(limit))
		assert.Equal(t, pagination.Offset, uint32(offset))
		assert.Equal(t, pagination.Count, uint32(0))
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetEventsComingSoonPaginated(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:           1,
		RefID:        refid.Must(model.NewEventRefID()),
		UserID:       user.ID,
		Name:         "event",
		Description:  "description",
		Archived:     false,
		StartTime:    ts,
		StartTimeTz:  util.Must(ParseTimeZone("Etc/UTC")),
		Created:      ts,
		LastModified: ts,
	}

	t.Run("get with results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		userID := 4
		limit := 5
		offset := 0
		currentCount := 1
		archivedCount := 3

		mock.ExpectQuery("^SELECT count(.+) FROM event_").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"current", "archived"}).
				AddRow(currentCount, archivedCount),
			)
		mock.ExpectQuery("^SELECT (.+) FROM event_").
			WithArgs(pgx.NamedArgs{
				"userID": userID,
				"limit":  limit,
				"offset": offset,
			}).
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

		results, pagination, err := svc.GetEventsComingSoonPaginated(ctx, userID, limit, offset)
		assert.NilError(t, err)
		assert.Equal(t, len(results), currentCount)
		assert.Equal(t, pagination.Limit, uint32(limit))
		assert.Equal(t, pagination.Offset, uint32(offset))
		assert.Equal(t, pagination.Count, uint32(currentCount))
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

		mock.ExpectQuery("^SELECT count(.+) FROM event_").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"current", "archived"}).
				AddRow(0, 0),
			)

		results, pagination, err := svc.GetEventsComingSoonPaginated(ctx, userID, limit, offset)
		assert.NilError(t, err)
		assert.Equal(t, len(results), 0)
		assert.Equal(t, pagination.Limit, uint32(limit))
		assert.Equal(t, pagination.Offset, uint32(offset))
		assert.Equal(t, pagination.Count, uint32(0))
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetEventsCount(t *testing.T) {
	t.Parallel()

	t.Run("get with results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		userID := 4
		currentCount := 1
		archivedCount := 3

		mock.ExpectQuery("^SELECT count(.+) FROM event_").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"current", "archived"}).
				AddRow(currentCount, archivedCount),
			)

		results, err := svc.GetEventsCount(ctx, userID)
		assert.NilError(t, err)
		assert.Equal(t, results.Archived, archivedCount)
		assert.Equal(t, results.Current, currentCount)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetEvents(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:           1,
		RefID:        refid.Must(model.NewEventRefID()),
		UserID:       user.ID,
		Name:         "event",
		Description:  "description",
		Archived:     false,
		StartTime:    ts,
		StartTimeTz:  util.Must(ParseTimeZone("Etc/UTC")),
		Created:      ts,
		LastModified: ts,
	}

	t.Run("get with results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		userID := 4
		archived := false
		count := 1

		mock.ExpectQuery("^SELECT (.+) FROM event_").
			WithArgs(pgx.NamedArgs{
				"userID":   userID,
				"archived": archived,
			}).
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

		results, err := svc.GetEvents(ctx, userID, archived)
		assert.NilError(t, err)
		assert.Equal(t, len(results), count)
		assert.Equal(t, results[0].RefID, event.RefID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get with empty results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		userID := 4
		archived := false
		count := 0

		mock.ExpectQuery("^SELECT (.+) FROM event_").
			WithArgs(pgx.NamedArgs{
				"userID":   userID,
				"archived": archived,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived",
				}),
			)

		results, err := svc.GetEvents(ctx, userID, archived)
		assert.NilError(t, err)
		assert.Equal(t, len(results), count)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
