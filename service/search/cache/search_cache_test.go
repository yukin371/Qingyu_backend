package cache

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSearchCache(t *testing.T) {
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})

	defaultTTL := 5 * time.Minute
	hotTTL := 10 * time.Minute

	cache, err := NewSearchCache(redisClient, defaultTTL, hotTTL)
	require.NoError(t, err)
	require.NotNil(t, cache)

	assert.Equal(t, defaultTTL, cache.defaultTTL)
	assert.Equal(t, hotTTL, cache.hotTTL)
	assert.True(t, cache.IsEnabled())
}

func TestGenerateCacheKey(t *testing.T) {
	tests := []struct {
		name      string
		searchType string
		query     string
		page      int
		pageSize  int
		filter    map[string]interface{}
	}{
		{
			name:      "basic search",
			searchType: "books",
			query:     "哈利波特",
			page:      1,
			pageSize:  20,
			filter:    map[string]interface{}{},
		},
		{
			name:      "search with filter",
			searchType: "projects",
			query:     "test",
			page:      1,
			pageSize:  10,
			filter:    map[string]interface{}{"user_id": "123"},
		},
		{
			name:      "second page",
			searchType: "documents",
			query:     "search",
			page:      2,
			pageSize:  10,
			filter:    map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := GenerateCacheKey(tt.searchType, tt.query, tt.page, tt.pageSize, tt.filter)
			assert.NotEmpty(t, key)
			assert.Contains(t, key, tt.searchType)
			assert.Contains(t, key, tt.query)
		})
	}
}

func TestCacheGetSet(t *testing.T) {
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})

	cache, err := NewSearchCache(redisClient, 5*time.Minute, 10*time.Minute)
	require.NoError(t, err)

	ctx := context.Background()
	testKey := "test:cache:key"
	testValue := []byte(`{"result": "test data"}`)

	// 设置缓存
	err = cache.Set(ctx, testKey, testValue, 1*time.Minute)
	require.NoError(t, err)

	// 获取缓存
	result, err := cache.Get(ctx, testKey)
	require.NoError(t, err)
	assert.Equal(t, testValue, result)

	// 验证统计信息
	stats, err := cache.Stats(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(1), stats.Hits)
	assert.Equal(t, int64(0), stats.Misses)
}

func TestCacheMiss(t *testing.T) {
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})

	cache, err := NewSearchCache(redisClient, 5*time.Minute, 10*time.Minute)
	require.NoError(t, err)

	ctx := context.Background()
	testKey := "nonexistent:key"

	// 获取不存在的缓存
	result, err := cache.Get(ctx, testKey)
	require.NoError(t, err)
	assert.Nil(t, result)

	// 验证统计信息
	stats, err := cache.Stats(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(0), stats.Hits)
	assert.Equal(t, int64(1), stats.Misses)
}

func TestCacheDelete(t *testing.T) {
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})

	cache, err := NewSearchCache(redisClient, 5*time.Minute, 10*time.Minute)
	require.NoError(t, err)

	ctx := context.Background()
	testKey := "test:delete:key"
	testValue := []byte(`{"result": "data"}`)

	// 设置缓存
	err = cache.Set(ctx, testKey, testValue, 1*time.Minute)
	require.NoError(t, err)

	// 删除缓存
	err = cache.Delete(ctx, testKey)
	require.NoError(t, err)

	// 验证缓存已删除
	result, err := cache.Get(ctx, testKey)
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestCacheExists(t *testing.T) {
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})

	cache, err := NewSearchCache(redisClient, 5*time.Minute, 10*time.Minute)
	require.NoError(t, err)

	ctx := context.Background()
	testKey := "test:exists:key"
	testValue := []byte(`{"result": "data"}`)

	// 不存在时
	exists, err := cache.Exists(ctx, testKey)
	require.NoError(t, err)
	assert.False(t, exists)

	// 设置缓存
	err = cache.Set(ctx, testKey, testValue, 1*time.Minute)
	require.NoError(t, err)

	// 存在时
	exists, err = cache.Exists(ctx, testKey)
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestCacheClear(t *testing.T) {
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})

	cache, err := NewSearchCache(redisClient, 5*time.Minute, 10*time.Minute)
	require.NoError(t, err)

	ctx := context.Background()

	// 设置多个缓存
	for i := 0; i < 5; i++ {
		key := "test:clear:key:" + string(rune('0'+i))
		value := []byte(`{"result": "data"}`)
		err = cache.Set(ctx, key, value, 1*time.Minute)
		require.NoError(t, err)
	}

	// 清空所有缓存
	err = cache.Clear(ctx)
	require.NoError(t, err)

	// 验证所有缓存已清空
	stats, err := cache.Stats(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(0), stats.Keys)
}

