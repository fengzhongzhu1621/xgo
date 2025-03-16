package db

import "context"

// Filter condition alias name
type Filter interface{}

type FindOpts struct {
	WithObjectID *bool
	WithCount    *bool
}

func NewFindOpts() *FindOpts {
	return &FindOpts{}
}

// SetWithObjectID TODO
func (f *FindOpts) SetWithObjectID(bl bool) *FindOpts {
	f.WithObjectID = &bl
	return f
}

// SetWithCount TODO
func (f *FindOpts) SetWithCount(bl bool) *FindOpts {
	f.WithCount = &bl
	return f
}

// Find find operation interface
type Find interface {
	// Fields 设置查询字段
	Fields(fields ...string) Find
	// Sort 设置查询排序
	Sort(sort string) Find
	// Start 设置限制查询上标
	Start(start uint64) Find
	// Limit 设置查询数量
	Limit(limit uint64) Find
	// All 查询多个
	All(ctx context.Context, result interface{}) error
	// One 查询单个
	One(ctx context.Context, result interface{}) error
	// Count 统计数量(非事务)
	Count(ctx context.Context) (uint64, error)
	// List 查询多个, start 等于0的时候，返回满足条件的行数
	List(ctx context.Context, result interface{}) (int64, error)

	Option(opts ...*FindOpts)
}
