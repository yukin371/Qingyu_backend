package reader

import (
	bookstoreModel "Qingyu_backend/models/bookstore"
	"Qingyu_backend/models/reader"
	"Qingyu_backend/repository/interfaces/infrastructure"
	baseInterface "Qingyu_backend/service/interfaces/base"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// =========================
// Mock Repository实现
// =========================

// MockReadingProgressRepository Mock阅读进度仓储
type MockReadingProgressRepository struct {
	mock.Mock
}

func (m *MockReadingProgressRepository) Create(ctx context.Context, progress *reader.ReadingProgress) error {
	args := m.Called(ctx, progress)
	return args.Error(0)
}

func (m *MockReadingProgressRepository) GetByID(ctx context.Context, id string) (*reader.ReadingProgress, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.ReadingProgress), args.Error(1)
}

func (m *MockReadingProgressRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockReadingProgressRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockReadingProgressRepository) GetByUserAndBook(ctx context.Context, userID, bookID string) (*reader.ReadingProgress, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.ReadingProgress), args.Error(1)
}

func (m *MockReadingProgressRepository) GetByUser(ctx context.Context, userID string) ([]*reader.ReadingProgress, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.ReadingProgress), args.Error(1)
}

func (m *MockReadingProgressRepository) GetRecentReadingByUser(ctx context.Context, userID string, limit int) ([]*reader.ReadingProgress, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.ReadingProgress), args.Error(1)
}

func (m *MockReadingProgressRepository) SaveProgress(ctx context.Context, userID, bookID, chapterID string, progress float64) error {
	args := m.Called(ctx, userID, bookID, chapterID, progress)
	return args.Error(0)
}

func (m *MockReadingProgressRepository) UpdateReadingTime(ctx context.Context, userID, bookID string, duration int64) error {
	args := m.Called(ctx, userID, bookID, duration)
	return args.Error(0)
}

func (m *MockReadingProgressRepository) UpdateLastReadAt(ctx context.Context, userID, bookID string) error {
	args := m.Called(ctx, userID, bookID)
	return args.Error(0)
}

func (m *MockReadingProgressRepository) BatchUpdateProgress(ctx context.Context, progresses []*reader.ReadingProgress) error {
	args := m.Called(ctx, progresses)
	return args.Error(0)
}

func (m *MockReadingProgressRepository) GetTotalReadingTime(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockReadingProgressRepository) GetReadingTimeByBook(ctx context.Context, userID, bookID string) (int64, error) {
	args := m.Called(ctx, userID, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockReadingProgressRepository) GetReadingTimeByPeriod(ctx context.Context, userID string, startTime, endTime time.Time) (int64, error) {
	args := m.Called(ctx, userID, startTime, endTime)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockReadingProgressRepository) CountReadingBooks(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockReadingProgressRepository) GetReadingHistory(ctx context.Context, userID string, limit, offset int) ([]*reader.ReadingProgress, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.ReadingProgress), args.Error(1)
}

func (m *MockReadingProgressRepository) GetUnfinishedBooks(ctx context.Context, userID string) ([]*reader.ReadingProgress, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.ReadingProgress), args.Error(1)
}

func (m *MockReadingProgressRepository) GetFinishedBooks(ctx context.Context, userID string) ([]*reader.ReadingProgress, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.ReadingProgress), args.Error(1)
}

func (m *MockReadingProgressRepository) SyncProgress(ctx context.Context, userID string, progresses []*reader.ReadingProgress) error {
	args := m.Called(ctx, userID, progresses)
	return args.Error(0)
}

func (m *MockReadingProgressRepository) GetProgressesByUser(ctx context.Context, userID string, updatedAfter time.Time) ([]*reader.ReadingProgress, error) {
	args := m.Called(ctx, userID, updatedAfter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.ReadingProgress), args.Error(1)
}

func (m *MockReadingProgressRepository) DeleteOldProgress(ctx context.Context, beforeTime time.Time) error {
	args := m.Called(ctx, beforeTime)
	return args.Error(0)
}

