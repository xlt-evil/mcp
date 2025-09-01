package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"hello-mcp-server/config"

	"github.com/redis/go-redis/v9"
)

// RedisManager Redis管理器
type RedisManager struct {
	config *config.RedisConfig
	client *redis.Client
}

// RedisResult Redis操作结果结构
type RedisResult struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// KeyInfo 键信息结构
type KeyInfo struct {
	Key      string `json:"key"`
	Type     string `json:"type"`
	TTL      int64  `json:"ttl"`
	Size     int64  `json:"size"`
	Encoding string `json:"encoding"`
}

// NewRedisManager 创建Redis管理器
func NewRedisManager(cfg *config.RedisConfig) *RedisManager {
	return &RedisManager{
		config: cfg,
	}
}

// Connect 连接Redis
func (rm *RedisManager) Connect() error {
	if !rm.config.IsValid() {
		return fmt.Errorf("invalid redis configuration")
	}

	log.Printf("Connecting to Redis: %s", rm.config.GetAddr())

	// 创建Redis客户端
	rm.client = redis.NewClient(&redis.Options{
		Addr:            rm.config.GetAddr(),
		Password:        rm.config.GetPassword(),
		DB:              rm.config.GetDB(),
		MaxRetries:      3,
		MinRetryBackoff: time.Millisecond * 100,
		MaxRetryBackoff: time.Second * 2,
		DialTimeout:     rm.config.Timeout.Connect,
		ReadTimeout:     rm.config.Timeout.Read,
		WriteTimeout:    rm.config.Timeout.Write,
		PoolSize:        rm.config.Pool.MaxActive,
		MinIdleConns:    rm.config.Pool.MaxIdle,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rm.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping redis: %v", err)
	}

	log.Printf("Successfully connected to Redis: %s", rm.config.GetAddr())
	return nil
}

// Close 关闭Redis连接
func (rm *RedisManager) Close() error {
	if rm.client != nil {
		return rm.client.Close()
	}
	return nil
}

// IsConnected 检查是否已连接
func (rm *RedisManager) IsConnected() bool {
	if rm.client == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return rm.client.Ping(ctx).Err() == nil
}

// Get 获取键值
func (rm *RedisManager) Get(key string) *RedisResult {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	val, err := rm.client.Get(ctx, key).Result()
	if err != nil {
		return &RedisResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to get key %s: %v", key, err),
		}
	}

	return &RedisResult{
		Success: true,
		Data:    val,
	}
}

// Set 设置键值
func (rm *RedisManager) Set(key, value string, expiration time.Duration) *RedisResult {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := rm.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return &RedisResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to set key %s: %v", key, err),
		}
	}

	return &RedisResult{
		Success: true,
		Data:    "OK",
	}
}

// Del 删除键
func (rm *RedisManager) Del(keys ...string) *RedisResult {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := rm.client.Del(ctx, keys...).Result()
	if err != nil {
		return &RedisResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to delete keys: %v", err),
		}
	}

	return &RedisResult{
		Success: true,
		Data:    result,
	}
}

// Keys 获取匹配的键
func (rm *RedisManager) Keys(pattern string) *RedisResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	keys, err := rm.client.Keys(ctx, pattern).Result()
	if err != nil {
		return &RedisResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to get keys with pattern %s: %v", pattern, err),
		}
	}

	return &RedisResult{
		Success: true,
		Data:    keys,
	}
}

// Type 获取键类型
func (rm *RedisManager) Type(key string) *RedisResult {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	keyType, err := rm.client.Type(ctx, key).Result()
	if err != nil {
		return &RedisResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to get type of key %s: %v", key, err),
		}
	}

	return &RedisResult{
		Success: true,
		Data:    keyType,
	}
}

// TTL 获取键的TTL
func (rm *RedisManager) TTL(key string) *RedisResult {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ttl, err := rm.client.TTL(ctx, key).Result()
	if err != nil {
		return &RedisResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to get TTL of key %s: %v", key, err),
		}
	}

	return &RedisResult{
		Success: true,
		Data:    ttl.Seconds(),
	}
}

// Info 获取Redis信息
func (rm *RedisManager) Info(section string) *RedisResult {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info, err := rm.client.Info(ctx, section).Result()
	if err != nil {
		return &RedisResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to get Redis info: %v", err),
		}
	}

	return &RedisResult{
		Success: true,
		Data:    info,
	}
}

// DBSize 获取数据库大小
func (rm *RedisManager) DBSize() *RedisResult {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	size, err := rm.client.DBSize(ctx).Result()
	if err != nil {
		return &RedisResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to get database size: %v", err),
		}
	}

	return &RedisResult{
		Success: true,
		Data:    size,
	}
}

// FlushDB 清空当前数据库
func (rm *RedisManager) FlushDB() *RedisResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := rm.client.FlushDB(ctx).Err()
	if err != nil {
		return &RedisResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to flush database: %v", err),
		}
	}

	return &RedisResult{
		Success: true,
		Data:    "OK",
	}
}

// ExecuteCommand 执行自定义命令
func (rm *RedisManager) ExecuteCommand(command string, args ...interface{}) *RedisResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 构建完整的参数列表
	allArgs := []interface{}{command}
	allArgs = append(allArgs, args...)
	result, err := rm.client.Do(ctx, allArgs...).Result()
	if err != nil {
		return &RedisResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to execute command %s: %v", command, err),
		}
	}

	return &RedisResult{
		Success: true,
		Data:    result,
	}
}
