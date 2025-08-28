# 简介
Gormigrate 专为 GORM（Go 流行 ORM）设计。如果项目已使用 GORM 进行数据库操作，Gormigrate 是天然之选，可在模型旁使用 Go 代码定义迁移。

* GORM 集成：利用 GORM 的 ORM 能力执行迁移；
* 编程式迁移：使用 Go 代码，而非 SQL；回滚支持：
* 内置回滚函数，轻松撤销迁移。

需要为每个迁移定义一个  Migrate  和一个  Rollback  函数。这使得你可以在迁移过程中使用 GORM 的功能，如自动迁移、添加外键等。

本身不提供命令行工具，但是你可以在你的 Go 项目中创建一个简单的命令行工具，以便在控制台中执行迁移和回滚操作。

# 示例
https://dev.to/kengowada/gorm-and-goose-migrations-1ec 
https://github.com/1Panel-dev/1Panel/blob/dev-v2/core/init/migration/migrate.go 
https://github.com/go-gormigrate/gormigrate/
