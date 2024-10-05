package ants

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/panjf2000/ants/v2"
)

const (
	DataSize    = 10000
	DataPerTask = 100
)

type Task struct {
	index int
	nums  []int
	sum   int
	wg    *sync.WaitGroup
}

// Do 将一个切片中的所有整数相加
func (t *Task) Do() {
	for _, num := range t.nums {
		t.sum += num
	}

	t.wg.Done()
}

func taskFunc(data interface{}) {
	task := data.(*Task)
	task.Do()
	fmt.Printf("task:%d sum:%d\n", task.index, task.sum)
}

func TestInvoke(t *testing.T) {
	// 创建 goroutine 池，注意池使用完后需要手动关闭
	// 第一个参数是池容量，即池中最多有 10 个 goroutine。第二个参数为每次执行任务的函数
	// 当调用 p.Invoke(data) 的时候，ants池会在其管理的 goroutine 中找出一个空闲的，让它执行函数taskFunc，并将data作为参数。
	p, _ := ants.NewPoolWithFunc(10, taskFunc)
	defer p.Release()

	// 随机生成一个 10000 长度的数组
	nums := make([]int, DataSize, DataSize)
	for i := range nums {
		nums[i] = rand.Intn(1000)
	}

	var wg sync.WaitGroup
	wg.Add(DataSize / DataPerTask)

	// 创建 100 个任务，每个任务处理 100 个数字
	tasks := make([]*Task, 0, DataSize/DataPerTask)
	for i := 0; i < DataSize/DataPerTask; i++ {
		// 每个任务处理 100 个数字的求和
		task := &Task{
			index: i + 1,
			nums:  nums[i*DataPerTask : (i+1)*DataPerTask],
			wg:    &wg,
		}
		// 发送任务到协程池
		tasks = append(tasks, task)
		p.Invoke(task)
	}

	// 等待协程池的所有任务执行完成
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", ants.Running())

	// 验证执行结果的正确性
	var sum int
	for _, task := range tasks {
		sum += task.sum
	}

	var expect int
	for _, num := range nums {
		expect += num
	}

	fmt.Printf("finish all tasks, result is %d expect:%d\n", sum, expect)
}
