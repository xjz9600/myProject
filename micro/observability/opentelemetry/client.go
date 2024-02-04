package opentelemetry

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ClientOtelBuilder struct {
	Tracer trace.Tracer
}

func (b *ClientOtelBuilder) Build() grpc.UnaryClientInterceptor {
	if b.Tracer == nil {
		b.Tracer = otel.GetTracerProvider().Tracer(instrumentationName)
	}
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		spanCtx, span := b.Tracer.Start(ctx, method, trace.WithSpanKind(trace.SpanKindClient))
		span.SetAttributes(attribute.String("address", "[:]client"))
		var err error
		defer func() {
			if err != nil {
				span.SetStatus(codes.Error, err.Error())
				span.RecordError(err)
			}
			span.End()
		}()
		spanCtx = b.inject(spanCtx)
		err = invoker(spanCtx, method, req, reply, cc, opts...)
		return err
	}
}

func (b *ClientOtelBuilder) inject(ctx context.Context) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	otel.GetTextMapPropagator().Inject(ctx, &metadataSupplier{metadata: md})
	return metadata.NewOutgoingContext(ctx, md)
}
