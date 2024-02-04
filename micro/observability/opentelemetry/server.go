package opentelemetry

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

const instrumentationName = "https://github.com/xjz9600/myProject/micro/observability/opentelemetry"

type TelServerBuilder struct {
	t trace.Tracer
}

func (s *TelServerBuilder) Build() grpc.UnaryServerInterceptor {
	if s.t == nil {
		s.t = otel.GetTracerProvider().Tracer(instrumentationName)
	}
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		ctx = s.extract(ctx)
		spanCtx, span := s.t.Start(ctx, info.FullMethod, trace.WithSpanKind(trace.SpanKindServer))
		span.SetAttributes(attribute.String("address", "[:]server"))
		time.Sleep(5 * time.Second)
		defer func() {
			if err != nil {
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
		return handler(spanCtx, req)
	}
}

func (s *TelServerBuilder) extract(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	return otel.GetTextMapPropagator().Extract(ctx, &metadataSupplier{metadata: md})
}
