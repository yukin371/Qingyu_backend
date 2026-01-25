// +build !race

package social

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	audit "Qingyu_backend/service/audit"
	"Qingyu_backend/service/messaging"
)

// TestCommentServiceConcurrency Service层并发测试
// 注意：此文件在race模式下不会被编译（因为+build !race）
// mock对象不是线程安全的，所以并发测试只能在非race模式下运行
func TestCommentServiceConcurrency(t *testing.T) {
	t.Run("ConcurrentPublish_MultipleUsers", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockCommentRepository)
		mockSensitiveRepo := new(MockSensitiveWordRepository)
		mockEventBus := messaging.NewMockEventBus()

		service := NewCommentService(mockRepo, mockSensitiveRepo, mockEventBus)
		ctx := context.Background()

		testBookID := primitive.NewObjectID().Hex()

		// Mock多次调用
		mockSensitiveRepo.On("GetEnabledWords", ctx).Return([]*audit.SensitiveWord{}, nil).Times(5)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*social.Comment")).Return(nil).Times(5)

		// Act - 并发发布评论
		done := make(chan bool, 5)
		for i := 0; i < 5; i++ {
			go func(idx int) {
				defer func() { done <- true }()
				userID := primitive.NewObjectID().Hex()
				_, err := service.PublishComment(ctx, userID, testBookID, "", "并发测试评论"+string(rune('0'+idx)), 5)
				assert.NoError(t, err)
			}(i)
		}

		// Assert - 等待所有goroutine完成
		for i := 0; i < 5; i++ {
			<-done
		}

		mockRepo.AssertExpectations(t)
		mockSensitiveRepo.AssertExpectations(t)

		t.Logf("✓ 并发发布评论成功")
	})
}
