package reflectutils

import (
	"bytes"
	"reflect"
	"testing"
)

type myType struct{ bytes.Buffer }

func (myType) valueMethod() {}
func (myType) ValueMethod() {}

func (*myType) pointerMethod() {}
func (*myType) PointerMethod() {}

func TestNameOf(t *testing.T) {
	tests := []struct {
		fnc  interface{}
		want string
	}{
		{TestNameOf, "reflectutils.TestNameOf"},
		{func() {}, "reflectutils.TestNameOf.func1"},
		{(myType).valueMethod, "reflectutils.myType.valueMethod"},
		{(myType).ValueMethod, "reflectutils.myType.ValueMethod"},
		{(myType{}).valueMethod, "reflectutils.myType.valueMethod"},
		{(myType{}).ValueMethod, "reflectutils.myType.ValueMethod"},
		{(*myType).valueMethod, "reflectutils.myType.valueMethod"},
		{(*myType).ValueMethod, "reflectutils.myType.ValueMethod"},
		{(&myType{}).valueMethod, "reflectutils.myType.valueMethod"},
		{(&myType{}).ValueMethod, "reflectutils.myType.ValueMethod"},
		{(*myType).pointerMethod, "reflectutils.myType.pointerMethod"},
		{(*myType).PointerMethod, "reflectutils.myType.PointerMethod"},
		{(&myType{}).pointerMethod, "reflectutils.myType.pointerMethod"},
		{(&myType{}).PointerMethod, "reflectutils.myType.PointerMethod"},
		{(*myType).Write, "reflectutils.myType.Write"},
		{(&myType{}).Write, "bytes.Buffer.Write"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := NameOf(reflect.ValueOf(tt.fnc))
			if got != tt.want {
				t.Errorf("NameOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
