package telemetry

import "time"

// Config holds the configuration for a telemetry Provider.
type Config struct {
	// ServiceName is the name of the instrumented service.
	ServiceName string
	// ServiceVersion is the version of the instrumented service.
	ServiceVersion string
	// Environment is the deployment environment (e.g. "production", "staging").
	Environment string
	// Endpoint is the OTLP gRPC collector endpoint (e.g. "localhost:4317").
	// When empty the provider runs in noop mode and no data is exported.
	Endpoint string
	// Insecure controls whether the gRPC connection uses TLS.
	Insecure bool
	// Headers are metadata headers sent with each OTLP export request.
	Headers map[string]string
	// SamplingRatio is the fraction of traces to sample [0.0, 1.0].
	// 1.0 means all traces are sampled (default), 0.0 means none.
	SamplingRatio float64
	// ExportInterval is how often metrics are pushed to the collector.
	ExportInterval time.Duration
}

// Option is a functional option for configuring a Provider.
type Option func(*Config)

func defaultConfig() Config {
	return Config{
		ServiceName:    "unknown",
		SamplingRatio:  1.0,
		ExportInterval: 30 * time.Second,
	}
}

// WithServiceName sets the service name reported in telemetry data.
func WithServiceName(name string) Option {
	return func(c *Config) { c.ServiceName = name }
}

// WithServiceVersion sets the service version reported in telemetry data.
func WithServiceVersion(version string) Option {
	return func(c *Config) { c.ServiceVersion = version }
}

// WithEnvironment sets the deployment environment reported in telemetry data.
func WithEnvironment(env string) Option {
	return func(c *Config) { c.Environment = env }
}

// WithEndpoint sets the OTLP gRPC collector endpoint and whether to use an
// insecure (plaintext) connection.
func WithEndpoint(endpoint string, insecure bool) Option {
	return func(c *Config) {
		c.Endpoint = endpoint
		c.Insecure = insecure
	}
}

// WithHeaders sets metadata headers sent with every OTLP export request.
func WithHeaders(headers map[string]string) Option {
	return func(c *Config) { c.Headers = headers }
}

// WithSamplingRatio sets the fraction of traces to sample in the range
// [0.0, 1.0]. Values outside that range are clamped to 0 or 1.
func WithSamplingRatio(ratio float64) Option {
	return func(c *Config) { c.SamplingRatio = ratio }
}

// WithExportInterval sets how often metric data is pushed to the collector.
func WithExportInterval(d time.Duration) Option {
	return func(c *Config) { c.ExportInterval = d }
}
