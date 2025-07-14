package main

import (
	"context"
	"reflect"
	"testing"

	pb "github.com/fengzhongzhu1621/xgo/trpc/trpcprotocol/helloworld"
	"go.uber.org/mock/gomock"
	_ "trpc.group/trpc-go/trpc-go/http"
)

//go:generate go mod tidy
//go:generate mockgen -destination=stub/github.com/fengzhongzhu1621/xgo/trpc/trpcprotocol/helloworld/helloworld_mock.go -package=helloworld -self_package=github.com/fengzhongzhu1621/xgo/trpc/trpcprotocol/helloworld --source=stub/github.com/fengzhongzhu1621/xgo/trpc/trpcprotocol/helloworld/helloworld.trpc.go

func Test_greeterHttpImpl_SayHello(t *testing.T) {
	// Start writing mock logic.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	greeterHttpService := pb.NewMockGreeterHttpService(ctrl)
	var inorderClient []any
	// Expected behavior.
	m := greeterHttpService.EXPECT().SayHello(gomock.Any(), gomock.Any()).AnyTimes()
	m.DoAndReturn(func(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
		s := &greeterHttpImpl{}
		return s.SayHello(ctx, req)
	})
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
			if rsp, err = greeterHttpService.SayHello(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("greeterHttpImpl.SayHello() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(rsp, tt.args.rsp) {
				t.Errorf("greeterHttpImpl.SayHello() rsp got = %v, want %v", rsp, tt.args.rsp)
			}
		})
	}
}
