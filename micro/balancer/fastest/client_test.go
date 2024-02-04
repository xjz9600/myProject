package fastest

import (
	"context"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
	"myProject/micro"
	"myProject/micro/grpc_demo/gen"
	"myProject/micro/registry/etcd"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	require.NoError(t, err)
	r, err := etcd.NewRegistry(etcdClient)
	require.NoError(t, err)

	client := micro.NewEtcdClient(micro.ClientWithInsecure(),
		micro.ClientWithRegister(r, time.Second*3),
		micro.ClientWithBalancer("prometheus", &BalancerBuilder{
			Point:    "http://localhost:9090",
			Duration: time.Second * 3,
			Query:    "mySever{kind=\"user-service\",quantile=\"0.5\"}",
		}))
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10000)
	defer cancel()
	cc, err := client.Dial(ctx, "registry:///user-service")
	uc := gen.NewUserServiceClient(cc)
	for i := 0; i < 10; i++ {
		resp, err := uc.GetById(ctx, &gen.GetByIdReq{Id: 13})
		require.NoError(t, err)
		t.Log(resp)
		time.Sleep(5 * time.Second)
	}
}
