package rpc

import (
	"context"
	"testing"
	"time"

	"github.com/dropwhile/refid/v2"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/twitchtv/twirp"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/convert"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/middleware/auth"
	"github.com/dropwhile/icbt/internal/util"
	"github.com/dropwhile/icbt/rpc/icbt"
)

func TestServer_ListEvents(t *testing.T) {
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

	t.Run("list events paginated should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())

		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(user.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"current", "archived"}).
				AddRow(1, 1),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(pgx.NamedArgs{
				"userID":   user.ID,
				"limit":    pgxmock.AnyArg(),
				"offset":   pgxmock.AnyArg(),
				"archived": false,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"user_id", "archived",
					"name", "description",
					"start_time", "start_time_tz",
					"created", "last_modified",
				}).
				AddRow(
					1, eventRefID,
					user.ID, false,
					"some name", "some description",
					tstTs, model.Must(model.ParseTimeZone("Etc/UTC")),
					tstTs, tstTs,
				),
			)

		request := &icbt.ListEventsRequest{
			Pagination: &icbt.PaginationRequest{Limit: 10, Offset: 0},
			Archived:   func(b bool) *bool { return &b }(false),
		}
		response, err := server.ListEvents(ctx, request)
		assert.NilError(t, err)

		assert.Check(t, len(response.Events) == 1)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("list events non-paginated should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{
			Db: mock,
		}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())

		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(pgx.NamedArgs{
				"userID":   user.ID,
				"archived": false,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"user_id", "archived",
					"name", "description",
					"start_time", "start_time_tz",
					"created", "last_modified",
				}).
				AddRow(
					1, eventRefID,
					user.ID, false,
					"some name", "some description",
					tstTs, model.Must(model.ParseTimeZone("Etc/UTC")),
					tstTs, tstTs,
				),
			)

		request := &icbt.ListEventsRequest{
			Archived: func(b bool) *bool { return &b }(false),
		}
		response, err := server.ListEvents(ctx, request)
		assert.NilError(t, err)

		assert.Check(t, len(response.Events) == 1)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestServer_GetEventDetails(t *testing.T) {
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

	t.Run("get event details should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventID := 1
		eventItemID := 3
		earmarkID := 4

		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(eventRefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"user_id", "archived",
					"name", "description",
					"start_time", "start_time_tz",
					"created", "last_modified",
				}).
				AddRow(
					eventID, eventRefID,
					user.ID, false,
					"some name", "some description",
					tstTs, model.Must(model.ParseTimeZone("Etc/UTC")),
					tstTs, tstTs,
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
					eventItemID, refid.Must(model.NewEventItemRefID()),
					eventID, "some description",
					tstTs, tstTs,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM earmark_").
			WithArgs(eventID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id",
					"event_item_id", "note",
					"created", "last_modified",
				}).
				AddRow(
					earmarkID, refid.Must(model.NewEarmarkRefID()),
					user.ID, eventItemID, "some note",
					tstTs, tstTs,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_item_").
			WithArgs(eventItemID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"event_id", "description",
					"created", "last_modified",
				}).
				AddRow(
					eventItemID, refid.Must(model.NewEventItemRefID()),
					eventID, "some description",
					tstTs, tstTs,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM user_").
			WithArgs(user.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "email", "name",
					"created", "last_modified",
				}).
				AddRow(
					user.ID, user.RefID, "user@example.com", "user",
					tstTs, tstTs,
				),
			)

		request := &icbt.GetEventDetailsRequest{
			RefId: eventRefID.String(),
		}
		response, err := server.GetEventDetails(ctx, request)
		assert.NilError(t, err)

		assert.Equal(t, response.Event.Name, "some name")
		assert.Equal(t, len(response.Items), 1)
		assert.Equal(t, len(response.Earmarks), 1)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get event details event not found should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())

		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(eventRefID).
			WillReturnError(pgx.ErrNoRows)

		request := &icbt.GetEventDetailsRequest{
			RefId: eventRefID.String(),
		}
		_, err := server.GetEventDetails(ctx, request)
		assertTwirpError(t, err, twirp.NotFound, "event not found")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
	t.Run("get event details with bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.GetEventDetailsRequest{
			RefId: "hodor",
		}
		_, err := server.GetEventDetails(ctx, request)
		assertTwirpError(t, err, twirp.InvalidArgument, "ref_id bad event ref-id")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestServer_CreateEvent(t *testing.T) {
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

	event := &model.Event{
		ID:           1,
		RefID:        refid.Must(model.NewEventRefID()),
		UserID:       user.ID,
		Name:         "event",
		Description:  "description",
		Archived:     false,
		StartTime:    tstTs,
		StartTimeTz:  model.Must(model.ParseTimeZone("Etc/UTC")),
		Created:      tstTs,
		LastModified: tstTs,
	}

	t.Run("create event should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)

		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectQuery("INSERT INTO event_ ").
			WithArgs(pgx.NamedArgs{
				"refID":       model.EventRefIDMatcher,
				"userID":      event.UserID,
				"name":        event.Name,
				"description": event.Description,
				"startTime": util.CloseTimeMatcher{
					Value: event.StartTime, Within: time.Minute,
				},
				"startTimeTz": event.StartTimeTz,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived", "start_time", "start_time_tz",
					"created", "last_modified",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name,
					event.Description, event.Archived,
					event.StartTime, event.StartTimeTz,
					tstTs, tstTs,
				))
		mock.ExpectCommit()
		mock.ExpectRollback()

		request := &icbt.CreateEventRequest{
			Name:        event.Name,
			Description: event.Description,
			When: convert.TimeToTimestampTZ(
				event.StartTime.In(event.StartTimeTz.Location)),
		}
		response, err := server.CreateEvent(ctx, request)
		assert.NilError(t, err)

		assert.Equal(t, response.Event.Name, event.Name)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create event with empty TZ should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.CreateEventRequest{
			Name:        event.Name,
			Description: event.Description,
			When: &icbt.TimestampTZ{
				Ts: convert.TimeToTimestamp(event.StartTime),
				Tz: "",
			},
		}
		_, err := server.CreateEvent(ctx, request)
		assertTwirpError(t, err, twirp.InvalidArgument, "bad empty value",
			map[string]string{"argument": "tz"})
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create event with empty ts should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.CreateEventRequest{
			Name:        event.Name,
			Description: event.Description,
			When: &icbt.TimestampTZ{
				Tz: "UTC",
			},
		}
		_, err := server.CreateEvent(ctx, request)
		assertTwirpError(t, err, twirp.InvalidArgument, "start_time bad empty value",
			map[string]string{"argument": "start_time"})
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create event with empty name should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.CreateEventRequest{
			Name:        "",
			Description: event.Description,
			When: convert.TimeToTimestampTZ(
				event.StartTime.In(event.StartTimeTz.Location)),
		}
		_, err := server.CreateEvent(ctx, request)
		assertTwirpError(t, err, twirp.InvalidArgument, "bad empty value",
			map[string]string{"argument": "name"})
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create event with empty description should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.CreateEventRequest{
			Name:        event.Name,
			Description: "",
			When: convert.TimeToTimestampTZ(
				event.StartTime.In(event.StartTimeTz.Location)),
		}
		_, err := server.CreateEvent(ctx, request)
		assertTwirpError(t, err, twirp.InvalidArgument, "bad empty value",
			map[string]string{"argument": "description"})
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestServer_UpdateEvent(t *testing.T) {
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

	event := &model.Event{
		ID:           1,
		RefID:        refid.Must(model.NewEventRefID()),
		UserID:       user.ID,
		Name:         "event",
		Description:  "description",
		Archived:     false,
		StartTime:    tstTs,
		StartTimeTz:  model.Must(model.ParseTimeZone("Etc/UTC")),
		Created:      tstTs,
		LastModified: tstTs,
	}

	t.Run("update event should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventID := 1

		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(eventRefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"user_id", "archived",
					"name", "description",
					"start_time", "start_time_tz",
					"created", "last_modified",
				}).
				AddRow(
					eventID, eventRefID,
					user.ID, false,
					"some name", "some description",
					tstTs, model.Must(model.ParseTimeZone("Etc/UTC")),
					tstTs, tstTs,
				),
			)
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectExec("UPDATE event_ ").
			WithArgs(pgx.NamedArgs{
				"eventID":       eventID,
				"name":          event.Name,
				"description":   event.Description,
				"itemSortOrder": pgxmock.AnyArg(),
				"startTime": util.CloseTimeMatcher{
					Value: event.StartTime, Within: time.Minute,
				},
				"startTimeTz": event.StartTimeTz,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		request := &icbt.UpdateEventRequest{
			RefId:       eventRefID.String(),
			Name:        &event.Name,
			Description: &event.Description,
			When: convert.TimeToTimestampTZ(
				event.StartTime.In(event.StartTimeTz.Location)),
		}
		response, err := server.UpdateEvent(ctx, request)
		assert.NilError(t, err)
		assert.Equal(t, response.Event.Name, event.Name)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update event with empty TZ should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventID := 1

		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(eventRefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"user_id", "archived",
					"name", "description",
					"start_time", "start_time_tz",
					"created", "last_modified",
				}).
				AddRow(
					eventID, eventRefID,
					user.ID, false,
					"some name", "some description",
					tstTs, model.Must(model.ParseTimeZone("Etc/UTC")),
					tstTs, tstTs,
				),
			)
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectExec("UPDATE event_ ").
			WithArgs(pgx.NamedArgs{
				"eventID":       eventID,
				"name":          event.Name,
				"description":   event.Description,
				"itemSortOrder": pgxmock.AnyArg(),
				"startTime": util.CloseTimeMatcher{
					Value: event.StartTime, Within: time.Minute,
				},
				"startTimeTz": event.StartTimeTz,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		request := &icbt.UpdateEventRequest{
			RefId:       eventRefID.String(),
			Name:        &event.Name,
			Description: &event.Description,
			When: &icbt.TimestampTZ{
				Ts: convert.TimeToTimestamp(event.StartTime),
				Tz: "",
			},
		}
		response, err := server.UpdateEvent(ctx, request)
		assert.NilError(t, err)
		assert.Equal(t, response.Event.Name, event.Name)
		// tz not updated/changed
		assert.Equal(t, response.Event.When.Tz, event.StartTimeTz.String())
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update event with empty Ts should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventID := 1

		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(eventRefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"user_id", "archived",
					"name", "description",
					"start_time", "start_time_tz",
					"created", "last_modified",
				}).
				AddRow(
					eventID, eventRefID,
					user.ID, false,
					"some name", "some description",
					tstTs, model.Must(model.ParseTimeZone("Etc/UTC")),
					tstTs, tstTs,
				),
			)
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectExec("UPDATE event_ ").
			WithArgs(pgx.NamedArgs{
				"eventID":       eventID,
				"name":          event.Name,
				"description":   event.Description,
				"itemSortOrder": pgxmock.AnyArg(),
				"startTime": util.CloseTimeMatcher{
					Value: event.StartTime, Within: time.Minute,
				},
				"startTimeTz": event.StartTimeTz,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		request := &icbt.UpdateEventRequest{
			RefId:       eventRefID.String(),
			Name:        &event.Name,
			Description: &event.Description,
			When: &icbt.TimestampTZ{
				Tz: event.StartTimeTz.String(),
			},
		}
		response, err := server.UpdateEvent(ctx, request)
		assert.NilError(t, err)
		assert.Equal(t, response.Event.Name, event.Name)
		// ts not updated/changed
		assert.Equal(t, response.Event.When.Ts.AsTime(), event.StartTime)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update event with empty name should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventID := 1

		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(eventRefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"user_id", "archived",
					"name", "description",
					"start_time", "start_time_tz",
					"created", "last_modified",
				}).
				AddRow(
					eventID, eventRefID,
					user.ID, false,
					event.Name, "some description",
					tstTs, model.Must(model.ParseTimeZone("Etc/UTC")),
					tstTs, tstTs,
				),
			)
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectExec("UPDATE event_ ").
			WithArgs(pgx.NamedArgs{
				"eventID":       eventID,
				"name":          event.Name,
				"description":   event.Description,
				"itemSortOrder": pgxmock.AnyArg(),
				"startTime": util.CloseTimeMatcher{
					Value: event.StartTime, Within: time.Minute,
				},
				"startTimeTz": event.StartTimeTz,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		request := &icbt.UpdateEventRequest{
			RefId:       eventRefID.String(),
			Description: &event.Description,
			When: &icbt.TimestampTZ{
				Tz: event.StartTimeTz.String(),
			},
		}
		response, err := server.UpdateEvent(ctx, request)
		assert.NilError(t, err)
		assert.Equal(t, response.Event.Name, event.Name)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update event with empty description should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventID := 1

		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(eventRefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"user_id", "archived",
					"name", "description",
					"start_time", "start_time_tz",
					"created", "last_modified",
				}).
				AddRow(
					eventID, eventRefID,
					user.ID, false,
					event.Name, event.Description,
					tstTs, model.Must(model.ParseTimeZone("Etc/UTC")),
					tstTs, tstTs,
				),
			)
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectExec("UPDATE event_ ").
			WithArgs(pgx.NamedArgs{
				"eventID":       eventID,
				"name":          event.Name,
				"description":   event.Description,
				"itemSortOrder": pgxmock.AnyArg(),
				"startTime": util.CloseTimeMatcher{
					Value: event.StartTime, Within: time.Minute,
				},
				"startTimeTz": event.StartTimeTz,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		request := &icbt.UpdateEventRequest{
			RefId: eventRefID.String(),
			Name:  &event.Name,
			When: &icbt.TimestampTZ{
				Tz: event.StartTimeTz.String(),
			},
		}
		response, err := server.UpdateEvent(ctx, request)
		assert.NilError(t, err)
		assert.Equal(t, response.Event.Description, event.Description)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update event with no data should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventID := 1

		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(eventRefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"user_id", "archived",
					"name", "description",
					"start_time", "start_time_tz",
					"created", "last_modified",
				}).
				AddRow(
					eventID, eventRefID,
					user.ID, false,
					event.Name, event.Description,
					tstTs, model.Must(model.ParseTimeZone("Etc/UTC")),
					tstTs, tstTs,
				),
			)

		request := &icbt.UpdateEventRequest{
			RefId: eventRefID.String(),
		}
		_, err := server.UpdateEvent(ctx, request)
		assertTwirpError(t, err, twirp.InvalidArgument, "missing fields")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update archived event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventID := 1

		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(eventRefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"user_id", "archived",
					"name", "description",
					"start_time", "start_time_tz",
					"created", "last_modified",
				}).
				AddRow(
					eventID, eventRefID,
					user.ID, true,
					event.Name, event.Description,
					tstTs, model.Must(model.ParseTimeZone("Etc/UTC")),
					tstTs, tstTs,
				),
			)

		request := &icbt.UpdateEventRequest{
			RefId:       eventRefID.String(),
			Description: &event.Description,
			When: &icbt.TimestampTZ{
				Tz: event.StartTimeTz.String(),
			},
		}
		_, err := server.UpdateEvent(ctx, request)
		assertTwirpError(t, err, twirp.PermissionDenied, "event is archived")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update event owned by other user should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)
		eventRefID := refid.Must(model.NewEventRefID())
		eventID := 1

		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(eventRefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"user_id", "archived",
					"name", "description",
					"start_time", "start_time_tz",
					"created", "last_modified",
				}).
				AddRow(
					eventID, eventRefID,
					33, false,
					event.Name, event.Description,
					tstTs, model.Must(model.ParseTimeZone("Etc/UTC")),
					tstTs, tstTs,
				),
			)

		request := &icbt.UpdateEventRequest{
			RefId:       eventRefID.String(),
			Description: &event.Description,
			When: &icbt.TimestampTZ{
				Tz: event.StartTimeTz.String(),
			},
		}
		_, err := server.UpdateEvent(ctx, request)
		assertTwirpError(t, err, twirp.PermissionDenied, "permission denied")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update event with bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := &Server{Db: mock}
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.UpdateEventRequest{
			RefId:       "hodor",
			Description: &event.Description,
			When: &icbt.TimestampTZ{
				Tz: event.StartTimeTz.String(),
			},
		}
		_, err := server.UpdateEvent(ctx, request)
		assertTwirpError(t, err, twirp.InvalidArgument, "ref_id bad event ref-id")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestServer_DeleteEvent(t *testing.T) {
}
