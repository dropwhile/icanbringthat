// Code generated by refidgen. DO NOT EDIT.

package modelx

import (
	"fmt"

	"github.com/dropwhile/refid"
)

const tagValUserRefID = 1

type UserRefID struct {
	refid.RefID
}

func (r *UserRefID) checkResult(err error) error {
	if err != nil {
		return err
	}
	if !r.RefID.HasTag(tagValUserRefID) {
		return fmt.Errorf("wrong refid type")
	}
	return nil
}

func (r *UserRefID) Scan(src interface{}) error {
	err := r.RefID.Scan(src)
	return r.checkResult(err)
}

func (r *UserRefID) UnmarshalJSON(b []byte) error {
	err := r.RefID.UnmarshalJSON(b)
	return r.checkResult(err)
}

func (r *UserRefID) UnmarshalBinary(b []byte) error {
	err := r.RefID.UnmarshalJSON(b)
	return r.checkResult(err)
}

func NewUserRefID() (UserRefID, error) {
	v, err := refid.NewTagged(tagValUserRefID)
	return UserRefID{v}, err
}

func ParseUserRefID(s string) (UserRefID, error) {
	v, err := refid.ParseWithRequire(s, refid.HasTag(tagValUserRefID))
	return UserRefID{v}, err
}

func ParseUserRefIDWithRequire(s string, reqs ...refid.Requirement) (UserRefID, error) {
	reqs = append(reqs, refid.HasTag(tagValUserRefID))
	v, err := refid.ParseWithRequire(s, reqs...)
	return UserRefID{v}, err
}

type NullUserRefID struct {
	refid.NullRefID
}

func (u *NullUserRefID) checkResult(err error) error {
	if err != nil {
		return err
	}
	n := u.NullRefID
	if n.Valid && !n.RefID.HasTag(tagValUserRefID) {
		return fmt.Errorf("wrong refid type")
	}
	return nil
}


func (u *NullUserRefID) Scan(src interface{}) error {
	err := u.NullRefID.Scan(src)
	return u.checkResult(err)
}

func (u *NullUserRefID) UnmarshalJSON(b []byte) error {
	err := u.NullRefID.UnmarshalJSON(b)
	return u.checkResult(err)
}

type UserRefIDMatcher struct{}

func (a UserRefIDMatcher) Match(v interface{}) bool {
	var r refid.RefID
	var err error
	switch x := v.(type) {
	case UserRefID:
		r = x.RefID
	case *UserRefID:
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
	return r.HasTag(tagValUserRefID)
}
