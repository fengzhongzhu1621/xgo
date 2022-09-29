package message

import (
	"testing"

	. "github.com/fengzhongzhu1621/xgo/utils/maps"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessage_Equals(t *testing.T) {
	withMetadata := func(msg *Message, metadata Metadata) *Message {
		msg.Metadata = metadata
		return msg
	}

	testCases := []struct {
		Name   string
		Msg1   *Message
		Msg2   *Message
		Equals bool
	}{
		{
			Name:   "equal",
			Msg1:   NewMessage("1", []byte("foo")),
			Msg2:   NewMessage("1", []byte("foo")),
			Equals: true,
		},
		{
			Name:   "different_uuid",
			Msg1:   NewMessage("1", []byte("foo")),
			Msg2:   NewMessage("2", []byte("foo")),
			Equals: false,
		},
		{
			Name:   "different_payload",
			Msg1:   NewMessage("1", []byte("foo")),
			Msg2:   NewMessage("1", []byte("bar")),
			Equals: false,
		},
		{
			Name:   "different_metadata",
			Msg1:   withMetadata(NewMessage("1", []byte("foo")), map[string]string{"foo": "1"}),
			Msg2:   withMetadata(NewMessage("1", []byte("foo")), map[string]string{"foo": "2"}),
			Equals: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Equals, tc.Msg1.Equals(tc.Msg2))
			assert.Equal(t, tc.Equals, tc.Msg2.Equals(tc.Msg1))
		})
	}
}

func TestMessage_Ack(t *testing.T) {
	msg := &Message{}
	require.True(t, msg.Ack())

	assertAcked(t, msg)
	assertNoNack(t, msg)
}

func TestMessage_Ack_idempotent(t *testing.T) {
	msg := &Message{}
	require.True(t, msg.Ack())
	require.True(t, msg.Ack())

	assertAcked(t, msg)
}

func TestMessage_Ack_already_Nack(t *testing.T) {
	msg := &Message{}
	require.True(t, msg.Nack())

	assert.False(t, msg.Ack())
}

func TestMessage_Nack(t *testing.T) {
	msg := &Message{}
	require.True(t, msg.Nack())

	assertNoAck(t, msg)
	assertNacked(t, msg)
}

func TestMessage_Nack_idempotent(t *testing.T) {
	msg := &Message{}
	// 两次ack都返回true
	require.True(t, msg.Nack())
	require.True(t, msg.Nack())

	assertNacked(t, msg)
}

func TestMessage_Nack_already_Ack(t *testing.T) {
	msg := &Message{}
	require.True(t, msg.Ack())

	assert.False(t, msg.Nack())
}

func TestMessage_Copy(t *testing.T) {
	msg := NewMessage("1", []byte("foo"))
	msgCopy := msg.Copy()

	require.True(t, msg.Ack())

	assertAcked(t, msg)
	assertNoAck(t, msgCopy)
	assert.True(t, msg.Equals(msgCopy))
}

func TestMessage_CopyMetadata(t *testing.T) {
	msg := NewMessage("1", []byte("foo"))
	msg.Metadata.Set("foo", "bar")
	msgCopy := msg.Copy()

	msg.Metadata.Set("foo", "baz")

	assert.Equal(t, msgCopy.Metadata.Get("foo"), "bar", "did not expect changing source message's metadata to alter copy's metadata")
}

func assertAcked(t *testing.T, msg *Message) {
	select {
	case <-msg.Acked():
		// ok
	default:
		t.Fatal("no ack received")
	}
}

func assertNacked(t *testing.T, msg *Message) {
	select {
	case <-msg.Nacked():
		// ok
	default:
		t.Fatal("no ack received")
	}
}

func assertNoAck(t *testing.T, msg *Message) {
	select {
	case <-msg.Acked():
		t.Fatal("nack should be not sent")
	default:
		// ok
	}
}

func assertNoNack(t *testing.T, msg *Message) {
	select {
	case <-msg.Nacked():
		t.Fatal("nack should be not sent")
	default:
		// ok
	}
}
