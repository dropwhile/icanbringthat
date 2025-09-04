// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package crypto

import (
	"testing"

	"github.com/dropwhile/assert"
)

func TestRoundTrip(t *testing.T) {
	pwBytes := []byte("test-credentialsµ!")
	pw, err := HashPW(pwBytes)
	assert.Nil(t, err)
	assert.True(t, len(pw) != 0)
	ok, err := CheckPWHash(pw, pwBytes)
	assert.Nil(t, err)
	assert.True(t, ok)

	pwBytes = append(pwBytes, 'x')
	ok, err = CheckPWHash(pw, pwBytes)
	assert.Nil(t, err)
	assert.True(t, !ok)
}
