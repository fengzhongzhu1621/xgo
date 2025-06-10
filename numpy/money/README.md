# 简介

go-money要求使用ISO 4217货币代码来设定货币类型，并为所有ISO 4217标准货币代码提供了常量。

Money 实例是不可变的，每次运算都会生成新的实例，这有助于保持数据的纯净性和一致性。

```sh
go get github.com/Rhymond/go-money
```

# 初始化

go-money允许通过两种方式初始化货币值：
* 使用货币最小单位（例如100代表1英镑）
* 使用直接金额（例如使用浮点数表示）

```go
// 使用最小单位初始化（100代表1英镑）
pound := money.New(100, money.GBP)

// 使用浮点数直接初始化
quarterEuro := money.NewFromFloat(0.25, money.EUR)
```

# 货币比较

比较操作必须在相同的货币单位之间进行。

* 等于（Equals）
* 大于（GreaterThan）
* 大于或等于（GreaterThanOrEqual）
* 小于（LessThan）
* 小于或等于（LessThanOrEqual）
* 比较（Compare）

# 基本运算

* 加法（Add）
* 减法（Subtract）
* 乘法（Multiply）
* 绝对值（Absolute）
* 负数值（Negative）
