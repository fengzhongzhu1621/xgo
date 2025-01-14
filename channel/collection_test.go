package channel

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectData(t *testing.T) {
	var m sync.Mutex
	numbers := []int{1, 2, 3, 4, 5}
	sum := 0
	ch := make(chan int)

	// 遍历数组，启动多个 routines
	for _, num := range numbers {
		go func(n int) {
			m.Lock()
			// 求和
			sum += n
			ch <- sum
			m.Unlock()
		}(num)
	}

	// 获得并行计算的结果
	var finalSum int
	for range numbers {
		finalSum = <-ch
		fmt.Println("result is:", finalSum)
	}
	assert.Equal(t, finalSum, 15)
}
