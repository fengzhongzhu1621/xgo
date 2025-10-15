package transport

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/fengzhongzhu1621/xgo/network/transport/options"
	"github.com/fengzhongzhu1621/xgo/network/transport/server_transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestST_UnixDomain(t *testing.T) {
	// Disable reuse port
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancel()
		time.Sleep(time.Millisecond * 100) // Ensure the unix listener is closed.
	})
	require.Nil(t, server_transport.NewServerTransport(
		options.WithReusePort(false),
	).ListenAndServe(
		ctx,
		options.WithListenNetwork("unix"),
		options.WithListenAddress(fmt.Sprintf("%s/test.sock", t.TempDir())),
		options.WithServerFramerBuilder(&framerBuilder{}),
	))

	// Enable reuse port
	require.Nil(t, server_transport.NewServerTransport(
		options.WithReusePort(true),
	).ListenAndServe(
		ctx,
		options.WithListenNetwork("unix"),
		options.WithListenAddress(fmt.Sprintf("%s/test.sock", t.TempDir())),
		options.WithServerFramerBuilder(&framerBuilder{}),
	))
}

func TestGetPassedListenerErr(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	assert.Nil(t, err)
	addr := listener.Addr().String()
	ln := listener.(*net.TCPListener)
	file, _ := ln.File()

	_ = os.Setenv(server_transport.EnvGraceFirstFd, fmt.Sprint(file.Fd()))
	_ = os.Setenv(server_transport.EnvGraceRestartFdNum, "1")

	_, err = server_transport.GetPassedListener("tcp", fmt.Sprintf("localhost:%d", savedListenerPort))
	assert.NotNil(t, err)

	// Simulate fd derived from environment.
	_, err = server_transport.GetPassedListener("tcp", addr)
	assert.Nil(t, err)

	_ = os.Setenv(server_transport.EnvGraceRestart, "true")
	fb := GetFramerBuilder("trpc")

	st := server_transport.NewServerTransport(options.WithReusePort(false))
	err = st.ListenAndServe(context.Background(),
		options.WithListenNetwork("tcp"),
		options.WithListenAddress(addr),
		options.WithServerFramerBuilder(fb))
	assert.NotNil(t, err)

	err = st.ListenAndServe(context.Background(),
		options.WithListenNetwork("udp"),
		options.WithListenAddress(addr),
		options.WithServerFramerBuilder(fb))
	assert.NotNil(t, err)

	_ = os.Setenv(server_transport.EnvGraceRestart, "")
}
