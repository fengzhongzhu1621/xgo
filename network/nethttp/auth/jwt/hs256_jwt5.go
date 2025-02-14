package jwt

import (
	"fmt"
	"time"

	jwt5 "github.com/golang-jwt/jwt/v5"
)

type HS256ClaimsV5 struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	// jwt包自带的jwt.RegisteredClaims只包含了官方字段
	jwt5.RegisteredClaims // 内嵌标准的声明
}

// HS256JWTV5Manager JWT（JSON Web Token）管理器
// HS256 使用对称密钥，适用于单服务器或信任的环境。如果需要在多个服务器或不可信环境中使用，建议使用非对称签名算法如 RS256。
type HS256JWTV5Manager struct {
	secretKey []byte        // 用于签名和验证 JWT 的密钥。应保持机密，避免泄露。
	expires   time.Duration // JWT 的有效期，表示从签发时间起多久后过期。
}

func NewHS256JWTV5Manager(secretKey string, expires time.Duration) *HS256JWTV5Manager {
	return &HS256JWTV5Manager{
		secretKey: []byte(secretKey),
		expires:   expires,
	}
}

// Generate 生成 jwt token
func (m *HS256JWTV5Manager) Generate(Issuer string, userID string, role string) (string, error) {
	// 创建一个新的 JWT，指定签名方法为 HS256（HMAC with SHA-256）。
	expirationTime := jwt5.NewNumericDate(time.Now().Add(m.expires))
	claims := &HS256ClaimsV5{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt5.RegisteredClaims{
			ExpiresAt: expirationTime, // 定义过期时间
			Issuer:    Issuer,         // 签发人
			IssuedAt:  jwt5.NewNumericDate(time.Now()),
		},
	}
	token := jwt5.NewWithClaims(jwt5.SigningMethodHS256, claims)

	// 对 JWT 进行签名，并返回生成的 JWT 字符串
	return token.SignedString(m.secretKey)
}

func (m *HS256JWTV5Manager) Verify(tokenStr string) (*HS256ClaimsV5, error) {
	// 解析传入的 JWT 字符串
	token, err := jwt5.ParseWithClaims(tokenStr, new(HS256ClaimsV5), func(token *jwt5.Token) (interface{}, error) {
		// 回调函数，用于验证签名方法和获取密钥
		// 检查 JWT 使用的签名方法是否为 HS256
		if _, ok := token.Method.(*jwt5.SigningMethodHMAC); !ok {
			// 如果签名方法不符合预期，返回错误，防止潜在的安全风险
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*HS256ClaimsV5); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
