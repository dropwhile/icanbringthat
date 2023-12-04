package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/someerr"
)

func GetEventItemsCount(
	ctx context.Context, db model.PgxHandle, userID int,
	eventIDs []int,
) ([]*model.EventItemCount, someerr.Error) {
	eventItemCounts, err := model.GetEventItemsCountByEventIDs(ctx, db, eventIDs)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("no rows for event items")
		eventItemCounts = []*model.EventItemCount{}
	case err != nil:
		return nil, someerr.Internal.Error("db error")
	}
	return eventItemCounts, nil
}
