package provider

import (
	"github.com/go-redis/redis"
)

// CreateRedisDB returns a connected redisdb.
func CreateRedisDB() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
