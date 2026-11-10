package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

const (
	testTraceID     = "0af7651916cd43dd8448eb211c80319c"
	testTraceparent = "00-" + testTraceID + "-b7ad6b7169203331-01"
	testBaggage     = "tenant=acme"
)

// otelTestEnv clears every variable the setup and the exporters read, so that
// the environment the tests run in cannot steer them, then sets env. It also puts the
// OpenTelemetry globals back afterwards.
func otelTestEnv(t *testing.T, env map[string]string) {
	t.Helper()
	for _, kv := range os.Environ() {
		if k, _, _ := strings.Cut(kv, "="); strings.HasPrefix(k, "OTEL_") || k == "TRACEPARENT" || k == "TRACESTATE" || k == "BAGGAGE" {
			clearEnv(t, k)
		}
	}
	for k, v := range env {
		t.Setenv(k, v)
	}
	provider, propagator := otel.GetTracerProvider(), otel.GetTextMapPropagator()
	t.Cleanup(func() {
		_ = otelAfter(context.Background(), nil)
		// Setting a global to the value it already has logs an error.
		if otel.GetTracerProvider() != provider {
			otel.SetTracerProvider(provider)
		}
		if !assert.ObjectsAreEqual(propagator, otel.GetTextMapPropagator()) {
			otel.SetTextMapPropagator(propagator)
		}
	})
}

// shortOtelShutdown keeps a test whose collector never answers from waiting the
// full flush timeout.
func shortOtelShutdown(t *testing.T) {
	t.Helper()
	saved := otelShutdownTimeout
	otelShutdownTimeout = 500 * time.Millisecond
	t.Cleanup(func() { otelShutdownTimeout = saved })
}

// stringValue is how OTLP encodes s as a whole string field: its length, then its bytes.
func stringValue(s string) string {
	return string(rune(len(s))) + s
}

// cliScope is how OTLP encodes the CLI's instrumentation scope: its name, then its version.
var cliScope = "\x0a" + stringValue("com.anthropic.sdk.cli") + "\x12" + stringValue(Version)

// otlpServer is an OTLP/HTTP collector that keeps the bodies posted to /v1/traces.
type otlpServer struct {
	*httptest.Server
	mu     sync.Mutex
	bodies [][]byte
}

func newOTLPServer(t *testing.T) *otlpServer {
	t.Helper()
	s := &otlpServer{}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/traces" {
			http.NotFound(w, r)
			return
		}
		body, _ := io.ReadAll(r.Body)
		s.mu.Lock()
		s.bodies = append(s.bodies, body)
		s.mu.Unlock()
		w.Header().Set("Content-Type", "application/x-protobuf")
	}))
	t.Cleanup(s.Close)
	return s
}

// exported returns everything posted so far. Span and resource strings are
// readable in it: protobuf stores them as plain bytes and the exporter does not
// compress by default.
func (s *otlpServer) exported() []byte {
	s.mu.Lock()
	defer s.mu.Unlock()
	return bytes.Join(s.bodies, nil)
}

func TestOtelInactive(t *testing.T) {
	for _, tc := range []struct {
		name string
		env  map[string]string
	}{
		{name: "no endpoint"},
		{name: "no endpoint, other variables set", env: map[string]string{"OTEL_SERVICE_NAME": "svc", "OTEL_TRACES_EXPORTER": "otlp", "TRACEPARENT": testTraceparent, "BAGGAGE": testBaggage}},
		{name: "sdk disabled", env: map[string]string{"OTEL_EXPORTER_OTLP_ENDPOINT": "http://127.0.0.1:1", "OTEL_SDK_DISABLED": "true"}},
		{name: "sdk disabled, upper case", env: map[string]string{"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT": "http://127.0.0.1:1", "OTEL_SDK_DISABLED": "TRUE"}},
		{name: "traces exporter none", env: map[string]string{"OTEL_EXPORTER_OTLP_ENDPOINT": "http://127.0.0.1:1", "OTEL_TRACES_EXPORTER": "none"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			otelTestEnv(t, tc.env)
			provider, propagator := otel.GetTracerProvider(), otel.GetTextMapPropagator()

			ctx, err := otelBefore(context.Background(), nil)

			require.NoError(t, err)
			assert.Nil(t, tracerProvider)
			assert.Equal(t, provider, otel.GetTracerProvider())
			assert.Equal(t, propagator, otel.GetTextMapPropagator())
			assert.False(t, trace.SpanContextFromContext(ctx).IsValid(), "TRACEPARENT must be ignored while tracing is off")
			assert.Zero(t, baggage.FromContext(ctx).Len(), "BAGGAGE must be ignored while tracing is off")
			require.NoError(t, otelAfter(ctx, nil))
		})
	}
}

