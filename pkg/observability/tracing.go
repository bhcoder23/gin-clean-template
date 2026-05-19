// Package observability wires optional telemetry providers.
package observability

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

var errUnsupportedTraceExporter = errors.New("unsupported trace exporter")

// TraceConfig configures tracing.
type TraceConfig struct {
	Enabled     bool
	Exporter    string
	ServiceName string
}

// InitTracing initializes an OpenTelemetry tracer provider.
func InitTracing(cfg TraceConfig) (func(context.Context) error, error) {
	if !cfg.Enabled {
		return func(context.Context) error { return nil }, nil
	}

	if cfg.ServiceName == "" {
		cfg.ServiceName = "gin-clean-template"
	}

	exporter, err := newTraceExporter(cfg.Exporter)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ServiceName),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp.Shutdown, nil
}

func newTraceExporter(exporter string) (sdktrace.SpanExporter, error) {
	switch exporter {
	case "", "stdout":
		return stdouttrace.New(stdouttrace.WithPrettyPrint())
	default:
		return nil, fmt.Errorf("observability - %w: %s", errUnsupportedTraceExporter, exporter)
	}
}
