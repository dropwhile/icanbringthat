package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
)

func GetEventItemsCount(
	ctx context.Context, db model.PgxHandle,
	eventIDs []int,
) ([]*model.EventItemCount, errs.Error) {
	if len(eventIDs) == 0 {
		return []*model.EventItemCount{}, nil
	}
	eventItemCounts, err := model.GetEventItemsCountByEventIDs(ctx, db, eventIDs)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		eventItemCounts = []*model.EventItemCount{}
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return eventItemCounts, nil
}

func GetEventItemsByEvent(
	ctx context.Context, db model.PgxHandle,
	refID model.EventRefID,
) ([]*model.EventItem, errs.Error) {
	event, err := model.GetEventByRefID(ctx, db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("event not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}

	return GetEventItemsByEventID(ctx, db, event.ID)
}

func GetEventItemsByEventID(
	ctx context.Context, db model.PgxHandle,
	eventID int,
) ([]*model.EventItem, errs.Error) {
	items, err := model.GetEventItemsByEvent(ctx, db, eventID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		items = []*model.EventItem{}
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return items, nil
}

func GetEventItemsByIDs(ctx context.Context, db model.PgxHandle,
	eventItemIDs []int,
) ([]*model.EventItem, errs.Error) {
	if len(eventItemIDs) == 0 {
		return []*model.EventItem{}, nil
	}
	items, err := model.GetEventItemsByIDs(ctx, db, eventItemIDs)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		items = []*model.EventItem{}
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return items, nil
}

func GetEventItem(
	ctx context.Context, db model.PgxHandle,
	eventItemRefID model.EventItemRefID,
) (*model.EventItem, errs.Error) {
	eventItem, err := model.GetEventItemByRefID(ctx, db, eventItemRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("event-item not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return eventItem, nil
}

func GetEventItemByID(
	ctx context.Context, db model.PgxHandle,
	eventItemID int,
) (*model.EventItem, errs.Error) {
	eventItem, err := model.GetEventItemByID(ctx, db, eventItemID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("event-item not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return eventItem, nil
}

func RemoveEventItem(
	ctx context.Context, db model.PgxHandle, userID int,
	eventItemRefID model.EventItemRefID,
	failIfChecks func(*model.EventItem) bool,
) errs.Error {
	eventItem, err := model.GetEventItemByRefID(ctx, db, eventItemRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return errs.NotFound.Error("event-item not found")
	case err != nil:
		return errs.Internal.Error("db error")
	}

	if failIfChecks != nil && failIfChecks(eventItem) {
		return errs.FailedPrecondition.Error("extra checks failed")
	}

	event, err := model.GetEventByID(ctx, db, eventItem.EventID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return errs.NotFound.Error("event not found")
	case err != nil:
		return errs.Internal.Error("db error")
	}

	if event.UserID != userID {
		return errs.PermissionDenied.Error("not event owner")
	}

	if event.Archived {
		return errs.PermissionDenied.Error("event is archived")
	}

	// this shouldn't ever be the case, but its a cheap check
	if eventItem.EventID != event.ID {
		return errs.NotFound.Error("event-item not found")
	}

	err = model.DeleteEventItem(ctx, db, eventItem.ID)
	if err != nil {
		return errs.Internal.Error("db error")
	}
	return nil
}

func AddEventItem(
	ctx context.Context, db model.PgxHandle, userID int,
	refID model.EventRefID, description string,
) (*model.EventItem, errs.Error) {
	event, err := model.GetEventByRefID(ctx, db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("event not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}

	if event.UserID != userID {
		return nil, errs.PermissionDenied.Error("not event owner")
	}

	if event.Archived {
		return nil, errs.PermissionDenied.Error("event is archived")
	}

	eventItem, err := model.NewEventItem(ctx, db, event.ID, description)
	if err != nil {
		return nil, errs.Internal.Error("db error")
	}

	return eventItem, nil
}

func UpdateEventItem(
	ctx context.Context, db model.PgxHandle, userID int,
	refID model.EventItemRefID, description string,
	failIfChecks func(*model.EventItem) bool,
) (*model.EventItem, errs.Error) {
	eventItem, err := model.GetEventItemByRefID(ctx, db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("event-item not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}

	if failIfChecks != nil && failIfChecks(eventItem) {
		return nil, errs.FailedPrecondition.Error("extra checks failed")
	}

	event, err := model.GetEventByID(ctx, db, eventItem.EventID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("event not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}

	if event.UserID != userID {
		return nil, errs.PermissionDenied.Error("not event owner")
	}

	if event.Archived {
		return nil, errs.PermissionDenied.Error("event is archived")
	}

	// check if earmark exists, and is marked by someone else
	// if so, disallow editing in that case.
	earmark, err := model.GetEarmarkByEventItem(ctx, db, eventItem.ID)
	switch {
	case err != nil && !errors.Is(err, pgx.ErrNoRows):
		return nil, errs.Internal.Error("db error")
	case err == nil && earmark.UserID != userID:
		slog.InfoContext(ctx, "user id mismatch",
			slog.Int("user.ID", userID),
			slog.Int("earmark.UserID", earmark.UserID),
		)
		return nil, errs.PermissionDenied.Error("earmarked by other user")
	}

	eventItem.Description = description
	err = model.UpdateEventItem(ctx, db, eventItem.ID, eventItem.Description)
	if err != nil {
		return nil, errs.Internal.Error("db error")
	}
	return eventItem, nil
}
