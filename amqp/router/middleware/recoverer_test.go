package middleware_test

import (
	"testing"

	"github.com/fengzhongzhu1621/xgo/amqp/message"
	"github.com/fengzhongzhu1621/xgo/amqp/router/middleware"
	"github.com/stretchr/testify/assert"
)

func TestRecoverer(t *testing.T) {
	h := middleware.Recoverer(func(msg *message.Message) (messages []*message.Message, e error) {
		panic("foo")
	})

	_, err := h(message.NewMessage("1", nil))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "amqp/router/middleware/recoverer.go") // stacktrace part
}
