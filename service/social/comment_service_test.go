package social

import (
	"Qingyu_backend/models/audit"
	"Qingyu_backend/models/social"
	"Qingyu_backend/repository/interfaces/infrastructure"
	"Qingyu_backend/service/base"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockCommentRepository Mock评论Repository
type MockCommentRepository struct {
	mock.Mock
}

func (m *MockCommentRepository) Create(ctx context.Context, comment *social.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *MockCommentRepository) GetByID(ctx context.Context, id string) (*social.Comment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*social.Comment), args.Error(1)
}

func (m *MockCommentRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockCommentRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCommentRepository) GetCommentsByBookID(ctx context.Context, bookID string, page, size int) ([]*social.Comment, int64, error) {
	args := m.Called(ctx, bookID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.Comment), args.Get(1).(int64), args.Error(2)
}

func (m *MockCommentRepository) GetCommentsByBookIDSorted(ctx context.Context, bookID string, sortBy string, page, size int) ([]*social.Comment, int64, error) {
	args := m.Called(ctx, bookID, sortBy, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.Comment), args.Get(1).(int64), args.Error(2)
}

func (m *MockCommentRepository) GetCommentsByChapterID(ctx context.Context, chapterID string, page, size int) ([]*social.Comment, int64, error) {
	args := m.Called(ctx, chapterID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.Comment), args.Get(1).(int64), args.Error(2)
}

func (m *MockCommentRepository) GetCommentsByIDs(ctx context.Context, ids []string) ([]*social.Comment, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*social.Comment), args.Error(1)
}

func (m *MockCommentRepository) GetCommentsByUserID(ctx context.Context, userID string, page, size int) ([]*social.Comment, int64, error) {
	args := m.Called(ctx, userID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.Comment), args.Get(1).(int64), args.Error(2)
}

func (m *MockCommentRepository) GetRepliesByCommentID(ctx context.Context, commentID string) ([]*social.Comment, error) {
	args := m.Called(ctx, commentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*social.Comment), args.Error(1)
}

func (m *MockCommentRepository) UpdateCommentStatus(ctx context.Context, id, status, reason string) error {
	args := m.Called(ctx, id, status, reason)
	return args.Error(0)
}

func (m *MockCommentRepository) GetPendingComments(ctx context.Context, page, size int) ([]*social.Comment, int64, error) {
	args := m.Called(ctx, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.Comment), args.Get(1).(int64), args.Error(2)
}

func (m *MockCommentRepository) IncrementLikeCount(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCommentRepository) DecrementLikeCount(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCommentRepository) IncrementReplyCount(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCommentRepository) DecrementReplyCount(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCommentRepository) GetBookRatingStats(ctx context.Context, bookID string) (map[string]interface{}, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockCommentRepository) DeleteCommentsByBookID(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockCommentRepository) GetCommentCount(ctx context.Context, bookID string) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCommentRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockSensitiveWordRepository Mock敏感词Repository
type MockSensitiveWordRepository struct {
	mock.Mock
}

func (m *MockSensitiveWordRepository) GetEnabledWords(ctx context.Context) ([]*audit.SensitiveWord, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*audit.SensitiveWord), args.Error(1)
}

func (m *MockSensitiveWordRepository) Create(ctx context.Context, word *audit.SensitiveWord) error {
	args := m.Called(ctx, word)
	return args.Error(0)
}

func (m *MockSensitiveWordRepository) BatchCreate(ctx context.Context, words []*audit.SensitiveWord) error {
	args := m.Called(ctx, words)
	return args.Error(0)
}

func (m *MockSensitiveWordRepository) BatchDelete(ctx context.Context, ids []string) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockSensitiveWordRepository) GetByID(ctx context.Context, id string) (*audit.SensitiveWord, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*audit.SensitiveWord), args.Error(1)
}

func (m *MockSensitiveWordRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockSensitiveWordRepository) BatchUpdate(ctx context.Context, ids []string, updates map[string]interface{}) error {
	args := m.Called(ctx, ids, updates)
	return args.Error(0)
}

