// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rpc

import (
	"context"
	"testing"

	"github.com/twitchtv/twirp"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
	"github.com/dropwhile/icanbringthat/internal/util"
	"github.com/dropwhile/icanbringthat/rpc/icbt"
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
					Limit:  uint32(limit),
					Offset: uint32(offset),
					Count:  1,
				}, nil,
			)

		request := &icbt.FavoriteListEventsRequest{
			Pagination: &icbt.PaginationRequest{Limit: 10, Offset: 0},
			Archived:   func(b bool) *bool { return &b }(false),
		}
		response, err := server.FavoriteListEvents(ctx, request)
		assert.NilError(t, err)

		assert.Equal(t, len(response.Events), 1)
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

		request := &icbt.FavoriteListEventsRequest{
			Archived: func(b bool) *bool { return &b }(false),
		}
		response, err := server.FavoriteListEvents(ctx, request)
		assert.NilError(t, err)
		assert.Equal(t, len(response.Events), 1)
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

		request := &icbt.FavoriteCreateRequest{
			EventRefId: eventRefID.String(),
		}
		response, err := server.FavoriteAdd(ctx, request)
		assert.NilError(t, err)
		assert.Equal(t, response.Favorite.EventRefId, eventRefID.String())
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

		request := &icbt.FavoriteCreateRequest{
			EventRefId: eventRefID.String(),
		}
		_, err := server.FavoriteAdd(ctx, request)
		errs.AssertError(t, err, twirp.PermissionDenied, "can't favorite own event")
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

		request := &icbt.FavoriteCreateRequest{
			EventRefId: eventRefID.String(),
		}
		_, err := server.FavoriteAdd(ctx, request)
		errs.AssertError(t, err, twirp.AlreadyExists, "favorite already exists")
	})

	t.Run("add favorite with bad event refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, _ := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.FavoriteCreateRequest{
			EventRefId: "hodor",
		}
		_, err := server.FavoriteAdd(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "ref_id incorrect value type")
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

		request := &icbt.FavoriteRemoveRequest{
			EventRefId: eventRefID.String(),
		}
		_, err := server.FavoriteRemove(ctx, request)
		assert.NilError(t, err)
	})

	t.Run("remove favorite with bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, _ := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.FavoriteRemoveRequest{
			EventRefId: "hodor",
		}
		_, err := server.FavoriteRemove(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "ref_id incorrect value type")
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

		request := &icbt.FavoriteRemoveRequest{
			EventRefId: eventRefID.String(),
		}
		_, err := server.FavoriteRemove(ctx, request)
		errs.AssertError(t, err, twirp.NotFound, "event not found")
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

		request := &icbt.FavoriteRemoveRequest{
			EventRefId: eventRefID.String(),
		}
		_, err := server.FavoriteRemove(ctx, request)
		errs.AssertError(t, err, twirp.NotFound, "favorite not found")
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

		request := &icbt.FavoriteRemoveRequest{
			EventRefId: eventRefID.String(),
		}
		_, err := server.FavoriteRemove(ctx, request)
		errs.AssertError(t, err, twirp.PermissionDenied, "permission denied")
	})
}
