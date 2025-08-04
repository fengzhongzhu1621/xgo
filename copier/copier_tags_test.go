package copier_test

import (
	"testing"

	"github.com/jinzhu/copier"
)

type EmployeeTags struct {
	// 如果字段不能拷贝，则 panic
	Name    string `copier:"must"`
	DOB     string
	Address string
	// 忽略拷贝这个字段
	ID int `copier:"-"`
}

type User1 struct {
	Name    string
	DOB     string
	Address string
	ID      int
}

type User2 struct {
	DOB     string
	Address string
	ID      int
}

// 测试自定字段不拷贝
func TestCopyTagIgnore(t *testing.T) {
	employee := EmployeeTags{ID: 100}
	user := User1{
		Name:    "Dexter Ledesma",
		DOB:     "1 November, 1970",
		Address: "21 Jump Street",
		ID:      12345,
	}
	// 不拷贝 ID 字段
	copier.Copy(&employee, user)
	if employee.ID == user.ID {
		t.Error("Was not expected to copy IDs")
	}
	if employee.ID != 100 {
		t.Error("Original ID was overwritten")
	}
}

// 字段不存在触发 panic
func TestCopyTagMust(t *testing.T) {
	employee := &EmployeeTags{}
	// User2 不包含 Name 字段
	user := &User2{DOB: "1 January 1970"}
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected a panic.")
		}
	}()
	copier.Copy(employee, user)
}

// 根据字段映射拷贝字段
func TestCopyTagFieldName(t *testing.T) {
	t.Run("another name field copy", func(t *testing.T) {
		type SrcTags struct {
			FieldA string
			FieldB string `copier:"Field2"`
			FieldC string `copier:"FieldTagMatch"`
		}

		type DestTags struct {
			Field1 string `copier:"FieldA"`
			Field2 string
			Field3 string `copier:"FieldTagMatch"`
		}

		dst := &DestTags{}
		src := &SrcTags{
			FieldA: "FieldA->Field1",
			FieldB: "FieldB->Field2",
			FieldC: "FieldC->Field3",
		}
		err := copier.Copy(dst, src)
		if err != nil {
			t.Fatal(err)
		}

		if dst.Field1 != src.FieldA {
			t.Error("Field1 no copy")
		}
		if dst.Field2 != src.FieldB {
			t.Error("Field2 no copy")
		}
		if dst.Field3 != src.FieldC {
			t.Error("Field3 no copy")
		}
	})

	t.Run("validate error flag name", func(t *testing.T) {
		type SrcTags struct {
			field string
		}

		type DestTags struct {
			Field1 string `copier:"field"`
		}

		dst := &DestTags{}
		src := &SrcTags{
			field: "field->Field1",
		}
		err := copier.Copy(dst, src)
		if err == nil {
			t.Fatal("must validate error")
		}
	})
}
