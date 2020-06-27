package name

import (
	"github.com/stretchr/testify/assert"
	"github.com/wxnacy/wgo/arrays"
	"testing"
)


func TestNameServiceAddressHttp(t *testing.T) {
	var address NameServiceAddress
	err := address.ParseAddress("http://www.example.com")
	assert.NoError(t, err)
	addr := address.GetAddress()
	assert.Equal(t, addr, "http://www.example.com")
}

func TestNameServiceAddressIp(t *testing.T) {
	var address NameServiceAddress
	err := address.ParseAddress("ip://1.1.1.1,2.2.2.2")
	assert.NoError(t, err)
	ip := address.GetAddress()
	index := arrays.ContainsString(address.IpList, ip[5:])
	assert.True(t, index >= 0)
}

func TestNameServiceAddressL5(t *testing.T) {
	var address NameServiceAddress
	err := address.ParseAddress("l5://111:222")
	assert.NoError(t, err)
	addr := address.GetAddress()
	assert.Equal(t, addr, "l5://111:222")
	assert.Equal(t, address.ModId, 111)
	assert.Equal(t, address.CmdId, 222)
}
