package cache

import (
	"context"
	"time"

	"github.com/Christian-007/fit-forge/internal/pkg/apperrors"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(options *redis.Options) (*RedisCache, error) {
	client := redis.NewClient(options)
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &RedisCache{
		client: client,
		ctx:    context.Background(),
	}, nil
}

func (r *RedisCache) Get(key string) (any, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return nil, apperrors.ErrRedisKeyNotFound
	}

	if err != nil {
		return nil, err
	}

	return val, nil
}

func (r *RedisCache) Set(key string, value any, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

func (r *RedisCache) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisCache) GetAllHashFields(key string) (map[string]string, error) {
	result, err := r.client.HGetAll(r.ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *RedisCache) SetHash(key string, values ...interface{}) error {
	return r.client.HSet(r.ctx, key, values).Err()
}

func (r *RedisCache) SetExpire(key string, expiration time.Duration) error {
	return r.client.Expire(r.ctx, key, expiration).Err()
}
