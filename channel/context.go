package channel

import "context"

type priorityKey struct{}

// WithPriority 传递额外的上下文信息（优先级）
func WithPriority(ctx context.Context, priority int) context.Context {
	return context.WithValue(ctx, priorityKey{}, priority)
}

// GetPriority 从上下文获取优先级
func GetPriority(ctx context.Context) (int, bool) {
	priority, ok := ctx.Value(priorityKey{}).(int)
	return priority, ok
}