func (m *MockReadingProgressRepository) DeleteByBook(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockReadingProgressRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockAnnotationRepository Mock标注仓储
type MockAnnotationRepository struct {
	mock.Mock
}

func (m *MockAnnotationRepository) Create(ctx context.Context, annotation *reader.Annotation) error {
	args := m.Called(ctx, annotation)
	return args.Error(0)
}

func (m *MockAnnotationRepository) GetByID(ctx context.Context, id string) (*reader.Annotation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockAnnotationRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAnnotationRepository) GetByUserAndBook(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) GetByUserAndChapter(ctx context.Context, userID, bookID, chapterID string) ([]*reader.Annotation, error) {
	args := m.Called(ctx, userID, bookID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) GetByType(ctx context.Context, userID, bookID string, annotationType string) ([]*reader.Annotation, error) {
	args := m.Called(ctx, userID, bookID, annotationType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) GetNotes(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) GetNotesByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*reader.Annotation, error) {
	args := m.Called(ctx, userID, bookID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) SearchNotes(ctx context.Context, userID string, keyword string) ([]*reader.Annotation, error) {
	args := m.Called(ctx, userID, keyword)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) GetBookmarks(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) GetBookmarkByPosition(ctx context.Context, userID, bookID, chapterID string, startOffset int) (*reader.Annotation, error) {
	args := m.Called(ctx, userID, bookID, chapterID, startOffset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) GetLatestBookmark(ctx context.Context, userID, bookID string) (*reader.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) GetHighlights(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) GetHighlightsByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*reader.Annotation, error) {
	args := m.Called(ctx, userID, bookID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) CountByUser(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAnnotationRepository) CountByBook(ctx context.Context, userID, bookID string) (int64, error) {
	args := m.Called(ctx, userID, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAnnotationRepository) CountByType(ctx context.Context, userID string, annotationType string) (int64, error) {
	args := m.Called(ctx, userID, annotationType)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAnnotationRepository) BatchCreate(ctx context.Context, annotations []*reader.Annotation) error {
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

func (m *MockAnnotationRepository) SyncAnnotations(ctx context.Context, userID string, annotations []*reader.Annotation) error {
	args := m.Called(ctx, userID, annotations)
	return args.Error(0)
}

func (m *MockAnnotationRepository) GetRecentAnnotations(ctx context.Context, userID string, limit int) ([]*reader.Annotation, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) GetPublicAnnotations(ctx context.Context, bookID, chapterID string) ([]*reader.Annotation, error) {
	args := m.Called(ctx, bookID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) GetSharedAnnotations(ctx context.Context, userID string) ([]*reader.Annotation, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Annotation), args.Error(1)
}

func (m *MockAnnotationRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockReadingSettingsRepository Mock阅读设置仓储
type MockReadingSettingsRepository struct {
	mock.Mock
}

func (m *MockReadingSettingsRepository) Create(ctx context.Context, settings *reader.ReadingSettings) error {
	args := m.Called(ctx, settings)
	return args.Error(0)
}

func (m *MockReadingSettingsRepository) GetByID(ctx context.Context, id string) (*reader.ReadingSettings, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.ReadingSettings), args.Error(1)
}

func (m *MockReadingSettingsRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockReadingSettingsRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockReadingSettingsRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*reader.ReadingSettings, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.ReadingSettings), args.Error(1)
}

func (m *MockReadingSettingsRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockReadingSettingsRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockReadingSettingsRepository) GetByUserID(ctx context.Context, userID string) (*reader.ReadingSettings, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.ReadingSettings), args.Error(1)
}

func (m *MockReadingSettingsRepository) UpdateByUserID(ctx context.Context, userID string, settings *reader.ReadingSettings) error {
	args := m.Called(ctx, userID, settings)
	return args.Error(0)
}

func (m *MockReadingSettingsRepository) CreateDefaultSettings(ctx context.Context, userID string) (*reader.ReadingSettings, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.ReadingSettings), args.Error(1)
}

func (m *MockReadingSettingsRepository) ExistsByUserID(ctx context.Context, userID string) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(bool), args.Error(1)
}

// MockChapterService Mock章节服务
type MockChapterService struct {
	mock.Mock
}

func (m *MockChapterService) CreateChapter(ctx context.Context, chapter *bookstoreModel.Chapter) error {
	args := m.Called(ctx, chapter)
	return args.Error(0)
}

func (m *MockChapterService) GetChapterByID(ctx context.Context, id string) (*bookstoreModel.Chapter, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.Chapter), args.Error(1)
}

func (m *MockChapterService) UpdateChapter(ctx context.Context, chapter *bookstoreModel.Chapter) error {
	args := m.Called(ctx, chapter)
	return args.Error(0)
}

func (m *MockChapterService) DeleteChapter(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockChapterService) GetChaptersByBookID(ctx context.Context, bookID string, page, pageSize int) ([]*bookstoreModel.Chapter, int64, error) {
	args := m.Called(ctx, bookID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstoreModel.Chapter), args.Get(1).(int64), args.Error(2)
}

func (m *MockChapterService) GetChapterByBookIDAndNum(ctx context.Context, bookID string, chapterNum int) (*bookstoreModel.Chapter, error) {
	args := m.Called(ctx, bookID, chapterNum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.Chapter), args.Error(1)
}

func (m *MockChapterService) GetChaptersByTitle(ctx context.Context, title string, page, pageSize int) ([]*bookstoreModel.Chapter, int64, error) {
	args := m.Called(ctx, title, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstoreModel.Chapter), args.Get(1).(int64), args.Error(2)
}

func (m *MockChapterService) GetFreeChaptersByBookID(ctx context.Context, bookID string, page, pageSize int) ([]*bookstoreModel.Chapter, int64, error) {
	args := m.Called(ctx, bookID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstoreModel.Chapter), args.Get(1).(int64), args.Error(2)
}

func (m *MockChapterService) GetPaidChaptersByBookID(ctx context.Context, bookID string, page, pageSize int) ([]*bookstoreModel.Chapter, int64, error) {
	args := m.Called(ctx, bookID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstoreModel.Chapter), args.Get(1).(int64), args.Error(2)
}

func (m *MockChapterService) GetPublishedChaptersByBookID(ctx context.Context, bookID string, page, pageSize int) ([]*bookstoreModel.Chapter, int64, error) {
	args := m.Called(ctx, bookID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstoreModel.Chapter), args.Get(1).(int64), args.Error(2)
}

func (m *MockChapterService) GetPreviousChapter(ctx context.Context, bookID string, chapterNum int) (*bookstoreModel.Chapter, error) {
	args := m.Called(ctx, bookID, chapterNum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.Chapter), args.Error(1)
}

func (m *MockChapterService) GetNextChapter(ctx context.Context, bookID string, chapterNum int) (*bookstoreModel.Chapter, error) {
	args := m.Called(ctx, bookID, chapterNum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.Chapter), args.Error(1)
}

func (m *MockChapterService) GetFirstChapter(ctx context.Context, bookID string) (*bookstoreModel.Chapter, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.Chapter), args.Error(1)
}

func (m *MockChapterService) GetLastChapter(ctx context.Context, bookID string) (*bookstoreModel.Chapter, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.Chapter), args.Error(1)
}

func (m *MockChapterService) GetChapterCountByBookID(ctx context.Context, bookID string) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterService) GetFreeChapterCountByBookID(ctx context.Context, bookID string) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterService) GetPaidChapterCountByBookID(ctx context.Context, bookID string) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterService) GetTotalWordCountByBookID(ctx context.Context, bookID string) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterService) GetChapterStats(ctx context.Context, bookID string) (map[string]interface{}, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockChapterService) GetChapterContent(ctx context.Context, chapterID, userID string) (string, error) {
	args := m.Called(ctx, chapterID, userID)
	return args.String(0), args.Error(1)
}

func (m *MockChapterService) UpdateChapterContent(ctx context.Context, chapterID string, content string) error {
	args := m.Called(ctx, chapterID, content)
	return args.Error(0)
}

func (m *MockChapterService) PublishChapter(ctx context.Context, chapterID string) error {
	args := m.Called(ctx, chapterID)
	return args.Error(0)
}

func (m *MockChapterService) UnpublishChapter(ctx context.Context, chapterID string) error {
	args := m.Called(ctx, chapterID)
	return args.Error(0)
}

func (m *MockChapterService) BatchUpdateChapterPrice(ctx context.Context, chapterIDs []string, price float64) error {
	args := m.Called(ctx, chapterIDs, price)
	return args.Error(0)
}

