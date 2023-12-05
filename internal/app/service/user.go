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

func GetUser(
	ctx context.Context, db model.PgxHandle, userID int,
	refID model.UserRefID,
) (*model.User, somerr.Error) {
	user, err := model.GetUserByRefID(ctx, db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, somerr.NotFound.Error("user not found")
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}
	return user, nil
}

func GetUserByEmail(
	ctx context.Context, db model.PgxHandle,
	email string,
) (*model.User, somerr.Error) {
	user, err := model.GetUserByEmail(ctx, db, email)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, somerr.NotFound.Error("user not found")
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}
	return user, nil
}

func UpdateUserSettings(
	ctx context.Context, db model.PgxHandle, userID int,
	pm *model.UserSettings,
) somerr.Error {
	err := model.UpdateUserSettings(ctx, db, pm, userID)
	if err != nil {
		return somerr.Internal.Error("db error")
	}
	return nil
}
