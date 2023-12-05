package service

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/somerr"
)

func GetEvent(
	ctx context.Context, db model.PgxHandle, userID int,
	refID model.EventRefID,
) (*model.Event, somerr.Error) {
	event, err := model.GetEventByRefID(ctx, db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, somerr.NotFound.Error("event not found")
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}
	return event, nil
}

func GetEventByID(
	ctx context.Context, db model.PgxHandle, userID int,
	ID int,
) (*model.Event, somerr.Error) {
	event, err := model.GetEventByID(ctx, db, ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, somerr.NotFound.Error("event not found")
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}
	return event, nil
}

func GetEventsByIDs(ctx context.Context, db model.PgxHandle,
	eventIDs []int,
) ([]*model.Event, somerr.Error) {
	elems, err := model.GetEventsByIDs(ctx, db, eventIDs)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return []*model.Event{}, nil
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}
	return elems, nil
}

func DeleteEvent(
	ctx context.Context, db model.PgxHandle, userID int,
	refID model.EventRefID,
) somerr.Error {
	event, err := model.GetEventByRefID(ctx, db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return somerr.NotFound.Error("event not found")
	case err != nil:
		return somerr.Internal.Error("db error")
	}

	if userID != event.UserID {
		return somerr.PermissionDenied.Error("permission denied")
	}

	err = model.DeleteEvent(ctx, db, event.ID)
	if err != nil {
		return somerr.Internal.Error("db error")
	}
	return nil
}

func UpdateEvent(
	ctx context.Context, db model.PgxHandle, userID int,
	refID model.EventRefID, name *string, description *string,
	startTime *time.Time, tz *string,
) (*model.Event, somerr.Error) {
	event, err := model.GetEventByRefID(ctx, db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, somerr.NotFound.Error("event not found")
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}

	if userID != event.UserID {
		return nil, somerr.PermissionDenied.Error("permission denied")
	}

	if name == nil && description == nil && startTime == nil && tz == nil {
		return nil, somerr.InvalidArgument.Error("missing fields")
	}

	if event.Archived {
		return nil, somerr.PermissionDenied.Error("event is archived")
	}

	changes := false
	if name != nil && *name != event.Name {
		if *name == "" {
			return nil, somerr.InvalidArgumentError("name", "cannot be empty")
		}
		event.Name = *name
		changes = true
	}
	if description != nil && *description != event.Description {
		if *description == "" {
			return nil, somerr.InvalidArgumentError("description", "cannot be empty")
		}
		event.Description = *description
		changes = true
	}
	if startTime != nil && *startTime != event.StartTime {
		if (*startTime).IsZero() {
			return nil, somerr.InvalidArgumentError("start_time", "bad value")
		}
		if (*startTime).Before(time.Now().UTC()) {
			return nil, somerr.InvalidArgumentError("start_time", "cannot be in the past")
		}
		event.StartTime = *startTime
		changes = true
	}
	if tz != nil {
		if *tz == "" {
			return nil, somerr.InvalidArgumentError("tz", "cannot be empty")
		}
		loc, err := time.LoadLocation(*tz)
		if err != nil {
			return nil, somerr.InvalidArgumentError("tz", "unrecognized timezone")
		}
		if loc != event.StartTimeTz.Location {
			event.StartTimeTz = &model.TimeZone{Location: loc}
			changes = true
		}
	}

	if !changes {
		return nil, somerr.FailedPrecondition.Error("no changes")
	}

	if err := model.UpdateEvent(
		ctx, db, event.ID,
		event.Name, event.Description, event.ItemSortOrder,
		event.StartTime, event.StartTimeTz,
	); err != nil {
		return nil, somerr.Internal.Error("db error")
	}
	return event, nil
}

func UpdateEventItemSorting(
	ctx context.Context, db model.PgxHandle, userID int,
	refID model.EventRefID, itemSortOrder []int,
) (*model.Event, somerr.Error) {
	event, err := model.GetEventByRefID(ctx, db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, somerr.NotFound.Error("event not found")
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}

	if userID != event.UserID {
		return nil, somerr.PermissionDenied.Error("permission denied")
	}

	if event.Archived {
		return nil, somerr.PermissionDenied.Error("event is archived")
	}

	if reflect.DeepEqual(event.ItemSortOrder, itemSortOrder) {
		return nil, somerr.FailedPrecondition.Error("no changes")
	}

	event.ItemSortOrder = itemSortOrder

	if err := model.UpdateEvent(
		ctx, db, event.ID,
		event.Name, event.Description, event.ItemSortOrder,
		event.StartTime, event.StartTimeTz,
	); err != nil {
		return nil, somerr.Internal.Error("db error")
	}
	return event, nil
}

func CreateEvent(
	ctx context.Context, db model.PgxHandle, user *model.User,
	name string, description string,
	when time.Time, tz string,
) (*model.Event, somerr.Error) {
	if !user.Verified {
		return nil, somerr.PermissionDenied.Error(
			"Account must be verified before event creation is allowed.")
	}
	if name == "" {
		return nil, somerr.InvalidArgumentError("name", "bad empty value")
	}
	if description == "" {
		return nil, somerr.InvalidArgumentError("description", "bad empty value")
	}
	if when.IsZero() {
		return nil, somerr.InvalidArgumentError("start_time", "bad empty value")
	}
	if tz == "" {
		return nil, somerr.InvalidArgumentError("tz", "bad empty value")
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, somerr.InvalidArgumentError("tz", "unrecognized timezone")
	}

	event, err := model.NewEvent(ctx, db, user.ID,
		name, description, when, &model.TimeZone{Location: loc})
	if err != nil {
		return nil, somerr.Internal.Error("db error")
	}
	return event, nil
}

func GetEventsPaginated(
	ctx context.Context, db model.PgxHandle, userID int,
	limit, offset int, archived bool,
) ([]*model.Event, *Pagination, somerr.Error) {
	eventCount, errx := GetEventsCount(ctx, db, userID)
	if errx != nil {
		return nil, nil, somerr.Internal.Error("db error")
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
			return nil, nil, somerr.Internal.Error("db error")
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
) ([]*model.Event, *Pagination, somerr.Error) {
	eventCount, errx := GetEventsCount(ctx, db, userID)
	if errx != nil {
		return nil, nil, somerr.Internal.Error("db error")
	}

	events := []*model.Event{}
	if eventCount.Current > 0 {
		evts, err := model.GetEventsComingSoonByUserPaginated(ctx, db, userID, limit, offset)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			evts = []*model.Event{}
		case err != nil:
			return nil, nil, somerr.Internal.Error("db error")
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
) (*model.BifurcatedRowCounts, somerr.Error) {
	count, err := model.GetEventCountsByUser(ctx, db, userID)
	if err != nil {
		return nil, somerr.Internal.Error("db error")
	}
	return count, nil
}

func GetEvents(
	ctx context.Context, db model.PgxHandle, userID int,
	archived bool,
) ([]*model.Event, somerr.Error) {
	elems, err := model.GetEventsByUserFiltered(ctx, db, userID, archived)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		elems = []*model.Event{}
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}
	return elems, nil
}
