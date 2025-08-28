package bind

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

// 定义需要绑定请求参数的结构体
type User struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func TestShouldBind(t *testing.T) {
	router := gin.Default()

	// 使用ShouldBind方法自动绑定请求参数到结构体，并进行校验
	router.POST("/user", func(c *gin.Context) {
		var user User
		// 根据req的content type 自动推断如何绑定,form/json/xml等格式
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"username": user.Username,
			"password": user.Password,
		})
	})

	router.POST("/submit", func(c *gin.Context) {
		var requestData RequestData
		if err := c.ShouldBind(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 处理数据
		c.JSON(http.StatusOK, gin.H{"message": "Data received", "data": requestData})
	})

	// 运行服务
	router.Run(":8080")
}
