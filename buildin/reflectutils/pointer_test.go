package reflectutils

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/convertor"
)

// Returns a pointer to passed value.
// func ToPointer[T any](value T) *T
func TestToPointer(t *testing.T) {
	result := convertor.ToPointer(123)
	fmt.Println(*result)

	// Output:
	// 123
}