func TestOtelActive(t *testing.T) {
	for _, tc := range []struct {
		name string
		env  func(endpoint string) map[string]string
	}{
		{name: "endpoint", env: func(endpoint string) map[string]string {
			return map[string]string{"OTEL_EXPORTER_OTLP_ENDPOINT": endpoint}
		}},
		{name: "traces endpoint", env: func(endpoint string) map[string]string {
			return map[string]string{"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT": endpoint + "/v1/traces"}
		}},
		{name: "sdk not disabled, exporter otlp", env: func(endpoint string) map[string]string {
			return map[string]string{"OTEL_EXPORTER_OTLP_ENDPOINT": endpoint, "OTEL_SDK_DISABLED": "false", "OTEL_TRACES_EXPORTER": "otlp", "OTEL_EXPORTER_OTLP_PROTOCOL": "http/protobuf"}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			collector := newOTLPServer(t)
			otelTestEnv(t, tc.env(collector.URL))

			ctx, err := otelBefore(context.Background(), nil)
			require.NoError(t, err)

			require.NotNil(t, tracerProvider)
			assert.Equal(t, cliScopeTracerProvider{inner: tracerProvider}, otel.GetTracerProvider())
			assert.False(t, trace.SpanContextFromContext(ctx).IsValid(), "no TRACEPARENT, so no parent")

			spanCtx, span := otel.Tracer("test").Start(ctx, "otel-test-span")
			headers := propagation.MapCarrier{}
			otel.GetTextMapPropagator().Inject(spanCtx, headers)
			assert.Contains(t, headers["traceparent"], span.SpanContext().TraceID().String())
			span.End()

			assert.Empty(t, collector.exported(), "the batch is only sent on shutdown here")
			require.NoError(t, otelAfter(ctx, nil))
			assert.Nil(t, tracerProvider)
			assert.Contains(t, string(collector.exported()), "otel-test-span")
		})
	}
}

func TestOtelAfterIgnoresCancelledContext(t *testing.T) {
	collector := newOTLPServer(t)
	otelTestEnv(t, map[string]string{"OTEL_EXPORTER_OTLP_ENDPOINT": collector.URL})
	ctx, cancel := context.WithCancel(context.Background())
	ctx, err := otelBefore(ctx, nil)
	require.NoError(t, err)
	_, span := otel.Tracer("test").Start(ctx, "otel-test-span")
	span.End()

	cancel()
	require.NoError(t, otelAfter(ctx, nil))

	assert.Contains(t, string(collector.exported()), "otel-test-span")
}

func TestOtelTraceparent(t *testing.T) {
	for _, tc := range []struct {
		name        string
		traceparent string
		tracestate  string
		wantTraceID string
	}{
		{name: "valid", traceparent: testTraceparent, tracestate: "vendor=value", wantTraceID: testTraceID},
		{name: "malformed is ignored", traceparent: "not-a-traceparent"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			collector := newOTLPServer(t)
			otelTestEnv(t, map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": collector.URL,
				"TRACEPARENT":                 tc.traceparent,
				"TRACESTATE":                  tc.tracestate,
			})

			ctx, err := otelBefore(context.Background(), nil)
			require.NoError(t, err)

			parent := trace.SpanContextFromContext(ctx)
			if tc.wantTraceID == "" {
				assert.False(t, parent.IsValid())
				return
			}
			assert.Equal(t, tc.wantTraceID, parent.TraceID().String())
			assert.True(t, parent.IsRemote())
			assert.True(t, parent.IsSampled())
			assert.Equal(t, tc.tracestate, parent.TraceState().String())

			_, span := otel.Tracer("test").Start(ctx, "child")
			assert.Equal(t, tc.wantTraceID, span.SpanContext().TraceID().String())
			span.End()
		})
	}
}

