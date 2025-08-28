package function

import (
	"fmt"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/function"
)

// Invoke function every duration time, until close the returned bool chan.
// func Schedule(d time.Duration, fn any, args ...any) chan bool
func TestSchedule(t *testing.T) {
	count := 0

	increase := func() {
		count++
	}

	stop := function.Schedule(2*time.Second, increase)

	time.Sleep(2 * time.Second)
	close(stop)

	fmt.Println(count)

	// Output:
	// 2
}
