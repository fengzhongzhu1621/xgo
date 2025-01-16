package randutils

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/random"
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
