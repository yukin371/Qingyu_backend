package strategies

import "time"

// TTLStrategy 固定TTL策略
type TTLStrategy struct {
	ttl time.Duration
}

// NewTTLStrategy 创建固定TTL策略
func NewTTLStrategy(ttl time.Duration) *TTLStrategy {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return &TTLStrategy{ttl: ttl}
}

// ShouldCache 判断是否需要缓存
func (s *TTLStrategy) ShouldCache(key string, value interface{}) bool {
	_ = key
	_ = value
	return true
}

// GetTTL 返回固定TTL
func (s *TTLStrategy) GetTTL(key string) time.Duration {
	_ = key
	return s.ttl
}

// OnMiss 缓存miss处理（基础策略无副作用）
func (s *TTLStrategy) OnMiss(key string) error {
	_ = key
	return nil
}