func TestOtelBaggage(t *testing.T) {
	for _, tc := range []struct {
		name string
		env  map[string]string
		want string
	}{
		{name: "with TRACEPARENT", env: map[string]string{"BAGGAGE": testBaggage, "TRACEPARENT": testTraceparent}, want: testBaggage},
		{name: "without TRACEPARENT", env: map[string]string{"BAGGAGE": testBaggage}, want: testBaggage},
		{name: "malformed is ignored", env: map[string]string{"BAGGAGE": "not-a-baggage"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			clearWorkerEnv(t)
			collector := newOTLPServer(t)
			otelTestEnv(t, tc.env)
			t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", collector.URL)
			t.Setenv("ANTHROPIC_CONFIG_DIR", t.TempDir())
			t.Setenv("ANTHROPIC_API_KEY", "k")
			var headers []http.Header
			api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				headers = append(headers, r.Header)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"data":[],"has_more":false,"first_id":null,"last_id":null}`)
			}))
			defer api.Close()

			require.NoError(t, Command.Run(context.Background(), []string{"ant", "--base-url", api.URL, "models", "list"}))
			api.Close()

			require.Len(t, headers, 1)
			assert.Equal(t, tc.want, headers[0].Get("baggage"))
		})
	}
}

func TestOtelServiceName(t *testing.T) {
	for _, tc := range []struct {
		name string
		env  map[string]string
		want string
	}{
		{name: "default", want: "ant"},
		{name: "OTEL_SERVICE_NAME", env: map[string]string{"OTEL_SERVICE_NAME": "from-service-name"}, want: "from-service-name"},
		{name: "OTEL_RESOURCE_ATTRIBUTES", env: map[string]string{"OTEL_RESOURCE_ATTRIBUTES": "service.name=from-attributes,team=x"}, want: "from-attributes"},
		{name: "OTEL_SERVICE_NAME beats OTEL_RESOURCE_ATTRIBUTES", env: map[string]string{"OTEL_SERVICE_NAME": "from-service-name", "OTEL_RESOURCE_ATTRIBUTES": "service.name=from-attributes"}, want: "from-service-name"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			collector := newOTLPServer(t)
			otelTestEnv(t, tc.env)
			t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", collector.URL)

			res, err := newOtelResource(context.Background())
			require.NoError(t, err)
			got, ok := res.Set().Value(attribute.Key("service.name"))
			require.True(t, ok)
			assert.Equal(t, tc.want, got.AsString())
			_, ok = res.Set().Value(attribute.Key("telemetry.sdk.language"))
			assert.True(t, ok)

			// The provider merges the environment under the resource it is given, so check what it exports too.
			ctx, err := otelBefore(context.Background(), nil)
			require.NoError(t, err)
			_, span := otel.Tracer("test").Start(ctx, "otel-test-span")
			span.End()
			require.NoError(t, otelAfter(ctx, nil))
			exported := string(collector.exported())
			assert.Contains(t, exported, stringValue(tc.want))
			if tc.want != "ant" {
				assert.NotContains(t, exported, stringValue("ant"), "the default service name must not be exported beside an override")
			}
		})
	}
}

// A collector nobody listens on and a malformed variable must cost the spans only.
func TestOtelSetupNeverFailsTheCommand(t *testing.T) {
	for _, tc := range []struct {
		name string
		env  map[string]string
	}{
		{name: "malformed resource attributes", env: map[string]string{"OTEL_EXPORTER_OTLP_ENDPOINT": "http://127.0.0.1:1", "OTEL_RESOURCE_ATTRIBUTES": "no-equals-sign"}},
		{name: "unreachable collector", env: map[string]string{"OTEL_EXPORTER_OTLP_ENDPOINT": "http://127.0.0.1:1"}},
		{name: "unreachable grpc collector", env: map[string]string{"OTEL_EXPORTER_OTLP_ENDPOINT": "http://127.0.0.1:1", "OTEL_EXPORTER_OTLP_PROTOCOL": "grpc"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			otelTestEnv(t, tc.env)
			shortOtelShutdown(t)
			ctx, err := otelBefore(context.Background(), nil)
			require.NoError(t, err)
			_, span := otel.Tracer("test").Start(ctx, "otel-test-span")
			span.End()
			require.NoError(t, otelAfter(ctx, nil))
		})
	}
}

// The grpc protocol must not reach an OTLP/HTTP collector: that would mean the
// protocol variables were ignored.
func TestOtelProtocol(t *testing.T) {
	for _, tc := range []struct {
		name     string
		env      map[string]string
		wantHTTP bool
	}{
		{name: "default is http", wantHTTP: true},
		{name: "grpc", env: map[string]string{"OTEL_EXPORTER_OTLP_PROTOCOL": "grpc"}},
		{name: "traces protocol wins", env: map[string]string{"OTEL_EXPORTER_OTLP_PROTOCOL": "grpc", "OTEL_EXPORTER_OTLP_TRACES_PROTOCOL": "http/protobuf"}, wantHTTP: true},
		{name: "traces protocol grpc", env: map[string]string{"OTEL_EXPORTER_OTLP_PROTOCOL": "http/protobuf", "OTEL_EXPORTER_OTLP_TRACES_PROTOCOL": "grpc"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			collector := newOTLPServer(t)
			otelTestEnv(t, tc.env)
			t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", collector.URL)
			if !tc.wantHTTP {
				shortOtelShutdown(t)
			}
			ctx, err := otelBefore(context.Background(), nil)
			require.NoError(t, err)
			_, span := otel.Tracer("test").Start(ctx, "otel-test-span")
			span.End()
			require.NoError(t, otelAfter(ctx, nil))
			assert.Equal(t, tc.wantHTTP, len(collector.exported()) > 0)
		})
	}
}

// TestOtelScope covers the provider that the Go SDK gets its tracers from: every
// span of an API command carries the CLI's instrumentation scope, and a tracer
// keeps the options that do not name the instrumentation.
func TestOtelScope(t *testing.T) {
	clearWorkerEnv(t)
	otelTestEnv(t, nil)
	t.Setenv("ANTHROPIC_CONFIG_DIR", t.TempDir())
	t.Setenv("ANTHROPIC_API_KEY", "k")
	recorded := tracetest.NewSpanRecorder()
	provider := cliScopeTracerProvider{inner: sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorded))}
	otel.SetTracerProvider(provider)
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"data":[],"has_more":false,"first_id":null,"last_id":null}`)
	}))
	defer api.Close()

	require.NoError(t, Command.Run(context.Background(), []string{"ant", "--base-url", api.URL, "models", "list"}))
	_, span := provider.Tracer("com.anthropic.sdk.go",
		trace.WithInstrumentationVersion("0.0.1"),
		trace.WithSchemaURL("https://example.com/schema"),
		trace.WithInstrumentationAttributes(attribute.String("key", "value")),
	).Start(context.Background(), "otel-test-span")
	span.End()

	spans := recorded.Ended()
	if len(spans) == 1 {
		t.Skip("the anthropic-sdk-go version this build links creates no spans; needs an SDK release with OpenTelemetry tracing")
	}
	for _, s := range spans {
		assert.Equal(t, "com.anthropic.sdk.cli", s.InstrumentationScope().Name, "%s", s.Name())
		assert.Equal(t, Version, s.InstrumentationScope().Version, "%s", s.Name())
	}
	assert.Equal(t, "anthropic.models.list", spans[0].Name())
	scope := spans[len(spans)-1].InstrumentationScope()
	assert.Equal(t, "https://example.com/schema", scope.SchemaURL)
	assert.Equal(t, []attribute.KeyValue{attribute.String("key", "value")}, scope.Attributes.ToSlice())
}

