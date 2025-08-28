package middleware_test

import (
	"errors"
	"testing"
	"time"

	"github.com/fengzhongzhu1621/xgo/amqp/message"
	"github.com/fengzhongzhu1621/xgo/amqp/router/middleware"
	"github.com/stretchr/testify/require"
)

func TestTimeout(t *testing.T) {
	timeout := middleware.Timeout(time.Millisecond * 10)

	h := timeout(func(msg *message.Message) ([]*message.Message, error) {
		delay := time.After(time.Millisecond * 100)

		select {
		case <-msg.Context().Done():
			return nil, nil
		case <-delay:
			return nil, errors.New("timeout did not occur")
		}
	})

	_, err := h(message.NewMessage("any-uuid", nil))
	require.NoError(t, err)
}
