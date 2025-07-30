# 定义 Scopes
允许您将常用的查询条件定义为可重用的方法。 这些作用域可以很容易地在查询中引用，从而使代码更加模块化和可读。
Scopes 被定义为被修改后返回一个 gorm.DB 实例的函数。 您可以根据您的应用程序的需要定义各种条件作为范围。

```go
// 用于筛选 amount 大于 1000 的记录范围
func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
  return db.Where("amount > ?", 1000)
}

// 用于 订单使用信用卡支付 的 Scope
func PaidWithCreditCard(db *gorm.DB) *gorm.DB {
  return db.Where("pay_mode_sign = ?", "C")
}

// 用于 订单货到付款(COD)的 Scope
func PaidWithCod(db *gorm.DB) *gorm.DB {
  return db.Where("pay_mode_sign = ?", "C")
}

//  用于按状态筛选订单的 Scope
func OrderStatus(status []string) func(db *gorm.DB) *gorm.DB {
  return func(db *gorm.DB) *gorm.DB {
    return db.Where("status IN (?)", status)
  }
}
```

# 在查询中使用 Scopes
```go
// 使用 scopes 来寻找所有的 金额大于1000的信用卡订单
db.Scopes(AmountGreaterThan1000, PaidWithCreditCard).Find(&orders)

// 使用 scopes 来寻找所有的 金额大于1000的货到付款（COD）订单
db.Scopes(AmountGreaterThan1000, PaidWithCod).Find(&orders)

//使用 scopes 来寻找所有的 具有特定状态且金额大于1000的订单
db.Scopes(AmountGreaterThan1000, OrderStatus([]string{"paid", "shipped"})).Find(&orders)
```
