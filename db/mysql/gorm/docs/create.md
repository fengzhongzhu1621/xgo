# 错误

## mysql.MySQLError
```go
import (
    "github.com/go-sql-driver/mysql"
    "gorm.io/gorm"
)

result := db.Create(&newRecord)
if result.Error != nil {
    if mysqlErr, ok := result.Error.(*mysql.MySQLError); ok {
        switch mysqlErr.Number {
        case 1062: // MySQL中表示重复条目的代码
            // 处理重复条目
        // 为其他特定错误代码添加案例
        default:
            // 处理其他错误
        }
    } else {
        // 处理非MySQL错误或未知错误
    }
}
```

## 方言转换错误
当启用TranslateError时，GORM可以返回与所使用的数据库方言相关的特定错误，GORM将数据库特有的错误转换为其自己的通用错误。

```go
db, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{TranslateError: true})
```

### ErrDuplicatedKey
当插入操作违反唯一约束时，会发生此错误

```go
result := db.Create(&newRecord)
if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
    // 处理重复键错误...
}
```

### ErrForeignKeyViolated
当违反外键约束时，会遇到此错误

```go
result := db.Create(&newRecord)
if errors.Is(result.Error, gorm.ErrForeignKeyViolated) {
    // 处理外键违规错误...
}
```


# 通过数据的指针来创建
```go
user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}

result := db.Create(&user) // 通过数据的指针来创建

user.ID             // 返回插入数据的主键
result.Error        // 返回 error
result.RowsAffected // 返回插入记录的条数
```

# 创建多项记录
要高效地插入大量记录，请将切片传递给Create方法。 GORM 将生成一条 SQL 来插入所有数据，以返回所有主键值，并触发 Hook 方法。 当这些记录可以被分割成多个批次时，GORM会开启一个事务来处理它们。

```go
users := []*User{
    {Name: "Jinzhu", Age: 18, Birthday: time.Now()},
    {Name: "Jackson", Age: 19, Birthday: time.Now()},
}

result := db.Create(users) // pass a slice to insert multiple row

result.Error        // returns error
result.RowsAffected // returns inserted records count

for _, user := range users {
  user.Name
}
```

通过db.CreateInBatches方法来指定批量插入的批次大小
```go
var users = []User{{Name: "jinzhu_1"}, ...., {Name: "jinzhu_10000"}}

// batch size 100
db.CreateInBatches(users, 100)
```


使用CreateBatchSize 选项初始化GORM实例，此后进行创建和关联操作时所有的INSERT行为都会遵循初始化时的配置。
```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
  CreateBatchSize: 1000,
})

db := db.Session(&gorm.Session{CreateBatchSize: 1000})

users = [5000]User{{Name: "jinzhu", Pets: []Pet{pet1, pet2, pet3}}...}

db.Create(&users)
// INSERT INTO users xxx (5 batches)
// INSERT INTO pets xxx (15 batches)
```

# 用指定的字段创建记录
```go
创建记录并为指定字段赋值。

db.Select("Name", "Age", "CreatedAt").Create(&user)
// INSERT INTO `users` (`name`,`age`,`created_at`) VALUES ("jinzhu", 18, "2020-07-04 11:05:21.775")
创建记录并忽略传递给 ‘Omit’ 的字段值

db.Omit("Name", "Age", "CreatedAt").Create(&user)
// INSERT INTO `users` (`birthday`,`updated_at`) VALUES ("2020-01-01 00:00:00.000", "2020-07-04 11:05:21.775")

// 创建用户时跳过字段“BillingAddress”
db.Omit("BillingAddress").Create(&user)

// 创建用户时跳过全部关联关系
db.Omit(clause.Associations).Create(&user)

// 跳过更新"Languages"关联
db.Omit("Languages.*").Create(&user)

// 跳过创建 'Languages' 关联及其引用
db.Omit("Languages").Create(&user)

user := User{
  Name:            "jinzhu",
  BillingAddress:  Address{Address1: "Billing Address - Address 1", Address2: "addr2"},
  ShippingAddress: Address{Address1: "Shipping Address - Address 1", Address2: "addr2"},
}

// 创建用户和他的账单地址,邮寄地址,只包括账单地址指定的字段
db.Select("BillingAddress.Address1", "BillingAddress.Address2").Create(&user)
// SQL: 只使用地址1和地址2来创建用户和账单地址

// 创建用户和账单地址,邮寄地址,但不包括账单地址的指定字段
db.Omit("BillingAddress.Address2", "BillingAddress.CreatedAt").Create(&user)
// SQL: 创建用户和账单地址,省略'地址2'和创建时间字段
```

