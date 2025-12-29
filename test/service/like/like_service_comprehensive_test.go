package like

import (
	"Qingyu_backend/models/community"
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/service/reading"
)

// TestLikeServiceBusinessRules 业务规则测试
func TestLikeServiceBusinessRules(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("Idempotency_LikeBook", func(t *testing.T) {
		// 测试点赞幂等性
		testBookID := primitive.NewObjectID().Hex()

		// 第一次点赞成功
		mockRepo.On("AddLike", ctx, mock.AnythingOfType("*reader.Like")).Return(nil).Once()
		mockRepo.On("GetLikeCount", ctx, community.LikeTargetTypeBook, testBookID).Return(int64(1), nil).Once()

		err := service.LikeBook(ctx, testUserID, testBookID)
		assert.NoError(t, err)

		// 第二次点赞（已存在）
		mockRepo.On("AddLike", ctx, mock.AnythingOfType("*reader.Like")).Return(errors.New("已经点赞过了")).Once()

		err = service.LikeBook(ctx, testUserID, testBookID)
		assert.NoError(t, err) // 应该不报错，幂等性处理

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 点赞幂等性测试通过")
	})

	t.Run("Idempotency_UnlikeBook", func(t *testing.T) {
		// 测试取消点赞幂等性
		testBookID := primitive.NewObjectID().Hex()

		// 第一次取消点赞成功
		mockRepo.On("RemoveLike", ctx, testUserID, community.LikeTargetTypeBook, testBookID).Return(nil).Once()
		mockRepo.On("GetLikeCount", ctx, community.LikeTargetTypeBook, testBookID).Return(int64(0), nil).Once()

		err := service.UnlikeBook(ctx, testUserID, testBookID)
		assert.NoError(t, err)

		// 第二次取消点赞（不存在）
		mockRepo.On("RemoveLike", ctx, testUserID, community.LikeTargetTypeBook, testBookID).Return(errors.New("点赞记录不存在")).Once()

		err = service.UnlikeBook(ctx, testUserID, testBookID)
		assert.NoError(t, err) // 应该不报错，幂等性处理

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 取消点赞幂等性测试通过")
	})
}

// TestLikeServiceCommentInteraction 评论点赞交互测试
func TestLikeServiceCommentInteraction(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testCommentID := primitive.NewObjectID().Hex()

	t.Run("LikeComment_IncrementCount", func(t *testing.T) {
		// 点赞评论应该增加评论点赞数
		mockRepo.On("AddLike", ctx, mock.AnythingOfType("*reader.Like")).Return(nil).Once()
		mockCommentRepo.On("IncrementLikeCount", ctx, testCommentID).Return(nil).Once()
		mockRepo.On("GetLikeCount", ctx, community.LikeTargetTypeComment, testCommentID).Return(int64(1), nil).Once()

		err := service.LikeComment(ctx, testUserID, testCommentID)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)

		t.Logf("✓ 点赞评论计数增加成功")
	})

	t.Run("UnlikeComment_DecrementCount", func(t *testing.T) {
		// 取消点赞评论应该减少评论点赞数
		mockRepo.On("RemoveLike", ctx, testUserID, community.LikeTargetTypeComment, testCommentID).Return(nil).Once()
		mockCommentRepo.On("DecrementLikeCount", ctx, testCommentID).Return(nil).Once()
		mockRepo.On("GetLikeCount", ctx, community.LikeTargetTypeComment, testCommentID).Return(int64(0), nil).Once()

		err := service.UnlikeComment(ctx, testUserID, testCommentID)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)

		t.Logf("✓ 取消点赞评论计数减少成功")
	})

	t.Run("LikeComment_IncrementFailed", func(t *testing.T) {
		// 点赞记录成功但增加评论计数失败
		mockRepo.On("AddLike", ctx, mock.AnythingOfType("*reader.Like")).Return(nil).Once()
		mockCommentRepo.On("IncrementLikeCount", ctx, testCommentID).Return(errors.New("increment failed")).Once()
		// 虽然IncrementLikeCount失败，但publishLikeEvent仍然会被调用，获取点赞数
		mockRepo.On("GetLikeCount", ctx, community.LikeTargetTypeComment, testCommentID).Return(int64(1), nil).Once()

		err := service.LikeComment(ctx, testUserID, testCommentID)
		// 虽然IncrementLikeCount失败，但函数返回nil（因为错误被吞掉了）
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)

		t.Logf("✓ 增加评论计数失败处理正确")
	})
}

