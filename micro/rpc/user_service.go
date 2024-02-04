package rpc

import (
	"context"
	"myProject/micro/rpc/serialize/proto/gen"
	"testing"
	"time"
)

type UserService struct {
	GetById      func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
	GetByIdProto func(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error)
}

type GetByIdReq struct {
	Id int
}

type GetByIdResp struct {
	Msg string
}

func (t *UserService) Name() string {
	return "user_service"
}

type UserServiceServer struct {
	Err error
	Msg string
}

func (u *UserServiceServer) GetByIdProto(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error) {
	return &gen.GetByIdResp{
		User: &gen.User{
			Id:   123,
			Name: u.Msg,
		},
	}, u.Err
}
func (u *UserServiceServer) GetById(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error) {
	return &GetByIdResp{
		Msg: u.Msg,
	}, u.Err
}

func (u *UserServiceServer) Name() string {
	return "user_service"
}

type UserServiceTimeout struct {
	Err   error
	Msg   string
	sleep time.Duration
	t     *testing.T
}

func (u *UserServiceTimeout) Name() string {
	return "user_service"
}

func (u *UserServiceTimeout) GetById(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error) {
	if _, ok := ctx.Deadline(); !ok {
		u.t.Fatal("没有设置超时")
	}
	time.Sleep(u.sleep)
	return &GetByIdResp{
		Msg: u.Msg,
	}, u.Err
}
