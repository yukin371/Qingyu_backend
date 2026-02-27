package admin

import (
	"bytes"
	"context"
	"encoding/csv"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/admin"
	"Qingyu_backend/models/bookstore"
	bookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
	base "Qingyu_backend/repository/interfaces/infrastructure"
)

// ==================== Mock Repositories ====================

// MockBookRepository 书籍仓储Mock
type MockBookRepository struct {
	mock.Mock
}

func (m *MockBookRepository) Create(ctx context.Context, entity *bookstore.Book) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockBookRepository) GetByID(ctx context.Context, id string) (*bookstore.Book, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockBookRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBookRepository) List(ctx context.Context, filter base.Filter) ([]*bookstore.Book, error) {
	args := m.Called(ctx, filter)
	// 获取返回值
	result := args.Get(0)
	if result == nil {
		return nil, nil
	}
	// 获取error - 检查是否有第二个参数
	errArg := args.Get(1)
	if errArg != nil {
		if err, ok := errArg.(error); ok {
			return result.([]*bookstore.Book), err
		}
	}
	// 没有error，返回nil
	return result.([]*bookstore.Book), nil
}

func (m *MockBookRepository) Count(ctx context.Context, filter base.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockBookRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// 实现 BookRepository 的其他方法
func (m *MockBookRepository) GetByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, categoryID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) GetByAuthor(ctx context.Context, author string, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, author, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) GetByAuthorID(ctx context.Context, authorID string, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, authorID, limit, offset)
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

func (m *MockBookRepository) GetNewReleases(ctx context.Context, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) GetFreeBooks(ctx context.Context, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) GetByPriceRange(ctx context.Context, minPrice, maxPrice float64, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, minPrice, maxPrice, limit, offset)
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

func (m *MockBookRepository) SearchWithFilter(ctx context.Context, filter *bookstore.BookFilter) ([]*bookstore.Book, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) CountByCategory(ctx context.Context, categoryID string) (int64, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRepository) CountByAuthor(ctx context.Context, author string) (int64, error) {
	args := m.Called(ctx, author)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRepository) CountByStatus(ctx context.Context, status bookstore.BookStatus) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRepository) CountByFilter(ctx context.Context, filter *bookstore.BookFilter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRepository) BatchUpdateStatus(ctx context.Context, bookIDs []string, status bookstore.BookStatus) error {
	args := m.Called(ctx, bookIDs, status)
	return args.Error(0)
}

func (m *MockBookRepository) BatchUpdateCategory(ctx context.Context, bookIDs []string, categoryIDs []string) error {
	args := m.Called(ctx, bookIDs, categoryIDs)
	return args.Error(0)
}

func (m *MockBookRepository) BatchUpdateRecommended(ctx context.Context, bookIDs []string, isRecommended bool) error {
	args := m.Called(ctx, bookIDs, isRecommended)
	return args.Error(0)
}

func (m *MockBookRepository) BatchUpdateFeatured(ctx context.Context, bookIDs []string, isFeatured bool) error {
	args := m.Called(ctx, bookIDs, isFeatured)
	return args.Error(0)
}

func (m *MockBookRepository) GetStats(ctx context.Context) (*bookstore.BookStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.BookStats), args.Error(1)
}

func (m *MockBookRepository) IncrementViewCount(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockBookRepository) GetYears(ctx context.Context) ([]int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]int), args.Error(1)
}

func (m *MockBookRepository) GetTags(ctx context.Context, categoryID *string) ([]string, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockBookRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

// MockChapterRepository 章节仓储Mock
type MockChapterRepository struct {
	mock.Mock
}

func (m *MockChapterRepository) Create(ctx context.Context, entity *bookstore.Chapter) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockChapterRepository) GetByID(ctx context.Context, id string) (*bookstore.Chapter, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockChapterRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockChapterRepository) List(ctx context.Context, filter base.Filter) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, filter)
	// 获取返回值
	result := args.Get(0)
	if result == nil {
		return nil, nil
	}
	// 获取error - 检查是否有第二个参数
	errArg := args.Get(1)
	if errArg != nil {
		if err, ok := errArg.(error); ok {
			return result.([]*bookstore.Chapter), err
		}
	}
	// 没有error，返回nil
	return result.([]*bookstore.Chapter), nil
}

