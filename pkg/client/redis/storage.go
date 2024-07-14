package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func New(redisCfg *Config) (*redis.Client, error) {
	host := fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: redisCfg.Password,
		DB:       redisCfg.DataBase,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return client, nil
}
