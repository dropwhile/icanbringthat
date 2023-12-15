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

func TestRpc_ListNotifications(t *testing.T) {
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
	notification := &model.Notification{
		ID:           2,
		RefID:        refid.Must(model.NewNotificationRefID()),
		UserID:       user.ID,
		Message:      "",
		Read:         false,
		Created:      tstTs,
		LastModified: tstTs,
	}

	t.Run("list notifications paginated should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := model.SetupDBMock(t, ctx)
		server := &Server{
			Db: mock,
		}
		ctx = auth.ContextSet(ctx, "user", user)

		mock.ExpectQuery("SELECT count(.+) FROM notification_").
			WithArgs(user.ID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{"count"}).
					AddRow(1),
			)
		mock.ExpectQuery("SELECT (.+) FROM notification_").
			WithArgs(pgx.NamedArgs{
				"userID": user.ID,
				"limit":  pgxmock.AnyArg(),
				"offset": pgxmock.AnyArg(),
			}).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id",
						"user_id", "message",
						"read", "created", "last_modified",
					}).
					AddRow(
						notification.ID, notification.RefID,
						user.ID, notification.Message,
						notification.Read, tstTs, tstTs,
					),
			)

		request := &icbt.ListNotificationsRequest{
			Pagination: &icbt.PaginationRequest{
				Limit:  10,
				Offset: 0,
			},
		}
		response, err := server.ListNotifications(ctx, request)
		assert.NilError(t, err)

		assert.Check(t, len(response.Notifications) == 1)
	})
}
