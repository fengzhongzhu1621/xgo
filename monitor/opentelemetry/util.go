package opentelemetry

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// UseOtlpMiddleware use OpenTelemetry Protocol middleware
func UseOtlpMiddleware(ws *gin.Engine) {
	if ws != nil && openTelemetryCfg.enable {
		ws.Use(otelgin.Middleware(serviceName()))
	}
}

// WrapperTraceClient wrapper client to record trace
func WrapperTraceClient(client *http.Client) {
	if client != nil && openTelemetryCfg.enable {
		client.Transport = otelhttp.NewTransport(client.Transport)
	}
}

func serviceName() string {
	return "xgo"
}
