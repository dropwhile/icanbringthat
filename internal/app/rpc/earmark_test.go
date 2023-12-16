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

func TestRpc_ListEarmarks(t *testing.T) {
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

	t.Run("list earmarks paginated should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{
			Db: mock,
		}
		ctx = auth.ContextSet(ctx, "user", user)
		earmarkRefID := refid.Must(model.NewEarmarkRefID())

		mock.ExpectQuery("SELECT count(.+) FROM earmark_").
			WithArgs(user.ID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{"current", "archived"}).
					AddRow(1, 1),
			)

		mock.ExpectQuery("SELECT (.+) FROM earmark_").
			WithArgs(pgx.NamedArgs{
				"userID":   user.ID,
				"limit":    pgxmock.AnyArg(),
				"offset":   pgxmock.AnyArg(),
				"archived": false,
			}).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "user_id",
						"event_item_id", "note",
						"created", "last_modified",
					}).
					AddRow(
						1, earmarkRefID, user.ID,
						12, "some note",
						tstTs, tstTs,
					),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_item_").
			WithArgs(12).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id",
						"event_id", "description",
						"created", "last_modified",
					}).
					AddRow(
						12, refid.Must(model.NewEventItemRefID()),
						10, "some description",
						tstTs, tstTs,
					),
			)
		mock.ExpectQuery("SELECT (.+) FROM user_").
			WithArgs(user.ID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "email", "name",
						"created", "last_modified",
					}).
					AddRow(
						user.ID, user.RefID, "user@example.com", "user",
						tstTs, tstTs,
					),
			)

		request := &icbt.ListEarmarksRequest{
			Pagination: &icbt.PaginationRequest{Limit: 10, Offset: 0},
			Archived:   func(b bool) *bool { return &b }(false),
		}
		response, err := server.ListEarmarks(ctx, request)
		assert.NilError(t, err)

		assert.Check(t, len(response.Earmarks) == 1)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("list earmarks non-paginated should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{
			Db: mock,
		}
		ctx = auth.ContextSet(ctx, "user", user)
		earmarkRefID := refid.Must(model.NewEarmarkRefID())

		mock.ExpectQuery("SELECT (.+) FROM earmark_").
			WithArgs(pgx.NamedArgs{
				"userID":   user.ID,
				"archived": false,
			}).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "user_id",
						"event_item_id", "note",
						"created", "last_modified",
					}).
					AddRow(
						1, earmarkRefID, user.ID,
						12, "some note",
						tstTs, tstTs,
					),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_item_").
			WithArgs(12).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id",
						"event_id", "description",
						"created", "last_modified",
					}).
					AddRow(
						12, refid.Must(model.NewEventItemRefID()),
						10, "some description",
						tstTs, tstTs,
					),
			)
		mock.ExpectQuery("SELECT (.+) FROM user_").
			WithArgs(user.ID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "email", "name",
						"created", "last_modified",
					}).
					AddRow(
						user.ID, user.RefID, "user@example.com", "user",
						tstTs, tstTs,
					),
			)

		request := &icbt.ListEarmarksRequest{
			Archived: func(b bool) *bool { return &b }(false),
		}
		response, err := server.ListEarmarks(ctx, request)
		assert.NilError(t, err)

		assert.Check(t, len(response.Earmarks) == 1)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestRpc_AddEarmark(t *testing.T) {
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

	t.Run("add earmark should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		earmarkRefID := refid.Must(model.NewEarmarkRefID())
		eventItemID := 33
		eventID := 22

		mock.ExpectQuery("SELECT (.+) FROM event_item_").
			WithArgs(eventItemRefID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id",
						"event_id", "description",
						"created", "last_modified",
					}).
					AddRow(
						eventItemID, eventItemRefID,
						eventID, "some description",
						tstTs, tstTs,
					),
			)
		mock.ExpectQuery("SELECT (.+) FROM earmark_").
			WithArgs(eventItemID).
			WillReturnError(pgx.ErrNoRows)

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO earmark_").
			WithArgs(pgx.NamedArgs{
				"userID":      user.ID,
				"eventItemID": eventItemID,
				"refID":       model.EarmarkRefIDMatcher,
				"note":        "some note",
			}).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "user_id",
						"event_item_id", "note",
						"created", "last_modified",
					}).
					AddRow(
						1, earmarkRefID, user.ID,
						eventItemID, "some note",
						tstTs, tstTs,
					),
			)
		mock.ExpectCommit()
		mock.ExpectRollback()
		mock.ExpectQuery("SELECT (.+) FROM event_item_").
			WithArgs(eventItemID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id",
						"event_id", "description",
						"created", "last_modified",
					}).
					AddRow(
						eventItemID, eventItemRefID,
						eventID, "some description",
						tstTs, tstTs,
					),
			)
		mock.ExpectQuery("SELECT (.+) FROM user_").
			WithArgs(user.ID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "email", "name",
						"created", "last_modified",
					}).
					AddRow(
						user.ID, user.RefID, "user@example.com", "user",
						tstTs, tstTs,
					),
			)

		request := &icbt.CreateEarmarkRequest{
			EventItemRefId: eventItemRefID.String(),
			Note:           "some note",
		}
		response, err := server.CreateEarmark(ctx, request)
		assert.NilError(t, err)

		assert.Equal(t, response.Earmark.RefId, earmarkRefID.String())
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("add earmark for already earmarked by self should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		earmarkRefID := refid.Must(model.NewEarmarkRefID())
		eventItemID := 33
		eventID := 22

		mock.ExpectQuery("SELECT (.+) FROM event_item_").
			WithArgs(eventItemRefID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id",
						"event_id", "description",
						"created", "last_modified",
					}).
					AddRow(
						eventItemID, eventItemRefID,
						eventID, "some description",
						tstTs, tstTs,
					),
			)
		mock.ExpectQuery("SELECT (.+) FROM earmark_").
			WithArgs(eventItemID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "user_id",
						"event_item_id", "note",
						"created", "last_modified",
					}).
					AddRow(
						1, earmarkRefID, user.ID,
						eventItemID, "some note",
						tstTs, tstTs,
					),
			)

		request := &icbt.CreateEarmarkRequest{
			EventItemRefId: eventItemRefID.String(),
			Note:           "some note",
		}
		_, err := server.CreateEarmark(ctx, request)
		assertTwirpError(t, err, twirp.PermissionDenied, "already earmarked")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("add earmark for already earmarked by other should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		earmarkRefID := refid.Must(model.NewEarmarkRefID())
		eventItemID := 33
		eventID := 22

		mock.ExpectQuery("SELECT (.+) FROM event_item_").
			WithArgs(eventItemRefID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id",
						"event_id", "description",
						"created", "last_modified",
					}).
					AddRow(
						eventItemID, eventItemRefID,
						eventID, "some description",
						tstTs, tstTs,
					),
			)
		mock.ExpectQuery("SELECT (.+) FROM earmark_").
			WithArgs(eventItemID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "user_id",
						"event_item_id", "note",
						"created", "last_modified",
					}).
					AddRow(
						1, earmarkRefID, 44,
						eventItemID, "some note",
						tstTs, tstTs,
					),
			)

		request := &icbt.CreateEarmarkRequest{
			EventItemRefId: eventItemRefID.String(),
			Note:           "some note",
		}
		_, err := server.CreateEarmark(ctx, request)
		assertTwirpError(t, err, twirp.PermissionDenied, "already earmarked by other user")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("add earmark with bad event item refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.CreateEarmarkRequest{
			EventItemRefId: "hodor",
			Note:           "some note",
		}
		_, err := server.CreateEarmark(ctx, request)
		assertTwirpError(t, err, twirp.InvalidArgument, "ref_id bad event-item ref-id")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestRpc_RemoveEarmark(t *testing.T) {
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

	t.Run("remove earmark should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		earmarkRefID := refid.Must(model.NewEarmarkRefID())
		eventItemID := 33
		earmarkID := 5
		eventID := 22

		mock.ExpectQuery("SELECT (.+) FROM earmark_").
			WithArgs(earmarkRefID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "user_id",
						"event_item_id", "note",
						"created", "last_modified",
					}).
					AddRow(
						earmarkID, earmarkRefID, user.ID,
						eventItemID, "some note",
						tstTs, tstTs,
					),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(earmarkID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "user_id", "name", "description",
						"archived", "created", "last_modified",
					}).
					AddRow(
						eventID, refid.Must(model.NewEventRefID()), user.ID,
						"event name", "event desc",
						false, tstTs, tstTs,
					),
			)
		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM earmark_").
			WithArgs(earmarkID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		request := &icbt.RemoveEarmarkRequest{
			RefId: earmarkRefID.String(),
		}
		_, err := server.RemoveEarmark(ctx, request)
		assert.NilError(t, err)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("remove earmark for another user should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		earmarkRefID := refid.Must(model.NewEarmarkRefID())
		eventItemID := 33
		earmarkID := 5

		mock.ExpectQuery("SELECT (.+) FROM earmark_").
			WithArgs(earmarkRefID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "user_id",
						"event_item_id", "note",
						"created", "last_modified",
					}).
					AddRow(
						earmarkID, earmarkRefID, 33,
						eventItemID, "some note",
						tstTs, tstTs,
					),
			)

		request := &icbt.RemoveEarmarkRequest{
			RefId: earmarkRefID.String(),
		}
		_, err := server.RemoveEarmark(ctx, request)
		assertTwirpError(t, err, twirp.PermissionDenied, "permission denied")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("remove earmark for bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.RemoveEarmarkRequest{
			RefId: "hodor",
		}
		_, err := server.RemoveEarmark(ctx, request)
		assertTwirpError(t, err, twirp.InvalidArgument, "ref_id bad earmark ref-id")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("remove earmark for archived event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		earmarkRefID := refid.Must(model.NewEarmarkRefID())
		eventItemID := 33
		earmarkID := 5
		eventID := 22

		mock.ExpectQuery("SELECT (.+) FROM earmark_").
			WithArgs(earmarkRefID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "user_id",
						"event_item_id", "note",
						"created", "last_modified",
					}).
					AddRow(
						earmarkID, earmarkRefID, user.ID,
						eventItemID, "some note",
						tstTs, tstTs,
					),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(earmarkID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "user_id", "name", "description",
						"archived", "created", "last_modified",
					}).
					AddRow(
						eventID, refid.Must(model.NewEventRefID()), user.ID,
						"event name", "event desc",
						true, tstTs, tstTs,
					),
			)

		request := &icbt.RemoveEarmarkRequest{
			RefId: earmarkRefID.String(),
		}
		_, err := server.RemoveEarmark(ctx, request)
		assertTwirpError(t, err, twirp.PermissionDenied, "event is archived")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