// TestOtelHooksOnRoot drives the production Command tree: the root's hooks must
// be the OpenTelemetry ones, the context the Before hook returns must reach a
// subcommand's action, and the After hook must flush when the action fails too.
func TestOtelHooksOnRoot(t *testing.T) {
	for _, tc := range []struct {
		name      string
		actionErr error
	}{
		{name: "action succeeds"},
		{name: "action fails", actionErr: errors.New("probe failed")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			collector := newOTLPServer(t)
			otelTestEnv(t, map[string]string{"OTEL_EXPORTER_OTLP_ENDPOINT": collector.URL, "TRACEPARENT": testTraceparent})
			t.Setenv("ANTHROPIC_CONFIG_DIR", t.TempDir())

			ran := false
			probe := &cli.Command{
				Name: "otel-probe",
				Action: func(ctx context.Context, _ *cli.Command) error {
					ran = true
					_, installed := otel.GetTracerProvider().(cliScopeTracerProvider)
					assert.True(t, installed, "the provider must be installed before the action runs")
					assert.Equal(t, testTraceID, trace.SpanContextFromContext(ctx).TraceID().String())
					_, span := otel.Tracer("test").Start(ctx, "otel-probe-span")
					span.End()
					return tc.actionErr
				},
			}
			saved := Command.Commands
			Command.Commands = append(append([]*cli.Command{}, saved...), probe)
			t.Cleanup(func() { Command.Commands = saved })

			err := Command.Run(context.Background(), []string{"ant", "otel-probe"})

			assert.Equal(t, tc.actionErr, err)
			require.True(t, ran)
			assert.Nil(t, tracerProvider, "the After hook must have shut the provider down")
			assert.Contains(t, string(collector.exported()), "otel-probe-span")
		})
	}
}

