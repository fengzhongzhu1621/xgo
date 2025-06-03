package once

import (
	"fmt"
	"sync"
	"testing"
)

var wrapper = sync.OnceFunc(printOnce)

func printOnce() {
	fmt.Println("This will be printed once")
}

func TestSyncOnceFunc(t *testing.T) {
	wrapper()
	wrapper()

	// Output:
	// This will be printed once
}
