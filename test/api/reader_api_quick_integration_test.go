package api

import (
	"Qingyu_backend/models/reader"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	readerAPI "Qingyu_backend/api/v1/reader"
	infrastructure "Qingyu_backend/repository/interfaces/infrastructure"
	readingRepo "Qingyu_backend/repository/interfaces/reading"
	"Qingyu_backend/service/interfaces/base"
	"Qingyu_backend/service/reading"
)

// === In-file minimal mocks for repositories and event bus ===

type stubEventBus struct{}

func (s *stubEventBus) Subscribe(eventType string, handler base.EventHandler) error { return nil }
func (s *stubEventBus) Unsubscribe(eventType string, handlerName string) error      { return nil }
func (s *stubEventBus) Publish(ctx context.Context, event base.Event) error         { return nil }
func (s *stubEventBus) PublishAsync(ctx context.Context, event base.Event) error    { return nil }

type mockChapterRepo struct{ contentByID map[string]string }

func (m *mockChapterRepo) Create(ctx context.Context, chapter *reader.Chapter) error { return nil }
func (m *mockChapterRepo) GetByID(ctx context.Context, id string) (*reader.Chapter, error) {
	return &reader.Chapter{ID: id, BookID: "book1", Title: "第一章", ChapterNum: 1}, nil
}
func (m *mockChapterRepo) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return nil
}
func (m *mockChapterRepo) Delete(ctx context.Context, id string) error { return nil }
func (m *mockChapterRepo) GetByBookID(ctx context.Context, bookID string) ([]*reader.Chapter, error) {
	return nil, nil
}
func (m *mockChapterRepo) GetByBookIDWithPagination(ctx context.Context, bookID string, limit, offset int64) ([]*reader.Chapter, error) {
	return []*reader.Chapter{{ID: "c1", BookID: bookID, Title: "第一章", ChapterNum: 1}}, nil
}
func (m *mockChapterRepo) GetByChapterNum(ctx context.Context, bookID string, chapterNum int) (*reader.Chapter, error) {
	return &reader.Chapter{ID: "c1", BookID: bookID, Title: "第一章", ChapterNum: chapterNum}, nil
}
func (m *mockChapterRepo) GetPrevChapter(ctx context.Context, bookID string, currentChapterNum int) (*reader.Chapter, error) {
	return nil, nil
}
func (m *mockChapterRepo) GetNextChapter(ctx context.Context, bookID string, currentChapterNum int) (*reader.Chapter, error) {
	return nil, nil
}
func (m *mockChapterRepo) GetFirstChapter(ctx context.Context, bookID string) (*reader.Chapter, error) {
	return &reader.Chapter{ID: "c1", BookID: bookID, Title: "第一章", ChapterNum: 1}, nil
}
func (m *mockChapterRepo) GetLastChapter(ctx context.Context, bookID string) (*reader.Chapter, error) {
	return nil, nil
}
func (m *mockChapterRepo) GetPublishedChapters(ctx context.Context, bookID string) ([]*reader.Chapter, error) {
	return nil, nil
}
func (m *mockChapterRepo) GetVIPChapters(ctx context.Context, bookID string) ([]*reader.Chapter, error) {
	return nil, nil
}
func (m *mockChapterRepo) GetFreeChapters(ctx context.Context, bookID string) ([]*reader.Chapter, error) {
	return nil, nil
}
func (m *mockChapterRepo) CountByBookID(ctx context.Context, bookID string) (int64, error) {
	return 1, nil
}
func (m *mockChapterRepo) CountByStatus(ctx context.Context, bookID string, status int) (int64, error) {
	return 0, nil
}
func (m *mockChapterRepo) CountVIPChapters(ctx context.Context, bookID string) (int64, error) {
	return 0, nil
}
func (m *mockChapterRepo) BatchCreate(ctx context.Context, chapters []*reader.Chapter) error {
	return nil
}
func (m *mockChapterRepo) BatchUpdateStatus(ctx context.Context, chapterIDs []string, status int) error {
	return nil
}
func (m *mockChapterRepo) BatchDelete(ctx context.Context, chapterIDs []string) error { return nil }
func (m *mockChapterRepo) CheckVIPAccess(ctx context.Context, chapterID string) (bool, error) {
	return false, nil
}
func (m *mockChapterRepo) GetChapterPrice(ctx context.Context, chapterID string) (int64, error) {
	return 0, nil
}
func (m *mockChapterRepo) GetChapterContent(ctx context.Context, chapterID string) (string, error) {
	if m.contentByID != nil {
		if v, ok := m.contentByID[chapterID]; ok {
			return v, nil
		}
	}
	return "章节内容", nil
}
func (m *mockChapterRepo) UpdateChapterContent(ctx context.Context, chapterID string, content string) error {
	return nil
}
func (m *mockChapterRepo) Health(ctx context.Context) error { return nil }

type mockProgressRepo struct {
	readingRepo.ReadingProgressRepository
}
type mockAnnotationRepo struct {
	readingRepo.AnnotationRepository
}

type mockSettingsRepo struct{}

