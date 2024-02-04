package broadcastfast

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"myProject/micro"
	"myProject/micro/grpc_demo/gen"
	"myProject/micro/registry/etcd"
	"testing"
	"time"
)

func TestUseBroadCast(t *testing.T) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	require.NoError(t, err)
	r, err := etcd.NewRegistry(etcdClient)
	require.NoError(t, err)

	var eg errgroup.Group
	var servers []*UserServiceServer
	for i := 0; i < 3; i++ {
		server := micro.NewEtcdServer("user-service", micro.WithRegister(r))
		require.NoError(t, err)
		us := &UserServiceServer{
			idx: i,
		}
		servers = append(servers, us)
		gen.RegisterUserServiceServer(server, us)
		// 启动 8081,8082, 8083 三个端口
		port := fmt.Sprintf(":808%d", i+1)
		eg.Go(func() error {
			return server.Start(port)
		})
	}
	time.Sleep(time.Second * 3)
	client := micro.NewEtcdClient(micro.ClientWithInsecure(),
		micro.ClientWithRegister(r, time.Second*3))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer func() {
		cancel()
	}()
	require.NoError(t, err)
	ctx, respChan := UseBroadCast(ctx)
	go func() {
		for rc := range respChan {
			t.Log(rc.Reply)
		}
	}()
	bd := NewClusterBuilder("user-service", r, grpc.WithInsecure())
	cc, err := client.Dial(ctx, "registry:///user-service", grpc.WithUnaryInterceptor(bd.BuildUnaryClientInterceptor()))
	require.NoError(t, err)
	uc := gen.NewUserServiceClient(cc)
	resp, err := uc.GetById(ctx, &gen.GetByIdReq{Id: 13})
	require.NoError(t, err)
	t.Log(resp)
}

type UserServiceServer struct {
	idx int
	cnt int
	gen.UnimplementedUserServiceServer
}

func (s *UserServiceServer) GetById(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error) {
	s.cnt++
	return &gen.GetByIdResp{
		User: &gen.User{
			Name: fmt.Sprintf("hello, world %d", s.idx),
		},
	}, nil
}
