package pointer

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/pointer"
)

// TestExtractPointer Returns the underlying value by the given interface type
// func ExtractPointer(value any) any
func TestExtractPointer(t *testing.T) {
	a := 1
	b := &a
	c := &b
	d := &c

	result1 := pointer.ExtractPointer(d)
	result2 := pointer.ExtractPointer(b)

	fmt.Println(result1) // 1
	fmt.Println(result2) // 1
}
