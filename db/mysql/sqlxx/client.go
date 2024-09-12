package sqlxx

import (
	"fmt"
	"time"

	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

const (
	defaultMaxOpenConns    = 100
	defaultMaxIdleConns    = 25
	defaultConnMaxLifetime = 10 * time.Minute
)

type DBClient struct {
	name string // db名称

	DB *sqlx.DB // db 连接对象

	dataSource string // db 连接字符串

	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
}

// TestConnection 测试 db 连接是否正常
func (db *DBClient) TestConnection() (err error) {
	// 创建 db 连接
	conn, err := sqlx.Connect("mysql", db.dataSource)
	if err != nil {
		return err
	}

	conn.Close()
	return nil
}

// Connect 连接数据库
func (db *DBClient) Connect() error {
	// 创建 db 连接
	var err error
	db.DB, err = sqlx.Connect("mysql", db.dataSource)
	if err != nil {
		return err
	}

	// 设置连接数
	db.DB.SetMaxOpenConns(db.maxOpenConns)
	db.DB.SetMaxIdleConns(db.maxIdleConns)
	db.DB.SetConnMaxLifetime(db.connMaxLifetime)

	log.Infof("connect to database: %s[maxOpenConns=%d, maxIdleConns=%d, connMaxLifetime=%s]",
		db.name, db.maxOpenConns, db.maxIdleConns, db.connMaxLifetime)

	return nil
}

// Close close db connection
func (db *DBClient) Close() {
	if db.DB != nil {
		db.DB.Close()
	}
}

// NewDBClient 根据db 配置创建数据库连接对象
func NewDBClient(cfg *config.Database) *DBClient {
	dataSource := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s&parseTime=True&interpolateParams=true&loc=%s",
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
		log.Errorf("error config for database %s, maxOpenConns should greater or equals to maxIdleConns, will"+
			"use the default [defaultMaxOpenConns=%d, defaultMaxIdleConns=%d]",
			cfg.Name, defaultMaxOpenConns, defaultMaxIdleConns)
		maxOpenConns = defaultMaxOpenConns
		maxIdleConns = defaultMaxIdleConns
	}

	connMaxLifetime := defaultConnMaxLifetime
	if cfg.ConnMaxLifetimeSecond > 0 {
		if cfg.ConnMaxLifetimeSecond >= 60 {
			connMaxLifetime = time.Duration(cfg.ConnMaxLifetimeSecond) * time.Second
		} else {
			log.Errorf("error config for database %s, connMaxLifetimeSeconds should be greater than 60 seconds"+
				"use the default [defaultConnMaxLifetime=%s]",
				cfg.Name, defaultConnMaxLifetime)
		}
	}

	return &DBClient{
		name:            cfg.Name,
		dataSource:      dataSource,
		maxOpenConns:    maxOpenConns,
		maxIdleConns:    maxIdleConns,
		connMaxLifetime: connMaxLifetime,
	}
}

// TestConnection 根据 db 配置测试数据库连接是否正常
func TestConnection(dbConfig *config.Database) error {
	c := NewDBClient(dbConfig)
	return c.TestConnection()
}
