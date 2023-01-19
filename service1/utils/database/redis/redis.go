package redis

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis/v9"
)

type RedisClient struct {
	Redis *redis.Client
}

var redisConnect = &RedisClient{}

func Connect2Redis(host, port, pass string) (*RedisClient, error) {
	var connectionStr string

	if host == "" && port == "" {
		return nil, errors.New("cannot estabished the redis connection")
	}

	if port == "REDIS_PORT" {
		port = "6379"
	}

	connectionStr = fmt.Sprintf("%v:%v", host, port)

	//connect redis
	dial := redis.NewClient(&redis.Options{
		Addr:     connectionStr,
		Password: pass, // no password set
		DB:       0,    // use default DB
	})

	redisConnect.Redis = dial
	return redisConnect, nil
}

func DisconnectRedis(db *redis.Client) {
	db.Close()
}
