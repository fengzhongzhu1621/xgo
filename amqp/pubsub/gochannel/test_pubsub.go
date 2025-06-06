package gochannel

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/fengzhongzhu1621/xgo/crypto/uuid"

	"github.com/fengzhongzhu1621/xgo/amqp/message"
	"github.com/fengzhongzhu1621/xgo/amqp/subscriber"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var defaultTimeout = time.Second * 15

const RaceEnabled = true

func init() {
	rand.Seed(3)
}

// TestPubSub is a universal test suite. Every Pub/Sub implementation should pass it
// before it's considered production ready.
//
// Execution of the tests may be a bit different for every Pub/Sub. You can configure it by changing provided Features.
func TestPubSub(
	t *testing.T,
	features Features,
	pubSubConstructor PubSubConstructor,
	consumerGroupPubSubConstructor ConsumerGroupPubSubConstructor,
) {
	testFuncs := []struct {
		Func        func(t *testing.T, tCtx TestContext, pubSubConstructor PubSubConstructor)
		NotParallel bool
	}{
		{Func: TestPublishSubscribe},
		{Func: TestConcurrentSubscribe},
		{Func: TestConcurrentSubscribeMultipleTopics},
		{Func: TestResendOnError},
		{Func: TestNoAck},
		{Func: TestContinueAfterSubscribeClose},
		{Func: TestConcurrentClose},
		{Func: TestContinueAfterErrors},
		{Func: TestPublishSubscribeInOrder},
		{Func: TestPublisherClose},
		{Func: TestTopic},
		{Func: TestMessageCtx},
		{Func: TestSubscribeCtx},
		{Func: TestNewSubscriberReceivesOldMessages},
		{
			Func:        TestReconnect,
			NotParallel: true,
		},
	}

	for i := range testFuncs {
		testFunc := testFuncs[i]

		// 执行通用的测试方法
		runTest(
			t,
			getTestName(testFunc.Func), // 获得测试函数的名称
			func(t *testing.T, testCtx TestContext) {
				testFunc.Func(t, testCtx, pubSubConstructor)
			},
			features,
			!testFunc.NotParallel,
		)
	}

	// 测试消费组
	runTest(
		t,
		getTestName(TestConsumerGroups),
		func(t *testing.T, testCtx TestContext) {
			TestConsumerGroups(
				t,
				testCtx,
				consumerGroupPubSubConstructor,
			)
		},
		features,
		true,
	)
}

// Features are used to configure Pub/Subs implementations behaviour.
// Different features set decides also which, and how tests should be run.
type Features struct {
	// ConsumerGroups should be true, if consumer groups are supported.
	ConsumerGroups bool

	// ExactlyOnceDelivery should be true, if exactly-once delivery is supported.
	ExactlyOnceDelivery bool

	// GuaranteedOrder should be true, if order of messages is guaranteed.
	GuaranteedOrder bool

	// Some Pub/Subs guarantee the order only when one subscriber is subscribed at a time.
	GuaranteedOrderWithSingleSubscriber bool

	// Persistent should be true, if messages are persistent between multiple instancees of a Pub/Sub
	// (in practice, only GoChannel doesn't support that).
	Persistent bool

	// RestartServiceCommand is a command to test reconnects. It should restart the message broker.
	// Example: []string{"docker", "restart", "rabbitmq"}
	RestartServiceCommand []string

	// RequireSingleInstance must be true,if a PubSub requires a single instance to work properly
	// (for example: GoChannel implementation).
	RequireSingleInstance bool

	// NewSubscriberReceivesOldMessages should be set to true if messages are persisted even
	// if they are already consumed (for example, like in Kafka).
	NewSubscriberReceivesOldMessages bool
}

// RunOnlyFastTests returns true if -short flag was provided -race was not provided.
// Useful for excluding some slow tests.
func RunOnlyFastTests() bool {
	return testing.Short() && !RaceEnabled
}

// PubSubConstructor is a function that creates a Publisher and a Subscriber.
type PubSubConstructor func(t *testing.T) (message.Publisher, message.Subscriber)

// ConsumerGroupPubSubConstructor is a function that creates a Publisher and a Subscriber that use given consumer group.
type ConsumerGroupPubSubConstructor func(t *testing.T, consumerGroup string) (message.Publisher, message.Subscriber)

// SimpleMessage is deprecated: not used anywhere internally
type SimpleMessage struct {
	Num int `json:"num"`
}

// 获得测试函数的名称
func getTestName(testFunc interface{}) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(testFunc).Pointer()).Name()
	nameSliced := strings.Split(fullName, ".")

	return nameSliced[len(nameSliced)-1]
}

// TestID is a unique ID of a test.
type TestID string

// NewTestID returns a new unique TestID.
func NewTestID() TestID {
	return TestID(uuid.NewUUID4())
}

// TestContext is a collection of values that belong to a single test.
type TestContext struct {
	// Unique ID of the test
	TestID TestID

	// PubSub features
	Features Features
}

