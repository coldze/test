package sources

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
)

type RedisWrap interface {
	Set(key string, data interface{}, ttl time.Duration) error
	Del(key string) error
	Get(key string) (interface{}, error)
	Close() error
}

type redisWrapImpl struct {
	client *redis.Client
}

func (r *redisWrapImpl) Set(key string, data interface{}, ttl time.Duration) error {
	return r.client.Set(key, data, ttl).Err()
}

func (r *redisWrapImpl) Del(key string) error {
	return r.client.Del(key).Err()
}

func (r *redisWrapImpl) Get(key string) (interface{}, error) {
	return r.client.Get(key).Result()
}

func (r *redisWrapImpl) Close() error {
	return r.client.Close()
}

func (r *redisWrapImpl) Ping() (string, error) {
	return r.client.Ping().Result()
}

func NewRedisWrap(cfg *redis.Options) (RedisWrap, error) {
	client := redis.NewClient(cfg)
	if client == nil {
		return nil, errors.New("internal error - redis client is nil")
	}
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return &redisWrapImpl{
		client: client,
	}, nil
}
