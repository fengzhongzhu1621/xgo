package client

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func scanKeyDemo(match string) {
	// WithTimeout返回WithDeadline(parent, time.Now(). add (timeout))。
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	// 根据前缀查询 Key
	iter := rdb.Scan(ctx, 0, match, 0).Iterator()

	for iter.Next(ctx) {
		fmt.Printf("key value: %v\n", iter.Val())
	}

	if err := iter.Err(); err != nil {
		fmt.Printf("rdb scan failed, err: %v\n", err)
		return
	}
}

func TestScan(t *testing.T) {
	if err := initRedisV9Client(); err != nil {
		fmt.Printf("initRedisV9Client failed: %v\n", err)
		return
	}
	fmt.Println("initRedisV9Client started successfully")
	defer rdb.Close()

	scanKeyDemo("l*")
}
