// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rpc

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/dropwhile/assert"
	"go.uber.org/mock/gomock"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
	"github.com/dropwhile/icanbringthat/internal/util"
	icbt "github.com/dropwhile/icanbringthat/rpc/icbt/rpc/v1"
)

var eventItemFailIfCheck service.FailIfCheckFunc[*model.EventItem]

func TestRpc_ListEventItems(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
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
		eventRefID := util.Must(model.NewEventRefID())
		eventID := 3

		mock.EXPECT().
			GetEventItemsByEvent(ctx, eventRefID).
			Return(
				[]*model.EventItem{{
					ID:          eventID,
					RefID:       util.Must(model.NewEventItemRefID()),
					EventID:     eventID,
					Description: "some desc",
				}}, nil,
			)

		request := icbt.EventListItemsRequest_builder{
			RefId: eventRefID.String(),
		}.Build()
		response, err := server.EventListItems(ctx, connect.NewRequest(request))
		assert.Nil(t, err)
		assert.Equal(t, len(response.Msg.GetItems()), 1)
	})

	t.Run("list event items with bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, _ := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		request := icbt.EventListItemsRequest_builder{
			RefId: "hodor",
		}.Build()
		_, err := server.EventListItems(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeInvalidArgument, "bad event ref-id")
	})

	t.Run("list event items with missing event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := util.Must(model.NewEventRefID())

		mock.EXPECT().
			GetEventItemsByEvent(ctx, eventRefID).
			Return(nil, errs.NotFound.Error("event not found"))

		request := icbt.EventListItemsRequest_builder{
			RefId: eventRefID.String(),
		}.Build()
		_, err := server.EventListItems(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeNotFound, "event not found")
	})
}

func TestRpc_RemoveEventItem(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
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
		eventItemRefID := util.Must(model.NewEventItemRefID())

		mock.EXPECT().
			RemoveEventItem(
				ctx, user.ID, eventItemRefID,
				gomock.AssignableToTypeOf(eventItemFailIfCheck),
			).
			Return(nil)

		request := icbt.EventRemoveItemRequest_builder{
			RefId: eventItemRefID.String(),
		}.Build()
		_, err := server.EventRemoveItem(ctx, connect.NewRequest(request))
		assert.Nil(t, err)
	})

	t.Run("remove event item with bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, _ := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		request := icbt.EventRemoveItemRequest_builder{
			RefId: "hodor",
		}.Build()
		_, err := server.EventRemoveItem(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeInvalidArgument, "bad event-item ref-id")
	})

	t.Run("remove event item not event owner should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := util.Must(model.NewEventItemRefID())

		mock.EXPECT().
			RemoveEventItem(
				ctx, user.ID, eventItemRefID,
				gomock.AssignableToTypeOf(eventItemFailIfCheck),
			).
			Return(errs.PermissionDenied.Error("not event owner"))

		request := icbt.EventRemoveItemRequest_builder{
			RefId: eventItemRefID.String(),
		}.Build()
		_, err := server.EventRemoveItem(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodePermissionDenied, "not event owner")
	})

	t.Run("remove event item with archived event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := util.Must(model.NewEventItemRefID())

		mock.EXPECT().
			RemoveEventItem(
				ctx, user.ID, eventItemRefID,
				gomock.AssignableToTypeOf(eventItemFailIfCheck),
			).
			Return(errs.PermissionDenied.Error("event is archived"))

		request := icbt.EventRemoveItemRequest_builder{
			RefId: eventItemRefID.String(),
		}.Build()
		_, err := server.EventRemoveItem(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodePermissionDenied, "event is archived")
	})

	t.Run("remove event item with event item not found should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := util.Must(model.NewEventItemRefID())

		mock.EXPECT().
			RemoveEventItem(
				ctx, user.ID, eventItemRefID,
				gomock.AssignableToTypeOf(eventItemFailIfCheck),
			).
			Return(errs.NotFound.Error("event-item not found"))

		request := icbt.EventRemoveItemRequest_builder{
			RefId: eventItemRefID.String(),
		}.Build()
		_, err := server.EventRemoveItem(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeNotFound, "event-item not found")
	})
}

func TestRpc_AddEventItem(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
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
		eventRefID := util.Must(model.NewEventRefID())
		eventID := 3
		description := "some description"

		mock.EXPECT().
			AddEventItem(ctx, user.ID, eventRefID, description).
			Return(
				&model.EventItem{
					ID:          eventID,
					RefID:       util.Must(model.NewEventItemRefID()),
					EventID:     eventID,
					Description: description,
				}, nil,
			)

		request := icbt.EventAddItemRequest_builder{
			EventRefId:  eventRefID.String(),
			Description: description,
		}.Build()
		response, err := server.EventAddItem(ctx, connect.NewRequest(request))
		assert.Nil(t, err)
		assert.Equal(t, response.Msg.GetEventItem().GetDescription(), description)
	})

	t.Run("add event item with empty description should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := util.Must(model.NewEventRefID())
		description := ""

		mock.EXPECT().
			AddEventItem(ctx, user.ID, eventRefID, description).
			Return(nil, errs.ArgumentError("description", "bad value"))

		request := icbt.EventAddItemRequest_builder{
			EventRefId:  eventRefID.String(),
			Description: description,
		}.Build()
		_, err := server.EventAddItem(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeInvalidArgument, "description bad value")
	})

	t.Run("add event item with user not event owner should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := util.Must(model.NewEventRefID())
		description := "some description"

		mock.EXPECT().
			AddEventItem(ctx, user.ID, eventRefID, description).
			Return(nil, errs.PermissionDenied.Error("not event owner"))

		request := icbt.EventAddItemRequest_builder{
			EventRefId:  eventRefID.String(),
			Description: description,
		}.Build()
		_, err := server.EventAddItem(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodePermissionDenied, "not event owner")
	})

	t.Run("add event item to archived event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := util.Must(model.NewEventRefID())
		description := "some description"

		mock.EXPECT().
			AddEventItem(ctx, user.ID, eventRefID, description).
			Return(nil, errs.PermissionDenied.Error("event is archived"))

		request := icbt.EventAddItemRequest_builder{
			EventRefId:  eventRefID.String(),
			Description: description,
		}.Build()
		_, err := server.EventAddItem(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodePermissionDenied, "event is archived")
	})

	t.Run("add event item to missing event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := util.Must(model.NewEventRefID())
		description := "some description"

		mock.EXPECT().
			AddEventItem(ctx, user.ID, eventRefID, description).
			Return(nil, errs.NotFound.Error("event not found"))

		request := icbt.EventAddItemRequest_builder{
			EventRefId:  eventRefID.String(),
			Description: description,
		}.Build()
		_, err := server.EventAddItem(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeNotFound, "event not found")
	})

	t.Run("add event item with bad event ref-id should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, _ := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		description := "some description"

		request := icbt.EventAddItemRequest_builder{
			EventRefId:  "hodor",
			Description: description,
		}.Build()
		_, err := server.EventAddItem(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeInvalidArgument, "bad event ref-id")
	})
}

