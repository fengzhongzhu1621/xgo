package maps

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/maputil"
)

// TestMapTo Rry to map any interface to struct or base type.
// func MapTo(src any, dst any) error
func TestMapTo(t *testing.T) {
	type (
		Address struct {
			Street string `json:"street"`
			Number int    `json:"number"`
		}

		Person struct {
			Name  string  `json:"name"`
			Age   int     `json:"age"`
			Phone string  `json:"phone"`
			Addr  Address `json:"address"`
		}
	)

	personInfo := map[string]interface{}{
		"name":  "Nothin",
		"age":   28,
		"phone": "123456789",
		"address": map[string]interface{}{
			"street": "test",
			"number": 1,
		},
	}

	var p Person
	err := maputil.MapTo(personInfo, &p)

	fmt.Println(err)
	fmt.Println(p)

	// Output:
	// <nil>
	// {Nothin 28 123456789 {test 1}}
}

// TestMapToStruct Converts map to struct
// func MapToStruct(m map[string]any, structObj any) error
func TestMapToStruct(t *testing.T) {
	personReqMap := map[string]any{
		"name":     "Nothin",
		"max_age":  35,
		"page":     1,
		"pageSize": 10,
	}

	type PersonReq struct {
		Name     string `json:"name"`
		MaxAge   int    `json:"max_age"`
		Page     int    `json:"page"`
		PageSize int    `json:"pageSize"`
	}
	var personReq PersonReq
	_ = maputil.MapToStruct(personReqMap, &personReq)
	fmt.Println(personReq)

	// Output:
	// {Nothin 35 1 10}
}
