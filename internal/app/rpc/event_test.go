package rpc

import (
	"context"
	"testing"
	"time"

	"github.com/samber/mo"
	"github.com/twitchtv/twirp"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/convert"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/middleware/auth"
	"github.com/dropwhile/icbt/internal/util"
	"github.com/dropwhile/icbt/rpc/icbt"
)

func TestRpc_ListEvents(t *testing.T) {
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

	t.Run("list events paginated should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := util.Must(model.NewEventRefID())

		limit := 10
		offset := 0
		archived := false

		mock.EXPECT().
			GetEventsPaginated(ctx, user.ID, limit, offset, archived).
			Return(
				[]*model.Event{{
					ID:            11,
					RefID:         eventRefID,
					UserID:        user.ID,
					Name:          "some name",
					Description:   "some desc",
					ItemSortOrder: []int{},
					StartTime:     tstTs,
					StartTimeTz:   util.Must(service.ParseTimeZone("UTC")),
					Archived:      archived,
				}},
				&service.Pagination{
					Limit:  uint32(limit),
					Offset: uint32(offset),
					Count:  1,
				}, nil,
			)

		request := &icbt.ListEventsRequest{
			Pagination: &icbt.PaginationRequest{Limit: 10, Offset: 0},
			Archived:   func(b bool) *bool { return &b }(false),
		}
		response, err := server.ListEvents(ctx, request)
		assert.NilError(t, err)
		assert.Equal(t, len(response.Events), 1)
	})

	t.Run("list events non-paginated should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := util.Must(model.NewEventRefID())

		archived := false

		mock.EXPECT().
			GetEvents(ctx, user.ID, archived).
			Return(
				[]*model.Event{{
					ID:            11,
					RefID:         eventRefID,
					UserID:        user.ID,
					Name:          "some name",
					Description:   "some desc",
					ItemSortOrder: []int{},
					StartTime:     tstTs,
					StartTimeTz:   util.Must(service.ParseTimeZone("UTC")),
					Archived:      archived,
				}}, nil,
			)

		request := &icbt.ListEventsRequest{
			Archived: func(b bool) *bool { return &b }(false),
		}
		response, err := server.ListEvents(ctx, request)
		assert.NilError(t, err)
		assert.Equal(t, len(response.Events), 1)
	})
}

func TestRpc_GetEventDetails(t *testing.T) {
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

	t.Run("get event details should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := util.Must(model.NewEventRefID())
		eventItemRefID := util.Must(model.NewEventItemRefID())
		eventID := 1
		eventItemID := 3
		earmarkID := 4

		mock.EXPECT().
			GetEvent(ctx, eventRefID).
			Return(
				&model.Event{
					ID:            eventID,
					RefID:         eventRefID,
					UserID:        user.ID,
					Name:          "some name",
					Description:   "some desc",
					ItemSortOrder: []int{},
					StartTime:     tstTs,
					StartTimeTz:   util.Must(service.ParseTimeZone("UTC")),
					Archived:      false,
				}, nil,
			)
		mock.EXPECT().
			GetEventItemsByEventID(ctx, eventID).
			Return(
				[]*model.EventItem{{
					ID:          eventID,
					RefID:       eventItemRefID,
					EventID:     eventID,
					Description: "some item",
				}}, nil,
			)
		mock.EXPECT().
			GetEarmarksByEventID(ctx, eventID).
			Return(
				[]*model.Earmark{{
					ID:          earmarkID,
					RefID:       util.Must(model.NewEarmarkRefID()),
					UserID:      user.ID,
					EventItemID: eventItemID,
					Note:        "some earmark",
				}}, nil,
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

		request := &icbt.GetEventDetailsRequest{
			RefId: eventRefID.String(),
		}
		response, err := server.GetEventDetails(ctx, request)
		assert.NilError(t, err)

		assert.Equal(t, response.Event.Name, "some name")
		assert.Equal(t, len(response.Items), 1)
		assert.Equal(t, len(response.Earmarks), 1)
	})

	t.Run("get event details event not found should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := util.Must(model.NewEventRefID())

		mock.EXPECT().
			GetEvent(ctx, eventRefID).
			Return(nil, errs.NotFound.Error("event not found"))

		request := &icbt.GetEventDetailsRequest{
			RefId: eventRefID.String(),
		}
		_, err := server.GetEventDetails(ctx, request)
		errs.AssertError(t, err, twirp.NotFound, "event not found")
	})

	t.Run("get event details with bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, _ := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.GetEventDetailsRequest{
			RefId: "hodor",
		}
		_, err := server.GetEventDetails(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "ref_id bad event ref-id")
	})
}

