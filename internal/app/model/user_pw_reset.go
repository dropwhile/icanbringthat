package model

import (
	"context"
	"time"

	"github.com/dropwhile/refid/v2/reftag"
	"github.com/jackc/pgx/v5"

	"github.com/dropwhile/icbt/internal/util"
)

type UserPWResetRefID = reftag.IDt5

type UserPWResetRefIDNull struct {
	reftag.NullIDt5
}

var NewUserPWResetRefID = reftag.New[UserPWResetRefID]

type UserPWReset struct {
	Created time.Time
	UserID  int              `db:"user_id"`
	RefID   UserPWResetRefID `db:"ref_id"`
}

const UserPWResetExpiry = 30 * time.Minute

func NewUserPWReset(ctx context.Context, db PgxHandle,
	userID int,
) (*UserPWReset, error) {
	refID := util.Must(NewUserPWResetRefID())
	return CreateUserPWReset(ctx, db, refID, userID)
}

func CreateUserPWReset(ctx context.Context, db PgxHandle,
	refID UserPWResetRefID, userID int,
) (*UserPWReset, error) {
	q := `
		INSERT INTO user_pw_reset_ (
			ref_id, user_id
		)
		VALUES (@refID, @userID)
		RETURNING *`
	args := pgx.NamedArgs{"refID": refID, "userID": userID}
	return QueryOneTx[UserPWReset](ctx, db, q, args)
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
