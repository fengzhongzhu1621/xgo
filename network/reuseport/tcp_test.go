package reuseport

import (
	"fmt"
	"html"
	"net/http"
	"os"
	"testing"

	"github.com/fengzhongzhu1621/xgo/tests/httptest"
	"github.com/gookit/goutil/testutil/assert"
)

const (
	httpServerOneResponse = "1"
	httpServerTwoResponse = "2"
)

var (
	httpServerOne = httptest.NewHTTPServer(httpServerOneResponse)
	httpServerTwo = httptest.NewHTTPServer(httpServerTwoResponse)
)

func TestNewReusablePortListener(t *testing.T) {
	listenerOne, err := NewReusablePortListener("tcp4", "localhost:10081")
	assert.Nil(t, err)
	defer listenerOne.Close()

	listenerTwo, err := NewReusablePortListener("tcp", "127.0.0.1:10081")
	assert.Nil(t, err)
	defer listenerTwo.Close()

	// devcloud ipv6地址无效
	_, err = NewReusablePortListener("tcp6", "[::x]:10081")
	if err == nil {
		t.Errorf("expect err, err[%v]", err)
	}

	listenerFour, err := NewReusablePortListener("tcp6", ":10081")
	assert.Nil(t, err)
	defer listenerFour.Close()

	listenerFive, err := NewReusablePortListener("tcp4", ":10081")
	assert.Nil(t, err)
	defer listenerFive.Close()

	listenerSix, err := NewReusablePortListener("tcp", ":10081")
	assert.Nil(t, err)
	defer listenerSix.Close()

	// proto invalid 非法协议
	_, err = NewReusablePortListener("xxx", "")
	if err == nil {
		t.Errorf("expect err")
	}
}

func TestListen(t *testing.T) {
	listenerOne, err := Listen("tcp4", "localhost:10081")
	assert.Nil(t, err)
	defer listenerOne.Close()

	listenerTwo, err := Listen("tcp", "127.0.0.1:10081")
	assert.Nil(t, err)
	defer listenerTwo.Close()

	listenerThree, err := Listen("tcp6", ":10081")
	assert.Nil(t, err)
	defer listenerThree.Close()

	listenerFour, err := Listen("tcp6", ":10081")
	assert.Nil(t, err)
	defer listenerFour.Close()

	listenerFive, err := Listen("tcp4", ":10081")
	assert.Nil(t, err)
	defer listenerFive.Close()

	listenerSix, err := Listen("tcp", ":10081")
	assert.Nil(t, err)
	defer listenerSix.Close()
}

func TestNewReusablePortServers(t *testing.T) {
	// 创建两个tcp server，监听同一个端口，测试是否可以同时启动
	listenerOne, err := NewReusablePortListener("tcp4", "localhost:10081")
	assert.Nil(t, err)
	defer listenerOne.Close()

	// listenerTwo, err := NewReusablePortListener("tcp6", ":10081")
	listenerTwo, err := NewReusablePortListener("tcp", "localhost:10081")
	assert.Nil(t, err)
	defer listenerTwo.Close()

	httpServerOne.Listener = listenerOne
	httpServerTwo.Listener = listenerTwo

	// 启动服务器
	httpServerOne.Start()
	httpServerTwo.Start()

	// Server One — First Response
	httptest.HttpGet(httpServerOne.URL, httpServerOneResponse, httpServerTwoResponse, t)

	// Server Two — First Response
	httptest.HttpGet(httpServerTwo.URL, httpServerOneResponse, httpServerTwoResponse, t)
	httpServerTwo.Close()

	// Server One — Second Response
	httptest.HttpGet(httpServerOne.URL, httpServerOneResponse, "", t)

	// Server One — Third Response
	httptest.HttpGet(httpServerOne.URL, httpServerOneResponse, "", t)
	httpServerOne.Close()
}

func ExampleNewReusablePortListener() {
	listener, err := NewReusablePortListener("tcp", ":8881")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	server := &http.Server{}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(os.Getgid())
		fmt.Fprintf(w, "Hello, %q\n", html.EscapeString(r.URL.Path))
	})

	panic(server.Serve(listener))
}

func BenchmarkNewReusablePortListener(b *testing.B) {
	for i := 0; i < b.N; i++ {
		listener, err := NewReusablePortListener("tcp", ":10081")

		if err != nil {
			b.Error(err)
		} else {
			listener.Close()
		}
	}
}
