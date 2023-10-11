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

func (upw *UserPWReset) IsExpired() bool {
	return upw.RefID.Time().Add(30 * time.Minute).Before(time.Now())
}

func (upw *UserPWReset) Insert(ctx context.Context, db PgxHandle) error {
	if upw.RefID.IsNil() {
		upw.RefID = refid.Must(NewUserPWResetRefID())
	}
	q := `INSERT INTO user_pw_reset_ (ref_id, user_id) VALUES ($1, $2) RETURNING *`
	res, err := QueryOneTx[UserPWReset](ctx, db, q, upw.RefID, upw.UserID)
	if err != nil {
		return err
	}
	upw.RefID = res.RefID
	return nil
}

func (upw *UserPWReset) Delete(ctx context.Context, db PgxHandle) error {
	q := `DELETE FROM user_pw_reset_ WHERE ref_id = $1`
	return ExecTx[UserPWReset](ctx, db, q, upw.RefID)
}

func NewUserPWReset(ctx context.Context, db PgxHandle, user *User) (*UserPWReset, error) {
	upw := &UserPWReset{
		UserID: user.ID,
	}
	err := upw.Insert(ctx, db)
	if err != nil {
		return nil, err
	}
	return upw, nil
}

func GetUserPWResetByRefID(ctx context.Context, db PgxHandle, refID UserPWResetRefID) (*UserPWReset, error) {
	q := `SELECT * FROM user_pw_reset_ WHERE ref_id = $1`
	return QueryOne[UserPWReset](ctx, db, q, refID)
}
