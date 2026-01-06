package comment

import (
	"Qingyu_backend/models/community"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/audit"
	"Qingyu_backend/service/reader"
)

// TestCommentServiceBusinessRules 业务规则测试
func TestCommentServiceBusinessRules(t *testing.T) {
	ctx := context.Background()
	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("ContentLength_Minimum", func(t *testing.T) {
		// 为这个子测试创建独立的Mock
		mockRepo := new(MockCommentRepository)
		mockSensitiveRepo := new(MockSensitiveWordRepository)
		mockEventBus := NewMockEventBus()
		service := reading.NewCommentService(mockRepo, mockSensitiveRepo, mockEventBus)

		// 测试最小长度（10字符）
		minContent := strings.Repeat("测", 5) // 10个字符
		mockSensitiveRepo.On("GetEnabledWords", ctx).Return([]*audit.SensitiveWord{}, nil).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*reader.Comment")).Return(nil).Once()

		comment, err := service.PublishComment(ctx, testUserID, testBookID, "", minContent, 3)
		assert.NoError(t, err)
		assert.NotNil(t, comment)

		t.Logf("✓ 最小长度内容验证通过")
	})

	t.Run("ContentLength_Maximum", func(t *testing.T) {
		mockRepo := new(MockCommentRepository)
		mockSensitiveRepo := new(MockSensitiveWordRepository)
		mockEventBus := NewMockEventBus()
		service := reading.NewCommentService(mockRepo, mockSensitiveRepo, mockEventBus)

		// 测试最大长度（500字节）- 使用ASCII字符
		longContent := strings.Repeat("a", 500) // 500个字节
		mockSensitiveRepo.On("GetEnabledWords", ctx).Return([]*audit.SensitiveWord{}, nil).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*reader.Comment")).Return(nil).Once()

		comment, err := service.PublishComment(ctx, testUserID, testBookID, "", longContent, 5)
		assert.NoError(t, err)
		assert.NotNil(t, comment)

		t.Logf("✓ 最大长度内容验证通过")
	})

	t.Run("Rating_BoundaryValues", func(t *testing.T) {
		mockRepo := new(MockCommentRepository)
		mockSensitiveRepo := new(MockSensitiveWordRepository)
		mockEventBus := NewMockEventBus()
		service := reading.NewCommentService(mockRepo, mockSensitiveRepo, mockEventBus)

		// 测试0分
		mockSensitiveRepo.On("GetEnabledWords", ctx).Return([]*audit.SensitiveWord{}, nil)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*reader.Comment")).Return(nil)

		comment, err := service.PublishComment(ctx, testUserID, testBookID, "", "这是一条测试评论这是一条测试评论", 0)
		assert.NoError(t, err)
		assert.Equal(t, 0, comment.Rating)

		// 测试5分
		comment, err = service.PublishComment(ctx, testUserID, testBookID, "", "这是一条测试评论这是一条测试评论", 5)
		assert.NoError(t, err)
		assert.Equal(t, 5, comment.Rating)

		// 测试无效评分
		_, err = service.PublishComment(ctx, testUserID, testBookID, "", "这是一条测试评论这是一条测试评论", -1)
		assert.Error(t, err)

		_, err = service.PublishComment(ctx, testUserID, testBookID, "", "这是一条测试评论这是一条测试评论", 6)
		assert.Error(t, err)

		t.Logf("✓ 评分边界值验证通过")
	})
}

