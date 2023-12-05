package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

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
		return nil, somerr.NotFound.Error("event not found")
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}
	return earmarks, nil
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
