package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
)

func GetEarmarksByEventID(
	ctx context.Context, db model.PgxHandle,
	eventID int,
) ([]*model.Earmark, errs.Error) {
	earmarks, err := model.GetEarmarksByEvent(ctx, db, eventID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return []*model.Earmark{}, nil
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return earmarks, nil
}

func GetEarmarkByEventItemID(
	ctx context.Context, db model.PgxHandle,
	eventItemID int,
) (*model.Earmark, errs.Error) {
	earmark, err := model.GetEarmarkByEventItem(ctx, db, eventItemID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("earmark not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return earmark, nil
}

func GetEarmarksCount(
	ctx context.Context, db model.PgxHandle, userID int,
) (*model.BifurcatedRowCounts, errs.Error) {
	bifurCount, err := model.GetEarmarkCountByUser(ctx, db, userID)
	if err != nil {
		return nil, errs.Internal.Error("db error")
	}
	return bifurCount, nil
}

func GetEarmarksPaginated(
	ctx context.Context, db model.PgxHandle, userID int,
	limit, offset int, archived bool,
) ([]*model.Earmark, *Pagination, errs.Error) {
	bifurCount, errx := GetEarmarksCount(ctx, db, userID)
	if errx != nil {
		return nil, nil, errs.Internal.Error("db error")
	}
	count := bifurCount.Current
	if archived {
		count = bifurCount.Archived
	}

	earmarks := []*model.Earmark{}
	if count > 0 {
		elems, err := model.GetEarmarksByUserPaginatedFiltered(
			ctx, db, userID, limit, offset, archived)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			elems = []*model.Earmark{}
		case err != nil:
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

func GetEarmarks(
	ctx context.Context, db model.PgxHandle, userID int,
	archived bool,
) ([]*model.Earmark, errs.Error) {
	elems, err := model.GetEarmarksByUserFiltered(ctx, db, userID, archived)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		elems = []*model.Earmark{}
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return elems, nil
}

func NewEarmark(ctx context.Context, db model.PgxHandle,
	eventItemID, userID int, note string,
) (*model.Earmark, errs.Error) {
	earmark, err := model.NewEarmark(ctx, db, eventItemID, userID, note)
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

func GetEarmark(ctx context.Context, db model.PgxHandle,
	refID model.EarmarkRefID,
) (*model.Earmark, errs.Error) {
	earmark, err := model.GetEarmarkByRefID(ctx, db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("earmark not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return earmark, nil
}

func DeleteEarmark(ctx context.Context, db model.PgxHandle, userID int,
	earmark *model.Earmark,
) errs.Error {
	event, err := model.GetEventByEarmarkID(ctx, db, earmark.ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return errs.NotFound.Error("event not found")
	case err != nil:
		return errs.Internal.Error("db error")
	}

	if event.Archived {
		return errs.PermissionDenied.Error("event is archived")
	}

	if earmark.UserID != userID {
		return errs.PermissionDenied.Error("permission denied")
	}

	err = model.DeleteEarmark(ctx, db, earmark.ID)
	if err != nil {
		return errs.Internal.Error("db error")
	}
	return nil
}

func DeleteEarmarkByRefID(ctx context.Context, db model.PgxHandle, userID int,
	refID model.EarmarkRefID,
) errs.Error {
	earmark, errx := GetEarmark(ctx, db, refID)
	if errx != nil {
		return errx
	}
	return DeleteEarmark(ctx, db, userID, earmark)
}
