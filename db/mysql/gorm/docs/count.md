# Count

```go
var count int64

// 计数 有着特定名字的 users
db.Model(&User{}).Where("name = ?", "jinzhu").Or("name = ?", "jinzhu 2").Count(&count)
// SQL: SELECT count(1) FROM users WHERE name = 'jinzhu' OR name = 'jinzhu 2'

// 计数 有着单一名字条件（single name condition）的 users
db.Model(&User{}).Where("name = ?", "jinzhu").Count(&count)
// SQL: SELECT count(1) FROM users WHERE name = 'jinzhu'

// 在不同的表中对记录计数
db.Table("deleted_users").Count(&count)
// SQL: SELECT count(1) FROM deleted_users

```

```go
// 为不同 name 计数
db.Model(&User{}).Distinct("name").Count(&count)
// SQL: SELECT COUNT(DISTINCT(`name`)) FROM `users`

// 使用自定义选择（custom select）计数不同的值
db.Table("deleted_users").Select("count(distinct(name))").Count(&count)
// SQL: SELECT count(distinct(name)) FROM deleted_users

// 分组记录计数
users := []User{
  {Name: "name1"},
  {Name: "name2"},
  {Name: "name3"},
  {Name: "name3"},
}

db.Model(&User{}).Group("name").Count(&count)
// 按名称分组后计数
// count => 3
```
