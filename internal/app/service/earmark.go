// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/dropwhile/refid/v2/reftag"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/errs"
)

var (
	EarmarkRefIDMatcher = reftag.NewMatcher[model.EarmarkRefID]()
	ParseEarmarkRefID   = reftag.Parse[model.EarmarkRefID]
)

func (s *Service) GetEarmarksByEventID(
	ctx context.Context, eventID int,
) ([]*model.Earmark, errs.Error) {
	earmarks, err := model.GetEarmarksByEvent(ctx, s.Db, eventID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return []*model.Earmark{}, nil
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return earmarks, nil
}

func (s *Service) GetEarmarkByEventItemID(
	ctx context.Context, eventItemID int,
) (*model.Earmark, errs.Error) {
	earmark, err := model.GetEarmarkByEventItem(ctx, s.Db, eventItemID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("earmark not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return earmark, nil
}

func (s *Service) GetEarmarksCount(
	ctx context.Context, userID int,
) (*model.BifurcatedRowCounts, errs.Error) {
	bifurCount, err := model.GetEarmarkCountByUser(ctx, s.Db, userID)
	if err != nil {
		return nil, errs.Internal.Error("db error")
	}
	return bifurCount, nil
}

func (s *Service) GetEarmarksPaginated(
	ctx context.Context, userID int,
	limit, offset int, archived bool,
) ([]*model.Earmark, *Pagination, errs.Error) {
	bifurCount, errx := s.GetEarmarksCount(ctx, userID)
	if errx != nil {
		slog.Error("db error", "error", errx)
		return nil, nil, errs.Internal.Error("db error")
	}
	count := bifurCount.Current
	if archived {
		count = bifurCount.Archived
	}

	earmarks := []*model.Earmark{}
	if count > 0 {
		elems, err := model.GetEarmarksByUserPaginatedFiltered(
			ctx, s.Db, userID, limit, offset, archived)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			elems = []*model.Earmark{}
		case err != nil:
			slog.Error("db error", "error", err)
			return nil, nil, errs.Internal.Error("db error")
		}
		earmarks = elems
	}
	pagination := &Pagination{
		Limit:  uint32(limit),
		Offset: uint32(offset),
		Count:  uint32(count),
	}
	return earmarks, pagination, nil
}

func (s *Service) GetEarmarks(
	ctx context.Context, userID int, archived bool,
) ([]*model.Earmark, errs.Error) {
	elems, err := model.GetEarmarksByUserFiltered(ctx, s.Db, userID, archived)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		elems = []*model.Earmark{}
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return elems, nil
}

func (s *Service) NewEarmark(
	ctx context.Context, user *model.User, eventItemID int, note string,
) (*model.Earmark, errs.Error) {
	// make sure no earmark exists yet
	em, errx := s.GetEarmarkByEventItemID(ctx, eventItemID)
	if errx != nil {
		switch errx.Code() {
		case errs.NotFound:
			// good. this is what we want
		default:
			slog.Error("db error", "error", errx)
			return nil, errs.Internal.Error("db error")
		}
	} else {
		// earmark already exists!
		errStr := "already earmarked"
		if em.UserID != user.ID {
			errStr += " by other user"
		}
		return nil, errs.AlreadyExists.Error(errStr)
	}

	// disallow earmarking archived event
	event, err := model.GetEventByEventItemID(ctx, s.Db, eventItemID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("event not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}

	if event.Archived {
		return nil, errs.PermissionDenied.Error("event is archived")
	}

	// non-owner must be verified before earmarking.
	// it is fine for owner to self-earmark though
	if !user.Verified && event.UserID != user.ID {
		return nil, errs.PermissionDenied.Error(
			"Account must be verified before earmarking is allowed.")
	}

	earmark, err := model.NewEarmark(ctx, s.Db, eventItemID, user.ID, note)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "earmark__event_item_id_key" {
				return nil, errs.AlreadyExists.Error("earmark already exists")
			}
		}
		return nil, errs.Internal.Errorf("error creating earmark: %w", err)
	}
	return earmark, nil
}

func (s *Service) GetEarmark(
	ctx context.Context, refID model.EarmarkRefID,
) (*model.Earmark, errs.Error) {
	earmark, err := model.GetEarmarkByRefID(ctx, s.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("earmark not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return earmark, nil
}

func (s *Service) DeleteEarmark(
	ctx context.Context, userID int, earmark *model.Earmark,
) errs.Error {
	if earmark.UserID != userID {
		return errs.PermissionDenied.Error("permission denied")
	}

	event, err := model.GetEventByEventItemID(ctx, s.Db, earmark.EventItemID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return errs.NotFound.Error("event not found")
	case err != nil:
		return errs.Internal.Error("db error")
	}

	if event.Archived {
		return errs.PermissionDenied.Error("event is archived")
	}

	err = model.DeleteEarmark(ctx, s.Db, earmark.ID)
	if err != nil {
		return errs.Internal.Error("db error")
	}
	return nil
}

func (s *Service) DeleteEarmarkByRefID(
	ctx context.Context, userID int, refID model.EarmarkRefID,
) errs.Error {
	earmark, errx := s.GetEarmark(ctx, refID)
	if errx != nil {
		return errx
	}

	return s.DeleteEarmark(ctx, userID, earmark)
}
