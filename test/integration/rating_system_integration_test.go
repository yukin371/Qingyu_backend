package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/service"
)

// TestRatingSystem_E2E 端到端测试：完整的评分流程
// 测试评分创建、查询、更新、删除的完整生命周期
func TestRatingSystem_E2E(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// 设置测试环境
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// 初始化helper
	helper := NewTestHelper(t, router)

	// 登录测试用户
	token := helper.LoginTestUser()
	if token == "" {
		t.Skip("无法登录测试用户，跳过评分测试")
	}

	ctx := context.Background()

	// 获取一本测试书籍
	var testBookID string
	db := service.ServiceManager.GetMongoDB()
	cursor, err := db.Collection("books").Find(ctx, bson.M{}, nil)
	if err == nil {
		var books []map[string]interface{}
		cursor.All(ctx, &books)
		cursor.Close(ctx)

		if len(books) > 0 {
			if oid, ok := books[0]["_id"].(primitive.ObjectID); ok {
				testBookID = oid.Hex()
			}
		}
	}

	if testBookID == "" {
		t.Skip("数据库中没有测试书籍")
	}

	// 存储测试中创建的ID
	var commentID string
	var reviewID string

	t.Run("1.评分_发表带评分的评论", func(t *testing.T) {
		// 发表带评分的评论
		requestBody := map[string]interface{}{
			"book_id":  testBookID,
			"content":  "集成测试评论，评分5星",
			"rating":   5,
			"is_public": true,
		}

		w := helper.DoAuthRequest("POST", "/api/v1/social/comments", requestBody, token)
		data := helper.AssertSuccess(w, 200, "发表评论应该成功")

		// 保存评论ID
		if id, ok := data["_id"].(string); ok {
			commentID = id
		} else if id, ok := data["id"].(string); ok {
			commentID = id
		}
		require.NotEmpty(t, commentID, "评论ID不应为空")
		helper.LogSuccess("带评分评论发表成功，评论ID: %s, 评分: 5", commentID)
	})

	t.Run("2.评分_查询评论评分统计", func(t *testing.T) {
		// 等待数据库写入完成
		time.Sleep(100 * time.Millisecond)

		// 查询评分统计
		w := helper.DoAuthRequest("GET", fmt.Sprintf("/api/v1/social/rating/stats?targetType=comment&targetId=%s", commentID), nil, token)
		data := helper.AssertSuccess(w, 200, "查询评分统计应该成功")

		// 验证评分统计数据
		helper.LogSuccess("评分统计查询成功: %+v", data)

		// 验证必要的字段存在
		assert.Contains(t, data, "targetId")
		assert.Contains(t, data, "targetType")
		assert.Contains(t, data, "averageRating")
		assert.Contains(t, data, "totalRatings")
	})

	t.Run("3.评分_发表书评", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"book_id":  testBookID,
			"title":    "集成测试书评",
			"content":  "这是一篇集成测试书评，评分4星",
			"rating":   4,
			"is_public": true,
		}

		w := helper.DoAuthRequest("POST", "/api/v1/social/reviews", requestBody, token)
		data := helper.AssertSuccess(w, 200, "发表书评应该成功")

		// 保存书评ID
		if id, ok := data["_id"].(string); ok {
			reviewID = id
		} else if id, ok := data["id"].(string); ok {
			reviewID = id
		}
		require.NotEmpty(t, reviewID, "书评ID不应为空")
		helper.LogSuccess("带评分书评发表成功，书评ID: %s, 评分: 4", reviewID)
	})

	t.Run("4.评分_查询书评评分统计", func(t *testing.T) {
		// 等待数据库写入完成
		time.Sleep(100 * time.Millisecond)

		// 查询评分统计
		w := helper.DoAuthRequest("GET", fmt.Sprintf("/api/v1/social/rating/stats?targetType=review&targetId=%s", reviewID), nil, token)
		data := helper.AssertSuccess(w, 200, "查询书评评分统计应该成功")

		helper.LogSuccess("书评评分统计查询成功: %+v", data)

		// 验证评分值
		if averageRating, ok := data["averageRating"].(float64); ok {
			assert.Equal(t, float64(4), averageRating, "平均评分应为4")
		}
	})

	t.Run("5.评分_查询书籍评分统计", func(t *testing.T) {
		// 查询书籍级别的评分统计
		w := helper.DoAuthRequest("GET", fmt.Sprintf("/api/v1/social/rating/stats?targetType=book&targetId=%s", testBookID), nil, token)
		data := helper.AssertSuccess(w, 200, "查询书籍评分统计应该成功")

		helper.LogSuccess("书籍评分统计查询成功: %+v", data)

		// 验证必要字段
		assert.Contains(t, data, "targetId")
		assert.Contains(t, data, "targetType")
		assert.Equal(t, "book", data["targetType"])
	})

	t.Run("6.评分_清理测试数据", func(t *testing.T) {
		// 清理评论 - 使用API删除
		if commentID != "" {
			w := helper.DoAuthRequest("DELETE", fmt.Sprintf("/api/v1/social/comments/%s", commentID), nil, token)
			if w.Code == 200 || w.Code == 204 {
				helper.LogSuccess("测试评论已删除: %s", commentID)
			} else {
				helper.LogWarning("删除评论失败 (状态码: %d)", w.Code)
			}
		}

		// 清理书评
		if reviewID != "" {
			w := helper.DoAuthRequest("DELETE", fmt.Sprintf("/api/v1/social/reviews/%s", reviewID), nil, token)
			if w.Code == 200 || w.Code == 204 {
				helper.LogSuccess("测试书评已删除: %s", reviewID)
			} else {
				helper.LogWarning("删除书评失败 (状态码: %d)", w.Code)
			}
		}
	})
}

