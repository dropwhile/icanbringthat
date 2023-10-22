// Code generated by refidgen. DO NOT EDIT.
// generated from: earmark.go

package model

import (
	"fmt"

	"github.com/dropwhile/refid"
)

const tagValEarmarkRefID = 4

type EarmarkRefID struct {
	refid.RefID
}

func (r *EarmarkRefID) checkResult(err error) error {
	if err != nil {
		return err
	}
	if !r.RefID.HasTag(tagValEarmarkRefID) {
		return fmt.Errorf("wrong refid type")
	}
	return nil
}

func (r *EarmarkRefID) Scan(src interface{}) error {
	err := r.RefID.Scan(src)
	return r.checkResult(err)
}

func (r *EarmarkRefID) UnmarshalJSON(b []byte) error {
	err := r.RefID.UnmarshalJSON(b)
	return r.checkResult(err)
}

func (r *EarmarkRefID) UnmarshalBinary(b []byte) error {
	err := r.RefID.UnmarshalBinary(b)
	return r.checkResult(err)
}

func NewEarmarkRefID() (EarmarkRefID, error) {
	v, err := refid.NewTagged(tagValEarmarkRefID)
	return EarmarkRefID{v}, err
}

func ParseEarmarkRefID(s string) (EarmarkRefID, error) {
	v, err := refid.ParseWithRequire(s, refid.HasTag(tagValEarmarkRefID))
	return EarmarkRefID{v}, err
}

func ParseEarmarkRefIDWithRequire(s string, reqs ...refid.Requirement) (EarmarkRefID, error) {
	reqs = append(reqs, refid.HasTag(tagValEarmarkRefID))
	v, err := refid.ParseWithRequire(s, reqs...)
	return EarmarkRefID{v}, err
}

func EarmarkRefIDFromBytes(input []byte) (EarmarkRefID, error) {
	var r EarmarkRefID
	err := r.UnmarshalBinary(input)
	return r, err
}

type NullEarmarkRefID struct {
	refid.NullRefID
}

func (u *NullEarmarkRefID) checkResult(err error) error {
	if err != nil {
		return err
	}
	n := u.NullRefID
	if n.Valid && !n.RefID.HasTag(tagValEarmarkRefID) {
		return fmt.Errorf("wrong refid type")
	}
	return nil
}


func (u *NullEarmarkRefID) Scan(src interface{}) error {
	err := u.NullRefID.Scan(src)
	return u.checkResult(err)
}

func (u *NullEarmarkRefID) UnmarshalJSON(b []byte) error {
	err := u.NullRefID.UnmarshalJSON(b)
	return u.checkResult(err)
}

type EarmarkRefIDMatcher struct{}

func (a EarmarkRefIDMatcher) Match(v interface{}) bool {
	var r refid.RefID
	var err error
	switch x := v.(type) {
	case EarmarkRefID:
		r = x.RefID
	case *EarmarkRefID:
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
	return r.HasTag(tagValEarmarkRefID)
}
