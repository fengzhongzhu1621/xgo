package transport

import (
	"encoding/binary"
	"encoding/json"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/fengzhongzhu1621/xgo/network/transport/options"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ServerTransport_UDP(t *testing.T) {
	var addr = getFreeAddr("udp")

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := ListenAndServe(
			options.WithListenNetwork("udp"),
			options.WithListenAddress(addr),
			options.WithHandler(&echoHandler{}),
			options.WithServerFramerBuilder(&framerBuilder{safe: true}),
		)
		assert.Nil(t, err)
		time.Sleep(20 * time.Millisecond)
	}()
	wg.Wait() // 等待服务启动成功

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
	reqDataBad := append(reqData, []byte("remain")...)
	pc, err := net.Dial("udp", addr)
	require.Nil(t, err)

	// Bad request, server will not response.
	pc.Write(reqDataBad)
	result := make([]byte, 20)
	pc.SetDeadline(time.Now().Add(100 * time.Millisecond))
	_, err = pc.Read(result)
	require.ErrorIs(t, err, os.ErrDeadlineExceeded)
	pc.SetDeadline(time.Time{})

	// Good request, server will response.
	pc.Write(reqData)
	_, err = pc.Read(result)
	require.Nil(t, err)
}
