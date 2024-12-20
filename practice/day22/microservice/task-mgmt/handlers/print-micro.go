package handlers

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

func CallPrintMicroService(c *gin.Context) {
	ctx, span := otel.Tracer("user-micro").Start(c.Request.Context(), "CallPrintMicroService")
	defer span.End()

	client := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	id := c.Param("id")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8089/print/task/"+id, nil)
	if err != nil {
		log.Printf("Failed to construct request for the Print micro service: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to call the print micro service: %v", err)
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response from print micro service: %v", err)
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	span.SetStatus(codes.Ok, "print micro service response received")
	log.Printf("Order service response: %s", string(b))
	c.String(http.StatusOK, string(b))
}
