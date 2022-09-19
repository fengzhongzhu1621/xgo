package middleware_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"xgo/amqp/message"
	"xgo/amqp/router/middleware"
)

func TestRecoverer(t *testing.T) {
	h := middleware.Recoverer(func(msg *message.Message) (messages []*message.Message, e error) {
		panic("foo")
	})

	_, err := h(message.NewMessage("1", nil))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "amqp/router/middleware/recoverer.go") // stacktrace part
}
