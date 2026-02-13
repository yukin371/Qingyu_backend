package reader

import (
	"context"
	"testing"
	"time"

	readerModel "Qingyu_backend/models/reader"
)

// ============================================================================
// 编译时检查
// ============================================================================

// 确保所有 Port 接口都包含 BaseService 的方法
var _ BaseService = (*MockReadingProgressPort)(nil)
var _ BaseService = (*MockAnnotationPort)(nil)
var _ BaseService = (*MockChapterContentPort)(nil)
var _ BaseService = (*MockReaderSettingsPort)(nil)
var _ BaseService = (*MockReaderSyncPort)(nil)

// ============================================================================
// Mock 实现 - 用于编译时检查
// ============================================================================

// MockReadingProgressPort Mock 实现用于编译时检查
type MockReadingProgressPort struct{}

func (m *MockReadingProgressPort) Initialize(ctx context.Context) error { return nil }
func (m *MockReadingProgressPort) Health(ctx context.Context) error     { return nil }
func (m *MockReadingProgressPort) Close(ctx context.Context) error      { return nil }
func (m *MockReadingProgressPort) GetServiceName() string               { return "MockReadingProgressPort" }
func (m *MockReadingProgressPort) GetVersion() string                   { return "1.0.0" }
func (m *MockReadingProgressPort) GetReadingProgress(ctx context.Context, userID, bookID string) (*readerModel.ReadingProgress, error) {
	return nil, nil
}
func (m *MockReadingProgressPort) SaveReadingProgress(ctx context.Context, req *SaveReadingProgressRequest) error {
	return nil
}
func (m *MockReadingProgressPort) UpdateReadingTime(ctx context.Context, req *UpdateReadingTimeRequest) error {
	return nil
}
func (m *MockReadingProgressPort) GetRecentReading(ctx context.Context, userID string, limit int) ([]*readerModel.ReadingProgress, error) {
	return nil, nil
}
func (m *MockReadingProgressPort) GetReadingHistory(ctx context.Context, req *GetReadingHistoryRequest) (*GetReadingHistoryResponse, error) {
	return nil, nil
}
func (m *MockReadingProgressPort) GetTotalReadingTime(ctx context.Context, userID string) (int64, error) {
	return 0, nil
}
func (m *MockReadingProgressPort) GetReadingTimeByPeriod(ctx context.Context, req *GetReadingTimeByPeriodRequest) (int64, error) {
	return 0, nil
}
func (m *MockReadingProgressPort) GetUnfinishedBooks(ctx context.Context, userID string) ([]*readerModel.ReadingProgress, error) {
	return nil, nil
}
func (m *MockReadingProgressPort) GetFinishedBooks(ctx context.Context, userID string) ([]*readerModel.ReadingProgress, error) {
	return nil, nil
}
func (m *MockReadingProgressPort) DeleteReadingProgress(ctx context.Context, userID, bookID string) error {
	return nil
}
func (m *MockReadingProgressPort) UpdateBookStatus(ctx context.Context, req *UpdateBookStatusRequest) error {
	return nil
}
func (m *MockReadingProgressPort) BatchUpdateBookStatus(ctx context.Context, req *BatchUpdateBookStatusRequest) error {
	return nil
}

// MockAnnotationPort Mock 实现用于编译时检查
type MockAnnotationPort struct{}