func (m *MockSensitiveWordRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSensitiveWordRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*audit.SensitiveWord, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*audit.SensitiveWord), args.Error(1)
}

func (m *MockSensitiveWordRepository) GetByWord(ctx context.Context, word string) (*audit.SensitiveWord, error) {
	args := m.Called(ctx, word)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*audit.SensitiveWord), args.Error(1)
}

func (m *MockSensitiveWordRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSensitiveWordRepository) FindWithPagination(ctx context.Context, filter infrastructure.Filter, pagination infrastructure.Pagination) (*infrastructure.PagedResult[audit.SensitiveWord], error) {
	args := m.Called(ctx, filter, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*infrastructure.PagedResult[audit.SensitiveWord]), args.Error(1)
}

func (m *MockSensitiveWordRepository) GetByCategory(ctx context.Context, category string) ([]*audit.SensitiveWord, error) {
	args := m.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*audit.SensitiveWord), args.Error(1)
}

func (m *MockSensitiveWordRepository) GetByLevel(ctx context.Context, minLevel int) ([]*audit.SensitiveWord, error) {
	args := m.Called(ctx, minLevel)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*audit.SensitiveWord), args.Error(1)
}

func (m *MockSensitiveWordRepository) CountByCategory(ctx context.Context) (map[string]int64, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *MockSensitiveWordRepository) CountByLevel(ctx context.Context) (map[int]int64, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[int]int64), args.Error(1)
}

func (m *MockSensitiveWordRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockEventBus Mock事件总线
type MockEventBus struct {
	events []base.Event
}

func NewMockEventBus() *MockEventBus {
	return &MockEventBus{
		events: make([]base.Event, 0),
	}
}

func (m *MockEventBus) Subscribe(eventType string, handler base.EventHandler) error {
	return nil
}

func (m *MockEventBus) Unsubscribe(eventType string, handlerName string) error {
	return nil
}

func (m *MockEventBus) Publish(ctx context.Context, event base.Event) error {
	m.events = append(m.events, event)
	return nil
}

func (m *MockEventBus) PublishAsync(ctx context.Context, event base.Event) error {
	m.events = append(m.events, event)
	return nil
}

func (m *MockEventBus) GetServiceName() string {
	return "MockEventBus"
}

func (m *MockEventBus) GetVersion() string {
	return "1.0.0"
}

func (m *MockEventBus) Initialize(ctx context.Context) error {
	return nil
}

// TestCommentService_PublishComment 发表评论测试
func TestCommentService_PublishComment(t *testing.T) {
	mockRepo := new(MockCommentRepository)
	mockSensitiveRepo := new(MockSensitiveWordRepository)
	mockEventBus := NewMockEventBus()

	service := NewCommentService(mockRepo, mockSensitiveRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("PublishComment_Success", func(t *testing.T) {
		// Arrange - Mock敏感词检测（通过）
		mockSensitiveRepo.On("GetEnabledWords", ctx).
			Return([]*audit.SensitiveWord{}, nil).Once()

		// Mock创建评论
		mockRepo.On("Create", ctx, mock.MatchedBy(func(c *social.Comment) bool {
			return c.AuthorID == testUserID && c.TargetID == testBookID
		})).
			Return(nil).Once()

		// Act - 发表评论
		comment, err := service.PublishComment(ctx, testUserID, testBookID, "", "这本书真不错", 5)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, comment)
		assert.Equal(t, social.CommentStateNormal, comment.State)
		assert.Equal(t, 5, comment.Rating)

		mockRepo.AssertExpectations(t)
		mockSensitiveRepo.AssertExpectations(t)

		t.Logf("✓ 发表评论成功")
	})

	t.Run("PublishComment_WithSensitiveWord", func(t *testing.T) {
		// Arrange - Mock敏感词检测（不通过）
		sensitiveWords := []*audit.SensitiveWord{
			{Word: "敏感词", Level: 2, IsEnabled: true},
		}
		mockSensitiveRepo.On("GetEnabledWords", ctx).
			Return(sensitiveWords, nil).Once()

		// Mock创建评论
		mockRepo.On("Create", ctx, mock.AnythingOfType("*social.Comment")).
			Return(nil).Once()

		// Act - 发表评论
		comment, err := service.PublishComment(ctx, testUserID, testBookID, "", "包含敏感词的内容", 4)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, comment)
		assert.NotEqual(t, social.CommentStateNormal, comment.State)
		assert.NotEmpty(t, comment.RejectReason)

		mockRepo.AssertExpectations(t)
		mockSensitiveRepo.AssertExpectations(t)

		t.Logf("✓ 敏感词检测通过，评论被拒绝")
	})

	t.Run("PublishComment_EmptyContent", func(t *testing.T) {
		// Act - 空内容
		_, err := service.PublishComment(ctx, testUserID, testBookID, "", "", 5)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "评论内容长度")

		t.Logf("✓ 空内容验证通过")
	})

	t.Run("PublishComment_InvalidRating", func(t *testing.T) {
		// Act - 无效评分
		_, err := service.PublishComment(ctx, testUserID, testBookID, "", "这是一条测试评论这是一条测试评论", 6)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "评分")

		t.Logf("✓ 评分范围验证通过")
	})
}

