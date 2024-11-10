package xorm

import (
	"fmt"
	"testing"
)

// TestDeleteAll 删除表中所有的数据
func TestDeleteAll(t *testing.T) {
	dbClient := GetDefaultXormDBClient()

	// 删除表中的所有数据，至少有一个条件
	affected, err := dbClient.DB.Where("1=1").Delete(&XormStudent{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d records deleted", affected)
}

// TestDeleteWhere 根据where 条件删除
func TestDeleteWhere(t *testing.T) {
	dbClient := GetDefaultXormDBClient()

	affected, _ := dbClient.DB.Where("name = ?", "lzy").Delete(&XormStudent{})
	fmt.Printf("%d records deleted", affected)
}

// TestDeleteById 根据主键 ID 删除
func TestDeleteById(t *testing.T) {
	dbClient := GetDefaultXormDBClient()

	affected, _ := dbClient.DB.ID(1).Delete(&XormStudent{})
	fmt.Printf("%d records deleted", affected)
}
