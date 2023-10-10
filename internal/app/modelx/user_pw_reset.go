package modelx

import (
	"time"

	"github.com/dropwhile/refid"
)

var UserPWResetRefIDT = refid.Tagger(5)

func (upw *UserPwReset) IsExpired() bool {
	return upw.RefID.Time().Add(30 * time.Minute).Before(time.Now())
}
