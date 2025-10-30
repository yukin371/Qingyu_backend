package document_test

import (
	"Qingyu_backend/models/writer"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/repository/interfaces/infrastructure"
	"Qingyu_backend/repository/interfaces/writing"
	documentService "Qingyu_backend/service/document"
	"Qingyu_backend/service/interfaces/base"
)

// ============ Mock Repositories ============

// MockDocumentRepository Mock文档Repository
type MockDocumentRepository struct {
	mock.Mock
	writing.DocumentRepository
}

func (m *MockDocumentRepository) Create(ctx context.Context, doc *writer.Document) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

func (m *MockDocumentRepository) GetByID(ctx context.Context, id string) (*writer.Document, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.Document), args.Error(1)
}

func (m *MockDocumentRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockDocumentRepository) UpdateByProject(ctx context.Context, documentID, projectID string, updates map[string]interface{}) error {
	args := m.Called(ctx, documentID, projectID, updates)
	return args.Error(0)
}

func (m *MockDocumentRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDocumentRepository) GetByProjectID(ctx context.Context, projectID string, limit, offset int64) ([]*writer.Document, error) {
	args := m.Called(ctx, projectID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Document), args.Error(1)
}

func (m *MockDocumentRepository) GetByParentID(ctx context.Context, parentID string) ([]*writer.Document, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Document), args.Error(1)
}

// MockDocumentContentRepository Mock文档内容Repository
type MockDocumentContentRepository struct {
	mock.Mock
	writing.DocumentContentRepository
}

func (m *MockDocumentContentRepository) Create(ctx context.Context, content *writer.DocumentContent) error {
	args := m.Called(ctx, content)
	return args.Error(0)
}

func (m *MockDocumentContentRepository) GetByDocumentID(ctx context.Context, docID string) (*writer.DocumentContent, error) {
	args := m.Called(ctx, docID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.DocumentContent), args.Error(1)
}

func (m *MockDocumentContentRepository) GetByID(ctx context.Context, id string) (*writer.DocumentContent, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.DocumentContent), args.Error(1)
}

func (m *MockDocumentContentRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockDocumentContentRepository) Delete(ctx context.Context, docID string) error {
	args := m.Called(ctx, docID)
	return args.Error(0)
}

func (m *MockDocumentContentRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*writer.DocumentContent, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.DocumentContent), args.Error(1)
}

func (m *MockDocumentContentRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDocumentContentRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDocumentContentRepository) UpdateWithVersion(ctx context.Context, documentID string, content string, expectedVersion int) error {
	args := m.Called(ctx, documentID, content, expectedVersion)
	return args.Error(0)
}

func (m *MockDocumentContentRepository) BatchUpdateContent(ctx context.Context, updates map[string]string) error {
	args := m.Called(ctx, updates)
	return args.Error(0)
}

func (m *MockDocumentContentRepository) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDocumentContentRepository) HealthCheck(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockProjectRepository Mock项目Repository
type MockProjectRepository struct {
	mock.Mock
	writing.ProjectRepository
}

func (m *MockProjectRepository) GetByID(ctx context.Context, id string) (*writer.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.Project), args.Error(1)
}

func (m *MockProjectRepository) UpdateTotalWords(ctx context.Context, id string, totalWords int64) error {
	args := m.Called(ctx, id, totalWords)
	return args.Error(0)
}

