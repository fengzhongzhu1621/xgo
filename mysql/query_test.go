package mysql

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
)

func TestQueryConnect(t *testing.T) {
	// create database xgo charset utf8mb4;
	// 创建数据库连接 用户名:密码@协议(IP:port)/数据库名?a=xxx&b=xxx
	// 创建一个默认的连接池
	conn, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/xgo?parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("gorm.Open err: ", err)
		return
	}

	defer conn.Close()

	conn.SingularTable(true)

	// 初始化数据
	var stu Student
	stu.Name = "username_1"
	stu.Age = 10
	fmt.Println(conn.Create(&stu).Error) // <nil>

	var stu2 Student
	stu.Name = "username_2"
	stu.Age = 20
	fmt.Println(conn.Create(&stu2).Error) // <nil>

	// select * from student order by id limit 1
	var stu3 Student
	conn.First(&stu3)
	fmt.Println(stu3) // {username_1 10 1}

	// select * from student order by id desc limit 1
	var stu4 Student
	conn.Last(&stu4)
	fmt.Println(stu4) // {username_1 10 1}

	// 获取表中的所有数据
	// select name, age from student;
	var stu5 []Student
	conn.Select("name", "age").Find(&stu5)
	fmt.Println(stu5) // {username_1 10 1}
}
