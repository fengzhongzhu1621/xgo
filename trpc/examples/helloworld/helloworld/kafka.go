package main

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"trpc.group/trpc-go/trpc-database/kafka"
	"trpc.group/trpc-go/trpc-go/log"
)

func handleKafkaProducer(ctx context.Context) error {
	topic := "quickstart-events"

	// 连接kafka
	proxy := kafka.NewClientProxy("trpc.kafka.producer.service")

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

	return nil
}
