package xorm

import (
	"fmt"
	"sync"
	"time"

	xormLog "xorm.io/xorm/log"

	"github.com/fengzhongzhu1621/xgo/db/mysql"
	log "github.com/sirupsen/logrus"
	"xorm.io/core"
	"xorm.io/xorm"
)

const (
	defaultMaxOpenConns    = 100
	defaultMaxIdleConns    = 25
	defaultConnMaxLifetime = 10 * time.Minute
)

var (
	once         sync.Once
	PingInterval = 1200 * time.Second
)

type XormDBClient struct {
	// db名称
	name string

	// 一个数据库引擎实例
	DB *xorm.Engine

	// db 连接字符串
	dataSource string

	// 最大连接数
	maxOpenConns int
	// 最大空闲连接数
	maxIdleConns int
	// 单个连接的最大生命周期
	connMaxLifetime time.Duration

	// 用于调试时打印sql
	debugMode bool
}

// SyncTable 将 Go 结构体映射到数据库中的表，并确保表结构是最新的
func (db *XormDBClient) SyncTable(args ...interface{}) error {
	return db.DB.Sync2(args...)
}

// ExecuteRawSQL 用于执行原生 SQL 语句。这个方法不会将结果集映射到任何结构体，而是直接执行 SQL 并返回受影响的行数以及可能出现的错误
// 用于执行任意的 SQL 语句，并返回一个 sql.Result 对象，该对象包含了执行 SQL 语句后的结果信息，如受影响的行数、最后插入的 ID 等。
func (db *XormDBClient) ExecuteRawSQL(args ...interface{}) error {
	_, err := db.DB.Exec(args...)
	return err
}

// QueryString 用于调试和查看生成的 SQL 字符串。这个方法主要用于调试和开发过程中查看查询结果的原始 SQL 字符串，而不是用于生产环境中的数据检索。
func (db *XormDBClient) QueryString(args ...interface{}) ([]map[string]string, error) {
	return db.DB.QueryString(args...)
}

// keepAlive 保持 db 的链接，防止被空闲回收
func (db *XormDBClient) KeepAlive() {
	for range time.Tick(PingInterval) {
		if db.DB != nil {
			// 检查与数据库的连接是否仍然有效
			// 在生产环境中，频繁地调用 Ping 可能会影响性能。因此，建议根据实际需求合理地使用它。
			err := db.DB.Ping()
			if err != nil {
				log.Errorf("Error pinging database: %v", err)
				return
			}
		}
	}
}

// TestConnection 测试数据库连接是否正常
func (db *XormDBClient) TestConnection() (err error) {
	err = db.DB.Ping()
	if err != nil {
		return err
	}

	return nil
}

// Close 关闭 db 连接
func (db *XormDBClient) Close() {
	if db.DB != nil {
		db.DB.Close()
	}
}

// Connect to db, and update some settings
func (db *XormDBClient) Connect() error {
	var (
		err error
	)

	// 创建一个新的 Engine 实例
	if db.DB, err = xorm.NewEngine("mysql", db.dataSource); err != nil {
		panic(fmt.Errorf("xorm.NewEngine error: %v", err))
	}
	// core.GonicMapper 是一个自定义的映射器，它会将结构体字段名保持原样（即驼峰命名法，camelCase），而不是转换为蛇形命名法。
	// core.SnakeCaseMapper 它会将结构体字段名转换为蛇形命名法（snake_case）的数据库列名。
	db.DB.SetMapper(core.GonicMapper{})

	// 设置连接数
	db.DB.SetMaxOpenConns(db.maxOpenConns)
	// 设置数据库连接池中最大空闲连接数。这个方法可以帮助你控制数据库连接的资源使用，优化应用程序的性能。
	// 如果没有显式设置最大空闲连接数，sqlx 会使用 database/sql 包的默认值，通常是 2
	// 假设你的应用程序在高并发时段需要处理大量的数据库请求，但在低峰时段请求量较少。
	// 在这种情况下，你可以设置一个较高的最大空闲连接数，以确保在高并发时段有足够的连接可用；而在低峰时段，多余的连接会自动关闭，释放资源。
	db.DB.SetMaxIdleConns(db.maxIdleConns)
	// 设置数据库连接池中单个连接的最大生命周期。可以控制数据库连接的复用时间，避免因长时间使用同一个连接而导致潜在的问题。
	db.DB.SetConnMaxLifetime(db.connMaxLifetime)

	err = db.DB.Ping()
	if err != nil {
		panic(fmt.Errorf("ping error: %v", err))
	}

	// 每次执行数据库操作时，将生成的 SQL 语句及其参数打印到控制台。
	if db.debugMode {
		dbLogger := logging.GetXormDBLogger()
		if dbLogger != nil {
			db.DB.SetLogger(xormLog.NewLoggerAdapter(dbLogger))
		}
		db.DB.ShowSQL(true)
	}

	go func() {
		db.KeepAlive()
	}()

	return nil
}

// NewXormDBClient 创建 mysql 客户端
func NewXormDBClient(config *mysql.Database) *XormDBClient {
	// parseTime 是一个布尔值, 是一个配置选项，可以在创建引擎时设置，以控制 XORM 如何处理数据库中的时间字段。
	// 当设置为 true 时，XORM 会在读取数据库中的时间字段时自动将其解析为 Go 语言中的 time.Time 类型。
	// 如果设置为 false，XORM 将时间字段作为字符串返回。
	// 如果数据库中的时间字段使用了特定的时区或格式，你可能需要自定义时间解析的行为。
	// 如果不需要自动解析时间字段，可以将 parseTime 设置为 false 或者从连接字符串中移除该参数。
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=False&interpolateParams=true&loc=%s",
		config.User, config.Password, config.Host, config.Port, config.Name, "utf8mb4", "UTC")

	// 获得最大连接数
	maxOpenConns := defaultMaxOpenConns
	if config.MaxOpenConns > 0 {
		maxOpenConns = config.MaxOpenConns
	}

	// 获得最大空闲连接数
	maxIdleConns := defaultMaxIdleConns
	if config.MaxIdleConns > 0 {
		maxIdleConns = config.MaxIdleConns
	}

	if maxOpenConns < maxIdleConns {
		maxOpenConns = defaultMaxOpenConns
		maxIdleConns = defaultMaxIdleConns
	}

	// 获得单个连接的最大生命周期
	connMaxLifetime := defaultConnMaxLifetime
	if config.ConnMaxLifetimeSecond > 0 {
		if config.ConnMaxLifetimeSecond >= 60 {
			connMaxLifetime = time.Duration(config.ConnMaxLifetimeSecond) * time.Second
		}
	}

	return &XormDBClient{
		name:            config.Name,
		dataSource:      dataSource,
		maxOpenConns:    maxOpenConns,
		maxIdleConns:    maxIdleConns,
		connMaxLifetime: connMaxLifetime,
		debugMode:       config.DebugMode,
	}
}

// TestConnection 根据 db 配置测试数据库连接是否正常
func TestConnection(dbConfig *mysql.Database) error {
	c := NewXormDBClient(dbConfig)
	return c.TestConnection()
}
