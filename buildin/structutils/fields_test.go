package structutils

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/duke-git/lancet/v2/structs"
	"github.com/stretchr/testify/assert"
)

// TestFields Get all fields of a given struct, that the fields are abstract struct field
// func (s *Struct) Fields() []*Field
// func (s *Struct) Field(name string) *Field
func TestFields(t *testing.T) {
	type People struct {
		Name string `json:"name"`
	}
	p1 := &People{Name: "abc"}
	s := structs.New(p1)
	fields := s.Fields()
	fmt.Println(len(fields)) // 1

	f, ok := s.Field("Name")
	assert.Equal(t, true, ok)
	assert.Equal(t, "abc", f.Value())

	_, ok1 := s.Field("Unknown")
	assert.Equal(t, false, ok1)
}

// TesIsEmbedded Check if the field is an embedded field
// func (f *Field) IsEmbedded() bool
func TestIsEmbedded(t *testing.T) {
	type Parent struct {
		Name string
	}
	type Child struct {
		Parent
		Age int
	}
	c1 := &Child{}
	c1.Name = "111"
	c1.Age = 11

	s := structs.New(c1)
	n, _ := s.Field("Name")
	a, _ := s.Field("Age")

	assert.Equal(t, true, n.IsEmbedded())
	assert.Equal(t, false, a.IsEmbedded())
}

// TestIsExported Check if the field is exported
// func (f *Field) IsExported() bool
func TestIsExported(t *testing.T) {
	type Parent struct {
		Name string
		age  int
	}
	p1 := &Parent{Name: "11", age: 11}
	s := structs.New(p1)
	n, _ := s.Field("Name")
	a, _ := s.Field("age")

	assert.Equal(t, true, n.IsExported())
	assert.Equal(t, false, a.IsExported())
}

// TestIsZero Check if the field is exported
// func (f *Field) IsZero() bool
func TestIsZero(t *testing.T) {
	type Parent struct {
		Name string
		Age  int
	}
	p1 := &Parent{Age: 11}
	s := structs.New(p1)
	n, _ := s.Field("Name")
	a, _ := s.Field("Age")

	assert.Equal(t, true, n.IsZero())
	assert.Equal(t, false, a.IsZero())
}

// TestIsSlice Check if the field is a slice
// func (f *Field) IsSlice() bool
func TestIsSlice(t *testing.T) {
	type Parent struct {
		Name string
		arr  []int
	}

	p1 := &Parent{arr: []int{1, 2, 3}}
	s := structs.New(p1)
	a, _ := s.Field("arr")

	assert.Equal(t, true, a.IsSlice())
}

// TestName Get the field name
// func (f *Field) Name() string
func TestName(t *testing.T) {
	type Parent struct {
		Name string
		Age  int
	}
	p1 := &Parent{Age: 11}
	s := structs.New(p1)
	n, _ := s.Field("Name")
	a, _ := s.Field("Age")

	assert.Equal(t, "Name", n.Name())
	assert.Equal(t, "Age", a.Name())
}

// TestKind Get the field's kind
// func (f *Field) Kind() reflect.Kind
func TestKind(t *testing.T) {
	type Parent struct {
		Name string
		Age  int
	}
	p1 := &Parent{Age: 11}
	s := structs.New(p1)
	n, _ := s.Field("Name")
	a, _ := s.Field("Age")

	assert.Equal(t, reflect.String, n.Kind())
	assert.Equal(t, reflect.Int, a.Kind())
}

// TestIsTargetType check if a struct field type is target type or not
// func (f *Field) IsTargetType(targetType reflect.Kind) bool
func TestIsTargetType(t *testing.T) {
	type Parent struct {
		Name string
		arr  []int
	}

	p1 := &Parent{arr: []int{1, 2, 3}}
	s := structs.New(p1)
	n, _ := s.Field("Name")
	a, _ := s.Field("arr")

	assert.Equal(t, true, n.IsTargetType(reflect.String))
	assert.Equal(t, true, a.IsTargetType(reflect.Slice))
}

// TestTag Get a `Tag` of the `Field`, `Tag` is a abstract struct field tag
// func (f *Field) Tag() *Tag
func TestTag(t *testing.T) {
	type Parent struct {
		Name string `json:"name,omitempty"`
	}
	p1 := &Parent{"111"}

	s := structs.New(p1)
	n, _ := s.Field("Name")
	tag := n.Tag()

	assert.Equal(t, "name", tag.Name)
}
