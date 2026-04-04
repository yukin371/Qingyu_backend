package document

import (
	"context"
	"testing"
	"time"

	servicemock "Qingyu_backend/service/mock"

	"Qingyu_backend/models/dto"
	"Qingyu_backend/models/writer"
	modelbase "Qingyu_backend/models/writer/base"
	writerRepo "Qingyu_backend/repository/interfaces/writer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDocumentService_AutoSaveDocument_HandlesOptimisticConflictBySentinel(t *testing.T) {
	docRepo := new(MockDocumentRepository)
	projectRepo := new(MockProjectRepository)
	contentRepo := new(servicemock.MockDocumentContentRepository)
	svc := NewDocumentService(docRepo, contentRepo, projectRepo, nil)

	userObjID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	projectOID := primitive.NewObjectID()
	docID := primitive.NewObjectID().Hex()

	project := &writer.Project{OwnedEntity: modelbase.OwnedEntity{AuthorID: userObjID}}
	doc := &writer.Document{ProjectID: projectOID}

	ctx := context.WithValue(context.Background(), "userId", userObjID.Hex())

	docRepo.On("GetByID", mock.Anything, docID).Return(doc, nil).Once()
	projectRepo.On("GetByID", mock.Anything, projectOID.Hex()).Return(project, nil).Once()
	contentRepo.On("GetByDocumentID", mock.Anything, docID).Return(&writer.DocumentContent{Version: 2}, nil).Once()
	contentRepo.On(
		"UpdateWithVersion",
		mock.Anything,
		docID,
		mock.MatchedBy(func(updates map[string]interface{}) bool {
			return updates["content"] == "abc" && updates["word_count"] == 3 && updates["char_count"] == 3
		}),
		2,
	).Return(writerRepo.ErrOptimisticLockConflict).Once()

	resp, err := svc.AutoSaveDocument(ctx, &dto.AutoSaveRequest{
		DocumentID:     docID,
		Content:        "abc",
		CurrentVersion: 2,
		SaveType:       "auto",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Saved)
	assert.True(t, resp.HasConflict)

	docRepo.AssertExpectations(t)
	projectRepo.AssertExpectations(t)
	contentRepo.AssertExpectations(t)
}

func TestDocumentService_UpdateDocumentContent_MapsSentinelConflictToBusinessError(t *testing.T) {
	docRepo := new(MockDocumentRepository)
	projectRepo := new(MockProjectRepository)
	contentRepo := new(servicemock.MockDocumentContentRepository)
	svc := NewDocumentService(docRepo, contentRepo, projectRepo, nil)

	userObjID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	projectOID := primitive.NewObjectID()
	docID := primitive.NewObjectID().Hex()

	project := &writer.Project{OwnedEntity: modelbase.OwnedEntity{AuthorID: userObjID}}
	doc := &writer.Document{ProjectID: projectOID}

	ctx := context.WithValue(context.Background(), "userId", userObjID.Hex())

	docRepo.On("GetByID", mock.Anything, docID).Return(doc, nil).Once()
	projectRepo.On("GetByID", mock.Anything, projectOID.Hex()).Return(project, nil).Once()
	contentRepo.On("GetByDocumentID", mock.Anything, docID).Return(&writer.DocumentContent{Version: 5}, nil).Once()
	contentRepo.On(
		"UpdateWithVersion",
		mock.Anything,
		docID,
		mock.AnythingOfType("map[string]interface {}"),
		5,
	).Return(writerRepo.ErrOptimisticLockConflict).Once()

	err := svc.UpdateDocumentContent(ctx, &dto.UpdateContentRequest{
		DocumentID: docID,
		Content:    "new-content",
		Version:    5,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "版本冲突")

	docRepo.AssertExpectations(t)
	projectRepo.AssertExpectations(t)
	contentRepo.AssertExpectations(t)
}

func TestDocumentService_ReplaceDocumentContents_PreparesDataBeforeRepositoryCreate(t *testing.T) {
	docRepo := new(MockDocumentRepository)
	projectRepo := new(MockProjectRepository)
	contentRepo := new(servicemock.MockDocumentContentRepository)
	svc := NewDocumentService(docRepo, contentRepo, projectRepo, nil)

	userObjID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	projectOID := primitive.NewObjectID()
	docID := primitive.NewObjectID().Hex()

	project := &writer.Project{OwnedEntity: modelbase.OwnedEntity{AuthorID: userObjID}}
	doc := &writer.Document{ProjectID: projectOID}

	ctx := context.WithValue(context.Background(), "userId", userObjID.Hex())

	docRepo.On("GetByID", mock.Anything, docID).Return(doc, nil).Once()
	projectRepo.On("GetByID", mock.Anything, projectOID.Hex()).Return(project, nil).Once()
	contentRepo.On("List", mock.Anything, mock.Anything).Return([]*writer.DocumentContent{}, nil).Once()
	contentRepo.On(
		"Create",
		mock.Anything,
		mock.MatchedBy(func(c *writer.DocumentContent) bool {
			if c == nil {
				return false
			}
			return !c.ID.IsZero() &&
				!c.CreatedAt.IsZero() &&
				!c.UpdatedAt.IsZero() &&
				!c.LastSavedAt.IsZero() &&
				c.Version >= 1 &&
				c.ContentType == "tiptap" &&
				c.DocumentID.Hex() == docID
		}),
	).Return(nil).Twice()
	docRepo.On("Update", mock.Anything, docID, mock.MatchedBy(func(updates map[string]interface{}) bool {
		return updates["word_count"] == 5
	})).Return(nil).Once()

	resp, err := svc.ReplaceDocumentContents(ctx, &dto.ReplaceDocumentContentsRequest{
		DocumentID: docID,
		Contents: []dto.ParagraphContent{
			{Order: 1, Content: "ab", ContentType: "tiptap"},
			{Order: 2, Content: "cde", ContentType: "tiptap"},
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, resp.Total)
	assert.Equal(t, 5, resp.WordCount)
	assert.WithinDuration(t, time.Now(), resp.UpdatedAt, 2*time.Second)

	docRepo.AssertExpectations(t)
	projectRepo.AssertExpectations(t)
	contentRepo.AssertExpectations(t)
}

func TestDocumentService_ReplaceDocumentContents_RejectsEmptyParagraph(t *testing.T) {
	docRepo := new(MockDocumentRepository)
	projectRepo := new(MockProjectRepository)
	contentRepo := new(servicemock.MockDocumentContentRepository)
	svc := NewDocumentService(docRepo, contentRepo, projectRepo, nil)

	_, err := svc.ReplaceDocumentContents(context.Background(), &dto.ReplaceDocumentContentsRequest{
		DocumentID: primitive.NewObjectID().Hex(),
		Contents: []dto.ParagraphContent{
			{Order: 1, Content: "   ", ContentType: "tiptap"},
		},
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "内容不能为空")
	docRepo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
	contentRepo.AssertNotCalled(t, "List", mock.Anything, mock.Anything)
}

func TestDocumentService_ReplaceDocumentContents_RejectsDuplicateParagraphID(t *testing.T) {
	docRepo := new(MockDocumentRepository)
	projectRepo := new(MockProjectRepository)
	contentRepo := new(servicemock.MockDocumentContentRepository)
	svc := NewDocumentService(docRepo, contentRepo, projectRepo, nil)

	pid := primitive.NewObjectID().Hex()
	_, err := svc.ReplaceDocumentContents(context.Background(), &dto.ReplaceDocumentContentsRequest{
		DocumentID: primitive.NewObjectID().Hex(),
		Contents: []dto.ParagraphContent{
			{ParagraphID: pid, Order: 1, Content: "a", ContentType: "tiptap"},
			{ParagraphID: pid, Order: 2, Content: "b", ContentType: "tiptap"},
		},
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "重复段落ID")
	docRepo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
	contentRepo.AssertNotCalled(t, "List", mock.Anything, mock.Anything)
}

func TestDocumentService_ReplaceDocumentContents_RejectsDuplicateOrder(t *testing.T) {
	docRepo := new(MockDocumentRepository)
	projectRepo := new(MockProjectRepository)
	contentRepo := new(servicemock.MockDocumentContentRepository)
	svc := NewDocumentService(docRepo, contentRepo, projectRepo, nil)

	_, err := svc.ReplaceDocumentContents(context.Background(), &dto.ReplaceDocumentContentsRequest{
		DocumentID: primitive.NewObjectID().Hex(),
		Contents: []dto.ParagraphContent{
			{Order: 1, Content: "a", ContentType: "tiptap"},
			{Order: 1, Content: "b", ContentType: "tiptap"},
		},
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "重复段落顺序")
	docRepo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
	contentRepo.AssertNotCalled(t, "List", mock.Anything, mock.Anything)
}
