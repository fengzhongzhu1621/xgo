# 1， 简介
pressly/goose 是一个用于 Go 语言的数据库迁移工具。它提供了一种简单且高效的方式来管理数据库 schema 的版本控制，
适用于数据库结构变更的管理和迁移。goose 的设计目标是让数据库迁移更加简便、安全，并且在多环境下可以轻松应用。

* 通过维护迁移文件（通常是 SQL 脚本或 Go 语言代码）来管理数据库的版本。这些迁移文件记录了数据库结构变更的具体步骤。
* 每个迁移文件都有一个版本号，goose 通过比对当前版本和目标版本来决定哪些迁移需要执行。
* 支持通过 SQL 文件进行数据库迁移，也支持通过 Go 代码进行迁移。
* 允许开发者回滚迁移操作，恢复到指定版本。

迁移脚本的命名规范
goose 默认会根据时间戳来为迁移文件命名，命名规则是：<timestamp>_<migration_name>.sql。其中 timestamp 是一个 13 位的时间戳（毫秒级），可以确保每个迁移文件都有唯一的名称。


# 2. 安装
```
go get github.com/pressly/goose/v3
go install github.com/pressly/goose/v3/cmd/goose@latest
```

# 3. 创建迁移
## 3.1 创建sql
```shell
goose -dir ./migrations create add_users_table sql
```

## 3.3 创建go文件
```shell
goose -dir ./migrations create rename_root go
goose -dir ./migrations create add_user_not_tx go
```

# 4. 执行迁移
## 4.1 运行迁移
up 命令会执行 migrations 目录中的所有待执行的迁移。

```sh
goose -dir ./migrations mysql "root:@tcp(127.0.0.1:3306)/xgo?parseTime=True&loc=Local" up
```

# 4.2 回滚迁移

回滚最近的迁移
```sh
goose -dir ./migrations mysql "root:@tcp(127.0.0.1:3306)/xgo?parseTime=True&loc=Local" down
```

回滚到倒数第二个版本
```sh
goose -dir ./migrations mysql "root:@tcp(127.0.0.1:3306)/xgo?parseTime=True&loc=Local" down 2
```

## 4.3 查看迁移历史
```sh
goose -dir ./migrations mysql "root:@tcp(127.0.0.1:3306)/xgo?parseTime=True&loc=Local" status
```


# 版本控制
goose 会在数据库中创建一个版本表（通常名为 goose_db_version），它存储了数据库的当前版本。
当你运行迁移时，goose 会检查此表以确定哪些迁移已经执行过，并按顺序执行未执行的迁移。

# 环境变量
```sh
export GOOSE_DRIVER=DRIVER
export GOOSE_DBSTRING=DBSTRING
export GOOSE_MIGRATION_DIR=MIGRATION_DIR
export GOOSE_TABLE=TABLENAME
```

```sh
//Via .env files with corresponding variables. .env file example:

GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://admin:admin@localhost:5432/admin_db
GOOSE_MIGRATION_DIR=./migrations
GOOSE_TABLE=custom.goose_migrations
```
