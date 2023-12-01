package rpc

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/rpc/dto"
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

	event, err := model.GetEventByRefID(ctx, s.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.NotFoundError("event not found")
	case err != nil:
		return nil, twirp.InternalError("db error")
	}

	items, err := model.GetEventItemsByEvent(ctx, s.Db, event.ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.NotFoundError("event not found")
	case err != nil:
		return nil, twirp.InternalError("db error")
	}

	response := &pb.ListEventItemsResponse{
		Items: dto.ToPbList(dto.ToPbEventItem, items),
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

	eventItem, err := model.GetEventItemByRefID(ctx, s.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.NotFoundError("event-item not found")
	case err != nil:
		return nil, twirp.InternalError("db error")
	}

	event, err := model.GetEventByID(ctx, s.Db, eventItem.EventID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.NotFoundError("event not found")
	case err != nil:
		return nil, twirp.InternalError("db error")
	}

	if event.UserID != user.ID {
		return nil, twirp.PermissionDenied.Error("not event owner")
	}

	if event.Archived {
		return nil, twirp.PermissionDenied.Error("event is archived")
	}

	err = model.DeleteEventItem(ctx, s.Db, eventItem.ID)
	if err != nil {
		return nil, twirp.InternalError("db error")
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

	event, err := model.GetEventByRefID(ctx, s.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.NotFoundError("event not found")
	case err != nil:
		return nil, twirp.InternalError("db error")
	}

	if event.UserID != user.ID {
		return nil, twirp.PermissionDenied.Error("not event owner")
	}

	if event.Archived {
		return nil, twirp.PermissionDenied.Error("event is archived")
	}

	eventItem, err := model.NewEventItem(ctx, s.Db, event.ID, r.Description)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}

	response := &pb.AddEventItemResponse{
		EventItem: dto.ToPbEventItem(eventItem),
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
		return nil, twirp.InvalidArgumentError("ref_id", "bad earmark ref-id")
	}

	eventItem, err := model.GetEventItemByRefID(ctx, s.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.NotFoundError("event-item not found")
	case err != nil:
		return nil, twirp.InternalError("db error")
	}

	event, err := model.GetEventByID(ctx, s.Db, eventItem.EventID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.NotFoundError("event not found")
	case err != nil:
		return nil, twirp.InternalError("db error")
	}

	if event.UserID != user.ID {
		return nil, twirp.PermissionDenied.Error("not event owner")
	}

	if event.Archived {
		return nil, twirp.PermissionDenied.Error("event is archived")
	}

	// check if earmark exists, and is marked by someone else
	// if so, disallow editing in that case.
	earmark, err := model.GetEarmarkByEventItem(ctx, s.Db, eventItem.ID)
	switch {
	case err != nil && !errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.InternalError("db error")
	case err == nil && earmark.UserID != user.ID:
		log.Info().
			Int("user.ID", user.ID).
			Int("earmark.UserID", earmark.UserID).
			Msg("user id mismatch")
		return nil, twirp.PermissionDenied.Error("earmarked by other user")
	}

	eventItem.Description = r.Description
	err = model.UpdateEventItem(ctx, s.Db, eventItem.ID, eventItem.Description)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}

	response := &pb.UpdateEventItemResponse{
		EventItem: dto.ToPbEventItem(eventItem),
	}
	return response, nil
}
