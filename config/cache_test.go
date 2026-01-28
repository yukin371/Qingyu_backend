package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultCacheConfig(t *testing.T) {
	config := DefaultCacheConfig()

	assert.False(t, config.Enabled, "默认应该关闭缓存")
	assert.Equal(t, 1*time.Second, config.DoubleDeleteDelay)
	assert.Equal(t, 30*time.Second, config.NullCacheTTL)
	assert.Equal(t, "@@NULL@@", config.NullCachePrefix)
	assert.Equal(t, uint32(3), config.BreakerMaxRequests)
	assert.Equal(t, 10*time.Second, config.BreakerInterval)
	assert.Equal(t, 30*time.Second, config.BreakerTimeout)
	assert.InDelta(t, 0.6, config.BreakerThreshold, 0.01)
}

func TestGetCacheConfig(t *testing.T) {
	// 重置全局配置
	globalCacheConfig = nil

	config := GetCacheConfig()
	require.NotNil(t, config)
	assert.False(t, config.Enabled, "默认应该关闭缓存")
}

func TestSetCacheConfig(t *testing.T) {
	// 重置全局配置
	globalCacheConfig = nil

	customConfig := &CacheConfig{
		Enabled:           true,
		DoubleDeleteDelay: 2 * time.Second,
		NullCacheTTL:      60 * time.Second,
		NullCachePrefix:   "##NULL##",
		BreakerMaxRequests: 5,
		BreakerInterval:    20 * time.Second,
		BreakerTimeout:     60 * time.Second,
		BreakerThreshold:   0.7,
	}

	SetCacheConfig(customConfig)

	retrieved := GetCacheConfig()
	assert.True(t, retrieved.Enabled)
	assert.Equal(t, 2*time.Second, retrieved.DoubleDeleteDelay)
	assert.Equal(t, 60*time.Second, retrieved.NullCacheTTL)
	assert.Equal(t, "##NULL##", retrieved.NullCachePrefix)
	assert.Equal(t, uint32(5), retrieved.BreakerMaxRequests)
}
