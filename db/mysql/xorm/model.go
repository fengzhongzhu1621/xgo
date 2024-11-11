package xorm

import (
	"time"
)

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

// CREATE TABLE `xorm_link` (
//
//	`id` int NOT NULL AUTO_INCREMENT,
//	`description` varchar(40) NOT NULL,
//	`url` varchar(80) NOT NULL,
//	`create_user` varchar(20) DEFAULT 'SYSTEM',
//	`create_date` datetime DEFAULT NULL,
//	`update_user` varchar(20) DEFAULT 'SYSTEM',
//	`update_date` datetime DEFAULT NULL,
//	`version` int DEFAULT NULL,
//	PRIMARY KEY (`id`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
type XormLink struct {
	Id          int       `form:"id" xorm:"int(3) pk not null autoincr"`
	Description string    `form:"description" xorm:"varchar(40) not null"`
	Url         string    `form:"url" xorm:"varchar(80) not null"`
	CreateUser  string    `xorm:"varchar(20) default 'SYSTEM'"`
	CreateDate  time.Time `xorm:"datetime created"`
	UpdateUser  string    `xorm:"varchar(20) default 'SYSTEM'"`
	UpdateDate  time.Time `xorm:"datetime updated"`
	Version     int       `form:"version" xorm:"int(11) version"`
	Page        `xorm:"-"`
}

// CREATE TABLE `xorm_shop` (
//
//	`id` bigint NOT NULL AUTO_INCREMENT,
//	`name` varchar(12) DEFAULT NULL,
//	`promotion_info` varchar(30) DEFAULT NULL,
//	`address` varchar(100) DEFAULT NULL,
//	`phone` varchar(11) DEFAULT NULL,
//	`status` tinyint DEFAULT NULL,
//	`longitude` double DEFAULT NULL,
//	`latitude` double DEFAULT NULL,
//	`image_path` varchar(255) DEFAULT NULL,
//	`is_new` tinyint(1) DEFAULT NULL,
//	`is_premium` tinyint(1) DEFAULT NULL,
//	`rating` float DEFAULT NULL,
//	`rating_count` int DEFAULT NULL,
//	`recent_order_num` int DEFAULT NULL,
//	`minimum_order_amount` int DEFAULT NULL,
//	`delivery_fee` int DEFAULT NULL,
//	`opening_hours` varchar(20) DEFAULT NULL,
//	PRIMARY KEY (`id`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
type XormShop struct {
	//id
	Id int64 `xorm:"pk autoincr" json:"id"`
	//商铺名称
	Name string `xorm:"varchar(12)" json:"name"`
	//宣传信息
	PromotionInfo string `xorm:"varchar(30)" json:"promotion_info"`
	//地址
	Address string `xorm:"varchar(100)" json:"address"`
	//联系电话
	Phone string `xorm:"varchar(11)" json:"phone"`
	//店铺营业状态
	Status int `xorm:"tinyint" json:"status"`

	//经度
	Longitude float64 `xorm:"" json:"longitude"`
	//纬度
	Latitude float64 `xorm:"" json:"latitude"`
	//店铺图标
	ImagePath string `xorm:"varchar(255)" json:"image_path"`

	IsNew     bool `xorm:"bool" json:"is_new"`
	IsPremium bool `xorm:"bool" json:"is_premium"`

	//商铺评分
	Rating float32 `xorm:"float" json:"rating"`
	//评分总数
	RatingCount int64 `xorm:"int" json:"rating_count"`
	//当前订单总数
	RecentOrderNum int64 `xorm:"int" json:"recent_order_num"`

	//配送起送价
	MinimumOrderAmount int32 `xorm:"int" json:"minimum_order_amount"`
	//配送费
	DeliveryFee int32 `xorm:"int" json:"delivery_fee"`

	//营业时间
	OpeningHours string `xorm:"varchar(20)" json:"opening_hours"`
}

// CREATE TABLE `xorm_db_device` (
//
//	`udid` varchar(255) NOT NULL,
//	`name` varchar(255) DEFAULT NULL,
//	`custom_name` varchar(255) DEFAULT NULL,
//	`provider_id` bigint DEFAULT NULL,
//	`json_info` varchar(255) DEFAULT NULL,
//	`width` int DEFAULT NULL,
//	`height` int DEFAULT NULL,
//	`click_width` int DEFAULT NULL,
//	`click_height` int DEFAULT NULL,
//	`wda_port` int DEFAULT NULL,
//	PRIMARY KEY (`udid`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
type XormDbDevice struct {
	Udid        string `xorm:"pk"`
	Name        string
	CustomName  string
	ProviderId  int64
	JsonInfo    string
	Width       int
	Height      int
	ClickWidth  int
	ClickHeight int
	Ready       string `xorm:"-"`
	WdaPort     int
}

// CREATE TABLE `xorm_app` (
//
//	`id` int unsigned NOT NULL,
//	`name_space` varchar(255) DEFAULT NULL,
//	`is_del` tinyint(1) DEFAULT NULL,
//	`version` int DEFAULT NULL,
//	`created_at` datetime DEFAULT NULL,
//	`updated_at` datetime DEFAULT NULL,
//	`deleted_at` datetime DEFAULT NULL,
//	PRIMARY KEY (`id`),
//	UNIQUE KEY `UQE_xorm_app_name_space` (`name_space`),
//	KEY `IDX_xorm_app_deleted_at` (`deleted_at`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
type XormApp struct {
	Id        uint       `xorm:"pk"`
	NameSpace string     `xorm:"unique"`
	IsDel     bool       `xorm:"index default 0"`
	Version   int        `xorm:"version"`
	CreatedAt time.Time  `xorm:"created"`
	UpdatedAt time.Time  `xorm:"updated"`
	DeletedAt *time.Time `xorm:"deleted index"`

	Operations []XormOperation `xorm:"- extends"`
	Traffics   []XormTraffic   `xorm:"- extends"`
}

