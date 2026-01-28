package config

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFeatureFlags_SetCacheEnabled(t *testing.T) {
	flags := NewFeatureFlags()

	// 默认应该启用缓存
	assert.True(t, flags.IsCacheEnabled())

	// 禁用缓存
	flags.SetCacheEnabled(false)
	assert.False(t, flags.IsCacheEnabled())

	// 重新启用缓存
	flags.SetCacheEnabled(true)
	assert.True(t, flags.IsCacheEnabled())
}

func TestFeatureFlags_ConcurrentAccess(t *testing.T) {
	flags := NewFeatureFlags()
	var wg sync.WaitGroup

	// 启动100个并发goroutine，同时读写
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			flags.IsCacheEnabled()
		}()
		go func(i int) {
			defer wg.Done()
			flags.SetCacheEnabled(i%2 == 0)
		}(i)
	}

	wg.Wait()
	// 如果没有panic和数据竞争，测试就通过
	assert.True(t, true)
}