func runTest(
	t *testing.T,
	name string, // 测试函数的名称
	fn func(t *testing.T, testCtx TestContext),
	features Features,
	parallel bool,
) {
	t.Run(name, func(t *testing.T) {
		if parallel {
			// 来标志某些测试也可以被安全的并发执行（和默认的一样，假设参数名为 t）。
			// 在并行执行的情况下，只有当那些被标记为并行的测试才会被并行执行，
			// 所以只有一个测试函数时是没意义的。它应该在测试函数体中第一个被调用（在任何需要跳过的条件之后），因为它会重置测试时间
			// 在并发情况下，同时运行的测试的数量默认取决于 GOMAXPROCS。它可以通过 -parallel n 被指定（go test -parallel 4）
			t.Parallel()
		}
		testID := NewTestID()

		t.Run(string(testID), func(t *testing.T) {
			tCtx := TestContext{
				TestID:   testID,
				Features: features,
			}

			fn(t, tCtx)
		})
	})
}

const defaultStressTestTestsCount = 10

// TestPubSubStressTest runs stress tests on a chosen Pub/Sub.
func TestPubSubStressTest(
	t *testing.T,
	features Features,
	pubSubConstructor PubSubConstructor,
	consumerGroupPubSubConstructor ConsumerGroupPubSubConstructor,
) {
	stressTestsCount, _ := strconv.ParseInt(os.Getenv("STRESS_TEST_COUNT"), 10, 64)
	if stressTestsCount == 0 {
		stressTestsCount = defaultStressTestTestsCount
	}

	for i := 0; i < int(stressTestsCount); i++ {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			TestPubSub(t, features, pubSubConstructor, consumerGroupPubSubConstructor)
		})
	}
}

// TestPublishSubscribe runs basic publish and subscribe tests on a chosen Pub/Sub.
func TestPublishSubscribe(
	t *testing.T,
	tCtx TestContext,
	pubSubConstructor PubSubConstructor,
) {
	// 模拟生产者和消费者
	pub, sub := pubSubConstructor(t)

	// 动态创建topic name
	topicName := testTopicName(tCtx.TestID)

	if subscribeInitializer, ok := sub.(message.SubscribeInitializer); ok {
		require.NoError(t, subscribeInitializer.SubscribeInitialize(topicName))
	}

	var messagesToPublish []*message.Message
	messagesPayloads := map[string][]byte{}
	messagesTestMetadata := map[string]string{}

	// 创建100个消息，放到发送缓存数组
	for i := 0; i < 100; i++ {
		id := uuid.NewUUID4()
		testMetadata := uuid.NewUUID4()

		// 创建一个消息
		payload := []byte(fmt.Sprintf("%d", i))
		msg := message.NewMessage(id, payload)

		// 设置消息的元数据
		msg.Metadata.Set("test", testMetadata)
		messagesTestMetadata[id] = testMetadata

		// 放到发送缓存数组
		messagesToPublish = append(messagesToPublish, msg)
		messagesPayloads[id] = payload
	}

	// 生产者发送消息
	err := publishWithRetry(pub, topicName, messagesToPublish...)
	require.NoError(t, err, "cannot publish message")

	messages, err := sub.Subscribe(context.Background(), topicName)
	require.NoError(t, err)

	receivedMessages, all := bulkRead(tCtx, messages, len(messagesToPublish), defaultTimeout*3)
	assert.True(t, all)

	AssertAllMessagesReceived(t, messagesToPublish, receivedMessages)
	AssertMessagesPayloads(t, messagesPayloads, receivedMessages)
	AssertMessagesMetadata(t, "test", messagesTestMetadata, receivedMessages)

	closePubSub(t, pub, sub)
	assertMessagesChannelClosed(t, messages)
}

// TestConcurrentSubscribe tests subscribing to messages by multiple concurrent subscribers.
func TestConcurrentSubscribe(
	t *testing.T,
	tCtx TestContext,
	pubSubConstructor PubSubConstructor,
) {
	pub, initSub := pubSubConstructor(t)
	defer closePubSub(t, pub, initSub)

	topicName := testTopicName(tCtx.TestID)

	messagesCount := 5000
	subscribersCount := 50

	if testing.Short() {
		messagesCount = 100
		subscribersCount = 10
	}

	if subscribeInitializer, ok := initSub.(message.SubscribeInitializer); ok {
		require.NoError(t, subscribeInitializer.SubscribeInitialize(topicName))
	}

	publishedMessages := AddSimpleMessagesParallel(t, messagesCount, pub, topicName, 50)

	var sub message.Subscriber
	if tCtx.Features.RequireSingleInstance {
		sub = initSub
	} else {
		sub = createMultipliedSubscriber(t, pubSubConstructor, subscribersCount)
	}

	defer closePubSub(t, pub, sub)

	messages, err := sub.Subscribe(context.Background(), topicName)
	require.NoError(t, err)

	receivedMessages, all := bulkRead(tCtx, messages, len(publishedMessages), defaultTimeout*3)
	assert.True(t, all)

	AssertAllMessagesReceived(t, publishedMessages, receivedMessages)
}

