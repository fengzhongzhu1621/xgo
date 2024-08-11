package gorm_2

import (
	"fmt"
	"log"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func rawSqlQuery(db *gorm.DB, query string) []Student {
	var students []Student
	result := db.Raw(query).Scan(&students)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	log.Println("rawSqlQuery completed successfully!")
	return students
}

func TestRawSql(t *testing.T) {
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
	createStudent(db, &stu)

	query := "select * from student"
	for i, student := range rawSqlQuery(db, query) {
		fmt.Println(i, student)
	}

	var students []Student
	db.Raw("select * from student where name = ?", "bob").Scan(&students)
	for i, student2 := range students {
		fmt.Println(i, student2)
	}
}
