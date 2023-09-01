// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package util

import (
	"sort"
	"testing"

	"gotest.tools/v3/assert"
)

func TestKeys(t *testing.T) {
	x := map[string]string{
		"one": "test1",
		"two": "test2",
	}

	strKeys := Keys(x)
	sort.Strings(strKeys)
	assert.DeepEqual(t, []string{"one", "two"}, strKeys)

	y := map[int]int{
		1: 21,
		2: 32,
	}

	intKeys := Keys(y)
	sort.Ints(intKeys)
	assert.DeepEqual(t, []int{1, 2}, intKeys)
}
