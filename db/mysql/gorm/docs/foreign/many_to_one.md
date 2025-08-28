# 声明
```go
// `User` 属于 `Company`，`CompanyID` 是外键
type User struct {
  gorm.Model
  Name      string
  // 注意，在 User 对象中，有一个和 Company 一样的 CompanyID。 默认情况下， CompanyID 被隐含地用来在 User 和 Company 之间创建一个外键关系， 因此必须包含在 User 结构体中才能填充 Company 内部结构体。
  CompanyID int
  Company   Company // 外键
}

type Company struct {
  ID   int
  Name string
}
```

# 重写外键
默认情况下，外键的名字，使用拥有者的类型名称加上表的主键的字段名字
例如，定义一个User实体属于Company实体，那么外键的名字一般使用CompanyID。

```go
type User struct {
  gorm.Model
  Name         string
  CompanyRefer int
  Company      Company `gorm:"foreignKey:CompanyRefer"`
  // 使用 CompanyRefer 作为外键
}

type Company struct {
  ID   int
  Name string
}
```

# 重写引用
```go
type User struct {
  gorm.Model
  Name      string
  CompanyID string
  Company   Company `gorm:"references:Code"` // 使用 Code 作为引用
}

type Company struct {
  ID   int
  Code string
  Name string
}
```

```go
type User struct {
  gorm.Model
  Name      string
  CompanyID string
  Company   Company `gorm:"references:CompanyID"` // 使用 Company.CompanyID 作为引用
}

type Company struct {
  CompanyID   int
  Code        string
  Name        string
}
```

# 外键约束
```go
type User struct {
  gorm.Model
  Name      string
  CompanyID int
    Company   Company `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Company struct {
  ID   int
  Name string
}
```

* OnUpdate
    - CASCADE
        当主表（被引用的表）中的某条记录的主键被更新时，从表（外键表）中所有关联该记录的外键字段会自动同步更新为新的主键值。
        确保数据一致性，避免因主键变更导致从表引用失效。例如，若主表用户ID从1改为101，从表中所有关联该用户的外键字段会从1自动变为101
* OnDelete
    - SET NULL
        当主表中的记录被删除时，从表中所有关联该记录的外键字段会被自动设为NULL（而非删除从表记录）
        保留从表数据的同时解除无效关联，适用于需要保留历史记录但允许外键为空的场景。例如，删除用户时保留其订单记录但清空用户ID字段

      从表的外键字段必须允许为NULL，否则会触发约束错误；必须使用支持外键的存储引擎（如InnoDB）。
    - CASCADE 会级联删除从表记录，可能导致数据丢失
    - RESTRICT 会阻止主表操作，需手动处理依赖
