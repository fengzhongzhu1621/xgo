# DSN

https://github.com/go-sql-driver/mysql#dsn-data-source-name

```go
db, err := gorm.Open(mysql.New(mysql.Config{
  DSN: "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local", // DSN data source name
  DefaultStringSize: 256, // string 类型字段的默认长度
  DisableDatetimePrecision: true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
  DontSupportRenameIndex: true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
  DontSupportRenameColumn: true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
  SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
}), &gorm.Config{})
```

GORM默认使用UTC时区，需要通过在DSN中设置parseTime=True&loc=Local或parseTime=True&loc=Asia%2FShanghai

* charset: charset=utf8mb4
* parseTime: parseTime=True
* loc: loc=Local 访问磁盘中/etc/localtime文件，从而知道当前系统是UTC+N


* DefaultStringSize
* DisableDatetimePrecision
* DontSupportRenameIndex
* DontSupportRenameColumn
* SkipInitializeWithVersion

# *sql.DB
```go
// 获取通用数据库对象 sql.DB，然后使用其提供的功能
sqlDB, err := db.DB()

// Ping
sqlDB.Ping()

// Close
sqlDB.Close()

// 返回数据库统计信息
sqlDB.Stats()
```

```go
// 连接池 GORM 使用 database/sql 来维护连接池
sqlDB, err := db.DB()

// SetMaxIdleConns 设置空闲连接池中连接的最大数量。
sqlDB.SetMaxIgleConns(10)
// SetMaxOpenConns 设置打开数据库连接的最大数量。
sqlDB.SetMaxOpenConns(100)
// SetConnMaxLifetime 设置了可以重新使用连接的最大时间。
sqlDB.SetConnMaxLifetime(time.Hour)
```

# QueryFields
```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
  QueryFields: true,
})

// 当 QueryFields 被设置为 true 时，此行为默认进行
db.Find(&user)
// SQL: SELECT `users`.`name`, `users`.`age`, ... FROM `users`

// 开启 QueryFields 并使用会话模式（Session mode）
db.Session(&gorm.Session{QueryFields: true}).Find(&user)
// SQL: SELECT `users`.`name`, `users`.`age`, ... FROM `users`

```

# TranslateError 方言转换错误
当启用TranslateError时，GORM可以返回与所使用的数据库方言相关的特定错误，GORM将数据库特有的错误转换为其自己的通用错误。

```go
db, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{TranslateError: true})
```

# Logger
```go
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	Logger: logger.Default.LogMode(logger.Info), // 输出详细日志
})
```
