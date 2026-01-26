package social

import (
	"Qingyu_backend/models/social"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockLikeRepository Mock点赞Repository
type MockLikeRepository struct {
	mock.Mock
}

func (m *MockLikeRepository) AddLike(ctx context.Context, like *social.Like) error {
	args := m.Called(ctx, like)
	return args.Error(0)
}

func (m *MockLikeRepository) RemoveLike(ctx context.Context, userID, targetType, targetID string) error {
	args := m.Called(ctx, userID, targetType, targetID)
	return args.Error(0)
}

func (m *MockLikeRepository) IsLiked(ctx context.Context, userID, targetType, targetID string) (bool, error) {
	args := m.Called(ctx, userID, targetType, targetID)
	return args.Bool(0), args.Error(1)
}

func (m *MockLikeRepository) GetLikeCount(ctx context.Context, targetType, targetID string) (int64, error) {
	args := m.Called(ctx, targetType, targetID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockLikeRepository) GetUserLikes(ctx context.Context, userID, targetType string, page, size int) ([]*social.Like, int64, error) {
	args := m.Called(ctx, userID, targetType, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.Like), args.Get(1).(int64), args.Error(2)
}

func (m *MockLikeRepository) CountUserLikes(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockLikeRepository) CountTargetLikes(ctx context.Context, targetType, targetID string) (int64, error) {
	args := m.Called(ctx, targetType, targetID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockLikeRepository) GetLikesCountBatch(ctx context.Context, targetType string, targetIDs []string) (map[string]int64, error) {
	args := m.Called(ctx, targetType, targetIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *MockLikeRepository) GetUserLikeStatusBatch(ctx context.Context, userID, targetType string, targetIDs []string) (map[string]bool, error) {
	args := m.Called(ctx, userID, targetType, targetIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]bool), args.Error(1)
}

func (m *MockLikeRepository) GetByID(ctx context.Context, id string) (*social.Like, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*social.Like), args.Error(1)
}

func (m *MockLikeRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// TestLikeServiceBusinessRules 业务规则测试
func TestLikeServiceBusinessRules(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("Idempotency_LikeBook", func(t *testing.T) {
		// Arrange - 测试点赞幂等性
		testBookID := primitive.NewObjectID().Hex()

		// 第一次点赞成功
		mockRepo.On("AddLike", ctx, mock.MatchedBy(func(l *social.Like) bool {
			return l.TargetType == social.LikeTargetTypeBook
		})).Return(nil).Once()
		mockRepo.On("GetLikeCount", ctx, social.LikeTargetTypeBook, testBookID).Return(int64(1), nil).Once()

		// Act - 第一次点赞
		err := service.LikeBook(ctx, testUserID, testBookID)
		assert.NoError(t, err)

		// 第二次点赞（已存在）
		mockRepo.On("AddLike", ctx, mock.MatchedBy(func(l *social.Like) bool {
			return l.TargetType == social.LikeTargetTypeBook
		})).Return(errors.New("已经点赞过了")).Once()

		// Act - 第二次点赞
		err = service.LikeBook(ctx, testUserID, testBookID)
		assert.NoError(t, err) // 应该不报错，幂等性处理

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 点赞幂等性测试通过")
	})

	t.Run("Idempotency_UnlikeBook", func(t *testing.T) {
		// Arrange - 测试取消点赞幂等性
		testBookID := primitive.NewObjectID().Hex()

		// 第一次取消点赞成功
		mockRepo.On("RemoveLike", ctx, testUserID, social.LikeTargetTypeBook, testBookID).Return(nil).Once()
		mockRepo.On("GetLikeCount", ctx, social.LikeTargetTypeBook, testBookID).Return(int64(0), nil).Once()

		// Act - 第一次取消点赞
		err := service.UnlikeBook(ctx, testUserID, testBookID)
		assert.NoError(t, err)

		// 第二次取消点赞（不存在）
		mockRepo.On("RemoveLike", ctx, testUserID, social.LikeTargetTypeBook, testBookID).Return(errors.New("点赞记录不存在")).Once()

		// Act - 第二次取消点赞
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

	service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testCommentID := primitive.NewObjectID().Hex()

	t.Run("LikeComment_IncrementCount", func(t *testing.T) {
		// Arrange - 点赞评论应该增加评论点赞数
		mockRepo.On("AddLike", ctx, mock.MatchedBy(func(l *social.Like) bool {
			return l.TargetType == social.LikeTargetTypeComment
		})).Return(nil).Once()
		mockCommentRepo.On("IncrementLikeCount", ctx, testCommentID).Return(nil).Once()
		mockRepo.On("GetLikeCount", ctx, social.LikeTargetTypeComment, testCommentID).Return(int64(1), nil).Once()

		// Act
		err := service.LikeComment(ctx, testUserID, testCommentID)

		// Assert
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)

		t.Logf("✓ 点赞评论计数增加成功")
	})

	t.Run("UnlikeComment_DecrementCount", func(t *testing.T) {
		// Arrange - 取消点赞评论应该减少评论点赞数
		mockRepo.On("RemoveLike", ctx, testUserID, social.LikeTargetTypeComment, testCommentID).Return(nil).Once()
		mockCommentRepo.On("DecrementLikeCount", ctx, testCommentID).Return(nil).Once()
		mockRepo.On("GetLikeCount", ctx, social.LikeTargetTypeComment, testCommentID).Return(int64(0), nil).Once()

		// Act
		err := service.UnlikeComment(ctx, testUserID, testCommentID)

		// Assert
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)

		t.Logf("✓ 取消点赞评论计数减少成功")
	})

	t.Run("LikeComment_IncrementFailed", func(t *testing.T) {
		// Arrange - 点赞记录成功但增加评论计数失败
		mockRepo.On("AddLike", ctx, mock.MatchedBy(func(l *social.Like) bool {
			return l.TargetType == social.LikeTargetTypeComment
		})).Return(nil).Once()
		mockCommentRepo.On("IncrementLikeCount", ctx, testCommentID).Return(errors.New("increment failed")).Once()
		// 虽然IncrementLikeCount失败，但publishLikeEvent仍然会被调用，获取点赞数
		mockRepo.On("GetLikeCount", ctx, social.LikeTargetTypeComment, testCommentID).Return(int64(1), nil).Once()

		// Act
		err := service.LikeComment(ctx, testUserID, testCommentID)

		// Assert - 虽然IncrementLikeCount失败，但函数返回nil（因为错误被吞掉了）
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

	service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("GetBooksLikeCount_Success", func(t *testing.T) {
		// Arrange
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

		mockRepo.On("GetLikesCountBatch", ctx, social.LikeTargetTypeBook, bookIDs).
			Return(expectedCounts, nil).Once()

		// Act
		counts, err := service.GetBooksLikeCount(ctx, bookIDs)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCounts, counts)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 批量获取点赞数成功")
	})

	t.Run("GetUserLikeStatus_Success", func(t *testing.T) {
		// Arrange
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

		mockRepo.On("GetUserLikeStatusBatch", ctx, testUserID, social.LikeTargetTypeBook, bookIDs).
			Return(expectedStatus, nil).Once()

		// Act
		status, err := service.GetUserLikeStatus(ctx, testUserID, bookIDs)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, status)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 批量获取用户点赞状态成功")
	})

	t.Run("BatchOperations_EmptyArray", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockLikeRepository)
		mockCommentRepo := new(MockCommentRepository)
		mockEventBus := NewMockEventBus()

		service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
		ctx := context.Background()

		emptyIDs := []string{}

		// Act - 空数组会导致service直接返回空map，不调用repository
		counts, err := service.GetBooksLikeCount(ctx, emptyIDs)

		// Assert
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

	service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("LikeBook_TargetType", func(t *testing.T) {
		// Arrange
		testBookID := primitive.NewObjectID().Hex()

		mockRepo.On("AddLike", ctx, mock.MatchedBy(func(l *social.Like) bool {
			return l.TargetType == social.LikeTargetTypeBook && l.TargetID == testBookID
		})).Return(nil).Once()
		mockRepo.On("GetLikeCount", ctx, social.LikeTargetTypeBook, testBookID).Return(int64(1), nil).Once()

		// Act
		err := service.LikeBook(ctx, testUserID, testBookID)

		// Assert
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 书籍点赞目标类型正确")
	})

	t.Run("LikeComment_TargetType", func(t *testing.T) {
		// Arrange
		testCommentID := primitive.NewObjectID().Hex()

		mockRepo.On("AddLike", ctx, mock.MatchedBy(func(l *social.Like) bool {
			return l.TargetType == social.LikeTargetTypeComment && l.TargetID == testCommentID
		})).Return(nil).Once()
		mockCommentRepo.On("IncrementLikeCount", ctx, testCommentID).Return(nil).Once()
		mockRepo.On("GetLikeCount", ctx, social.LikeTargetTypeComment, testCommentID).Return(int64(1), nil).Once()

		// Act
		err := service.LikeComment(ctx, testUserID, testCommentID)

		// Assert
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

	service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("RepositoryError_AddLike", func(t *testing.T) {
		// Arrange
		mockRepo.On("AddLike", ctx, mock.MatchedBy(func(l *social.Like) bool {
			return l.TargetType == social.LikeTargetTypeBook
		})).
			Return(errors.New("database connection error")).Once()

		// Act
		err := service.LikeBook(ctx, testUserID, testBookID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "点赞")

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 添加点赞失败处理正确")
	})

	t.Run("RepositoryError_RemoveLike", func(t *testing.T) {
		// Arrange
		mockRepo.On("RemoveLike", ctx, testUserID, social.LikeTargetTypeBook, testBookID).
			Return(errors.New("database error")).Once()

		// Act
		err := service.UnlikeBook(ctx, testUserID, testBookID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "取消")

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 移除点赞失败处理正确")
	})

	t.Run("InvalidParameters", func(t *testing.T) {
		// Act - 空用户ID
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

	service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("GetUserLikeStats_AllTypes", func(t *testing.T) {
		// Arrange
		mockRepo.On("CountUserLikes", ctx, testUserID).Return(int64(100), nil).Once()

		// Act
		stats, err := service.GetUserLikeStats(ctx, testUserID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, int64(100), stats["total_likes"])

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 获取用户点赞统计成功")
	})

	t.Run("GetBookLikeCount_Zero", func(t *testing.T) {
		// Arrange
		testBookID := primitive.NewObjectID().Hex()

		mockRepo.On("GetLikeCount", ctx, social.LikeTargetTypeBook, testBookID).Return(int64(0), nil).Once()

		// Act
		count, err := service.GetBookLikeCount(ctx, testBookID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 获取零点赞数成功")
	})

	t.Run("GetBookLikeCount_Large", func(t *testing.T) {
		// Arrange
		testBookID := primitive.NewObjectID().Hex()

		mockRepo.On("GetLikeCount", ctx, social.LikeTargetTypeBook, testBookID).Return(int64(999999), nil).Once()

		// Act
		count, err := service.GetBookLikeCount(ctx, testBookID)

		// Assert
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

	service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("GetUserLikedBooks_FirstPage", func(t *testing.T) {
		// Arrange
		likes := []*social.Like{
			{ID: primitive.NewObjectID(), TargetID: "book1"},
			{ID: primitive.NewObjectID(), TargetID: "book2"},
		}
		mockRepo.On("GetUserLikes", ctx, testUserID, social.LikeTargetTypeBook, 1, 20).
			Return(likes, int64(50), nil).Once()

		// Act
		result, total, err := service.GetUserLikedBooks(ctx, testUserID, 1, 20)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, int64(50), total)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 获取第一页点赞列表成功")
	})

	t.Run("GetUserLikedBooks_EmptyResult", func(t *testing.T) {
		// Arrange
		mockRepo.On("GetUserLikes", ctx, testUserID, social.LikeTargetTypeBook, 10, 20).
			Return([]*social.Like{}, int64(0), nil).Once()

		// Act
		result, total, err := service.GetUserLikedBooks(ctx, testUserID, 10, 20)

		// Assert
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

	service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("LikeBook_EventPublished", func(t *testing.T) {
		// Arrange
		mockRepo.On("AddLike", ctx, mock.MatchedBy(func(l *social.Like) bool {
			return l.TargetType == social.LikeTargetTypeBook
		})).Return(nil).Once()
		mockRepo.On("GetLikeCount", ctx, social.LikeTargetTypeBook, testBookID).Return(int64(1), nil).Once()

		// Act
		err := service.LikeBook(ctx, testUserID, testBookID)

		// Assert
		assert.NoError(t, err)
		// 验证事件已发布
		assert.Greater(t, len(mockEventBus.events), 0)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 点赞书籍事件发布成功，事件数: %d", len(mockEventBus.events))
	})

	t.Run("UnlikeBook_EventPublished", func(t *testing.T) {
		// Arrange
		mockEventBus = NewMockEventBus() // 新建EventBus以清空之前的事件
		service = NewLikeService(mockRepo, mockCommentRepo, mockEventBus)

		mockRepo.On("RemoveLike", ctx, testUserID, social.LikeTargetTypeBook, testBookID).Return(nil).Once()
		mockRepo.On("GetLikeCount", ctx, social.LikeTargetTypeBook, testBookID).Return(int64(0), nil).Once()

		// Act
		err := service.UnlikeBook(ctx, testUserID, testBookID)

		// Assert
		assert.NoError(t, err)
		// 验证事件已发布
		assert.Greater(t, len(mockEventBus.events), 0)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 取消点赞书籍事件发布成功")
	})
}

// 注意：并发测试已移至 like_service_concurrency_test.go
// 该文件使用 +build !race tag，在race模式下不会被编译
// 因为mock对象不是线程安全的，并发测试只在非race模式下运行

// TestLikeServiceTableDrivenComprehensive 表格驱动综合测试
func TestLikeServiceTableDrivenComprehensive(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	tests := []struct {
		name        string
		setupMock   func()
		action      func() error
		wantErr     bool
		errContains string
		checkEvent  bool
		eventType   string
	}{
		{
			name: "成功点赞书籍",
			setupMock: func() {
				mockRepo.On("AddLike", ctx, mock.AnythingOfType("*social.Like")).Return(nil).Once()
				mockRepo.On("GetLikeCount", ctx, social.LikeTargetTypeBook, mock.Anything).Return(int64(1), nil).Once()
			},
			action: func() error {
				return service.LikeBook(ctx, testUserID, primitive.NewObjectID().Hex())
			},
			wantErr:    false,
			checkEvent: true,
			eventType:  "like.book.added",
		},
		{
			name: "成功点赞评论",
			setupMock: func() {
				mockRepo.On("AddLike", ctx, mock.AnythingOfType("*social.Like")).Return(nil).Once()
				mockCommentRepo.On("IncrementLikeCount", ctx, mock.Anything).Return(nil).Once()
				mockRepo.On("GetLikeCount", ctx, social.LikeTargetTypeComment, mock.Anything).Return(int64(1), nil).Once()
			},
			action: func() error {
				return service.LikeComment(ctx, testUserID, primitive.NewObjectID().Hex())
			},
			wantErr:    false,
			checkEvent: true,
			eventType:  "like.comment.added",
		},
		{
			name: "数据库错误",
			setupMock: func() {
				mockRepo.On("AddLike", ctx, mock.AnythingOfType("*social.Like")).Return(errors.New("db error")).Once()
			},
			action: func() error {
				return service.LikeBook(ctx, testUserID, primitive.NewObjectID().Hex())
			},
			wantErr:     true,
			errContains: "点赞",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMock()

			// Act
			err := tt.action()

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockCommentRepo.AssertExpectations(t)

			t.Logf("✓ %s", tt.name)
		})
	}
}
