package opentelemetry

import (
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/fengzhongzhu1621/xgo/config/server_option"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/emicklei/go-restful/otelrestful"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// AddOtlpFilter add OpenTelemetry Protocol filter
func AddOtlpFilter(container *restful.Container) {
	if container != nil && openTelemetryCfg.enable {
		container.Filter(otelrestful.OTelFilter(serviceName()))
	}
}

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
	return fmt.Sprintf("%s_%s", "xgo", server_option.GetIdentification())
}
