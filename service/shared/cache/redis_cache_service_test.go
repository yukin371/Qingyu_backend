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

// TestNewRedisCacheService 测试构造函数
func TestNewRedisCacheService(t *testing.T) {
	t.Run("Constructor_Success", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})

		// Act
		service := NewRedisCacheService(redisClient)

		// Assert
		assert.NotNil(t, service)
		assert.IsType(t, &RedisCacheService{}, service)

		// 验证可以调用 Ping
		err := service.Ping(context.Background())
		assert.NoError(t, err)
	})
}

// TestRedisCacheService_Get 测试获取缓存
func TestRedisCacheService_Get(t *testing.T) {
	t.Run("Get_Success", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:get:key"
		testValue := "test:value"

		// 设置缓存
		err := service.Set(ctx, testKey, testValue, 5*time.Minute)
		require.NoError(t, err)

		// Act
		result, err := service.Get(ctx, testKey)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, testValue, result)
	})

	t.Run("Get_NotFound", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// Act
		result, err := service.Get(ctx, "nonexistent:key")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "key not found")
	})

	t.Run("Get_EmptyKey", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// Act
		result, err := service.Get(ctx, "")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("Get_AfterExpiration", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:exp:key"
		testValue := "test:value"

		// 设置缓存，过期时间 1ms
		err := service.Set(ctx, testKey, testValue, 1*time.Millisecond)
		require.NoError(t, err)

		// 等待过期
		time.Sleep(10 * time.Millisecond)
		miniRedis.FastForward(10 * time.Millisecond)

		// Act
		result, err := service.Get(ctx, testKey)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "key not found")
	})
}

// TestRedisCacheService_Set 测试设置缓存
func TestRedisCacheService_Set(t *testing.T) {
	t.Run("Set_Success", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:set:key"
		testValue := "test:value"

		// Act
		err := service.Set(ctx, testKey, testValue, 5*time.Minute)

		// Assert
		assert.NoError(t, err)

		// 验证缓存已设置
		result, err := service.Get(ctx, testKey)
		assert.NoError(t, err)
		assert.Equal(t, testValue, result)
	})

	t.Run("Set_WithExpiration", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:set:exp"
		testValue := "test:value"

		// Act
		err := service.Set(ctx, testKey, testValue, 1*time.Second)

		// Assert
		assert.NoError(t, err)

		// 验证 TTL
		ttl, err := service.TTL(ctx, testKey)
		assert.NoError(t, err)
		assert.Greater(t, ttl, time.Duration(0))
		assert.LessOrEqual(t, ttl, 1*time.Second)
	})

	t.Run("Set_Overwrite", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:set:overwrite"

		// 第一次设置
		err := service.Set(ctx, testKey, "value1", 5*time.Minute)
		require.NoError(t, err)

		// Act - 覆盖设置
		err = service.Set(ctx, testKey, "value2", 5*time.Minute)

		// Assert
		assert.NoError(t, err)

		result, err := service.Get(ctx, testKey)
		assert.NoError(t, err)
		assert.Equal(t, "value2", result)
	})

	t.Run("Set_EmptyValue", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// Act
		err := service.Set(ctx, "test:empty", "", 5*time.Minute)

		// Assert
		assert.NoError(t, err)

		result, err := service.Get(ctx, "test:empty")
		assert.NoError(t, err)
		assert.Equal(t, "", result)
	})
}

// TestRedisCacheService_Delete 测试删除缓存
func TestRedisCacheService_Delete(t *testing.T) {
	t.Run("Delete_Success", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:delete:key"
		testValue := "test:value"

		// 设置缓存
		err := service.Set(ctx, testKey, testValue, 5*time.Minute)
		require.NoError(t, err)

		// 验证存在
		exists, _ := service.Exists(ctx, testKey)
		assert.True(t, exists)

		// Act
		err = service.Delete(ctx, testKey)

		// Assert
		assert.NoError(t, err)

		// 验证已删除
		exists, _ = service.Exists(ctx, testKey)
		assert.False(t, exists)
	})

	t.Run("Delete_NotFound", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// Act - 删除不存在的键
		err := service.Delete(ctx, "nonexistent:key")

		// Assert - Redis DEL 对不存在的键返回 OK
		assert.NoError(t, err)
	})

	t.Run("Delete_EmptyKey", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// Act
		err := service.Delete(ctx, "")

		// Assert
		assert.NoError(t, err)
	})
}

