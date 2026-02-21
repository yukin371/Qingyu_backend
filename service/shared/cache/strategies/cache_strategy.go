package strategies

import (
	"strings"
	"sync"
	"time"
)

// CacheStrategy 缓存策略接口
type CacheStrategy interface {
	ShouldCache(key string, value interface{}) bool
	GetTTL(key string) time.Duration
	OnMiss(key string) error
}

// CacheManager 策略管理器
type CacheManager interface {
	RegisterStrategy(pattern string, strategy CacheStrategy)
	GetStrategy(key string) CacheStrategy
}

// StrategyManager CacheManager 默认实现
type StrategyManager struct {
	mu         sync.RWMutex
	defaultTTL time.Duration
	strategies map[string]CacheStrategy
}

// NewStrategyManager 创建策略管理器
func NewStrategyManager(defaultTTL time.Duration) *StrategyManager {
	if defaultTTL <= 0 {
		defaultTTL = 5 * time.Minute
	}
	return &StrategyManager{
		defaultTTL: defaultTTL,
		strategies: make(map[string]CacheStrategy),
	}
}

// RegisterStrategy 注册键模式对应策略
func (m *StrategyManager) RegisterStrategy(pattern string, strategy CacheStrategy) {
	if strings.TrimSpace(pattern) == "" || strategy == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.strategies[pattern] = strategy
}

// GetStrategy 获取与key匹配的策略，未匹配返回默认TTL策略
func (m *StrategyManager) GetStrategy(key string) CacheStrategy {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for pattern, strategy := range m.strategies {
		if strings.HasPrefix(key, pattern) {
			return strategy
		}
	}

	return NewTTLStrategy(m.defaultTTL)
}
