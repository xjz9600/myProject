package leastactive

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"sync/atomic"
)

type BalancerBuilder struct {
}

func (w *BalancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conn := make([]*activeConn, len(info.ReadySCs))
	for k := range info.ReadySCs {
		wc := &activeConn{
			c: k,
		}
		conn = append(conn, wc)
	}
	return &Balancer{connects: conn}
}

type Balancer struct {
	connects []*activeConn
}

func (w *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(w.connects) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	var conn *activeConn
	for _, cn := range w.connects {
		if conn == nil || atomic.LoadUint32(&cn.cnt) < conn.cnt {
			conn = cn
		}
	}
	atomic.AddUint32(&conn.cnt, 1)
	return balancer.PickResult{
		SubConn: conn.c,
		Done: func(info balancer.DoneInfo) {
			atomic.AddUint32(&conn.cnt, -1)
		},
	}, nil
}

type activeConn struct {
	c   balancer.SubConn
	cnt uint32
}
