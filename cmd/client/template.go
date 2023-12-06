package main

import (
	"fmt"
	"html/template"
	"reflect"
	"strconv"

	"github.com/rs/zerolog/log"
)

var funcMap = template.FuncMap{
	// isset is a helper func from hugo
	"isset": func(ac, kv reflect.Value) (bool, error) {
	SWITCH:
		switch ac.Kind() {
		case reflect.Array, reflect.Slice:
			k := 0
			switch kv.Kind() {
			case reflect.Int | reflect.Int8 | reflect.Int16 | reflect.Int32 | reflect.Int64:
				k = int(kv.Int())
			case reflect.Uint | reflect.Uint8 | reflect.Uint16 | reflect.Uint32 | reflect.Uint64:
				k = int(kv.Uint())
			case reflect.String:
				v, err := strconv.ParseInt(kv.String(), 0, 0)
				if err != nil {
					return false, fmt.Errorf("unable to cast %#v of type %T to int64", kv, kv)
				}
				k = int(v)
			default:
				return false, fmt.Errorf("unable to cast %#v of type %T to int", kv, kv)
			}
			if ac.Len() > k {
				return true, nil
			}
		case reflect.Ptr:
			ac = ac.Elem()
			goto SWITCH
		case reflect.Struct:
			ac.FieldByName(kv.String()).IsValid()
		case reflect.Map:
			if kv.Type() == ac.Type().Key() {
				return ac.MapIndex(kv).IsValid(), nil
			}
		default:
			log.Info().
				Msgf("calling IsSet with unsupported type %q (%T) will always return false", ac.Kind(), ac)
		}
		return false, nil
	},
}