func (m *MockChapterService) BatchPublishChapters(ctx context.Context, chapterIDs []string) error {
	args := m.Called(ctx, chapterIDs)
	return args.Error(0)
}

func (m *MockChapterService) BatchDeleteChapters(ctx context.Context, chapterIDs []string) error {
	args := m.Called(ctx, chapterIDs)
	return args.Error(0)
}

func (m *MockChapterService) BatchDeleteChaptersByBookID(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockChapterService) SearchChapters(ctx context.Context, keyword string, page, pageSize int) ([]*bookstoreModel.Chapter, int64, error) {
	args := m.Called(ctx, keyword, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstoreModel.Chapter), args.Get(1).(int64), args.Error(2)
}

// MockReaderCacheService Mock阅读器缓存服务
type MockReaderCacheService struct {
	mock.Mock
}

func (m *MockReaderCacheService) GetReadingSettings(ctx context.Context, userID string) (*reader.ReadingSettings, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.ReadingSettings), args.Error(1)
}

func (m *MockReaderCacheService) SetReadingSettings(ctx context.Context, userID string, settings *reader.ReadingSettings, expiration time.Duration) error {
	args := m.Called(ctx, userID, settings, expiration)
	return args.Error(0)
}

func (m *MockReaderCacheService) InvalidateReadingSettings(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockReaderCacheService) GetReadingProgress(ctx context.Context, userID, bookID string) (*reader.ReadingProgress, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.ReadingProgress), args.Error(1)
}

func (m *MockReaderCacheService) SetReadingProgress(ctx context.Context, userID, bookID string, progress *reader.ReadingProgress, expiration time.Duration) error {
	args := m.Called(ctx, userID, bookID, progress, expiration)
	return args.Error(0)
}

func (m *MockReaderCacheService) InvalidateReadingProgress(ctx context.Context, userID, bookID string) error {
	args := m.Called(ctx, userID, bookID)
	return args.Error(0)
}

func (m *MockReaderCacheService) GetAnnotationsByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*reader.Annotation, error) {
	args := m.Called(ctx, userID, bookID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.Annotation), args.Error(1)
}

func (m *MockReaderCacheService) SetAnnotationsByChapter(ctx context.Context, userID, bookID, chapterID string, annotations []*reader.Annotation, expiration time.Duration) error {
	args := m.Called(ctx, userID, bookID, chapterID, annotations, expiration)
	return args.Error(0)
}

func (m *MockReaderCacheService) InvalidateAnnotationsByChapter(ctx context.Context, userID, bookID, chapterID string) error {
	args := m.Called(ctx, userID, bookID, chapterID)
	return args.Error(0)
}

func (m *MockReaderCacheService) InvalidateUserData(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// MockVIPPermissionService MockVIP权限服务
type MockVIPPermissionService struct {
	mock.Mock
}

func (m *MockVIPPermissionService) CheckVIPAccess(ctx context.Context, userID, chapterID string, isVIPChapter bool) (bool, error) {
	args := m.Called(ctx, userID, chapterID, isVIPChapter)
	return args.Bool(0), args.Error(1)
}

func (m *MockVIPPermissionService) CheckUserVIPStatus(ctx context.Context, userID string) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockVIPPermissionService) CheckChapterPurchased(ctx context.Context, userID, chapterID string) (bool, error) {
	args := m.Called(ctx, userID, chapterID)
	return args.Bool(0), args.Error(1)
}

func (m *MockVIPPermissionService) GrantVIPAccess(ctx context.Context, userID string, duration time.Duration) error {
	args := m.Called(ctx, userID, duration)
	return args.Error(0)
}

func (m *MockVIPPermissionService) GrantChapterAccess(ctx context.Context, userID, chapterID string) error {
	args := m.Called(ctx, userID, chapterID)
	return args.Error(0)
}

// MockEventBus Mock事件总线
type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Publish(ctx context.Context, event baseInterface.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventBus) PublishAsync(ctx context.Context, event baseInterface.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventBus) Subscribe(eventType string, handler baseInterface.EventHandler) error {
	args := m.Called(eventType, handler)
	return args.Error(0)
}

func (m *MockEventBus) Unsubscribe(eventType string, handlerName string) error {
	args := m.Called(eventType, handlerName)
	return args.Error(0)
}

// =========================
// 测试辅助函数
// =========================

// setupReaderService 创建ReaderService实例用于测试
func setupReaderService() (*ReaderService, *MockReadingProgressRepository, *MockAnnotationRepository, *MockReadingSettingsRepository, *MockChapterService, *MockEventBus, *MockReaderCacheService) {
	mockProgressRepo := new(MockReadingProgressRepository)
	mockAnnotationRepo := new(MockAnnotationRepository)
	mockSettingsRepo := new(MockReadingSettingsRepository)
	mockChapterService := new(MockChapterService)
	mockEventBus := new(MockEventBus)
	mockCacheService := new(MockReaderCacheService)
	mockVIPService := new(MockVIPPermissionService)

	service := NewReaderService(
		mockProgressRepo,
		mockAnnotationRepo,
		mockSettingsRepo,
		mockChapterService,
		mockEventBus,
		mockCacheService,
		mockVIPService,
	)

	return service, mockProgressRepo, mockAnnotationRepo, mockSettingsRepo, mockChapterService, mockEventBus, mockCacheService
}

// =========================
// 阅读进度相关测试
// =========================

// TestReaderService_GetReadingProgress_Success 测试获取阅读进度成功场景
func TestReaderService_GetReadingProgress_Success(t *testing.T) {
	// Arrange
	service, mockProgressRepo, _, _, _, _, _ := setupReaderService()
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()

	progressID := primitive.NewObjectID()
	userObjectID, _ := primitive.ObjectIDFromHex(userID)
	bookObjectID, _ := primitive.ObjectIDFromHex(bookID)

	expectedProgress := &reader.ReadingProgress{
		UserID:   userObjectID,
		BookID:   bookObjectID,
		Progress: 0.5,
	}
	expectedProgress.ID = progressID

	mockProgressRepo.On("GetByUserAndBook", ctx, userID, bookID).
		Return(expectedProgress, nil)

	// Act
	progress, err := service.GetReadingProgress(ctx, userID, bookID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedProgress, progress)
	mockProgressRepo.AssertExpectations(t)
}

