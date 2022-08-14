package utils

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
