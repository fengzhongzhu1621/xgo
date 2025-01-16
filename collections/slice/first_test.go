package slice

import (
	"reflect"
	"testing"

	"github.com/araujo88/lambda-go/pkg/utils"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestHead(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		want  int
	}{
		{"non-empty slice", []int{1, 2, 3}, 1},
		{"single element", []int{5}, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.Head(tt.slice); got != tt.want {
				t.Errorf("Head() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTail(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		want  []int
	}{
		{"non-empty slice", []int{1, 2, 3}, []int{2, 3}},
		{"single element", []int{5}, []int{}},
		{"empty slice", []int{}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.Tail(tt.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFirst(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1, ok1 := lo.First([]int{1, 2, 3})
	result2, ok2 := lo.First([]int{})

	is.Equal(result1, 1)
	is.Equal(ok1, true)
	is.Equal(result2, 0)
	is.Equal(ok2, false)

}

func TestFirstOrEmpty(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.FirstOrEmpty([]int{1, 2, 3})
	result2 := lo.FirstOrEmpty([]int{})
	result3 := lo.FirstOrEmpty([]string{})

	is.Equal(result1, 1)
	is.Equal(result2, 0)
	is.Equal(result3, "")
}

func TestFirstOr(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.FirstOr([]int{1, 2, 3}, 63)
	result2 := lo.FirstOr([]int{}, 23)
	result3 := lo.FirstOr([]string{}, "test")

	is.Equal(result1, 1)
	is.Equal(result2, 23)
	is.Equal(result3, "test")
}

// 找到第一个非空值
func TestCoalesce(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	newStr := func(v string) *string { return &v }
	var nilStr *string
	str1 := newStr("str1")
	str2 := newStr("str2")

	type structType struct {
		field1 int
		field2 float64
	}
	var zeroStruct structType
	struct1 := structType{1, 1.0}
	struct2 := structType{2, 2.0}

	result1, ok1 := lo.Coalesce[int]()
	result2, ok2 := lo.Coalesce(3)
	result3, ok3 := lo.Coalesce(nil, nilStr)
	result4, ok4 := lo.Coalesce(nilStr, str1)
	result5, ok5 := lo.Coalesce(nilStr, str1, str2)
	result6, ok6 := lo.Coalesce(str1, str2, nilStr)
	result7, ok7 := lo.Coalesce(0, 1, 2, 3)
	result8, ok8 := lo.Coalesce(zeroStruct)
	result9, ok9 := lo.Coalesce(zeroStruct, struct1)
	result10, ok10 := lo.Coalesce(zeroStruct, struct1, struct2)

	is.Equal(0, result1)
	is.False(ok1)

	is.Equal(3, result2)
	is.True(ok2)

	is.Nil(result3)
	is.False(ok3)

	is.Equal(str1, result4)
	is.True(ok4)

	is.Equal(str1, result5)
	is.True(ok5)

	is.Equal(str1, result6)
	is.True(ok6)

	is.Equal(result7, 1)
	is.True(ok7)

	is.Equal(result8, zeroStruct)
	is.False(ok8)

	is.Equal(result9, struct1)
	is.True(ok9)

	is.Equal(result10, struct1)
	is.True(ok10)
}

// 尝试获得第一个非零的值，如果没有则返回零值
func TestCoalesceOrEmpty(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	newStr := func(v string) *string { return &v }
	var nilStr *string
	str1 := newStr("str1")
	str2 := newStr("str2")

	type structType struct {
		field1 int
		field2 float64
	}
	var zeroStruct structType
	struct1 := structType{1, 1.0}
	struct2 := structType{2, 2.0}

	result1 := lo.CoalesceOrEmpty[int]()
	result2 := lo.CoalesceOrEmpty(3)
	result3 := lo.CoalesceOrEmpty(nil, nilStr)
	result4 := lo.CoalesceOrEmpty(nilStr, str1)
	result5 := lo.CoalesceOrEmpty(nilStr, str1, str2)
	result6 := lo.CoalesceOrEmpty(str1, str2, nilStr)
	result7 := lo.CoalesceOrEmpty(0, 1, 2, 3)
	result8 := lo.CoalesceOrEmpty(zeroStruct)
	result9 := lo.CoalesceOrEmpty(zeroStruct, struct1)
	result10 := lo.CoalesceOrEmpty(zeroStruct, struct1, struct2)

	is.Equal(0, result1)
	is.Equal(3, result2)
	is.Nil(result3)
	is.Equal(str1, result4)
	is.Equal(str1, result5)
	is.Equal(str1, result6)
	is.Equal(result7, 1)
	is.Equal(result8, zeroStruct)
	is.Equal(result9, struct1)
	is.Equal(result10, struct1)
}
