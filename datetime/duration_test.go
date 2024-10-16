package datetime

import (
	"fmt"
	"testing"
)

func TestDuration(t *testing.T) {
	duration, err := Duration("1h30m")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Duration:", duration) // 输出: Duration: 1h30m0s

	duration, _ = Duration("1s")
	fmt.Println("Duration:", duration) // 输出: Duration: 1s
}
