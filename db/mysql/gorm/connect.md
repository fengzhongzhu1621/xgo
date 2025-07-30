# 连接配置

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
* charset
* parseTime
* loc


* DefaultStringSize
* DisableDatetimePrecision
* DontSupportRenameIndex
* DontSupportRenameColumn
* SkipInitializeWithVersion

GORM 使用 database/sql 来维护连接池

```go
sqlDB, err := db.DB()

// SetMaxIdleConns 设置空闲连接池中连接的最大数量。
sqlDB.SetMaxIgleConns(10)
// SetMaxOpenConns 设置打开数据库连接的最大数量。
sqlDB.SetMaxOpenConns(100)
// SetConnMaxLifetime 设置了可以重新使用连接的最大时间。
sqlDB.SetConnMaxLifetime(time.Hour)
```
