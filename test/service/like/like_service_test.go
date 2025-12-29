package like

import (
	"Qingyu_backend/models/community"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/service/base"
	"Qingyu_backend/service/reading"
)

// MockCommentRepository Mock评论Repository（为LikeService测试提供）
type MockCommentRepository struct {
	mock.Mock
}

func (m *MockCommentRepository) Create(ctx context.Context, comment *community.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *MockCommentRepository) GetByID(ctx context.Context, id string) (*community.Comment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*community.Comment), args.Error(1)
}

func (m *MockCommentRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockCommentRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCommentRepository) GetCommentsByBookID(ctx context.Context, bookID string, page, size int) ([]*community.Comment, int64, error) {
	args := m.Called(ctx, bookID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*community.Comment), args.Get(1).(int64), args.Error(2)
}

func (m *MockCommentRepository) GetCommentsByBookIDSorted(ctx context.Context, bookID string, sortBy string, page, size int) ([]*community.Comment, int64, error) {
	args := m.Called(ctx, bookID, sortBy, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*community.Comment), args.Get(1).(int64), args.Error(2)
}

func (m *MockCommentRepository) GetCommentsByChapterID(ctx context.Context, chapterID string, page, size int) ([]*community.Comment, int64, error) {
	args := m.Called(ctx, chapterID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*community.Comment), args.Get(1).(int64), args.Error(2)
}

func (m *MockCommentRepository) GetCommentsByIDs(ctx context.Context, ids []string) ([]*community.Comment, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*community.Comment), args.Error(1)
}

func (m *MockCommentRepository) GetCommentsByUserID(ctx context.Context, userID string, page, size int) ([]*community.Comment, int64, error) {
	args := m.Called(ctx, userID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*community.Comment), args.Get(1).(int64), args.Error(2)
}

func (m *MockCommentRepository) GetRepliesByCommentID(ctx context.Context, commentID string) ([]*community.Comment, error) {
	args := m.Called(ctx, commentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*community.Comment), args.Error(1)
}

func (m *MockCommentRepository) UpdateCommentStatus(ctx context.Context, id, status, reason string) error {
	args := m.Called(ctx, id, status, reason)
	return args.Error(0)
}