// TestConcurrentSubscribeMultipleTopics tests subscribing to messages by concurrent subscribers on multiple topics.
func TestConcurrentSubscribeMultipleTopics(
	t *testing.T,
	tCtx TestContext,
	pubSubConstructor PubSubConstructor,
) {
	pub, sub := pubSubConstructor(t)
	defer closePubSub(t, pub, sub)

	messagesCount := 100
	topicsCount := 50

	if testing.Short() {
		messagesCount = 50
		topicsCount = 10
	}

	var messagesToPublish []*message.Message
	for i := 0; i < messagesCount; i++ {
		id := uuid.NewUUID4()

		msg := message.NewMessage(id, nil)
		messagesToPublish = append(messagesToPublish, msg)
	}

	subsWg := sync.WaitGroup{}
	subsWg.Add(topicsCount)

	receivedMessagesCh := make(chan message.Messages, topicsCount)

	for i := 0; i < topicsCount; i++ {
		topicName := testTopicName(tCtx.TestID) + fmt.Sprintf("_%d", i)

		go func() {
			defer subsWg.Done()

			if subscribeInitializer, ok := sub.(message.SubscribeInitializer); ok {
				err := subscribeInitializer.SubscribeInitialize(topicName)
				if err != nil {
					t.Fatal(err)
				}
			}

			err := publishWithRetry(pub, topicName, messagesToPublish...)
			if err != nil {
				t.Fatal(err)
			}

			messages, err := sub.Subscribe(context.Background(), topicName)
			if err != nil {
				t.Fatal(err)
			}
			topicMessages, _ := bulkRead(tCtx, messages, len(messagesToPublish), defaultTimeout)

			receivedMessagesCh <- topicMessages
		}()
	}

	subsWg.Wait()
	close(receivedMessagesCh)

	topicsReceivedMessages := 0

	for msgs := range receivedMessagesCh {
		AssertAllMessagesReceived(t, messagesToPublish, msgs)
		topicsReceivedMessages++
	}

	assert.Equal(t, topicsCount, topicsReceivedMessages)
}

// TestPublishSubscribeInOrder tests if published messages are received in a proper order.
// This test is skipped for Pub/Subs that don't support GuaranteedOrder feature.
func TestPublishSubscribeInOrder(
	t *testing.T,
	tCtx TestContext,
	pubSubConstructor PubSubConstructor,
) {
	if !tCtx.Features.GuaranteedOrder {
		t.Skipf("order is not guaranteed")
	}

	messagesCount := 1000
	if testing.Short() {
		messagesCount = 100
	}

	pub, initSub := pubSubConstructor(t)
	defer closePubSub(t, pub, initSub)

	topicName := testTopicName(tCtx.TestID)

	if subscribeInitializer, ok := initSub.(message.SubscribeInitializer); ok {
		require.NoError(t, subscribeInitializer.SubscribeInitialize(topicName))
	}

	var messagesToPublish []*message.Message
	expectedMessages := map[string][]string{}

	for i := 0; i < messagesCount; i++ {
		id := uuid.NewUUID4()
		msgType := fmt.Sprintf("%d", i%16)

		msg := message.NewMessage(id, []byte(msgType))

		messagesToPublish = append(messagesToPublish, msg)

		if _, ok := expectedMessages[msgType]; !ok {
			expectedMessages[msgType] = []string{}
		}
		expectedMessages[msgType] = append(expectedMessages[msgType], msg.UUID)
	}

	err := publishWithRetry(pub, topicName, messagesToPublish...)
	require.NoError(t, err)

	var sub message.Subscriber
	if tCtx.Features.RequireSingleInstance {
		sub = initSub
	} else {
		subscribersCount := 10
		if tCtx.Features.GuaranteedOrderWithSingleSubscriber {
			subscribersCount = 1
		}

		sub = createMultipliedSubscriber(t, pubSubConstructor, subscribersCount)
		defer require.NoError(t, sub.Close())
	}

	messages, err := sub.Subscribe(context.Background(), topicName)
	require.NoError(t, err)

	receivedMessages, all := bulkRead(tCtx, messages, len(messagesToPublish), defaultTimeout)
	require.True(t, all, "not all messages received (%d of %d)", len(receivedMessages), len(messagesToPublish))

	receivedMessagesByType := map[string][]string{}
	for _, msg := range receivedMessages {
		if _, ok := receivedMessagesByType[string(msg.Payload)]; !ok {
			receivedMessagesByType[string(msg.Payload)] = []string{}
		}
		receivedMessagesByType[string(msg.Payload)] = append(receivedMessagesByType[string(msg.Payload)], msg.UUID)
	}

	require.Equal(t, len(receivedMessagesByType), len(expectedMessages))
	require.Equal(t, len(receivedMessages), len(messagesToPublish))

	for key, ids := range expectedMessages {
		assert.Equal(t, ids, receivedMessagesByType[key])
	}
}

// TestResendOnError tests if messages are re-delivered after the subscriber fails to process them.
func TestResendOnError(
	t *testing.T,
	tCtx TestContext,
	pubSubConstructor PubSubConstructor,
) {
	pub, sub := pubSubConstructor(t)
	defer closePubSub(t, pub, sub)

	topicName := testTopicName(tCtx.TestID)

	if subscribeInitializer, ok := sub.(message.SubscribeInitializer); ok {
		require.NoError(t, subscribeInitializer.SubscribeInitialize(topicName))
	}

	messagesToSend := 100
	nacksCount := 2

	var publishedMessages message.Messages
	allMessagesSent := make(chan struct{})

	publishedMessages = PublishSimpleMessages(t, messagesToSend, pub, topicName)
	close(allMessagesSent)

	messages, err := sub.Subscribe(context.Background(), topicName)
	require.NoError(t, err)

NackLoop:
	for i := 0; i < nacksCount; i++ {
		select {
		case msg, closed := <-messages:
			if !closed {
				t.Fatal("messages channel closed before all received")
			}

			log.Println("sending err for ", msg.UUID)
			msg.Nack()
		case <-time.After(defaultTimeout):
			break NackLoop
		}
	}

	receivedMessages, _ := bulkRead(tCtx, messages, messagesToSend, defaultTimeout)

	<-allMessagesSent
	AssertAllMessagesReceived(t, publishedMessages, receivedMessages)
}

