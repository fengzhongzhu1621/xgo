package client

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func zSetDemo() {
	// key
	zSetKey := "language_rank"
	// value Z表示有序集合的成员。
	languages := []redis.Z{
		{Score: 90.0, Member: "Golang"},
		{Score: 95.0, Member: "Python"},
		{Score: 97.0, Member: "Rust"},
		{Score: 99.0, Member: "C/C++"},
		{Score: 88.0, Member: "Java"},
	}
	// WithTimeout返回WithDeadline(parent, time.Now(). add (timeout))。
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// ZAdd Redis `ZADD key score member [score member ...]` command.
	num, err := rdb.ZAdd(ctx, zSetKey, languages...).Result()
	if err != nil {
		fmt.Printf("zadd failed, err:%v\n", err)
		return
	}
	fmt.Printf("zadd successful num: %v\n", num)

	// ZIncrBy 给某一个元素添加分数值 把Golang的分数加 10
	newScore, err := rdb.ZIncrBy(ctx, zSetKey, 10.0, "Golang").Result()
	if err != nil {
		fmt.Printf("ZIncrBy failed, err:%v\n", err)
		return
	}
	fmt.Printf("ZIncrBy success Golang's score is %f now.\n", newScore)

	// 取分数最高的3个   适用于 排行榜、充值榜...
	ret, err := rdb.ZRevRangeWithScores(ctx, zSetKey, 0, 2).Result()
	if err != nil {
		fmt.Printf("zRevRangeWithScores failed, err: %v\n", err)
		return
	}
	for _, z := range ret {
		fmt.Printf("z.Member: %v, z.Score: %v\n", z.Member, z.Score)
	}

	// 取95~100分的
	op := &redis.ZRangeBy{
		Min: "95",
		Max: "100",
	}
	ret, err = rdb.ZRangeByScoreWithScores(ctx, zSetKey, op).Result()
	if err != nil {
		fmt.Printf("zrangebyscore failed, err:%v\n", err)
		return
	}
	fmt.Printf("zrangebyscore returned %v\n", ret)
	for _, z := range ret {
		fmt.Printf("ZRangeByScoreWithScores success Member: %v, Score: %v\n", z.Member, z.Score)
	}
}

func TestSortedSet(t *testing.T) {
	if err := initRedisV9Client(); err != nil {
		fmt.Printf("initRedisV9Client failed: %v\n", err)
		return
	}
	fmt.Println("initRedisV9Client started successfully")
	defer rdb.Close() // Close 关闭客户端，释放所有打开的资源。关闭客户端是很少见的，因为客户端是长期存在的，并在许多例程之间共享。

	zSetDemo()
}
