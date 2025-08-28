package sqlxx

import (
	"fmt"
	"time"

	"github.com/fengzhongzhu1621/xgo/db/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	defaultMaxOpenConns    = 100
	defaultMaxIdleConns    = 25
	defaultConnMaxLifetime = 10 * time.Minute
)

type SqlxDBClient struct {
	name string // db名称

	DB *sqlx.DB // db 连接对象

	dataSource string // db 连接字符串

	// 最大连接数
	maxOpenConns int
	// 最大空闲连接数
	maxIdleConns int
	// 单个连接的最大生命周期
	connMaxLifetime time.Duration
}

// TestConnection 测试 db 连接是否正常
func (db *SqlxDBClient) TestConnection() (err error) {
	// 创建 db 连接
	conn, err := sqlx.Connect("mysql", db.dataSource)
	if err != nil {
		return err
	}

	conn.Close()
	return nil
}

// Connect 连接数据库
func (db *SqlxDBClient) Connect() error {
	// 创建 db 连接
	var err error
	db.DB, err = sqlx.Connect("mysql", db.dataSource)
	if err != nil {
		panic(fmt.Errorf("sqlx.Connect error: %v", err))
	}

	// 设置连接数
	db.DB.SetMaxOpenConns(db.maxOpenConns)
	// 设置数据库连接池中最大空闲连接数。这个方法可以帮助你控制数据库连接的资源使用，优化应用程序的性能。
	// 如果没有显式设置最大空闲连接数，sqlx 会使用 database/sql 包的默认值，通常是 2
	// 假设你的应用程序在高并发时段需要处理大量的数据库请求，但在低峰时段请求量较少。
	// 在这种情况下，你可以设置一个较高的最大空闲连接数，以确保在高并发时段有足够的连接可用；而在低峰时段，多余的连接会自动关闭，释放资源。
	db.DB.SetMaxIdleConns(db.maxIdleConns)
	// 设置数据库连接池中单个连接的最大生命周期。可以控制数据库连接的复用时间，避免因长时间使用同一个连接而导致潜在的问题。
	db.DB.SetConnMaxLifetime(db.connMaxLifetime)

	return nil
}

// Close close db connection
func (db *SqlxDBClient) Close() {
	if db.DB != nil {
		db.DB.Close()
	}
}

// NewSqlxDBClient 根据db 配置创建数据库连接对象
func NewSqlxDBClient(cfg *mysql.Database) *SqlxDBClient {
	dataSource := fmt.Sprintf(
		"%s:%s@(%s:%d)/%s?charset=%s&parseTime=True&interpolateParams=true&loc=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		"utf8",
		"UTC",
	)

	maxOpenConns := defaultMaxOpenConns
	if cfg.MaxOpenConns > 0 {
		maxOpenConns = cfg.MaxOpenConns
	}

	maxIdleConns := defaultMaxIdleConns
	if cfg.MaxIdleConns > 0 {
		maxIdleConns = cfg.MaxIdleConns
	}

	if maxOpenConns < maxIdleConns {
		maxOpenConns = defaultMaxOpenConns
		maxIdleConns = defaultMaxIdleConns
	}

	connMaxLifetime := defaultConnMaxLifetime
	if cfg.ConnMaxLifetimeSecond > 0 {
		if cfg.ConnMaxLifetimeSecond >= 60 {
			connMaxLifetime = time.Duration(cfg.ConnMaxLifetimeSecond) * time.Second
		}
	}

	return &SqlxDBClient{
		name:            cfg.Name,
		dataSource:      dataSource,
		maxOpenConns:    maxOpenConns,
		maxIdleConns:    maxIdleConns,
		connMaxLifetime: connMaxLifetime,
	}
}

// TestConnection 根据 db 配置测试数据库连接是否正常
func TestConnection(dbConfig *mysql.Database) error {
	c := NewSqlxDBClient(dbConfig)
	return c.TestConnection()
}
