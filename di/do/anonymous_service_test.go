package do

import (
	"database/sql"
	"testing"

	"github.com/samber/do"
)

type DBService struct {
	db *sql.DB
}

func (c *DBService) Ping() {
	println("pong")
}

func TestAnonymousService(t *testing.T) {
	injector := do.New()

	do.Provide[*DBService](injector, func(i *do.Injector) (*DBService, error) {
		// 补齐 sql.Open 的参数
		db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/dbname")
		if err != nil {
			return nil, err
		}

		// 可选：测试连接是否成功（生产环境建议添加）
		if err := db.Ping(); err != nil {
			return nil, err
		}

		return &DBService{db: db}, nil
	})

	injector.Shutdown()
}
