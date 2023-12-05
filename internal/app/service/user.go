package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

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
	ctx context.Context, db model.PgxHandle,
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

func GetUserByID(
	ctx context.Context, db model.PgxHandle,
	ID int,
) (*model.User, somerr.Error) {
	user, err := model.GetUserByID(ctx, db, ID)
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

func GetApiKeyByUser(ctx context.Context, db model.PgxHandle,
	userID int,
) (*model.ApiKey, somerr.Error) {
	apiKey, err := model.GetApiKeyByUser(ctx, db, userID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, somerr.NotFound.Error("user not found")
	case err != nil:
		return nil, somerr.Internal.Error("db error")
	}
	return apiKey, nil
}

func RotateApiKey(ctx context.Context, db model.PgxHandle,
	userID int,
) (*model.ApiKey, somerr.Error) {
	apiKey, err := model.RotateApiKey(ctx, db, userID)
	if err != nil {
		return nil, somerr.Internal.Errorf("db error: %w", err)
	}
	return apiKey, nil
}

func NewUser(ctx context.Context, db model.PgxHandle,
	email, name string, rawPass []byte,
) (*model.User, somerr.Error) {
	user, err := model.NewUser(ctx, db, email, name, rawPass)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "user_email_idx" {
				return nil, somerr.AlreadyExists.Error("user already exists")
			}
		}
		return nil, somerr.Internal.Errorf("error creating user: %w", err)
	}
	return user, nil
}

func UpdateUser(ctx context.Context, db model.PgxHandle,
	email, name string, pwHash []byte, verified bool,
	pwAuth, apiAccess, webAuthn bool, userID int,
) somerr.Error {
	err := model.UpdateUser(ctx, db,
		email, name, pwHash, verified,
		pwAuth, apiAccess, webAuthn, userID,
	)
	if err != nil {
		return somerr.Internal.Errorf("db error: %w", err)
	}
	return nil
}

func DeleteUser(ctx context.Context, db model.PgxHandle,
	userID int,
) somerr.Error {
	err := model.DeleteUser(ctx, db, userID)
	if err != nil {
		return somerr.Internal.Errorf("db error: %w", err)
	}
	return nil
}
