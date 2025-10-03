package config

import (
	"context"
	"fmt"
	"sync"
)

type mockResponse struct {
	val string
}

func (r *mockResponse) Value() string {
	return r.val
}

func (r *mockResponse) MetaData() map[string]string {
	return nil
}

func (r *mockResponse) Event() EventType {
	return EventTypeNull
}

type mockKV struct {
	mu sync.RWMutex
	db map[string]string
}

// Put mocks putting key and value into kv storage.
func (kv *mockKV) Put(ctx context.Context, key, val string, opts ...Option) error {
	kv.mu.Lock()
	kv.db[key] = val
	kv.mu.Unlock()
	return nil
}

// Get mocks getting value from kv storage by key.
func (kv *mockKV) Get(ctx context.Context, key string, opts ...Option) (IResponse, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	if val, ok := kv.db[key]; ok {
		v := &mockResponse{val: val}
		return v, nil
	}

	return nil, fmt.Errorf("invalid key")
}

// Watch makes mockKV satisfy the KV interface, this method
// is empty.
func (kv *mockKV) Watch(ctx context.Context, key string, opts ...Option) (<-chan IResponse, error) {
	return nil, nil
}

func (kv *mockKV) Name() string {
	return "mock"
}

// Del makes mockKV satisfy the KV interface, this method
// is empty.
func (kv *mockKV) Del(ctx context.Context, key string, opts ...Option) error {
	return nil
}

type mockValue struct {
	Age  int
	Name string
}
