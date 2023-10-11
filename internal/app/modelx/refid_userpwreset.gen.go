// Code generated by refidgen. DO NOT EDIT.

package modelx

import (
	"fmt"

	"github.com/dropwhile/refid"
)

const tagValUserPwResetRefID = 5

type UserPwResetRefID struct {
	refid.RefID
}

func (r *UserPwResetRefID) checkResult(err error) error {
	if err != nil {
		return err
	}
	if !r.RefID.HasTag(tagValUserPwResetRefID) {
		return fmt.Errorf("wrong refid type")
	}
	return nil
}

func (r *UserPwResetRefID) Scan(src interface{}) error {
	err := r.RefID.Scan(src)
	return r.checkResult(err)
}

func (r *UserPwResetRefID) UnmarshalJSON(b []byte) error {
	err := r.RefID.UnmarshalJSON(b)
	return r.checkResult(err)
}

func (r *UserPwResetRefID) UnmarshalBinary(b []byte) error {
	err := r.RefID.UnmarshalJSON(b)
	return r.checkResult(err)
}

func NewUserPwResetRefID() (UserPwResetRefID, error) {
	v, err := refid.NewTagged(tagValUserPwResetRefID)
	return UserPwResetRefID{v}, err
}

func ParseUserPwResetRefID(s string) (UserPwResetRefID, error) {
	v, err := refid.ParseWithRequire(s, refid.HasTag(tagValUserPwResetRefID))
	return UserPwResetRefID{v}, err
}

func ParseUserPwResetRefIDWithRequire(s string, reqs ...refid.Requirement) (UserPwResetRefID, error) {
	reqs = append(reqs, refid.HasTag(tagValUserPwResetRefID))
	v, err := refid.ParseWithRequire(s, reqs...)
	return UserPwResetRefID{v}, err
}

type NullUserPwResetRefID struct {
	refid.NullRefID
}

func (u *NullUserPwResetRefID) checkResult(err error) error {
	if err != nil {
		return err
	}
	n := u.NullRefID
	if n.Valid && !n.RefID.HasTag(tagValUserPwResetRefID) {
		return fmt.Errorf("wrong refid type")
	}
	return nil
}


func (u *NullUserPwResetRefID) Scan(src interface{}) error {
	err := u.NullRefID.Scan(src)
	return u.checkResult(err)
}

func (u *NullUserPwResetRefID) UnmarshalJSON(b []byte) error {
	err := u.NullRefID.UnmarshalJSON(b)
	return u.checkResult(err)
}

type UserPwResetRefIDMatcher struct{}

func (a UserPwResetRefIDMatcher) Match(v interface{}) bool {
	var r refid.RefID
	var err error
	switch x := v.(type) {
	case UserPwResetRefID:
		r = x.RefID
	case *UserPwResetRefID:
		r = x.RefID
	case string:
		r, err = refid.Parse(x)
	case []byte:
		r, err = refid.FromBytes(x)
	default:
		return false
	}
	if err != nil {
		return false
	}
	return r.HasTag(tagValUserPwResetRefID)
}