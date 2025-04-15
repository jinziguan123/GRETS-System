package cache

import (
	"container/list"
	"sync"
	"time"
)

// CacheItem 表示缓存中的一个项目
type CacheItem struct {
	Key        string
	Value      interface{}
	Size       int       // 项目大小（千字节）
	ExpireTime time.Time // 过期时间
	Created    time.Time // 创建时间
}

// HybridCache 实现一个混合缓存系统，同时使用LRU和TTL策略
// 优化点：
// 1. LRU策略管理内存容量限制
// 2. TTL策略管理数据新鲜度
// 3. 异步清理过期项目以提高性能
// 4. 并发安全设计
type HybridCache struct {
	capacity     int                      // 总容量（千字节）
	items        map[string]*list.Element // 存储缓存项的哈希表
	evictionList *list.List               // LRU淘汰列表
	mu           sync.RWMutex             // 并发控制读写锁
	currentSize  int                      // 当前使用的大小
	cleanupTimer *time.Ticker             // 清理定时器
	stopCleanup  chan bool                // 停止清理的信号
}

// NewHybridCache 创建一个新的混合缓存系统
func NewHybridCache(capacityKB int, cleanupInterval time.Duration) *HybridCache {
	cache := &HybridCache{
		capacity:     capacityKB,
		items:        make(map[string]*list.Element),
		evictionList: list.New(),
		currentSize:  0,
		stopCleanup:  make(chan bool),
	}

	// 启动异步过期项目清理
	cache.startCleanupTimer(cleanupInterval)
	return cache
}

// startCleanupTimer 启动定期清理过期缓存项的定时器
func (c *HybridCache) startCleanupTimer(cleanupInterval time.Duration) {
	c.cleanupTimer = time.NewTicker(cleanupInterval)
	go func() {
		for {
			select {
			case <-c.cleanupTimer.C:
				c.Cleanup()
			case <-c.stopCleanup:
				c.cleanupTimer.Stop()
				return
			}
		}
	}()
}

// Set 添加或更新缓存项
func (c *HybridCache) Set(key string, value interface{}, sizeKB int, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	expireTime := now.Add(ttl)

	// 如果键已存在，则更新并移至队列前端
	if element, exists := c.items[key]; exists {
		c.evictionList.MoveToFront(element)
		item := element.Value.(*CacheItem)
		// 更新当前容量
		c.currentSize -= item.Size
		c.currentSize += sizeKB
		// 更新项
		item.Value = value
		item.Size = sizeKB
		item.ExpireTime = expireTime
	} else {
		// 添加新项
		item := &CacheItem{
			Key:        key,
			Value:      value,
			Size:       sizeKB,
			ExpireTime: expireTime,
			Created:    now,
		}
		element := c.evictionList.PushFront(item)
		c.items[key] = element
		c.currentSize += sizeKB
	}

	// 如果超出容量，删除最不常用的项直到容量足够
	for c.currentSize > c.capacity && c.evictionList.Len() > 0 {
		c.evictOldest()
	}
}

// Get 获取缓存项，如果存在且未过期则返回
func (c *HybridCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	element, exists := c.items[key]
	if !exists {
		c.mu.RUnlock()
		return nil, false
	}

	item := element.Value.(*CacheItem)
	// 检查是否过期
	if time.Now().After(item.ExpireTime) {
		c.mu.RUnlock()
		// 异步删除过期项
		go c.Remove(key)
		return nil, false
	}

	c.mu.RUnlock()

	// 更新项的位置（需要写锁）
	c.mu.Lock()
	defer c.mu.Unlock()

	// 再次检查项目是否存在（可能在获取写锁期间被其他协程删除）
	if _, stillExists := c.items[key]; stillExists {
		c.evictionList.MoveToFront(element)
		return item.Value, true
	}

	return nil, false
}

// Remove 从缓存中移除指定项
func (c *HybridCache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, exists := c.items[key]; exists {
		item := element.Value.(*CacheItem)
		c.currentSize -= item.Size
		delete(c.items, key)
		c.evictionList.Remove(element)
	}
}

// Cleanup 清理所有过期项
func (c *HybridCache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, element := range c.items {
		item := element.Value.(*CacheItem)
		if now.After(item.ExpireTime) {
			c.currentSize -= item.Size
			delete(c.items, key)
			c.evictionList.Remove(element)
		}
	}
}

// evictOldest 淘汰最旧（最不常用）的缓存项
func (c *HybridCache) evictOldest() {
	if element := c.evictionList.Back(); element != nil {
		item := element.Value.(*CacheItem)
		c.currentSize -= item.Size
		delete(c.items, item.Key)
		c.evictionList.Remove(element)
	}
}

// Close 关闭缓存系统，停止后台任务
func (c *HybridCache) Close() {
	close(c.stopCleanup)
}

// Size 返回当前缓存中的项目数
func (c *HybridCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// CurrentMemoryUsage 返回当前缓存使用的内存量（KB）
func (c *HybridCache) CurrentMemoryUsage() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentSize
}
