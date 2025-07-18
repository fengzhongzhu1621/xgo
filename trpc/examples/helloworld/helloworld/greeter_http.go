package main

import (
	"context"
	"encoding/json"
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

	// 默认打包序列化方式和请求时的“Content-Type”保持一致
	// msg := trpc.Message(ctx)
	// msg.WithSerializationType(codec.SerializationTypeJSON)

	// 服务内调用 client
	// 创建一个客户端调用代理，该操作很轻量不会创建连接，可以每次请求创建，也可以全局初始化一个 proxy
	// proxy 不要每次创建，这里只是演示
	proxy := pb.NewGreeterClientProxy(
	//client.WithTarget("ip://127.0.0.1:8001"),
	//client.WithProtocol("trpc"),
	)
	reply, err := proxy.SayHello(ctx, &pb.HelloRequest{
		Msg: "hello",
	})
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	rsp := &pb.HelloReply{
		Msg: "NewGreeterClientProxy: " + reply.Msg,
	}

	// 处理mysql
	handleGormMysql(ctx, req)

	// 处理redis
	handleRedis(ctx, req)

	// 为响应报文设置 Cookie
	cookie := &http.Cookie{Name: "admin", Value: "admin", HttpOnly: false}
	http.SetCookie(head.Response, cookie)

	// 为响应报文头添加 reply 字段
	head.Response.Header().Add("reply", "for test")

	return rsp, nil
}

type Response struct {
	Result  bool            `json:"result"`
	Code    int32           `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// 自定义返回数据处理函数
func init() {
	thttp.DefaultServerCodec.RspHandler = func(w http.ResponseWriter, r *http.Request, rspbody []byte) error {
		if len(rspbody) == 0 {
			return nil
		}
		bs, err := json.Marshal(&Response{Result: true, Code: 0, Message: "OK", Data: rspbody})
		if err != nil {
			return err
		}
		_, err = w.Write(bs)
		return err
	}
}

// 自定义错误处理函数
func init() {
	thttp.DefaultServerCodec.ErrHandler = func(w http.ResponseWriter, r *http.Request, e *errs.Error) {
		// 填充指定格式错误信息到 HTTP Body
		bs, err := json.Marshal(&Response{Result: false, Code: int32(e.Code), Message: e.Msg})
		if err != nil {
			return
		}
		_, err = w.Write(bs)
		//w.Write([]byte(fmt.Sprintf(`{"retcode": %d, "retmsg": "%s"}`, e.Code, e.Msg)))
	}
}
