package middleware

import (
	"fmt"
	"net/http"

	"github.com/fengzhongzhu1621/xgo/ginx/constant"
	"github.com/fengzhongzhu1621/xgo/ginx/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

// UserAuth 用户身份认证中间件
func TokenAuth(username string, loginUrl string, getUserInfo func(string) (string, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 构造重定向地址
		scheme := lo.Ternary(c.Request.TLS != nil, "https", "http")
		referUrl := fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, c.Request.RequestURI)

		// 获取token 的值
		userToken, err := c.Request.Cookie(constant.UserTokenKey)
		if err != nil {
			c.HTML(http.StatusUnauthorized, "401.html", gin.H{"loginUrl": fmt.Sprintf("%s?login_url=%s", loginUrl, referUrl)})
			c.Abort()
			return
		}
		token := userToken.Value

		// 在 session 中匹配，匹配成功则 session 中获取用户
		session := sessions.Default(c)
		if token == session.Get(constant.UserTokenKey) {
			utils.SetUserID(c, session.Get(constant.UserIDKey).(string))
			c.Next()
			return
		}

		userId, err := getUserInfo(token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		utils.SetUserID(c, userId)
		session.Set(constant.UserTokenKey, token)
		session.Set(constant.UserIDKey, userId)
		_ = session.Save()

		c.Next()
	}
}
