package param

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func TestQuery() {
	router := gin.Default()

	// 使用GET请求并获取参数
	router.GET("/user", func(c *gin.Context) {
		// 使用Query获取请求参数
		username := c.Query("username")
		password := c.Query("password")

		// 返回请求参数
		c.JSON(200, gin.H{
			"username": username,
			"password": password,
		})
	})

	// 获取query参数示例：GET /user?uid=20&name=jack&page=1
	router.GET("/user", func(c *gin.Context) {
		// 获取参数
		// Query获取参数
		uid := c.Query("uid")
		username := c.Query("name")
		// DefaultQuery获取参数，可以设置默认值：也就是如果没有该参数，则使用默认值
		page := c.DefaultQuery("page", "1")

		// 返回JSON结果
		c.JSON(http.StatusOK, gin.H{
			"uid":      uid,
			"username": username,
			"page":     page,
		})
	})

	// 运行服务
	router.Run(":8080")
}
