package gorm_2

import (
	"log"
	"testing"

	"github.com/go-playground/assert/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func getStudent(db *gorm.DB, id int) *Student {
	var student Student
	// 根据主键查询
	result := db.First(&student, id)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	return &student
}

func paginateStudents(db *gorm.DB, page, pageSize int) []Student {
	var activities []Student
	db.Where("name = ?", "bob").Or("name = ?", "foo").Order("created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&activities)
	return activities
}

func TestQueryConnect(t *testing.T) {
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

	// 查询一条数据
	var id = stu.Id
	row := getStudent(db, id)
	assert.Equal(t, row.Name, "bob")

	// 分页查询
	paginateStudents(db, 1, 10)
}
