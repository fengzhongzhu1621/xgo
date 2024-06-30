package mysql

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
)


func TestUpdate(t *testing.T) {
	conn, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/xgo?parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("gorm.Open err: ", err)
		return
	}

	defer conn.Close()

	// 创建非复数表名，默认创建复数表名
	var stu Student
	stu.Name = "bob"
	stu.Age = 10
	stu.Id = 1

	// 创建非复数表名，默认创建复数表名
	// 必须使用，需要和表创建时的参数保持一致
	conn.SingularTable(true)

	// Save() 根据主键更新，如果数据没有指定主键，则为 insert 操作
	fmt.Println(conn.Save(&stu).Error)

	// 更新一个字段
	fmt.Println(conn.Model(new(Student)).Where("name = ?", "username_a").Update("age", 11).Error)

	// 更新多个字段
	fmt.Println(conn.Model(new(Student)).Where("name = ?", "username_a").Updates(map[string]interface{} {
		"name": "bob",
		"age": 11,
	}).Error)
}
