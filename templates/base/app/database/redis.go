package database

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// RedisClient returns the service-local cache client.
func RedisClient() *redis.Client {

	return newClient("REDIS_HOST", "REDIS_PORT", "REDIS_DATABASE_NUMBER", "REDIS_PASSWORD")
}

// GlobalRedisClient returns the shared platform-wide client. Auth tokens are
// written here by identity-service, so auth.Authenticate must be given this
// client and not the service-local one.
func GlobalRedisClient() *redis.Client {

	return newClient("GLOBAL_REDIS_HOST", "GLOBAL_REDIS_PORT", "GLOBAL_REDIS_DATABASE_NUMBER", "GLOBAL_REDIS_PASSWORD")
}

func newClient(hostKey, portKey, dbKey, passwordKey string) *redis.Client {

	host := os.Getenv(hostKey)
	port := os.Getenv(portKey)
	db := os.Getenv(dbKey)
	auth := os.Getenv(passwordKey)

	dbNumber, err := strconv.Atoi(db)
	if err != nil {
		dbNumber = 1
	}

	uri := fmt.Sprintf("%s:%s", host, port)

	opts := redis.Options{
		MinIdleConns: 10,
		IdleTimeout:  60 * time.Second,
		PoolSize:     1000,
		Addr:         uri,
		DB:           dbNumber, // use default DB
	}

	if len(auth) > 0 {

		opts.Password = auth
	}

	client := redis.NewClient(&opts)

	return client
}
