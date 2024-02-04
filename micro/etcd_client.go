package registry

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	resolver2 "google.golang.org/grpc/resolver"
	"time"
)

type EtcdClient struct {
	insecure bool
	resolver resolver2.Builder
	balancer balancer.Builder
}

type EtcdClientOpt func(client *EtcdClient)

func NewEtcdClient(opts ...EtcdClientOpt) *EtcdClient {
	res := &EtcdClient{}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func (e *EtcdClient) Dial(ctx context.Context, severName string) (*grpc.ClientConn, error) {
	var opt []grpc.DialOption
	if e.insecure {
		opt = append(opt, grpc.WithInsecure())
	}
	if e.resolver != nil {
		opt = append(opt, grpc.WithResolvers(e.resolver))
	}
	if e.balancer != nil {
		opt = append(opt, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, e.balancer.Name())))
	}
	return grpc.DialContext(ctx, severName, opt...)
}

func ClientWithInsecure() EtcdClientOpt {
	return func(client *EtcdClient) {
		client.insecure = true
	}
}

func ClientWithRegister(r Registry, timeout time.Duration) EtcdClientOpt {
	return func(client *EtcdClient) {
		resolverBuilder := NewResolverBuilder(r, timeout)
		client.resolver = resolverBuilder
	}
}

func ClientWithBalancer(pb base.PickerBuilder) EtcdClientOpt {
	return func(client *EtcdClient) {
		builder := base.NewBalancerBuilder("DEMO_ROUND_ROBIN", pb, base.Config{HealthCheck: true})
		balancer.Register(builder)
		client.balancer = builder
	}
}
