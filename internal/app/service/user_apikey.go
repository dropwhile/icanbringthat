package service

import (
	"context"
	"errors"
	"log/slog"

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
		slog.ErrorContext(ctx,
			"error getting api key by user", "error", err)
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
		slog.ErrorContext(ctx,
			"error getting user by key", "error", err)
		return nil, errs.Internal.Error("db error")
	}
	return user, nil
}

func NewApiKey(ctx context.Context, db model.PgxHandle,
	userID int,
) (*model.ApiKey, errs.Error) {
	apikey, err := model.NewApiKey(ctx, db, userID)
	if err != nil {
		slog.ErrorContext(ctx,
			"error generating new api key", "error", err)
		return nil, errs.Internal.Errorf("db error: %w", err)
	}
	return apikey, nil
}

func NewApiKeyIfNotExists(ctx context.Context, db model.PgxHandle,
	userID int,
) (*model.ApiKey, errs.Error) {
	apikey, errx := GetApiKeyByUser(ctx, db, userID)
	if errx == nil {
		return apikey, nil
	} else if errx.Code() != errs.NotFound {
		return nil, errx
	}

	apikey, err := model.NewApiKey(ctx, db, userID)
	if err != nil {
		return nil, errs.Internal.Errorf("db error: %w", err)
	}
	return apikey, nil
}
