package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

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
	trace2 "go.opentelemetry.io/otel/trace"
)

// In-memory mock database for user data
var users = map[string]string{
	"1": "Alice",
	"2": "Bob",
	"3": "Charlie",
}

func main() {
	// Step 1: Initialize OpenTelemetry
	traceProvider, err := initOpenTelemetry()
	if err != nil {
		panic(err)
	}
	defer traceProvider.Shutdown(context.Background())

	// Step 2: Create a Gin router
	r := gin.Default()

	// Step 3: Add OpenTelemetry middleware to Gin
	// This will automatically trace all incoming HTTP requests handled by Gin.
	r.Use(otelgin.Middleware("user-micro"))

	// Step 4: Define the `/user` endpoint
	r.GET("/user/:id", GetUser)
	r.GET("/call-order-service", CallOrderService)

	// Step 5: Start the server on port 8086
	if err := r.Run(":8086"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func CallOrderService(c *gin.Context) {
	ctx, span := otel.Tracer("user-micro").Start(c.Request.Context(), "CallOrderService")
	defer span.End()

	client := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8089/order", nil)
	if err != nil {
		log.Printf("Failed to construct request for the order service: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to call the order service: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to call order service: %v", err)
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}

	span.SetStatus(codes.Ok, "order service response received")
	log.Printf("Order service response: %s", string(b))
	c.String(http.StatusOK, string(b))
}

// Initialize OpenTelemetry for distributed tracing
func initOpenTelemetry() (*trace.TracerProvider, error) {
	// Set up the OTLP trace exporter to send tracing data to the OpenTelemetry Collector
	traceExporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithInsecure(),                 // No TLS for local development
		otlptracehttp.WithEndpoint("localhost:4318"), // Collector/Jaeger endpoint
	)
	if err != nil {
		return nil, err
	}

	// Configure a TracerProvider
	traceProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()), // Sample all traces
		trace.WithBatcher(traceExporter),        // Batch traces in export
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("user-micro"), // Set the service name for tracing
		)),
	)

	// Register the global TracerProvider and propagators
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return traceProvider, nil
}

// Handler for the `/user` endpoint
// Retrieves a user by ID based on the query parameter `id`
func GetUser(c *gin.Context) {
	// Step 6: Start a new span for the handler
	tracer := otel.Tracer("user-micro")
	ctx, span := tracer.Start(c.Request.Context(), "GetUser")
	defer span.End()

	// Extract the user ID from the query parameters
	userId := c.Param("id")
	traceId := span.SpanContext().TraceID().String()

	// Retrieve user data from the mock database
	user, err := GetUserById(ctx, userId, tracer)
	if err != nil {
		// Handle and record any errors in the span
		span.SetAttributes(semconv.HTTPResponseStatusCodeKey.Int(http.StatusBadRequest)) // HTTP 400
		span.SetAttributes(attribute.String("user_id", userId))                          // Attach user ID
		span.SetAttributes(attribute.String("traceId", traceId))
		span.AddEvent("USER NOT FOUND")                            // Record event in tracing span// Attach trace ID
		span.SetStatus(codes.Error, err.Error())                   // Mark span as error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Respond with 400 error
		return
	}

	// Mark the span as successful
	span.SetAttributes(semconv.HTTPResponseStatusCodeKey.Int(http.StatusOK)) // HTTP 200
	c.String(http.StatusOK, user)                                            // Respond with the user
}

// Simulate retrieving a user by ID (mock database operation)
func GetUserById(ctx context.Context, userId string, tracer trace2.Tracer) (string, error) {
	// Step 7: Create a child span for this database operation
	_, span := tracer.Start(ctx, "GetUserById")
	defer span.End()

	// Log the trace ID for debugging
	traceId := span.SpanContext().TraceID().String()
	log.Printf("traceID: %s - fetching user from database", traceId)

	// Check the in-memory database for the user ID
	if user, exists := users[userId]; exists {
		// Record success in the span
		span.SetStatus(codes.Ok, "user found")
		return user, nil
	}

	// Record failure in the span
	span.SetStatus(codes.Error, "user not found")
	return "", errors.New("user not found")
}
