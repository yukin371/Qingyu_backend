package reader

import (
	"context"
	"testing"

	readeriface "Qingyu_backend/service/interfaces/reader"
	"Qingyu_backend/models/reader"
	readermigration "Qingyu_backend/service/reader/_migration"
)

// 导入 ReaderServiceAdapter 类型别名，避免长包名
type ReaderServiceAdapter = readermigration.ReaderServiceAdapter
type ReaderChapterServiceAdapter = readermigration.ReaderChapterServiceAdapter

// TestNewReaderServiceFactory 测试创建工厂
func TestNewReaderServiceFactory(t *testing.T) {
	factory := NewReaderServiceFactory()
	if factory == nil {
		t.Fatal("NewReaderServiceFactory() returned nil")
	}
}

// TestReaderServiceFactory_CreateWithPorts 测试使用 Port 接口创建服务
func TestReaderServiceFactory_CreateWithPorts(t *testing.T) {
	factory := NewReaderServiceFactory()

	// 创建 Mock Port 实现
	progressPort := &MockFactoryReadingProgressPort{}
	annotationPort := &MockFactoryAnnotationPort{}
	chapterPort := &MockFactoryChapterContentPort{}
	settingsPort := &MockFactoryReaderSettingsPort{}
	syncPort := &MockFactoryReaderSyncPort{}

	// 创建服务
	service := factory.CreateWithPorts(
		progressPort,
		annotationPort,
		chapterPort,
		settingsPort,
		syncPort,
	)

	if service == nil {
		t.Fatal("CreateWithPorts() returned nil")
	}

	// 验证返回的是 ReaderServiceAdapter 类型
	if _, ok := service.(*ReaderServiceAdapter); !ok {
		t.Errorf("CreateWithPorts() returned %T, want *ReaderServiceAdapter", service)
	}
}

// TestReaderServiceFactory_CreateChapterServiceWithPorts 测试创建章节服务
func TestReaderServiceFactory_CreateChapterServiceWithPorts(t *testing.T) {
	factory := NewReaderServiceFactory()

	// 创建 Mock Port 实现
	chapterPort := &MockFactoryChapterContentPort{}

	// 创建服务
	service := factory.CreateChapterServiceWithPorts(chapterPort)

	if service == nil {
		t.Fatal("CreateChapterServiceWithPorts() returned nil")
	}

	// 验证返回的是 ReaderChapterServiceAdapter 类型
	if _, ok := service.(*ReaderChapterServiceAdapter); !ok {
		t.Errorf("CreateChapterServiceWithPorts() returned %T, want *ReaderChapterServiceAdapter", service)
	}
}

// TestReaderServiceFactory_CreateFromImplementations 测试从结构体创建服务
func TestReaderServiceFactory_CreateFromImplementations(t *testing.T) {
	factory := NewReaderServiceFactory()

	// 创建 Port 实现结构体
	ports := PortImplementations{
		ProgressPort:   &MockFactoryReadingProgressPort{},
		AnnotationPort: &MockFactoryAnnotationPort{},
		ChapterPort:    &MockFactoryChapterContentPort{},
		SettingsPort:   &MockFactoryReaderSettingsPort{},
		SyncPort:       &MockFactoryReaderSyncPort{},
	}

	// 创建服务
	service := factory.CreateFromImplementations(ports)

	if service == nil {
		t.Fatal("CreateFromImplementations() returned nil")
	}

	// 验证返回的是 ReaderServiceAdapter 类型
	if _, ok := service.(*ReaderServiceAdapter); !ok {
		t.Errorf("CreateFromImplementations() returned %T, want *ReaderServiceAdapter", service)
	}
}

// TestReaderServiceFactory_CreateChapterServiceFromImplementations 测试从结构体创建章节服务
func TestReaderServiceFactory_CreateChapterServiceFromImplementations(t *testing.T) {
	factory := NewReaderServiceFactory()

	// 创建章节 Port 实现结构体
	ports := ChapterPortImplementations{
		ChapterPort: &MockFactoryChapterContentPort{},
	}

	// 创建服务
	service := factory.CreateChapterServiceFromImplementations(ports)

	if service == nil {
		t.Fatal("CreateChapterServiceFromImplementations() returned nil")
	}

	// 验证返回的是 ReaderChapterServiceAdapter 类型
	if _, ok := service.(*ReaderChapterServiceAdapter); !ok {
		t.Errorf("CreateChapterServiceFromImplementations() returned %T, want *ReaderChapterServiceAdapter", service)
	}
}

// ============================================================================
// Mock 实现 - 用于工厂测试
// ============================================================================

// MockFactoryReadingProgressPort Mock 实现
type MockFactoryReadingProgressPort struct{}