func (m *MockAnnotationPort) Initialize(ctx context.Context) error { return nil }
func (m *MockAnnotationPort) Health(ctx context.Context) error     { return nil }
func (m *MockAnnotationPort) Close(ctx context.Context) error      { return nil }
func (m *MockAnnotationPort) GetServiceName() string               { return "MockAnnotationPort" }
func (m *MockAnnotationPort) GetVersion() string                   { return "1.0.0" }
func (m *MockAnnotationPort) CreateAnnotation(ctx context.Context, annotation *readerModel.Annotation) error {
	return nil
}
func (m *MockAnnotationPort) UpdateAnnotation(ctx context.Context, annotationID string, updates map[string]interface{}) error {
	return nil
}
func (m *MockAnnotationPort) DeleteAnnotation(ctx context.Context, annotationID string) error {
	return nil
}
func (m *MockAnnotationPort) GetAnnotationsByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*readerModel.Annotation, error) {
	return nil, nil
}
func (m *MockAnnotationPort) GetAnnotationsByBook(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error) {
	return nil, nil
}
func (m *MockAnnotationPort) GetNotes(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error) {
	return nil, nil
}
func (m *MockAnnotationPort) SearchNotes(ctx context.Context, userID, keyword string) ([]*readerModel.Annotation, error) {
	return nil, nil
}
func (m *MockAnnotationPort) GetBookmarks(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error) {
	return nil, nil
}
func (m *MockAnnotationPort) GetLatestBookmark(ctx context.Context, userID, bookID string) (*readerModel.Annotation, error) {
	return nil, nil
}
func (m *MockAnnotationPort) GetHighlights(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error) {
	return nil, nil
}
func (m *MockAnnotationPort) GetRecentAnnotations(ctx context.Context, userID string, limit int) ([]*readerModel.Annotation, error) {
	return nil, nil
}
func (m *MockAnnotationPort) GetPublicAnnotations(ctx context.Context, bookID, chapterID string) ([]*readerModel.Annotation, error) {
	return nil, nil
}
func (m *MockAnnotationPort) GetAnnotationStats(ctx context.Context, userID, bookID string) (map[string]interface{}, error) {
	return nil, nil
}
func (m *MockAnnotationPort) BatchCreateAnnotations(ctx context.Context, annotations []*readerModel.Annotation) error {
	return nil
}
func (m *MockAnnotationPort) BatchDeleteAnnotations(ctx context.Context, annotationIDs []string) error {
	return nil
}

// MockChapterContentPort Mock 实现用于编译时检查
type MockChapterContentPort struct{}

func (m *MockChapterContentPort) Initialize(ctx context.Context) error { return nil }
func (m *MockChapterContentPort) Health(ctx context.Context) error     { return nil }
func (m *MockChapterContentPort) Close(ctx context.Context) error      { return nil }
func (m *MockChapterContentPort) GetServiceName() string               { return "MockChapterContentPort" }
func (m *MockChapterContentPort) GetVersion() string                   { return "1.0.0" }
func (m *MockChapterContentPort) GetChapterContent(ctx context.Context, userID, chapterID string) (string, error) {
	return "", nil
}
func (m *MockChapterContentPort) GetChapterByID(ctx context.Context, chapterID string) (interface{}, error) {
	return nil, nil
}
func (m *MockChapterContentPort) GetBookChapters(ctx context.Context, bookID string, page, size int) (interface{}, int64, error) {
	return nil, 0, nil
}
func (m *MockChapterContentPort) GetChapterContentWithProgress(ctx context.Context, req *GetChapterContentRequest) (*ChapterContentResponse, error) {
	return nil, nil
}
func (m *MockChapterContentPort) GetChapterByNumber(ctx context.Context, userID, bookID string, chapterNum int) (*ChapterContentResponse, error) {
	return nil, nil
}
func (m *MockChapterContentPort) GetNextChapter(ctx context.Context, userID, bookID, chapterID string) (*ChapterInfo, error) {
	return nil, nil
}
func (m *MockChapterContentPort) GetPreviousChapter(ctx context.Context, userID, bookID, chapterID string) (*ChapterInfo, error) {
	return nil, nil
}
func (m *MockChapterContentPort) GetChapterList(ctx context.Context, userID, bookID string, page, size int) (*ChapterListResponse, error) {
	return nil, nil
}
func (m *MockChapterContentPort) GetChapterInfo(ctx context.Context, userID, chapterID string) (*ChapterInfo, error) {
	return nil, nil
}

// MockReaderSettingsPort Mock 实现用于编译时检查
type MockReaderSettingsPort struct{}

