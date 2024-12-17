package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"time"
)

type key string

const TraceIdKey key = "1"

func LoggerV0(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		next(c)
	}
}

func mid(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestStartTime := time.Now()

		// get a trace id
		traceId := uuid.NewString()

		//fetching the context container from the request.context()
		// putting the traceId in the context
		ctx := context.WithValue(c.Request.Context(), TraceIdKey, traceId)

		// this creates a new copy of the request with the updated context
		c.Request = c.Request.WithContext(ctx)

		slog.Info("STARTED", slog.String("TRACE ID", traceId),
			slog.String("Method", c.Request.Method), slog.Any("URL Path", c.Request.URL.Path))

		//we use c.Next only when we are using r.Use() method to assign middlewares
		c.Next() // call next thing in the chain

		slog.Info("COMPLETED", slog.String("TRACE ID", traceId),
			slog.String("Method", c.Request.Method), slog.Any("URL Path", c.Request.URL.Path),
			slog.Int("Status Code", c.Writer.Status()), slog.Int64("duration μs,",
				time.Since(requestStartTime).Milliseconds()))

	}
}
