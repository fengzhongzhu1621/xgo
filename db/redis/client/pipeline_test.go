package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
)

// 管道里的命令是批量执行的，一个命令出错不会影响其他命令。
// 比如，如果 SET 失败，GET 还是会继续执行。
// 所以，执行完管道后，要逐一检查每个命令的 Err()，确保没漏掉问题。
func TestPipeline(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// 创建管道
	pipe := rdb.Pipeline()

	// 添加命令到管道
	inc := pipe.Incr(ctx, "counter")        // 自增 counter
	set := pipe.Set(ctx, "key", "value", 0) // 设置 key=value
	get := pipe.Get(ctx, "key")             // 获取 key 的值
	setCmd := pipe.Set(ctx, "name", "Redis", 0)
	getCmd := pipe.Get(ctx, "name")
	incrCmd := pipe.Incr(ctx, "visits")
	hsetCmd := pipe.HSet(ctx, "user:1", "age", 30)

	// 执行管道
	_, err := pipe.Exec(ctx)
	if err != nil {
		panic(err)
	}

	// 查看结果
	fmt.Println("counter 自增后:", inc.Val())
	fmt.Println("set 命令错误:", set.Err())
	fmt.Println("key 的值:", get.Val())
	// counter 自增后: 9
	// set 命令错误: <nil>
	// key 的值: value

	// 检查结果
	fmt.Println("SET 命令错误:", setCmd.Err())
	fmt.Println("GET 命令结果:", getCmd.Val())
	fmt.Println("INCR 命令结果:", incrCmd.Val())
	fmt.Println("HSET 命令错误:", hsetCmd.Err())
	// SET 命令错误: <nil>
	// GET 命令结果: Redis
	// INCR 命令结果: 5
	// HSET 命令错误: <nil>

	// 添加命令
	pipe2 := rdb.Pipeline()
	setCmd2 := pipe2.Set(ctx, "key1", "value1", 0)
	getCmd2 := pipe2.Get(ctx, "nonexistent") // 不存在的 key
	incrCmd2 := pipe2.Incr(ctx, "counter")

	///////////////////////////////////////////////////////////////////
	// 执行管道：存在失败的命令
	_, err = pipe2.Exec(ctx)
	if err != nil {
		fmt.Println("pipe2 执行失败:", err) // pipe2 执行失败: redis: nil
	}

	// 检查每个命令的结果
	if setCmd2.Err() != nil {
		fmt.Println("SET 命令失败:", setCmd2.Err())
	}
	if getCmd2.Err() != nil {
		if getCmd2.Err() == redis.Nil {
			fmt.Println("GET 命令: key 不存在")
		} else {
			fmt.Println("GET 命令失败:", getCmd2.Err())
		}
	} else {
		fmt.Println("GET 命令结果:", getCmd2.Val())
	}
	if incrCmd2.Err() != nil {
		fmt.Println("INCR 命令失败:", incrCmd2.Err())
	} else {
		fmt.Println("INCR 命令结果:", incrCmd2.Val())
	}
	// pipe2 执行失败: redis: nil
	// GET 命令: key 不存在
	// INCR 命令结果: 10
}

// 批量设置多个 key-value 对
func TestPipelineBatchSet(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// 准备数据
	data := map[string]interface{}{
		"user:1": "Alice",
		"user:2": "Bob",
		"user:3": "Charlie",
	}

	// 创建管道
	pipe := rdb.Pipeline()

	// 添加 SET 命令到管道
	for k, v := range data {
		pipe.Set(ctx, k, v, 0) // 0 表示永不过期
	}

	// 执行管道中的所有命令，并返回结果。
	// 第一个返回值（这里用_忽略）是每个命令的结果，第二个返回值是错误。
	_, err := pipe.Exec(ctx)
	if err != nil {
		panic(err)
	}

	// 验证结果
	for k := range data {
		// .Result()：获取命令的执行结果（返回值和错误）
		val, err := rdb.Get(ctx, k).Result()
		if err != nil {
			fmt.Printf("获取 %s 失败: %v\n", k, err)
		} else {
			fmt.Printf("%s: %s\n", k, val)
		}
	}
}

// 批量获取多个 key 的值
func TestPipelineBatchGet(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// 假设已有数据
	keys := []string{"item:1:price", "item:2:price", "item:3:price"}

	// 创建管道
	pipe := rdb.Pipeline()

	// 添加 GET 命令到管道
	getResultList := make([]*redis.StringCmd, len(keys))
	for i, k := range keys {
		// Get()返回的是 *redis.StringCmd 类型
		getResultList[i] = pipe.Get(ctx, k)
	}

	// 执行管道
	_, err := pipe.Exec(ctx)
	if err != nil {
		panic(err)
	}

	// 获取结果
	for i, getResult := range getResultList {
		if getResult.Err() != nil {
			fmt.Printf("获取 %s 失败: %v\n", keys[i], getResult.Err())
		} else {
			fmt.Printf("%s: %s\n", keys[i], getResult.Val())
		}
	}
}
