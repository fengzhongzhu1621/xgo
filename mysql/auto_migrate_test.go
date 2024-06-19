package mysql

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql" // _ 表示不直接使用，在 main()调用前执行，自动执行 init() 函数
	"github.com/jinzhu/gorm"
)

type Student struct {
	Name string
	Age  int
	Id   int
}

func TestAutoMigrate(t *testing.T) {
	// create database xgo charset=utf8;
	// 创建数据库连接 用户名:密码@协议(IP:port)/数据库名?xxx
	conn, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/xgo")
	if err != nil {
		fmt.Println("gorm.Open err: ", err)
		return
	}

	defer conn.Close()

	// 创建非复数表名，默认创建复数表名
	conn.SingularTable(true)

	// 创建数据表
	fmt.Println(conn.AutoMigrate(new(Student)).Error)

	// CREATE TABLE `student` (
	//   `name` varchar(255) DEFAULT NULL,
	//   `age` int DEFAULT NULL,
	//   `id` int NOT NULL AUTO_INCREMENT,
	//   PRIMARY KEY (`id`)
	// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

}
