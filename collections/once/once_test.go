package once

import (
	"fmt"
	"sync"
	"testing"
)

var once sync.Once

func TestSyncOnce(t *testing.T) {
	once.Do(func() {
		fmt.Println("This will be printed once")
	})

	once.Do(func() {
		fmt.Println("This will not be printed")
	})

	// Output:
	// This will be printed once
}