// TestNoAck tests if no new messages are received by the subscriber until the previous message is acknowledged.
// This test is skipped for Pub/Subs that don't support GuaranteedOrder feature.
func TestNoAck(
	t *testing.T,
	tCtx TestContext,
	pubSubConstructor PubSubConstructor,
) {
	if !tCtx.Features.GuaranteedOrder {
		t.Skip("guaranteed order is required for this test")
	}

	pub, sub := pubSubConstructor(t)
	defer closePubSub(t, pub, sub)

	topicName := testTopicName(tCtx.TestID)

	if subscribeInitializer, ok := sub.(message.SubscribeInitializer); ok {
		require.NoError(t, subscribeInitializer.SubscribeInitialize(topicName))
	}

	for i := 0; i < 2; i++ {
		id := uuid.NewUUID4()
		log.Printf("sending %s", id)

		msg := message.NewMessage(id, nil)

		err := publishWithRetry(pub, topicName, msg)
		require.NoError(t, err)
	}

	messages, err := sub.Subscribe(context.Background(), topicName)
	require.NoError(t, err)

	receivedMessage := make(chan struct{})
	unlockAck := make(chan struct{}, 1)
	go func() {
		msg := <-messages
		receivedMessage <- struct{}{}
		<-unlockAck
		msg.Ack()
	}()

	select {
	case <-receivedMessage:
	// ok
	case <-time.After(defaultTimeout):
		t.Fatal("timeouted")
	}

	select {
	case msg := <-messages:
		t.Fatalf("messages channel should be blocked since Ack() was not sent, received %s", msg.UUID)
	case <-time.After(time.Millisecond * 100):
		// ok
	}

	unlockAck <- struct{}{}

	select {
	case msg := <-messages:
		msg.Ack()
	case <-time.After(time.Second * 5):
		t.Fatal("messages channel should be unblocked after Ack()")
	}

	select {
	case <-messages:
		t.Fatal("msg should be not sent again")
	case <-time.After(time.Millisecond * 50):
		// ok
	}
}

// TestContinueAfterSubscribeClose checks, that we don't lose messages after closing subscriber.
func TestContinueAfterSubscribeClose(
	t *testing.T,
	tCtx TestContext,
	createPubSub PubSubConstructor,
) {
	if !tCtx.Features.Persistent {
		t.Skip("Non-Persistent is not supported yet")
	}

	if tCtx.Features.ExactlyOnceDelivery {
		t.Skip("ExactlyOnceDelivery test is not supported yet")
	}

	totalMessagesCount := 5000
	batches := 5
	if testing.Short() {
		totalMessagesCount = 50
		batches = 2
	}
	batchSize := int(totalMessagesCount / batches)
	readAttempts := batches * 4

	pub, sub := createPubSub(t)
	defer closePubSub(t, pub, sub)

	topicName := testTopicName(tCtx.TestID)
	if subscribeInitializer, ok := sub.(message.SubscribeInitializer); ok {
		require.NoError(t, subscribeInitializer.SubscribeInitialize(topicName))
	}

	publishedMessages := AddSimpleMessagesParallel(t, totalMessagesCount, pub, topicName, 50)

	receivedMessages := map[string]*message.Message{}
	for i := 0; i < readAttempts; i++ {

		pub, sub := createPubSub(t)

		messages, err := sub.Subscribe(context.Background(), topicName)
		require.NoError(t, err)

		messagesToRead := batchSize
		messagesLeft := totalMessagesCount - len(receivedMessages)

		if messagesToRead > messagesLeft {
			messagesToRead = messagesLeft
		}

		receivedMessagesBatch, _ := bulkRead(tCtx, messages, messagesToRead, defaultTimeout)
		closePubSub(t, pub, sub)

		for _, msg := range receivedMessagesBatch {
			receivedMessages[msg.UUID] = msg
		}

		closePubSub(t, pub, sub)

		if len(receivedMessages) >= totalMessagesCount {
			break
		}
	}

	// to make this test more robust - let's consume all missing messages
	// (we care here if we didn't lost any message, not if we received duplicated)
	missingMessagesCount := totalMessagesCount - len(receivedMessages)
	if missingMessagesCount > 0 && !tCtx.Features.ExactlyOnceDelivery {
		messages, err := sub.Subscribe(context.Background(), topicName)
		require.NoError(t, err)
		defer closePubSub(t, pub, sub)

		timeout := time.After(defaultTimeout)

	MessagesLoop:
		for len(receivedMessages) < totalMessagesCount {
			select {
			case msg, ok := <-messages:
				if !ok {
					break MessagesLoop
				}

				receivedMessages[msg.UUID] = msg
				msg.Ack()
			case <-timeout:
				break MessagesLoop
			}
		}
	}

	// we need to deduplicate messages, because bulkRead will deduplicate only per one batch
	uniqueReceivedMessages := message.Messages{}
	for _, msg := range receivedMessages {
		uniqueReceivedMessages = append(uniqueReceivedMessages, msg)
	}

	AssertAllMessagesReceived(t, publishedMessages, uniqueReceivedMessages)
}

