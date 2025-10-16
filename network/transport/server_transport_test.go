package transport

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/fengzhongzhu1621/xgo/network/ip"
	"github.com/fengzhongzhu1621/xgo/network/transport/options"
	"github.com/fengzhongzhu1621/xgo/network/transport/server_transport"
	"github.com/stretchr/testify/assert"
)

func TestTCPListenAndServe(t *testing.T) {
	var addr = ip.GetFreeAddr("tcp4")

	// Wait until server transport is ready.
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		st := server_transport.NewServerTransport(options.WithKeepAlivePeriod(time.Minute))
		err := st.ListenAndServe(context.Background(),
			options.WithListenNetwork("tcp4"),
			options.WithListenAddress(addr),
			options.WithHandler(&errorHandler{}), // 模拟服务端错误，会关闭连接池中的 socket 连接
			options.WithServerFramerBuilder(&framerBuilder{}),
			options.WithServiceName("test name"),
		)

		if err != nil {
			t.Logf("ListenAndServe fail:%v", err)
		}
	}()
	wg.Wait() // 等待服务启动

	// 构造发送的协议包
	req := &helloRequest{
		Name: "trpc",
		Msg:  "HelloWorld",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json marshal fail:%v", err)
	}
	lenData := make([]byte, 4)
	binary.BigEndian.PutUint32(lenData, uint32(len(data))) // 帧长度（4 个字节）

	// +----------------+---------------------+
	// | 4字节长度头    | 实际数据内容        |
	// | (大端序uint32) | (长度由头部指定)    |
	// +----------------+---------------------+
	reqData := append(lenData, data...)

	ctx, f := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer f()

	// 发送协议包给服务端
	_, err = RoundTrip(ctx, reqData,
		options.WithDialNetwork("tcp4"),
		options.WithDialAddress(addr),
		options.WithClientFramerBuilder(&framerBuilder{}))

	// "tcp client transport ReadFrame: EOF"
	assert.NotNil(t, err)
}
