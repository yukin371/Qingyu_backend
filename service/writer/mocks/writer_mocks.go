package mocks

import (
	"context"
	"io"
	"time"

	"github.com/stretchr/testify/mock"

	"Qingyu_backend/models/writer"
	serviceInterfaces "Qingyu_backend/service/interfaces"
)

// MockDocumentRepository Mock文档仓储
type MockDocumentRepository struct {
	mock.Mock
}

func (m *MockDocumentRepository) FindByID(ctx context.Context, id string) (*writer.Document, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.Document), args.Error(1)
}

func (m *MockDocumentRepository) FindByProjectID(ctx context.Context, projectID string) ([]*writer.Document, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Document), args.Error(1)
}

// MockDocumentContentRepository Mock文档内容仓储
type MockDocumentContentRepository struct {
	mock.Mock
}

func (m *MockDocumentContentRepository) FindByID(ctx context.Context, id string) (*writer.DocumentContent, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.DocumentContent), args.Error(1)
}

// MockProjectRepository Mock项目仓储
type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) FindByID(ctx context.Context, id string) (*writer.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.Project), args.Error(1)
}

// MockExportTaskRepository Mock导出任务仓储
type MockExportTaskRepository struct {
	mock.Mock
}

func (m *MockExportTaskRepository) Create(ctx context.Context, task *serviceInterfaces.ExportTask) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockExportTaskRepository) FindByID(ctx context.Context, id string) (*serviceInterfaces.ExportTask, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*serviceInterfaces.ExportTask), args.Error(1)
}

func (m *MockExportTaskRepository) FindByProjectID(ctx context.Context, projectID string, page, pageSize int) ([]*serviceInterfaces.ExportTask, int64, error) {
	args := m.Called(ctx, projectID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*serviceInterfaces.ExportTask), args.Get(1).(int64), args.Error(2)
}

func (m *MockExportTaskRepository) Update(ctx context.Context, task *serviceInterfaces.ExportTask) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockExportTaskRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockExportTaskRepository) FindByUser(ctx context.Context, userID string, page, pageSize int) ([]*serviceInterfaces.ExportTask, int64, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*serviceInterfaces.ExportTask), args.Get(1).(int64), args.Error(2)
}

// MockFileStorage Mock文件存储
type MockFileStorage struct {
	mock.Mock
}

func (m *MockFileStorage) Upload(ctx context.Context, filename string, content io.Reader, mimeType string) (string, error) {
	args := m.Called(ctx, filename, content, mimeType)
	return args.String(0), args.Error(1)
}

func (m *MockFileStorage) Download(ctx context.Context, url string) (io.ReadCloser, error) {
	args := m.Called(ctx, url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockFileStorage) Delete(ctx context.Context, url string) error {
	args := m.Called(ctx, url)
	return args.Error(0)
}

func (m *MockFileStorage) GetSignedURL(ctx context.Context, url string, expiration time.Duration) (string, error) {
	args := m.Called(ctx, url, expiration)
	return args.String(0), args.Error(1)
}

// MockPublicationRepository Mock发布记录仓储
type MockPublicationRepository struct {
	mock.Mock
}

func (m *MockPublicationRepository) Create(ctx context.Context, record *serviceInterfaces.PublicationRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockPublicationRepository) FindByID(ctx context.Context, id string) (*serviceInterfaces.PublicationRecord, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*serviceInterfaces.PublicationRecord), args.Error(1)
}

func (m *MockPublicationRepository) FindByProjectID(ctx context.Context, projectID string, page, pageSize int) ([]*serviceInterfaces.PublicationRecord, int64, error) {
	args := m.Called(ctx, projectID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*serviceInterfaces.PublicationRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockPublicationRepository) FindByResourceID(ctx context.Context, resourceID string) (*serviceInterfaces.PublicationRecord, error) {
	args := m.Called(ctx, resourceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*serviceInterfaces.PublicationRecord), args.Error(1)
}

func (m *MockPublicationRepository) Update(ctx context.Context, record *serviceInterfaces.PublicationRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockPublicationRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPublicationRepository) FindPublishedByProjectID(ctx context.Context, projectID string) (*serviceInterfaces.PublicationRecord, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*serviceInterfaces.PublicationRecord), args.Error(1)
}

// MockEventBus Mock事件总线
type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) PublishAsync(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}