func (m *MockReaderSettingsPort) Initialize(ctx context.Context) error { return nil }
func (m *MockReaderSettingsPort) Health(ctx context.Context) error     { return nil }
func (m *MockReaderSettingsPort) Close(ctx context.Context) error      { return nil }
func (m *MockReaderSettingsPort) GetServiceName() string               { return "MockReaderSettingsPort" }
func (m *MockReaderSettingsPort) GetVersion() string                   { return "1.0.0" }
func (m *MockReaderSettingsPort) GetReadingSettings(ctx context.Context, userID string) (*readerModel.ReadingSettings, error) {
	return nil, nil
}
func (m *MockReaderSettingsPort) SaveReadingSettings(ctx context.Context, settings *readerModel.ReadingSettings) error {
	return nil
}
func (m *MockReaderSettingsPort) UpdateReadingSettings(ctx context.Context, userID string, updates map[string]interface{}) error {
	return nil
}

// MockReaderSyncPort Mock 实现用于编译时检查
type MockReaderSyncPort struct{}

func (m *MockReaderSyncPort) Initialize(ctx context.Context) error { return nil }
func (m *MockReaderSyncPort) Health(ctx context.Context) error     { return nil }
func (m *MockReaderSyncPort) Close(ctx context.Context) error      { return nil }
func (m *MockReaderSyncPort) GetServiceName() string               { return "MockReaderSyncPort" }
func (m *MockReaderSyncPort) GetVersion() string                   { return "1.0.0" }
func (m *MockReaderSyncPort) SyncAnnotations(ctx context.Context, req *SyncAnnotationsRequest) (*SyncAnnotationsResponse, error) {
	return nil, nil
}

// ============================================================================
// 测试 Port 接口的方法签名
// ============================================================================

func TestPortInterfacesMethodSignatures(t *testing.T) {
	t.Run("ReadingProgressPort has 13 methods", func(t *testing.T) {
		var port ReadingProgressPort = &MockReadingProgressPort{}

		// 测试所有方法都能正常调用
		ctx := context.Background()

		_ = port.Initialize(ctx)
		_ = port.Health(ctx)
		_ = port.Close(ctx)
		_ = port.GetServiceName()
		_ = port.GetVersion()
		_, _ = port.GetReadingProgress(ctx, "user1", "book1")
		_ = port.SaveReadingProgress(ctx, &SaveReadingProgressRequest{})
		_ = port.UpdateReadingTime(ctx, &UpdateReadingTimeRequest{})
		_, _ = port.GetRecentReading(ctx, "user1", 10)
		_, _ = port.GetReadingHistory(ctx, &GetReadingHistoryRequest{})
		_, _ = port.GetTotalReadingTime(ctx, "user1")
		_, _ = port.GetReadingTimeByPeriod(ctx, &GetReadingTimeByPeriodRequest{})
		_, _ = port.GetUnfinishedBooks(ctx, "user1")
		_, _ = port.GetFinishedBooks(ctx, "user1")
		_ = port.DeleteReadingProgress(ctx, "user1", "book1")
		_ = port.UpdateBookStatus(ctx, &UpdateBookStatusRequest{})
		_ = port.BatchUpdateBookStatus(ctx, &BatchUpdateBookStatusRequest{})
	})

	t.Run("AnnotationPort has 16 methods", func(t *testing.T) {
		var port AnnotationPort = &MockAnnotationPort{}
		ctx := context.Background()

		port.Initialize(ctx)
		port.Health(ctx)
		port.Close(ctx)
		_ = port.GetServiceName()
		_ = port.GetVersion()
		port.CreateAnnotation(ctx, &readerModel.Annotation{})
		port.UpdateAnnotation(ctx, "id1", map[string]interface{}{})
		port.DeleteAnnotation(ctx, "id1")
		port.GetAnnotationsByChapter(ctx, "user1", "book1", "chapter1")
		port.GetAnnotationsByBook(ctx, "user1", "book1")
		port.GetNotes(ctx, "user1", "book1")
		port.SearchNotes(ctx, "user1", "keyword")
		port.GetBookmarks(ctx, "user1", "book1")
		port.GetLatestBookmark(ctx, "user1", "book1")
		port.GetHighlights(ctx, "user1", "book1")
		port.GetRecentAnnotations(ctx, "user1", 10)
		port.GetPublicAnnotations(ctx, "book1", "chapter1")
		port.GetAnnotationStats(ctx, "user1", "book1")
		port.BatchCreateAnnotations(ctx, []*readerModel.Annotation{})
		port.BatchDeleteAnnotations(ctx, []string{})
	})

	t.Run("ChapterContentPort has 12 methods", func(t *testing.T) {
		var port ChapterContentPort = &MockChapterContentPort{}
		ctx := context.Background()

		port.Initialize(ctx)
		port.Health(ctx)
		port.Close(ctx)
		_ = port.GetServiceName()
		_ = port.GetVersion()
		port.GetChapterContent(ctx, "user1", "chapter1")
		port.GetChapterByID(ctx, "chapter1")
		port.GetBookChapters(ctx, "book1", 1, 10)
		port.GetChapterContentWithProgress(ctx, &GetChapterContentRequest{})
		port.GetChapterByNumber(ctx, "user1", "book1", 1)
		port.GetNextChapter(ctx, "user1", "book1", "chapter1")
		port.GetPreviousChapter(ctx, "user1", "book1", "chapter1")
		port.GetChapterList(ctx, "user1", "book1", 1, 10)
		port.GetChapterInfo(ctx, "user1", "chapter1")
	})

	t.Run("ReaderSettingsPort has 6 methods", func(t *testing.T) {
		var port ReaderSettingsPort = &MockReaderSettingsPort{}
		ctx := context.Background()

		port.Initialize(ctx)
		port.Health(ctx)
		port.Close(ctx)
		_ = port.GetServiceName()
		_ = port.GetVersion()
		port.GetReadingSettings(ctx, "user1")
		port.SaveReadingSettings(ctx, &readerModel.ReadingSettings{})
		port.UpdateReadingSettings(ctx, "user1", map[string]interface{}{})
	})

	t.Run("ReaderSyncPort has 6 methods", func(t *testing.T) {
		var port ReaderSyncPort = &MockReaderSyncPort{}
		ctx := context.Background()

		port.Initialize(ctx)
		port.Health(ctx)
		port.Close(ctx)
		_ = port.GetServiceName()
		_ = port.GetVersion()
		port.SyncAnnotations(ctx, &SyncAnnotationsRequest{})
	})
}