func (m *mockSettingsRepo) Create(ctx context.Context, settings *reader.ReadingSettings) error {
	return nil
}
func (m *mockSettingsRepo) GetByID(ctx context.Context, id string) (*reader.ReadingSettings, error) {
	return nil, nil
}
func (m *mockSettingsRepo) GetByUserID(ctx context.Context, userID string) (*reader.ReadingSettings, error) {
	return &reader.ReadingSettings{UserID: userID, FontSize: 16, FontFamily: "serif", LineHeight: 1.8, Theme: "light"}, nil
}
func (m *mockSettingsRepo) CreateDefaultSettings(ctx context.Context, userID string) (*reader.ReadingSettings, error) {
	return &reader.ReadingSettings{UserID: userID, FontSize: 16, FontFamily: "serif", LineHeight: 1.8, Theme: "light"}, nil
}
func (m *mockSettingsRepo) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return nil
}
func (m *mockSettingsRepo) UpdateByUserID(ctx context.Context, userID string, settings *reader.ReadingSettings) error {
	return nil
}
func (m *mockSettingsRepo) Delete(ctx context.Context, id string) error { return nil }
func (m *mockSettingsRepo) ExistsByUserID(ctx context.Context, userID string) (bool, error) {
	return true, nil
}
func (m *mockSettingsRepo) Health(ctx context.Context) error { return nil }

// implement base CRUDRepository extra methods: List, Count, Exists
func (m *mockSettingsRepo) List(ctx context.Context, filter infrastructure.Filter) ([]*reader.ReadingSettings, error) {
	return []*reader.ReadingSettings{}, nil
}
func (m *mockSettingsRepo) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	return 0, nil
}
func (m *mockSettingsRepo) Exists(ctx context.Context, id string) (bool, error) { return true, nil }

// Minimal no-op cache and vip services to satisfy ReaderService constructor

type noCache struct{}

func (n *noCache) GetChapterContent(ctx context.Context, chapterID string) (string, error) {
	return "", nil
}
func (n *noCache) SetChapterContent(ctx context.Context, chapterID string, content string, expiration time.Duration) error {
	return nil
}

// match interface name InvalidateChapterContent
func (n *noCache) InvalidateChapterContent(ctx context.Context, chapterID string) error { return nil }
func (n *noCache) GetReadingSettings(ctx context.Context, userID string) (*reader.ReadingSettings, error) {
	return nil, nil
}
func (n *noCache) SetReadingSettings(ctx context.Context, userID string, settings *reader.ReadingSettings, expiration time.Duration) error {
	return nil
}
func (n *noCache) InvalidateReadingSettings(ctx context.Context, userID string) error { return nil }

// complete ReaderCacheService methods used by interface
func (n *noCache) GetChapter(ctx context.Context, chapterID string) (*reader.Chapter, error) {
	return nil, nil
}
func (n *noCache) SetChapter(ctx context.Context, chapterID string, chapter *reader.Chapter, expiration time.Duration) error {
	return nil
}
func (n *noCache) InvalidateChapter(ctx context.Context, chapterID string) error { return nil }
func (n *noCache) GetReadingProgress(ctx context.Context, userID, bookID string) (*reader.ReadingProgress, error) {
	return nil, nil
}
func (n *noCache) SetReadingProgress(ctx context.Context, userID, bookID string, progress *reader.ReadingProgress, expiration time.Duration) error {
	return nil
}
func (n *noCache) InvalidateReadingProgress(ctx context.Context, userID, bookID string) error {
	return nil
}
func (n *noCache) InvalidateBookChapters(ctx context.Context, bookID string) error { return nil }
func (n *noCache) InvalidateUserData(ctx context.Context, userID string) error     { return nil }

// VIP service stub

type noVIP struct{}

func (n *noVIP) CheckVIPAccess(ctx context.Context, userID, chapterID string, isVIPChapter bool) (bool, error) {
	return true, nil
}
func (n *noVIP) CheckUserVIPStatus(ctx context.Context, userID string) (bool, error) {
	return true, nil
}
func (n *noVIP) CheckChapterPurchased(ctx context.Context, userID, chapterID string) (bool, error) {
	return true, nil
}
func (n *noVIP) GrantVIPAccess(ctx context.Context, userID string, duration time.Duration) error {
	return nil
}
func (n *noVIP) GrantChapterAccess(ctx context.Context, userID, chapterID string) error { return nil }

// === Quick integration test ===

func Test_QuickIntegration_GetChapterContent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Assemble real ReaderService with mocks
	chapRepo := &mockChapterRepo{contentByID: map[string]string{"c1": "集成测试内容"}}
	progRepo := &mockProgressRepo{}
	annoRepo := &mockAnnotationRepo{}
	setRepo := &mockSettingsRepo{}
	eventBus := &stubEventBus{}
	cache := &noCache{}
	vip := &noVIP{}

	readerSvc := reading.NewReaderService(chapRepo, progRepo, annoRepo, setRepo, eventBus, cache, vip)

	api := readerAPI.NewChaptersAPI(readerSvc)
	r := gin.New()
	// Inject fake auth (set userId in context)
	r.GET("/api/v1/reader/chapters/:id/content", func(c *gin.Context) { c.Set("userId", "u1"); c.Next() }, api.GetChapterContent)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reader/chapters/c1/content", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
