package gotel

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"io"
	"log/slog"
	"os"
	"time"
)

const (
	TraceExporterEndpoint  = "OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"
	ExporterEndpoint       = "OTEL_EXPORTER_OTLP_ENDPOINT"
	DefaultExporterTimeout = 5 * time.Second
)

func getFirstEnvironmentValue(varNames ...string) (result string) {
	for _, varName := range varNames {
		if result = os.Getenv(varName); len(result) > 0 {
			return result
		}
	}

	return
}

// Option specifies configuration options for tracing
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

type config struct {
	serviceName          string
	serviceNamespace     string
	traceExporter        trace.SpanExporter
	endpoint             string
	exporterBatchTimeout *time.Duration
	sampler              trace.Sampler
	propagators          propagation.TextMapPropagator
	traceLogLevel        slog.Level
}

type ShutdownFunc func(ctx context.Context) error

func WithServiceName(serviceName string) Option {
	return optionFunc(func(config *config) {
		config.serviceName = serviceName
	})
}

func WithServiceNamespace(serviceNamespace string) Option {
	return optionFunc(func(config *config) {
		config.serviceNamespace = serviceNamespace
	})
}

func WithExporter(exporter trace.SpanExporter) Option {
	return optionFunc(func(config *config) {
		config.traceExporter = exporter
	})
}

func WithEndpoint(endpoint string) Option {
	return optionFunc(func(config *config) {
		config.endpoint = endpoint
	})
}

func WithBatchTimeout(dur time.Duration) Option {
	return optionFunc(func(config *config) {
		config.exporterBatchTimeout = &dur
	})
}

func WithDefaultBatchTimeout() Option {
	return optionFunc(func(config *config) {
		var dur = DefaultExporterTimeout
		config.exporterBatchTimeout = &dur
	})
}

func WithSampler(sampler trace.Sampler) Option {
	return optionFunc(func(config *config) {
		config.sampler = sampler
	})
}

func WithTextMapPropagator(propagator propagation.TextMapPropagator) Option {
	return optionFunc(func(config *config) {
		config.propagators = propagator
	})
}

func WithDefaultTextMapPropagator() Option {
	return optionFunc(func(config *config) {
		config.propagators = propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)
	})
}

func WithTraceLogLevel(level slog.Level) Option {
	return optionFunc(func(config *config) {
		config.traceLogLevel = level
	})
}

func newConfig(ctx context.Context, options ...Option) (*config, error) {
	config := &config{
		serviceName:      "<unspecified>",
		serviceNamespace: "<unspecified>",
		traceLogLevel:    slog.LevelInfo,
	}
	for _, option := range options {
		option.apply(config)
	}

	if config.sampler == nil {
		config.sampler = trace.AlwaysSample()
	}

	if config.propagators == nil {
		WithDefaultTextMapPropagator().apply(config)
	}

	// if a Trace Exporter is not provided, use the endpoints(s) specified
	// via the environment (recommended), or just discard everything. Note that
	// if you want the otel default (localhost:4317) then you must use WithEndpoint
	// as our default is io.Discard
	if config.traceExporter == nil {
		var err error
		if len(config.endpoint) > 0 {
			config.traceExporter, err = otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(config.endpoint))
		} else if getFirstEnvironmentValue(TraceExporterEndpoint, ExporterEndpoint) != "" {
			config.traceExporter, err = otlptracehttp.New(ctx)
		} else {
			config.traceExporter, err = stdouttrace.New(stdouttrace.WithWriter(io.Discard))
		}
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func initTracing(ctx context.Context, config *config) (ShutdownFunc, error) {
	hostName := os.Getenv("HOSTNAME")
	if len(hostName) == 0 {
		hostName = "<unknown>"
	}

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.HostName(hostName),
			semconv.ServiceName(config.serviceName),
			semconv.ServiceNamespace(config.serviceNamespace),
		),
	)

	if err != nil {
		return nil, err
	}

	var batchOption trace.TracerProviderOption

	if config.exporterBatchTimeout != nil {
		batchOption = trace.WithBatcher(
			config.traceExporter,
			trace.WithBatchTimeout(*config.exporterBatchTimeout))
	} else {
		batchOption = trace.WithBatcher(config.traceExporter)
	}

	// create the traceProvider and batch all span exports
	traceProvider := trace.NewTracerProvider(
		batchOption,
		trace.WithResource(r),
		trace.WithSampler(config.sampler))

	// set the text map propagator
	otel.SetTextMapPropagator(config.propagators)

	// set our trace provider as the global trace provider
	otel.SetTracerProvider(traceProvider)

	// return a function to shut down the trace provider
	return traceProvider.Shutdown, nil
}

func InitTracing(ctx context.Context, options ...Option) (ShutdownFunc, error) {
	config, err := newConfig(ctx, options...)
	if err != nil {
		return nil, err
	}

	return initTracing(ctx, config)
}
