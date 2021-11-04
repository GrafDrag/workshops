package redis

import (
	"calendar/internal/config"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()
)

type session struct {
	redis *redis.Client
}

func NewSession(config config.SessionConfig) (*session, error) {
	rdb, err := newRedisConnect(config.Host, config.Port, config.Password, config.DB)
	if err != nil {
		return nil, err
	}
	return &session{
		redis: rdb,
	}, nil
}

func (s *session) Get(key string) (string, error) {
	v, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed find key %s", key)
	}

	return v, nil
}

func (s *session) Set(key string, value string) error {
	ctx = context.Background()
	return s.redis.Set(ctx, key, value, 0).Err()
}

func (s *session) Remove(key string) (string, error) {
	val, err := s.Get(key)
	if err != nil {
		return "", fmt.Errorf("failed delete key %s", key)
	}
	pipe := s.redis.Pipeline()

	pipe.Del(ctx, key)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return "", fmt.Errorf("failed delete key %s", key)
	}

	return val, nil
}

func (s *session) Flash() error {
	pipe := s.redis.Pipeline()

	pipe.FlushAll(ctx)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed flash")
	}

	return nil
}

func newRedisConnect(host string, port int, password string, db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password, // no password set
		DB:       db,       // use default DB
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}