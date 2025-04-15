package cache

import (
	"math"
	"sync"
	"time"
)

// AdaptiveCache 实现一个自适应缓存系统
// 优化点：
// 1. 动态调整TTL基于访问频率
// 2. 热点数据预测和管理
// 3. 自动负载均衡
// 4. 缓存命中率监控与优化
type AdaptiveCache struct {
	baseCache     *HybridCache          // 底层混合缓存实现
	accessStats   map[string]AccessStat // 访问统计
	statsMutex    sync.RWMutex          // 统计数据的读写锁
	hitCount      int                   // 缓存命中次数
	missCount     int                   // 缓存未命中次数
	adaptInterval time.Duration         // 自适应调整的时间间隔
	stopAdapting  chan bool             // 停止自适应调整的信号
}

// AccessStat 存储一个键的访问统计信息
type AccessStat struct {
	AccessCount    int           // 访问次数
	LastAccessTime time.Time     // 最后访问时间
	HitRate        float64       // 命中率
	CurrentTTL     time.Duration // 当前TTL
}

// NewAdaptiveCache 创建一个新的自适应缓存系统
func NewAdaptiveCache(capacityKB int, baseCleanupInterval, adaptInterval time.Duration) *AdaptiveCache {
	cache := &AdaptiveCache{
		baseCache:     NewHybridCache(capacityKB, baseCleanupInterval),
		accessStats:   make(map[string]AccessStat),
		adaptInterval: adaptInterval,
		stopAdapting:  make(chan bool),
	}

	// 启动自适应缓存调整协程
	go cache.startAdaptiveAdjustment()
	return cache
}

// startAdaptiveAdjustment 启动定期自适应调整
func (c *AdaptiveCache) startAdaptiveAdjustment() {
	ticker := time.NewTicker(c.adaptInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.adjustCachingStrategy()
		case <-c.stopAdapting:
			return
		}
	}
}

// adjustCachingStrategy 根据访问模式调整缓存策略
func (c *AdaptiveCache) adjustCachingStrategy() {
	c.statsMutex.Lock()
	defer c.statsMutex.Unlock()

	currentTime := time.Now()
	// 计算热点数据，更新TTL
	for key, stat := range c.accessStats {
		// 检查是否在合理时间内有访问
		if currentTime.Sub(stat.LastAccessTime) < c.adaptInterval*3 {
			// 基于访问频率动态调整TTL
			newTTL := c.calculateOptimalTTL(stat)
			if newTTL != stat.CurrentTTL {
				// 如果有值，尝试重新设置TTL
				if val, exists := c.baseCache.Get(key); exists {
					// 获取当前大小
					size := 0
					c.baseCache.mu.RLock()
					if element, ok := c.baseCache.items[key]; ok {
						item := element.Value.(*CacheItem)
						size = item.Size
					}
					c.baseCache.mu.RUnlock()

					if size > 0 {
						c.baseCache.Set(key, val, size, newTTL)
						stat.CurrentTTL = newTTL
						c.accessStats[key] = stat
					}
				}
			}
		}
	}

	// 清理长时间未访问的统计信息
	for key, stat := range c.accessStats {
		if currentTime.Sub(stat.LastAccessTime) > c.adaptInterval*5 {
			delete(c.accessStats, key)
		}
	}
}

// calculateOptimalTTL 基于访问模式计算最佳TTL
func (c *AdaptiveCache) calculateOptimalTTL(stat AccessStat) time.Duration {
	baseTTL := 5 * time.Minute

	// 访问频率越高，TTL越长
	accessFactor := math.Log1p(float64(stat.AccessCount))
	accessWeight := math.Min(3.0, accessFactor/2.0)

	// 命中率高的数据保留更长时间
	hitRateWeight := stat.HitRate

	// 组合权重
	weightedTTL := baseTTL * time.Duration(1.0+accessWeight+hitRateWeight)

	// 限制最大TTL
	if weightedTTL > 30*time.Minute {
		return 30 * time.Minute
	}

	return weightedTTL
}

// Set 添加或更新缓存项
func (c *AdaptiveCache) Set(key string, value interface{}, size int, ttl time.Duration) {
	c.baseCache.Set(key, value, size, ttl)

	// 更新访问统计
	c.statsMutex.Lock()
	defer c.statsMutex.Unlock()

	stat, exists := c.accessStats[key]
	if exists {
		stat.LastAccessTime = time.Now()
		stat.CurrentTTL = ttl
		c.accessStats[key] = stat
	} else {
		c.accessStats[key] = AccessStat{
			AccessCount:    0,
			LastAccessTime: time.Now(),
			HitRate:        0.0,
			CurrentTTL:     ttl,
		}
	}
}

// Get 获取缓存项，同时更新访问统计
func (c *AdaptiveCache) Get(key string) (interface{}, bool) {
	value, exists := c.baseCache.Get(key)

	// 更新访问统计
	c.statsMutex.Lock()
	defer c.statsMutex.Unlock()

	if exists {
		c.hitCount++
	} else {
		c.missCount++
	}

	totalAccess := c.hitCount + c.missCount
	overallHitRate := float64(0)
	if totalAccess > 0 {
		overallHitRate = float64(c.hitCount) / float64(totalAccess)
	}

	stat, statsExist := c.accessStats[key]
	if statsExist {
		stat.AccessCount++
		stat.LastAccessTime = time.Now()
		if exists {
			// 逐渐调整命中率，加权平均
			stat.HitRate = 0.7*stat.HitRate + 0.3*1.0
		} else {
			stat.HitRate = 0.7*stat.HitRate + 0.3*0.0
		}
		c.accessStats[key] = stat
	} else if exists {
		// 首次命中，创建新的统计记录
		c.baseCache.mu.RLock()
		element, ok := c.baseCache.items[key]
		ttl := 5 * time.Minute // 默认TTL
		if ok {
			item := element.Value.(*CacheItem)
			ttl = item.ExpireTime.Sub(time.Now())
		}
		c.baseCache.mu.RUnlock()

		c.accessStats[key] = AccessStat{
			AccessCount:    1,
			LastAccessTime: time.Now(),
			HitRate:        overallHitRate,
			CurrentTTL:     ttl,
		}
	}

	return value, exists
}

// Remove 从缓存中移除指定项
func (c *AdaptiveCache) Remove(key string) {
	c.baseCache.Remove(key)

	// 可选择是否删除统计数据
	c.statsMutex.Lock()
	delete(c.accessStats, key)
	c.statsMutex.Unlock()
}

// GetHitRate 获取当前缓存命中率
func (c *AdaptiveCache) GetHitRate() float64 {
	c.statsMutex.RLock()
	defer c.statsMutex.RUnlock()

	total := c.hitCount + c.missCount
	if total == 0 {
		return 0.0
	}
	return float64(c.hitCount) / float64(total)
}

// Close 关闭缓存系统
func (c *AdaptiveCache) Close() {
	close(c.stopAdapting)
	c.baseCache.Close()
}

// Size 返回当前缓存中的项目数
func (c *AdaptiveCache) Size() int {
	return c.baseCache.Size()
}

// CurrentMemoryUsage 返回当前缓存使用的内存量（KB）
func (c *AdaptiveCache) CurrentMemoryUsage() int {
	return c.baseCache.CurrentMemoryUsage()
}