// TestOtelWorker covers the self-hosted worker, which builds its own client and
// derives its context from the one the Before hook returns. `run` serves one
// session, so its requests continue the trace from TRACEPARENT and carry the
// baggage from BAGGAGE. `poll` serves many, one after another, so its requests
// must not. Stopping either must flush its spans.
func TestOtelWorker(t *testing.T) {
	for _, tc := range []struct {
		name          string
		args          []string
		wantSpan      string
		continueTrace bool
		wantBaggage   string
	}{
		{
			name:     "poll",
			args:     []string{"beta:worker", "poll", "--environment-id", "env_t", "--environment-key", "ek_t"},
			wantSpan: "anthropic.environments.work.poll",
		},
		{
			name:          "run",
			args:          []string{"beta:worker", "run", "--environment-id", "env_t", "--environment-key", "ek_t", "--work-id", "work_t", "--session-id", "sesn_t"},
			wantSpan:      "anthropic.environments.work.heartbeat",
			continueTrace: true,
			wantBaggage:   testBaggage,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			clearWorkerEnv(t)
			collector := newOTLPServer(t)
			otelTestEnv(t, map[string]string{"OTEL_EXPORTER_OTLP_ENDPOINT": collector.URL, "TRACEPARENT": testTraceparent, "BAGGAGE": testBaggage})
			t.Setenv("ANTHROPIC_CONFIG_DIR", t.TempDir())

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			headersCh := make(chan http.Header, 1)
			api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				select {
				case headersCh <- r.Header:
					cancel()
				default:
				}
				w.WriteHeader(http.StatusServiceUnavailable)
			}))
			defer api.Close()

			_ = Command.Run(ctx, append([]string{"ant"}, append(tc.args, "--workdir", t.TempDir(), "--base-url", api.URL)...))

			var headers http.Header
			select {
			case headers = <-headersCh:
			default:
				t.Fatal("worker made no requests")
			}
			traceparent := headers.Get("traceparent")
			assert.Nil(t, tracerProvider, "the After hook must have shut the provider down")
			exported := string(collector.exported())
			if traceparent == "" && exported == "" {
				t.Skip("the anthropic-sdk-go version this build links creates no spans; needs an SDK release with OpenTelemetry tracing")
			}
			assert.Regexp(t, "^00-[0-9a-f]{32}-[0-9a-f]{16}-01$", traceparent)
			assert.Equal(t, tc.continueTrace, strings.Contains(traceparent, testTraceID))
			assert.Equal(t, tc.wantBaggage, headers.Get("baggage"))
			assert.Contains(t, exported, tc.wantSpan)
			assert.Contains(t, exported, cliScope)
			assert.NotContains(t, exported, "com.anthropic.sdk.go")
		})
	}
}

// workItemJSON is a claimed work item for the session sesn_t, as the API returns it.
const workItemJSON = `{"id":"work_t","type":"work","environment_id":"env_t","state":"queued","secret":"","metadata":{},
	"actor":{"type":"api_key","api_key_id":"apk_t"},"data":{"type":"session","id":"sesn_t","mode":"private"},
	"created_at":"2026-05-11T12:00:00Z","acknowledged_at":"2026-05-11T12:00:00Z","started_at":"2026-05-11T12:00:00Z",
	"latest_heartbeat_at":"2026-05-11T12:00:00Z","stop_requested_at":"","stopped_at":""}`

