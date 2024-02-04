package micro

import (
	"context"
	"google.golang.org/grpc"
	"myProject/micro/registry"
	"net"
	"time"
)

type EtcdServer struct {
	register        registry.Registry
	registerTimeout time.Duration
	name            string
	listener        net.Listener
	*grpc.Server
	weight uint32
	group  string
}

type etcdServerOpt func(server *EtcdServer)

func WithServerWeight(weight uint32) etcdServerOpt {
	return func(server *EtcdServer) {
		server.weight = weight
	}
}

func WithUnaryServerInterceptor(serverOption ...grpc.ServerOption) etcdServerOpt {
	return func(server *EtcdServer) {
		server.Server = grpc.NewServer(serverOption...)
	}
}

func WithServerGroup(group string) etcdServerOpt {
	return func(server *EtcdServer) {
		server.group = group
	}
}

func WithRegister(register registry.Registry) etcdServerOpt {
	return func(server *EtcdServer) {
		server.register = register
	}
}

func NewEtcdServer(name string, opts ...etcdServerOpt) *EtcdServer {
	res := &EtcdServer{
		name:            name,
		Server:          grpc.NewServer(),
		registerTimeout: time.Second * 10,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func (e *EtcdServer) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	e.listener = listener
	if e.register != nil {
		ctx, cancel := context.WithTimeout(context.Background(), e.registerTimeout)
		defer cancel()
		err = e.register.Register(ctx, registry.ServiceInstance{
			Name:   e.name,
			Addr:   listener.Addr().String(),
			Weight: e.weight,
			Group:  e.group,
		})
		if err != nil {
			return err
		}
	}
	defer func() {
		_ = e.Close()
	}()
	return e.Serve(listener)
}

func (e *EtcdServer) Close() error {
	if e.register != nil {
		err := e.register.Close()
		if err != nil {
			return err
		}
	}
	e.GracefulStop()
	return nil
}
