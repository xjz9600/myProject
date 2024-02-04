package round_robin

import (
	"fmt"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math"
	"sync"
	"sync/atomic"
)

type WeightBalancerBuilder struct {
}

func (w *WeightBalancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conn := make([]*weightConn, len(info.ReadySCs))
	for k, v := range info.ReadySCs {
		weight, ok := v.Address.Attributes.Value("weight").(uint32)
		if !ok {
			panic(fmt.Sprintf("micro：节点 %s 没有设置权重", v.Address.Addr))
		}
		wc := &weightConn{
			c:               k,
			weight:          weight,
			currentWeight:   weight,
			efficientWeight: weight,
		}
		conn = append(conn, wc)
	}
	return &WeightBalancer{connects: conn}
}

type WeightBalancer struct {
	connects []*weightConn
	mutex    sync.Mutex
}

func (w *WeightBalancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(w.connects) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	var totalWeight uint32
	var conn *weightConn
	w.mutex.Lock()
	for _, cn := range w.connects {
		totalWeight = totalWeight + cn.efficientWeight
		cn.currentWeight = cn.currentWeight + cn.efficientWeight
		if conn == nil || conn.currentWeight < cn.currentWeight {
			conn = cn
		}
	}
	conn.currentWeight = conn.currentWeight - totalWeight
	w.mutex.Unlock()
	return balancer.PickResult{
		SubConn: conn.c,
		Done: func(info balancer.DoneInfo) {
			for {
				weight := atomic.LoadUint32(&conn.efficientWeight)
				if info.Err == nil && weight == math.MaxUint32 {
					return
				}
				if info.Err != nil && weight == 0 {
					return
				}
				newWeight := weight
				if info.Err != nil {
					newWeight--
				} else {
					newWeight++
				}
				if atomic.CompareAndSwapUint32(&conn.efficientWeight, weight, newWeight) {
					return
				}
			}
		},
	}, nil
}

type weightConn struct {
	c               balancer.SubConn
	weight          uint32
	currentWeight   uint32
	efficientWeight uint32
}
