package compare

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type S1 struct {
	Field int
}
type S2 struct {
	Field int
}

type S struct {
	Field1 int
	field2 string
}

func TestCompareStringSliceReflect(t *testing.T) {
	// 不同类型的值不会深度相等
	assert.Equal(t, reflect.DeepEqual(S1{1}, S2{1}), false)

	// 当两个数组的元素对应深度相等时，两个数组深度相等
	Array1 := []string{"hello1", "hello2"}
	Array2 := []string{"hello1", "hello2"}
	assert.Equal(t, reflect.DeepEqual(Array1, Array2), true)

	// 当两个相同结构体的所有字段对应深度相等的时候，两个结构体深度相等
	s1 := S{Field1: 1, field2: "hello"}
	s2 := S{Field1: 1, field2: "hello"}
	assert.Equal(t, reflect.DeepEqual(s1, s2), true)

	// 当两个函数都为nil时，两个函数深度相等，其他情况不相等（相同函数也不相等）
	f1 := func(a int) int {
		return a * 2
	}
	assert.Equal(t, reflect.DeepEqual(f1, f1), false)
	f1 = nil
	assert.Equal(t, reflect.DeepEqual(f1, f1), true)

	// 当两个interface的真实值深度相等时，两个interface深度相等
	var i1 interface{} = "hello"
	var i2 interface{} = "hello"
	assert.Equal(t, reflect.DeepEqual(i1, i2), true)

	// go中map的比较需要同时满足以下几个
	// 1.两个map都为nil或者都不为nil，并且长度要相等
	// 2.相同的map对象或者所有key要对应相同
	// 3.map对应的value也要深度相等
	m1 := map[string]int{
		"a": 1,
		"b": 2,
	}
	m2 := map[string]int{
		"a": 1,
		"b": 2,
	}
	assert.Equal(t, reflect.DeepEqual(m1, m2), true)

	// 指针，满足以下其一即是深度相等
	// 1.两个指针满足go的==操作符
	// 2.两个指针指向的值是深度相等的
	m3 := map[string]int{
		"a": 1,
		"b": 2,
	}
	m4 := map[string]int{
		"a": 1,
		"b": 2,
	}
	M3 := &m3
	M4 := &m4
	assert.Equal(t, reflect.DeepEqual(M3, M4), true)

	// 切片，需要同时满足以下几点才是深度相等
	// 1.两个切片都为nil或者都不为nil，并且长度要相等
	// 2.两个切片底层数据指向的第一个位置要相同或者底层的元素要深度相等
	// 空的切片跟nil切片是不深度相等的
	s3 := []int{1, 2, 3, 4, 5}
	s4 := s3[0:3]
	s5 := s3[0:3]
	assert.Equal(t, reflect.DeepEqual(s4, s5), true)
	s6 := s3[1:4]
	assert.Equal(t, reflect.DeepEqual(s4, s6), false)
	s7 := []byte{}
	s8 := []byte(nil)
	assert.Equal(t, reflect.DeepEqual(s7, s8), false)

	// 其他类型的值（numbers, bools, strings, channels）如果满足go的==操作符，则是深度相等的
	// 要注意不是所有的值都深度相等于自己，例如函数，以及嵌套包含这些值的结构体，数组等
}

func BenchmarkCompareStringSliceReflect(b *testing.B) {
	sliceA := []string{"a", "b", "c", "d", "e"}
	sliceB := []string{"e", "d", "c", "b", "a"}
	for n := 0; n < b.N; n++ {
		CompareStringSliceReflect(sliceA, sliceB)
	}
}

func BenchmarkCompareStringSlice(b *testing.B) {
	sliceA := []string{"a", "b", "c", "d", "e"}
	sliceB := []string{"e", "d", "c", "b", "a"}
	for n := 0; n < b.N; n++ {
		CompareStringSlice(sliceA, sliceB)
	}
}
