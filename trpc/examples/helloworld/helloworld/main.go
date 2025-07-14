package main

import (
	pb "github.com/fengzhongzhu1621/xgo/trpc/trpcprotocol/helloworld"
	_ "trpc.group/trpc-go/trpc-filter/debuglog"
	_ "trpc.group/trpc-go/trpc-filter/recovery"
	_ "trpc.group/trpc-go/trpc-filter/validation"
	trpc "trpc.group/trpc-go/trpc-go"

	// thttp "trpc.group/trpc-go/trpc-go/http"
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
	pb.RegisterGreeterHttpService(s.Service("trpc.examples.helloworld.GreeterHttp"), &greeterHttpImpl{})

	// 泛 HTTP 标准服务
	// 1. URL 注册模式
	// func HandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request) error)
	// func RegisterNoProtocolService(s server.Service)
	//
	// 2. Mux 注册模式
	// func RegisterNoProtocolServiceMux(s server.Service, mux http.Handler)

	// 启动服务，并阻塞
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
