# 简介
是一个功能强大且易于使用的内存缓存库，适用于单机应用程序。
```go
import (
    gocache "github.com/patrickmn/go-cache"
)
```

# 特性
* 过期时间与自动清理：每个缓存项可以设置独立的过期时间，支持永久存储（NoExpiration）和默认过期时间（DefaultExpiration）。库内部通过定时任务自动清理过期数据，避免内存泄漏。
* 线程安全与高性能：通过 sync.RWMutex 实现读写锁，确保多个 goroutine 可以安全地并发访问缓存。尽管在高并发场景下可能存在锁竞争，但其性能仍优于许多其他缓存方案。

# 语法

## 创建缓存实例
```go
const (
// For use with functions that take an expiration time.
// 表示缓存条目永不过期。
NoExpiration time.Duration = -1
// For use with functions that take an expiration time. Equivalent to
// passing in the same expiration duration as was given to New() or
// NewFrom() when the cache was created (e.g. 5 minutes.)
// 表示使用缓存实例创建时定义的默认过期时间（例如 5*time.Minute）。
DefaultExpiration time.Duration = 0
)

// 默认过期时间（第一个参数）：缓存条目在没有显式设置过期时间时，默认的存活时间。
// 清理间隔（第二个参数）：缓存后台清理过期条目的时间间隔。
c := gocache.New(12*time.Hour, 5*time.Minute)
```

## 存储数据
```go
// 使用默认的过期时间
c.Set("key","value", cache.DefaultExpiration)

// 使用默认的过期时间
c.SetDefault("key", "value")

// 使用自定义过期时间
c.Set("key", "value", 10*time.Minute)

// 不过期
c.Set("key", "value", cache.NoExpiration)


// 存储结构体（推荐存储指针，提高性能）
type User struct {
    ID int
    Name string
}
user := &User{ID: 1, Name: "name"}
c.Set("user:1", user, cache.DefaultExpiration)
if x, found := c.Get("user:1"); found {
    u := x.(*User)
    fmt.Println("user:", u.Name)
}
```

## 读取数据
```go
value, ok := c.Get("key")
if ok {
    fmt.Println("value:", value)// value: abc
}
```

## 删除数据
```go
c.Delete("key")
```

## 获取过期时间
```go
// 从缓存中获取指定键的值及其过期时间
if value, expiration, found := c.GetWithExpiration("foo"); found {
    fmt.Printf("Key: foo, Value: %v, Expiration: %v
", value, expiration)
} else {
    fmt.Println("Key 'foo' not found or expired")
}
```

## 递增/减
```go
// 将 int 类型 value 值递增/减
c.IncrementInt("key", 10)
c.DecrementInt("key", 10)
```

## 获取所有缓存数据
```go
for k, v := rang ec.Items() {
    fmt.Printf("Key: %s, Value: %v\n", k, v.Object)
}
```
