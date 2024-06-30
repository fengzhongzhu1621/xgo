package mysql

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
)

// 测试 insert 操作
func TestInsert(t *testing.T) {
	// create database xgo charset utf8mb4;
	// 创建数据库连接 用户名:密码@协议(IP:port)/数据库名?a=xxx&b=xxx
	// 创建一个默认的连接池
	conn, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/xgo?parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("gorm.Open err: ", err)
		return
	}

	defer conn.Close()

	// 创建非复数表名，默认创建复数表名
	var stu Student
	stu.Name = "bob"
	stu.Age = 10

	// 创建非复数表名，默认创建复数表名
	// 必须使用，需要和表创建时的参数保持一致
	conn.SingularTable(true)
	// 参数必须使用 &
	fmt.Println(conn.Create(&stu).Error)
}