# 根据 Map 创建
注意当使用map来创建时，钩子方法不会执行，关联不会被保存且不会回写主键。

```go
db.Model(&User{}).Create(map[string]interface{}{
  "Name": "jinzhu", "Age": 18,
})

// batch insert from `[]map[string]interface{}{}`
db.Model(&User{}).Create([]map[string]interface{}{
  {"Name": "jinzhu_1", "Age": 18},
  {"Name": "jinzhu_2", "Age": 20},
})
```

# 使用 SQL 表达式、Context Valuer 创建记录
```go
// Create from map
db.Model(User{}).Create(map[string]interface{}{
  "Name": "jinzhu",
  "Location": clause.Expr{SQL: "ST_PointFromText(?)", Vars: []interface{}{"POINT(100 100)"}},
})
// INSERT INTO `users` (`name`,`location`) VALUES ("jinzhu",ST_PointFromText("POINT(100 100)"));

// Create from customized data type
type Location struct {
    X, Y int
}

// Scan implements the sql.Scanner interface
func (loc *Location) Scan(v interface{}) error {
  // Scan a value into struct from database driver
}

func (loc Location) GormDataType() string {
  return "geometry"
}

func (loc Location) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
  return clause.Expr{
    SQL:  "ST_PointFromText(?)",
    Vars: []interface{}{fmt.Sprintf("POINT(%d %d)", loc.X, loc.Y)},
  }
}

type User struct {
  Name     string
  Location Location
}

db.Create(&User{
  Name:     "jinzhu",
  Location: Location{X: 100, Y: 100},
})
// INSERT INTO `users` (`name`,`location`) VALUES ("jinzhu",ST_PointFromText("POINT(100 100)"))
```

# 默认值
这些默认值会被当作结构体字段的零值插入到数据库中
```go
type User struct {
  ID   int64
  Name string `gorm:"default:galeone"`
  Age  int64  `gorm:"default:18"`
}
```

注意，当结构体的字段默认值是零值的时候比如 0, '', false，这些字段值将不会被保存到数据库中，你可以使用指针类型或者Scanner/Valuer来避免这种情况。

```go
type User struct {
  gorm.Model
  Name string
  // Age 字段使用指针类型（*int）而非普通 int 类型，主要是为了解决 默认值与零值冲突 的问题，同时实现对字段的 空值（NULL）支持。
  Age  *int           `gorm:"default:18"`
  // sql.NullBool 通过 Valid 字段区分 false（有效值）和 NULL（无效值），与指针类型的设计目的一致
  Active sql.NullBool `gorm:"default:true"` // 默认值为 true，支持 NULL
}
```

```go
// 区分「未赋值」与「显式赋零值」

// 问题背景：
//  当字段定义为非指针类型（如 int）时，GORM 无法区分以下两种情况：
//  - 用户未显式赋值（希望使用默认值 18）。
//  - 用户显式赋值为 0（Go 中 int 的零值）。
// 指针类型的作用：
//  - 指针的默认值为 nil，表示「未赋值」，此时 GORM 会应用 default:18。
//  - 若显式赋值为 0（如 Age: new(int)），指针指向 0，GORM 会插入 0 而非默认值。

// 非指针类型（问题场景）
type User1 struct {
    Age int `gorm:"default:18"` // 无法区分未赋值和赋零值
}
u1 := User1{Age: 0} // 实际会插入 0，但无法确定用户意图是「用默认值」还是「显式设 0」

// 指针类型（解决方案）
type User2 struct {
    Age *int `gorm:"default:18"`
}
u2 := User2{Age: nil}   // 插入 18（默认值）
u3 := User2{Age: &[]int{0}[0]} // 插入 0
```

