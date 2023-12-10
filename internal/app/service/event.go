package service

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/jackc/pgx/v5"

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

func UpdateEvent(
	ctx context.Context, db model.PgxHandle, userID int,
	refID model.EventRefID, name *string, description *string,
	startTime *time.Time, tz *string,
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

	if name == nil && description == nil && startTime == nil && tz == nil {
		return nil, errs.InvalidArgument.Error("missing fields")
	}

	if event.Archived {
		return nil, errs.PermissionDenied.Error("event is archived")
	}

	changes := false
	if name != nil && *name != event.Name {
		if *name == "" {
			return nil, errs.InvalidArgumentError("name", "cannot be empty")
		}
		event.Name = *name
		changes = true
	}
	if description != nil && *description != event.Description {
		if *description == "" {
			return nil, errs.InvalidArgumentError("description", "cannot be empty")
		}
		event.Description = *description
		changes = true
	}
	if startTime != nil && *startTime != event.StartTime {
		if (*startTime).IsZero() {
			return nil, errs.InvalidArgumentError("start_time", "bad value")
		}
		if (*startTime).Before(time.Now().UTC()) {
			return nil, errs.InvalidArgumentError("start_time", "cannot be in the past")
		}
		event.StartTime = *startTime
		changes = true
	}
	if tz != nil {
		if *tz == "" {
			return nil, errs.InvalidArgumentError("tz", "cannot be empty")
		}
		loc, err := time.LoadLocation(*tz)
		if err != nil {
			return nil, errs.InvalidArgumentError("tz", "unrecognized timezone")
		}
		if loc != event.StartTimeTz.Location {
			event.StartTimeTz = &model.TimeZone{Location: loc}
			changes = true
		}
	}

	if !changes {
		return nil, errs.FailedPrecondition.Error("no changes")
	}

	if err := model.UpdateEvent(
		ctx, db, event.ID,
		event.Name, event.Description, event.ItemSortOrder,
		event.StartTime, event.StartTimeTz,
	); err != nil {
		return nil, errs.Internal.Error("db error")
	}
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
		event.Name, event.Description, event.ItemSortOrder,
		event.StartTime, event.StartTimeTz,
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
	if name == "" {
		return nil, errs.InvalidArgumentError("name", "bad empty value")
	}
	if description == "" {
		return nil, errs.InvalidArgumentError("description", "bad empty value")
	}
	if when.IsZero() {
		return nil, errs.InvalidArgumentError("start_time", "bad empty value")
	}
	if tz == "" {
		return nil, errs.InvalidArgumentError("tz", "bad empty value")
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