// TestReaderService_GetReadingProgress_NotFound 测试获取阅读进度-不存在时返回空进度
func TestReaderService_GetReadingProgress_NotFound(t *testing.T) {
	// Arrange
	service, mockProgressRepo, _, _, _, _, _ := setupReaderService()
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()

	mockProgressRepo.On("GetByUserAndBook", ctx, userID, bookID).
		Return(nil, nil)

	// Act
	progress, err := service.GetReadingProgress(ctx, userID, bookID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, progress)
	assert.Equal(t, userID, progress.UserID)
	assert.Equal(t, bookID, progress.BookID)
	assert.Equal(t, 0.0, progress.Progress)
	mockProgressRepo.AssertExpectations(t)
}

// TestReaderService_GetReadingProgress_RepositoryError 测试获取阅读进度-仓储错误
func TestReaderService_GetReadingProgress_RepositoryError(t *testing.T) {
	// Arrange
	service, mockProgressRepo, _, _, _, _, _ := setupReaderService()
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()

	mockProgressRepo.On("GetByUserAndBook", ctx, userID, bookID).
		Return(nil, errors.New("数据库连接失败"))

	// Act
	progress, err := service.GetReadingProgress(ctx, userID, bookID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, progress)
	assert.Contains(t, err.Error(), "获取阅读进度失败")
	mockProgressRepo.AssertExpectations(t)
}

// TestReaderService_SaveReadingProgress_Success 测试保存阅读进度成功
func TestReaderService_SaveReadingProgress_Success(t *testing.T) {
	// Arrange
	service, mockProgressRepo, _, _, _, mockEventBus, _ := setupReaderService()
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapterID := "chapter123"
	progress := 0.5

	mockProgressRepo.On("SaveProgress", ctx, userID, bookID, chapterID, progress).
		Return(nil)
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil).Maybe()

	// Act
	err := service.SaveReadingProgress(ctx, userID, bookID, chapterID, progress)

	// Assert
	require.NoError(t, err)
	mockProgressRepo.AssertExpectations(t)
}

// TestReaderService_SaveReadingProgress_InvalidProgress 测试保存阅读进度-无效进度值
func TestReaderService_SaveReadingProgress_InvalidProgress(t *testing.T) {
	// Arrange
	service, mockProgressRepo, _, _, _, _, _ := setupReaderService()
	ctx := context.Background()

	tests := []struct {
		name     string
		progress float64
		wantErr  string
	}{
		{"进度小于0", -0.1, "进度值必须在0-1之间"},
		{"进度大于1", 1.1, "进度值必须在0-1之间"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := service.SaveReadingProgress(ctx, "user123", "book123", "chapter123", tt.progress)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
			mockProgressRepo.AssertNotCalled(t, "SaveProgress", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		})
	}
}

// TestReaderService_SaveReadingProgress_RepositoryError 测试保存阅读进度-仓储错误
func TestReaderService_SaveReadingProgress_RepositoryError(t *testing.T) {
	// Arrange
	service, mockProgressRepo, _, _, _, _, _ := setupReaderService()
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapterID := "chapter123"
	progress := 0.5

	mockProgressRepo.On("SaveProgress", ctx, userID, bookID, chapterID, progress).
		Return(errors.New("保存失败"))

	// Act
	err := service.SaveReadingProgress(ctx, userID, bookID, chapterID, progress)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "保存阅读进度失败")
	mockProgressRepo.AssertExpectations(t)
}

// TestReaderService_UpdateReadingTime_Success 测试更新阅读时长成功
func TestReaderService_UpdateReadingTime_Success(t *testing.T) {
	// Arrange
	service, mockProgressRepo, _, _, _, _, _ := setupReaderService()
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	duration := int64(3600) // 1小时

	mockProgressRepo.On("UpdateReadingTime", ctx, userID, bookID, duration).
		Return(nil)

	// Act
	err := service.UpdateReadingTime(ctx, userID, bookID, duration)

	// Assert
	require.NoError(t, err)
	mockProgressRepo.AssertExpectations(t)
}

