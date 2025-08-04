package bind

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RequestData struct {
	Name  string `json:"name"  binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

func TestShouldBindJSON(t *testing.T) {
	r := gin.Default()

	r.POST("/submit", func(c *gin.Context) {
		var requestData RequestData
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 处理数据
		c.JSON(http.StatusOK, gin.H{"message": "Data received", "data": requestData})
	})

	r.Run(":8080")
}

func TestShouldBindJSON2(t *testing.T) {
	r := gin.Default()

	r.POST("/login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			// 返回校验错误
			errs := err.(validator.ValidationErrors)
			c.JSON(http.StatusBadRequest, gin.H{"error": errs.Error()})
			return
		}

		// 校验通过
		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
	})

	r.Run()
}
