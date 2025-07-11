package main

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/fengzhongzhu1621/xgo/tests"
	pb "github.com/fengzhongzhu1621/xgo/trpc/trpcprotocol/helloworld"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	_ "trpc.group/trpc-go/trpc-filter/validation"
	trpc "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/client"
)

func Test_greeterImpl_SayHello(t *testing.T) {
	// Start writing mock logic.
	// 使用 GreeterService 的 mock 服务
	ctrl := gomock.NewController(t)
	// 确保在测试结束时调用 Finish()，它会检查是否有未验证的 Mock 期望（expectation），如果没有被调用会报错。
	defer ctrl.Finish()
	greeterService := pb.NewMockGreeterService(ctrl)

	// 定义一个切片，用于存放按顺序执行的 Mock 调用（InOrder 机制）。
	// 但目前这个切片是空的，所以 gomock.InOrder(inorderClient...) 实际上没有起到任何作用。
	var inorderClient []any

	// Expected behavior.
	// gomock.Any()：表示对参数不做具体限制，可以是任意值。
	// .AnyTimes()：表示这个方法可以被调用任意次数（0 次或多次）。
	m := greeterService.EXPECT().SayHello(gomock.Any(), gomock.Any()).AnyTimes()

	// 允许你在 Mock 方法被调用时执行自定义逻辑，并返回指定的值。
	m.DoAndReturn(func(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
		// 在 Mock GreeterService 的 SayHello 方法时，内部又调用了 greeterImpl 的 SayHello 方法。
		// 这实际上是在用真实实现替代了 Mock，失去了 Mock 的意义。
		s := &greeterImpl{}
		return s.SayHello(ctx, req)
	})
	// 用于确保多个 Mock 调用按照指定顺序执行。
	gomock.InOrder(inorderClient...)

	// Start writing unit test logic.
	type args struct {
		ctx context.Context
		req *pb.HelloRequest
		rsp *pb.HelloReply
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rsp *pb.HelloReply
			var err error
			if rsp, err = greeterService.SayHello(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("greeterImpl.SayHello() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(rsp, tt.args.rsp) {
				t.Errorf("greeterImpl.SayHello() rsp got = %v, want %v", rsp, tt.args.rsp)
			}
		})
	}
}

// 测试无法连接服务端
func TestHelloworld(t *testing.T) {
	proxy := pb.NewGreeterClientProxy()

	req := &pb.HelloRequest{
		Msg: "trpc-go-client",
	}
	rsp, err := proxy.SayHello(trpc.BackgroundContext(), req)
	assert.NotNil(t, err)
	assert.Nil(t, rsp)
}

func TestValidate(t *testing.T) {
	proxy := pb.NewGreeterClientProxy(
		client.WithTarget("ip://127.0.0.1:8001"),
		client.WithProtocol("trpc"),
	)

	req := &pb.HelloRequest{
		Msg: "",
	}
	rsp, err := proxy.SayHello(trpc.BackgroundContext(), req)
	assert.NotNil(t, err)
	// 错误类型：*errs.Error，错误信息：type:business, code:51, msg:invalid HelloRequest.Msg: value length must be at least 1 runes
	fmt.Printf("错误类型：%T，错误信息：%v\n", err, err)
	// {
	// 	"Type": 2,
	// 	"Code": 51,
	// 	"Msg": "invalid HelloRequest.Msg: value length must be at least 1 runes",
	// 	"Desc": ""
	// }
	tests.PrintStruct(err)
	assert.Nil(t, rsp)
}
