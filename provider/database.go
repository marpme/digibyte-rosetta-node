package provider

import (
	"github.com/go-redis/redis"
	"github.com/marpme/digibyte-rosetta-node/configuration"
)

// DBType is a strict type to determine different db layers from each other
type DBType = int

const (
	// BlockDB identifies the database and it's usage for blocks only
	BlockDB DBType = 0
	// TxDB identifies the database and it's usage for transactions only
	TxDB DBType = 1
)

// CreateRedisDB returns a connected redisdb.
func CreateRedisDB(cfg *configuration.Config, typeOfDB DBType) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       typeOfDB,
	})
}
