package model

import "time"

type BifurcatedRowCounts struct {
	Current  int
	Archived int
}

type Timer interface {
	Time() time.Time
}

func IsExpired(tm Timer, expiry time.Duration) bool {
	return tm.Time().Add(expiry).Before(time.Now())
}
