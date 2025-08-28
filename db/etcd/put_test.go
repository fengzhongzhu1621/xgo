package etcd

import (
	"context"
	"fmt"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestPut(t *testing.T) {
	// 创建etcd客户端
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"}, // etcd地址
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// 写入共享数据到etcd
	_, err = client.Put(context.Background(), "shared_data", "Hello, Distributed Shared Memory!")
	if err != nil {
		panic(err)
	}

	// 从etcd读取共享数据
	resp, err := client.Get(context.Background(), "shared_data")
	if err != nil {
		panic(err)
	}
	for _, ev := range resp.Kvs {
		fmt.Println("Shared Data:", string(ev.Value))
	}
}
