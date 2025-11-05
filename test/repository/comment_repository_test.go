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

// TestCommentRepository 评论Repository测试
func TestCommentRepository(t *testing.T) {
	// 设置测试数据库
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := reading.NewMongoCommentRepository(db)
	ctx := context.Background()

	// 测试用数据
	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("Create_Success", func(t *testing.T) {
		comment := &reader.Comment{
			UserID:    testUserID,
			BookID:    testBookID,
			Content:   "这是一条测试评论",
			Rating:    5,
			Status:    "approved",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, comment)
		assert.NoError(t, err)
		assert.False(t, comment.ID.IsZero(), "评论ID应该被设置")

		t.Logf("✓ 创建评论成功，ID: %s", comment.ID.Hex())
	})

	t.Run("GetByID_Success", func(t *testing.T) {
		// 先创建一条评论
		comment := &reader.Comment{
			UserID:    testUserID,
			BookID:    testBookID,
			Content:   "测试获取评论",
			Rating:    4,
			Status:    "approved",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, comment)
		assert.NoError(t, err)

		// 获取评论
		found, err := repo.GetByID(ctx, comment.ID.Hex())
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, comment.Content, found.Content)

		t.Logf("✓ 获取评论成功")
	})

	t.Run("GetByID_NotFound", func(t *testing.T) {
		fakeID := primitive.NewObjectID().Hex()
		found, err := repo.GetByID(ctx, fakeID)
		assert.Error(t, err)
		assert.Nil(t, found)

		t.Logf("✓ 不存在的评论正确返回错误")
	})

	t.Run("GetCommentsByBookID_WithPagination", func(t *testing.T) {
		// 创建多条评论
		for i := 0; i < 5; i++ {
			comment := &reader.Comment{
				UserID:    testUserID,
				BookID:    testBookID,
				Content:   "测试评论" + string(rune(i)),
				Rating:    5,
				Status:    "approved",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err := repo.Create(ctx, comment)
			assert.NoError(t, err)
		}

		// 查询
		comments, total, err := repo.GetCommentsByBookID(ctx, testBookID, 1, 3)
		assert.NoError(t, err)
		assert.Greater(t, total, int64(0))
		assert.LessOrEqual(t, len(comments), 3)

		t.Logf("✓ 分页查询成功，总数: %d, 本页: %d", total, len(comments))
	})

	t.Run("Update_Success", func(t *testing.T) {
		// 创建评论
		comment := &reader.Comment{
			UserID:    testUserID,
			BookID:    testBookID,
			Content:   "原始内容",
			Rating:    3,
			Status:    "approved",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, comment)
		assert.NoError(t, err)

		// 更新评论
		updates := map[string]interface{}{
			"content": "更新后的内容",
			"rating":  5,
		}
		err = repo.Update(ctx, comment.ID.Hex(), updates)
		assert.NoError(t, err)

		// 验证更新
		found, err := repo.GetByID(ctx, comment.ID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, "更新后的内容", found.Content)
		assert.Equal(t, 5, found.Rating)

		t.Logf("✓ 更新评论成功")
	})

	t.Run("Delete_Success", func(t *testing.T) {
		// 创建评论
		comment := &reader.Comment{
			UserID:    testUserID,
			BookID:    testBookID,
			Content:   "待删除的评论",
			Status:    "approved",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, comment)
		assert.NoError(t, err)

		// 删除评论
		err = repo.Delete(ctx, comment.ID.Hex())
		assert.NoError(t, err)

		// 验证已删除
		found, err := repo.GetByID(ctx, comment.ID.Hex())
		assert.Error(t, err)
		assert.Nil(t, found)

		t.Logf("✓ 删除评论成功")
	})

	t.Run("IncrementLikeCount_Success", func(t *testing.T) {
		// 创建评论
		comment := &reader.Comment{
			UserID:    testUserID,
			BookID:    testBookID,
			Content:   "点赞测试",
			LikeCount: 0,
			Status:    "approved",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, comment)
		assert.NoError(t, err)

		// 增加点赞数
		err = repo.IncrementLikeCount(ctx, comment.ID.Hex())
		assert.NoError(t, err)

		// 验证点赞数
		found, err := repo.GetByID(ctx, comment.ID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, 1, found.LikeCount)

		t.Logf("✓ 增加点赞数成功")
	})

	t.Run("DecrementLikeCount_Success", func(t *testing.T) {
		// 创建评论
		comment := &reader.Comment{
			UserID:    testUserID,
			BookID:    testBookID,
			Content:   "取消点赞测试",
			LikeCount: 5,
			Status:    "approved",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, comment)
		assert.NoError(t, err)

		// 减少点赞数
		err = repo.DecrementLikeCount(ctx, comment.ID.Hex())
		assert.NoError(t, err)

		// 验证点赞数
		found, err := repo.GetByID(ctx, comment.ID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, 4, found.LikeCount)

		t.Logf("✓ 减少点赞数成功")
	})

	t.Run("UpdateCommentStatus_Success", func(t *testing.T) {
		// 创建待审核评论
		comment := &reader.Comment{
			UserID:    testUserID,
			BookID:    testBookID,
			Content:   "待审核评论",
			Status:    "pending",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, comment)
		assert.NoError(t, err)

		// 更新审核状态
		err = repo.UpdateCommentStatus(ctx, comment.ID.Hex(), "rejected", "包含敏感词")
		assert.NoError(t, err)

		// 验证状态
		found, err := repo.GetByID(ctx, comment.ID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, "rejected", found.Status)
		assert.Equal(t, "包含敏感词", found.RejectReason)

		t.Logf("✓ 更新审核状态成功")
	})

	t.Run("GetRepliesByCommentID_Success", func(t *testing.T) {
		// 创建父评论
		parentComment := &reader.Comment{
			UserID:    testUserID,
			BookID:    testBookID,
			Content:   "父评论",
			Status:    "approved",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, parentComment)
		assert.NoError(t, err)

		// 创建回复
		reply := &reader.Comment{
			UserID:    testUserID,
			BookID:    testBookID,
			Content:   "回复内容",
			ParentID:  parentComment.ID.Hex(),
			Status:    "approved",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = repo.Create(ctx, reply)
		assert.NoError(t, err)

		// 获取回复列表
		replies, err := repo.GetRepliesByCommentID(ctx, parentComment.ID.Hex())
		assert.NoError(t, err)
		assert.Greater(t, len(replies), 0)

		t.Logf("✓ 获取回复列表成功，回复数: %d", len(replies))
	})

	t.Run("Health_Success", func(t *testing.T) {
		err := repo.Health(ctx)
		assert.NoError(t, err)

		t.Logf("✓ 健康检查通过")
	})
}

// TestCommentRepositoryStatistics 评论统计测试
func TestCommentRepositoryStatistics(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := reading.NewMongoCommentRepository(db)
	ctx := context.Background()

	testBookID := primitive.NewObjectID().Hex()
	testUserID := primitive.NewObjectID().Hex()

	t.Run("GetBookRatingStats_Success", func(t *testing.T) {
		// 创建不同评分的评论
		ratings := []int{5, 5, 4, 3, 5}
		for _, rating := range ratings {
			comment := &reader.Comment{
				UserID:    testUserID,
				BookID:    testBookID,
				Content:   "测试评论",
				Rating:    rating,
				Status:    "approved",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err := repo.Create(ctx, comment)
			assert.NoError(t, err)
		}

		// 获取评分统计
		stats, err := repo.GetBookRatingStats(ctx, testBookID)
		assert.NoError(t, err)
		assert.NotNil(t, stats)

		t.Logf("✓ 评分统计: %+v", stats)
	})
}
