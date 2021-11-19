package redisclient

import (
	"context"
	"github.com/go-redis/redis"
	"time"
)

type Redis struct {
	client   *redis.Client
	Addr     string
	Password string
	DB       int
}

func (r *Redis) Init() error {
	r.client = redis.NewClient(
		&redis.Options{
			Addr:     r.Addr,
			Password: r.Password,
			DB:       r.DB,
		})
	ctx := context.Background()
	_, err := r.client.Ping(ctx).Result()
	return err
}

func (r *Redis) Setex(key string, value string, duration time.Duration) error {
	ctx := context.Background()
	_, err := r.client.Set(ctx, key, value, duration).Result()
	return err
}

func (r *Redis) GetKeyValue(key string) (string, error) {
	ctx := context.Background()
	return r.client.Get(ctx, key).Result()
}

func (r *Redis) GetKeys(pattern string) ([]string, error) {
	ctx := context.Background()
	return r.client.Keys(ctx, pattern).Result()
}

func (r *Redis) IsExists(key string) (bool, error) {
	ctx := context.Background()
	value, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return value == 1, nil
}

func (r *Redis) Delete(key string) error {
	ctx := context.Background()
	_, err := r.client.Del(ctx, key).Result()
	return err
}
