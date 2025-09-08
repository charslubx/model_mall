package svc

import (
	"context"
	"fmt"
	"model_mall_backend/backend/internal/config"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// RedisHelper 封装了 go-zero 的 redis 客户端
type RedisHelper struct {
	client *redis.Redis
	logger *LogHelper // 保留自定义 logger
}

// NewRedisHelper 创建一个新的 RedisHelper 实例
// 适配 config.Redis (Password) 到 go-zero redis.RedisConf (Pass)
func NewRedisHelper(ctx context.Context, conf *config.Config, logger *LogHelper) (*RedisHelper, error) {
	// 创建一个符合 go-zero redis.RedisConf 结构的实例
	// 假设 conf.Redis.DB 是整数形式的数据库索引
	redisConf := redis.RedisConf{
		Host: fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port), // Host 包含地址和端口
		Pass: conf.Redis.Password,                                    // 适配 Password 字段
		Type: "node",                                                 // 单节点模式，根据实际情况调整
		Tls:  false,                                                  // 根据实际情况调整
	}

	// 使用 go-zero 的 redis.MustNewRedis 创建客户端
	// 注意：MustNewRedis 在连接失败时会 panic，如果需要更优雅的错误处理，
	// 可以使用 redis.NewRedis 并检查返回的 error。
	// 但通常在服务初始化阶段，连接失败是致命错误。
	rds := redis.MustNewRedis(redisConf, redis.WithPass(redisConf.Pass))

	// 测试连接 - 使用 go-zero 的上下文感知方法
	status := rds.PingCtx(ctx)
	if !status {
		logger.Debug("Failed to connect to redis")
		// 注意：如果 MustNewRedis 失败会 panic，这里主要是检查 Ping
		return nil, fmt.Errorf("failed to connect to redis")
	}

	return &RedisHelper{
		client: rds,
		logger: logger,
	}, nil
}

// Close 关闭Redis连接
func (r *RedisHelper) Close() {
	// go-zero 的 redis.Redis 通常由内部连接池管理。
	// 如果 MustNewRedis 或 NewRedis 返回的实例有 Close 方法，可以调用。
	// 根据 go-zero 源码，*redis.Redis 本身没有 Close() 方法，但其内部 client 可能有。
	// 为了安全起见，可以尝试调用，或者留空。这里假设没有直接的 Close 方法需要调用。
	// 如果未来 go-zero 版本变更需要显式关闭，再添加。
	// r.logger.Warn("RedisHelper.Close() called, standard procedure for go-zero redis client.")
	// 如果确实需要关闭底层资源（不常见），可能需要类型断言访问内部 client。
}

// GetClient 获取Redis客户端实例
func (r *RedisHelper) GetClient() *redis.Redis {
	return r.client
}

// Set 设置键值对
func (r *RedisHelper) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	start := time.Now()

	expirySeconds := int(expiration / time.Second)
	if expiration > 0 && expirySeconds == 0 {
		expirySeconds = 1
	}

	var err error
	if strVal, ok := value.(string); ok {
		// 使用 go-zero 的 SetexCtx 方法设置带过期时间的键值
		err = r.client.SetexCtx(ctx, key, strVal, expirySeconds)
	} else {
		// 简单处理非字符串类型，转换为字符串
		err = r.client.SetexCtx(ctx, key, fmt.Sprintf("%v", value), expirySeconds)
		// 如果需要序列化（如 JSON），请在此处添加逻辑
		// import "encoding/json"
		// var data []byte
		// data, err = json.Marshal(value)
		// if err != nil { ... }
		// err = r.client.SetexCtx(ctx, key, string(data), expirySeconds)
	}

	duration := time.Since(start)

	if duration > time.Second {
		logx.WithContext(ctx).WithDuration(duration).Slowf("Slow Redis SET operation: key=%s", key)
	}

	if err != nil {
		logx.WithContext(ctx).Errorf("Failed to SET key: %s, error: %v", key, err)
		return err
	}

	logx.WithContext(ctx).Infof("Redis SET successful: key=%s", key)
	return nil
}

