package opentelemetry

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// InitTracer init tracer to report trace information
func InitTracer(ctx context.Context) error {
	if !openTelemetryCfg.enable {
		return nil
	}

	// 配置http上报地址
	option := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(openTelemetryCfg.endpoint),
		otlptracehttp.WithInsecure(),
	}

	if openTelemetryCfg.tlsConf != nil {
		option = append(option, otlptracehttp.WithTLSClientConfig(openTelemetryCfg.tlsConf))
	}
	traceExporter, err := otlptracehttp.New(ctx, option...)
	if err != nil {
		return err
	}

	// 设置resource配置 服务名称 bk_data_id
	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName()),
	)

	// 初始化Trace配置
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(resources),
	)

	// 注入trace配置
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{},
		propagation.Baggage{}))

	// 开启协程在应用结束时关闭tracer
	go func() {
		select {
		case <-ctx.Done():
			if err := tp.Shutdown(ctx); err != nil {
				log.Errorf("Error shutting down tracer provider: %v", err)
			}
		}
	}()
	return nil
}