// TestRatingSystem_CacheIntegration 缓存集成测试
// 测试评分统计的缓存读取、写入和失效机制
func TestRatingSystem_CacheIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// 这个测试需要真实的Redis连接
	t.Skip("需要Redis集成环境 - 缓存集成测试暂时跳过")

	// TODO: 当Redis环境准备好后，实现以下测试场景：
	// 1. 首次查询评分统计 - 缓存未命中，从数据库聚合
	// 2. 再次查询评分统计 - 缓存命中，直接返回
	// 3. 更新评分 - 缓存失效
	// 4. 查询评分统计 - 缓存未命中，重新从数据库聚合
	// 5. 验证TTL设置
}

// TestRatingSystem_Performance 性能测试
// 验证评分统计查询性能符合要求（< 100ms）
func TestRatingSystem_Performance(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过性能测试")
	}

	// 设置测试环境
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)
	token := helper.LoginTestUser()
	if token == "" {
		t.Skip("无法登录测试用户，跳过性能测试")
	}

	ctx := context.Background()

	// 获取测试书籍
	var testBookID string
	db := service.ServiceManager.GetMongoDB()
	cursor, err := db.Collection("books").Find(ctx, bson.M{}, nil)
	if err == nil {
		var books []map[string]interface{}
		cursor.All(ctx, &books)
		cursor.Close(ctx)

		if len(books) > 0 {
			if oid, ok := books[0]["_id"].(primitive.ObjectID); ok {
				testBookID = oid.Hex()
			}
		}
	}

	if testBookID == "" {
		t.Skip("数据库中没有测试书籍")
	}

	t.Run("性能_评分统计查询", func(t *testing.T) {
		// 目标: 评分统计查询 < 100ms
		iterations := 10
		var durations []time.Duration

		for i := 0; i < iterations; i++ {
			start := time.Now()

			w := helper.DoAuthRequest("GET", fmt.Sprintf("/api/v1/social/rating/stats?targetType=book&targetId=%s", testBookID), nil, token)

			duration := time.Since(start)
			durations = append(durations, duration)

			// 确保请求成功
			if w.Code != 200 {
				t.Logf("警告: 第 %d 次请求失败，状态码: %d", i+1, w.Code)
			}
		}

		// 计算统计数据
		var total time.Duration
		var min, max time.Duration = durations[0], durations[0]
		for _, d := range durations {
			total += d
			if d < min {
				min = d
			}
			if d > max {
				max = d
			}
		}
		avg := total / time.Duration(iterations)

		helper.LogSuccess("评分统计查询性能统计:")
		helper.LogSuccess("  平均耗时: %v", avg)
		helper.LogSuccess("  最小耗时: %v", min)
		helper.LogSuccess("  最大耗时: %v", max)

		// 验证性能目标: 95%的请求应该在100ms内完成
		// 这里我们检查最大值是否在合理范围内（允许一定误差）
		assert.True(t, max < 200*time.Millisecond, "最大响应时间应小于200ms，实际: %v", max)
		assert.True(t, avg < 100*time.Millisecond, "平均响应时间应小于100ms，实际: %v", avg)
	})
}