func (m *MockFactoryReadingProgressPort) Initialize(ctx context.Context) error                           { return nil }
func (m *MockFactoryReadingProgressPort) Health(ctx context.Context) error                                { return nil }
func (m *MockFactoryReadingProgressPort) Close(ctx context.Context) error                                 { return nil }
func (m *MockFactoryReadingProgressPort) GetServiceName() string                                     { return "MockProgressPort" }
func (m *MockFactoryReadingProgressPort) GetVersion() string                                         { return "1.0.0" }
func (m *MockFactoryReadingProgressPort) GetReadingProgress(ctx context.Context, userID, bookID string) (*reader.ReadingProgress, error) {
	return nil, nil
}
func (m *MockFactoryReadingProgressPort) SaveReadingProgress(ctx context.Context, req *readeriface.SaveReadingProgressRequest) error {
	return nil
}
func (m *MockFactoryReadingProgressPort) UpdateReadingTime(ctx context.Context, req *readeriface.UpdateReadingTimeRequest) error {
	return nil
}
func (m *MockFactoryReadingProgressPort) GetRecentReading(ctx context.Context, userID string, limit int) ([]*reader.ReadingProgress, error) {
	return nil, nil
}
func (m *MockFactoryReadingProgressPort) GetReadingHistory(ctx context.Context, req *readeriface.GetReadingHistoryRequest) (*readeriface.GetReadingHistoryResponse, error) {
	return nil, nil
}
func (m *MockFactoryReadingProgressPort) GetTotalReadingTime(ctx context.Context, userID string) (int64, error) { return 0, nil }
func (m *MockFactoryReadingProgressPort) GetReadingTimeByPeriod(ctx context.Context, req *readeriface.GetReadingTimeByPeriodRequest) (int64, error) {
	return 0, nil
}
func (m *MockFactoryReadingProgressPort) GetUnfinishedBooks(ctx context.Context, userID string) ([]*reader.ReadingProgress, error) {
	return nil, nil
}
func (m *MockFactoryReadingProgressPort) GetFinishedBooks(ctx context.Context, userID string) ([]*reader.ReadingProgress, error) {
	return nil, nil
}
func (m *MockFactoryReadingProgressPort) DeleteReadingProgress(ctx context.Context, userID, bookID string) error { return nil }
func (m *MockFactoryReadingProgressPort) UpdateBookStatus(ctx context.Context, req *readeriface.UpdateBookStatusRequest) error {
	return nil
}
func (m *MockFactoryReadingProgressPort) BatchUpdateBookStatus(ctx context.Context, req *readeriface.BatchUpdateBookStatusRequest) error {
	return nil
}

// MockFactoryAnnotationPort Mock 实现
type MockFactoryAnnotationPort struct{}

func (m *MockFactoryAnnotationPort) Initialize(ctx context.Context) error { return nil }
func (m *MockFactoryAnnotationPort) Health(ctx context.Context) error      { return nil }
func (m *MockFactoryAnnotationPort) Close(ctx context.Context) error       { return nil }
func (m *MockFactoryAnnotationPort) GetServiceName() string           { return "MockAnnotationPort" }
func (m *MockFactoryAnnotationPort) GetVersion() string               { return "1.0.0" }
func (m *MockFactoryAnnotationPort) CreateAnnotation(ctx context.Context, annotation *reader.Annotation) error {
	return nil
}
func (m *MockFactoryAnnotationPort) UpdateAnnotation(ctx context.Context, annotationID string, updates map[string]interface{}) error {
	return nil
}
func (m *MockFactoryAnnotationPort) DeleteAnnotation(ctx context.Context, annotationID string) error { return nil }
func (m *MockFactoryAnnotationPort) GetAnnotationsByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*reader.Annotation, error) {
	return nil, nil
}
func (m *MockFactoryAnnotationPort) GetAnnotationsByBook(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	return nil, nil
}
func (m *MockFactoryAnnotationPort) GetNotes(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	return nil, nil
}
func (m *MockFactoryAnnotationPort) SearchNotes(ctx context.Context, userID, keyword string) ([]*reader.Annotation, error) {
	return nil, nil
}
func (m *MockFactoryAnnotationPort) GetBookmarks(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	return nil, nil
}
func (m *MockFactoryAnnotationPort) GetLatestBookmark(ctx context.Context, userID, bookID string) (*reader.Annotation, error) {
	return nil, nil
}
func (m *MockFactoryAnnotationPort) GetHighlights(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	return nil, nil
}
func (m *MockFactoryAnnotationPort) GetRecentAnnotations(ctx context.Context, userID string, limit int) ([]*reader.Annotation, error) {
	return nil, nil
}
func (m *MockFactoryAnnotationPort) GetPublicAnnotations(ctx context.Context, bookID, chapterID string) ([]*reader.Annotation, error) {
	return nil, nil
}
func (m *MockFactoryAnnotationPort) GetAnnotationStats(ctx context.Context, userID, bookID string) (map[string]interface{}, error) {
	return nil, nil
}
func (m *MockFactoryAnnotationPort) BatchCreateAnnotations(ctx context.Context, annotations []*reader.Annotation) error { return nil }
func (m *MockFactoryAnnotationPort) BatchDeleteAnnotations(ctx context.Context, annotationIDs []string) error {
	return nil
}

