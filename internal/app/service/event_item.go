package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/dropwhile/refid/v2/reftag"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/validate"
)

var (
	ParseEventItemRefID     = reftag.Parse[model.EventItemRefID]
	EventItemRefIDMatcher   = reftag.NewMatcher[model.EventItemRefID]()
	EventItemRefIDFromBytes = reftag.FromBytes[model.EventItemRefID]
)

func (s *Service) GetEventItemsCount(
	ctx context.Context, eventIDs []int,
) ([]*model.EventItemCount, errs.Error) {
	if len(eventIDs) == 0 {
		return []*model.EventItemCount{}, nil
	}
	eventItemCounts, err := model.GetEventItemsCountByEventIDs(ctx, s.Db, eventIDs)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		eventItemCounts = []*model.EventItemCount{}
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return eventItemCounts, nil
}

func (s *Service) GetEventItemsByEvent(
	ctx context.Context, refID model.EventRefID,
) ([]*model.EventItem, errs.Error) {
	event, err := model.GetEventByRefID(ctx, s.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("event not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}

	return s.GetEventItemsByEventID(ctx, event.ID)
}

func (s *Service) GetEventItemsByEventID(
	ctx context.Context, eventID int,
) ([]*model.EventItem, errs.Error) {
	items, err := model.GetEventItemsByEvent(ctx, s.Db, eventID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		items = []*model.EventItem{}
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return items, nil
}

func (s *Service) GetEventItemsByIDs(
	ctx context.Context, eventItemIDs []int,
) ([]*model.EventItem, errs.Error) {
	if len(eventItemIDs) == 0 {
		return []*model.EventItem{}, nil
	}
	items, err := model.GetEventItemsByIDs(ctx, s.Db, eventItemIDs)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		items = []*model.EventItem{}
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return items, nil
}

func (s *Service) GetEventItem(
	ctx context.Context, eventItemRefID model.EventItemRefID,
) (*model.EventItem, errs.Error) {
	eventItem, err := model.GetEventItemByRefID(ctx, s.Db, eventItemRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("event-item not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return eventItem, nil
}

func (s *Service) GetEventItemByID(
	ctx context.Context, eventItemID int,
) (*model.EventItem, errs.Error) {
	eventItem, err := model.GetEventItemByID(ctx, s.Db, eventItemID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("event-item not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return eventItem, nil
}

func (s *Service) RemoveEventItem(
	ctx context.Context, userID int, eventItemRefID model.EventItemRefID,
	failIfChecks FailIfCheckFunc[*model.EventItem],
) errs.Error {
	eventItem, err := model.GetEventItemByRefID(ctx, s.Db, eventItemRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return errs.NotFound.Error("event-item not found")
	case err != nil:
		return errs.Internal.Error("db error")
	}

	if failIfChecks != nil && failIfChecks(eventItem) {
		return errs.FailedPrecondition.Error("extra checks failed")
	}

	event, err := model.GetEventByID(ctx, s.Db, eventItem.EventID)
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

	err = model.DeleteEventItem(ctx, s.Db, eventItem.ID)
	if err != nil {
		return errs.Internal.Error("db error")
	}
	return nil
}

func (s *Service) AddEventItem(
	ctx context.Context, userID int,
	refID model.EventRefID, description string,
) (*model.EventItem, errs.Error) {
	err := validate.Validate.VarCtx(ctx, description, "required,notblank")
	if err != nil {
		slog.
			With("field", "description").
			With("error", err).
			Info("bad field value")
		return nil, errs.InvalidArgumentError("description", "bad value")
	}

	event, err := model.GetEventByRefID(ctx, s.Db, refID)
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

	eventItem, err := model.NewEventItem(ctx, s.Db, event.ID, description)
	if err != nil {
		return nil, errs.Internal.Error("db error")
	}

	return eventItem, nil
}

func (s *Service) UpdateEventItem(
	ctx context.Context, userID int,
	refID model.EventItemRefID, description string,
	failIfChecks func(*model.EventItem) bool,
) (*model.EventItem, errs.Error) {
	err := validate.Validate.VarCtx(ctx, description, "required,notblank")
	if err != nil {
		slog.
			With("field", "description").
			With("error", err).
			Info("bad field value")
		return nil, errs.InvalidArgumentError("description", "bad value")
	}

	eventItem, err := model.GetEventItemByRefID(ctx, s.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("event-item not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}

	if failIfChecks != nil && failIfChecks(eventItem) {
		return nil, errs.FailedPrecondition.Error("extra checks failed")
	}

	event, err := model.GetEventByID(ctx, s.Db, eventItem.EventID)
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
	earmark, err := model.GetEarmarkByEventItem(ctx, s.Db, eventItem.ID)
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
	err = model.UpdateEventItem(ctx, s.Db, eventItem.ID, eventItem.Description)
	if err != nil {
		return nil, errs.Internal.Error("db error")
	}
	return eventItem, nil
}
