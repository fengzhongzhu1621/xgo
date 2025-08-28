package main

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
)

func TestTotp(t *testing.T) {
	// 1. 生成共享密钥
	// 生成的密钥包含一个 URI，格式类似于 otpauth://totp/Github:a@b?secret=...&issuer=Github
	// 可以用于配置 TOTP 客户端(如 Google Authenticator)
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Github",
		AccountName: "a@b",
	})
	if err != nil {
		log.Fatalf("Failed to generate TOTP key: %v", err)
	}

	// otpauth://totp/Github:a@b?algorithm=SHA1&digits=6&issuer=Github&period=30&secret=AMAROVJB5EHAVOC7AJN4ZFELMEGW2WTM
	fmt.Println("Generated TOTP Key URI:")
	fmt.Println(key.URL())

	// 2. 获取当前时间戳生成的 TOTP
	// TOTP 依赖于客户端和服务器的时间同步
	// 默认情况下，TOTP 代码每 30 秒更新一次
	code, err := totp.GenerateCode(key.Secret(), time.Now())
	if err != nil {
		log.Fatalf("Failed to generate TOTP code: %v", err)
	}

	// Generated TOTP Code: 348182
	fmt.Printf("Generated TOTP Code: %s\n", code)

	// 3. 验证 TOTP
	// 如果客户端和服务器时间不同步，验证可能会失败
	isValid := totp.Validate(code, key.Secret())
	if isValid {
		fmt.Println("TOTP Code is valid!")
	} else {
		fmt.Println("Invalid TOTP Code.")
	}
}