// TestReaderService_UpdateReadingTime_InvalidDuration 测试更新阅读时长-无效时长
func TestReaderService_UpdateReadingTime_InvalidDuration(t *testing.T) {
	// Arrange
	service, _, _, _, _, _, _ := setupReaderService()
	ctx := context.Background()

	tests := []struct {
		name     string
		duration int64
		wantErr  string
	}{
		{"时长为0", 0, "阅读时长必须大于0"},
		{"时长为负数", -100, "阅读时长必须大于0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := service.UpdateReadingTime(ctx, "user123", "book123", tt.duration)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

// TestReaderService_GetRecentReading_Success 测试获取最近阅读记录成功
func TestReaderService_GetRecentReading_Success(t *testing.T) {
	// Arrange
	service, mockProgressRepo, _, _, _, _, _ := setupReaderService()
	ctx := context.Background()
	userID := "user123"
	limit := 10

	progress1ID := primitive.NewObjectID()
	progress2ID := primitive.NewObjectID()
	userObjectID, _ := primitive.ObjectIDFromHex(userID)
	book1ObjectID, _ := primitive.ObjectIDFromHex("book1")
	book2ObjectID, _ := primitive.ObjectIDFromHex("book2")

	progress1 := &reader.ReadingProgress{
		UserID: userObjectID,
		BookID: book1ObjectID,
	}
	progress1.ID = progress1ID

	progress2 := &reader.ReadingProgress{
		UserID: userObjectID,
		BookID: book2ObjectID,
	}
	progress2.ID = progress2ID

	expectedProgresses := []*reader.ReadingProgress{
		progress1,
		progress2,
	}

	mockProgressRepo.On("GetRecentReadingByUser", ctx, userID, limit).
		Return(expectedProgresses, nil)

	// Act
	progresses, err := service.GetRecentReading(ctx, userID, limit)

	// Assert
	require.NoError(t, err)
	assert.Len(t, progresses, 2)
	mockProgressRepo.AssertExpectations(t)
}

// TestReaderService_GetRecentReading_InvalidLimit 测试获取最近阅读记录-无效限制参数
func TestReaderService_GetRecentReading_InvalidLimit(t *testing.T) {
	// Arrange
	service, mockProgressRepo, _, _, _, _, _ := setupReaderService()
	ctx := context.Background()
	userID := "user123"

	expectedProgresses := []*reader.ReadingProgress{}

	tests := []struct {
		name  string
		limit int
		want  int
	}{
		{"limit为0，使用默认值20", 0, 20},
		{"limit为负数，使用默认值20", -1, 20},
		{"limit超过100，使用默认值20", 101, 20},
		{"limit为10，使用10", 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProgressRepo.On("GetRecentReadingByUser", ctx, userID, tt.want).
				Return(expectedProgresses, nil)

			// Act
			_, err := service.GetRecentReading(ctx, userID, tt.limit)

			// Assert
			require.NoError(t, err)
			mockProgressRepo.AssertExpectations(t)

			// Reset mock for next iteration
			mockProgressRepo.ExpectedCalls = nil
		})
	}
}

// TestReaderService_GetReadingHistory_Success 测试获取阅读历史成功
func TestReaderService_GetReadingHistory_Success(t *testing.T) {
	// Arrange
	service, mockProgressRepo, _, _, _, _, _ := setupReaderService()
	ctx := context.Background()
	userID := "user123"
	page := 1
	size := 20

	progressID := primitive.NewObjectID()
	userObjectID, _ := primitive.ObjectIDFromHex(userID)
	bookObjectID, _ := primitive.ObjectIDFromHex("book1")

	progress := &reader.ReadingProgress{
		UserID: userObjectID,
		BookID: bookObjectID,
	}
	progress.ID = progressID

	expectedProgresses := []*reader.ReadingProgress{
		progress,
	}
	total := int64(100)

	mockProgressRepo.On("GetReadingHistory", ctx, userID, size, 0).
		Return(expectedProgresses, nil)
	mockProgressRepo.On("CountReadingBooks", ctx, userID).
		Return(total, nil)

	// Act
	progresses, totalResult, err := service.GetReadingHistory(ctx, userID, page, size)

	// Assert
	require.NoError(t, err)
	assert.Len(t, progresses, 1)
	assert.Equal(t, total, totalResult)
	mockProgressRepo.AssertExpectations(t)
}

// TestReaderService_GetTotalReadingTime_Success 测试获取总阅读时长成功
func TestReaderService_GetTotalReadingTime_Success(t *testing.T) {
	// Arrange
	service, mockProgressRepo, _, _, _, _, _ := setupReaderService()
	ctx := context.Background()
	userID := "user123"
	expectedTime := int64(7200) // 2小时

	mockProgressRepo.On("GetTotalReadingTime", ctx, userID).
		Return(expectedTime, nil)

	// Act
	total, err := service.GetTotalReadingTime(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedTime, total)
	mockProgressRepo.AssertExpectations(t)
}

// TestReaderService_GetReadingTimeByPeriod_Success 测试获取时间段阅读时长成功
func TestReaderService_GetReadingTimeByPeriod_Success(t *testing.T) {
	// Arrange
	service, mockProgressRepo, _, _, _, _, _ := setupReaderService()
	ctx := context.Background()
	userID := "user123"
	startTime := time.Now().AddDate(0, -1, 0)
	endTime := time.Now()
	expectedTime := int64(3600)

	mockProgressRepo.On("GetReadingTimeByPeriod", ctx, userID, startTime, endTime).
		Return(expectedTime, nil)

	// Act
	total, err := service.GetReadingTimeByPeriod(ctx, userID, startTime, endTime)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedTime, total)
	mockProgressRepo.AssertExpectations(t)
}

// TestReaderService_DeleteReadingProgress_Success 测试删除阅读进度成功
func TestReaderService_DeleteReadingProgress_Success(t *testing.T) {
	// Arrange
	service, mockProgressRepo, _, _, _, _, mockCacheService := setupReaderService()
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	userIDObj, _ := primitive.ObjectIDFromHex(userID)
	bookIDObj, _ := primitive.ObjectIDFromHex(bookID)

	existingProgress := &reader.ReadingProgress{
		UserID: userIDObj,
		BookID: bookIDObj,
	}
	existingProgress.ID = primitive.NewObjectID()

	mockProgressRepo.On("GetByUserAndBook", ctx, userID, bookID).
		Return(existingProgress, nil)
	mockProgressRepo.On("Delete", ctx, existingProgress.ID.Hex()).
		Return(nil)
	mockCacheService.On("InvalidateReadingProgress", ctx, userID, bookID).
		Return(nil)

	// Act
	err := service.DeleteReadingProgress(ctx, userID, bookID)

	// Assert
	require.NoError(t, err)
	mockProgressRepo.AssertExpectations(t)
	mockCacheService.AssertExpectations(t)
}

// TestReaderService_DeleteReadingProgress_EmptyIDs 测试删除阅读进度-空ID
func TestReaderService_DeleteReadingProgress_EmptyIDs(t *testing.T) {
	// Arrange
	service, mockProgressRepo, _, _, _, _, _ := setupReaderService()
	ctx := context.Background()

	tests := []struct {
		name    string
		userID  string
		bookID  string
		wantErr string
	}{
		{"空用户ID", "", "book123", "用户ID和书籍ID不能为空"},
		{"空书籍ID", "user123", "", "用户ID和书籍ID不能为空"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := service.DeleteReadingProgress(ctx, tt.userID, tt.bookID)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
			mockProgressRepo.AssertNotCalled(t, "GetByUserAndBook", mock.Anything, mock.Anything, mock.Anything)
		})
	}
}

// TestReaderService_DeleteReadingProgress_NotFound 测试删除阅读进度-记录不存在
func TestReaderService_DeleteReadingProgress_NotFound(t *testing.T) {
	// Arrange
	service, mockProgressRepo, _, _, _, _, _ := setupReaderService()
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()

	mockProgressRepo.On("GetByUserAndBook", ctx, userID, bookID).
		Return(nil, nil)

	// Act
	err := service.DeleteReadingProgress(ctx, userID, bookID)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "阅读进度记录不存在")
	mockProgressRepo.AssertExpectations(t)
}

// =========================
// 标注相关测试
// =========================