func (m *MockChapterRepository) Count(ctx context.Context, filter base.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockChapterRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// 实现 ChapterRepository 的其他方法
func (m *MockChapterRepository) GetByBookID(ctx context.Context, bookID string, limit, offset int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetByBookIDAndChapterNum(ctx context.Context, bookID string, chapterNum int) (*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, chapterNum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetByTitle(ctx context.Context, title string, limit, offset int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, title, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetFreeChapters(ctx context.Context, bookID string, limit, offset int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetPaidChapters(ctx context.Context, bookID string, limit, offset int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetPublishedChapters(ctx context.Context, bookID string, limit, offset int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetChapterRange(ctx context.Context, bookID string, startChapter, endChapter int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, startChapter, endChapter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, keyword, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) SearchByFilter(ctx context.Context, filter *bookstoreRepo.ChapterFilter) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) CountByBookID(ctx context.Context, bookID string) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) CountFreeChapters(ctx context.Context, bookID string) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) CountPaidChapters(ctx context.Context, bookID string) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) CountPublishedChapters(ctx context.Context, bookID string) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) GetTotalWordCount(ctx context.Context, bookID string) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) GetPreviousChapter(ctx context.Context, bookID string, chapterNum int) (*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, chapterNum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetNextChapter(ctx context.Context, bookID string, chapterNum int) (*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, chapterNum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetFirstChapter(ctx context.Context, bookID string) (*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetLastChapter(ctx context.Context, bookID string) (*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) BatchUpdatePrice(ctx context.Context, chapterIDs []string, price float64) error {
	args := m.Called(ctx, chapterIDs, price)
	return args.Error(0)
}

func (m *MockChapterRepository) BatchDelete(ctx context.Context, chapterIDs []string) error {
	args := m.Called(ctx, chapterIDs)
	return args.Error(0)
}

func (m *MockChapterRepository) BatchUpdateFreeStatus(ctx context.Context, chapterIDs []string, isFree bool) error {
	args := m.Called(ctx, chapterIDs, isFree)
	return args.Error(0)
}

func (m *MockChapterRepository) BatchUpdatePublishTime(ctx context.Context, chapterIDs []string, publishTime time.Time) error {
	args := m.Called(ctx, chapterIDs, publishTime)
	return args.Error(0)
}

func (m *MockChapterRepository) DeleteByBookID(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockChapterRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

// ==================== Mock Export History Repository ====================

// MockExportHistoryRepositoryForAPI 导出历史仓储Mock（用于API测试）
type MockExportHistoryRepositoryForAPI struct {
	mock.Mock
}

func (m *MockExportHistoryRepositoryForAPI) Create(ctx context.Context, history *admin.ExportHistory) error {
	args := m.Called(ctx, history)
	return args.Error(0)
}

func (m *MockExportHistoryRepositoryForAPI) GetByID(ctx context.Context, id string) (*admin.ExportHistory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*admin.ExportHistory), args.Error(1)
}

func (m *MockExportHistoryRepositoryForAPI) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockExportHistoryRepositoryForAPI) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockExportHistoryRepositoryForAPI) List(ctx context.Context, filter base.Filter) ([]*admin.ExportHistory, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*admin.ExportHistory), args.Error(1)
}

func (m *MockExportHistoryRepositoryForAPI) Count(ctx context.Context, filter base.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockExportHistoryRepositoryForAPI) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockExportHistoryRepositoryForAPI) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockExportHistoryRepositoryForAPI) ListByUser(ctx context.Context, adminID string, page, pageSize int) ([]*admin.ExportHistory, int64, error) {
	args := m.Called(ctx, adminID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*admin.ExportHistory), args.Get(1).(int64), args.Error(2)
}