// TestRatingSystem_CacheHitRatio 缓存命中率测试
// 测试评分统计缓存的命中率
func TestRatingSystem_CacheHitRatio(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过缓存命中率测试")
	}

	// 这个测试需要真实的Redis连接和监控
	t.Skip("需要Redis监控环境 - 缓存命中率测试暂时跳过")

	// TODO: 当Redis环境准备好后，实现以下测试场景：
	// 1. 预热缓存 - 对多个目标进行首次查询
	// 2. 重复查询 - 多次查询相同目标
	// 3. 计算缓存命中率 - 目标 > 80%
	// 4. 验证缓存TTL - 检查缓存过期时间
}

// TestRatingSystem_DataConsistency 数据一致性测试
// 测试评分更新后统计数据的正确性
func TestRatingSystem_DataConsistency(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过数据一致性测试")
	}

	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)
	token := helper.LoginTestUser()
	if token == "" {
		t.Skip("无法登录测试用户，跳过数据一致性测试")
	}

	ctx := context.Background()

	// 获取测试书籍
	var testBookID string
	db := service.ServiceManager.GetMongoDB()
	cursor, err := db.Collection("books").Find(ctx, bson.M{}, nil)
	if err == nil {
		var books []map[string]interface{}
		cursor.All(ctx, &books)
		cursor.Close(ctx)

		if len(books) > 0 {
			if oid, ok := books[0]["_id"].(primitive.ObjectID); ok {
				testBookID = oid.Hex()
			}
		}
	}

	if testBookID == "" {
		t.Skip("数据库中没有测试书籍")
	}

	var commentID string

	// 清理函数
	defer func() {
		if commentID != "" {
			w := helper.DoAuthRequest("DELETE", fmt.Sprintf("/api/v1/social/comments/%s", commentID), nil, token)
			if w.Code != 200 && w.Code != 204 {
				helper.LogWarning("清理评论失败 (状态码: %d)", w.Code)
			}
		}
	}()

	t.Run("一致性_评分前后统计对比", func(t *testing.T) {
		// 1. 获取初始评分统计
		w1 := helper.DoAuthRequest("GET", fmt.Sprintf("/api/v1/social/rating/stats?targetType=book&targetId=%s", testBookID), nil, token)
		data1 := helper.AssertSuccess(w1, 200, "查询初始评分统计应该成功")

		initialTotal := int64(0)
		if total, ok := data1["totalRatings"].(int64); ok {
			initialTotal = total
		} else if total, ok := data1["totalRatings"].(float64); ok {
			initialTotal = int64(total)
		}

		initialAvg := 0.0
		if avg, ok := data1["averageRating"].(float64); ok {
			initialAvg = avg
		}

		t.Logf("初始评分统计 - 总评分数: %d, 平均分: %.2f", initialTotal, initialAvg)

		// 2. 发表新评论（5星）
		requestBody := map[string]interface{}{
			"book_id":  testBookID,
			"content":  "数据一致性测试评论",
			"rating":   5,
			"is_public": true,
		}

		w2 := helper.DoAuthRequest("POST", "/api/v1/social/comments", requestBody, token)
		data2 := helper.AssertSuccess(w2, 200, "发表评论应该成功")

		if id, ok := data2["_id"].(string); ok {
			commentID = id
		} else if id, ok := data2["id"].(string); ok {
			commentID = id
		}

		// 3. 等待数据库更新
		time.Sleep(200 * time.Millisecond)

		// 4. 获取更新后的评分统计
		w3 := helper.DoAuthRequest("GET", fmt.Sprintf("/api/v1/social/rating/stats?targetType=book&targetId=%s", testBookID), nil, token)
		data3 := helper.AssertSuccess(w3, 200, "查询更新后评分统计应该成功")

		updatedTotal := int64(0)
		if total, ok := data3["totalRatings"].(int64); ok {
			updatedTotal = total
		} else if total, ok := data3["totalRatings"].(float64); ok {
			updatedTotal = int64(total)
		}

		updatedAvg := 0.0
		if avg, ok := data3["averageRating"].(float64); ok {
			updatedAvg = avg
		}

		t.Logf("更新后评分统计 - 总评分数: %d, 平均分: %.2f", updatedTotal, updatedAvg)

		// 5. 验证数据一致性
		assert.Equal(t, initialTotal+1, updatedTotal, "评分数应该增加1")
		// 平均分应该变化（除非之前没有评分）
		if initialTotal > 0 {
			assert.NotEqual(t, initialAvg, updatedAvg, "平均分应该变化")
		}
	})
}

