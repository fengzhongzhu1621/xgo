package gorm

import (
	"log"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Student struct {
	Name string
	Age  int
	Id   int
}

// TableName 解决gorm表明映射
func (Student) TableName() string {
	return "Student"
}

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

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		SingularTable: true, // 使用单数表名
	}})
	if err != nil {
		log.Fatal(err)
	}

	var stu Student
	stu.Name = "bob"
	stu.Age = 10

	// 参数必须使用 &
	createStudent(db, &stu)
}
