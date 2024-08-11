package gorm_1

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
)

var GlobalConn *gorm.DB

// 测试创建数据库连接，并设置最大连接数
func TestConnect(t *testing.T) {
	// create database xgo charset utf8mb4;
	// 创建数据库连接 用户名:密码@协议(IP:port)/数据库名?a=xxx&b=xxx
	// 创建一个默认的连接池
	// parseTime=True&loc=Local 访问数据库使用北京时区
	conn, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/xgo?parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("gorm.Open err: ", err)
		return
	}

	defer conn.Close()

	// 设置连接池初始的属性
	GlobalConn = conn
	GlobalConn.DB().SetMaxIdleConns(10)  // 初始连接数
	GlobalConn.DB().SetMaxOpenConns(100) // 最大连接数

}
