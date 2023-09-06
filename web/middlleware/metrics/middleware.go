package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"myProject/web"
	"strconv"
	"time"
)

type MiddlewareBuilder struct {
	namespace string
	subsystem string
	name      string
	help      string
}

func NewMiddleware(namespace, subsystem, name, help string) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		name:      name,
		subsystem: subsystem,
		namespace: namespace,
		help:      help,
	}
}

func (m *MiddlewareBuilder) Build() web.Middleware {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      m.name,
		Namespace: m.namespace,
		Help:      m.help,
		Subsystem: m.subsystem,
	}, []string{"pattern", "method", "status"})
	prometheus.MustRegister(vector)
	return func(handleFunc web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			startTime := time.Now()
			defer func() {
				duration := time.Now().Sub(startTime).Milliseconds()
				pattern := ctx.MathRoute
				if pattern == "" {
					pattern = "unknown"
				}
				vector.WithLabelValues(pattern, ctx.Req.Method, strconv.Itoa(ctx.RespStatusCode)).Observe(float64(duration))
			}()
			handleFunc(ctx)
		}
	}
}