// Get 获取键值
func (r *RedisHelper) Get(ctx context.Context, key string) (string, error) {
	start := time.Now()
	val, err := r.client.GetCtx(ctx, key)
	duration := time.Since(start)

	if duration > time.Second {
		logx.WithContext(ctx).WithDuration(duration).Slowf("Slow Redis GET operation: key=%s", key)
	}

	if err == redis.Nil {
		logx.WithContext(ctx).Infof("Redis key not found: %s", key)
		return "", nil
	} else if err != nil {
		logx.WithContext(ctx).Errorf("Failed to GET key: %s, error: %v", key, err)
		return "", err
	}

	logx.WithContext(ctx).Infof("Redis GET successful: key=%s", key)
	return val, nil
}

// Del 删除键
func (r *RedisHelper) Del(ctx context.Context, key string) error {
	start := time.Now()
	// DelCtx 返回删除的键数量和 error
	_, err := r.client.DelCtx(ctx, key)
	duration := time.Since(start)

	if duration > time.Second {
		logx.WithContext(ctx).WithDuration(duration).Slowf("Slow Redis DEL operation: key=%s", key)
	}

	if err != nil {
		logx.WithContext(ctx).Errorf("Failed to DEL key: %s, error: %v", key, err)
		return err
	}

	logx.WithContext(ctx).Infof("Redis DEL successful: key=%s", key)
	return nil
}

// Exists 检查键是否存在
func (r *RedisHelper) Exists(ctx context.Context, key string) (bool, error) {
	start := time.Now()
	// ExistsCtx 返回存在的键数量 (0 or 1 for single key) 和 error
	exists, err := r.client.ExistsCtx(ctx, key)
	duration := time.Since(start)

	if duration > time.Second {
		logx.WithContext(ctx).WithDuration(duration).Slowf("Slow Redis EXISTS operation: key=%s", key)
	}

	if err != nil {
		logx.WithContext(ctx).Errorf("Failed to check EXISTS for key: %s, error: %v", key, err)
		return false, err
	}

	logx.WithContext(ctx).Infof("Redis EXISTS check successful: key=%s, exists=%v", key, exists)
	return exists, nil
}

// SetNX 当键不存在时设置值
func (r *RedisHelper) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	start := time.Now()

	expirySeconds := int(expiration / time.Second)
	if expiration > 0 && expirySeconds == 0 {
		expirySeconds = 1
	}

	var setResult bool
	var err error
	if strVal, ok := value.(string); ok {
		// SetnxCtx 只设置键值，不直接处理过期时间
		setResult, err = r.client.SetnxCtx(ctx, key, strVal)
	} else {
		// 简单处理非字符串类型
		setResult, err = r.client.SetnxCtx(ctx, key, fmt.Sprintf("%v", value))
		// 序列化处理同 Set 方法
	}

	// 如果 SetNX 成功，并且设置了过期时间，则设置 TTL
	if err == nil && setResult && expirySeconds > 0 {
		expireErr := r.client.ExpireCtx(ctx, key, expirySeconds)
		if expireErr != nil {
			logx.WithContext(ctx).Errorf("Failed to set expiration after SETNX for key: %s, error: %v", key, expireErr)
			// 将过期设置失败作为主要错误返回，或根据业务决定
			err = expireErr
		}
	}

	duration := time.Since(start)

	if duration > time.Second {
		logx.WithContext(ctx).WithDuration(duration).Slowf("Slow Redis SETNX operation: key=%s", key)
	}

	if err != nil {
		logx.WithContext(ctx).Errorf("Failed to SETNX key: %s, error: %v", key, err)
		return false, err
	}

	logx.WithContext(ctx).Infof("Redis SETNX successful: key=%s, set=%v", key, setResult)
	return setResult, nil
}
