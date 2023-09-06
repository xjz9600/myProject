package opentelemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"myProject/web"
)

type MiddlewareBuilder struct {
	tracer trace.Tracer
}

func NewMiddleware(tracer trace.Tracer) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		tracer: tracer,
	}
}

func (m *MiddlewareBuilder) Build() web.Middleware {
	return func(handleFunc web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			reqCtx := ctx.Req.Context()
			reqCtx = otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(ctx.Req.Header))

			context, span := m.tracer.Start(reqCtx, "unknown")
			defer span.End()
			span.SetAttributes(attribute.String("http.method", ctx.Req.Method))
			span.SetAttributes(attribute.String("http.utl", ctx.Req.URL.String()))
			span.SetAttributes(attribute.String("http.host", ctx.Req.Host))
			ctx.Req = ctx.Req.WithContext(context)
			handleFunc(ctx)
			span.SetName(ctx.MathRoute)
		}
	}
}
