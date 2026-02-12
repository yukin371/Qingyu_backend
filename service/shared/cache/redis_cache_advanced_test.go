package cache

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// ============ 批量操作测试 ============

// TestRedisCacheService_MGet 测试批量获取
func TestRedisCacheService_MGet(t *testing.T) {
	t.Run("MGet_Success", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// 设置多个缓存
		_ = service.Set(ctx, "key1", "value1", 5*time.Minute)
		_ = service.Set(ctx, "key2", "value2", 5*time.Minute)
		_ = service.Set(ctx, "key3", "value3", 5*time.Minute)

		// Act
		results, err := service.MGet(ctx, "key1", "key2", "key3")

		// Assert
		assert.NoError(t, err)
		assert.Len(t, results, 3)
		assert.Equal(t, "value1", results[0])
		assert.Equal(t, "value2", results[1])
		assert.Equal(t, "value3", results[2])
	})

	t.Run("MGet_PartialMiss", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// 只设置部分键
		_ = service.Set(ctx, "key1", "value1", 5*time.Minute)
		_ = service.Set(ctx, "key3", "value3", 5*time.Minute)

		// Act
		results, err := service.MGet(ctx, "key1", "key2", "key3")

		// Assert
		assert.NoError(t, err)
		assert.Len(t, results, 3)
		assert.Equal(t, "value1", results[0])
		assert.Equal(t, "", results[1]) // key2 不存在
		assert.Equal(t, "value3", results[2])
	})

	t.Run("MGet_Empty", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// Act - Redis MGET 命令不接受空参数，所以传递一个不存在的键
		results, err := service.MGet(ctx, "nonexistent:empty:key")

		// Assert
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "", results[0]) // 不存在的键返回空字符串
	})
}

// TestRedisCacheService_MSet 测试批量设置
func TestRedisCacheService_MSet(t *testing.T) {
	t.Run("MSet_Success", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		kvPairs := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}

		// Act
		err := service.MSet(ctx, kvPairs, 5*time.Minute)

		// Assert
		assert.NoError(t, err)

		// 验证所有键都已设置
		val1, _ := service.Get(ctx, "key1")
		val2, _ := service.Get(ctx, "key2")
		val3, _ := service.Get(ctx, "key3")
		assert.Equal(t, "value1", val1)
		assert.Equal(t, "value2", val2)
		assert.Equal(t, "value3", val3)
	})

	t.Run("MSet_Empty", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// Act
		err := service.MSet(ctx, map[string]string{}, 5*time.Minute)

		// Assert
		assert.NoError(t, err)
	})
}

// TestRedisCacheService_MDelete 测试批量删除
func TestRedisCacheService_MDelete(t *testing.T) {
	t.Run("MDelete_Success", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// 设置多个键
		_ = service.Set(ctx, "key1", "value1", 5*time.Minute)
		_ = service.Set(ctx, "key2", "value2", 5*time.Minute)
		_ = service.Set(ctx, "key3", "value3", 5*time.Minute)

		// Act
		err := service.MDelete(ctx, "key1", "key2", "key3")

		// Assert
		assert.NoError(t, err)

		// 验证所有键都已删除
		exists1, _ := service.Exists(ctx, "key1")
		exists2, _ := service.Exists(ctx, "key2")
		exists3, _ := service.Exists(ctx, "key3")
		assert.False(t, exists1)
		assert.False(t, exists2)
		assert.False(t, exists3)
	})

	t.Run("MDelete_Empty", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// Act
		err := service.MDelete(ctx)

		// Assert
		assert.NoError(t, err)
	})
}

// ============ 高级操作测试 ============

// TestRedisCacheService_Expire 测试设置过期时间
func TestRedisCacheService_Expire(t *testing.T) {
	t.Run("Expire_Success", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:expire"

		// 设置缓存（无过期时间）
		_ = service.Set(ctx, testKey, "value", 0)

		// Act
		err := service.Expire(ctx, testKey, 1*time.Minute)

		// Assert
		assert.NoError(t, err)

		// 验证 TTL 已设置
		ttl, _ := service.TTL(ctx, testKey)
		assert.Greater(t, ttl, time.Duration(0))
	})

	t.Run("Expire_NotFound", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// Act - 对不存在的键设置过期时间
		err := service.Expire(ctx, "nonexistent:key", 1*time.Minute)

		// Assert
		// Redis EXPIRE 对不存在的键返回 0，但不报错
		assert.NoError(t, err)
	})
}

