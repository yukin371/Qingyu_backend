package shared

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/service/shared/messaging"
)

// 注意：这些测试需要真实的Redis实例
// 可以通过设置环境变量 REDIS_ADDR 来指定Redis地址
// 默认使用 localhost:6379

func getTestRedisClient(t *testing.T) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15, // 使用DB 15 进行测试，避免影响生产数据
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	return client
}

func cleanupTestStream(t *testing.T, client *redis.Client, stream string) {
	ctx := context.Background()
	client.Del(ctx, stream)
}

// TestRedisQueueClientPublish 测试发布消息
func TestRedisQueueClientPublish(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	redisClient := getTestRedisClient(t)
	defer redisClient.Close()

	queueClient := messaging.NewRedisQueueClient(redisClient)
	ctx := context.Background()
	stream := "test_stream_publish"
	defer cleanupTestStream(t, redisClient, stream)

	t.Run("成功发布消息", func(t *testing.T) {
		data := map[string]interface{}{
			"topic":   "test",
			"message": "Hello World",
		}

		messageID, err := queueClient.Publish(ctx, stream, data)

		assert.NoError(t, err)
		assert.NotEmpty(t, messageID)

		// 验证消息是否真的被添加到Stream
		length, err := redisClient.XLen(ctx, stream).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(1), length)
	})

	t.Run("发布多条消息", func(t *testing.T) {
		cleanupTestStream(t, redisClient, stream)

		for i := 0; i < 5; i++ {
			data := map[string]interface{}{
				"index": i,
				"data":  "message",
			}
			_, err := queueClient.Publish(ctx, stream, data)
			assert.NoError(t, err)
		}

		length, err := redisClient.XLen(ctx, stream).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(5), length)
	})
}

// TestRedisQueueClientSubscribe 测试订阅消息
func TestRedisQueueClientSubscribe(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	redisClient := getTestRedisClient(t)
	defer redisClient.Close()

	queueClient := messaging.NewRedisQueueClient(redisClient)
	ctx := context.Background()
	stream := "test_stream_subscribe"
	defer cleanupTestStream(t, redisClient, stream)

	t.Run("创建消费者组", func(t *testing.T) {
		err := queueClient.CreateGroup(ctx, stream, "test_group")
		assert.NoError(t, err)

		// 再次创建应该不报错（幂等）
		err = queueClient.CreateGroup(ctx, stream, "test_group")
		assert.NoError(t, err)
	})

	t.Run("读取消息", func(t *testing.T) {
		cleanupTestStream(t, redisClient, stream)

		// 创建组
		err := queueClient.CreateGroup(ctx, stream, "read_group")
		require.NoError(t, err)

		// 发布消息
		data := map[string]interface{}{
			"topic": "test",
			"data":  "read test",
		}
		_, err = queueClient.Publish(ctx, stream, data)
		require.NoError(t, err)

		// 读取消息
		messages, err := queueClient.Subscribe(ctx, stream, "read_group", "consumer1", 10)
		assert.NoError(t, err)
		assert.Len(t, messages, 1)
		assert.Equal(t, "test", messages[0].Data["topic"])
	})

	t.Run("消息确认", func(t *testing.T) {
		cleanupTestStream(t, redisClient, stream)

		// 创建组
		err := queueClient.CreateGroup(ctx, stream, "ack_group")
		require.NoError(t, err)

		// 发布消息
		data := map[string]interface{}{
			"topic": "ack_test",
		}
		_, err = queueClient.Publish(ctx, stream, data)
		require.NoError(t, err)

		// 读取消息
		messages, err := queueClient.Subscribe(ctx, stream, "ack_group", "consumer1", 10)
		require.NoError(t, err)
		require.Len(t, messages, 1)

		// 确认消息
		err = queueClient.Ack(ctx, stream, "ack_group", messages[0].ID)
		assert.NoError(t, err)
	})

	t.Run("没有新消息时返回空列表", func(t *testing.T) {
		cleanupTestStream(t, redisClient, stream)

		err := queueClient.CreateGroup(ctx, stream, "empty_group")
		require.NoError(t, err)

		messages, err := queueClient.Subscribe(ctx, stream, "empty_group", "consumer1", 10)
		assert.NoError(t, err)
		assert.Empty(t, messages)
	})
}

