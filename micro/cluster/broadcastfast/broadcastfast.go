package broadcastfast

import (
	"context"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"myProject/micro/registry"
	"reflect"
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
		resp, ok := isBroadcast(ctx)
		if !ok {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		instances, err := c.r.ListServices(ctx, c.serverName)
		if err != nil {
			return err
		}
		replyType := reflect.TypeOf(reply).Elem()
		var eg errgroup.Group
		for _, in := range instances {
			addr := in.Addr
			eg.Go(func() error {
				conn, er := grpc.Dial(addr, c.dialOptions...)
				var newResp Resp
				if er != nil {
					newResp = Resp{
						Err: er,
					}
				} else {
					newReply := reflect.New(replyType).Interface()
					er = invoker(ctx, method, req, newReply, conn, opts...)
					newResp = Resp{
						Err:   er,
						Reply: newReply,
					}
				}
				select {
				case resp <- newResp:
					return nil
				default:
					return nil
				}
			})
		}
		return eg.Wait()
	}
}

type broadcastKey struct {
}

func UseBroadCast(ctx context.Context) (context.Context, <-chan Resp) {
	ch := make(chan Resp)
	return context.WithValue(ctx, broadcastKey{}, ch), ch
}

func isBroadcast(ctx context.Context) (chan<- Resp, bool) {
	val, ok := ctx.Value(broadcastKey{}).(chan Resp)
	return val, ok
}

type Resp struct {
	Err   error
	Reply any
}
