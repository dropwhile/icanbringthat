package rpc

import (
	"context"
	"log/slog"

	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/app/convert"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/middleware/auth"
	"github.com/dropwhile/icbt/rpc/icbt"
)

func (s *Server) ListEventEarmarks(ctx context.Context,
	r *icbt.ListEventEarmarksRequest,
) (*icbt.ListEventEarmarksResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := model.ParseEventRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad event ref-id")
	}

	event, errx := service.GetEvent(ctx, s.Db, refID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	earmarks, errx := service.GetEarmarksByEventID(ctx, s.Db, event.ID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	pbEarmarks, err := convert.ToPbListWithDb(convert.ToPbEarmark, s.Db, earmarks)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}
	response := &icbt.ListEventEarmarksResponse{
		Earmarks: pbEarmarks,
	}
	return response, nil
}

func (s *Server) ListEarmarks(ctx context.Context,
	r *icbt.ListEarmarksRequest,
) (*icbt.ListEarmarksResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	showArchived := false
	if r.Archived != nil && *r.Archived {
		showArchived = true
	}

	var paginationResult *icbt.PaginationResult
	var earmarks []*model.Earmark
	if r.Pagination != nil {
		ems, pgResult, errx := service.GetEarmarksPaginated(
			ctx, s.Db, user.ID,
			int(r.Pagination.Limit),
			int(r.Pagination.Offset),
			showArchived)
		if errx != nil {
			return nil, convert.ToTwirpError(errx)
		}
		paginationResult = convert.ToPbPagination(pgResult)
		earmarks = ems
	} else {
		var errx errs.Error
		earmarks, errx = service.GetEarmarks(
			ctx, s.Db, user.ID, showArchived)
		if errx != nil {
			return nil, convert.ToTwirpError(errx)
		}
	}

	pbEarmarks, err := convert.ToPbListWithDb(convert.ToPbEarmark, s.Db, earmarks)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}
	response := &icbt.ListEarmarksResponse{
		Earmarks:   pbEarmarks,
		Pagination: paginationResult,
	}
	return response, nil
}

func (s *Server) CreateEarmark(ctx context.Context,
	r *icbt.CreateEarmarkRequest,
) (*icbt.CreateEarmarkResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	eventItemRefID, err := model.ParseEventItemRefID(r.EventItemRefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad event-item ref-id")
	}

	eventItem, errx := service.GetEventItem(ctx, s.Db, eventItemRefID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	// make sure no earmark exists yet
	earmark, errx := service.GetEarmarkByEventItemID(ctx, s.Db, eventItem.ID)
	if errx != nil {
		switch errx.Code() {
		case errs.NotFound:
			// good. this is what we want
		default:
			slog.Error("db error", "error", errx)
			return nil, convert.ToTwirpError(errx)
		}
	} else {
		// earmark already exists!
		errStr := "already earmarked"
		if earmark.UserID != user.ID {
			errStr += " by other user"
		}
		return nil, twirp.PermissionDenied.Error(errStr)
	}

	earmark, errx = service.NewEarmark(ctx, s.Db, user, eventItem.ID, r.Note)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	pbEarmark, err := convert.ToPbEarmark(s.Db, earmark)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}

	response := &icbt.CreateEarmarkResponse{
		Earmark: pbEarmark,
	}
	return response, nil
}

func (s *Server) GetEarmarkDetails(ctx context.Context,
	r *icbt.GetEarmarkDetailsRequest,
) (*icbt.GetEarmarkDetailsResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := model.ParseEarmarkRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad earmark ref-id")
	}

	earmark, errx := service.GetEarmark(ctx, s.Db, refID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	eventItem, errx := service.GetEventItemByID(ctx, s.Db, earmark.EventItemID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	event, errx := service.GetEventByID(ctx, s.Db, eventItem.EventID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	pbEarmark, err := convert.ToPbEarmark(s.Db, earmark)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}
	response := &icbt.GetEarmarkDetailsResponse{
		Earmark:    pbEarmark,
		EventRefId: event.RefID.String(),
	}
	return response, nil
}

func (s *Server) RemoveEarmark(ctx context.Context,
	r *icbt.RemoveEarmarkRequest,
) (*icbt.RemoveEarmarkResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := model.ParseEarmarkRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad earmark ref-id")
	}

	errx := service.DeleteEarmarkByRefID(ctx, s.Db, user.ID, refID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	return &icbt.RemoveEarmarkResponse{}, nil
}
