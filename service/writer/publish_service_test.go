package writer

import (
	"context"
	"errors"
	"testing"
	"time"

	serviceInterfaces "Qingyu_backend/service/interfaces"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type stubPublicationRepository struct {
	updateFn func(ctx context.Context, record *serviceInterfaces.PublicationRecord) error
}

func (s *stubPublicationRepository) Create(ctx context.Context, record *serviceInterfaces.PublicationRecord) error {
	return nil
}
func (s *stubPublicationRepository) FindByID(ctx context.Context, id string) (*serviceInterfaces.PublicationRecord, error) {
	return nil, errors.New("not implemented")
}
func (s *stubPublicationRepository) FindByProjectID(ctx context.Context, projectID string, page, pageSize int) ([]*serviceInterfaces.PublicationRecord, int64, error) {
	return nil, 0, nil
}
func (s *stubPublicationRepository) FindPending(ctx context.Context, page, pageSize int) ([]*serviceInterfaces.PublicationRecord, int64, error) {
	return nil, 0, nil
}
func (s *stubPublicationRepository) FindByResourceID(ctx context.Context, resourceID string) (*serviceInterfaces.PublicationRecord, error) {
	return nil, errors.New("not implemented")
}
func (s *stubPublicationRepository) Update(ctx context.Context, record *serviceInterfaces.PublicationRecord) error {
	if s.updateFn != nil {
		return s.updateFn(ctx, record)
	}
	return nil
}
func (s *stubPublicationRepository) Delete(ctx context.Context, id string) error { return nil }
func (s *stubPublicationRepository) FindPublishedByProjectID(ctx context.Context, projectID string) (*serviceInterfaces.PublicationRecord, error) {
	return nil, errors.New("not implemented")
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
