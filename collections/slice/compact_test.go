package slice

import (
	"testing"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// TestCompact 创建一些切片，其中删除了所有假值。值false、nil、0和""是假值。
// func Compact[T comparable](slice []T) []T
func TestCompact(t *testing.T) {
	{
		result1 := slice.Compact([]int{0})
		result2 := slice.Compact([]int{0, 1, 2, 3})
		result3 := slice.Compact([]string{"", "a", "b", "0"})
		result4 := slice.Compact([]bool{false, true, true})

		assert.Equal(t, []int{}, result1)
		assert.Equal(t, []int{1, 2, 3}, result2)
		assert.Equal(t, []string{"a", "b", "0"}, result3)
		assert.Equal(t, []bool{true, true}, result4)
	}

	{
		t.Parallel()
		is := assert.New(t)

		r1 := lo.Compact([]int{2, 0, 4, 0})

		is.Equal(r1, []int{2, 4})

		r2 := lo.Compact([]string{"", "foo", "", "bar", ""})

		is.Equal(r2, []string{"foo", "bar"})

		r3 := lo.Compact([]bool{true, false, true, false})

		is.Equal(r3, []bool{true, true})

		type foo struct {
			bar int
			baz string
		}

		// slice of structs
		// If all fields of an element are zero values, Compact removes it.

		r4 := lo.Compact([]foo{
			{bar: 1, baz: "a"}, // all fields are non-zero values
			{bar: 0, baz: ""},  // all fields are zero values
			{bar: 2, baz: ""},  // bar is non-zero
		})

		is.Equal(r4, []foo{{bar: 1, baz: "a"}, {bar: 2, baz: ""}})

		// slice of pointers to structs
		// If an element is nil, Compact removes it.

		e1, e2, e3 := foo{bar: 1, baz: "a"}, foo{bar: 0, baz: ""}, foo{bar: 2, baz: ""}
		// NOTE: e2 is a zero value of foo, but its pointer &e2 is not a zero value of *foo.
		r5 := lo.Compact([]*foo{&e1, &e2, nil, &e3})

		is.Equal(r5, []*foo{&e1, &e2, &e3})

		type myStrings []string
		allStrings := myStrings{"", "foo", "bar"}
		nonempty := lo.Compact(allStrings)
		is.IsType(nonempty, allStrings, "type preserved")
	}
}

func TestWithoutEmpty(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.Compact([]int{0, 1, 2})
	result2 := lo.Compact([]int{1, 2})
	result3 := lo.Compact([]int{})
	is.Equal(result1, []int{1, 2})
	is.Equal(result2, []int{1, 2})
	is.Equal(result3, []int{})

	type myStrings []string
	allStrings := myStrings{"", "foo", "bar"}
	nonempty := lo.Compact(allStrings)
	is.IsType(nonempty, allStrings, "type preserved")
}
