package service

import (
	"context"
	"testing"

	"github.com/dropwhile/refid/v2"
	"github.com/pashagolub/pgxmock/v3"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
)

func TestService_NotifyUsersPendingEvents(t *testing.T) {
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

	t.Run("get user verify should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		refID := refid.Must(model.NewUserVerifyRefID())

		mock.ExpectQuery("^SELECT (.+) FROM user_verify_").
			WithArgs(refID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"ref_id", "user_id"}).
				AddRow(refID, user.ID),
			)

		result, err := svc.GetUserVerifyByRefID(ctx, refID)
		assert.NilError(t, err)
		assert.Equal(t, result.RefID, refID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
