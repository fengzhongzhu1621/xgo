package gorm

import (
	"fmt"
	"log"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// CREATE TABLE `user_has_ones` (
//
//	`id` bigint unsigned NOT NULL AUTO_INCREMENT,
//	`name` longtext,
//	PRIMARY KEY (`id`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
type UserHasOne struct {
	ID      uint
	Name    string
	Profile Profile `gorm:"foreignKey:UserID;references:ID"` // 显式定义 HasOne 关系，Profile的外键字段是UserID，与UserHasOne的ID关联
}

// CREATE TABLE `profiles` (
//
//	`id` bigint unsigned NOT NULL AUTO_INCREMENT,
//	`user_id` bigint unsigned DEFAULT NULL,
//	`bio` longtext,
//	PRIMARY KEY (`id`),
//	KEY `fk_user_has_ones_profile` (`user_id`),
//	CONSTRAINT `fk_user_has_ones_profile` FOREIGN KEY (`user_id`) REFERENCES `user_has_ones` (`id`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
type Profile struct {
	ID     uint
	UserID uint // 外键，关联到 UserHasOne 的 ID
	Bio    string
}

func TestOneToOne(t *testing.T) {
	dsn := "root:@tcp(127.0.0.1:3306)/xgo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 输出详细日志
	})
	if err != nil {
		log.Fatal(err)
	}

	// INSERT INTO `user_has_ones` (`name`) VALUES ('Alice')
	// INSERT INTO `profiles` (`user_id`,`bio`) VALUES (1,'Hello, I am Alice.') ON DUPLICATE KEY UPDATE `user_id`=VALUES(`user_id`)
	user := UserHasOne{
		Name: "Alice",
		Profile: Profile{
			Bio: "Hello, I am Alice.",
		},
	}
	db.Create(&user)

	// 预加载
	var user2 UserHasOne
	// SELECT * FROM `profiles` WHERE `profiles`.`user_id` = 1
	// SELECT * FROM `user_has_ones` WHERE `user_has_ones`.`id` = 1 ORDER BY `user_has_ones`.`id` LIMIT 1
	db.Preload("Profile").First(&user2, 1) // 预加载Profile信息
	fmt.Printf(
		"User: %+v\n",
		user2,
	) // {ID:1 Name:Alice Profile:{ID:1 UserID:1 Bio:Hello, I am Alice.}}
	fmt.Printf("Profile: %+v\n", user2.Profile) // {ID:1 UserID:1 Bio:Hello, I am Alice.}
}
