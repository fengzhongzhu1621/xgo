package xorm

import (
	"sync"

	"github.com/dlmiddlecote/sqlstats"
	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/db/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	DefaultXormDBClient *XormDBClient
)

var defaultXormDBClientOnce sync.Once

// InitXormDBClient 初始化 db 连接，只能初始化一次
func InitXormDBClient(dbConfig *mysql.Database) {
	if DefaultXormDBClient == nil {
		defaultXormDBClientOnce.Do(func() {
			// 连接数据库
			DefaultXormDBClient = NewXormDBClient(dbConfig)
			DefaultXormDBClient.Connect()
			// 创建一个新的 SQL 统计信息收集器。这个收集器可以用于跟踪和分析 SQL 查询的性能和执行情况。
			// 具体的实现和功能可能会根据你所使用的库或框架而有所不同，但一般来说，它的目的是提供一种机制来收集、存储和报告与 SQL 查询相关的统计信息。
			collector := sqlstats.NewStatsCollector(dbConfig.Name, DefaultXormDBClient.DB.DB())
			// 注册到 prometheus
			prometheus.MustRegister(collector)
		})
	}
}

// GetDefaultXormDBClient 获取默认的DB实例
func GetDefaultXormDBClient() *XormDBClient {
	return DefaultXormDBClient
}

func InitDatabase() {
	// 格式化配置
	cfg := config.GetGlobalConfig()
	cfg.DatabaseMap = make(map[string]mysql.Database)
	for _, db := range cfg.Databases {
		cfg.DatabaseMap[db.ID] = db
	}
	if len(cfg.DatabaseMap) == 0 {
		panic("database cannot be empty")
	}

	// 初始化 DB 连接
	defaultDBConfig, ok := cfg.DatabaseMap["default"]
	if !ok {
		panic("database default should be configured")
	}
	InitXormDBClient(&defaultDBConfig)

	log.Info("init Database success")
}

func init() {
	InitDatabase()
}
