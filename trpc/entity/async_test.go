package entity

import (
	"bytes"
	"context"
	"log"
	"testing"
	"time"
)

// 测试正常执行的任务
func TestDo_NormalExecution(t *testing.T) {
	done := make(chan struct{})
	fn := func(ctx context.Context) {
		defer close(done)
		time.Sleep(50 * time.Millisecond) // 模拟任务耗时
	}

	timeout := 200 * time.Millisecond
	go Do(fn, timeout)

	select {
	case <-done:
		// 任务正常完成
	case <-time.After(timeout + 50*time.Millisecond):
		t.Errorf("任务超时")
	}
}

// 测试 panic 场景
func TestDo_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("测试函数自身 panic: %v", r)
		}
	}()

	panicFn := func(ctx context.Context) {
		panic("test panic")
	}

	var buf bytes.Buffer
	log.SetOutput(&buf) // 重定向日志到内存缓冲区

	go Do(panicFn, 100*time.Millisecond)

	time.Sleep(200 * time.Millisecond) // 等待任务完成
}
