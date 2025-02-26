// pkg/redis/utils.go

package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisUtil struct {
    client *redis.Client
}

func NewRedisUtil(client *redis.Client) *RedisUtil {
    return &RedisUtil{client: client}
}

// SetJSON 存储JSON数据
func (r *RedisUtil) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
    json, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return r.client.Set(ctx, key, json, expiration).Err()
}

// GetJSON 获取JSON数据
func (r *RedisUtil) GetJSON(ctx context.Context, key string, dest interface{}) error {
    val, err := r.client.Get(ctx, key).Result()
    if err != nil {
        return err
    }
    return json.Unmarshal([]byte(val), dest)
}

// Lock 获取分布式锁
func (r *RedisUtil) Lock(ctx context.Context, key string, expiration time.Duration) bool {
    return r.client.SetNX(ctx, key, 1, expiration).Val()
}

// Unlock 释放分布式锁
func (r *RedisUtil) Unlock(ctx context.Context, key string) error {
    return r.client.Del(ctx, key).Err()
}

// CacheAside 缓存旁路模式
func (r *RedisUtil) CacheAside(ctx context.Context, key string, dest interface{}, 
    queryFunc func() (interface{}, error), expiration time.Duration) error {
    // 1. 查询缓存
    err := r.GetJSON(ctx, key, dest)
    if err == nil {
        return nil
    }

    // 2. 缓存未命中，查询数据库
    data, err := queryFunc()
    if err != nil {
        return err
    }

    // 3. 写入缓存
    if err := r.SetJSON(ctx, key, data, expiration); err != nil {
        return err
    }

    // 4. 将结果赋值给目标变量
    bytes, _ := json.Marshal(data)
    return json.Unmarshal(bytes, dest)
}

// Pipeline 批量操作
func (r *RedisUtil) Pipeline(ctx context.Context, fn func(pipe redis.Pipeliner) error) error {
    pipe := r.client.Pipeline()
    if err := fn(pipe); err != nil {
        return err
    }
    _, err := pipe.Exec(ctx)
    return err
}