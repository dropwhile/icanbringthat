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

type PasswdUpdate struct {
	NewPass []byte `validate:"omitempty,notblank,gt=0"`
	OldPass []byte `validate:"omitempty,notblank,gt=0"`
}

type UserUpdateValues struct {
	Name      mo.Option[string] `validate:"omitnil,notblank"`
	Email     mo.Option[string] `validate:"omitnil,notblank,email"`
	Verified  mo.Option[bool]
	PWAuth    mo.Option[bool]
	ApiAccess mo.Option[bool]
	WebAuthn  mo.Option[bool]
	PwUpdate  mo.Option[*PasswdUpdate] `validate:"omitnil"`
}

func (s *Service) UpdateUser(
	ctx context.Context, user *model.User, euvs *UserUpdateValues,
) errs.Error {
	err := validate.Validate.StructCtx(ctx, euvs)
	if err != nil {
		badField := validate.GetErrorField(err)
		slog.
			With("field", badField).
			With("error", err).
			Info("bad field value")
		return errs.InvalidArgumentError(badField, "bad value")
	}

	pwHash := mo.None[[]byte]()
	if pwup, ok := euvs.PwUpdate.Get(); ok && pwup != nil {
		if err := validate.Validate.VarCtx(ctx, pwup.OldPass, "notblank,gt=0"); err != nil {
			return errs.InvalidArgumentError("OldPass", "bad value")
		}
		if err := validate.Validate.VarCtx(ctx, pwup.NewPass, "notblank,gt=0"); err != nil {
			return errs.InvalidArgumentError("Passwd", "bad value")
		}
		if ok, err := model.CheckPass(ctx, user.PWHash, pwup.OldPass); !ok || err != nil {
			return errs.InvalidArgumentError("OldPass", "bad value")
		}
		hpw, err := model.HashPass(ctx, pwup.NewPass)
		if err != nil {
			return errs.Internal.Error("failed to set password")
		}
		pwHash = mo.Some(hpw)
	}

	err = model.UpdateUser(ctx, s.Db, user.ID, &model.UserUpdateModelValues{
		Name:      euvs.Name,
		Email:     euvs.Email,
		PWHash:    pwHash,
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

func CheckPass(ctx context.Context, pwHash []byte, passwd []byte) bool {
	ok, err := model.CheckPass(ctx, pwHash, passwd)
	return err == nil && ok
}
