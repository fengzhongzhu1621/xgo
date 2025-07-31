# 关联创建
创建关联数据时，如果关联值非零，这些关联会被upsert，并且它们的Hooks方法也会被调用。

```go
user := User{
  Name:            "jinzhu",
  BillingAddress:  Address{Address1: "Billing Address - Address 1"},
  ShippingAddress: Address{Address1: "Shipping Address - Address 1"},
  Emails:          []Email{
    {Email: "jinzhu@example.com"},
    {Email: "jinzhu-2@example.com"},
  },
  Languages:       []Language{
    {Name: "ZH"},
    {Name: "EN"},
  },
}

// 创建用户及其关联的地址、电子邮件和语言
db.Create(&user)
// BEGIN TRANSACTION;
// INSERT INTO "addresses" (address1) VALUES ("Billing Address - Address 1"), ("Shipping Address - Address 1") ON DUPLICATE KEY DO NOTHING;
// INSERT INTO "users" (name,billing_address_id,shipping_address_id) VALUES ("jinzhu", 1, 2);
// INSERT INTO "emails" (user_id,email) VALUES (111, "jinzhu@example.com"), (111, "jinzhu-2@example.com") ON DUPLICATE KEY DO NOTHING;
// INSERT INTO "languages" ("name") VALUES ('ZH'), ('EN') ON DUPLICATE KEY DO NOTHING;
// INSERT INTO "user_languages" ("user_id","language_id") VALUES (111, 1), (111, 2) ON DUPLICATE KEY DO NOTHING;
// COMMIT;

db.Save(&user)
```

```go
type CreditCard struct {
  gorm.Model
  Number   string
  UserID   uint
}

type User struct {
  gorm.Model
  Name       string
  CreditCard CreditCard
}

db.Create(&User{
  Name: "jinzhu",
  CreditCard: CreditCard{Number: "411111111111"},
})
// INSERT INTO `users` ...
// INSERT INTO `credit_cards` ...
```

Omit方法来跳过关联更新
```go
db.Omit("CreditCard").Create(&user)

// skip all associations
db.Omit(clause.Associations).Create(&user)
```

# 查询关联
```go
// Simple find
db.Model(&user).Association("Languages").Find(&languages)

// Find with conditions
codes := []string{"zh-CN", "en-US", "ja-JP"}
db.Model(&user).Where("code IN ?", codes).Association("Languages").Find(&languages)
```

# 关联计数
获取当前关联的数量，可以带有或不带有条件。
```go
// Count all languages
db.Model(&user).Association("Languages").Count()

// Count with conditions
codes := []string{"zh-CN", "en-US", "ja-JP"}
db.Model(&user).Where("code IN ?", codes).Association("Languages").Count()
```

# 替换关联
用新的关联替换当前的关联。
```go
// Replace existing languages
db.Model(&user).Association("Languages").Replace([]Language{languageZH, languageEN})
db.Model(&user).Association("Languages").Replace(Language{Name: "DE"}, languageEN)
```

# 删除关联
删除源模型和目标模型之间的关系，只删除引用。

**注意：**请务必注意，仅当删除记录的主键不为零时，才会删除关联。 GORM 使用这些主键作为删除所选关联的条件。


```go
// Delete specific languages
db.Model(&user).Association("Languages").Delete([]Language{languageZH, languageEN})
db.Model(&user).Association("Languages").Delete(languageZH, languageEN)
```

```go
//删除用户时，也删除用户的帐户
db.Select("Account").Delete(&user)

// 删除 user 时，也删除 user 的 Orders、CreditCards 关联记录
db.Select("Orders", "CreditCards").Delete(&user)

// 删除用户 时，也删除用户的所有一对一、一对多和多对多关联
db.Select(clause.Associations).Delete(&user)

// 删除多个用户时，也同步删除每个用户的账户信息。
db.Select("Account").Delete(&users)

// 这样无法实现预期效果
db.Select("Account").Where("name = ?", "jinzhu").Delete(&User{})
// SQL: 删除所有名为 'jinzhu' 的用户，但这些用户的关联账户不会被删除

// 正确删除用户及其账户的方式
db.Select("Account").Where("name = ?", "jinzhu").Delete(&User{ID: 1})
// SQL: 删除名为 'jinzhu' 且 ID 为 '1' 的用户及其关联账户

// 删除指定 ID 用户及其账户db.Select("Account").Delete(&User{ID: 1})
// SQL: 删除 ID 为 '1' 的用户及其关联账户
```

# 清空关联
删除源模型和关联之间的所有引用。
```go
// Clear all languages
db.Model(&user).Association("Languages").Clear()
```

# 预加载
```go
type User struct {
  gorm.Model
  Username string
  Orders   []Order
}

type Order struct {
  gorm.Model
  UserID uint
  Price  float64
}

// 查找 user 时预加载相关 Order
db.Preload("Orders").Find(&users)
// SELECT * FROM users;
// SELECT * FROM orders WHERE user_id IN (1,2,3,4);

db.Preload("Orders").Preload("Profile").Preload("Role").Find(&users)
// SELECT * FROM users;
// SELECT * FROM orders WHERE user_id IN (1,2,3,4); // has many
// SELECT * FROM profiles WHERE user_id IN (1,2,3,4); // has one
// SELECT * FROM roles WHERE id IN (4,5,6); // belongs to
```