```go
// 若数据库字段有默认值（如 default:true），当 Valid=false 时，GORM 会插入 NULL 而非默认值
// sql.NullBool 默认不会自动处理 JSON 的 null，需自定义序列化逻辑或使用第三方库（如 guregu/null）
// Go 的 bool 零值是 false，而 sql.NullBool 的零值是 {false, false}（表示 NULL）
// 若需要更灵活的 NULL 值处理（如 JSON 兼容），可使用第三方库如 guregu/null

// 情况1：显式设置为 true（非 NULL）
user1 := User{
    Name:   "Alice",
    Active: sql.NullBool{Bool: true, Valid: true},
}

// 情况2：显式设置为 false（非 NULL）
user2 := User{
    Name:   "Bob",
    Active: sql.NullBool{Bool: false, Valid: true},
}

// 情况3：不设置值（数据库存 NULL，使用默认值 true）
user3 := User{
    Name: "Charlie",
    // Active 字段未赋值，Valid 默认为 false，数据库存 NULL
}
db.Create(&user1)
db.Create(&user2)
db.Create(&user3)

var user User
db.First(&user, 1) // 查询 ID=1 的用户

if user.Active.Valid {
    fmt.Printf("User's active status: %v\n", user.Active.Bool)
} else {
    fmt.Println("Active status is NULL (unknown)")
}

// 将 Active 更新为 NULL
db.Model(&user).Update("Active", sql.NullBool{Valid: false})

// 将 Active 更新为 false
db.Model(&user).Update("Active", sql.NullBool{Bool: false, Valid: true})
```

注意，若要让字段在数据库中拥有默认值则必须使用defaultTag来为结构体字段设置默认值。如果想要在数据库迁移的时候跳过默认值，可以使用 default:(-)
```go
type User struct {
  ID        string `gorm:"default:uuid_generate_v3()"` // db func
  FirstName string
  LastName  string
  Age       uint8
  FullName  string `gorm:"->;type:GENERATED ALWAYS AS (concat(firstname,' ',lastname));default:(-);"`
}
```

# Upsert 及冲突
GORM为不同数据库提供了对Upsert的兼容性支持。
```go
import "gorm.io/gorm/clause"

// Do nothing on conflict
db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)

// Update columns to default value on `id` conflict
db.Clauses(clause.OnConflict{
  Columns:   []clause.Column{{Name: "id"}},
  DoUpdates: clause.Assignments(map[string]interface{}{"role": "user"}),
}).Create(&users)
// MERGE INTO "users" USING *** WHEN NOT MATCHED THEN INSERT *** WHEN MATCHED THEN UPDATE SET ***; SQL Server
// INSERT INTO `users` *** ON DUPLICATE KEY UPDATE ***; MySQL

// Use SQL expression
db.Clauses(clause.OnConflict{
  Columns:   []clause.Column{{Name: "id"}},
  DoUpdates: clause.Assignments(map[string]interface{}{"count": gorm.Expr("GREATEST(count, VALUES(count))")}),
}).Create(&users)
// INSERT INTO `users` *** ON DUPLICATE KEY UPDATE `count`=GREATEST(count, VALUES(count));

// Update columns to new value on `id` conflict
db.Clauses(clause.OnConflict{
  Columns:   []clause.Column{{Name: "id"}},
  DoUpdates: clause.AssignmentColumns([]string{"name", "age"}),
}).Create(&users)
// MERGE INTO "users" USING *** WHEN NOT MATCHED THEN INSERT *** WHEN MATCHED THEN UPDATE SET "name"="excluded"."name"; SQL Server
// INSERT INTO "users" *** ON CONFLICT ("id") DO UPDATE SET "name"="excluded"."name", "age"="excluded"."age"; PostgreSQL
// INSERT INTO `users` *** ON DUPLICATE KEY UPDATE `name`=VALUES(name),`age`=VALUES(age); MySQL

// Update all columns to new value on conflict except primary keys and those columns having default values from sql func
db.Clauses(clause.OnConflict{
  UpdateAll: true,
}).Create(&users)
// INSERT INTO "users" *** ON CONFLICT ("id") DO UPDATE SET "name"="excluded"."name", "age"="excluded"."age", ...;
// INSERT INTO `users` *** ON DUPLICATE KEY UPDATE `name`=VALUES(name),`age`=VALUES(age), ...; MySQL
```

# IGNORE
```go
db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(&user)
// INSERT IGNORE INTO users (name,age...) VALUES ("jinzhu",18...);
```
