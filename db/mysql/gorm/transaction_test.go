package gorm_2

import (
	"log"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func transactionalUpdate(db *gorm.DB, id int, newName string) {
	err := db.Transaction(func(tx *gorm.DB) error {
		var student Student
		if err := tx.First(&student, id).Error; err != nil {
			return err
		}
		student.Name = newName
		if err := tx.Save(&student).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Transactional update completed successfully!")
}

func TestTransaction(t *testing.T) {
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
	transactionalUpdate(db, stu.Id, "foo")

}
