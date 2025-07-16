package main

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/fengzhongzhu1621/xgo/tests"
	"trpc.group/trpc-go/trpc-database/kafka"
	"trpc.group/trpc-go/trpc-go/log"
)

func init() {
	// AsyncProducerSuccCallback asynchronous production success callback,
	// no processing is done by default, the user can rewrite the callback function to achieve sending success capture.
	// var AsyncProducerSuccCallback = func(topic string, key, value []byte, headers []sarama.RecordHeader) {}
	kafka.AsyncProducerSuccCallback = func(topic string, key, value []byte, headers []sarama.RecordHeader) {
		log.Infof("async producer success. topic=%s, key=%s, value=%s", topic, string(key), string(value))
	}
}

// handleKafkaProducer 模拟生产者发送消息给kafka
func handleKafkaProducer(ctx context.Context) error {
	topic := "quickstart-events"

	// 连接kafka
	proxy := kafka.NewClientProxy("trpc.kafka.producer.service")
	// proxy := kafka.NewClientProxy("trpc.kafka.server.service",
	// 	client.WithTarget("kafka://127.0.0.1:9092?clientid=test_producer&partitioner=hash"))

	// 发送消息
	key := "key"
	value := "value"
	for i := range 3 {
		err := proxy.Produce(ctx, []byte(key), []byte(value))
		if err != nil {
			log.Fatal(i, key, value, err)
			time.Sleep(time.Second)
			continue
		}
		log.Info(i, key, value)
		break
	}

	// 生产原生 sarama 消息，返回 offset、partition
	key2 := "key2"
	value2 := "value2"
	partition, offset, err := proxy.SendSaramaMessage(ctx, sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key2),
		Value: sarama.ByteEncoder(value2),
	})
	log.Infof("partition=%d, offset=%d, err=%v, key=%s, value=%s", partition, offset, err, key2, value2)

	// 生产非 saram 消息，返回 offset、partition
	key3 := "key3"
	value3 := "value3"
	partition, offset, err = proxy.SendMessage(ctx, topic, []byte(key3), []byte(value3))
	log.Infof("partition=%d, offset=%d, err=%v, key=%s, value=%s", partition, offset, err, key3, value3)

	// asynchronous producing
	// key4 := "key4"
	// value4 := "value4"
	// err = proxy.AsyncSendMessage(ctx, topic, []byte(key4), []byte(value4))
	// log.Infof("partition=%d, offset=%d, err=%v, key=%s, value=%s", partition, offset, err, key4, value4)

	return nil
}

// Consumer is the consumer
type Consumer struct{}

// Handle handle function
func (Consumer) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
	if rawContext, ok := kafka.GetRawSaramaContext(ctx); ok {
		log.Infof("InitialOffset:%d", rawContext.Claim.InitialOffset())
	}
	// log.Infof("get kafka message: %+v", msg)
	log.Infof("get kafka message: %s", tests.ToString(msg))

	// {
	// 	"Headers": [],
	// 	"Timestamp": "2025-07-16T17:07:50.016+08:00",
	// 	"BlockTimestamp": "0001-01-01T00:00:00Z",
	// 	"Key": "a2V5Mw==",
	// 	"Value": "dmFsdWUz",
	// 	"Topic": "quickstart-events",
	// 	"Partition": 0,
	// 	"Offset": 5797
	// }

	// Successful consumption is confirmed only when returning nil.
	// 当返回 nil 时，插件才会确认消费成功
	// 当返回非 nil时，插件会休眠 3s 后重新消费。
	//
	// 返回 nil 会导致重试之后可能会继续返回 nil，导致出无限重试消费，消费者会不断重复消费同一个消息，导致其它的消息被阻塞
	// 建议：由业务逻辑实现者自行重试，重试失败后放到失败队列单独处理，或者触发告警通知；业务代码自行处理错误重试，无论何种情况不要将任何错误返回给框架
	// 消费者业务逻辑尽量不要耗时太长，耗时超过阈值会触发超时异常，也会返回nil值，导致重复消费
	//
	return nil
}
