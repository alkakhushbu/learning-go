package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func main() {
	tp, err := initOpenTelemetry()
	if err != nil {
		panic(err)
	}
	defer tp.Shutdown(context.Background())

	router := gin.Default()
	
	router.GET("/print/task/:id", PrintTask)
	router.Run(":8089")
}

func PrintTask(c *gin.Context) {
	otel.GetTextMapPropagator()

	id, _ := strconv.Atoi(c.Param("id"))
	log.Println("Task added")
	task := struct {
		id       int
		msg      string
		taskType string
	}{
		id:  id,
		msg: "Task added",
	}
	if id%2 == 0 {
		task.taskType = "Even Task"
	} else {
		task.taskType = "Odd Task"
	}
	c.JSON(http.StatusOK, gin.H{"response": task})
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
