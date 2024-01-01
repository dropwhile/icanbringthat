package service

import (
	"context"
	"errors"
	"log/slog"
	"reflect"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/samber/mo"

	"github.com/dropwhile/refid/v2/reftag"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/validate"
)

var (
	EventRefIDMatcher   = reftag.NewMatcher[model.EventRefID]()
	EventRefIDFromBytes = reftag.FromBytes[model.EventRefID]
	ParseEventRefID     = reftag.Parse[model.EventRefID]
)

func (s *Service) GetEvent(
	ctx context.Context, refID model.EventRefID,
) (*model.Event, errs.Error) {
	event, err := model.GetEventByRefID(ctx, s.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("event not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return event, nil
}

func (s *Service) GetEventByID(
	ctx context.Context, ID int,
) (*model.Event, errs.Error) {
	event, err := model.GetEventByID(ctx, s.Db, ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("event not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return event, nil
}

func (s *Service) GetEventsByIDs(
	ctx context.Context, eventIDs []int,
) ([]*model.Event, errs.Error) {
	if len(eventIDs) == 0 {
		return []*model.Event{}, nil
	}
	elems, err := model.GetEventsByIDs(ctx, s.Db, eventIDs)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return []*model.Event{}, nil
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return elems, nil
}

func (s *Service) DeleteEvent(
	ctx context.Context, userID int, refID model.EventRefID,
) errs.Error {
	event, err := model.GetEventByRefID(ctx, s.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return errs.NotFound.Error("event not found")
	case err != nil:
		return errs.Internal.Error("db error")
	}

	if userID != event.UserID {
		return errs.PermissionDenied.Error("permission denied")
	}

	err = model.DeleteEvent(ctx, s.Db, event.ID)
	if err != nil {
		return errs.Internal.Error("db error")
	}
	return nil
}

type EventUpdateValues struct {
	StartTime     mo.Option[time.Time] `validate:"omitempty"`
	Name          mo.Option[string]    `validate:"omitempty,notblank"`
	Description   mo.Option[string]    `validate:"omitempty,notblank"`
	Tz            mo.Option[string]    `validate:"omitempty,timezone"`
	ItemSortOrder mo.Option[[]int]     `validate:"omitempty,gt=0"`
}

func (s *Service) UpdateEvent(
	ctx context.Context, userID int,
	refID model.EventRefID, euvs *EventUpdateValues,
) errs.Error {
	// if no values, error
	if euvs.Name.IsAbsent() &&
		euvs.Description.IsAbsent() &&
		euvs.ItemSortOrder.IsAbsent() &&
		euvs.StartTime.IsAbsent() &&
		euvs.Tz.IsAbsent() {
		return errs.InvalidArgument.Error("missing fields")
	}

	err := validate.Validate.StructCtx(ctx, euvs)
	if err != nil {
		badField := validate.GetErrorField(err)
		slog.
			With("field", badField).
			With("error", err).
			Info("bad field value")
		return errs.InvalidArgumentError(badField, "bad value")
	}

	if val, ok := euvs.StartTime.Get(); ok {
		if val.IsZero() {
			return errs.InvalidArgumentError("start_time", "bad value")
		}
		/*
			if val.Before(time.Now().UTC().Add(-30 * time.Minute)) {
				return nil, errs.InvalidArgumentError("start_time", "cannot be in the past")
			}
		*/
	}

	var loc *model.TimeZone
	var maybeLoc mo.Option[*model.TimeZone]
	if val, ok := euvs.Tz.Get(); ok {
		loc, err = ParseTimeZone(val)
		if err != nil {
			return errs.InvalidArgumentError("tz", "unrecognized timezone")
		}
		maybeLoc = mo.Some(loc)
	}

	// get event
	event, err := model.GetEventByRefID(ctx, s.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return errs.NotFound.Error("event not found")
	case err != nil:
		return errs.Internal.Error("db error")
	}

	// check general condition requirements
	if userID != event.UserID {
		return errs.PermissionDenied.Error("permission denied")
	}

	if event.Archived {
		return errs.PermissionDenied.Error("event is archived")
	}

	// do update
	err = model.UpdateEvent(ctx, s.Db, event.ID, &model.EventUpdateModelValues{
		Name:          euvs.Name,
		Description:   euvs.Description,
		ItemSortOrder: euvs.ItemSortOrder,
		StartTime:     euvs.StartTime,
		Tz:            maybeLoc,
	})
	if err != nil {
		slog.With("error", err).Error("db error")
		return errs.Internal.Error("db error")
	}
	return nil
}

func (s *Service) UpdateEventItemSorting(
	ctx context.Context, userID int,
	refID model.EventRefID, itemSortOrder []int,
) (*model.Event, errs.Error) {
	event, err := model.GetEventByRefID(ctx, s.Db, refID)
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
		ctx, s.Db, event.ID, &model.EventUpdateModelValues{
			ItemSortOrder: mo.Some(event.ItemSortOrder),
		},
	); err != nil {
		return nil, errs.Internal.Error("db error")
	}
	return event, nil
}

func (s *Service) CreateEvent(
	ctx context.Context, user *model.User,
	name string, description string,
	when time.Time, tz string,
) (*model.Event, errs.Error) {
	if !user.Verified {
		return nil, errs.PermissionDenied.Error(
			"Account must be verified before event creation is allowed.")
	}

	err := validate.Validate.VarCtx(ctx, name, "required,notblank")
	if err != nil {
		slog.
			With("field", "name").
			With("error", err).
			Info("bad field value")
		return nil, errs.InvalidArgumentError("name", "bad value")
	}

	err = validate.Validate.VarCtx(ctx, description, "required,notblank")
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

	err = validate.Validate.VarCtx(ctx, tz, "required,timezone")
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

	event, err := model.NewEvent(ctx, s.Db, user.ID,
		name, description, when, &model.TimeZone{Location: loc})
	if err != nil {
		return nil, errs.Internal.Error("db error")
	}
	return event, nil
}

func (s *Service) GetEventsPaginated(
	ctx context.Context, userID int,
	limit, offset int, archived bool,
) ([]*model.Event, *Pagination, errs.Error) {
	eventCount, errx := s.GetEventsCount(ctx, userID)
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
			ctx, s.Db, userID, limit, offset, archived)
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

func (s *Service) GetEventsComingSoonPaginated(
	ctx context.Context, userID int,
	limit, offset int,
) ([]*model.Event, *Pagination, errs.Error) {
	eventCount, errx := s.GetEventsCount(ctx, userID)
	if errx != nil {
		return nil, nil, errs.Internal.Error("db error")
	}

	events := []*model.Event{}
	if eventCount.Current > 0 {
		evts, err := model.GetEventsComingSoonByUserPaginated(ctx, s.Db, userID, limit, offset)
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

func (s *Service) GetEventsCount(
	ctx context.Context, userID int,
) (*model.BifurcatedRowCounts, errs.Error) {
	count, err := model.GetEventCountsByUser(ctx, s.Db, userID)
	if err != nil {
		return nil, errs.Internal.Error("db error")
	}
	return count, nil
}

func (s *Service) GetEvents(
	ctx context.Context, userID int,
	archived bool,
) ([]*model.Event, errs.Error) {
	elems, err := model.GetEventsByUserFiltered(ctx, s.Db, userID, archived)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		elems = []*model.Event{}
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return elems, nil
}
