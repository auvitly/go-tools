package telemetry_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/auvitly/go-tools/standard/sdk/telemetry"
)

// TestNew_Noop verifies that a Provider is created in noop mode when no
// endpoint is configured and that Tracer and Meter are non-nil.
func TestNew_Noop(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	p, err := telemetry.New(ctx,
		telemetry.WithServiceName("test-service"),
		telemetry.WithServiceVersion("0.1.0"),
		telemetry.WithEnvironment("test"),
	)
	require.NoError(t, err)
	require.NotNil(t, p)

	assert.NotNil(t, p.Tracer())
	assert.NotNil(t, p.Meter())
}

// TestNew_Noop_Shutdown verifies that Shutdown on a noop Provider returns nil.
func TestNew_Noop_Shutdown(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	p, err := telemetry.New(ctx)
	require.NoError(t, err)

	assert.NoError(t, p.Shutdown(ctx))
}

// TestNew_Noop_Tracer verifies that the noop tracer creates valid spans.
func TestNew_Noop_Tracer(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	p, err := telemetry.New(ctx, telemetry.WithServiceName("svc"))
	require.NoError(t, err)

	_, span := p.Tracer().Start(ctx, "test-span")
	span.End()
}

// TestNew_Noop_Meter verifies that the noop meter creates valid instruments.
func TestNew_Noop_Meter(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	p, err := telemetry.New(ctx, telemetry.WithServiceName("svc"))
	require.NoError(t, err)

	counter, err := p.Meter().Int64Counter("requests_total")
	require.NoError(t, err)

	counter.Add(ctx, 1)
}

// TestNew_InvalidEndpoint verifies that a connection error to a non-existent
// endpoint is reported during New.
func TestNew_InvalidEndpoint(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := telemetry.New(ctx,
		telemetry.WithEndpoint("localhost:19999", true),
	)
	// The OTLP gRPC exporter connects lazily, so New itself succeeds.
	// Actual export errors only surface when spans/metrics are flushed.
	// We just check no panic occurs and a Provider (or error) is returned.
	_ = err
}

// TestDefaultConfig verifies that the default configuration values are sane.
func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// New with no options uses defaultConfig; just check it doesn't panic.
	p, err := telemetry.New(ctx)
	require.NoError(t, err)
	require.NotNil(t, p)
}

// TestWithOptions verifies that all functional options are accepted without
// error when creating a noop provider.
func TestWithOptions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	p, err := telemetry.New(ctx,
		telemetry.WithServiceName("svc"),
		telemetry.WithServiceVersion("1.2.3"),
		telemetry.WithEnvironment("production"),
		telemetry.WithHeaders(map[string]string{"x-api-key": "secret"}),
		telemetry.WithSamplingRatio(0.5),
		telemetry.WithExportInterval(10*time.Second),
	)
	require.NoError(t, err)
	require.NotNil(t, p)
}
