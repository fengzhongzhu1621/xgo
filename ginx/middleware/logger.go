package middleware

import (
	"bytes"
	"fmt"
	"time"

	"github.com/TencentBlueKing/gopkg/stringx"
	"github.com/fengzhongzhu1621/xgo/ginx/utils"
	"github.com/fengzhongzhu1621/xgo/logging"
	"github.com/fengzhongzhu1621/xgo/monitor/sentry"
	"github.com/fengzhongzhu1621/xgo/network/nethttp"
	"github.com/fengzhongzhu1621/xgo/str/slice"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write ...
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// AppLogger ...
func AppLogger() gin.HandlerFunc {
	logger := logging.GetAppLogger()

	return func(c *gin.Context) {
		fields := logContextFields(c)
		logger.Info("-", fields...)
	}
}

func logContextFields(c *gin.Context) []zap.Field {
	start := time.Now()

	var body string
	// 获取响应结果 response body
	requestBody, err := nethttp.ReadRequestBody(c.Request)
	if err != nil {
		body = ""
	} else {
		// 截断响应结果，防止 body 过大
		body = slice.TruncateBytesToString(requestBody, 1024)
	}

	newWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = newWriter

	c.Next()

	// 获得响应耗时
	duration := time.Since(start)
	// always add 1ms, in case the 0ms in log
	latency := float64(duration/time.Millisecond) + 1

	// 判断请求是否报错
	e, hasError := utils.GetError(c)
	if !hasError {
		e = ""
	}

	// 获得输入参数（需要截断）
	params := stringx.Truncate(c.Request.URL.RawQuery, 1024)

	// 构造日志字段
	fields := []zap.Field{
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("params", params),
		zap.String("body", body),
		zap.Int("status", c.Writer.Status()),
		zap.Float64("latency", latency),
		zap.String("request_id", utils.GetRequestID(c)),
		zap.String("client_id", utils.GetClientID(c)),
		zap.String("client_ip", c.ClientIP()),
		zap.Any("error", e),
	}
	if hasError {
		fields = append(fields, zap.String("response_body", newWriter.body.String()))
	} else {
		fields = append(fields, zap.String("response_body", stringx.Truncate(newWriter.body.String(), 1024)))
	}

	// 发送 sentry 报告
	if hasError && e != nil {
		sentry.ReportToSentry(
			fmt.Sprintf("%s %s error", c.Request.Method, c.Request.URL.Path),
			map[string]interface{}{
				"fields": fields,
			},
		)
	}

	return fields
}
