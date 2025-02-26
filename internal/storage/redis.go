package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/aurelius15/product-reviews/internal/config"
)

type CacheStore interface {
	Set(key string, value any, expiration time.Duration) error
	Get(key string, dest any) error
	Lock(key string, ttl time.Duration) (bool, error)
	Unlock(key string)
	Close()
}

type RedisStorage struct {
	conn *redis.Client
	ctx  context.Context
}

var cacheInstance *RedisStorage

func NewRedisStorage(ctx context.Context, cnf *config.RedisCnf) (CacheStore, error) {
	if cacheInstance != nil {
		return cacheInstance, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr: cnf.Host(),
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	cacheInstance = &RedisStorage{
		conn: client,
		ctx:  ctx,
	}

	return cacheInstance, nil
}

func (r *RedisStorage) Set(key string, value any, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.conn.Set(r.ctx, key, data, expiration).Err()
}

func (r *RedisStorage) Lock(key string, ttl time.Duration) (bool, error) {
	key = fmt.Sprintf("lock:%s", key)

	return r.conn.SetNX(r.ctx, key, "", ttl).Result()
}

func (r *RedisStorage) Unlock(key string) {
	key = fmt.Sprintf("lock:%s", key)
	_ = r.conn.Del(r.ctx, key).Err()
}

func (r *RedisStorage) Get(key string, dest any) error {
	result := r.conn.Get(r.ctx, key)
	if err := result.Err(); err != nil {
		return err
	}

	data, err := result.Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

func (r *RedisStorage) Close() {
	if err := r.conn.Close(); err != nil {
		slog.Error("redis connection is not closed", slog.Any("err", err))
	}
}
