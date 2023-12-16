package rpc

import (
	"context"
	"testing"

	"github.com/dropwhile/refid/v2"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/twitchtv/twirp"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/middleware/auth"
	"github.com/dropwhile/icbt/rpc/icbt"
)

func TestRpc_ListFavoriteEvents(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      tstTs,
		LastModified: tstTs,
	}

	t.Run("list favorite events paginated should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := model.SetupDBMock(t, ctx)
		server := &Server{
			Db: mock,
		}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())

		mock.ExpectQuery("SELECT (.+) FROM favorite_").
			WithArgs(user.ID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{"current", "archived"}).
					AddRow(1, 1),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(pgx.NamedArgs{
				"userID":   user.ID,
				"limit":    pgxmock.AnyArg(),
				"offset":   pgxmock.AnyArg(),
				"archived": false,
			}).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id",
						"user_id", "archived",
						"name", "description",
						"start_time", "start_time_tz",
						"created", "last_modified",
					}).
					AddRow(
						1, eventRefID,
						user.ID, false,
						"some name", "some description",
						tstTs, model.Must(model.ParseTimeZone("Etc/UTC")),
						tstTs, tstTs,
					),
			)

		request := &icbt.ListFavoriteEventsRequest{
			Pagination: &icbt.PaginationRequest{Limit: 10, Offset: 0},
			Archived:   func(b bool) *bool { return &b }(false),
		}
		response, err := server.ListFavoriteEvents(ctx, request)
		assert.NilError(t, err)

		assert.Check(t, len(response.Events) == 1)
	})

	t.Run("list favorite events non-paginated should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := model.SetupDBMock(t, ctx)
		server := &Server{
			Db: mock,
		}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())

		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(pgx.NamedArgs{
				"userID":   user.ID,
				"archived": false,
			}).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id",
						"user_id", "archived",
						"name", "description",
						"start_time", "start_time_tz",
						"created", "last_modified",
					}).
					AddRow(
						1, eventRefID,
						user.ID, false,
						"some name", "some description",
						tstTs, model.Must(model.ParseTimeZone("Etc/UTC")),
						tstTs, tstTs,
					),
			)

		request := &icbt.ListFavoriteEventsRequest{
			Archived: func(b bool) *bool { return &b }(false),
		}
		response, err := server.ListFavoriteEvents(ctx, request)
		assert.NilError(t, err)

		assert.Check(t, len(response.Events) == 1)
	})
}

func TestRpc_AddFavorite(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      tstTs,
		LastModified: tstTs,
	}

	t.Run("add favorite should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := model.SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())

		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(eventRefID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id",
						"user_id", "archived",
						"name", "description",
						"start_time", "start_time_tz",
						"created", "last_modified",
					}).
					AddRow(
						1, eventRefID,
						33, false,
						"some name", "some description",
						tstTs, model.Must(model.ParseTimeZone("Etc/UTC")),
						tstTs, tstTs,
					),
			)
		mock.ExpectQuery("SELECT (.+) FROM favorite_").
			WithArgs(pgx.NamedArgs{
				"userID":  user.ID,
				"eventID": pgxmock.AnyArg(),
			}).
			WillReturnError(pgx.ErrNoRows)

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO favorite_").
			WithArgs(pgx.NamedArgs{
				"userID":  user.ID,
				"eventID": pgxmock.AnyArg(),
			}).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "user_id",
						"event_id", "created",
					}).
					AddRow(
						1, user.ID,
						1, tstTs,
					),
			)
		mock.ExpectCommit()
		mock.ExpectRollback()

		request := &icbt.CreateFavoriteRequest{
			EventRefId: eventRefID.String(),
		}
		response, err := server.AddFavorite(ctx, request)
		assert.NilError(t, err)

		assert.Equal(t, response.Favorite.EventRefId, eventRefID.String())
	})

	t.Run("add favorite for own event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := model.SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())

		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(eventRefID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id",
						"user_id", "archived",
						"name", "description",
						"start_time", "start_time_tz",
						"created", "last_modified",
					}).
					AddRow(
						1, eventRefID,
						user.ID, false,
						"some name", "some description",
						tstTs, model.Must(model.ParseTimeZone("Etc/UTC")),
						tstTs, tstTs,
					),
			)

		request := &icbt.CreateFavoriteRequest{
			EventRefId: eventRefID.String(),
		}
		_, err := server.AddFavorite(ctx, request)
		assertTwirpError(t, err, twirp.PermissionDenied, "can't favorite own event")
	})

	t.Run("add favorite for already favorited should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := model.SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())

		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(eventRefID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id",
						"user_id", "archived",
						"name", "description",
						"start_time", "start_time_tz",
						"created", "last_modified",
					}).
					AddRow(
						1, eventRefID,
						33, false,
						"some name", "some description",
						tstTs, model.Must(model.ParseTimeZone("Etc/UTC")),
						tstTs, tstTs,
					),
			)
		mock.ExpectQuery("SELECT (.+) FROM favorite_").
			WithArgs(pgx.NamedArgs{
				"userID":  user.ID,
				"eventID": pgxmock.AnyArg(),
			}).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "user_id",
						"event_id", "created",
					}).
					AddRow(
						1, user.ID,
						1, tstTs,
					),
			)

		request := &icbt.CreateFavoriteRequest{
			EventRefId: eventRefID.String(),
		}
		_, err := server.AddFavorite(ctx, request)
		assertTwirpError(t, err, twirp.AlreadyExists, "favorite already exists")
	})
}
