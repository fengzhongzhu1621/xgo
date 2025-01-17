package cast

import (
	"fmt"
	"testing"

	"github.com/gookit/goutil/arrutil"

	"github.com/duke-git/lancet/v2/slice"

	"github.com/samber/lo"

	"github.com/duke-git/lancet/v2/convertor"

	"github.com/duke-git/lancet/v2/structs"
	"github.com/stretchr/testify/assert"
)

func TestStr2map(t *testing.T) {
	s := "a=1&b=2&c="
	actual := Str2map(s, "&", "=")
	expect := map[string]string{"a": "1", "b": "2", "c": ""}
	assert.Equal(t, expect, actual)
}

// TestToMap convert a valid struct to a map
// func (s *Struct) ToMap() (map[string]any, error)
// func ToMap(v any) (map[string]any, error)
func TestToMap(t *testing.T) {
	type People struct {
		Name string `json:"name"`
	}
	p1 := &People{Name: "11"}
	// use constructor function
	s1 := structs.New(p1)
	m1, _ := s1.ToMap()
	fmt.Println(m1) // map[name:11]

	// use static function
	m2, _ := structs.ToMap(p1)
	fmt.Println(m2) // map[name:11]
}

// TestKeyBy Converts a slice to a map based on a callback function.
// func KeyBy[T any, U comparable](slice []T, iteratee func(item T) U) map[U]T
func TestKeyBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)
	{
		result1 := lo.KeyBy([]string{"a", "aa", "aaa"}, func(str string) int {
			return len(str)
		})
		is.Equal(result1, map[int]string{1: "a", 2: "aa", 3: "aaa"})
	}

	{
		result := slice.KeyBy([]string{"a", "ab", "abc"}, func(str string) int {
			return len(str)
		})
		assert.Equal(t, map[int]string{1: "a", 2: "ab", 3: "abc"}, result)
	}

}

// func StructToMap(value any) (map[string]any, error)
func TestStructToMap(t *testing.T) {
	t.Parallel()

	{
		type foo struct {
			baz string
			bar int
		}
		transform := func(f *foo) (string, int) {
			return f.baz, f.bar
		}
		testCases := []struct {
			in     []*foo
			expect map[string]int
		}{
			{
				in:     []*foo{{baz: "apple", bar: 1}},
				expect: map[string]int{"apple": 1},
			},
			{
				in:     []*foo{{baz: "apple", bar: 1}, {baz: "banana", bar: 2}},
				expect: map[string]int{"apple": 1, "banana": 2},
			},
			{
				in:     []*foo{{baz: "apple", bar: 1}, {baz: "apple", bar: 2}},
				expect: map[string]int{"apple": 2},
			},
		}
		for i, testCase := range testCases {
			t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
				is := assert.New(t)
				is.Equal(lo.Associate(testCase.in, transform), testCase.expect)
			})
		}
	}

	{
		type People struct {
			Name string `json:"name"`
			age  int
		}
		p := People{
			"test",
			100,
		}
		pm, _ := convertor.StructToMap(p)

		fmt.Println(pm)

		// Output:
		// map[name:test]
	}
}

func TestSliceToMap(t *testing.T) {
	t.Parallel()

	type foo struct {
		baz string
		bar int
	}
	transform := func(f *foo) (string, int) {
		return f.baz, f.bar
	}
	testCases := []struct {
		in     []*foo
		expect map[string]int
	}{
		{
			in:     []*foo{{baz: "apple", bar: 1}},
			expect: map[string]int{"apple": 1},
		},
		{
			in:     []*foo{{baz: "apple", bar: 1}, {baz: "banana", bar: 2}},
			expect: map[string]int{"apple": 1, "banana": 2},
		},
		{
			in:     []*foo{{baz: "apple", bar: 1}, {baz: "apple", bar: 2}},
			expect: map[string]int{"apple": 2},
		},
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			is := assert.New(t)
			is.Equal(lo.SliceToMap(testCase.in, transform), testCase.expect)
		})
	}
}

// TestFrequency Counts the frequency of each element in the slice.
// func Frequency[T comparable](slice []T) map[T]int
func TestFrequency(t *testing.T) {
	{
		strs := []string{"a", "b", "b", "c", "c", "c"}
		result := slice.Frequency(strs)

		assert.Equal(t, map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		}, result)
	}

	{
		t.Parallel()
		is := assert.New(t)

		is.Equal(map[int]int{}, lo.CountValues([]int{}))
		is.Equal(map[int]int{1: 1, 2: 1}, lo.CountValues([]int{1, 2}))
		is.Equal(map[int]int{1: 1, 2: 2}, lo.CountValues([]int{1, 2, 2}))
		is.Equal(map[string]int{"": 1, "foo": 1, "bar": 1}, lo.CountValues([]string{"foo", "bar", ""}))
		is.Equal(map[string]int{"foo": 1, "bar": 2}, lo.CountValues([]string{"foo", "bar", "bar"}))
	}
}

func TestCountValuesBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	oddEven := func(v int) bool {
		return v%2 == 0
	}
	length := func(v string) int {
		return len(v)
	}

	result1 := lo.CountValuesBy([]int{}, oddEven)
	result2 := lo.CountValuesBy([]int{1, 2}, oddEven)
	result3 := lo.CountValuesBy([]int{1, 2, 2}, oddEven)
	result4 := lo.CountValuesBy([]string{"foo", "bar", ""}, length)
	result5 := lo.CountValuesBy([]string{"foo", "bar", "bar"}, length)

	is.Equal(map[bool]int{}, result1)
	is.Equal(map[bool]int{true: 1, false: 1}, result2)
	is.Equal(map[bool]int{true: 2, false: 1}, result3)
	is.Equal(map[int]int{0: 1, 3: 2}, result4)
	is.Equal(map[int]int{3: 3}, result5)
}

func TestCombineToMap(t *testing.T) {
	keys := []string{"key0", "key1"}

	mp := arrutil.CombineToMap(keys, []int{1, 2})
	assert.Len(t, mp, 2)
	assert.Equal(t, 1, mp["key0"])
	assert.Equal(t, 2, mp["key1"])

	mp = arrutil.CombineToMap(keys, []int{1})
	assert.Len(t, mp, 1)
	assert.Equal(t, 1, mp["key0"])
}

func TestCombineToSMap(t *testing.T) {
	keys := []string{"key0", "key1"}

	mp := arrutil.CombineToSMap(keys, []string{"val0", "val1"})
	assert.Len(t, mp, 2)
	assert.Equal(t, "val0", mp["key0"])

	mp = arrutil.CombineToSMap(keys, []string{"val0"})
	assert.Len(t, mp, 2)
	assert.Equal(t, "val0", mp["key0"])
	assert.Equal(t, "", mp["key1"])
}
