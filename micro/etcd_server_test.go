package micro

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/sync/errgroup"
	"myProject/micro/grpc_demo/gen"
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
	var eg errgroup.Group
	for i := 0; i < 3; i++ {
		group := "A"
		if i == 1 {
			group = "B"
		}
		etcdServer := NewEtcdServer("user-service", WithRegister(register), WithServerGroup(group))
		defer client.Close()
		gen.RegisterUserServiceServer(etcdServer, &ServerTest{group: group})
		defer etcdServer.Close()
		port := fmt.Sprintf("localhost:809%d", i)
		eg.Go(func() error {
			return etcdServer.Start(port)
		})
	}
	eg.Wait()
}

type ServerTest struct {
	gen.UnimplementedUserServiceServer
	group string
}

func (c *ServerTest) GetById(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error) {
	fmt.Println(c.group)
	return &gen.GetByIdResp{
		User: &gen.User{
			Name: "JunZeXie",
			Id:   5,
		},
	}, nil
}
