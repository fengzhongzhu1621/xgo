package cast

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/util/gconv"

	"github.com/duke-git/lancet/v2/maputil"
)

// Converts map to struct
// func MapTo(src any, dst any) error
// func MapToStruct(m map[string]any, structObj any) error
func TestMapToStruct(t *testing.T) {
	{
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
		// personInfo(map) -> p(struct)
		err := maputil.MapTo(personInfo, &p)

		fmt.Println(err)
		fmt.Println(p)

		// Output:
		// <nil>
		// {Nothin 28 123456789 {test 1}}
	}

	{
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
		// personReqMap(map) -> personReq(struct)
		_ = maputil.MapToStruct(personReqMap, &personReq)
		fmt.Println(personReq)

		// Output:
		// {Nothin 35 1 10}
	}

	{
		type User struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		userMap := map[string]interface{}{
			"name": "Alice",
			"age":  30,
		}

		// map 转换成结构体
		var user User
		err := gconv.Struct(userMap, &user)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println(user) // 输出：{Alice 30}
	}
}

func TestMapSliceToStruct(t *testing.T) {
	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	usersSlice := []map[string]interface{}{
		{"name": "Tom", "age": 30},
		{"name": "Jerry", "age": 25},
	}

	// map 切片 转换成 结构体 切片
	var users []User
	err := gconv.Structs(usersSlice, &users)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(users) // 输出：[{Tom 30} {Jerry 25}]
}
