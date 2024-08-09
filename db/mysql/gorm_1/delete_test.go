package gorm_1

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
)

func TestDelete(t *testing.T) {
	conn, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/xgo?parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("gorm.Open err: ", err)
		return
	}

	defer conn.Close()

	/*
		type Model struct {
			ID        uint `gorm:"primary_key"`
			CreatedAt time.Time
			UpdatedAt time.Time
			DeletedAt *time.Time `sql:"index"`
		}

		CREATE TABLE `student_gorms` (
			`id` int unsigned NOT NULL AUTO_INCREMENT,
			`created_at` datetime DEFAULT NULL,
			`updated_at` datetime DEFAULT NULL,
			`deleted_at` datetime DEFAULT NULL,
			`name` varchar(255) DEFAULT NULL,
			`age` int DEFAULT NULL,
			PRIMARY KEY (`id`),
			KEY `idx_student_gorms_deleted_at` (`deleted_at`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
	*/
	type StudentGorm struct {
		gorm.Model
		Name string
		Age  int
	}
	fmt.Println(conn.AutoMigrate(new(StudentGorm)).Error)

	// 插入一条数据
	/*
		mysql> select * from student_gorms;
		+----+---------------------+---------------------+------------+------+------+
		| id | created_at          | updated_at          | deleted_at | name | age  |
		+----+---------------------+---------------------+------------+------+------+
		|  1 | 2024-06-30 18:36:29 | 2024-06-30 18:36:29 | NULL       | bob  |   10 |
		+----+---------------------+---------------------+------------+------+------+
	*/
	var stu StudentGorm
	stu.Name = "bob"
	stu.Age = 10
	fmt.Println(conn.Create(&stu).Error)

	// 软删除
	/*
		mysql> select * from student_gorms;
		+----+---------------------+---------------------+---------------------+------+------+
		| id | created_at          | updated_at          | deleted_at          | name | age  |
		+----+---------------------+---------------------+---------------------+------+------+
		|  1 | 2024-06-30 18:36:29 | 2024-06-30 18:36:29 | 2024-06-30 18:38:27 | bob  |   10 |
		+----+---------------------+---------------------+---------------------+------+------+
	*/
	fmt.Println(conn.Where("name = ?", "bob").Delete(new(StudentGorm)).Error)

	// 查询被删除的数据，没有查询到
	var stu2 []StudentGorm
	conn.Find(&stu2)
	fmt.Println(stu2) // []

	// 查询被软删除的数据
	conn.Unscoped().Find(&stu2)
	fmt.Println(&stu2)

	// 物理删除
	conn.Unscoped().Where("name = ?", "bob").Delete(new(StudentGorm))
}
