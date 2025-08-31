package reuseport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func moreCaseNewReusablePortPacketConn(t *testing.T) {
	listenerFour, err := NewReusablePortListener("udp6", ":10081")
	assert.Nil(t, err)
	defer listenerFour.Close()

	listenerFive, err := NewReusablePortListener("udp4", ":10081")
	assert.Nil(t, err)
	defer listenerFive.Close()

	listenerSix, err := NewReusablePortListener("udp", ":10081")
	assert.Nil(t, err)
	defer listenerSix.Close()
}

func TestNewReusablePortPacketConn(t *testing.T) {
	listenerOne, err := NewReusablePortPacketConn("udp4", "localhost:10082")
	assert.Nil(t, err)
	defer listenerOne.Close()

	listenerTwo, err := NewReusablePortPacketConn("udp", "127.0.0.1:10082")
	assert.Nil(t, err)
	defer listenerTwo.Close()

	listenerThree, err := NewReusablePortPacketConn("udp6", ":10082")
	assert.Nil(t, err)
	defer listenerThree.Close()

	moreCaseNewReusablePortPacketConn(t)
}

func BenchmarkNewReusableUDPPortListener(b *testing.B) {
	for i := 0; i < b.N; i++ {
		listener, err := NewReusablePortPacketConn("udp4", "localhost:10082")

		if err != nil {
			b.Error(err)
		} else {
			listener.Close()
		}
	}
}
