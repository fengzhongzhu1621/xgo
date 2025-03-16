package db

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrTransactionStated   = errors.New("transaction already started")
	ErrTransactionNotFound = errors.New("not in transaction environment")
	ErrDocumentNotFound    = errors.New("document not found")
	ErrDuplicated          = errors.New("duplicated")
	ErrSessionNotStarted   = errors.New("session is not started")

	UpdateOpAddToSet = "addToSet"
	UpdateOpPull     = "pull"
)

// Index define the DB index struct
type Index struct {
	Keys                    bson.D                 `json:"keys" bson:"key"`
	Name                    string                 `json:"name" bson:"name"`
	Unique                  bool                   `json:"unique" bson:"unique"`
	Background              bool                   `json:"background" bson:"background"`
	ExpireAfterSeconds      int32                  `json:"expire_after_seconds" bson:"expire_after_seconds,omitempty"`
	PartialFilterExpression map[string]interface{} `json:"partialFilterExpression" bson:"partialFilterExpression"`
}

// ModeUpdate  根据不同的操作符去更新数据
type ModeUpdate struct {
	Op  string
	Doc interface{}
}

// AggregateOpts TODO
type AggregateOpts struct {
	AllowDiskUse *bool
}

// NewAggregateOpts TODO
func NewAggregateOpts() *AggregateOpts {
	return &AggregateOpts{}
}

// SetAllowDiskUse TODO
func (a *AggregateOpts) SetAllowDiskUse(bl bool) *AggregateOpts {
	a.AllowDiskUse = &bl
	return a
}

type Table interface {
	// Find 查询多个并反序列化到 Result
	Find(filter Filter, opts ...*FindOpts) Find
	// AggregateOne 聚合查询
	AggregateOne(ctx context.Context, pipeline interface{}, result interface{}) error
	AggregateAll(ctx context.Context, pipeline interface{}, result interface{}, opts ...*AggregateOpts) error
	// Insert 插入数据, docs 可以为 单个数据 或者 多个数据
	Insert(ctx context.Context, docs interface{}) error
	// Update 更新数据
	Update(ctx context.Context, filter Filter, doc interface{}) error
	// Upsert TODO
	// update or insert data
	Upsert(ctx context.Context, filter Filter, doc interface{}) error
	// UpdateMultiModel  data based on operators.
	UpdateMultiModel(ctx context.Context, filter Filter, updateModel ...ModeUpdate) error

	// Delete 删除数据
	Delete(ctx context.Context, filter Filter) error

	// CreateIndex 创建索引
	CreateIndex(ctx context.Context, index Index) error
	// BatchCreateIndexes 批量创建索引
	BatchCreateIndexes(ctx context.Context, index []Index) error

	// DropIndex 移除索引
	DropIndex(ctx context.Context, indexName string) error
	// Indexes 查询索引
	Indexes(ctx context.Context) ([]Index, error)

	// AddColumn 添加字段
	AddColumn(ctx context.Context, column string, value interface{}) error
	// RenameColumn 重命名字段
	RenameColumn(ctx context.Context, filter Filter, oldName, newColumn string) error
	// DropColumn 移除字段
	DropColumn(ctx context.Context, field string) error
	// DropColumns 根据条件移除字段
	DropColumns(ctx context.Context, filter Filter, fields []string) error

	// DropDocsColumn remove a column by the name for doc use filter
	DropDocsColumn(ctx context.Context, field string, filter Filter) error

	// Distinct Finds the distinct values for a specified field across a single collection or view and returns the results in an
	// field the field for which to return distinct values.
	// filter query that specifies the documents from which to retrieve the distinct values.
	Distinct(ctx context.Context, field string, filter Filter) ([]interface{}, error)

	// DeleteMany delete document, return number of documents that were deleted.
	DeleteMany(ctx context.Context, filter Filter) (uint64, error)
	// UpdateMany update document, return number of documents that were modified.
	UpdateMany(ctx context.Context, filter Filter, doc interface{}) (uint64, error)
}