// TestConcurrentClose tests if the Pub/Sub works correctly when subscribers are being closed concurrently.
func TestConcurrentClose(
	t *testing.T,
	tCtx TestContext,
	createPubSub PubSubConstructor,
) {
	if tCtx.Features.ExactlyOnceDelivery {
		t.Skip("ExactlyOnceDelivery test is not supported yet")
	}

	topicName := testTopicName(tCtx.TestID)
	createPub, createSub := createPubSub(t)
	if subscribeInitializer, ok := createSub.(message.SubscribeInitializer); ok {
		require.NoError(t, subscribeInitializer.SubscribeInitialize(topicName))
	}
	closePubSub(t, createPub, createSub)

	totalMessagesCount := 50

	closeWg := sync.WaitGroup{}
	closeWg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer closeWg.Done()

			pub, sub := createPubSub(t)
			defer closePubSub(t, pub, sub)

			_, err := sub.Subscribe(context.Background(), topicName)
			require.NoError(t, err)
		}()
	}

	closeWg.Wait()

	pub, sub := createPubSub(t)
	expectedMessages := PublishSimpleMessages(t, totalMessagesCount, pub, topicName)
	closePubSub(t, pub, sub)

	pub, sub = createPubSub(t)
	messages, err := sub.Subscribe(context.Background(), topicName)
	require.NoError(t, err)

	receivedMessages, all := bulkRead(tCtx, messages, len(expectedMessages), defaultTimeout*3)
	assert.True(t, all)

	AssertAllMessagesReceived(t, expectedMessages, receivedMessages)
	closePubSub(t, pub, sub)
}

// TestContinueAfterErrors tests if messages are processed again after an initial failure.
func TestContinueAfterErrors(
	t *testing.T,
	tCtx TestContext,
	createPubSub PubSubConstructor,
) {
	pub, sub := createPubSub(t)
	defer closePubSub(t, pub, sub)

	topicName := testTopicName(tCtx.TestID)
	if subscribeInitializer, ok := sub.(message.SubscribeInitializer); ok {
		require.NoError(t, subscribeInitializer.SubscribeInitialize(topicName))
	}

	totalMessagesCount := 50
	subscribersToNack := 3
	nacksPerSubscriber := 100

	if testing.Short() {
		subscribersToNack = 1
		nacksPerSubscriber = 5
	}

	messagesToPublish := PublishSimpleMessages(t, totalMessagesCount, pub, topicName)

	for i := 0; i < subscribersToNack; i++ {
		var errorsPub message.Publisher
		var errorsSub message.Subscriber

		if !tCtx.Features.Persistent {
			errorsPub = pub
			errorsSub = sub
		} else {
			errorsPub, errorsSub = createPubSub(t)
		}

		messages, err := errorsSub.Subscribe(context.Background(), topicName)
		require.NoError(t, err)

		for j := 0; j < nacksPerSubscriber; j++ {
			select {
			case msg := <-messages:
				msg.Nack()
			case <-time.After(defaultTimeout):
				t.Fatal("no messages left, probably seek after error doesn't work")
			}
		}

		if tCtx.Features.Persistent {
			closePubSub(t, errorsPub, errorsSub)
		}
	}

	messages, err := sub.Subscribe(context.Background(), topicName)
	require.NoError(t, err)

	// only nacks was sent, so all messages should be consumed
	receivedMessages, _ := bulkRead(tCtx, messages, totalMessagesCount, defaultTimeout*6)
	AssertAllMessagesReceived(t, messagesToPublish, receivedMessages)
}

// TestConsumerGroups tests if the consumer groups feature behaves correctly.
// This test is skipped for Pub/Sub that don't support ConsumerGroups feature.
func TestConsumerGroups(
	t *testing.T,
	tCtx TestContext,
	pubSubConstructor ConsumerGroupPubSubConstructor,
) {
	if !tCtx.Features.ConsumerGroups {
		t.Skip("consumer groups are not supported")
	}

	publisherPub, publisherSub := pubSubConstructor(t, "test_"+uuid.NewUUID4())
	defer closePubSub(t, publisherPub, publisherSub)

	topicName := testTopicName(tCtx.TestID)
	if subscribeInitializer, ok := publisherSub.(message.SubscribeInitializer); ok {
		require.NoError(t, subscribeInitializer.SubscribeInitialize(topicName))
	}
	totalMessagesCount := 50

	group1 := generateConsumerGroup(t, pubSubConstructor, topicName)
	group2 := generateConsumerGroup(t, pubSubConstructor, topicName)

	messagesToPublish := PublishSimpleMessages(t, totalMessagesCount, publisherPub, topicName)

	assertConsumerGroupReceivedMessages(t, tCtx, pubSubConstructor, group1, topicName, messagesToPublish)
	assertConsumerGroupReceivedMessages(t, tCtx, pubSubConstructor, group2, topicName, messagesToPublish)
}

// TestPublisherClose sends big amount of messages and them run close to ensure that messages are not lost during adding.
func TestPublisherClose(
	t *testing.T,
	tCtx TestContext,
	pubSubConstructor PubSubConstructor,
) {
	pub, sub := pubSubConstructor(t)
	defer closePubSub(t, pub, sub)

	topicName := testTopicName(tCtx.TestID)
	if subscribeInitializer, ok := sub.(message.SubscribeInitializer); ok {
		require.NoError(t, subscribeInitializer.SubscribeInitialize(topicName))
	}

	messagesCount := 10000
	if testing.Short() {
		messagesCount = 100
	}

	producedMessages := AddSimpleMessagesParallel(t, messagesCount, pub, topicName, 20)

	messages, err := sub.Subscribe(context.Background(), topicName)
	require.NoError(t, err)
	receivedMessages, _ := bulkRead(tCtx, messages, messagesCount, defaultTimeout*3)

	AssertAllMessagesReceived(t, producedMessages, receivedMessages)
}

