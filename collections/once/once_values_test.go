package once

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

var wrapper3 = sync.OnceValues(printOnce3)

func printOnce3() (int, error) {
	fmt.Println("This will be printed once")
	return 1, errors.New("printOnce3 error")
}

func TestSyncOnceValues(t *testing.T) {
	value1, err1 := wrapper3()
	value2, err2 := wrapper3()

	fmt.Println("value1 = ", value1)
	fmt.Println("err1 = ", err1)

	fmt.Println("value2 = ", value2)
	fmt.Println("err2 = ", err2)

	// Output:
	// This will be printed once
	// value1 =  1
	// err1 =  printOnce3 error
	// value2 =  1
	// err2 =  printOnce3 error
}
