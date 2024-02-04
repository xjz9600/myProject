package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"time"
)

type ServerMetricsBuilder struct {
	namespace string
	subsystem string
	help      string
}

func NewServerMetricsBuilder(namespace, subsystem, help string) *ServerMetricsBuilder {
	return &ServerMetricsBuilder{
		namespace: namespace,
		subsystem: subsystem,
		help:      help,
	}
}

func (s *ServerMetricsBuilder) Build() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		reqGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "active_request_cnt",
			Subsystem: s.subsystem,
			Namespace: s.namespace,
			Help:      s.help,
		}, []string{"kind"})

		errCnt := prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "err_request_cnt",
			Subsystem: s.subsystem,
			Namespace: s.namespace,
			Help:      s.help,
		}, []string{"kind"})

		costTime := prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Name:      "err_request_cnt",
			Subsystem: s.subsystem,
			Namespace: s.namespace,
			Help:      s.help,
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		}, []string{"kind"})
		prometheus.MustRegister(reqGauge, errCnt, costTime)
		startTime := time.Now()
		reqGauge.WithLabelValues(info.FullMethod).Add(1)
		defer func() {
			reqGauge.WithLabelValues(info.FullMethod).Add(-1)
			costTime.WithLabelValues(info.FullMethod).Observe(float64(time.Now().Sub(startTime).Milliseconds()))
			if err != nil {
				errCnt.WithLabelValues(info.FullMethod).Inc()
			}
		}()
		return handler(ctx, req)
	}
}
