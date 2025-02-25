package jwt

import (
	"crypto/rsa"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/logging"
)

// JWTRsa JWTRsa结构，用于存储RSA算法的私钥和公钥
type JWTRsa struct {
	PrivateKey string
	PublicKey  string
}

// CustomJwtClaimsOption 自定义JwtClaims选项
type CustomJwtClaimsOption struct {
	// 操作人
	Operator string
	// 额外的信息或元数据
	Extra string
	// jwt 认证过期时间
	Expire time.Duration
	// 全局配置对象
	Cfg *config.Config
	// 后台服务之间的 jwt 认证 token，用于HS256算法的密钥
	HS256Key string
	// 一个JWTRsa结构体，包含用于RS256算法的RSA私钥和公钥
	RS256Key JWTRsa
}

// CustomJwtClaims 自定义JwtClaims内容
type CustomJwtClaims struct {
	Operator    string          `json:"operator"`
	Extra       string          `json:"extra"`
	logger      *zap.Logger     `json:"-"`
	HS256Key    string          `json:"-"`
	RS256PriKey *rsa.PrivateKey `json:"-"`
	RS256PubKey *rsa.PublicKey  `json:"-"`

	// JWT标准声明的一部分，包含如iss（发行者）、sub（主题）、exp（过期时间）等字段
	jwt.StandardClaims
}

// NewCustomJwtClaims 创建JWT断言，支持多种加密算法
func NewCustomJwtClaims(option *CustomJwtClaimsOption) (*CustomJwtClaims, error) {
	var (
		err error
	)
	// 检查option中的密钥配置。如果HS256密钥和RS256密钥（私钥和公钥）都为空，则触发一个panic，表示没有提供有效的JWT根令牌密钥
	if option.HS256Key == "" && option.RS256Key.PrivateKey == "" && option.RS256Key.PublicKey == "" {
		panic(ErrJwtRootTokenNone)
	}
	claims := &CustomJwtClaims{
		Operator: option.Operator,
		Extra:    option.Extra,
		logger:   logging.GetAppLogger(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(option.Expire)),
			IssuedAt:  jwt.At(time.Now()),
			Issuer:    JwtIssuer,
		},
	}

	// 设置HS256密钥
	if option.HS256Key != "" {
		claims.HS256Key = option.HS256Key
	}

	// 解析并设置RS256私钥
	if option.RS256Key.PrivateKey != "" {
		claims.RS256PriKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(option.RS256Key.PrivateKey))
		if err != nil {
			return nil, err
		}
	}

	// 解析并设置RS256公钥
	if option.RS256Key.PublicKey != "" {
		claims.RS256PubKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(option.RS256Key.PublicKey))
		if err != nil {
			return nil, err
		}
	}

	return claims, nil
}

// GenerateHS256JwtToken 生成 JWT-token HS256 令牌字符串
func (c *CustomJwtClaims) GenerateHS256JwtToken() (string, error) {
	// 创建了一个新的 JWT Token，并指定了签名方法为 HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 生成签名字符串
	ss, err := token.SignedString([]byte(c.HS256Key))
	return ss, err
}

// GenerateRsaJwtToken 生成 RsaJwtToken
func (c *CustomJwtClaims) GenerateRsaJwtToken() (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, c).SignedString(c.RS256PriKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetOperator 获得 Operator
func (c *CustomJwtClaims) GetOperator() string {
	return c.Operator
}

// ParseHS256JwtToken 解析一个 JWT（JSON Web Token），并将其转换为 CustomJwtClaims 类型的声明
func (c *CustomJwtClaims) ParseHS256JwtToken(jwtStr string) (*CustomJwtClaims, error) {
	token, err := jwt.ParseWithClaims(jwtStr, &CustomJwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.HS256Key), nil
	})
	if err != nil {
		return nil, err
	}
	// 验证 Token 有效性并转换声明类型
	if claims, ok := token.Claims.(*CustomJwtClaims); ok && token.Valid {
		return claims, err
	} else {
		return nil, err
	}
}

// ParseRsaJwtToken 解析RsaJwtToken
func (c *CustomJwtClaims) ParseRsaJwtToken(jwtStr string) (*CustomJwtClaims, error) {
	var (
		key = c.RS256PubKey
	)

	token, err := jwt.ParseWithClaims(jwtStr, &CustomJwtClaims{}, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			// 判断签名算法是否一致
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := token.Claims.(*CustomJwtClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("validate: invalid")
	}

	return claims, nil
}
