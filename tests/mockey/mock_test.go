package mockey

import (
	"fmt"
	"testing"

	. "github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"
)

// Simple function
func Foo(in string) string {
	return in
}

// Function method (value receiver)
type A struct{}

func (a A) Foo(in string) string { return in }

// Function method (pointer receiver)
type B struct{}

func (b *B) Foo(in string) string { return in }

// Value
var Bar = 0

func TestMockXXX(t *testing.T) {

	PatchConvey("Function mocking", t, func() {
		Mock(Foo).Return("c").Build()         // mock function
		So(Foo("anything"), ShouldEqual, "c") // assert `Foo` is mocked
	})

	PatchConvey("Method mocking (value receiver)", t, func() {
		Mock(A.Foo).Return("c").Build()              // mock method
		So(new(A).Foo("anything"), ShouldEqual, "c") // assert `A.Foo` is mocked
	})

	PatchConvey("Method mocking (pointer receiver)", t, func() {
		Mock((*B).Foo).Return("c").Build() // mock method
		b := &B{}
		So(b.Foo("anything"), ShouldEqual, "c") // assert `*B.Foo` is mocked
	})

	PatchConvey("Variable mocking", t, func() {
		MockValue(&Bar).To(1)   // mock variable
		So(Bar, ShouldEqual, 1) // assert `Bar` is mocked
	})

	// the mocks are released automatically outside `PatchConvey`
	fmt.Println(Foo("a"))        // a
	fmt.Println(new(A).Foo("b")) // b
	fmt.Println(Bar)             // 0
}
