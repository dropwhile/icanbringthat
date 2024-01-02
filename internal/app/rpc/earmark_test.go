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
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		earmarkRefID := refid.Must(model.NewEarmarkRefID())

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
					Limit:  uint32(limit),
					Offset: uint32(offset),
					Count:  1,
				},
				nil,
			).
			Once()
		mock.EXPECT().
			GetEventItemByID(ctx, eventItemID).
			Return(
				&model.EventItem{
					ID:           eventItemID,
					RefID:        refid.Must(model.NewEventItemRefID()),
					EventID:      eventID,
					Description:  "some desc",
					Created:      tstTs,
					LastModified: tstTs,
				}, nil,
			).
			Once()
		mock.EXPECT().
			GetUserByID(ctx, user.ID).
			Return(user, nil).
			Once()

		request := &icbt.ListEarmarksRequest{
			Pagination: &icbt.PaginationRequest{Limit: 10, Offset: 0},
			Archived:   func(b bool) *bool { return &b }(false),
		}
		response, err := server.ListEarmarks(ctx, request)
		assert.NilError(t, err)

		assert.Equal(t, len(response.Earmarks), 1)
		mock.AssertExpectations(t)
	})

	t.Run("list earmarks non-paginated should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		earmarkRefID := refid.Must(model.NewEarmarkRefID())

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
			).
			Once()
		mock.EXPECT().
			GetEventItemByID(ctx, eventItemID).
			Return(
				&model.EventItem{
					ID:           eventItemID,
					RefID:        refid.Must(model.NewEventItemRefID()),
					EventID:      eventID,
					Description:  "some desc",
					Created:      tstTs,
					LastModified: tstTs,
				}, nil,
			).
			Once()
		mock.EXPECT().
			GetUserByID(ctx, user.ID).
			Return(user, nil).
			Once()

		request := &icbt.ListEarmarksRequest{
			Archived: func(b bool) *bool { return &b }(false),
		}
		response, err := server.ListEarmarks(ctx, request)
		assert.NilError(t, err)

		assert.Equal(t, len(response.Earmarks), 1)
		mock.AssertExpectations(t)
	})
}

func TestRpc_ListEventEarmarks(t *testing.T) {
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

	t.Run("list event earmarks should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		earmarkRefID := refid.Must(model.NewEarmarkRefID())
		eventRefID := refid.Must(model.NewEventRefID())

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
			).
			Once()
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
			).
			Once()
		mock.EXPECT().
			GetEventItemByID(ctx, eventItemID).
			Return(
				&model.EventItem{
					ID:           eventItemID,
					RefID:        refid.Must(model.NewEventItemRefID()),
					EventID:      eventID,
					Description:  "some desc",
					Created:      tstTs,
					LastModified: tstTs,
				}, nil,
			).
			Once()
		mock.EXPECT().
			GetUserByID(ctx, user.ID).
			Return(user, nil).
			Once()

		request := &icbt.ListEventEarmarksRequest{
			RefId: eventRefID.String(),
		}
		response, err := server.ListEventEarmarks(ctx, request)
		assert.NilError(t, err)

		assert.Equal(t, len(response.Earmarks), 1)
		mock.AssertExpectations(t)
	})

	t.Run("list event earmarks event not found should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())

		mock.EXPECT().
			GetEvent(ctx, eventRefID).
			Return(nil, errs.NotFound.Error("event not found")).
			Once()

		request := &icbt.ListEventEarmarksRequest{
			RefId: eventRefID.String(),
		}
		_, err := server.ListEventEarmarks(ctx, request)
		errs.AssertError(t, err, twirp.NotFound, "event not found")
		mock.AssertExpectations(t)
	})

	t.Run("list event earmarks bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.ListEventEarmarksRequest{
			RefId: "hodor",
		}
		_, err := server.ListEventEarmarks(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "ref_id bad event ref-id")
		mock.AssertExpectations(t)
	})
}