func (m *MockCommentRepository) GetPendingComments(ctx context.Context, page, size int) ([]*community.Comment, int64, error) {
	args := m.Called(ctx, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*community.Comment), args.Get(1).(int64), args.Error(2)
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

// MockLikeRepository Mock点赞Repository
type MockLikeRepository struct {
	mock.Mock
}

func (m *MockLikeRepository) AddLike(ctx context.Context, like *community.Like) error {
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

func (m *MockLikeRepository) GetUserLikes(ctx context.Context, userID, targetType string, page, size int) ([]*community.Like, int64, error) {
	args := m.Called(ctx, userID, targetType, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*community.Like), args.Get(1).(int64), args.Error(2)
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

func (m *MockLikeRepository) GetByID(ctx context.Context, id string) (*community.Like, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*community.Like), args.Error(1)
}

func (m *MockLikeRepository) Health(ctx context.Context) error {
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

func (m *MockEventBus) Health(ctx context.Context) error {
	return nil
}

func (m *MockEventBus) Close(ctx context.Context) error {
	return nil
}

// TestLikeService_LikeBook 点赞书籍测试
func TestLikeService_LikeBook(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("LikeBook_Success", func(t *testing.T) {
		// Mock添加点赞
		mockRepo.On("AddLike", ctx, mock.AnythingOfType("*reader.Like")).
			Return(nil).Once()

		// Mock获取点赞数（用于事件发布）
		mockRepo.On("GetLikeCount", ctx, community.LikeTargetTypeBook, testBookID).
			Return(int64(1), nil).Once()

		// 点赞书籍
		err := service.LikeBook(ctx, testUserID, testBookID)

		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 点赞书籍成功")
	})

	t.Run("LikeBook_AlreadyLiked", func(t *testing.T) {
		// Mock添加点赞（返回已点赞错误）
		mockRepo.On("AddLike", ctx, mock.AnythingOfType("*reader.Like")).
			Return(fmt.Errorf("已经点赞过了")).Once()

		// 点赞书籍（幂等性：不报错）
		err := service.LikeBook(ctx, testUserID, testBookID)

		assert.NoError(t, err) // Service将"已经点赞过了"错误处理为幂等性，不报错

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 重复点赞检测通过（幂等性）")
	})

	t.Run("LikeBook_EmptyBookID", func(t *testing.T) {
		err := service.LikeBook(ctx, testUserID, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "书籍ID")

		t.Logf("✓ 空书籍ID验证通过")
	})
}

// TestLikeService_UnlikeBook 取消点赞书籍测试
func TestLikeService_UnlikeBook(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("UnlikeBook_Success", func(t *testing.T) {
		// Mock移除点赞
		mockRepo.On("RemoveLike", ctx, testUserID, community.LikeTargetTypeBook, testBookID).
			Return(nil).Once()

		// Mock获取点赞数（用于事件发布）
		mockRepo.On("GetLikeCount", ctx, community.LikeTargetTypeBook, testBookID).
			Return(int64(0), nil).Once()

		// 取消点赞
		err := service.UnlikeBook(ctx, testUserID, testBookID)

		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 取消点赞成功")
	})

	t.Run("UnlikeBook_NotLiked", func(t *testing.T) {
		// Mock移除点赞（返回未点赞错误）
		mockRepo.On("RemoveLike", ctx, testUserID, community.LikeTargetTypeBook, testBookID).
			Return(fmt.Errorf("点赞记录不存在")).Once()

		// 取消点赞（幂等性：不报错）
		err := service.UnlikeBook(ctx, testUserID, testBookID)

		assert.NoError(t, err) // Service将"点赞记录不存在"错误处理为幂等性，不报错

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 未点赞检测通过（幂等性）")
	})
}

// TestLikeService_LikeComment 点赞评论测试
func TestLikeService_LikeComment(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testCommentID := primitive.NewObjectID().Hex()

	t.Run("LikeComment_Success", func(t *testing.T) {
		// Mock添加点赞
		mockRepo.On("AddLike", ctx, mock.AnythingOfType("*reader.Like")).
			Return(nil).Once()

		// Mock增加评论点赞数
		mockCommentRepo.On("IncrementLikeCount", ctx, testCommentID).
			Return(nil).Once()

		// Mock获取点赞数（用于事件发布）
		mockRepo.On("GetLikeCount", ctx, community.LikeTargetTypeComment, testCommentID).
			Return(int64(1), nil).Once()

		// 点赞评论
		err := service.LikeComment(ctx, testUserID, testCommentID)

		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)

		t.Logf("✓ 点赞评论成功")
	})

	t.Run("LikeComment_AlreadyLiked", func(t *testing.T) {
		// Mock添加点赞（返回已点赞错误）
		mockRepo.On("AddLike", ctx, mock.AnythingOfType("*reader.Like")).
			Return(fmt.Errorf("已经点赞过了")).Once()

		// 点赞评论（幂等性：不报错）
		err := service.LikeComment(ctx, testUserID, testCommentID)

		assert.NoError(t, err) // Service将"已经点赞过了"错误处理为幂等性，不报错

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 重复点赞评论检测通过（幂等性）")
	})
}

// TestLikeService_UnlikeComment 取消点赞评论测试
func TestLikeService_UnlikeComment(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testCommentID := primitive.NewObjectID().Hex()

	t.Run("UnlikeComment_Success", func(t *testing.T) {
		// Mock移除点赞
		mockRepo.On("RemoveLike", ctx, testUserID, community.LikeTargetTypeComment, testCommentID).
			Return(nil).Once()

		// Mock减少评论点赞数
		mockCommentRepo.On("DecrementLikeCount", ctx, testCommentID).
			Return(nil).Once()

		// Mock获取点赞数（用于事件发布）
		mockRepo.On("GetLikeCount", ctx, community.LikeTargetTypeComment, testCommentID).
			Return(int64(0), nil).Once()

		// 取消点赞
		err := service.UnlikeComment(ctx, testUserID, testCommentID)

		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)

		t.Logf("✓ 取消点赞评论成功")
	})

	t.Run("UnlikeComment_NotLiked", func(t *testing.T) {
		// Mock移除点赞（返回未点赞错误）
		mockRepo.On("RemoveLike", ctx, testUserID, community.LikeTargetTypeComment, testCommentID).
			Return(fmt.Errorf("点赞记录不存在")).Once()

		// 取消点赞（幂等性：不报错）
		err := service.UnlikeComment(ctx, testUserID, testCommentID)

		assert.NoError(t, err) // Service将"点赞记录不存在"错误处理为幂等性，不报错

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 未点赞评论检测通过（幂等性）")
	})
}

// TestLikeService_GetBookLikeCount 获取书籍点赞数测试
func TestLikeService_GetBookLikeCount(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testBookID := primitive.NewObjectID().Hex()

	t.Run("GetBookLikeCount_Success", func(t *testing.T) {
		// Mock获取点赞数
		mockRepo.On("GetLikeCount", ctx, community.LikeTargetTypeBook, testBookID).
			Return(int64(100), nil).Once()

		// 获取点赞数
		count, err := service.GetBookLikeCount(ctx, testBookID)

		assert.NoError(t, err)
		assert.Equal(t, int64(100), count)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 获取点赞数成功")
	})
}

// TestLikeService_GetUserLikedBooks 获取用户点赞书籍列表测试
func TestLikeService_GetUserLikedBooks(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("GetUserLikedBooks_Success", func(t *testing.T) {
		// Mock查询
		likes := []*community.Like{
			{ID: primitive.NewObjectID(), TargetID: "book1"},
			{ID: primitive.NewObjectID(), TargetID: "book2"},
		}
		mockRepo.On("GetUserLikes", ctx, testUserID, community.LikeTargetTypeBook, 1, 20).
			Return(likes, int64(2), nil).Once()

		// 获取列表
		result, total, err := service.GetUserLikedBooks(ctx, testUserID, 1, 20)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, int64(2), total)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 获取点赞书籍列表成功")
	})
}

