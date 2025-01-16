package forwarder

import (
	"context"
	"testing"

	"github.com/fengzhongzhu1621/xgo/crypto/uuid"

	"github.com/fengzhongzhu1621/xgo/amqp/message"
	"github.com/fengzhongzhu1621/xgo/collections/maps"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvelope(t *testing.T) {
	expectedUUID := uuid.NewUUID()
	expectedPayload := message.Payload("msg content")
	expectedMetadata := maps.Metadata{"key": "value"}
	expectedDestinationTopic := "dest_topic"

	ctx := context.WithValue(context.Background(), "key", "value")

	msg := message.NewMessage(expectedUUID, expectedPayload)
	msg.Metadata = expectedMetadata
	msg.SetContext(ctx)

	wrappedMsg, err := wrapMessageInEnvelope(expectedDestinationTopic, msg)
	require.NoError(t, err)
	require.NotNil(t, wrappedMsg)
	v, ok := wrappedMsg.Context().Value("key").(string)
	require.True(t, ok)
	require.Equal(t, "value", v)

	destinationTopic, unwrappedMsg, err := unwrapMessageFromEnvelope(wrappedMsg)
	require.NoError(t, err)
	require.NotNil(t, unwrappedMsg)
	assert.Equal(t, expectedUUID, unwrappedMsg.UUID)
	assert.Equal(t, expectedPayload, unwrappedMsg.Payload)
	assert.Equal(t, expectedMetadata, unwrappedMsg.Metadata)
	assert.Equal(t, expectedDestinationTopic, destinationTopic)
}