// TestReaderService_CreateAnnotation_Success 测试创建标注成功
func TestReaderService_CreateAnnotation_Success(t *testing.T) {
	// Arrange
	service, _, mockAnnotationRepo, _, _, mockEventBus, _ := setupReaderService()
	ctx := context.Background()

	userObjectID, _ := primitive.ObjectIDFromHex("user123")
	bookObjectID, _ := primitive.ObjectIDFromHex("book123")
	chapterObjectID, _ := primitive.ObjectIDFromHex("chapter123")

	annotation := &reader.Annotation{
		UserID:    userObjectID,
		BookID:    bookObjectID,
		ChapterID: chapterObjectID,
		Type:      "note",
		Note:      "这是笔记内容",
	}

	mockAnnotationRepo.On("Create", ctx, annotation).
		Return(nil)
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil).Maybe()

	// Act
	err := service.CreateAnnotation(ctx, annotation)

	// Assert
	require.NoError(t, err)
	mockAnnotationRepo.AssertExpectations(t)
}

// TestReaderService_CreateAnnotation_ValidationFailed 测试创建标注-验证失败
func TestReaderService_CreateAnnotation_ValidationFailed(t *testing.T) {
	// Arrange
	service, _, mockAnnotationRepo, _, _, _, _ := setupReaderService()
	ctx := context.Background()

	tests := []struct {
		name       string
		annotation *reader.Annotation
		wantErrSub string
	}{
		{
			name: "空用户ID",
			annotation: &reader.Annotation{
				UserID:    primitive.NilObjectID,
				BookID:    primitive.NewObjectID(),
				ChapterID: primitive.NewObjectID(),
				Type:      "note",
			},
			wantErrSub: "用户ID不能为空",
		},
		{
			name: "空书籍ID",
			annotation: &reader.Annotation{
				UserID:    primitive.NewObjectID(),
				BookID:    primitive.NilObjectID,
				ChapterID: primitive.NewObjectID(),
				Type:      "note",
			},
			wantErrSub: "书籍ID不能为空",
		},
		{
			name: "空章节ID",
			annotation: &reader.Annotation{
				UserID:    primitive.NewObjectID(),
				BookID:    primitive.NewObjectID(),
				ChapterID: primitive.NilObjectID,
				Type:      "note",
			},
			wantErrSub: "章节ID不能为空",
		},
		{
			name: "空类型",
			annotation: &reader.Annotation{
				UserID:    primitive.NewObjectID(),
				BookID:    primitive.NewObjectID(),
				ChapterID: primitive.NewObjectID(),
				Type:      "",
			},
			wantErrSub: "标注类型不能为空",
		},
		{
			name: "无效类型",
			annotation: &reader.Annotation{
				UserID:    primitive.NewObjectID(),
				BookID:    primitive.NewObjectID(),
				ChapterID: primitive.NewObjectID(),
				Type:      "invalid",
			},
			wantErrSub: "标注类型必须是bookmark(书签)、highlight(高亮)或note(笔记)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := service.CreateAnnotation(ctx, tt.annotation)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErrSub)
			mockAnnotationRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
		})
	}
}

// TestReaderService_GetAnnotationsByChapter_Success 测试获取章节标注成功
func TestReaderService_GetAnnotationsByChapter_Success(t *testing.T) {
	// Arrange
	service, _, mockAnnotationRepo, _, _, _, _ := setupReaderService()
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapterID := "chapter123"
	userIDObj, _ := primitive.ObjectIDFromHex(userID)
	bookIDObj, _ := primitive.ObjectIDFromHex(bookID)

	expectedAnnotations := []*reader.Annotation{
		{ID: primitive.NewObjectID(), UserID: userIDObj, BookID: bookIDObj, ChapterID: primitive.NewObjectID()},
		{ID: primitive.NewObjectID(), UserID: userIDObj, BookID: bookIDObj, ChapterID: primitive.NewObjectID()},
	}

	mockAnnotationRepo.On("GetByUserAndChapter", ctx, userID, bookID, chapterID).
		Return(expectedAnnotations, nil)

	// Act
	annotations, err := service.GetAnnotationsByChapter(ctx, userID, bookID, chapterID)

	// Assert
	require.NoError(t, err)
	assert.Len(t, annotations, 2)
	mockAnnotationRepo.AssertExpectations(t)
}

// TestReaderService_SearchNotes_Success 测试搜索笔记成功
func TestReaderService_SearchNotes_Success(t *testing.T) {
	// Arrange
	service, _, mockAnnotationRepo, _, _, _, _ := setupReaderService()
	ctx := context.Background()
	userID := "user123"
	userIDObj, _ := primitive.ObjectIDFromHex(userID)
	keyword := "重要"

	expectedNotes := []*reader.Annotation{
		{ID: primitive.NewObjectID(), UserID: userIDObj, Type: "note", Note: "重要内容"},
	}

	mockAnnotationRepo.On("SearchNotes", ctx, userID, keyword).
		Return(expectedNotes, nil)

	// Act
	notes, err := service.SearchNotes(ctx, userID, keyword)

	// Assert
	require.NoError(t, err)
	assert.Len(t, notes, 1)
	mockAnnotationRepo.AssertExpectations(t)
}

// TestReaderService_SearchNotes_EmptyKeyword 测试搜索笔记-空关键词
func TestReaderService_SearchNotes_EmptyKeyword(t *testing.T) {
	// Arrange
	service, _, mockAnnotationRepo, _, _, _, _ := setupReaderService()
	ctx := context.Background()

	// Act
	notes, err := service.SearchNotes(ctx, "user123", "")

	// Assert
	require.Error(t, err)
	assert.Nil(t, notes)
	assert.Contains(t, err.Error(), "搜索关键词不能为空")
	mockAnnotationRepo.AssertNotCalled(t, "SearchNotes", mock.Anything, mock.Anything, mock.Anything)
}

// TestReaderService_GetLatestBookmark_Success 测试获取最新书签成功
func TestReaderService_GetLatestBookmark_Success(t *testing.T) {
	// Arrange
	service, _, mockAnnotationRepo, _, _, _, _ := setupReaderService()
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	userIDObj, _ := primitive.ObjectIDFromHex(userID)
	bookIDObj, _ := primitive.ObjectIDFromHex(bookID)

	expectedBookmark := &reader.Annotation{
		ID:        primitive.NewObjectID(),
		UserID:    userIDObj,
		BookID:    bookIDObj,
		Type:      "bookmark",
		ChapterID: primitive.NewObjectID(),
	}

	mockAnnotationRepo.On("GetLatestBookmark", ctx, userID, bookID).
		Return(expectedBookmark, nil)

	// Act
	bookmark, err := service.GetLatestBookmark(ctx, userID, bookID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedBookmark, bookmark)
	mockAnnotationRepo.AssertExpectations(t)
}

