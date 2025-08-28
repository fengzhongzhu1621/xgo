package gorm

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

	// INSERT INTO `product` (`created_at`,`updated_at`,`deleted_at`,`code`,`price`) VALUES ('2025-07-31 15:39:37.677','2025-07-31 15:39:37.677',NULL,'A1',100)
	db.Create(&Product{Code: "A1", Price: 100})

	var product Product
	// SELECT * FROM `product` WHERE `product`.`id` = 1 AND `product`.`deleted_at` IS NULL ORDER BY `product`.`id` LIMIT 1
	db.First(&product, 1)
	// SELECT * FROM `product` WHERE code = 'A1' AND `product`.`deleted_at` IS NULL ORDER BY `product`.`id` LIMIT 1
	db.First(&product, "code = ?", "A1")

	// UPDATE `product` SET `price`=200,`updated_at`='2025-07-31 15:39:37.679' WHERE `product`.`deleted_at` IS NULL
	db.Model(&product).Update("Price", 200)
	// UPDATE `product` SET `price`=0,`updated_at`='2025-07-31 15:42:15.808' WHERE `product`.`deleted_at` IS NULL
	db.Model(&product).Update("Price", 0) // 可以更新非零值

	// UPDATE `product` SET `updated_at`='2025-07-31 15:42:42.601',`code`='A2' WHERE `product`.`deleted_at` IS NULL
	db.Model(&product).Updates(Product{Price: 0, Code: "A2"}) // 仅更新非零值字段
	// UPDATE `product` SET `updated_at`='2025-07-31 15:57:23.516',`code`='A2',`price`=0 WHERE `product`.`deleted_at` IS NULL
	db.Model(&product).Select("Price", "Code").Updates(Product{Price: 0, Code: "A2"}) // 可以更新非零值
	// UPDATE `product` SET `updated_at`='2025-07-31 15:58:18.767',`price`=0 WHERE `product`.`deleted_at` IS NULL
	db.Model(&product).Select("Price").Updates(Product{Price: 0, Code: "A2"}) // 可以更新非零值
	// UPDATE `product` SET `code`='A3',`price`=0,`updated_at`='2025-07-31 15:43:14.475' WHERE `product`.`deleted_at` IS NULL
	db.Model(&product).Updates(map[string]interface{}{"Price": 0, "Code": "A3"}) // 可以更新非零值

	// UPDATE `product` SET `deleted_at`='2025-07-31 15:43:14.478' WHERE `product`.`id` = 1 AND `product`.`deleted_at` IS NULL
	db.Delete(&product, 1)
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
