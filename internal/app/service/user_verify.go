package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
)

func GetUserVerifyByRefID(ctx context.Context, db model.PgxHandle,
	refID model.UserVerifyRefID,
) (*model.UserVerify, errs.Error) {
	verify, err := model.GetUserVerifyByRefID(ctx, db, refID)
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

func NewUserVerify(ctx context.Context, db model.PgxHandle,
	user *model.User,
) (*model.UserVerify, errs.Error) {
	verify, err := model.NewUserVerify(ctx, db, user)
	if err != nil {
		return nil, errs.Internal.Error("db error")
	}
	return verify, nil
}

func SetUserVerified(ctx context.Context, db model.PgxHandle,
	user *model.User, verifier *model.UserVerify,
) errs.Error {
	errx := TxnFunc(ctx, db, func(tx pgx.Tx) error {
		innerErr := model.UpdateUser(ctx, tx,
			user.Email, user.Name, user.PWHash,
			user.Verified, user.PWAuth, user.ApiAccess,
			user.WebAuthn, user.ID,
		)
		if innerErr != nil {
			log.Debug().Err(innerErr).Msg("inner db error saving user")
			return innerErr
		}

		innerErr = model.DeleteUserVerify(ctx, tx, verifier.RefID)
		if innerErr != nil {
			log.Debug().Err(innerErr).Msg("inner db error cleaning up verifier token")
			return innerErr
		}
		return nil
	})
	return errx
}