func TestRpc_UpdateEventItem(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
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
		eventItemRefID := util.Must(model.NewEventItemRefID())
		eventItemID := 5
		eventID := 3
		description := "some new description"

		mock.EXPECT().
			UpdateEventItem(
				ctx, user.ID, eventItemRefID, description,
				gomock.AssignableToTypeOf(eventItemFailIfCheck),
			).
			Return(
				&model.EventItem{
					ID:          eventItemID,
					RefID:       eventItemRefID,
					EventID:     eventID,
					Description: description,
				}, nil,
			)

		request := icbt.EventUpdateItemRequest_builder{
			RefId:       eventItemRefID.String(),
			Description: description,
		}.Build()
		response, err := server.EventUpdateItem(ctx, connect.NewRequest(request))
		assert.Nil(t, err)
		assert.Equal(t, response.Msg.GetEventItem().GetDescription(), description)
	})

	t.Run("update event item with bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, _ := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		request := icbt.EventUpdateItemRequest_builder{
			RefId:       "hodor",
			Description: "some nonsense",
		}.Build()
		_, err := server.EventUpdateItem(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeInvalidArgument, "bad event-item ref-id")
	})

	t.Run("update event item with archived event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := util.Must(model.NewEventItemRefID())
		description := "some new description"

		mock.EXPECT().
			UpdateEventItem(
				ctx, user.ID, eventItemRefID, description,
				gomock.AssignableToTypeOf(eventItemFailIfCheck),
			).
			Return(nil, errs.PermissionDenied.Error("event is archived"))

		request := icbt.EventUpdateItemRequest_builder{
			RefId:       eventItemRefID.String(),
			Description: description,
		}.Build()
		_, err := server.EventUpdateItem(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodePermissionDenied, "event is archived")
	})

	t.Run("update event item with user not event owner should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := util.Must(model.NewEventItemRefID())
		description := "some new description"

		mock.EXPECT().
			UpdateEventItem(
				ctx, user.ID, eventItemRefID, description,
				gomock.AssignableToTypeOf(eventItemFailIfCheck),
			).
			Return(nil, errs.PermissionDenied.Error("not event owner"))

		request := icbt.EventUpdateItemRequest_builder{
			RefId:       eventItemRefID.String(),
			Description: description,
		}.Build()
		_, err := server.EventUpdateItem(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodePermissionDenied, "not event owner")
	})

	t.Run("update event item with earmarked by other should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := util.Must(model.NewEventItemRefID())
		description := "some new description"

		mock.EXPECT().
			UpdateEventItem(
				ctx, user.ID, eventItemRefID, description,
				gomock.AssignableToTypeOf(eventItemFailIfCheck),
			).
			Return(nil, errs.PermissionDenied.Error("earmarked by other user"))

		request := icbt.EventUpdateItemRequest_builder{
			RefId:       eventItemRefID.String(),
			Description: description,
		}.Build()
		_, err := server.EventUpdateItem(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodePermissionDenied, "earmarked by other user")
	})

	t.Run("update event item with bad description should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := util.Must(model.NewEventItemRefID())
		description := ""

		mock.EXPECT().
			UpdateEventItem(
				ctx, user.ID, eventItemRefID, description,
				gomock.AssignableToTypeOf(eventItemFailIfCheck),
			).
			Return(nil, errs.ArgumentError("description", "bad value"))

		request := icbt.EventUpdateItemRequest_builder{
			RefId:       eventItemRefID.String(),
			Description: description,
		}.Build()
		_, err := server.EventUpdateItem(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeInvalidArgument, "description bad value")
	})

	t.Run("update event item with event item not found should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := util.Must(model.NewEventItemRefID())
		description := ""

		mock.EXPECT().
			UpdateEventItem(
				ctx, user.ID, eventItemRefID, description,
				gomock.AssignableToTypeOf(eventItemFailIfCheck),
			).
			Return(nil, errs.NotFound.Error("event-item not found"))

		request := icbt.EventUpdateItemRequest_builder{
			RefId:       eventItemRefID.String(),
			Description: description,
		}.Build()
		_, err := server.EventUpdateItem(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeNotFound, "event-item not found")
	})
}
