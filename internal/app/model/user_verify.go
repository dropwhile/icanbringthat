package model

import (
	"context"
	"fmt"
	"time"

	"github.com/dropwhile/refid"
)

var VerifyRefIDT = refid.Tagger(6)

type UserVerify struct {
	RefID   refid.RefID `db:"ref_id"`
	UserID  int         `db:"user_id"`
	Created time.Time
}

func (uv *UserVerify) IsExpired() bool {
	return uv.RefID.Time().Add(30 * time.Minute).Before(time.Now())
}

func (uv *UserVerify) Insert(ctx context.Context, db PgxHandle) error {
	if uv.RefID.IsNil() {
		uv.RefID = refid.Must(VerifyRefIDT.New())
	}
	q := `INSERT INTO user_verify_ (ref_id, user_id) VALUES ($1, $2) RETURNING *`
	res, err := QueryOneTx[UserVerify](ctx, db, q, uv.RefID, uv.UserID)
	if err != nil {
		return err
	}
	uv.RefID = res.RefID
	return nil
}

func (uv *UserVerify) Delete(ctx context.Context, db PgxHandle) error {
	q := `DELETE FROM user_verify_ WHERE ref_id = $1`
	return ExecTx[UserVerify](ctx, db, q, uv.RefID)
}

func NewUserVerify(ctx context.Context, db PgxHandle, user *User) (*UserVerify, error) {
	uv := &UserVerify{
		UserID: user.ID,
	}
	err := uv.Insert(ctx, db)
	if err != nil {
		return nil, err
	}
	return uv, nil
}

func GetUserVerifyByRefID(ctx context.Context, db PgxHandle, refID refid.RefID) (*UserVerify, error) {
	if !VerifyRefIDT.HasCorrectTag(refID) {
		err := fmt.Errorf(
			"bad refid type: got %d expected %d",
			refID.Tag(), VerifyRefIDT.Tag(),
		)
		return nil, err
	}
	q := `SELECT * FROM user_verify_ WHERE ref_id = $1`
	return QueryOne[UserVerify](ctx, db, q, refID)
}
