package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConstants 测试常量定义
func TestConstants(t *testing.T) {
	assert.Equal(t, "token:blacklist:", DefaultTokenBlacklistPrefix)
	assert.Equal(t, "1", TokenBlacklistValue)
}

// TestDefaultTokenBlacklistConfig 测试默认配置
func TestDefaultTokenBlacklistConfig(t *testing.T) {
	config := DefaultTokenBlacklistConfig()

	assert.NotNil(t, config)
	assert.Equal(t, DefaultTokenBlacklistPrefix, config.Prefix)
}

// TestTokenBlacklistRepository_NewTokenBlacklistRepository 测试创建 Repository
func TestTokenBlacklistRepository_NewTokenBlacklistRepository(t *testing.T) {
	// 使用 nil 客户端创建（仅测试结构，不执行实际操作）
	repo := NewTokenBlacklistRepository(nil)
	assert.NotNil(t, repo)

	redisRepo, ok := repo.(*TokenBlacklistRepositoryRedis)
	assert.True(t, ok, "应为 *TokenBlacklistRepositoryRedis 类型")
	assert.Equal(t, DefaultTokenBlacklistPrefix, redisRepo.config.Prefix)
	assert.Nil(t, redisRepo.client)
}

// TestTokenBlacklistRepository_NewTokenBlacklistRepositoryWithConfig 测试使用自定义配置创建 Repository
func TestTokenBlacklistRepository_NewTokenBlacklistRepositoryWithConfig(t *testing.T) {
	customPrefix := "custom:prefix:"
	config := &TokenBlacklistConfig{
		Prefix: customPrefix,
	}

	repo := NewTokenBlacklistRepositoryWithConfig(nil, config)
	assert.NotNil(t, repo)

	redisRepo, ok := repo.(*TokenBlacklistRepositoryRedis)
	assert.True(t, ok)
	assert.Equal(t, customPrefix, redisRepo.config.Prefix)
}

// TestTokenBlacklistRepository_NilConfig 测试 nil 配置
func TestTokenBlacklistRepository_NilConfig(t *testing.T) {
	repo := NewTokenBlacklistRepositoryWithConfig(nil, nil)
	assert.NotNil(t, repo)

	redisRepo, ok := repo.(*TokenBlacklistRepositoryRedis)
	assert.True(t, ok)
	// nil 配置应该使用默认配置
	assert.Equal(t, DefaultTokenBlacklistPrefix, redisRepo.config.Prefix)
}

// TestTokenBlacklistRepository_EmptyPrefix 测试空前缀
func TestTokenBlacklistRepository_EmptyPrefix(t *testing.T) {
	config := &TokenBlacklistConfig{Prefix: ""}
	repo := NewTokenBlacklistRepositoryWithConfig(nil, config)

	redisRepo, ok := repo.(*TokenBlacklistRepositoryRedis)
	assert.True(t, ok)
	// 空前缀应该被替换为默认前缀
	assert.Equal(t, DefaultTokenBlacklistPrefix, redisRepo.config.Prefix)
}

// TestTokenBlacklistRepository_AddToBlacklist_EmptyToken 测试添加空 Token
func TestTokenBlacklistRepository_AddToBlacklist_EmptyToken(t *testing.T) {
	repo := NewTokenBlacklistRepository(nil)
	ctx := context.Background()

	err := repo.AddToBlacklist(ctx, "", time.Hour)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token不能为空")
}

// TestTokenBlacklistRepository_IsBlacklisted_EmptyToken 测试检查空 Token
func TestTokenBlacklistRepository_IsBlacklisted_EmptyToken(t *testing.T) {
	repo := NewTokenBlacklistRepository(nil)
	ctx := context.Background()

	isBlacklisted, err := repo.IsBlacklisted(ctx, "")
	assert.Error(t, err)
	assert.False(t, isBlacklisted)
	assert.Contains(t, err.Error(), "token不能为空")
}

// TestTokenBlacklistRepository_RemoveFromBlacklist_EmptyToken 测试移除空 Token
func TestTokenBlacklistRepository_RemoveFromBlacklist_EmptyToken(t *testing.T) {
	repo := NewTokenBlacklistRepository(nil)
	ctx := context.Background()

	err := repo.RemoveFromBlacklist(ctx, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token不能为空")
}

// TestTokenBlacklistRepository_ClearExpiredTokens 测试清理过期 Token
func TestTokenBlacklistRepository_ClearExpiredTokens(t *testing.T) {
	repo := NewTokenBlacklistRepository(nil)
	ctx := context.Background()

	// ClearExpiredTokens 应该总是返回 nil（Redis 自动处理过期）
	err := repo.ClearExpiredTokens(ctx)
	assert.NoError(t, err)
}

// TestTokenBlacklistRepository_Health_NilClient 测试健康检查（nil 客户端）
func TestTokenBlacklistRepository_Health_NilClient(t *testing.T) {
	repo := NewTokenBlacklistRepository(nil)
	ctx := context.Background()

	// nil 客户端应该返回错误
	err := repo.Health(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Redis客户端未初始化")
}

// TestTokenBlacklistRepository_KeyPrefix 测试键前缀
func TestTokenBlacklistRepository_KeyPrefix(t *testing.T) {
	tests := []struct {
		name           string
		config         *TokenBlacklistConfig
		expectedPrefix string
	}{
		{
			name:           "默认前缀",
			config:         DefaultTokenBlacklistConfig(),
			expectedPrefix: DefaultTokenBlacklistPrefix,
		},
		{
			name: "自定义前缀",
			config: &TokenBlacklistConfig{
				Prefix: "custom:blacklist:",
			},
			expectedPrefix: "custom:blacklist:",
		},
		{
			name:           "nil 配置使用默认前缀",
			config:         nil,
			expectedPrefix: DefaultTokenBlacklistPrefix,
		},
		{
			name:           "空前缀使用默认前缀",
			config:         &TokenBlacklistConfig{Prefix: ""},
			expectedPrefix: DefaultTokenBlacklistPrefix,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewTokenBlacklistRepositoryWithConfig(nil, tt.config)
			redisRepo, ok := repo.(*TokenBlacklistRepositoryRedis)
			require.True(t, ok)
			assert.Equal(t, tt.expectedPrefix, redisRepo.config.Prefix)
		})
	}
}

// TestTokenBlacklistRepository_ErrorMessages 测试错误消息
func TestTokenBlacklistRepository_ErrorMessages(t *testing.T) {
	repo := NewTokenBlacklistRepository(nil)
	ctx := context.Background()

	t.Run("AddToBlacklist 错误消息", func(t *testing.T) {
		err := repo.AddToBlacklist(ctx, "", time.Hour)
		assert.Contains(t, err.Error(), "token不能为空")
	})

	t.Run("IsBlacklisted 错误消息", func(t *testing.T) {
		_, err := repo.IsBlacklisted(ctx, "")
		assert.Contains(t, err.Error(), "token不能为空")
	})

	t.Run("RemoveFromBlacklist 错误消息", func(t *testing.T) {
		err := repo.RemoveFromBlacklist(ctx, "")
		assert.Contains(t, err.Error(), "token不能为空")
	})

	t.Run("Health 错误消息", func(t *testing.T) {
		err := repo.Health(ctx)
		assert.Contains(t, err.Error(), "Redis客户端未初始化")
	})
}
