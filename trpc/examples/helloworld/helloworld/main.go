package main

import (
	pb "github.com/fengzhongzhu1621/xgo/trpc/trpcprotocol/helloworld"
	"trpc.group/trpc-go/trpc-database/kafka"
	"trpc.group/trpc-go/trpc-database/timer"
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

	// ///////////////////////////////////////////////////////////////////////////////////////////////////////
	// timer
	// 注册本地定时器服务
	timer.RegisterHandlerService(s.Service("trpc.examples.helloworld.time_local"), handleLocalTimer)

	// 注册调度策略
	timer.RegisterScheduler("use_redis", &scheduler{})
	// 注册分布式互斥定时器服务
	timer.RegisterHandlerService(s.Service("trpc.examples.helloworld.time_distributed"), handleDistributedTimer)

	///////////////////////////////////////////////////////////////////////////////////////////////////////
	// kafka
	// 注册一个本地定时服务，模拟生产者向kafka发送数据
	timer.RegisterHandlerService(s.Service("trpc.examples.helloworld.kafka_produer"), handleKafkaProducer)

	// If you're using a custom address configuration,
	// configure it before starting the server.
	/*
		cfg := kafka.GetDefaultConfig()
		cfg.ClientID = "newClientID"
		kafka.RegisterAddrConfig("address", cfg)
	*/
	// default service name is trpc.kafka.consumer.service
	// kafka.RegisterKafkaConsumerService(s, &Consumer{})
	// The parameter of s.Service should be the same as the name of the service in
	// the configuration file. The configuration is loaded based on this parameter.
	kafka.RegisterKafkaConsumerService(s.Service("trpc.examples.helloworld.kafka-consumer-1"), &Consumer{})

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
