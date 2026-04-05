<div align="center">
  <img width="100" height="100" src="https://img.icons8.com/clouds/200/graph.png" alt="telemetry"/>
  <h3 align="center">Telemetry</h3>
  <h4> <a href="../../README.md" align="center"> github.com/auvitly/go-tools/standard/sdk </a> > <b>telemetry</b></h4>
  <p align="center">Observe your services with ease!</p>
  <img src="https://img.shields.io/badge/version-0.1.0-yellow?style=for-the-badge" alt="version">
</div>

---

Package `telemetry` provides a thin wrapper around the [OpenTelemetry Go SDK](https://opentelemetry.io/docs/languages/go/) to simplify setting up distributed tracing and metrics for a service.

A single `Provider` manages both a `TracerProvider` and a `MeterProvider`. When no OTLP endpoint is configured the package falls back to noop implementations so that instrumented code remains valid even in environments where telemetry is not needed.

---

## Installation

```bash
go get github.com/auvitly/go-tools/standard/sdk/telemetry
```

---

## Quick start

```go
package main

import (
    "context"
    "log"

    "github.com/auvitly/go-tools/standard/sdk/telemetry"
)

func main() {
    ctx := context.Background()

    p, err := telemetry.New(ctx,
        telemetry.WithServiceName("my-service"),
        telemetry.WithServiceVersion("1.0.0"),
        telemetry.WithEnvironment("production"),
        telemetry.WithEndpoint("otel-collector:4317", false),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer p.Shutdown(ctx)

    // Tracing
    tracer := p.Tracer()
    ctx, span := tracer.Start(ctx, "my-operation")
    defer span.End()

    // Metrics
    meter := p.Meter()
    counter, _ := meter.Int64Counter("requests_total")
    counter.Add(ctx, 1)
}
```

---

## Configuration options

| Option | Default | Description |
|---|---|---|
| `WithServiceName(name)` | `"unknown"` | Name of the instrumented service |
| `WithServiceVersion(version)` | `""` | Version of the service |
| `WithEnvironment(env)` | `""` | Deployment environment (e.g. `production`) |
| `WithEndpoint(addr, insecure)` | `""` (noop) | OTLP gRPC collector address |
| `WithHeaders(headers)` | `nil` | Metadata headers sent with export requests |
| `WithSamplingRatio(ratio)` | `1.0` | Fraction of traces to sample `[0.0, 1.0]` |
| `WithExportInterval(d)` | `30s` | How often metrics are pushed to the collector |

When `WithEndpoint` is not provided (or the address is empty) the Provider runs in **noop mode** – all spans and metric instruments are valid but no data is exported.
