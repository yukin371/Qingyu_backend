package repository

import (
	"Qingyu_backend/models/community"
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/repository/mongodb/reading"
)

// TestLikeRepositoryAdvanced 高级点赞Repository测试
func TestLikeRepositoryAdvanced(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := reading.NewMongoLikeRepository(db)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("GetByID_Success", func(t *testing.T) {
		testBookID := primitive.NewObjectID().Hex()

		// 添加点赞
		like := &community.Like{
			UserID:     testUserID,
			TargetType: community.LikeTargetTypeBook,
			TargetID:   testBookID,
			CreatedAt:  time.Now(),
		}
		err := repo.AddLike(ctx, like)
		assert.NoError(t, err)

		// 通过ID获取
		found, err := repo.GetByID(ctx, like.ID.Hex())
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, testUserID, found.UserID)
		assert.Equal(t, testBookID, found.TargetID)

		t.Logf("✓ 通过ID获取点赞记录成功")
	})

	t.Run("GetByID_NotFound", func(t *testing.T) {
		fakeID := primitive.NewObjectID().Hex()

		found, err := repo.GetByID(ctx, fakeID)
		assert.Error(t, err)
		assert.Nil(t, found)
		assert.Contains(t, err.Error(), "not found")

		t.Logf("✓ 不存在的点赞记录正确返回错误")
	})

	t.Run("CountTargetLikes_Success", func(t *testing.T) {
		testBookID := primitive.NewObjectID().Hex()

		// 添加多个用户的点赞
		for i := 0; i < 5; i++ {
			userID := primitive.NewObjectID().Hex()
			like := &community.Like{
				UserID:     userID,
				TargetType: community.LikeTargetTypeBook,
				TargetID:   testBookID,
				CreatedAt:  time.Now(),
			}
			err := repo.AddLike(ctx, like)
			assert.NoError(t, err)
		}

		// 统计点赞数
		count, err := repo.CountTargetLikes(ctx, community.LikeTargetTypeBook, testBookID)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)

		t.Logf("✓ 统计目标点赞数成功: %d", count)
	})

	t.Run("GetUserLikes_AllTypes", func(t *testing.T) {
		user2 := primitive.NewObjectID().Hex()

		// 添加不同类型的点赞
		like1 := &community.Like{
			UserID:     user2,
			TargetType: community.LikeTargetTypeBook,
			TargetID:   primitive.NewObjectID().Hex(),
			CreatedAt:  time.Now(),
		}
		err := repo.AddLike(ctx, like1)
		assert.NoError(t, err)

		like2 := &community.Like{
			UserID:     user2,
			TargetType: community.LikeTargetTypeComment,
			TargetID:   primitive.NewObjectID().Hex(),
			CreatedAt:  time.Now(),
		}
		err = repo.AddLike(ctx, like2)
		assert.NoError(t, err)

		// 查询所有类型（不过滤targetType）
		likes, total, err := repo.GetUserLikes(ctx, user2, "", 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Equal(t, 2, len(likes))

		t.Logf("✓ 查询用户所有类型点赞成功")
	})

	t.Run("GetUserLikes_FilterByType", func(t *testing.T) {
		user3 := primitive.NewObjectID().Hex()

		// 添加不同类型的点赞
		for i := 0; i < 3; i++ {
			like := &community.Like{
				UserID:     user3,
				TargetType: community.LikeTargetTypeBook,
				TargetID:   primitive.NewObjectID().Hex(),
				CreatedAt:  time.Now(),
			}
			err := repo.AddLike(ctx, like)
			assert.NoError(t, err)
		}

		for i := 0; i < 2; i++ {
			like := &community.Like{
				UserID:     user3,
				TargetType: community.LikeTargetTypeComment,
				TargetID:   primitive.NewObjectID().Hex(),
				CreatedAt:  time.Now(),
			}
			err := repo.AddLike(ctx, like)
			assert.NoError(t, err)
		}

		// 只查询书籍点赞
		likes, total, err := repo.GetUserLikes(ctx, user3, community.LikeTargetTypeBook, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Equal(t, 3, len(likes))

		// 验证所有都是书籍类型
		for _, like := range likes {
			assert.Equal(t, community.LikeTargetTypeBook, like.TargetType)
		}

		t.Logf("✓ 按类型过滤用户点赞成功")
	})

	t.Run("ChapterLike_Success", func(t *testing.T) {
		user4 := primitive.NewObjectID().Hex()
		chapterID := primitive.NewObjectID().Hex()

		// 点赞章节
		like := &community.Like{
			UserID:     user4,
			TargetType: community.LikeTargetTypeChapter,
			TargetID:   chapterID,
			CreatedAt:  time.Now(),
		}
		err := repo.AddLike(ctx, like)
		assert.NoError(t, err)

		// 验证点赞状态
		isLiked, err := repo.IsLiked(ctx, user4, community.LikeTargetTypeChapter, chapterID)
		assert.NoError(t, err)
		assert.True(t, isLiked)

		// 取消点赞
		err = repo.RemoveLike(ctx, user4, community.LikeTargetTypeChapter, chapterID)
		assert.NoError(t, err)

		// 再次验证
		isLiked, err = repo.IsLiked(ctx, user4, community.LikeTargetTypeChapter, chapterID)
		assert.NoError(t, err)
		assert.False(t, isLiked)

		t.Logf("✓ 章节点赞功能测试通过")
	})
}

// TestLikeRepositoryBoundary 边界测试
func TestLikeRepositoryBoundary(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := reading.NewMongoLikeRepository(db)
	ctx := context.Background()

	t.Run("AddLike_WithEmptyParams", func(t *testing.T) {
		// 空参数
		like := &community.Like{
			UserID:     "",
			TargetType: "",
			TargetID:   "",
			CreatedAt:  time.Now(),
		}

		err := repo.AddLike(ctx, like)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不完整")

		t.Logf("✓ 空参数正确返回错误")
	})

	t.Run("RemoveLike_WithEmptyParams", func(t *testing.T) {
		err := repo.RemoveLike(ctx, "", "", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不完整")

		t.Logf("✓ 取消点赞空参数正确返回错误")
	})

	t.Run("IsLiked_WithEmptyParams", func(t *testing.T) {
		isLiked, err := repo.IsLiked(ctx, "", "", "")
		assert.Error(t, err)
		assert.False(t, isLiked)
		assert.Contains(t, err.Error(), "不完整")

		t.Logf("✓ 检查点赞空参数正确返回错误")
	})

	t.Run("GetByID_WithInvalidID", func(t *testing.T) {
		found, err := repo.GetByID(ctx, "invalid_id")
		assert.Error(t, err)
		assert.Nil(t, found)
		assert.Contains(t, err.Error(), "invalid")

		t.Logf("✓ 无效ID正确返回错误")
	})

	t.Run("GetUserLikes_WithEmptyUserID", func(t *testing.T) {
		likes, total, err := repo.GetUserLikes(ctx, "", community.LikeTargetTypeBook, 1, 10)
		assert.Error(t, err)
		assert.Nil(t, likes)
		assert.Equal(t, int64(0), total)
		assert.Contains(t, err.Error(), "不能为空")

		t.Logf("✓ 空用户ID正确返回错误")
	})

	t.Run("GetLikeCount_WithEmptyParams", func(t *testing.T) {
		count, err := repo.GetLikeCount(ctx, "", "")
		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
		assert.Contains(t, err.Error(), "不完整")

		t.Logf("✓ 获取点赞数空参数正确返回错误")
	})

	t.Run("GetLikesCountBatch_WithEmptyArray", func(t *testing.T) {
		counts, err := repo.GetLikesCountBatch(ctx, community.LikeTargetTypeBook, []string{})
		assert.NoError(t, err)
		assert.NotNil(t, counts)
		assert.Equal(t, 0, len(counts))

		t.Logf("✓ 空数组批量查询正常处理")
	})

	t.Run("GetUserLikeStatusBatch_WithEmptyArray", func(t *testing.T) {
		status, err := repo.GetUserLikeStatusBatch(ctx, primitive.NewObjectID().Hex(), community.LikeTargetTypeBook, []string{})
		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.Equal(t, 0, len(status))

		t.Logf("✓ 空数组批量状态查询正常处理")
	})

	t.Run("CountUserLikes_WithEmptyUserID", func(t *testing.T) {
		count, err := repo.CountUserLikes(ctx, "")
		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
		assert.Contains(t, err.Error(), "不能为空")

		t.Logf("✓ 统计用户点赞空ID正确返回错误")
	})

	t.Run("Pagination_WithLargeNumbers", func(t *testing.T) {
		userID := primitive.NewObjectID().Hex()

		// 大页码和大页面大小
		likes, total, err := repo.GetUserLikes(ctx, userID, "", 1000, 10000)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(0))
		assert.NotNil(t, likes)

		t.Logf("✓ 大分页参数正常处理")
	})
}

// TestLikeRepositoryConcurrency 并发测试
func TestLikeRepositoryConcurrency(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := reading.NewMongoLikeRepository(db)
	ctx := context.Background()

	t.Run("ConcurrentAddLike_SameTarget", func(t *testing.T) {
		testBookID := primitive.NewObjectID().Hex()

		// 多个用户并发点赞同一本书
		const goroutines = 20
		var wg sync.WaitGroup
		wg.Add(goroutines)

		successCount := 0
		var mu sync.Mutex

		for i := 0; i < goroutines; i++ {
			go func(index int) {
				defer wg.Done()
				userID := primitive.NewObjectID().Hex()
				like := &community.Like{
					UserID:     userID,
					TargetType: community.LikeTargetTypeBook,
					TargetID:   testBookID,
					CreatedAt:  time.Now(),
				}
				err := repo.AddLike(ctx, like)
				if err == nil {
					mu.Lock()
					successCount++
					mu.Unlock()
				}
			}(i)
		}

		wg.Wait()

		// 验证所有点赞都成功
		assert.Equal(t, goroutines, successCount)

		// 验证点赞总数
		count, err := repo.GetLikeCount(ctx, community.LikeTargetTypeBook, testBookID)
		assert.NoError(t, err)
		assert.Equal(t, int64(goroutines), count)

		t.Logf("✓ 并发点赞测试通过，成功数: %d, 总点赞数: %d", successCount, count)
	})

	t.Run("ConcurrentAddLike_SameUser", func(t *testing.T) {
		testUser := primitive.NewObjectID().Hex()
		testBookID := primitive.NewObjectID().Hex()

		// 同一用户并发点赞（应该只成功一次）
		const goroutines = 10
		var wg sync.WaitGroup
		wg.Add(goroutines)

		successCount := 0
		errorCount := 0
		var mu sync.Mutex

		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				like := &community.Like{
					UserID:     testUser,
					TargetType: community.LikeTargetTypeBook,
					TargetID:   testBookID,
					CreatedAt:  time.Now(),
				}
				err := repo.AddLike(ctx, like)
				mu.Lock()
				if err == nil {
					successCount++
				} else {
					errorCount++
				}
				mu.Unlock()
			}()
		}

		wg.Wait()

		// 由于唯一索引，应该只有一次成功
		assert.Equal(t, 1, successCount)
		assert.Equal(t, goroutines-1, errorCount)

		// 验证只有一条点赞记录
		isLiked, err := repo.IsLiked(ctx, testUser, community.LikeTargetTypeBook, testBookID)
		assert.NoError(t, err)
		assert.True(t, isLiked)

		t.Logf("✓ 并发重复点赞测试通过，成功: %d, 失败: %d", successCount, errorCount)
	})

	t.Run("ConcurrentRemoveLike", func(t *testing.T) {
		testUser := primitive.NewObjectID().Hex()
		testBookID := primitive.NewObjectID().Hex()

		// 先添加点赞
		like := &community.Like{
			UserID:     testUser,
			TargetType: community.LikeTargetTypeBook,
			TargetID:   testBookID,
			CreatedAt:  time.Now(),
		}
		err := repo.AddLike(ctx, like)
		assert.NoError(t, err)

		// 并发取消点赞
		const goroutines = 10
		var wg sync.WaitGroup
		wg.Add(goroutines)

		successCount := 0
		var mu sync.Mutex

		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				err := repo.RemoveLike(ctx, testUser, community.LikeTargetTypeBook, testBookID)
				if err == nil {
					mu.Lock()
					successCount++
					mu.Unlock()
				}
			}()
		}

		wg.Wait()

		// 应该只有一次成功
		assert.Equal(t, 1, successCount)

		// 验证已取消
		isLiked, err := repo.IsLiked(ctx, testUser, community.LikeTargetTypeBook, testBookID)
		assert.NoError(t, err)
		assert.False(t, isLiked)

		t.Logf("✓ 并发取消点赞测试通过")
	})

	t.Run("ConcurrentIsLiked", func(t *testing.T) {
		testUser := primitive.NewObjectID().Hex()
		testBookID := primitive.NewObjectID().Hex()

		// 添加点赞
		like := &community.Like{
			UserID:     testUser,
			TargetType: community.LikeTargetTypeBook,
			TargetID:   testBookID,
			CreatedAt:  time.Now(),
		}
		err := repo.AddLike(ctx, like)
		assert.NoError(t, err)

		// 并发查询点赞状态
		const goroutines = 50
		var wg sync.WaitGroup
		wg.Add(goroutines)

		allTrue := true
		var mu sync.Mutex

		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				isLiked, err := repo.IsLiked(ctx, testUser, community.LikeTargetTypeBook, testBookID)
				assert.NoError(t, err)
				if !isLiked {
					mu.Lock()
					allTrue = false
					mu.Unlock()
				}
			}()
		}

		wg.Wait()

		assert.True(t, allTrue, "所有并发查询都应该返回true")

		t.Logf("✓ 并发查询点赞状态测试通过")
	})
}

