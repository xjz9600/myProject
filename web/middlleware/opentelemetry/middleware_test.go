//go:build e2e

package opentelemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"log"
	"myProject/web"
	"os"
	"testing"
	"time"
)

const instrumentationName = "github.com/xjz9600/myProject/web/middleware/tracing"

func TestNewMiddleWare(t *testing.T) {
	tracer := otel.GetTracerProvider().Tracer(instrumentationName)
	tracing := NewMiddleware(tracer)
	initZipkin(t)
	server := web.NewServer(web.WithMiddleWare(tracing.Build()))
	server.GET("/a/b/*", func(ctx *web.Context) {
		c, span := tracer.Start(ctx.Req.Context(), "first_layer")
		defer span.End()
		c, second := tracer.Start(c, "second_layer")
		time.Sleep(time.Second)
		c, third1 := tracer.Start(c, "third_layer_1")
		time.Sleep(200 * time.Millisecond)
		third1.End()
		time.Sleep(100 * time.Millisecond)
		c, third2 := tracer.Start(c, "third_layer_2")
		third2.AddEvent("出什么事了？？？")
		time.Sleep(300 * time.Millisecond)
		third2.End()
		second.End()
		ctx.RespJson(&User{Name: "JunZeXie"})
	})
	server.Start(":8085")
}

type User struct {
	Name string `json:"name"`
}

func initJeager(t *testing.T) {
	url := "http://localhost:14268/api/traces"
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		t.Fatal(err)
	}
	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithBatcher(exp),
		// Record information about this application in a Resource.
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("opentelemetry-demo"),
			attribute.String("environment", "dev"),
			attribute.Int64("ID", 1),
		)),
	)

	otel.SetTracerProvider(tp)
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
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{}))
}
