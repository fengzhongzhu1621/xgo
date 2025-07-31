# FOR UPDATE
在事务（transaction）中锁定选中行（selected rows）
当你准备在事务（transaction）中更新（update）一些行（rows）时，并且想要在本事务完成前，阻止（prevent）其他的事务（other transactions）修改你准备更新的选中行。
```go
// 基本的 FOR UPDATE 锁
db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&users)
// SQL: SELECT * FROM `users` FOR UPDATE
```

# FOR SHARE OF
只允许其他事务读取（read）被锁定的内容，而无法修改（update）或者删除（delete）。
Table选项用于指定将要被锁定的表。 这在你想要 join 多个表，并且锁定其一时非常有用。
```go
db.Clauses(clause.Locking{
  Strength: "SHARE",
  Table: clause.Table{Name: clause.CurrentTable},
}).Find(&users)
// SQL: SELECT * FROM `users` FOR SHARE OF `users`
```

# FOR UPDATE NOWAIT
尝试获取一个锁，如果锁不可用，导致了获取失败，函数将会立即返回一个error。 当一个事务等待其他事务释放它们的锁时，此Options（Nowait）可以阻止这种行为
```go
db.Clauses(clause.Locking{
  Strength: "UPDATE",
  Options: "NOWAIT",
}).Find(&users)
// SQL: SELECT * FROM `users` FOR UPDATE NOWAIT
```

# 乐观锁
1. 在 table 中增加一列，用于记录此行数据的版本号
2. 更新数据前，先读取当前数据行的版本号
3. 更新时，对 UPDATE 语句作两处调整
    * WHERE 语句中加入版本号的比较条件，确保只有当前版本号与数据库中的版本号一致时才执行更新
    ```
    WHERE ... and version = [current version]
    ```

    * UPDATE 语句中递增版本号以保证每次更新后版本号都会变化
    ```
    UPDATE set  ..., version = version + 1
    ```
4. SQL 执行以后需要检查更新行数是否为0，如果为0则说明有更新冲突，需要重试直到成功为止
