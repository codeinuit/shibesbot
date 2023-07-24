package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisDB struct {
	client *redis.Client
}

type RedisOptions struct {
	Port     int32
	Address  string
	Password string
	Database string
}

// NewRedisCache returns a Redis implemenation based on Cache interface
func NewRedisCache(opt RedisOptions) (*RedisDB, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", opt.Address, opt.Port),
		Password: opt.Password,
		DB:       0,
	})

	ctx := context.TODO()
	if err := client.Ping(ctx).Err(); err != nil {
		return &RedisDB{}, err
	}

	return &RedisDB{client: client}, nil
}

// Get is the implementation of Get Redis function
func (r *RedisDB) Get(ctx context.Context, k string) (any, error) {
	return r.client.Get(ctx, k).Result()
}

// Set is the implementation of Set Redis function
func (r *RedisDB) Set(ctx context.Context, k string, v any) (any, error) {
	return r.client.Set(ctx, k, v, 0).Result()
}

// Incr is the implementation of Inch Redis function
func (r *RedisDB) Incr(ctx context.Context, k string) (any, error) {
	return r.client.Incr(ctx, k).Result()
}

// SetNX is the implementation of SetNX Redis function
func (r *RedisDB) SetNX(ctx context.Context, k string, v any) (bool, error) {
	return r.client.SetNX(ctx, k, v, 0).Result()
}
