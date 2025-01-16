package randutils

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/random"
)

// func RandBytes(length int) []byte
func TestRandBytes(t *testing.T) {
	randBytes := random.RandBytes(4)
	fmt.Println(randBytes)
}

// TestRandFromGivenSlice Generate a random slice of length num from given slice.
// func RandSliceFromGivenSlice[T any](slice []T, num int, repeatable bool) []T
func TestRandSliceFromGivenSlice(t *testing.T) {
	goods := []string{"apple", "banana", "cherry", "elderberry", "fig", "grape", "honeydew", "kiwi", "lemon", "mango", "nectarine", "orange"}

	chosen3goods := random.RandSliceFromGivenSlice(goods, 3, false)

	fmt.Println(chosen3goods)
}
