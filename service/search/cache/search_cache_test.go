package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSearchCache(t *testing.T) {
	defaultTTL := 5 * time.Minute
	hotTTL := 10 * time.Minute

	cache, err := NewSearchCache(defaultTTL, hotTTL)
	require.NoError(t, err)
	require.NotNil(t, cache)

	assert.Equal(t, defaultTTL, cache.defaultTTL)
	assert.Equal(t, hotTTL, cache.hotTTL)
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

func TestGetTTLForPage(t *testing.T) {
	defaultTTL := 5 * time.Minute
	hotTTL := 10 * time.Minute

	cache, err := NewSearchCache(defaultTTL, hotTTL)
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
