// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package resources

import (
	"testing"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestPagination(t *testing.T) {
	t.Parallel()

	expected := []struct {
		numPages    int
		pageSize    int
		currentPage int
		eIsCurrent  int
		output      []string
	}{
		{
			200, 10, 10, 7,
			[]string{
				"1", "2", "3", "...", "8", "9", "10", "11", "12", "...", "18", "19", "20",
			},
		},
		{
			200, 10, 11, 7,
			[]string{
				"1", "2", "3", "...", "9", "10", "11", "12", "13", "...", "18", "19", "20",
			},
		},
		{
			200, 10, 20, 7,
			[]string{
				"1", "2", "3", "...", "18", "19", "20",
			},
		},
		{
			200, 10, 19, 7,
			[]string{
				"1", "2", "3", "...", "17", "18", "19", "20",
			},
		},
		{
			200, 10, 18, 7,
			[]string{
				"1", "2", "3", "...", "16", "17", "18", "19", "20",
			},
		},
		{
			200, 10, 14, 7,
			[]string{
				"1", "2", "3", "...", "12", "13", "14", "15", "16", "...", "18", "19", "20",
			},
		},
		{
			200, 10, 1, 1,
			[]string{
				"1", "2", "3", "...", "18", "19", "20",
			},
		},
		{
			200, 10, 3, 3,
			[]string{
				"1", "2", "3", "4", "5", "...", "18", "19", "20",
			},
		},
		{
			200, 10, 4, 4,
			[]string{
				"1", "2", "3", "4", "5", "6", "...", "18", "19", "20",
			},
		},
		{
			200, 10, 5, 5,
			[]string{
				"1", "2", "3", "4", "5", "6", "7", "...", "18", "19", "20",
			},
		},
		{
			200, 10, 6, 6,
			[]string{
				"1", "2", "3", "4", "5", "6", "7", "8", "...", "18", "19", "20",
			},
		},
		{
			200, 10, 7, 7,
			[]string{
				"1", "2", "3", "...", "5", "6", "7", "8", "9", "...", "18", "19", "20",
			},
		},
		{
			40, 10, 1, 1,
			[]string{
				"1", "2", "3", "4",
			},
		},
		{
			50, 10, 1, 1,
			[]string{
				"1", "2", "3", "4", "5",
			},
		},
		{
			60, 10, 1, 1,
			[]string{
				"1", "2", "3", "4", "5", "6",
			},
		},
		{
			70, 10, 1, 1,
			[]string{
				"1", "2", "3", "...", "5", "6", "7",
			},
		},
	}

	for _, v := range expected {
		pages := NewPaginator(2, 3, 3).AddPages(v.numPages, v.pageSize).Paginate(v.currentPage)
		output := make([]string, 0)
		for i, pg := range pages {
			output = append(output, pg.String())
			if i+1 == v.eIsCurrent {
				assert.Equal(t, pg.IsCurrent, true,
					"page %d unexpectedly not IsCurrent for current=%d", i+1, v.currentPage,
				)
			} else {
				assert.Equal(
					t, pg.IsCurrent, false,
					"page %d unexpectedly IsCurrent for current=%d", i+1, v.currentPage)
			}
		}
		assert.Assert(
			t, cmp.DeepEqual(v.output, output),
			"unexpectedly mismatch for namePages=%d, pageSize=%d, currentPage=%d",
			v.numPages, v.pageSize, v.currentPage,
		)
	}
}
