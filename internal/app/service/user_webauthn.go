package service

import (
	"context"
	"errors"

	"github.com/dropwhile/refid/v2/reftag"
	"github.com/jackc/pgx/v5"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
)

var (
	CredentialRefIDMatcher   = reftag.NewMatcher[model.CredentialRefID]()
	CredentialRefIDFromBytes = reftag.FromBytes[model.CredentialRefID]
	ParseCredentialRefID     = reftag.Parse[model.CredentialRefID]
)

func (s *Service) GetUserCredentialByRefID(
	ctx context.Context, refID model.CredentialRefID,
) (*model.UserCredential, errs.Error) {
	cred, err := model.GetUserCredentialByRefID(ctx, s.Db, refID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, errs.NotFound.Error("credential")
		default:
			return nil, errs.Internal.Errorf("db error: %w", err)
		}
	}
	return cred, nil
}

func (s *Service) GetUserCredentialsByUser(
	ctx context.Context, userID int,
) ([]*model.UserCredential, errs.Error) {
	creds, err := model.GetUserCredentialsByUser(ctx, s.Db, userID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return []*model.UserCredential{}, nil
		default:
			return nil, errs.Internal.Errorf("db error: %w", err)
		}
	}
	return creds, nil
}

func (s *Service) GetUserCredentialCountByUser(
	ctx context.Context, userID int,
) (int, errs.Error) {
	count, err := model.GetUserCredentialCountByUser(ctx, s.Db, userID)
	if err != nil {
		return 0, errs.Internal.Errorf("db error: %w", err)
	}
	return count, nil
}

func (s *Service) DeleteUserCredential(
	ctx context.Context, ID int,
) errs.Error {
	err := model.DeleteUserCredential(ctx, s.Db, ID)
	if err != nil {
		return errs.Internal.Error("db error")
	}
	return nil
}

func (s *Service) NewUserCredential(
	ctx context.Context, userID int, keyName string, credential []byte,
) (*model.UserCredential, errs.Error) {
	userCred, err := model.NewUserCredential(
		ctx, s.Db, userID, keyName, credential,
	)
	if err != nil {
		return nil, errs.Internal.Error("db error")
	}
	return userCred, nil
}
