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
	ctx context.Context, userID int, refID model.CredentialRefID,
) (*model.UserCredential, errs.Error) {
	cred, err := model.GetUserCredentialByRefID(ctx, s.Db, refID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, errs.NotFound.Error("credential not found")
		default:
			return nil, errs.Internal.Errorf("db error: %w", err)
		}
	}

	if cred.UserID != userID {
		return nil, errs.PermissionDenied.Error("permission denied")
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
	ctx context.Context, user *model.User, refID model.CredentialRefID,
) errs.Error {
	credential, errx := s.GetUserCredentialByRefID(ctx, user.ID, refID)
	if errx != nil {
		return errx
	}

	if credential.UserID != user.ID {
		return errs.PermissionDenied.Error("permission denied")
	}

	count, errx := s.GetUserCredentialCountByUser(ctx, user.ID)
	if errx != nil {
		return errx
	}

	if count == 1 && user.WebAuthn {
		return errs.FailedPrecondition.Error(
			"refusing to remove last passkey when password auth disabled",
		)
	}

	err := model.DeleteUserCredential(ctx, s.Db, credential.ID)
	if err != nil {
		return errs.Internal.Error("db error")
	}
	return nil
}

func (s *Service) NewUserCredential(
	ctx context.Context, userID int, keyName string, credential []byte,
) (*model.UserCredential, errs.Error) {
	if keyName == "" {
		return nil, errs.InvalidArgumentError("keyName", "bad value")
	}
	if len(credential) == 0 {
		return nil, errs.InvalidArgumentError("credential", "bad value")
	}
	userCred, err := model.NewUserCredential(
		ctx, s.Db, userID, keyName, credential,
	)
	if err != nil {
		return nil, errs.Internal.Error("db error")
	}
	return userCred, nil
}
