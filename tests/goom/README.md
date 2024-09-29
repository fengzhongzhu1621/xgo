# 简介
GOOM单测Mock框架

# 功能特性

1. mock过程中调用原函数(线程安全, 支持并发单测)
2. 异常注入，对函数调用支持异常注入，延迟模拟等稳定性测试
3. 所有操作都是并发安全的
4. 未导出(未导出)函数(或方法)的mock(不建议使用, 对于未导出函数的Mock 通常都是因为代码设计可能有问题, 此功能会在未来版本中废弃)
5. 支持M1 mac环境运行，支持IDE debug，函数、方法mock，接口mock，未导出函数mock，等能力均可在arm64架构上使用

# 注意
注意: 按照go编译规则，短函数会被内联优化，导致无法mock的情况，编译参数需要加上 -gcflags=all=-l 关闭内联
例如: go test -gcflags=all=-l hello.go


# 引用
```go
// 在需要使用mock的测试文件import
import "github.com/tencent/goom"
```

## 函数mock

