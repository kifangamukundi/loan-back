package rates

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

var redisInitialized bool

func InitRedis() {
	if redisInitialized {
		log.Println("Redis is already initialized. Skipping reinitialization.")
		return
	}
	env := os.Getenv("GIN_MODE")
	var options *redis.Options

	if env == "true" {
		options = &redis.Options{
			Addr:     os.Getenv("REDIS_ADDRESS"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		}
		log.Println("Using internal Redis instance...")
	} else {
		options = &redis.Options{
			Addr:     "",
			Password: "",
			DB:       0,
		}
		log.Println("Using Local Redis instance...")
	}

	client := redis.NewClient(options)

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	} else {
		log.Println("Successfully connected to Redis")
	}

	RedisClient = client
	redisInitialized = true
}