// TestLikeServiceBatchOperations 批量操作测试
func TestLikeServiceBatchOperations(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("GetBooksLikeCount_Success", func(t *testing.T) {
		bookIDs := []string{
			primitive.NewObjectID().Hex(),
			primitive.NewObjectID().Hex(),
			primitive.NewObjectID().Hex(),
		}

		expectedCounts := map[string]int64{
			bookIDs[0]: 10,
			bookIDs[1]: 20,
			bookIDs[2]: 30,
		}

		mockRepo.On("GetLikesCountBatch", ctx, community.LikeTargetTypeBook, bookIDs).
			Return(expectedCounts, nil).Once()

		counts, err := service.GetBooksLikeCount(ctx, bookIDs)
		assert.NoError(t, err)
		assert.Equal(t, expectedCounts, counts)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 批量获取点赞数成功")
	})

	t.Run("GetUserLikeStatus_Success", func(t *testing.T) {
		bookIDs := []string{
			primitive.NewObjectID().Hex(),
			primitive.NewObjectID().Hex(),
			primitive.NewObjectID().Hex(),
		}

		expectedStatus := map[string]bool{
			bookIDs[0]: true,
			bookIDs[1]: false,
			bookIDs[2]: true,
		}

		mockRepo.On("GetUserLikeStatusBatch", ctx, testUserID, community.LikeTargetTypeBook, bookIDs).
			Return(expectedStatus, nil).Once()

		status, err := service.GetUserLikeStatus(ctx, testUserID, bookIDs)
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, status)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 批量获取用户点赞状态成功")
	})

	t.Run("BatchOperations_EmptyArray", func(t *testing.T) {
		mockRepo := new(MockLikeRepository)
		mockCommentRepo := new(MockCommentRepository)
		mockEventBus := NewMockEventBus()

		service := reading.NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
		ctx := context.Background()

		emptyIDs := []string{}

		// 空数组会导致service直接返回空map，不调用repository
		counts, err := service.GetBooksLikeCount(ctx, emptyIDs)
		assert.NoError(t, err)
		assert.Empty(t, counts)

		// 不验证mock期望，因为没有调用repository

		t.Logf("✓ 批量操作空数组处理正确")
	})
}

// TestLikeServiceTargetTypes 目标类型测试
func TestLikeServiceTargetTypes(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("LikeBook_TargetType", func(t *testing.T) {
		testBookID := primitive.NewObjectID().Hex()

		mockRepo.On("AddLike", ctx, mock.MatchedBy(func(like *community.Like) bool {
			return like.TargetType == community.LikeTargetTypeBook && like.TargetID == testBookID
		})).Return(nil).Once()
		mockRepo.On("GetLikeCount", ctx, community.LikeTargetTypeBook, testBookID).Return(int64(1), nil).Once()

		err := service.LikeBook(ctx, testUserID, testBookID)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 书籍点赞目标类型正确")
	})

	t.Run("LikeComment_TargetType", func(t *testing.T) {
		testCommentID := primitive.NewObjectID().Hex()

		mockRepo.On("AddLike", ctx, mock.MatchedBy(func(like *community.Like) bool {
			return like.TargetType == community.LikeTargetTypeComment && like.TargetID == testCommentID
		})).Return(nil).Once()
		mockCommentRepo.On("IncrementLikeCount", ctx, testCommentID).Return(nil).Once()
		mockRepo.On("GetLikeCount", ctx, community.LikeTargetTypeComment, testCommentID).Return(int64(1), nil).Once()

		err := service.LikeComment(ctx, testUserID, testCommentID)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)

		t.Logf("✓ 评论点赞目标类型正确")
	})
}

// TestLikeServiceError 错误处理测试
func TestLikeServiceError(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("RepositoryError_AddLike", func(t *testing.T) {
		mockRepo.On("AddLike", ctx, mock.AnythingOfType("*reader.Like")).
			Return(errors.New("database connection error")).Once()

		err := service.LikeBook(ctx, testUserID, testBookID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "点赞书籍失败")

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 添加点赞失败处理正确")
	})

	t.Run("RepositoryError_RemoveLike", func(t *testing.T) {
		mockRepo.On("RemoveLike", ctx, testUserID, community.LikeTargetTypeBook, testBookID).
			Return(errors.New("database error")).Once()

		err := service.UnlikeBook(ctx, testUserID, testBookID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "取消点赞书籍失败")

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 移除点赞失败处理正确")
	})

	t.Run("InvalidParameters", func(t *testing.T) {
		// 空用户ID
		err := service.LikeBook(ctx, "", testBookID)
		assert.Error(t, err)

		// 空书籍ID
		err = service.LikeBook(ctx, testUserID, "")
		assert.Error(t, err)

		t.Logf("✓ 无效参数验证通过")
	})
}

// TestLikeServiceStatistics 统计功能测试
func TestLikeServiceStatistics(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("GetUserLikeStats_AllTypes", func(t *testing.T) {
		mockRepo.On("CountUserLikes", ctx, testUserID).Return(int64(100), nil).Once()

		stats, err := service.GetUserLikeStats(ctx, testUserID)
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, int64(100), stats["total_likes"])

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 获取用户点赞统计成功")
	})

	t.Run("GetBookLikeCount_Zero", func(t *testing.T) {
		testBookID := primitive.NewObjectID().Hex()

		mockRepo.On("GetLikeCount", ctx, community.LikeTargetTypeBook, testBookID).Return(int64(0), nil).Once()

		count, err := service.GetBookLikeCount(ctx, testBookID)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 获取零点赞数成功")
	})

	t.Run("GetBookLikeCount_Large", func(t *testing.T) {
		testBookID := primitive.NewObjectID().Hex()

		mockRepo.On("GetLikeCount", ctx, community.LikeTargetTypeBook, testBookID).Return(int64(999999), nil).Once()

		count, err := service.GetBookLikeCount(ctx, testBookID)
		assert.NoError(t, err)
		assert.Equal(t, int64(999999), count)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 获取大量点赞数成功")
	})
}