func TestRpc_GetEarmarkDetails(t *testing.T) {
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

	t.Run("get earmark details should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		earmarkRefID := refid.Must(model.NewEarmarkRefID())
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
			).
			Once()
		mock.EXPECT().
			GetEventItemByID(ctx, eventItemID).
			Return(
				&model.EventItem{
					ID:           eventItemID,
					RefID:        refid.Must(model.NewEventItemRefID()),
					EventID:      eventID,
					Description:  "some desc",
					Created:      tstTs,
					LastModified: tstTs,
				}, nil,
			).
			Once()
		mock.EXPECT().
			GetEventByID(ctx, eventID).
			Return(
				&model.Event{
					ID:            eventID,
					RefID:         refid.Must(model.NewEventRefID()),
					UserID:        user.ID,
					Name:          "some event",
					Description:   "some desc",
					ItemSortOrder: []int{},
					Archived:      false,
				}, nil,
			).
			Once()
		mock.EXPECT().
			GetEventItemByID(ctx, eventItemID).
			Return(
				&model.EventItem{
					ID:           eventItemID,
					RefID:        refid.Must(model.NewEventItemRefID()),
					EventID:      eventID,
					Description:  "some desc",
					Created:      tstTs,
					LastModified: tstTs,
				}, nil,
			).
			Once()
		mock.EXPECT().
			GetUserByID(ctx, user.ID).
			Return(user, nil).
			Once()

		request := &icbt.GetEarmarkDetailsRequest{
			RefId: earmarkRefID.String(),
		}
		response, err := server.GetEarmarkDetails(ctx, request)
		assert.NilError(t, err)

		assert.Equal(t, response.Earmark.RefId, earmarkRefID.String())
		mock.AssertExpectations(t)
	})

	t.Run("get earmark details bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.GetEarmarkDetailsRequest{
			RefId: "hodor",
		}
		_, err := server.GetEarmarkDetails(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "ref_id bad earmark ref-id")
		mock.AssertExpectations(t)
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
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		earmarkRefID := refid.Must(model.NewEarmarkRefID())
		eventItemID := 33
		eventID := 22
		note := "some note"

		mock.EXPECT().
			GetEventItem(ctx, eventItemRefID).
			Return(
				&model.EventItem{
					ID:           eventItemID,
					RefID:        refid.Must(model.NewEventItemRefID()),
					EventID:      eventID,
					Description:  "some desc",
					Created:      tstTs,
					LastModified: tstTs,
				}, nil,
			).
			Once()
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
			).
			Once()
		mock.EXPECT().
			GetEventItemByID(ctx, eventItemID).
			Return(
				&model.EventItem{
					ID:           eventItemID,
					RefID:        refid.Must(model.NewEventItemRefID()),
					EventID:      eventID,
					Description:  "some desc",
					Created:      tstTs,
					LastModified: tstTs,
				}, nil,
			).
			Once()
		mock.EXPECT().
			GetUserByID(ctx, user.ID).
			Return(user, nil).
			Once()

		request := &icbt.CreateEarmarkRequest{
			EventItemRefId: eventItemRefID.String(),
			Note:           "some note",
		}
		response, err := server.CreateEarmark(ctx, request)
		assert.NilError(t, err)

		assert.Equal(t, response.Earmark.RefId, earmarkRefID.String())
		mock.AssertExpectations(t)
	})

	t.Run("add earmark for already earmarked by self should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		eventItemID := 33
		eventID := 22
		note := "some note"

		mock.EXPECT().
			GetEventItem(ctx, eventItemRefID).
			Return(
				&model.EventItem{
					ID:           eventItemID,
					RefID:        refid.Must(model.NewEventItemRefID()),
					EventID:      eventID,
					Description:  "some desc",
					Created:      tstTs,
					LastModified: tstTs,
				}, nil,
			).
			Once()
		mock.EXPECT().
			NewEarmark(ctx, user, eventItemID, note).
			Return(nil, errs.AlreadyExists.Error("already earmarked")).
			Once()

		request := &icbt.CreateEarmarkRequest{
			EventItemRefId: eventItemRefID.String(),
			Note:           "some note",
		}
		_, err := server.CreateEarmark(ctx, request)
		errs.AssertError(t, err, twirp.AlreadyExists, "already earmarked")
		mock.AssertExpectations(t)
	})

	t.Run("add earmark user not verified should fail", func(t *testing.T) {
		t.Parallel()

		user := &model.User{
			ID:           1,
			RefID:        refid.Must(model.NewUserRefID()),
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
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		eventItemID := 33
		eventID := 22
		note := "some note"

		mock.EXPECT().
			GetEventItem(ctx, eventItemRefID).
			Return(
				&model.EventItem{
					ID:           eventItemID,
					RefID:        refid.Must(model.NewEventItemRefID()),
					EventID:      eventID,
					Description:  "some desc",
					Created:      tstTs,
					LastModified: tstTs,
				}, nil,
			).
			Once()
		mock.EXPECT().
			NewEarmark(ctx, user, eventItemID, note).
			Return(nil, errs.PermissionDenied.Error("Account must be verified before earmarking is allowed.")).
			Once()

		request := &icbt.CreateEarmarkRequest{
			EventItemRefId: eventItemRefID.String(),
			Note:           "some note",
		}
		_, err := server.CreateEarmark(ctx, request)
		errs.AssertError(t, err, twirp.PermissionDenied,
			"Account must be verified before earmarking is allowed.")
		mock.AssertExpectations(t)
	})

	t.Run("add earmark for already earmarked by other should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		eventItemID := 33
		eventID := 22
		note := "some note"

		mock.EXPECT().
			GetEventItem(ctx, eventItemRefID).
			Return(
				&model.EventItem{
					ID:           eventItemID,
					RefID:        refid.Must(model.NewEventItemRefID()),
					EventID:      eventID,
					Description:  "some desc",
					Created:      tstTs,
					LastModified: tstTs,
				}, nil,
			).
			Once()
		mock.EXPECT().
			NewEarmark(ctx, user, eventItemID, note).
			Return(nil, errs.AlreadyExists.Error("already earmarked by other user")).
			Once()

		request := &icbt.CreateEarmarkRequest{
			EventItemRefId: eventItemRefID.String(),
			Note:           "some note",
		}
		_, err := server.CreateEarmark(ctx, request)
		errs.AssertError(t, err, twirp.AlreadyExists, "already earmarked by other user")
		mock.AssertExpectations(t)
	})

	t.Run("add earmark with bad event item refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.CreateEarmarkRequest{
			EventItemRefId: "hodor",
			Note:           "some note",
		}
		_, err := server.CreateEarmark(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "ref_id bad event-item ref-id")
		mock.AssertExpectations(t)
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
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		earmarkRefID := refid.Must(model.NewEarmarkRefID())

		mock.EXPECT().
			DeleteEarmarkByRefID(ctx, user.ID, earmarkRefID).
			Return(nil).
			Once()

		request := &icbt.RemoveEarmarkRequest{
			RefId: earmarkRefID.String(),
		}
		_, err := server.RemoveEarmark(ctx, request)
		assert.NilError(t, err)
		mock.AssertExpectations(t)
	})

	t.Run("remove earmark for another user should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		earmarkRefID := refid.Must(model.NewEarmarkRefID())

		mock.EXPECT().
			DeleteEarmarkByRefID(ctx, user.ID, earmarkRefID).
			Return(errs.PermissionDenied.Error("permission denied")).
			Once()

		request := &icbt.RemoveEarmarkRequest{
			RefId: earmarkRefID.String(),
		}
		_, err := server.RemoveEarmark(ctx, request)
		errs.AssertError(t, err, twirp.PermissionDenied, "permission denied")
		mock.AssertExpectations(t)
	})

	t.Run("remove earmark for bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.RemoveEarmarkRequest{
			RefId: "hodor",
		}
		_, err := server.RemoveEarmark(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "ref_id bad earmark ref-id")
		mock.AssertExpectations(t)
	})

	t.Run("remove earmark for archived event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		earmarkRefID := refid.Must(model.NewEarmarkRefID())

		mock.EXPECT().
			DeleteEarmarkByRefID(ctx, user.ID, earmarkRefID).
			Return(errs.PermissionDenied.Error("event is archived")).
			Once()

		request := &icbt.RemoveEarmarkRequest{
			RefId: earmarkRefID.String(),
		}
		_, err := server.RemoveEarmark(ctx, request)
		errs.AssertError(t, err, twirp.PermissionDenied, "event is archived")
		mock.AssertExpectations(t)
	})
}
