package gredis

import (
	"context"
	"fmt"
	"gin-starter/internal/config"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

//单例模式：保证全局只有一个 Redis 连接池
//配置灵活：通过配置文件管理 Redis 连接参数
//功能完整：封装了常用的 Redis 操作
//上下文支持：支持自定义上下文处理超时等场景
//分布式锁：实现了基本的分布式锁功能
//缓存装饰器：提供了便捷的缓存模式
//类型安全：使用泛型和强类型接口
//错误处理：统一的错误处理方式

type RedisClient struct {
	Client *redis.Client
	ctx    context.Context
}

var (
	redisClient *RedisClient
	once        sync.Once
)

// GetRedisClient 获取Redis客户端单例
func GetRedisClient() *RedisClient {
	once.Do(func() {
		redisClient = newRedisClient()
	})
	return redisClient
}

// newRedisClient 创建新的Redis客户端
func newRedisClient() *RedisClient {
	conf := config.GloConfig.Redis

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password:     conf.Password,
		DB:           conf.DB,
		PoolSize:     conf.PoolSize,
		MinIdleConns: conf.MinIdleConns,
		IdleTimeout:  time.Duration(conf.IdleTimeout) * time.Second,
	})

	ctx := context.Background()

	// 测试连接
	if err := client.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Redis connection failed: %v", err))
	}

	return &RedisClient{
		Client: client,
		ctx:    ctx,
	}
}

// Get 获取键值
func (r *RedisClient) Get(key string) (string, error) {
	return r.Client.Get(r.ctx, key).Result()
}

// Set 设置键值
func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(r.ctx, key, value, expiration).Err()
}

// Delete 删除键
func (r *RedisClient) Delete(keys ...string) error {
	return r.Client.Del(r.ctx, keys...).Err()
}

// Exists 检查键是否存在
func (r *RedisClient) Exists(keys ...string) (bool, error) {
	n, err := r.Client.Exists(r.ctx, keys...).Result()
	return n > 0, err
}

// SetNX 仅当键不存在时设置
func (r *RedisClient) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.Client.SetNX(r.ctx, key, value, expiration).Result()
}

// HSet Hash设置
func (r *RedisClient) HSet(key string, values ...interface{}) error {
	return r.Client.HSet(r.ctx, key, values...).Err()
}

// HGet Hash获取
func (r *RedisClient) HGet(key, field string) (string, error) {
	return r.Client.HGet(r.ctx, key, field).Result()
}

// HGetAll 获取整个Hash
func (r *RedisClient) HGetAll(key string) (map[string]string, error) {
	return r.Client.HGetAll(r.ctx, key).Result()
}

// Incr 自增
func (r *RedisClient) Incr(key string) (int64, error) {
	return r.Client.Incr(r.ctx, key).Result()
}

// ZAdd 添加有序集合成员
func (r *RedisClient) ZAdd(key string, members ...*redis.Z) error {
	return r.Client.ZAdd(r.ctx, key, members...).Err()
}

// ZRange 获取有序集合范围
func (r *RedisClient) ZRange(key string, start, stop int64) ([]string, error) {
	return r.Client.ZRange(r.ctx, key, start, stop).Result()
}

// Pipeline 管道操作
func (r *RedisClient) Pipeline() redis.Pipeliner {
	return r.Client.Pipeline()
}

// WithContext 使用自定义上下文
func (r *RedisClient) WithContext(ctx context.Context) *RedisClient {
	return &RedisClient{
		Client: r.Client,
		ctx:    ctx,
	}
}

// Close 关闭连接
func (r *RedisClient) Close() error {
	return r.Client.Close()
}

// Lock 分布式锁实现
func (r *RedisClient) Lock(key string, expiration time.Duration) (bool, error) {
	return r.Client.SetNX(r.ctx, "lock:"+key, 1, expiration).Result()
}

// Unlock 释放分布式锁
func (r *RedisClient) Unlock(key string) error {
	return r.Client.Del(r.ctx, "lock:"+key).Err()
}

// TryLock 尝试获取锁带重试
func (r *RedisClient) TryLock(key string, expiration time.Duration, retryTimes int, retryDelay time.Duration) (bool, error) {
	for i := 0; i < retryTimes; i++ {
		locked, err := r.Lock(key, expiration)
		if err != nil {
			return false, err
		}
		if locked {
			return true, nil
		}
		time.Sleep(retryDelay)
	}
	return false, nil
}

// Cache 缓存装饰器
type CacheOption struct {
	Expiration  time.Duration
	ForceUpdate bool
}

// GetOrSet 获取缓存，不存在则设置
func (r *RedisClient) GetOrSet(key string, fn func() (interface{}, error), opt *CacheOption) (string, error) {
	if !opt.ForceUpdate {
		// 尝试从缓存获取
		if val, err := r.Get(key); err == nil {
			return val, nil
		}
	}

	// 获取新数据
	val, err := fn()
	if err != nil {
		return "", err
	}

	// 设置缓存
	if err := r.Set(key, val, opt.Expiration); err != nil {
		return "", err
	}

	return fmt.Sprint(val), nil
}
