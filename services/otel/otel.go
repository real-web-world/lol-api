package otel

import (
	"context"
	"errors"
	"time"

	apiProj "github.com/real-web-world/lol-api"
	"github.com/real-web-world/lol-api/global"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

func InitOtel(ctx context.Context) (shutdown func(ctx context.Context) error, err error) {
	var shutdownFunctions []func(context.Context) error
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFunctions {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFunctions = nil
		return err
	}
	otel.SetTextMapPropagator(newPropagator())

	tracerProvider, err := newTraceProvider(ctx)
	if err != nil {
		return
	}
	shutdownFunctions = append(shutdownFunctions, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)
	return
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	cfg := global.Conf
	sampler := sdktrace.NeverSample()
	if global.IsProdMode() {
		sampler = sdktrace.AlwaysSample()
	}
	res, err := newResource()
	if err != nil {
		return nil, err
	}
	retryCfg := otlptracehttp.RetryConfig{
		Enabled:         true,
		InitialInterval: time.Second,
		MaxInterval:     time.Second * 5,
		MaxElapsedTime:  30 * time.Minute,
	}
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(cfg.Otel.Endpoint),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithRetry(retryCfg),
	)
	if err != nil {
		return nil, err
	}
	bsp := sdktrace.NewBatchSpanProcessor(exporter, sdktrace.WithBatchTimeout(time.Second))
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	return tracerProvider, nil
}
func newResource() (*resource.Resource, error) {
	cfg := global.Conf
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(cfg.ProjectName),
			semconv.ServiceVersion(apiProj.APIVersion),
			attribute.String("buff.commitID", apiProj.Commit),
			attribute.String("buff.mode", global.Conf.Mode),
		))

}