// TestRedisQueueClientStreamManagement 测试流管理
func TestRedisQueueClientStreamManagement(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	redisClient := getTestRedisClient(t)
	defer redisClient.Close()

	queueClient := messaging.NewRedisQueueClient(redisClient)
	ctx := context.Background()

	t.Run("删除流", func(t *testing.T) {
		stream := "test_stream_delete"

		// 创建流
		data := map[string]interface{}{"test": "data"}
		_, err := queueClient.Publish(ctx, stream, data)
		require.NoError(t, err)

		// 删除流
		err = queueClient.DeleteStream(ctx, stream)
		assert.NoError(t, err)

		// 验证流已删除
		length, err := redisClient.XLen(ctx, stream).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), length)
	})

	t.Run("列出所有流", func(t *testing.T) {
		// 清理
		redisClient.FlushDB(ctx)

		// 创建几个测试流
		streams := []string{"stream_a", "stream_b", "stream_c"}
		for _, stream := range streams {
			data := map[string]interface{}{"test": "data"}
			_, err := queueClient.Publish(ctx, stream, data)
			require.NoError(t, err)
		}

		// 列出所有流
		foundStreams, err := queueClient.ListStreams(ctx)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(foundStreams), 3)

		// 清理
		for _, stream := range streams {
			queueClient.DeleteStream(ctx, stream)
		}
	})
}

// TestRedisQueueClientConcurrency 测试并发场景
func TestRedisQueueClientConcurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	redisClient := getTestRedisClient(t)
	defer redisClient.Close()

	queueClient := messaging.NewRedisQueueClient(redisClient)
	ctx := context.Background()
	stream := "test_stream_concurrent"
	defer cleanupTestStream(t, redisClient, stream)

	t.Run("并发发布消息", func(t *testing.T) {
		cleanupTestStream(t, redisClient, stream)

		const numMessages = 100
		errors := make(chan error, numMessages)

		for i := 0; i < numMessages; i++ {
			go func(index int) {
				data := map[string]interface{}{
					"index": index,
				}
				_, err := queueClient.Publish(ctx, stream, data)
				errors <- err
			}(i)
		}

		// 收集错误
		for i := 0; i < numMessages; i++ {
			err := <-errors
			assert.NoError(t, err)
		}

		// 验证所有消息都已发布
		length, err := redisClient.XLen(ctx, stream).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(numMessages), length)
	})

	t.Run("多消费者并发读取", func(t *testing.T) {
		cleanupTestStream(t, redisClient, stream)

		// 创建消费者组
		err := queueClient.CreateGroup(ctx, stream, "multi_consumer_group")
		require.NoError(t, err)

		// 发布一些消息
		for i := 0; i < 10; i++ {
			data := map[string]interface{}{
				"message": i,
			}
			_, err := queueClient.Publish(ctx, stream, data)
			require.NoError(t, err)
		}

		// 多个消费者并发读取
		const numConsumers = 3
		messagesChan := make(chan int, numConsumers)

		for i := 0; i < numConsumers; i++ {
			go func(consumerID int) {
				consumerName := "consumer_" + string(rune(consumerID))
				messages, err := queueClient.Subscribe(ctx, stream, "multi_consumer_group", consumerName, 10)
				if err == nil {
					messagesChan <- len(messages)
				} else {
					messagesChan <- 0
				}
			}(i)
		}

		// 收集结果
		totalRead := 0
		for i := 0; i < numConsumers; i++ {
			count := <-messagesChan
			totalRead += count
		}

		// 至少应该读取到一些消息（具体数量取决于消费者的读取顺序）
		assert.Greater(t, totalRead, 0)
	})
}

// BenchmarkRedisQueueClientPublish 基准测试发布性能
func BenchmarkRedisQueueClientPublish(b *testing.B) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15,
	})
	defer redisClient.Close()

	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		b.Skip("Redis not available")
	}

	queueClient := messaging.NewRedisQueueClient(redisClient)
	stream := "benchmark_stream"
	defer redisClient.Del(ctx, stream)

	data := map[string]interface{}{
		"topic": "benchmark",
		"data":  "test data",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = queueClient.Publish(ctx, stream, data)
	}
}

// BenchmarkRedisQueueClientSubscribe 基准测试订阅性能
func BenchmarkRedisQueueClientSubscribe(b *testing.B) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15,
	})
	defer redisClient.Close()

	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		b.Skip("Redis not available")
	}

	queueClient := messaging.NewRedisQueueClient(redisClient)
	stream := "benchmark_stream_sub"
	defer redisClient.Del(ctx, stream)

	// 准备数据
	queueClient.CreateGroup(ctx, stream, "bench_group")
	for i := 0; i < 100; i++ {
		data := map[string]interface{}{"index": i}
		queueClient.Publish(ctx, stream, data)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = queueClient.Subscribe(ctx, stream, "bench_group", "consumer1", 10)
	}
}
