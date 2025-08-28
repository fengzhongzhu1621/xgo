package cast

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/condition"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestBool(t *testing.T) {
	// bool
	result1 := condition.Bool(false)
	result2 := condition.Bool(true)
	fmt.Println(result1) // false
	fmt.Println(result2) // true

	// integer
	result3 := condition.Bool(0)
	result4 := condition.Bool(1)
	fmt.Println(result3) // false
	fmt.Println(result4) // true

	// string
	result5 := condition.Bool("")
	result6 := condition.Bool(" ")
	fmt.Println(result5) // false
	fmt.Println(result6) // true

	// slice
	nums := []int{}
	result7 := condition.Bool(nums)
	fmt.Println(result7) // false

	nums = append(nums, 1, 2)
	result8 := condition.Bool(nums)
	fmt.Println(result8) // true

	// struct
	result9 := condition.Bool(struct{}{})
	fmt.Println(result9) // false
}

// func And[T, U any](a T, b U) bool
func TestAnd(t *testing.T) {
	fmt.Println(condition.And(0, 1))   // false
	fmt.Println(condition.And(0, ""))  // false
	fmt.Println(condition.And(0, "0")) // false
	fmt.Println(condition.And(1, "0")) // true
}

// func Or[T, U any](a T, b U) bool
func TestOr(t *testing.T) {
	fmt.Println(condition.Or(0, ""))  // false
	fmt.Println(condition.Or(0, 1))   // true
	fmt.Println(condition.Or(0, "0")) // true
	fmt.Println(condition.Or(1, "0")) // true
}

// func Xor[T, U any](a T, b U) bool
func TestXor(t *testing.T) {
	fmt.Println(condition.Xor(0, 0)) // false
	fmt.Println(condition.Xor(0, 1)) // true
	fmt.Println(condition.Xor(1, 0)) // true
	fmt.Println(condition.Xor(1, 1)) // false
}

// func Nor[T, U any](a T, b U) bool
func TestNor(t *testing.T) {
	fmt.Println(condition.Nor(0, 0)) // true
	fmt.Println(condition.Nor(0, 1)) // false
	fmt.Println(condition.Nor(1, 0)) // false
	fmt.Println(condition.Nor(1, 1)) // false
}

// func Xnor[T, U any](a T, b U) bool
func TestXnor(t *testing.T) {
	fmt.Println(condition.Xnor(0, 0)) // true
	fmt.Println(condition.Xnor(0, 1)) // false
	fmt.Println(condition.Xnor(1, 0)) // false
	fmt.Println(condition.Xnor(1, 1)) // true
}

// func Nand[T, U any](a T, b U) bool
func TestNand(t *testing.T) {
	fmt.Println(condition.Nand(0, 0)) // true
	fmt.Println(condition.Nand(0, 1)) // true
	fmt.Println(condition.Nand(1, 0)) // true
	fmt.Println(condition.Nand(1, 1)) // false
}

// func ToBool(s string) (bool, error)
func TestToBool(t *testing.T) {
	cases := []string{"1", "true", "True", "false", "False", "0", "123", "0.0", "abc"}

	for i := 0; i < len(cases); i++ {
		result, _ := convertor.ToBool(cases[i])
		fmt.Println(result)
	}

	// Output:
	// true
	// true
	// true
	// false
	// false
	// false
	// false
	// false
	// false
}

// func Ternary[T, U any](isTrue T, ifValue U, elseValue U) U
func TestTernary(t *testing.T) {
	conditionTrue := 2 > 1
	result1 := condition.Ternary(conditionTrue, 0, 1)

	conditionFalse := 2 > 3
	result2 := condition.Ternary(conditionFalse, 0, 1)

	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// 0
	// 1
}

func TestTernary2(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.Ternary(true, "a", "b")
	result2 := lo.Ternary(false, "a", "b")

	is.Equal(result1, "a")
	is.Equal(result2, "b")
}

