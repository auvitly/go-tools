// Package telemetry provides tools for application observability using
// OpenTelemetry. It wraps the OpenTelemetry SDK to offer a single Provider
// that manages both a TracerProvider and a MeterProvider, exporting data to
// an OTLP gRPC endpoint.
//
// Basic usage:
//
//	p, err := telemetry.New(ctx,
//	    telemetry.WithServiceName("my-service"),
//	    telemetry.WithServiceVersion("1.0.0"),
//	    telemetry.WithEndpoint("localhost:4317", true),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer p.Shutdown(ctx)
//
//	tracer := p.Tracer()
//	meter  := p.Meter()
package telemetry

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc/credentials"
)

// Provider manages OpenTelemetry tracing and metrics resources for a service.
// Call Shutdown when the application exits to flush any pending telemetry.
type Provider struct {
	tracer    trace.Tracer
	meter     metric.Meter
	shutdowns []func(context.Context) error
}

// New creates a new Provider configured by the supplied options.
// When no endpoint is configured the Provider runs in noop mode – all spans
// and metric instruments are valid but no data is exported.
func New(ctx context.Context, opts ...Option) (*Provider, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.Endpoint == "" {
		return newNoopProvider(cfg), nil
	}

	res, err := newResource(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("telemetry: create resource: %w", err)
	}

	tp, err := newTraceProvider(ctx, res, cfg)
	if err != nil {
		return nil, fmt.Errorf("telemetry: create trace provider: %w", err)
	}

	mp, err := newMeterProvider(ctx, res, cfg)
	if err != nil {
		_ = tp.Shutdown(ctx)
		return nil, fmt.Errorf("telemetry: create meter provider: %w", err)
	}

	return &Provider{
		tracer: tp.Tracer(cfg.ServiceName),
		meter:  mp.Meter(cfg.ServiceName),
		shutdowns: []func(context.Context) error{
			tp.Shutdown,
			mp.Shutdown,
		},
	}, nil
}

// Tracer returns the OpenTelemetry Tracer for creating spans.
func (p *Provider) Tracer() trace.Tracer { return p.tracer }

// Meter returns the OpenTelemetry Meter for creating metric instruments.
func (p *Provider) Meter() metric.Meter { return p.meter }

// Shutdown gracefully shuts down the Provider, flushing any buffered telemetry.
// It returns a combined error if any shutdown step fails.
func (p *Provider) Shutdown(ctx context.Context) error {
	var errs []error
	for _, fn := range p.shutdowns {
		if err := fn(ctx); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// newNoopProvider returns a Provider whose tracer and meter discard all data.
func newNoopProvider(cfg Config) *Provider {
	return &Provider{
		tracer: tracenoop.NewTracerProvider().Tracer(cfg.ServiceName),
		meter:  metricnoop.NewMeterProvider().Meter(cfg.ServiceName),
	}
}

// newResource builds an OpenTelemetry Resource with standard service attributes.
func newResource(ctx context.Context, cfg Config) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.DeploymentEnvironment(cfg.Environment),
		),
	)
}

// traceGRPCOptions returns OTLP gRPC options for the trace exporter.
func traceGRPCOptions(cfg Config) []otlptracegrpc.Option {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
	}

	if cfg.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	} else {
		opts = append(opts, otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")))
	}

	if len(cfg.Headers) > 0 {
		opts = append(opts, otlptracegrpc.WithHeaders(cfg.Headers))
	}

	return opts
}

// metricGRPCOptions returns OTLP gRPC options for the metric exporter.
func metricGRPCOptions(cfg Config) []otlpmetricgrpc.Option {
	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(cfg.Endpoint),
	}

	if cfg.Insecure {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	} else {
		opts = append(opts, otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")))
	}

	if len(cfg.Headers) > 0 {
		opts = append(opts, otlpmetricgrpc.WithHeaders(cfg.Headers))
	}

	return opts
}

// newTraceProvider creates an SDK TracerProvider that exports via OTLP gRPC.
func newTraceProvider(ctx context.Context, res *resource.Resource, cfg Config) (*sdktrace.TracerProvider, error) {
	exp, err := otlptracegrpc.New(ctx, traceGRPCOptions(cfg)...)
	if err != nil {
		return nil, err
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(buildSampler(cfg.SamplingRatio)),
	), nil
}

// newMeterProvider creates an SDK MeterProvider that exports via OTLP gRPC.
func newMeterProvider(ctx context.Context, res *resource.Resource, cfg Config) (*sdkmetric.MeterProvider, error) {
	exp, err := otlpmetricgrpc.New(ctx, metricGRPCOptions(cfg)...)
	if err != nil {
		return nil, err
	}

	return sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp,
			sdkmetric.WithInterval(cfg.ExportInterval),
		)),
		sdkmetric.WithResource(res),
	), nil
}

// buildSampler returns an appropriate sampler for the given ratio.
func buildSampler(ratio float64) sdktrace.Sampler {
	switch {
	case ratio >= 1.0:
		return sdktrace.AlwaysSample()
	case ratio <= 0.0:
		return sdktrace.NeverSample()
	default:
		return sdktrace.TraceIDRatioBased(ratio)
	}
}
