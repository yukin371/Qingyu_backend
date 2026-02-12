package reader

import (
	"context"
	"errors"
	"testing"
	"time"

	readerModel "Qingyu_backend/models/reader"
	readeriface "Qingyu_backend/service/interfaces/reader"
)

// ============================================================================
// Mock 实现
// ============================================================================

// MockReadingProgressPort Mock 实现
type MockReadingProgressPort struct {
	serviceName string
	version     string
	initCalled  bool
	closeCalled bool
}

func (m *MockReadingProgressPort) Initialize(ctx context.Context) error {
	m.initCalled = true
	return nil
}

func (m *MockReadingProgressPort) Health(ctx context.Context) error { return nil }
func (m *MockReadingProgressPort) Close(ctx context.Context) error {
	m.closeCalled = true
	return nil
}
func (m *MockReadingProgressPort) GetServiceName() string { return m.serviceName }
func (m *MockReadingProgressPort) GetVersion() string     { return m.version }
func (m *MockReadingProgressPort) GetReadingProgress(ctx context.Context, userID, bookID string) (*readerModel.ReadingProgress, error) {
	return &readerModel.ReadingProgress{}, nil
}
func (m *MockReadingProgressPort) SaveReadingProgress(ctx context.Context, req *SaveReadingProgressRequest) error {
	return nil
}
func (m *MockReadingProgressPort) UpdateReadingTime(ctx context.Context, req *UpdateReadingTimeRequest) error {
	return nil
}
func (m *MockReadingProgressPort) GetRecentReading(ctx context.Context, userID string, limit int) ([]*readerModel.ReadingProgress, error) {
	return []*readerModel.ReadingProgress{}, nil
}
func (m *MockReadingProgressPort) GetReadingHistory(ctx context.Context, req *GetReadingHistoryRequest) (*GetReadingHistoryResponse, error) {
	return &GetReadingHistoryResponse{
		Progresses: []*readerModel.ReadingProgress{},
		Total:      0,
		Page:       req.Page,
		Size:       req.Size,
		TotalPages: 0,
	}, nil
}
func (m *MockReadingProgressPort) GetTotalReadingTime(ctx context.Context, userID string) (int64, error) { return 0, nil }
func (m *MockReadingProgressPort) GetReadingTimeByPeriod(ctx context.Context, req *GetReadingTimeByPeriodRequest) (int64, error) {
	return 0, nil
}
func (m *MockReadingProgressPort) GetUnfinishedBooks(ctx context.Context, userID string) ([]*readerModel.ReadingProgress, error) {
	return []*readerModel.ReadingProgress{}, nil
}
func (m *MockReadingProgressPort) GetFinishedBooks(ctx context.Context, userID string) ([]*readerModel.ReadingProgress, error) {
	return []*readerModel.ReadingProgress{}, nil
}
func (m *MockReadingProgressPort) DeleteReadingProgress(ctx context.Context, userID, bookID string) error { return nil }
func (m *MockReadingProgressPort) UpdateBookStatus(ctx context.Context, req *UpdateBookStatusRequest) error     { return nil }
func (m *MockReadingProgressPort) BatchUpdateBookStatus(ctx context.Context, req *BatchUpdateBookStatusRequest) error {
	return nil
}

// MockAnnotationPort Mock 实现
type MockAnnotationPort struct {
	serviceName string
	version     string
	initCalled  bool
	closeCalled bool
}

func (m *MockAnnotationPort) Initialize(ctx context.Context) error {
	m.initCalled = true
	return nil
}

