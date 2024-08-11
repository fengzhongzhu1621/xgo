package channel

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

func TestWithCancel(t *testing.T) {
	// 创建可取消的Context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("子协程收到上下文取消通知")
				return
			default:
				time.Sleep(time.Second)
			}
		}
	}()

	time.Sleep(3 * time.Second)
	// // 取消goroutine，子协成会接收到取消通知
	cancel()
	time.Sleep(2 * time.Second) // 等待 子goroutine 退出
}

func TestWithDeadline(t *testing.T) {
	// 创建带有截止时间的Context
	deadline := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("子协程收到上下文取消通知")
				return
			default:
				time.Sleep(time.Second)
			}
		}
	}()

	time.Sleep(6 * time.Second) // 等待超过截止时间
}

func TestWithTimeout(t *testing.T) {
	// 创建带有超时的 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("子协程收到上下文取消通知")
				return
			default:
				time.Sleep(time.Second)
			}
		}
	}()

	time.Sleep(6 * time.Second) // 等待超过超时时间
}

func TestWithValue(t *testing.T) {
	// 创建携带键值对的 Context
	type key int
	const userKey key = 0
	ctx := context.WithValue(context.Background(), userKey, "xxx")

	go func(ctx context.Context) {
		// 获得父协程传递的值
		if userId, ok := ctx.Value(userKey).(string); ok {
			assert.Equal(t, userId, "xxx")
		} else {
			fmt.Println("没有获得父协程传递的值")
		}
	}(ctx)

	// 等待子协程执行完毕
	time.Sleep(time.Second)
}

// 启动多个协程执行多任务
func WorkerPool(ctx context.Context, tasks <-chan int) {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case task, ok := <-tasks:
					if !ok {
						return
					}
					// 执行子任务
					processTask(ctx, task)
				case <-ctx.Done():
					// 接收到父协程取消信号
					return
				}
			}
		}()
	}
	wg.Wait()
}

// 处理任务
func processTask(ctx context.Context, task int) {
	select {
	case <-time.After(time.Second):
		// 阻塞操作直到 1 秒之后执行
		fmt.Printf("完成任务 %d\n", task)
	case <-ctx.Done():
		fmt.Printf("任务 %d 被父协程取消\n", task)
	}
}

func TestWaitGroup(t *testing.T) {
	// 创建一个 3 秒后超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 创建多个任务，任务执行完毕后关闭管道
	tasks := make(chan int, 10)
	go func() {
		for i := 0; i < 20; i++ {
			tasks <- i
		}
		close(tasks)
	}()

	// 执行多个任务
	WorkerPool(ctx, tasks)
}
