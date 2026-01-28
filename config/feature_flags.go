package config

import "sync"

// FeatureFlags 功能开关，用于控制功能的启用/禁用
type FeatureFlags struct {
	mu           sync.RWMutex
	EnableCache  bool `yaml:"enable_cache" json:"enable_cache"`
}

// NewFeatureFlags 创建默认启用的FeatureFlags
func NewFeatureFlags() *FeatureFlags {
	return &FeatureFlags{EnableCache: true}
}

// SetCacheEnabled 设置缓存开关
func (f *FeatureFlags) SetCacheEnabled(enabled bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.EnableCache = enabled
}

// IsCacheEnabled 获取缓存开关状态
func (f *FeatureFlags) IsCacheEnabled() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.EnableCache
}
