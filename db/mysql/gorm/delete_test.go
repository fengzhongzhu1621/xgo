package gorm

import (
	"log"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func deleteStudent(db *gorm.DB, id int) {
	// 根据主键删除
	result := db.Delete(&Student{}, id)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	// 根据记录对象删除
	var student Student
	db.First(&student, id)
	result = db.Delete(&student)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
}

func TestDelete(t *testing.T) {
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

	// 删除记录
	deleteStudent(db, stu.Id)
}
