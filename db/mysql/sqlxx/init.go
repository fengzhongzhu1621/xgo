package sqlxx

import (
	"sync"

	"github.com/fengzhongzhu1621/xgo/db/mysql"
)

var (
	DefaultSqlxDBClient *SqlxDBClient
)

var defaultSqlxDBClientOnce sync.Once

// InitSqlxDBClient 初始化 db 连接，只能初始化一次
func InitSqlxDBClient(dbConfig *mysql.Database) {
	if DefaultSqlxDBClient == nil {
		defaultSqlxDBClientOnce.Do(func() {
			// 连接数据库
			DefaultSqlxDBClient = NewSqlxDBClient(dbConfig)
			DefaultSqlxDBClient.Connect()
		})
	}
}

// GetDefaultSqlxDBClient 获取默认的DB实例
func GetDefaultSqlxDBClient() *SqlxDBClient {
	return DefaultSqlxDBClient
}
