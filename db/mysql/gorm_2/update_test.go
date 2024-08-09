package gorm_2

import (
	"log"
	"testing"

	"github.com/go-playground/assert/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func updateStudent(db *gorm.DB, id int, newName string) {
	var student Student
	result := db.First(&student, id)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	student.Name = newName
	db.Save(&student)
}

func batchUpdateStudents(db *gorm.DB, students []Student) {
	for _, student := range students {
		db.Model(&student).Updates(student)
	}
}

func TestUpdate(t *testing.T) {
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
	updateStudent(db, stu.Id, "bar")

	db.First(&stu, stu.Id)
	assert.Equal(t, stu.Name, "bar")

}