func (m *MockExportHistoryRepositoryForAPI) ListByDateRange(ctx context.Context, startDate, endDate time.Time, page, pageSize int) ([]*admin.ExportHistory, int64, error) {
	args := m.Called(ctx, startDate, endDate, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*admin.ExportHistory), args.Get(1).(int64), args.Error(2)
}

func (m *MockExportHistoryRepositoryForAPI) CleanOldRecords(ctx context.Context, beforeDate time.Time) (int64, error) {
	args := m.Called(ctx, beforeDate)
	return args.Get(0).(int64), args.Error(1)
}

// ==================== 测试辅助函数 ====================

// setupContentExportAPITestRouter 设置内容导出测试路由
func setupContentExportAPITestRouter(bookRepo *MockBookRepository, chapterRepo *MockChapterRepository, exportHistoryRepo *MockExportHistoryRepositoryForAPI) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	api := NewContentExportAPI(bookRepo, chapterRepo, exportHistoryRepo)

	v1 := r.Group("/api/v1/admin")
	{
		v1.GET("/content/books/export", api.ExportBooks)
		v1.GET("/content/chapters/export", api.ExportChapters)
		v1.GET("/content/books/export/template", api.GetBookExportTemplate)
		v1.GET("/content/chapters/export/template", api.GetChapterExportTemplate)
	}

	return r
}

// 创建测试书籍数据
func createTestBooks(count int) []*bookstore.Book {
	books := make([]*bookstore.Book, count)
	now := time.Now()
	for i := 0; i < count; i++ {
		books[i] = &bookstore.Book{
			Title:        "测试书籍",
			Author:       "测试作者",
			Introduction: "这是测试简介",
			Status:       bookstore.BookStatusOngoing,
			WordCount:    10000,
			ChapterCount: 10,
		}
		books[i].ID = primitive.NewObjectID()
		books[i].CreatedAt = now
		books[i].UpdatedAt = now
	}
	return books
}

// 创建测试章节数据
func createTestChapters(count int) []*bookstore.Chapter {
	chapters := make([]*bookstore.Chapter, count)
	now := time.Now()
	for i := 0; i < count; i++ {
		chapters[i] = &bookstore.Chapter{
			ID:         "chapter-id",
			BookID:     "book-id",
			Title:      "测试章节",
			ChapterNum: i + 1,
			WordCount:  1000,
			IsFree:     true,
			Price:      0,
			PublishTime: now,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
	}
	return chapters
}

// ==================== TestContentExportAPI_ExportBooks_Success ====================

func TestContentExportAPI_ExportBooksToCSV_Success(t *testing.T) {
	// Given
	mockBookRepo := new(MockBookRepository)
	mockChapterRepo := new(MockChapterRepository)
	mockExportHistoryRepo := new(MockExportHistoryRepositoryForAPI)
	router := setupContentExportAPITestRouter(mockBookRepo, mockChapterRepo, mockExportHistoryRepo)

	books := createTestBooks(2)
	mockBookRepo.On("List", mock.Anything, mock.Anything).Return(books, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/content/books/export?format=csv", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/csv")

	// 验证CSV内容
	reader := csv.NewReader(strings.NewReader(w.Body.String()))
	records, err := reader.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(records)) // Header + 2 rows
	assert.Equal(t, "ID", records[0][0])
	assert.Equal(t, "书名", records[0][1])

	mockBookRepo.AssertExpectations(t)
}

// ==================== TestContentExportAPI_ExportBooksToExcel_Success ====================

func TestContentExportAPI_ExportBooksToExcel_Success(t *testing.T) {
	// Given
	mockBookRepo := new(MockBookRepository)
	mockChapterRepo := new(MockChapterRepository)
	mockExportHistoryRepo := new(MockExportHistoryRepositoryForAPI)
	router := setupContentExportAPITestRouter(mockBookRepo, mockChapterRepo, mockExportHistoryRepo)

	books := createTestBooks(2)
	mockBookRepo.On("List", mock.Anything, mock.Anything).Return(books, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/content/books/export?format=excel", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", w.Header().Get("Content-Type"))

	// 验证Excel文件
	_, err := excelize.OpenReader(bytes.NewReader(w.Body.Bytes()))
	assert.NoError(t, err)

	mockBookRepo.AssertExpectations(t)
}

