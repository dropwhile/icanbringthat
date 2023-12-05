package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/somerr"
)

func GetEventItemsCount(
	ctx context.Context, db model.PgxHandle,
	eventIDs []int,
) ([]*model.EventItemCount, somerr.Error) {
	eventItemCounts, err := model.GetEventItemsCountByEventIDs(ctx, db, eventIDs)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("no rows for event items")
		eventItemCounts = []*model.EventItemCount{}
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}
	return eventItemCounts, nil
}

func GetEventItemsByEvent(
	ctx context.Context, db model.PgxHandle,
	refID model.EventRefID,
) ([]*model.EventItem, somerr.Error) {
	event, err := model.GetEventByRefID(ctx, db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, somerr.NotFound.Error("event not found")
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}

	return GetEventItemsByEventID(ctx, db, event.ID)
}

func GetEventItemsByEventID(
	ctx context.Context, db model.PgxHandle,
	eventID int,
) ([]*model.EventItem, somerr.Error) {
	items, err := model.GetEventItemsByEvent(ctx, db, eventID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		items = []*model.EventItem{}
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}
	return items, nil
}

func GetEventItemsByIDs(ctx context.Context, db model.PgxHandle,
	eventItemIDs []int,
) ([]*model.EventItem, somerr.Error) {
	items, err := model.GetEventItemsByIDs(ctx, db, eventItemIDs)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		items = []*model.EventItem{}
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}
	return items, nil
}

func GetEventItem(
	ctx context.Context, db model.PgxHandle,
	eventItemRefID model.EventItemRefID,
) (*model.EventItem, somerr.Error) {
	eventItem, err := model.GetEventItemByRefID(ctx, db, eventItemRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, somerr.NotFound.Error("event-item not found")
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}
	return eventItem, nil
}

func GetEventItemByID(
	ctx context.Context, db model.PgxHandle,
	eventItemID int,
) (*model.EventItem, somerr.Error) {
	eventItem, err := model.GetEventItemByID(ctx, db, eventItemID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, somerr.NotFound.Error("event-item not found")
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}
	return eventItem, nil
}

func RemoveEventItem(
	ctx context.Context, db model.PgxHandle, userID int,
	eventItemRefID model.EventItemRefID,
	failIfChecks func(*model.EventItem) bool,
) somerr.Error {
	eventItem, err := model.GetEventItemByRefID(ctx, db, eventItemRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return somerr.NotFound.Error("event-item not found")
	case err != nil:
		return somerr.Internal.Error("db error")
	}

	if failIfChecks != nil && failIfChecks(eventItem) {
		return somerr.FailedPrecondition.Error("extra checks failed")
	}

	event, err := model.GetEventByID(ctx, db, eventItem.EventID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return somerr.NotFound.Error("event not found")
	case err != nil:
		return somerr.Internal.Error("db error")
	}

	if event.UserID != userID {
		return somerr.PermissionDenied.Error("not event owner")
	}

	if event.Archived {
		return somerr.PermissionDenied.Error("event is archived")
	}

	// this shouldn't ever be the case, but its a cheap check
	if eventItem.EventID != event.ID {
		return somerr.NotFound.Error("event-item not found")
	}

	err = model.DeleteEventItem(ctx, db, eventItem.ID)
	if err != nil {
		return somerr.Internal.Error("db error")
	}
	return nil
}

func AddEventItem(
	ctx context.Context, db model.PgxHandle, userID int,
	refID model.EventRefID, description string,
) (*model.EventItem, somerr.Error) {
	event, err := model.GetEventByRefID(ctx, db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, somerr.NotFound.Error("event not found")
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}

	if event.UserID != userID {
		return nil, somerr.PermissionDenied.Error("not event owner")
	}

	if event.Archived {
		return nil, somerr.PermissionDenied.Error("event is archived")
	}

	eventItem, err := model.NewEventItem(ctx, db, event.ID, description)
	if err != nil {
		return nil, somerr.Internal.Error("db error")
	}

	return eventItem, nil
}

func UpdateEventItem(
	ctx context.Context, db model.PgxHandle, userID int,
	refID model.EventItemRefID, description string,
	failIfChecks func(*model.EventItem) bool,
) (*model.EventItem, somerr.Error) {
	eventItem, err := model.GetEventItemByRefID(ctx, db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, somerr.NotFound.Error("event-item not found")
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}

	log.Debug().Msg("here1")
	if failIfChecks != nil && failIfChecks(eventItem) {
		return nil, somerr.FailedPrecondition.Error("extra checks failed")
	}

	log.Debug().Msg("here")
	event, err := model.GetEventByID(ctx, db, eventItem.EventID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, somerr.NotFound.Error("event not found")
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}

	log.Debug().Msg("here-after")
	if event.UserID != userID {
		return nil, somerr.PermissionDenied.Error("not event owner")
	}

	if event.Archived {
		return nil, somerr.PermissionDenied.Error("event is archived")
	}

	// check if earmark exists, and is marked by someone else
	// if so, disallow editing in that case.
	earmark, err := model.GetEarmarkByEventItem(ctx, db, eventItem.ID)
	switch {
	case err != nil && !errors.Is(err, pgx.ErrNoRows):
		return nil, somerr.Internal.Error("db error")
	case err == nil && earmark.UserID != userID:
		log.Info().
			Int("user.ID", userID).
			Int("earmark.UserID", earmark.UserID).
			Msg("user id mismatch")
		return nil, somerr.PermissionDenied.Error("earmarked by other user")
	}

	eventItem.Description = description
	err = model.UpdateEventItem(ctx, db, eventItem.ID, eventItem.Description)
	if err != nil {
		return nil, somerr.Internal.Error("db error")
	}
	return eventItem, nil
}
