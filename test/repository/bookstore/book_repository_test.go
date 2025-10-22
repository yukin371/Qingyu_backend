package bookstore_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/reading/bookstore"
	bookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
)

// MockBookRepository Mock实现
type MockBookRepository struct {
	mock.Mock
	bookstoreRepo.BookRepository // 嵌入接口避免实现所有方法
}

func (m *MockBookRepository) Create(ctx context.Context, book *bookstore.Book) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

func (m *MockBookRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore.Book, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockBookRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBookRepository) List(ctx context.Context, filter interface{}) ([]*bookstore.Book, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) GetByCategory(ctx context.Context, categoryID primitive.ObjectID, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, categoryID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) GetByStatus(ctx context.Context, status bookstore.BookStatus, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, keyword, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) GetRecommended(ctx context.Context, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) GetFeatured(ctx context.Context, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) GetHotBooks(ctx context.Context, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) CountByCategory(ctx context.Context, categoryID primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRepository) CountByStatus(ctx context.Context, status bookstore.BookStatus) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRepository) BatchUpdateStatus(ctx context.Context, bookIDs []primitive.ObjectID, status bookstore.BookStatus) error {
	args := m.Called(ctx, bookIDs, status)
	return args.Error(0)
}

func (m *MockBookRepository) IncrementViewCount(ctx context.Context, bookID primitive.ObjectID) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockBookRepository) GetStats(ctx context.Context) (*bookstore.BookStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.BookStats), args.Error(1)
}

// 测试助手函数
func createTestBook(id primitive.ObjectID, title, author string) *bookstore.Book {
	now := time.Now()
	return &bookstore.Book{
		ID:            id,
		Title:         title,
		Author:        author,
		Introduction:  "测试简介",
		Cover:         "https://example.com/cover.jpg",
		Status:        bookstore.BookStatusPublished,
		WordCount:     100000,
		ChapterCount:  100,
		Price:         19.99,
		IsFree:        false,
		IsRecommended: false,
		IsFeatured:    false,
		IsHot:         false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// TestBookRepository_Create 测试创建书籍
func TestBookRepository_Create(t *testing.T) {
	mockRepo := new(MockBookRepository)
	ctx := context.Background()

	book := createTestBook(primitive.NewObjectID(), "测试书籍", "测试作者")

	// 设置Mock期望
	mockRepo.On("Create", ctx, book).Return(nil)

	// 执行测试
	err := mockRepo.Create(ctx, book)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestBookRepository_GetByID_Success 测试根据ID获取书籍成功
func TestBookRepository_GetByID_Success(t *testing.T) {
	mockRepo := new(MockBookRepository)
	ctx := context.Background()

	bookID := primitive.NewObjectID()
	expectedBook := createTestBook(bookID, "测试书籍", "测试作者")

	// 设置Mock期望
	mockRepo.On("GetByID", ctx, bookID).Return(expectedBook, nil)

	// 执行测试
	book, err := mockRepo.GetByID(ctx, bookID)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, expectedBook.ID, book.ID)
	assert.Equal(t, expectedBook.Title, book.Title)
	mockRepo.AssertExpectations(t)
}

// TestBookRepository_GetByID_NotFound 测试根据ID获取书籍不存在
func TestBookRepository_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockBookRepository)
	ctx := context.Background()

	bookID := primitive.NewObjectID()

	// 设置Mock期望：书籍不存在
	mockRepo.On("GetByID", ctx, bookID).Return(nil, errors.New("book not found"))

	// 执行测试
	book, err := mockRepo.GetByID(ctx, bookID)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, book)
	mockRepo.AssertExpectations(t)
}

// TestBookRepository_GetByCategory 测试根据分类获取书籍列表
func TestBookRepository_GetByCategory(t *testing.T) {
	mockRepo := new(MockBookRepository)
	ctx := context.Background()

	categoryID := primitive.NewObjectID()
	books := []*bookstore.Book{
		createTestBook(primitive.NewObjectID(), "书籍1", "作者1"),
		createTestBook(primitive.NewObjectID(), "书籍2", "作者2"),
	}

	// 设置Mock期望
	mockRepo.On("GetByCategory", ctx, categoryID, 10, 0).Return(books, nil)

	// 执行测试
	result, err := mockRepo.GetByCategory(ctx, categoryID, 10, 0)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	mockRepo.AssertExpectations(t)
}

// TestBookRepository_GetByStatus 测试根据状态获取书籍列表
func TestBookRepository_GetByStatus(t *testing.T) {
	mockRepo := new(MockBookRepository)
	ctx := context.Background()

	books := []*bookstore.Book{
		createTestBook(primitive.NewObjectID(), "已发布书籍", "作者1"),
	}

	// 设置Mock期望
	mockRepo.On("GetByStatus", ctx, bookstore.BookStatusPublished, 10, 0).Return(books, nil)

	// 执行测试
	result, err := mockRepo.GetByStatus(ctx, bookstore.BookStatusPublished, 10, 0)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, bookstore.BookStatusPublished, result[0].Status)
	mockRepo.AssertExpectations(t)
}

// TestBookRepository_Search 测试搜索书籍
func TestBookRepository_Search(t *testing.T) {
	mockRepo := new(MockBookRepository)
	ctx := context.Background()

	keyword := "测试"
	books := []*bookstore.Book{
		createTestBook(primitive.NewObjectID(), "测试书籍1", "作者1"),
		createTestBook(primitive.NewObjectID(), "测试书籍2", "作者2"),
	}

	// 设置Mock期望
	mockRepo.On("Search", ctx, keyword, 10, 0).Return(books, nil)

	// 执行测试
	result, err := mockRepo.Search(ctx, keyword, 10, 0)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	mockRepo.AssertExpectations(t)
}

