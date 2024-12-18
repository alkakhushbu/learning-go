package main

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"log"
	"runtime/debug"

	// Importing OpenTelemetry packages for HTTP instrumentation and tracing
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	trace2 "go.opentelemetry.io/otel/trace"
	"net/http"
)

var users = map[string]string{
	"1": "Alice",
	"2": "Bob",
	"3": "Charlie",
}

func main() {
	mux := http.NewServeMux()

	// Step 1: Initialize OpenTelemetry TracerProvider.
	// This sets up the tracing system, which will record information about the application's operations.
	// It enables capturing traces and sending them to a collector or backend for debugging,
	//performance monitoring, or distributed tracing.
	traceProvider, err := initOpenTelemetry()
	if err != nil {

		panic(err)
	}
	// Ensure the tracing system shuts down cleanly when the application stops.
	defer traceProvider.Shutdown(context.Background())

	// Step 2: Wrap the HTTP multiplexer (mux) with OpenTelemetry middleware.
	// The `otelMid` function applies instrumentation to automatically create spans for HTTP requests.
	wrappedMux := otelMid(mux)

	// Step 3: Register the "/user" endpoint, which is traced via OpenTelemetry.
	mux.HandleFunc("/user", GetUser)

	// Start an HTTP server on port 8086.
	// All requests going through `wrappedMux` will automatically have tracing information.
	panic(http.ListenAndServe(":8086", wrappedMux))
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
			semconv.ServiceNameKey.String("user-micro"), // Set the service name for identification in the tracing backends.
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

func otelMid(next http.Handler) http.Handler {
	// Step 8: Add OpenTelemetry middleware to trace HTTP requests.
	// This middleware intercepts all incoming HTTP requests to automatically start spans,
	// collect trace details (like HTTP method, URL), and associate them with the requested handler.
	// The second argument ("user-micro") specifies the name of the operation being traced (use your service name here).
	return otelhttp.NewHandler(next, "user-micro")
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Step 9: Start a new trace span for the "/user" request.
	// This span represents a unit of work for this operation
	//and will capture timing, metadata, and errors (if any).
	tracer := otel.Tracer("user-micro")               // Obtain a tracer associated with your service.
	ctx, span := tracer.Start(r.Context(), "GetUser") // Create a span named "GetUser".
	defer span.End()                                  // Ensure the span ends when the operation is complete.

	// Step 10: Attach custom attributes (key-value pairs) to the span.
	// These attributes make it easier to debug or filter traces in the backend.
	traceId := span.SpanContext().TraceID().String() // Extract the trace ID for debugging purposes.

	// Extract "id" parameter from the request and continue the process...
	// Tracing continues in the `GetUserById` function.
	userId := r.URL.Query().Get("id")             // Extract "id" query parameter.
	user, err := GetUserById(ctx, userId, tracer) // Pass context (with trace info) to downstream function.

	if err != nil {
		// If an error occurs, mark it in the span.
		// These details will show up in the tracing backend, making debugging easier.
		span.SetAttributes(semconv.HTTPResponseStatusCodeKey.Int(http.StatusBadRequest)) // Add HTTP 400 status.
		span.SetAttributes(attribute.String("user_id", userId))                          // Add user ID for context.
		span.SetAttributes(attribute.String("traceId", traceId))                         // Add trace ID for easier debugging.
		span.SetStatus(codes.Error, err.Error())                                         // Mark the span as an error.

		http.Error(w, err.Error(), http.StatusBadRequest) // Notify the client with the error.
		return
	}

	// If there's no error, mark the span as successful.
	span.SetAttributes(semconv.HTTPResponseStatusCodeKey.Int(http.StatusOK)) // Add HTTP 200 status.
	w.Write([]byte(user))                                                    // Send the retrieved user's name in response.
}

func GetUserById(ctx context.Context, userId string, tracer trace2.Tracer) (string, error) {
	// Step 11: Start a child span for the "database operation".
	// A child span allows breaking down the parent operation ("GetUser") into smaller steps for detailed analysis.
	_, span := tracer.Start(ctx, "GetUserById") // Create a new span for this specific task.
	defer span.End()                            // Ensure the span ends after the operation.

	// Attach additional information like the trace ID for debugging purposes.
	traceId := span.SpanContext().TraceID().String()
	log.Printf("traceId: %s hitting db\n", traceId) // Log the trace ID (simulates a "database call").

	// Check if the user exists in the mock database (map).
	if user, exists := users[userId]; exists {
		span.SetStatus(codes.Ok, "user found in database") // Mark span as successful if user is found.
		return user, nil
	}

	// If user does not exist, mark the span as an error.
	span.AddEvent(string(debug.Stack()))
	span.SetStatus(codes.Error, "user not found in database")
	return "", errors.New("user not found") // Return an error.
}
