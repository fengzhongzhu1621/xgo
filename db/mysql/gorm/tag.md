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

## primarykey
```go
ID        uint `gorm:"primarykey"`
```

## index
```go
DeletedAt DeletedAt `gorm:"index"`
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

## softDelete
```go
DeletedAt soft_delete.DeletedAt `gorm:"softDelete:milli"`
DeletedAt soft_delete.DeletedAt `gorm:"softDelete:nano"`
IsDel soft_delete.DeletedAt `gorm:"softDelete:flag"`
IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"` // use `1` `0`
IsDel     soft_delete.DeletedAt `gorm:"softDelete:,DeletedAtField:DeletedAt"` // use `unix second`
IsDel     soft_delete.DeletedAt `gorm:"softDelete:nano,DeletedAtField:DeletedAt"` // use `unix nano second`
```

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
