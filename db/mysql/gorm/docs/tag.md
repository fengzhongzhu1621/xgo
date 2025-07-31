# 约定

```go
type User struct {
  ID           uint           // Standard field for the primary key
  Name         string         // A regular string field
  Email        *string        // A pointer to a string, allowing for null values
  Age          uint8          // An unsigned 8-bit integer
  Birthday     *time.Time     // A pointer to time.Time, can be null
  MemberNumber sql.NullString // Uses sql.NullString to handle nullable strings
  ActivatedAt  sql.NullTime   // Uses sql.NullTime for nullable time fields
  CreatedAt    time.Time      // Automatically managed by GORM for creation time
  UpdatedAt    time.Time      // Automatically managed by GORM for update time
  ignored      string         // fields that aren't exported are ignored
}
```
* 具体数字类型如 uint、string和 uint8 直接使用。
* 指向 *string 和 *time.Time 类型的指针表示可空字段。
* 来自 database/sql 包的 sql.NullString 和 sql.NullTime 用于具有更多控制的可空字段。
* CreatedAt 和 UpdatedAt 是特殊字段，当记录被创建或更新时，GORM 会自动向内填充当前时间。


* 主键：GORM 使用一个名为ID 的字段作为每个模型的默认主键。
* 表名：默认情况下，GORM 将结构体名称转换为 snake_case 并为表名加上复数形式。 For instance, a User struct becomes users in the database, and a GormUserName becomes gorm_user_names.
* 列名：GORM 自动将结构体字段名称转换为 snake_case 作为数据库中的列名。
* 时间戳字段：GORM使用字段 CreatedAt 和 UpdatedAt 来自动跟踪记录的创建和更新时间。

# gorm.Model

预定义的结构体，名为gorm.Model，其中包含常用字段

```go
// gorm.Model 的定义
type Model struct {
  ID        uint           `gorm:"primaryKey"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt gorm.DeletedAt `gorm:"index"`
}
```
* 每个记录的唯一标识符（主键）
* 在创建记录时自动设置为当前时间。
* 每当记录更新时，自动更新为当前时间。
* 用于软删除（将记录标记为已删除，而实际上并未从数据库中删除）。

# 主键
```go
type User struct {
  ID   string // 默认情况下，名为 `ID` 的字段会作为表的主键
  Name string
}
```

```go
// 可以通过标签 primaryKey 将其它字段设为主键
// 将 `UUID` 设为主键
type Animal struct {
  ID     int64
  UUID   string `gorm:"primaryKey"`
  Name   string
  Age    int64
}
```

默认情况下，整型 PrioritizedPrimaryField 启用了 AutoIncrement，要禁用它，您需要为整型字段关闭 autoIncrement
```go
type Product struct {
  CategoryID uint64 `gorm:"primaryKey;autoIncrement:false"`
  TypeID     uint64 `gorm:"primaryKey;autoIncrement:false"`
}
```

# 命名策略
```go
type User struct {
  ID        uint      // 列名是 `id`
  Name      string    // 列名是 `name`
  Birthday  time.Time // 列名是 `birthday`
  CreatedAt time.Time // 列名是 `created_at`
}
```

```go
// 使用 column 标签或 命名策略 来覆盖列名
type Animal struct {
  AnimalID int64     `gorm:"column:beast_id"`         // 将列名设为 `beast_id`
  Birthday time.Time `gorm:"column:day_of_the_beast"` // 将列名设为 `day_of_the_beast`
  Age      int64     `gorm:"column:age_of_the_beast"` // 将列名设为 `age_of_the_beast`
}
```

# 权限控制

使用 GORM Migrator 创建表时，不会创建被忽略的字段

```go
type User struct {
  Name string `gorm:"<-:create"` // 允许读和创建
  Name string `gorm:"<-:update"` // 允许读和更新
  Name string `gorm:"<-"`        // 允许读和写（创建和更新）
  Name string `gorm:"<-:false"`  // 允许读，禁止写
  Name string `gorm:"->"`        // 只读（除非有自定义配置，否则禁止写）
  Name string `gorm:"->;<-:create"` // 允许读和写
  Name string `gorm:"->:false;<-:create"` // 仅创建（禁止从 db 读）
  Name string `gorm:"-"`  // 通过 struct 读写会忽略该字段
  Name string `gorm:"-:all"`        // 通过 struct 读写、迁移会忽略该字段
  Name string `gorm:"-:migration"`  // 通过 struct 迁移会忽略该字段
}
```

# 时间
GORM 约定使用 CreatedAt、UpdatedAt 追踪创建/更新时间。如果您定义了这种字段，GORM 在创建、更新时会自动填充 当前时间

* 要使用不同名称的字段，您可以配置 autoCreateTime、autoUpdateTime 标签。
* 如果您想要保存 UNIX（毫/纳）秒时间戳，而不是 time，您只需简单地将 time.Time 修改为 int 即可

```
type User struct {
  CreatedAt time.Time // 在创建时，如果该字段值为零值，则使用当前时间填充
  UpdatedAt int       // 在创建时该字段值为零值或者在更新时，使用当前时间戳秒数填充
  Updated   int64 `gorm:"autoUpdateTime:nano"` // 使用时间戳纳秒数填充更新时间
  Updated   int64 `gorm:"autoUpdateTime:milli"` // 使用时间戳毫秒数填充更新时间
  Created   int64 `gorm:"autoCreateTime"`      // 使用时间戳秒数填充创建时间
}
```

# 复数表名
GORM 使用结构体名的 蛇形命名 作为表名。对于结构体 User，根据约定，其表名为 users

注意： TableName 不支持动态变化，它会被缓存下来以便后续使用。想要使用动态表名，你可以使用 Scopes

```go
type Tabler interface {
    TableName() string
}

