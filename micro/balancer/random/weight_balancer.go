package random

import (
	"fmt"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math/rand"
)

type WeightBalancerBuilder struct {
}

func (w *WeightBalancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conn := make([]*weightConn, len(info.ReadySCs))
	var totalWeight uint32
	for k, v := range info.ReadySCs {
		weight, ok := v.Address.Attributes.Value("weight").(uint32)
		if !ok {
			panic(fmt.Sprintf("micro：节点 %s 没有设置权重", v.Address.Addr))
		}
		wc := &weightConn{
			c:      k,
			weight: weight,
		}
		conn = append(conn, wc)
		totalWeight += weight
	}
	return &WeightBalancer{connects: conn, totalWeight: totalWeight}
}

type WeightBalancer struct {
	connects    []*weightConn
	totalWeight uint32
}

func (w *WeightBalancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(w.connects) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	var idx int
	tw := rand.Intn(int(w.totalWeight)) + 1
	for i, cn := range w.connects {
		tw = tw - int(cn.weight)
		if tw <= 0 {
			idx = i
			break
		}
	}
	return balancer.PickResult{
		SubConn: w.connects[idx].c,
		Done: func(info balancer.DoneInfo) {
		},
	}, nil
}

type weightConn struct {
	c      balancer.SubConn
	weight uint32
}
