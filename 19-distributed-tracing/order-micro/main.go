package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"log"
	"net/http"
	"runtime/debug"
)

func main() {
	tp, err := initOpenTelemetry()
	if err != nil {
		panic(err)
	}
	defer tp.Shutdown(context.Background())
	// Step 2: Create a Gin router
	r := gin.Default()

	// Step 3: Add OpenTelemetry middleware to Gin
	// This will automatically trace all incoming HTTP requests handled by Gin.
	r.Use(otelgin.Middleware("order-micro"))

	// Step 4: Define the `/user` endpoint
	r.GET("/order", GetOrder)

	// Step 5: Start the server on port 8086
	if err := r.Run(":8089"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func GetOrder(c *gin.Context) {
	propogator := otel.GetTextMapPropagator()
	extractedCtx := propogator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

	fmt.Printf("viewing headers")
	for k, v := range c.Request.Header {
		fmt.Println(k, v)
	}

	ctx, span := otel.Tracer("order-micro").Start(extractedCtx, "GetOrder Handler")
	defer span.End()

	err := FetchOrder(ctx)
	if err != nil {
		span.SetStatus(codes.Error, "order not found")
		span.AddEvent("order not found")
		span.SetAttributes(attribute.Int("order.id failed", n))
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error()})
		return
	}

	span.SetAttributes(semconv.HTTPResponseStatusCodeKey.Int(http.StatusOK))
	c.JSON(http.StatusOK, gin.H{
		"message": "order found",
	})

}

var n = 0

func FetchOrder(ctx context.Context) error {
	n++
	_, span := otel.Tracer("order-micro").Start(ctx, "FetchOrder")
	defer span.End() // if you forget to end the span then your trace won't show up
	fmt.Println(span.SpanContext().TraceID().String())
	if n%2 == 0 {
		span.SetAttributes(attribute.Int("order.id", n))
		span.AddEvent(string(debug.Stack()))
		return fmt.Errorf("order not found")
	}
	span.SetAttributes(
		attribute.Int("order.id", n),
		attribute.String("order.status", "completed"),
	)

	return nil
}

func initOpenTelemetry() (*trace.TracerProvider, error) {
	// Step 4: Set up a trace exporter.
	// The trace exporter sends captured trace data to an OpenTelemetry collector or a backend.
	// Here, we configure an OTLP exporter that uses HTTP to send data to a collector
	//running at `localhost:4318`.
	traceExporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint("localhost:4318"), // Specify the OpenTelemetry collector endpoint.
	)
	if err != nil {

		return nil, err
	}

	// Step 5: Configure the TracerProvider.
	// The TracerProvider manages all spans (units of trace data) created in your application.
	traceProvider := trace.NewTracerProvider(
		// Set a sampling strategy. `trace.WithSampler(trace.AlwaysSample())` ensures all requests are traced.
		// In production, you can adjust this for cost and performance, e.g., to sample a percentage of traces.
		trace.WithSampler(trace.AlwaysSample()),
		//0.0 would mean sampling 0% of requests (never sampling).
		//0.5 would mean sampling 50% of requests (half of them).
		//1.0 would mean sampling 100% of requests (all of them).
		//
		//trace.WithSampler(trace.TraceIDRatioBased(0.1)), // 10% of traces

		// Use a batch processor to optimize the performance of exporting traces. This sends trace data in batches.
		trace.WithBatcher(traceExporter),

		// Define resources (metadata) associated with traces, such as service name.
		// These attributes can help group or filter traces in the backend (e.g., Jaeger, Zipkin).
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL, // OpenTelemetry semantic conventions schema URL.
			semconv.ServiceNameKey.String("order-micro"), // Set the service name for identification in the tracing backends.
		)),
	)

	// Step 6: Register the TracerProvider globally for the application.
	// This means any tracing operation in your code will use this TracerProvider by default.
	otel.SetTracerProvider(traceProvider)

	// Step 7: Set a propagator for distributed tracing.
	// Propagators ensure that trace information (like trace IDs) is transmitted between different services.
	// The `TraceContext` propagator manages the trace context over HTTP headers.
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Return the configured TracerProvider so it can control the trace lifecycle.
	return traceProvider, nil
}
