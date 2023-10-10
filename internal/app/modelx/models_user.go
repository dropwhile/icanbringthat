package modelx

import (
	"context"
	"fmt"

	"github.com/dropwhile/icbt/internal/util"
)

func (user *User) SetPass(ctx context.Context, rawPass []byte) error {
	pwHash, err := util.HashPW([]byte(rawPass))
	if err != nil {
		return fmt.Errorf("error hashing pw: %w", err)
	}
	user.PwHash = pwHash
	return nil
}

func (user *User) CheckPass(ctx context.Context, rawPass []byte) (bool, error) {
	ok, err := util.CheckPWHash(user.PwHash, rawPass)
	if err != nil {
		return false, fmt.Errorf("error when comparing pass")
	}
	return ok, nil
}

func (q *Queries) NewUser(ctx context.Context, email, name string, rawPass []byte) (*User, error) {
	refID, err := NewUserRefID()
	if err != nil {
		return nil, err
	}

	pwHash, err := util.HashPW([]byte(rawPass))
	if err != nil {
		return nil, fmt.Errorf("error hashing pw: %w", err)
	}

	user, err := q.CreateUser(ctx, CreateUserParams{
		RefID:  refID,
		Email:  email,
		Name:   name,
		PwHash: pwHash,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
