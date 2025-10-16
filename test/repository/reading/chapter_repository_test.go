package reading

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"Qingyu_backend/models/reading/reader"
)

// MockChapterRepository Mock实现
type MockChapterRepository struct {
	mock.Mock
}

func (m *MockChapterRepository) Create(ctx context.Context, chapter *reader.Chapter) error {
	args := m.Called(ctx, chapter)
	return args.Error(0)
}

func (m *MockChapterRepository) GetByID(ctx context.Context, id string) (*reader.Chapter, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.Chapter), args.Error(1)
}

func (m *MockChapterRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockChapterRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockChapterRepository) GetByBookID(ctx context.Context, bookID string) ([]*reader.Chapter, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetByBookIDWithPagination(ctx context.Context, bookID string, limit, offset int64) ([]*reader.Chapter, error) {
	args := m.Called(ctx, bookID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetByChapterNum(ctx context.Context, bookID string, chapterNum int) (*reader.Chapter, error) {
	args := m.Called(ctx, bookID, chapterNum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetPrevChapter(ctx context.Context, bookID string, currentChapterNum int) (*reader.Chapter, error) {
	args := m.Called(ctx, bookID, currentChapterNum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetNextChapter(ctx context.Context, bookID string, currentChapterNum int) (*reader.Chapter, error) {
	args := m.Called(ctx, bookID, currentChapterNum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetFirstChapter(ctx context.Context, bookID string) (*reader.Chapter, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetLastChapter(ctx context.Context, bookID string) (*reader.Chapter, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetPublishedChapters(ctx context.Context, bookID string) ([]*reader.Chapter, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetVIPChapters(ctx context.Context, bookID string) ([]*reader.Chapter, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetFreeChapters(ctx context.Context, bookID string) ([]*reader.Chapter, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Chapter), args.Error(1)
}

func (m *MockChapterRepository) CountByBookID(ctx context.Context, bookID string) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) CountByStatus(ctx context.Context, bookID string, status int) (int64, error) {
	args := m.Called(ctx, bookID, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) CountVIPChapters(ctx context.Context, bookID string) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) BatchCreate(ctx context.Context, chapters []*reader.Chapter) error {
	args := m.Called(ctx, chapters)
	return args.Error(0)
}

func (m *MockChapterRepository) BatchUpdateStatus(ctx context.Context, chapterIDs []string, status int) error {
	args := m.Called(ctx, chapterIDs, status)
	return args.Error(0)
}

func (m *MockChapterRepository) BatchDelete(ctx context.Context, chapterIDs []string) error {
	args := m.Called(ctx, chapterIDs)
	return args.Error(0)
}

func (m *MockChapterRepository) CheckVIPAccess(ctx context.Context, chapterID string) (bool, error) {
	args := m.Called(ctx, chapterID)
	return args.Bool(0), args.Error(1)
}

func (m *MockChapterRepository) GetChapterPrice(ctx context.Context, chapterID string) (int64, error) {
	args := m.Called(ctx, chapterID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) GetChapterContent(ctx context.Context, chapterID string) (string, error) {
	args := m.Called(ctx, chapterID)
	return args.String(0), args.Error(1)
}

func (m *MockChapterRepository) UpdateChapterContent(ctx context.Context, chapterID string, content string) error {
	args := m.Called(ctx, chapterID, content)
	return args.Error(0)
}

func (m *MockChapterRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// TestChapterRepository 测试基本操作
func TestChapterRepository_GetByID(t *testing.T) {
	// 创建mock
	mockRepo := new(MockChapterRepository)
	ctx := context.Background()

	// 准备测试数据
	expectedChapter := &reader.Chapter{
		ID:         "chapter123",
		BookID:     "book123",
		Title:      "第一章",
		ChapterNum: 1,
		Content:    "章节内容",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 设置期望
	mockRepo.On("GetByID", ctx, "chapter123").Return(expectedChapter, nil)

	// 执行测试
	chapter, err := mockRepo.GetByID(ctx, "chapter123")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, chapter)
	assert.Equal(t, "chapter123", chapter.ID)
	assert.Equal(t, "第一章", chapter.Title)
	mockRepo.AssertExpectations(t)
}

func TestChapterRepository_GetByBookID(t *testing.T) {
	mockRepo := new(MockChapterRepository)
	ctx := context.Background()

	// 准备测试数据
	expectedChapters := []*reader.Chapter{
		{ID: "ch1", BookID: "book123", Title: "第一章", ChapterNum: 1},
		{ID: "ch2", BookID: "book123", Title: "第二章", ChapterNum: 2},
	}

	// 设置期望
	mockRepo.On("GetByBookID", ctx, "book123").Return(expectedChapters, nil)

	// 执行测试
	chapters, err := mockRepo.GetByBookID(ctx, "book123")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, chapters)
	assert.Len(t, chapters, 2)
	assert.Equal(t, "第一章", chapters[0].Title)
	mockRepo.AssertExpectations(t)
}

func TestChapterRepository_GetNextChapter(t *testing.T) {
	mockRepo := new(MockChapterRepository)
	ctx := context.Background()

	// 准备测试数据
	nextChapter := &reader.Chapter{
		ID:         "ch2",
		BookID:     "book123",
		Title:      "第二章",
		ChapterNum: 2,
	}

	// 设置期望
	mockRepo.On("GetNextChapter", ctx, "book123", 1).Return(nextChapter, nil)

	// 执行测试
	chapter, err := mockRepo.GetNextChapter(ctx, "book123", 1)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, chapter)
	assert.Equal(t, 2, chapter.ChapterNum)
	mockRepo.AssertExpectations(t)
}

func TestChapterRepository_CheckVIPAccess(t *testing.T) {
	mockRepo := new(MockChapterRepository)
	ctx := context.Background()

	// 设置期望 - VIP章节
	mockRepo.On("CheckVIPAccess", ctx, "vip_chapter").Return(true, nil)

	// 执行测试
	isVIP, err := mockRepo.CheckVIPAccess(ctx, "vip_chapter")

	// 验证结果
	assert.NoError(t, err)
	assert.True(t, isVIP)
	mockRepo.AssertExpectations(t)
}

func TestChapterRepository_CountByBookID(t *testing.T) {
	mockRepo := new(MockChapterRepository)
	ctx := context.Background()

	// 设置期望
	mockRepo.On("CountByBookID", ctx, "book123").Return(int64(10), nil)

	// 执行测试
	count, err := mockRepo.CountByBookID(ctx, "book123")

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, int64(10), count)
	mockRepo.AssertExpectations(t)
}