func TestTernaryF(t *testing.T) {
	is := assert.New(t)

	result1 := lo.TernaryF(true, func() string { return "a" }, func() string { return "b" })
	result2 := lo.TernaryF(false, func() string { return "a" }, func() string { return "b" })

	is.Equal(result1, "a")
	is.Equal(result2, "b")
}

func TestIfElse(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.If(true, 1).ElseIf(false, 2).Else(3)
	result2 := lo.If(true, 1).ElseIf(true, 2).Else(3)
	result3 := lo.If(false, 1).ElseIf(true, 2).Else(3)
	result4 := lo.If(false, 1).ElseIf(false, 2).Else(3)

	is.Equal(result1, 1)
	is.Equal(result2, 1)
	is.Equal(result3, 2)
	is.Equal(result4, 3)
}

func TestIfFElseF(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.IfF(true, func() int { return 1 }).
		ElseIfF(false, func() int { return 2 }).
		ElseF(func() int { return 3 })
	result2 := lo.IfF(true, func() int { return 1 }).
		ElseIfF(true, func() int { return 2 }).
		ElseF(func() int { return 3 })
	result3 := lo.IfF(false, func() int { return 1 }).
		ElseIfF(true, func() int { return 2 }).
		ElseF(func() int { return 3 })
	result4 := lo.IfF(false, func() int { return 1 }).
		ElseIfF(false, func() int { return 2 }).
		ElseF(func() int { return 3 })

	is.Equal(result1, 1)
	is.Equal(result2, 1)
	is.Equal(result3, 2)
	is.Equal(result4, 3)
}

func TestSwitchCase(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.Switch[int, int](42).Case(42, 1).Case(1, 2).Default(3)
	result2 := lo.Switch[int, int](42).Case(42, 1).Case(42, 2).Default(3)
	result3 := lo.Switch[int, int](42).Case(1, 1).Case(42, 2).Default(3)
	result4 := lo.Switch[int, int](42).Case(1, 1).Case(1, 2).Default(3)

	is.Equal(result1, 1)
	is.Equal(result2, 1)
	is.Equal(result3, 2)
	is.Equal(result4, 3)
}

func TestSwitchCaseF(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.Switch[int, int](
		42,
	).CaseF(42, func() int { return 1 }).
		CaseF(1, func() int { return 2 }).
		DefaultF(func() int { return 3 })
	result2 := lo.Switch[int, int](
		42,
	).CaseF(42, func() int { return 1 }).
		CaseF(42, func() int { return 2 }).
		DefaultF(func() int { return 3 })
	result3 := lo.Switch[int, int](
		42,
	).CaseF(1, func() int { return 1 }).
		CaseF(42, func() int { return 2 }).
		DefaultF(func() int { return 3 })
	result4 := lo.Switch[int, int](
		42,
	).CaseF(1, func() int { return 1 }).
		CaseF(1, func() int { return 2 }).
		DefaultF(func() int { return 3 })

	is.Equal(result1, 1)
	is.Equal(result2, 1)
	is.Equal(result3, 2)
	is.Equal(result4, 3)
}

func TestToBoolE(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect bool
		iserr  bool
	}{
		{0, false, false},
		{nil, false, false},
		{"false", false, false},
		{"FALSE", false, false},
		{"False", false, false},
		{"f", false, false},
		{"F", false, false},
		{false, false, false},

		{"true", true, false},
		{"TRUE", true, false},
		{"True", true, false},
		{"t", true, false},
		{"T", true, false},
		{1, true, false},
		{true, true, false},
		{-1, true, false},

		// errors
		{"test", false, true},
		{testing.T{}, false, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToBoolE(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, v, errMsg)

		// Non-E test
		v = ToBool(test.input)
		assert.Equal(t, test.expect, v, errMsg)
	}
}

func BenchmarkTooBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if !ToBool(true) {
			b.Fatal("ToBool returned false")
		}
	}
}
