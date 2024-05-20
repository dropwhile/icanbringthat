// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rpc

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icanbringthat/internal/app/convert"
	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
	"github.com/dropwhile/icanbringthat/rpc/icbt"
)

func (s *Server) EventListEarmarks(ctx context.Context,
	r *icbt.EventListEarmarksRequest,
) (*icbt.EventListEarmarksResponse, error) {
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

	pbEarmarks, err := convert.ToPbListWithService(ctx, convert.ToPbEarmark, s.svc, earmarks)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}
	response := &icbt.EventListEarmarksResponse{
		Earmarks: pbEarmarks,
	}
	return response, nil
}

func (s *Server) EarmarksList(ctx context.Context,
	r *icbt.EarmarksListRequest,
) (*icbt.EarmarksListResponse, error) {
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

	pbEarmarks, err := convert.ToPbListWithService(ctx, convert.ToPbEarmark, s.svc, earmarks)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}
	response := &icbt.EarmarksListResponse{
		Earmarks:   pbEarmarks,
		Pagination: paginationResult,
	}
	return response, nil
}

func (s *Server) EarmarkCreate(ctx context.Context,
	r *icbt.EarmarkCreateRequest,
) (*icbt.EarmarkCreateResponse, error) {
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

	pbEarmark, err := convert.ToPbEarmark(ctx, s.svc, earmark)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}

	response := &icbt.EarmarkCreateResponse{
		Earmark: pbEarmark,
	}
	return response, nil
}

func (s *Server) EarmarkGetDetails(ctx context.Context,
	r *icbt.EarmarkGetDetailsRequest,
) (*icbt.EarmarkGetDetailsResponse, error) {
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

	pbEarmark, err := convert.ToPbEarmark(ctx, s.svc, earmark)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}
	response := &icbt.EarmarkGetDetailsResponse{
		Earmark:    pbEarmark,
		EventRefId: event.RefID.String(),
	}
	return response, nil
}

func (s *Server) EarmarkRemove(ctx context.Context,
	r *icbt.EarmarkRemoveRequest,
) (*icbt.Empty, error) {
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

	return &icbt.Empty{}, nil
}
