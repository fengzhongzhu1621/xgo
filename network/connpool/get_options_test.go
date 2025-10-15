package connpool

import (
	"context"
	"testing"
	"time"

	"github.com/fengzhongzhu1621/xgo/buildin/buffer"
	"github.com/stretchr/testify/assert"
)

func TestWithGetOptions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	fb := &noopFramerBuilder{}
	opts := &GetOptions{CustomReader: buffer.NewReader,
		FramerBuilder: fb,
		Ctx:           ctx,
	}

	localAddr := "127.0.0.1:8080"
	opts.WithLocalAddr(localAddr)
	protocol := "xxx-protocol"
	opts.WithProtocol(protocol)
	opts.WithCustomReader(buffer.NewReader)

	assert.Equal(t, opts.FramerBuilder, fb)
	assert.Equal(t, opts.Ctx, ctx)
	assert.Equal(t, opts.LocalAddr, localAddr)
	assert.Equal(t, protocol, opts.Protocol)
	assert.NotNil(t, opts.CustomReader)
}
