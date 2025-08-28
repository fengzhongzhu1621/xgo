package get

import "github.com/gin-gonic/gin"

func main() {
	// 创建一个默认配置的 Gin 引擎
	r := gin.Default()

	// 路由参数示例
	r.GET("/user/:name", func(c *gin.Context) {
		// 从路由路径中获取参数 :name
		name := c.Param("name")

		// 返回 HTTP 200 状态码和格式化响应
		c.String(200, "Hello %s", name)
	})

	// 查询字符串参数示例
	r.GET("/welcome", func(c *gin.Context) {
		// 获取查询参数 firstname，如果没有则使用默认值 "Guest"
		firstname := c.DefaultQuery("firstname", "Guest")

		// 获取查询参数 lastname，如果没有则返回空字符串
		lastname := c.Query("lastname")

		// 返回 HTTP 200 状态码和格式化响应
		c.String(200, "Hello %s %s", firstname, lastname)
	})

	// 启动 Web 服务器，监听 8080 端口
	r.Run(":8080")
}
