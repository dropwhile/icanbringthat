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

func TestRpc_ListEarmarks(t *testing.T) {
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

	t.Run("list earmarks paginated should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		earmarkRefID := util.Must(model.NewEarmarkRefID())

		limit := 10
		offset := 0
		archived := false
		earmarkID := 13
		eventID := 3
		eventItemID := 4

		mock.EXPECT().
			GetEarmarksPaginated(ctx, user.ID, limit, offset, archived).
			Return(
				[]*model.Earmark{{
					ID:           earmarkID,
					RefID:        earmarkRefID,
					EventItemID:  eventItemID,
					UserID:       user.ID,
					Note:         "some note",
					Created:      tstTs,
					LastModified: tstTs,
				}},
				&service.Pagination{
					Limit:  limit,
					Offset: offset,
					Count:  1,
				},
				nil,
			)
		mock.EXPECT().
			GetEventItemByID(ctx, eventItemID).
			Return(
				&model.EventItem{
					ID:           eventItemID,
					RefID:        util.Must(model.NewEventItemRefID()),
					EventID:      eventID,
					Description:  "some desc",
					Created:      tstTs,
					LastModified: tstTs,
				}, nil,
			)
		mock.EXPECT().
			GetUserByID(ctx, user.ID).
			Return(user, nil)

		request := icbt.EarmarksListRequest_builder{
			Pagination: icbt.PaginationRequest_builder{Limit: 10, Offset: 0}.Build(),
			Archived:   func(b bool) *bool { return &b }(false),
		}.Build()
		response, err := server.EarmarksList(ctx, connect.NewRequest(request))
		assert.Nil(t, err)

		assert.Equal(t, len(response.Msg.GetEarmarks()), 1)
	})

	t.Run("list earmarks non-paginated should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		earmarkRefID := util.Must(model.NewEarmarkRefID())

		archived := false
		earmarkID := 13
		eventID := 3
		eventItemID := 4

		mock.EXPECT().
			GetEarmarks(ctx, user.ID, archived).
			Return(
				[]*model.Earmark{{
					ID:           earmarkID,
					RefID:        earmarkRefID,
					EventItemID:  eventItemID,
					UserID:       user.ID,
					Note:         "some note",
					Created:      tstTs,
					LastModified: tstTs,
				}},
				nil,
			)
		mock.EXPECT().
			GetEventItemByID(ctx, eventItemID).
			Return(
				&model.EventItem{
					ID:           eventItemID,
					RefID:        util.Must(model.NewEventItemRefID()),
					EventID:      eventID,
					Description:  "some desc",
					Created:      tstTs,
					LastModified: tstTs,
				}, nil,
			)
		mock.EXPECT().
			GetUserByID(ctx, user.ID).
			Return(user, nil)

		request := icbt.EarmarksListRequest_builder{
			Archived: func(b bool) *bool { return &b }(false),
		}.Build()
		response, err := server.EarmarksList(ctx, connect.NewRequest(request))
		assert.Nil(t, err)

		assert.Equal(t, len(response.Msg.GetEarmarks()), 1)
	})
}

func TestRpc_ListEventEarmarks(t *testing.T) {
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

	t.Run("list event earmarks should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		earmarkRefID := util.Must(model.NewEarmarkRefID())
		eventRefID := util.Must(model.NewEventRefID())

		eventID := 1
		eventItemID := 4
		earmarkID := 13

		mock.EXPECT().
			GetEvent(ctx, eventRefID).
			Return(
				&model.Event{
					ID:            eventID,
					RefID:         eventRefID,
					UserID:        user.ID,
					Name:          "some event",
					Description:   "some desc",
					ItemSortOrder: []int{},
					Archived:      false,
				},
				nil,
			)
		mock.EXPECT().
			GetEarmarksByEventID(ctx, eventID).
			Return(
				[]*model.Earmark{{
					ID:          earmarkID,
					RefID:       earmarkRefID,
					EventItemID: eventItemID,
					UserID:      user.ID,
					Note:        "some note",
				}},
				nil,
			)
		mock.EXPECT().
			GetEventItemByID(ctx, eventItemID).
			Return(
				&model.EventItem{
					ID:           eventItemID,
					RefID:        util.Must(model.NewEventItemRefID()),
					EventID:      eventID,
					Description:  "some desc",
					Created:      tstTs,
					LastModified: tstTs,
				}, nil,
			)
		mock.EXPECT().
			GetUserByID(ctx, user.ID).
			Return(user, nil)

		request := icbt.EventListEarmarksRequest_builder{
			RefId: eventRefID.String(),
		}.Build()
		response, err := server.EventListEarmarks(ctx, connect.NewRequest(request))
		assert.Nil(t, err)

		assert.Equal(t, len(response.Msg.GetEarmarks()), 1)
	})

	t.Run("list event earmarks event not found should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := util.Must(model.NewEventRefID())

		mock.EXPECT().
			GetEvent(ctx, eventRefID).
			Return(nil, errs.NotFound.Error("event not found"))

		request := icbt.EventListEarmarksRequest_builder{
			RefId: eventRefID.String(),
		}.Build()
		_, err := server.EventListEarmarks(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeNotFound, "event not found")
	})

	t.Run("list event earmarks bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, _ := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		request := icbt.EventListEarmarksRequest_builder{
			RefId: "hodor",
		}.Build()
		_, err := server.EventListEarmarks(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeInvalidArgument, "bad event ref-id")
	})
}

