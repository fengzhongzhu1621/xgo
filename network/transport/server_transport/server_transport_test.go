package server_transport

import (
	"testing"
	"time"

	"github.com/fengzhongzhu1621/xgo/network/transport/options"
	"github.com/stretchr/testify/assert"
)

func TestNewServerTransport(t *testing.T) {
	st := NewServerTransport(options.WithKeepAlivePeriod(time.Minute))
	assert.NotNil(t, st)
}
