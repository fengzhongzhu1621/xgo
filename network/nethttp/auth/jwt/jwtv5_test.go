package jwt

import (
	"fmt"
	"testing"
	"time"

	jwt5 "github.com/golang-jwt/jwt/v5"
)

// 生成Token
func TestJwtV5GenerateToken(t *testing.T) {
	// 生成Token
	key := []byte("your-secret-key")
	claims := jwt5.MapClaims{
		"user": "example",
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt5.NewWithClaims(jwt5.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(key) // 使用HMAC-SHA256签名
	// fmt.Println(value)

	// 验证Token
	token, _ = jwt5.ParseWithClaims(
		tokenString,
		&jwt5.MapClaims{},
		func(t *jwt5.Token) (interface{}, error) {
			return key, nil // 返回密钥用于验证
		},
	)
	if claims, ok := token.Claims.(*jwt5.MapClaims); ok && token.Valid {
		fmt.Println(claims) // &map[exp:1.753765223e+09 user:example]
	}
}