func (m *MockAnnotationPort) Health(ctx context.Context) error { return nil }
func (m *MockAnnotationPort) Close(ctx context.Context) error {
	m.closeCalled = true
	return nil
}
func (m *MockAnnotationPort) GetServiceName() string { return m.serviceName }
func (m *MockAnnotationPort) GetVersion() string     { return m.version }
func (m *MockAnnotationPort) CreateAnnotation(ctx context.Context, annotation *readerModel.Annotation) error {
	return nil
}
func (m *MockAnnotationPort) UpdateAnnotation(ctx context.Context, annotationID string, updates map[string]interface{}) error {
	return nil
}
func (m *MockAnnotationPort) DeleteAnnotation(ctx context.Context, annotationID string) error { return nil }
func (m *MockAnnotationPort) GetAnnotationsByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*readerModel.Annotation, error) {
	return []*readerModel.Annotation{}, nil
}
func (m *MockAnnotationPort) GetAnnotationsByBook(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error) {
	return []*readerModel.Annotation{}, nil
}
func (m *MockAnnotationPort) GetNotes(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error) {
	return []*readerModel.Annotation{}, nil
}
func (m *MockAnnotationPort) SearchNotes(ctx context.Context, userID, keyword string) ([]*readerModel.Annotation, error) {
	return []*readerModel.Annotation{}, nil
}
func (m *MockAnnotationPort) GetBookmarks(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error) {
	return []*readerModel.Annotation{}, nil
}
func (m *MockAnnotationPort) GetLatestBookmark(ctx context.Context, userID, bookID string) (*readerModel.Annotation, error) {
	return &readerModel.Annotation{}, nil
}
func (m *MockAnnotationPort) GetHighlights(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error) {
	return []*readerModel.Annotation{}, nil
}
func (m *MockAnnotationPort) GetRecentAnnotations(ctx context.Context, userID string, limit int) ([]*readerModel.Annotation, error) {
	return []*readerModel.Annotation{}, nil
}
func (m *MockAnnotationPort) GetPublicAnnotations(ctx context.Context, bookID, chapterID string) ([]*readerModel.Annotation, error) {
	return []*readerModel.Annotation{}, nil
}
func (m *MockAnnotationPort) GetAnnotationStats(ctx context.Context, userID, bookID string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
func (m *MockAnnotationPort) BatchCreateAnnotations(ctx context.Context, annotations []*readerModel.Annotation) error {
	return nil
}
func (m *MockAnnotationPort) BatchDeleteAnnotations(ctx context.Context, annotationIDs []string) error { return nil }

// MockChapterContentPort Mock 实现
type MockChapterContentPort struct {
	serviceName string
	version     string
	initCalled  bool
	closeCalled bool
}

func (m *MockChapterContentPort) Initialize(ctx context.Context) error {
	m.initCalled = true
	return nil
}

func (m *MockChapterContentPort) Health(ctx context.Context) error { return nil }
func (m *MockChapterContentPort) Close(ctx context.Context) error {
	m.closeCalled = true
	return nil
}
func (m *MockChapterContentPort) GetServiceName() string { return m.serviceName }
func (m *MockChapterContentPort) GetVersion() string     { return m.version }
func (m *MockChapterContentPort) GetChapterContent(ctx context.Context, userID, chapterID string) (string, error) {
	return "内容", nil
}
func (m *MockChapterContentPort) GetChapterByID(ctx context.Context, chapterID string) (interface{}, error) {
	return map[string]interface{}{}, nil
}
func (m *MockChapterContentPort) GetBookChapters(ctx context.Context, bookID string, page, size int) (interface{}, int64, error) {
	return []interface{}{}, 0, nil
}
func (m *MockChapterContentPort) GetChapterContentWithProgress(ctx context.Context, req *GetChapterContentRequest) (*readeriface.ChapterContentResponse, error) {
	return &readeriface.ChapterContentResponse{}, nil
}
func (m *MockChapterContentPort) GetChapterByNumber(ctx context.Context, userID, bookID string, chapterNum int) (*readeriface.ChapterInfo, error) {
	return &readeriface.ChapterInfo{}, nil
}
func (m *MockChapterContentPort) GetNextChapter(ctx context.Context, userID, bookID, chapterID string) (*readeriface.ChapterInfo, error) {
	return &readeriface.ChapterInfo{}, nil
}
func (m *MockChapterContentPort) GetPreviousChapter(ctx context.Context, userID, bookID, chapterID string) (*readeriface.ChapterInfo, error) {
	return &readeriface.ChapterInfo{}, nil
}
func (m *MockChapterContentPort) GetChapterList(ctx context.Context, userID, bookID string, page, size int) (*readeriface.ChapterListResponse, error) {
	return &readeriface.ChapterListResponse{}, nil
}
func (m *MockChapterContentPort) GetChapterInfo(ctx context.Context, userID, chapterID string) (*readeriface.ChapterInfo, error) {
	return &readeriface.ChapterInfo{}, nil
}

// MockReaderSettingsPort Mock 实现
type MockReaderSettingsPort struct {
	serviceName string
	version     string
	initCalled  bool
	closeCalled bool
}

func (m *MockReaderSettingsPort) Initialize(ctx context.Context) error {
	m.initCalled = true
	return nil
}

func (m *MockReaderSettingsPort) Health(ctx context.Context) error { return nil }
func (m *MockReaderSettingsPort) Close(ctx context.Context) error {
	m.closeCalled = true
	return nil
}
func (m *MockReaderSettingsPort) GetServiceName() string { return m.serviceName }
func (m *MockReaderSettingsPort) GetVersion() string     { return m.version }
func (m *MockReaderSettingsPort) GetReadingSettings(ctx context.Context, userID string) (*readerModel.ReadingSettings, error) {
	return &readerModel.ReadingSettings{}, nil
}
func (m *MockReaderSettingsPort) SaveReadingSettings(ctx context.Context, settings *readerModel.ReadingSettings) error {
	return nil
}
func (m *MockReaderSettingsPort) UpdateReadingSettings(ctx context.Context, userID string, updates map[string]interface{}) error {
	return nil
}

// MockReaderSyncPort Mock 实现
type MockReaderSyncPort struct {
	serviceName string
	version     string
	initCalled  bool
	closeCalled bool
}

func (m *MockReaderSyncPort) Initialize(ctx context.Context) error {
	m.initCalled = true
	return nil
}

func (m *MockReaderSyncPort) Health(ctx context.Context) error { return nil }
func (m *MockReaderSyncPort) Close(ctx context.Context) error {
	m.closeCalled = true
	return nil
}
func (m *MockReaderSyncPort) GetServiceName() string { return m.serviceName }
func (m *MockReaderSyncPort) GetVersion() string     { return m.version }
func (m *MockReaderSyncPort) SyncAnnotations(ctx context.Context, req *SyncAnnotationsRequest) (*SyncAnnotationsResponse, error) {
	return &SyncAnnotationsResponse{
		NewAnnotations:  []*readerModel.Annotation{},
		SyncTime:        time.Now().Unix(),
		UploadedCount:   0,
		DownloadedCount: 0,
	}, nil
}

// ============================================================================
// ReaderServiceAdapter 测试
// ============================================================================

// TestNewReaderServiceAdapter 测试创建 ReaderServiceAdapter
func TestNewReaderServiceAdapter(t *testing.T) {
	progressPort := &MockReadingProgressPort{serviceName: "ProgressPort", version: "1.0.0"}
	annotationPort := &MockAnnotationPort{serviceName: "AnnotationPort", version: "1.0.0"}
	chapterPort := &MockChapterContentPort{serviceName: "ChapterPort", version: "1.0.0"}
	settingsPort := &MockReaderSettingsPort{serviceName: "SettingsPort", version: "1.0.0"}
	syncPort := &MockReaderSyncPort{serviceName: "SyncPort", version: "1.0.0"}

	adapter := NewReaderServiceAdapter(
		progressPort,
		annotationPort,
		chapterPort,
		settingsPort,
		syncPort,
	)

	if adapter == nil {
		t.Fatal("NewReaderServiceAdapter() returned nil")
	}

	if adapter.progressPort != progressPort {
		t.Error("progressPort not set correctly")
	}
	if adapter.annotationPort != annotationPort {
		t.Error("annotationPort not set correctly")
	}
	if adapter.chapterPort != chapterPort {
		t.Error("chapterPort not set correctly")
	}
	if adapter.settingsPort != settingsPort {
		t.Error("settingsPort not set correctly")
	}
	if adapter.syncPort != syncPort {
		t.Error("syncPort not set correctly")
	}
}

// TestReaderServiceAdapter_BaseService 测试 BaseService 接口实现
func TestReaderServiceAdapter_BaseService(t *testing.T) {
	progressPort := &MockReadingProgressPort{serviceName: "ProgressPort", version: "1.0.0"}
	annotationPort := &MockAnnotationPort{serviceName: "AnnotationPort", version: "1.0.0"}
	chapterPort := &MockChapterContentPort{serviceName: "ChapterPort", version: "1.0.0"}
	settingsPort := &MockReaderSettingsPort{serviceName: "SettingsPort", version: "1.0.0"}
	syncPort := &MockReaderSyncPort{serviceName: "SyncPort", version: "1.0.0"}

	adapter := NewReaderServiceAdapter(
		progressPort,
		annotationPort,
		chapterPort,
		settingsPort,
		syncPort,
	)

	ctx := context.Background()

	t.Run("Initialize delegates to progressPort", func(t *testing.T) {
		if err := adapter.Initialize(ctx); err != nil {
			t.Errorf("Initialize() error = %v", err)
		}
		if !progressPort.initCalled {
			t.Error("progressPort.Initialize was not called")
		}
	})

	t.Run("Health delegates to progressPort", func(t *testing.T) {
		if err := adapter.Health(ctx); err != nil {
			t.Errorf("Health() error = %v", err)
		}
	})

	t.Run("Close closes all ports", func(t *testing.T) {
		if err := adapter.Close(ctx); err != nil {
			t.Errorf("Close() error = %v", err)
		}
		if !progressPort.closeCalled {
			t.Error("progressPort.Close was not called")
		}
		if !annotationPort.closeCalled {
			t.Error("annotationPort.Close was not called")
		}
		if !chapterPort.closeCalled {
			t.Error("chapterPort.Close was not called")
		}
		if !settingsPort.closeCalled {
			t.Error("settingsPort.Close was not called")
		}
		if !syncPort.closeCalled {
			t.Error("syncPort.Close was not called")
		}
	})

	t.Run("GetServiceName delegates to progressPort", func(t *testing.T) {
		name := adapter.GetServiceName()
		if name != "ProgressPort" {
			t.Errorf("GetServiceName() = %v, want %v", name, "ProgressPort")
		}
	})

	t.Run("GetVersion delegates to progressPort", func(t *testing.T) {
		version := adapter.GetVersion()
		if version != "1.0.0" {
			t.Errorf("GetVersion() = %v, want %v", version, "1.0.0")
		}
	})
}

// TestReaderServiceAdapter_ReadingProgress 测试阅读进度方法委托
func TestReaderServiceAdapter_ReadingProgress(t *testing.T) {
	progressPort := &MockReadingProgressPort{serviceName: "ProgressPort", version: "1.0.0"}
	annotationPort := &MockAnnotationPort{serviceName: "AnnotationPort", version: "1.0.0"}
	chapterPort := &MockChapterContentPort{serviceName: "ChapterPort", version: "1.0.0"}
	settingsPort := &MockReaderSettingsPort{serviceName: "SettingsPort", version: "1.0.0"}
	syncPort := &MockReaderSyncPort{serviceName: "SyncPort", version: "1.0.0"}

	adapter := NewReaderServiceAdapter(
		progressPort,
		annotationPort,
		chapterPort,
		settingsPort,
		syncPort,
	)

	ctx := context.Background()

	t.Run("GetReadingProgress", func(t *testing.T) {
		progress, err := adapter.GetReadingProgress(ctx, "user1", "book1")
		if err != nil {
			t.Errorf("GetReadingProgress() error = %v", err)
		}
		if progress == nil {
			t.Error("GetReadingProgress() returned nil progress")
		}
	})

	t.Run("SaveReadingProgress", func(t *testing.T) {
		err := adapter.SaveReadingProgress(ctx, "user1", "book1", "chapter1", 0.5)
		if err != nil {
			t.Errorf("SaveReadingProgress() error = %v", err)
		}
	})

	t.Run("UpdateReadingTime", func(t *testing.T) {
		err := adapter.UpdateReadingTime(ctx, "user1", "book1", 100)
		if err != nil {
			t.Errorf("UpdateReadingTime() error = %v", err)
		}
	})

	t.Run("GetRecentReading", func(t *testing.T) {
		progresses, err := adapter.GetRecentReading(ctx, "user1", 10)
		if err != nil {
			t.Errorf("GetRecentReading() error = %v", err)
		}
		if progresses == nil {
			t.Error("GetRecentReading() returned nil progresses")
		}
	})

	t.Run("GetReadingHistory", func(t *testing.T) {
		progresses, total, err := adapter.GetReadingHistory(ctx, "user1", 1, 10)
		if err != nil {
			t.Errorf("GetReadingHistory() error = %v", err)
		}
		if progresses == nil {
			t.Error("GetReadingHistory() returned nil progresses")
		}
		if total != 0 {
			t.Errorf("GetReadingHistory() total = %v, want 0", total)
		}
	})

	t.Run("GetTotalReadingTime", func(t *testing.T) {
		total, err := adapter.GetTotalReadingTime(ctx, "user1")
		if err != nil {
			t.Errorf("GetTotalReadingTime() error = %v", err)
		}
		if total != 0 {
			t.Errorf("GetTotalReadingTime() = %v, want 0", total)
		}
	})

	t.Run("GetReadingTimeByPeriod", func(t *testing.T) {
		startTime := time.Now().Add(-24 * time.Hour)
		endTime := time.Now()
		total, err := adapter.GetReadingTimeByPeriod(ctx, "user1", startTime, endTime)
		if err != nil {
			t.Errorf("GetReadingTimeByPeriod() error = %v", err)
		}
		if total != 0 {
			t.Errorf("GetReadingTimeByPeriod() = %v, want 0", total)
		}
	})

	t.Run("GetUnfinishedBooks", func(t *testing.T) {
		books, err := adapter.GetUnfinishedBooks(ctx, "user1")
		if err != nil {
			t.Errorf("GetUnfinishedBooks() error = %v", err)
		}
		if books == nil {
			t.Error("GetUnfinishedBooks() returned nil books")
		}
	})

	t.Run("GetFinishedBooks", func(t *testing.T) {
		books, err := adapter.GetFinishedBooks(ctx, "user1")
		if err != nil {
			t.Errorf("GetFinishedBooks() error = %v", err)
		}
		if books == nil {
			t.Error("GetFinishedBooks() returned nil books")
		}
	})

	t.Run("DeleteReadingProgress", func(t *testing.T) {
		err := adapter.DeleteReadingProgress(ctx, "user1", "book1")
		if err != nil {
			t.Errorf("DeleteReadingProgress() error = %v", err)
		}
	})

	t.Run("UpdateBookStatus", func(t *testing.T) {
		err := adapter.UpdateBookStatus(ctx, "user1", "book1", "reading")
		if err != nil {
			t.Errorf("UpdateBookStatus() error = %v", err)
		}
	})

	t.Run("BatchUpdateBookStatus", func(t *testing.T) {
		err := adapter.BatchUpdateBookStatus(ctx, "user1", []string{"book1", "book2"}, "reading")
		if err != nil {
			t.Errorf("BatchUpdateBookStatus() error = %v", err)
		}
	})
}

// TestReaderServiceAdapter_Annotation 测试标注方法委托
func TestReaderServiceAdapter_Annotation(t *testing.T) {
	progressPort := &MockReadingProgressPort{serviceName: "ProgressPort", version: "1.0.0"}
	annotationPort := &MockAnnotationPort{serviceName: "AnnotationPort", version: "1.0.0"}
	chapterPort := &MockChapterContentPort{serviceName: "ChapterPort", version: "1.0.0"}
	settingsPort := &MockReaderSettingsPort{serviceName: "SettingsPort", version: "1.0.0"}
	syncPort := &MockReaderSyncPort{serviceName: "SyncPort", version: "1.0.0"}

	adapter := NewReaderServiceAdapter(
		progressPort,
		annotationPort,
		chapterPort,
		settingsPort,
		syncPort,
	)

	ctx := context.Background()

	t.Run("CreateAnnotation", func(t *testing.T) {
		err := adapter.CreateAnnotation(ctx, &readerModel.Annotation{})
		if err != nil {
			t.Errorf("CreateAnnotation() error = %v", err)
		}
	})

	t.Run("UpdateAnnotation", func(t *testing.T) {
		err := adapter.UpdateAnnotation(ctx, "annotation1", map[string]interface{}{"content": "new content"})
		if err != nil {
			t.Errorf("UpdateAnnotation() error = %v", err)
		}
	})

	t.Run("DeleteAnnotation", func(t *testing.T) {
		err := adapter.DeleteAnnotation(ctx, "annotation1")
		if err != nil {
			t.Errorf("DeleteAnnotation() error = %v", err)
		}
	})

	t.Run("GetAnnotationsByChapter", func(t *testing.T) {
		annotations, err := adapter.GetAnnotationsByChapter(ctx, "user1", "book1", "chapter1")
		if err != nil {
			t.Errorf("GetAnnotationsByChapter() error = %v", err)
		}
		if annotations == nil {
			t.Error("GetAnnotationsByChapter() returned nil annotations")
		}
	})

	t.Run("GetAnnotationsByBook", func(t *testing.T) {
		annotations, err := adapter.GetAnnotationsByBook(ctx, "user1", "book1")
		if err != nil {
			t.Errorf("GetAnnotationsByBook() error = %v", err)
		}
		if annotations == nil {
			t.Error("GetAnnotationsByBook() returned nil annotations")
		}
	})

	t.Run("GetNotes", func(t *testing.T) {
		notes, err := adapter.GetNotes(ctx, "user1", "book1")
		if err != nil {
			t.Errorf("GetNotes() error = %v", err)
		}
		if notes == nil {
			t.Error("GetNotes() returned nil notes")
		}
	})

	t.Run("SearchNotes", func(t *testing.T) {
		notes, err := adapter.SearchNotes(ctx, "user1", "keyword")
		if err != nil {
			t.Errorf("SearchNotes() error = %v", err)
		}
		if notes == nil {
			t.Error("SearchNotes() returned nil notes")
		}
	})

	t.Run("GetBookmarks", func(t *testing.T) {
		bookmarks, err := adapter.GetBookmarks(ctx, "user1", "book1")
		if err != nil {
			t.Errorf("GetBookmarks() error = %v", err)
		}
		if bookmarks == nil {
			t.Error("GetBookmarks() returned nil bookmarks")
		}
	})

	t.Run("GetLatestBookmark", func(t *testing.T) {
		bookmark, err := adapter.GetLatestBookmark(ctx, "user1", "book1")
		if err != nil {
			t.Errorf("GetLatestBookmark() error = %v", err)
		}
		if bookmark == nil {
			t.Error("GetLatestBookmark() returned nil bookmark")
		}
	})

	t.Run("GetHighlights", func(t *testing.T) {
		highlights, err := adapter.GetHighlights(ctx, "user1", "book1")
		if err != nil {
			t.Errorf("GetHighlights() error = %v", err)
		}
		if highlights == nil {
			t.Error("GetHighlights() returned nil highlights")
		}
	})

	t.Run("GetRecentAnnotations", func(t *testing.T) {
		annotations, err := adapter.GetRecentAnnotations(ctx, "user1", 10)
		if err != nil {
			t.Errorf("GetRecentAnnotations() error = %v", err)
		}
		if annotations == nil {
			t.Error("GetRecentAnnotations() returned nil annotations")
		}
	})

	t.Run("GetPublicAnnotations", func(t *testing.T) {
		annotations, err := adapter.GetPublicAnnotations(ctx, "book1", "chapter1")
		if err != nil {
			t.Errorf("GetPublicAnnotations() error = %v", err)
		}
		if annotations == nil {
			t.Error("GetPublicAnnotations() returned nil annotations")
		}
	})

	t.Run("GetAnnotationStats", func(t *testing.T) {
		stats, err := adapter.GetAnnotationStats(ctx, "user1", "book1")
		if err != nil {
			t.Errorf("GetAnnotationStats() error = %v", err)
		}
		if stats == nil {
			t.Error("GetAnnotationStats() returned nil stats")
		}
	})

	t.Run("BatchCreateAnnotations", func(t *testing.T) {
		err := adapter.BatchCreateAnnotations(ctx, []*readerModel.Annotation{})
		if err != nil {
			t.Errorf("BatchCreateAnnotations() error = %v", err)
		}
	})

	t.Run("BatchDeleteAnnotations", func(t *testing.T) {
		err := adapter.BatchDeleteAnnotations(ctx, []string{"annotation1", "annotation2"})
		if err != nil {
			t.Errorf("BatchDeleteAnnotations() error = %v", err)
		}
	})
}

// TestReaderServiceAdapter_Settings 测试阅读设置方法委托
func TestReaderServiceAdapter_Settings(t *testing.T) {
	progressPort := &MockReadingProgressPort{serviceName: "ProgressPort", version: "1.0.0"}
	annotationPort := &MockAnnotationPort{serviceName: "AnnotationPort", version: "1.0.0"}
	chapterPort := &MockChapterContentPort{serviceName: "ChapterPort", version: "1.0.0"}
	settingsPort := &MockReaderSettingsPort{serviceName: "SettingsPort", version: "1.0.0"}
	syncPort := &MockReaderSyncPort{serviceName: "SyncPort", version: "1.0.0"}

	adapter := NewReaderServiceAdapter(
		progressPort,
		annotationPort,
		chapterPort,
		settingsPort,
		syncPort,
	)

	ctx := context.Background()

	t.Run("GetReadingSettings", func(t *testing.T) {
		settings, err := adapter.GetReadingSettings(ctx, "user1")
		if err != nil {
			t.Errorf("GetReadingSettings() error = %v", err)
		}
		if settings == nil {
			t.Error("GetReadingSettings() returned nil settings")
		}
	})

	t.Run("SaveReadingSettings", func(t *testing.T) {
		err := adapter.SaveReadingSettings(ctx, &readerModel.ReadingSettings{})
		if err != nil {
			t.Errorf("SaveReadingSettings() error = %v", err)
		}
	})

	t.Run("UpdateReadingSettings", func(t *testing.T) {
		err := adapter.UpdateReadingSettings(ctx, "user1", map[string]interface{}{"font_size": 18})
		if err != nil {
			t.Errorf("UpdateReadingSettings() error = %v", err)
		}
	})
}

// TestReaderServiceAdapter_Sync 测试同步方法委托
func TestReaderServiceAdapter_Sync(t *testing.T) {
	progressPort := &MockReadingProgressPort{serviceName: "ProgressPort", version: "1.0.0"}
	annotationPort := &MockAnnotationPort{serviceName: "AnnotationPort", version: "1.0.0"}
	chapterPort := &MockChapterContentPort{serviceName: "ChapterPort", version: "1.0.0"}
	settingsPort := &MockReaderSettingsPort{serviceName: "SettingsPort", version: "1.0.0"}
	syncPort := &MockReaderSyncPort{serviceName: "SyncPort", version: "1.0.0"}

	adapter := NewReaderServiceAdapter(
		progressPort,
		annotationPort,
		chapterPort,
		settingsPort,
		syncPort,
	)

	ctx := context.Background()

	t.Run("SyncAnnotations", func(t *testing.T) {
		// 创建旧的 sync 请求结构
		type legacySyncReq struct {
			BookID           string
			LastSyncTime     int64
			LocalAnnotations []*readerModel.Annotation
		}
		req := legacySyncReq{
			BookID:           "book1",
			LastSyncTime:     time.Now().Unix(),
			LocalAnnotations: []*readerModel.Annotation{},
		}

		result, err := adapter.SyncAnnotations(ctx, "user1", req)
		if err != nil {
			t.Errorf("SyncAnnotations() error = %v", err)
		}
		if result == nil {
			t.Error("SyncAnnotations() returned nil result")
		}
	})
}

// ============================================================================
// ReaderChapterServiceAdapter 测试
// ============================================================================

// TestNewReaderChapterServiceAdapter 测试创建 ReaderChapterServiceAdapter
func TestNewReaderChapterServiceAdapter(t *testing.T) {
	chapterPort := &MockChapterContentPort{serviceName: "ChapterPort", version: "1.0.0"}

	adapter := NewReaderChapterServiceAdapter(chapterPort)

	if adapter == nil {
		t.Fatal("NewReaderChapterServiceAdapter() returned nil")
	}

	if adapter.chapterPort != chapterPort {
		t.Error("chapterPort not set correctly")
	}
}

// TestReaderChapterServiceAdapter_AllMethods 测试所有方法委托
func TestReaderChapterServiceAdapter_AllMethods(t *testing.T) {
	chapterPort := &MockChapterContentPort{serviceName: "ChapterPort", version: "1.0.0"}

	adapter := NewReaderChapterServiceAdapter(chapterPort)

	ctx := context.Background()

	t.Run("GetChapterContent", func(t *testing.T) {
		result, err := adapter.GetChapterContent(ctx, "user1", "book1", "chapter1")
		if err != nil {
			t.Errorf("GetChapterContent() error = %v", err)
		}
		if result == nil {
			t.Error("GetChapterContent() returned nil result")
		}
	})

	t.Run("GetChapterByNumber", func(t *testing.T) {
		result, err := adapter.GetChapterByNumber(ctx, "user1", "book1", 1)
		if err != nil {
			t.Errorf("GetChapterByNumber() error = %v", err)
		}
		if result == nil {
			t.Error("GetChapterByNumber() returned nil result")
		}
	})

	t.Run("GetNextChapter", func(t *testing.T) {
		result, err := adapter.GetNextChapter(ctx, "user1", "book1", "chapter1")
		if err != nil {
			t.Errorf("GetNextChapter() error = %v", err)
		}
		if result == nil {
			t.Error("GetNextChapter() returned nil result")
		}
	})

	t.Run("GetPreviousChapter", func(t *testing.T) {
		result, err := adapter.GetPreviousChapter(ctx, "user1", "book1", "chapter1")
		if err != nil {
			t.Errorf("GetPreviousChapter() error = %v", err)
		}
		if result == nil {
			t.Error("GetPreviousChapter() returned nil result")
		}
	})

	t.Run("GetChapterList", func(t *testing.T) {
		result, err := adapter.GetChapterList(ctx, "user1", "book1", 1, 10)
		if err != nil {
			t.Errorf("GetChapterList() error = %v", err)
		}
		if result == nil {
			t.Error("GetChapterList() returned nil result")
		}
	})

	t.Run("GetChapterInfo", func(t *testing.T) {
		result, err := adapter.GetChapterInfo(ctx, "user1", "chapter1")
		if err != nil {
			t.Errorf("GetChapterInfo() error = %v", err)
		}
		if result == nil {
			t.Error("GetChapterInfo() returned nil result")
		}
	})
}

// TestAdapter_ErrorHandling 测试适配器错误处理
func TestAdapter_ErrorHandling(t *testing.T) {
	// 创建一个会返回错误的 Mock
	errorPort := &MockReadingProgressPort{
		serviceName: "ErrorPort",
		version:     "1.0.0",
	}

	// 覆盖方法使其返回错误
	originalGet := errorPort.GetReadingProgress
	errorPort.GetReadingProgress = func(ctx context.Context, userID, bookID string) (*readerModel.ReadingProgress, error) {
		return nil, errors.New("mock error")
	}

	progressPort := errorPort
	annotationPort := &MockAnnotationPort{serviceName: "AnnotationPort", version: "1.0.0"}
	chapterPort := &MockChapterContentPort{serviceName: "ChapterPort", version: "1.0.0"}
	settingsPort := &MockReaderSettingsPort{serviceName: "SettingsPort", version: "1.0.0"}
	syncPort := &MockReaderSyncPort{serviceName: "SyncPort", version: "1.0.0"}

	adapter := NewReaderServiceAdapter(
		progressPort,
		annotationPort,
		chapterPort,
		settingsPort,
		syncPort,
	)

	ctx := context.Background()

	t.Run("GetReadingProgress propagates error", func(t *testing.T) {
		_, err := adapter.GetReadingProgress(ctx, "user1", "book1")
		if err == nil {
			t.Error("Expected error but got nil")
		}
		if err.Error() != "mock error" {
			t.Errorf("Expected 'mock error' but got '%v'", err)
		}
	})

	// 恢复原始方法
	errorPort.GetReadingProgress = originalGet
}
