package reading

import (
	reader2 "Qingyu_backend/models/reader"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"Qingyu_backend/service/base"
	"Qingyu_backend/service/reading"
)

// Mock EventBus
type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Subscribe(eventType string, handler base.EventHandler) error {
	args := m.Called(eventType, handler)
	return args.Error(0)
}

func (m *MockEventBus) Unsubscribe(eventType string, handlerName string) error {
	args := m.Called(eventType, handlerName)
	return args.Error(0)
}

func (m *MockEventBus) Publish(ctx context.Context, event base.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventBus) PublishAsync(ctx context.Context, event base.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// Mock ReadingProgressRepository
type MockProgressRepository struct {
	mock.Mock
}

func (m *MockProgressRepository) Create(ctx context.Context, progress *reader2.ReadingProgress) error {
	args := m.Called(ctx, progress)
	return args.Error(0)
}

func (m *MockProgressRepository) GetByID(ctx context.Context, id string) (*reader2.ReadingProgress, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader2.ReadingProgress), args.Error(1)
}

func (m *MockProgressRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockProgressRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProgressRepository) GetByUserAndBook(ctx context.Context, userID, bookID string) (*reader2.ReadingProgress, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader2.ReadingProgress), args.Error(1)
}

func (m *MockProgressRepository) GetByUser(ctx context.Context, userID string) ([]*reader2.ReadingProgress, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.ReadingProgress), args.Error(1)
}

func (m *MockProgressRepository) GetRecentReadingByUser(ctx context.Context, userID string, limit int) ([]*reader2.ReadingProgress, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.ReadingProgress), args.Error(1)
}

func (m *MockProgressRepository) SaveProgress(ctx context.Context, userID, bookID, chapterID string, progress float64) error {
	args := m.Called(ctx, userID, bookID, chapterID, progress)
	return args.Error(0)
}

func (m *MockProgressRepository) UpdateReadingTime(ctx context.Context, userID, bookID string, duration int64) error {
	args := m.Called(ctx, userID, bookID, duration)
	return args.Error(0)
}

func (m *MockProgressRepository) UpdateLastReadAt(ctx context.Context, userID, bookID string) error {
	args := m.Called(ctx, userID, bookID)
	return args.Error(0)
}

func (m *MockProgressRepository) BatchUpdateProgress(ctx context.Context, progresses []*reader2.ReadingProgress) error {
	args := m.Called(ctx, progresses)
	return args.Error(0)
}

func (m *MockProgressRepository) GetTotalReadingTime(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProgressRepository) GetReadingTimeByBook(ctx context.Context, userID, bookID string) (int64, error) {
	args := m.Called(ctx, userID, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProgressRepository) GetReadingTimeByPeriod(ctx context.Context, userID string, startTime, endTime time.Time) (int64, error) {
	args := m.Called(ctx, userID, startTime, endTime)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProgressRepository) CountReadingBooks(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProgressRepository) GetReadingHistory(ctx context.Context, userID string, limit, offset int) ([]*reader2.ReadingProgress, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.ReadingProgress), args.Error(1)
}

func (m *MockProgressRepository) GetUnfinishedBooks(ctx context.Context, userID string) ([]*reader2.ReadingProgress, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.ReadingProgress), args.Error(1)
}

func (m *MockProgressRepository) GetFinishedBooks(ctx context.Context, userID string) ([]*reader2.ReadingProgress, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.ReadingProgress), args.Error(1)
}

func (m *MockProgressRepository) SyncProgress(ctx context.Context, userID string, progresses []*reader2.ReadingProgress) error {
	args := m.Called(ctx, userID, progresses)
	return args.Error(0)
}

func (m *MockProgressRepository) GetProgressesByUser(ctx context.Context, userID string, updatedAfter time.Time) ([]*reader2.ReadingProgress, error) {
	args := m.Called(ctx, userID, updatedAfter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.ReadingProgress), args.Error(1)
}

func (m *MockProgressRepository) DeleteOldProgress(ctx context.Context, beforeTime time.Time) error {
	args := m.Called(ctx, beforeTime)
	return args.Error(0)
}

