package integration

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/config"
	"Qingyu_backend/pkg/cache"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRedisIntegration Redis集成测试
// 需要本地Redis服务运行
func TestRedisIntegration(t *testing.T) {
	// 跳过CI环境
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// 创建Redis客户端
	cfg := config.DefaultRedisConfig()
	client, err := cache.NewRedisClient(cfg)
	require.NoError(t, err, "创建Redis客户端失败")
	defer client.Close()

	ctx := context.Background()

	t.Run("健康检查", func(t *testing.T) {
		err := client.Ping(ctx)
		assert.NoError(t, err)
	})

	t.Run("基本读写", func(t *testing.T) {
		key := "test:integration:basic"
		value := "test_value_" + time.Now().Format("20060102150405")

		// 写入
		err := client.Set(ctx, key, value, 1*time.Minute)
		assert.NoError(t, err)

		// 读取
		result, err := client.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, result)

		// 删除
		err = client.Delete(ctx, key)
		assert.NoError(t, err)

		// 验证删除
		_, err = client.Get(ctx, key)
		assert.Error(t, err) // 应该返回错误
	})

	t.Run("过期时间", func(t *testing.T) {
		key := "test:integration:expire"
		value := "expire_value"

		// 设置短过期时间
		err := client.Set(ctx, key, value, 2*time.Second)
		assert.NoError(t, err)

		// 立即读取应该成功
		result, err := client.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, result)

		// 等待过期
		time.Sleep(3 * time.Second)

		// 读取应该失败
		_, err = client.Get(ctx, key)
		assert.Error(t, err)
	})

	t.Run("Hash操作", func(t *testing.T) {
		key := "test:integration:hash"

		// 设置Hash字段
		err := client.HSet(ctx, key, "field1", "value1", "field2", "value2")
		assert.NoError(t, err)

		// 获取单个字段
		val, err := client.HGet(ctx, key, "field1")
		assert.NoError(t, err)
		assert.Equal(t, "value1", val)

		// 获取所有字段
		all, err := client.HGetAll(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, "value1", all["field1"])
		assert.Equal(t, "value2", all["field2"])

		// 删除字段
		err = client.HDel(ctx, key, "field1")
		assert.NoError(t, err)

		// 验证删除
		all, err = client.HGetAll(ctx, key)
		assert.NoError(t, err)
		assert.NotContains(t, all, "field1")
		assert.Contains(t, all, "field2")

		// 清理
		err = client.Delete(ctx, key)
		assert.NoError(t, err)
	})

	t.Run("Set集合操作", func(t *testing.T) {
		key := "test:integration:set"

		// 添加成员
		err := client.SAdd(ctx, key, "member1", "member2", "member3")
		assert.NoError(t, err)

		// 获取所有成员
		members, err := client.SMembers(ctx, key)
		assert.NoError(t, err)
		assert.Len(t, members, 3)
		assert.Contains(t, members, "member1")

		// 删除成员
		err = client.SRem(ctx, key, "member1")
		assert.NoError(t, err)

		// 验证删除
		members, err = client.SMembers(ctx, key)
		assert.NoError(t, err)
		assert.Len(t, members, 2)
		assert.NotContains(t, members, "member1")

		// 清理
		err = client.Delete(ctx, key)
		assert.NoError(t, err)
	})

	t.Run("原子操作", func(t *testing.T) {
		key := "test:integration:counter"

		// 自增
		val, err := client.Incr(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), val)

		// 增加指定值
		val, err = client.IncrBy(ctx, key, 5)
		assert.NoError(t, err)
		assert.Equal(t, int64(6), val)

		// 自减
		val, err = client.Decr(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), val)

		// 减少指定值
		val, err = client.DecrBy(ctx, key, 3)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), val)

		// 清理
		err = client.Delete(ctx, key)
		assert.NoError(t, err)
	})

	t.Run("批量操作", func(t *testing.T) {
		keys := []string{
			"test:integration:batch1",
			"test:integration:batch2",
			"test:integration:batch3",
		}

		// 批量设置
		err := client.MSet(ctx, keys[0], "value1", keys[1], "value2", keys[2], "value3")
		assert.NoError(t, err)

		// 批量获取
		values, err := client.MGet(ctx, keys...)
		assert.NoError(t, err)
		assert.Len(t, values, 3)

		// 批量删除
		err = client.Delete(ctx, keys...)
		assert.NoError(t, err)
	})

	t.Run("TTL操作", func(t *testing.T) {
		key := "test:integration:ttl"
		value := "ttl_value"

		// 设置带过期时间的键
		err := client.Set(ctx, key, value, 10*time.Minute)
		assert.NoError(t, err)

		// 获取TTL
		ttl, err := client.TTL(ctx, key)
		assert.NoError(t, err)
		assert.True(t, ttl > 0 && ttl <= 10*time.Minute)

		// 更新过期时间
		err = client.Expire(ctx, key, 5*time.Minute)
		assert.NoError(t, err)

		// 验证新的TTL
		ttl, err = client.TTL(ctx, key)
		assert.NoError(t, err)
		assert.True(t, ttl > 0 && ttl <= 5*time.Minute)

		// 清理
		err = client.Delete(ctx, key)
		assert.NoError(t, err)
	})

	t.Run("Exists检查", func(t *testing.T) {
		key := "test:integration:exists"

		// 不存在的键
		count, err := client.Exists(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)

		// 创建键
		err = client.Set(ctx, key, "value", 1*time.Minute)
		assert.NoError(t, err)

		// 存在的键
		count, err = client.Exists(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// 清理
		err = client.Delete(ctx, key)
		assert.NoError(t, err)
	})
}

// TestRedisServiceContainerIntegration 测试Redis与服务容器集成
func TestRedisServiceContainerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// 加载配置
	cfg := &config.Config{
		Redis: config.DefaultRedisConfig(),
	}
	config.GlobalConfig = cfg

	// 注意：这里需要完整的服务容器初始化流程
	// 实际使用中，服务容器会在应用启动时初始化
	t.Run("配置加载", func(t *testing.T) {
		redisCfg := config.GetRedisConfig()
		assert.NotNil(t, redisCfg)
		assert.Equal(t, "localhost", redisCfg.Host)
		assert.Equal(t, 6379, redisCfg.Port)
	})
}

// BenchmarkRedisOperations Redis操作性能测试
func BenchmarkRedisOperations(b *testing.B) {
	if testing.Short() {
		b.Skip("跳过性能测试")
	}

	cfg := config.DefaultRedisConfig()
	client, err := cache.NewRedisClient(cfg)
	if err != nil {
		b.Fatalf("创建Redis客户端失败: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	b.Run("Set", func(b *testing.B) {
		key := "bench:set"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = client.Set(ctx, key, "value", 1*time.Minute)
		}
	})

	b.Run("Get", func(b *testing.B) {
		key := "bench:get"
		_ = client.Set(ctx, key, "value", 1*time.Minute)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = client.Get(ctx, key)
		}
	})

	b.Run("Incr", func(b *testing.B) {
		key := "bench:incr"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = client.Incr(ctx, key)
		}
	})

	b.Run("HSet", func(b *testing.B) {
		key := "bench:hset"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = client.HSet(ctx, key, "field", "value")
		}
	})
}
