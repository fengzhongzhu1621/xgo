package numpy

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/mathutil"
)

func TestFibonacci(t *testing.T) {
	result1 := mathutil.Fibonacci(1, 1, 1)
	result2 := mathutil.Fibonacci(1, 1, 2)
	result3 := mathutil.Fibonacci(1, 1, 5)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)

	// Output:
	// 1
	// 1
	// 5
}

func TestFibonacci2(t *testing.T) {
	c := make(chan int)
	go fibonacci(10, c)
	for num := range c {
		fmt.Println(num)
	}
}
