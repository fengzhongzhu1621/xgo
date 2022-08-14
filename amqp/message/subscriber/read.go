package subscriber

import (
	"time"

	"xgo/amqp/message"
)

// BulkRead reads provided amount of messages from the provided channel, until a timeout occurrs or the limit is reached.
// 批量消费消息
func BulkRead(messagesCh <-chan *message.Message, limit int, timeout time.Duration) (receivedMessages message.Messages, all bool) {
MessagesLoop:
	for len(receivedMessages) < limit {
		select {
		case msg, ok := <-messagesCh:
			if !ok {
				// 管道关闭
				break MessagesLoop
			}
			// 获取消息并回应
			receivedMessages = append(receivedMessages, msg)
			msg.Ack()
		case <-time.After(timeout):
			// 超时后停止获取消息
			break MessagesLoop
		}
	}
	// 返回获取的消息，数量可能小于limit
	return receivedMessages, len(receivedMessages) == limit
}

// BulkReadWithDeduplication reads provided number of messages from the provided channel, ignoring duplicates,
// until a timeout occurrs or the limit is reached.
// 批量消费消息，保证没有重复的消息
func BulkReadWithDeduplication(messagesCh <-chan *message.Message, limit int, timeout time.Duration) (receivedMessages message.Messages, all bool) {
	receivedIDs := map[string]struct{}{}

MessagesLoop:
	for len(receivedMessages) < limit {
		select {
		case msg, ok := <-messagesCh:
			if !ok {
				// 管道关闭
				break MessagesLoop
			}
			// 去重消息
			if _, ok := receivedIDs[msg.UUID]; !ok {
				receivedIDs[msg.UUID] = struct{}{}
				receivedMessages = append(receivedMessages, msg)
			}
			msg.Ack()
		case <-time.After(timeout):
			break MessagesLoop
		}
	}

	return receivedMessages, len(receivedMessages) == limit
}
