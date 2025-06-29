package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
)

func CheckConnectionStatus(ctx context.Context, db *sql.DB) (int, map[string]string) {

	res := make(map[string]string)
	status := http.StatusOK

	err := db.PingContext(ctx)
	if err == nil {

		res["database"] = "database - sent successful ping"

	} else {

		res["database"] = fmt.Sprintf("database error - %s", err.Error())
		status = http.StatusInternalServerError
	}

	redisClient := RedisClient()
	defer redisClient.Close()

	resp, err := redisClient.Ping().Result()
	if err == nil {

		res["redis"] = resp

	} else {

		res["redis"] = fmt.Sprintf("redis error - %s", err.Error())
		status = http.StatusInternalServerError

	}

	return status, res

}