// MockFactoryChapterContentPort Mock 实现
type MockFactoryChapterContentPort struct{}

func (m *MockFactoryChapterContentPort) Initialize(ctx context.Context) error { return nil }
func (m *MockFactoryChapterContentPort) Health(ctx context.Context) error    { return nil }
func (m *MockFactoryChapterContentPort) Close(ctx context.Context) error     { return nil }
func (m *MockFactoryChapterContentPort) GetServiceName() string         { return "MockChapterContentPort" }
func (m *MockFactoryChapterContentPort) GetVersion() string             { return "1.0.0" }
func (m *MockFactoryChapterContentPort) GetChapterContent(ctx context.Context, userID, chapterID string) (string, error) {
	return "", nil
}
func (m *MockFactoryChapterContentPort) GetChapterByID(ctx context.Context, chapterID string) (interface{}, error) {
	return nil, nil
}
func (m *MockFactoryChapterContentPort) GetBookChapters(ctx context.Context, bookID string, page, size int) (interface{}, int64, error) {
	return nil, 0, nil
}
func (m *MockFactoryChapterContentPort) GetChapterContentWithProgress(ctx context.Context, req *readeriface.GetChapterContentRequest) (*readeriface.ChapterContentResponse, error) {
	return nil, nil
}
func (m *MockFactoryChapterContentPort) GetChapterByNumber(ctx context.Context, userID, bookID string, chapterNum int) (*readeriface.ChapterContentResponse, error) {
	return nil, nil
}
func (m *MockFactoryChapterContentPort) GetNextChapter(ctx context.Context, userID, bookID, chapterID string) (*readeriface.ChapterInfo, error) {
	return nil, nil
}
func (m *MockFactoryChapterContentPort) GetPreviousChapter(ctx context.Context, userID, bookID, chapterID string) (*readeriface.ChapterInfo, error) {
	return nil, nil
}
func (m *MockFactoryChapterContentPort) GetChapterList(ctx context.Context, userID, bookID string, page, size int) (*readeriface.ChapterListResponse, error) {
	return nil, nil
}
func (m *MockFactoryChapterContentPort) GetChapterInfo(ctx context.Context, userID, chapterID string) (*readeriface.ChapterInfo, error) {
	return nil, nil
}

// MockFactoryReaderSettingsPort Mock 实现
type MockFactoryReaderSettingsPort struct{}

func (m *MockFactoryReaderSettingsPort) Initialize(ctx context.Context) error { return nil }
func (m *MockFactoryReaderSettingsPort) Health(ctx context.Context) error      { return nil }
func (m *MockFactoryReaderSettingsPort) Close(ctx context.Context) error       { return nil }
func (m *MockFactoryReaderSettingsPort) GetServiceName() string           { return "MockReaderSettingsPort" }
func (m *MockFactoryReaderSettingsPort) GetVersion() string               { return "1.0.0" }
func (m *MockFactoryReaderSettingsPort) GetReadingSettings(ctx context.Context, userID string) (*reader.ReadingSettings, error) {
	return nil, nil
}
func (m *MockFactoryReaderSettingsPort) SaveReadingSettings(ctx context.Context, settings *reader.ReadingSettings) error {
	return nil
}
func (m *MockFactoryReaderSettingsPort) UpdateReadingSettings(ctx context.Context, userID string, updates map[string]interface{}) error {
	return nil
}

// MockFactoryReaderSyncPort Mock 实现
type MockFactoryReaderSyncPort struct{}

func (m *MockFactoryReaderSyncPort) Initialize(ctx context.Context) error { return nil }
func (m *MockFactoryReaderSyncPort) Health(ctx context.Context) error      { return nil }
func (m *MockFactoryReaderSyncPort) Close(ctx context.Context) error       { return nil }
func (m *MockFactoryReaderSyncPort) GetServiceName() string           { return "MockReaderSyncPort" }
func (m *MockFactoryReaderSyncPort) GetVersion() string               { return "1.0.0" }
func (m *MockFactoryReaderSyncPort) SyncAnnotations(ctx context.Context, req *readeriface.SyncAnnotationsRequest) (*readeriface.SyncAnnotationsResponse, error) {
	return nil, nil
}
