package once

import (
	"fmt"
	"sync"
	"testing"
)

var wrapper2 = sync.OnceValue(printOnce2)

func printOnce2() int {
	fmt.Println("This will be printed once")
	return 1
}

func TestSyncOnceOne(t *testing.T) {
	value1 := wrapper2()
	value2 := wrapper2()

	fmt.Println("value1 = ", value1)
	fmt.Println("value2 = ", value2)

	// Output:
	// This will be printed once
	// value1 =  1
	// value2 =  1
}
