package randutils

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/random"
	"github.com/duke-git/lancet/v2/slice"
)

// TestRandBool Generate a random boolean value (true or false).
// func RandBool() bool
func TestRandBool(t *testing.T) {
	result := random.RandBool()
	fmt.Println(result) // true or false (random)
}

// TestRandBoolSlice Generates a random boolean slice of specified length.
// func RandBoolSlice(length int) []bool
func TestRandBoolSlice(t *testing.T) {
	result := random.RandBoolSlice(2)
	fmt.Println(result) // [true false] (random)
}

// 生成随机int, 范围[min, max)
// Generate random int between min and max, may contain min, not max.
// func RandInt(min, max int) int
func TestRandomRandInt(t *testing.T) {
	rInt := random.RandInt(1, 10)
	fmt.Println(rInt)
}

// TestRandIntSlice Generate a slice of random int. Number range in [min, max)
// func RandIntSlice(length, min, max int) []int
func TestRandIntSlice(t *testing.T) {
	result := random.RandIntSlice(5, 0, 10)
	fmt.Println(result) //[1 4 7 1 5] (random)
}

// TestRandUniqueIntSlice Generate a slice of random int of length that do not repeat. Number range in [min, max)
// func RandIntSlice(length, min, max int) []int
func TestRandUniqueIntSlice(t *testing.T) {
	result := random.RandUniqueIntSlice(5, 0, 10)
	fmt.Println(result) //[0 4 7 1 5] (random)
}

// TestRandNumberOfLength Generates a random int number of specified length.
// func RandNumberOfLength(len int) int
func TestRandNumberOfLength(t *testing.T) {
	i := random.RandNumberOfLength(2)
	fmt.Println(i) // 21 (random number of length 2)
}

// TestRandFloat Generate a random float64 number between [min, max) with specific precision.
// func RandFloat(min, max float64, precision int) float64
func TestRandFloat(t *testing.T) {
	floatNumber := random.RandFloat(1.0, 5.0, 2)
	fmt.Println(floatNumber) //2.14 (random number)
}

// TestRandFloats Generate a slice of random float64 numbers of length n that do not repeat. Number range in [min, max)
// func RandFloats(length int, min, max float64, precision int) []float64
func TestRandFloats(t *testing.T) {
	floatNumbers := random.RandFloats(5, 1.0, 5.0, 2)
	fmt.Println(floatNumbers) //[3.42 3.99 1.3 2.38 4.23] (random)
}

// Generate random given length string. only contains letter (a-zA-Z)
// func RandString(length int) string
func TestRandomRandString(t *testing.T) {
	randStr := random.RandString(6)
	fmt.Println(randStr) // ACMbJt
}

// TestRandStringSlice Generate a slice of random string of length strLen based on charset.
// chartset should be one of the following:
// random.Numeral, random.LowwerLetters, random.UpperLetters random.Letters, random.SymbolChars, random.AllChars. or a combination of them.
// func RandStringSlice(charset string, sliceLen, strLen int) []string
func TestRandStringSlice(t *testing.T) {
	strs := random.RandStringSlice(random.Letters, 4, 6)
	fmt.Println(strs)

	// output random string slice like below:
	//[CooSMq RUFjDz FAeMPf heRyGv]
}

// TestRandString Generate a random upper case string
// func RandUpper(length int) string
func TestRandUpper(t *testing.T) {
	randStr := random.RandUpper(6)
	fmt.Println(randStr)
}

// TestRandLower Generate a random lower case string
// func RandLower(length int) string
func TestRandLower(t *testing.T) {
	randStr := random.RandLower(6)
	fmt.Println(randStr)
}

// TestRandNumeral Generate a random numeral string
// func RandNumeral(length int) string
func TestRandNumeral(t *testing.T) {
	randStr := random.RandNumeral(6)
	fmt.Println(randStr)
}

// TestRandNumeralOrLetter generate a random numeral or letter string
// func RandNumeralOrLetter(length int) string
func TestRandNumeralOrLetter(t *testing.T) {
	randStr := random.RandNumeralOrLetter(6)
	fmt.Println(randStr) //0aW7cQ
}

// TestRandSymbolChar Generate a random symbol char of specified length.
// func RandSymbolChar(length int) string
func TestRandSymbolChar(t *testing.T) {
	randStr := random.RandSymbolChar(6)
	fmt.Println(randStr) // random string like: @#(_")
}

// TestSliceRandom Random get a random item of slice, return idx=-1 when slice is empty.
// func Random[T any](slice []T) (val T, idx int)
func TestSliceRandom(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}

	val, idx := slice.Random(nums)
	if idx >= 0 && idx < len(nums) && slice.Contain(nums, val) {
		fmt.Println("okk")
	}
}

// func RandBytes(length int) []byte
func TestRandBytes(t *testing.T) {
	randBytes := random.RandBytes(4)
	fmt.Println(randBytes)
}

// TestRandFromGivenSlice Generate a random element from given slice.
// func RandFromGivenSlice[T any](slice []T) T
func TestRandFromGivenSlice(t *testing.T) {
	randomSet := []any{"a", 8, "hello", true, 1.1}
	randElm := random.RandFromGivenSlice(randomSet)
	fmt.Println(randElm)
}

// TestRandFromGivenSlice Generate a random slice of length num from given slice.
// func RandSliceFromGivenSlice[T any](slice []T, num int, repeatable bool) []T
func TestRandSliceFromGivenSlice(t *testing.T) {
	goods := []string{"apple", "banana", "cherry", "elderberry", "fig", "grape", "honeydew", "kiwi", "lemon", "mango", "nectarine", "orange"}

	chosen3goods := random.RandSliceFromGivenSlice(goods, 3, false)

	fmt.Println(chosen3goods)
}
