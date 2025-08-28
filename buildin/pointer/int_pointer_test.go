package pointer

import (
	"fmt"
	"testing"
)

func swapInt(a, b *int) {
	temp := *a
	*a = *b
	*b = temp
}

func TestSwapInt(t *testing.T) {
	x := 10
	y := 20
	fmt.Println("Before swap:", x, y)
	swapInt(&x, &y)
	fmt.Println("After swap:", x, y)
}
