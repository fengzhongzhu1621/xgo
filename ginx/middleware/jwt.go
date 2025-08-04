package middleware

import (
	"fmt"

	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/ginx/utils"
	"github.com/fengzhongzhu1621/xgo/network/nethttp/auth/jwt"
	"github.com/gin-gonic/gin"
)

// NewClientAuthMiddleware 从配置文件获取网关公钥，并根据公钥从 jwt token 解析用户信息
func NewClientAuthMiddleware(c *config.Config) gin.HandlerFunc {
	var apiGatewayPublicKey []byte
	// 从配置文件获取蓝鲸网关公钥
	apigwCrypto, ok := c.Cryptos["apigateway_public_key"]
	if ok {
		apiGatewayPublicKey = []byte(apigwCrypto.Key)
	}

	// 从 jwt 解析用户信息
	return ClientAuthMiddleware(apiGatewayPublicKey)
}

// ClientAuthMiddleware 根据公钥从 jwt token 解析用户信息
func ClientAuthMiddleware(apiGatewayPublicKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		var clientID, clientUsername string

		// 1. 获得 jwt
		jwtToken := c.GetHeader("X-Bkapi-JWT")
		if jwtToken == "" {
			utils.UnauthorizedJSONResponse(
				c,
				"request from apigateway jwt token should not be empty!",
			)
			c.Abort()
			return
		}
		if len(apiGatewayPublicKey) == 0 {
			utils.UnauthorizedJSONResponse(
				c,
				"apigateway public key is not configured, not support request from apigateway",
			)
			c.Abort()
			return
		}

		// 2. 从 jwt 中解析出 app_code 和 username
		var err error
		clientID, clientUsername, err = jwt.GetClientIDFromJWTToken(jwtToken, apiGatewayPublicKey)
		if err != nil {
			message := fmt.Sprintf("request from apigateway jwt token invalid! err=%s", err.Error())
			utils.UnauthorizedJSONResponse(c, message)
			c.Abort()
			return
		}

		// 3. set client_id
		utils.SetClientID(c, clientID)
		utils.SetClientUsername(c, clientUsername)

		c.Next()
	}
}
