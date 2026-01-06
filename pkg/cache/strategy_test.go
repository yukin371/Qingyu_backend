package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestNewCacheStrategy 测试缓存策略管理器创建
func TestNewCacheStrategy(t *testing.T) {
	mockClient := new(MockRedisClient)

	strategy := NewCacheStrategy(mockClient)

	assert.NotNil(t, strategy)
	assert.NotNil(t, strategy.policies)
	assert.NotNil(t, strategy.client)
}

// TestCacheStrategy_GetPolicy 测试获取缓存策略
func TestCacheStrategy_GetPolicy(t *testing.T) {
	mockClient := new(MockRedisClient)
	strategy := NewCacheStrategy(mockClient)

	tests := []struct {
		name           string
		key            string
		expectFound    bool
		expectedTTL    time.Duration
		expectedStrTTL time.Duration // 至少应该有的TTL
	}{
		{
			name:           "用户会话策略",
			key:            "user:session:123",
			expectFound:    true,
			expectedTTL:    7 * 24 * time.Hour,
			expectedStrTTL: 6 * 24 * time.Hour,
		},
		{
			name:           "用户信息策略",
			key:            "user:info:456",
			expectFound:    true,
			expectedTTL:    30 * time.Minute,
			expectedStrTTL: 25 * time.Minute,
		},
		{
			name:           "书籍详情策略",
			key:            "book:detail:789",
			expectFound:    true,
			expectedTTL:    1 * time.Hour,
			expectedStrTTL: 50 * time.Minute,
		},
		{
			name:           "热门列表策略",
			key:            "book:hot:list",
			expectFound:    true,
			expectedTTL:    15 * time.Minute,
			expectedStrTTL: 12 * time.Minute,
		},
		{
			name:           "未定义的策略（使用默认）",
			key:            "unknown:key:value",
			expectFound:    false,
			expectedTTL:    5 * time.Minute,
			expectedStrTTL: 4 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy, found := strategy.GetPolicy(tt.key)

			if tt.expectFound {
				assert.True(t, found, "应该找到预定义策略")
				assert.GreaterOrEqual(t, policy.TTL, tt.expectedStrTTL, "TTL应该大于等于最小值")
			} else {
				assert.False(t, found, "不应该找到预定义策略")
				assert.Equal(t, tt.expectedTTL, policy.TTL, "应该使用默认TTL")
			}
		})
	}
}

// TestCacheStrategy_Set 测试设置缓存
func TestCacheStrategy_Set(t *testing.T) {
	mockClient := new(MockRedisClient)
	strategy := NewCacheStrategy(mockClient)

	ctx := context.Background()
	testKey := "book:detail:123"
	testValue := map[string]interface{}{
		"id":    "123",
		"title": "测试书籍",
	}

	// 设置Mock预期
	mockClient.On("Set", ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// 执行测试
	err := strategy.Set(ctx, testKey, testValue)

	// 验证结果
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

// TestCacheStrategy_Get 测试获取缓存
func TestCacheStrategy_Get(t *testing.T) {
	mockClient := new(MockRedisClient)
	strategy := NewCacheStrategy(mockClient)

	ctx := context.Background()
	testKey := "book:detail:123"
	testValue := `{"id":"123","title":"测试书籍"}`

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "成功获取缓存",
			setup: func() {
				mockClient.On("Get", ctx, testKey).Return(testValue, nil)
			},
			wantErr: false,
		},
		{
			name: "缓存不存在",
			setup: func() {
				mockClient.On("Get", ctx, testKey).Return("", ErrRedisNil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 重置Mock
			mockClient.ExpectedCalls = nil

			// 设置预期
			tt.setup()

			// 执行测试
			var result map[string]interface{}
			err := strategy.Get(ctx, testKey, &result)

			// 验证结果
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "123", result["id"])
				assert.Equal(t, "测试书籍", result["title"])
			}

			mockClient.AssertExpectations(t)
		})
	}
}

// TestCacheStrategy_Delete 测试删除缓存
func TestCacheStrategy_Delete(t *testing.T) {
	mockClient := new(MockRedisClient)
	strategy := NewCacheStrategy(mockClient)

	ctx := context.Background()
	testKeys := []string{"book:detail:123", "book:detail:456"}

	// 设置Mock预期
	mockClient.On("Delete", ctx, testKeys).Return(nil)

	// 执行测试
	err := strategy.Delete(ctx, testKeys...)

	// 验证结果
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

// TestBuildKey 测试缓存Key构建
func TestBuildKey(t *testing.T) {
	tests := []struct {
		name     string
		parts    []string
		expected string
	}{
		{
			name:     "单个部分",
			parts:    []string{"user"},
			expected: "user",
		},
		{
			name:     "多个部分",
			parts:    []string{"user", "info", "123"},
			expected: "user:info:123",
		},
		{
			name:     "书籍key",
			parts:    []string{"book", "detail", "456"},
			expected: "book:detail:456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildKey(tt.parts...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestBuildUserKey 测试用户Key构建
func TestBuildUserKey(t *testing.T) {
	result := BuildUserKey("123", "info")
	assert.Equal(t, "user:info:123", result)
}

// TestBuildBookKey 测试书籍Key构建
func TestBuildBookKey(t *testing.T) {
	result := BuildBookKey("456", "detail")
	assert.Equal(t, "book:detail:456", result)
}

// TestBuildChapterKey 测试章节Key构建
func TestBuildChapterKey(t *testing.T) {
	result := BuildChapterKey("789", 1)
	assert.Equal(t, "book:chapter:789:1", result)
}

// TestBuildSearchKey 测试搜索Key构建
func TestBuildSearchKey(t *testing.T) {
	result := BuildSearchKey("测试", 1)
	assert.Equal(t, "search:result:测试:page:1", result)
}

// TestCacheStrategy_calculateTTL 测试TTL计算（包含随机值）
func TestCacheStrategy_calculateTTL(t *testing.T) {
	mockClient := new(MockRedisClient)
	strategy := NewCacheStrategy(mockClient)

	baseTTL := 1 * time.Hour
	randomTTL := 10 * time.Minute

	// 测试无随机TTL的情况
	ttl := strategy.calculateTTL(baseTTL, 0)
	assert.Equal(t, baseTTL, ttl, "无随机TTL时应返回基础TTL")

	// 测试有随机TTL的情况（多次调用，验证基本功能）
	for i := 0; i < 5; i++ {
		ttl := strategy.calculateTTL(baseTTL, randomTTL)
		assert.NotZero(t, ttl, "TTL不应为零")
	}
}

// BenchmarkCacheStrategy_Set 性能基准测试
func BenchmarkCacheStrategy_Set(b *testing.B) {
	mockClient := new(MockRedisClient)
	mockClient.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	strategy := NewCacheStrategy(mockClient)
	ctx := context.Background()
	testKey := "book:detail:123"
	testValue := map[string]interface{}{
		"id":    "123",
		"title": "测试书籍",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strategy.Set(ctx, testKey, testValue)
	}
}
