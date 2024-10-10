package cron

import (
	"fmt"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
)

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
