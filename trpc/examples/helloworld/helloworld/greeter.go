package main

import (
	"context"

	pb "github.com/fengzhongzhu1621/xgo/trpc/trpcprotocol/helloworld"
)

type greeterImpl struct {
	pb.UnimplementedGreeter
}

func (s *greeterImpl) SayHello(
	ctx context.Context,
	req *pb.HelloRequest,
) (*pb.HelloReply, error) {
	rsp := &pb.HelloReply{
		Msg: "reply",
	}

	// 调用下游服务错误的时，如果返回 error，rsp 的内容将不再被返回
	return rsp, nil
}
