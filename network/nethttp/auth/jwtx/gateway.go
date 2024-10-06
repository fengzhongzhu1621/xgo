package jwtx

import (
	"fmt"
	"strings"

	"github.com/fengzhongzhu1621/xgo"
	"github.com/fengzhongzhu1621/xgo/cache/cacheimpls"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// GetJwtTokenFromHeader 从请求头获取 jwt 的值
func GetJwtTokenFromHeader(c *gin.Context) (string, error) {
	// 找不到认证头部
	header := c.Request.Header.Get("Authorization")
	if len(header) == 0 {
		return "", xgo.JwtTokenNoneErr
	}
	// 获得 Bearer 的值
	strs := strings.SplitN(header, " ", 2)
	if len(strs) != 2 || strs[0] != "Bearer" {
		return "", xgo.JwtTokenInvalidErr
	}

	return strs[1], nil
}

// GetClientIDFromJWTToken 从 jwt 中解析出 app_code 和 username
func GetClientIDFromJWTToken(jwtToken string, apiGatewayPublicKey []byte) (clientID, clientUsername string, err error) {
	// 尝试从缓存中获取用户 ID 和名称
	clientID, clientUsername, err = cacheimpls.GetJWTTokenClientIDAndUsername(jwtToken)
	if err == nil {
		return
	}

	// 从 jwt token 解析出 app_code 和 username
	clientID, clientUsername, err = VerifyClientAndUsername(jwtToken, apiGatewayPublicKey)
	if err != nil {
		return "", "", err
	}
	cacheimpls.SetJWTTokenClientIDAndUsername(jwtToken, clientID, clientUsername)
	return
}

// VerifyClientAndUsername 从 jwt token 解析出 app_code 和 username
func VerifyClientAndUsername(jwtToken string, publicKey []byte) (clientID, userName string, err error) {
	var (
		claims jwt.MapClaims
	)
	// 根据公钥从 jwt token 中解码出用户信息
	claims, err = parseBKJWTToken(jwtToken, publicKey)
	if err != nil {
		return
	}

	appInfo, ok := claims["app"]
	if !ok {
		err = ErrAPIGatewayJWTMissingApp
		return
	}

	app, ok := appInfo.(map[string]interface{})
	if !ok {
		err = ErrAPIGatewayJWTAppInfoParseFail
		return
	}

	verifiedRaw, ok := app["verified"]
	if !ok {
		err = ErrAPIGatewayJWTAppInfoNoVerified
		return
	}

	verified, ok := verifiedRaw.(bool)
	if !ok {
		err = ErrAPIGatewayJWTAppInfoVerifiedNotBool
		return
	}

	if !verified {
		err = ErrAPIGatewayJWTAppNotVerified
		return
	}

	appCode, ok := app["app_code"]
	if !ok {
		err = ErrAPIGatewayJWTAppInfoNoAppCode
		return
	}

	clientID, ok = appCode.(string)
	if !ok {
		err = ErrAPIGatewayJWTAppCodeNotString
		return
	}

	usernameMap, ok := claims["user"]
	if !ok {
		err = ErrAPIGatewayJWTAppInfoNoUsername
		return
	}

	usernameSMap, ok := usernameMap.(map[string]interface{})
	if !ok {
		err = ErrAPIGatewayJWTAppInfoNoUsername
		return
	}

	userName, ok = usernameSMap["username"].(string)
	if !ok {
		err = ErrAPIGatewayJWTAppInfoNoUsername
		return
	}

	return clientID, userName, nil
}

// parseBKJWTToken 根据公钥从 jwt token 中解码出用户信息
func parseBKJWTToken(tokenString string, publicKey []byte) (jwt.MapClaims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		pubKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
		if err != nil {
			//  jwt parse fail, err=invalid key: Key must be a PEM encoded PKCS1 or PKCS8 key
			return pubKey, fmt.Errorf("jwt parse fail, err=%w", err)
		}
		return pubKey, nil
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		if verr, ok := err.(*jwt.ValidationError); ok {
			switch {
			case verr.Errors&jwt.ValidationErrorExpired > 0:
				return nil, ErrExpired
			case verr.Errors&jwt.ValidationErrorIssuedAt > 0:
				return nil, ErrIATInvalid
			case verr.Errors&jwt.ValidationErrorNotValidYet > 0:
				return nil, ErrNBFInvalid
			}
		}
		return nil, err
	}

	if !token.Valid {
		return nil, ErrUnauthorized
	}

	return claims, nil
}
