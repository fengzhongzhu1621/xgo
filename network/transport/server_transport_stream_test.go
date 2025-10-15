package transport

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/fengzhongzhu1621/xgo/codec"
	"github.com/fengzhongzhu1621/xgo/network/transport/client_transport"
	"github.com/fengzhongzhu1621/xgo/network/transport/options"
	"github.com/fengzhongzhu1621/xgo/network/transport/server_transport"
	"github.com/stretchr/testify/assert"
)

// TestStreamTCPListenAndServe tests listen and send.
func TestStreamTCPListenAndServe(t *testing.T) {
	st := server_transport.NewServerStreamTransport()

	// 启动一个 tcp 服务
	go func() {
		err := st.ListenAndServe(context.Background(),
			options.WithListenNetwork("tcp"),
			options.WithListenAddress(":12013"),
			options.WithHandler(&echoHandler{}),                          // 模拟业务处理逻辑，返回成功
			options.WithServerFramerBuilder(&multiplexedFramerBuilder{}), // 帧构建器
		)
		if err != nil {
			t.Logf("ListenAndServe fail:%v", err)
		}
	}()

	ctx, f := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer f()

	req := &helloRequest{
		Name: "trpc",
		Msg:  "HelloWorld",
	}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json marshal fail:%v", err)
	}
	headData := make([]byte, 8)
	binary.BigEndian.PutUint32(headData[:4], defaultStreamID)
	binary.BigEndian.PutUint32(headData[4:8], uint32(len(data)))
	reqData := append(headData, data...)

	ctx, msg := codec.WithNewMessage(ctx)
	msg.WithStreamID(defaultStreamID)

	time.Sleep(time.Millisecond * 20)

	// 创建流式客户端
	ct := client_transport.NewClientStreamTransport()
	err = ct.Init(ctx, options.WithDialNetwork("tcp"), options.WithDialAddress(":12013"),
		options.WithClientFramerBuilder(&multiplexedFramerBuilder{}),
		options.WithMsg(msg))
	assert.Nil(t, err)

	// 发送流消息
	err = ct.Send(ctx, reqData)
	assert.Nil(t, err)
	err = st.Send(ctx, reqData)
	assert.NotNil(t, err)

	// 接受流消息
	rsp, err := ct.Recv(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, rsp)

	// 关闭流
	ct.Close(ctx)

	err = ct.Send(ctx, reqData)
	assert.NotNil(t, err)

}

// TestStreamTCPListenAndServeFail tests listen and send failures.
func TestStreamTCPListenAndServeFail(t *testing.T) {
	st := server_transport.NewServerStreamTransport()
	go func() {
		err := st.ListenAndServe(context.Background(),
			options.WithListenNetwork("tcp"),
			options.WithListenAddress(":12014"),
			options.WithHandler(&echoHandler{}),
			options.WithServerFramerBuilder(&multiplexedFramerBuilder{}),
		)
		if err != nil {
			t.Logf("ListenAndServe fail:%v", err)
		}
	}()

	ctx, f := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer f()
	req := &helloRequest{
		Name: "trpc",
		Msg:  "HelloWorld",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json marshal fail:%v", err)
	}
	headData := make([]byte, 8)
	binary.BigEndian.PutUint32(headData[:4], defaultStreamID)
	binary.BigEndian.PutUint32(headData[4:8], uint32(len(data)))
	reqData := append(headData, data...)

	ctx, msg := codec.WithNewMessage(ctx)
	msg.WithStreamID(defaultStreamID)

	time.Sleep(time.Millisecond * 20)
	ct := client_transport.NewClientStreamTransport()
	err = ct.Init(ctx, options.WithDialNetwork("tcp"), options.WithDialAddress(":12015"),
		options.WithClientFramerBuilder(&multiplexedFramerBuilder{}))
	assert.NotNil(t, err)
	err = ct.Send(ctx, reqData)
	assert.NotNil(t, err)
	_, err = ct.Recv(ctx)
	assert.NotNil(t, err)
	ct.Close(ctx)

	// Test opts pool is nil.
	err = ct.Init(ctx, options.WithDialPool(nil))
	assert.NotNil(t, err)

	// Test frame builder is nil.
	err = ct.Init(ctx)
	assert.NotNil(t, err)

	// test context.
	ct = client_transport.NewClientStreamTransport()
	err = ct.Init(ctx, options.WithDialNetwork("tcp"), options.WithDialAddress(":12014"),
		options.WithClientFramerBuilder(&multiplexedFramerBuilder{}))
	assert.NotNil(t, err)

	ctx = context.Background()
	ctx, msg = codec.WithNewMessage(ctx)
	msg.WithStreamID(defaultStreamID)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	_, err = ct.Recv(ctx)
	// type:framework, code:161, msg:tcp client transport canceled before Write: context canceled
	assert.NotNil(t, err)

	ctx = context.Background()
	ctx, msg = codec.WithNewMessage(ctx)
	msg.WithStreamID(defaultStreamID)
	ctx, cancel = context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()
	_, err = ct.Recv(ctx)
	// type:framework, code:101, msg:tcp client transport timeout before Write: context deadline exceeded
	assert.NotNil(t, err)

}

// TestStreamTCPListenAndServeSend tests listen and send failures.
func TestStreamTCPListenAndServeSend(t *testing.T) {
	lnAddr := "127.0.0.1:12016"
	st := server_transport.NewServerStreamTransport()
	go func() {
		err := st.ListenAndServe(context.Background(),
			options.WithListenNetwork("tcp"),
			options.WithListenAddress(lnAddr),
			options.WithHandler(&echoStreamHandler{}),
			options.WithServerFramerBuilder(&multiplexedFramerBuilder{}),
		)
		if err != nil {
			t.Logf("ListenAndServe fail:%v", err)
		}
	}()
	time.Sleep(20 * time.Millisecond)
	req := &helloRequest{
		Name: "trpc",
		Msg:  "HelloWorld",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json marshal fail:%v", err)
	}
	headData := make([]byte, 8)
	binary.BigEndian.PutUint32(headData[:4], defaultStreamID)
	binary.BigEndian.PutUint32(headData[4:8], uint32(len(data)))
	reqData := append(headData, data...)

	ctx := context.Background()
	ctx, msg := codec.WithNewMessage(ctx)
	msg.WithStreamID(defaultStreamID)
	fb := &multiplexedFramerBuilder{}

	// Test IO EOF.
	port := getFreeAddr("tcp")
	la := "127.0.0.1" + port
	ct := client_transport.NewClientStreamTransport()
	err = ct.Init(ctx, options.WithDialNetwork("tcp"), options.WithDialAddress(lnAddr),
		options.WithClientFramerBuilder(fb), options.WithMsg(msg), options.WithLocalAddr(la))
	assert.Nil(t, err)
	time.Sleep(100 * time.Millisecond)
	raddr, err := net.ResolveTCPAddr("tcp", la)
	assert.Nil(t, err)
	laddr, err := net.ResolveTCPAddr("tcp", lnAddr)
	assert.Nil(t, err)
	msg.WithRemoteAddr(raddr)
	msg.WithLocalAddr(laddr)
	err = st.Send(ctx, reqData)
	assert.Nil(t, err)
}
