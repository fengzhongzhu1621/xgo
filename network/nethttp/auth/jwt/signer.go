package jwt

import (
	"fmt"
	"time"

	jwt5 "github.com/golang-jwt/jwt/v5"
)

// ErrInvalidToken 非法 token，表示 Token 验证失败时返回的错误
var ErrInvalidToken = fmt.Errorf("invalid token")

// Signer 数字签名器接口
type Signer interface {
	// 将任意自定义数据（custom）签名生成 JWT Token，返回字符串形式的 Token。
	Sign(custom interface{}) (string, error)
	// 验证一个 JWT Token 是否合法，并解析出其中自定义的数据部分。
	Verify(token string) (interface{}, error)
}

// jwtSign jwt 用户身份数字签名器
type jwtSign struct {
	Secret  []byte        // 用于签名和验证的密钥（也叫 Signing Key），必须是保密的。
	Expired time.Duration // Token 的有效期，比如 1 小时、1 天等。
	Issuer  string        // Token 的发行者，通常用于标识服务端身份，在 JWT 的标准字段 iss 中体现。
}

// NewJwtSign 构造 jwt 签名
func NewJwtSign(secret []byte, expired time.Duration, issuer string) Signer {
	return &jwtSign{
		Secret:  secret,
		Expired: expired,
		Issuer:  issuer,
	}
}

// claims 元数据
// 定义了一个自定义的 JWT Claims（载荷）结构体，它嵌入了 jwt 标准的 StandardClaims，并额外添加了一个 Custom 字段用于存放任意自定义数据。
type claims struct {
	jwt5.RegisteredClaims
	// 用于存放业务相关的自定义数据，比如用户 ID、角色等。通过 JSON 标签，它在 JWT 中的字段名为 custom。
	Custom interface{} `json:"custom,omitempty"`
}

// 生成 JWT Token
func (t *jwtSign) Sign(custom interface{}) (string, error) {
	now := time.Now()
	cl := claims{
		RegisteredClaims: jwt5.RegisteredClaims{
			ExpiresAt: jwt5.NewNumericDate(
				now.Add(t.Expired),
			), // 当前时间 + 有效期（now.Add(t.Expired)），通过 jwt.At() 转换为 UNIX 时间戳。
			IssuedAt: jwt5.NewNumericDate(now), // 当前时间，表示 Token 的签发时间。
			Issuer:   t.Issuer,                 // 设置为结构体中的 t.Issuer。
		},
		Custom: custom,
	}
	// 创建一个新的 JWT Token，指定签名算法为 jwt.SigningMethodHS512（即 HMAC-SHA512）
	// HS512 签名算法，是一种对称加密算法，安全性依赖于 Secret 的保密性
	token := jwt5.NewWithClaims(jwt5.SigningMethodHS512, cl)
	// 生成最终的 JWT Token 字符串
	return token.SignedString(t.Secret)
}

// 验证一个 JWT Token 是否合法，并提取其中的自定义数据部分（Custom）
func (t *jwtSign) Verify(tokenStr string) (interface{}, error) {
	// 解析传入的 JWT Token 字符串
	token, err := jwt5.ParseWithClaims(tokenStr, &claims{},
		func(token *jwt5.Token) (i interface{}, err error) {
			// 用于提供签名验证的密钥，在这里直接返回结构体中的 t.Secret
			return t.Secret, nil
		})
	if err != nil {
		// 如果解析过程中出错（比如 Token 被篡改、过期等），则直接返回错误
		return nil, err
	}

	// 如果解析成功，通过类型断言判断 Claims 是否是我们定义的 *claims 类型
	if claim, ok := token.Claims.(*claims); ok && token.Valid {
		return claim.Custom, nil
	}
	return nil, ErrInvalidToken
}
