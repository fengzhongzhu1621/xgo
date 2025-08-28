package do

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/samber/do"
)

func TestNamedService(t *testing.T) {
	injector := do.New()

	do.ProvideNamed(injector, "dbconn", func(i *do.Injector) (*DBService, error) {
		// 补齐 sql.Open 的参数
		db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/xgo")
		if err != nil {
			return nil, err
		}

		// 可选：测试连接是否成功（生产环境建议添加）
		if err := db.Ping(); err != nil {
			return nil, err
		}

		return &DBService{db: db}, nil
	})

	s := do.MustInvokeNamed[*DBService](injector, "dbconn")
	s.Ping() // Pong

	injector.Shutdown()
}
