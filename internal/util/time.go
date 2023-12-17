package util

import "time"

func MustParseTime(layout, value string) time.Time {
	ts, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return ts
}

type CloseTimeMatcher struct {
	Value  time.Time
	Within time.Duration
}

// Match satisfies sqlmock.Argument interface
func (a CloseTimeMatcher) Match(v interface{}) bool {
	ts, ok := v.(time.Time)
	if !ok {
		return false
	}
	if ts.Equal(a.Value) {
		return true
	}
	if ts.Before(a.Value) {
		return !ts.Add(a.Within).Before(a.Value)
	}
	if ts.After(a.Value) {
		return !ts.Add(-a.Within).After(a.Value)
	}
	return true
}
