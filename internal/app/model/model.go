package model

import "time"

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

type Timer interface {
	Time() time.Time
}

func IsExpired(tm Timer, expiry time.Duration) bool {
	return tm.Time().Add(expiry).Before(time.Now())
}
