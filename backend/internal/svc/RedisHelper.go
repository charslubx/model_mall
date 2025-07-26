package svc

import (
	"context"
	"fmt"
	"model_mall_backend/backend/internal/config"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisHelper struct {
	client *redis.Client
	logger *LogHelper
	ctx    context.Context
}

func NewRedisHelper(ctx context.Context, conf *config.Config, logger *LogHelper) (*RedisHelper, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port),
		Password:     conf.Redis.Password,
		DB:           conf.Redis.DB,
		PoolSize:     10,
		MinIdleConns: 2,
		MaxRetries:   3,
	})

	// 测试连接
	if err := client.Ping(ctx).Err(); err != nil {
		logger.Error(err, "Failed to connect to redis")
		return nil, fmt.Errorf("failed to connect to redis: %v", err)
	}

	return &RedisHelper{
		client: client,
		logger: logger,
		ctx:    ctx,
	}, nil
}

// Close 关闭Redis连接
func (r *RedisHelper) Close() {
	if r.client != nil {
		r.client.Close()
	}
}

// GetClient 获取Redis客户端实例
func (r *RedisHelper) GetClient() *redis.Client {
	return r.client
}

// Set 设置键值对
func (r *RedisHelper) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	start := time.Now()
	err := r.client.Set(ctx, key, value, expiration).Err()
	duration := time.Since(start)

	if duration > time.Second {
		r.logger.Slow("Slow Redis SET operation (took %v): key=%s", duration, key)
	}

	if err != nil {
		r.logger.Error(err, "Failed to SET key: %s", key)
		return err
	}

	r.logger.Debug("Redis SET successful: key=%s", key)
	return nil
}

// Get 获取键值
func (r *RedisHelper) Get(ctx context.Context, key string) (string, error) {
	start := time.Now()
	val, err := r.client.Get(ctx, key).Result()
	duration := time.Since(start)

	if duration > time.Second {
		r.logger.Slow("Slow Redis GET operation (took %v): key=%s", duration, key)
	}

	if err == redis.Nil {
		r.logger.Debug("Redis key not found: %s", key)
		return "", nil
	} else if err != nil {
		r.logger.Error(err, "Failed to GET key: %s", key)
		return "", err
	}

	r.logger.Debug("Redis GET successful: key=%s", key)
	return val, nil
}

// Del 删除键
func (r *RedisHelper) Del(ctx context.Context, key string) error {
	start := time.Now()
	err := r.client.Del(ctx, key).Err()
	duration := time.Since(start)

	if duration > time.Second {
		r.logger.Slow("Slow Redis DEL operation (took %v): key=%s", duration, key)
	}

	if err != nil {
		r.logger.Error(err, "Failed to DEL key: %s", key)
		return err
	}

	r.logger.Debug("Redis DEL successful: key=%s", key)
	return nil
}

// Exists 检查键是否存在
func (r *RedisHelper) Exists(ctx context.Context, key string) (bool, error) {
	start := time.Now()
	result, err := r.client.Exists(ctx, key).Result()
	duration := time.Since(start)

	if duration > time.Second {
		r.logger.Slow("Slow Redis EXISTS operation (took %v): key=%s", duration, key)
	}

	if err != nil {
		r.logger.Error(err, "Failed to check EXISTS for key: %s", key)
		return false, err
	}

	exists := result > 0
	r.logger.Debug("Redis EXISTS check successful: key=%s, exists=%v", key, exists)
	return exists, nil
}

// SetNX 当键不存在时设置值
func (r *RedisHelper) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	start := time.Now()
	result, err := r.client.SetNX(ctx, key, value, expiration).Result()
	duration := time.Since(start)

	if duration > time.Second {
		r.logger.Slow("Slow Redis SETNX operation (took %v): key=%s", duration, key)
	}

	if err != nil {
		r.logger.Error(err, "Failed to SETNX key: %s", key)
		return false, err
	}

	r.logger.Debug("Redis SETNX successful: key=%s, set=%v", key, result)
	return result, nil
}
