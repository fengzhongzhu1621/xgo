package amqp

import (
	"context"
	"testing"
	"time"

	"github.com/fengzhongzhu1621/xgo/crypto/uuid"

	"github.com/fengzhongzhu1621/xgo/amqp/message"
	"github.com/fengzhongzhu1621/xgo/amqp/pubsub/gochannel"
	"github.com/fengzhongzhu1621/xgo/logging"
)

var logger = logging.JwwLogger{}

func TestGochannle(t *testing.T) {
	pubSub := gochannel.NewGoChannel(
		gochannel.Config{},
		logging.NewStdLogger(false, false, "[watermill] "),
	)
	// 创建并注册一个订阅者并返回接收消息的缓存队列
	messages, err := pubSub.Subscribe(context.Background(), "example.topic")
	if err != nil {
		panic(err)
	}

	// 处理订阅者接收到的消息
	go process(messages)

	publishMessages(pubSub)
}

func process(messages <-chan *message.Message) {
	for msg := range messages {
		logger.Info("received message: %s, payload: %s", msg.UUID, string(msg.Payload))

		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}

// 生产者发送消息
func publishMessages(publisher message.Publisher) {
	for i := 0; i < 10; i++ {
		msg := message.NewMessage(uuid.NewUUID(), []byte("Hello, world!"))
		// 生产者发送消息给所有注册的订阅者，并启动协程模拟订阅者接受消息到消息缓存队列，并等待消息返回ack
		if err := publisher.Publish("example.topic", msg); err != nil {
			panic(err)
		}

		time.Sleep(time.Second)
	}
}
