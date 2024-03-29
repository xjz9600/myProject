package opentelemetry

import (
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/metadata"
)

var _ propagation.TextMapCarrier = &metadataSupplier{}

type metadataSupplier struct {
	metadata metadata.MD
}

func (s *metadataSupplier) Get(key string) string {
	values := s.metadata.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (s *metadataSupplier) Set(key string, value string) {
	s.metadata.Set(key, value)
}

func (s *metadataSupplier) Keys() []string {
	out := make([]string, 0, len(s.metadata))
	for key := range s.metadata {
		out = append(out, key)
	}
	return out
}
