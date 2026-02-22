package database

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	db *redis.Client
}

func NewRedis() (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Redis{db: client}, nil
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	cached, err := r.db.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return cached, nil
}

func (r *Redis) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.db.Set(ctx, key, value, ttl).Err()
}
