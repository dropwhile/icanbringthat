package model

import (
	"context"
	"fmt"
	"time"

	"github.com/dropwhile/refid/v2"
	"github.com/dropwhile/refid/v2/reftag"
	"github.com/jackc/pgx/v5"
)

type UserVerifyRefID = reftag.IDt6

type UserVerifyRefIDNull struct {
	reftag.NullIDt6
}

var (
	NewUserVerifyRefID       = reftag.New[UserVerifyRefID]
	UserVerifyRefIDMatcher   = reftag.NewMatcher[UserVerifyRefID]()
	UserVerifyRefIDFromBytes = reftag.FromBytes[UserVerifyRefID]
	ParseUserVerifyRefID     = reftag.Parse[UserVerifyRefID]
)

type UserVerify struct {
	Created time.Time
	UserID  int             `db:"user_id"`
	RefID   UserVerifyRefID `db:"ref_id"`
}

const UserVerifyExpiry = 30 * time.Minute

func NewUserVerify(ctx context.Context, db PgxHandle,
	user *User,
) (*UserVerify, error) {
	if user == nil {
		return nil, fmt.Errorf("nil user supplied")
	}
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
