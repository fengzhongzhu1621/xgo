package gochannel

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/fengzhongzhu1621/xgo/crypto/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fengzhongzhu1621/xgo/amqp/message"
	"github.com/fengzhongzhu1621/xgo/amqp/subscriber"
	"github.com/fengzhongzhu1621/xgo/logging"
)

// 默认生产者和消费者
func createPersistentPubSub(t *testing.T) (message.Publisher, message.Subscriber) {
	pubSub := NewGoChannel(
		Config{
			OutputChannelBuffer: 10000,
			Persistent:          true,
		},
		logging.NewStdLogger(true, true, "[watermill] "),
	)
	return pubSub, pubSub
}

func TestPublishSubscribe_persistent(t *testing.T) {
	TestPubSub(
		t,
		Features{
			ConsumerGroups:        false,
			ExactlyOnceDelivery:   true,
			GuaranteedOrder:       false,
			Persistent:            false,
			RequireSingleInstance: true,
		},
		createPersistentPubSub,
		nil,
	)
}

func TestPublishSubscribe_not_persistent(t *testing.T) {
	messagesCount := 100
	pubSub := NewGoChannel(
		Config{OutputChannelBuffer: int64(messagesCount)},
		logging.NewStdLogger(true, true, "[watermill] "),
	)
	topicName := "test_topic_" + uuid.NewUUID4()

	msgs, err := pubSub.Subscribe(context.Background(), topicName)
	require.NoError(t, err)

	sendMessages := PublishSimpleMessages(t, messagesCount, pubSub, topicName)
	receivedMsgs, _ := subscriber.BulkRead(msgs, messagesCount, time.Second)

	AssertAllMessagesReceived(t, sendMessages, receivedMsgs)

	assert.NoError(t, pubSub.Close())
}

func TestPublishSubscribe_block_until_ack(t *testing.T) {
	pubSub := NewGoChannel(
		Config{BlockPublishUntilSubscriberAck: true},
		logging.NewStdLogger(true, true, "[watermill] "),
	)
	topicName := "test_topic_" + uuid.NewUUID4()

	msgs, err := pubSub.Subscribe(context.Background(), topicName)
	require.NoError(t, err)

	published := make(chan struct{})
	go func() {
		err := pubSub.Publish(topicName, message.NewMessage("1", nil))
		require.NoError(t, err)
		close(published)
	}()

	msg1 := <-msgs
	select {
	case <-published:
		t.Fatal("publish should be blocked until ack")
	default:
		// ok
	}

	msg1.Nack()
	select {
	case <-published:
		t.Fatal("publish should be blocked after nack")
	default:
		// ok
	}

	msg2 := <-msgs
	msg2.Ack()

	select {
	case <-published:
		// ok
	case <-time.After(time.Second):
		t.Fatal("publish should be not blocked after ack")
	}
}

func TestPublishSubscribe_race_condition_on_subscribe(t *testing.T) {
	testsCount := 15
	if testing.Short() {
		testsCount = 3
	}

	for i := 0; i < testsCount; i++ {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			testPublishSubscribeSubRace(t)
		})
	}
}

func TestSubscribe_race_condition_when_closing(t *testing.T) {
	testsCount := 15
	if testing.Short() {
		testsCount = 3
	}

	for i := 0; i < testsCount; i++ {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			pubSub := NewGoChannel(
				Config{},
				logging.NewStdLogger(true, false, "[watermill] "),
			)
			go func() {
				err := pubSub.Close()
				require.NoError(t, err)
			}()
			_, err := pubSub.Subscribe(context.Background(), "topic")
			require.NoError(t, err)
		})
	}
}

func TestPublish_race_condition_when_closing(t *testing.T) {
	testsCount := 15
	if testing.Short() {
		testsCount = 3
	}

	for i := 0; i < testsCount; i++ {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			pubSub := NewGoChannel(
				Config{},
				logging.NewStdLogger(true, false, "[watermill] "),
			)
			go func() {
				_ = pubSub.Publish("topic", message.NewMessage(strconv.Itoa(i), nil))
			}()

			err := pubSub.Close()
			require.NoError(t, err)
		})
	}
}

func TestPublishSubscribe_do_not_block_other_subscribers(t *testing.T) {
	pubSub := NewGoChannel(
		Config{},
		logging.NewStdLogger(true, true, "[watermill] "),
	)
	topicName := "test_topic_" + uuid.NewUUID4()

	msgsFromSubscriber1, err := pubSub.Subscribe(context.Background(), topicName)
	require.NoError(t, err)

	_, err = pubSub.Subscribe(context.Background(), topicName)
	require.NoError(t, err)

	msgsFromSubscriber3, err := pubSub.Subscribe(context.Background(), topicName)
	require.NoError(t, err)

	err = pubSub.Publish(topicName, message.NewMessage("1", nil))
	require.NoError(t, err)

	received := make(chan struct{})
	go func() {
		msg := <-msgsFromSubscriber1
		msg.Ack()

		msg = <-msgsFromSubscriber3
		msg.Ack()

		close(received)
	}()

	select {
	case <-received:
		// ok
	case <-time.After(5 * time.Second):
		t.Fatal("subscriber which didn't ack a message blocked other subscribers from receiving it")
	}
}

func testPublishSubscribeSubRace(t *testing.T) {
	t.Helper()

	messagesCount := 500
	subscribersCount := 200
	if testing.Short() {
		messagesCount = 200
		subscribersCount = 20
	}

	pubSub := NewGoChannel(
		Config{
			OutputChannelBuffer: int64(messagesCount),
			Persistent:          true,
		},
		logging.NewStdLogger(true, false, "[watermill] "),
	)

	allSent := sync.WaitGroup{}
	allSent.Add(messagesCount)
	allReceived := sync.WaitGroup{}

	sentMessages := message.Messages{}
	go func() {
		for i := 0; i < messagesCount; i++ {
			msg := message.NewMessage(uuid.NewUUID4(), nil)
			sentMessages = append(sentMessages, msg)

			go func() {
				require.NoError(t, pubSub.Publish("topic", msg))
				allSent.Done()
			}()
		}
	}()

	subscriberReceivedCh := make(chan message.Messages, subscribersCount)
	for i := 0; i < subscribersCount; i++ {
		allReceived.Add(1)

		go func() {
			msgs, err := pubSub.Subscribe(context.Background(), "topic")
			require.NoError(t, err)

			received, _ := subscriber.BulkRead(msgs, messagesCount, time.Second*10)
			subscriberReceivedCh <- received

			allReceived.Done()
		}()
	}

	logging.Info("waiting for all sent")
	allSent.Wait()

	logging.Info("waiting for all received")
	allReceived.Wait()

	close(subscriberReceivedCh)

	logging.Info("asserting")

	for subMsgs := range subscriberReceivedCh {
		AssertAllMessagesReceived(t, sentMessages, subMsgs)
	}
}
