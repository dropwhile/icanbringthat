package rpc

import (
	"context"
	"testing"

	"github.com/dropwhile/refid/v2"
	mox "github.com/stretchr/testify/mock"
	"github.com/twitchtv/twirp"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/middleware/auth"
	"github.com/dropwhile/icbt/rpc/icbt"
)

func TestRpc_ListEventItems(t *testing.T) {
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
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventID := 3

		mock.EXPECT().
			GetEventItemsByEvent(ctx, eventRefID).
			Return(
				[]*model.EventItem{{
					ID:          eventID,
					RefID:       refid.Must(model.NewEventItemRefID()),
					EventID:     eventID,
					Description: "some desc",
				}}, nil,
			).
			Once()

		request := &icbt.ListEventItemsRequest{
			RefId: eventRefID.String(),
		}
		response, err := server.ListEventItems(ctx, request)
		assert.NilError(t, err)
		assert.Equal(t, len(response.Items), 1)
		mock.AssertExpectations(t)
	})

	t.Run("list event items with bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.ListEventItemsRequest{
			RefId: "hodor",
		}
		_, err := server.ListEventItems(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "ref_id bad event ref-id")
		mock.AssertExpectations(t)
	})

	t.Run("list event items with missing event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())

		mock.EXPECT().
			GetEventItemsByEvent(ctx, eventRefID).
			Return(nil, errs.NotFound.Error("event not found")).
			Once()

		request := &icbt.ListEventItemsRequest{
			RefId: eventRefID.String(),
		}
		_, err := server.ListEventItems(ctx, request)
		errs.AssertError(t, err, twirp.NotFound, "event not found")
		mock.AssertExpectations(t)
	})
}

func TestRpc_RemoveEventItem(t *testing.T) {
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
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := refid.Must(model.NewEventItemRefID())

		mock.EXPECT().
			RemoveEventItem(
				ctx, user.ID, eventItemRefID,
				mox.AnythingOfType("func(*model.EventItem) bool"),
			).
			Return(nil)

		request := &icbt.RemoveEventItemRequest{
			RefId: eventItemRefID.String(),
		}
		_, err := server.RemoveEventItem(ctx, request)
		assert.NilError(t, err)
		mock.AssertExpectations(t)
	})

	t.Run("remove event item with bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.RemoveEventItemRequest{
			RefId: "hodor",
		}
		_, err := server.RemoveEventItem(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "ref_id bad event-item ref-id")
		mock.AssertExpectations(t)
	})

	t.Run("remove event item not event owner should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := refid.Must(model.NewEventItemRefID())

		mock.EXPECT().
			RemoveEventItem(
				ctx, user.ID, eventItemRefID,
				mox.AnythingOfType("func(*model.EventItem) bool"),
			).
			Return(errs.PermissionDenied.Error("not event owner"))

		request := &icbt.RemoveEventItemRequest{
			RefId: eventItemRefID.String(),
		}
		_, err := server.RemoveEventItem(ctx, request)
		errs.AssertError(t, err, twirp.PermissionDenied, "not event owner")
		mock.AssertExpectations(t)
	})

	t.Run("remove event item with archived event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := refid.Must(model.NewEventItemRefID())

		mock.EXPECT().
			RemoveEventItem(
				ctx, user.ID, eventItemRefID,
				mox.AnythingOfType("func(*model.EventItem) bool"),
			).
			Return(errs.PermissionDenied.Error("event is archived"))

		request := &icbt.RemoveEventItemRequest{
			RefId: eventItemRefID.String(),
		}
		_, err := server.RemoveEventItem(ctx, request)
		errs.AssertError(t, err, twirp.PermissionDenied, "event is archived")
		mock.AssertExpectations(t)
	})

	t.Run("remove event item with event item not found should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := refid.Must(model.NewEventItemRefID())

		mock.EXPECT().
			RemoveEventItem(
				ctx, user.ID, eventItemRefID,
				mox.AnythingOfType("func(*model.EventItem) bool"),
			).
			Return(errs.NotFound.Error("event-item not found"))

		request := &icbt.RemoveEventItemRequest{
			RefId: eventItemRefID.String(),
		}
		_, err := server.RemoveEventItem(ctx, request)
		errs.AssertError(t, err, twirp.NotFound, "event-item not found")
		mock.AssertExpectations(t)
	})
}

