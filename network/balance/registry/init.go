package registry

import (
	"github.com/go-playground/validator/v10"
)

// g_validator 声明了一个全局的 validator.Validate 类型的指针变量 g_validator。这个变量将在程序的其他部分用于执行验证操作。
var g_validator *validator.Validate

func init() {
	g_validator = validator.New()
}