// TestRedisCacheService_TTL 测试获取剩余过期时间
func TestRedisCacheService_TTL(t *testing.T) {
	t.Run("TTL_Success", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:ttl"

		// 设置带过期时间的缓存
		_ = service.Set(ctx, testKey, "value", 1*time.Minute)

		// Act
		ttl, err := service.TTL(ctx, testKey)

		// Assert
		assert.NoError(t, err)
		assert.Greater(t, ttl, time.Duration(0))
		assert.LessOrEqual(t, ttl, 1*time.Minute)
	})

	t.Run("TTL_NoKey", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// Act
		ttl, err := service.TTL(ctx, "nonexistent:key")

		// Assert
		assert.NoError(t, err)
		// Redis 对不存在的键返回 -2ns
		assert.Equal(t, time.Duration(-2), ttl)
	})

	t.Run("TTL_NoExpiration", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:ttl:noexp"

		// 设置不带过期时间的缓存
		_ = service.Set(ctx, testKey, "value", 0)

		// Act
		ttl, err := service.TTL(ctx, testKey)

		// Assert
		assert.NoError(t, err)
		// Redis 对永不过期的键返回 -1ns
		assert.Equal(t, time.Duration(-1), ttl)
	})
}

// TestRedisCacheService_Increment 测试递增
func TestRedisCacheService_Increment(t *testing.T) {
	t.Run("Increment_Success", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:increment"

		// 设置初始值
		_ = service.Set(ctx, testKey, "10", 5*time.Minute)

		// Act
		result, err := service.Increment(ctx, testKey)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(11), result)

		// 验证值已更新
		val, _ := service.Get(ctx, testKey)
		assert.Equal(t, "11", val)
	})

	t.Run("Increment_NewKey", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// Act - 对不存在的键递增
		result, err := service.Increment(ctx, "new:increment:key")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(1), result)
	})
}

// TestRedisCacheService_Decrement 测试递减
func TestRedisCacheService_Decrement(t *testing.T) {
	t.Run("Decrement_Success", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:decrement"

		// 设置初始值
		_ = service.Set(ctx, testKey, "10", 5*time.Minute)

		// Act
		result, err := service.Decrement(ctx, testKey)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(9), result)

		// 验证值已更新
		val, _ := service.Get(ctx, testKey)
		assert.Equal(t, "9", val)
	})

	t.Run("Decrement_NewKey", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// Act - 对不存在的键递减
		result, err := service.Decrement(ctx, "new:decrement:key")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(-1), result)
	})
}

// ============ 哈希操作测试 ============

// TestRedisCacheService_HGet_HSet 测试哈希字段读写
func TestRedisCacheService_HGet_HSet(t *testing.T) {
	t.Run("HGet_HSet_Success", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:hash"
		testField := "field1"
		testValue := "value1"

		// Act
		err := service.HSet(ctx, testKey, testField, testValue)

		// Assert
		assert.NoError(t, err)

		// 验证可以读取
		result, err := service.HGet(ctx, testKey, testField)
		assert.NoError(t, err)
		assert.Equal(t, testValue, result)
	})

	t.Run("HGet_NotFound", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// Act
		result, err := service.HGet(ctx, "test:hash", "nonexistent:field")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "field not found")
	})
}

// TestRedisCacheService_HGetAll 测试获取所有哈希字段
func TestRedisCacheService_HGetAll(t *testing.T) {
	t.Run("HGetAll_Success", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:getall"

		// 设置多个字段
		_ = service.HSet(ctx, testKey, "field1", "value1")
		_ = service.HSet(ctx, testKey, "field2", "value2")
		_ = service.HSet(ctx, testKey, "field3", "value3")

		// Act
		result, err := service.HGetAll(ctx, testKey)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, result, 3)
		assert.Equal(t, "value1", result["field1"])
		assert.Equal(t, "value2", result["field2"])
		assert.Equal(t, "value3", result["field3"])
	})

	t.Run("HGetAll_Empty", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// Act
		result, err := service.HGetAll(ctx, "nonexistent:hash")

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, result)
	})
}

