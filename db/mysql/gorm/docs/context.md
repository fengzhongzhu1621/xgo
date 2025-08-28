# 声明
```go
db.WithContext(ctx).Find(&users)
```

# 持续会话模式
持续会话模式非常适合执行一系列相关的操作。 它在这些操作之间保持上下文，对于事务等场景特别有用。

```go
tx := db.WithContext(ctx)
tx.First(&user, 1)
tx.Model(&user).Update("role", "admin")
```

# Context超时
```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

db.WithContext(ctx).Find(&users)
```

# Hooks/Callbacks 中的 Context
```go
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
  ctx := tx.Statement.Context
  // ...使用ctx做一些你想做的事
  return
}
```
