package cast

import (
	"fmt"
	"log"
	"testing"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/mitchellh/mapstructure"
)

// Converts map to struct
// func MapTo(src any, dst any) error
// func MapToStruct(m map[string]any, structObj any) error
func TestMapToStruct(t *testing.T) {
	{
		type Address struct {
			Street  string `mapstructure:"street"`
			City    string `mapstructure:"city"`
			Zipcode string `mapstructure:"zipcode"`
		}

		type User struct {
			ID      string   `mapstructure:"id"`
			Name    string   `mapstructure:"name"`
			Email   string   `mapstructure:"email"`
			Age     int      `mapstructure:"age"`
			Active  bool     `mapstructure:"active"`
			Address Address  `mapstructure:"address"`
			Roles   []string `mapstructure:"roles"`
		}

		data := map[string]interface{}{
			"id":     "12345",
			"name":   "李四",
			"email":  "lisi@example.com",
			"age":    25,
			"active": false,
			"address": map[string]interface{}{
				"street":  "长安街",
				"city":    "北京",
				"zipcode": "100000",
			},
			"roles": []interface{}{"admin", "editor"},
		}

		var user User
		err := mapstructure.Decode(data, &user)
		if err != nil {
			log.Fatalf("无法解码数据: %v", err)
		}

		fmt.Printf("用户信息: %+v\n", user)
	}
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
