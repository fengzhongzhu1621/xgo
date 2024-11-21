package middleware

import (
	"net/http"

	"github.com/fengzhongzhu1621/xgo/ginx/utils"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

// UseAccessControl 实现基于用户 ID 的简单访问控制。
// 它检查当前请求的用户 ID 是否在允许访问的用户列表中，如果在，则允许请求继续处理；如果不在，则拒绝请求并返回 403 状态码。
// allowedUsers 允许访问受保护资源的用户 ID
func UseAccessControl(allowedUsers []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果 allowedUsers 切片为空，意味着没有特定的用户被限制访问。
		// 在这种情况下，中间件应该直接调用 c.Next()，将控制权传递给下一个中间件或处理函数，并立即返回。
		if len(allowedUsers) == 0 {
			c.Next()
			return
		}

		// 获得登录用户的名称
		userID := utils.GetUserID(c)

		// allowedUsers 切片中是否包含当前用户的 ID。如果包含，说明当前用户是被允许访问的，中间件应该调用 c.Next() 将控制权传递给下一个中间件或处理函数，并立即返回。
		if lo.Contains(allowedUsers, userID) {
			c.Next()
			return
		}

		c.AbortWithStatus(http.StatusForbidden)
	}
}
