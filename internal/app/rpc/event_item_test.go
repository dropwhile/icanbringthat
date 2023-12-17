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

func TestServer_ListEventItems(t *testing.T) {
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

	t.Run("list event items should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventID := 3

		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(eventRefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived", "created", "last_modified",
				}).
				AddRow(
					eventID, eventRefID, user.ID,
					"event name", "event desc",
					false, tstTs, tstTs,
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
					12, refid.Must(model.NewEventItemRefID()),
					10, "some description",
					tstTs, tstTs,
				),
			)

		request := &icbt.ListEventItemsRequest{
			RefId: eventRefID.String(),
		}
		response, err := server.ListEventItems(ctx, request)
		assert.NilError(t, err)

		assert.Check(t, len(response.Items) == 1)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("list event items with bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.ListEventItemsRequest{
			RefId: "hodor",
		}
		_, err := server.ListEventItems(ctx, request)
		assertTwirpError(t, err, twirp.InvalidArgument, "ref_id bad event ref-id")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestServer_RemoveEventItem(t *testing.T) {
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

	t.Run("remove event item should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		eventItemID := 11
		eventID := 3

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(eventItemRefID).
			WillReturnRows(pgxmock.NewRows(
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
		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(eventID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived", "created", "last_modified",
				}).
				AddRow(
					eventID, eventRefID, user.ID,
					"event name", "event desc",
					false, tstTs, tstTs,
				),
			)
		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM event_item_").
			WithArgs(eventItemID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		request := &icbt.RemoveEventItemRequest{
			RefId: eventItemRefID.String(),
		}
		_, err := server.RemoveEventItem(ctx, request)
		assert.NilError(t, err)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("remove event item with bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.RemoveEventItemRequest{
			RefId: "hodor",
		}
		_, err := server.RemoveEventItem(ctx, request)
		assertTwirpError(t, err, twirp.InvalidArgument, "ref_id bad event-item ref-id")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("remove event item not event owner should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		eventItemID := 11
		eventID := 3

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(eventItemRefID).
			WillReturnRows(pgxmock.NewRows(
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
		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(eventID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived", "created", "last_modified",
				}).
				AddRow(
					eventID, eventRefID, 99,
					"event name", "event desc",
					false, tstTs, tstTs,
				),
			)

		request := &icbt.RemoveEventItemRequest{
			RefId: eventItemRefID.String(),
		}
		_, err := server.RemoveEventItem(ctx, request)
		assertTwirpError(t, err, twirp.PermissionDenied, "not event owner")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("remove event item with archived event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		eventItemID := 11
		eventID := 3

		mock.ExpectQuery("SELECT (.+) FROM event_item_ ").
			WithArgs(eventItemRefID).
			WillReturnRows(pgxmock.NewRows(
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
		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(eventID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived", "created", "last_modified",
				}).
				AddRow(
					eventID, eventRefID, user.ID,
					"event name", "event desc",
					true, tstTs, tstTs,
				),
			)

		request := &icbt.RemoveEventItemRequest{
			RefId: eventItemRefID.String(),
		}
		_, err := server.RemoveEventItem(ctx, request)
		assertTwirpError(t, err, twirp.PermissionDenied, "event is archived")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestServer_AddEventItem(t *testing.T) {
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

	t.Run("add event item should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventID := 3
		description := "some description"

		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(eventRefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived", "created", "last_modified",
				}).
				AddRow(
					eventID, eventRefID, user.ID,
					"event name", "event desc",
					false, tstTs, tstTs,
				),
			)
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO event_item_").
			WithArgs(
				pgx.NamedArgs{
					"refID":       model.EventItemRefIDMatcher,
					"eventID":     eventID,
					"description": description,
				}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"event_id", "description",
					"created", "last_modified",
				}).
				AddRow(
					99, refid.Must(model.NewEventItemRefID()),
					eventID, description,
					tstTs, tstTs,
				),
			)
		mock.ExpectCommit()
		mock.ExpectRollback()

		request := &icbt.AddEventItemRequest{
			EventRefId:  eventRefID.String(),
			Description: description,
		}
		response, err := server.AddEventItem(ctx, request)
		assert.NilError(t, err)
		assert.Equal(t, response.EventItem.Description, description)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("add event item with user not event owner should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventID := 3
		description := "some description"

		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(eventRefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived", "created", "last_modified",
				}).
				AddRow(
					eventID, eventRefID, 99,
					"event name", "event desc",
					false, tstTs, tstTs,
				),
			)

		request := &icbt.AddEventItemRequest{
			EventRefId:  eventRefID.String(),
			Description: description,
		}
		_, err := server.AddEventItem(ctx, request)
		assertTwirpError(t, err, twirp.PermissionDenied, "not event owner")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("add event item to archived event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventID := 3
		description := "some description"

		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(eventRefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived", "created", "last_modified",
				}).
				AddRow(
					eventID, eventRefID, user.ID,
					"event name", "event desc",
					true, tstTs, tstTs,
				),
			)

		request := &icbt.AddEventItemRequest{
			EventRefId:  eventRefID.String(),
			Description: description,
		}
		_, err := server.AddEventItem(ctx, request)
		assertTwirpError(t, err, twirp.PermissionDenied, "event is archived")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("add event item with bad event ref-id should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		description := "some description"

		request := &icbt.AddEventItemRequest{
			EventRefId:  "hodor",
			Description: description,
		}
		_, err := server.AddEventItem(ctx, request)
		assertTwirpError(t, err, twirp.InvalidArgument, "ref_id bad event ref-id")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestServer_UpdateEventItem(t *testing.T) {
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

	t.Run("update event item should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		eventItemID := 5
		eventID := 3
		newDescription := "some new description"

		mock.ExpectQuery("SELECT (.+) FROM event_item_").
			WithArgs(eventItemRefID).
			WillReturnRows(pgxmock.NewRows(
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
		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(eventID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived", "created", "last_modified",
				}).
				AddRow(
					eventID, eventRefID, user.ID,
					"event name", "event desc",
					false, tstTs, tstTs,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM earmark_").
			WithArgs(eventItemID).
			WillReturnError(pgx.ErrNoRows)
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE event_item_").
			WithArgs(
				pgx.NamedArgs{
					"eventItemID": eventItemID,
					"description": newDescription,
				}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		request := &icbt.UpdateEventItemRequest{
			RefId:       eventItemRefID.String(),
			Description: newDescription,
		}
		response, err := server.UpdateEventItem(ctx, request)
		assert.NilError(t, err)
		assert.Equal(t, response.EventItem.Description, newDescription)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update event item with bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		request := &icbt.UpdateEventItemRequest{
			RefId:       "hodor",
			Description: "some nonsense",
		}
		_, err := server.UpdateEventItem(ctx, request)
		assertTwirpError(t, err, twirp.InvalidArgument, "ref_id bad event-item ref-id")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update event item with archived event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		eventItemID := 5
		eventID := 3
		newDescription := "some new description"

		mock.ExpectQuery("SELECT (.+) FROM event_item_").
			WithArgs(eventItemRefID).
			WillReturnRows(pgxmock.NewRows(
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
		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(eventID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived", "created", "last_modified",
				}).
				AddRow(
					eventID, eventRefID, user.ID,
					"event name", "event desc",
					true, tstTs, tstTs,
				),
			)

		request := &icbt.UpdateEventItemRequest{
			RefId:       eventItemRefID.String(),
			Description: newDescription,
		}
		_, err := server.UpdateEventItem(ctx, request)
		assertTwirpError(t, err, twirp.PermissionDenied, "event is archived")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update event item with user not event owner should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		eventItemID := 5
		eventID := 3
		newDescription := "some new description"

		mock.ExpectQuery("SELECT (.+) FROM event_item_").
			WithArgs(eventItemRefID).
			WillReturnRows(pgxmock.NewRows(
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
		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(eventID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived", "created", "last_modified",
				}).
				AddRow(
					eventID, eventRefID, 33,
					"event name", "event desc",
					false, tstTs, tstTs,
				),
			)

		request := &icbt.UpdateEventItemRequest{
			RefId:       eventItemRefID.String(),
			Description: newDescription,
		}
		_, err := server.UpdateEventItem(ctx, request)
		assertTwirpError(t, err, twirp.PermissionDenied, "not event owner")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update event item with earmarked by self should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		eventItemID := 5
		eventID := 3
		newDescription := "some new description"

		mock.ExpectQuery("SELECT (.+) FROM event_item_").
			WithArgs(eventItemRefID).
			WillReturnRows(pgxmock.NewRows(
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
		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(eventID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived", "created", "last_modified",
				}).
				AddRow(
					eventID, eventRefID, user.ID,
					"event name", "event desc",
					false, tstTs, tstTs,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM earmark_").
			WithArgs(eventItemID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id",
					"event_item_id", "note",
					"created", "last_modified",
				}).
				AddRow(
					99, refid.Must(model.NewEarmarkRefID()), user.ID,
					eventItemID, "some note",
					tstTs, tstTs,
				),
			)
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE event_item_").
			WithArgs(
				pgx.NamedArgs{
					"eventItemID": eventItemID,
					"description": newDescription,
				}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		request := &icbt.UpdateEventItemRequest{
			RefId:       eventItemRefID.String(),
			Description: newDescription,
		}
		response, err := server.UpdateEventItem(ctx, request)
		assert.NilError(t, err)
		assert.Equal(t, response.EventItem.Description, newDescription)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update event item with earmarked by other should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		eventItemID := 5
		eventID := 3
		newDescription := "some new description"

		mock.ExpectQuery("SELECT (.+) FROM event_item_").
			WithArgs(eventItemRefID).
			WillReturnRows(pgxmock.NewRows(
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
		mock.ExpectQuery("SELECT (.+) FROM event_ ").
			WithArgs(eventID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived", "created", "last_modified",
				}).
				AddRow(
					eventID, eventRefID, user.ID,
					"event name", "event desc",
					false, tstTs, tstTs,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM earmark_").
			WithArgs(eventItemID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id",
					"event_item_id", "note",
					"created", "last_modified",
				}).
				AddRow(
					99, refid.Must(model.NewEarmarkRefID()), 99,
					eventItemID, "some note",
					tstTs, tstTs,
				),
			)

		request := &icbt.UpdateEventItemRequest{
			RefId:       eventItemRefID.String(),
			Description: newDescription,
		}
		_, err := server.UpdateEventItem(ctx, request)
		assertTwirpError(t, err, twirp.PermissionDenied, "earmarked by other user")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
