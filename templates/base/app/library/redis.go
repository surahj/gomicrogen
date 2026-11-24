package library

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func GetRedisKey(ctx context.Context, conn *redis.Client, key string) (string, error) {

	data, err := conn.Get(ctx, key).Result()
	if err != nil {

		return data, fmt.Errorf("error getting key %s: %v", key, err)
	}

	return data, nil
}

func SetRedisKey(ctx context.Context, conn *redis.Client, key string, value string) error {

	_, err := conn.Set(ctx, key, value, 0).Result()
	if err != nil {

		return fmt.Errorf("error setting key %s to %s: %v", key, truncate(value), err)
	}

	return nil
}

func SetRedisKeyWithExpiry(ctx context.Context, conn *redis.Client, key string, value string, seconds int) error {

	_, err := conn.Set(ctx, key, value, time.Second*time.Duration(seconds)).Result()
	if err != nil {

		log.Printf("error saving redisKey %s error %s", key, err.Error())
		return fmt.Errorf("error setting key %s to %s: %v", key, truncate(value), err)
	}

	return nil
}

func IncRedisKey(ctx context.Context, conn *redis.Client, key string) (int64, error) {

	data, err := conn.Incr(ctx, key).Result()
	if err != nil {

		return data, fmt.Errorf("error incrementing key %s: %v", key, err)
	}

	return data, nil
}

// truncate keeps values out of error messages beyond a short prefix
func truncate(value string) string {

	if len(value) > 15 {

		return value[0:12] + "..."
	}

	return value
}
