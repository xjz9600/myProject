package opentelemetry

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"myProject/micro/grpc_demo/gen"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	initZipkin(t)
	lis, err := net.Listen("tcp", "localhost:8082")
	defer lis.Close()
	assert.NoError(t, err)
	builder := &TelServerBuilder{}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(builder.Build()))
	gen.RegisterUserServiceServer(grpcServer, &Client{})
	grpcServer.Serve(lis)
}

type Client struct {
	gen.UnimplementedUserServiceServer
}

func (c *Client) GetById(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error) {
	return &gen.GetByIdResp{
		User: &gen.User{
			Name: "JunZeXie",
			Id:   5,
		},
	}, nil
}
