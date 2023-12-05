package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
)

func GetUserPWResetByRefID(ctx context.Context, db model.PgxHandle,
	refID model.UserPWResetRefID,
) (*model.UserPWReset, errs.Error) {
	pwreset, err := model.GetUserPWResetByRefID(ctx, db, refID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, errs.NotFound.Error("pwreset not found")
		default:
			return nil, errs.Internal.Error("db error")
		}
	}
	return pwreset, nil
}

func NewUserPWReset(ctx context.Context, db model.PgxHandle,
	userID int,
) (*model.UserPWReset, errs.Error) {
	pwreset, err := model.NewUserPWReset(ctx, db, userID)
	if err != nil {
		return nil, errs.Internal.Error("db error")
	}
	return pwreset, nil
}

func UpdateUserPWReset(ctx context.Context, db model.PgxHandle,
	user *model.User, upw *model.UserPWReset,
) errs.Error {
	errx := TxnFunc(ctx, db, func(tx pgx.Tx) error {
		innerErr := model.UpdateUser(ctx, tx,
			user.Email, user.Name, user.PWHash,
			user.Verified, user.PWAuth, user.ApiAccess,
			user.WebAuthn, user.ID,
		)
		if innerErr != nil {
			log.Debug().
				Err(innerErr).
				Msg("inner db error saving user")
			return innerErr
		}

		innerErr = model.DeleteUserPWReset(ctx, tx, upw.RefID)
		if innerErr != nil {
			log.Debug().
				Err(innerErr).
				Msg("inner db error cleaning up pw reset token")
			return innerErr
		}
		return nil
	})
	return errx
}
