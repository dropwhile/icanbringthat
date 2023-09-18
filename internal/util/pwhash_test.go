package util

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestRoundTrip(t *testing.T) {
	pwBytes := []byte("test-credentialsµ!")
	pw, err := HashPW(pwBytes)
	assert.NilError(t, err)
	assert.Assert(t, len(pw) != 0)
	ok, err := CheckPWHash(pw, pwBytes)
	assert.NilError(t, err)
	assert.Assert(t, ok)

	pwBytes = append(pwBytes, 'x')
	ok, err = CheckPWHash(pw, pwBytes)
	assert.NilError(t, err)
	assert.Assert(t, !ok)
}
