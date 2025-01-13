package function

import (
	"fmt"
	"strings"
	"testing"

	"github.com/duke-git/lancet/v2/function"
)

// 返回一个组合谓词，该谓词表示一系列谓词的逻辑与（AND）关系。只有当对于给定值所有谓词的计算结果都为真时，该组合谓词的计算结果才为真。
// func And[T any](predicates ...func(T) bool) func(T) bool
func TestAnd(t *testing.T) {
	isNumericAndLength5 := function.And(
		func(s string) bool { return strings.ContainsAny(s, "0123456789") },
		func(s string) bool { return len(s) == 5 },
	)

	fmt.Println(isNumericAndLength5("12345"))
	fmt.Println(isNumericAndLength5("1234"))
	fmt.Println(isNumericAndLength5("abcde"))

	// Output:
	// true
	// false
	// false
}

// func Or[T any](predicates ...func(T) bool) func(T) bool
func TestOr(t *testing.T) {
	containsDigitOrSpecialChar := function.Or(
		func(s string) bool { return strings.ContainsAny(s, "0123456789") },
		func(s string) bool { return strings.ContainsAny(s, "!@#$%") },
	)

	fmt.Println(containsDigitOrSpecialChar("hello!"))
	fmt.Println(containsDigitOrSpecialChar("hello"))

	// Output:
	// true
	// false
}

// 返回一个代表此谓词逻辑非（NOT）关系的谓词。
// func Negate[T any](predicate func(T) bool) func(T) bool
func TestNegate(t *testing.T) {
	// Define some simple predicates for demonstration
	isUpperCase := func(s string) bool {
		return strings.ToUpper(s) == s
	}
	isLowerCase := func(s string) bool {
		return strings.ToLower(s) == s
	}
	isMixedCase := function.Negate(function.Or(isUpperCase, isLowerCase))

	fmt.Println(isMixedCase("ABC"))
	fmt.Println(isMixedCase("AbC"))

	// Output:
	// false
	// true
}

// 返回一个组合谓词，该谓词表示一系列谓词的逻辑或非（NOR）关系。只有当对于给定值所有谓词的计算结果都为假时，该组合谓词的计算结果才为真。
// func Nor[T any](predicates ...func(T) bool) func(T) bool
func TestNor(t *testing.T) {
	match := function.Nor(
		func(s string) bool { return strings.ContainsAny(s, "0123456789") },
		func(s string) bool { return len(s) == 5 },
	)

	fmt.Println(match("dbcdckkeee"))

	match = function.Nor(
		func(s string) bool { return strings.ContainsAny(s, "0123456789") },
		func(s string) bool { return len(s) == 5 },
	)

	fmt.Println(match("0123456789"))

	// Output:
	// true
	// false
}

// func Nand[T any](predicates ...func(T) bool) func(T) bool
func TestNand(t *testing.T) {
	isNumericAndLength5 := function.Nand(
		func(s string) bool { return strings.ContainsAny(s, "0123456789") },
		func(s string) bool { return len(s) == 5 },
	)

	fmt.Println(isNumericAndLength5("12345"))
	fmt.Println(isNumericAndLength5("1234"))
	fmt.Println(isNumericAndLength5("abcdef"))

	// Output:
	// false
	// false
	// true
}

// func Xnor[T any](predicates ...func(T) bool) func(T) bool
func TestXnor(t *testing.T) {
	isEven := func(i int) bool { return i%2 == 0 }
	isPositive := func(i int) bool { return i > 0 }

	match := function.Xnor(isEven, isPositive)

	fmt.Println(match(2))
	fmt.Println(match(-3))
	fmt.Println(match(3))

	// Output:
	// true
	// true
	// false
}
