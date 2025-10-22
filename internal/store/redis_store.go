package store

import (
    "context"
    "time"

    "github.com/go-redis/redis/v9"
)

type RedisStore struct {
    cli     *redis.Client
    backend Store
    ttl     time.Duration
}

func NewRedisStore(cli *redis.Client, backend Store, ttl time.Duration) *RedisStore {
    return &RedisStore{cli: cli, backend: backend, ttl: ttl}
}

func (r *RedisStore) Put(short string, long string) error {
    // write-through: backend durable first
    if err := r.backend.Put(short, long); err != nil {
        return err
    }
    ctx := context.Background()
    _ = r.cli.Set(ctx, short, long, r.ttl).Err()
    return nil
}

func (r *RedisStore) Get(short string) (string, error) {
    ctx := context.Background()
    v, err := r.cli.Get(ctx, short).Result()
    if err == nil {
        return v, nil
    }
    // fallback
    long, err2 := r.backend.Get(short)
    if err2 == nil {
        _ = r.cli.Set(ctx, short, long, r.ttl).Err()
        return long, nil
    }
    return "", ErrNotFound
}
