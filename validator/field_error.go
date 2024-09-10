package validator

import (
	"fmt"
	"io"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

type ValidationFieldError struct {
	Err validator.FieldError
}

func (v ValidationFieldError) String() string {
	e := v.Err

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", e.Field(), e.Param())
	case "min":
		return fmt.Sprintf("%s must be longer than %s", e.Field(), e.Param())
	case "email":
		return "Invalid email format"
	case "len":
		return fmt.Sprintf("%s must be %s characters long", e.Field(), e.Param())
	case "gt":
		return fmt.Sprintf("%s must greater than %s", e.Field(), e.Param())
	case "gte":
		return fmt.Sprintf("%s must greater or equals to %s", e.Field(), e.Param())
	case "lt":
		return fmt.Sprintf("%s must less than %s", e.Field(), e.Param())
	case "lte":
		return fmt.Sprintf("%s must less or equals to %s", e.Field(), e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of '%s'", e.Field(), e.Param())
	}

	return fmt.Sprintf("%s is not valid, condition: %s", e.Field(), e.ActualTag())
}

// ValidationErrorMessage 接受一个 error 类型的参数 err，并返回一个字符串。
// 这个函数的主要目的是根据传入的错误类型生成相应的错误信息。
func ValidationErrorMessage(err error) string {
	// 遇到了文件结束符，这通常意味着 JSON 解码失败。
	if err == io.EOF {
		return "EOF, json decode fail"
	}

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		// 错误可能是 JSON 解码失败或其他非验证错误。
		message := fmt.Sprintf("json decode or validate fail, err=%s", err)
		log.Info(message)
		return message
	}

	// 错误是由验证器产生的
	// currently, only return the first error
	for _, fieldErr := range validationErrs {
		return ValidationFieldError{fieldErr}.String()
	}

	return "validationErrs with no error message"
}
