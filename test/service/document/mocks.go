package document_test

import (
	"context"

	"github.com/stretchr/testify/mock"

	writerModels "Qingyu_backend/models/writer"
	"Qingyu_backend/models/writer/base"
	modelBase "Qingyu_backend/models/writer/base"
	"Qingyu_backend/repository/interfaces/infrastructure"
	serviceBase "Qingyu_backend/service/interfaces/base"
)

// MockDocumentRepository Mock文档Repository
type MockDocumentRepository struct {
	mock.Mock
}

func (m *MockDocumentRepository) Create(ctx context.Context, doc *writerModels.Document) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

func (m *MockDocumentRepository) GetByID(ctx context.Context, id string) (*writerModels.Document, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writerModels.Document), args.Error(1)
}

func (m *MockDocumentRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockDocumentRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDocumentRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*writerModels.Document, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writerModels.Document), args.Error(1)
}

func (m *MockDocumentRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDocumentRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDocumentRepository) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDocumentRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDocumentRepository) GetByProjectID(ctx context.Context, projectID string, limit, offset int64) ([]*writerModels.Document, error) {
	args := m.Called(ctx, projectID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writerModels.Document), args.Error(1)
}

func (m *MockDocumentRepository) GetByProjectAndType(ctx context.Context, projectID, documentType string, limit, offset int64) ([]*writerModels.Document, error) {
	args := m.Called(ctx, projectID, documentType, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writerModels.Document), args.Error(1)
}

func (m *MockDocumentRepository) UpdateByProject(ctx context.Context, documentID, projectID string, updates map[string]interface{}) error {
	args := m.Called(ctx, documentID, projectID, updates)
	return args.Error(0)
}

func (m *MockDocumentRepository) DeleteByProject(ctx context.Context, documentID, projectID string) error {
	args := m.Called(ctx, documentID, projectID)
	return args.Error(0)
}

func (m *MockDocumentRepository) RestoreByProject(ctx context.Context, documentID, projectID string) error {
	args := m.Called(ctx, documentID, projectID)
	return args.Error(0)
}

func (m *MockDocumentRepository) IsProjectMember(ctx context.Context, documentID, projectID string) (bool, error) {
	args := m.Called(ctx, documentID, projectID)
	return args.Bool(0), args.Error(1)
}

func (m *MockDocumentRepository) SoftDelete(ctx context.Context, documentID, projectID string) error {
	args := m.Called(ctx, documentID, projectID)
	return args.Error(0)
}

func (m *MockDocumentRepository) HardDelete(ctx context.Context, documentID string) error {
	args := m.Called(ctx, documentID)
	return args.Error(0)
}

func (m *MockDocumentRepository) CountByProject(ctx context.Context, projectID string) (int64, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDocumentRepository) CreateWithTransaction(ctx context.Context, document *writerModels.Document, callback func(ctx context.Context) error) error {
	args := m.Called(ctx, document, callback)
	return args.Error(0)
}

func (m *MockDocumentRepository) GetByParentID(ctx context.Context, parentID string) ([]*writerModels.Document, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writerModels.Document), args.Error(1)
}

// MockDocumentContentRepository Mock文档内容Repository
type MockDocumentContentRepository struct {
	mock.Mock
}

func (m *MockDocumentContentRepository) Create(ctx context.Context, content *writerModels.DocumentContent) error {
	args := m.Called(ctx, content)
	return args.Error(0)
}

func (m *MockDocumentContentRepository) GetByID(ctx context.Context, id string) (*writerModels.DocumentContent, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writerModels.DocumentContent), args.Error(1)
}

func (m *MockDocumentContentRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockDocumentContentRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDocumentContentRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*writerModels.DocumentContent, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writerModels.DocumentContent), args.Error(1)
}

func (m *MockDocumentContentRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDocumentContentRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDocumentContentRepository) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDocumentContentRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDocumentContentRepository) GetByDocumentID(ctx context.Context, documentID string) (*writerModels.DocumentContent, error) {
	args := m.Called(ctx, documentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writerModels.DocumentContent), args.Error(1)
}

