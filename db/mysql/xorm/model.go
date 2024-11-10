package xorm

import "time"

// CREATE TABLE `xorm_student` (
//
//	`name` varchar(64) NOT NULL,
//	`age` int NOT NULL DEFAULT '0',
//	`id` bigint NOT NULL AUTO_INCREMENT,
//	PRIMARY KEY (`id`),
//	UNIQUE KEY `UQE_xorm_student_name` (`name`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
type XormStudent struct {
	Name string `xorm:"not null unique VARCHAR(64)"`
	Age  int    `xorm:"not null default 0 INT(10)"`
	Id   int64  `xorm:"bigint(20) pk autoincr"`
}

// CREATE TABLE `xorm_user` (
//
//	`id` bigint NOT NULL AUTO_INCREMENT,
//	`name` varchar(255) DEFAULT NULL,
//	PRIMARY KEY (`id`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
type XormUser struct {
	ID   int64  `xorm:"pk autoincr"` // 主键，自增
	Name string `xorm:"varchar(255)"`
}

// CREATE TABLE `xorm_post` (
//
//	`id` bigint NOT NULL AUTO_INCREMENT,
//	`title` varchar(255) DEFAULT NULL,
//	`user_id` bigint DEFAULT NULL,
//	PRIMARY KEY (`id`),
//	KEY `IDX_xorm_post_user_id` (`user_id`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
type XormPost struct {
	ID     int64    `xorm:"pk autoincr"` // 主键，自增
	Title  string   `xorm:"varchar(255)"`
	UserID int64    `xorm:"index"` // 用户 ID，用于关联查询
	User   XormUser `xorm:"-"`     // 不映射到数据库的字段，XormUser 和 XormPost 是一对多的关系
}

// CREATE TABLE `xorm_user2` (
//
//	`id` bigint NOT NULL AUTO_INCREMENT,
//	`usr_name` varchar(255) NOT NULL COMMENT 'NickName',
//	PRIMARY KEY (`id`),
//	UNIQUE KEY `UQE_xorm_user2_usr_name` (`usr_name`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
type XormUser2 struct {
	ID   int64  `xorm:"pk autoincr"` // 主键，自增
	Name string `xorm:"varchar(255) not null unique 'usr_name' comment('NickName')"`
}

// CREATE TABLE `xorm_user3` (
//
//	`id` bigint NOT NULL AUTO_INCREMENT,
//	`name` varchar(255) DEFAULT NULL,
//	`salt` varchar(255) DEFAULT NULL,
//	`age` int DEFAULT NULL,
//	`passwd` varchar(200) DEFAULT NULL,
//	`created` datetime DEFAULT NULL,
//	`updated` datetime DEFAULT NULL,
//	PRIMARY KEY (`id`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
type XormUser3 struct {
	Id      int64
	Name    string
	Salt    string
	Age     int
	Passwd  string    `xorm:"varchar(200)"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
}

// Using card with address like this doesnt seem to work with json b
// CREATE TABLE `xorm_card_s` (
//
//	`addr` text NOT NULL,
//	`id` varchar(255) NOT NULL,
//	`is_default` tinyint(1) DEFAULT NULL,
//	`nickname` text NOT NULL,
//	`number_last_4` text NOT NULL,
//	PRIMARY KEY (`id`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
type XormCardS struct {
	Addr        Address `xorm:"json notnull 'addr'" json:"addr"`
	Id          string  `xorm:"pk notnull 'id'" json:"id"`
	IsDefault   bool    `xorm:"'is_default'" json:"isDefault"`
	Nickname    string  `xorm:"text notnull 'nickname'" json:"nickname"`
	NumberLast4 string  `xorm:"text notnull 'number_last_4'" json:"numberLast4"`
}

// CREATE TABLE `xorm_card_m` (
//
//	`addr` text NOT NULL,
//	`id` varchar(255) NOT NULL,
//	`nickname` text NOT NULL,
//	`number_last_4` text NOT NULL,
//	PRIMARY KEY (`id`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
type XormCardM struct {
	Addr        map[string]interface{} `xorm:"json notnull 'addr'" json:"addr"`
	Id          string                 `xorm:"pk notnull 'id'" json:"id"`
	Nickname    string                 `xorm:"text notnull 'nickname'" json:"nickname"`
	NumberLast4 string                 `xorm:"text notnull 'number_last_4'" json:"numberLast4"`
}
