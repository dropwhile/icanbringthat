package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
)

func GetUser(
	ctx context.Context, db model.PgxHandle,
	refID model.UserRefID,
) (*model.User, errs.Error) {
	user, err := model.GetUserByRefID(ctx, db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("user not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return user, nil
}

func GetUserByEmail(
	ctx context.Context, db model.PgxHandle,
	email string,
) (*model.User, errs.Error) {
	user, err := model.GetUserByEmail(ctx, db, email)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("user not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return user, nil
}

func GetUserByID(
	ctx context.Context, db model.PgxHandle,
	ID int,
) (*model.User, errs.Error) {
	user, err := model.GetUserByID(ctx, db, ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errs.NotFound.Error("user not found")
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return user, nil
}

func GetUsersByIDs(ctx context.Context, db model.PgxHandle,
	userIDs []int,
) ([]*model.User, errs.Error) {
	if len(userIDs) == 0 {
		return []*model.User{}, nil
	}

	users, err := model.GetUsersByIDs(ctx, db, userIDs)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		users = []*model.User{}
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return users, nil
}

func NewUser(ctx context.Context, db model.PgxHandle,
	email, name string, rawPass []byte,
) (*model.User, errs.Error) {
	err := validate.VarCtx(ctx, name, "required,notblank")
	if err != nil {
		slog.
			With("field", "name").
			With("error", err).
			Info("bad field value")
		return nil, errs.InvalidArgumentError("name", "bad value")
	}
	err = validate.VarCtx(ctx, email, "required,notblank,email")
	if err != nil {
		slog.
			With("field", "email").
			With("error", err).
			Info("bad field value")
		return nil, errs.InvalidArgumentError("email", "bad value")
	}

	user, err := model.NewUser(ctx, db, email, name, rawPass)
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

func UpdateUser(ctx context.Context, db model.PgxHandle,
	email, name string, pwHash []byte, verified bool,
	pwAuth, apiAccess, webAuthn bool, userID int,
) errs.Error {
	err := validate.VarCtx(ctx, name, "required,notblank")
	if err != nil {
		slog.
			With("field", "name").
			With("error", err).
			Info("bad field value")
		return errs.InvalidArgumentError("name", "bad value")
	}
	err = validate.VarCtx(ctx, email, "required,notblank,email")
	if err != nil {
		slog.
			With("field", "email").
			With("error", err).
			Info("bad field value")
		return errs.InvalidArgumentError("email", "bad value")
	}

	err = model.UpdateUser(ctx, db,
		email, name, pwHash, verified,
		pwAuth, apiAccess, webAuthn, userID,
	)
	if err != nil {
		return errs.Internal.Errorf("db error: %w", err)
	}
	return nil
}

func UpdateUserSettings(
	ctx context.Context, db model.PgxHandle, userID int,
	pm *model.UserSettings,
) errs.Error {
	err := model.UpdateUserSettings(ctx, db, pm, userID)
	if err != nil {
		return errs.Internal.Error("db error")
	}
	return nil
}

func DeleteUser(ctx context.Context, db model.PgxHandle,
	userID int,
) errs.Error {
	err := model.DeleteUser(ctx, db, userID)
	if err != nil {
		return errs.Internal.Errorf("db error: %w", err)
	}
	return nil
}
