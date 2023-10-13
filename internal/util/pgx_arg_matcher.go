package util

import (
	"reflect"
	"slices"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
)

type pgxNamedArgsArgument struct {
	args pgx.NamedArgs
}

func NewPgxNamedArgsMatcher(args pgx.NamedArgs) pgxmock.Argument {
	return pgxNamedArgsArgument{args}
}

func (a pgxNamedArgsArgument) Match(x interface{}) bool {
	keys := KeysSorted(a.args)
	if namedArgs, ok := x.(pgx.NamedArgs); ok {
		// match all keys
		nk := KeysSorted(namedArgs)
		if !slices.Equal(nk, keys) {
			return false
		}
		// match all values
		for k := range a.args {
			if matcher, ok := a.args[k].(pgxmock.Argument); ok {
				if !matcher.Match(namedArgs[k]) {
					return false
				}
				continue
			}
			if !reflect.DeepEqual(a.args[k], namedArgs[k]) {
				return false
			}
		}
		return true
	}
	return false
}