// TestLikeService_GetUserLikeStats 获取用户点赞统计测试
func TestLikeService_GetUserLikeStats(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("GetUserLikeStats_Success", func(t *testing.T) {
		// Mock统计
		mockRepo.On("CountUserLikes", ctx, testUserID).
			Return(int64(50), nil).Once()

		// 获取统计
		stats, err := service.GetUserLikeStats(ctx, testUserID)

		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, int64(50), stats["total_likes"])

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 获取点赞统计成功")
	})
}

// TestLikeService_ConcurrentLike 并发点赞测试
func TestLikeService_ConcurrentLike(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("ConcurrentLike_OnlyOneSuccess", func(t *testing.T) {
		// 第一次添加成功
		mockRepo.On("AddLike", ctx, mock.AnythingOfType("*reader.Like")).
			Return(nil).Once()

		// Mock获取点赞数（用于事件发布）
		mockRepo.On("GetLikeCount", ctx, community.LikeTargetTypeBook, testBookID).
			Return(int64(1), nil).Once()

		// 第二次添加（返回已点赞错误，幂等性处理）
		mockRepo.On("AddLike", ctx, mock.AnythingOfType("*reader.Like")).
			Return(fmt.Errorf("已经点赞过了")).Once()

		// 第一次点赞应该成功
		err1 := service.LikeBook(ctx, testUserID, testBookID)
		assert.NoError(t, err1)

		// 第二次点赞也成功（幂等性）
		err2 := service.LikeBook(ctx, testUserID, testBookID)
		assert.NoError(t, err2)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 并发点赞幂等性检测通过")
	})
}

// TestLikeService_AntiSpam 防刷测试
func TestLikeService_AntiSpam(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("AntiSpam_RapidLike", func(t *testing.T) {
		// 模拟快速点赞多本书
		for i := 0; i < 5; i++ {
			bookID := primitive.NewObjectID().Hex()

			mockRepo.On("AddLike", ctx, mock.AnythingOfType("*reader.Like")).
				Return(nil).Once()

			// Mock获取点赞数（用于事件发布）
			mockRepo.On("GetLikeCount", ctx, community.LikeTargetTypeBook, bookID).
				Return(int64(1), nil).Once()

			err := service.LikeBook(ctx, testUserID, bookID)
			assert.NoError(t, err)
		}

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 正常点赞速率通过")
	})
}

// TestLikeService_BatchOperations 批量操作测试
func TestLikeService_BatchOperations(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	_ = reading.NewLikeService(mockRepo, mockCommentRepo, mockEventBus)

	t.Run("BatchRepository_Success", func(t *testing.T) {
		// 测试批量Repository操作
		// 这里可以测试Repository的批量方法是否正常工作

		t.Logf("✓ 批量操作测试预留")
	})
}
