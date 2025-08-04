package slice

import (
	"reflect"
	"testing"

	"github.com/araujo88/lambda-go/pkg/utils"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gookit/goutil/arrutil"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestReverse(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		result1 := lo.Reverse([]int{0, 1, 2, 3, 4, 5})
		result2 := lo.Reverse([]int{0, 1, 2, 3, 4, 5, 6})
		result3 := lo.Reverse([]int{})

		is.Equal(result1, []int{5, 4, 3, 2, 1, 0})
		is.Equal(result2, []int{6, 5, 4, 3, 2, 1, 0})
		is.Equal(result3, []int{})

		type myStrings []string
		allStrings := myStrings{"", "foo", "bar"}
		nonempty := lo.Reverse(allStrings)
		is.IsType(nonempty, allStrings, "type preserved")
	}

	{
		strs := []string{"a", "b", "c", "d"}

		slice.Reverse(strs)
		assert.Equal(t, []string{"d", "c", "b", "a"}, strs)
	}

	{
		slice1 := []int{1, 2, 3}
		reversed := utils.Reverse(slice1)
		assert.Equal(t, []int{3, 2, 1}, reversed)

		tests := []struct {
			name  string
			slice []int
			want  []int
		}{
			{"reverse non-empty", []int{1, 2, 3}, []int{3, 2, 1}},
			{"reverse empty", []int{}, []int{}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := utils.Reverse(tt.slice); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Reverse() = %v, want %v", got, tt.want)
				}
			})
		}
	}

	{
		ss := []string{"a", "b", "c"}
		arrutil.Reverse(ss)
		assert.Equal(t, []string{"c", "b", "a"}, ss)

		ints := []int{1, 2, 3}
		arrutil.Reverse(ints)
		assert.Equal(t, []int{3, 2, 1}, ints)
	}
}

func TestReflectReverseSlice(t *testing.T) {
	names := []string{"a", "b", "c", "d", "e", "f", "g"}
	ReflectReverseSlice(names)
	expected := []string{"g", "f", "e", "d", "c", "b", "a"}
	assert.Equal(t, expected, names)
}

func TestReverseSliceGetNew(t *testing.T) {
	names := []string{"a", "b", "c", "d", "e", "f", "g"}
	newNames := ReverseSliceGetNew(names)
	expected := []string{"g", "f", "e", "d", "c", "b", "a"}
	assert.Equal(t, expected, newNames)
}

func TestReverseSlice(t *testing.T) {
	names := []string{"a", "b", "c", "d", "e", "f", "g"}
	ReverseSlice(names)
	expected := []string{"g", "f", "e", "d", "c", "b", "a"}
	assert.Equal(t, expected, names)
}

func BenchmarkReverseReflectSlice(b *testing.B) {
	names := []string{"a", "b", "c", "d", "e", "f", "g"}
	for i := 0; i < b.N; i++ {
		ReflectReverseSlice(names)
	}
}

func BenchmarkReverseSlice(b *testing.B) {
	names := []string{"a", "b", "c", "d", "e", "f", "g"}
	for i := 0; i < b.N; i++ {
		ReverseSlice(names)
	}
}

func BenchmarkReverseSliceNew(b *testing.B) {
	names := []string{"a", "b", "c", "d", "e", "f", "g"}
	for i := 0; i < b.N; i++ {
		ReverseSliceGetNew(names)
	}
}
