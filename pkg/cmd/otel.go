package cmd

// OpenTelemetry export for the spans the Go SDK creates. The CLI creates no
// spans of its own: it only registers a tracer provider and a propagator with
// the OpenTelemetry API, which the SDK picks up, when the standard OTLP
// environment variables ask for it.

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/urfave/cli/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
)

// Set by otelBefore and flushed by otelAfter, which urfave/cli runs around the
// command's action in that order.
var tracerProvider *sdktrace.TracerProvider

// Long enough for one batch to reach a collector, short enough that an
// unreachable collector does not hold the process at exit. A variable so that
// tests of unreachable collectors can shorten it.
var otelShutdownTimeout = 5 * time.Second

// cliScopeTracerProvider is the provider the Go SDK gets its tracers from: the
// spans carry the CLI's own name and version, as its request headers do.
type cliScopeTracerProvider struct {
	embedded.TracerProvider
	inner trace.TracerProvider
}

func (p cliScopeTracerProvider) Tracer(_ string, opts ...trace.TracerOption) trace.Tracer {
	return p.inner.Tracer("com.anthropic.sdk.cli", append(slices.Clip(opts), trace.WithInstrumentationVersion(Version))...)
}

// otelEnabled reports whether the environment asks for OTLP trace export. An
// endpoint is required so that a plain `ant` never dials localhost:4318.
func otelEnabled() bool {
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") == "" && os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT") == "" {
		return false
	}
	if strings.EqualFold(os.Getenv("OTEL_SDK_DISABLED"), "true") {
		return false
	}
	return !strings.EqualFold(os.Getenv("OTEL_TRACES_EXPORTER"), "none")
}

// newOtelExporter leaves the endpoint, headers, timeout, compression and TLS
// variables to the exporters, which read OTEL_EXPORTER_OTLP_* themselves. Every
// protocol but grpc gets OTLP/HTTP with protobuf, the only other one the Go
// exporters speak.
func newOtelExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	// The traces variable wins over the general one, as for every OTLP variable.
	protocol := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_PROTOCOL")
	if protocol == "" {
		protocol = os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL")
	}
	if strings.EqualFold(protocol, "grpc") {
		return otlptracegrpc.New(ctx)
	}
	return otlptracehttp.New(ctx)
}

// newOtelResource puts the default service name first: later detectors win, so
// OTEL_SERVICE_NAME and OTEL_RESOURCE_ATTRIBUTES override it.
func newOtelResource(ctx context.Context) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(semconv.ServiceName("ant")),
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
	)
}

// otelBefore is the root command's Before hook. It never fails the command: a
// broken telemetry setup costs the spans, not the API call.
func otelBefore(ctx context.Context, _ *cli.Command) (context.Context, error) {
	if !otelEnabled() {
		return ctx, nil
	}
	exporter, err := newOtelExporter(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: OpenTelemetry tracing is off: %v\n", err)
		return ctx, nil
	}
	res, err := newOtelResource(ctx)
	if err != nil {
		// The resource is still usable: only the failing detector's attributes are missing.
		fmt.Fprintf(os.Stderr, "warning: incomplete OpenTelemetry resource: %v\n", err)
	}
	// The sampler and the batch limits come from OTEL_TRACES_SAMPLER and OTEL_BSP_*, which the OpenTelemetry SDK reads.
	tracerProvider = sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(res))
	otel.SetTracerProvider(cliScopeTracerProvider{inner: tracerProvider})
	// The API's default propagator is a no-op, which would leave requests without a traceparent header.
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)

	// A parent process hands its trace over in TRACEPARENT and its baggage in
	// BAGGAGE, so that the spans of this run join that trace. `beta:worker poll`
	// leaves both again (handleWorkerPoll).
	return propagator.Extract(ctx, propagation.MapCarrier{
		"traceparent": os.Getenv("TRACEPARENT"),
		"tracestate":  os.Getenv("TRACESTATE"),
		"baggage":     os.Getenv("BAGGAGE"),
	}), nil
}

// otelAfter is the root command's After hook, which runs whether or not the
// action failed. It sends the spans still in the batch.
func otelAfter(context.Context, *cli.Command) error {
	if tracerProvider == nil {
		return nil
	}
	// Not the command's context, which has no deadline and may already be cancelled.
	ctx, cancel := context.WithTimeout(context.Background(), otelShutdownTimeout)
	defer cancel()
	if err := tracerProvider.Shutdown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "warning: exporting OpenTelemetry spans failed: %v\n", err)
	}
	tracerProvider = nil
	return nil
}
