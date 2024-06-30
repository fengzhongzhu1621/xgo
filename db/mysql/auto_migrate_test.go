package mysql

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql" // _ 表示不直接使用，在 main()调用前执行，自动执行 init() 函数
	"github.com/jinzhu/gorm"
)

type Student struct {
	Name string
	Age  int
	Id   int
}

type StudentSize struct {
	Name string `gorm:"size:100;default:''"`
	Age  int `gorm:not null"`
	Id   int
}

type StudentTimestamp struct {
	Name string `gorm:"size:100;default:''"`
	Age  int `gorm:not null"`
	Id   int
	CreateAt time.Time
	Join time.Time `gorm:"type:timestamp"`
}

// 根据结构体创建数据表
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
	// CREATE TABLE `student` (
	//   `name` varchar(255) DEFAULT NULL,
	//   `age` int DEFAULT NULL,
	//   `id` int NOT NULL AUTO_INCREMENT,
	//   PRIMARY KEY (`id`)
	// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

	fmt.Println(conn.AutoMigrate(new(Student)).Error)


	// 设置字段属性，生效场景
	// 1. 第一次创建表时
	// 2. 表增加新字段时
	// 其他场景下，修改表属性，在 gorm 操作中，是无效的
	/*
	CREATE TABLE `student_size` (
		`name` varchar(100) DEFAULT '',
		`age` int DEFAULT NULL,
		`id` int NOT NULL AUTO_INCREMENT,
		PRIMARY KEY (`id`)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
	*/
	fmt.Println(conn.AutoMigrate(new(StudentSize)).Error)
	
	// 设置时间字段格式
	/*
	CREATE TABLE `student_timestamp` (
		`name` varchar(100) DEFAULT '',
		`age` int DEFAULT NULL,
		`id` int NOT NULL AUTO_INCREMENT,
		`create_at` datetime DEFAULT NULL,
		`join` timestamp NULL DEFAULT NULL,
		PRIMARY KEY (`id`)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
	*/
	fmt.Println(conn.AutoMigrate(new(StudentTimestamp)).Error)
}
