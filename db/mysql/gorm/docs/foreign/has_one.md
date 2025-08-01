# 声明
has one 与另一个模型建立一对一的关联，但它和一对一关系有些许不同。 这种关联表明一个模型的每个实例都包含或拥有另一个模型的一个实例。

对于 has one 关系，同样必须存在外键字段。拥有者将把属于它的模型的主键保存到这个字段。
这个字段的名称通常由 has one 模型的类型加上其 主键 生成，对于上面的例子，它是 UserID。
为 user 添加 credit card 时，它会将 user 的 ID 保存到自己的 UserID 字段。

```go
// User 有一张 CreditCard，UserID 是外键
type User struct {
  gorm.Model
  CreditCard CreditCard // 将 user 的 ID 保存到CreditCard的 UserID 字段
}

type CreditCard struct {
  gorm.Model
  Number string
  UserID uint
}
```

# 检索
```go
// 检索用户列表并预加载信用卡
func GetAll(db *gorm.DB) ([]User, error) {
    var users []User
    err := db.Model(&User{}).Preload("CreditCard").Find(&users).Error
    return users, err
}
```

# 重写外键

如果你想要使用另一个字段来保存该关系，你同样可以使用标签 foreignKey 来更改它
```go
type User struct {
  gorm.Model
  CreditCard CreditCard `gorm:"foreignKey:UserName"` // 使用 UserName 作为外键
}

type CreditCard struct {
  gorm.Model
  Number   string
  UserName string
}
```

# 重写引用

默认情况下，拥有者实体会将 has one 对应模型的主键保存为外键，您也可以修改它，用另一个字段来保存，例如下面这个使用 Name 来保存的例子

```go
type User struct {
  gorm.Model
  Name       string     `gorm:"index"`
  CreditCard CreditCard `gorm:"foreignKey:UserName;references:Name"`
}

type CreditCard struct {
  gorm.Model
  Number   string
  UserName string
}
```

# 外键约束

```go
type User struct {
  gorm.Model
  CreditCard CreditCard `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type CreditCard struct {
  gorm.Model
  Number string
  UserID uint
}
```


# 自引用 Has One
```go
type User struct {
  gorm.Model
  Name      string
  ManagerID *uint
  Manager   *User
}
```
