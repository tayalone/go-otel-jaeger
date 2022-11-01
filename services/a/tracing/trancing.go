package tracing

import (
	"os"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

/*JaegerProvider create Jeager Provider for otel*/
func JaegerProvider() (*sdktrace.TracerProvider, error) {
	ep := os.Getenv("JEAGER_ENDPOINT")

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(ep)))
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("a-service"),
			semconv.DeploymentEnvironmentKey.String("production"),
		)),
	)

	return tp, nil
}
