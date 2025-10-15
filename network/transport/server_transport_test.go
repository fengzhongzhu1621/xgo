package transport

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/fengzhongzhu1621/xgo/network/transport/options"
	"github.com/fengzhongzhu1621/xgo/network/transport/server_transport"
	"github.com/stretchr/testify/assert"
)

func TestTCPListenAndServe(t *testing.T) {
	var addr = getFreeAddr("tcp4")

	// Wait until server transport is ready.
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		st := server_transport.NewServerTransport(options.WithKeepAlivePeriod(time.Minute))
		err := st.ListenAndServe(context.Background(),
			options.WithListenNetwork("tcp4"),
			options.WithListenAddress(addr),
			options.WithHandler(&errorHandler{}),
			options.WithServerFramerBuilder(&framerBuilder{}),
			options.WithServiceName("test name"),
		)

		if err != nil {
			t.Logf("ListenAndServe fail:%v", err)
		}
	}()
	wg.Wait() // 等待服务启动

	// Round trip.
	req := &helloRequest{
		Name: "trpc",
		Msg:  "HelloWorld",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json marshal fail:%v", err)
	}
	lenData := make([]byte, 4)
	binary.BigEndian.PutUint32(lenData, uint32(len(data)))

	reqData := append(lenData, data...)

	ctx, f := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer f()

	_, err = RoundTrip(ctx, reqData,
		options.WithDialNetwork("tcp4"),
		options.WithDialAddress(addr),
		options.WithClientFramerBuilder(&framerBuilder{}))
	assert.NotNil(t, err)
}
