package gocron

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func TestNewJob(t *testing.T) {
	s, err := gocron.NewScheduler()
	if err != nil {
		panic("error")
	}

	// 创建10秒的定时任务
	j, err := s.NewJob(
		gocron.DurationJob(10*time.Second),
		// 定义任务执行函数
		gocron.NewTask(func(a string, b int) {
			fmt.Println(a, b)
		},
			"hello", 1,
		),
	)
	if err != nil {
		panic("error")
	}

	// 启动定时任务
	fmt.Println(j.ID())
	s.Start()

	// when you're done, shut it down
	err = s.Shutdown()
	if err != nil {
		panic("error")
	}
}
