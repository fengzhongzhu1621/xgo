# USE INDEX
```go
import "gorm.io/hints"

// 对指定索引提供建议
db.Clauses(hints.UseIndex("idx_user_name")).Find(&User{})
// SQL: SELECT * FROM `users` USE INDEX (`idx_user_name`)

// 强制对JOIN操作使用某些索引
db.Clauses(hints.ForceIndex("idx_user_name", "idx_user_id").ForJoin()).Find(&User{})
// SQL: SELECT * FROM `users` FORCE INDEX FOR JOIN (`idx_user_name`,`idx_user_id`)
```
