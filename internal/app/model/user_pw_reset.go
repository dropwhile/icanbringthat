package model

import (
	"context"
	"time"

	"github.com/dropwhile/refid"
)

//go:generate go run ../../../cmd/refidgen -t UserPWReset -v 5

type UserPWReset struct {
	RefID   UserPWResetRefID `db:"ref_id"`
	UserID  int              `db:"user_id"`
	Created time.Time
}

const UserPWResetExpiry = 30 * time.Minute

func NewUserPWReset(ctx context.Context, db PgxHandle,
	userID int,
) (*UserPWReset, error) {
	refID := refid.Must(NewUserPWResetRefID())
	return CreateUserPWReset(ctx, db, refID, userID)
}

func CreateUserPWReset(ctx context.Context, db PgxHandle,
	refID UserPWResetRefID, userID int,
) (*UserPWReset, error) {
	q := `
		INSERT INTO user_pw_reset_ (
			ref_id, user_id
		)
		VALUES ($1, $2)
		RETURNING *`
	return QueryOneTx[UserPWReset](ctx, db, q, refID, userID)
}

func DeleteUserPWReset(ctx context.Context, db PgxHandle,
	refID UserPWResetRefID,
) error {
	q := `DELETE FROM user_pw_reset_ WHERE ref_id = $1`
	return ExecTx[UserPWReset](ctx, db, q, refID)
}

func GetUserPWResetByRefID(ctx context.Context, db PgxHandle,
	refID UserPWResetRefID,
) (*UserPWReset, error) {
	q := `SELECT * FROM user_pw_reset_ WHERE ref_id = $1`
	return QueryOne[UserPWReset](ctx, db, q, refID)
}
