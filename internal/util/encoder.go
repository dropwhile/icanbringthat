package util

import "encoding/base32"

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
