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
	"trpc.group/trpc-go/trpc-go/admin"
	"trpc.group/trpc-go/trpc-go/log"
)

func init() {
	// 开启trace级别的日志 或 使用环境变量设置 export TRPC_LOG_TRACE=1
	// log.EnableTrace()
}

func init() {
	// 注册自定义命令
	admin.HandleFunc("/cmds/custom", customAdminCmd)
}

func main() {
	// 注册自定义日志插件
	// 注意：plugin.Register 要在 trpc.NewServer 之前执行
	// plugin.Register("customlog", log.DefaultLogFactory)
	// 设置 ctx 的 Logger 为 custom
	// trpc.Message(ctx).WithLogger(log.Get("customlog"))
	// 使用 Context 类型的日志接口
	// log.InfoContext(ctx, "custom log msg")

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

	// 需要配置自己的参数
	// If you're using a custom address configuration,
	// configure it before starting the server.
	/*
		cfg := kafka.GetDefaultConfig()
		cfg.ClientID = "newClientID"
		kafka.RegisterAddrConfig("address", cfg)
	*/

	// 如何注入自定义配置 (远端配置)
	// 在trpc_go.yaml中配置fake_address，然后配合kafka.RegisterAddrConfig方法注入 trpc_go.yaml配置如下
	// address: fake_address
	//
	// 在服务启动前，注入自定义配置
	// func main() {
	//   s := trpc.NewServer()
	//   // 使用自定义 addr，需在启动 server 前注入
	//   cfg := kafka.GetDefaultConfig()
	//   cfg.Brokers = []string{"127.0.0.1:9092"}
	//   cfg.Topics = []string{"test_topic"}
	//   kafka.RegisterAddrConfig("fake_address", cfg)
	//   kafka.RegisterKafkaConsumerService(s, &Consumer{})
	//   s.Serve()
	// }

	// default service name is trpc.kafka.consumer.service
	// kafka.RegisterKafkaConsumerService(s, &Consumer{})
	// The parameter of s.Service should be the same as the name of the service in
	// the configuration file. The configuration is loaded based on this parameter.
	// 启动多个消费者的情况，可以配置多个 service，然后这里任意匹配 kafka.RegisterHandlerService(s.Service("name"), handle)，
	// 没有指定 name 的情况，代表所有 service 共用同一个 handler
	kafka.RegisterKafkaConsumerService(s.Service("trpc.examples.helloworld.kafka-consumer-1"), &KafkaConsumer{})

	kafka.RegisterBatchHandlerService(s.Service("trpc.examples.helloworld.kafka-consumer-2"), kafkaBatchHandle)

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
