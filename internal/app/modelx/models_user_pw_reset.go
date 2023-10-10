package modelx

import (
	"context"
	"time"
)

func (upw *UserPwReset) IsExpired() bool {
	return upw.RefID.Time().Add(30 * time.Minute).Before(time.Now())
}

func (q *Queries) NewUserPWReset(ctx context.Context, userID int32) (*UserPwReset, error) {
	refID, err := NewUserPwResetRefID()
	if err != nil {
		return nil, err
	}
	return q.CreateUserPWReset(ctx, refID, userID)
}
