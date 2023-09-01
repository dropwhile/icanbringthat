// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/errs"
)

func (s *Service) AddFavorite(
	ctx context.Context, userID int, refID model.EventRefID,
) (*model.Event, errs.Error) {
	event, err := model.GetEventByRefID(ctx, s.Db, refID)
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
	_, err = model.GetFavoriteByUserEvent(ctx, s.Db, userID, event.ID)
	if err == nil {
		return nil, errs.AlreadyExists.Error("favorite already exists")
	}

	_, err = model.CreateFavorite(ctx, s.Db, userID, event.ID)
	if err != nil {
		slog.Error("db error", "error", err)
		return nil, errs.Internal.Error("db error")
	}

	return event, nil
}

func (s *Service) RemoveFavorite(
	ctx context.Context, userID int, refID model.EventRefID,
) errs.Error {
	event, err := model.GetEventByRefID(ctx, s.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return errs.NotFound.Error("event not found")
	case err != nil:
		return errs.Internal.Error("db error")
	}

	favorite, err := model.GetFavoriteByUserEvent(ctx, s.Db, userID, event.ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return errs.NotFound.Error("favorite not found")
	case err != nil:
		return errs.Internal.Error("db error")
	}

	// no need to double check (for permissions) userID == favorite.UserID here,
	// since we fetch favorite by user+event above

	err = model.DeleteFavorite(ctx, s.Db, favorite.ID)
	if err != nil {
		return errs.Internal.Error("db error")
	}
	return nil
}

func (s *Service) GetFavoriteEventsPaginated(
	ctx context.Context, userID int,
	limit, offset int, archived bool,
) ([]*model.Event, *Pagination, errs.Error) {
	favCount, errx := s.GetFavoriteEventsCount(ctx, userID)
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
			ctx, s.Db, userID, limit, offset, archived)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			favs = []*model.Event{}
		case err != nil:
			slog.Error("db error", "error", err)
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

func (s *Service) GetFavoriteEventsCount(
	ctx context.Context, userID int,
) (*model.BifurcatedRowCounts, errs.Error) {
	favCount, err := model.GetFavoriteCountByUser(ctx, s.Db, userID)
	if err != nil {
		return nil, errs.Internal.Error("db error")
	}
	return favCount, nil
}

func (s *Service) GetFavoriteEvents(
	ctx context.Context, userID int, archived bool,
) ([]*model.Event, errs.Error) {
	events, err := model.GetFavoriteEventsByUserFiltered(ctx, s.Db, userID, archived)
	if err != nil {
		return nil, errs.Internal.Error("db error")
	}
	return events, nil
}

func (s *Service) GetFavoriteByUserEvent(
	ctx context.Context, userID int, eventID int,
) (*model.Favorite, errs.Error) {
	favorite, err := model.GetFavoriteByUserEvent(ctx, s.Db, userID, eventID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("favorite not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return favorite, nil
}
