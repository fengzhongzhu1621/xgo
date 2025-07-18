// Package main is originally generated by trpc-cmdline v1.0.9.
// It is located at `project/cmd/client`.
// Run this file by executing `go run cmd/client/main.go` in the project directory.
package main

import (
	pb "github.com/fengzhongzhu1621/xgo/trpc/trpcprotocol/helloworld"
	_ "trpc.group/trpc-go/trpc-filter/debuglog"
	trpc "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/client"
	"trpc.group/trpc-go/trpc-go/log"
	trpcpb "trpc.group/trpc/trpc-protocol/pb/go/trpc"
)

func callGreeterSayHello() {
	proxy := pb.NewGreeterClientProxy(
	// client.WithTarget("ip://127.0.0.1:8001"),
	// client.WithProtocol("trpc"),
	)
	ctx := trpc.BackgroundContext()

	// 构造响应 header
	head := &trpcpb.ResponseProtocol{}

	// 使用 trpc-trans-info透传字段
	// client 透传数据到 server
	options := []client.Option{
		client.WithMetaData("key1", []byte("val1")),
		client.WithMetaData("key2", []byte("val2")),
		client.WithMetaData("key3", []byte("val3")),
		client.WithRspHead(head),
	}

	// Example usage of unary client.
	reply, err := proxy.SayHello(ctx, &pb.HelloRequest{
		Msg: "hello",
	}, options...)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	// 打印响应结构体
	log.Debugf("simple rpc receive: %+v", reply)

	// 打印下游服务器返回的透传字段
	log.Info("head.TransInfo = %+v", head)
}

func main() {
	// 读取本地默认配置文件
	// Load configuration following the logic in trpc.NewServer.
	cfg, err := trpc.LoadConfig(trpc.ServerConfigPath)
	if err != nil {
		panic("load config fail: " + err.Error())
	}

	// 设置为全局配置
	trpc.SetGlobalConfig(cfg)

	// 执行插件和初始化客户端
	if err := trpc.Setup(cfg); err != nil {
		panic("setup plugin fail: " + err.Error())
	}

	// 执行业务逻辑
	callGreeterSayHello()
}
