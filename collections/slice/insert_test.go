package slice

import (
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// TestInsertAt insert the element into slice at index.
// func InsertAt[T any](slice []T, index int, value any) []T
func TestInsertAt(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		result1 := slice.InsertAt([]string{"a", "b", "c"}, 0, "1")
		result2 := slice.InsertAt([]string{"a", "b", "c"}, 1, "1")
		result3 := slice.InsertAt([]string{"a", "b", "c"}, 2, "1")
		result4 := slice.InsertAt([]string{"a", "b", "c"}, 3, "1")
		result5 := slice.InsertAt([]string{"a", "b", "c"}, 0, []string{"1", "2", "3"})

		assert.Equal(t, []string{"1", "a", "b", "c"}, result1)
		assert.Equal(t, []string{"a", "1", "b", "c"}, result2)
		assert.Equal(t, []string{"a", "b", "1", "c"}, result3)
		assert.Equal(t, []string{"a", "b", "c", "1"}, result4)
		assert.Equal(t, []string{"1", "2", "3", "a", "b", "c"}, result5)
	}

	{
		sample := []string{"a", "b", "c", "d", "e", "f", "g"}

		// normal case
		results := lo.Splice(sample, 1, "1", "2")
		is.Equal([]string{"a", "b", "c", "d", "e", "f", "g"}, sample)
		is.Equal([]string{"a", "1", "2", "b", "c", "d", "e", "f", "g"}, results)

		// check there is no side effect
		results = lo.Splice(sample, 1)
		results[0] = "b"
		is.Equal([]string{"a", "b", "c", "d", "e", "f", "g"}, sample)

		// positive overflow
		results = lo.Splice(sample, 42, "1", "2")
		is.Equal([]string{"a", "b", "c", "d", "e", "f", "g"}, sample)
		is.Equal(results, []string{"a", "b", "c", "d", "e", "f", "g", "1", "2"})

		// negative overflow
		results = lo.Splice(sample, -42, "1", "2")
		is.Equal([]string{"a", "b", "c", "d", "e", "f", "g"}, sample)
		is.Equal(results, []string{"1", "2", "a", "b", "c", "d", "e", "f", "g"})

		// backward
		results = lo.Splice(sample, -2, "1", "2")
		is.Equal([]string{"a", "b", "c", "d", "e", "f", "g"}, sample)
		is.Equal(results, []string{"a", "b", "c", "d", "e", "1", "2", "f", "g"})

		results = lo.Splice(sample, -7, "1", "2")
		is.Equal([]string{"a", "b", "c", "d", "e", "f", "g"}, sample)
		is.Equal(results, []string{"1", "2", "a", "b", "c", "d", "e", "f", "g"})

		// other
		is.Equal([]string{"1", "2"}, lo.Splice([]string{}, 0, "1", "2"))
		is.Equal([]string{"1", "2"}, lo.Splice([]string{}, 1, "1", "2"))
		is.Equal([]string{"1", "2"}, lo.Splice([]string{}, -1, "1", "2"))
		is.Equal([]string{"1", "2", "0"}, lo.Splice([]string{"0"}, 0, "1", "2"))
		is.Equal([]string{"0", "1", "2"}, lo.Splice([]string{"0"}, 1, "1", "2"))
		is.Equal([]string{"1", "2", "0"}, lo.Splice([]string{"0"}, -1, "1", "2"))

		// type preserved
		type myStrings []string
		allStrings := myStrings{"", "foo", "bar"}
		nonempty := lo.Splice(allStrings, 1, "1", "2")
		is.IsType(nonempty, allStrings, "type preserved")
	}
}
