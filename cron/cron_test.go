package cron

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
)

// TestNew 创建一个分钟任务
func TestNew(t *testing.T) {
	c := cron.New()

	counter := 0
	stopChan := make(chan struct{})

	// 添加一个每分钟执行一次的作业
	_, err := c.AddFunc("* * * * *", func() {
		fmt.Println("每分钟执行一次的任务")

		// 递增计数器
		counter++

		if counter == 1 {
			fmt.Println("执行 10 次，停止调度器")
			c.Stop() // 调用 c.Stop() 停止调度器
			stopChan <- struct{}{}
		}
	})
	if err != nil {
		fmt.Println("添加作业失败:", err)
		return
	}

	// 启动 cron 作业调度器
	c.Start()

	// 防止主程序立即退出
	select {
	case <-stopChan:
		break
	}
}

// TestWithSeconds 创建一个秒级任务
func TestWithSeconds(t *testing.T) {
	// 创建新的 Cron 实例，启用秒级调度（默认 cron 不支持秒级，需要额外启用）
	c := cron.New(cron.WithSeconds())

	// 任务的执行次数
	counter := 0

	stopChan := make(chan struct{})
	// 添加一个任务，每秒执行一次
	// Cron表达式 "*/1 * * * * *" 表示每秒执行一次任务
	c.AddFunc("*/1 * * * * *", func() {
		// 打印当前时间和任务执行次数
		fmt.Printf("任务执行时间: %s - 每秒执行一次任务 - 第 %d 次\n", time.Now().Format(time.DateTime), counter+1)

		// 递增计数器
		counter++

		// 在第 10 次执行时停止调度器
		if counter == 3 {
			fmt.Println("执行 10 次，停止调度器")
			c.Stop() // 调用 c.Stop() 停止调度器
			stopChan <- struct{}{}
		}
	})

	// 启动调度器
	c.Start()

	// 阻塞主线程，防止程序退出
	select {
	case <-stopChan:
		return
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 定义一个自定义任务，实现 cron.Job 接口
type MyTask struct {
	Name string
}

// 实现 Run 方法，定义任务执行的具体内容
func (t MyTask) Run() {
	fmt.Printf("任务 %s 被执行，执行时间：%s\n", t.Name, time.Now().Format(time.DateTime))
}

func TestAddJob(t *testing.T) {
	// 创建新的 Cron 实例，启用秒级调度（默认 cron 不支持秒级，需要额外启用）
	c := cron.New(cron.WithSeconds())

	// 任务的执行次数
	counter := 0

	stopChan := make(chan struct{})
	// 添加一个任务，每秒执行一次
	// Cron表达式 "*/1 * * * * *" 表示每秒执行一次任务
	c.AddFunc("*/1 * * * * *", func() {
		// 打印当前时间和任务执行次数
		fmt.Printf("任务执行时间: %s - 每秒执行一次任务 - 第 %d 次\n", time.Now().Format(time.DateTime), counter+1)

		// 递增计数器
		counter++

		// 在第 10 次执行时停止调度器
		if counter == 10 {
			fmt.Println("执行 10 次，停止调度器")
			c.Stop() // 调用 c.Stop() 停止调度器
			stopChan <- struct{}{}
		}
	})

	// 添加第二个任务，每隔 3 秒执行一次
	task := MyTask{Name: "第二个任务"}
	c.AddJob("*/3 * * * * *", task)

	// 启动调度器
	c.Start()

	// 阻塞主线程，防止程序退出
	select {
	case <-stopChan:
		return
	}
}

// TestNewParser 解析crontab字符串，并将任务添加到调度器中
func TestNewParser(t *testing.T) {
	// cron.NewParser 可以指定哪些时间字段被解析，包含秒字段的表达式需要额外支持秒字段
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

	// 解析表达式，得到一个调度计划
	spec := "*/1 * * * * *"
	schedule, err := parser.Parse(spec)
	if err != nil {
		// 如果解析失败，输出错误并退出
		fmt.Println("解析 cron 表达式出错：", err)
		return
	}

	// 计算下次任务触发时间
	now := time.Now()              // 获取当前时间
	nextTime := schedule.Next(now) // 根据当前时间计算下次执行时间
	fmt.Printf("当前时间：%s\n", now.Format(time.DateTime))
	fmt.Printf("下次任务触发时间：%s\n", nextTime.Format(time.DateTime))

	// 创建一个新的 Cron 调度器，并启用秒级调度
	// cron.WithSeconds() 选项让调度器支持秒级别的任务调度
	c := cron.New(cron.WithSeconds())

	// 任务的执行次数
	counter := 0
	stopChan := make(chan struct{})

	// 将任务添加到调度器中，定时执行
	// 调度器会根据上面解析的 schedule 来触发任务
	c.Schedule(schedule, cron.FuncJob(func() {
		// 每次任务执行时，打印当前的时间
		fmt.Printf("任务被执行，执行时间：%s\n", time.Now().Format(time.DateTime))

		// 递增计数器
		counter++

		// 在第 10 次执行时停止调度器
		if counter == 3 {
			fmt.Println("执行 10 次，停止调度器")
			c.Stop() // 调用 c.Stop() 停止调度器
			stopChan <- struct{}{}
		}
	}))

	// 启动调度器
	c.Start()

	// 阻塞主线程，防止程序退出
	select {
	case <-stopChan:
		return
	}
}

// TestLock 并发添加多个定时任务到调度器中
func TestLock(t *testing.T) {
	c := cron.New(cron.WithSeconds())

	var wg sync.WaitGroup

	// 启动调度器
	c.Start()
	fmt.Println("Cron 调度器已启动，正在添加任务...")

	// 并发添加多个定时任务到调度器中
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			// 构造 cron 表达式，每个任务的执行间隔不同
			// 例如，"*/2 * * * * *" 表示每 2 秒执行一次
			spec := fmt.Sprintf("*/%d * * * * *", i*2) // 每 i*2 秒执行一次

			// 添加任务到调度器
			id, err := c.AddFunc(spec, func() {
				fmt.Printf("任务 %d 被执行，执行时间：%s\n", i, time.Now().Format("2006-01-02 15:04:05"))
			})

			// 检查任务是否成功添加
			if err != nil {
				fmt.Printf("添加任务 %d 出错：%v\n", i, err)
			} else {
				fmt.Printf("任务 %d 已成功添加到调度器，ID：%d，执行间隔：每 %d 秒\n", i, id, i*2)
			}
		}(i)
	}

	// 等待所有任务成功添加
	wg.Wait()
	fmt.Println("所有任务已成功添加到调度器")

	// 主线程等待一段时间以观察任务执行
	time.Sleep(10 * time.Second) // 等待 20 秒，让调度器执行几轮任务

	// 调度器停止后，不会再触发新的任务，但已在执行的任务会继续完成
	c.Stop()
	fmt.Println("Cron 调度器已停止，所有任务已终止")

	// 额外等待一段时间，确保所有正在执行的任务完成
	time.Sleep(3 * time.Second)
	fmt.Println("程序已安全退出")
}