func (m *MockDocumentContentRepository) UpdateWithVersion(ctx context.Context, documentID string, content string, expectedVersion int) error {
	args := m.Called(ctx, documentID, content, expectedVersion)
	return args.Error(0)
}

func (m *MockDocumentContentRepository) BatchUpdateContent(ctx context.Context, updates map[string]string) error {
	args := m.Called(ctx, updates)
	return args.Error(0)
}

func (m *MockDocumentContentRepository) GetContentStats(ctx context.Context, documentID string) (int, int, error) {
	args := m.Called(ctx, documentID)
	return args.Int(0), args.Int(1), args.Error(2)
}

func (m *MockDocumentContentRepository) StoreToGridFS(ctx context.Context, documentID string, content []byte) (string, error) {
	args := m.Called(ctx, documentID, content)
	return args.String(0), args.Error(1)
}

func (m *MockDocumentContentRepository) LoadFromGridFS(ctx context.Context, gridFSID string) ([]byte, error) {
	args := m.Called(ctx, gridFSID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockDocumentContentRepository) CreateWithTransaction(ctx context.Context, content *writerModels.DocumentContent, callback func(ctx context.Context) error) error {
	args := m.Called(ctx, content, callback)
	return args.Error(0)
}

// MockProjectRepository Mock项目Repository
type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) GetByID(ctx context.Context, id string) (*writerModels.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writerModels.Project), args.Error(1)
}

func (m *MockProjectRepository) Create(ctx context.Context, project *writerModels.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockProjectRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProjectRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*writerModels.Project, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writerModels.Project), args.Error(1)
}

func (m *MockProjectRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProjectRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockProjectRepository) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockProjectRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockProjectRepository) GetListByOwnerID(ctx context.Context, ownerID string, limit, offset int64) ([]*writerModels.Project, error) {
	args := m.Called(ctx, ownerID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writerModels.Project), args.Error(1)
}

func (m *MockProjectRepository) GetByOwnerAndStatus(ctx context.Context, ownerID, status string, limit, offset int64) ([]*writerModels.Project, error) {
	args := m.Called(ctx, ownerID, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writerModels.Project), args.Error(1)
}

func (m *MockProjectRepository) UpdateByOwner(ctx context.Context, projectID, ownerID string, updates map[string]interface{}) error {
	args := m.Called(ctx, projectID, ownerID, updates)
	return args.Error(0)
}

func (m *MockProjectRepository) IsOwner(ctx context.Context, projectID, ownerID string) (bool, error) {
	args := m.Called(ctx, projectID, ownerID)
	return args.Bool(0), args.Error(1)
}

func (m *MockProjectRepository) SoftDelete(ctx context.Context, projectID, ownerID string) error {
	args := m.Called(ctx, projectID, ownerID)
	return args.Error(0)
}

func (m *MockProjectRepository) HardDelete(ctx context.Context, projectID string) error {
	args := m.Called(ctx, projectID)
	return args.Error(0)
}

func (m *MockProjectRepository) Restore(ctx context.Context, projectID, ownerID string) error {
	args := m.Called(ctx, projectID, ownerID)
	return args.Error(0)
}

func (m *MockProjectRepository) CountByOwner(ctx context.Context, ownerID string) (int64, error) {
	args := m.Called(ctx, ownerID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProjectRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProjectRepository) CreateWithTransaction(ctx context.Context, project *writerModels.Project, callback func(ctx context.Context) error) error {
	args := m.Called(ctx, project, callback)
	return args.Error(0)
}

func (m *MockProjectRepository) UpdateTotalWords(ctx context.Context, id string, totalWords int64) error {
	args := m.Called(ctx, id, totalWords)
	return args.Error(0)
}

// MockEventBus Mock事件总线
type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Publish(ctx context.Context, event serviceBase.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventBus) PublishAsync(ctx context.Context, event serviceBase.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventBus) Subscribe(eventType string, handler serviceBase.EventHandler) error {
	args := m.Called(eventType, handler)
	return args.Error(0)
}

func (m *MockEventBus) Unsubscribe(eventType string, handlerName string) error {
	args := m.Called(eventType, handlerName)
	return args.Error(0)
}
