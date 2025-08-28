package middleware

import (
	"strconv"
	"time"

	"github.com/fengzhongzhu1621/xgo/ginx/utils"
	"github.com/fengzhongzhu1621/xgo/monitor/metrics"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		clientID := utils.GetClientID(c)
		status := strconv.Itoa(c.Writer.Status())

		e := "0"
		if _, hasError := utils.GetError(c); hasError {
			e = "1"
		}

		metrics.RequestCount.With(prometheus.Labels{
			"method":    c.Request.Method,
			"path":      c.Request.URL.Path,
			"status":    status,
			"error":     e,
			"client_id": clientID,
		}).Inc()

		metrics.RequestDuration.With(prometheus.Labels{
			"method":    c.Request.Method,
			"path":      c.Request.URL.Path,
			"status":    status,
			"client_id": clientID,
		}).Observe(float64(duration / time.Millisecond))
	}
}
