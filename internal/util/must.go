package util

import "io"

func MustReadAll(r io.ReadCloser) []byte {
	d, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return d
}
