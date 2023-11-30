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
		earmarks = []*model.Earmark{}
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
