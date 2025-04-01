package trace

import (
	"context"
	"fmt"
	"net/http"

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
	)
	return ctx, span
}

func StartRedisOperation(ctx context.Context, operation string, key string) (context.Context, trace.Span) {
	tracer := otel.Tracer("redis")
	ctx, span := tracer.Start(ctx, fmt.Sprintf("redis.%s", operation))
	if key != "" {
		span.SetAttributes(
			attribute.String("key", key),
		)
	}
	return ctx, span
}

// StartHTTPOperation starts a new span for an HTTP operation
func StartHTTPOperation(ctx context.Context, method, url string) (context.Context, trace.Span) {
	tracer := otel.Tracer("http")
	ctx, span := tracer.Start(ctx, fmt.Sprintf("http.%s", method))
	span.SetAttributes(
		attribute.String("http.method", method),
		attribute.String("http.url", url),
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

// RecordDBQuery records database query details in the span
func RecordDBQuery(span trace.Span, query string, args ...interface{}) {
	span.SetAttributes(
		attribute.String("db.query", query),
		attribute.String("db.args", fmt.Sprintf("%v", args)),
	)
}

func RecordRedisQuery(span trace.Span, op string, args ...interface{}) {
	span.SetAttributes(
		attribute.String("redis.operation", op),
		attribute.String("redis.key", fmt.Sprintf("%v", args)),
	)
}
