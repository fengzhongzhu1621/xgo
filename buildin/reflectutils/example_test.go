package reflectutils

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExample(t *testing.T) {
	var x float64 = 3.4

	t1 := reflect.TypeOf(x)  //返回给定值的类型
	v1 := reflect.ValueOf(x) // 返回给定值的值

	fmt.Println("Type: ", t1)  // Type:  float64
	fmt.Println("Value: ", v1) // Value:  3.4

	v2 := v1.Kind()
	fmt.Println("Kind: ", v2) // Kind:  float64

	t2 := v1.Type()
	fmt.Println("v1.Type: ", t2) // v1.Type:  float64

	assert.Equal(t, v1.Kind(), reflect.Float64)
	assert.Equal(t, reflect.TypeOf(x), reflect.TypeOf(float64(0)))
}

func TestElemSetFloat(t *testing.T) {
	var x float64 = 3.14

	// Elem(): 它允许你访问指针指向的值或接口包含的值。
	// 接受一个 reflect.Value 类型的参数，并返回一个新的 reflect.Value，表示指针指向的值或接口包含的值。
	// 如果传入的值不是指针或接口类型，Elem() 函数将引发 panic。
	v := reflect.ValueOf(&x).Elem() // 返回给定值的值
	v.SetFloat(3.1415926)
	fmt.Println("x: ", v) // x:  3.1415926
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type User struct {
	Name string `validate:"required,min=3"`
	Age  int    `validate:"required,min=0,max=120"`
}

func TestField(t *testing.T) {
	p := Person{Name: "bob", Age: 10}
	t1 := reflect.TypeOf(p)
	fmt.Println("Type: ", t1) // Type:  reflectutils.Person

	// 遍历结构体的所有字段
	for i := 0; i < t1.NumField(); i++ {
		// 按索引获得字段
		field := t1.Field(i)
		// Field: Name, Tag: name
		// Field: Age, Tag: age
		fmt.Printf("Field: %s, Tag: %s\n", field.Name, field.Tag.Get("json"))
	}

	// 打印字段
	// interface{} 是一个特殊的类型，称为空接口。空接口没有任何方法，因此所有类型都实现了空接口。
	// 这意味着你可以将任何类型的值赋给空接口类型的变量。这使得空接口成为处理不确定类型数据的强大工具。
	printFields := func(s interface{}) {
		v := reflect.ValueOf(s)

		// Name: bob
		// Age: 10
		if v.Kind() == reflect.Struct {
			t2 := v.Type()
			for i := 0; i < v.NumField(); i++ {
				field := t2.Field(i)
				value := v.Field(i)
				fmt.Printf("%s: %v\n", field.Name, value)
			}
		} else {
			fmt.Println("Expected a struct")
		}
	}

	printFields(p)
}

type Calculator struct{}

func (c Calculator) Add(a, b int) int { return a + b }

func TestMethodByName(t *testing.T) {
	// 转换为 reflect.Value
	c := Calculator{}
	v := reflect.ValueOf(c)

	// 获取对象的方法
	method := v.MethodByName("Add")

	// 构造函数执行参数
	args := []reflect.Value{
		reflect.ValueOf(1),
		reflect.ValueOf(2),
	}
	// 调用方法
	result := method.Call(args)

	// 结果: 3
	fmt.Println("结果:", result[0].Int())
}

func TestNumField(t *testing.T) {
	u := User{"张三", 20}

	t1 := reflect.TypeOf(u)
	v := reflect.ValueOf(u)

	// 类型: reflectutils.User
	fmt.Printf("类型: %v\n", t1)
	// 值: {张三 20}
	fmt.Printf("值: %v\n", v)

	// 遍历结构体字段
	// Name: 张三
	// Age: 20
	for i := 0; i < t1.NumField(); i++ {
		field := t1.Field(i)
		value := v.Field(i)
		fmt.Printf("%s: %1v\n", field.Name, value)
	}
}

type Validator interface {
	Validate() error
}

func TestValidateTag(t *testing.T) {
	s := User{
		Name: "张",
		Age:  -1,
	}

	var errors []string
	// 获取结构体的反射值
	v := reflect.ValueOf(s)
	// 如果是指针，获取其元素
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	// 确保是结构体
	if v.Kind() != reflect.Struct {
		return
	}

	t1 := v.Type()
	// 遍历所有字段
	for i := 0; i < t1.NumField(); i++ {
		field := t1.Field(i)
		value := v.Field(i)
		// 获取验证规则
		tags := strings.Split(field.Tag.Get("validate"), ",")
		for _, tag := range tags {
			// 处理required标签
			if tag == "required" {
				// 判断是否是默认值
				if value.Interface() == reflect.Zero(value.Type()).Interface() {
					errors = append(errors, fmt.Sprintf("%s 是必填字段", field.Name))
				}
			}
			// 处理min标签
			if strings.HasPrefix(tag, "min=") {
				// 这里简化处理，实际应该根据字段类型分别处理
				if value.Kind() == reflect.Int {
					if value.Int() < 0 {
						errors = append(errors, fmt.Sprintf("%s 不能小于0", field.Name))
					}
				}
			}
		}
	}

	if len(errors) > 0 {
		fmt.Println("验证错误:")
		for _, err := range errors {
			fmt.Printf("- %s\n", err)
		}
	} else {
		fmt.Println("验证通过")
	}
}
