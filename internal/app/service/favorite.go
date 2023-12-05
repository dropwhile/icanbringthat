package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
)

func AddFavorite(
	ctx context.Context, db model.PgxHandle, userID int,
	refID model.EventRefID,
) (*model.Event, errs.Error) {
	event, err := model.GetEventByRefID(ctx, db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("event not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}

	// can't favorite your own event
	if userID == event.UserID {
		return nil, errs.PermissionDenied.Error("can't favorite own event")
	}

	// check if favorite already exists
	_, err = model.GetFavoriteByUserEvent(ctx, db, userID, event.ID)
	if err == nil {
		return nil, errs.AlreadyExists.Error("favorite already exists")
	}

	_, err = model.CreateFavorite(ctx, db, userID, event.ID)
	if err != nil {
		return nil, errs.Internal.Error("db error")
	}

	return event, nil
}

func RemoveFavorite(
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

	favorite, err := model.GetFavoriteByUserEvent(ctx, db, userID, event.ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return errs.NotFound.Error("favorite not found")
	case err != nil:
		return errs.Internal.Error("db error")
	}

	// superfluous check, but fine to leave in
	if userID != favorite.UserID {
		return errs.PermissionDenied.Error("permission denied")
	}

	err = model.DeleteFavorite(ctx, db, favorite.ID)
	if err != nil {
		return errs.Internal.Error("db error")
	}
	return nil
}

func GetFavoriteEventsPaginated(
	ctx context.Context, db model.PgxHandle, userID int,
	limit, offset int, archived bool,
) ([]*model.Event, *Pagination, errs.Error) {
	favCount, errx := GetFavoriteEventsCount(ctx, db, userID)
	if errx != nil {
		return nil, nil, errs.Internal.Error("db error")
	}
	count := favCount.Current
	if archived {
		count = favCount.Archived
	}

	events := []*model.Event{}
	if count > 0 {
		favs, err := model.GetFavoriteEventsByUserPaginatedFiltered(
			ctx, db, userID, limit, offset, archived)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			favs = []*model.Event{}
		case err != nil:
			return nil, nil, errs.Internal.Error("db error")
		}
		events = favs
	}
	pagination := &Pagination{
		Limit:  uint32(limit),
		Offset: uint32(offset),
		Count:  uint32(count),
	}
	return events, pagination, nil
}

func GetFavoriteEventsCount(
	ctx context.Context, db model.PgxHandle, userID int,
) (*model.BifurcatedRowCounts, errs.Error) {
	favCount, err := model.GetFavoriteCountByUser(ctx, db, userID)
	if err != nil {
		return nil, errs.Internal.Error("db error")
	}
	return favCount, nil
}

func GetFavoriteEvents(
	ctx context.Context, db model.PgxHandle, userID int,
	archived bool,
) ([]*model.Event, errs.Error) {
	events, err := model.GetFavoriteEventsByUserFiltered(ctx, db, userID, archived)
	if err != nil {
		return nil, errs.Internal.Error("db error")
	}
	return events, nil
}

func GetFavoriteByUserEvent(
	ctx context.Context, db model.PgxHandle, userID int,
	eventID int,
) (*model.Favorite, errs.Error) {
	favorite, err := model.GetFavoriteByUserEvent(ctx, db, userID, eventID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("favorite not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return favorite, nil
}
