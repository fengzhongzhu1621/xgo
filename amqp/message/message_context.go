package message

import (
	"context"
)

type ctxKey string

// 定义上下文中的key
const (
	HandlerNameKey    ctxKey = "handler_name"
	PublisherNameKey  ctxKey = "publisher_name"
	SubscriberNameKey ctxKey = "subscriber_name"
	SubscribeTopicKey ctxKey = "subscribe_topic"
	PublishTopicKey   ctxKey = "publish_topic"
)

// valFromCtx 从上下文根据key获取value的值
func valFromCtx(ctx context.Context, key ctxKey) string {
	val, ok := ctx.Value(key).(string)
	if !ok {
		return ""
	}
	return val
}

// HandlerNameFromCtx returns the name of the message handler in the router that consumed the message.
func HandlerNameFromCtx(ctx context.Context) string {
	return valFromCtx(ctx, HandlerNameKey)
}

// PublisherNameFromCtx returns the name of the message publisher type that published the message in the router.
// For example, for Kafka it will be `kafka.Publisher`.
func PublisherNameFromCtx(ctx context.Context) string {
	return valFromCtx(ctx, PublisherNameKey)
}

// SubscriberNameFromCtx returns the name of the message subscriber type that subscribed to the message in the router.
// For example, for Kafka it will be `kafka.Subscriber`.
func SubscriberNameFromCtx(ctx context.Context) string {
	return valFromCtx(ctx, SubscriberNameKey)
}

// SubscribeTopicFromCtx returns the topic from which message was received in the router.
func SubscribeTopicFromCtx(ctx context.Context) string {
	return valFromCtx(ctx, SubscribeTopicKey)
}

// PublishTopicFromCtx returns the topic to which message will be published by the router.
func PublishTopicFromCtx(ctx context.Context) string {
	return valFromCtx(ctx, PublishTopicKey)
}
