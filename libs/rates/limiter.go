package rates

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	redisStore "github.com/ulule/limiter/v3/drivers/store/redis"
)

func CreateRateLimiter(format string) gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted(format)
	if err != nil {
		panic(err)
	}

	rdb := RedisClient

	store, err := redisStore.NewStore(rdb)
	if err != nil {
		log.Fatalf("Could not create Redis store: %v", err)
	} else {
		log.Println("Successfully created Redis store for rate limiter")
	}

	instance := limiter.New(store, rate)

	return ginlimiter.NewMiddleware(instance, ginlimiter.WithErrorHandler(func(c *gin.Context, err error) {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error":   "Too many requests, please try again later",
			"details": err.Error(),
		})
		c.Abort()
	}))
}
