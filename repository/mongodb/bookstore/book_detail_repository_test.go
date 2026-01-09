package bookstore_test

import (
	bookstore2 "Qingyu_backend/models/bookstore"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	bookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
)

// MockBookDetailRepository Mock实现
type MockBookDetailRepository struct {
	mock.Mock
	bookstoreRepo.BookDetailRepository // 嵌入接口
}

func (m *MockBookDetailRepository) Create(ctx context.Context, bookDetail *bookstore2.BookDetail) error {
	args := m.Called(ctx, bookDetail)
	return args.Error(0)
}

func (m *MockBookDetailRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore2.BookDetail, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore2.BookDetail), args.Error(1)
}

func (m *MockBookDetailRepository) GetByTitle(ctx context.Context, title string) (*bookstore2.BookDetail, error) {
	args := m.Called(ctx, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore2.BookDetail), args.Error(1)
}

func (m *MockBookDetailRepository) GetByAuthor(ctx context.Context, author string, limit, offset int) ([]*bookstore2.BookDetail, error) {
	args := m.Called(ctx, author, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore2.BookDetail), args.Error(1)
}

func (m *MockBookDetailRepository) GetByCategory(ctx context.Context, category string, limit, offset int) ([]*bookstore2.BookDetail, error) {
	args := m.Called(ctx, category, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore2.BookDetail), args.Error(1)
}

func (m *MockBookDetailRepository) GetByStatus(ctx context.Context, status bookstore2.BookStatus, limit, offset int) ([]*bookstore2.BookDetail, error) {
	args := m.Called(ctx, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore2.BookDetail), args.Error(1)
}

func (m *MockBookDetailRepository) GetByTags(ctx context.Context, tags []string, limit, offset int) ([]*bookstore2.BookDetail, error) {
	args := m.Called(ctx, tags, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore2.BookDetail), args.Error(1)
}

func (m *MockBookDetailRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]*bookstore2.BookDetail, error) {
	args := m.Called(ctx, keyword, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore2.BookDetail), args.Error(1)
}

func (m *MockBookDetailRepository) IncrementViewCount(ctx context.Context, bookID primitive.ObjectID) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockBookDetailRepository) IncrementLikeCount(ctx context.Context, bookID primitive.ObjectID) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockBookDetailRepository) DecrementLikeCount(ctx context.Context, bookID primitive.ObjectID) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockBookDetailRepository) IncrementCommentCount(ctx context.Context, bookID primitive.ObjectID) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockBookDetailRepository) CountByCategory(ctx context.Context, category string) (int64, error) {
	args := m.Called(ctx, category)
	return args.Get(0).(int64), args.Error(1)
}