func (m *MockProjectRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

// MockEventBus Mock事件总线
type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Publish(ctx context.Context, event base.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventBus) PublishAsync(ctx context.Context, event base.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventBus) Subscribe(eventType string, handler base.EventHandler) error {
	args := m.Called(eventType, handler)
	return args.Error(0)
}

func (m *MockEventBus) Unsubscribe(eventType string, handlerName string) error {
	args := m.Called(eventType, handlerName)
	return args.Error(0)
}

// ============ Test Cases ============

func TestDocumentService_CreateDocument_Success(t *testing.T) {
	// Setup
	mockDocRepo := new(MockDocumentRepository)
	mockContentRepo := new(MockDocumentContentRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	service := documentService.NewDocumentService(mockDocRepo, mockContentRepo, mockProjectRepo, mockEventBus)

	ctx := context.WithValue(context.Background(), "userID", "test-user-id")
	projectID := primitive.NewObjectID().Hex()

	// Mock project exists and user has permission
	mockProject := &writer.Project{
		ID:       projectID,
		Title:    "Test Project",
		AuthorID: "test-user-id",
		Status:   "active",
	}
	mockProjectRepo.On("GetByID", ctx, projectID).Return(mockProject, nil)

	// Mock document creation
	mockDocRepo.On("Create", ctx, mock.AnythingOfType("*writer.Document")).Return(nil)

	// Mock GetByProjectID for statistics update (async)
	mockDocRepo.On("GetByProjectID", mock.Anything, projectID, mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).Return([]*writer.Document{}, nil).Maybe()

	// Mock project update for statistics (async)
	mockProjectRepo.On("Update", mock.Anything, projectID, mock.Anything).Return(nil).Maybe()

	// Mock event publishing
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil).Maybe()

	// Request
	req := &documentService.CreateDocumentRequest{
		ProjectID: projectID,
		Title:     "Test Document",
		Type:      "chapter",
		Order:     1,
	}

	// Execute
	resp, err := service.CreateDocument(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Test Document", resp.Title)
	assert.Equal(t, "chapter", resp.Type)
	mockDocRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)

	// Give async operations time to complete
	time.Sleep(100 * time.Millisecond)
}

func TestDocumentService_CreateDocument_ProjectNotFound(t *testing.T) {
	// Setup
	mockDocRepo := new(MockDocumentRepository)
	mockContentRepo := new(MockDocumentContentRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	service := documentService.NewDocumentService(mockDocRepo, mockContentRepo, mockProjectRepo, mockEventBus)

	ctx := context.WithValue(context.Background(), "userID", "test-user-id")
	projectID := primitive.NewObjectID().Hex()

	// Mock project not found
	mockProjectRepo.On("GetByID", ctx, projectID).Return(nil, nil)

	// Request
	req := &documentService.CreateDocumentRequest{
		ProjectID: projectID,
		Title:     "Test Document",
		Type:      "chapter",
	}

	// Execute
	resp, err := service.CreateDocument(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "项目不存在")
	mockProjectRepo.AssertExpectations(t)
}

func TestDocumentService_CreateDocument_NoPermission(t *testing.T) {
	// Setup
	mockDocRepo := new(MockDocumentRepository)
	mockContentRepo := new(MockDocumentContentRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	service := documentService.NewDocumentService(mockDocRepo, mockContentRepo, mockProjectRepo, mockEventBus)

	ctx := context.WithValue(context.Background(), "userID", "other-user-id")
	projectID := primitive.NewObjectID().Hex()

	// Mock project exists but user has no permission
	mockProject := &writer.Project{
		ID:       projectID,
		Title:    "Test Project",
		AuthorID: "test-user-id", // Different owner
		Status:   "active",
	}
	mockProjectRepo.On("GetByID", ctx, projectID).Return(mockProject, nil)

	// Request
	req := &documentService.CreateDocumentRequest{
		ProjectID: projectID,
		Title:     "Test Document",
		Type:      "chapter",
	}

	// Execute
	resp, err := service.CreateDocument(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "无权限")
	mockProjectRepo.AssertExpectations(t)
}

func TestDocumentService_GetDocument_Success(t *testing.T) {
	// Setup
	mockDocRepo := new(MockDocumentRepository)
	mockContentRepo := new(MockDocumentContentRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	service := documentService.NewDocumentService(mockDocRepo, mockContentRepo, mockProjectRepo, mockEventBus)

	ctx := context.WithValue(context.Background(), "userID", "test-user-id")
	docID := primitive.NewObjectID().Hex()
	projectID := primitive.NewObjectID().Hex()

	// Mock document exists
	mockDoc := &writer.Document{
		ID:        docID,
		ProjectID: projectID,
		Title:     "Test Document",
		Type:      "chapter",
		WordCount: 100,
		Status:    "draft",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockDocRepo.On("GetByID", ctx, docID).Return(mockDoc, nil)

	// Mock project exists and user has permission
	mockProject := &writer.Project{
		ID:       projectID,
		Title:    "Test Project",
		AuthorID: "test-user-id",
		Status:   "active",
	}
	mockProjectRepo.On("GetByID", ctx, projectID).Return(mockProject, nil)

	// Execute
	resp, err := service.GetDocument(ctx, docID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Test Document", resp.Title)
	assert.Equal(t, writer.DocumentType("chapter"), resp.Type)
	assert.Equal(t, 100, resp.WordCount)
	mockDocRepo.AssertExpectations(t)
}

func TestDocumentService_GetDocument_NotFound(t *testing.T) {
	// Setup
	mockDocRepo := new(MockDocumentRepository)
	mockContentRepo := new(MockDocumentContentRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	service := documentService.NewDocumentService(mockDocRepo, mockContentRepo, mockProjectRepo, mockEventBus)

	ctx := context.Background()
	docID := primitive.NewObjectID().Hex()

	// Mock document not found
	mockDocRepo.On("GetByID", ctx, docID).Return(nil, nil)

	// Execute
	resp, err := service.GetDocument(ctx, docID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "文档不存在")
	mockDocRepo.AssertExpectations(t)
}

func TestDocumentService_UpdateDocument_Success(t *testing.T) {
	// Setup
	mockDocRepo := new(MockDocumentRepository)
	mockContentRepo := new(MockDocumentContentRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	service := documentService.NewDocumentService(mockDocRepo, mockContentRepo, mockProjectRepo, mockEventBus)

	ctx := context.WithValue(context.Background(), "userID", "test-user-id")
	docID := primitive.NewObjectID().Hex()
	projectID := primitive.NewObjectID().Hex()

	// Mock document exists
	mockDoc := &writer.Document{
		ID:        docID,
		ProjectID: projectID,
		Title:     "Old Title",
		Type:      "chapter",
		Status:    "draft",
	}
	mockDocRepo.On("GetByID", ctx, docID).Return(mockDoc, nil)

	// Mock project exists and user has permission
	mockProject := &writer.Project{
		ID:       projectID,
		AuthorID: "test-user-id",
		Status:   "active",
	}
	mockProjectRepo.On("GetByID", ctx, projectID).Return(mockProject, nil)

	// Mock update by project
	mockDocRepo.On("UpdateByProject", ctx, docID, mockProject.ID, mock.AnythingOfType("map[string]interface {}")).Return(nil)

	// Mock event
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil)

	// Request
	req := &documentService.UpdateDocumentRequest{
		Title: "New Title",
	}

	// Execute
	err := service.UpdateDocument(ctx, docID, req)

	// Assert
	require.NoError(t, err)
	mockDocRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)
}

func TestDocumentService_DeleteDocument_Success(t *testing.T) {
	// Setup
	mockDocRepo := new(MockDocumentRepository)
	mockContentRepo := new(MockDocumentContentRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	service := documentService.NewDocumentService(mockDocRepo, mockContentRepo, mockProjectRepo, mockEventBus)

	ctx := context.WithValue(context.Background(), "userID", "test-user-id")
	docID := primitive.NewObjectID().Hex()
	projectID := primitive.NewObjectID().Hex()

	// Mock document exists
	mockDoc := &writer.Document{
		ID:        docID,
		ProjectID: projectID,
		Title:     "Test Document",
		Type:      "chapter",
		Status:    "draft",
	}
	mockDocRepo.On("GetByID", ctx, docID).Return(mockDoc, nil)

	// Mock project exists and user has permission
	mockProject := &writer.Project{
		ID:       projectID,
		AuthorID: "test-user-id",
		Status:   "active",
	}
	mockProjectRepo.On("GetByID", ctx, projectID).Return(mockProject, nil)

	// Mock no children
	mockDocRepo.On("GetByParentID", ctx, docID).Return([]*writer.Document{}, nil)

	// Mock soft delete
	mockDocRepo.On("SoftDelete", ctx, docID, projectID).Return(nil)

	// Mock event
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil)

	// Execute
	err := service.DeleteDocument(ctx, docID)

	// Assert
	require.NoError(t, err)
	mockDocRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)
}

func TestDocumentService_DeleteDocument_HasChildren(t *testing.T) {
	// Setup
	mockDocRepo := new(MockDocumentRepository)
	mockContentRepo := new(MockDocumentContentRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	service := documentService.NewDocumentService(mockDocRepo, mockContentRepo, mockProjectRepo, mockEventBus)

	ctx := context.WithValue(context.Background(), "userID", "test-user-id")
	docID := primitive.NewObjectID().Hex()
	projectID := primitive.NewObjectID().Hex()

	// Mock document exists
	mockDoc := &writer.Document{
		ID:        docID,
		ProjectID: projectID,
		Title:     "Test Document",
		Type:      "volume",
		Status:    "draft",
	}
	mockDocRepo.On("GetByID", ctx, docID).Return(mockDoc, nil)

	// Mock project exists and user has permission
	mockProject := &writer.Project{
		ID:       projectID,
		AuthorID: "test-user-id",
		Status:   "active",
	}
	mockProjectRepo.On("GetByID", ctx, projectID).Return(mockProject, nil)

	// Mock has children
	children := []*writer.Document{
		{ID: primitive.NewObjectID().Hex(), Title: "Child 1"},
	}
	mockDocRepo.On("GetByParentID", ctx, docID).Return(children, nil)

	// Execute
	err := service.DeleteDocument(ctx, docID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "请先删除子文档")
	mockDocRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)
}

func TestDocumentService_ListDocuments_Success(t *testing.T) {
	// Setup
	mockDocRepo := new(MockDocumentRepository)
	mockContentRepo := new(MockDocumentContentRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	service := documentService.NewDocumentService(mockDocRepo, mockContentRepo, mockProjectRepo, mockEventBus)

	ctx := context.Background()
	projectID := primitive.NewObjectID().Hex()

	// Mock documents exist
	mockDocs := []*writer.Document{
		{
			ID:        primitive.NewObjectID().Hex(),
			ProjectID: projectID,
			Title:     "Document 1",
			Type:      "chapter",
			Order:     1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID().Hex(),
			ProjectID: projectID,
			Title:     "Document 2",
			Type:      "chapter",
			Order:     2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	mockDocRepo.On("GetByProjectID", ctx, projectID, mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).Return(mockDocs, nil)

	// Execute
	req := &documentService.ListDocumentsRequest{
		ProjectID: projectID,
	}
	resp, err := service.ListDocuments(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Documents, 2)
	assert.Equal(t, "Document 1", resp.Documents[0].Title)
	assert.Equal(t, "Document 2", resp.Documents[1].Title)
	mockDocRepo.AssertExpectations(t)
}

func TestDocumentService_ListDocuments_Error(t *testing.T) {
	// Setup
	mockDocRepo := new(MockDocumentRepository)
	mockContentRepo := new(MockDocumentContentRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	service := documentService.NewDocumentService(mockDocRepo, mockContentRepo, mockProjectRepo, mockEventBus)

	ctx := context.Background()
	projectID := primitive.NewObjectID().Hex()

	// Mock error
	mockDocRepo.On("GetByProjectID", ctx, projectID, mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).Return(nil, errors.New("database error"))

	// Execute
	req := &documentService.ListDocumentsRequest{
		ProjectID: projectID,
	}
	resp, err := service.ListDocuments(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	mockDocRepo.AssertExpectations(t)
}
