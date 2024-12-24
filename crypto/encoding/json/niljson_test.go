package json

import (
	"encoding/json"
	"fmt"
	"github.com/fengzhongzhu1621/xgo/tests"
	"github.com/stretchr/testify/assert"
	"github.com/wneessen/niljson"
	"testing"
)

type JSONType struct {
	Bool1   niljson.NilBoolean `json:"bool1"`             // "bool1":true
	Bool2   niljson.NilBoolean `json:"bool2"`             // "bool2":null
	Bool3   niljson.NilBoolean `json:"bool3,omitempty"`   // "bool3":null omitempty 不起作用
	Float32 niljson.NilFloat32 `json:"float32,omitempty"` // omitempty 不起作用
	Float64 niljson.NilFloat64 `json:"float64"`
	// int
	Int1  niljson.NilInt   `json:"int1"`           // "int1":123
	Int2  niljson.NilInt   `json:"int2"`           // "int2":null
	Int3  int              `json:"int3"`           // "int3":0
	Int4  int              `json:"int4,omitempty"` // omitempty 不起作用
	Int5  niljson.NilInt   `json:"int5,omitempty"` // "int5":null
	Int64 niljson.NilInt64 `json:"int64"`          // "int64":1234567890
	// nil
	NullString niljson.NilString `json:"nil"` // "nil":null
	// string
	String1  niljson.NilString `json:"string1"`            // str
	String2  niljson.NilString `json:"string2"`            // ""
	String3  niljson.NilString `json:"string3"`            // null
	String4  string            `json:"string4"`            // str
	String5  string            `json:"string5"`            // ""
	String6  string            `json:"string6"`            // ""
	String7  string            `json:"string7,omitempty"`  // 忽略 null 值
	String8  niljson.NilString `json:"string8,omitempty"`  // "string8":""
	String9  string            `json:"string9,omitempty"`  // 忽略空值
	String10 niljson.NilString `json:"string10,omitempty"` // "string10":null
}

func TestNilJson(t *testing.T) {
	data := []byte(`{
		"bool1": true,
		"bool2": null,
		"bool3": null,
		"float32": null,
		"float64": 0,
		"int1": 123,
		"int2": null,
		"int3": null,
		"int4": null,
		"int5": null,
		"int64": 1234567890,
		"nil": null,
		"string1": "str",
		"string2": "",
		"string3": null,
		"string4": "str",
		"string5": "",
		"string6": null,
		"string7": "",
		"string8": "",
		"string9": null,
		"string10": null
 }`)

	// 将字符串转换为结构体对象
	var example JSONType
	if err := json.Unmarshal(data, &example); err != nil {
		fmt.Println("failed to unmarshal JSON:", err)
		return
	}

	if example.Bool1.NotNil() {
		// Bool1 is: true
		fmt.Printf("Bool1 is: %t \n", example.Bool1.Value())
		// Bool2 is: false
		fmt.Printf("Bool2 is: %t \n", example.Bool2.Value())
		// Bool3 is: false
		fmt.Printf("Bool3 is: %t \n", example.Bool3.Value())
	}
	if example.Float32.IsNil() {
		// Float 32 is nil, value is 0.000000
		fmt.Printf("Float 32 is nil, value is %f\n", example.Float32.Value())
	}
	if example.Float64.NotNil() {
		// Float 64 is: 0.000000
		fmt.Printf("Float 64 is: %f \n", example.Float64.Value())
	}
	if example.String1.NotNil() {
		// String is: str
		fmt.Printf("String1 is: %s \n", example.String1.Value())
		// String2 is:
		fmt.Printf("String2 is: %s \n", example.String2.Value())
		// String3 is:
		fmt.Printf("String3 is: %s \n", example.String3.Value())
		// String4 is: str
		fmt.Printf("String4 is: %s \n", example.String4)
		// String5 is:
		fmt.Printf("String5 is: %s \n", example.String5)
		// String6 is:
		fmt.Printf("String6 is: %s \n", example.String6)
	}

	data, err := json.Marshal(&example)
	if err != nil {
		fmt.Printf("failed to marshal JSON: %s", err)
		return
	}
	expect := `{"bool1":true,"bool2":null,"bool3":null,"float32":null,"float64":0,"int1":123,"int2":null,"int3":0,"int5":null,"int64":1234567890,"nil":null,"string1":"str","string2":"","string3":null,"string4":"str","string5":"","string6":"","string8":"","string10":null}`
	assert.Equal(t, expect, string(data))

	tests.PrintStruct(example)
}
