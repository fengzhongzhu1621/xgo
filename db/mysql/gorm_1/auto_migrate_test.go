package gorm_1

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql" // _ 表示不直接使用，在 main()调用前执行，自动执行 init() 函数
	"github.com/jinzhu/gorm"
)

type Student struct {
	Name string `gorm:"size:255;uniqueIndex"`
	Age  int
	Id   int `gorm:"primaryKey"`
}

type StudentSize struct {
	Name string `gorm:"size:100;default:''"`
	Age  int    `gorm:not null"`
	Id   int
}

type StudentTimestamp struct {
	Name     string `gorm:"size:100;default:''"`
	Age      int    `gorm:not null"`
	Id       int
	CreateAt time.Time
	Join     time.Time `gorm:"type:timestamp"`
}

type User struct {
	ID       int
	Name     string    `gorm:"index"`
	Students []Student `gorm:"many2many:user_students;"`
}

type User2 struct {
	ID        uint   `gorm:"column:user_id"`
	Name      string `gorm:"column:user_name"`
	Email     string `gorm:"column:user_email"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User2) TableName() string {
	return "users"
}

// one2one
type CreditCard struct {
	ID     uint
	Number string
	UserID uint
	User   User
}

type User3 struct {
	ID          uint
	Name        string
	Email       string
	CreditCard3 CreditCard
}

// one2many
type Order struct {
	ID     uint
	Amount float64
	UserID uint
}

type User4 struct {
	ID     uint
	Name   string
	Email  string
	Orders []Order
}

// many2many
type Product struct {
	ID    uint
	Name  string
	Users []User `gorm:"many2many:user_products;"`
}

type User5 struct {
	ID       uint
	Name     string
	Email    string
	Products []Product `gorm:"many2many:user_products;"`
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
	/*
		CREATE TABLE `student` (
			`name` varchar(255) DEFAULT NULL,
			`age` int DEFAULT NULL,
			`id` int NOT NULL AUTO_INCREMENT,
			PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
	*/
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

	/*
		CREATE TABLE `user` (
			`id` int NOT NULL AUTO_INCREMENT,
			`name` varchar(255) DEFAULT NULL,
			PRIMARY KEY (`id`),
			KEY `idx_user_name` (`name`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

		CREATE TABLE `user_students` (
			`user_id` int NOT NULL,
			`student_id` int NOT NULL,
			PRIMARY KEY (`user_id`,`student_id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
	*/
	fmt.Println(conn.AutoMigrate(new(User)).Error)

	/*
		CREATE TABLE `users` (
			`user_id` int unsigned DEFAULT NULL,
			`user_name` varchar(255) DEFAULT NULL,
			`user_email` varchar(255) DEFAULT NULL,
			`created_at` datetime DEFAULT NULL,
			`updated_at` datetime DEFAULT NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
	*/
	fmt.Println(conn.AutoMigrate(new(User2)).Error)

	/*
		CREATE TABLE `credit_card` (
			`id` int unsigned NOT NULL AUTO_INCREMENT,
			`number` varchar(255) DEFAULT NULL,
			`user_id` int unsigned DEFAULT NULL,
			PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci


		CREATE TABLE `user3` (
			`id` int unsigned NOT NULL AUTO_INCREMENT,
			`name` varchar(255) DEFAULT NULL,
			`email` varchar(255) DEFAULT NULL,
			PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
	*/
	fmt.Println(conn.AutoMigrate(new(CreditCard)).Error)
	fmt.Println(conn.AutoMigrate(new(User3)).Error)

	/*
		CREATE TABLE `order` (
			`id` int unsigned NOT NULL AUTO_INCREMENT,
			`amount` double DEFAULT NULL,
			`user_id` int unsigned DEFAULT NULL,
			PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

		CREATE TABLE `user4` (
			`id` int unsigned NOT NULL AUTO_INCREMENT,
			`name` varchar(255) DEFAULT NULL,
			`email` varchar(255) DEFAULT NULL,
			PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
	*/
	fmt.Println(conn.AutoMigrate(new(Order)).Error)
	fmt.Println(conn.AutoMigrate(new(User4)).Error)

	/*
		CREATE TABLE `product` (
			`id` int unsigned NOT NULL AUTO_INCREMENT,
			`name` varchar(255) DEFAULT NULL,
			PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

		CREATE TABLE `user5` (
			`id` int unsigned NOT NULL AUTO_INCREMENT,
			`name` varchar(255) DEFAULT NULL,
			`email` varchar(255) DEFAULT NULL,
			PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

		 CREATE TABLE `user_products` (
			`product_id` int unsigned NOT NULL,
			`user_id` int NOT NULL,
			PRIMARY KEY (`product_id`,`user_id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
	*/
	fmt.Println(conn.AutoMigrate(new(Product)).Error)
	fmt.Println(conn.AutoMigrate(new(User5)).Error)

}