// TestReaderService_GetAnnotationStats_Success 测试获取标注统计成功
func TestReaderService_GetAnnotationStats_Success(t *testing.T) {
	// Arrange
	service, _, mockAnnotationRepo, _, _, _, _ := setupReaderService()
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()

	expectedAnnotations := []*reader.Annotation{
		{ID: primitive.NewObjectID(), Type: "bookmark"},
		{ID: primitive.NewObjectID(), Type: "highlight"},
		{ID: primitive.NewObjectID(), Type: "note"},
		{ID: primitive.NewObjectID(), Type: "highlight"},
	}

	mockAnnotationRepo.On("GetByUserAndBook", ctx, userID, bookID).
		Return(expectedAnnotations, nil)

	// Act
	stats, err := service.GetAnnotationStats(ctx, userID, bookID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 4, stats["totalCount"])
	assert.Equal(t, 1, stats["bookmarkCount"])
	assert.Equal(t, 2, stats["highlightCount"])
	assert.Equal(t, 1, stats["noteCount"])
	mockAnnotationRepo.AssertExpectations(t)
}

// =========================
// 阅读设置相关测试
// =========================

// TestReaderService_GetReadingSettings_Success 测试获取阅读设置成功
func TestReaderService_GetReadingSettings_Success(t *testing.T) {
	// Arrange
	service, _, _, mockSettingsRepo, _, _, mockCacheService := setupReaderService()
	ctx := context.Background()
	userID := "user123"

	expectedSettings := &reader.ReadingSettings{
		UserID:     userID,
		FontSize:   18,
		FontFamily: "sans-serif",
	}

	// 缓存未命中
	mockCacheService.On("GetReadingSettings", ctx, userID).
		Return(nil, nil)
	// 从数据库获取
	mockSettingsRepo.On("GetByUserID", ctx, userID).
		Return(expectedSettings, nil)
	// 缓存设置（服务会自动缓存结果）
	mockCacheService.On("SetReadingSettings", ctx, userID, expectedSettings, time.Hour).
		Return(nil)

	// Act
	settings, err := service.GetReadingSettings(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedSettings, settings)
	mockSettingsRepo.AssertExpectations(t)
	mockCacheService.AssertExpectations(t)
}

// TestReaderService_GetReadingSettings_DefaultSettings 测试获取阅读设置-返回默认设置
func TestReaderService_GetReadingSettings_DefaultSettings(t *testing.T) {
	// Arrange
	service, _, _, mockSettingsRepo, _, _, mockCacheService := setupReaderService()
	ctx := context.Background()
	userID := "user123"

	// 缓存未命中
	mockCacheService.On("GetReadingSettings", ctx, userID).
		Return(nil, nil)
	// 数据库返回nil，应该返回默认设置
	mockSettingsRepo.On("GetByUserID", ctx, userID).
		Return(nil, nil)

	// Act
	settings, err := service.GetReadingSettings(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, settings)
	assert.Equal(t, userID, settings.UserID)
	assert.Equal(t, 16, settings.FontSize) // 默认字体大小
	assert.Equal(t, "serif", settings.FontFamily)
	mockSettingsRepo.AssertExpectations(t)
	mockCacheService.AssertExpectations(t)
}

// TestReaderService_SaveReadingSettings_CreateNew 测试保存阅读设置-创建新设置
func TestReaderService_SaveReadingSettings_CreateNew(t *testing.T) {
	// Arrange
	service, _, _, mockSettingsRepo, _, _, mockCacheService := setupReaderService()
	ctx := context.Background()

	settings := &reader.ReadingSettings{
		UserID:     "user123",
		FontSize:   20,
		FontFamily: "monospace",
	}

	mockSettingsRepo.On("ExistsByUserID", ctx, settings.UserID).
		Return(false, nil)
	mockSettingsRepo.On("Create", ctx, settings).
		Return(nil)
	mockCacheService.On("SetReadingSettings", ctx, settings.UserID, settings, time.Hour).
		Return(nil).Maybe()

	// Act
	err := service.SaveReadingSettings(ctx, settings)

	// Assert
	require.NoError(t, err)
	mockSettingsRepo.AssertExpectations(t)
}

// TestReaderService_SaveReadingSettings_UpdateExisting 测试保存阅读设置-更新现有设置
func TestReaderService_SaveReadingSettings_UpdateExisting(t *testing.T) {
	// Arrange
	service, _, _, mockSettingsRepo, _, _, mockCacheService := setupReaderService()
	ctx := context.Background()

	settings := &reader.ReadingSettings{
		UserID:     "user123",
		FontSize:   20,
		FontFamily: "monospace",
	}

	mockSettingsRepo.On("ExistsByUserID", ctx, settings.UserID).
		Return(true, nil)
	mockSettingsRepo.On("UpdateByUserID", ctx, settings.UserID, settings).
		Return(nil)
	mockCacheService.On("SetReadingSettings", ctx, settings.UserID, settings, time.Hour).
		Return(nil).Maybe()

	// Act
	err := service.SaveReadingSettings(ctx, settings)

	// Assert
	require.NoError(t, err)
	mockSettingsRepo.AssertExpectations(t)
}

// TestReaderService_SaveReadingSettings_EmptyUserID 测试保存阅读设置-空用户ID
func TestReaderService_SaveReadingSettings_EmptyUserID(t *testing.T) {
	// Arrange
	service, _, _, mockSettingsRepo, _, _, _ := setupReaderService()
	ctx := context.Background()

	settings := &reader.ReadingSettings{
		UserID: "",
	}

	// Act
	err := service.SaveReadingSettings(ctx, settings)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "用户ID不能为空")
	mockSettingsRepo.AssertNotCalled(t, "ExistsByUserID", mock.Anything, mock.Anything)
}

