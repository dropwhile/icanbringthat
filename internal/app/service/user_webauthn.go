package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/somerr"
)

func GetUserCredentialByRefID(ctx context.Context, db model.PgxHandle,
	refID model.CredentialRefID,
) (*model.UserCredential, somerr.Error) {
	cred, err := model.GetUserCredentialByRefID(ctx, db, refID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, somerr.NotFound.Error("credential")
		default:
			return nil, somerr.Internal.Errorf("db error: %w", err)
		}
	}
	return cred, nil
}

func GetUserCredentialsByUser(ctx context.Context, db model.PgxHandle, userID int,
) ([]*model.UserCredential, somerr.Error) {
	creds, err := model.GetUserCredentialsByUser(ctx, db, userID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return []*model.UserCredential{}, nil
		default:
			return nil, somerr.Internal.Errorf("db error: %w", err)
		}
	}
	return creds, nil
}

func GetUserCredentialCountByUser(ctx context.Context, db model.PgxHandle, userID int,
) (int, somerr.Error) {
	count, err := model.GetUserCredentialCountByUser(ctx, db, userID)
	if err != nil {
		return 0, somerr.Internal.Errorf("db error: %w", err)
	}
	return count, nil
}

func DeleteUserCredential(ctx context.Context, db model.PgxHandle,
	ID int,
) somerr.Error {
	err := model.DeleteUserCredential(ctx, db, ID)
	if err != nil {
		return somerr.Internal.Error("db error")
	}
	return nil
}
