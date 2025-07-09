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
	rsp := &pb.HelloReply{}
	return rsp, nil
}
