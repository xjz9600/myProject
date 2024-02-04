package balancer

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)

type Filter func(info balancer.PickInfo, addr resolver.Address) bool

type GroupFilterBuilder struct {
}

func (g GroupFilterBuilder) Build() Filter {
	return func(info balancer.PickInfo, addr resolver.Address) bool {
		target := addr.Attributes.Value("group")
		input := info.Ctx.Value("group")
		return target == input
	}
}
