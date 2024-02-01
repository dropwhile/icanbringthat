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

func TestService_GetFavoriteEventsCount(t *testing.T) {
	t.Parallel()

	t.Run("count with results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		userID := 4
		currentCount := 4
		archivedCount := 2

		mock.ExpectQuery("^SELECT (.+) FROM favorite_").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"current", "archived"}).
				AddRow(currentCount, archivedCount),
			)

		result, err := svc.GetFavoriteEventsCount(ctx, userID)
		assert.NilError(t, err)
		assert.Equal(t, result.Current, currentCount)
		assert.Equal(t, result.Archived, archivedCount)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetFavoriteEventsPaginated(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:       1,
		RefID:    util.Must(model.NewUserRefID()),
		Email:    "user@example.com",
		Name:     "user",
		PWHash:   []byte("00x00"),
		Verified: true,
	}
	event := &model.Event{
		ID:            2,
		RefID:         util.Must(model.NewEventRefID()),
		UserID:        user.ID,
		Name:          "event",
		Description:   "description",
		Archived:      false,
		ItemSortOrder: []int{1, 2, 3},
		StartTime:     tstTs,
		StartTimeTz:   util.Must(ParseTimeZone("Etc/UTC")),
	}

	t.Run("get with results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		limit := 5
		offset := 0
		archived := false
		currentCount := 1
		archivedCount := 3

		mock.ExpectQuery("^SELECT count(.+) FROM favorite_").
			WithArgs(user.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"current", "archived"}).
				AddRow(currentCount, archivedCount),
			)
		mock.ExpectQuery("^SELECT (.+) FROM event_").
			WithArgs(pgx.NamedArgs{
				"userID":   user.ID,
				"limit":    limit,
				"offset":   offset,
				"archived": archived,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description", "archived",
					"item_sort_order", "start_time", "start_time_tz",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name,
					event.Description, event.Archived, event.ItemSortOrder,
					event.StartTime, event.StartTimeTz,
				),
			)

		events, pagination, err := svc.GetFavoriteEventsPaginated(ctx, user.ID, limit, offset, archived)
		assert.NilError(t, err)
		assert.Equal(t, len(events), currentCount)
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

		limit := 5
		offset := 0
		archived := true
		currentCount := 0
		archivedCount := 1

		mock.ExpectQuery("^SELECT count(.+) FROM favorite_").
			WithArgs(user.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"current", "archived"}).
				AddRow(currentCount, archivedCount),
			)
		mock.ExpectQuery("^SELECT (.+) FROM event_").
			WithArgs(pgx.NamedArgs{
				"userID":   user.ID,
				"limit":    limit,
				"offset":   offset,
				"archived": archived,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description", "archived",
					"item_sort_order", "start_time", "start_time_tz",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name,
					event.Description, true, event.ItemSortOrder,
					event.StartTime, event.StartTimeTz,
				),
			)

		events, pagination, err := svc.GetFavoriteEventsPaginated(ctx, user.ID, limit, offset, archived)
		assert.NilError(t, err)
		assert.Equal(t, len(events), archivedCount)
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

		limit := 5
		offset := 0
		archived := false
		currentCount := 0
		archivedCount := 0

		mock.ExpectQuery("^SELECT count(.+) FROM favorite_").
			WithArgs(user.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"current", "archived"}).
				AddRow(currentCount, archivedCount),
			)

		events, pagination, err := svc.GetFavoriteEventsPaginated(ctx, user.ID, limit, offset, archived)
		assert.NilError(t, err)
		assert.Equal(t, len(events), currentCount)
		assert.Equal(t, pagination.Limit, uint32(limit))
		assert.Equal(t, pagination.Offset, uint32(offset))
		assert.Equal(t, pagination.Count, uint32(currentCount))
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetFavoriteEvents(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:       1,
		RefID:    util.Must(model.NewUserRefID()),
		Email:    "user@example.com",
		Name:     "user",
		PWHash:   []byte("00x00"),
		Verified: true,
	}
	event := &model.Event{
		ID:            2,
		RefID:         util.Must(model.NewEventRefID()),
		UserID:        user.ID,
		Name:          "event",
		Description:   "description",
		Archived:      false,
		ItemSortOrder: []int{1, 2, 3},
		StartTime:     tstTs,
		StartTimeTz:   util.Must(ParseTimeZone("Etc/UTC")),
	}

	t.Run("get with results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		archived := false
		count := 1

		mock.ExpectQuery("^SELECT (.+) FROM event_").
			WithArgs(pgx.NamedArgs{
				"userID":   user.ID,
				"archived": archived,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description", "archived",
					"item_sort_order", "start_time", "start_time_tz",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name,
					event.Description, event.Archived, event.ItemSortOrder,
					event.StartTime, event.StartTimeTz,
				),
			)

		events, err := svc.GetFavoriteEvents(ctx, user.ID, archived)
		assert.NilError(t, err)
		assert.Equal(t, len(events), count)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get with archived results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		archived := true
		count := 1

		mock.ExpectQuery("^SELECT (.+) FROM event_").
			WithArgs(pgx.NamedArgs{
				"userID":   user.ID,
				"archived": archived,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description", "archived",
					"item_sort_order", "start_time", "start_time_tz",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name,
					event.Description, true, event.ItemSortOrder,
					event.StartTime, event.StartTimeTz,
				),
			)

		events, err := svc.GetFavoriteEvents(ctx, user.ID, archived)
		assert.NilError(t, err)
		assert.Equal(t, len(events), count)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get with no results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		archived := false
		count := 0

		mock.ExpectQuery("^SELECT (.+) FROM event_").
			WithArgs(pgx.NamedArgs{
				"userID":   user.ID,
				"archived": archived,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description", "archived",
					"item_sort_order", "start_time", "start_time_tz",
				}),
			)

		events, err := svc.GetFavoriteEvents(ctx, user.ID, archived)
		assert.NilError(t, err)
		assert.Equal(t, len(events), count)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetFavoriteByUserEvent(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:       1,
		RefID:    util.Must(model.NewUserRefID()),
		Email:    "user@example.com",
		Name:     "user",
		PWHash:   []byte("00x00"),
		Verified: true,
	}
	event := &model.Event{
		ID:            2,
		RefID:         util.Must(model.NewEventRefID()),
		UserID:        user.ID,
		Name:          "event",
		Description:   "description",
		Archived:      false,
		ItemSortOrder: []int{1, 2, 3},
		StartTime:     tstTs,
		StartTimeTz:   util.Must(ParseTimeZone("Etc/UTC")),
	}

	t.Run("get with result should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM favorite_").
			WithArgs(pgx.NamedArgs{
				"userID":  user.ID,
				"eventID": event.ID,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{"id", "user_id", "event_id"}).
				AddRow(1, user.ID, event.ID),
			)

		favorite, err := svc.GetFavoriteByUserEvent(ctx, user.ID, event.ID)
		assert.NilError(t, err)
		assert.Equal(t, favorite.EventID, event.ID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get with no result should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM favorite_").
			WithArgs(pgx.NamedArgs{
				"userID":  user.ID,
				"eventID": event.ID,
			}).
			WillReturnError(pgx.ErrNoRows)

		_, err := svc.GetFavoriteByUserEvent(ctx, user.ID, event.ID)
		errs.AssertError(t, err, errs.NotFound, "favorite not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_AddFavorite(t *testing.T) {
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

	t.Run("add favorite should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"start_time", "start_time_tz", "created", "last_modified",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID+1, event.Name, event.Description,
					event.StartTime, event.StartTimeTz, ts, ts,
				),
			)
		mock.ExpectQuery("^SELECT (.+) FROM favorite_").
			WithArgs(pgx.NamedArgs{
				"userID":  user.ID,
				"eventID": event.ID,
			}).
			WillReturnError(pgx.ErrNoRows)
		mock.ExpectBegin()
		mock.ExpectQuery("^INSERT INTO favorite_").
			WithArgs(pgx.NamedArgs{
				"userID":  user.ID,
				"eventID": event.ID,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{"id", "user_id", "event_id"}).
				AddRow(1, user.ID, event.ID),
			)
		mock.ExpectCommit()
		mock.ExpectRollback()

		event, err := svc.AddFavorite(ctx, user.ID, event.RefID)
		assert.NilError(t, err)
		assert.Equal(t, event.RefID, event.RefID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("add favorite with missing event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnError(pgx.ErrNoRows)

		_, err := svc.AddFavorite(ctx, user.ID, event.RefID)
		errs.AssertError(t, err, errs.NotFound, "event not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("add favorite to user owned event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
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

		_, err := svc.AddFavorite(ctx, user.ID, event.RefID)
		errs.AssertError(t, err, errs.PermissionDenied, "can't favorite own event")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("add favorite already exists should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"start_time", "start_time_tz", "created", "last_modified",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID+1, event.Name, event.Description,
					event.StartTime, event.StartTimeTz, ts, ts,
				),
			)
		mock.ExpectQuery("^SELECT (.+) FROM favorite_").
			WithArgs(pgx.NamedArgs{
				"userID":  user.ID,
				"eventID": event.ID,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{"id", "user_id", "event_id"}).
				AddRow(1, user.ID, event.ID),
			)

		_, err := svc.AddFavorite(ctx, user.ID, event.RefID)
		errs.AssertError(t, err, errs.AlreadyExists, "favorite already exists")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_RemoveFavorite(t *testing.T) {
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

	t.Run("remove favorite should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		favID := 44

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"start_time", "start_time_tz", "created", "last_modified",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID+1, event.Name, event.Description,
					event.StartTime, event.StartTimeTz, ts, ts,
				),
			)
		mock.ExpectQuery("^SELECT (.+) FROM favorite_").
			WithArgs(pgx.NamedArgs{
				"userID":  user.ID,
				"eventID": event.ID,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{"id", "user_id", "event_id"}).
				AddRow(favID, user.ID, event.ID),
			)
		mock.ExpectBegin()
		mock.ExpectExec("^DELETE FROM favorite_").
			WithArgs(favID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		err := svc.RemoveFavorite(ctx, user.ID, event.RefID)
		assert.NilError(t, err)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("remove favorite with missing event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnError(pgx.ErrNoRows)

		err := svc.RemoveFavorite(ctx, user.ID, event.RefID)
		errs.AssertError(t, err, errs.NotFound, "event not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("remove favorite not exists should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"start_time", "start_time_tz", "created", "last_modified",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID+1, event.Name, event.Description,
					event.StartTime, event.StartTimeTz, ts, ts,
				),
			)
		mock.ExpectQuery("^SELECT (.+) FROM favorite_").
			WithArgs(pgx.NamedArgs{
				"userID":  user.ID,
				"eventID": event.ID,
			}).
			WillReturnError(pgx.ErrNoRows)

		err := svc.RemoveFavorite(ctx, user.ID, event.RefID)
		errs.AssertError(t, err, errs.NotFound, "favorite not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
