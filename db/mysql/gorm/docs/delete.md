# 删除一条记录
```go
db.Delete(&email)
// DELETE from emails where id = 10;
```

# 根据主键删除记录
```go
db.Delete(&User{}, 10)
// DELETE FROM users WHERE id = 10;

db.Delete(&User{}, "10")
// DELETE FROM users WHERE id = 10;

db.Delete(&users, []int{1,2,3})
// DELETE FROM users WHERE id IN (1,2,3);

var users = []User{{ID: 1}, {ID: 2}, {ID: 3}}
db.Delete(&users)
// DELETE FROM users WHERE id IN (1,2,3);

db.Delete(&users, "name LIKE ?", "%jinzhu%")
// DELETE FROM users WHERE name LIKE "%jinzhu%" AND id IN (1,2,3);
```

# 批量删除
```go
// Delete with additional conditions
db.Where("name = ?", "jinzhu").Delete(&email)
// DELETE from emails where id = 10 AND name = "jinzhu";
//
db.Where("email LIKE ?", "%jinzhu%").Delete(&Email{})
// DELETE from emails where email LIKE "%jinzhu%";
//
db.Delete(&Email{}, "email LIKE ?", "%jinzhu%")
// DELETE from emails where email LIKE "%jinzhu%";
```

# Delete Hooks
```go
func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
  if u.Role == "admin" {
    return errors.New("admin user not allowed to delete")
  }
  return
}
```

# Block Global Delete
```go
// If you perform a batch delete without any conditions, GORM WON’T run it, and will return ErrMissingWhereClause error
db.Delete(&User{}).Error // gorm.ErrMissingWhereClause

db.Delete(&[]User{{Name: "jinzhu1"}, {Name: "jinzhu2"}}).Error // gorm.ErrMissingWhereClause

db.Where("1 = 1").Delete(&User{})
// DELETE FROM `users` WHERE 1=1

db.Exec("DELETE FROM users")
// DELETE FROM users

db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&User{})
// DELETE FROM users
```

# Soft Delete
## gorm.DeletedAt
```go
// By default, gorm.Model uses *time.Time as the value for the DeletedAt field,
//  and it provides other data formats support with plugin gorm.io/plugin/soft_delete
type User struct {
  ID      int
  Deleted gorm.DeletedAt
  Name    string
}
```

```go
// If your model includes a gorm.DeletedAt field (which is included in gorm.Model), it will get soft delete ability automatically!
// When calling Delete, the record WON’T be removed from the database, but GORM will set the DeletedAt‘s value to the current time, and the data is not findable with normal Query methods anymore.
//
// user's ID is `111`
db.Delete(&user)
// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE id = 111;

// Batch Delete
db.Where("age = ?", 20).Delete(&User{})
// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE age = 20;

// Soft deleted records will be ignored when querying
db.Where("age = 20").Find(&user)
// SELECT * FROM users WHERE age = 20 AND deleted_at IS NULL;
```

## soft_delete.DeletedAt
```go
// Use unix second as delete flag
import "gorm.io/plugin/soft_delete"

type User struct {
  ID        uint
  Name      string
  DeletedAt soft_delete.DeletedAt
}

type User struct {
  ID    uint
  Name  string
  DeletedAt soft_delete.DeletedAt `gorm:"softDelete:milli"`
  // DeletedAt soft_delete.DeletedAt `gorm:"softDelete:nano"`
}

// Use 1 / 0 AS Delete Flag
type User struct {
  ID    uint
  Name  string
  IsDel soft_delete.DeletedAt `gorm:"softDelete:flag"`
}

// Mixed Mode
type User struct {
  ID        uint
  Name      string
  DeletedAt time.Time
  IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"` // use `1` `0`
  // IsDel     soft_delete.DeletedAt `gorm:"softDelete:,DeletedAtField:DeletedAt"` // use `unix second`
  // IsDel     soft_delete.DeletedAt `gorm:"softDelete:nano,DeletedAtField:DeletedAt"` // use `unix nano second`
}
```

# Find soft deleted records
```go
db.Unscoped().Where("age = 20").Find(&users)
// SELECT * FROM users WHERE age = 20;
```

# Delete permanently
```go
db.Unscoped().Delete(&order)
// DELETE FROM orders WHERE id=10;
```
