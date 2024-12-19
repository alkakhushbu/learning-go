package midware

import (
	"context"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type key string

var TraceId key = "traceId"

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestStartTime := time.Now().UTC()
		traceId := uuid.NewString()

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, TraceId, traceId)

		// this creates a new copy of the request with the updated context
		c.Request = c.Request.WithContext(ctx)

		slog.Info("Started",
			slog.String("Trace ID", traceId),
			slog.String("Method", c.Request.Method),
			slog.String("URL", c.Request.RequestURI))

		// we use c.Next only when we are using r.Use() method to assign middlewares
		c.Next() //call next thing in the chain

		slog.Info("Completed",
			slog.String("Trace ID", traceId),
			slog.String("Method", c.Request.Method),
			slog.String("URL", c.Request.RequestURI),
			slog.Int("Status Code", c.Writer.Status()),
			slog.Int64("Duration in ms", time.Since(requestStartTime).Milliseconds()))
	}
}
