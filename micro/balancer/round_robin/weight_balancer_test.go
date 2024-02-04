package round_robin

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/balancer"
	"testing"
)

func Test_BalancerPick(t *testing.T) {
	testCases := []struct {
		name      string
		balance   *Balancer
		wantErr   error
		wantIndex int32
		wantConn  testConn
	}{
		{
			name: "start",
			balance: &Balancer{
				index: -1,
				connects: []balancer.SubConn{
					testConn{
						name: "127.0.0.1:8080",
					},
					testConn{
						name: "127.0.0.1:8081",
					},
				},
			},
			wantIndex: 0,
			wantConn: testConn{
				name: "127.0.0.1:8080",
			},
		},
		{
			name: "end",
			balance: &Balancer{
				index: 1,
				connects: []balancer.SubConn{
					testConn{
						name: "127.0.0.1:8080",
					},
					testConn{
						name: "127.0.0.1:8081",
					},
				},
			},
			wantIndex: 2,
			wantConn: testConn{
				name: "127.0.0.1:8080",
			},
		},
		{
			name: "empty",
			balance: &Balancer{
				index: -1,
			},
			wantErr: balancer.ErrNoSubConnAvailable,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.balance.Pick(balancer.PickInfo{})
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.balance.index, tc.wantIndex)
			assert.Equal(t, tc.wantConn.name, res.SubConn.(testConn).name)
		})
	}
}

type testConn struct {
	name string
	balancer.SubConn
}
