package registry

import (
	"context"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"time"
)

type ResolverBuilder struct {
	r       Registry
	timeout time.Duration
}

func NewResolverBuilder(r Registry, timeout time.Duration) *ResolverBuilder {
	return &ResolverBuilder{
		r:       r,
		timeout: timeout,
	}
}

func (b *ResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	res := &Resolver{
		cc:      cc,
		r:       b.r,
		target:  target,
		timeout: b.timeout,
		close:   make(chan struct{}),
	}
	res.resolve()
	go res.watch()
	return res, nil
}

func (b *ResolverBuilder) Scheme() string {
	return "registry"
}

type Resolver struct {
	cc      resolver.ClientConn
	r       Registry
	target  resolver.Target
	timeout time.Duration
	close   chan struct{}
}

func (r *Resolver) watch() {
	event, err := r.r.Subscribe(r.target.Endpoint())
	if err != nil {
		r.cc.ReportError(err)
		return
	}
	for {
		select {
		case <-event:
			r.resolve()
		case <-r.close:
			return
		}
	}
}

func (r *Resolver) resolve() {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	instances, err := r.r.ListServices(ctx, r.target.Endpoint())
	if err != nil {
		r.cc.ReportError(err)
		return
	}
	address := make([]resolver.Address, len(instances))
	for _, in := range instances {
		addr := resolver.Address{
			Addr:       in.Addr,
			Attributes: attributes.New("weight", in.Weight).WithValue("group", in.Group),
		}
		address = append(address, addr)
	}
	err = r.cc.UpdateState(resolver.State{
		Addresses: address,
	})
	if err != nil {
		r.cc.ReportError(err)
	}
}

func (r *Resolver) ResolveNow(options resolver.ResolveNowOptions) {
	r.resolve()
}

func (r *Resolver) Close() {
	close(r.close)
}
