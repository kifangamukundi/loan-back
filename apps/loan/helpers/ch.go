package helpers

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"time"

	"github.com/kifangamukundi/gm/libs/rates"

	"github.com/gin-gonic/gin"
)

func StoreCache(redisCache *rates.RedisCache, cacheKey string, data interface{}, ttl time.Duration) {
	response, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal data for key %s: %v", cacheKey, err)
		return
	}

	err = redisCache.SetCache(cacheKey, string(response), ttl)
	if err != nil {
		log.Printf("Failed to cache data for key %s: %v", cacheKey, err)
	} else {
		log.Printf("Cache set: Data successfully stored for key %s", cacheKey)
	}
}

func CheckCache(redisCache *rates.RedisCache, cacheKey string, c *gin.Context, returnFunc func(*gin.Context, interface{})) bool {
	cachedData, err := redisCache.CheckCache(cacheKey)
	if err == nil && cachedData != "" {
		log.Println("Cache hit: Returning cached data")

		var decodedJSON string
		if err := json.Unmarshal([]byte(cachedData), &decodedJSON); err != nil {
			log.Printf("Failed to unmarshal the outer string: %v", err)
			return false
		}

		var innerDecoded map[string]interface{}
		if err := json.Unmarshal([]byte(decodedJSON), &innerDecoded); err != nil {
			log.Printf("Failed to unmarshal the inner string: %v", err)
			return false
		} else {
			returnFunc(c, innerDecoded)
		}
		return true
	}

	return false
}

func StoreBinaryCache(redisCache *rates.RedisCache, cacheKey string, data interface{}, ttl time.Duration) {
	binaryData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal data for key %s: %v", cacheKey, err)
		return
	}

	err = redisCache.SetBinaryCache(cacheKey, binaryData, ttl)
	if err != nil {
		log.Printf("Failed to cache data for key %s: %v", cacheKey, err)
	} else {
		log.Printf("Cache set: Data successfully stored for key %s", cacheKey)
	}
}

func CheckBinaryCache(redisCache *rates.RedisCache, cacheKey string, c *gin.Context, returnFunc func(*gin.Context, interface{})) bool {
	cachedData, err := redisCache.CheckBinaryCache(cacheKey)
	if err == nil && cachedData != nil && len(cachedData) > 0 {
		log.Println("Cache hit: Returning cached data")

		var decodedJSON string
		if err := json.Unmarshal(cachedData, &decodedJSON); err != nil {
			return false
		}

		decodedBytes, err := base64.StdEncoding.DecodeString(decodedJSON)
		if err != nil {
			log.Printf("Failed to decode Base64 data: %v", err)
			return false
		}

		var innerDecoded map[string]interface{}
		if err := json.Unmarshal(decodedBytes, &innerDecoded); err != nil {
			log.Printf("Failed to unmarshal the inner JSON data: %v", err)
			return false
		} else {
			returnFunc(c, innerDecoded)
		}
		return true
	}

	return false
}