// TestBookRepository_GetRecommended 测试获取推荐书籍
func TestBookRepository_GetRecommended(t *testing.T) {
	mockRepo := new(MockBookRepository)
	ctx := context.Background()

	recommendedBook := createTestBook(primitive.NewObjectID(), "推荐书籍", "推荐作者")
	recommendedBook.IsRecommended = true
	books := []*bookstore.Book{recommendedBook}

	// 设置Mock期望
	mockRepo.On("GetRecommended", ctx, 10, 0).Return(books, nil)

	// 执行测试
	result, err := mockRepo.GetRecommended(ctx, 10, 0)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
	assert.True(t, result[0].IsRecommended)
	mockRepo.AssertExpectations(t)
}

// TestBookRepository_GetFeatured 测试获取精选书籍
func TestBookRepository_GetFeatured(t *testing.T) {
	mockRepo := new(MockBookRepository)
	ctx := context.Background()

	featuredBook := createTestBook(primitive.NewObjectID(), "精选书籍", "精选作者")
	featuredBook.IsFeatured = true
	books := []*bookstore.Book{featuredBook}

	// 设置Mock期望
	mockRepo.On("GetFeatured", ctx, 10, 0).Return(books, nil)

	// 执行测试
	result, err := mockRepo.GetFeatured(ctx, 10, 0)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
	assert.True(t, result[0].IsFeatured)
	mockRepo.AssertExpectations(t)
}

// TestBookRepository_GetHotBooks 测试获取热门书籍
func TestBookRepository_GetHotBooks(t *testing.T) {
	mockRepo := new(MockBookRepository)
	ctx := context.Background()

	hotBook := createTestBook(primitive.NewObjectID(), "热门书籍", "热门作者")
	hotBook.IsHot = true
	books := []*bookstore.Book{hotBook}

	// 设置Mock期望
	mockRepo.On("GetHotBooks", ctx, 10, 0).Return(books, nil)

	// 执行测试
	result, err := mockRepo.GetHotBooks(ctx, 10, 0)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
	assert.True(t, result[0].IsHot)
	mockRepo.AssertExpectations(t)
}

// TestBookRepository_CountByCategory 测试统计分类下的书籍数量
func TestBookRepository_CountByCategory(t *testing.T) {
	mockRepo := new(MockBookRepository)
	ctx := context.Background()

	categoryID := primitive.NewObjectID()
	expectedCount := int64(25)

	// 设置Mock期望
	mockRepo.On("CountByCategory", ctx, categoryID).Return(expectedCount, nil)

	// 执行测试
	count, err := mockRepo.CountByCategory(ctx, categoryID)

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	mockRepo.AssertExpectations(t)
}

// TestBookRepository_CountByStatus 测试统计指定状态的书籍数量
func TestBookRepository_CountByStatus(t *testing.T) {
	mockRepo := new(MockBookRepository)
	ctx := context.Background()

	expectedCount := int64(150)

	// 设置Mock期望
	mockRepo.On("CountByStatus", ctx, bookstore.BookStatusPublished).Return(expectedCount, nil)

	// 执行测试
	count, err := mockRepo.CountByStatus(ctx, bookstore.BookStatusPublished)

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	mockRepo.AssertExpectations(t)
}

// TestBookRepository_BatchUpdateStatus 测试批量更新书籍状态
func TestBookRepository_BatchUpdateStatus(t *testing.T) {
	mockRepo := new(MockBookRepository)
	ctx := context.Background()

	bookIDs := []primitive.ObjectID{
		primitive.NewObjectID(),
		primitive.NewObjectID(),
		primitive.NewObjectID(),
	}
	newStatus := bookstore.BookStatusCompleted

	// 设置Mock期望
	mockRepo.On("BatchUpdateStatus", ctx, bookIDs, newStatus).Return(nil)

	// 执行测试
	err := mockRepo.BatchUpdateStatus(ctx, bookIDs, newStatus)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestBookRepository_IncrementViewCount 测试增加书籍浏览次数
func TestBookRepository_IncrementViewCount(t *testing.T) {
	mockRepo := new(MockBookRepository)
	ctx := context.Background()

	bookID := primitive.NewObjectID()

	// 设置Mock期望
	mockRepo.On("IncrementViewCount", ctx, bookID).Return(nil)

	// 执行测试
	err := mockRepo.IncrementViewCount(ctx, bookID)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestBookRepository_GetStats 测试获取书籍统计信息
func TestBookRepository_GetStats(t *testing.T) {
	mockRepo := new(MockBookRepository)
	ctx := context.Background()

	expectedStats := &bookstore.BookStats{
		TotalBooks:       1000,
		PublishedBooks:   800,
		DraftBooks:       200,
		RecommendedBooks: 50,
		FeaturedBooks:    30,
	}

	// 设置Mock期望
	mockRepo.On("GetStats", ctx).Return(expectedStats, nil)

	// 执行测试
	stats, err := mockRepo.GetStats(ctx)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, expectedStats.TotalBooks, stats.TotalBooks)
	assert.Equal(t, expectedStats.PublishedBooks, stats.PublishedBooks)
	mockRepo.AssertExpectations(t)
}