// TestLocation 测试时区
func TestLocation(t *testing.T) {
	// 加载时区信息
	nyLocation, err := time.LoadLocation("America/New_York")
	if err != nil {
		fmt.Println("加载时区失败：", err)
		return
	}

	// 创建一个使用纽约时区的 Cron 调度器
	c := cron.New(cron.WithLocation(nyLocation))

	// 添加任务，在每天的凌晨 1 点执行（纽约时间）
	c.AddFunc("0 1 * * *", func() {
		fmt.Printf("任务在纽约时间 %s 被执行\n", time.Now().In(nyLocation).Format("2006-01-02 15:04:05"))
	})

	c.Start()

	// 打印当前的本地时间和纽约时间
	fmt.Printf("本地时间：%s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("纽约时间：%s\n", time.Now().In(nyLocation).Format("2006-01-02 15:04:05"))

	// 调度器停止后，不会再触发新的任务，但已在执行的任务会继续完成
	c.Stop()
	fmt.Println("Cron 调度器已停止，所有任务已终止")

	// 额外等待一段时间，确保所有正在执行的任务完成
	time.Sleep(3 * time.Second)
	fmt.Println("程序已安全退出")
}

// TestSignal 测试信号
func TestSignal(t *testing.T) {
	c := cron.New()

	// 添加一个任务，每分钟的第 0 秒执行
	// Cron 表达式 "0 * * * *" 表示每小时的第 0 分钟执行任务
	c.AddFunc("0 * * * *", func() {
		fmt.Printf("任务被执行，执行时间：%s\n", time.Now().Format("2006-01-02 15:04:05"))
	})

	c.Start()
	fmt.Println("Cron 调度器已启动，任务将在每小时的第 0 分钟执行...")

	// 创建一个通道，用于接收系统信号 (SIGINT, SIGTERM, SIGHUP 等)
	sigs := make(chan os.Signal, 1)
	// 捕获常见的信号，如 Ctrl+C (SIGINT) 或 kill 信号 (SIGTERM)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// 等待信号到来
	select {
	case sig := <-sigs:
		fmt.Println("\n接收到信号：", sig)
		c.Stop()
		fmt.Println("Cron 调度器已停止，正在安全退出程序...")
	case <-time.After(5 * time.Second):
		fmt.Println("\n超时")
		c.Stop()
		fmt.Println("Cron 调度器已停止，正在安全退出程序...")
	}

	// 延迟退出，模拟清理工作
	time.Sleep(1 * time.Second)
	fmt.Println("程序已安全退出")
}
