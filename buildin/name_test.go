package buildin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct{}

type stringerStruct struct{}

func (stringerStruct) String() string {
	return "stringer"
}

func TestStructName(t *testing.T) {
	testCases := []struct {
		Name         string
		Struct       interface{}
		ExpectedName string
	}{
		{
			Name:         "simple_struct",
			Struct:       testStruct{},
			ExpectedName: "utils.testStruct",
		},
		{
			Name:         "pointer_struct",
			Struct:       &testStruct{},
			ExpectedName: "utils.testStruct",
		},
		{
			Name:         "stringer",
			Struct:       stringerStruct{},
			ExpectedName: "stringer",
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			s := StructName(c.Struct)
			assert.Equal(t, c.ExpectedName, s)
		})
	}
}

func TestFullyQualifiedStructName(t *testing.T) {
	type Object struct{}

	assert.Equal(t, "utils.Object", FullyQualifiedStructName(Object{}))
	assert.Equal(t, "utils.Object", FullyQualifiedStructName(&Object{}))
}

func BenchmarkFullyQualifiedStructName(b *testing.B) {
	type Object struct{}
	o := Object{}

	for i := 0; i < b.N; i++ {
		FullyQualifiedStructName(o)
	}
}

func TestStructName2(t *testing.T) {
	type Object struct{}

	assert.Equal(t, "Object", RawStructName(Object{}))
	assert.Equal(t, "Object", RawStructName(&Object{}))
}

func TestNamedStruct(t *testing.T) {
	assert.Equal(t, "named object", NamedStruct(RawStructName)(namedObject{}))
	assert.Equal(t, "named object", NamedStruct(RawStructName)(&namedObject{}))

	// Test fallback
	type Object struct{}

	assert.Equal(t, "Object", NamedStruct(RawStructName)(Object{}))
	assert.Equal(t, "Object", NamedStruct(RawStructName)(&Object{}))
}

type namedObject struct{}

func (namedObject) Name() string {
	return "named object"
}
