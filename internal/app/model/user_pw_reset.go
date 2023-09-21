package model

import (
	"context"

	"github.com/dropwhile/icbt/internal/util/refid"
)

type UserPWReset struct {
	RefId  refid.RefId `db:"ref_id"`
	UserId int         `db:"user_id"`
}

func (upw *UserPWReset) Insert(ctx context.Context, db PgxHandle) error {
	if upw.RefId.IsNil() {
		upw.RefId = UserRefIdT.MustNew()
	}
	q := `INSERT INTO user_pw_reset_ (ref_id, user_id) VALUES ($1, $2) RETURNING *`
	res, err := QueryOneTx[UserPWReset](ctx, db, q, upw.RefId, upw.UserId)
	if err != nil {
		return err
	}
	upw.RefId = res.RefId
	return nil
}

func (upw *UserPWReset) Delete(ctx context.Context, db PgxHandle) error {
	q := `DELETE FROM user_pw_reset_ WHERE ref_id = $1`
	return ExecTx[UserPWReset](ctx, db, q, upw.RefId)
}

func NewUserPWReset(ctx context.Context, db PgxHandle, user *User) (*UserPWReset, error) {
	upw := &UserPWReset{
		UserId: user.Id,
	}
	err := upw.Insert(ctx, db)
	if err != nil {
		return nil, err
	}
	return upw, nil
}

func GetUserPWResetByRefId(ctx context.Context, db PgxHandle, refId refid.RefId) (*UserPWReset, error) {
	q := `SELECT * FROM user_pw_reset_ WHERE ref_id = $1`
	return QueryOne[UserPWReset](ctx, db, q, refId)
}
