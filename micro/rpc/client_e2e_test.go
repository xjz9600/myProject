package rpc

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"myProject/micro/rpc/compress"
	"myProject/micro/rpc/serialize/proto"
	"myProject/micro/rpc/serialize/proto/gen"
	"testing"
	"time"
)

func TestInitService(t *testing.T) {
	server := NewServer()
	service := &UserServiceServer{}
	server.Register(service)
	go func() {
		err := server.Serve("tcp", ":8089")
		t.Log(err)
	}()
	time.Sleep(time.Second * 3)
	us := &UserService{}
	client, err := NewClient("tcp", ":8089")
	require.NoError(t, err)
	err = client.InitService(us)
	require.NoError(t, err)
	testCases := []struct {
		name string
		mock func()

		wantErr  error
		wantResp *GetByIdResp
	}{
		{
			name: "no error",
			mock: func() {
				service.Err = nil
				service.Msg = "hello, world"
			},
			wantResp: &GetByIdResp{
				Msg: "hello, world",
			},
		},
		{
			name: "error",
			mock: func() {
				service.Msg = ""
				service.Err = errors.New("mock error")
			},
			wantResp: &GetByIdResp{},
			wantErr:  errors.New("mock error"),
		},

		{
			name: "both",
			mock: func() {
				service.Msg = "hello, world"
				service.Err = errors.New("mock error")
			},
			wantResp: &GetByIdResp{
				Msg: "hello, world",
			},
			wantErr: errors.New("mock error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			resp, err := us.GetById(context.Background(), &GetByIdReq{Id: 123})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantResp, resp)

		})
	}
}

func TestInitServiceProto(t *testing.T) {
	server := NewServer()
	server.RegisterSerialize(proto.Serializer{})
	service := &UserServiceServer{}
	server.Register(service)
	go func() {
		err := server.Serve("tcp", ":8089")
		t.Log(err)
	}()
	time.Sleep(time.Second * 3)
	us := &UserService{}
	client, err := NewClient("tcp", ":8089", WithSerializer(proto.Serializer{}))
	require.NoError(t, err)
	err = client.InitService(us)
	require.NoError(t, err)
	testCases := []struct {
		name string
		mock func()

		wantErr  error
		wantResp *gen.GetByIdResp
	}{
		{
			name: "no error",
			mock: func() {
				service.Err = nil
				service.Msg = "hello, world"
			},
			wantResp: &gen.GetByIdResp{
				User: &gen.User{
					Name: "hello, world",
				},
			},
		},
		{
			name: "error",
			mock: func() {
				service.Msg = ""
				service.Err = errors.New("mock error")
			},
			wantResp: &gen.GetByIdResp{
				User: &gen.User{},
			},
			wantErr: errors.New("mock error"),
		},

		{
			name: "both",
			mock: func() {
				service.Msg = "hello, world"
				service.Err = errors.New("mock error")
			},
			wantResp: &gen.GetByIdResp{
				User: &gen.User{
					Name: "hello, world",
				},
			},
			wantErr: errors.New("mock error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			resp, err := us.GetByIdProto(context.Background(), &gen.GetByIdReq{Id: 123})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantResp.User.Name, resp.User.Name)

		})
	}
}

func TestOneway(t *testing.T) {
	server := NewServer()
	service := &UserServiceServer{}
	server.Register(service)
	go func() {
		err := server.Serve("tcp", ":8089")
		t.Log(err)
	}()
	time.Sleep(time.Second * 3)

	usClient := &UserService{}
	client, err := NewClient("tcp", ":8089")
	require.NoError(t, err)
	err = client.InitService(usClient)
	require.NoError(t, err)
	testCases := []struct {
		name string
		mock func()

		wantErr  error
		ctx      context.Context
		wantResp *GetByIdResp
	}{
		{
			name: "oneway",
			mock: func() {
				service.Err = errors.New("mock error")
				service.Msg = "hello, world"
			},
			ctx:     CtxWithOneway(context.Background()),
			wantErr: errors.New("mirco：这是一个 oneway 调用，你不应该处理结果"),
		},
		{
			name: "not-one-way",
			mock: func() {
				service.Err = errors.New("mock error")
				service.Msg = "hello, world"
			},
			ctx:     context.Background(),
			wantErr: errors.New("mock error"),
			wantResp: &GetByIdResp{
				Msg: "hello, world",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			resp, er := usClient.GetById(tc.ctx, &GetByIdReq{Id: 123})
			assert.Equal(t, tc.wantErr, er)
			assert.Equal(t, tc.wantResp, resp)
		})
	}
}

func TestTimeout(t *testing.T) {
	server := NewServer()
	service := &UserServiceTimeout{t: t, sleep: time.Second * 2}
	server.Register(service)
	go func() {
		err := server.Serve("tcp", ":8089")
		t.Log(err)
	}()
	time.Sleep(time.Second * 3)
	usClient := &UserService{}
	client, err := NewClient("tcp", ":8089")
	require.NoError(t, err)
	err = client.InitService(usClient)
	require.NoError(t, err)
	testCases := []struct {
		name     string
		mock     func()
		wantErr  error
		ctx      func() context.Context
		wantResp *GetByIdResp
	}{
		{
			name: "timeout",
			mock: func() {
				service.Err = errors.New("mock error")
				service.Msg = "hello, world"
			},
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), time.Second)
				return ctx
			},
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			resp, er := usClient.GetById(tc.ctx(), &GetByIdReq{Id: 123})
			assert.Equal(t, tc.wantErr, er)
			assert.Equal(t, tc.wantResp, resp)
		})
	}
}

func TestCompress(t *testing.T) {
	server := NewServer()
	service := &UserServiceServer{}
	server.RegisterCompress(compress.NoCompress{})
	server.Register(service)
	go func() {
		err := server.Serve("tcp", ":8089")
		t.Log(err)
	}()
	time.Sleep(time.Second * 3)
	us := &UserService{}
	client, err := NewClient("tcp", ":8089", WithCompress(compress.NoCompress{}))
	require.NoError(t, err)
	err = client.InitService(us)
	require.NoError(t, err)
	testCases := []struct {
		name string
		mock func()

		wantErr  error
		wantResp *GetByIdResp
	}{
		{
			name: "no error",
			mock: func() {
				service.Err = nil
				service.Msg = "hello, world"
			},
			wantResp: &GetByIdResp{
				Msg: "hello, world",
			},
		},
		{
			name: "error",
			mock: func() {
				service.Msg = ""
				service.Err = errors.New("mock error")
			},
			wantResp: &GetByIdResp{},
			wantErr:  errors.New("mock error"),
		},

		{
			name: "both",
			mock: func() {
				service.Msg = "hello, world"
				service.Err = errors.New("mock error")
			},
			wantResp: &GetByIdResp{
				Msg: "hello, world",
			},
			wantErr: errors.New("mock error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			resp, err := us.GetById(context.Background(), &GetByIdReq{Id: 123})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantResp, resp)

		})
	}
}
