package rpc

import (
	"context"
	"testing"

	"github.com/dropwhile/refid/v2"
	"github.com/twitchtv/twirp"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/middleware/auth"
	"github.com/dropwhile/icbt/internal/util"
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

		ctx := context.TODO()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())

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
			).
			Once()

		request := &icbt.ListFavoriteEventsRequest{
			Pagination: &icbt.PaginationRequest{Limit: 10, Offset: 0},
			Archived:   func(b bool) *bool { return &b }(false),
		}
		response, err := server.ListFavoriteEvents(ctx, request)
		assert.NilError(t, err)

		assert.Equal(t, len(response.Events), 1)
		mock.AssertExpectations(t)
	})

	t.Run("list favorite events non-paginated should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())

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
			).
			Once()

		request := &icbt.ListFavoriteEventsRequest{
			Archived: func(b bool) *bool { return &b }(false),
		}
		response, err := server.ListFavoriteEvents(ctx, request)
		assert.NilError(t, err)
		assert.Equal(t, len(response.Events), 1)
		mock.AssertExpectations(t)
	})
}

func TestRpc_AddFavorite(t *testing.T) {
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

	t.Run("add favorite should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())

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
			}, nil).
			Once()

		request := &icbt.CreateFavoriteRequest{
			EventRefId: eventRefID.String(),
		}
		response, err := server.AddFavorite(ctx, request)
		assert.NilError(t, err)
		assert.Equal(t, response.Favorite.EventRefId, eventRefID.String())
		mock.AssertExpectations(t)
	})

	t.Run("add favorite for own event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())

		mock.EXPECT().
			AddFavorite(ctx, user.ID, eventRefID).
			Return(nil, errs.PermissionDenied.Error("can't favorite own event")).
			Once()

		request := &icbt.CreateFavoriteRequest{
			EventRefId: eventRefID.String(),
		}
		_, err := server.AddFavorite(ctx, request)
		errs.AssertError(t, err, twirp.PermissionDenied, "can't favorite own event")
		mock.AssertExpectations(t)
	})

	t.Run("add favorite for already favorited should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())

		mock.EXPECT().
			AddFavorite(ctx, user.ID, eventRefID).
			Return(nil, errs.AlreadyExists.Error("favorite already exists")).
			Once()

		request := &icbt.CreateFavoriteRequest{
			EventRefId: eventRefID.String(),
		}
		_, err := server.AddFavorite(ctx, request)
		errs.AssertError(t, err, twirp.AlreadyExists, "favorite already exists")
		mock.AssertExpectations(t)
	})

	t.Run("add favorite with bad event refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.CreateFavoriteRequest{
			EventRefId: "hodor",
		}
		_, err := server.AddFavorite(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "ref_id incorrect value type")
		mock.AssertExpectations(t)
	})
}

func TestRpc_RemoveFavorite(t *testing.T) {
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

	t.Run("remove favorite should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())

		mock.EXPECT().
			RemoveFavorite(ctx, user.ID, eventRefID).
			Return(nil)

		request := &icbt.RemoveFavoriteRequest{
			EventRefId: eventRefID.String(),
		}
		_, err := server.RemoveFavorite(ctx, request)
		assert.NilError(t, err)
		mock.AssertExpectations(t)
	})

	t.Run("remove favorite with bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.RemoveFavoriteRequest{
			EventRefId: "hodor",
		}
		_, err := server.RemoveFavorite(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "ref_id incorrect value type")
		mock.AssertExpectations(t)
	})

	t.Run("remove favorite with event not found should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())

		mock.EXPECT().
			RemoveFavorite(ctx, user.ID, eventRefID).
			Return(errs.NotFound.Error("event not found"))

		request := &icbt.RemoveFavoriteRequest{
			EventRefId: eventRefID.String(),
		}
		_, err := server.RemoveFavorite(ctx, request)
		errs.AssertError(t, err, twirp.NotFound, "event not found")
		mock.AssertExpectations(t)
	})

	t.Run("remove favorite with favorite not found should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())

		mock.EXPECT().
			RemoveFavorite(ctx, user.ID, eventRefID).
			Return(errs.NotFound.Error("favorite not found"))

		request := &icbt.RemoveFavoriteRequest{
			EventRefId: eventRefID.String(),
		}
		_, err := server.RemoveFavorite(ctx, request)
		errs.AssertError(t, err, twirp.NotFound, "favorite not found")
		mock.AssertExpectations(t)
	})

	t.Run("remove favorite not owned should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())

		mock.EXPECT().
			RemoveFavorite(ctx, user.ID, eventRefID).
			Return(errs.PermissionDenied.Error("permission denied"))

		request := &icbt.RemoveFavoriteRequest{
			EventRefId: eventRefID.String(),
		}
		_, err := server.RemoveFavorite(ctx, request)
		errs.AssertError(t, err, twirp.PermissionDenied, "permission denied")
		mock.AssertExpectations(t)
	})
}
