package form

import (
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
)

// form:"Code"这种写法通常用于表单验证库
// 此库允许你在结构体字段上添加标签（tags），以便在验证表单数据时执行自定义规则。
type User struct {
	Name string `form:"Name"`
	Code string `form:"Code" validate:"required,min=4,max=10"`
}

func TestValidator(t *testing.T) {
	validate := validator.New()

	user := User{
		Name: "Alice",
		Code: "1234",
	}

	err := validate.Struct(user)
	if err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation succeeded")
	}
}
