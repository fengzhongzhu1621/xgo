package ip

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddrToKey(t *testing.T) {
	laddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:10000")
	require.Nil(t, err)
	raddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:10001")
	require.Nil(t, err)
	key := AddrToKey(laddr, raddr)

	assert.Equal(t, "tcp", laddr.Network())
	assert.Equal(t, "127.0.0.1:10000", laddr.String())
	assert.Equal(t, "127.0.0.1:10001", raddr.String())
	require.Equal(t, key, laddr.Network()+"_"+laddr.String()+"_"+raddr.String())
}
