package rpc

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/app/convert"
	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/errs"
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

		earmarkCounts, errx := service.GetEarmarksCount(ctx, s.Db, user.ID)
		if errx != nil {
			return nil, convert.ToTwirpError(errx)
		}

		count := earmarkCounts.Current
		if showArchived {
			count = earmarkCounts.Archived
		}

		if count > 0 {
			earmarks, _, errx = service.GetEarmarksPaginated(
				ctx, s.Db, user.ID, limit, offset, showArchived)
			if errx != nil {
				return nil, convert.ToTwirpError(errx)
			}
		}
		paginationResult = &pb.PaginationResult{
			Limit:  uint32(limit),
			Offset: uint32(offset),
			Count:  uint32(count),
		}
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

	earmark, errx = service.NewEarmark(ctx, s.Db, eventItem.ID, user.ID, r.Note)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	pbEarmark, err := convert.ToPbEarmark(s.Db, earmark)
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

	errx := service.DeleteEarmarkByRefID(ctx, s.Db, user.ID, refID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	return &pb.RemoveEarmarkResponse{}, nil
}