func TestRpc_AddEventItem(t *testing.T) {
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
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventID := 3
		description := "some description"

		mock.EXPECT().
			AddEventItem(ctx, user.ID, eventRefID, description).
			Return(
				&model.EventItem{
					ID:          eventID,
					RefID:       refid.Must(model.NewEventItemRefID()),
					EventID:     eventID,
					Description: description,
				}, nil,
			).
			Once()

		request := &icbt.AddEventItemRequest{
			EventRefId:  eventRefID.String(),
			Description: description,
		}
		response, err := server.AddEventItem(ctx, request)
		assert.NilError(t, err)
		assert.Equal(t, response.EventItem.Description, description)
		mock.AssertExpectations(t)
	})

	t.Run("add event item with empty description should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		description := ""

		mock.EXPECT().
			AddEventItem(ctx, user.ID, eventRefID, description).
			Return(nil, errs.InvalidArgumentError("description", "bad value")).
			Once()

		request := &icbt.AddEventItemRequest{
			EventRefId:  eventRefID.String(),
			Description: description,
		}
		_, err := server.AddEventItem(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "description bad value")
		mock.AssertExpectations(t)
	})

	t.Run("add event item with user not event owner should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		description := "some description"

		mock.EXPECT().
			AddEventItem(ctx, user.ID, eventRefID, description).
			Return(nil, errs.PermissionDenied.Error("not event owner")).
			Once()

		request := &icbt.AddEventItemRequest{
			EventRefId:  eventRefID.String(),
			Description: description,
		}
		_, err := server.AddEventItem(ctx, request)
		errs.AssertError(t, err, twirp.PermissionDenied, "not event owner")
		mock.AssertExpectations(t)
	})

	t.Run("add event item to archived event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		description := "some description"

		mock.EXPECT().
			AddEventItem(ctx, user.ID, eventRefID, description).
			Return(nil, errs.PermissionDenied.Error("event is archived")).
			Once()

		request := &icbt.AddEventItemRequest{
			EventRefId:  eventRefID.String(),
			Description: description,
		}
		_, err := server.AddEventItem(ctx, request)
		errs.AssertError(t, err, twirp.PermissionDenied, "event is archived")
		mock.AssertExpectations(t)
	})

	t.Run("add event item to missing event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		description := "some description"

		mock.EXPECT().
			AddEventItem(ctx, user.ID, eventRefID, description).
			Return(nil, errs.NotFound.Error("event not found")).
			Once()

		request := &icbt.AddEventItemRequest{
			EventRefId:  eventRefID.String(),
			Description: description,
		}
		_, err := server.AddEventItem(ctx, request)
		errs.AssertError(t, err, twirp.NotFound, "event not found")
		mock.AssertExpectations(t)
	})

	t.Run("add event item with bad event ref-id should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		description := "some description"

		request := &icbt.AddEventItemRequest{
			EventRefId:  "hodor",
			Description: description,
		}
		_, err := server.AddEventItem(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "ref_id bad event ref-id")
		mock.AssertExpectations(t)
	})
}

func TestRpc_UpdateEventItem(t *testing.T) {
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
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		eventItemID := 5
		eventID := 3
		description := "some new description"

		mock.EXPECT().
			UpdateEventItem(
				ctx, user.ID, eventItemRefID, description,
				mox.AnythingOfType("func(*model.EventItem) bool"),
			).
			Return(
				&model.EventItem{
					ID:          eventItemID,
					RefID:       eventItemRefID,
					EventID:     eventID,
					Description: description,
				}, nil,
			).
			Once()

		request := &icbt.UpdateEventItemRequest{
			RefId:       eventItemRefID.String(),
			Description: description,
		}
		response, err := server.UpdateEventItem(ctx, request)
		assert.NilError(t, err)
		assert.Equal(t, response.EventItem.Description, description)
		mock.AssertExpectations(t)
	})

	t.Run("update event item with bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		request := &icbt.UpdateEventItemRequest{
			RefId:       "hodor",
			Description: "some nonsense",
		}
		_, err := server.UpdateEventItem(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "ref_id bad event-item ref-id")
		mock.AssertExpectations(t)
	})

	t.Run("update event item with archived event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		description := "some new description"

		mock.EXPECT().
			UpdateEventItem(
				ctx, user.ID, eventItemRefID, description,
				mox.AnythingOfType("func(*model.EventItem) bool"),
			).
			Return(nil, errs.PermissionDenied.Error("event is archived")).
			Once()

		request := &icbt.UpdateEventItemRequest{
			RefId:       eventItemRefID.String(),
			Description: description,
		}
		_, err := server.UpdateEventItem(ctx, request)
		errs.AssertError(t, err, twirp.PermissionDenied, "event is archived")
		mock.AssertExpectations(t)
	})

	t.Run("update event item with user not event owner should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		description := "some new description"

		mock.EXPECT().
			UpdateEventItem(
				ctx, user.ID, eventItemRefID, description,
				mox.AnythingOfType("func(*model.EventItem) bool"),
			).
			Return(nil, errs.PermissionDenied.Error("not event owner")).
			Once()

		request := &icbt.UpdateEventItemRequest{
			RefId:       eventItemRefID.String(),
			Description: description,
		}
		_, err := server.UpdateEventItem(ctx, request)
		errs.AssertError(t, err, twirp.PermissionDenied, "not event owner")
		mock.AssertExpectations(t)
	})

	t.Run("update event item with earmarked by other should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		description := "some new description"

		mock.EXPECT().
			UpdateEventItem(
				ctx, user.ID, eventItemRefID, description,
				mox.AnythingOfType("func(*model.EventItem) bool"),
			).
			Return(nil, errs.PermissionDenied.Error("earmarked by other user")).
			Once()

		request := &icbt.UpdateEventItemRequest{
			RefId:       eventItemRefID.String(),
			Description: description,
		}
		_, err := server.UpdateEventItem(ctx, request)
		errs.AssertError(t, err, twirp.PermissionDenied, "earmarked by other user")
		mock.AssertExpectations(t)
	})

	t.Run("update event item with bad description should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		description := ""

		mock.EXPECT().
			UpdateEventItem(
				ctx, user.ID, eventItemRefID, description,
				mox.AnythingOfType("func(*model.EventItem) bool"),
			).
			Return(nil, errs.InvalidArgumentError("description", "bad value")).
			Once()

		request := &icbt.UpdateEventItemRequest{
			RefId:       eventItemRefID.String(),
			Description: description,
		}
		_, err := server.UpdateEventItem(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "description bad value")
		mock.AssertExpectations(t)
	})

	t.Run("update event item with event item not found should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		description := ""

		mock.EXPECT().
			UpdateEventItem(
				ctx, user.ID, eventItemRefID, description,
				mox.AnythingOfType("func(*model.EventItem) bool"),
			).
			Return(nil, errs.NotFound.Error("event-item not found")).
			Once()

		request := &icbt.UpdateEventItemRequest{
			RefId:       eventItemRefID.String(),
			Description: description,
		}
		_, err := server.UpdateEventItem(ctx, request)
		errs.AssertError(t, err, twirp.NotFound, "event-item not found")
		mock.AssertExpectations(t)
	})
}
