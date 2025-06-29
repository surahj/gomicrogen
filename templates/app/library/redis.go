package library

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"time"
)

func GetRedisKey(conn *redis.Client, key string) (string, error) {

	//BOOKING:CODE

	//AUTHORIZATION
	//if strings.HasPrefix(key,"PROFILE:") || strings.HasPrefix(key,"BOOKING:") || strings.HasPrefix(key,"AUTHORIZATION:")  {

	var data string
	data, err := conn.Get(key).Result()
	if err != nil {

		return data, fmt.Errorf("error getting key %s: %v", key, err)
	}

	return data, err
	//}

	//return "",errors.New("redis stopped")

}

func SetRedisKey(conn *redis.Client, key string, value string) error {

	_, err := conn.Set(key, value, time.Second*time.Duration(0)).Result()
	if err != nil {

		v := string(value)

		if len(v) > 15 {

			v = v[0:12] + "..."
		}

		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	return err
}

func SetRedisKeyWithExpiry(conn *redis.Client, key string, value string, seconds int) error {

	_, err := conn.Set(key, value, time.Second*time.Duration(seconds)).Result()
	if err != nil {

		v := string(value)

		if len(v) > 15 {

			v = v[0:12] + "..."
		}

		log.Printf("error saving redisKey %s error %s", key, err.Error())
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}

	return err
}

func IncRedisKey(conn *redis.Client, key string) (int64, error) {

	var data int64
	data, err := conn.Incr(key).Result()

	if err != nil {

		return data, fmt.Errorf("error getting key %s: %v", key, err)
	}

	return data, err
}