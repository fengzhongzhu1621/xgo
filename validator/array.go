package validator

import (
	"fmt"

	"github.com/fengzhongzhu1621/xgo/cast"
	"github.com/gin-gonic/gin/binding"
)

// ValidateArray 接受一个任意类型的参数 data，并返回一个布尔值和一个字符串。
// 主要目的是验证传入的数据是否是一个非空数组，并且数组中的每个元素都通过结构体验证。
func ValidateArray(data any) (bool, string) {
	// 将 data 转换为数组切片
	array, err := cast.ToSlice2(data)
	if err != nil {
		return false, err.Error()
	}

	// 判断是否为空数组
	if len(array) == 0 {
		return false, "the array should contain at least 1 item"
	}

	// 非空数组
	for index, item := range array {
		// binding.Validator.ValidateStruct 验证传入的对象是否符合预定义的结构或结构体指针
		// 实现了一个通用的结构体验证方法，它可以处理结构体、结构体指针、切片和数组类型的对象。
		// 对于非结构体类型的对象，它会忽略验证。对于结构体、结构体指针、切片和数组类型的对象，它会递归地进行验证，并收集所有的验证错误。
		// 如果所有元素都通过验证，则返回nil；否则，返回包含所有错误的validateRrt切片。
		if err := binding.Validator.ValidateStruct(item); err != nil {
			message := fmt.Sprintf("data in array[%d], %s", index, ValidationErrorMessage(err))
			return false, message
		}
	}

	return true, "valid"
}
