package main

import (
	"context"

	pb "github.com/fengzhongzhu1621/xgo/trpc/trpcprotocol/helloworld"
)

type greeterHttpImpl struct {
	pb.UnimplementedGreeterHttp
}

func (s *greeterHttpImpl) SayHello(
	ctx context.Context,
	req *pb.HelloRequest,
) (*pb.HelloReply, error) {
	rsp := &pb.HelloReply{
		Msg: "world",
	}
	return rsp, nil
}
