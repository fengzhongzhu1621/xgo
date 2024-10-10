package param

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TestParam() {
	router := gin.Default()

	// 使用URL参数并获取参数
	router.GET("/user/:username/:password", func(c *gin.Context) {
		// 使用Param获取URL参数
		username := c.Param("username")
		password := c.Param("password")

		// 返回请求参数
		c.JSON(200, gin.H{
			"username": username,
			"password": password,
		})
	})

	// GET 获取path路径参数
	router.GET("/book/:bid", func(c *gin.Context) {
		// 获取path参数
		bid := c.Param("bid")
		// 返回响应信息
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("bid=%s", bid),
		})
	})

	// 运行服务
	router.Run(":8080")
}