// TableName 会将 User 的表名重写为 `profiles`
func (User) TableName() string {
  return "profiles"
}
```

# CreatedAt
对于有 CreatedAt 字段的模型，创建记录时，如果该字段值为零值，则将该字段的值设为当前时间

```go
db.Create(&user) // 将 `CreatedAt` 设为当前时间

user2 := User{Name: "jinzhu", CreatedAt: time.Now()}
db.Create(&user2) // user2 的 `CreatedAt` 不会被修改
```

```go
// 想要修改该值，您可以使用 `Update`
db.Model(&user).Update("CreatedAt", time.Now())
你可以通过将 autoCreateTime 标签置为 false 来禁用时间戳追踪，例如：

type User struct {
  CreatedAt time.Time `gorm:"autoCreateTime:false"`
}
```

# UpdatedAt
对于有 UpdatedAt 字段的模型，更新记录时，将该字段的值设为当前时间。创建记录时，如果该字段值为零值，则将该字段的值设为当前时间

```go
db.Save(&user) // 将 `UpdatedAt` 设为当前时间

db.Model(&user).Update("name", "jinzhu") // 会将 `UpdatedAt` 设为当前时间

db.Model(&user).UpdateColumn("name", "jinzhu") // `UpdatedAt` 不会被修改

user2 := User{Name: "jinzhu", UpdatedAt: time.Now()}
db.Create(&user2) // 创建记录时，user2 的 `UpdatedAt` 不会被修改

user3 := User{Name: "jinzhu", UpdatedAt: time.Now()}
db.Save(&user3) // 更新时，user3 的 `UpdatedAt` 会修改为当前时间
```

```go
可以通过将 autoUpdateTime 标签置为 false 来禁用时间戳追踪，例如：

type User struct {
  UpdatedAt time.Time `gorm:"autoUpdateTime:false"`
}
```

# 嵌入结构体
* 对于匿名字段，GORM 会将其字段包含在父结构体中
* 对于正常的结构体字段，可以通过标签 embedded 将其嵌入

```go
type Author struct {
    Name  string
    Email string
}

type Blog struct {
  ID      int
  Author  Author `gorm:"embedded"`
  Upvotes int32
}
// 等效于
type Blog struct {
  ID    int64
  Name  string
  Email string
  Upvotes  int32
}
```

```go
type Author struct {
    Name  string
    Email string
}
type Blog struct {
  ID      int
  Author  Author `gorm:"embedded;embeddedPrefix:author_"`
  Upvotes int32
}
// 等效于
type Blog struct {
  ID          int64
  AuthorName string
  AuthorEmail string
  Upvotes     int32
}
```

# 标签

## 字段名
```go
AnimalID int64     `gorm:"column:beast_id"`         // 将列名设为 `beast_id`
```

## primarykey
```go
ID        uint `gorm:"primarykey"`
```

## index
```go
DeletedAt DeletedAt `gorm:"index"`
Name3     string    `gorm:"index:,sort:desc,collate:utf8,type:btree,length:10"`
```

```go
type User struct {
    Name  string `gorm:"index"`
    Name2 string `gorm:"index:idx_name,unique"`
    Name3 string `gorm:"index:,sort:desc,collate:utf8,type:btree,length:10,where:name3 != 'jinzhu'"`
    Name4 string `gorm:"uniqueIndex"`
    Age   int64  `gorm:"index:,class:FULLTEXT,comment:hello \\, world,where:age > 10"`
    Age2  int64  `gorm:"index:,expression:ABS(age)"`
}

// MySQL 选项
type User struct {
    Name string `gorm:"index:,class:FULLTEXT,option:WITH PARSER ngram INVISIBLE"`
}

