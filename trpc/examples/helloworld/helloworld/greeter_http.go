package main

import (
	"context"
	"net/http"

	pb "github.com/fengzhongzhu1621/xgo/trpc/trpcprotocol/helloworld"
	"trpc.group/trpc-go/trpc-go/errs"
	thttp "trpc.group/trpc-go/trpc-go/http"
	"trpc.group/trpc-go/trpc-go/log"
)

type greeterHttpImpl struct {
	pb.UnimplementedGreeterHttp
}

func (s *greeterHttpImpl) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	head := thttp.Head(ctx)

	// 判断请求报文是否为泛 http 协议
	if head == nil {
		// 使用业务自定义错误码
		return nil, errs.New(10000, "not http request")
	}

	// 获取请求报文头里的 request 字段
	reqHead := head.Request.Header.Get("request")
	// 获取请求报文头里的 Cookie 字段
	cookieStr := head.Request.Header.Get("Cookie")

	log.InfoContextf(ctx, "Msg: %s, reqHead: %s, cookie is: %s", req.Msg, reqHead, cookieStr)

	rsp := &pb.HelloReply{
		Msg: "Hello, World!",
	}

	// 为响应报文设置 Cookie
	cookie := &http.Cookie{Name: "admin", Value: "admin", HttpOnly: false}
	http.SetCookie(head.Response, cookie)

	// 为响应报文头添加 reply 字段
	head.Response.Header().Add("reply", "for test")

	return rsp, nil
}
