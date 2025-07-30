package gorm_2

import (
	"log"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// 删除一条记录
// Delete a Record
// Email's ID is `10`
// db.Delete(&email)
// DELETE from emails where id = 10;

// ////////////////////////////////////////////////////////////////////////////////////////////////////////
// 根据主键删除记录
// db.Delete(&User{}, 10)
// DELETE FROM users WHERE id = 10;
//
// db.Delete(&User{}, "10")
// DELETE FROM users WHERE id = 10;
//
// db.Delete(&users, []int{1,2,3})
// DELETE FROM users WHERE id IN (1,2,3);
//
// var users = []User{{ID: 1}, {ID: 2}, {ID: 3}}
// db.Delete(&users)
// DELETE FROM users WHERE id IN (1,2,3);
//
// db.Delete(&users, "name LIKE ?", "%jinzhu%")
// DELETE FROM users WHERE name LIKE "%jinzhu%" AND id IN (1,2,3);

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// 批量删除
// Delete with additional conditions
// db.Where("name = ?", "jinzhu").Delete(&email)
// DELETE from emails where id = 10 AND name = "jinzhu";
//
// db.Where("email LIKE ?", "%jinzhu%").Delete(&Email{})
// DELETE from emails where email LIKE "%jinzhu%";
//
// db.Delete(&Email{}, "email LIKE ?", "%jinzhu%")
// DELETE from emails where email LIKE "%jinzhu%";

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// Delete Hooks
// func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
//   if u.Role == "admin" {
//     return errors.New("admin user not allowed to delete")
//   }
//   return
// }

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// Block Global Delete
// If you perform a batch delete without any conditions, GORM WON’T run it, and will return ErrMissingWhereClause error
// db.Delete(&User{}).Error // gorm.ErrMissingWhereClause
//
// db.Delete(&[]User{{Name: "jinzhu1"}, {Name: "jinzhu2"}}).Error // gorm.ErrMissingWhereClause
//
// db.Where("1 = 1").Delete(&User{})
// DELETE FROM `users` WHERE 1=1
//
// db.Exec("DELETE FROM users")
// DELETE FROM users
//
// db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&User{})
// DELETE FROM users

// ////////////////////////////////////////////////////////////////////////////////////////////////////////
// Soft Delete
// If your model includes a gorm.DeletedAt field (which is included in gorm.Model), it will get soft delete ability automatically!
// When calling Delete, the record WON’T be removed from the database, but GORM will set the DeletedAt‘s value to the current time, and the data is not findable with normal Query methods anymore.
//
// user's ID is `111`
// db.Delete(&user)
// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE id = 111;

// Batch Delete
// db.Where("age = ?", 20).Delete(&User{})
// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE age = 20;

// Soft deleted records will be ignored when querying
// db.Where("age = 20").Find(&user)
// SELECT * FROM users WHERE age = 20 AND deleted_at IS NULL;

// ////////////////////////////////////////////////////////////////////////////////////////////////////////
// gorm.DeletedAt
// By default, gorm.Model uses *time.Time as the value for the DeletedAt field,
//  and it provides other data formats support with plugin gorm.io/plugin/soft_delete
// type User struct {
//   ID      int
//   Deleted gorm.DeletedAt
//   Name    string
// }
//
// Use unix second as delete flag
// import "gorm.io/plugin/soft_delete"

// type User struct {
//   ID        uint
//   Name      string
//   DeletedAt soft_delete.DeletedAt
// }
//
// type User struct {
//   ID    uint
//   Name  string
//   DeletedAt soft_delete.DeletedAt `gorm:"softDelete:milli"`
//   // DeletedAt soft_delete.DeletedAt `gorm:"softDelete:nano"`
// }
//
// Use 1 / 0 AS Delete Flag
// type User struct {
//   ID    uint
//   Name  string
//   IsDel soft_delete.DeletedAt `gorm:"softDelete:flag"`
// }
//
// Mixed Mode
// type User struct {
//   ID        uint
//   Name      string
//   DeletedAt time.Time
//   IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"` // use `1` `0`
//   // IsDel     soft_delete.DeletedAt `gorm:"softDelete:,DeletedAtField:DeletedAt"` // use `unix second`
//   // IsDel     soft_delete.DeletedAt `gorm:"softDelete:nano,DeletedAtField:DeletedAt"` // use `unix nano second`
// }

// ////////////////////////////////////////////////////////////////////////////////////////////////////////
// Find soft deleted records
// db.Unscoped().Where("age = 20").Find(&users)
// SELECT * FROM users WHERE age = 20;

// ////////////////////////////////////////////////////////////////////////////////////////////////////////
// Delete permanently
// db.Unscoped().Delete(&order)
// DELETE FROM orders WHERE id=10;

func deleteStudent(db *gorm.DB, id int) {
	// 根据主键删除
	result := db.Delete(&Student{}, id)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	// 根据记录对象删除
	var student Student
	db.First(&student, id)
	result = db.Delete(&student)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
}

func TestDelete(t *testing.T) {
	dsn := "root:@tcp(127.0.0.1:3306)/xgo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: logger.Default.LogMode(logger.Info), // 输出详细日志
	})
	if err != nil {
		log.Fatal(err)
	}

	var stu Student
	stu.Name = "bob"
	stu.Age = 10

	createStudent(db, &stu)

	// 删除记录
	deleteStudent(db, stu.Id)
}