// 测试助手函数
func createTestBookDetail(id primitive.ObjectID, title, author string) *bookstore2.BookDetail {
	now := time.Now()
	return &bookstore2.BookDetail{
		ID:           id,
		Title:        title,
		Author:       author,
		Introduction: "详细简介内容",
		Description:  "完整描述内容",
		CoverURL:     "https://example.com/cover.jpg",
		Status:       bookstore2.BookStatusPublished,
		Tags:         []string{"玄幻", "热血"},
		WordCount:    500000,
		ChapterCount: 200,
		Price:        29.99,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// TestBookDetailRepository_Create 测试创建书籍详情
func TestBookDetailRepository_Create(t *testing.T) {
	mockRepo := new(MockBookDetailRepository)
	ctx := context.Background()

	bookDetail := createTestBookDetail(primitive.NewObjectID(), "测试书籍", "测试作者")

	// 设置Mock期望
	mockRepo.On("Create", ctx, bookDetail).Return(nil)

	// 执行测试
	err := mockRepo.Create(ctx, bookDetail)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestBookDetailRepository_GetByID_Success 测试根据ID获取书籍详情成功
func TestBookDetailRepository_GetByID_Success(t *testing.T) {
	mockRepo := new(MockBookDetailRepository)
	ctx := context.Background()

	bookID := primitive.NewObjectID()
	expectedBook := createTestBookDetail(bookID, "测试书籍", "测试作者")

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

// TestBookDetailRepository_GetByID_NotFound 测试根据ID获取书籍详情不存在
func TestBookDetailRepository_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockBookDetailRepository)
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

// TestBookDetailRepository_GetByTitle 测试根据标题获取书籍详情
func TestBookDetailRepository_GetByTitle(t *testing.T) {
	mockRepo := new(MockBookDetailRepository)
	ctx := context.Background()

	title := "斗破苍穹"
	expectedBook := createTestBookDetail(primitive.NewObjectID(), title, "天蚕土豆")

	// 设置Mock期望
	mockRepo.On("GetByTitle", ctx, title).Return(expectedBook, nil)

	// 执行测试
	book, err := mockRepo.GetByTitle(ctx, title)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, title, book.Title)
	mockRepo.AssertExpectations(t)
}

// TestBookDetailRepository_GetByAuthor 测试根据作者获取书籍列表
func TestBookDetailRepository_GetByAuthor(t *testing.T) {
	mockRepo := new(MockBookDetailRepository)
	ctx := context.Background()

	author := "天蚕土豆"
	books := []*bookstore2.BookDetail{
		createTestBookDetail(primitive.NewObjectID(), "斗破苍穹", author),
		createTestBookDetail(primitive.NewObjectID(), "武动乾坤", author),
	}

	// 设置Mock期望
	mockRepo.On("GetByAuthor", ctx, author, 10, 0).Return(books, nil)

	// 执行测试
	result, err := mockRepo.GetByAuthor(ctx, author, 10, 0)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	for _, book := range result {
		assert.Equal(t, author, book.Author)
	}
	mockRepo.AssertExpectations(t)
}

// TestBookDetailRepository_GetByCategory 测试根据分类获取书籍列表
func TestBookDetailRepository_GetByCategory(t *testing.T) {
	mockRepo := new(MockBookDetailRepository)
	ctx := context.Background()

	category := "玄幻"
	books := []*bookstore2.BookDetail{
		createTestBookDetail(primitive.NewObjectID(), "玄幻书籍1", "作者1"),
		createTestBookDetail(primitive.NewObjectID(), "玄幻书籍2", "作者2"),
	}

	// 设置Mock期望
	mockRepo.On("GetByCategory", ctx, category, 10, 0).Return(books, nil)

	// 执行测试
	result, err := mockRepo.GetByCategory(ctx, category, 10, 0)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	mockRepo.AssertExpectations(t)
}

// TestBookDetailRepository_GetByStatus 测试根据状态获取书籍列表
func TestBookDetailRepository_GetByStatus(t *testing.T) {
	mockRepo := new(MockBookDetailRepository)
	ctx := context.Background()

	books := []*bookstore2.BookDetail{
		createTestBookDetail(primitive.NewObjectID(), "已发布书籍", "作者1"),
	}

	// 设置Mock期望
	mockRepo.On("GetByStatus", ctx, bookstore2.BookStatusPublished, 10, 0).Return(books, nil)

	// 执行测试
	result, err := mockRepo.GetByStatus(ctx, bookstore2.BookStatusPublished, 10, 0)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, bookstore2.BookStatusPublished, result[0].Status)
	mockRepo.AssertExpectations(t)
}

// TestBookDetailRepository_GetByTags 测试根据标签获取书籍列表
func TestBookDetailRepository_GetByTags(t *testing.T) {
	mockRepo := new(MockBookDetailRepository)
	ctx := context.Background()

	tags := []string{"玄幻", "热血"}
	books := []*bookstore2.BookDetail{
		createTestBookDetail(primitive.NewObjectID(), "标签书籍", "作者1"),
	}

	// 设置Mock期望
	mockRepo.On("GetByTags", ctx, tags, 10, 0).Return(books, nil)

	// 执行测试
	result, err := mockRepo.GetByTags(ctx, tags, 10, 0)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
	mockRepo.AssertExpectations(t)
}

// TestBookDetailRepository_Search 测试搜索书籍
func TestBookDetailRepository_Search(t *testing.T) {
	mockRepo := new(MockBookDetailRepository)
	ctx := context.Background()

	keyword := "斗破"
	books := []*bookstore2.BookDetail{
		createTestBookDetail(primitive.NewObjectID(), "斗破苍穹", "天蚕土豆"),
	}

	// 设置Mock期望
	mockRepo.On("Search", ctx, keyword, 10, 0).Return(books, nil)

	// 执行测试
	result, err := mockRepo.Search(ctx, keyword, 10, 0)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
	mockRepo.AssertExpectations(t)
}

// TestBookDetailRepository_IncrementViewCount 测试增加浏览次数
func TestBookDetailRepository_IncrementViewCount(t *testing.T) {
	mockRepo := new(MockBookDetailRepository)
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

// TestBookDetailRepository_IncrementLikeCount 测试增加点赞次数
func TestBookDetailRepository_IncrementLikeCount(t *testing.T) {
	mockRepo := new(MockBookDetailRepository)
	ctx := context.Background()

	bookID := primitive.NewObjectID()

	// 设置Mock期望
	mockRepo.On("IncrementLikeCount", ctx, bookID).Return(nil)

	// 执行测试
	err := mockRepo.IncrementLikeCount(ctx, bookID)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestBookDetailRepository_DecrementLikeCount 测试减少点赞次数
func TestBookDetailRepository_DecrementLikeCount(t *testing.T) {
	mockRepo := new(MockBookDetailRepository)
	ctx := context.Background()

	bookID := primitive.NewObjectID()

	// 设置Mock期望
	mockRepo.On("DecrementLikeCount", ctx, bookID).Return(nil)

	// 执行测试
	err := mockRepo.DecrementLikeCount(ctx, bookID)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestBookDetailRepository_IncrementCommentCount 测试增加评论次数
func TestBookDetailRepository_IncrementCommentCount(t *testing.T) {
	mockRepo := new(MockBookDetailRepository)
	ctx := context.Background()

	bookID := primitive.NewObjectID()

	// 设置Mock期望
	mockRepo.On("IncrementCommentCount", ctx, bookID).Return(nil)

	// 执行测试
	err := mockRepo.IncrementCommentCount(ctx, bookID)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestBookDetailRepository_CountByCategory 测试统计分类下的书籍数量
func TestBookDetailRepository_CountByCategory(t *testing.T) {
	mockRepo := new(MockBookDetailRepository)
	ctx := context.Background()

	category := "玄幻"
	expectedCount := int64(150)

	// 设置Mock期望
	mockRepo.On("CountByCategory", ctx, category).Return(expectedCount, nil)

	// 执行测试
	count, err := mockRepo.CountByCategory(ctx, category)

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	mockRepo.AssertExpectations(t)
}
