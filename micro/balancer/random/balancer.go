package random

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math/rand"
)

type BalancerBuilder struct {
}

func (b *BalancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conn := make([]balancer.SubConn, 0, len(info.ReadySCs))
	for k := range info.ReadySCs {
		conn = append(conn, k)
	}
	return &Balancer{
		connects: conn,
	}
}

type Balancer struct {
	connects []balancer.SubConn
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(b.connects) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	idx := rand.Intn(len(b.connects))
	return balancer.PickResult{
		SubConn: b.connects[idx],
		Done: func(info balancer.DoneInfo) {

		},
	}, nil
}
