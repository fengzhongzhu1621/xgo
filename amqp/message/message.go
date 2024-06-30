package message

import (
	"bytes"
	"context"
	"sync"

	. "github.com/fengzhongzhu1621/xgo/collections/maps"
)

// HandlerFunc is function called when message is received.
//
// msg.Ack() is called automatically when HandlerFunc doesn't return error.
// When HandlerFunc returns error, msg.Nack() is called.
// When msg.Ack() was called in handler and HandlerFunc returns error,
// msg.Nack() will be not sent because Ack was already sent.
//
// HandlerFunc's are executed parallel when multiple messages was received
// (because msg.Ack() was sent in HandlerFunc or Subscriber supports multiple consumers).
type HandlerFunc func(msg *Message) ([]*Message, error)

var closedchan = make(chan struct{})

func init() {
	// 加载时关闭channel
	close(closedchan)
}

// Payload is the Message's payload.
type Payload []byte

// Message is the basic transfer unit.
// Messages are emitted by Publishers and received by Subscribers.
type Message struct {
	// UUID is an unique identifier of message.
	//
	// It is only used by Watermill for debugging.
	// UUID can be empty.
	// 消息唯一标识
	UUID string

	// Metadata contains the message metadata.
	//
	// Can be used to store data which doesn't require unmarshaling the entire payload.
	// It is something similar to HTTP request's headers.
	//
	// Metadata is marshaled and will be saved to the PubSub.
	// 消息元数据字典
	Metadata Metadata

	// Payload is the message's payload.
	Payload Payload

	// ack is closed, when acknowledge is received.
	ack chan struct{}
	// noACk is closed, when negative acknowledge is received.
	noAck chan struct{}

	ackMutex    sync.Mutex
	ackSentType ackType // 消息类型，3种，默认值是noAckSent

	ctx context.Context
}

// NewMessage creates a new Message with given uuid and payload.
func NewMessage(uuid string, payload Payload) *Message {
	return &Message{
		UUID:     uuid,
		Metadata: make(map[string]string),
		Payload:  payload,
		ack:      make(chan struct{}),
		noAck:    make(chan struct{}),
	}
}

type ackType int

const (
	noAckSent ackType = iota
	ack
	nack
)

// Equals compare, that two messages are equal. Acks/Nacks are not compared.
func (m *Message) Equals(toCompare *Message) bool {
	// 字符串比较
	if m.UUID != toCompare.UUID {
		return false
	}
	// 字典比较，先比较大小，然后比较value
	if len(m.Metadata) != len(toCompare.Metadata) {
		return false
	}
	for key, value := range m.Metadata {
		if value != toCompare.Metadata[key] {
			return false
		}
	}
	// 字节比较
	return bytes.Equal(m.Payload, toCompare.Payload)
}

// Ack sends message's acknowledgement.
//
// Ack is not blocking.
// Ack is idempotent.
// False is returned, if Nack is already sent.
func (m *Message) Ack() bool {
	m.ackMutex.Lock()
	defer m.ackMutex.Unlock()

	if m.ackSentType == nack {
		return false
	}
	if m.ackSentType != noAckSent {
		// 第n(n>1)次ack，然后true
		// 消息已经执行过ack了
		return true
	}
	// 第一次ack
	// ackSentType: noAckSent -> ack
	m.ackSentType = ack
	if m.ack == nil {
		m.ack = closedchan
	} else {
		// 收到ack后，关闭应答功能
		close(m.ack)
	}

	return true
}

// Nack sends message's negative acknowledgement.
//
// Nack is not blocking.
// Nack is idempotent.
// False is returned, if Ack is already sent.
func (m *Message) Nack() bool {
	m.ackMutex.Lock()
	defer m.ackMutex.Unlock()

	if m.ackSentType == ack {
		return false
	}
	if m.ackSentType != noAckSent {
		// 第一次ack，然后true
		return true
	}
	// ackSentType: noAckSent -> nack
	m.ackSentType = nack

	if m.noAck == nil {
		m.noAck = closedchan
	} else {
		close(m.noAck)
	}

	return true
}

// Acked returns channel which is closed when acknowledgement is sent.
//
// Usage:
// 		select {
//		case <-message.Acked():
// 			// ack received
//		case <-message.Nacked():
//			// nack received
//		}
func (m *Message) Acked() <-chan struct{} {
	return m.ack
}

// Nacked returns channel which is closed when negative acknowledgement is sent.
//
// Usage:
// 		select {
//		case <-message.Acked():
// 			// ack received
//		case <-message.Nacked():
//			// nack received
//		}
func (m *Message) Nacked() <-chan struct{} {
	return m.noAck
}

// Context returns the message's context. To change the context, use
// SetContext.
//
// The returned context is always non-nil; it defaults to the
// background context.
func (m *Message) Context() context.Context {
	if m.ctx != nil {
		return m.ctx
	}
	return context.Background()
}

// SetContext sets provided context to the message.
func (m *Message) SetContext(ctx context.Context) {
	m.ctx = ctx
}

// Copy copies all message without Acks/Nacks.
// The context is not propagated to the copy.
func (m *Message) Copy() *Message {
	msg := NewMessage(m.UUID, m.Payload)
	for k, v := range m.Metadata {
		msg.Metadata.Set(k, v)
	}
	return msg
}
