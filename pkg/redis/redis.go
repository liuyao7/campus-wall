// pkg/redis/redis.go

package redis

import (
	"context"
	"strconv"
	"time"

	"campus-wall/pkg/config"

	"github.com/go-redis/redis/v8"
)

// type Config struct {
//     Host         string        `mapstructure:"host"`
//     Port         string        `mapstructure:"port"`
//     Password     string        `mapstructure:"password"`
//     DB           int          `mapstructure:"db"`
//     PoolSize     int          `mapstructure:"pool_size"`
//     MinIdleConns int          `mapstructure:"min_idle_conns"`
//     MaxRetries   int          `mapstructure:"max_retries"`
//     DialTimeout  time.Duration `mapstructure:"dial_timeout"`
//     ReadTimeout  time.Duration `mapstructure:"read_timeout"`
//     WriteTimeout time.Duration `mapstructure:"write_timeout"`
// }

func NewRedisClient(cfg *config.RedisConfig) (*redis.Client, error) {
    client := redis.NewClient(&redis.Options{
        Addr:         cfg.Host + ":" + strconv.Itoa(cfg.Port),
        Password:     cfg.Password,
        DB:           cfg.DB,
        PoolSize:     cfg.PoolSize,     // 连接池最大连接数
        MinIdleConns: cfg.MinIdleConns, // 最小空闲连接数
        MaxRetries:   cfg.MaxRetries,   // 最大重试次数
        DialTimeout:  cfg.DialTimeout,  // 连接超时
        ReadTimeout:  cfg.ReadTimeout,  // 读取超时
        WriteTimeout: cfg.WriteTimeout, // 写入超时
    })

    // 测试连接
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := client.Ping(ctx).Err(); err != nil {
        return nil, err
    }

    return client, nil
}