package round_robin

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/balancer"
	"testing"
)

func Test_WeightBalancerPick(t *testing.T) {
	weightBalancer := &WeightBalancer{
		connects: []*weightConn{
			&weightConn{
				c: testConn{
					name: "weight-5",
				},
				weight:          5,
				currentWeight:   5,
				efficientWeight: 5,
			},
			&weightConn{
				c: testConn{
					name: "weight-4",
				},
				weight:          4,
				currentWeight:   4,
				efficientWeight: 4,
			},
			&weightConn{
				c: testConn{
					name: "weight-3",
				},
				weight:          3,
				currentWeight:   3,
				efficientWeight: 3,
			},
		},
	}
	res, err := weightBalancer.Pick(balancer.PickInfo{})
	assert.NoError(t, err)
	assert.Equal(t, res.SubConn.(testConn).name, "weight-5")

	res, err = weightBalancer.Pick(balancer.PickInfo{})
	assert.NoError(t, err)
	assert.Equal(t, res.SubConn.(testConn).name, "weight-4")

	res, err = weightBalancer.Pick(balancer.PickInfo{})
	assert.NoError(t, err)
	assert.Equal(t, res.SubConn.(testConn).name, "weight-3")

	res, err = weightBalancer.Pick(balancer.PickInfo{})
	assert.NoError(t, err)
	assert.Equal(t, res.SubConn.(testConn).name, "weight-5")

	res, err = weightBalancer.Pick(balancer.PickInfo{})
	assert.NoError(t, err)
	assert.Equal(t, res.SubConn.(testConn).name, "weight-4")
}
