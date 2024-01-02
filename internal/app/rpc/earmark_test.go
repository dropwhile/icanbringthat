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
		mock := SetupDBMock(t, ctx)
		server := NewTestServerOld(mock)
		ctx = auth.ContextSet(ctx, "user", user)
		earmarkRefID := refid.Must(model.NewEarmarkRefID())
		eventID := 1

		mock.ExpectQuery("SELECT (.+) FROM earmark_").
			WithArgs(earmarkRefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id",
					"event_item_id", "note",
					"created", "last_modified",
				}).
				AddRow(
					1, earmarkRefID, user.ID,
					12, "some note",
					tstTs, tstTs,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_item_").
			WithArgs(12).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"event_id", "description",
					"created", "last_modified",
				}).
				AddRow(
					12, refid.Must(model.NewEventItemRefID()),
					eventID, "some description",
					tstTs, tstTs,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(eventID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived", "created", "last_modified",
				}).
				AddRow(
					eventID, refid.Must(model.NewEventRefID()), user.ID,
					"event name", "event desc",
					false, tstTs, tstTs,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_item_").
			WithArgs(12).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"event_id", "description",
					"created", "last_modified",
				}).
				AddRow(
					12, refid.Must(model.NewEventItemRefID()),
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

		request := &icbt.GetEarmarkDetailsRequest{
			RefId: earmarkRefID.String(),
		}
		response, err := server.GetEarmarkDetails(ctx, request)
		assert.NilError(t, err)

		assert.Equal(t, response.Earmark.RefId, earmarkRefID.String())
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get earmark details bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := NewTestServerOld(mock)
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.GetEarmarkDetailsRequest{
			RefId: "hodor",
		}
		_, err := server.GetEarmarkDetails(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "ref_id bad earmark ref-id")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
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
		mock := SetupDBMock(t, ctx)
		server := NewTestServerOld(mock)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		earmarkRefID := refid.Must(model.NewEarmarkRefID())
		eventItemID := 33
		eventID := 22

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
		mock.ExpectQuery("SELECT (.+) FROM earmark_").
			WithArgs(eventItemID).
			WillReturnError(pgx.ErrNoRows)

		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(eventItemID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived", "created", "last_modified",
				}).
				AddRow(
					eventID, refid.Must(model.NewEventRefID()), user.ID,
					"event name", "event desc",
					false, tstTs, tstTs,
				),
			)

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO earmark_").
			WithArgs(pgx.NamedArgs{
				"userID":      user.ID,
				"eventItemID": eventItemID,
				"refID":       service.EarmarkRefIDMatcher,
				"note":        "some note",
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id",
					"event_item_id", "note",
					"created", "last_modified",
				}).
				AddRow(
					1, earmarkRefID, user.ID,
					eventItemID, "some note",
					tstTs, tstTs,
				),
			)
		mock.ExpectCommit()
		mock.ExpectRollback()
		mock.ExpectQuery("SELECT (.+) FROM event_item_").
			WithArgs(eventItemID).
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

		request := &icbt.CreateEarmarkRequest{
			EventItemRefId: eventItemRefID.String(),
			Note:           "some note",
		}
		response, err := server.CreateEarmark(ctx, request)
		assert.NilError(t, err)

		assert.Equal(t, response.Earmark.RefId, earmarkRefID.String())
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("add earmark for already earmarked by self should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := NewTestServerOld(mock)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		earmarkRefID := refid.Must(model.NewEarmarkRefID())
		eventItemID := 33
		eventID := 22

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
		mock.ExpectQuery("SELECT (.+) FROM earmark_").
			WithArgs(eventItemID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id",
					"event_item_id", "note",
					"created", "last_modified",
				}).
				AddRow(
					1, earmarkRefID, user.ID,
					eventItemID, "some note",
					tstTs, tstTs,
				),
			)

		request := &icbt.CreateEarmarkRequest{
			EventItemRefId: eventItemRefID.String(),
			Note:           "some note",
		}
		_, err := server.CreateEarmark(ctx, request)
		errs.AssertError(t, err, twirp.AlreadyExists, "already earmarked")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
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
		mock := SetupDBMock(t, ctx)
		server := NewTestServerOld(mock)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		// earmarkRefID := refid.Must(model.NewEarmarkRefID())
		eventItemID := 33
		eventID := 22

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
		mock.ExpectQuery("SELECT (.+) FROM earmark_").
			WithArgs(eventItemID).
			WillReturnError(pgx.ErrNoRows)
		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(eventItemID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived", "created", "last_modified",
				}).
				AddRow(
					eventID, refid.Must(model.NewEventRefID()), 33,
					"event name", "event desc",
					false, tstTs, tstTs,
				),
			)

		request := &icbt.CreateEarmarkRequest{
			EventItemRefId: eventItemRefID.String(),
			Note:           "some note",
		}
		_, err := server.CreateEarmark(ctx, request)
		errs.AssertError(t, err, twirp.PermissionDenied,
			"Account must be verified before earmarking is allowed.")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
	t.Run("add earmark for already earmarked by other should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := NewTestServerOld(mock)
		ctx = auth.ContextSet(ctx, "user", user)
		eventItemRefID := refid.Must(model.NewEventItemRefID())
		earmarkRefID := refid.Must(model.NewEarmarkRefID())
		eventItemID := 33
		eventID := 22

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
		mock.ExpectQuery("SELECT (.+) FROM earmark_").
			WithArgs(eventItemID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id",
					"event_item_id", "note",
					"created", "last_modified",
				}).
				AddRow(
					1, earmarkRefID, 44,
					eventItemID, "some note",
					tstTs, tstTs,
				),
			)

		request := &icbt.CreateEarmarkRequest{
			EventItemRefId: eventItemRefID.String(),
			Note:           "some note",
		}
		_, err := server.CreateEarmark(ctx, request)
		errs.AssertError(t, err, twirp.AlreadyExists, "already earmarked by other user")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("add earmark with bad event item refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := NewTestServerOld(mock)
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.CreateEarmarkRequest{
			EventItemRefId: "hodor",
			Note:           "some note",
		}
		_, err := server.CreateEarmark(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "ref_id bad event-item ref-id")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
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
		mock := SetupDBMock(t, ctx)
		server := NewTestServerOld(mock)
		ctx = auth.ContextSet(ctx, "user", user)
		earmarkRefID := refid.Must(model.NewEarmarkRefID())
		eventItemID := 33
		earmarkID := 5
		eventID := 22

		mock.ExpectQuery("SELECT (.+) FROM earmark_").
			WithArgs(earmarkRefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id",
					"event_item_id", "note",
					"created", "last_modified",
				}).
				AddRow(
					earmarkID, earmarkRefID, user.ID,
					eventItemID, "some note",
					tstTs, tstTs,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(eventItemID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived", "created", "last_modified",
				}).
				AddRow(
					eventID, refid.Must(model.NewEventRefID()), user.ID,
					"event name", "event desc",
					false, tstTs, tstTs,
				),
			)
		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM earmark_").
			WithArgs(earmarkID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		request := &icbt.RemoveEarmarkRequest{
			RefId: earmarkRefID.String(),
		}
		_, err := server.RemoveEarmark(ctx, request)
		assert.NilError(t, err)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("remove earmark for another user should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := NewTestServerOld(mock)
		ctx = auth.ContextSet(ctx, "user", user)
		earmarkRefID := refid.Must(model.NewEarmarkRefID())
		eventItemID := 33
		earmarkID := 5

		mock.ExpectQuery("SELECT (.+) FROM earmark_").
			WithArgs(earmarkRefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id",
					"event_item_id", "note",
					"created", "last_modified",
				}).
				AddRow(
					earmarkID, earmarkRefID, 33,
					eventItemID, "some note",
					tstTs, tstTs,
				),
			)

		request := &icbt.RemoveEarmarkRequest{
			RefId: earmarkRefID.String(),
		}
		_, err := server.RemoveEarmark(ctx, request)
		errs.AssertError(t, err, twirp.PermissionDenied, "permission denied")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("remove earmark for bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := NewTestServerOld(mock)
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.RemoveEarmarkRequest{
			RefId: "hodor",
		}
		_, err := server.RemoveEarmark(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "ref_id bad earmark ref-id")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("remove earmark for archived event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := NewTestServerOld(mock)
		ctx = auth.ContextSet(ctx, "user", user)
		earmarkRefID := refid.Must(model.NewEarmarkRefID())
		eventItemID := 33
		earmarkID := 5
		eventID := 22

		mock.ExpectQuery("SELECT (.+) FROM earmark_").
			WithArgs(earmarkRefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id",
					"event_item_id", "note",
					"created", "last_modified",
				}).
				AddRow(
					earmarkID, earmarkRefID, user.ID,
					eventItemID, "some note",
					tstTs, tstTs,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(eventItemID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description",
					"archived", "created", "last_modified",
				}).
				AddRow(
					eventID, refid.Must(model.NewEventRefID()), user.ID,
					"event name", "event desc",
					true, tstTs, tstTs,
				),
			)

		request := &icbt.RemoveEarmarkRequest{
			RefId: earmarkRefID.String(),
		}
		_, err := server.RemoveEarmark(ctx, request)
		errs.AssertError(t, err, twirp.PermissionDenied, "event is archived")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
