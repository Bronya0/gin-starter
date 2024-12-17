package gredis

import (
	"testing"
	"time"
)

func Test(t *testing.T) {

	// 获取Redis客户端实例
	rdb := GetRedisClient()

	// 基本操作
	err := rdb.Set("key", "value", 1*time.Hour)
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get("key")
	if err != nil {
		panic(err)
	}
	t.Log(val)

	// 分布式锁
	locked, err := rdb.Lock("mylock", 10*time.Second)
	if err != nil {
		panic(err)
	}
	if locked {
		defer rdb.Unlock("mylock")
		// 执行需要加锁的操作
	}

	// 缓存装饰器使用
	val, err = rdb.GetOrSet("cached_key", func() (interface{}, error) {
		// 这里是获取数据的函数
		return "cached_value", nil
	}, &CacheOption{
		Expiration:  1 * time.Hour,
		ForceUpdate: false,
	})
}
