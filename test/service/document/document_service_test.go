package document_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	writerModels "Qingyu_backend/models/writer"
	"Qingyu_backend/models/writer/base"
	modelBase "Qingyu_backend/models/writer/base"
	documentService "Qingyu_backend/service/writer/document"
)

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
	mockProject := &writerModels.Project{
		IdentifiedEntity: base.IdentifiedEntity{ID: projectID},
		TitledEntity:     base.TitledEntity{Title: "Test Project"},
		OwnedEntity:      base.OwnedEntity{AuthorID: "test-user-id"},
	}
	mockProjectRepo.On("GetByID", ctx, projectID).Return(mockProject, nil)

	// Mock document creation
	mockDocRepo.On("Create", ctx, mock.Anything).Return(nil)

	// Mock GetByProjectID for statistics update (async)
	mockDocRepo.On("GetByProjectID", mock.Anything, projectID, mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).Return([]*writerModels.Document{}, nil).Maybe()

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
	mockProject := &writerModels.Project{
		IdentifiedEntity: base.IdentifiedEntity{ID: projectID},
		TitledEntity:     base.TitledEntity{Title: "Test Project"},
		OwnedEntity:      base.OwnedEntity{AuthorID: "test-user-id"}, // Different owner
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
	mockDoc := &writerModels.Document{
		IdentifiedEntity:     modelBase.IdentifiedEntity{ID: docID},
		ProjectScopedEntity: modelBase.ProjectScopedEntity{ProjectID: projectID},
		TitledEntity:         modelBase.TitledEntity{Title: "Test Document"},
		Timestamps:           modelBase.Timestamps{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Type:                 "chapter",
		WordCount:            100,
		Status:               "draft",
	}
	mockDocRepo.On("GetByID", ctx, docID).Return(mockDoc, nil)

	// Mock project exists and user has permission
	mockProject := &writerModels.Project{
		IdentifiedEntity: base.IdentifiedEntity{ID: projectID},
		TitledEntity:     base.TitledEntity{Title: "Test Project"},
		OwnedEntity:      base.OwnedEntity{AuthorID: "test-user-id"},
	}
	mockProjectRepo.On("GetByID", ctx, projectID).Return(mockProject, nil)

	// Execute
	resp, err := service.GetDocument(ctx, docID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Test Document", resp.Title)
	assert.Equal(t, "chapter", resp.Type)
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
	mockDoc := &writerModels.Document{
		IdentifiedEntity:     modelBase.IdentifiedEntity{ID: docID},
		ProjectScopedEntity: modelBase.ProjectScopedEntity{ProjectID: projectID},
		TitledEntity:         modelBase.TitledEntity{Title: "Old Title"},
		Type:                 "chapter",
		Status:               "draft",
	}
	mockDocRepo.On("GetByID", ctx, docID).Return(mockDoc, nil)

	// Mock project exists and user has permission
	mockProject := &writerModels.Project{
		IdentifiedEntity: base.IdentifiedEntity{ID: projectID},
		TitledEntity:     base.TitledEntity{Title: "Test Project"},
		OwnedEntity:      base.OwnedEntity{AuthorID: "test-user-id"},
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
	mockDoc := &writerModels.Document{
		IdentifiedEntity:     modelBase.IdentifiedEntity{ID: docID},
		ProjectScopedEntity: modelBase.ProjectScopedEntity{ProjectID: projectID},
		TitledEntity:         modelBase.TitledEntity{Title: "Test Document"},
		Type:                 "chapter",
		Status:               "draft",
	}
	mockDocRepo.On("GetByID", ctx, docID).Return(mockDoc, nil)

	// Mock project exists and user has permission
	mockProject := &writerModels.Project{
		IdentifiedEntity: base.IdentifiedEntity{ID: projectID},
		TitledEntity:     base.TitledEntity{Title: "Test Project"},
		OwnedEntity:      base.OwnedEntity{AuthorID: "test-user-id"},
	}
	mockProjectRepo.On("GetByID", ctx, projectID).Return(mockProject, nil)

	// Mock no children
	mockDocRepo.On("GetByParentID", ctx, docID).Return([]*writerModels.Document{}, nil)

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
	mockDoc := &writerModels.Document{
		IdentifiedEntity:     modelBase.IdentifiedEntity{ID: docID},
		ProjectScopedEntity: modelBase.ProjectScopedEntity{ProjectID: projectID},
		TitledEntity:         modelBase.TitledEntity{Title: "Test Document"},
		Type:                 "volume",
		Status:               "draft",
	}
	mockDocRepo.On("GetByID", ctx, docID).Return(mockDoc, nil)

	// Mock project exists and user has permission
	mockProject := &writerModels.Project{
		IdentifiedEntity: base.IdentifiedEntity{ID: projectID},
		TitledEntity:     base.TitledEntity{Title: "Test Project"},
		OwnedEntity:      base.OwnedEntity{AuthorID: "test-user-id"},
	}
	mockProjectRepo.On("GetByID", ctx, projectID).Return(mockProject, nil)

	// Mock has children
	children := []*writerModels.Document{
		{IdentifiedEntity: modelBase.IdentifiedEntity{ID: primitive.NewObjectID().Hex()}, TitledEntity: modelBase.TitledEntity{Title: "Child 1"}},
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
	mockDocs := []*writerModels.Document{
		{
			IdentifiedEntity:     modelBase.IdentifiedEntity{ID: primitive.NewObjectID().Hex()},
			ProjectScopedEntity: modelBase.ProjectScopedEntity{ProjectID: projectID},
			TitledEntity:         modelBase.TitledEntity{Title: "Document 1"},
			Timestamps:           modelBase.Timestamps{CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Type:                 "chapter",
			Order:                1,
		},
		{
			IdentifiedEntity:     modelBase.IdentifiedEntity{ID: primitive.NewObjectID().Hex()},
			ProjectScopedEntity: modelBase.ProjectScopedEntity{ProjectID: projectID},
			TitledEntity:         modelBase.TitledEntity{Title: "Document 2"},
			Timestamps:           modelBase.Timestamps{CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Type:                 "chapter",
			Order:                2,
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
