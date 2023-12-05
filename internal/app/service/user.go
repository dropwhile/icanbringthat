package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/somerr"
)

func GetUsersByIDs(
	ctx context.Context, db model.PgxHandle, userID int,
	userIDs []int,
) ([]*model.User, somerr.Error) {
	if len(userIDs) == 0 {
		return []*model.User{}, nil
	}

	users, err := model.GetUsersByIDs(ctx, db, userIDs)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		users = []*model.User{}
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}
	return users, nil
}
