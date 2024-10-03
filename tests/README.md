# go test

## go test
执行 package 下所有的测试用例

```
go test ./...
```

## -v
显示每个用例的测试结果
```
go test -v
```

## -cover
测试覆盖率

```bash
go test -cover ./user
go test . -cover -v
```

生成测试覆盖率的 profile 文件
```
go test ./... -coverprofile=cover.out
```

利用 profile 文件生成可视化界面
```
go tool cover -html=cover.out
```


## 运行其中的一个用例

```
go test -run TestAdd -v 运行其中的一个用例
go test -timeout 30s -run ^TestAdd$ github.com/fengzhongzhu1621/xgo/tests -v
```

## 子测试(Subtests)

```
go test -run TestMul/pos -v
go test -timeout 30s -run ^TestMul$/^pos$ github.com/fengzhongzhu1621/xgo/tests -v
```


## helper 函数
```
t.Helper()
```
Go 语言在 1.9 版本中引入了 t.Helper()，用于标注该函数是帮助函数，报错时将输出帮助函数调用者的信息，而不是帮助函数的内部信息。

* 不要返回错误， 帮助函数内部直接使用 t.Error 或 t.Fatal 即可，在用例主逻辑中不会因为太多的错误处理代码，影响可读性。
* 调用 t.Helper() 让报错信息更准确，有助于定位。

## 基准测试

```
go test -benchmem -bench .
go test -benchmem -run=^$ -bench ^BenchmarkHello$ github.com/fengzhongzhu1621/xgo/tests -v
```

```go
type BenchmarkResult struct {
    N         int           // 迭代次数
    T         time.Duration // 基准测试花费的时间
    Bytes     int64         // 一次迭代处理的字节数
    MemAllocs uint64        // 总的分配内存的次数
    MemBytes  uint64        // 总的分配内存的字节数
}
```