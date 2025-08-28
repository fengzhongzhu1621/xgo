package gorm

import (
	"errors"
	"fmt"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/optimisticlock"
)

// CREATE TABLE `product_optimistic_locks` (
//
//	`id` bigint unsigned NOT NULL AUTO_INCREMENT,
//	`name` longtext,
//	`quantity` bigint DEFAULT NULL,
//	`version` bigint DEFAULT NULL,
//	PRIMARY KEY (`id`)
//
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
type ProductOptimisticLock struct {
	ID       uint
	Name     string
	Quantity int
	Version  optimisticlock.Version // 关键：乐观锁版本字段[3](@ref)[6](@ref)
}

type UserOptimisticlock struct {
	ID      int
	Name    string
	Age     uint
	Version optimisticlock.Version
}

func ReduceInventory(db *gorm.DB, productID uint, quantity int) error {
	maxRetries := 3 // 最大重试次数
	for range maxRetries {
		var product ProductOptimisticLock
		// 1. 查询最新数据（含当前版本号）
		// ELECT * FROM `product_optimistic_locks` WHERE `product_optimistic_locks`.`id` = 1 ORDER BY `product_optimistic_locks`.`id` LIMIT 1
		if err := db.First(&product, productID).Error; err != nil {
			return err
		}

		// 2. 检查库存是否充足
		if product.Quantity < quantity {
			return errors.New("库存不足")
		}

		// 3. 更新库存（自动检查版本号）
		// UPDATE `product_optimistic_locks` SET `quantity`=90,`version`=`version`+1 WHERE version = 1 AND `product_optimistic_locks`.`version` = 1 AND `id` = 1
		result := db.Model(&product).
			Where("version = ?", product.Version). // 关键：WHERE 条件检查版本[3](@ref)[4](@ref)
			UpdateColumns(map[string]interface{}{
				"quantity": product.Quantity - quantity,
				"version":  gorm.Expr("version + 1"), // 版本号+1
			})

		// 4. 处理更新结果
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 { // 冲突检测：影响行数为0表示版本冲突[3](@ref)[4](@ref)
			continue // 触发重试
		}
		return nil // 更新成功
	}
	return errors.New("操作冲突，重试失败")
}

func TestProductOptimisticLock(t *testing.T) {
	// 建立连接
	dsn := "root:@tcp(127.0.0.1:3306)/xgo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 输出详细日志
	})
	if err != nil {
		panic("failed to connect database")
	}

	// 创建测试商品
	// INSERT INTO `product_optimistic_locks` (`name`,`quantity`,`version`) VALUES ('iPhone',100,1)
	db.Create(&ProductOptimisticLock{Name: "iPhone", Quantity: 100})

	// 并发减少库存（模拟10个并发请求）
	err = ReduceInventory(db, 1, 10) // 每次减少10个库存
	if err != nil {
		fmt.Println("操作失败:", err)
	}
}

func TestUserOptimisticlock(t *testing.T) {
	dsn := "root:@tcp(127.0.0.1:3306)/xgo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 输出详细日志
	})
	if err != nil {
		panic("failed to connect database")
	}

	// SELECT * FROM `user_optimisticlocks` ORDER BY `user_optimisticlocks`.`id` LIMIT 1
	var user UserOptimisticlock
	db.First(&user)

	db.Model(&user).Update("age", 18)
	// UPDATE `users` SET `age`=18,`version`=`version`+1 WHERE `users`.`version` = 1 AND `id` = 1

	// Ignoring the optimistic lock check.
	db.Unscoped().Model(&user).Update("age", 18)
	// UPDATE `users` SET `age`=18,`version`=`version`+1 WHERE `id` = 1

	// Ignoring the passed Version value.
	db.Model(&user).Updates(&UserOptimisticlock{Age: 18, Version: optimisticlock.Version{Int64: 1}})
	// UPDATE `users` SET `age`=18,`version`=`version`+1 WHERE `users`.`version` = 3 AND `id` = 1

	// If the Model's Version value is zero, Without considering optimistic lock check.
	db.Model(&UserOptimisticlock{}).Where("id = 1").Update("age", 12)
	// UPDATE `users` SET `age`=12,`version`=`version`+1 WHERE id = 1

	// When want to use GORM's Save method, need to call Select. Otherwise, will return primary key duplicate error.
	// The Select param is the fields that you want to update or "*".
	// INSERT INTO `user_optimisticlocks` (`name`,`age`,`version`) VALUES ('',18,1)
	db.Model(&user).Select("*").Save(&UserOptimisticlock{Age: 18})
}