func TestRpc_GetEarmarkDetails(t *testing.T) {
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

	t.Run("get earmark details should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		earmarkRefID := util.Must(model.NewEarmarkRefID())
		earmarkID := 5
		eventItemID := 7
		eventID := 1

		mock.EXPECT().
			GetEarmark(ctx, earmarkRefID).
			Return(
				&model.Earmark{
					ID:          earmarkID,
					RefID:       earmarkRefID,
					EventItemID: eventItemID,
					UserID:      user.ID,
					Note:        "some note",
				}, nil,
			)
		mock.EXPECT().
			GetEventItemByID(ctx, eventItemID).
			Return(
				&model.EventItem{
					ID:           eventItemID,
					RefID:        util.Must(model.NewEventItemRefID()),
					EventID:      eventID,
					Description:  "some desc",
					Created:      tstTs,
					LastModified: tstTs,
				}, nil,
			)
		mock.EXPECT().
			GetEventByID(ctx, eventID).
			Return(
				&model.Event{
					ID:            eventID,
					RefID:         util.Must(model.NewEventRefID()),
					UserID:        user.ID,
					Name:          "some event",
					Description:   "some desc",
					ItemSortOrder: []int{},
					Archived:      false,
				}, nil,
			)
		mock.EXPECT().
			GetEventItemByID(ctx, eventItemID).
			Return(
				&model.EventItem{
					ID:           eventItemID,
					RefID:        util.Must(model.NewEventItemRefID()),
					EventID:      eventID,
					Description:  "some desc",
					Created:      tstTs,
					LastModified: tstTs,
				}, nil,
			)
		mock.EXPECT().
			GetUserByID(ctx, user.ID).
			Return(user, nil)

		request := icbt.EarmarkGetDetailsRequest_builder{
			RefId: earmarkRefID.String(),
		}.Build()
		response, err := server.EarmarkGetDetails(ctx, connect.NewRequest(request))
		assert.Nil(t, err)

		assert.Equal(t, response.Msg.GetEarmark().GetRefId(), earmarkRefID.String())
	})

	t.Run("get earmark details bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, _ := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		request := icbt.EarmarkGetDetailsRequest_builder{
			RefId: "hodor",
		}.Build()
		_, err := server.EarmarkGetDetails(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeInvalidArgument, "bad earmark ref-id")
	})
}

