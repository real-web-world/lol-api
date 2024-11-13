package trace

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"

	app "github.com/real-web-world/lol-api"
	"github.com/real-web-world/lol-api/conf"
	"github.com/real-web-world/lol-api/global"
)

func newJaegerTracer(serviceName, jaegerAddr string) (*tracesdk.TracerProvider, error) {
	ctx := context.Background()
	exp, err := jaeger.New(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerAddr)),
	)
	if err != nil {
		return nil, err
	}
	var sampler tracesdk.Sampler
	if global.IsDevMode() {
		sampler = tracesdk.AlwaysSample()
	} else {
		sampler = tracesdk.TraceIDRatioBased(1)
	}
	res, _ := resource.New(ctx,
		resource.WithAttributes(
			semconv.HTTPSchemeKey.String(semconv.SchemaURL),
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(app.Commit),
			attribute.String("environment", global.GetEnv()),
		),
	)
	bsp := tracesdk.NewBatchSpanProcessor(exp,
		tracesdk.WithBatchTimeout(time.Second))
	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(sampler),
		tracesdk.WithResource(res),
		tracesdk.WithSpanProcessor(bsp),
	)
	return tracerProvider, nil
}
func InitJeager(cfg *conf.TraceConf) error {
	serverName := cfg.ServerName
	jaegerAddr := cfg.JaegerAddr
	tp, err := newJaegerTracer(serverName, jaegerAddr)
	if err != nil {
		return err
	}
	global.SetCleanup(global.JaegerCleanupKey, func() error {
		c, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		return tp.Shutdown(c)
	})
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return nil
}
