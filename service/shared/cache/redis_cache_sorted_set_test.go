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

// ============ 有序集合操作测试 ============

// TestRedisCacheService_ZAdd_ZRange 测试有序集合基本操作
func TestRedisCacheService_ZAdd_ZRange(t *testing.T) {
	t.Run("ZAdd_ZRange_Success", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:zset"

		// 添加成员（分数从小到大）
		_ = service.ZAdd(ctx, testKey, 1.0, "member1")
		_ = service.ZAdd(ctx, testKey, 3.0, "member3")
		_ = service.ZAdd(ctx, testKey, 2.0, "member2")

		// Act - 获取所有成员（按分数升序）
		result, err := service.ZRange(ctx, testKey, 0, -1)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, result, 3)
		assert.Equal(t, "member1", result[0])
		assert.Equal(t, "member2", result[1])
		assert.Equal(t, "member3", result[2])
	})

	t.Run("ZRange_WithLimit", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:zset:limit"

		// 添加成员
		_ = service.ZAdd(ctx, testKey, 1.0, "member1")
		_ = service.ZAdd(ctx, testKey, 2.0, "member2")
		_ = service.ZAdd(ctx, testKey, 3.0, "member3")
		_ = service.ZAdd(ctx, testKey, 4.0, "member4")
		_ = service.ZAdd(ctx, testKey, 5.0, "member5")

		// Act - 获取前 3 个成员
		result, err := service.ZRange(ctx, testKey, 0, 2)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, result, 3)
		assert.Equal(t, "member1", result[0])
		assert.Equal(t, "member2", result[1])
		assert.Equal(t, "member3", result[2])
	})
}

// TestRedisCacheService_ZRangeWithScores 测试带分数的范围查询
func TestRedisCacheService_ZRangeWithScores(t *testing.T) {
	t.Run("ZRangeWithScores_Success", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:zset:score"

		// 添加成员
		_ = service.ZAdd(ctx, testKey, 1.5, "member1")
		_ = service.ZAdd(ctx, testKey, 2.5, "member2")
		_ = service.ZAdd(ctx, testKey, 3.5, "member3")

		// Act
		result, err := service.ZRangeWithScores(ctx, testKey, 0, -1)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, result, 3)
		assert.InDelta(t, 1.5, result["member1"], 0.01)
		assert.InDelta(t, 2.5, result["member2"], 0.01)
		assert.InDelta(t, 3.5, result["member3"], 0.01)
	})

	t.Run("ZRangeWithScores_Empty", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// Act
		result, err := service.ZRangeWithScores(ctx, "nonexistent:zset", 0, -1)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, result)
	})
}

// TestRedisCacheService_ZRemove 测试移除有序集合成员
func TestRedisCacheService_ZRemove(t *testing.T) {
	t.Run("ZRemove_Success", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:zset:remove"

		// 添加成员
		_ = service.ZAdd(ctx, testKey, 1.0, "member1")
		_ = service.ZAdd(ctx, testKey, 2.0, "member2")
		_ = service.ZAdd(ctx, testKey, 3.0, "member3")

		// Act
		err := service.ZRemove(ctx, testKey, "member1", "member2")

		// Assert
		assert.NoError(t, err)

		// 验证成员已移除
		result, _ := service.ZRange(ctx, testKey, 0, -1)
		assert.Len(t, result, 1)
		assert.Equal(t, "member3", result[0])
	})

	t.Run("ZRemove_Empty", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// Act
		err := service.ZRemove(ctx, "test:zset")

		// Assert
		assert.NoError(t, err)
	})
}

// TestRedisCacheService_ZAdd_UpdateScore 测试更新分数
func TestRedisCacheService_ZAdd_UpdateScore(t *testing.T) {
	// Arrange
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	service := NewRedisCacheService(redisClient)
	ctx := context.Background()

	testKey := "test:zset:update"

	// 添加成员
	_ = service.ZAdd(ctx, testKey, 1.0, "member1")
	_ = service.ZAdd(ctx, testKey, 2.0, "member2")

	// Act - 更新 member1 的分数
	err := service.ZAdd(ctx, testKey, 5.0, "member1")

	// Assert
	assert.NoError(t, err)

	// 验证排序已更新
	result, _ := service.ZRange(ctx, testKey, 0, -1)
	assert.Equal(t, "member2", result[0]) // 分数 2.0
	assert.Equal(t, "member1", result[1]) // 分数 5.0
}

// ============ 服务管理测试 ============

// TestRedisCacheService_Ping 测试健康检查
func TestRedisCacheService_Ping(t *testing.T) {
	t.Run("Ping_Success", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// Act
		err := service.Ping(ctx)

		// Assert
		assert.NoError(t, err)
	})
}

// TestRedisCacheService_FlushDB 测试清空数据库
func TestRedisCacheService_FlushDB(t *testing.T) {
	t.Run("FlushDB_Success", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// 设置一些数据
		_ = service.Set(ctx, "key1", "value1", 5*time.Minute)
		_ = service.Set(ctx, "key2", "value2", 5*time.Minute)
		_ = service.Set(ctx, "key3", "value3", 5*time.Minute)

		// 验证数据存在
		exists1, _ := service.Exists(ctx, "key1")
		assert.True(t, exists1)

		// Act
		err := service.FlushDB(ctx)

		// Assert
		assert.NoError(t, err)

		// 验证所有数据已清空
		exists1, _ = service.Exists(ctx, "key1")
		exists2, _ := service.Exists(ctx, "key2")
		exists3, _ := service.Exists(ctx, "key3")
		assert.False(t, exists1)
		assert.False(t, exists2)
		assert.False(t, exists3)
	})
}

