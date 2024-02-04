package fastest

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"myProject/micro"
	"myProject/micro/grpc_demo/gen"
	"myProject/micro/registry/etcd"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	require.NoError(t, err)
	r, err := etcd.NewRegistry(etcdClient)
	require.NoError(t, err)
	var eg errgroup.Group
	for i := 0; i < 3; i++ {
		port := fmt.Sprintf(":809%d", i+1)
		prom := NewPrometheusBuilder("user-service", "mySever", port, "测试")
		server := micro.NewEtcdServer("user-service", micro.WithRegister(r), micro.WithUnaryServerInterceptor(grpc.UnaryInterceptor(prom.Build())))
		us := &Server{}
		gen.RegisterUserServiceServer(server, us)
		require.NoError(t, err)
		eg.Go(func() error {
			return server.Start(port)
		})
	}
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9097", nil)
	}()
	eg.Wait()
	t.Log(err)
}

type Server struct {
	gen.UnimplementedUserServiceServer
}

func (s *Server) GetById(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error) {
	return &gen.GetByIdResp{
		User: &gen.User{
			Name: "hello,world",
		},
	}, nil
}
