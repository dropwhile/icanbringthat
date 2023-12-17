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

func TestServer_ListEvents(t *testing.T) {
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

	t.Run("list events paginated should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())

		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(user.ID).
			WillReturnRows(pgxmock.NewRows(
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
			WillReturnRows(pgxmock.NewRows(
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

		request := &icbt.ListEventsRequest{
			Pagination: &icbt.PaginationRequest{Limit: 10, Offset: 0},
			Archived:   func(b bool) *bool { return &b }(false),
		}
		response, err := server.ListEvents(ctx, request)
		assert.NilError(t, err)

		assert.Check(t, len(response.Events) == 1)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("list events non-paginated should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
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
			WillReturnRows(pgxmock.NewRows(
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

		request := &icbt.ListEventsRequest{
			Archived: func(b bool) *bool { return &b }(false),
		}
		response, err := server.ListEvents(ctx, request)
		assert.NilError(t, err)

		assert.Check(t, len(response.Events) == 1)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestServer_GetEventDetails(t *testing.T) {
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

	t.Run("get event details should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventID := 1
		eventItemID := 3
		earmarkID := 4

		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(eventRefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"user_id", "archived",
					"name", "description",
					"start_time", "start_time_tz",
					"created", "last_modified",
				}).
				AddRow(
					eventID, eventRefID,
					user.ID, false,
					"some name", "some description",
					tstTs, model.Must(model.ParseTimeZone("Etc/UTC")),
					tstTs, tstTs,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_item_").
			WithArgs(eventID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"event_id", "description",
					"created", "last_modified",
				}).
				AddRow(
					eventItemID, refid.Must(model.NewEventItemRefID()),
					eventID, "some description",
					tstTs, tstTs,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM earmark_").
			WithArgs(eventID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id",
					"event_item_id", "note",
					"created", "last_modified",
				}).
				AddRow(
					earmarkID, refid.Must(model.NewEarmarkRefID()),
					user.ID, eventItemID, "some note",
					tstTs, tstTs,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_item_").
			WithArgs(eventItemID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"event_id", "description",
					"created", "last_modified",
				}).
				AddRow(
					eventItemID, refid.Must(model.NewEventItemRefID()),
					eventID, "some description",
					tstTs, tstTs,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM user_").
			WithArgs(user.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "email", "name",
					"created", "last_modified",
				}).
				AddRow(
					user.ID, user.RefID, "user@example.com", "user",
					tstTs, tstTs,
				),
			)

		request := &icbt.GetEventDetailsRequest{
			RefId: eventRefID.String(),
		}
		response, err := server.GetEventDetails(ctx, request)
		assert.NilError(t, err)

		assert.Equal(t, response.Event.Name, "some name")
		assert.Equal(t, len(response.Items), 1)
		assert.Equal(t, len(response.Earmarks), 1)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get event details event not found should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())

		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(eventRefID).
			WillReturnError(pgx.ErrNoRows)

		request := &icbt.GetEventDetailsRequest{
			RefId: eventRefID.String(),
		}
		_, err := server.GetEventDetails(ctx, request)
		assertTwirpError(t, err, twirp.NotFound, "event not found")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
	t.Run("get event details with bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.GetEventDetailsRequest{
			RefId: "hodor",
		}
		_, err := server.GetEventDetails(ctx, request)
		assertTwirpError(t, err, twirp.InvalidArgument, "ref_id bad event ref-id")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestServer_CreateEvent(t *testing.T) {
}

func TestServer_UpdateEvent(t *testing.T) {
}

func TestServer_DeleteEvent(t *testing.T) {
}
