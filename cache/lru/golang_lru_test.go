package lru

import (
	"fmt"
	"testing"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/hashicorp/golang-lru/v2/expirable"
)

func TestGolongLRUCache(t *testing.T) {
	// 创建一个最大容量为128的LRU缓存
	l, err := lru.New[int, any](128)
	if err != nil {
		panic(fmt.Sprintf("failed to create LRU cache: %v", err))
	}

	// 向缓存中添加256个元素（超过容量限制）
	for i := 0; i < 256; i++ {
		l.Add(i, nil)
	}

	// 验证缓存长度是否符合预期（应该只保留最近使用的128个元素）
	if l.Len() != 128 {
		panic(fmt.Sprintf("bad len: %v", l.Len()))
	}
}

func TestGolongLRUCacheWithExpirable(t *testing.T) {
	// 创建一个带有10毫秒TTL和最大5个键的LRU缓存
	cache := expirable.NewLRU[string, string](5, nil, time.Millisecond*10)

	// 在键"key1"下设置值
	cache.Add("key1", "val1")

	// 从键"key1"获取值
	r, ok := cache.Get("key1")

	// 检查是否成功获取值
	if ok {
		fmt.Printf("value before expiration is found: %v, value: %q\n", ok, r)
	}

	// 等待缓存过期
	time.Sleep(time.Millisecond * 12)

	// 在键"key1"过期后获取值
	r, ok = cache.Get("key1")
	fmt.Printf("value after expiration is found: %v, value: %q\n", ok, r)

	// 在键"key2"下设置值，由于旧条目已过期，这将淘汰旧条目
	cache.Add("key2", "val2")

	// 打印当前缓存长度
	fmt.Printf("Cache len: %d\n", cache.Len())
}
