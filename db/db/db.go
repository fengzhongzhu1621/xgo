package db

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// DB db operation interface
type DB interface {
	// Table collection 操作
	Table(collection string) Table

	// NextSequence 获取新序列号(非事务)
	NextSequence(ctx context.Context, sequenceName string) (uint64, error)

	// NextSequences 批量获取新序列号(非事务)
	NextSequences(ctx context.Context, sequenceName string, num int) ([]uint64, error)

	// Ping 健康检查
	Ping() error // 健康检查

	// HasTable 判断是否存在集合
	HasTable(ctx context.Context, name string) (bool, error)
	// ListTables 获取所有的表名
	ListTables(ctx context.Context) ([]string, error)
	// DropTable 移除集合
	DropTable(ctx context.Context, name string) error
	// CreateTable 创建集合
	CreateTable(ctx context.Context, name string) error
	// RenameTable 更新集合名称
	RenameTable(ctx context.Context, prevName, currName string) error

	IsDuplicatedError(error) bool
	IsNotFoundError(error) bool

	Close() error

	// CommitTransaction 提交事务
	CommitTransaction(context.Context, *TxnCapable) error
	// AbortTransaction 取消事务
	AbortTransaction(context.Context, *TxnCapable) (bool, error)

	// InitTxnManager TxnID management of initial transaction
	InitTxnManager(r *redis.Client) error
}
