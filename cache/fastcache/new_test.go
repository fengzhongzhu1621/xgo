package fastcache

import (
	"fmt"
	"testing"

	"github.com/VictoriaMetrics/fastcache"
)

func TestFastcacheNew(t *testing.T) { // 创建一个最大容量为 1GB 的缓存
	cache := fastcache.New(1 * 1024 * 1024 * 1024) // 1GB
	key := []byte("myKey")
	value := []byte("myValue")

	// 设置缓存
	cache.Set(key, value)
	// 获取缓存
	got := cache.Get(nil, key)
	if got != nil {
		fmt.Printf("Cache hit: %s\n", got)
	} else {
		fmt.Println("Cache miss")
	}

	// 删除缓存（不是强制的）
	cache.Del(key)
}
