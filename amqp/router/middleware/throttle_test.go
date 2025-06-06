package middleware_test

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/fengzhongzhu1621/xgo/amqp/message"
	"github.com/fengzhongzhu1621/xgo/amqp/router/middleware"

	"github.com/stretchr/testify/assert"
)

const (
	perSecond          = 10
	testTimeout        = time.Second
	concurrentHandlers = 10
)

func TestThrottle_Middleware(t *testing.T) {
	// 每秒10个请求
	throttle := middleware.NewThrottle(perSecond, testTimeout)
	// 1秒后超时
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	producedMessagesChannel := make(chan struct{})

	producedMessagesCounter := 0

	for i := 0; i < concurrentHandlers; i++ {
		go func() {
			for {
				producedMessages := []*message.Message{message.NewMessage("produced", nil)}
				producedErr := errors.New("produced err")

				produced, err := throttle.Middleware(func(msg *message.Message) ([]*message.Message, error) {
					return producedMessages, producedErr
				})(
					message.NewMessage("uuid", nil),
				)

				assert.Equal(t, producedMessages, produced)
				assert.Equal(t, producedErr, err)

				go func() {
					// non blocking counting
					producedMessagesChannel <- struct{}{}
				}()

				select {
				case <-ctx.Done():
					break
				default:
				}
			}
		}()
	}

CounterLoop:
	for {
		select {
		case <-ctx.Done():
			break CounterLoop
		case <-producedMessagesChannel:
			producedMessagesCounter++
		}
	}

	t.Logf("produced %d messages in %d seconds, at rate of total %d messages per second",
		producedMessagesCounter,
		int(testTimeout.Seconds()),
		perSecond,
	)

	assert.True(t, producedMessagesCounter <= int(perSecond*testTimeout.Seconds()))
	assert.True(t, producedMessagesCounter > 0)
}