// TestOtelWorkerPollOnWork covers `poll --on-work`: each work item gets a root
// span, which the script's exit status does not fail.
func TestOtelWorkerPollOnWork(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the on-work script is a shell script")
	}
	clearWorkerEnv(t)
	// Tracing is set up here, in memory, so that the test can read the spans.
	otelTestEnv(t, map[string]string{"TRACEPARENT": testTraceparent, "TRACESTATE": "vendor=value"})
	t.Setenv("ANTHROPIC_CONFIG_DIR", t.TempDir())
	recorded := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(cliScopeTracerProvider{inner: sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorded))})
	otel.SetTextMapPropagator(propagation.TraceContext{})

	script := filepath.Join(t.TempDir(), "on-work.sh")
	require.NoError(t, os.WriteFile(script, []byte("#!/bin/sh\nexit 3\n"), 0o700))

	// The context that the Before hook makes of the worker's own TRACEPARENT.
	ctx := propagation.TraceContext{}.Extract(context.Background(), propagation.MapCarrier{"traceparent": testTraceparent})
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var mu sync.Mutex
	polls := 0
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/work/poll"):
			mu.Lock()
			polls++
			first := polls == 1
			mu.Unlock()
			if !first {
				// The first item is done with by the time the worker polls again.
				cancel()
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, workItemJSON)
		case strings.HasSuffix(r.URL.Path, "/ack"):
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, workItemJSON)
		case strings.HasSuffix(r.URL.Path, "/stop"):
			w.WriteHeader(http.StatusNoContent)
		default:
			http.NotFound(w, r)
		}
	}))
	defer api.Close()

	require.NoError(t, Command.Run(ctx, []string{"ant", "beta:worker", "poll",
		"--environment-id", "env_t", "--environment-key", "ek_t", "--base-url", api.URL, "--on-work", script}))

	var work sdktrace.ReadOnlySpan
	for _, s := range recorded.Ended() {
		assert.NotEqual(t, testTraceID, s.SpanContext().TraceID().String(), "%s must not continue the worker's own TRACEPARENT", s.Name())
		assert.Equal(t, "com.anthropic.sdk.cli", s.InstrumentationScope().Name, "%s", s.Name())
		assert.Equal(t, Version, s.InstrumentationScope().Version, "%s", s.Name())
		if s.Name() == "anthropic.environment.work" {
			work = s
		}
	}
	require.NotNil(t, work, "the work item has a span")
	assert.False(t, work.Parent().IsValid(), "the span of a work item is a root span")
	assert.Equal(t, codes.Ok, work.Status().Code, "the worker is done with the item when it asks for the next one")
}

// TestOtelEndToEnd runs the real binary against a fake API and a fake collector:
// the API request must carry the trace context and the collector must receive
// the span of the SDK method.
func TestOtelEndToEnd(t *testing.T) {
	collector := newOTLPServer(t)
	var traceparents []string
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceparents = append(traceparents, r.Header.Get("traceparent"))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"data":[],"has_more":false,"first_id":null,"last_id":null}`)
	}))
	defer api.Close()

	cmd := exec.Command("go", "run", "./cmd/ant", "--base-url", api.URL, "models", "list")
	cmd.Dir = "../.."
	// Dropped rather than overridden: a signal-specific OTLP variable in the
	// environment of the test run would otherwise win over the endpoint set here.
	env := slices.DeleteFunc(cmd.Environ(), func(kv string) bool {
		return strings.HasPrefix(kv, "OTEL_") || strings.HasPrefix(kv, "TRACESTATE=")
	})
	cmd.Env = append(env,
		"ANTHROPIC_API_KEY=k",
		"ANTHROPIC_CONFIG_DIR="+t.TempDir(),
		"OTEL_EXPORTER_OTLP_ENDPOINT="+collector.URL,
		"OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf",
		"OTEL_SDK_DISABLED=false",
		"OTEL_TRACES_EXPORTER=otlp",
		"OTEL_TRACES_SAMPLER=always_on",
		"OTEL_SERVICE_NAME=otel-end-to-end",
		"TRACEPARENT="+testTraceparent,
	)
	out, err := cmd.CombinedOutput()
	api.Close()
	require.NoError(t, err, "%s", out)
	require.Len(t, traceparents, 1, "%s", out)

	exported := string(collector.exported())
	if traceparents[0] == "" && exported == "" {
		t.Skip("the anthropic-sdk-go version this build links creates no spans; needs an SDK release with OpenTelemetry tracing")
	}
	assert.Contains(t, traceparents[0], "00-"+testTraceID+"-", "the API request must continue the trace from TRACEPARENT")
	assert.Contains(t, exported, "anthropic.models.list")
	assert.Contains(t, exported, cliScope)
	assert.NotContains(t, exported, "com.anthropic.sdk.go")
	assert.Contains(t, exported, "otel-end-to-end")
}
