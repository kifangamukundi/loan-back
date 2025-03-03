package rates

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	CheckCache(cacheKey string) (string, error)
	SetCache(cacheKey string, data interface{}, ttl time.Duration) error
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache() *RedisCache {
	client := RedisClient

	return &RedisCache{client: client}
}

func (r *RedisCache) SetCache(cacheKey string, data interface{}, ttl time.Duration) error {
	response, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.client.Set(context.Background(), cacheKey, string(response), ttl).Err()
}

func (r *RedisCache) CheckCache(cacheKey string) (string, error) {
	return r.client.Get(context.Background(), cacheKey).Result()
}

func (r *RedisCache) SetBinaryCache(cacheKey string, data interface{}, ttl time.Duration) error {
	binaryData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.client.Set(context.Background(), cacheKey, binaryData, ttl).Err()
}

func (r *RedisCache) CheckBinaryCache(cacheKey string) ([]byte, error) {
	return r.client.Get(context.Background(), cacheKey).Bytes()
}
