package model

import (
	"context"
	"fmt"
	"time"

	"github.com/dropwhile/refid"
)

var UserPWResetRefIdT = refid.RefIdTagger(5)

type UserPWReset struct {
	RefId   refid.RefId `db:"ref_id"`
	UserId  int         `db:"user_id"`
	Created time.Time
}

func (upw *UserPWReset) IsExpired() bool {
	return upw.RefId.Time().Add(30 * time.Minute).Before(time.Now())
}

func (upw *UserPWReset) Insert(ctx context.Context, db PgxHandle) error {
	if upw.RefId.IsNil() {
		upw.RefId = UserPWResetRefIdT.MustNew()
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
	if !UserPWResetRefIdT.HasCorrectTag(refId) {
		err := fmt.Errorf(
			"bad refid type: got %d expected %d",
			refId.Tag(), UserPWResetRefIdT.Tag(),
		)
		return nil, err
	}
	q := `SELECT * FROM user_pw_reset_ WHERE ref_id = $1`
	return QueryOne[UserPWReset](ctx, db, q, refId)
}