# Joins 预加载
```go
db.Joins("Company").Joins("Manager").Joins("Account").First(&user, 1)
db.Joins("Company").Joins("Manager").Joins("Account").First(&user, "users.name = ?", "jinzhu")
db.Joins("Company").Joins("Manager").Joins("Account").Find(&users, "users.id IN ?", []int{1,2,3,4,5})
```

```go
db.Joins("Company", DB.Where(&Company{Alive: true})).Find(&users)
// SELECT `users`.`id`,`users`.`name`,`users`.`age`,`Company`.`id` AS `Company__id`,`Company`.`name` AS `Company__name` FROM `users` LEFT JOIN `companies` AS `Company` ON `users`.`company_id` = `Company`.`id` AND `Company`.`alive` = true;
```

```go
db.Joins("Manager").Joins("Manager.Company").Find(&users)
// SELECT "users"."id","users"."created_at","users"."updated_at","users"."deleted_at","users"."name","users"."age","users"."birthday","users"."company_id","users"."manager_id","users"."active","Manager"."id" AS "Manager__id","Manager"."created_at" AS "Manager__created_at","Manager"."updated_at" AS "Manager__updated_at","Manager"."deleted_at" AS "Manager__deleted_at","Manager"."name" AS "Manager__name","Manager"."age" AS "Manager__age","Manager"."birthday" AS "Manager__birthday","Manager"."company_id" AS "Manager__company_id","Manager"."manager_id" AS "Manager__manager_id","Manager"."active" AS "Manager__active","Manager__Company"."id" AS "Manager__Company__id","Manager__Company"."name" AS "Manager__Company__name" FROM "users" LEFT JOIN "users" "Manager" ON "users"."manager_id" = "Manager"."id" AND "Manager"."deleted_at" IS NULL LEFT JOIN "companies" "Manager__Company" ON "Manager"."company_id" = "Manager__Company"."id" WHERE "users"."deleted_at" IS NULL
```

# 预加载全部
```go
type User struct {
  gorm.Model
  Name       string
  CompanyID  uint
  Company    Company
  Role       Role
  Orders     []Order
}

db.Preload(clause.Associations).Find(&users)

db.Preload("Orders.OrderItems.Product").Preload(clause.Associations).Find(&users)
```

# 条件预加载
```go
// Preload Orders with conditions
db.Preload("Orders", "state NOT IN (?)", "cancelled").Find(&users)
// SELECT * FROM users;
// SELECT * FROM orders WHERE user_id IN (1,2,3,4) AND state NOT IN ('cancelled');

db.Where("state = ?", "active").Preload("Orders", "state NOT IN (?)", "cancelled").Find(&users)
// SELECT * FROM users WHERE state = 'active';
// SELECT * FROM orders WHERE user_id IN (1,2) AND state NOT IN ('cancelled');
```

# 自定义预加载 SQL
```go
db.Preload("Orders", func(db *gorm.DB) *gorm.DB {
  return db.Order("orders.amount DESC")
}).Find(&users)
// SELECT * FROM users;
// SELECT * FROM orders WHERE user_id IN (1,2,3,4) order by orders.amount DESC;
```

# 嵌套预加载
```go
db.Preload("Orders.OrderItems.Product").Preload("CreditCard").Find(&users)

// Customize Preload conditions for `Orders`
// And GORM won't preload unmatched order's OrderItems then
db.Preload("Orders", "state = ?", "paid").Preload("Orders.OrderItems").Find(&users)
```

# Embedded Preloading
```go
type Address struct {
    CountryID int
    Country   Country
}

type Org struct {
    PostalAddress   Address `gorm:"embedded;embeddedPrefix:postal_address_"`
    VisitingAddress Address `gorm:"embedded;embeddedPrefix:visiting_address_"`
    Address         struct {
        ID int
        Address
    }
}

// Only preload Org.Address and Org.Address.Country
db.Preload("Address.Country")  // "Address" is has_one, "Country" is belongs_to (nested association)

// Only preload Org.VisitingAddress
db.Preload("PostalAddress.Country") // "PostalAddress.Country" is belongs_to (embedded association)

// Only preload Org.NestedAddress
db.Preload("NestedAddress.Address.Country") // "NestedAddress.Address.Country" is belongs_to (embedded association)

// All preloaded include "Address" but exclude "Address.Country", because it won't preload nested associations.
db.Preload(clause.Associations)
```

```go
// We can omit embedded part when there is no ambiguity.
type Address struct {
    CountryID int
    Country   Country
}

type Org struct {
    Address Address `gorm:"embedded"`
}

db.Preload("Address.Country")
db.Preload("Country") // omit "Address" because there is no ambiguity
```
