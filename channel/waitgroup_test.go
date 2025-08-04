package channel

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/remeh/sizedwaitgroup"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

// 支持等待超时
func TestWaitGroupTimeout(t *testing.T) {
	// sync.WaitGroup 是 Go 语言中的一个同步原语，用于等待一组协程（goroutine）完成执行。
	// 它提供了一种简单的方式来同步多个并发执行的协程，确保在所有协程完成之前，主协程不会退出。
	wg := &sync.WaitGroup{}

	isTimeout := WaitGroupTimeout(wg, time.Millisecond*100)
	assert.False(t, isTimeout)

	wg2 := &sync.WaitGroup{}
	// 增加等待组的计数器。delta 参数表示要增加的值。通常在启动一个新的协程之前调用此方法。
	wg2.Add(1)

	// 因为没有执行 wg.Done() // 减少等待组计数器，wg 会超时退出
	isTimeout2 := WaitGroupTimeout(wg2, time.Millisecond*100)
	assert.True(t, isTimeout2)
}

// 限制 waitgroup 的大小
func TestSizedWaitGroup(t *testing.T) {
	var wg sync.WaitGroup
	// 创建一个容量为 3 的 SizedWaitGroup, 最多只能有 3 个协程同时运行
	swg := sizedwaitgroup.New(3)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			swg.Add()        // 增加 SizedWaitGroup 的计数器
			defer swg.Done() // 减少 SizedWaitGroup 的计数器

			// 当一个协程完成任务并调用 swg.Done() 时，SizedWaitGroup 的计数器会减少，从而允许其他等待中的协程开始执行。
			// 模拟一个耗时任务
			time.Sleep(time.Second)
			fmt.Printf("Task %d completed\n", i)
		}(i)
	}

	// 还需要使用标准库中的 sync.WaitGroup 来确保主协程等待所有任务完成。
	// 这是因为 sizedwaitgroup 只负责限制并发数量，而不负责等待所有任务完成。
	wg.Wait() // 等待所有协程完成
}

// 等待所有 goroutine 完成并返回第一个错误，所有的goroutine 都会执行
// sync.WaitGroup 只负责等待 goroutine 完成，不处理 goroutine 的返回值或错误。
// errgroup.Group 虽然目前也不能直接处理 goroutine 的返回值，但在 goroutine 返回错误时，可以立即取消其他正在运行的 goroutine，并在 Wait 方法中返回第一个非 nil 的错误。
func TestErrgroup(t *testing.T) {
	urls := []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"xxx", // 这是一个错误的 URL，会导致任务失败
	}

	var g errgroup.Group

	for _, url := range urls {
		// 启动一个 goroutine 来获取 URL
		g.Go(func() error {
			resp, err := http.Get(url)
			if err != nil {
				return err // 发生错误，返回该错误
			}
			defer resp.Body.Close()
			fmt.Printf("fetch url %s status %s\n", url, resp.Status)
			return nil // 返回 nil 表示成功
		})
	}

	// 等待所有 goroutine 完成并返回第一个错误（如果有）
	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	// Output:
	// fetch url http://www.google.com/ status 200 OK
	// fetch url http://www.golang.org/ status 200 OK
	// Error: Get "xxx": unsupported protocol scheme ""
}

// 等待所有 goroutine 完成并返回第一个错误，如果有错误发生，则立即取消其他正在运行的 goroutine
// errgroup.WithContext 可以与 context.Context 配合使用，支持在某个 goroutine 出现错误时自动取消其他 goroutine，这样可以更好地控制资源，避免不必要的工作。
func TestErrGroupWithContext(t *testing.T) {
	urls := []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"xxx", // 这是一个错误的 URL，会导致任务失败
	}

	// 任何一个 goroutine 返回非 nil 的错误，或 Wait() 等待所有 goroutine 完成后，context 都会被取消
	g, ctx := errgroup.WithContext(context.Background())

	// 创建一个 map 来保存结果
	var result sync.Map

	for _, url := range urls {
		// 启动一个 goroutine 来获取 URL
		g.Go(func() error {
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				return err // 发生错误，返回该错误
			}

			// 发起请求
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return err // 发生错误，返回该错误
			}
			defer resp.Body.Close()

			// 保存每个 URL 的响应状态码
			result.Store(url, resp.Status)
			return nil // 返回 nil 表示成功
		})
	}

	// 等待所有 goroutine 完成并返回第一个错误（如果有）
	if err := g.Wait(); err != nil {
		fmt.Println("Error: ", err)
	}

	// 所有 goroutine 都执行完成，遍历并打印成功的结果
	result.Range(func(key, value any) bool {
		fmt.Printf("fetch url %s status %s\n", key, value)
		return true
	})

	// Output:
	// Error:  Get "xxx": unsupported protocol scheme ""
}

// 限制并发数量
// errgroup 提供了便捷的接口来限制并发 goroutine 的数量，避免过载，而 sync.WaitGroup 没有这样的功能。
func TestErrGroupSetLimit(t *testing.T) {
	// 创建一个 errgroup.Group
	var g errgroup.Group
	// 设置最大并发限制为 3
	g.SetLimit(3)

	// 启动 10 个 goroutine
	for i := 1; i <= 10; i++ {
		g.Go(func() error {
			// 打印正在运行的 goroutine
			fmt.Printf("Goroutine %d is starting\n", i)
			time.Sleep(2 * time.Second) // 模拟任务耗时
			fmt.Printf("Goroutine %d is done\n", i)
			return nil
		})
	}

	// 等待所有 goroutine 完成
	if err := g.Wait(); err != nil {
		fmt.Printf("Encountered an error: %v\n", err)
	}

	fmt.Println("All goroutines complete.")
}

// TestTryGo 尝试启动一个任务，它返回一个 bool 值，标识任务是否启动成功，true 表示成功，false 表示失败。
// 需要搭配 errgroup.SetLimit 一同使用，因为如果不限制并发数量，那么 errgroup.TryGo 始终返回 true，
// 注意在调用 errgroup.Go 或 errgroup.TryGo 方法前调用 errgroup.SetLimit，以防程序出现 panic
// 当达到最大并发数量限制时，errgroup.TryGo 返回 false。
func TestTryGo(t *testing.T) {
	// 创建一个 errgroup.Group
	var g errgroup.Group
	// 设置最大并发限制为 3
	g.SetLimit(3)

	// 启动 10 个 goroutine
	for i := 1; i <= 10; i++ {
		if g.TryGo(func() error {
			// 打印正在运行的 goroutine
			fmt.Printf("Goroutine %d is starting\n", i)
			time.Sleep(2 * time.Second) // 模拟工作
			fmt.Printf("Goroutine %d is done\n", i)
			return nil
		}) {
			// 如果成功启动，打印提示
			fmt.Printf("Goroutine %d started successfully\n", i)
		} else {
			// 如果达到并发限制，打印提示
			fmt.Printf("Goroutine %d waiting\n", i)
		}
	}

	// 等待所有 goroutine 完成
	if err := g.Wait(); err != nil {
		fmt.Printf("Encountered an error: %v\n", err)
	}

	fmt.Println("All goroutines complete.")
}
