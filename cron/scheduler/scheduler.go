package scheduler

import (
	"context"
	"sync"

	"github.com/robfig/cron/v3"
)

// Scheduler 表示任务调度器，负责管理多个 Task 的调度和执行。
type Scheduler struct {
	c      *cron.Cron              // 底层使用的 Cron 调度器
	mu     sync.Mutex              // 互斥锁，用于保护并发访问（如 tasks、funcs、status 的修改）
	tasks  map[string]cron.EntryID // 存储任务 ID 到 Cron EntryID 的映射（用于取消或查找任务）
	funcs  map[string]*Task        // 存储任务 ID 到 Task 结构体的映射（用于获取任务详情）
	status map[string]string       // 存储任务的状态（如 "running"、"pending"、"failed"）
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		c:      cron.New(cron.WithSeconds()), // 创建新的 Cron 实例，启用秒级调度（默认 cron 不支持秒级，需要额外启用）
		tasks:  make(map[string]cron.EntryID),
		funcs:  make(map[string]*Task),
		status: make(map[string]string),
	}
}

// Start 启动调度器
func (s *Scheduler) Start() {
	s.c.Start()
}

// Stop 停止调度器
// 调度器停止后，不会再触发新的任务，但已在执行的任务会继续完成
func (s *Scheduler) Stop() {
	s.c.Stop()
}

// AddTask 添加定时任务
func (s *Scheduler) AddTask(t *Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	id, err := s.c.AddFunc(t.Schedule, func() {
		ctx := context.Background()
		t.JobFunc(ctx)
	})
	if err != nil {
		return err
	}

	s.tasks[t.ID] = id
	s.funcs[t.ID] = t
	s.status[t.ID] = scheduleTaskRunning

	return nil
}

// RunNow 立即执行指定的任务（一次性执行 once）
func (s *Scheduler) RunNow(taskID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if t, ok := s.funcs[taskID]; ok {
		go t.JobFunc(context.Background())
	}
}

// PauseTask 暂停指定的任务
func (s *Scheduler) PauseTask(taskID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if id, ok := s.tasks[taskID]; ok {
		s.c.Remove(id)
		s.status[taskID] = scheduleTaskPaused
	}
}

// ResumeTask 将使用PauseTask暂停的任务进行恢复
func (s *Scheduler) ResumeTask(taskID string) error {
	if s.Status(taskID) != scheduleTaskPaused {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if t, ok := s.funcs[taskID]; ok {
		return s.AddTask(t)
	}

	return nil
}

func (s *Scheduler) Status(taskID string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.status[taskID]
}
