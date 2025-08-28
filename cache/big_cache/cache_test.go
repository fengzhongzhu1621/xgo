package big_cache

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/allegro/bigcache/v3"
)

// This function is based on [BigCache](https://github.com/allegro/bigcache)
func TestSimple(t *testing.T) {
	// 创建一个新的 BigCache 实例。使用默认配置，并设置缓存项的默认过期时间为 10 分钟。
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	if err != nil {
		log.Fatalf("Failed to create BigCache: %v", err)
	}
	defer cache.Close() // 确保在程序结束时关闭缓存

	// 设置缓存项，并处理可能的错误
	if err := cache.Set("my-unique-key", []byte("value")); err != nil {
		log.Fatalf("Failed to set cache item: %v", err)
	}

	// 获取存在的缓存项，并处理可能的错误
	key := "my-unique-key"
	entry, err := cache.Get(key)
	if err != nil {
		// 在 v3 中，没有 ErrNotFound，需要通过错误消息判断
		if err.Error() == "key not found" {
			fmt.Printf("键 '%s' 在缓存中未找到\n", key)
		} else {
			log.Fatalf("获取缓存项失败: %v", err)
		}
	} else {
		fmt.Printf("键 '%s' 的值为: %s\n", key, string(entry))
	}

	// 尝试获取一个不存在的键，并处理错误
	nonExistentKey := "non-existent-key"
	entry, err = cache.Get(nonExistentKey)
	if err != nil {
		// 检查错误消息是否为 "key not found"
		if err.Error() == "key not found" {
			fmt.Printf("键 '%s' 在缓存中未找到\n", nonExistentKey)
		} else {
			log.Fatalf("获取缓存项失败: %v", err)
		}
	} else {
		fmt.Printf("键 '%s' 的值为: %s\n", nonExistentKey, string(entry))
	}

	fmt.Println(string(entry))
}

// This function is based on [BigCache](https://github.com/allegro/bigcache)
func TestConfig(t *testing.T) {
	config := bigcache.Config{
		// number of shards (must be a power of 2)
		Shards: 1024,

		// time after which entry can be evicted
		LifeWindow: 10 * time.Minute,

		// Interval between removing expired entries (clean up).
		// If set to <= 0 then no action is performed.
		// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
		CleanWindow: 5 * time.Minute,

		// rps * lifeWindow, used only in initial memory allocation
		MaxEntriesInWindow: 1000 * 10 * 60,

		// max entry size in bytes, used only in initial memory allocation
		MaxEntrySize: 500,

		// prints information about additional memory allocation
		Verbose: true,

		// cache will not allocate more memory than this limit, value in MB
		// if value is reached then the oldest entries can be overridden for the new ones
		// 0 value means no size limit
		HardMaxCacheSize: 8192,

		// callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A bitmask representing the reason will be returned.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		OnRemove: nil,

		// OnRemoveWithReason is a callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A constant representing the reason will be passed through.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		// Ignored if OnRemove is specified.
		OnRemoveWithReason: nil,
	}

	cache, initErr := bigcache.New(context.Background(), config)
	if initErr != nil {
		log.Fatal(initErr)
	}

	cache.Set("my-unique-key", []byte("value"))

	if entry, err := cache.Get("my-unique-key"); err == nil {
		fmt.Println(string(entry))
	}
}
