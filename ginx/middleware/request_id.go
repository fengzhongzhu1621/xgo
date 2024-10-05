package middleware

import (
	"encoding/hex"

	"github.com/fengzhongzhu1621/xgo/ginx/constant"
	"github.com/fengzhongzhu1621/xgo/ginx/utils"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
)

// RequestID 用于处理HTTP请求并为其生成或提取一个唯一的X-Request-ID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debug("Middleware: RequestID")

		// 尝试从HTTP请求头中获取X-Request-ID
		requestID := c.GetHeader(constant.RequestIDHeaderKey)
		// 如果该头不存在或其长度不是32（通常UUID的长度），则生成一个新的UUID，并将其转换为十六进制字符串
		if requestID == "" || len(requestID) != 32 {
			requestID = hex.EncodeToString(uuid.Must(uuid.NewV4()).Bytes())
		}

		// 设置Request ID，将生成的或提取的X-Request-ID存储在Gin的上下文中，以便后续的处理函数可以访问它
		utils.SetRequestID(c, requestID)

		// 将该ID设置回HTTP响应头中，这样客户端也可以接收到这个ID
		c.Writer.Header().Set(constant.RequestIDHeaderKey, requestID)

		c.Next()
	}
}
