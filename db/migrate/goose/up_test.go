package goose

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
)

func TestUp(t *testing.T) {
	// 连接数据库
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/xgo?parseTime=True&loc=Local&multiStatements=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 执行迁移
	if err := goose.Up(db, "./migrations"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database migrated successfully")
}