// TestPortRequestResponseStructures 测试请求和响应结构体
func TestPortRequestResponseStructures(t *testing.T) {
	t.Run("SaveReadingProgressRequest", func(t *testing.T) {
		req := &SaveReadingProgressRequest{
			UserID:    "user1",
			BookID:    "book1",
			ChapterID: "chapter1",
			Progress:  0.5,
		}
		if req.UserID != "user1" {
			t.Errorf("Expected UserID to be user1, got %s", req.UserID)
		}
	})

	t.Run("GetReadingHistoryRequest", func(t *testing.T) {
		req := &GetReadingHistoryRequest{
			UserID: "user1",
			Page:   1,
			Size:   10,
		}
		if req.Page != 1 {
			t.Errorf("Expected Page to be 1, got %d", req.Page)
		}
	})

	t.Run("GetReadingHistoryResponse", func(t *testing.T) {
		resp := &GetReadingHistoryResponse{
			Progresses: []*readerModel.ReadingProgress{},
			Total:      100,
			Page:       1,
			Size:       10,
			TotalPages: 10,
		}
		if resp.Total != 100 {
			t.Errorf("Expected Total to be 100, got %d", resp.Total)
		}
	})

	t.Run("SyncAnnotationsRequest and Response", func(t *testing.T) {
		req := &SyncAnnotationsRequest{
			UserID:           "user1",
			BookID:           "book1",
			LastSyncTime:     time.Now().Unix(),
			LocalAnnotations: []*readerModel.Annotation{},
		}

		resp := &SyncAnnotationsResponse{
			NewAnnotations:  []*readerModel.Annotation{},
			SyncTime:        time.Now().Unix(),
			UploadedCount:   5,
			DownloadedCount: 10,
		}

		if req.UserID != "user1" {
			t.Errorf("Expected UserID to be user1, got %s", req.UserID)
		}
		if resp.UploadedCount != 5 {
			t.Errorf("Expected UploadedCount to be 5, got %d", resp.UploadedCount)
		}
	})
}
