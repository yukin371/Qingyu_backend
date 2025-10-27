package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/reader"
	"Qingyu_backend/repository/mongodb/reading"
)

// TestLikeRepository 点赞Repository测试
func TestLikeRepository(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := reading.NewMongoLikeRepository(db)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("AddLike_Success", func(t *testing.T) {
		like := &reader.Like{
			UserID:     testUserID,
			TargetType: reader.LikeTargetTypeBook,
			TargetID:   testBookID,
			CreatedAt:  time.Now(),
		}

		err := repo.AddLike(ctx, like)
		assert.NoError(t, err)
		assert.False(t, like.ID.IsZero(), "点赞ID应该被设置")

		t.Logf("✓ 添加点赞成功，ID: %s", like.ID.Hex())
	})

	t.Run("AddLike_Duplicate", func(t *testing.T) {
		like := &reader.Like{
			UserID:     testUserID,
			TargetType: reader.LikeTargetTypeBook,
			TargetID:   testBookID,
			CreatedAt:  time.Now(),
		}

		// 第二次点赞应该失败
		err := repo.AddLike(ctx, like)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "已经点赞")

		t.Logf("✓ 重复点赞检测通过")
	})

	t.Run("IsLiked_True", func(t *testing.T) {
		isLiked, err := repo.IsLiked(ctx, testUserID, reader.LikeTargetTypeBook, testBookID)
		assert.NoError(t, err)
		assert.True(t, isLiked, "应该已经点赞")

		t.Logf("✓ 检查点赞状态成功")
	})

	t.Run("IsLiked_False", func(t *testing.T) {
		fakeBookID := primitive.NewObjectID().Hex()
		isLiked, err := repo.IsLiked(ctx, testUserID, reader.LikeTargetTypeBook, fakeBookID)
		assert.NoError(t, err)
		assert.False(t, isLiked, "应该未点赞")

		t.Logf("✓ 未点赞状态正确")
	})

	t.Run("GetLikeCount_Success", func(t *testing.T) {
		count, err := repo.GetLikeCount(ctx, reader.LikeTargetTypeBook, testBookID)
		assert.NoError(t, err)
		assert.Greater(t, count, int64(0))

		t.Logf("✓ 获取点赞数成功: %d", count)
	})

	t.Run("RemoveLike_Success", func(t *testing.T) {
		err := repo.RemoveLike(ctx, testUserID, reader.LikeTargetTypeBook, testBookID)
		assert.NoError(t, err)

		// 验证已取消
		isLiked, err := repo.IsLiked(ctx, testUserID, reader.LikeTargetTypeBook, testBookID)
		assert.NoError(t, err)
		assert.False(t, isLiked, "应该已取消点赞")

		t.Logf("✓ 取消点赞成功")
	})

	t.Run("RemoveLike_NotFound", func(t *testing.T) {
		fakeBookID := primitive.NewObjectID().Hex()
		err := repo.RemoveLike(ctx, testUserID, reader.LikeTargetTypeBook, fakeBookID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不存在")

		t.Logf("✓ 取消不存在的点赞正确返回错误")
	})

	t.Run("GetUserLikes_WithPagination", func(t *testing.T) {
		// 添加多个点赞
		for i := 0; i < 5; i++ {
			bookID := primitive.NewObjectID().Hex()
			like := &reader.Like{
				UserID:     testUserID,
				TargetType: reader.LikeTargetTypeBook,
				TargetID:   bookID,
				CreatedAt:  time.Now(),
			}
			err := repo.AddLike(ctx, like)
			assert.NoError(t, err)
		}

		// 查询
		likes, total, err := repo.GetUserLikes(ctx, testUserID, reader.LikeTargetTypeBook, 1, 3)
		assert.NoError(t, err)
		assert.Greater(t, total, int64(0))
		assert.LessOrEqual(t, len(likes), 3)

		t.Logf("✓ 分页查询用户点赞成功，总数: %d, 本页: %d", total, len(likes))
	})

	t.Run("CountUserLikes_Success", func(t *testing.T) {
		count, err := repo.CountUserLikes(ctx, testUserID)
		assert.NoError(t, err)
		assert.Greater(t, count, int64(0))

		t.Logf("✓ 统计用户点赞数成功: %d", count)
	})

	t.Run("Health_Success", func(t *testing.T) {
		err := repo.Health(ctx)
		assert.NoError(t, err)

		t.Logf("✓ 健康检查通过")
	})
}

