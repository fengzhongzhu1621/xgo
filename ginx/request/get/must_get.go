package get

import (
	"log"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMuatGet(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		example := c.MustGet("example").(string)
		log.Println(example)
	})

	r.Run(":8080")
}
