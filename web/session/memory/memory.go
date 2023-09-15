package memory

import (
	"context"
	"github.com/patrickmn/go-cache"
	"myProject/web/session"
	"sync"
	"time"
)

type Store struct {
	session         *cache.Cache
	expiration      time.Duration
	cleanupInterval time.Duration
}

type storeOpt func(*Store)

func NewStore(opts ...storeOpt) *Store {
	res := &Store{
		session: cache.New(time.Minute, time.Second),
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	newSession := &Session{
		id: id,
	}
	s.session.Set(id, newSession, cache.DefaultExpiration)
	return newSession, nil
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	sess, ok := s.session.Get(id)
	if !ok {
		return nil, session.ErrSessionKeyNotFound
	}
	return sess.(session.Session), nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	s.session.Delete(id)
	return nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	val, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	s.session.Set(id, val, cache.DefaultExpiration)
	return nil
}

type Session struct {
	id     string
	values sync.Map
}

func (s *Session) Get(ctx context.Context, key string) (any, error) {
	val, ok := s.values.Load(key)
	if !ok {
		return nil, session.ErrSessionKeyNotFound
	}
	return val, nil
}

func (s *Session) Set(ctx context.Context, key string, val any) error {
	s.values.Store(key, val)
	return nil
}

func (s *Session) ID() string {
	return s.id
}