// TestLikeRepositoryBatch 批量操作测试
func TestLikeRepositoryBatch(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := reading.NewMongoLikeRepository(db)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	// 创建测试数据
	bookIDs := make([]string, 5)
	for i := 0; i < 5; i++ {
		bookIDs[i] = primitive.NewObjectID().Hex()
		// 为前3本书添加点赞
		if i < 3 {
			like := &reader.Like{
				UserID:     testUserID,
				TargetType: reader.LikeTargetTypeBook,
				TargetID:   bookIDs[i],
				CreatedAt:  time.Now(),
			}
			err := repo.AddLike(ctx, like)
			assert.NoError(t, err)
		}
	}

	t.Run("GetLikesCountBatch_Success", func(t *testing.T) {
		counts, err := repo.GetLikesCountBatch(ctx, reader.LikeTargetTypeBook, bookIDs)
		assert.NoError(t, err)
		assert.Equal(t, len(bookIDs), len(counts))

		// 验证前3本书有点赞数
		for i := 0; i < 3; i++ {
			assert.Greater(t, counts[bookIDs[i]], int64(0), "应该有点赞")
		}
		// 验证后2本书没有点赞
		for i := 3; i < 5; i++ {
			assert.Equal(t, int64(0), counts[bookIDs[i]], "应该没有点赞")
		}

		t.Logf("✓ 批量获取点赞数成功")
	})

	t.Run("GetUserLikeStatusBatch_Success", func(t *testing.T) {
		status, err := repo.GetUserLikeStatusBatch(ctx, testUserID, reader.LikeTargetTypeBook, bookIDs)
		assert.NoError(t, err)
		assert.Equal(t, len(bookIDs), len(status))

		// 验证前3本书已点赞
		for i := 0; i < 3; i++ {
			assert.True(t, status[bookIDs[i]], "应该已点赞")
		}
		// 验证后2本书未点赞
		for i := 3; i < 5; i++ {
			assert.False(t, status[bookIDs[i]], "应该未点赞")
		}

		t.Logf("✓ 批量检查点赞状态成功")
	})
}

// TestLikeRepositoryCommentLike 评论点赞测试
func TestLikeRepositoryCommentLike(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := reading.NewMongoLikeRepository(db)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testCommentID := primitive.NewObjectID().Hex()

	t.Run("LikeComment_Success", func(t *testing.T) {
		like := &reader.Like{
			UserID:     testUserID,
			TargetType: reader.LikeTargetTypeComment,
			TargetID:   testCommentID,
			CreatedAt:  time.Now(),
		}

		err := repo.AddLike(ctx, like)
		assert.NoError(t, err)

		// 验证点赞状态
		isLiked, err := repo.IsLiked(ctx, testUserID, reader.LikeTargetTypeComment, testCommentID)
		assert.NoError(t, err)
		assert.True(t, isLiked)

		t.Logf("✓ 点赞评论成功")
	})

	t.Run("UnlikeComment_Success", func(t *testing.T) {
		err := repo.RemoveLike(ctx, testUserID, reader.LikeTargetTypeComment, testCommentID)
		assert.NoError(t, err)

		// 验证已取消
		isLiked, err := repo.IsLiked(ctx, testUserID, reader.LikeTargetTypeComment, testCommentID)
		assert.NoError(t, err)
		assert.False(t, isLiked)

		t.Logf("✓ 取消点赞评论成功")
	})
}
