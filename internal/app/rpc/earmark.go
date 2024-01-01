package rpc

import (
	"context"

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

	refID, err := service.ParseEventRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad event ref-id")
	}

	event, errx := s.svc.GetEvent(ctx, refID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	earmarks, errx := s.svc.GetEarmarksByEventID(ctx, event.ID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	pbEarmarks, err := convert.ToPbListWithService(convert.ToPbEarmark, s.svc, earmarks)
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
		ems, pgResult, errx := s.svc.GetEarmarksPaginated(
			ctx, user.ID,
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
		earmarks, errx = s.svc.GetEarmarks(
			ctx, user.ID, showArchived)
		if errx != nil {
			return nil, convert.ToTwirpError(errx)
		}
	}

	pbEarmarks, err := convert.ToPbListWithService(convert.ToPbEarmark, s.svc, earmarks)
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

	eventItemRefID, err := service.ParseEventItemRefID(r.EventItemRefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad event-item ref-id")
	}

	eventItem, errx := s.svc.GetEventItem(ctx, eventItemRefID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	earmark, errx := s.svc.NewEarmark(ctx, user, eventItem.ID, r.Note)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	pbEarmark, err := convert.ToPbEarmark(s.svc, earmark)
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

	refID, err := service.ParseEarmarkRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad earmark ref-id")
	}

	earmark, errx := s.svc.GetEarmark(ctx, refID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	eventItem, errx := s.svc.GetEventItemByID(ctx, earmark.EventItemID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	event, errx := s.svc.GetEventByID(ctx, eventItem.EventID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	pbEarmark, err := convert.ToPbEarmark(s.svc, earmark)
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

	refID, err := service.ParseEarmarkRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad earmark ref-id")
	}

	errx := s.svc.DeleteEarmarkByRefID(ctx, user.ID, refID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	return &icbt.RemoveEarmarkResponse{}, nil
}
