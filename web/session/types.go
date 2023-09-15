package session

import (
	"context"
	"errors"
	"net/http"
)

var (
	ErrSessionKeyNotFound = errors.New("session: 找不到 key")
)

// Store 管理session
type Store interface {
	Generate(ctx context.Context, id string) (Session, error)
	Get(ctx context.Context, id string) (Session, error)
	Remove(ctx context.Context, id string) error
	Refresh(ctx context.Context, id string) error
}

// Session 本身
type Session interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, val any) error
	ID() string
}

// Propagator Session写入请求跟从请求中读出
type Propagator interface {
	// Inject 注入sessionID到请求中
	Inject(writer http.ResponseWriter, id string) error
	// Extract 从请求中拿到sessionID
	Extract(req *http.Request) (string, error)
	Remove(writer http.ResponseWriter, id string) error
}
