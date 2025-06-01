package channel

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/petermattis/goid"
)

// 获取当前 goroutine 的 ID
func TestGetGoroutineId(t *testing.T) {
	// 设置日志输出格式
	log.SetFlags(0) // 禁用时间戳，保持输出简单
	log.SetOutput(os.Stdout)

	// 在不同的 goroutine 中打印日志
	go func() {
		log.Printf("goroutine[%d]: Hello from goroutine 1", goid.Get())
	}()

	go func() {
		log.Printf("goroutine[%d]: Hello from goroutine 2", goid.Get())
	}()

	// 主 goroutine
	log.Printf("goroutine[%d]: Hello from main goroutine", goid.Get())

	// 防止程序过快退出
	time.Sleep(100 * time.Microsecond)
}
