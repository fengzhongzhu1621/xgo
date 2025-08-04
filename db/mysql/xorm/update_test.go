package xorm

import (
	"fmt"
	"testing"

	"xorm.io/xorm"
)

// NOTE:
// 更新通过engine.Update()实现，可以传入结构指针或map[string]interface{}。对于传入结构体指针的情况，xorm只会更新非空的字段。
// 如果一定要更新空字段，需要使用Cols()方法显示指定更新的列。使用Cols()方法指定列后，即使字段为空也会更新

// updateUserAge 更新单个字段
func updateUserAge(engine *xorm.Engine, userId int64, newAge int) error {
	affected, err := engine.ID(userId).Update(&XormUser4{Age: newAge})
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("没有找到 ID 为 %d 的用户", userId)
	}
	return nil
}

// updateUserDetails 更新多个字段
func updateUserDetails(engine *xorm.Engine, userId int64, newName string, newEmail string) error {
	affected, err := engine.ID(userId).
		Cols("Name", "Email").
		Update(&XormUser4{Name: newName, Email: newEmail})
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("没有找到 ID 为 %d 的用户", userId)
	}
	return nil
}

// updateUserByEmail 使用条件更新
func updateUserByEmail(engine *xorm.Engine, email string, newAge int) error {
	affected, err := engine.Where("email = ?", email).Update(&XormUser4{Age: newAge})
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("没有找到 Email 为 %s 的用户", email)
	}
	return nil
}

func TestXormUpdate(t *testing.T) {
	dbClient := GetDefaultXormDBClient()
	engine := dbClient.DB
	defer engine.Close()

	_ = engine.Sync2(new(XormUser4))
	affected, err := engine.Where("1=1").Delete(&XormUser4{})

	// 插入单个用户
	newUser := &XormUser4{
		Name:  "username_a",
		Age:   30,
		Email: "a@example.com",
	}
	affected, _ = engine.Insert(newUser)
	fmt.Printf("插入了 %d 条记录，新用户的 ID 是 %d\n", affected, newUser.Id)

	// 更新用户年龄
	err = updateUserAge(engine, newUser.Id, 31)
	if err != nil {
		fmt.Printf("更新用户年龄失败: %v", err)
	} else {
		fmt.Println("用户年龄更新成功")
	}

	// 更新用户名和邮箱
	err = updateUserDetails(engine, newUser.Id, "bob", "bob@example.com")
	if err != nil {
		fmt.Printf("更新用户详情失败: %v", err)
	} else {
		fmt.Println("用户详情更新成功")
	}

	// 根据邮箱更新年龄
	err = updateUserByEmail(engine, "bob@example.com", 32)
	if err != nil {
		fmt.Printf("根据邮箱更新用户年龄失败: %v", err)
	} else {
		fmt.Println("根据邮箱更新用户年龄成功")
	}
}
