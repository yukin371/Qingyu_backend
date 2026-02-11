package reader

import (
	"testing"

	readeriface "Qingyu_backend/service/interfaces/reader"
)

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

func (m *MockFactoryReadingProgressPort) Initialize(ctx interface{}) error                           { return nil }
func (m *MockFactoryReadingProgressPort) Health(ctx interface{}) error                                { return nil }
func (m *MockFactoryReadingProgressPort) Close(ctx interface{}) error                                 { return nil }
func (m *MockFactoryReadingProgressPort) GetServiceName() string                                     { return "MockProgressPort" }
func (m *MockFactoryReadingProgressPort) GetVersion() string                                         { return "1.0.0" }
func (m *MockFactoryReadingProgressPort) GetReadingProgress(ctx interface{}, userID, bookID string) (interface{}, error) {
	return nil, nil
}
func (m *MockFactoryReadingProgressPort) SaveReadingProgress(ctx interface{}, req *readeriface.SaveReadingProgressRequest) error {
	return nil
}
func (m *MockFactoryReadingProgressPort) UpdateReadingTime(ctx interface{}, req *readeriface.UpdateReadingTimeRequest) error {
	return nil
}
func (m *MockFactoryReadingProgressPort) GetRecentReading(ctx interface{}, userID string, limit int) (interface{}, error) {
	return nil, nil
}
func (m *MockFactoryReadingProgressPort) GetReadingHistory(ctx interface{}, req *readeriface.GetReadingHistoryRequest) (*readeriface.GetReadingHistoryResponse, error) {
	return nil, nil
}
func (m *MockFactoryReadingProgressPort) GetTotalReadingTime(ctx interface{}, userID string) (int64, error) { return 0, nil }
func (m *MockFactoryReadingProgressPort) GetReadingTimeByPeriod(ctx interface{}, req *readeriface.GetReadingTimeByPeriodRequest) (int64, error) {
	return 0, nil
}
func (m *MockFactoryReadingProgressPort) GetUnfinishedBooks(ctx interface{}, userID string) (interface{}, error) {
	return nil, nil
}
func (m *MockFactoryReadingProgressPort) GetFinishedBooks(ctx interface{}, userID string) (interface{}, error) {
	return nil, nil
}
func (m *MockFactoryReadingProgressPort) DeleteReadingProgress(ctx interface{}, userID, bookID string) error { return nil }
func (m *MockFactoryReadingProgressPort) UpdateBookStatus(ctx interface{}, req *readeriface.UpdateBookStatusRequest) error {
	return nil
}
func (m *MockFactoryReadingProgressPort) BatchUpdateBookStatus(ctx interface{}, req *readeriface.BatchUpdateBookStatusRequest) error {
	return nil
}

// MockFactoryAnnotationPort Mock 实现
type MockFactoryAnnotationPort struct{}

func (m *MockFactoryAnnotationPort) Initialize(ctx interface{}) error { return nil }
func (m *MockFactoryAnnotationPort) Health(ctx interface{}) error      { return nil }
func (m *MockFactoryAnnotationPort) Close(ctx interface{}) error       { return nil }
func (m *MockFactoryAnnotationPort) GetServiceName() string           { return "MockAnnotationPort" }
func (m *MockFactoryAnnotationPort) GetVersion() string               { return "1.0.0" }
func (m *MockFactoryAnnotationPort) CreateAnnotation(ctx interface{}, annotation interface{}) error {
	return nil
}
func (m *MockFactoryAnnotationPort) UpdateAnnotation(ctx interface{}, annotationID string, updates map[string]interface{}) error {
	return nil
}
func (m *MockFactoryAnnotationPort) DeleteAnnotation(ctx interface{}, annotationID string) error { return nil }
func (m *MockFactoryAnnotationPort) GetAnnotationsByChapter(ctx interface{}, userID, bookID, chapterID string) (interface{}, error) {
	return nil, nil
}
func (m *MockFactoryAnnotationPort) GetAnnotationsByBook(ctx interface{}, userID, bookID string) (interface{}, error) {
	return nil, nil
}
func (m *MockFactoryAnnotationPort) GetNotes(ctx interface{}, userID, bookID string) (interface{}, error) {
	return nil, nil
}
func (m *MockFactoryAnnotationPort) SearchNotes(ctx interface{}, userID, keyword string) (interface{}, error) {
	return nil, nil
}
func (m *MockFactoryAnnotationPort) GetBookmarks(ctx interface{}, userID, bookID string) (interface{}, error) {
	return nil, nil
}
func (m *MockFactoryAnnotationPort) GetLatestBookmark(ctx interface{}, userID, bookID string) (interface{}, error) {
	return nil, nil
}
func (m *MockFactoryAnnotationPort) GetHighlights(ctx interface{}, userID, bookID string) (interface{}, error) {
	return nil, nil
}
func (m *MockFactoryAnnotationPort) GetRecentAnnotations(ctx interface{}, userID string, limit int) (interface{}, error) {
	return nil, nil
}
func (m *MockFactoryAnnotationPort) GetPublicAnnotations(ctx interface{}, bookID, chapterID string) (interface{}, error) {
	return nil, nil
}
func (m *MockFactoryAnnotationPort) GetAnnotationStats(ctx interface{}, userID, bookID string) (map[string]interface{}, error) {
	return nil, nil
}
func (m *MockFactoryAnnotationPort) BatchCreateAnnotations(ctx interface{}, annotations interface{}) error { return nil }
func (m *MockFactoryAnnotationPort) BatchDeleteAnnotations(ctx interface{}, annotationIDs []string) error {
	return nil
}

