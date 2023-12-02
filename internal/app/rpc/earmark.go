package rpc

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/rpc/dto"
	pb "github.com/dropwhile/icbt/rpc"
)

func (s *Server) ListEventEarmarks(ctx context.Context,
	r *pb.ListEventEarmarksRequest,
) (*pb.ListEventEarmarksResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := model.ParseEventRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad event ref-id")
	}

	event, err := model.GetEventByRefID(ctx, s.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.NotFoundError("event not found")
	case err != nil:
		return nil, twirp.InternalError("db error")
	}

	earmarks, err := model.GetEarmarksByEvent(ctx, s.Db, event.ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.NotFoundError("event not found")
	case err != nil:
		return nil, twirp.InternalError("db error")
	}

	pbEarmarks, err := dto.ToPbListWithDb(dto.ToPbEarmark, s.Db, earmarks)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}
	response := &pb.ListEventEarmarksResponse{
		Earmarks: pbEarmarks,
	}
	return response, nil
}

func (s *Server) ListEarmarks(ctx context.Context,
	r *pb.ListEarmarksRequest,
) (*pb.ListEarmarksResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	showArchived := false
	if r.Archived != nil && *r.Archived {
		showArchived = true
	}

	var paginationResult *pb.PaginationResult
	var earmarks []*model.Earmark
	if r.Pagination != nil {
		limit := int(r.Pagination.Limit)
		offset := int(r.Pagination.Offset)

		earmarkCounts, err := model.GetEarmarkCountByUser(ctx, s.Db, user.ID)
		if err != nil {
			return nil, twirp.InternalError("db error")
		}

		count := earmarkCounts.Current
		if showArchived {
			count = earmarkCounts.Archived
		}

		if count > 0 {
			earmarks, err = model.GetEarmarksByUserPaginatedFiltered(
				ctx, s.Db, user.ID, limit, offset, showArchived)
			switch {
			case errors.Is(err, pgx.ErrNoRows):
				earmarks = []*model.Earmark{}
			case err != nil:
				return nil, twirp.InternalError("db error")
			}
		}
		paginationResult = &pb.PaginationResult{
			Limit:  uint32(limit),
			Offset: uint32(offset),
			Count:  uint32(count),
		}
	} else {
		earmarks, err = model.GetEarmarksByUserFiltered(
			ctx, s.Db, user.ID, showArchived)
		if err != nil {
			return nil, twirp.InternalError("db error")
		}

	}

	pbEarmarks, err := dto.ToPbListWithDb(dto.ToPbEarmark, s.Db, earmarks)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}
	response := &pb.ListEarmarksResponse{
		Earmarks:   pbEarmarks,
		Pagination: paginationResult,
	}
	return response, nil
}

func (s *Server) CreateEarmark(ctx context.Context,
	r *pb.CreateEarmarkRequest,
) (*pb.CreateEarmarkResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	eventItemRefID, err := model.ParseEventItemRefID(r.EventItemRefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad event-item ref-id")
	}

	eventItem, err := model.GetEventItemByRefID(ctx, s.Db, eventItemRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.NotFoundError("event-item not found")
	case err != nil:
		return nil, twirp.InternalError("db error")
	}

	// make sure no earmark exists yet
	earmark, err := model.GetEarmarkByEventItem(ctx, s.Db, eventItem.ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		// good. this is what we want
	case err == nil:
		// earmark already exists!
		errStr := "already earmarked"
		if earmark.UserID != user.ID {
			errStr += " by other user"
		}
		return nil, twirp.PermissionDenied.Error(errStr)
	default:
		return nil, twirp.InternalError("db error")
	}

	earmark, err = model.NewEarmark(ctx, s.Db, eventItem.ID, user.ID, r.Note)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}

	pbEarmark, err := dto.ToPbEarmark(s.Db, earmark)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}

	response := &pb.CreateEarmarkResponse{
		Earmark: pbEarmark,
	}
	return response, nil
}

func (s *Server) GetEarmarkDetails(ctx context.Context,
	r *pb.GetEarmarkDetailsRequest,
) (*pb.GetEarmarkDetailsResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := model.ParseEarmarkRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad earmark ref-id")
	}

	earmark, err := model.GetEarmarkByRefID(ctx, s.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.NotFoundError("earmark not found")
	case err != nil:
		return nil, twirp.InternalError("db error")
	}

	eventItem, err := model.GetEventItemByID(ctx, s.Db, earmark.EventItemID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.NotFoundError("earmark not found")
	case err != nil:
		return nil, twirp.InternalError("db error")
	}

	event, err := model.GetEventByID(ctx, s.Db, eventItem.EventID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.NotFoundError("earmark not found")
	case err != nil:
		return nil, twirp.InternalError("db error")
	}

	pbEarmark, err := dto.ToPbEarmark(s.Db, earmark)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}
	response := &pb.GetEarmarkDetailsResponse{
		Earmark:    pbEarmark,
		EventRefId: event.RefID.String(),
	}
	return response, nil
}

func (s *Server) RemoveEarmark(ctx context.Context,
	r *pb.RemoveEarmarkRequest,
) (*pb.RemoveEarmarkResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := model.ParseEarmarkRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad earmark ref-id")
	}

	earmark, err := model.GetEarmarkByRefID(ctx, s.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.NotFoundError("earmark not found")
	case err != nil:
		return nil, twirp.InternalError("db error")
	}

	if earmark.UserID != user.ID {
		return nil, twirp.PermissionDenied.Error("permission denied")
	}

	err = model.DeleteEarmark(ctx, s.Db, earmark.ID)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}

	return &pb.RemoveEarmarkResponse{}, nil
}
