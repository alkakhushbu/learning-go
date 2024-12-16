package ctxmanage

import (
	"book-store/middlewares"
	"context"
	"log/slog"
)

func GetTraceID(ctx context.Context) string {

	// fetching value from the context
	// fetching traceId present in the context
	traceId, ok := ctx.Value(middlewares.TraceIdKey).(string)
	if !ok {
		slog.Error("traceId not found in the context")
		return "UNKNOWN"
	}
	return traceId
}
