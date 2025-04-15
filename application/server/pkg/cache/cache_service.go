package cache

import (
	"encoding/json"
	"fmt"
	"grets_server/pkg/utils"
	"sync"
	"time"
)

// 缓存键前缀常量
const (
	RealtyPrefix      = "realty:"
	UserPrefix        = "user:"
	TransactionPrefix = "transaction:"
	ContractPrefix    = "contract:"
)

// CacheService 缓存服务接口
type CacheService interface {
	// Set 存储对象到缓存
	Set(key string, value interface{}, size int, ttl time.Duration)
	// Get 从缓存获取对象
	Get(key string, result interface{}) bool
	// Remove 从缓存中删除对象
	Remove(key string)
	// GetHitRate 获取命中率
	GetHitRate() float64
	// Close 关闭服务
	Close()
}

// cacheServiceImpl 缓存服务实现
type cacheServiceImpl struct {
	cache          *AdaptiveCache // 自适应缓存
	bytesConverter func(interface{}) ([]byte, error)
}

var (
	globalCacheService CacheService
	cacheServiceOnce   sync.Once
)

// GetCacheService 获取全局缓存服务实例
func GetCacheService() CacheService {
	cacheServiceOnce.Do(func() {
		capacityKB := 10 * 1024 // 10MB 默认容量
		cleanupInterval := 5 * time.Minute
		adaptInterval := 30 * time.Second

		cache := NewAdaptiveCache(capacityKB, cleanupInterval, adaptInterval)
		globalCacheService = &cacheServiceImpl{
			cache:          cache,
			bytesConverter: json.Marshal,
		}

		utils.Log.Info("缓存服务初始化完成")
	})
	return globalCacheService
}

// Set 将对象存储到缓存
func (s *cacheServiceImpl) Set(key string, value interface{}, size int, ttl time.Duration) {
	// 如果size为0，尝试估算对象大小
	if size == 0 {
		bytes, err := s.bytesConverter(value)
		if err != nil {
			utils.Log.Error(fmt.Sprintf("序列化对象失败: %v", err))
			return
		}
		size = len(bytes) / 1024 // 转换为KB
		if size < 1 {
			size = 1 // 至少1KB
		}
	}

	// 存储到缓存
	s.cache.Set(key, value, size, ttl)
}

// Get 从缓存获取对象
func (s *cacheServiceImpl) Get(key string, result interface{}) bool {
	value, exists := s.cache.Get(key)
	if !exists {
		return false
	}

	// 根据不同的类型处理结果
	switch result.(type) {
	case *string:
		// 字符串类型
		if strValue, ok := value.(string); ok {
			*(result.(*string)) = strValue
			return true
		}
	default:
		// 复杂类型，尝试转换
		bytes, err := json.Marshal(value)
		if err != nil {
			return false
		}

		if err := json.Unmarshal(bytes, result); err != nil {
			return false
		}
		return true
	}

	return false
}

// Remove 从缓存中删除对象
func (s *cacheServiceImpl) Remove(key string) {
	s.cache.Remove(key)
}

// GetHitRate 获取缓存命中率
func (s *cacheServiceImpl) GetHitRate() float64 {
	return s.cache.GetHitRate()
}

// Close 关闭缓存服务
func (s *cacheServiceImpl) Close() {
	s.cache.Close()
}
