package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/somerr"
)

func GetEarmarksByEventID(
	ctx context.Context, db model.PgxHandle, userID int,
	eventID int,
) ([]*model.Earmark, somerr.Error) {
	earmarks, err := model.GetEarmarksByEvent(ctx, db, eventID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return []*model.Earmark{}, nil
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}
	return earmarks, nil
}

func GetEarmarkByEventItemID(
	ctx context.Context, db model.PgxHandle, userID int,
	eventItemID int,
) (*model.Earmark, somerr.Error) {
	earmark, err := model.GetEarmarkByEventItem(ctx, db, eventItemID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, somerr.NotFound.Error("earmark not found")
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}
	return earmark, nil
}

func GetEarmarksCount(
	ctx context.Context, db model.PgxHandle, userID int,
) (*model.BifurcatedRowCounts, somerr.Error) {
	bifurCount, err := model.GetEarmarkCountByUser(ctx, db, userID)
	if err != nil {
		return nil, somerr.Internal.Error("db error")
	}
	return bifurCount, nil
}

func GetEarmarksPaginated(
	ctx context.Context, db model.PgxHandle, userID int,
	limit, offset int, archived bool,
) ([]*model.Earmark, *Pagination, somerr.Error) {
	bifurCount, errx := GetEarmarksCount(ctx, db, userID)
	if errx != nil {
		return nil, nil, somerr.Internal.Error("db error")
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
			return nil, nil, somerr.Internal.Error("db error")
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
) ([]*model.Earmark, somerr.Error) {
	elems, err := model.GetEarmarksByUserFiltered(ctx, db, userID, archived)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		elems = []*model.Earmark{}
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}
	return elems, nil
}

func NewEarmark(ctx context.Context, db model.PgxHandle,
	eventItemID, userID int, note string,
) (*model.Earmark, somerr.Error) {
	earmark, err := model.NewEarmark(ctx, db, eventItemID, userID, note)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "earmark__event_item_id_key" {
				return nil, somerr.AlreadyExists.Error("earmark already exists")
			}
		}
		return nil, somerr.Internal.Errorf("error creating earmark: %w", err)
	}
	return earmark, nil
}

func GetEarmark(ctx context.Context, db model.PgxHandle,
	refID model.EarmarkRefID,
) (*model.Earmark, somerr.Error) {
	earmark, err := model.GetEarmarkByRefID(ctx, db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, somerr.NotFound.Error("earmark not found")
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}
	return earmark, nil
}

func DeleteEarmark(ctx context.Context, db model.PgxHandle,
	userID int, earmark *model.Earmark,
) somerr.Error {
	eventItem, err := model.GetEventItemByID(ctx, db, earmark.EventItemID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return somerr.NotFound.Error("event-item not found")
	case err != nil:
		return somerr.Internal.Error("db error")
	}

	event, err := model.GetEventByID(ctx, db, eventItem.EventID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return somerr.NotFound.Error("event not found")
	case err != nil:
		return somerr.Internal.Error("db error")
	}

	if event.Archived {
		return somerr.PermissionDenied.Error("event is archived")
	}

	if earmark.UserID != userID {
		return somerr.PermissionDenied.Error("permission denied")
	}

	err = model.DeleteEarmark(ctx, db, earmark.ID)
	if err != nil {
		return somerr.Internal.Error("db error")
	}
	return nil
}

func DeleteEarmarkByRefID(ctx context.Context, db model.PgxHandle,
	userID int, refID model.EarmarkRefID,
) somerr.Error {
	earmark, errx := GetEarmark(ctx, db, refID)
	if errx != nil {
		return errx
	}

	eventItem, err := model.GetEventItemByID(ctx, db, earmark.EventItemID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return somerr.NotFound.Error("event-item not found")
	case err != nil:
		return somerr.Internal.Error("db error")
	}

	event, err := model.GetEventByID(ctx, db, eventItem.EventID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return somerr.NotFound.Error("event not found")
	case err != nil:
		return somerr.Internal.Error("db error")
	}

	if event.Archived {
		return somerr.PermissionDenied.Error("event is archived")
	}

	if earmark.UserID != userID {
		return somerr.PermissionDenied.Error("permission denied")
	}

	err = model.DeleteEarmark(ctx, db, earmark.ID)
	if err != nil {
		return somerr.Internal.Error("db error")
	}
	return nil
}
