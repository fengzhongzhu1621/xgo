package jwtx

import (
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

var (
	JwtIssuer = "xgo"
)

type UserIdJwtClaims struct {
	UserId string `json:"user_id"`
	Salt   string `json:"salt"`
	Data   string `json:"data"`
	jwt.StandardClaims
}

// NewJwtClaims 创建Jwt断言信息
func NewJwtClaims(token_expired time.Time, uid, data, salt string) *UserIdJwtClaims {
	return &UserIdJwtClaims{
		UserId: uid,
		Salt:   salt,
		Data:   data,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(token_expired),
			IssuedAt:  jwt.At(time.Now()),
			Issuer:    JwtIssuer,
		},
	}
}

// GenerateJwtToken 使用 HS256 算法生成JWT-token
func GenerateJwtToken(c *UserIdJwtClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	ss, err := token.SignedString([]byte(secret))
	return ss, err
}

// ParseJwtToken 使用密钥解析 jwt token
func ParseJwtToken(jwtStr, secret string) (*UserIdJwtClaims, error) {
	token, err := jwt.ParseWithClaims(jwtStr, &UserIdJwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*UserIdJwtClaims); ok && token.Valid {
		return claims, err
	} else {
		return nil, err
	}
}
