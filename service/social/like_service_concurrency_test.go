// +build !race

package social

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/social"
	"Qingyu_backend/service/messaging"
)

// TestLikeServiceConcurrency Service层并发测试
// 注意：此文件在race模式下不会被编译（因为+build !race）
// mock对象不是线程安全的，所以并发测试只能在非race模式下运行
func TestLikeServiceConcurrency(t *testing.T) {
	t.Run("ConcurrentLike_DifferentUsers", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockLikeRepository)
		mockCommentRepo := new(MockCommentRepository)
		mockEventBus := messaging.NewMockEventBus()

		service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
		ctx := context.Background()

		testBookID := primitive.NewObjectID().Hex()

		// Mock多次调用
		mockRepo.On("AddLike", ctx, mock.MatchedBy(func(l *social.Like) bool {
			return l.TargetType == social.LikeTargetTypeBook
		})).Return(nil).Times(10)
		mockRepo.On("GetLikeCount", ctx, social.LikeTargetTypeBook, testBookID).Return(int64(1), nil).Times(10)

		// Act - 10个用户并发点赞同一本书
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
		// Arrange
		mockRepo := new(MockLikeRepository)
		mockCommentRepo := new(MockCommentRepository)
		mockEventBus := messaging.NewMockEventBus()

		service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
		ctx := context.Background()

		testBookID := primitive.NewObjectID().Hex()

		// Mock多次调用
		mockRepo.On("RemoveLike", ctx, mock.AnythingOfType("string"), social.LikeTargetTypeBook, testBookID).Return(nil).Times(10)
		mockRepo.On("GetLikeCount", ctx, social.LikeTargetTypeBook, testBookID).Return(int64(0), nil).Times(10)

		// Act - 10个用户并发取消点赞同一本书
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
