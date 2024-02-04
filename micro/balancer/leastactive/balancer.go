package round_robin

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"sync/atomic"
)

type BalancerBuilder struct {
	size int32
}

func (b *BalancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conn := make([]balancer.SubConn, 0, len(info.ReadySCs))
	for k := range info.ReadySCs {
		conn = append(conn, k)
	}
	return &Balancer{
		connects: conn,
		index:    -1,
	}
}

type Balancer struct {
	connects []balancer.SubConn
	index    int32
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(b.connects) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	idx := atomic.AddInt32(&b.index, 1)
	c := idx % int32(len(b.connects))
	return balancer.PickResult{
		SubConn: b.connects[c],
		Done: func(info balancer.DoneInfo) {

		},
	}, nil
}
