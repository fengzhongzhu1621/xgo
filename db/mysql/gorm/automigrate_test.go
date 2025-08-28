package gorm

import (
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// CREATE TABLE `products` (
//   `id` bigint unsigned NOT NULL AUTO_INCREMENT,
//   `created_at` datetime(3) DEFAULT NULL,
//   `updated_at` datetime(3) DEFAULT NULL,
//   `deleted_at` datetime(3) DEFAULT NULL,
//   `code` longtext,
//   `price` bigint unsigned DEFAULT NULL,
//   PRIMARY KEY (`id`),
//   KEY `idx_products_deleted_at` (`deleted_at`)
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func TestAutoMigrate(t *testing.T) {
	// 建立连接
	dsn := "root:@tcp(127.0.0.1:3306)/xgo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 输出详细日志
	})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	// 注意：必须调用Set设置 CHARSET 和 COLLATE；dsn中配置的charset和collation 对AutoMigrate不生效
	db = db.Set("gorm:table_options", "CHARSET=utf8mb4")
	db.AutoMigrate(&Product{})
	db.AutoMigrate(&Student{})

	db.AutoMigrate(&UserHasOne{})
	db.AutoMigrate(&Profile{})

	db.AutoMigrate(&UserHasMany{})
	db.AutoMigrate(&Order{})

	db.AutoMigrate(&UserManyToMany{})
	db.AutoMigrate(&Language{})

	db.AutoMigrate(&ProductOptimisticLock{})
	db.AutoMigrate(&UserOptimisticlock{})

	db.AutoMigrate(&Company{})
	db.AutoMigrate(&UserWithCompany{})
}
