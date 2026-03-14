package v8

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v8"
)

// Config Redis 客户端配置。
type Config struct {
	Addr         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DefaultConfig 返回默认 Redis 配置。
func DefaultConfig() Config {
	return Config{
		PoolSize:     10,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
}

// NewClient 创建 Redis v8 客户端并执行连通性校验。
func NewClient(ctx context.Context, cfg Config) (*redis.Client, error) {
	if cfg.Addr == "" {
		return nil, fmt.Errorf("redis address is empty")
	}
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})
	pingCtx := ctx
	if pingCtx == nil {
		tmp, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		pingCtx = tmp
	}
	if err := client.Ping(pingCtx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect redis: %w", err)
	}
	return client, nil
}

// SetString 写字符串缓存。
func SetString(ctx context.Context, cli *redis.Client, key string, val string, ttl time.Duration) error {
	if cli == nil {
		return nil
	}
	return cli.Set(ctx, key, val, ttl).Err()
}

// GetString 读取字符串缓存。
func GetString(ctx context.Context, cli *redis.Client, key string) (string, bool, error) {
	if cli == nil {
		return "", false, nil
	}
	raw, err := cli.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return raw, true, nil
}

// Del 删除缓存键。
func Del(ctx context.Context, cli *redis.Client, keys ...string) error {
	if cli == nil || len(keys) == 0 {
		return nil
	}
	return cli.Del(ctx, keys...).Err()
}

// IncrWithExpire 自增计数并在首次创建时设置过期时间。
func IncrWithExpire(ctx context.Context, cli *redis.Client, key string, ttl time.Duration) (int64, error) {
	if cli == nil {
		return 0, nil
	}
	n, err := cli.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if n == 1 && ttl > 0 {
		_ = cli.Expire(ctx, key, ttl).Err()
	}
	return n, nil
}

// SetJSON 序列化并写入 JSON 缓存。
func SetJSON(ctx context.Context, cli *redis.Client, key string, val any, ttl time.Duration) error {
	if cli == nil {
		return nil
	}
	raw, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return cli.Set(ctx, key, raw, ttl).Err()
}

// GetJSON 读取并反序列化 JSON 缓存。
func GetJSON(ctx context.Context, cli *redis.Client, key string, out any) (bool, error) {
	if cli == nil {
		return false, nil
	}
	raw, err := cli.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if len(raw) == 0 {
		return false, nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return false, err
	}
	return true, nil
}

// RedisFromContext 从上下文按 key 提取 Redis 客户端。
func RedisFromContext(ctx context.Context, key any) *redis.Client {
	if ctx == nil {
		return nil
	}
	if cli, ok := ctx.Value(key).(*redis.Client); ok {
		return cli
	}
	return nil
}