// TestTopic tests if different topics work correctly in a Pub/Sub.
func TestTopic(
	t *testing.T,
	tCtx TestContext,
	pubSubConstructor PubSubConstructor,
) {
	pub, sub := pubSubConstructor(t)
	defer closePubSub(t, pub, sub)

	topic1 := testTopicName(tCtx.TestID) + "_1"
	topic2 := testTopicName(tCtx.TestID) + "_2"

	if subscribeInitializer, ok := sub.(message.SubscribeInitializer); ok {
		require.NoError(t, subscribeInitializer.SubscribeInitialize(topic1))
	}
	if subscribeInitializer, ok := sub.(message.SubscribeInitializer); ok {
		require.NoError(t, subscribeInitializer.SubscribeInitialize(topic2))
	}

	topic1Msg := message.NewMessage(uuid.NewUUID4(), nil)
	topic2Msg := message.NewMessage(uuid.NewUUID4(), nil)

	require.NoError(t, publishWithRetry(pub, topic1, topic1Msg))
	require.NoError(t, publishWithRetry(pub, topic2, topic2Msg))

	messagesTopic1, err := sub.Subscribe(context.Background(), topic1)
	require.NoError(t, err)

	messagesTopic2, err := sub.Subscribe(context.Background(), topic2)
	require.NoError(t, err)

	messagesConsumedTopic1, received := bulkRead(tCtx, messagesTopic1, 1, defaultTimeout)
	require.True(t, received, "no messages received in topic %s", topic1)

	messagesConsumedTopic2, received := bulkRead(tCtx, messagesTopic2, 1, defaultTimeout)
	require.True(t, received, "no messages received in topic %s", topic2)

	assert.Equal(t, messagesConsumedTopic1.IDs()[0], topic1Msg.UUID)
	assert.Equal(t, messagesConsumedTopic2.IDs()[0], topic2Msg.UUID)
}

// TestMessageCtx tests if the Message's Context works correctly.
func TestMessageCtx(
	t *testing.T,
	tCtx TestContext,
	pubSubConstructor PubSubConstructor,
) {
	if tCtx.Features.ExactlyOnceDelivery {
		// with ExactlyOnce delivery (at least as implemented by NATS jetstream)
		// the second message will never be received because the broker deduplicates
		// by message ID.
		t.Skip("ExactlyOnceDelivery test is not supported yet")
	}

	pub, sub := pubSubConstructor(t)
	defer closePubSub(t, pub, sub)

	topicName := testTopicName(tCtx.TestID)
	if subscribeInitializer, ok := sub.(message.SubscribeInitializer); ok {
		require.NoError(t, subscribeInitializer.SubscribeInitialize(topicName))
	}

	msg := message.NewMessage(uuid.NewUUID4(), nil)

	// ensuring that context is not propagated via pub/sub
	ctx, ctxCancel := context.WithCancel(context.Background())
	ctxCancel()
	msg.SetContext(ctx)

	require.NoError(t, publishWithRetry(pub, topicName, msg))
	// this might actually be an error in some pubsubs (http), because we close the subscriber without ACK.
	_ = pub.Publish(topicName, msg)

	messages, err := sub.Subscribe(context.Background(), topicName)
	require.NoError(t, err)

	select {
	case msg := <-messages:
		ctx := msg.Context()

		select {
		case <-ctx.Done():
			t.Fatal("context should not be canceled")
		default:
			// ok
		}

		require.True(t, msg.Ack())

		select {
		case <-ctx.Done():
			// ok
		case <-time.After(defaultTimeout):
			t.Fatal("context should be canceled after Ack")
		}
	case <-time.After(defaultTimeout):
		t.Fatal("no message received")
	}

	select {
	case msg := <-messages:
		ctx := msg.Context()

		select {
		case <-ctx.Done():
			t.Fatal("context should not be canceled")
		default:
			// ok
		}

		go closePubSub(t, pub, sub)

		select {
		case <-ctx.Done():
			// ok
		case <-time.After(defaultTimeout):
			t.Fatal("context should be canceled after pubSub.Close()")
		}
	case <-time.After(defaultTimeout):
		t.Fatal("no message received")
	}
}

// TestSubscribeCtx tests if the Subscriber's Context works correctly.
func TestSubscribeCtx(
	t *testing.T,
	tCtx TestContext,
	pubSubConstructor PubSubConstructor,
) {
	pub, sub := pubSubConstructor(t)
	defer closePubSub(t, pub, sub)

	const messagesCount = 20

	ctxWithCancel, cancel := context.WithCancel(context.Background())
	ctxWithCancel = context.WithValue(ctxWithCancel, "foo", "bar")

	topicName := testTopicName(tCtx.TestID)
	if subscribeInitializer, ok := sub.(message.SubscribeInitializer); ok {
		require.NoError(t, subscribeInitializer.SubscribeInitialize(topicName))
	}
	publishedMessages := PublishSimpleMessages(t, messagesCount, pub, topicName)

	msgsToCancel, err := sub.Subscribe(ctxWithCancel, topicName)
	require.NoError(t, err)
	cancel()

	timeout := time.After(defaultTimeout)

ClosedLoop:
	for {
		select {
		case msg, open := <-msgsToCancel:
			if !open {
				break ClosedLoop
			}
			msg.Nack()
		case <-timeout:
			t.Fatal("messages channel is not closed after ", defaultTimeout)
			t.FailNow()
		}
		time.Sleep(time.Millisecond * 100)
	}

	ctx := context.WithValue(context.Background(), "foo", "bar")
	msgs, err := sub.Subscribe(ctx, topicName)
	require.NoError(t, err)

	receivedMessages, _ := bulkRead(tCtx, msgs, messagesCount, defaultTimeout)
	AssertAllMessagesReceived(t, publishedMessages, receivedMessages)

	for _, msg := range receivedMessages {
		assert.EqualValues(t, "bar", msg.Context().Value("foo"))
	}
}

