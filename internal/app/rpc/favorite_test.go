package rpc

import (
	"context"
	"testing"

	"github.com/dropwhile/refid/v2"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
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
