package micro

import (
	"context"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"myProject/micro/grpc_demo/gen"
	"myProject/micro/registry"
	"myProject/micro/registry/etcd"
	"testing"
)

func TestServer(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	assert.NoError(t, err)
	register, err := etcd.NewRegistry(client)
	assert.NoError(t, err)
	etcdServer := registry.NewEtcdServer("user-service", registry.WithRegister(register))
	defer client.Close()
	gen.RegisterUserServiceServer(etcdServer, &ServerTest{})
	defer etcdServer.Close()
	etcdServer.Start("localhost:8085")
}

type ServerTest struct {
	gen.UnimplementedUserServiceServer
}

func (c *ServerTest) GetById(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error) {
	return &gen.GetByIdResp{
		User: &gen.User{
			Name: "JunZeXie",
			Id:   5,
		},
	}, nil
}