// TestReconnect tests if reconnecting to a Pub/Sub works correctly.
func TestReconnect(
	t *testing.T,
	tCtx TestContext,
	pubSubConstructor PubSubConstructor,
) {
	if len(tCtx.Features.RestartServiceCommand) == 0 {
		t.Skip("no RestartServiceCommand provided, cannot test reconnect")
	}

	pub, sub := pubSubConstructor(t)

	topicName := testTopicName(tCtx.TestID)
	if subscribeInitializer, ok := sub.(message.SubscribeInitializer); ok {
		require.NoError(t, subscribeInitializer.SubscribeInitialize(topicName))
	}

	const messagesCount = 10000
	const publishersCount = 100

	restartAfterMessages := map[int]struct{}{
		messagesCount / 3: {}, // restart at 1/3 of messages
		messagesCount / 2: {}, // restart at 1/2 of messages
	}

	messages, err := sub.Subscribe(context.Background(), topicName)
	require.NoError(t, err)

	var publishedMessages message.Messages
	messagePublished := make(chan *message.Message, messagesCount)
	publishMessage := make(chan struct{})

	go func() {
		for i := 0; i < messagesCount; i++ {
			publishMessage <- struct{}{}

			if _, shouldRestart := restartAfterMessages[i]; shouldRestart {
				go restartServer(t, tCtx.Features)
			}
		}
		close(publishMessage)
	}()

	go func() {
		for msg := range messagePublished {
			publishedMessages = append(publishedMessages, msg)
		}
	}()

	for i := 0; i < publishersCount; i++ {
		go func() {
			for range publishMessage {
				id := uuid.NewUUID4()
				msg := message.NewMessage(id, nil)

				for {
					fmt.Println("publishing message")

					// some randomization in sending
					if rand.Int31n(10) == 0 {
						time.Sleep(time.Millisecond * 500)
					}

					if err := publishWithRetry(pub, topicName, msg); err == nil {
						break
					}

					fmt.Printf("cannot publish message %s, trying again, err: %s\n", msg.UUID, err)
					time.Sleep(time.Millisecond * 500)
				}

				messagePublished <- msg
			}
		}()
	}

	receivedMessages, allMessages := bulkRead(tCtx, messages, messagesCount, defaultTimeout*4)
	assert.True(t, allMessages, "not all messages received (has %d of %d)", len(receivedMessages), messagesCount)

	AssertAllMessagesReceived(t, publishedMessages, receivedMessages)

	closePubSub(t, pub, sub)
}

// TestNewSubscriberReceivesOldMessages tests if a new subscriber receives previously published messages.
func TestNewSubscriberReceivesOldMessages(
	t *testing.T,
	tCtx TestContext,
	pubSubConstructor PubSubConstructor,
) {
	if !tCtx.Features.NewSubscriberReceivesOldMessages {
		t.Skip("only subscribers with TestNewSubscriberReceivesOldMessages are supported")
	}

	publishedMessages := message.Messages{}

	pub, sub := pubSubConstructor(t)

	topicName := testTopicName(tCtx.TestID)
	if subscribeInitializer, ok := sub.(message.SubscribeInitializer); ok {
		require.NoError(t, subscribeInitializer.SubscribeInitialize(topicName))
	}
	require.NoError(t, sub.Close())

	var publishMessage = func() {
		publishedMessages = append(publishedMessages, PublishSimpleMessages(t, 1, pub, topicName)...)
	}
	publishMessage()

	type Subscriber struct {
		Msgs             <-chan *message.Message
		Subscriber       message.Subscriber
		ConsumedMessages int
	}

	var subscribers []*Subscriber
	defer func() {
		for _, sub := range subscribers {
			require.NoError(t, sub.Subscriber.Close())
		}
	}()

	var addSubscriber = func() {
		pub, sub := pubSubConstructor(t)
		require.NoError(t, pub.Close())

		msgs, err := sub.Subscribe(context.Background(), topicName)
		require.NoError(t, err)

		subscribers = append(subscribers, &Subscriber{
			Msgs:             msgs,
			Subscriber:       sub,
			ConsumedMessages: 0,
		})
	}

	var consumeMessages = func() {
		for i, sub := range subscribers {
			toConsume := len(publishedMessages) - sub.ConsumedMessages
			receivedMessages, all := bulkRead(tCtx, sub.Msgs, toConsume, defaultTimeout)

			require.True(t, all, "subscriber %d not received all messages (%d/%d)", i, len(receivedMessages), toConsume)

			fmt.Printf("subscriber no %d consumed %d messages\n", i, toConsume)
			sub.ConsumedMessages += toConsume
		}
	}

	publishMessage()
	addSubscriber()
	consumeMessages()

	publishMessage()
	addSubscriber()
	consumeMessages()

	publishMessage()
	addSubscriber()
	consumeMessages()
}

