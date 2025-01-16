package randutils

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/random"
)

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
