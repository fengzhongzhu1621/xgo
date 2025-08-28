package gorm

import (
	"log"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// CREATE TABLE `Student` (
//   `name` longtext,
//   `age` bigint DEFAULT NULL,
//   `id` bigint NOT NULL AUTO_INCREMENT,
//   PRIMARY KEY (`id`)
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4

type Student struct {
	Name string
	Age  int
	Id   int
}

// TableName 解决gorm表明映射
func (Student) TableName() string {
	return "Student"
}

type User struct {
	ID       int
	Name     string
	Students []Student `gorm:"many2many:user_students;"`
}

type Product2 struct {
	ID         uint `gorm:"primaryKey"`
	Name       string
	CategoryID int     `gorm:"index:idx_category"`
	Price      float64 `gorm:"index:idx_price"`
	CreatedAt  time.Time
}

func TestCreateConnect(t *testing.T) {
	dsn := "root:@tcp(127.0.0.1:3306)/xgo?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 输出详细日志
	})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB.SetMaxIdleConns(10)           // 设置最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 设置最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置连接最大存活时间

	log.Println("Database connected successfully!")
}
