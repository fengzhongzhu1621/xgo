package gorm

import (
	"fmt"
	"log"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 实现方式1：has many （一对多）
// CREATE TABLE `user_has_manies` (
//
//	`id` bigint unsigned NOT NULL AUTO_INCREMENT,
//	`name` longtext,
//	PRIMARY KEY (`id`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
type UserHasMany struct {
	ID     uint
	Name   string
	Orders []Order `gorm:"foreignKey:UserID;references:ID"` // 显式定义外键和引用关系
}

// CREATE TABLE `orders` (
//
//	`id` bigint unsigned NOT NULL AUTO_INCREMENT,
//	`user_id` bigint unsigned DEFAULT NULL,
//	`price` double DEFAULT NULL,
//	PRIMARY KEY (`id`),
//	KEY `fk_user_has_manies_orders` (`user_id`),
//	CONSTRAINT `fk_user_has_manies_orders` FOREIGN KEY (`user_id`) REFERENCES `user_has_manies` (`id`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
type Order struct {
	ID     uint
	UserID uint // 外键
	Price  float64
}

// 实现方式2：belong to （多对一）
// CREATE TABLE `companies` (
//
//	`id` bigint NOT NULL AUTO_INCREMENT,
//	`name` longtext,
//	PRIMARY KEY (`id`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
type Company struct {
	ID   int
	Name string
}

// CREATE TABLE `user_with_companies` (
//
//	`id` bigint unsigned NOT NULL AUTO_INCREMENT,
//	`created_at` datetime(3) DEFAULT NULL,
//	`updated_at` datetime(3) DEFAULT NULL,
//	`deleted_at` datetime(3) DEFAULT NULL,
//	`name` longtext,
//	`company_id` bigint DEFAULT NULL,
//	PRIMARY KEY (`id`),
//	KEY `idx_user_with_companies_deleted_at` (`deleted_at`),
//	KEY `fk_user_with_companies_company` (`company_id`),
//	CONSTRAINT `fk_user_with_companies_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
type UserWithCompany struct {
	gorm.Model
	Name      string
	CompanyID int
	Company   Company `gorm:"foreignKey:CompanyID;references:ID"` // 用户属于某个公司
}

func TestManyToOne(t *testing.T) {
	dsn := "root:@tcp(127.0.0.1:3306)/xgo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 输出详细日志
	})
	if err != nil {
		log.Fatal(err)
	}

	// INSERT INTO `orders` (`user_id`,`price`) VALUES (1,9.9),(1,19.9) ON DUPLICATE KEY UPDATE `user_id`=VALUES(`user_id`)
	// INSERT INTO `user_has_manies` (`name`) VALUES ('Bob')
	user := UserHasMany{
		Name: "Bob",
		Orders: []Order{
			{Price: 9.9},
			{Price: 19.9},
		},
	}
	db.Create(&user)

	// SELECT * FROM `orders` WHERE `orders`.`user_id` = 1
	// SELECT * FROM `user_has_manies` WHERE `user_has_manies`.`id` = 1 ORDER BY `user_has_manies`.`id` LIMIT 1
	var user2 UserHasMany
	db.Preload("Orders").First(&user2, 1)
	fmt.Printf(
		"User: %+v\n",
		user2,
	) // {ID:1 Name:Bob Orders:[{ID:1 UserID:1 Price:9.9} {ID:2 UserID:1 Price:19.9}]}
	for _, order := range user2.Orders {
		// {ID:1 UserID:1 Price:9.9}
		// {ID:2 UserID:1 Price:19.9}
		fmt.Printf("Order: %+v\n", order)
	}
}
