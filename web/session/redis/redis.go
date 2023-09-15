package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"myProject/web/session"
	"time"
)

type Store struct {
	redis      redis.Cmdable
	prefix     string
	expiration time.Duration
}

type storeOpt func(*Store)

func NewStore(redis redis.Cmdable, opts ...storeOpt) *Store {
	res := &Store{
		redis:      redis,
		prefix:     "session",
		expiration: 10 * time.Minute,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func redisKey(prefix, id string) string {
	return fmt.Sprintf("%s-%s", prefix, id)
}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	redisKey := redisKey(s.prefix, id)
	_, err := s.redis.HSet(ctx, redisKey, id, id).Result()
	if err != nil {
		return nil, err
	}
	_, err = s.redis.Expire(ctx, redisKey, s.expiration).Result()
	if err != nil {
		return nil, err
	}
	session := &Session{
		redis:    s.redis,
		redisKey: redisKey,
	}
	return session, nil
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	redisKey := redisKey(s.prefix, id)
	cnt, err := s.redis.Exists(ctx, redisKey).Result()
	if err != nil {
		return nil, err
	}
	// 没有的话返回为0
	if cnt != 1 {
		return nil, session.ErrSessionKeyNotFound
	}
	return &Session{
		redis:    s.redis,
		redisKey: redisKey,
	}, nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	redisKey := redisKey(s.prefix, id)
	cnt, err := s.redis.Del(ctx, redisKey).Result()
	if err != nil {
		return err
	}
	if cnt != 1 {
		return session.ErrSessionKeyNotFound
	}
	return nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	redisKey := redisKey(s.prefix, id)
	ok, err := s.redis.Expire(ctx, redisKey, s.expiration).Result()
	if err != nil {
		return err
	}
	if !ok {
		return session.ErrSessionKeyNotFound
	}
	return nil
}

type Session struct {
	redis    redis.Cmdable
	redisKey string
}

func (s *Session) Get(ctx context.Context, key string) (any, error) {
	val, err := s.redis.HGet(ctx, s.redisKey, key).Result()
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (s *Session) Set(ctx context.Context, key string, val any) error {
	const lua = `
if redis.call("exists", KEYS[1])
then
	return redis.call("hset", KEYS[1], ARGV[1], ARGV[2])
else
	return -1
end
`
	res, err := s.redis.Eval(ctx, lua, []string{s.redisKey}, key, val).Int()
	if err != nil {
		return err
	}
	if res != 1 {
		return session.ErrSessionKeyNotFound
	}
	return nil
}

func (s *Session) ID() string {
	return s.redisKey
}