func TestRpc_CreateEvent(t *testing.T) {
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

	event := &model.Event{
		ID:           1,
		RefID:        util.Must(model.NewEventRefID()),
		UserID:       user.ID,
		Name:         "event",
		Description:  "description",
		Archived:     false,
		StartTime:    tstTs,
		StartTimeTz:  util.Must(service.ParseTimeZone("Etc/UTC")),
		Created:      tstTs,
		LastModified: tstTs,
	}

	t.Run("create event should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			CreateEvent(
				ctx, user, event.Name, event.Description, event.StartTime,
				event.StartTimeTz.Location.String(),
			).
			Return(event, nil)

		request := &icbt.CreateEventRequest{
			Name:        event.Name,
			Description: event.Description,
			When: convert.TimeToTimestampTZ(
				event.StartTime.In(event.StartTimeTz.Location)),
		}
		response, err := server.CreateEvent(ctx, request)
		assert.NilError(t, err)
		assert.Equal(t, response.Event.Name, event.Name)
	})

	t.Run("create event with empty TZ should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			CreateEvent(
				ctx, user, event.Name, event.Description, event.StartTime,
				"",
			).
			Return(nil, errs.InvalidArgumentError("tz", "bad value"))

		request := &icbt.CreateEventRequest{
			Name:        event.Name,
			Description: event.Description,
			When: &icbt.TimestampTZ{
				Ts: convert.TimeToTimestamp(event.StartTime),
				Tz: "",
			},
		}
		_, err := server.CreateEvent(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "tz bad value",
			map[string]string{"argument": "tz"})
	})

	t.Run("create event with empty ts should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, _ := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.CreateEventRequest{
			Name:        event.Name,
			Description: event.Description,
			When: &icbt.TimestampTZ{
				Tz: "UTC",
			},
		}
		_, err := server.CreateEvent(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "start_time bad empty value",
			map[string]string{"argument": "start_time"})
	})

	t.Run("create event with empty name should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			CreateEvent(
				ctx, user, "", event.Description, event.StartTime,
				event.StartTimeTz.Location.String(),
			).
			Return(nil, errs.InvalidArgumentError("name", "bad value"))

		request := &icbt.CreateEventRequest{
			Name:        "",
			Description: event.Description,
			When: convert.TimeToTimestampTZ(
				event.StartTime.In(event.StartTimeTz.Location)),
		}
		_, err := server.CreateEvent(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "name bad value",
			map[string]string{"argument": "name"})
	})

	t.Run("create event with empty description should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			CreateEvent(
				ctx, user, event.Name, "", event.StartTime,
				event.StartTimeTz.Location.String(),
			).
			Return(nil, errs.InvalidArgumentError("description", "bad value"))

		request := &icbt.CreateEventRequest{
			Name:        event.Name,
			Description: "",
			When: convert.TimeToTimestampTZ(
				event.StartTime.In(event.StartTimeTz.Location)),
		}
		_, err := server.CreateEvent(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "description bad value",
			map[string]string{"argument": "description"})
	})
}

