package modelx

import (
	"context"
	"time"
)

func (uv *UserVerify) IsExpired() bool {
	return uv.RefID.Time().Add(30 * time.Minute).Before(time.Now())
}

func (q *Queries) NewUserVerify(ctx context.Context, userID int32) (*UserVerify, error) {
	refID, err := NewVerifyRefID()
	if err != nil {
		return nil, err
	}
	return q.CreateUserVerify(ctx, refID, userID)
}
