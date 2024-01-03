package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/dropwhile/refid/v2/reftag"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samber/mo"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/validate"
)

var (
	ParseUserRefID     = reftag.Parse[model.UserRefID]
	UserRefIDFromBytes = reftag.FromBytes[model.UserRefID]
	UserRefIDMatcher   = reftag.NewMatcher[model.UserRefID]()
)

func (s *Service) GetUser(
	ctx context.Context, refID model.UserRefID,
) (*model.User, errs.Error) {
	user, err := model.GetUserByRefID(ctx, s.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("user not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return user, nil
}

func (s *Service) GetUserByEmail(
	ctx context.Context, email string,
) (*model.User, errs.Error) {
	user, err := model.GetUserByEmail(ctx, s.Db, email)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("user not found")
	case err != nil:
		slog.
			With("error", err).
			Info("error getting user")
		return nil, errs.Internal.Error("db error")
	}
	return user, nil
}

func (s *Service) GetUserByID(
	ctx context.Context, ID int,
) (*model.User, errs.Error) {
	user, err := model.GetUserByID(ctx, s.Db, ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("user not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return user, nil
}

func (s *Service) GetUsersByIDs(
	ctx context.Context, userIDs []int,
) ([]*model.User, errs.Error) {
	if len(userIDs) == 0 {
		return []*model.User{}, nil
	}

	users, err := model.GetUsersByIDs(ctx, s.Db, userIDs)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		users = []*model.User{}
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return users, nil
}

func (s *Service) NewUser(
	ctx context.Context, email, name string, rawPass []byte,
) (*model.User, errs.Error) {
	err := validate.Validate.VarCtx(ctx, name, "required,notblank")
	if err != nil {
		slog.
			With("field", "name").
			With("error", err).
			Info("bad field value")
		return nil, errs.InvalidArgumentError("name", "bad value")
	}
	err = validate.Validate.VarCtx(ctx, email, "required,notblank,email")
	if err != nil {
		slog.
			With("field", "email").
			With("error", err).
			Info("bad field value")
		return nil, errs.InvalidArgumentError("email", "bad value")
	}

	user, err := model.NewUser(ctx, s.Db, email, name, rawPass)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "user_email_idx" {
				return nil, errs.AlreadyExists.Error("user already exists")
			}
		}
		return nil, errs.Internal.Errorf("error creating user: %w", err)
	}
	return user, nil
}

type UserUpdateValues struct {
	Name      mo.Option[string] `validate:"omitempty,notblank"`
	Email     mo.Option[string] `validate:"omitempty,notblank,email"`
	PWHash    mo.Option[[]byte] `validate:"omitempty,gt=0"`
	Verified  mo.Option[bool]
	PWAuth    mo.Option[bool]
	ApiAccess mo.Option[bool]
	WebAuthn  mo.Option[bool]
}

func (s *Service) UpdateUser(
	ctx context.Context, userID int, euvs *UserUpdateValues,
) errs.Error {
	// buggy: see https://github.com/go-playground/validator/issues/1209
	err := validate.Validate.StructCtx(ctx, euvs)
	if err != nil {
		badField := validate.GetErrorField(err)
		slog.
			With("field", badField).
			With("error", err).
			Info("bad field value")
		return errs.InvalidArgumentError(badField, "bad value")
	}

	err = model.UpdateUser(ctx, s.Db, userID, &model.UserUpdateModelValues{
		Name:      euvs.Name,
		Email:     euvs.Email,
		PWHash:    euvs.PWHash,
		Verified:  euvs.Verified,
		PWAuth:    euvs.PWAuth,
		ApiAccess: euvs.ApiAccess,
		WebAuthn:  euvs.WebAuthn,
	})
	if err != nil {
		return errs.Internal.Errorf("db error: %w", err)
	}
	return nil
}

func (s *Service) UpdateUserSettings(
	ctx context.Context, userID int, pm *model.UserSettings,
) errs.Error {
	return s.updateUserSettings(ctx, s.Db, userID, pm)
}

func (s *Service) updateUserSettings(
	ctx context.Context, db model.PgxHandle,
	userID int, pm *model.UserSettings,
) errs.Error {
	err := model.UpdateUserSettings(ctx, db, userID, pm)
	if err != nil {
		return errs.Internal.Error("db error")
	}
	return nil
}

func (s *Service) DeleteUser(
	ctx context.Context, userID int,
) errs.Error {
	err := model.DeleteUser(ctx, s.Db, userID)
	if err != nil {
		return errs.Internal.Errorf("db error: %w", err)
	}
	return nil
}
