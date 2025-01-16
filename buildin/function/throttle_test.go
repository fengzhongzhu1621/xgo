package function

// 节流（Throttle）创建一个经过节流处理的（所提供的）函数版本。返回的函数保证每个时间间隔最多只会被调用一次。
//func Throttle(fn func(), interval time.Duration) func()
import (
	"fmt"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/function"
)

func TestThrottle(t *testing.T) {
	callCount := 0

	fn := func() {
		callCount++
	}

	// 一秒最多执行一次
	throttledFn := function.Throttle(fn, 1*time.Second)

	for i := 0; i < 5; i++ {
		throttledFn()
	}

	time.Sleep(1 * time.Second)
	fmt.Println(callCount)

	// Output:
	// 1
}