// PostgreSQL 选项
type User struct {
    Name string `gorm:"index:,option:CONCURRENTLY"`
}
```

uniqueIndex 标签的作用与 index 类似，它等效于 index:,unique
```go
Refer    uint      `gorm:"index:,unique"`
Name     string    `gorm:"unique;not null;size:32"`
Name2     string   `gorm:"index:idx_name,unique"`

type User struct {
    Name1 string `gorm:"uniqueIndex"`
    Name2 string `gorm:"uniqueIndex:idx_name,sort:desc"`
}
```

两个字段使用同一个索引名将创建复合索引
```go
// create composite index `idx_member` with columns `name`, `number`
type User struct {
    Name   string `gorm:"index:idx_member"`
    Number string `gorm:"index:idx_member"`
}
```

使用 priority 指定顺序，默认优先级值是 10，如果优先级值相同，则顺序取决于模型结构体字段的顺序
```go
type User struct {
    Name   string `gorm:"index:idx_member"`
    Number string `gorm:"index:idx_member"`
}
// column order: name, number

type User struct {
    Name   string `gorm:"index:idx_member,priority:2"`
    Number string `gorm:"index:idx_member,priority:1"`
}
// column order: number, name

type User struct {
    Name   string `gorm:"index:idx_member,priority:12"`
    Number string `gorm:"index:idx_member"`
}
// column order: number, name
```

一个字段接受多个 index、uniqueIndex 标签，这会在一个字段上创建多个索引
```go
type UserIndex struct {
    OID          int64  `gorm:"index:idx_id;index:idx_oid,unique"`
    MemberNumber string `gorm:"index:idx_id"`
}
```

共享复合索引，通过命名策略生成索引名称
```go
type Foo struct {
  IndexA int `gorm:"index:,unique,composite:myname"`
  IndexB int `gorm:"index:,unique,composite:myname"`
}

type Bar0 struct {
  Foo
}

type Bar1 struct {
  Foo
}
// 复合索引的名称分别是 idx_bar0_myname 和 idx_bar1_myname。
// 复合 只能在指定索引名称时使用。
```

## size
```go
Email        string  `gorm:"default:'xx@x.com';size:32"`
BizType     string `gorm:"size:32"`
Env         string `gorm:"size:32;index"`
```

## type
```go
Text string `gorm:"type:text"`
CheckLine string `gorm:"type:longtext"` // 映射为 LONGTEXT
Content string `gorm:"type:text;not null;default:'待编辑'"`
Sno       int64          `gorm:"column:sno;type:int;primaryKey" json:"sno"`
Score     float32        `gorm:"column:score;type:float;not null" json:"score"`
Cname     string         `gorm:"column:cname;type:varchar(255);not null" json:"cname"`
CreatedAt time.Time      `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
```

## comment
```go
Age          uint8   `gorm:"default:0;comment:'user age'"`
```

## default
```go
Name string `gorm:"default:galeone"`
Age  int64  `gorm:"default:18"`
// Age 字段使用指针类型（*int）而非普通 int 类型，主要是为了解决 默认值与零值冲突 的问题，同时实现对字段的 空值（NULL）支持。
// 当字段定义为非指针类型（如 int）时，GORM 无法区分以下两种情况：
//  * 用户未显式赋值（希望使用默认值 18）。
//  * 用户显式赋值为 0（Go 中 int 的零值）

Age  *int           `gorm:"default:18"`
Active sql.NullBool `gorm:"default:true"`
ID        string `gorm:"default:uuid_generate_v3()"` // db func
FullName  string `gorm:"->;type:GENERATED ALWAYS AS (concat(firstname,' ',lastname));default:(-);"`
```

## embedded
```go
Author  Author `gorm:"embedded"`
Author  Author `gorm:"embedded;embeddedPrefix:author_"`
```

## 序列化
```go
Name        []byte                 `gorm:"serializer:json"`
JobInfo     Job                    `gorm:"type:bytes;serializer:gob"`
CreatedTime int64                  `gorm:"serializer:unixtime;type:time"` // 将 int 作为日期时间存储到数据库中
```

## softDelete
```go
DeletedAt soft_delete.DeletedAt `gorm:"softDelete:milli"`
DeletedAt soft_delete.DeletedAt `gorm:"softDelete:nano"`
IsDel soft_delete.DeletedAt `gorm:"softDelete:flag"`
IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"` // use `1` `0`
IsDel     soft_delete.DeletedAt `gorm:"softDelete:,DeletedAtField:DeletedAt"` // use `unix second`
IsDel     soft_delete.DeletedAt `gorm:"softDelete:nano,DeletedAtField:DeletedAt"` // use `unix nano second`
```

# 外键
```go
Company   Company `gorm:"foreignKey:CompanyRefer"`
Company   Company `gorm:"references:Code"` // 使用 Code 作为引用
Company   Company `gorm:"references:CompanyID"` // 使用 Company.CompanyID 作为引用
Company   Company `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
CreditCard CreditCard `gorm:"foreignKey:UserName;references:Name"`
Languages []Language `gorm:"many2many:user_languages;"`
LocaleTags []Tag `gorm:"many2many:locale_blog_tags;ForeignKey:id,locale;References:id"`
SharedTags []Tag `gorm:"many2many:shared_blog_tags;ForeignKey:id;References:id"`
Toys []Toy `gorm:"polymorphic:Owner;"`
Toys []Toy `gorm:"polymorphic:Owner;"`
Toys []Toy `gorm:"polymorphicType:Kind;polymorphicId:OwnerID;polymorphicValue:master"`
```