func TestRpc_UpdateEvent(t *testing.T) {
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

	event := &model.Event{
		ID:           1,
		RefID:        util.Must(model.NewEventRefID()),
		UserID:       user.ID,
		Name:         "event",
		Description:  "description",
		Archived:     false,
		StartTime:    tstTs,
		StartTimeTz:  util.Must(service.ParseTimeZone("Etc/UTC")),
		Created:      tstTs,
		LastModified: tstTs,
	}

	t.Run("update event should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			UpdateEvent(ctx, user.ID, event.RefID, &service.EventUpdateValues{
				Name:        mo.Some(event.Name),
				Description: mo.Some(event.Description),
				StartTime:   mo.Some(event.StartTime),
				Tz:          mo.Some(event.StartTimeTz.Location.String()),
			}).
			Return(nil)

		request := &icbt.UpdateEventRequest{
			RefId:       event.RefID.String(),
			Name:        &event.Name,
			Description: &event.Description,
			When: convert.TimeToTimestampTZ(
				event.StartTime.In(event.StartTimeTz.Location)),
		}
		_, err := server.UpdateEvent(ctx, request)
		assert.NilError(t, err)
	})

	t.Run("update event with empty TZ should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			UpdateEvent(ctx, user.ID, event.RefID, &service.EventUpdateValues{
				Name:        mo.Some(event.Name),
				Description: mo.Some(event.Description),
				StartTime:   mo.Some(event.StartTime),
				Tz:          mo.None[string](),
			}).
			Return(nil)

		request := &icbt.UpdateEventRequest{
			RefId:       event.RefID.String(),
			Name:        &event.Name,
			Description: &event.Description,
			When: &icbt.TimestampTZ{
				Ts: convert.TimeToTimestamp(event.StartTime),
				Tz: "",
			},
		}
		_, err := server.UpdateEvent(ctx, request)
		assert.NilError(t, err)
	})

	t.Run("update event with empty Ts should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			UpdateEvent(ctx, user.ID, event.RefID, &service.EventUpdateValues{
				Name:        mo.Some(event.Name),
				Description: mo.Some(event.Description),
				StartTime:   mo.None[time.Time](),
				Tz:          mo.Some(event.StartTimeTz.Location.String()),
			}).
			Return(nil)

		request := &icbt.UpdateEventRequest{
			RefId:       event.RefID.String(),
			Name:        &event.Name,
			Description: &event.Description,
			When: &icbt.TimestampTZ{
				Tz: event.StartTimeTz.String(),
			},
		}
		_, err := server.UpdateEvent(ctx, request)
		assert.NilError(t, err)
	})

	t.Run("update event with empty name should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			UpdateEvent(ctx, user.ID, event.RefID, &service.EventUpdateValues{
				Description: mo.Some(event.Description),
				Tz:          mo.Some(event.StartTimeTz.Location.String()),
			}).
			Return(nil)

		request := &icbt.UpdateEventRequest{
			RefId:       event.RefID.String(),
			Description: &event.Description,
			When: &icbt.TimestampTZ{
				Tz: event.StartTimeTz.String(),
			},
		}
		_, err := server.UpdateEvent(ctx, request)
		assert.NilError(t, err)
	})

	t.Run("update event with empty description should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			UpdateEvent(ctx, user.ID, event.RefID, &service.EventUpdateValues{
				Name: mo.Some(event.Name),
				Tz:   mo.Some(event.StartTimeTz.Location.String()),
			}).
			Return(nil)

		request := &icbt.UpdateEventRequest{
			RefId: event.RefID.String(),
			Name:  &event.Name,
			When: &icbt.TimestampTZ{
				Tz: event.StartTimeTz.String(),
			},
		}
		_, err := server.UpdateEvent(ctx, request)
		assert.NilError(t, err)
	})

	t.Run("update event with no data should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			UpdateEvent(ctx, user.ID, event.RefID, &service.EventUpdateValues{}).
			Return(errs.InvalidArgument.Error("missing fields"))
		request := &icbt.UpdateEventRequest{
			RefId: event.RefID.String(),
		}
		_, err := server.UpdateEvent(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "missing fields")
	})

	t.Run("update archived event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			UpdateEvent(ctx, user.ID, event.RefID, &service.EventUpdateValues{
				Description: mo.Some(event.Description),
				Tz:          mo.Some(event.StartTimeTz.Location.String()),
			}).
			Return(errs.PermissionDenied.Error("event is archived"))

		request := &icbt.UpdateEventRequest{
			RefId:       event.RefID.String(),
			Description: &event.Description,
			When: &icbt.TimestampTZ{
				Tz: event.StartTimeTz.String(),
			},
		}
		_, err := server.UpdateEvent(ctx, request)
		errs.AssertError(t, err, twirp.PermissionDenied, "event is archived")
	})

	t.Run("update event owned by other user should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			UpdateEvent(ctx, user.ID, event.RefID, &service.EventUpdateValues{
				Description: mo.Some(event.Description),
				Tz:          mo.Some(event.StartTimeTz.Location.String()),
			}).
			Return(errs.PermissionDenied.Error("permission denied"))

		request := &icbt.UpdateEventRequest{
			RefId:       event.RefID.String(),
			Description: &event.Description,
			When: &icbt.TimestampTZ{
				Tz: event.StartTimeTz.String(),
			},
		}
		_, err := server.UpdateEvent(ctx, request)
		errs.AssertError(t, err, twirp.PermissionDenied, "permission denied")
	})

	t.Run("update event with bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, _ := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.UpdateEventRequest{
			RefId:       "hodor",
			Description: &event.Description,
			When: &icbt.TimestampTZ{
				Tz: event.StartTimeTz.String(),
			},
		}
		_, err := server.UpdateEvent(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "ref_id bad event ref-id")
	})
}

func TestRpc_DeleteEvent(t *testing.T) {
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

	t.Run("delete event should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := util.Must(model.NewEventRefID())

		mock.EXPECT().
			DeleteEvent(ctx, user.ID, eventRefID).
			Return(nil)

		request := &icbt.DeleteEventRequest{
			RefId: eventRefID.String(),
		}
		_, err := server.DeleteEvent(ctx, request)
		assert.NilError(t, err)
	})

	t.Run("delete event with bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, _ := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.DeleteEventRequest{
			RefId: "hodor",
		}
		_, err := server.DeleteEvent(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "ref_id bad event ref-id")
	})

	t.Run("delete event owned by other user should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := util.Must(model.NewEventRefID())

		mock.EXPECT().
			DeleteEvent(ctx, user.ID, eventRefID).
			Return(errs.PermissionDenied.Error("permission denied"))

		request := &icbt.DeleteEventRequest{
			RefId: eventRefID.String(),
		}
		_, err := server.DeleteEvent(ctx, request)
		errs.AssertError(t, err, twirp.PermissionDenied, "permission denied")
	})
}