// MockFactoryChapterContentPort Mock 实现
type MockFactoryChapterContentPort struct{}

func (m *MockFactoryChapterContentPort) Initialize(ctx interface{}) error { return nil }
func (m *MockFactoryChapterContentPort) Health(ctx interface{}) error    { return nil }
func (m *MockFactoryChapterContentPort) Close(ctx interface{}) error     { return nil }
func (m *MockFactoryChapterContentPort) GetServiceName() string         { return "MockChapterContentPort" }
func (m *MockFactoryChapterContentPort) GetVersion() string             { return "1.0.0" }
func (m *MockFactoryChapterContentPort) GetChapterContent(ctx interface{}, userID, chapterID string) (string, error) {
	return "", nil
}
func (m *MockFactoryChapterContentPort) GetChapterByID(ctx interface{}, chapterID string) (interface{}, error) {
	return nil, nil
}
func (m *MockFactoryChapterContentPort) GetBookChapters(ctx interface{}, bookID string, page, size int) (interface{}, int64, error) {
	return nil, 0, nil
}
func (m *MockFactoryChapterContentPort) GetChapterContentWithProgress(ctx interface{}, req *readeriface.GetChapterContentRequest) (*readeriface.ChapterContentResponse, error) {
	return nil, nil
}
func (m *MockFactoryChapterContentPort) GetChapterByNumber(ctx interface{}, userID, bookID string, chapterNum int) (*readeriface.ChapterContentResponse, error) {
	return nil, nil
}
func (m *MockFactoryChapterContentPort) GetNextChapter(ctx interface{}, userID, bookID, chapterID string) (*readeriface.ChapterInfo, error) {
	return nil, nil
}
func (m *MockFactoryChapterContentPort) GetPreviousChapter(ctx interface{}, userID, bookID, chapterID string) (*readeriface.ChapterInfo, error) {
	return nil, nil
}
func (m *MockFactoryChapterContentPort) GetChapterList(ctx interface{}, userID, bookID string, page, size int) (*readeriface.ChapterListResponse, error) {
	return nil, nil
}
func (m *MockFactoryChapterContentPort) GetChapterInfo(ctx interface{}, userID, chapterID string) (*readeriface.ChapterInfo, error) {
	return nil, nil
}

// MockFactoryReaderSettingsPort Mock 实现
type MockFactoryReaderSettingsPort struct{}

func (m *MockFactoryReaderSettingsPort) Initialize(ctx interface{}) error { return nil }
func (m *MockFactoryReaderSettingsPort) Health(ctx interface{}) error      { return nil }
func (m *MockFactoryReaderSettingsPort) Close(ctx interface{}) error       { return nil }
func (m *MockFactoryReaderSettingsPort) GetServiceName() string           { return "MockReaderSettingsPort" }
func (m *MockFactoryReaderSettingsPort) GetVersion() string               { return "1.0.0" }
func (m *MockFactoryReaderSettingsPort) GetReadingSettings(ctx interface{}, userID string) (interface{}, error) {
	return nil, nil
}
func (m *MockFactoryReaderSettingsPort) SaveReadingSettings(ctx interface{}, settings interface{}) error {
	return nil
}
func (m *MockFactoryReaderSettingsPort) UpdateReadingSettings(ctx interface{}, userID string, updates map[string]interface{}) error {
	return nil
}

// MockFactoryReaderSyncPort Mock 实现
type MockFactoryReaderSyncPort struct{}

func (m *MockFactoryReaderSyncPort) Initialize(ctx interface{}) error { return nil }
func (m *MockFactoryReaderSyncPort) Health(ctx interface{}) error      { return nil }
func (m *MockFactoryReaderSyncPort) Close(ctx interface{}) error       { return nil }
func (m *MockFactoryReaderSyncPort) GetServiceName() string           { return "MockReaderSyncPort" }
func (m *MockFactoryReaderSyncPort) GetVersion() string               { return "1.0.0" }
func (m *MockFactoryReaderSyncPort) SyncAnnotations(ctx interface{}, req *readeriface.SyncAnnotationsRequest) (*readeriface.SyncAnnotationsResponse, error) {
	return nil, nil
}