// TestCommentServiceSensitiveWord 敏感词检测测试
func TestCommentServiceSensitiveWord(t *testing.T) {
	mockRepo := new(MockCommentRepository)
	mockSensitiveRepo := new(MockSensitiveWordRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewCommentService(mockRepo, mockSensitiveRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("SensitiveWord_HighLevel", func(t *testing.T) {
		// 高级别敏感词（直接拒绝）
		sensitiveWords := []*audit.SensitiveWord{
			{Word: "严重敏感词", Level: 3, IsEnabled: true},
		}
		mockSensitiveRepo.On("GetEnabledWords", ctx).Return(sensitiveWords, nil).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*reader.Comment")).Return(nil).Once()

		comment, err := service.PublishComment(ctx, testUserID, testBookID, "", "包含严重敏感词的内容", 4)
		assert.NoError(t, err)
		assert.Equal(t, "rejected", comment.Status)
		assert.Contains(t, comment.RejectReason, "敏感词")

		mockRepo.AssertExpectations(t)
		mockSensitiveRepo.AssertExpectations(t)

		t.Logf("✓ 高级别敏感词拒绝成功")
	})

	t.Run("SensitiveWord_LowLevel", func(t *testing.T) {
		// 低级别敏感词（待审核）
		sensitiveWords := []*audit.SensitiveWord{
			{Word: "轻微敏感词", Level: 1, IsEnabled: true},
		}
		mockSensitiveRepo.On("GetEnabledWords", ctx).Return(sensitiveWords, nil).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*reader.Comment")).Return(nil).Once()

		comment, err := service.PublishComment(ctx, testUserID, testBookID, "", "包含轻微敏感词的内容", 4)
		assert.NoError(t, err)
		// 低级别敏感词可能进入pending状态
		assert.Contains(t, []string{"pending", "rejected"}, comment.Status)

		mockRepo.AssertExpectations(t)
		mockSensitiveRepo.AssertExpectations(t)

		t.Logf("✓ 低级别敏感词检测成功")
	})

	t.Run("SensitiveWord_MultipleWords", func(t *testing.T) {
		// 多个敏感词
		sensitiveWords := []*audit.SensitiveWord{
			{Word: "敏感词1", Level: 2, IsEnabled: true},
			{Word: "敏感词2", Level: 2, IsEnabled: true},
		}
		mockSensitiveRepo.On("GetEnabledWords", ctx).Return(sensitiveWords, nil).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*reader.Comment")).Return(nil).Once()

		comment, err := service.PublishComment(ctx, testUserID, testBookID, "", "包含敏感词1和敏感词2的内容", 4)
		assert.NoError(t, err)
		assert.NotEqual(t, "approved", comment.Status)

		mockRepo.AssertExpectations(t)
		mockSensitiveRepo.AssertExpectations(t)

		t.Logf("✓ 多敏感词检测成功")
	})
}

// TestCommentServiceReplyChain 回复链测试
func TestCommentServiceReplyChain(t *testing.T) {
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testCommentID := primitive.NewObjectID().Hex()
	testRootID := primitive.NewObjectID().Hex()

	t.Run("Reply_ToRootComment", func(t *testing.T) {
		// 为每个子测试创建独立的Mock
		mockRepo := new(MockCommentRepository)
		mockSensitiveRepo := new(MockSensitiveWordRepository)
		mockEventBus := NewMockEventBus()

		service := reading.NewCommentService(mockRepo, mockSensitiveRepo, mockEventBus)

		// 回复根评论 - 使用实际的testCommentID
		localCommentID := primitive.NewObjectID().Hex()
		objID, _ := primitive.ObjectIDFromHex(localCommentID)
		parentComment := &community.Comment{
			ID:       objID,
			UserID:   primitive.NewObjectID().Hex(),
			BookID:   primitive.NewObjectID().Hex(),
			Status:   "approved",
			ParentID: "",
			RootID:   "",
		}
		mockRepo.On("GetByID", ctx, localCommentID).Return(parentComment, nil).Once()
		mockSensitiveRepo.On("GetEnabledWords", ctx).Return([]*audit.SensitiveWord{}, nil).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*reader.Comment")).Return(nil).Once()
		mockRepo.On("IncrementReplyCount", ctx, localCommentID).Return(nil).Once()

		reply, err := service.ReplyComment(ctx, testUserID, localCommentID, "回复根评论的内容")
		assert.NoError(t, err)
		assert.NotNil(t, reply)
		assert.Equal(t, localCommentID, reply.ParentID)
		assert.Equal(t, localCommentID, reply.RootID) // RootID应该等于ParentID

		mockRepo.AssertExpectations(t)
		mockSensitiveRepo.AssertExpectations(t)

		t.Logf("✓ 回复根评论成功")
	})

	t.Run("Reply_ToNestedComment", func(t *testing.T) {
		// 为每个子测试创建独立的Mock
		mockRepo := new(MockCommentRepository)
		mockSensitiveRepo := new(MockSensitiveWordRepository)
		mockEventBus := NewMockEventBus()

		service := reading.NewCommentService(mockRepo, mockSensitiveRepo, mockEventBus)

		// 回复嵌套评论
		nestedComment := &community.Comment{
			ID:       primitive.NewObjectID(),
			UserID:   primitive.NewObjectID().Hex(),
			BookID:   primitive.NewObjectID().Hex(),
			Status:   "approved",
			ParentID: testCommentID,
			RootID:   testRootID,
		}
		mockRepo.On("GetByID", ctx, testCommentID).Return(nestedComment, nil).Once()
		mockSensitiveRepo.On("GetEnabledWords", ctx).Return([]*audit.SensitiveWord{}, nil).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*reader.Comment")).Return(nil).Once()
		mockRepo.On("IncrementReplyCount", ctx, testCommentID).Return(nil).Once()

		reply, err := service.ReplyComment(ctx, testUserID, testCommentID, "回复嵌套评论的内容")
		assert.NoError(t, err)
		assert.NotNil(t, reply)
		assert.Equal(t, testCommentID, reply.ParentID)
		assert.Equal(t, testRootID, reply.RootID) // RootID应该继承

		mockRepo.AssertExpectations(t)
		mockSensitiveRepo.AssertExpectations(t)

		t.Logf("✓ 回复嵌套评论成功")
	})

	t.Run("Reply_ToDeletedComment", func(t *testing.T) {
		// 为每个子测试创建独立的Mock
		mockRepo := new(MockCommentRepository)
		mockSensitiveRepo := new(MockSensitiveWordRepository)
		mockEventBus := NewMockEventBus()

		service := reading.NewCommentService(mockRepo, mockSensitiveRepo, mockEventBus)

		// 回复已删除评论
		deletedComment := &community.Comment{
			ID:     primitive.NewObjectID(),
			Status: "deleted",
		}
		mockRepo.On("GetByID", ctx, testCommentID).Return(deletedComment, nil).Once()

		_, err := service.ReplyComment(ctx, testUserID, testCommentID, "回复已删除评论的内容")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "无法回复")

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 回复已删除评论拒绝成功")
	})
}

// TestCommentServiceStatistics 统计功能测试
func TestCommentServiceStatistics(t *testing.T) {
	ctx := context.Background()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("GetBookCommentStats_Success", func(t *testing.T) {
		mockRepo := new(MockCommentRepository)
		mockSensitiveRepo := new(MockSensitiveWordRepository)
		mockEventBus := NewMockEventBus()
		service := reading.NewCommentService(mockRepo, mockSensitiveRepo, mockEventBus)

		stats := map[string]interface{}{
			"total_count": int64(100),
			"average":     4.5,
		}
		mockRepo.On("GetBookRatingStats", ctx, testBookID).Return(stats, nil).Once()
		mockRepo.On("GetCommentCount", ctx, testBookID).Return(int64(100), nil).Once()

		result, err := service.GetBookCommentStats(ctx, testBookID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(100), result["comment_count"])

		t.Logf("✓ 获取书籍评论统计成功")
	})
}

// TestCommentServiceError 错误处理测试
func TestCommentServiceError(t *testing.T) {
	mockRepo := new(MockCommentRepository)
	mockSensitiveRepo := new(MockSensitiveWordRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewCommentService(mockRepo, mockSensitiveRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("SensitiveWordCheck_RepositoryError", func(t *testing.T) {
		mockSensitiveRepo.On("GetEnabledWords", ctx).Return(nil, errors.New("database error")).Once()

		_, err := service.PublishComment(ctx, testUserID, testBookID, "", "正常的评论内容", 5)
		assert.Error(t, err)

		mockSensitiveRepo.AssertExpectations(t)

		t.Logf("✓ 敏感词检测失败处理正确")
	})

	t.Run("CreateComment_RepositoryError", func(t *testing.T) {
		mockSensitiveRepo.On("GetEnabledWords", ctx).Return([]*audit.SensitiveWord{}, nil).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*reader.Comment")).Return(errors.New("insert error")).Once()

		_, err := service.PublishComment(ctx, testUserID, testBookID, "", "正常的评论内容", 5)
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
		mockSensitiveRepo.AssertExpectations(t)

		t.Logf("✓ 创建评论失败处理正确")
	})

	t.Run("EmptyParameters", func(t *testing.T) {
		// 空用户ID
		_, err := service.PublishComment(ctx, "", testBookID, "", "正常的评论内容", 5)
		assert.Error(t, err)

		// 空书籍ID
		_, err = service.PublishComment(ctx, testUserID, "", "", "正常的评论内容", 5)
		assert.Error(t, err)

		t.Logf("✓ 空参数验证通过")
	})
}

// TestCommentServiceEventPublishing 事件发布测试
func TestCommentServiceEventPublishing(t *testing.T) {
	mockRepo := new(MockCommentRepository)
	mockSensitiveRepo := new(MockSensitiveWordRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewCommentService(mockRepo, mockSensitiveRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("PublishComment_EventTriggered", func(t *testing.T) {
		mockSensitiveRepo.On("GetEnabledWords", ctx).Return([]*audit.SensitiveWord{}, nil).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*reader.Comment")).Return(nil).Once()

		_, err := service.PublishComment(ctx, testUserID, testBookID, "", "测试事件发布的评论内容", 5)
		assert.NoError(t, err)

		// 验证事件已发布
		assert.Greater(t, len(mockEventBus.events), 0)
		assert.Equal(t, "comment.created", mockEventBus.events[0].GetEventType())

		mockRepo.AssertExpectations(t)
		mockSensitiveRepo.AssertExpectations(t)

		t.Logf("✓ 评论创建事件发布成功，事件数: %d", len(mockEventBus.events))
	})
}

// TestCommentServiceConcurrency Service层并发测试
func TestCommentServiceConcurrency(t *testing.T) {
	t.Run("ConcurrentPublish_MultipleUsers", func(t *testing.T) {
		mockRepo := new(MockCommentRepository)
		mockSensitiveRepo := new(MockSensitiveWordRepository)
		mockEventBus := NewMockEventBus()

		service := reading.NewCommentService(mockRepo, mockSensitiveRepo, mockEventBus)
		ctx := context.Background()

		testBookID := primitive.NewObjectID().Hex()

		// Mock多次调用
		mockSensitiveRepo.On("GetEnabledWords", ctx).Return([]*audit.SensitiveWord{}, nil).Times(5)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*reader.Comment")).Return(nil).Times(5)

		// 并发发布评论
		done := make(chan bool, 5)
		for i := 0; i < 5; i++ {
			go func(idx int) {
				userID := primitive.NewObjectID().Hex()
				_, err := service.PublishComment(ctx, userID, testBookID, "", "并发测试评论"+string(rune(idx)), 5)
				assert.NoError(t, err)
				done <- true
			}(i)
		}

		// 等待所有goroutine完成
		for i := 0; i < 5; i++ {
			<-done
		}

		mockRepo.AssertExpectations(t)
		mockSensitiveRepo.AssertExpectations(t)

		t.Logf("✓ 并发发布评论成功")
	})
}
