package gorm

import (
	"fmt"
	"log"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// CREATE TABLE `user_many_to_manies` (
//
//	`id` bigint unsigned NOT NULL AUTO_INCREMENT,
//	`name` longtext,
//	PRIMARY KEY (`id`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
type UserManyToMany struct {
	ID        uint
	Name      string
	Languages []Language `gorm:"many2many:user_languages;"` // 通过标签表示多对多关系
}

// CREATE TABLE `languages` (
//
//	`id` bigint unsigned NOT NULL AUTO_INCREMENT,
//	`name` longtext,
//	PRIMARY KEY (`id`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
type Language struct {
	ID   uint
	Name string
}

// CREATE TABLE `user_languages` (
//   `user_many_to_many_id` bigint unsigned NOT NULL,
//   `language_id` bigint unsigned NOT NULL,
//   PRIMARY KEY (`user_many_to_many_id`,`language_id`),
//   KEY `fk_user_languages_language` (`language_id`),
//   CONSTRAINT `fk_user_languages_language` FOREIGN KEY (`language_id`) REFERENCES `languages` (`id`),
//   CONSTRAINT `fk_user_languages_user_many_to_many` FOREIGN KEY (`user_many_to_many_id`) REFERENCES `user_many_to_manies` (`id`)
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4

func TestManyToMany(t *testing.T) {
	dsn := "root:@tcp(127.0.0.1:3306)/xgo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 输出详细日志
	})
	if err != nil {
		log.Fatal(err)
	}

	// INSERT INTO `user_many_to_manies` (`name`) VALUES ('bob')
	// INSERT INTO `languages` (`name`) VALUES ('English')
	// INSERT INTO `languages` (`name`) VALUES ('Chinese')
	// INSERT INTO `languages` (`name`,`id`) VALUES ('English',3),('Chinese',4) ON DUPLICATE KEY UPDATE `id`=`id`
	// INSERT INTO `user_languages` (`user_many_to_many_id`,`language_id`) VALUES (2,3),(2,4) ON DUPLICATE KEY UPDATE `user_many_to_many_id`=`user_many_to_many_id`
	user := UserManyToMany{Name: "bob"}
	language1 := Language{Name: "English"}
	language2 := Language{Name: "Chinese"}
	db.Create(&user)
	db.Create(&language1)
	db.Create(&language2)
	db.Model(&user).Association("Languages").Append([]Language{language1, language2})

	// SELECT * FROM `user_languages` WHERE `user_languages`.`user_many_to_many_id` = 1
	// SELECT * FROM `languages` WHERE `languages`.`id` IN (1,2)
	// SELECT * FROM `user_many_to_manies` WHERE `user_many_to_manies`.`id` = 1 ORDER BY `user_many_to_manies`.`id` LIMIT 1
	var user2 UserManyToMany
	db.Preload("Languages").First(&user2, 1)
	fmt.Printf(
		"User: %+v\n",
		user2,
	) // {ID:1 Name:Charlie Languages:[{ID:1 Name:English} {ID:2 Name:Chinese}]}
	for _, lang := range user2.Languages {
		// Language: {ID:1 Name:English}
		// Language: {ID:2 Name:Chinese}
		fmt.Printf("Language: %+v\n", lang)
	}
}