// TestCommentService_ReplyComment 回复评论测试
func TestCommentService_ReplyComment(t *testing.T) {
	mockRepo := new(MockCommentRepository)
	mockSensitiveRepo := new(MockSensitiveWordRepository)
	mockEventBus := NewMockEventBus()

	service := NewCommentService(mockRepo, mockSensitiveRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testCommentID := primitive.NewObjectID().Hex()

	t.Run("ReplyComment_Success", func(t *testing.T) {
		// Arrange - Mock获取父评论
		parentComment := &social.Comment{
			IdentifiedEntity: social.IdentifiedEntity{ID: testCommentID},
			AuthorID:         primitive.NewObjectID().Hex(),
			TargetID:         primitive.NewObjectID().Hex(),
			State:            social.CommentStateNormal,
		}
		mockRepo.On("GetByID", ctx, testCommentID).
			Return(parentComment, nil).Once()

		// Mock敏感词检测
		mockSensitiveRepo.On("GetEnabledWords", ctx).
			Return([]*audit.SensitiveWord{}, nil).Once()

		// Mock创建回复
		mockRepo.On("Create", ctx, mock.AnythingOfType("*social.Comment")).
			Return(nil).Once()

		// Mock增加回复数
		mockRepo.On("IncrementReplyCount", ctx, testCommentID).
			Return(nil).Once()

		// Act - 回复评论
		reply, err := service.ReplyComment(ctx, testUserID, testCommentID, "回复内容")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, reply)
		if reply.ReplyToCommentID != nil {
			assert.Equal(t, testCommentID, *reply.ReplyToCommentID)
		}

		mockRepo.AssertExpectations(t)
		mockSensitiveRepo.AssertExpectations(t)

		t.Logf("✓ 回复评论成功")
	})

	t.Run("ReplyComment_ParentNotFound", func(t *testing.T) {
		// Arrange
		mockRepo.On("GetByID", ctx, testCommentID).
			Return(nil, errors.New("comment not found")).Once()

		// Act
		_, err := service.ReplyComment(ctx, testUserID, testCommentID, "回复内容")

		// Assert
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 父评论不存在检测通过")
	})
}

