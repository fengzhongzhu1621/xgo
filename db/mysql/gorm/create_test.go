package gorm_2

import (
	"log"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func createStudent(db *gorm.DB, activity *Student) {
	result := db.Create(activity)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	log.Println("Student created successfully!")
}

// 测试 insert 操作
func TestInsert(t *testing.T) {
	dsn := "root:@tcp(127.0.0.1:3306)/xgo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: logger.Default.LogMode(logger.Info), // 输出详细日志
	})
	if err != nil {
		log.Fatal(err)
	}

	var stu Student
	stu.Name = "bob"
	stu.Age = 10

	createStudent(db, &stu)
}

func batchInsertStudents(db *gorm.DB, students []Student) {
	result := db.Create(&students)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
}

func TestBulkInsert(t *testing.T) {
	dsn := "root:@tcp(127.0.0.1:3306)/xgo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		SingularTable: true, // 使用单数表名
	}})
	if err != nil {
		log.Fatal(err)
	}

	var stu1 Student
	stu1.Name = "foo"
	stu1.Age = 10

	var stu2 Student
	stu2.Name = "bar"
	stu2.Age = 11

	batchInsertStudents(db, []Student{stu1, stu2})
}
