package cast

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/convertor"

	"github.com/duke-git/lancet/v2/condition"
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

// func ToBytes(data any) ([]byte, error)
func TestToBytes(t *testing.T) {
	bytesData, err := convertor.ToBytes("abc")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(bytesData)

	// Output:
	// [97 98 99]
}
