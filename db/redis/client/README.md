# 语法
## TLS
```go
rdb := redis.NewClient(&redis.Options{
    TLSConfig: &tls.Config{
        MinVersion: tls.VersionTLS12,
        ServerName:  "you domain",
        //Certificates: []tls.Certificate{cert}
    },
})
```

## Do()

Do()  方法返回  Cmd  类型，使用下面的命令获取想要的类型

```go
s, err := cmd.Text()
flag, err := cmd.Bool()

num, err := cmd.Int()
num, err := cmd.Int64()
num, err := cmd.Uint64()
num, err := cmd.Float32()
num, err := cmd.Float64()

ss, err := cmd.StringSlice()
ns, err := cmd.Int64Slice()
ns, err := cmd.Uint64Slice()
fs, err := cmd.Float32Slice()
fs, err := cmd.Float64Slice()
bs, err := cmd.BoolSlice()
```

```go
val, err := rdb.Do(ctx,  "get",  "key").Result()
if  err !=  nil  {
    if  err == redis.Nil {
        fmt.Println("key does not exists")
        return
    }
    panic(err)
}
fmt.Println(val.(string))
```

## redis.Nil
表一种状态，例如你使用 Get 命令获取 key 的值，当 key 不存在时，返回  redis.Nil。需要自行区分。

```go
val, err := rdb.Get(ctx,  "key").Result()
switch  {
    case  err == redis.Nil:
        fmt.Println("key不存在")
    case  err !=  nil:
        fmt.Println("错误", err)
    case  val ==  "":
        fmt.Println("值是空字符串")
}
```


