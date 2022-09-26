package gochannel

import (
	"context"
	"sync"
	"xgo/amqp/message"
	"xgo/log"
)

// 订阅者
type GoChannelSubscriber struct {
	ctx context.Context

	uuid string

	sending       sync.Mutex
	outputChannel chan *message.Message

	logger  log.LoggerAdapter
	closed  bool
	closing chan struct{}
}

// Close 关闭订阅者
func (s *GoChannelSubscriber) Close() {
	if s.closed {
		return
	}
	// 发送关闭信号
	close(s.closing)

	s.logger.Debug("Closing subscriber, waiting for sending lock", nil)

	// ensuring that we are not sending to closed channel
	s.sending.Lock()
	defer s.sending.Unlock()

	// 设置关闭状态
	s.logger.Debug("GoChannel Pub/Sub Subscriber closed", nil)
	s.closed = true

	// 关闭订阅者缓存队列
	close(s.outputChannel)
}

// sendMessageToSubscriber 订阅者接收消息到缓存队列
func (s *GoChannelSubscriber) SendMessageToSubscriber(msg *message.Message, logFields log.LogFields) {
	s.sending.Lock()
	defer s.sending.Unlock()

	ctx, cancelCtx := context.WithCancel(s.ctx)
	defer cancelCtx()

SendToSubscriber:
	for {
		// copy the message to prevent ack/nack propagation to other consumers
		// also allows to make retries on a fresh copy of the original message
		msgToSend := msg.Copy()
		msgToSend.SetContext(ctx)

		s.logger.Trace("Sending msg to subscriber", logFields)

		if s.closed {
			s.logger.Info("Pub/Sub closed, discarding msg", logFields)
			return
		}

		select {
		case s.outputChannel <- msgToSend:
			// 发给订阅者缓存队列
			s.logger.Trace("Sent message to subscriber", logFields)
		case <-s.closing:
			// 接收到关闭信号
			s.logger.Trace("Closing, message discarded", logFields)
			return
		}

		// 等待监听outputChannel的协程返回ack

		select {
		case <-msgToSend.Acked():
			// 消息处理完毕
			s.logger.Trace("Message acked", logFields)
			return
		case <-msgToSend.Nacked():
			// 消息取消处理，重新发给订阅者
			s.logger.Trace("Nack received, resending message", logFields)
			continue SendToSubscriber
		case <-s.closing:
			// 关闭订阅者
			s.logger.Trace("Closing, message discarded", logFields)
			return
		}
	}
}
