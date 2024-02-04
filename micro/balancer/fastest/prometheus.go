package fastest

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"time"
)

type Prometheus struct {
	Help       string
	Name       string
	ServerName string
	Port       string
}

func NewPrometheusBuilder(serverName, name, port, help string) *Prometheus {
	return &Prometheus{
		Name:       name,
		Port:       port,
		Help:       help,
		ServerName: serverName,
	}
}

func (p Prometheus) Build() grpc.UnaryServerInterceptor {
	addr := "[::]"
	if p.Port != "" {
		addr += p.Port
	} else {
		addr += ":80"
	}
	vector := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       p.Name,
			Help:       p.Help,
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			ConstLabels: map[string]string{
				"addr": addr,
			},
		},
		[]string{"kind"},
	)
	prometheus.MustRegister(vector)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		startTime := time.Now()
		defer func() {
			duration := time.Now().Sub(startTime).Nanoseconds()
			vector.WithLabelValues(p.ServerName).Observe(float64(duration))
		}()
		return handler(ctx, req)
	}
}
