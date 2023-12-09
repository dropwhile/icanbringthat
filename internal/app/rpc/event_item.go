package rpc

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/app/convert"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/middleware/auth"
	pb "github.com/dropwhile/icbt/rpc"
)

func (s *Server) ListEventItems(ctx context.Context,
	r *pb.ListEventItemsRequest,
) (*pb.ListEventItemsResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := model.ParseEventRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad event ref-id")
	}

	items, errx := service.GetEventItemsByEvent(ctx, s.Db, refID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	response := &pb.ListEventItemsResponse{
		Items: convert.ToPbList(convert.ToPbEventItem, items),
	}
	return response, nil
}

func (s *Server) RemoveEventItem(ctx context.Context,
	r *pb.RemoveEventItemRequest,
) (*pb.RemoveEventItemResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := model.ParseEventItemRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad earmark ref-id")
	}

	errx := service.RemoveEventItem(ctx, s.Db, user.ID, refID, nil)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	return &pb.RemoveEventItemResponse{}, nil
}

func (s *Server) AddEventItem(ctx context.Context,
	r *pb.AddEventItemRequest,
) (*pb.AddEventItemResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := model.ParseEventRefID(r.EventRefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad event ref-id")
	}

	eventItem, errx := service.AddEventItem(
		ctx, s.Db, user.ID, refID, r.Description,
	)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	response := &pb.AddEventItemResponse{
		EventItem: convert.ToPbEventItem(eventItem),
	}
	return response, nil
}

func (s *Server) UpdateEventItem(ctx context.Context,
	r *pb.UpdateEventItemRequest,
) (*pb.UpdateEventItemResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := model.ParseEventItemRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad event-item ref-id")
	}

	eventItem, errx := service.UpdateEventItem(
		ctx, s.Db, user.ID, refID, r.Description, nil,
	)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	response := &pb.UpdateEventItemResponse{
		EventItem: convert.ToPbEventItem(eventItem),
	}
	return response, nil
}
