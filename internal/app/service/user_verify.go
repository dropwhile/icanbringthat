package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/dropwhile/refid/v2/reftag"
	"github.com/jackc/pgx/v5"
	"github.com/samber/mo"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/logger"
)

var (
	UserVerifyRefIDMatcher   = reftag.NewMatcher[model.UserVerifyRefID]()
	UserVerifyRefIDFromBytes = reftag.FromBytes[model.UserVerifyRefID]
	ParseUserVerifyRefID     = reftag.Parse[model.UserVerifyRefID]
)

func (s *Service) GetUserVerifyByRefID(
	ctx context.Context, refID model.UserVerifyRefID,
) (*model.UserVerify, errs.Error) {
	verify, err := model.GetUserVerifyByRefID(ctx, s.Db, refID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, errs.NotFound.Error("verify not found")
		default:
			return nil, errs.Internal.Error("db error")
		}
	}
	return verify, nil
}

func (s *Service) NewUserVerify(
	ctx context.Context, userID int,
) (*model.UserVerify, errs.Error) {
	verify, err := model.NewUserVerify(ctx, s.Db, userID)
	if err != nil {
		return nil, errs.Internal.Error("db error")
	}
	return verify, nil
}

func (s *Service) SetUserVerified(
	ctx context.Context, user *model.User, verifier *model.UserVerify,
) errs.Error {
	errx := TxnFunc(ctx, s.Db, func(tx pgx.Tx) error {
		innerErr := model.UpdateUser(ctx, tx, user.ID,
			&model.UserUpdateModelValues{
				Verified: mo.Some(user.Verified),
			},
		)
		if innerErr != nil {
			slog.DebugContext(ctx, "inner db error saving user",
				logger.Err(innerErr))
			return innerErr
		}

		innerErr = model.DeleteUserVerify(ctx, tx, verifier.RefID)
		if innerErr != nil {
			slog.DebugContext(ctx, "inner db error cleaning up verifier token",
				logger.Err(innerErr))
			return innerErr
		}
		return nil
	})
	return errx
}