// ==================== TestContentExportAPI_ExportChapters_Success ====================

func TestContentExportAPI_ExportChaptersToCSV_Success(t *testing.T) {
	// Given
	mockBookRepo := new(MockBookRepository)
	mockChapterRepo := new(MockChapterRepository)
	mockExportHistoryRepo := new(MockExportHistoryRepositoryForAPI)
	router := setupContentExportAPITestRouter(mockBookRepo, mockChapterRepo, mockExportHistoryRepo)

	chapters := createTestChapters(2)
	mockChapterRepo.On("List", mock.Anything, mock.Anything).Return(chapters, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/content/chapters/export?format=csv", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/csv")

	// 验证CSV内容
	reader := csv.NewReader(strings.NewReader(w.Body.String()))
	records, err := reader.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(records)) // Header + 2 rows

	mockChapterRepo.AssertExpectations(t)
}

// ==================== TestContentExportAPI_ExportWithDateRange_Success ====================

func TestContentExportAPI_ExportBooksFromDateRange_Success(t *testing.T) {
	// Given
	mockBookRepo := new(MockBookRepository)
	mockChapterRepo := new(MockChapterRepository)
	mockExportHistoryRepo := new(MockExportHistoryRepositoryForAPI)
	router := setupContentExportAPITestRouter(mockBookRepo, mockChapterRepo, mockExportHistoryRepo)

	books := createTestBooks(1)
	mockBookRepo.On("List", mock.Anything, mock.MatchedBy(func(f base.Filter) bool {
		// 验证时间范围过滤器
		conditions := f.GetConditions()
		_, hasStart := conditions["createdAt.$gte"]
		_, hasEnd := conditions["createdAt.$lte"]
		return hasStart && hasEnd
	})).Return(books, int64(1), nil)

	// When
	startDate := time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")
	req, _ := http.NewRequest("GET", "/api/v1/admin/content/books/export?format=csv&start_date="+startDate+"&end_date="+endDate, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockBookRepo.AssertExpectations(t)
}

// ==================== TestContentExportAPI_ExportEmpty_Success ====================

func TestContentExportAPI_ExportEmptyBooks_Success(t *testing.T) {
	// Given
	mockBookRepo := new(MockBookRepository)
	mockChapterRepo := new(MockChapterRepository)
	mockExportHistoryRepo := new(MockExportHistoryRepositoryForAPI)
	router := setupContentExportAPITestRouter(mockBookRepo, mockChapterRepo, mockExportHistoryRepo)

	emptyBooks := []*bookstore.Book{}
	mockBookRepo.On("List", mock.Anything, mock.Anything).Return(emptyBooks, int64(0), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/content/books/export?format=csv", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/csv")

	// 验证CSV只有header
	reader := csv.NewReader(strings.NewReader(w.Body.String()))
	records, err := reader.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(records)) // Only header

	mockBookRepo.AssertExpectations(t)
}

// ==================== TestContentExportAPI_ExportBooksInvalidFormat_Error ====================

func TestContentExportAPI_ExportBooksInvalidFormat_Error(t *testing.T) {
	// Given
	mockBookRepo := new(MockBookRepository)
	mockChapterRepo := new(MockChapterRepository)
	mockExportHistoryRepo := new(MockExportHistoryRepositoryForAPI)
	router := setupContentExportAPITestRouter(mockBookRepo, mockChapterRepo, mockExportHistoryRepo)

	// When - 无效格式，不应该调用repository
	req, _ := http.NewRequest("GET", "/api/v1/admin/content/books/export?format=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==================== TestContentExportAPI_GetBookExportTemplate_Success ====================

func TestContentExportAPI_GetBookExportTemplate_Success(t *testing.T) {
	// Given
	mockBookRepo := new(MockBookRepository)
	mockChapterRepo := new(MockChapterRepository)
	mockExportHistoryRepo := new(MockExportHistoryRepositoryForAPI)
	router := setupContentExportAPITestRouter(mockBookRepo, mockChapterRepo, mockExportHistoryRepo)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/content/books/export/template", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/csv")

	// 验证CSV包含表头
	reader := csv.NewReader(strings.NewReader(w.Body.String()))
	records, err := reader.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(records)) // Only header
	assert.Greater(t, len(records[0]), 0)
}

