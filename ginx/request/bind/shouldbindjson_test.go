package bind

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestData struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

func TestShouldBindJSON() {
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