package main

import (
	"context"

	pb "github.com/fengzhongzhu1621/xgo/trpc/trpcprotocol/helloworld"
	trpc "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
)

// greeterImpl GreeterService接口的实现
type greeterImpl struct {
	pb.UnimplementedGreeter
}

func (s *greeterImpl) SayHello(
	ctx context.Context,
	req *pb.HelloRequest,
) (*pb.HelloReply, error) {
	// 获得客户端的透传字段
	value1 := trpc.GetMetaData(ctx, "key1")
	log.Infof("GetMetaData key1: %s", string(value1))

	// 设置透传字段返回给上游调用方
	trpc.SetMetaData(ctx, "key1", []byte("val1"))

	rsp := &pb.HelloReply{
		Msg: "Hello, World!",
	}

	// 调用下游服务错误的时，如果返回 error，rsp 的内容将不再被返回
	return rsp, nil
}
