package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/dropwhile/refid/v2/reftag"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
)

var (
	ApiKeyRefIDMatcher   = reftag.NewMatcher[model.ApiKeyRefID]()
	ApiKeyRefIDFromBytes = reftag.FromBytes[model.ApiKeyRefID]
	ParseApiKeyRefID     = reftag.Parse[model.ApiKeyRefID]
)

func (s *Service) GetApiKeyByUser(
	ctx context.Context, userID int,
) (*model.ApiKey, errs.Error) {
	apiKey, err := model.GetApiKeyByUser(ctx, s.Db, userID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("user-api-key not found")
	case err != nil:
		slog.ErrorContext(ctx,
			"error getting api key by user", "error", err)
		return nil, errs.Internal.Error("db error")
	}
	return apiKey, nil
}

func (s *Service) GetUserByApiKey(
	ctx context.Context, token string,
) (*model.User, errs.Error) {
	user, err := model.GetUserByApiKey(ctx, s.Db, token)
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

func (s *Service) NewApiKey(
	ctx context.Context, userID int,
) (*model.ApiKey, errs.Error) {
	apikey, err := model.NewApiKey(ctx, s.Db, userID)
	if err != nil {
		slog.ErrorContext(ctx,
			"error generating new api key", "error", err)
		return nil, errs.Internal.Errorf("db error: %w", err)
	}
	return apikey, nil
}

func (s *Service) NewApiKeyIfNotExists(
	ctx context.Context, userID int,
) (*model.ApiKey, errs.Error) {
	apikey, errx := s.GetApiKeyByUser(ctx, userID)
	if errx == nil {
		return apikey, nil
	} else if errx.Code() != errs.NotFound {
		return nil, errx
	}

	apikey, err := model.NewApiKey(ctx, s.Db, userID)
	if err != nil {
		return nil, errs.Internal.Errorf("db error: %w", err)
	}
	return apikey, nil
}
