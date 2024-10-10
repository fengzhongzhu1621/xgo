package param

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestPostForm() {
	router := gin.Default()

	// 使用POST请求并获取参数
	router.POST("/user", func(c *gin.Context) {
		// 使用PostForm获取请求参数
		username := c.PostForm("username")
		password := c.PostForm("password")

		// 返回请求参数
		c.JSON(200, gin.H{
			"username": username,
			"password": password,
		})
	})

	// 运行服务
	router.Run(":8080")
}

func TestPostJson() {
	router := gin.Default()

	// 使用POST请求并获取参数
	router.POST("/user", func(c *gin.Context) {
		// 将JSON格式请求参数绑定到结构体上
		var user User
		if err := c.BindJSON(&user); err != nil {
			// 返回错误信息
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 返回请求参数
		c.JSON(200, gin.H{
			"username": user.Username,
			"password": user.Password,
		})
	})

	// 运行服务
	router.Run(":8080")
}
