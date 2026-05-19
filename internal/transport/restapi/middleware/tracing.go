package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

// Tracing starts a server span for each HTTP request.
func Tracing(serviceName string) gin.HandlerFunc {
	tracer := otel.Tracer(serviceName + "/restapi")
	propagator := otel.GetTextMapPropagator()

	return func(ctx *gin.Context) {
		requestCtx := propagator.Extract(ctx.Request.Context(), propagation.HeaderCarrier(ctx.Request.Header))

		requestCtx, span := tracer.Start(requestCtx, ctx.FullPath())
		defer span.End()

		ctx.Request = ctx.Request.WithContext(requestCtx)
		ctx.Next()

		status := ctx.Writer.Status()
		span.SetAttributes(
			attribute.String("http.method", ctx.Request.Method),
			attribute.String("http.route", ctx.FullPath()),
			attribute.Int("http.status_code", status),
			attribute.String("http.response.size", strconv.Itoa(ctx.Writer.Size())),
		)

		if status >= http.StatusInternalServerError {
			span.SetStatus(codes.Error, "server error")
		}
	}
}
