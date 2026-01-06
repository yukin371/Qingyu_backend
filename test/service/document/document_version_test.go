package document_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	writerModels "Qingyu_backend/models/writer"
	"Qingyu_backend/models/writer/base"
	modelBase "Qingyu_backend/models/writer/base"
	documentService "Qingyu_backend/service/writer/document"
)

// ============ Test Cases ============

func TestDocumentService_AutoSave_Success(t *testing.T) {
	// Setup
	mockDocRepo := new(MockDocumentRepository)
	mockContentRepo := new(MockDocumentContentRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	service := documentService.NewDocumentService(mockDocRepo, mockContentRepo, mockProjectRepo, mockEventBus)

	ctx := context.Background()
	docID := primitive.NewObjectID().Hex()

	// Mock document exists
	mockDoc := &writerModels.Document{
		IdentifiedEntity:  modelBase.IdentifiedEntity{ID: docID},
		TitledEntity:       modelBase.TitledEntity{Title: "Test Document"},
		Timestamps:         modelBase.Timestamps{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Type:               "chapter",
		Status:             "draft",
	}
	mockDocRepo.On("GetByID", ctx, docID).Return(mockDoc, nil)

	// Mock content update with version
	mockContentRepo.On("UpdateWithVersion", ctx, docID, "New content", 1).Return(nil)

	// Mock document update
	mockDocRepo.On("Update", ctx, docID, mock.Anything).Return(nil)

	// Mock event publishing
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil).Maybe()

	// Request
	req := &documentService.AutoSaveRequest{
		DocumentID:   docID,
		Content:      "New content",
		WordCount:    100,
		ExpectedVersion: 1,
	}

	// Execute
	resp, err := service.AutoSaveDocument(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	mockDocRepo.AssertExpectations(t)
	mockContentRepo.AssertExpectations(t)
}

func TestDocumentService_AutoSave_VersionConflict(t *testing.T) {
	// Setup
	mockDocRepo := new(MockDocumentRepository)
	mockContentRepo := new(MockDocumentContentRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockEventBus := new(MockEventBus)

	service := documentService.NewDocumentService(mockDocRepo, mockContentRepo, mockProjectRepo, mockEventBus)

	ctx := context.Background()
	docID := primitive.NewObjectID().Hex()

	// Mock document exists
	mockDoc := &writerModels.Document{
		IdentifiedEntity:  modelBase.IdentifiedEntity{ID: docID},
		TitledEntity:       modelBase.TitledEntity{Title: "Test Document"},
		Timestamps:         modelBase.Timestamps{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Type:               "chapter",
		Status:             "draft",
	}
	mockDocRepo.On("GetByID", ctx, docID).Return(mockDoc, nil)

	// Mock version conflict
	mockContentRepo.On("UpdateWithVersion", ctx, docID, "New content", 1).Return(errors.New("version conflict"))

	// Request
	req := &documentService.AutoSaveRequest{
		DocumentID:   docID,
		Content:      "New content",
		WordCount:    100,
		ExpectedVersion: 1,
	}

	// Execute
	resp, err := service.AutoSaveDocument(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "版本冲突")
}
