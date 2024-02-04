package hash

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
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
	// 在这个地方你拿不到请求，无法做根据请求特性做负载均衡
	//idx := info.Ctx.Value("user_id")
	//idx := info.Ctx.Value("hash_code")
	return balancer.PickResult{
		SubConn: b.connects[0],
		Done: func(info balancer.DoneInfo) {

		},
	}, nil
}