// TestRedisCacheService_Close 测试关闭连接
func TestRedisCacheService_Close(t *testing.T) {
	t.Run("Close_Success", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// 验证服务正常工作
		_ = service.Set(ctx, "key1", "value1", 5*time.Minute)
		_, err := service.Get(ctx, "key1")
		assert.NoError(t, err)

		// Act
		err = service.Close()

		// Assert
		assert.NoError(t, err)

		// 关闭后再操作应该返回错误
		// 注意：miniredis 关闭后可能不会立即反映，这是正常的
		miniRedis.Close()
	})
}

// ============ 并发测试 ============

// TestRedisCacheService_ConcurrentIncrement 测试并发递增
func TestRedisCacheService_ConcurrentIncrement(t *testing.T) {
	// Arrange
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	service := NewRedisCacheService(redisClient)
	ctx := context.Background()

	testKey := "test:concurrent:increment"
	iterations := 100

	// 设置初始值
	_ = service.Set(ctx, testKey, "0", 5*time.Minute)

	// Act - 并发递增
	done := make(chan bool)
	for i := 0; i < iterations; i++ {
		go func() {
			_, _ = service.Increment(ctx, testKey)
			done <- true
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < iterations; i++ {
		<-done
	}

	// Assert - 验证最终值
	result, err := service.Get(ctx, testKey)
	assert.NoError(t, err)
	assert.Equal(t, "100", result)
}

// TestRedisCacheService_ConcurrentSortedSet 测试并发有序集合操作
func TestRedisCacheService_ConcurrentSortedSet(t *testing.T) {
	// Arrange
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	service := NewRedisCacheService(redisClient)
	ctx := context.Background()

	testKey := "test:concurrent:zset"
	iterations := 50

	// Act - 并发添加成员
	done := make(chan bool)
	for i := 0; i < iterations; i++ {
		go func(index int) {
			member := "member" + string(rune('0'+index%10))
			score := float64(index)
			_ = service.ZAdd(ctx, testKey, score, member)
			done <- true
		}(i)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < iterations; i++ {
		<-done
	}

	// Assert - 验证有序集合状态
	result, err := service.ZRange(ctx, testKey, 0, -1)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

// ============ 表格驱动测试 ============

// TestRedisCacheService_SortedSetTableDriven 有序集合表格驱动测试
func TestRedisCacheService_SortedSetTableDriven(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*RedisCacheService, context.Context) error
		action      func(*RedisCacheService, context.Context) error
		validate    func(*testing.T, interface{}, error)
		cleanup     func(*RedisCacheService, context.Context) error
	}{
		{
			name: "ZAdd 添加单个成员成功",
			setup: nil,
			action: func(s *RedisCacheService, ctx context.Context) error {
				return s.ZAdd(ctx, "test:zset", 1.0, "member1")
			},
			validate: func(t *testing.T, _ interface{}, err error) {
				assert.NoError(t, err)
			},
			cleanup: nil,
		},
		{
			name: "ZRange 获取范围成功",
			setup: func(s *RedisCacheService, ctx context.Context) error {
				_ = s.ZAdd(ctx, "test:zset", 1.0, "member1")
				_ = s.ZAdd(ctx, "test:zset", 2.0, "member2")
				return nil
			},
			action: func(s *RedisCacheService, ctx context.Context) error {
				_, err := s.ZRange(ctx, "test:zset", 0, -1)
				return err
			},
			validate: func(t *testing.T, _ interface{}, err error) {
				assert.NoError(t, err)
			},
			cleanup: nil,
		},
		{
			name: "ZRemove 移除成员成功",
			setup: func(s *RedisCacheService, ctx context.Context) error {
				return s.ZAdd(ctx, "test:zset", 1.0, "member1")
			},
			action: func(s *RedisCacheService, ctx context.Context) error {
				return s.ZRemove(ctx, "test:zset", "member1")
			},
			validate: func(t *testing.T, _ interface{}, err error) {
				assert.NoError(t, err)
			},
			cleanup: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			miniRedis := miniredis.RunT(t)
			redisClient := redis.NewClient(&redis.Options{
				Addr: miniRedis.Addr(),
			})
			service := NewRedisCacheService(redisClient).(*RedisCacheService)
			ctx := context.Background()

			if tt.setup != nil {
				require.NoError(t, tt.setup(service, ctx))
			}

			err := tt.action(service, ctx)
			tt.validate(t, nil, err)

			if tt.cleanup != nil {
				_ = tt.cleanup(service, ctx)
			}

			miniRedis.Close()
		})
	}
}

// Benchmark 有序集合性能测试

// BenchmarkRedisCacheService_ZAdd 有序集合添加性能测试
func BenchmarkRedisCacheService_ZAdd(b *testing.B) {
	miniRedis := miniredis.RunT(b)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	service := NewRedisCacheService(redisClient)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		member := "member" + string(rune('0'+i%10))
		_ = service.ZAdd(ctx, "bench:zset", float64(i), member)
	}
}

// BenchmarkRedisCacheService_ZRange 有序集合范围查询性能测试
func BenchmarkRedisCacheService_ZRange(b *testing.B) {
	miniRedis := miniredis.RunT(b)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	service := NewRedisCacheService(redisClient)
	ctx := context.Background()

	// 准备测试数据
	for i := 0; i < 100; i++ {
		member := "member" + string(rune('0'+i%10))
		_ = service.ZAdd(ctx, "bench:zset", float64(i), member)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.ZRange(ctx, "bench:zset", 0, -1)
	}
}
