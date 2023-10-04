package model

import (
	"context"
	"fmt"
	"time"

	"github.com/dropwhile/refid"
)

var UserPWResetRefIDT = refid.Tagger(5)

type UserPWReset struct {
	RefID   refid.RefID `db:"ref_id"`
	UserId  int         `db:"user_id"`
	Created time.Time
}

func (upw *UserPWReset) IsExpired() bool {
	return upw.RefID.Time().Add(30 * time.Minute).Before(time.Now())
}

func (upw *UserPWReset) Insert(ctx context.Context, db PgxHandle) error {
	if upw.RefID.IsNil() {
		upw.RefID = refid.Must(UserPWResetRefIDT.New())
	}
	q := `INSERT INTO user_pw_reset_ (ref_id, user_id) VALUES ($1, $2) RETURNING *`
	res, err := QueryOneTx[UserPWReset](ctx, db, q, upw.RefID, upw.UserId)
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
		UserId: user.Id,
	}
	err := upw.Insert(ctx, db)
	if err != nil {
		return nil, err
	}
	return upw, nil
}

func GetUserPWResetByRefID(ctx context.Context, db PgxHandle, refId refid.RefID) (*UserPWReset, error) {
	if !UserPWResetRefIDT.HasCorrectTag(refId) {
		err := fmt.Errorf(
			"bad refid type: got %d expected %d",
			refId.Tag(), UserPWResetRefIDT.Tag(),
		)
		return nil, err
	}
	q := `SELECT * FROM user_pw_reset_ WHERE ref_id = $1`
	return QueryOne[UserPWReset](ctx, db, q, refId)
}
