package pointer

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/pointer"
)

// TestOf Returns a pointer to the pass value `v`.
// func Of[T any](v T) *T
func TestOf(t *testing.T) {
	result1 := pointer.Of(123)
	result2 := pointer.Of("abc")

	fmt.Println(*result1)
	fmt.Println(*result2)

	// Output:
	// 123
	// abc
}
