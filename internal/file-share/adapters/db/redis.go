package db

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDB struct {
	Conn *redis.Client
}

func NewRedis(addr string) (*RedisDB, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisDB{Conn: rdb}, nil
}
