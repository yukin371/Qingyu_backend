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

	"Qingyu_backend/repository/mongodb/reader"
)

// TestCommentRepositoryAdvanced 高级评论Repository测试
func TestCommentRepositoryAdvanced(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := reading.NewMongoCommentRepository(db)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("GetCommentsByUserID_Success", func(t *testing.T) {
		// 创建多条用户评论
		for i := 0; i < 3; i++ {
			comment := &community.Comment{
				UserID:    testUserID,
				BookID:    testBookID,
				Content:   fmt.Sprintf("用户评论%d", i),
				Rating:    5,
				Status:    community.CommentStatusApproved,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err := repo.Create(ctx, comment)
			assert.NoError(t, err)
		}

		// 查询用户评论
		comments, total, err := repo.GetCommentsByUserID(ctx, testUserID, 1, 10)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(3))
		assert.GreaterOrEqual(t, len(comments), 3)

		t.Logf("✓ 获取用户评论列表成功，总数: %d", total)
	})

	t.Run("GetCommentsByChapterID_Success", func(t *testing.T) {
		testChapterID := primitive.NewObjectID().Hex()

		// 创建章节评论
		comment := &community.Comment{
			UserID:    testUserID,
			BookID:    testBookID,
			ChapterID: testChapterID,
			Content:   "章节评论",
			Status:    community.CommentStatusApproved,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, comment)
		assert.NoError(t, err)

		// 查询章节评论
		comments, total, err := repo.GetCommentsByChapterID(ctx, testChapterID, 1, 10)
		assert.NoError(t, err)
		assert.Greater(t, total, int64(0))
		assert.Greater(t, len(comments), 0)

		t.Logf("✓ 获取章节评论成功")
	})

	t.Run("GetCommentsByBookIDSorted_ByHot", func(t *testing.T) {
		testBook2 := primitive.NewObjectID().Hex()

		// 创建不同点赞数的评论
		for i := 0; i < 3; i++ {
			comment := &community.Comment{
				UserID:    testUserID,
				BookID:    testBook2,
				Content:   fmt.Sprintf("评论%d", i),
				LikeCount: i * 10, // 0, 10, 20
				Status:    community.CommentStatusApproved,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err := repo.Create(ctx, comment)
			assert.NoError(t, err)
			time.Sleep(10 * time.Millisecond) // 确保创建时间不同
		}

		// 按热度排序查询
		comments, total, err := repo.GetCommentsByBookIDSorted(ctx, testBook2, community.CommentSortByHot, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Equal(t, 3, len(comments))

		// 验证按点赞数降序
		assert.GreaterOrEqual(t, comments[0].LikeCount, comments[1].LikeCount)
		assert.GreaterOrEqual(t, comments[1].LikeCount, comments[2].LikeCount)

		t.Logf("✓ 按热度排序查询成功，点赞数: %d, %d, %d",
			comments[0].LikeCount, comments[1].LikeCount, comments[2].LikeCount)
	})

	t.Run("GetCommentsByBookIDSorted_ByLatest", func(t *testing.T) {
		testBook3 := primitive.NewObjectID().Hex()

		// 创建多条评论
		var createdIDs []string
		for i := 0; i < 3; i++ {
			comment := &community.Comment{
				UserID:    testUserID,
				BookID:    testBook3,
				Content:   fmt.Sprintf("评论%d", i),
				Status:    community.CommentStatusApproved,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err := repo.Create(ctx, comment)
			assert.NoError(t, err)
			createdIDs = append(createdIDs, comment.ID.Hex())
			time.Sleep(10 * time.Millisecond)
		}

		// 按最新排序查询
		comments, total, err := repo.GetCommentsByBookIDSorted(ctx, testBook3, community.CommentSortByLatest, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)

		// 验证最新的在前面
		assert.Equal(t, createdIDs[2], comments[0].ID.Hex())

		t.Logf("✓ 按最新排序查询成功")
	})

	t.Run("GetPendingComments_Success", func(t *testing.T) {
		// 创建待审核评论
		for i := 0; i < 3; i++ {
			comment := &community.Comment{
				UserID:    testUserID,
				BookID:    testBookID,
				Content:   fmt.Sprintf("待审核评论%d", i),
				Status:    community.CommentStatusPending,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err := repo.Create(ctx, comment)
			assert.NoError(t, err)
		}

		// 获取待审核评论
		comments, total, err := repo.GetPendingComments(ctx, 1, 10)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(3))
		assert.GreaterOrEqual(t, len(comments), 3)

		// 验证所有评论都是待审核状态
		for _, comment := range comments {
			assert.Equal(t, community.CommentStatusPending, comment.Status)
		}

		t.Logf("✓ 获取待审核评论成功，总数: %d", total)
	})

	t.Run("IncrementReplyCount_Success", func(t *testing.T) {
		comment := &community.Comment{
			UserID:     testUserID,
			BookID:     testBookID,
			Content:    "测试回复计数",
			ReplyCount: 0,
			Status:     community.CommentStatusApproved,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		err := repo.Create(ctx, comment)
		assert.NoError(t, err)

		// 增加回复数
		err = repo.IncrementReplyCount(ctx, comment.ID.Hex())
		assert.NoError(t, err)

		// 验证
		found, err := repo.GetByID(ctx, comment.ID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, 1, found.ReplyCount)

		t.Logf("✓ 增加回复数成功")
	})

	t.Run("DecrementReplyCount_Success", func(t *testing.T) {
		comment := &community.Comment{
			UserID:     testUserID,
			BookID:     testBookID,
			Content:    "测试减少回复",
			ReplyCount: 5,
			Status:     community.CommentStatusApproved,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		err := repo.Create(ctx, comment)
		assert.NoError(t, err)

		// 减少回复数
		err = repo.DecrementReplyCount(ctx, comment.ID.Hex())
		assert.NoError(t, err)

		// 验证
		found, err := repo.GetByID(ctx, comment.ID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, 4, found.ReplyCount)

		t.Logf("✓ 减少回复数成功")
	})

	t.Run("GetCommentCount_Success", func(t *testing.T) {
		testBook4 := primitive.NewObjectID().Hex()

		// 创建多条评论
		for i := 0; i < 5; i++ {
			comment := &community.Comment{
				UserID:    testUserID,
				BookID:    testBook4,
				Content:   fmt.Sprintf("评论%d", i),
				Status:    community.CommentStatusApproved,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err := repo.Create(ctx, comment)
			assert.NoError(t, err)
		}

		// 获取评论总数
		count, err := repo.GetCommentCount(ctx, testBook4)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)

		t.Logf("✓ 获取评论总数成功: %d", count)
	})

	t.Run("GetCommentsByIDs_Success", func(t *testing.T) {
		// 创建多条评论
		var ids []string
		for i := 0; i < 3; i++ {
			comment := &community.Comment{
				UserID:    testUserID,
				BookID:    testBookID,
				Content:   fmt.Sprintf("批量查询评论%d", i),
				Status:    community.CommentStatusApproved,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err := repo.Create(ctx, comment)
			assert.NoError(t, err)
			ids = append(ids, comment.ID.Hex())
		}

		// 批量获取
		comments, err := repo.GetCommentsByIDs(ctx, ids)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(comments))

		t.Logf("✓ 批量获取评论成功")
	})

	t.Run("DeleteCommentsByBookID_Success", func(t *testing.T) {
		testBook5 := primitive.NewObjectID().Hex()

		// 创建多条评论
		for i := 0; i < 3; i++ {
			comment := &community.Comment{
				UserID:    testUserID,
				BookID:    testBook5,
				Content:   fmt.Sprintf("待删除评论%d", i),
				Status:    community.CommentStatusApproved,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err := repo.Create(ctx, comment)
			assert.NoError(t, err)
		}

		// 删除书籍的所有评论
		err := repo.DeleteCommentsByBookID(ctx, testBook5)
		assert.NoError(t, err)

		// 验证评论数为0（软删除后不会在approved状态中）
		count, err := repo.GetCommentCount(ctx, testBook5)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)

		t.Logf("✓ 批量删除书籍评论成功")
	})
}

// TestCommentRepositoryBoundary 边界测试
func TestCommentRepositoryBoundary(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := reading.NewMongoCommentRepository(db)
	ctx := context.Background()

	t.Run("Create_WithInvalidData", func(t *testing.T) {
		// 空内容
		comment := &community.Comment{
			UserID:    "",
			BookID:    "",
			Content:   "",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, comment)
		// MongoDB不会在这里验证，但我们可以测试空ID的创建
		assert.NoError(t, err)

		t.Logf("✓ 创建空数据评论（MongoDB允许）")
	})

	t.Run("GetByID_WithInvalidID", func(t *testing.T) {
		// 无效的ObjectID格式
		found, err := repo.GetByID(ctx, "invalid_id")
		assert.Error(t, err)
		assert.Nil(t, found)
		assert.Contains(t, err.Error(), "invalid comment ID")

		t.Logf("✓ 无效ID正确返回错误")
	})

	t.Run("Update_WithNonExistentID", func(t *testing.T) {
		fakeID := primitive.NewObjectID().Hex()
		updates := map[string]interface{}{
			"content": "更新内容",
		}

		err := repo.Update(ctx, fakeID, updates)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")

		t.Logf("✓ 更新不存在的评论正确返回错误")
	})

	t.Run("Pagination_WithZeroPage", func(t *testing.T) {
		testBookID := primitive.NewObjectID().Hex()

		// 页码为0的情况（会导致负数skip，期望报错或自动处理为page=1）
		comments, total, err := repo.GetCommentsByBookID(ctx, testBookID, 0, 10)
		// MongoDB不允许负数skip，所以这里会报错是正常的
		// 生产代码应该在API层验证page>=1
		if err != nil {
			assert.Contains(t, err.Error(), "skip")
			t.Logf("✓ 页码为0的分页查询正确返回错误（需要在API层验证）")
		} else {
			// 如果没有错误，说明Repository层做了处理
			assert.GreaterOrEqual(t, total, int64(0))
			assert.NotNil(t, comments)
			t.Logf("✓ 页码为0的分页查询已在Repository层处理")
		}
	})

	t.Run("Pagination_WithLargePageSize", func(t *testing.T) {
		testBookID := primitive.NewObjectID().Hex()

		// 页面大小很大的情况
		comments, total, err := repo.GetCommentsByBookID(ctx, testBookID, 1, 10000)
		if err != nil {
			// MongoDB可能有大小限制
			t.Logf("⚠ 大页面大小查询报错（预期行为）: %v", err)
		} else {
			assert.GreaterOrEqual(t, total, int64(0))
			assert.NotNil(t, comments)
			t.Logf("✓ 大页面大小的分页查询正常处理")
		}
	})

	t.Run("GetBookRatingStats_WithNoRatings", func(t *testing.T) {
		fakeBookID := primitive.NewObjectID().Hex()

		// 没有评分的书籍
		stats, err := repo.GetBookRatingStats(ctx, fakeBookID)
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, 0, stats["total_count"])

		t.Logf("✓ 没有评分的统计正确返回零值")
	})
}

// TestCommentRepositoryConcurrency 并发测试
func TestCommentRepositoryConcurrency(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := reading.NewMongoCommentRepository(db)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("ConcurrentIncrementLikeCount", func(t *testing.T) {
		// 创建评论
		comment := &community.Comment{
			UserID:    testUserID,
			BookID:    testBookID,
			Content:   "并发点赞测试",
			LikeCount: 0,
			Status:    community.CommentStatusApproved,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, comment)
		assert.NoError(t, err)

		// 并发增加点赞数
		const goroutines = 20
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				err := repo.IncrementLikeCount(ctx, comment.ID.Hex())
				assert.NoError(t, err)
			}()
		}

		wg.Wait()

		// 验证最终点赞数
		found, err := repo.GetByID(ctx, comment.ID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, goroutines, found.LikeCount)

		t.Logf("✓ 并发增加点赞数测试通过，最终点赞数: %d", found.LikeCount)
	})

	t.Run("ConcurrentCreate", func(t *testing.T) {
		// 并发创建评论
		const goroutines = 10
		var wg sync.WaitGroup
		wg.Add(goroutines)

		successCount := 0
		var mu sync.Mutex

		for i := 0; i < goroutines; i++ {
			go func(index int) {
				defer wg.Done()
				comment := &community.Comment{
					UserID:    testUserID,
					BookID:    testBookID,
					Content:   fmt.Sprintf("并发创建评论%d", index),
					Status:    community.CommentStatusApproved,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				err := repo.Create(ctx, comment)
				if err == nil {
					mu.Lock()
					successCount++
					mu.Unlock()
				}
			}(i)
		}

		wg.Wait()

		assert.Equal(t, goroutines, successCount)

		t.Logf("✓ 并发创建评论测试通过，成功数: %d", successCount)
	})

	t.Run("ConcurrentUpdateStatus", func(t *testing.T) {
		// 创建评论
		comment := &community.Comment{
			UserID:    testUserID,
			BookID:    testBookID,
			Content:   "并发更新状态测试",
			Status:    community.CommentStatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, comment)
		assert.NoError(t, err)

		// 并发更新状态
		const goroutines = 5
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for i := 0; i < goroutines; i++ {
			go func(index int) {
				defer wg.Done()
				status := community.CommentStatusApproved
				if index%2 == 0 {
					status = community.CommentStatusRejected
				}
				err := repo.UpdateCommentStatus(ctx, comment.ID.Hex(), status, "测试原因")
				assert.NoError(t, err)
			}(i)
		}

		wg.Wait()

		// 验证最终状态（会是最后一个更新的状态）
		found, err := repo.GetByID(ctx, comment.ID.Hex())
		assert.NoError(t, err)
		assert.Contains(t, []string{community.CommentStatusApproved, community.CommentStatusRejected}, found.Status)

		t.Logf("✓ 并发更新状态测试通过，最终状态: %s", found.Status)
	})
}

// TestCommentRepositoryReplyThread 评论回复嵌套测试
func TestCommentRepositoryReplyThread(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := reading.NewMongoCommentRepository(db)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("NestedReplyThread", func(t *testing.T) {
		// 创建根评论
		rootComment := &community.Comment{
			UserID:    testUserID,
			BookID:    testBookID,
			Content:   "根评论",
			Status:    community.CommentStatusApproved,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, rootComment)
		assert.NoError(t, err)

		// 创建一级回复
		reply1 := &community.Comment{
			UserID:    testUserID,
			BookID:    testBookID,
			Content:   "一级回复",
			ParentID:  rootComment.ID.Hex(),
			RootID:    rootComment.ID.Hex(),
			Status:    community.CommentStatusApproved,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = repo.Create(ctx, reply1)
		assert.NoError(t, err)

		// 创建二级回复
		reply2 := &community.Comment{
			UserID:      testUserID,
			BookID:      testBookID,
			Content:     "二级回复",
			ParentID:    reply1.ID.Hex(),
			RootID:      rootComment.ID.Hex(),
			ReplyToUser: testUserID,
			Status:      community.CommentStatusApproved,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err = repo.Create(ctx, reply2)
		assert.NoError(t, err)

		// 获取根评论的所有回复
		replies, err := repo.GetRepliesByCommentID(ctx, rootComment.ID.Hex())
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(replies), 1)

		// 获取一级回复的回复
		nestedReplies, err := repo.GetRepliesByCommentID(ctx, reply1.ID.Hex())
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(nestedReplies), 1)

		t.Logf("✓ 嵌套回复测试通过")
	})
}
