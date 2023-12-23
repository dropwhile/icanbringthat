package service

import (
	"context"
	"errors"
	"log/slog"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/samber/mo"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
)

func GetEvent(
	ctx context.Context, db model.PgxHandle,
	refID model.EventRefID,
) (*model.Event, errs.Error) {
	event, err := model.GetEventByRefID(ctx, db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("event not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return event, nil
}

func GetEventByID(
	ctx context.Context, db model.PgxHandle,
	ID int,
) (*model.Event, errs.Error) {
	event, err := model.GetEventByID(ctx, db, ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("event not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return event, nil
}

func GetEventsByIDs(ctx context.Context, db model.PgxHandle,
	eventIDs []int,
) ([]*model.Event, errs.Error) {
	if len(eventIDs) == 0 {
		return []*model.Event{}, nil
	}
	elems, err := model.GetEventsByIDs(ctx, db, eventIDs)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return []*model.Event{}, nil
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return elems, nil
}

func DeleteEvent(
	ctx context.Context, db model.PgxHandle, userID int,
	refID model.EventRefID,
) errs.Error {
	event, err := model.GetEventByRefID(ctx, db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return errs.NotFound.Error("event not found")
	case err != nil:
		return errs.Internal.Error("db error")
	}

	if userID != event.UserID {
		return errs.PermissionDenied.Error("permission denied")
	}

	err = model.DeleteEvent(ctx, db, event.ID)
	if err != nil {
		return errs.Internal.Error("db error")
	}
	return nil
}

type EventUpdateValues struct {
	Name          mo.Option[string]    `validate:"omitempty,notblank"`
	Description   mo.Option[string]    `validate:"omitempty,notblank"`
	ItemSortOrder mo.Option[[]int]     `validate:"omitempty,gt=0"`
	StartTime     mo.Option[time.Time] `validate:"omitempty"`
	Tz            mo.Option[string]    `validate:"omitempty,timezone"`
}

func UpdateEvent(
	ctx context.Context, db model.PgxHandle, userID int,
	refID model.EventRefID, euvs *EventUpdateValues,
) (*model.Event, errs.Error) {
	// if no values, error
	if euvs.Name.IsAbsent() &&
		euvs.Description.IsAbsent() &&
		euvs.ItemSortOrder.IsAbsent() &&
		euvs.StartTime.IsAbsent() &&
		euvs.Tz.IsAbsent() {
		return nil, errs.InvalidArgument.Error("missing fields")
	}

	err := validate.StructCtx(ctx, euvs)
	if err != nil {
		badField := err.(validator.ValidationErrors)[0].Field()
		slog.
			With("field", badField).
			With("error", err).
			Info("bad field value")
		return nil, errs.InvalidArgumentError(badField, "bad value")
	}

	if val, ok := euvs.StartTime.Get(); ok {
		if val.IsZero() {
			return nil, errs.InvalidArgumentError("start_time", "bad value")
		}
		if val.Before(time.Now().UTC().Add(-30 * time.Minute)) {
			return nil, errs.InvalidArgumentError("start_time", "cannot be in the past")
		}
	}

	var loc *model.TimeZone
	var maybeLoc mo.Option[*model.TimeZone]
	if val, ok := euvs.Tz.Get(); ok {
		loc, err = model.ParseTimeZone(val)
		if err != nil {
			return nil, errs.InvalidArgumentError("tz", "unrecognized timezone")
		}
		maybeLoc = mo.Some(loc)
	}

	// get event
	event, err := model.GetEventByRefID(ctx, db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("event not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}

	// check general condition requirements
	if userID != event.UserID {
		return nil, errs.PermissionDenied.Error("permission denied")
	}

	if event.Archived {
		return nil, errs.PermissionDenied.Error("event is archived")
	}

	// do update
	if err := model.UpdateEvent(
		ctx, db, event.ID,
		euvs.Name, euvs.Description, euvs.ItemSortOrder,
		euvs.StartTime, maybeLoc,
	); err != nil {
		return nil, errs.Internal.Error("db error")
	}

	event.Name = euvs.Name.OrElse(event.Name)
	event.Description = euvs.Description.OrElse(event.Description)
	event.ItemSortOrder = euvs.ItemSortOrder.OrElse(event.ItemSortOrder)
	event.StartTime = euvs.StartTime.OrElse(event.StartTime)
	event.StartTimeTz = maybeLoc.OrElse(event.StartTimeTz)
	return event, nil
}

func UpdateEventItemSorting(
	ctx context.Context, db model.PgxHandle, userID int,
	refID model.EventRefID, itemSortOrder []int,
) (*model.Event, errs.Error) {
	event, err := model.GetEventByRefID(ctx, db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("event not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}

	if userID != event.UserID {
		return nil, errs.PermissionDenied.Error("permission denied")
	}

	if event.Archived {
		return nil, errs.PermissionDenied.Error("event is archived")
	}

	if reflect.DeepEqual(event.ItemSortOrder, itemSortOrder) {
		return nil, errs.FailedPrecondition.Error("no changes")
	}

	event.ItemSortOrder = itemSortOrder

	if err := model.UpdateEvent(
		ctx, db, event.ID,
		mo.None[string](), mo.None[string](),
		mo.Some(event.ItemSortOrder),
		mo.None[time.Time](),
		mo.None[*model.TimeZone](),
	); err != nil {
		return nil, errs.Internal.Error("db error")
	}
	return event, nil
}

func CreateEvent(
	ctx context.Context, db model.PgxHandle, user *model.User,
	name string, description string,
	when time.Time, tz string,
) (*model.Event, errs.Error) {
	if !user.Verified {
		return nil, errs.PermissionDenied.Error(
			"Account must be verified before event creation is allowed.")
	}

	err := validate.VarCtx(ctx, name, "required,notblank")
	if err != nil {
		slog.
			With("field", "name").
			With("error", err).
			Info("bad field value")
		return nil, errs.InvalidArgumentError("name", "bad value")
	}

	err = validate.VarCtx(ctx, description, "required,notblank")
	if err != nil {
		slog.
			With("field", "description").
			With("error", err).
			Info("bad field value")
		return nil, errs.InvalidArgumentError("description", "bad value")
	}

	// check for zero time
	if when.IsZero() {
		return nil, errs.InvalidArgumentError("start_time", "bad empty value")
	}
	// check for unix epoch
	if when.UTC().Equal(time.Unix(0, 0).UTC()) {
		return nil, errs.InvalidArgumentError("start_time", "bad value")
	}

	err = validate.VarCtx(ctx, tz, "required,timezone")
	if err != nil {
		slog.
			With("field", "tz").
			With("error", err).
			Info("bad field value")
		return nil, errs.InvalidArgumentError("tz", "bad value")
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, errs.InvalidArgumentError("tz", "unrecognized timezone")
	}

	event, err := model.NewEvent(ctx, db, user.ID,
		name, description, when, &model.TimeZone{Location: loc})
	if err != nil {
		return nil, errs.Internal.Error("db error")
	}
	return event, nil
}

func GetEventsPaginated(
	ctx context.Context, db model.PgxHandle, userID int,
	limit, offset int, archived bool,
) ([]*model.Event, *Pagination, errs.Error) {
	eventCount, errx := GetEventsCount(ctx, db, userID)
	if errx != nil {
		return nil, nil, errs.Internal.Error("db error")
	}
	count := eventCount.Current
	if archived {
		count = eventCount.Archived
	}

	events := []*model.Event{}
	if count > 0 {
		evts, err := model.GetEventsByUserPaginatedFiltered(
			ctx, db, userID, limit, offset, archived)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			evts = []*model.Event{}
		case err != nil:
			return nil, nil, errs.Internal.Error("db error")
		}
		events = evts
	}
	pagination := &Pagination{
		Limit:  uint32(limit),
		Offset: uint32(offset),
		Count:  uint32(count),
	}
	return events, pagination, nil
}

func GetEventsComingSoonPaginated(
	ctx context.Context, db model.PgxHandle, userID int,
	limit, offset int,
) ([]*model.Event, *Pagination, errs.Error) {
	eventCount, errx := GetEventsCount(ctx, db, userID)
	if errx != nil {
		return nil, nil, errs.Internal.Error("db error")
	}

	events := []*model.Event{}
	if eventCount.Current > 0 {
		evts, err := model.GetEventsComingSoonByUserPaginated(ctx, db, userID, limit, offset)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			evts = []*model.Event{}
		case err != nil:
			return nil, nil, errs.Internal.Error("db error")
		}
		events = evts
	}
	pagination := &Pagination{
		Limit:  uint32(limit),
		Offset: uint32(offset),
		Count:  uint32(eventCount.Current),
	}
	return events, pagination, nil
}

func GetEventsCount(
	ctx context.Context, db model.PgxHandle, userID int,
) (*model.BifurcatedRowCounts, errs.Error) {
	count, err := model.GetEventCountsByUser(ctx, db, userID)
	if err != nil {
		return nil, errs.Internal.Error("db error")
	}
	return count, nil
}

func GetEvents(
	ctx context.Context, db model.PgxHandle, userID int,
	archived bool,
) ([]*model.Event, errs.Error) {
	elems, err := model.GetEventsByUserFiltered(ctx, db, userID, archived)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		elems = []*model.Event{}
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return elems, nil
}
