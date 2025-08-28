package main

import (
	"context"

	pb "github.com/fengzhongzhu1621/xgo/trpc/trpcprotocol/helloworld"
	"trpc.group/trpc-go/trpc-database/gorm"
	"trpc.group/trpc-go/trpc-go/log"
)

// https://gorm.io/zh_CN/docs/index.html

// User is the model struct.
type User struct {
	ID       int
	Username string
}

func handleGormMysql(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	// 简单使用方法，只支持mysql, gormDB是原生的gorm.DB类型指针
	// gormDB := gorm.NewClientProxy("trpc.mysql.test.test")
	// gormDB.Where("current_owners = ?", "xxxx").Where("id < ?", xxxx).Find(&owners)

	// proxy是并发安全的，可以在程序入口定义一个全局变量，如
	// var mysqlproxy = mysql.NewClientProxy("xxx")，
	// 也可以每次请求都 NewClientProxy。
	// 更推荐的做法是定义成 impl struct 里面的成员变量，方便依赖注入和 mock 测试。

	// 读写分离
	// reader := mysql.NewClientProxy("trpc.mysql.xgo.read")
	// writer := mysql.NewClientProxy("trpc.mysql.xgo.write")

	cli, err := gorm.NewClientProxy("trpc.mysql.server.service")
	if err != nil {
		panic(err)
	}

	// Create record
	insertUser := User{Username: "gorm-client"}
	result := cli.Create(&insertUser)
	log.Infof("inserted data's primary key: %d, err: %v", insertUser.ID, result.Error)

	// Query record
	var queryUser User
	if err := cli.First(&queryUser).Error; err != nil {
		panic(err)
	}
	log.Infof("query user: %+v", queryUser)

	// Delete record
	// WithContext(ctx) 让请求带 TraceID 和全链路超时信息
	deleteUser := User{ID: insertUser.ID}
	if err := cli.WithContext(ctx).Delete(&deleteUser).Error; err != nil {
		panic(err)
	}
	log.Info("delete record succeed")

	return nil, nil
}
