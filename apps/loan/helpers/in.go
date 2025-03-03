package helpers

import (
	"log"

	"github.com/kifangamukundi/gm/libs/rates"
)

// InitializeRedisCache initializes the Redis cache and returns the client
func InitializeRedisCache() *rates.RedisCache {
	redisCache := rates.NewRedisCache()
	if redisCache == nil {
		log.Fatal("Failed to connect to Redis. Exiting application.")
	}
	return redisCache
}