func TestGetTTLForPage(t *testing.T) {
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})

	defaultTTL := 5 * time.Minute
	hotTTL := 10 * time.Minute

	cache, err := NewSearchCache(redisClient, defaultTTL, hotTTL)
	require.NoError(t, err)

	// 第一页应该使用 hotTTL
	ttl := cache.GetTTLForPage(1)
	assert.Equal(t, hotTTL, ttl)

	// 其他页应该使用 defaultTTL
	ttl = cache.GetTTLForPage(2)
	assert.Equal(t, defaultTTL, ttl)

	ttl = cache.GetTTLForPage(10)
	assert.Equal(t, defaultTTL, ttl)
}

func TestCacheStats(t *testing.T) {
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})

	cache, err := NewSearchCache(redisClient, 5*time.Minute, 10*time.Minute)
	require.NoError(t, err)

	ctx := context.Background()

	// 初始统计
	stats, err := cache.Stats(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(0), stats.Hits)
	assert.Equal(t, int64(0), stats.Misses)
	assert.Equal(t, 0.0, stats.HitRate)

	// 设置并获取缓存
	testKey := "test:stats:key"
	testValue := []byte(`{"result": "data"}`)
	err = cache.Set(ctx, testKey, testValue, 1*time.Minute)
	require.NoError(t, err)

	_, _ = cache.Get(ctx, testKey)
	_, _ = cache.Get(ctx, "nonexistent:key")

	// 更新后的统计
	stats, err = cache.Stats(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(1), stats.Hits)
	assert.Equal(t, int64(1), stats.Misses)
	assert.Equal(t, 0.5, stats.HitRate)
}

func TestCacheReset(t *testing.T) {
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})

	cache, err := NewSearchCache(redisClient, 5*time.Minute, 10*time.Minute)
	require.NoError(t, err)

	ctx := context.Background()

	// 设置并获取缓存
	testKey := "test:reset:key"
	testValue := []byte(`{"result": "data"}`)
	err = cache.Set(ctx, testKey, testValue, 1*time.Minute)
	require.NoError(t, err)

	_, _ = cache.Get(ctx, testKey)

	// 重置统计
	err = cache.Reset(ctx)
	require.NoError(t, err)

	// 验证统计已重置
	stats, err := cache.Stats(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(0), stats.Hits)
	assert.Equal(t, int64(0), stats.Misses)
}

func TestSetEnabled(t *testing.T) {
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})

	cache, err := NewSearchCache(redisClient, 5*time.Minute, 10*time.Minute)
	require.NoError(t, err)

	// 默认启用
	assert.True(t, cache.IsEnabled())

	// 禁用缓存
	cache.SetEnabled(false)
	assert.False(t, cache.IsEnabled())

	// 启用缓存
	cache.SetEnabled(true)
	assert.True(t, cache.IsEnabled())
}

func TestCacheDisabledBehavior(t *testing.T) {
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})

	cache, err := NewSearchCache(redisClient, 5*time.Minute, 10*time.Minute)
	require.NoError(t, err)

	ctx := context.Background()
	testKey := "test:disabled:key"
	testValue := []byte(`{"result": "data"}`)

	// 禁用缓存
	cache.SetEnabled(false)

	// 设置缓存应该不报错
	err = cache.Set(ctx, testKey, testValue, 1*time.Minute)
	assert.NoError(t, err)

	// 获取缓存应该返回 disabled 错误
	result, err := cache.Get(ctx, testKey)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPing(t *testing.T) {
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})

	cache, err := NewSearchCache(redisClient, 5*time.Minute, 10*time.Minute)
	require.NoError(t, err)

	ctx := context.Background()

	// 正常情况下 ping 应该成功
	err = cache.Ping(ctx)
	assert.NoError(t, err)
}