func TestRpc_AddEarmark(t *testing.T) {
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

	t.Run("add earmark should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := util.Must(model.NewEventItemRefID())
		earmarkRefID := util.Must(model.NewEarmarkRefID())
		eventItemID := 33
		eventID := 22
		note := "some note"

		mock.EXPECT().
			GetEventItem(ctx, eventItemRefID).
			Return(
				&model.EventItem{
					ID:           eventItemID,
					RefID:        util.Must(model.NewEventItemRefID()),
					EventID:      eventID,
					Description:  "some desc",
					Created:      tstTs,
					LastModified: tstTs,
				}, nil,
			)
		mock.EXPECT().
			NewEarmark(ctx, user, eventItemID, note).
			Return(
				&model.Earmark{
					ID:          eventItemID,
					RefID:       earmarkRefID,
					EventItemID: eventItemID,
					UserID:      user.ID,
					Note:        note,
				}, nil,
			)
		mock.EXPECT().
			GetEventItemByID(ctx, eventItemID).
			Return(
				&model.EventItem{
					ID:           eventItemID,
					RefID:        util.Must(model.NewEventItemRefID()),
					EventID:      eventID,
					Description:  "some desc",
					Created:      tstTs,
					LastModified: tstTs,
				}, nil,
			)
		mock.EXPECT().
			GetUserByID(ctx, user.ID).
			Return(user, nil)

		request := icbt.EarmarkCreateRequest_builder{
			EventItemRefId: eventItemRefID.String(),
			Note:           "some note",
		}.Build()
		response, err := server.EarmarkCreate(ctx, connect.NewRequest(request))
		assert.Nil(t, err)

		assert.Equal(t, response.Msg.GetEarmark().GetRefId(), earmarkRefID.String())
	})

	t.Run("add earmark for already earmarked by self should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := util.Must(model.NewEventItemRefID())
		eventItemID := 33
		eventID := 22
		note := "some note"

		mock.EXPECT().
			GetEventItem(ctx, eventItemRefID).
			Return(
				&model.EventItem{
					ID:           eventItemID,
					RefID:        util.Must(model.NewEventItemRefID()),
					EventID:      eventID,
					Description:  "some desc",
					Created:      tstTs,
					LastModified: tstTs,
				}, nil,
			)
		mock.EXPECT().
			NewEarmark(ctx, user, eventItemID, note).
			Return(nil, errs.AlreadyExists.Error("already earmarked"))

		request := icbt.EarmarkCreateRequest_builder{
			EventItemRefId: eventItemRefID.String(),
			Note:           "some note",
		}.Build()
		_, err := server.EarmarkCreate(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeAlreadyExists, "already earmarked")
	})

	t.Run("add earmark user not verified should fail", func(t *testing.T) {
		t.Parallel()

		user := &model.User{
			ID:           1,
			RefID:        util.Must(model.NewUserRefID()),
			Email:        "user@example.com",
			Name:         "user",
			PWHash:       []byte("00x00"),
			Verified:     false,
			Created:      tstTs,
			LastModified: tstTs,
		}

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := util.Must(model.NewEventItemRefID())
		eventItemID := 33
		eventID := 22
		note := "some note"

		mock.EXPECT().
			GetEventItem(ctx, eventItemRefID).
			Return(
				&model.EventItem{
					ID:           eventItemID,
					RefID:        util.Must(model.NewEventItemRefID()),
					EventID:      eventID,
					Description:  "some desc",
					Created:      tstTs,
					LastModified: tstTs,
				}, nil,
			)
		mock.EXPECT().
			NewEarmark(ctx, user, eventItemID, note).
			Return(nil, errs.PermissionDenied.Error("Account must be verified before earmarking is allowed."))

		request := icbt.EarmarkCreateRequest_builder{
			EventItemRefId: eventItemRefID.String(),
			Note:           "some note",
		}.Build()
		_, err := server.EarmarkCreate(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodePermissionDenied,
			"Account must be verified before earmarking is allowed.")
	})

	t.Run("add earmark for already earmarked by other should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := util.Must(model.NewEventItemRefID())
		eventItemID := 33
		eventID := 22
		note := "some note"

		mock.EXPECT().
			GetEventItem(ctx, eventItemRefID).
			Return(
				&model.EventItem{
					ID:           eventItemID,
					RefID:        util.Must(model.NewEventItemRefID()),
					EventID:      eventID,
					Description:  "some desc",
					Created:      tstTs,
					LastModified: tstTs,
				}, nil,
			)
		mock.EXPECT().
			NewEarmark(ctx, user, eventItemID, note).
			Return(nil, errs.AlreadyExists.Error("already earmarked by other user"))

		request := icbt.EarmarkCreateRequest_builder{
			EventItemRefId: eventItemRefID.String(),
			Note:           "some note",
		}.Build()
		_, err := server.EarmarkCreate(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeAlreadyExists, "already earmarked by other user")
	})

	t.Run("add earmark with bad event item refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, _ := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		request := icbt.EarmarkCreateRequest_builder{
			EventItemRefId: "hodor",
			Note:           "some note",
		}.Build()
		_, err := server.EarmarkCreate(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeInvalidArgument, "bad event-item ref-id")
	})
}

func TestRpc_RemoveEarmark(t *testing.T) {
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

	t.Run("remove earmark should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		earmarkRefID := util.Must(model.NewEarmarkRefID())

		mock.EXPECT().
			DeleteEarmarkByRefID(ctx, user.ID, earmarkRefID).
			Return(nil)

		request := icbt.EarmarkRemoveRequest_builder{
			RefId: earmarkRefID.String(),
		}.Build()
		_, err := server.EarmarkRemove(ctx, connect.NewRequest(request))
		assert.Nil(t, err)
	})

	t.Run("remove earmark for another user should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		earmarkRefID := util.Must(model.NewEarmarkRefID())

		mock.EXPECT().
			DeleteEarmarkByRefID(ctx, user.ID, earmarkRefID).
			Return(errs.PermissionDenied.Error("permission denied"))

		request := icbt.EarmarkRemoveRequest_builder{
			RefId: earmarkRefID.String(),
		}.Build()

		_, err := server.EarmarkRemove(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodePermissionDenied, "permission denied")
	})

	t.Run("remove earmark for bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, _ := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		request := icbt.EarmarkRemoveRequest_builder{
			RefId: "hodor",
		}.Build()
		_, err := server.EarmarkRemove(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeInvalidArgument, "bad earmark ref-id")
	})

	t.Run("remove earmark for archived event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		earmarkRefID := util.Must(model.NewEarmarkRefID())

		mock.EXPECT().
			DeleteEarmarkByRefID(ctx, user.ID, earmarkRefID).
			Return(errs.PermissionDenied.Error("event is archived"))

		request := icbt.EarmarkRemoveRequest_builder{
			RefId: earmarkRefID.String(),
		}.Build()
		_, err := server.EarmarkRemove(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodePermissionDenied, "event is archived")
	})
}
