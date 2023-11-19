package encoder

import (
	"encoding/base32"
	"reflect"
)

var (
	Alphabet         = "0123456789abcdefghjkmnpqrstvwxyz"
	WordSafeEncoding = base32.NewEncoding(Alphabet).WithPadding(base32.NoPadding)
)

func Base32EncodeToString(src []byte) string {
	return WordSafeEncoding.EncodeToString(src)
}

func Base32DecodeString(src string) ([]byte, error) {
	return WordSafeEncoding.DecodeString(src)
}

func StructToMap(s interface{}) map[string]interface{} {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(reflect.ValueOf(s))
	}
	values := make(map[string]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanInterface() {
			name := v.Type().Field(i).Name
			val := v.Field(i).Interface()
			values[name] = val
		}
	}

	return values
}