# mysql 表字段映射
```go
Id   int    ->  `id` bigint NOT NULL AUTO_INCREMENT,
ID   uint   -> `id` bigint unsigned NOT NULL AUTO_INCREMENT
uint        -> bigint unsigned DEFAULT NULL
int         -> bigint DEFAULT NULL
string      -> longtext
float64     -> double DEFAULT NULL
optimisticlock.Version -> bigint DEFAULT NULL
```


# 标签参数
| 标签名                  | 说明                                                                                                                                                                                                                                                                                                                                 |
|-------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `column`                | 指定 db 列名                                                                                                                                                                                                                                                                                                                        |
| `type`                  | 列数据类型，推荐使用兼容性好的通用类型，例如：所有数据库都支持 `bool`、`int`、`uint`、`float`、`string`、`time`、`bytes` 并且可以和其他标签一起使用，例如：`not null`、`size`、`autoIncrement`。像 `varbinary(8)` 这样指定数据库数据类型也是支持的。在使用指定数据库数据类型时，它需要是完整的数据库数据类型，如：`MEDIUMINT UNSIGNED not NULL AUTO_INCREMENT` |
| `serializer`            | 指定将数据序列化或反序列化到数据库中的序列化器，例如：`serializer:json/gob/unixtime`                                                                                                                                                                                                                                               |
| `size`                  | 定义列数据类型的大小或长度，例如 `size: 256`                                                                                                                                                                                                                                                                                        |
| `primaryKey`            | 将列定义为主键                                                                                                                                                                                                                                                                                                                     |
| `unique`                | 将列定义为唯一键                                                                                                                                                                                                                                                                                                                   |
| `default`               | 定义列的默认值                                                                                                                                                                                                                                                                                                                     |
| `precision`             | 指定列的精度                                                                                                                                                                                                                                                                                                                       |
| `scale`                 | 指定列大小                                                                                                                                                                                                                                                                                                                         |
| `not null`              | 指定列为 `NOT NULL`                                                                                                                                                                                                                                                                                                                |
| `autoIncrement`         | 指定列为自动增长                                                                                                                                                                                                                                                                                                                   |
| `autoIncrementIncrement`| 自动步长，控制连续记录之间的间隔                                                                                                                                                                                                                                                                                                   |
| `embedded`              | 嵌套字段                                                                                                                                                                                                                                                                                                                           |
| `embeddedPrefix`        | 嵌入字段的列名前缀                                                                                                                                                                                                                                                                                                                 |
| `autoCreateTime`        | 创建时追踪当前时间，对于 `int` 字段，它会追踪时间戳秒数，您可以使用 `nano/milli` 来追踪纳秒、毫秒时间戳，例如：`autoCreateTime:nano`                                                                                                                                                                                              |
| `autoUpdateTime`        | 创建/更新时追踪当前时间，对于 `int` 字段，它会追踪时间戳秒数，您可以使用 `nano/milli` 来追踪纳秒、毫秒时间戳，例如：`autoUpdateTime:milli`                                                                                                                                                                                        |
| `index`                 | 根据参数创建索引，多个字段使用相同的名称则创建复合索引                                                                                                                                                                                                                                                                             |
| `uniqueIndex`           | 与 `index` 相同，但创建的是唯一索引                                                                                                                                                                                                                                                                                                |
| `check`                 | 创建检查约束，例如 `check:age > 13`                                                                                                                                                                                                                                                                                                |
| `<-`                    | 设置字段写入的权限，`<-:create` 只创建、`<-:update` 只更新、`<-:false` 无写入权限、`<-` 创建和更新权限                                                                                                                                                                                                                              |
| `->`                    | 设置字段读的权限，`->:false` 无读权限                                                                                                                                                                                                                                                                                              |
| `-`                     | 忽略该字段，`-` 表示无读写，`-:migration` 表示无迁移权限，`-:all` 表示无读写迁移权限                                                                                                                                                                                                                                                |
| `comment`               | 迁移时为字段添加注释                                                                                                                                                                                                                                                                                                               |