// TestLikeRepositoryBatchOperations 批量操作性能测试
func TestLikeRepositoryBatchOperations(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := reading.NewMongoLikeRepository(db)
	ctx := context.Background()

	t.Run("BatchOperations_LargeDataset", func(t *testing.T) {
		// 创建大量点赞数据
		const targetCount = 100
		targetIDs := make([]string, targetCount)

		for i := 0; i < targetCount; i++ {
			targetIDs[i] = primitive.NewObjectID().Hex()

			// 为每个目标添加随机数量的点赞
			likeCount := (i % 10) + 1
			for j := 0; j < likeCount; j++ {
				userID := primitive.NewObjectID().Hex()
				like := &community.Like{
					UserID:     userID,
					TargetType: community.LikeTargetTypeBook,
					TargetID:   targetIDs[i],
					CreatedAt:  time.Now(),
				}
				err := repo.AddLike(ctx, like)
				assert.NoError(t, err)
			}
		}

		// 批量获取点赞数
		startTime := time.Now()
		counts, err := repo.GetLikesCountBatch(ctx, community.LikeTargetTypeBook, targetIDs)
		duration := time.Since(startTime)

		assert.NoError(t, err)
		assert.Equal(t, targetCount, len(counts))

		// 验证点赞数正确
		for i, targetID := range targetIDs {
			expectedCount := (i % 10) + 1
			assert.Equal(t, int64(expectedCount), counts[targetID])
		}

		t.Logf("✓ 批量获取%d个目标的点赞数，耗时: %v", targetCount, duration)
	})

	t.Run("BatchOperations_UserStatus", func(t *testing.T) {
		testUser := primitive.NewObjectID().Hex()

		// 创建测试数据
		const targetCount = 50
		targetIDs := make([]string, targetCount)

		for i := 0; i < targetCount; i++ {
			targetIDs[i] = primitive.NewObjectID().Hex()

			// 为前一半添加点赞
			if i < targetCount/2 {
				like := &community.Like{
					UserID:     testUser,
					TargetType: community.LikeTargetTypeBook,
					TargetID:   targetIDs[i],
					CreatedAt:  time.Now(),
				}
				err := repo.AddLike(ctx, like)
				assert.NoError(t, err)
			}
		}

		// 批量获取用户点赞状态
		startTime := time.Now()
		status, err := repo.GetUserLikeStatusBatch(ctx, testUser, community.LikeTargetTypeBook, targetIDs)
		duration := time.Since(startTime)

		assert.NoError(t, err)
		assert.Equal(t, targetCount, len(status))

		// 验证前一半是true，后一半是false
		for i, targetID := range targetIDs {
			if i < targetCount/2 {
				assert.True(t, status[targetID], fmt.Sprintf("目标%d应该已点赞", i))
			} else {
				assert.False(t, status[targetID], fmt.Sprintf("目标%d应该未点赞", i))
			}
		}

		t.Logf("✓ 批量获取%d个目标的用户点赞状态，耗时: %v", targetCount, duration)
	})
}

// TestLikeRepositoryHealth 健康检查测试
func TestLikeRepositoryHealth(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := reading.NewMongoLikeRepository(db)
	ctx := context.Background()

	t.Run("Health_DatabaseConnected", func(t *testing.T) {
		err := repo.Health(ctx)
		assert.NoError(t, err)

		t.Logf("✓ 健康检查通过")
	})
}
