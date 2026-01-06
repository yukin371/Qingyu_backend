package bookstore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/bookstore"
	bookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
	bookstoreService "Qingyu_backend/service/bookstore"
)

// MockBookDetailRepository Mock书籍详情Repository
type MockBookDetailRepository struct {
	mock.Mock
}

func (m *MockBookDetailRepository) GetByAuthorID(ctx context.Context, authorID primitive.ObjectID, limit, offset int) ([]*bookstore.BookDetail, error) {
	args := m.Called(ctx, authorID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.BookDetail), args.Error(1)
}

func (m *MockBookDetailRepository) CountByAuthorID(ctx context.Context, authorID primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, authorID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookDetailRepository) SearchByFilter(ctx context.Context, filter *bookstoreRepo.BookDetailFilter) ([]*bookstore.BookDetail, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.BookDetail), args.Error(1)
}

func (m *MockBookDetailRepository) CountByFilter(ctx context.Context, filter *bookstoreRepo.BookDetailFilter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookDetailRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore.BookDetail, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.BookDetail), args.Error(1)
}

func (m *MockBookDetailRepository) Create(ctx context.Context, detail *bookstore.BookDetail) error {
	args := m.Called(ctx, detail)
	return args.Error(0)
}

func (m *MockBookDetailRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockBookDetailRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBookDetailRepository) List(ctx context.Context, filter interface{}) ([]*bookstore.BookDetail, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.BookDetail), args.Error(1)
}

func (m *MockBookDetailRepository) Count(ctx context.Context, filter interface{}) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookDetailRepository) Exists(ctx context.Context, id primitive.ObjectID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockBookDetailRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockBookDetailRepository) BatchUpdateCategories(ctx context.Context, bookIDs []primitive.ObjectID, categoryIDs []string) error {
	args := m.Called(ctx, bookIDs, categoryIDs)
	return args.Error(0)
}

// TestCountByAuthorID 测试按作者ID统计
func TestCountByAuthorID(t *testing.T) {
	ctx := context.Background()

	t.Run("成功获取作者的书籍数量", func(t *testing.T) {
		mockRepo := new(MockBookDetailRepository)
		service := bookstoreService.NewBookDetailService(mockRepo, nil)

		authorID := primitive.NewObjectID()

		// Mock数据
		books := []*bookstore.BookDetail{
			{ID: primitive.NewObjectID(), Title: "Book 1", AuthorID: authorID},
			{ID: primitive.NewObjectID(), Title: "Book 2", AuthorID: authorID},
			{ID: primitive.NewObjectID(), Title: "Book 3", AuthorID: authorID},
		}

		mockRepo.On("GetByAuthorID", ctx, authorID, 20, 0).Return(books, nil)
		mockRepo.On("CountByAuthorID", ctx, authorID).Return(int64(3), nil)

		result, total, err := service.GetBookDetailsByAuthorID(ctx, authorID, 1, 20)

		assert.NoError(t, err)
		assert.Len(t, result, 3)
		assert.Equal(t, int64(3), total)

		mockRepo.AssertExpectations(t)
	})

	t.Run("作者没有书籍", func(t *testing.T) {
		mockRepo := new(MockBookDetailRepository)
		service := bookstoreService.NewBookDetailService(mockRepo, nil)

		authorID := primitive.NewObjectID()

		mockRepo.On("GetByAuthorID", ctx, authorID, 20, 0).Return([]*bookstore.BookDetail{}, nil)
		mockRepo.On("CountByAuthorID", ctx, authorID).Return(int64(0), nil)

		result, total, err := service.GetBookDetailsByAuthorID(ctx, authorID, 1, 20)

		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.Equal(t, int64(0), total)

		mockRepo.AssertExpectations(t)
	})

	t.Run("CountByAuthorID失败时使用降级方案", func(t *testing.T) {
		mockRepo := new(MockBookDetailRepository)
		service := bookstoreService.NewBookDetailService(mockRepo, nil)

		authorID := primitive.NewObjectID()

		books := []*bookstore.BookDetail{
			{ID: primitive.NewObjectID(), Title: "Book 1"},
			{ID: primitive.NewObjectID(), Title: "Book 2"},
		}

		mockRepo.On("GetByAuthorID", ctx, authorID, 20, 0).Return(books, nil)
		mockRepo.On("CountByAuthorID", ctx, authorID).Return(int64(0), assert.AnError)

		result, total, err := service.GetBookDetailsByAuthorID(ctx, authorID, 1, 20)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, int64(2), total) // 使用列表长度作为降级方案

		mockRepo.AssertExpectations(t)
	})

	t.Run("分页功能", func(t *testing.T) {
		mockRepo := new(MockBookDetailRepository)
		service := bookstoreService.NewBookDetailService(mockRepo, nil)

		authorID := primitive.NewObjectID()

		// 第2页数据
		books := []*bookstore.BookDetail{
			{ID: primitive.NewObjectID(), Title: "Book 21"},
			{ID: primitive.NewObjectID(), Title: "Book 22"},
		}

		mockRepo.On("GetByAuthorID", ctx, authorID, 20, 20).Return(books, nil)
		mockRepo.On("CountByAuthorID", ctx, authorID).Return(int64(25), nil)

		result, total, err := service.GetBookDetailsByAuthorID(ctx, authorID, 2, 20)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, int64(25), total)

		mockRepo.AssertExpectations(t)
	})
}

