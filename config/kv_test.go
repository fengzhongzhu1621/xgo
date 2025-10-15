package config

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalKV(t *testing.T) {
	kv := GlobalKV()

	mock := "foo"

	err := kv.Put(context.Background(), "mockString", mock)
	assert.Nil(t, err)

	val, err := kv.Get(context.Background(), "mockString")
	assert.NotNil(t, err)
	assert.Nil(t, val)

	err = kv.Del(context.Background(), "mockString")
	assert.Nil(t, err)

	SetGlobalKV(&mockKV{})
	assert.NotNil(t, GlobalKV())
}