// TestCommentService_UpdateComment 更新评论测试
func TestCommentService_UpdateComment(t *testing.T) {
	mockRepo := new(MockCommentRepository)
	mockSensitiveRepo := new(MockSensitiveWordRepository)
	mockEventBus := NewMockEventBus()

	service := NewCommentService(mockRepo, mockSensitiveRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testCommentID := primitive.NewObjectID().Hex()

	t.Run("UpdateComment_Success", func(t *testing.T) {
		// Arrange - Mock获取评论
		comment := &social.Comment{
			IdentifiedEntity: social.IdentifiedEntity{ID: testCommentID},
			AuthorID:         testUserID,
			Timestamps:       social.Timestamps{CreatedAt: time.Now().Add(-10 * time.Minute)},
			State:            social.CommentStateNormal,
		}
		mockRepo.On("GetByID", ctx, testCommentID).
			Return(comment, nil).Once()

		// Mock敏感词检测
		mockSensitiveRepo.On("GetEnabledWords", ctx).
			Return([]*audit.SensitiveWord{}, nil).Once()

		// Mock更新
		mockRepo.On("Update", ctx, testCommentID, mock.AnythingOfType("map[string]interface {}")).
			Return(nil).Once()

		// Act - 更新评论
		err := service.UpdateComment(ctx, testUserID, testCommentID, "更新后的内容")

		// Assert
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockSensitiveRepo.AssertExpectations(t)

		t.Logf("✓ 更新评论成功")
	})

	t.Run("UpdateComment_NotOwner", func(t *testing.T) {
		// Arrange - Mock获取评论（不是自己的）
		comment := &social.Comment{
			IdentifiedEntity: social.IdentifiedEntity{ID: testCommentID},
			AuthorID:         primitive.NewObjectID().Hex(),
			Timestamps:       social.Timestamps{CreatedAt: time.Now()},
		}
		mockRepo.On("GetByID", ctx, testCommentID).
			Return(comment, nil).Once()

		// Act
		err := service.UpdateComment(ctx, testUserID, testCommentID, "更新内容更新内容更新内容")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "权限")

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 权限检查通过")
	})
}

// TestCommentService_DeleteComment 删除评论测试
func TestCommentService_DeleteComment(t *testing.T) {
	mockRepo := new(MockCommentRepository)
	mockSensitiveRepo := new(MockSensitiveWordRepository)
	mockEventBus := NewMockEventBus()

	service := NewCommentService(mockRepo, mockSensitiveRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testCommentID := primitive.NewObjectID().Hex()

	t.Run("DeleteComment_Success", func(t *testing.T) {
		// Arrange - Mock获取评论
		comment := &social.Comment{
			IdentifiedEntity: social.IdentifiedEntity{ID: testCommentID},
			AuthorID:         testUserID,
		}
		mockRepo.On("GetByID", ctx, testCommentID).
			Return(comment, nil).Once()

		// Mock删除
		mockRepo.On("Delete", ctx, testCommentID).
			Return(nil).Once()

		// Act - 删除评论
		err := service.DeleteComment(ctx, testUserID, testCommentID)

		// Assert
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 删除评论成功")
	})

	t.Run("DeleteComment_NotOwner", func(t *testing.T) {
		// Arrange - Mock获取评论（不是自己的）
		comment := &social.Comment{
			IdentifiedEntity: social.IdentifiedEntity{ID: testCommentID},
			AuthorID:         primitive.NewObjectID().Hex(),
		}
		mockRepo.On("GetByID", ctx, testCommentID).
			Return(comment, nil).Once()

		// Act
		err := service.DeleteComment(ctx, testUserID, testCommentID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "权限")

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 删除权限检查通过")
	})
}

// TestCommentService_GetCommentList 获取评论列表测试
func TestCommentService_GetCommentList(t *testing.T) {
	mockRepo := new(MockCommentRepository)
	mockSensitiveRepo := new(MockSensitiveWordRepository)
	mockEventBus := NewMockEventBus()

	service := NewCommentService(mockRepo, mockSensitiveRepo, mockEventBus)
	ctx := context.Background()

	testBookID := primitive.NewObjectID().Hex()

	t.Run("GetCommentList_Success", func(t *testing.T) {
		// Arrange - Mock查询
		comments := []*social.Comment{
			{IdentifiedEntity: social.IdentifiedEntity{ID: primitive.NewObjectID().Hex()}, Content: "评论1"},
			{IdentifiedEntity: social.IdentifiedEntity{ID: primitive.NewObjectID().Hex()}, Content: "评论2"},
		}
		mockRepo.On("GetCommentsByBookIDSorted", ctx, testBookID, "latest", 1, 20).
			Return(comments, int64(2), nil).Once()

		// Act - 获取列表
		result, total, err := service.GetCommentList(ctx, testBookID, "latest", 1, 20)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, int64(2), total)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 获取评论列表成功")
	})
}
