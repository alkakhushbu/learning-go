package pkg

import (
	"context"
	"log/slog"
	"task-mgmt/midware"
)

func GetTraceId(ctx context.Context) string {
	traceId, ok := ctx.Value(midware.TraceId).(string)
	if !ok {
		slog.Error("traceId not found in the context")
		return "UNKNOWN"
	}

	return traceId
}