// ==================== TestContentExportAPI_GetChapterExportTemplate_Success ====================

func TestContentExportAPI_GetChapterExportTemplate_Success(t *testing.T) {
	// Given
	mockBookRepo := new(MockBookRepository)
	mockChapterRepo := new(MockChapterRepository)
	mockExportHistoryRepo := new(MockExportHistoryRepositoryForAPI)
	router := setupContentExportAPITestRouter(mockBookRepo, mockChapterRepo, mockExportHistoryRepo)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/content/chapters/export/template", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/csv")

	// 验证CSV包含表头
	reader := csv.NewReader(strings.NewReader(w.Body.String()))
	records, err := reader.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(records)) // Only header
}

// ==================== TestContentExportAPI_ExportLargeDataset_Success ====================

func TestContentExportAPI_ExportLargeBooksDataset_Success(t *testing.T) {
	// Given
	mockBookRepo := new(MockBookRepository)
	mockChapterRepo := new(MockChapterRepository)
	mockExportHistoryRepo := new(MockExportHistoryRepositoryForAPI)
	router := setupContentExportAPITestRouter(mockBookRepo, mockChapterRepo, mockExportHistoryRepo)

	// 创建1000条数据
	books := createTestBooks(1000)
	mockBookRepo.On("List", mock.Anything, mock.Anything).Return(books, int64(1000), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/content/books/export?format=csv", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证CSV内容
	reader := csv.NewReader(strings.NewReader(w.Body.String()))
	records, err := reader.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, 1001, len(records)) // Header + 1000 rows

	mockBookRepo.AssertExpectations(t)
}

// ==================== TestContentExportAPI_ExportWithFilter_Success ====================

func TestContentExportAPI_ExportBooksWithFilter_Success(t *testing.T) {
	// Given
	mockBookRepo := new(MockBookRepository)
	mockChapterRepo := new(MockChapterRepository)
	mockExportHistoryRepo := new(MockExportHistoryRepositoryForAPI)
	router := setupContentExportAPITestRouter(mockBookRepo, mockChapterRepo, mockExportHistoryRepo)

	books := createTestBooks(1)
	mockBookRepo.On("List", mock.Anything, mock.MatchedBy(func(f base.Filter) bool {
		conditions := f.GetConditions()
		_, hasStatus := conditions["status"]
		return hasStatus
	})).Return(books, int64(1), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/content/books/export?format=csv&status=ongoing", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockBookRepo.AssertExpectations(t)
}

// ==================== TestContentExportAPI_ExportChaptersByBook_Success ====================

func TestContentExportAPI_ExportChaptersByBook_Success(t *testing.T) {
	// Given
	mockBookRepo := new(MockBookRepository)
	mockChapterRepo := new(MockChapterRepository)
	mockExportHistoryRepo := new(MockExportHistoryRepositoryForAPI)
	router := setupContentExportAPITestRouter(mockBookRepo, mockChapterRepo, mockExportHistoryRepo)

	chapters := createTestChapters(2)
	mockChapterRepo.On("List", mock.Anything, mock.MatchedBy(func(f base.Filter) bool {
		conditions := f.GetConditions()
		bookID, hasBookID := conditions["book_id"]
		return hasBookID && bookID == "test-book-id"
	})).Return(chapters, int64(2), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/content/chapters/export?format=csv&book_id=test-book-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockChapterRepo.AssertExpectations(t)
}
