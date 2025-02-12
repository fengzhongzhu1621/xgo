package password

import (
	"errors"
	"unicode"
)

// PasswordGeneratePolicy 密码生成策略
type PasswordGeneratePolicy struct {
	MinLength  int // 密码最小长度
	MinNumbers int // 密码中包含数字的数量
	MinSymbols int // 密码中包含符号或标点符号的数量
	MinUpper   int // 密码中包含大写字符的数量
}

func (pp *PasswordGeneratePolicy) ValidatePassword(password string) error {
	var numbers, symbols, upper int
	for _, char := range password {
		switch {
		case unicode.IsNumber(char):
			numbers++
		case unicode.IsSymbol(char) || unicode.IsPunct(char):
			symbols++
		case unicode.IsUpper(char):
			upper++
		}
	}

	// 密码长度限制
	if len(password) < pp.MinLength {
		return errors.New("password too short")
	}

	// 密码中包含的数字数量限制
	if numbers < pp.MinNumbers {
		return errors.New("not enough numbers")
	}

	// 密码中包含符号或标点符号的数量
	if symbols < pp.MinSymbols {
		return errors.New("not enough symbols")
	}

	if upper < pp.MinUpper {
		return errors.New("not enough uppercase letters")
	}

	return nil
}