// TestRedisCacheService_HDelete 测试删除哈希字段
func TestRedisCacheService_HDelete(t *testing.T) {
	t.Run("HDelete_Success", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:hdelete"

		// 设置多个字段
		_ = service.HSet(ctx, testKey, "field1", "value1")
		_ = service.HSet(ctx, testKey, "field2", "value2")
		_ = service.HSet(ctx, testKey, "field3", "value3")

		// Act
		err := service.HDelete(ctx, testKey, "field1", "field2")

		// Assert
		assert.NoError(t, err)

		// 验证字段已删除
		allFields, _ := service.HGetAll(ctx, testKey)
		assert.Len(t, allFields, 1)
		assert.Contains(t, allFields, "field3")
	})

	t.Run("HDelete_Empty", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// Act
		err := service.HDelete(ctx, "test:hdelete")

		// Assert
		assert.NoError(t, err)
	})
}

// ============ 集合操作测试 ============

// TestRedisCacheService_SAdd_SMembers 测试集合操作
func TestRedisCacheService_SAdd_SMembers(t *testing.T) {
	t.Run("SAdd_SMembers_Success", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:set"

		// Act
		err := service.SAdd(ctx, testKey, "member1", "member2", "member3")

		// Assert
		assert.NoError(t, err)

		// 验证可以获取所有成员
		members, err := service.SMembers(ctx, testKey)
		assert.NoError(t, err)
		assert.Len(t, members, 3)
		assert.Contains(t, members, "member1")
		assert.Contains(t, members, "member2")
		assert.Contains(t, members, "member3")
	})

	t.Run("SAdd_Duplicate", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:set:dup"

		// 添加成员
		_ = service.SAdd(ctx, testKey, "member1", "member2")

		// Act - 添加重复成员
		err := service.SAdd(ctx, testKey, "member2", "member3")

		// Assert
		assert.NoError(t, err)

		// 验证集合去重
		members, _ := service.SMembers(ctx, testKey)
		assert.Len(t, members, 3)
	})

	t.Run("SAdd_Empty", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// Act
		err := service.SAdd(ctx, "test:set")

		// Assert
		assert.NoError(t, err)
	})
}

// TestRedisCacheService_SIsMember 测试检查集合成员
func TestRedisCacheService_SIsMember(t *testing.T) {
	t.Run("SIsMember_True", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:set:member"

		// 添加成员
		_ = service.SAdd(ctx, testKey, "member1", "member2")

		// Act
		exists, err := service.SIsMember(ctx, testKey, "member1")

		// Assert
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("SIsMember_False", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		// Act
		exists, err := service.SIsMember(ctx, "test:set:member", "nonexistent")

		// Assert
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

// TestRedisCacheService_SRemove 测试移除集合成员
func TestRedisCacheService_SRemove(t *testing.T) {
	t.Run("SRemove_Success", func(t *testing.T) {
		// Arrange
		miniRedis := miniredis.RunT(t)
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})
		service := NewRedisCacheService(redisClient)
		ctx := context.Background()

		testKey := "test:set:remove"

		// 添加成员
		_ = service.SAdd(ctx, testKey, "member1", "member2", "member3")

		// Act
		err := service.SRemove(ctx, testKey, "member1", "member2")

		// Assert
		assert.NoError(t, err)

		// 验证成员已移除
		members, _ := service.SMembers(ctx, testKey)
		assert.Len(t, members, 1)
		assert.Contains(t, members, "member3")
	})
}

// Benchmark 基准测试

// BenchmarkRedisCacheService_MGet 批量获取性能测试
func BenchmarkRedisCacheService_MGet(b *testing.B) {
	miniRedis := miniredis.RunT(b)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	service := NewRedisCacheService(redisClient)
	ctx := context.Background()

	// 准备测试数据
	keys := make([]string, 100)
	for i := 0; i < 100; i++ {
		keys[i] = "bench:keys:" + string(rune('0'+i%10))
		_ = service.Set(ctx, keys[i], "value", 5*time.Minute)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.MGet(ctx, keys...)
	}
}

// BenchmarkRedisCacheService_MSet 批量设置性能测试
func BenchmarkRedisCacheService_MSet(b *testing.B) {
	miniRedis := miniredis.RunT(b)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	service := NewRedisCacheService(redisClient)
	ctx := context.Background()

	kvPairs := make(map[string]string)
	for i := 0; i < 100; i++ {
		key := "bench:mset:" + string(rune('0'+i%10))
		kvPairs[key] = "value"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.MSet(ctx, kvPairs, 5*time.Minute)
	}
}
