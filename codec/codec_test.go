package codec

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v -coverprofile=cover.out
// go tool cover -func=cover.out

// Fake is a fake codec for test
type Fake struct {
}

func (c *Fake) Encode(message IMsg, inbody []byte) (outbuf []byte, err error) {
	return nil, nil
}

func (c *Fake) Decode(message IMsg, inbuf []byte) (outbody []byte, err error) {
	return nil, nil
}

// TestCodec is unit test for the register logic of codec.
func TestCodec(t *testing.T) {
	f := &Fake{}

	Register("fake", f, f)

	serverCodec := GetServer("NoExists")
	assert.Nil(t, serverCodec)

	clientCodec := GetClient("NoExists")
	assert.Nil(t, clientCodec)

	serverCodec = GetServer("fake")
	assert.Equal(t, f, serverCodec)

	clientCodec = GetClient("fake")
	assert.Equal(t, f, clientCodec)
}

// GOMAXPROCS=1 go test -bench=WithNewMessage -benchmem -benchtime=10s
// -memprofile mem.out -cpuprofile cpu.out codec_test.go

// BenchmarkWithNewMessage is the benchmark test of codec
func BenchmarkWithNewMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WithNewMessage(context.Background())
	}
}
