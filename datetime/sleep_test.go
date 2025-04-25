package datetime

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestSleep(t *testing.T) {
	ctx := context.Background()
	dur := time.Second
	err := Sleep(ctx, dur)
	if err != nil {
		t.Errorf("Sleep error: %v", err)
	}
}

func TestSleepWithTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	dur := time.Second
	err := Sleep(ctx, dur)
	if err != nil {
		t.Errorf("Sleep error: %v", err)
	}
}

func TestSleepWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dur := time.Second
	err := Sleep(ctx, dur)
	if err != nil {
		t.Errorf("Sleep error: %v", err)
	}
}

// 周期性定时器(Ticker)
func TestNewTicker(t *testing.T) {
	ticker := time.NewTicker(1 * time.Second)

	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
			}
		}
	}()

	// 运行5秒后停止
	time.Sleep(5 * time.Second)
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")
}
