package scheduler

import "context"

// Task 表示一个可调度的任务，包含任务的基本信息和执行逻辑。
type Task struct {
	ID       string                    // 任务的唯一标识符（如 "task-1"），用于在调度器中区分不同任务
	Name     string                    // 任务的名称（如 "Daily Report Generator"），便于人类阅读
	Schedule string                    // Cron 表达式（如 "0 0 * * *"），定义任务的执行时间
	JobFunc  func(ctx context.Context) // 任务的具体执行逻辑（函数），接收 context.Context 参数
	Retry    int                       // 任务失败时的重试次数（如 3 表示最多重试 3 次）
}
