package channel

import "context"

// ///////////////////////////////////////////////////////////////////////////
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

// ///////////////////////////////////////////////////////////////////////////
type Context struct {
	ctx context.Context
}

type IContextInterface interface {
	WithCancel() (context.Context, context.CancelFunc)
}

func NewContext() IContextInterface {
	return &Context{
		ctx: context.Background(),
	}
}

func (c *Context) WithCancel() (context.Context, context.CancelFunc) {
	return context.WithCancel(c.ctx)
}

// ///////////////////////////////////////////////////////////////////////////
// NewContextWithValues will use the valuesCtx's Value function.
// Effects of the returned context:
//
//	Whether it has timed out or canceled: decided by ctx.
//	Retrieve value using key: first use valuesCtx.Value, then ctx.Value.
//
// 创建一个新的 valueCtx，将传入的 ctx 和 valuesCtx 合并。新 context 的超时和取消行为由 ctx 决定，但值的查找会优先从 valuesCtx 开始。
func NewContextWithValues(ctx, valuesCtx context.Context) context.Context {
	return &valueCtx{Context: ctx, values: valuesCtx}
}

type valueCtx struct {
	context.Context
	values context.Context
}

// Value re-implements context.Value, valueCtx.values.Value has the highest priority.
// valuesCtx 的值优先级高于 ctx，适合需要动态覆盖值的场景。
func (c *valueCtx) Value(key interface{}) interface{} {
	if v := c.values.Value(key); v != nil {
		return v
	}
	return c.Context.Value(key)
}
