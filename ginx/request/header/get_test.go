package header

import "github.com/gin-gonic/gin"

func HeaderGet() {
	router := gin.Default()

	// 使用请求头部参数并获取参数
	router.GET("/user", func(c *gin.Context) {
		// 使用Request获取请求头部参数
		username := c.Request.Header.Get("username")
		password := c.Request.Header.Get("password")

		// 返回请求参数
		c.JSON(200, gin.H{
			"username": username,
			"password": password,
		})
	})

	// 运行服务
	router.Run(":8080")
}
