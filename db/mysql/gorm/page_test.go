package gorm

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func TestScopes(t *testing.T) {
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

	// 查询一条数据，支持分页
	id2 := stu.Id
	row := getStudent2(db, id2)
	assert.Equal(t, row.Name, "bob")

	// 查询一条数据，支持分页
	id3 := stu.Id
	row = getStudent3(db, id3)
	assert.Equal(t, row.Name, "bob")
}

// getStudent2 使用 Scopes 支持分页
func getStudent2(db *gorm.DB, id int) *Student {
	var student Student
	// 根据主键查询
	page := 0
	pageSize := 10
	result := db.Scopes(Paginate(page, pageSize)).First(&student, id)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	return &student
}

func getStudent3(db *gorm.DB, id int) *Student {
	var student Student

	pagination := Pagination{
		page:     0,
		pageSize: 20,
	}
	db.Clauses(&pagination).Find(&student, "id = ?", id)

	return &student
}