// CREATE TABLE `xorm_operation` (
//
//	`id` int unsigned NOT NULL,
//	`app_id` int unsigned DEFAULT NULL,
//	`end_point` varchar(255) DEFAULT NULL,
//	`is_del` tinyint(1) DEFAULT NULL,
//	`version` int DEFAULT NULL,
//	`created_at` datetime DEFAULT NULL,
//	`updated_at` datetime DEFAULT NULL,
//	`deleted_at` datetime DEFAULT NULL,
//	PRIMARY KEY (`id`),
//	KEY `IDX_xorm_operation_app_id` (`app_id`),
//	KEY `IDX_xorm_operation_deleted_at` (`deleted_at`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
type XormOperation struct {
	Id        uint `xorm:"pk"`
	AppId     uint `xorm:"index"`
	EndPoint  string
	IsDel     bool       `xorm:"index default 0"`
	Version   int        `xorm:"version"`
	CreatedAt time.Time  `xorm:"created"`
	UpdatedAt time.Time  `xorm:"updated"`
	DeletedAt *time.Time `xorm:"deleted index"`

	App XormApp `xorm:"- extends"`
}

// CREATE TABLE `xorm_traffic` (
//
//	`id` int unsigned NOT NULL AUTO_INCREMENT,
//	`app_id` int unsigned DEFAULT NULL,
//	`unit` varchar(10) DEFAULT 'd',
//	`val` int unsigned DEFAULT NULL,
//	`seq` int unsigned DEFAULT NULL,
//	`created_at` datetime DEFAULT NULL,
//	`updated_at` datetime DEFAULT NULL,
//	`deleted_at` datetime DEFAULT NULL,
//	PRIMARY KEY (`id`),
//	KEY `IDX_xorm_traffic_with_seq` (`app_id`,`seq`),
//	KEY `IDX_xorm_traffic_app_id` (`app_id`),
//	KEY `IDX_xorm_traffic_unit` (`unit`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
type XormTraffic struct {
	Id        uint   `xorm:"pk autoincr"`
	AppId     uint   `xorm:"index index(with_seq)"`
	Unit      string `xorm:"varchar(10) index default 'd'"` //트래픽 단위(min:분, hour:시간, day:1일, month:1달)
	Val       uint
	Seq       uint      `xorm:"index(with_seq)"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt *time.Time

	App XormApp `xorm:"- extends"`
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

type XormCategory struct {
	ID      int64       `json:"id" xorm:"bigint(20) pk autoincr"`
	Name    string      `json:"name" xorm:"not null unique VARCHAR(32)"`
	Entries []XormEntry `json:"entries" xorm:"- extends"`

	BaseTimeModelWithoutSoftDelete `xorm:"extends"`
}

// XormEntry 条目
type XormEntry struct {
	CategoryID int64        `json:"categoryID" xorm:"not null"`
	Category   XormCategory `json:"category" xorm:"- extends"`

	ID    int64   `json:"id" xorm:"pk autoincr"`
	Name  string  `json:"name" xorm:"varchar(64) unique not null"`
	Desc  string  `json:"desc" xorm:"text null"`
	Price float32 `json:"price" xorm:"not null"`

	BaseTimeModelWithoutSoftDelete `xorm:"extends"`
}

// XormTask 后台任务
// CREATE TABLE `xorm_task` (
//
//	`id` bigint NOT NULL AUTO_INCREMENT,
//	`name` varchar(128) NOT NULL,
//	`args` text,
//	`result` text,
//	`started_at` datetime DEFAULT NULL,
//	`duration` bigint DEFAULT NULL,
//	`created_at` datetime(6) NOT NULL,
//	`created_by` varchar(32) NOT NULL,
//	`updated_at` datetime(6) NOT NULL,
//	`updated_by` varchar(32) NOT NULL,
//	PRIMARY KEY (`id`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
type XormTask struct {
	ID        int64                  `json:"id"  xorm:"pk autoincr"`
	Name      string                 `json:"name" xorm:"varchar(128) not null"`
	Args      map[string]interface{} `json:"args" xorm:"json"`
	Result    map[string]interface{} `json:"result" xorm:"json"`
	StartedAt time.Time              `json:"startedAt" xorm:"datetime default null"`
	Duration  time.Duration          `json:"duration" xorm:"bigint default null"`

	BaseTimeModelWithoutSoftDelete `xorm:"extends"`
}

// XormPeriodicTask 周期任务
// CREATE TABLE `xorm_periodic_task` (
//
//	`id` bigint NOT NULL AUTO_INCREMENT,
//	`cron` varchar(32) NOT NULL,
//	`name` varchar(128) NOT NULL,
//	`args` text,
//	`enabled` tinyint(1) NOT NULL DEFAULT '1',
//	`created_at` datetime(6) NOT NULL,
//	`created_by` varchar(32) NOT NULL,
//	`updated_at` datetime(6) NOT NULL,
//	`updated_by` varchar(32) NOT NULL,
//	PRIMARY KEY (`id`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
type XormPeriodicTask struct {
	ID      int64                  `json:"id" xorm:"pk autoincr"`
	Cron    string                 `json:"cron" xorm:"varchar(32) not null"`
	Name    string                 `json:"name" xorm:"varchar(128) not null"`
	Args    map[string]interface{} `json:"args" xorm:"json"`
	Enabled bool                   `json:"enabled" xorm:"not null default true"`

	BaseTimeModelWithoutSoftDelete `xorm:"extends"`
}
