package model

import (
	"context"
	"time"

	"github.com/dropwhile/refid"
	"github.com/jackc/pgx/v5"
)

//go:generate go run ../../../cmd/refidgen -t UserVerify -v 6

type UserVerify struct {
	RefID   UserVerifyRefID `db:"ref_id"`
	UserID  int             `db:"user_id"`
	Created time.Time
}

const UserVerifyExpiry = 30 * time.Minute

func NewUserVerify(ctx context.Context, db PgxHandle,
	user *User,
) (*UserVerify, error) {
	refID := refid.Must(NewUserVerifyRefID())
	return CreateUserVerify(ctx, db, refID, user.ID)
}

func CreateUserVerify(ctx context.Context, db PgxHandle,
	refID UserVerifyRefID, userID int,
) (*UserVerify, error) {
	q := `
		INSERT INTO user_verify_ (
			ref_id, user_id
		)
		VALUES (@refID, @userID)
		RETURNING *`
	args := pgx.NamedArgs{"refID": refID, "userID": userID}
	return QueryOneTx[UserVerify](ctx, db, q, args)
}

func DeleteUserVerify(ctx context.Context, db PgxHandle,
	refID UserVerifyRefID,
) error {
	q := `DELETE FROM user_verify_ WHERE ref_id = $1`
	return ExecTx[UserVerify](ctx, db, q, refID)
}

func GetUserVerifyByRefID(ctx context.Context, db PgxHandle,
	refID UserVerifyRefID,
) (*UserVerify, error) {
	q := `SELECT * FROM user_verify_ WHERE ref_id = $1`
	return QueryOne[UserVerify](ctx, db, q, refID)
}
