package maps

import (
	"fmt"

	"github.com/duke-git/lancet/v2/maputil"
)

// TestMapTo Rry to map any interface to struct or base type.
// func MapTo(src any, dst any) error
func TestMapTo() {
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
