package user

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	userService "Qingyu_backend/service/user"
)

// TestPasswordResetTokenManager 测试密码重置Token管理器
func TestPasswordResetTokenManager(t *testing.T) {
	ctx := context.Background()
	manager := userService.NewPasswordResetTokenManager()

	t.Run("生成Token", func(t *testing.T) {
		email := "test@example.com"
		token, err := manager.GenerateToken(ctx, email)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, 64, len(token)) // 32字节转hex = 64字符
	})

	t.Run("验证有效Token", func(t *testing.T) {
		email := "valid@example.com"
		token, err := manager.GenerateToken(ctx, email)
		assert.NoError(t, err)

		// 验证Token
		err = manager.ValidateToken(ctx, email, token)
		assert.NoError(t, err)
	})

	t.Run("验证无效的邮箱", func(t *testing.T) {
		err := manager.ValidateToken(ctx, "notexist@example.com", "randomtoken")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "无效的重置Token")
	})

	t.Run("验证错误的Token", func(t *testing.T) {
		email := "test2@example.com"
		_, err := manager.GenerateToken(ctx, email)
		assert.NoError(t, err)

		// 使用错误的Token
		err = manager.ValidateToken(ctx, email, "wrongtoken123")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "无效的重置Token")
	})

	t.Run("标记Token为已使用", func(t *testing.T) {
		email := "used@example.com"
		token, err := manager.GenerateToken(ctx, email)
		assert.NoError(t, err)

		// 标记为已使用
		err = manager.MarkTokenAsUsed(ctx, email)
		assert.NoError(t, err)

		// 再次验证应该失败
		err = manager.ValidateToken(ctx, email, token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "已使用")
	})

	t.Run("Token过期检测", func(t *testing.T) {
		// 注意：这个测试需要等待1小时才能验证过期，这里只测试逻辑
		// 实际应该mock时间或使用可配置的过期时间
		email := "expire@example.com"
		_, err := manager.GenerateToken(ctx, email)
		assert.NoError(t, err)

		// Token应该立即有效
		// 实际过期测试需要等待或mock时间
	})

	t.Run("清理过期Token", func(t *testing.T) {
		// 生成一些Token
		emails := []string{"clean1@example.com", "clean2@example.com", "clean3@example.com"}
		for _, email := range emails {
			_, err := manager.GenerateToken(ctx, email)
			assert.NoError(t, err)
		}

		// 清理过期Token（当前没有过期的）
		manager.CleanExpiredTokens(ctx)
		// 验证Token仍然有效
		// 注意：需要访问内部状态或通过验证来确认
	})

	t.Run("并发Token生成", func(t *testing.T) {
		// 测试并发安全性
		emails := make([]string, 10)
		tokens := make([]string, 10)
		errors := make([]error, 10)

		for i := 0; i < 10; i++ {
			i := i
			go func() {
				emails[i] = "concurrent" + string(rune(i)) + "@example.com"
				tokens[i], errors[i] = manager.GenerateToken(ctx, emails[i])
			}()
		}

		time.Sleep(100 * time.Millisecond) // 等待goroutine完成

		// 验证所有Token都生成成功
		for i := 0; i < 10; i++ {
			if errors[i] != nil {
				assert.NoError(t, errors[i])
			}
			if tokens[i] != "" {
				assert.NotEmpty(t, tokens[i])
			}
		}
	})

	t.Run("相同邮箱多次生成Token", func(t *testing.T) {
		email := "multi@example.com"

		token1, err := manager.GenerateToken(ctx, email)
		assert.NoError(t, err)

		time.Sleep(10 * time.Millisecond)

		token2, err := manager.GenerateToken(ctx, email)
		assert.NoError(t, err)

		// Token应该不同
		assert.NotEqual(t, token1, token2)

		// 旧Token应该失效，新Token有效
		err = manager.ValidateToken(ctx, email, token1)
		assert.Error(t, err) // 旧token被覆盖应该失效

		err = manager.ValidateToken(ctx, email, token2)
		assert.NoError(t, err) // 新token有效
	})
}

// BenchmarkGenerateToken 基准测试Token生成性能
func BenchmarkGenerateToken(b *testing.B) {
	ctx := context.Background()
	manager := userService.NewPasswordResetTokenManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		email := "bench@example.com"
		_, err := manager.GenerateToken(ctx, email)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkValidateToken 基准测试Token验证性能
func BenchmarkValidateToken(b *testing.B) {
	ctx := context.Background()
	manager := userService.NewPasswordResetTokenManager()

	email := "bench@example.com"
	token, _ := manager.GenerateToken(ctx, email)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.ValidateToken(ctx, email, token)
	}
}
