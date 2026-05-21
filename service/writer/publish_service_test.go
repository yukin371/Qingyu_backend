package writer

import (
	"context"
	"errors"
	"testing"
	"time"

	"Qingyu_backend/models/dto"
	writerModel "Qingyu_backend/models/writer"
	serviceInterfaces "Qingyu_backend/service/interfaces"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type stubPublicationRepository struct {
	createFn                   func(ctx context.Context, record *serviceInterfaces.PublicationRecord) error
	findByIDFn                 func(ctx context.Context, id string) (*serviceInterfaces.PublicationRecord, error)
	findByProjectIDFn          func(ctx context.Context, projectID string, page, pageSize int) ([]*serviceInterfaces.PublicationRecord, int64, error)
	findPendingFn              func(ctx context.Context, page, pageSize int) ([]*serviceInterfaces.PublicationRecord, int64, error)
	findByResourceIDFn         func(ctx context.Context, resourceID string) (*serviceInterfaces.PublicationRecord, error)
	findPublishedByProjectIDFn func(ctx context.Context, projectID string) (*serviceInterfaces.PublicationRecord, error)
	updateFn                   func(ctx context.Context, record *serviceInterfaces.PublicationRecord) error
	records                    map[string]*serviceInterfaces.PublicationRecord
}

func (s *stubPublicationRepository) Create(ctx context.Context, record *serviceInterfaces.PublicationRecord) error {
	if s.createFn != nil {
		return s.createFn(ctx, record)
	}
	if s.records != nil {
		s.records[record.ID] = record
	}
	return nil
}
func (s *stubPublicationRepository) FindByID(ctx context.Context, id string) (*serviceInterfaces.PublicationRecord, error) {
	if s.findByIDFn != nil {
		return s.findByIDFn(ctx, id)
	}
	if s.records != nil {
		record, ok := s.records[id]
		if ok {
			return record, nil
		}
	}
	return nil, errors.New("not implemented")
}
func (s *stubPublicationRepository) FindByProjectID(ctx context.Context, projectID string, page, pageSize int) ([]*serviceInterfaces.PublicationRecord, int64, error) {
	if s.findByProjectIDFn != nil {
		return s.findByProjectIDFn(ctx, projectID, page, pageSize)
	}
	return nil, 0, nil
}
func (s *stubPublicationRepository) FindPending(ctx context.Context, page, pageSize int) ([]*serviceInterfaces.PublicationRecord, int64, error) {
	if s.findPendingFn != nil {
		return s.findPendingFn(ctx, page, pageSize)
	}
	return nil, 0, nil
}
func (s *stubPublicationRepository) FindByResourceID(ctx context.Context, resourceID string) (*serviceInterfaces.PublicationRecord, error) {
	if s.findByResourceIDFn != nil {
		return s.findByResourceIDFn(ctx, resourceID)
	}
	if s.records != nil {
		for _, record := range s.records {
			if record.ResourceID == resourceID {
				return record, nil
			}
		}
	}
	return nil, errors.New("not implemented")
}
func (s *stubPublicationRepository) Update(ctx context.Context, record *serviceInterfaces.PublicationRecord) error {
	if s.updateFn != nil {
		return s.updateFn(ctx, record)
	}
	if s.records != nil {
		s.records[record.ID] = record
	}
	return nil
}
func (s *stubPublicationRepository) Delete(ctx context.Context, id string) error { return nil }
func (s *stubPublicationRepository) FindPublishedByProjectID(ctx context.Context, projectID string) (*serviceInterfaces.PublicationRecord, error) {
	if s.findPublishedByProjectIDFn != nil {
		return s.findPublishedByProjectIDFn(ctx, projectID)
	}
	if s.records != nil {
		for _, record := range s.records {
			if record.Type == "project" && record.ResourceID == projectID && record.Status == serviceInterfaces.PublicationStatusPublished {
				return record, nil
			}
		}
	}
	return nil, errors.New("not implemented")
}

type stubProjectRepository struct {
	findByIDFn func(ctx context.Context, id string) (*writerModel.Project, error)
}

func (s *stubProjectRepository) FindByID(ctx context.Context, id string) (*writerModel.Project, error) {
	if s.findByIDFn != nil {
		return s.findByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

type stubDocumentRepository struct {
	findByIDFn        func(ctx context.Context, id string) (*writerModel.Document, error)
	findByProjectIDFn func(ctx context.Context, projectID string) ([]*writerModel.Document, error)
}

func (s *stubDocumentRepository) FindByID(ctx context.Context, id string) (*writerModel.Document, error) {
	if s.findByIDFn != nil {
		return s.findByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (s *stubDocumentRepository) FindByProjectID(ctx context.Context, projectID string) ([]*writerModel.Document, error) {
	if s.findByProjectIDFn != nil {
		return s.findByProjectIDFn(ctx, projectID)
	}
	return nil, nil
}

type stubEventBus struct {
	mock.Mock
}

func (s *stubEventBus) PublishAsync(ctx context.Context, event interface{}) error {
	args := s.Called(ctx, event)
	return args.Error(0)
}

func TestPublishEventFailureIsRecordedOnProjectPublish(t *testing.T) {
	bookstoreClient := new(MockBookstoreClient)
	eventBus := new(stubEventBus)
	updateCalls := 0
	publicationRepo := &stubPublicationRepository{
		updateFn: func(ctx context.Context, record *serviceInterfaces.PublicationRecord) error {
			updateCalls++
			if updateCalls == 1 {
				assert.Equal(t, serviceInterfaces.PublicationStatusPublished, record.Status)
				assert.Empty(t, record.ReviewNote)
			}
			if updateCalls == 2 {
				assert.Contains(t, record.ReviewNote, "event dispatch failed for project.published")
			}
			return nil
		},
	}

	service := NewPublishService(nil, nil, publicationRepo, bookstoreClient, eventBus).(*PublishService)
	record := &serviceInterfaces.PublicationRecord{
		ID:         "record-1",
		Type:       "project",
		ResourceID: "project-1",
		Status:     serviceInterfaces.PublicationStatusPending,
	}
	project := createTestProject("author-1", "测试项目")
	req := &serviceInterfaces.PublishProjectRequest{
		BookstoreID: "local",
		CategoryID:  "cat-1",
	}

	bookstoreClient.On("PublishProject", mock.Anything, mock.AnythingOfType("*writer.BookstorePublishProjectRequest")).
		Return(&BookstorePublishResponse{
			Success:       true,
			BookstoreID:   "local",
			BookstoreName: "Local",
			ExternalID:    "book-1",
		}, nil).Once()
	eventBus.On("PublishAsync", mock.Anything, mock.Anything).Return(errors.New("mock event failure")).Once()

	service.executeProjectPublish(context.Background(), record, project, req)

	assert.Equal(t, 2, updateCalls)
	assert.Contains(t, record.ReviewNote, "event dispatch failed for project.published")
	bookstoreClient.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

func TestPublishEventFailureIsRecordedOnDocumentUnpublish(t *testing.T) {
	bookstoreClient := new(MockBookstoreClient)
	eventBus := new(stubEventBus)
	updateCalls := 0
	publicationRepo := &stubPublicationRepository{
		updateFn: func(ctx context.Context, record *serviceInterfaces.PublicationRecord) error {
			updateCalls++
			if updateCalls == 1 {
				assert.Equal(t, serviceInterfaces.PublicationStatusUnpublished, record.Status)
				assert.NotNil(t, record.UnpublishTime)
				assert.Empty(t, record.ReviewNote)
			}
			if updateCalls == 2 {
				assert.Contains(t, record.ReviewNote, "event dispatch failed for document.unpublished")
			}
			return nil
		},
	}

	service := NewPublishService(nil, nil, publicationRepo, bookstoreClient, eventBus).(*PublishService)
	record := &serviceInterfaces.PublicationRecord{
		ID:          "record-1",
		Type:        "document",
		ResourceID:  "document-1",
		BookstoreID: "local",
		Status:      serviceInterfaces.PublicationStatusPublished,
	}

	bookstoreClient.On("UnpublishChapter", mock.Anything, "document-1", "local").Return(nil).Once()
	eventBus.On("PublishAsync", mock.Anything, mock.Anything).Return(errors.New("mock event failure")).Once()

	err := service.unpublishDocument(context.Background(), record, "cleanup")
	assert.NoError(t, err)
	assert.Equal(t, 2, updateCalls)
	assert.Contains(t, record.ReviewNote, "event dispatch failed for document.unpublished")
	bookstoreClient.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

func TestPublishEventFailureNoteIsNotDuplicated(t *testing.T) {
	eventBus := new(stubEventBus)
	publicationRepo := &stubPublicationRepository{}
	service := NewPublishService(nil, nil, publicationRepo, nil, eventBus).(*PublishService)
	record := &serviceInterfaces.PublicationRecord{
		ID:         "record-1",
		Status:     serviceInterfaces.PublicationStatusPublished,
		ReviewNote: "event dispatch failed for project.published: mock event failure",
		UpdatedAt:  time.Now(),
	}

	eventBus.On("PublishAsync", mock.Anything, mock.Anything).Return(errors.New("mock event failure")).Once()
	service.publishEventWithAudit(context.Background(), record, map[string]interface{}{"eventType": "project.published"}, "project.published")

	assert.Equal(t, "event dispatch failed for project.published: mock event failure", record.ReviewNote)
	eventBus.AssertExpectations(t)
}

func TestReviewPublicationRejectsPendingRecord(t *testing.T) {
	record := &serviceInterfaces.PublicationRecord{
		ID:         "record-1",
		Type:       "project",
		ResourceID: "project-1",
		Status:     serviceInterfaces.PublicationStatusPending,
	}
	publicationRepo := &stubPublicationRepository{
		records: map[string]*serviceInterfaces.PublicationRecord{
			record.ID: record,
		},
	}
	service := NewPublishService(nil, nil, publicationRepo, nil, nil).(*PublishService)

	updated, err := service.ReviewPublication(context.Background(), record.ID, "reviewer-1", false, "资料不完整")

	assert.NoError(t, err)
	assert.Equal(t, serviceInterfaces.PublicationStatusRejected, updated.Status)
	assert.Equal(t, "reviewer-1", updated.ReviewedBy)
	assert.Equal(t, "资料不完整", updated.ReviewNote)
	assert.NotNil(t, updated.ReviewedAt)
}

func TestReviewPublicationApprovesDocumentAndPublishes(t *testing.T) {
	authorID := "507f1f77bcf86cd799439011"
	project := createTestProject(authorID, "测试项目")
	projectID := project.ID.Hex()
	document := createTestDocument(projectID, "第1章")
	record := &serviceInterfaces.PublicationRecord{
		ID:          "record-1",
		Type:        "document",
		ResourceID:  document.ID.Hex(),
		BookstoreID: "store-1",
		Status:      serviceInterfaces.PublicationStatusPending,
		Metadata: dto.PublicationMetadata{
			ChapterTitle:  "第1章",
			ChapterNumber: 1,
			IsFree:        true,
			AuthorNote:    "首章免费",
		},
	}

	publicationRepo := &stubPublicationRepository{
		records: map[string]*serviceInterfaces.PublicationRecord{
			record.ID: record,
		},
	}
	projectRepo := &stubProjectRepository{
		findByIDFn: func(ctx context.Context, id string) (*writerModel.Project, error) {
			assert.Equal(t, projectID, id)
			return project, nil
		},
	}
	documentRepo := &stubDocumentRepository{
		findByIDFn: func(ctx context.Context, id string) (*writerModel.Document, error) {
			assert.Equal(t, document.ID.Hex(), id)
			return document, nil
		},
	}
	bookstoreClient := new(MockBookstoreClient)
	bookstoreClient.
		On("PublishChapter", mock.Anything, mock.MatchedBy(func(req *BookstorePublishChapterRequest) bool {
			return req.ProjectID == projectID &&
				req.DocumentID == document.ID.Hex() &&
				req.ChapterTitle == "第1章" &&
				req.ChapterNumber == 1 &&
				req.IsFree
		})).
		Return(&BookstorePublishResponse{
			Success:       true,
			BookstoreID:   "store-1",
			BookstoreName: "测试书城",
			ExternalID:    "chapter-1",
		}, nil).
		Once()

	service := NewPublishService(projectRepo, documentRepo, publicationRepo, bookstoreClient, nil).(*PublishService)

	updated, err := service.ReviewPublication(context.Background(), record.ID, "reviewer-1", true, "通过")

	assert.NoError(t, err)
	assert.Equal(t, serviceInterfaces.PublicationStatusPublished, updated.Status)
	assert.Equal(t, "store-1", updated.BookstoreID)
	assert.Equal(t, "测试书城", updated.BookstoreName)
	assert.Equal(t, "chapter-1", updated.ExternalID)
	assert.Equal(t, "通过", updated.ReviewNote)
	assert.Equal(t, "reviewer-1", updated.ReviewedBy)
	assert.NotNil(t, updated.PublishTime)
	bookstoreClient.AssertExpectations(t)
}

func TestGetPublicationRecordsNormalizesInvalidPagination(t *testing.T) {
	var capturedPage int
	var capturedPageSize int
	publicationRepo := &stubPublicationRepository{
		findByProjectIDFn: func(ctx context.Context, projectID string, page, pageSize int) ([]*serviceInterfaces.PublicationRecord, int64, error) {
			capturedPage = page
			capturedPageSize = pageSize
			return []*serviceInterfaces.PublicationRecord{}, 0, nil
		},
	}
	service := NewPublishService(nil, nil, publicationRepo, nil, nil).(*PublishService)

	_, _, err := service.GetPublicationRecords(context.Background(), "project-1", 0, 500)

	assert.NoError(t, err)
	assert.Equal(t, 1, capturedPage)
	assert.Equal(t, 20, capturedPageSize)
}