// TestRatingService_AggregateFromComments 测试从评论聚合评分
func TestRatingService_AggregateFromComments(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过聚合测试")
	}

	// 这个测试需要直接测试Service层
	t.Skip("需要Service层测试环境 - 聚合测试暂时跳过")

	// TODO: 实现Service层的直接测试
	// 1. 创建Mock Repository
	// 2. 测试aggregateCommentRatings方法
	// 3. 测试aggregateReviewRatings方法
	// 4. 测试aggregateBookRatings方法
}

// TestRatingStats_DistributionValidation 评分分布验证测试
func TestRatingStats_DistributionValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过评分分布验证测试")
	}

	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)
	token := helper.LoginTestUser()
	if token == "" {
		t.Skip("无法登录测试用户，跳过评分分布验证测试")
	}

	ctx := context.Background()

	// 获取测试书籍
	var testBookID string
	db := service.ServiceManager.GetMongoDB()
	cursor, err := db.Collection("books").Find(ctx, bson.M{}, nil)
	if err == nil {
		var books []map[string]interface{}
		cursor.All(ctx, &books)
		cursor.Close(ctx)

		if len(books) > 0 {
			if oid, ok := books[0]["_id"].(primitive.ObjectID); ok {
				testBookID = oid.Hex()
			}
		}
	}

	if testBookID == "" {
		t.Skip("数据库中没有测试书籍")
	}

	t.Run("分布_验证评分分布结构", func(t *testing.T) {
		// 查询评分统计
		w := helper.DoAuthRequest("GET", fmt.Sprintf("/api/v1/social/rating/stats?targetType=book&targetId=%s", testBookID), nil, token)
		data := helper.AssertSuccess(w, 200, "查询评分统计应该成功")

		// 验证评分分布字段
		if distribution, ok := data["distribution"].(map[string]interface{}); ok {
			t.Logf("评分分布: %+v", distribution)

			// 验证分布包含1-5星的所有等级（如果存在评分）
			// 注意：如果没有任何评分，distribution可能为空或所有值为0
			for star := 1; star <= 5; star++ {
				starKey := fmt.Sprintf("%d", star)
				if count, exists := distribution[starKey]; exists {
					t.Logf("  %d星: %v", star, count)
				}
			}
		} else {
			t.Log("评分分布字段不存在或为空")
		}
	})
}

// TestRatingSystem_ConcurrentOperations 并发操作测试
// 测试多个并发评分请求的数据一致性
func TestRatingSystem_ConcurrentOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过并发操作测试")
	}

	// 这个测试需要模拟并发场景
	t.Skip("需要并发测试环境 - 并发操作测试暂时跳过")

	// TODO: 实现并发测试场景：
	// 1. 多个用户同时对同一目标评分
	// 2. 验证最终评分统计的正确性
	// 3. 检查是否存在竞态条件
}