func restartServer(t *testing.T, features Features) {
	fmt.Println("restarting server with:", features.RestartServiceCommand)
	cmd := exec.Command(features.RestartServiceCommand[0], features.RestartServiceCommand[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

	fmt.Println("server restarted")
}

func assertConsumerGroupReceivedMessages(
	t *testing.T,
	tCtx TestContext,
	pubSubConstructor ConsumerGroupPubSubConstructor,
	consumerGroup string,
	topicName string,
	expectedMessages []*message.Message,
) {
	pub, sub := pubSubConstructor(t, consumerGroup)
	defer closePubSub(t, pub, sub)

	messages, err := sub.Subscribe(context.Background(), topicName)
	require.NoError(t, err)

	receivedMessages, all := bulkRead(tCtx, messages, len(expectedMessages), defaultTimeout)
	assert.True(t, all)

	AssertAllMessagesReceived(t, expectedMessages, receivedMessages)
}

func testTopicName(testID TestID) string {
	return "topic_" + string(testID)
}

func closePubSub(t *testing.T, pub message.Publisher, sub message.Subscriber) {
	err := pub.Close()
	require.NoError(t, err)

	err = sub.Close()
	require.NoError(t, err)
}

func generateConsumerGroup(t *testing.T, pubSubConstructor ConsumerGroupPubSubConstructor, topicName string) string {
	groupName := "cg_" + uuid.NewUUID4()

	// create a pubsub to ensure that the consumer group exists
	// for those providers that require subscription before publishing messages (e.g. Google Cloud PubSub)
	pub, sub := pubSubConstructor(t, groupName)
	if subInitializer, ok := sub.(message.SubscribeInitializer); ok {
		require.NoError(t, subInitializer.SubscribeInitialize(topicName))
	}
	_, err := sub.Subscribe(context.Background(), topicName)
	require.NoError(t, err)
	closePubSub(t, pub, sub)

	return groupName
}

// PublishSimpleMessages publishes provided number of simple messages without a payload.
func PublishSimpleMessages(t *testing.T, messagesCount int, publisher message.Publisher, topicName string) message.Messages {
	var messagesToPublish []*message.Message

	for i := 0; i < messagesCount; i++ {
		id := uuid.NewUUID4()

		msg := message.NewMessage(id, nil)
		messagesToPublish = append(messagesToPublish, msg)

		err := publishWithRetry(publisher, topicName, msg)
		require.NoError(t, err, "cannot publish messages")
	}

	return messagesToPublish
}

// AddSimpleMessagesParallel publishes provided number of simple messages without a payload
// using the provided number of publishers (goroutines).
func AddSimpleMessagesParallel(t *testing.T, messagesCount int, publisher message.Publisher, topicName string, publishers int) message.Messages {
	var messagesToPublish []*message.Message
	publishMsg := make(chan *message.Message)

	wg := sync.WaitGroup{}
	wg.Add(messagesCount)

	for i := 0; i < publishers; i++ {
		go func() {
			for msg := range publishMsg {
				err := publishWithRetry(publisher, topicName, msg)
				require.NoError(t, err, "cannot publish messages")
				wg.Done()
			}
		}()
	}

	for i := 0; i < messagesCount; i++ {
		id := uuid.NewUUID4()

		msg := message.NewMessage(id, nil)
		messagesToPublish = append(messagesToPublish, msg)

		publishMsg <- msg
	}
	close(publishMsg)

	wg.Wait()

	return messagesToPublish
}

func assertMessagesChannelClosed(t *testing.T, messages <-chan *message.Message) bool {
	select {
	case _, open := <-messages:
		return assert.False(t, open)
	default:
		t.Error("messages channel is not closed (blocked)")
		return false
	}
}

func publishWithRetry(publisher message.Publisher, topic string, messages ...*message.Message) error {
	retries := 5

	for {
		err := publisher.Publish(topic, messages...)
		if err == nil {
			return nil
		}
		// 失败重发
		retries--

		fmt.Printf("error on publish: %s, %d retries left\n", err, retries)

		if retries == 0 {
			return err
		}
	}
}

func bulkRead(testCtx TestContext, messagesCh <-chan *message.Message, limit int, timeout time.Duration) (receivedMessages message.Messages, all bool) {
	start := time.Now()

	defer func() {
		duration := time.Since(start)

		logMsg := "all messages (%d/%d) received in bulk read after %s of %s (test ID: %s)\n"
		if !all {
			logMsg = "not " + logMsg
		}

		log.Printf(logMsg, len(receivedMessages), limit, duration, timeout, testCtx.TestID)
	}()

	if !testCtx.Features.ExactlyOnceDelivery {
		return subscriber.BulkReadWithDeduplication(messagesCh, limit, timeout)
	}

	return subscriber.BulkRead(messagesCh, limit, timeout)
}

func createMultipliedSubscriber(t *testing.T, pubSubConstructor PubSubConstructor, subscribersCount int) message.Subscriber {
	return NewMultiplier(
		func() (message.Subscriber, error) {
			pub, sub := pubSubConstructor(t)
			require.NoError(t, pub.Close()) // pub is not needed

			return sub, nil
		},
		subscribersCount,
	)
}
