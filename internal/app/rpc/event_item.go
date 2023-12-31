package rpc

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/app/convert"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/middleware/auth"
	"github.com/dropwhile/icbt/rpc/icbt"
)

func (s *Server) ListEventItems(ctx context.Context,
	r *icbt.ListEventItemsRequest,
) (*icbt.ListEventItemsResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := service.ParseEventRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad event ref-id")
	}

	items, errx := s.service.GetEventItemsByEvent(ctx, refID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	response := &icbt.ListEventItemsResponse{
		Items: convert.ToPbList(convert.ToPbEventItem, items),
	}
	return response, nil
}

func (s *Server) RemoveEventItem(ctx context.Context,
	r *icbt.RemoveEventItemRequest,
) (*icbt.RemoveEventItemResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := service.ParseEventItemRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad event-item ref-id")
	}

	errx := s.service.RemoveEventItem(ctx, user.ID, refID, nil)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	return &icbt.RemoveEventItemResponse{}, nil
}

func (s *Server) AddEventItem(ctx context.Context,
	r *icbt.AddEventItemRequest,
) (*icbt.AddEventItemResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := service.ParseEventRefID(r.EventRefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad event ref-id")
	}

	eventItem, errx := s.service.AddEventItem(
		ctx, user.ID, refID, r.Description,
	)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	response := &icbt.AddEventItemResponse{
		EventItem: convert.ToPbEventItem(eventItem),
	}
	return response, nil
}

func (s *Server) UpdateEventItem(ctx context.Context,
	r *icbt.UpdateEventItemRequest,
) (*icbt.UpdateEventItemResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := service.ParseEventItemRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad event-item ref-id")
	}

	eventItem, errx := s.service.UpdateEventItem(
		ctx, user.ID, refID, r.Description, nil,
	)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	response := &icbt.UpdateEventItemResponse{
		EventItem: convert.ToPbEventItem(eventItem),
	}
	return response, nil
}