// TestSearchWithFilter 测试搜索过滤
func TestSearchWithFilter(t *testing.T) {
	ctx := context.Background()

	t.Run("按标题搜索", func(t *testing.T) {
		mockRepo := new(MockBookDetailRepository)
		service := bookstoreService.NewBookDetailService(mockRepo, nil)

		filter := &bookstoreService.BookDetailFilter{
			Title: "测试",
		}

		books := []*bookstore.BookDetail{
			{ID: primitive.NewObjectID(), Title: "测试书籍1"},
			{ID: primitive.NewObjectID(), Title: "测试书籍2"},
		}

		mockRepo.On("SearchByFilter", ctx, mock.AnythingOfType("*bookstore.BookDetailFilter")).
			Return(books, nil)
		mockRepo.On("CountByFilter", ctx, mock.AnythingOfType("*bookstore.BookDetailFilter")).
			Return(int64(2), nil)

		result, total, err := service.SearchBookDetailsWithFilter(ctx, filter, 1, 20)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, int64(2), total)

		mockRepo.AssertExpectations(t)
	})

	t.Run("按作者搜索", func(t *testing.T) {
		mockRepo := new(MockBookDetailRepository)
		service := bookstoreService.NewBookDetailService(mockRepo, nil)

		filter := &bookstoreService.BookDetailFilter{
			Author: "张三",
		}

		books := []*bookstore.BookDetail{
			{ID: primitive.NewObjectID(), Title: "书籍A", Author: "张三"},
			{ID: primitive.NewObjectID(), Title: "书籍B", Author: "张三"},
			{ID: primitive.NewObjectID(), Title: "书籍C", Author: "张三"},
		}

		mockRepo.On("SearchByFilter", ctx, mock.AnythingOfType("*bookstore.BookDetailFilter")).
			Return(books, nil)
		mockRepo.On("CountByFilter", ctx, mock.AnythingOfType("*bookstore.BookDetailFilter")).
			Return(int64(3), nil)

		result, total, err := service.SearchBookDetailsWithFilter(ctx, filter, 1, 20)

		assert.NoError(t, err)
		assert.Len(t, result, 3)
		assert.Equal(t, int64(3), total)

		mockRepo.AssertExpectations(t)
	})

	t.Run("按状态过滤", func(t *testing.T) {
		mockRepo := new(MockBookDetailRepository)
		service := bookstoreService.NewBookDetailService(mockRepo, nil)

		status := bookstore.BookStatusOngoing
		filter := &bookstoreService.BookDetailFilter{
			Status: &status,
		}

		books := []*bookstore.BookDetail{
			{ID: primitive.NewObjectID(), Status: bookstore.BookStatusOngoing},
		}

		mockRepo.On("SearchByFilter", ctx, mock.AnythingOfType("*bookstore.BookDetailFilter")).
			Return(books, nil)
		mockRepo.On("CountByFilter", ctx, mock.AnythingOfType("*bookstore.BookDetailFilter")).
			Return(int64(1), nil)

		result, _, err := service.SearchBookDetailsWithFilter(ctx, filter, 1, 20)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, bookstore.BookStatusOngoing, result[0].Status)

		mockRepo.AssertExpectations(t)
	})

	t.Run("空搜索结果", func(t *testing.T) {
		mockRepo := new(MockBookDetailRepository)
		service := bookstoreService.NewBookDetailService(mockRepo, nil)

		filter := &bookstoreService.BookDetailFilter{
			Title: "不存在的书籍",
		}

		mockRepo.On("SearchByFilter", ctx, mock.AnythingOfType("*bookstore.BookDetailFilter")).
			Return([]*bookstore.BookDetail{}, nil)
		mockRepo.On("CountByFilter", ctx, mock.AnythingOfType("*bookstore.BookDetailFilter")).
			Return(int64(0), nil)

		result, total, err := service.SearchBookDetailsWithFilter(ctx, filter, 1, 20)

		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.Equal(t, int64(0), total)

		mockRepo.AssertExpectations(t)
	})

	t.Run("排序功能", func(t *testing.T) {
		mockRepo := new(MockBookDetailRepository)
		service := bookstoreService.NewBookDetailService(mockRepo, nil)

		filter := &bookstoreService.BookDetailFilter{
			SortBy:    "created_at",
			SortOrder: "desc",
		}

		books := []*bookstore.BookDetail{
			{ID: primitive.NewObjectID(), Title: "最新书籍"},
			{ID: primitive.NewObjectID(), Title: "较新书籍"},
			{ID: primitive.NewObjectID(), Title: "旧书籍"},
		}

		mockRepo.On("SearchByFilter", ctx, mock.MatchedBy(func(f *bookstoreRepo.BookDetailFilter) bool {
			return f.SortBy == "created_at" && f.SortOrder == "desc"
		})).Return(books, nil)
		mockRepo.On("CountByFilter", ctx, mock.AnythingOfType("*bookstore.BookDetailFilter")).
			Return(int64(3), nil)

		result, total, err := service.SearchBookDetailsWithFilter(ctx, filter, 1, 20)

		assert.NoError(t, err)
		assert.Len(t, result, 3)
		assert.Equal(t, "最新书籍", result[0].Title)

		mockRepo.AssertExpectations(t)
	})

	t.Run("分页功能", func(t *testing.T) {
		mockRepo := new(MockBookDetailRepository)
		service := bookstoreService.NewBookDetailService(mockRepo, nil)

		filter := &bookstoreService.BookDetailFilter{}

		// 生成30个书籍（测试分页）
		allBooks := make([]*bookstore.BookDetail, 30)
		for i := 0; i < 30; i++ {
			allBooks[i] = &bookstore.BookDetail{
				ID:    primitive.NewObjectID(),
				Title: "Book " + string(rune(i)),
			}
		}

		mockRepo.On("SearchByFilter", ctx, mock.AnythingOfType("*bookstore.BookDetailFilter")).
			Return(allBooks, nil)
		mockRepo.On("CountByFilter", ctx, mock.AnythingOfType("*bookstore.BookDetailFilter")).
			Return(int64(30), nil)

		// 获取第2页，每页10条
		result, _, err := service.SearchBookDetailsWithFilter(ctx, filter, 2, 10)

		assert.NoError(t, err)
		assert.Len(t, result, 10) // 第2页应该有10条

		mockRepo.AssertExpectations(t)
	})
}

// BenchmarkSearchBookDetailsWithFilter 基准测试
func BenchmarkSearchBookDetailsWithFilter(b *testing.B) {
	ctx := context.Background()
	mockRepo := new(MockBookDetailRepository)
	service := bookstoreService.NewBookDetailService(mockRepo, nil)

	filter := &bookstoreService.BookDetailFilter{
		Title: "测试",
	}

	books := []*bookstore.BookDetail{
		{ID: primitive.NewObjectID(), Title: "测试书籍"},
	}

	mockRepo.On("SearchByFilter", ctx, mock.AnythingOfType("*bookstore.BookDetailFilter")).
		Return(books, nil)
	mockRepo.On("CountByFilter", ctx, mock.AnythingOfType("*bookstore.BookDetailFilter")).
		Return(int64(1), nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = service.SearchBookDetailsWithFilter(ctx, filter, 1, 20)
	}
}