// TestReaderService_UpdateReadingSettings_Success 测试更新阅读设置成功
func TestReaderService_UpdateReadingSettings_Success(t *testing.T) {
	// Arrange
	service, _, _, mockSettingsRepo, _, _, mockCacheService := setupReaderService()
	ctx := context.Background()
	userID := "user123"

	existingSettings := &reader.ReadingSettings{
		UserID:     userID,
		FontSize:   16,
		FontFamily: "serif",
		LineHeight: 1.5,
		Theme:      "light",
		Background: "#FFFFFF",
		PageMode:   1,
		AutoScroll: false,
	}

	updates := map[string]interface{}{
		"font_size":   20,
		"font_family": "sans-serif",
		"theme":       "dark",
	}

	mockSettingsRepo.On("GetByUserID", ctx, userID).
		Return(existingSettings, nil)
	mockSettingsRepo.On("UpdateByUserID", ctx, userID, mock.MatchedBy(func(s *reader.ReadingSettings) bool {
		return s.FontSize == 20 && s.FontFamily == "sans-serif" && s.Theme == "dark"
	})).Return(nil)
	mockCacheService.On("SetReadingSettings", ctx, userID, mock.Anything, time.Hour).
		Return(nil).Maybe()

	// Act
	err := service.UpdateReadingSettings(ctx, userID, updates)

	// Assert
	require.NoError(t, err)
	mockSettingsRepo.AssertExpectations(t)
}

// TestReaderService_UpdateReadingSettings_NotFound 测试更新阅读设置-设置不存在
func TestReaderService_UpdateReadingSettings_NotFound(t *testing.T) {
	// Arrange
	service, _, _, mockSettingsRepo, _, _, _ := setupReaderService()
	ctx := context.Background()
	userID := "user123"

	mockSettingsRepo.On("GetByUserID", ctx, userID).
		Return(nil, nil)

	updates := map[string]interface{}{
		"font_size": 20,
	}

	// Act
	err := service.UpdateReadingSettings(ctx, userID, updates)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "阅读设置不存在")
	mockSettingsRepo.AssertExpectations(t)
}

// =========================
// 批量操作测试
// =========================

// TestReaderService_BatchCreateAnnotations_Success 测试批量创建标注成功
func TestReaderService_BatchCreateAnnotations_Success(t *testing.T) {
	// Arrange
	service, _, mockAnnotationRepo, _, _, mockEventBus, _ := setupReaderService()
	ctx := context.Background()

	annotations := []*reader.Annotation{
		{UserID: "user123", BookID: "book123", ChapterID: "chapter123", Type: "note"},
		{UserID: "user123", BookID: "book123", ChapterID: "chapter124", Type: "bookmark"},
	}

	mockAnnotationRepo.On("Create", ctx, mock.AnythingOfType("*reader.Annotation")).
		Return(nil).Times(2)
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil).Times(2)

	// Act
	err := service.BatchCreateAnnotations(ctx, annotations)

	// Assert
	require.NoError(t, err)
	mockAnnotationRepo.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)
}

// TestReaderService_BatchCreateAnnotations_TooMany 测试批量创建标注-数量过多
func TestReaderService_BatchCreateAnnotations_TooMany(t *testing.T) {
	// Arrange
	service, _, mockAnnotationRepo, _, _, _, _ := setupReaderService()
	ctx := context.Background()

	annotations := make([]*reader.Annotation, 51)
	for i := range annotations {
		annotations[i] = &reader.Annotation{
			UserID:    "user123",
			BookID:    "book123",
			ChapterID: "chapter123",
			Type:      "note",
		}
	}

	// Act
	err := service.BatchCreateAnnotations(ctx, annotations)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "批量创建注记数量不能超过50个")
	mockAnnotationRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// TestReaderService_BatchDeleteAnnotations_Success 测试批量删除标注成功
func TestReaderService_BatchDeleteAnnotations_Success(t *testing.T) {
	// Arrange
	service, _, mockAnnotationRepo, _, _, _, _ := setupReaderService()
	ctx := context.Background()

	annotationIDs := []string{"annotation1", "annotation2", "annotation3"}

	mockAnnotationRepo.On("Delete", ctx, mock.AnythingOfType("string")).
		Return(nil).Times(3)

	// Act
	err := service.BatchDeleteAnnotations(ctx, annotationIDs)

	// Assert
	require.NoError(t, err)
	mockAnnotationRepo.AssertExpectations(t)
}

// TestReaderService_BatchDeleteAnnotations_TooMany 测试批量删除标注-数量过多
func TestReaderService_BatchDeleteAnnotations_TooMany(t *testing.T) {
	// Arrange
	service, _, mockAnnotationRepo, _, _, _, _ := setupReaderService()
	ctx := context.Background()

	annotationIDs := make([]string, 101)
	for i := range annotationIDs {
		annotationIDs[i] = primitive.NewObjectID().Hex()
	}

	// Act
	err := service.BatchDeleteAnnotations(ctx, annotationIDs)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "批量删除注记数量不能超过100个")
	mockAnnotationRepo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

// =========================
// BaseService接口测试
// =========================

// TestReaderService_Health_Success 测试健康检查成功
func TestReaderService_Health_Success(t *testing.T) {
	// Arrange
	service, mockProgressRepo, mockAnnotationRepo, _, _, _, _ := setupReaderService()
	ctx := context.Background()

	mockProgressRepo.On("Health", ctx).Return(nil)
	mockAnnotationRepo.On("Health", ctx).Return(nil)

	// Act
	err := service.Health(ctx)

	// Assert
	require.NoError(t, err)
	mockProgressRepo.AssertExpectations(t)
	mockAnnotationRepo.AssertExpectations(t)
}

// TestReaderService_Health_ProgressRepoError 测试健康检查-进度仓储错误
func TestReaderService_Health_ProgressRepoError(t *testing.T) {
	// Arrange
	service, mockProgressRepo, _, _, _, _, _ := setupReaderService()
	ctx := context.Background()

	mockProgressRepo.On("Health", ctx).Return(errors.New("进度仓储健康检查失败"))

	// Act
	err := service.Health(ctx)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "进度Repository健康检查失败")
	mockProgressRepo.AssertExpectations(t)
}

// TestReaderService_GetServiceName 测试获取服务名称
func TestReaderService_GetServiceName(t *testing.T) {
	// Arrange
	service, _, _, _, _, _, _ := setupReaderService()

	// Act
	name := service.GetServiceName()

	// Assert
	assert.Equal(t, "ReaderService", name)
}

// TestReaderService_GetVersion 测试获取服务版本
func TestReaderService_GetVersion(t *testing.T) {
	// Arrange
	service, _, _, _, _, _, _ := setupReaderService()

	// Act
	version := service.GetVersion()

	// Assert
	assert.Equal(t, "1.0.0", version)
}
