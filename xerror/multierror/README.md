#  简介
go-multierror 是 HashiCorp 提供的一个 Go 语言库，它能将多个错误合并成一个标准 error，非常适合需要收集和统一管理多个错误的场景。

兼容 Go 1.13+ 错误链（errors.Is, errors.As, errors.Unwrap），支持错误格式化，提供并发安全方案。

go-multierror 的核心在于 Append 函数，它的行为类似 Go 内置的 append，可以智能地处理 nil 错误和嵌套的多错误。


# 与标准库 errors.Join 的比较
Go 1.20 引入了 errors.Join 函数，它也能返回一个包含多个错误的 error。两者的主要区别在于：

错误格式：errors.Join 返回的错误字符串是所有错误的简单拼接。go-multierror 默认提供结构化的列表格式，并且支持自定义。
功能丰富性：go-multierror 提供了如错误展平（Flatten）、添加上下文前缀（Prefix）等更多辅助功能。