func (m *MockProgressRepository) DeleteByBook(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockProgressRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Mock AnnotationRepository
type MockAnnotationRepository struct {
	mock.Mock
}

func (m *MockAnnotationRepository) Create(ctx context.Context, annotation *reader2.Annotation) error {
	args := m.Called(ctx, annotation)
	return args.Error(0)
}

func (m *MockAnnotationRepository) GetByID(ctx context.Context, id string) (*reader2.Annotation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader2.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockAnnotationRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAnnotationRepository) GetByUserAndBook(ctx context.Context, userID, bookID string) ([]*reader2.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) GetByUserAndChapter(ctx context.Context, userID, bookID, chapterID string) ([]*reader2.Annotation, error) {
	args := m.Called(ctx, userID, bookID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) GetByType(ctx context.Context, userID, bookID string, annotationType int) ([]*reader2.Annotation, error) {
	args := m.Called(ctx, userID, bookID, annotationType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) GetNotes(ctx context.Context, userID, bookID string) ([]*reader2.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) GetNotesByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*reader2.Annotation, error) {
	args := m.Called(ctx, userID, bookID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) SearchNotes(ctx context.Context, userID, keyword string) ([]*reader2.Annotation, error) {
	args := m.Called(ctx, userID, keyword)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) GetBookmarks(ctx context.Context, userID, bookID string) ([]*reader2.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) GetBookmarkByPosition(ctx context.Context, userID, bookID, chapterID string, startOffset int) (*reader2.Annotation, error) {
	args := m.Called(ctx, userID, bookID, chapterID, startOffset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader2.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) GetLatestBookmark(ctx context.Context, userID, bookID string) (*reader2.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader2.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) GetHighlights(ctx context.Context, userID, bookID string) ([]*reader2.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) GetHighlightsByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*reader2.Annotation, error) {
	args := m.Called(ctx, userID, bookID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) CountByUser(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAnnotationRepository) CountByBook(ctx context.Context, userID, bookID string) (int64, error) {
	args := m.Called(ctx, userID, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAnnotationRepository) CountByType(ctx context.Context, userID string, annotationType int) (int64, error) {
	args := m.Called(ctx, userID, annotationType)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAnnotationRepository) BatchCreate(ctx context.Context, annotations []*reader2.Annotation) error {
	args := m.Called(ctx, annotations)
	return args.Error(0)
}

func (m *MockAnnotationRepository) BatchDelete(ctx context.Context, annotationIDs []string) error {
	args := m.Called(ctx, annotationIDs)
	return args.Error(0)
}

func (m *MockAnnotationRepository) DeleteByBook(ctx context.Context, userID, bookID string) error {
	args := m.Called(ctx, userID, bookID)
	return args.Error(0)
}

func (m *MockAnnotationRepository) DeleteByChapter(ctx context.Context, userID, bookID, chapterID string) error {
	args := m.Called(ctx, userID, bookID, chapterID)
	return args.Error(0)
}

func (m *MockAnnotationRepository) SyncAnnotations(ctx context.Context, userID string, annotations []*reader2.Annotation) error {
	args := m.Called(ctx, userID, annotations)
	return args.Error(0)
}

func (m *MockAnnotationRepository) GetRecentAnnotations(ctx context.Context, userID string, limit int) ([]*reader2.Annotation, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) GetPublicAnnotations(ctx context.Context, bookID, chapterID string) ([]*reader2.Annotation, error) {
	args := m.Called(ctx, bookID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) GetSharedAnnotations(ctx context.Context, userID string) ([]*reader2.Annotation, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// TestReaderService_SaveReadingProgress 测试保存阅读进度
func TestReaderService_SaveReadingProgress(t *testing.T) {
	// 创建mocks
	mockProgressRepo := new(MockProgressRepository)
	mockEventBus := new(MockEventBus)

	// 创建service
	service := reading.NewReaderService(
		nil, // chapterRepo不需要
		mockProgressRepo,
		nil, // annotationRepo不需要
		nil, // settingsRepo不需要
		mockEventBus,
		nil, // cacheService不需要
		nil, // vipService不需要
	)

	ctx := context.Background()

	// 设置期望
	mockProgressRepo.On("SaveProgress", ctx, "user123", "book123", "chapter1", 0.5).Return(nil)
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil)

	// 执行测试
	err := service.SaveReadingProgress(ctx, "user123", "book123", "chapter1", 0.5)

	// 验证结果
	assert.NoError(t, err)
	mockProgressRepo.AssertExpectations(t)
}

// TestReaderService_GetTotalReadingTime 测试获取总阅读时长
func TestReaderService_GetTotalReadingTime(t *testing.T) {
	mockProgressRepo := new(MockProgressRepository)

	service := reading.NewReaderService(nil, mockProgressRepo, nil, nil, nil, nil, nil)

	ctx := context.Background()

	// 设置期望 - 返回3600秒(1小时)
	mockProgressRepo.On("GetTotalReadingTime", ctx, "user123").Return(int64(3600), nil)

	// 执行测试
	totalTime, err := service.GetTotalReadingTime(ctx, "user123")

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, int64(3600), totalTime)
	mockProgressRepo.AssertExpectations(t)
}

// TestReaderService_CreateAnnotation 测试创建标注
func TestReaderService_CreateAnnotation(t *testing.T) {
	mockAnnotationRepo := new(MockAnnotationRepository)
	mockEventBus := new(MockEventBus)

	service := reading.NewReaderService(nil, nil, mockAnnotationRepo, nil, mockEventBus, nil, nil)

	ctx := context.Background()

	// 准备测试数据
	annotation := &reader2.Annotation{
		UserID:    "user123",
		BookID:    "book123",
		ChapterID: "chapter1",
		Type:      "bookmark",
		Text:      "测试标注",
	}

	// 设置期望
	mockAnnotationRepo.On("Create", ctx, annotation).Return(nil)
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil)

	// 执行测试
	err := service.CreateAnnotation(ctx, annotation)

	// 验证结果
	assert.NoError(t, err)
	mockAnnotationRepo.AssertExpectations(t)
}

// TestReaderService_ServiceInfo 测试服务信息方法
func TestReaderService_ServiceInfo(t *testing.T) {
	service := reading.NewReaderService(nil, nil, nil, nil, nil, nil, nil)

	// 验证服务名称和版本
	assert.Equal(t, "ReaderService", service.GetServiceName())
	assert.Equal(t, "1.0.0", service.GetVersion())
}
