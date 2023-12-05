package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
)

func GetApiKeyByUser(ctx context.Context, db model.PgxHandle,
	userID int,
) (*model.ApiKey, errs.Error) {
	apiKey, err := model.GetApiKeyByUser(ctx, db, userID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("user not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return apiKey, nil
}

func GetUserByApiKey(ctx context.Context, db model.PgxHandle,
	token string,
) (*model.User, errs.Error) {
	user, err := model.GetUserByApiKey(ctx, db, token)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("user not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return user, nil
}

func RotateApiKey(ctx context.Context, db model.PgxHandle,
	userID int,
) (*model.ApiKey, errs.Error) {
	apiKey, err := model.RotateApiKey(ctx, db, userID)
	if err != nil {
		return nil, errs.Internal.Errorf("db error: %w", err)
	}
	return apiKey, nil
}
