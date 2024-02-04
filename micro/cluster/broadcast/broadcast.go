package broadcast

import (
	"context"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"myProject/micro/registry"
)

type ClusterBuilder struct {
	r           registry.Registry
	serverName  string
	dialOptions []grpc.DialOption
}

func NewClusterBuilder(serverName string, r registry.Registry, opts ...grpc.DialOption) ClusterBuilder {
	return ClusterBuilder{
		serverName:  serverName,
		r:           r,
		dialOptions: opts,
	}
}

func (c ClusterBuilder) BuildUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if !isBroadCast(ctx) {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		instances, err := c.r.ListServices(ctx, c.serverName)
		if err != nil {
			return err
		}
		var eg errgroup.Group
		for _, ins := range instances {
			addr := ins.Addr
			eg.Go(func() error {
				insCC, er := grpc.Dial(addr, c.dialOptions...)
				if er != nil {
					return er
				}
				return invoker(ctx, method, req, reply, insCC, opts...)
			})
		}
		return eg.Wait()
	}
}

type broadCastKey struct {
}

func UseBroadCast(ctx context.Context) context.Context {
	return context.WithValue(ctx, broadCastKey{}, true)
}

func isBroadCast(ctx context.Context) bool {
	val, ok := ctx.Value(broadCastKey{}).(bool)
	return ok && val
}
