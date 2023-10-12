package util

import (
	"reflect"
	"slices"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
)

type pgxNamedArgsArgument struct {
	args pgx.NamedArgs
	keys []string
}

func NewPgxNamedArgsMatcher(args pgx.NamedArgs) pgxmock.Argument {
	keys := Keys(args)
	slices.Sort(keys)
	return pgxNamedArgsArgument{args, keys}
}

func (a pgxNamedArgsArgument) Match(x interface{}) bool {
	if namedArgs, ok := x.(pgx.NamedArgs); ok {
		// match all keys
		nk := Keys(namedArgs)
		slices.Sort(nk)
		if !slices.Equal(nk, a.keys) {
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
