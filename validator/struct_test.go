package validator

import (
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
)

// validate:"required,dive,required" 是一个组合标签，它告诉验证器对结构体中的某个字段执行一系列的验证规则。
// dive dive 是一个特殊的验证关键字，它用于对切片、数组或映射中的每个元素执行后续指定的验证规则。
// 当在 dive 后面再跟一个验证规则时，比如这里的第二个 required，它意味着你需要对集合中的每一个元素单独应用这个规则。
// 1. Names 字段本身必须是存在的（即非nil），这是第一个 required 规则的作用。
// 2. Names 切片中的每一个元素（即每一个字符串）都必须是非空的，这是 dive 和第二个 required 规则共同作用的结果。
type User struct {
	Names []string `validate:"required,dive,required"`
}

type User2 struct {
	Username string `validate:"required"`
	Age      int    `validate:"gte=0,lte=130"`
	Email    string `validate:"required,email"`
}

func TestUser2(t *testing.T) {
	user := &User2{
		Username: "username_1",
		Age:      18,
		Email:    "username_1@example.com",
	}

	err := validatorer.Struct(user)
	if err != nil {
		// 判断是否是 validator 第三方库的校验错误
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			fmt.Println(validationErrors)
		}
	} else {
		fmt.Println("Validation passed")
	}
}
