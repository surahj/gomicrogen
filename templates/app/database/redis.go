package database

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func RedisClient() *redis.Client {

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	db := os.Getenv("REDIS_DATABASE_NUMBER")
	auth := os.Getenv("REDIS_PASSWORD")

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
