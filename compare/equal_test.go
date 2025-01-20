package compare

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/compare"
	"github.com/gookit/goutil/arrutil"
	"github.com/stretchr/testify/assert"
)

// TestEqual 检查两个值是否相等（同时检查类型和值）。
// func Equal(left, right any) bool
func TestEqual(t *testing.T) {
	result1 := compare.Equal(1, 1)
	result2 := compare.Equal("1", "1")
	result3 := compare.Equal([]int{1, 2, 3}, []int{1, 2, 3})
	result4 := compare.Equal(map[int]string{1: "a", 2: "b"}, map[int]string{1: "a", 2: "b"})

	result5 := compare.Equal(1, "1")
	result6 := compare.Equal(1, int64(1))
	result7 := compare.Equal([]int{1, 2}, []int{1, 2, 3})

	fmt.Println(result1) // true
	fmt.Println(result2) // true
	fmt.Println(result3) // true
	fmt.Println(result4) // true
	fmt.Println(result5) // false
	fmt.Println(result6) // false
	fmt.Println(result7) // false
}

// StringEqualComparer tests
func TestStringEqualsComparer(t *testing.T) {
	assert.Equal(t, 0, arrutil.StringEqualsComparer("a", "a"))
	assert.Equal(t, -1, arrutil.StringEqualsComparer("a", "b"))
}

func TestValueEqualsComparer(t *testing.T) {
	assert.Equal(t, 0, arrutil.ValueEqualsComparer("1", "1"))
	assert.Equal(t, -1, arrutil.ValueEqualsComparer(1, 2))
}

// ReflectEqualsComparer tests
func TestReflectEqualsComparer(t *testing.T) {
	assert.Equal(t, 0, arrutil.ReflectEqualsComparer(1, 1))
	assert.Equal(t, -1, arrutil.ReflectEqualsComparer(1, 2))
}

// ElemTypeEqualCompareFunc
func TestElemTypeEqualCompareFuncShouldEquals(t *testing.T) {
	var c = 1
	assert.Equal(t, 0, arrutil.ElemTypeEqualsComparer(c, c))
	assert.Equal(t, 0, arrutil.ElemTypeEqualsComparer(1, 1))

	var a, b any
	a = 1
	b = "2"
	assert.Equal(t, -1, arrutil.ElemTypeEqualsComparer(a, b))
}
