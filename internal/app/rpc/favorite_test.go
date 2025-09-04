// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rpc

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/dropwhile/assert"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
	"github.com/dropwhile/icanbringthat/internal/util"
	icbt "github.com/dropwhile/icanbringthat/rpc/icbt/rpc/v1"
)

func TestRpc_ListFavoriteEvents(t *testing.T) {
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

	t.Run("list favorite events paginated should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := util.Must(model.NewEventRefID())

		limit := 10
		offset := 0
		archived := false

		mock.EXPECT().
			GetFavoriteEventsPaginated(ctx, user.ID, limit, offset, archived).
			Return(
				[]*model.Event{{
					Created:       tstTs,
					LastModified:  tstTs,
					StartTime:     tstTs,
					StartTimeTz:   util.Must(service.ParseTimeZone("UTC")),
					Name:          "name",
					Description:   "desc",
					ItemSortOrder: []int{},
					Archived:      archived,
					UserID:        user.ID,
					ID:            1,
					RefID:         eventRefID,
				}},
				&service.Pagination{
					Limit:  limit,
					Offset: offset,
					Count:  1,
				}, nil,
			)

		request := icbt.FavoriteListEventsRequest_builder{
			Pagination: icbt.PaginationRequest_builder{
				Limit:  10,
				Offset: 0,
			}.Build(),
			Archived: func(b bool) *bool { return &b }(false),
		}.Build()
		response, err := server.FavoriteListEvents(ctx, connect.NewRequest(request))
		assert.Nil(t, err)

		assert.Equal(t, len(response.Msg.GetEvents()), 1)
	})

	t.Run("list favorite events non-paginated should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := util.Must(model.NewEventRefID())

		archived := false

		mock.EXPECT().
			GetFavoriteEvents(ctx, user.ID, archived).
			Return(
				[]*model.Event{{
					Created:       tstTs,
					LastModified:  tstTs,
					StartTime:     tstTs,
					StartTimeTz:   util.Must(service.ParseTimeZone("UTC")),
					Name:          "name",
					Description:   "desc",
					ItemSortOrder: []int{},
					Archived:      archived,
					UserID:        user.ID,
					ID:            1,
					RefID:         eventRefID,
				}}, nil,
			)

		request := icbt.FavoriteListEventsRequest_builder{
			Archived: func(b bool) *bool { return &b }(false),
		}.Build()
		response, err := server.FavoriteListEvents(ctx, connect.NewRequest(request))
		assert.Nil(t, err)
		assert.Equal(t, len(response.Msg.GetEvents()), 1)
	})
}

func TestRpc_AddFavorite(t *testing.T) {
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

	t.Run("add favorite should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := util.Must(model.NewEventRefID())

		mock.EXPECT().
			AddFavorite(ctx, user.ID, eventRefID).
			Return(&model.Event{
				ID:            1,
				RefID:         eventRefID,
				UserID:        user.ID,
				ItemSortOrder: []int{},
				StartTime:     tstTs,
				StartTimeTz:   util.Must(service.ParseTimeZone("UTC")),
				Archived:      false,
			}, nil)

		request := icbt.FavoriteAddRequest_builder{
			EventRefId: eventRefID.String(),
		}.Build()
		response, err := server.FavoriteAdd(ctx, connect.NewRequest(request))
		assert.Nil(t, err)
		assert.Equal(t, response.Msg.GetFavorite().GetEventRefId(), eventRefID.String())
	})

	t.Run("add favorite for own event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := util.Must(model.NewEventRefID())

		mock.EXPECT().
			AddFavorite(ctx, user.ID, eventRefID).
			Return(nil, errs.PermissionDenied.Error("can't favorite own event"))

		request := icbt.FavoriteAddRequest_builder{
			EventRefId: eventRefID.String(),
		}.Build()
		_, err := server.FavoriteAdd(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodePermissionDenied, "can't favorite own event")
	})

	t.Run("add favorite for already favorited should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := util.Must(model.NewEventRefID())

		mock.EXPECT().
			AddFavorite(ctx, user.ID, eventRefID).
			Return(nil, errs.AlreadyExists.Error("favorite already exists"))

		request := icbt.FavoriteAddRequest_builder{
			EventRefId: eventRefID.String(),
		}.Build()
		_, err := server.FavoriteAdd(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeAlreadyExists, "favorite already exists")
	})

	t.Run("add favorite with bad event refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, _ := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		request := icbt.FavoriteAddRequest_builder{
			EventRefId: "hodor",
		}.Build()
		_, err := server.FavoriteAdd(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeInvalidArgument, "bad event ref-id")
	})
}

func TestRpc_RemoveFavorite(t *testing.T) {
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

	t.Run("remove favorite should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := util.Must(model.NewEventRefID())

		mock.EXPECT().
			RemoveFavorite(ctx, user.ID, eventRefID).
			Return(nil)

		request := icbt.FavoriteRemoveRequest_builder{
			EventRefId: eventRefID.String(),
		}.Build()
		_, err := server.FavoriteRemove(ctx, connect.NewRequest(request))
		assert.Nil(t, err)
	})

	t.Run("remove favorite with bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, _ := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		request := icbt.FavoriteRemoveRequest_builder{
			EventRefId: "hodor",
		}.Build()
		_, err := server.FavoriteRemove(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeInvalidArgument, "bad event ref-id")
	})

	t.Run("remove favorite with event not found should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := util.Must(model.NewEventRefID())

		mock.EXPECT().
			RemoveFavorite(ctx, user.ID, eventRefID).
			Return(errs.NotFound.Error("event not found"))

		request := icbt.FavoriteRemoveRequest_builder{
			EventRefId: eventRefID.String(),
		}.Build()
		_, err := server.FavoriteRemove(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeNotFound, "event not found")
	})

	t.Run("remove favorite with favorite not found should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := util.Must(model.NewEventRefID())

		mock.EXPECT().
			RemoveFavorite(ctx, user.ID, eventRefID).
			Return(errs.NotFound.Error("favorite not found"))

		request := icbt.FavoriteRemoveRequest_builder{
			EventRefId: eventRefID.String(),
		}.Build()
		_, err := server.FavoriteRemove(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeNotFound, "favorite not found")
	})

	t.Run("remove favorite not owned should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := util.Must(model.NewEventRefID())

		mock.EXPECT().
			RemoveFavorite(ctx, user.ID, eventRefID).
			Return(errs.PermissionDenied.Error("permission denied"))

		request := icbt.FavoriteRemoveRequest_builder{
			EventRefId: eventRefID.String(),
		}.Build()
		_, err := server.FavoriteRemove(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodePermissionDenied, "permission denied")
	})
}
