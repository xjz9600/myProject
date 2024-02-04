package opentelemetry

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"google.golang.org/grpc"
	"log"
	"myProject/micro/grpc_demo/gen"
	"os"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	initZipkin(t)
	builder := &ClientOtelBuilder{}
	conn, err := grpc.Dial("localhost:8082", grpc.WithInsecure(), grpc.WithUnaryInterceptor(builder.Build()))
	assert.NoError(t, err)
	us := gen.NewUserServiceClient(conn)
	res, err := us.GetById(context.Background(), &gen.GetByIdReq{})
	assert.NoError(t, err)
	assert.Equal(t, res.User.Name, "JunZeXie")
	assert.Equal(t, res.User.Id, int64(5))
	time.Sleep(10 * time.Second)
}

func initZipkin(t *testing.T) {
	// 要注意这个端口，和 docker-compose 中的保持一致
	exporter, err := zipkin.New(
		"http://localhost:9411/api/v2/spans",
		zipkin.WithLogger(log.New(os.Stderr, "opentelemetry-demo", log.Ldate|log.Ltime|log.Llongfile)),
	)
	if err != nil {
		t.Fatal(err)
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("opentelemetry-demo"),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}))
}
