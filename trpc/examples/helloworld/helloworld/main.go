package main

import (
	pb "github.com/fengzhongzhu1621/xgo/trpc/trpcprotocol/helloworld"
	_ "trpc.group/trpc-go/trpc-filter/debuglog"
	_ "trpc.group/trpc-go/trpc-filter/recovery"
	trpc "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
)

func main() {
	// 创建一个服务对象，底层会自动读取服务配置及初始化插件，
	// 必须放在 main 函数首行，
	// 业务初始化逻辑必须放在 NewServer 后面
	s := trpc.NewServer()

	// 注册当前实现到服务对象中
	// trpc.examples.helloworld.Greeter trpc_go.yaml中配置的路由标识
	pb.RegisterGreeterService(s.Service("trpc.examples.helloworld.Greeter"), &greeterImpl{})

	// 启动服务，并阻塞
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
