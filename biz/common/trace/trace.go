package trace

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hertz-contrib/logger/zap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// StartDBOperation starts a new span for a database operation
func StartDBOperation(ctx context.Context, operation string) (context.Context, trace.Span) {
	tracer := otel.Tracer("database")
	ctx, span := tracer.Start(ctx, fmt.Sprintf("db.%s", operation))
	span.SetAttributes(
		attribute.String("db.system", "mysql"),
		attribute.String("X-Request-ID", ctx.Value(zap.ExtraKey("X-Request-ID")).(string)),
	)
	return ctx, span
}

func StartRedisOperation(ctx context.Context, operation string, key string) (context.Context, trace.Span) {
	tracer := otel.Tracer("redis")
	ctx, span := tracer.Start(ctx, fmt.Sprintf("redis.%s", operation))
	if key != "" {
		span.SetAttributes(
			attribute.String("key", key),
			attribute.String("X-Request-ID", ctx.Value(zap.ExtraKey("X-Request-ID")).(string)),
		)
	}
	return ctx, span
}

// StartHTTPOperation starts a new span for an HTTP operation
func StartHTTPOperation(ctx context.Context, req *http.Request) (context.Context, trace.Span) {
	tracer := otel.Tracer("http")
	ctx, span := tracer.Start(ctx, req.URL.Path)
	span.SetAttributes(
		attribute.String("http.url", req.URL.String()),
		attribute.String("http.method", req.Method),
		attribute.String("X-Request-ID", ctx.Value(zap.ExtraKey("X-Request-ID")).(string)),
	)
	return ctx, span
}

// EndSpan ends a span and records any error
func EndSpan(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "success")
	}
	span.End()
}

// RecordHTTPResponse records HTTP response details in the span
func RecordHTTPResponse(span trace.Span, resp *http.Response) {
	if resp != nil {
		span.SetAttributes(
			attribute.Int("http.status_code", resp.StatusCode),
			attribute.String("http.status_text", resp.Status),
		)
	}
}
