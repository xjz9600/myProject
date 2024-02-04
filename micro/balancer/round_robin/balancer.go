package round_robin

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
	mybalancer "myProject/micro/balancer"
	"sync/atomic"
)

type BalancerBuilder struct {
	Filter mybalancer.Filter
}

func (b *BalancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conn := make([]subConn, 0, len(info.ReadySCs))
	for k, v := range info.ReadySCs {
		cn := subConn{
			c:    k,
			addr: v.Address,
		}
		conn = append(conn, cn)
	}
	filter := func(info balancer.PickInfo, addr resolver.Address) bool {
		return true
	}
	if b.Filter != nil {
		filter = b.Filter
	}
	return &Balancer{
		connects: conn,
		index:    -1,
		filter:   filter,
	}
}

type Balancer struct {
	connects []subConn
	index    int32
	filter   mybalancer.Filter
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	var candidates []subConn
	for _, cn := range b.connects {
		if !b.filter(info, cn.addr) {
			continue
		}
		candidates = append(candidates, cn)
	}
	if len(candidates) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	idx := atomic.AddInt32(&b.index, 1)
	c := idx % int32(len(candidates))
	return balancer.PickResult{
		SubConn: candidates[c].c,
		Done: func(info balancer.DoneInfo) {

		},
	}, nil
}

type subConn struct {
	c    balancer.SubConn
	addr resolver.Address
}
