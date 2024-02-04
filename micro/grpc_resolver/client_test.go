package grpc_resolver

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"myProject/micro/grpc_demo/gen"
	"testing"
)

func TestClient(t *testing.T) {
	conn, err := grpc.Dial("registry:///localhost:8082", grpc.WithInsecure(), grpc.WithResolvers(&Builder{}))
	assert.NoError(t, err)
	us := gen.NewUserServiceClient(conn)
	res, err := us.GetById(context.Background(), &gen.GetByIdReq{})
	assert.NoError(t, err)
	assert.Equal(t, res.User.Name, "JunZeXie")
	assert.Equal(t, res.User.Id, int64(5))
}
