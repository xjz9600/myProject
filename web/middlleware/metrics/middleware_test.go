package metrics

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"myProject/web"
	"net/http"
	"testing"
	"time"
)

func TestNewMiddleWare(t *testing.T) {
	metrics := NewMiddleware("http_response", "web", "test", "firstPrometheus")
	server := web.NewServer(web.WithMiddleWare(metrics.Build()))
	server.GET("/test/metric", func(ctx *web.Context) {
		random := rand.Intn(1000)
		time.Sleep(time.Duration(random) * time.Millisecond)
		ctx.RespJson(&User{
			Name: "xieJunZe",
		})
	})
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8086", nil)
	}()
	server.Start(":8085")
}

type User struct {
	Name string `json:"name"`
}
