package map_util

import "sync"

// ConcurrentMap 是一个泛型并发安全的映射
type ConcurrentMap[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

// NewConcurrentMap 创建一个新的并发安全映射
func NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{m: make(map[K]V)}
}

// Store 添加或更新一个键值对
func (cm *ConcurrentMap[K, V]) Store(key K, value V) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.m[key] = value
}

// Load 获取一个键对应的值，第二个返回值表示是否存在该键
func (cm *ConcurrentMap[K, V]) Load(key K) (V, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	value, ok := cm.m[key]
	return value, ok
}

// Delete 删除一个键值对
func (cm *ConcurrentMap[K, V]) Delete(key K) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.m, key)
}

// Range 遍历映射中的所有键值对
func (cm *ConcurrentMap[K, V]) Range(f func(key K, value V) bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	for k, v := range cm.m {
		if !f(k, v) {
			break
		}
	}
}

// Clear 清空 ConcurrentMap
func (cm *ConcurrentMap[K, V]) Clear() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.m = make(map[K]V)
}

//	// 示例：使用字符串键和整数值
//	myMap := NewConcurrentMap[string, int]()
//
//	// 写操作
//	myMap.Store("key1", 100)
//	myMap.Store("key2", 200)
//
//	// 读操作
//	value, ok := myMap.Load("key1")
//	fmt.Printf("key1: %d, exists? %v\n", value, ok)
//
//	// 遍历
//	myMap.Range(func(key string, value int) bool {
//		fmt.Printf("%s: %d\n", key, value)
//		return true
//	})
