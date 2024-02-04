package micro

import (
	"context"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"myProject/micro/balancer"
	"myProject/micro/balancer/round_robin"
	"myProject/micro/grpc_demo/gen"
	"myProject/micro/registry/etcd"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	defer client.Close()
	assert.NoError(t, err)
	register, err := etcd.NewRegistry(client)
	assert.NoError(t, err)
	groupFilter := balancer.GroupFilterBuilder{}
	etcdClient := NewEtcdClient(ClientWithInsecure(), ClientWithRegister(register, time.Second*10), ClientWithBalancer("group", &round_robin.BalancerBuilder{Filter: groupFilter.Build()}))
	conn, err := etcdClient.Dial(context.Background(), "registry:///user-service")
	defer conn.Close()
	assert.NoError(t, err)
	us := gen.NewUserServiceClient(conn)
	for i := 0; i < 10; i++ {
		res, err := us.GetById(context.WithValue(context.Background(), "group", "A"), &gen.GetByIdReq{})
		assert.NoError(t, err)
		assert.Equal(t, res.User.Name, "JunZeXie")
		assert.Equal(t, res.User.Id, int64(5))
	}
}
