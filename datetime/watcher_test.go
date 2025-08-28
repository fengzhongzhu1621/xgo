package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/datetime"
	"github.com/duke-git/lancet/v2/function"
)

// 观察者（Watcher）用于记录代码执行时间。可以启动/停止/重置观察计时器（监测计时器）。获取函数执行的已耗时间。
// Watcher is used for record code execution time. can start/stop/reset the watch timer. get the elapsed time of function execution.
//
//	type Watcher struct {
//		startTime int64
//		stopTime  int64
//		excuting  bool
//	}
//
// func NewWatcher() *Watcher
// func (w *Watcher) Start() //start the watcher
// func (w *Watcher) Stop() //stop the watcher
// func (w *Watcher) Reset() //reset the watcher
// func (w *Watcher) GetElapsedTime() time.Duration //get the elapsed time of function execution
func TestWatcher(t *testing.T) {
	// 创建观察者
	w := function.NewWatcher()

	// 启动观察计时器
	w.Start()

	// 模拟耗时操作
	longRunningTask()
	fmt.Println(w)

	// 关闭观察计时器
	w.Stop()

	// 获得函数执行已耗时间
	elapsedTime := w.GetElapsedTime().Milliseconds()
	fmt.Println(elapsedTime)

	// 重置观察计时器s
	w.Reset()
}

// Tracks function execution time.
// func TrackFuncTime(pre time.Time) func()
func TestTrackFuncTime(t *testing.T) {
	defer datetime.TrackFuncTime(time.Now())()

	var n int
	for i := 0; i < 5000000; i++ {
		n++
	}

	fmt.Println(1) // Function main execution time:     1.460287ms
}

func longRunningTask() {
	var slice []int64
	for i := 0; i < 10000000; i++ {
		slice = append(slice, int64(i))
	}
}