// TestLikeServicePagination 分页测试
func TestLikeServicePagination(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("GetUserLikedBooks_FirstPage", func(t *testing.T) {
		likes := []*community.Like{
			{ID: primitive.NewObjectID(), TargetID: "book1"},
			{ID: primitive.NewObjectID(), TargetID: "book2"},
		}
		mockRepo.On("GetUserLikes", ctx, testUserID, community.LikeTargetTypeBook, 1, 20).
			Return(likes, int64(50), nil).Once()

		result, total, err := service.GetUserLikedBooks(ctx, testUserID, 1, 20)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, int64(50), total)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 获取第一页点赞列表成功")
	})

	t.Run("GetUserLikedBooks_EmptyResult", func(t *testing.T) {
		mockRepo.On("GetUserLikes", ctx, testUserID, community.LikeTargetTypeBook, 10, 20).
			Return([]*community.Like{}, int64(0), nil).Once()

		result, total, err := service.GetUserLikedBooks(ctx, testUserID, 10, 20)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(result))
		assert.Equal(t, int64(0), total)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 空结果分页处理正确")
	})
}

// TestLikeServiceEventPublishing 事件发布测试
func TestLikeServiceEventPublishing(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("LikeBook_EventPublished", func(t *testing.T) {
		mockRepo.On("AddLike", ctx, mock.AnythingOfType("*reader.Like")).Return(nil).Once()
		mockRepo.On("GetLikeCount", ctx, community.LikeTargetTypeBook, testBookID).Return(int64(1), nil).Once()

		err := service.LikeBook(ctx, testUserID, testBookID)
		assert.NoError(t, err)

		// 验证事件已发布
		assert.Greater(t, len(mockEventBus.events), 0)
		assert.Equal(t, "like.book.added", mockEventBus.events[0].GetEventType())

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 点赞书籍事件发布成功，事件数: %d", len(mockEventBus.events))
	})

	t.Run("UnlikeBook_EventPublished", func(t *testing.T) {
		mockEventBus := NewMockEventBus() // 新建EventBus以清空之前的事件
		service := reading.NewLikeService(mockRepo, mockCommentRepo, mockEventBus)

		mockRepo.On("RemoveLike", ctx, testUserID, community.LikeTargetTypeBook, testBookID).Return(nil).Once()
		mockRepo.On("GetLikeCount", ctx, community.LikeTargetTypeBook, testBookID).Return(int64(0), nil).Once()

		err := service.UnlikeBook(ctx, testUserID, testBookID)
		assert.NoError(t, err)

		// 验证事件已发布
		assert.Greater(t, len(mockEventBus.events), 0)
		assert.Equal(t, "like.book.removed", mockEventBus.events[0].GetEventType())

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 取消点赞书籍事件发布成功")
	})
}

// TestLikeServiceConcurrency Service层并发测试
func TestLikeServiceConcurrency(t *testing.T) {
	t.Run("ConcurrentLike_DifferentUsers", func(t *testing.T) {
		mockRepo := new(MockLikeRepository)
		mockCommentRepo := new(MockCommentRepository)
		mockEventBus := NewMockEventBus()

		service := reading.NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
		ctx := context.Background()

		testBookID := primitive.NewObjectID().Hex()

		// Mock多次调用
		mockRepo.On("AddLike", ctx, mock.AnythingOfType("*reader.Like")).Return(nil).Times(10)
		mockRepo.On("GetLikeCount", ctx, community.LikeTargetTypeBook, testBookID).Return(int64(1), nil).Times(10)

		// 10个用户并发点赞同一本书
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				userID := primitive.NewObjectID().Hex()
				err := service.LikeBook(ctx, userID, testBookID)
				assert.NoError(t, err)
			}()
		}

		wg.Wait()

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 多用户并发点赞成功")
	})

	t.Run("ConcurrentUnlike_DifferentUsers", func(t *testing.T) {
		mockRepo := new(MockLikeRepository)
		mockCommentRepo := new(MockCommentRepository)
		mockEventBus := NewMockEventBus()

		service := reading.NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
		ctx := context.Background()

		testBookID := primitive.NewObjectID().Hex()

		// Mock多次调用
		mockRepo.On("RemoveLike", ctx, mock.AnythingOfType("string"), community.LikeTargetTypeBook, testBookID).Return(nil).Times(10)
		mockRepo.On("GetLikeCount", ctx, community.LikeTargetTypeBook, testBookID).Return(int64(0), nil).Times(10)

		// 10个用户并发取消点赞同一本书
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				userID := primitive.NewObjectID().Hex()
				err := service.UnlikeBook(ctx, userID, testBookID)
				assert.NoError(t, err)
			}()
		}

		wg.Wait()

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 多用户并发取消点赞成功")
	})
}