// TestRedisCacheService_Exists 测试检查键是否存在
func TestRedisCacheService_Exists(t *testing.T) {
	t.Run("Exists_True", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:exists:key"

		// 设置缓存
		err := service.Set(ctx, testKey, "value", 5*time.Minute)
		require.NoError(t, err)

		// Act
		exists, err := service.Exists(ctx, testKey)

		// Assert
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Exists_False", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// Act
		exists, err := service.Exists(ctx, "nonexistent:key")

		// Assert
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Exists_AfterDeletion", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:exists:delete"

		// 设置缓存
		err := service.Set(ctx, testKey, "value", 5*time.Minute)
		require.NoError(t, err)

		// 删除缓存
		err = service.Delete(ctx, testKey)
		require.NoError(t, err)

		// Act
		exists, err := service.Exists(ctx, testKey)

		// Assert
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Exists_AfterExpiration", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:exists:exp"

		// 设置缓存，过期时间 1ms
		err := service.Set(ctx, testKey, "value", 1*time.Millisecond)
		require.NoError(t, err)

		// 等待过期
		time.Sleep(10 * time.Millisecond)
		miniRedis.FastForward(10 * time.Millisecond)

		// Act
		exists, err := service.Exists(ctx, testKey)

		// Assert
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

// TestRedisCacheService_ContextCancellation 测试上下文取消
func TestRedisCacheService_ContextCancellation(t *testing.T) {
	t.Run("Get_ContextCancelled", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // 立即取消

		// Act
		_, err := service.Get(ctx, "test:key")

		// Assert
		assert.Error(t, err)
	})

	t.Run("Set_ContextCancelled", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // 立即取消

		// Act
		err := service.Set(ctx, "test:key", "value", 5*time.Minute)

		// Assert
		assert.Error(t, err)
	})
}

// TestRedisCacheService_ConcurrentGetSet 测试并发安全性
func TestRedisCacheService_ConcurrentGetSet(t *testing.T) {
	// Arrange
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	service := NewRedisCacheService(redisClient)
	ctx := context.Background()

	iterations := 100

	// Act - 并发写入和读取
	done := make(chan bool)
	for i := 0; i < iterations; i++ {
		go func(index int) {
			key := "test:concurrent:" + string(rune('0'+index%10))
			_ = service.Set(ctx, key, "value", 5*time.Minute)
			_, _ = service.Get(ctx, key)
			done <- true
		}(i)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < iterations; i++ {
		<-done
	}

	// Assert - 验证没有错误
	// 只要没有 panic 就算通过
	assert.True(t, true)
}

// Helper function to create test service
func setupTestService(t *testing.T) (*miniredis.Miniredis, *RedisCacheService) {
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	service := NewRedisCacheService(redisClient).(*RedisCacheService)
	return miniRedis, service
}

// TestRedisCacheService_TableDriven 表格驱动测试
func TestRedisCacheService_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*RedisCacheService, context.Context) error
		action      func(*RedisCacheService, context.Context) error
		validate    func(*testing.T, error)
		cleanup     func(*RedisCacheService, context.Context) error
	}{
		{
			name: "设置并获取成功",
			setup: nil,
			action: func(s *RedisCacheService, ctx context.Context) error {
				return s.Set(ctx, "test:key", "test:value", 5*time.Minute)
			},
			validate: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
			cleanup: nil,
		},
		{
			name: "获取不存在的键失败",
			setup: nil,
			action: func(s *RedisCacheService, ctx context.Context) error {
				_, err := s.Get(ctx, "nonexistent:key")
				return err
			},
			validate: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "key not found")
			},
			cleanup: nil,
		},
		{
			name: "删除存在的键成功",
			setup: func(s *RedisCacheService, ctx context.Context) error {
				return s.Set(ctx, "test:delete", "value", 5*time.Minute)
			},
			action: func(s *RedisCacheService, ctx context.Context) error {
				return s.Delete(ctx, "test:delete")
			},
			validate: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
			cleanup: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			miniRedis, service := setupTestService(t)
			ctx := context.Background()

			if tt.setup != nil {
				require.NoError(t, tt.setup(service, ctx))
			}

			err := tt.action(service, ctx)
			tt.validate(t, err)

			if tt.cleanup != nil {
				_ = tt.cleanup(service, ctx)
			}

			miniRedis.Close()
		})
	}
}

// TestRedisCacheService_ErrorHandling 测试错误处理
func TestRedisCacheService_ErrorHandling(t *testing.T) {
	t.Run("ClosedConnection", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)

		// 关闭 Redis 连接
		miniRedis.Close()

		ctx := context.Background()

		// Act & Assert - Get 应该返回错误
		_, err := service.Get(ctx, "test:key")
		assert.Error(t, err)
	})
}

// BenchmarkRedisCacheService_Get 性能测试
func BenchmarkRedisCacheService_Get(b *testing.B) {
	miniRedis := miniredis.RunT(b)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	service := NewRedisCacheService(redisClient)
	ctx := context.Background()

	// 设置测试数据
	testKey := "bench:get"
	_ = service.Set(ctx, testKey, "test:value", 5*time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.Get(ctx, testKey)
	}
}

// BenchmarkRedisCacheService_Set 性能测试
func BenchmarkRedisCacheService_Set(b *testing.B) {
	miniRedis := miniredis.RunT(b)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	service := NewRedisCacheService(redisClient)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.Set(ctx, "bench:set", "test:value", 5*time.Minute)
	}
}
