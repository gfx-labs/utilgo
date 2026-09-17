package gotel

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestInitTracingResource(t *testing.T) {
	// InitTracing changes process-wide state, so this test must not run in parallel.
	previousProvider := otel.GetTracerProvider()
	previousPropagator := otel.GetTextMapPropagator()
	t.Cleanup(func() {
		otel.SetTracerProvider(previousProvider)
		otel.SetTextMapPropagator(previousPropagator)
	})
	t.Setenv("HOSTNAME", "test-host")
	ctx := context.Background()
	exporter := tracetest.NewInMemoryExporter()
	shutdown, err := InitTracing(ctx,
		WithServiceName("princessapi"),
		WithServiceNamespace("oku"),
		WithExporter(exporter),
	)
	if err != nil {
		t.Fatalf("InitTracing with SDK schema %s: %v", resource.Default().SchemaURL(), err)
	}
	t.Cleanup(func() {
		if err := shutdown(ctx); err != nil {
			t.Errorf("shutdown: %v", err)
		}
	})

	_, span := otel.Tracer("gotel-test").Start(ctx, "resource-test")
	span.End()
	provider := otel.GetTracerProvider().(*sdktrace.TracerProvider)
	if err := provider.ForceFlush(ctx); err != nil {
		t.Fatalf("flush: %v", err)
	}
	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("exported %d spans, want 1", len(spans))
	}
	r := spans[0].Resource
	if r.SchemaURL() != resource.Default().SchemaURL() {
		t.Errorf("schema = %q, want %q", r.SchemaURL(), resource.Default().SchemaURL())
	}
	attrs := r.Set()
	for key, want := range map[attribute.Key]string{
		"host.name":          "test-host",
		"service.name":       "princessapi",
		"service.namespace":  "oku",
		"telemetry.sdk.name": "opentelemetry",
	} {
		got, ok := attrs.Value(key)
		if !ok || got.AsString() != want {
			t.Errorf("%s = %v, want %q", key, got, want)
		}
	}
}
