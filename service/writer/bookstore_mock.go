package writer

import (
	"context"

	"github.com/stretchr/testify/mock"

	"Qingyu_backend/models/dto"
)

// MockBookstoreClient Mock书城客户端
type MockBookstoreClient struct {
	mock.Mock
}

func (m *MockBookstoreClient) PublishProject(ctx context.Context, req *BookstorePublishProjectRequest) (*BookstorePublishResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*BookstorePublishResponse), args.Error(1)
}

func (m *MockBookstoreClient) UnpublishProject(ctx context.Context, projectID, bookstoreID string) error {
	args := m.Called(ctx, projectID, bookstoreID)
	return args.Error(0)
}

func (m *MockBookstoreClient) PublishChapter(ctx context.Context, req *BookstorePublishChapterRequest) (*BookstorePublishResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*BookstorePublishResponse), args.Error(1)
}

func (m *MockBookstoreClient) UnpublishChapter(ctx context.Context, chapterID, bookstoreID string) error {
	args := m.Called(ctx, chapterID, bookstoreID)
	return args.Error(0)
}

func (m *MockBookstoreClient) UpdateChapter(ctx context.Context, req *BookstoreUpdateChapterRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockBookstoreClient) GetStatistics(ctx context.Context, projectID, bookstoreID string) (*dto.PublicationStatistics, error) {
	args := m.Called(ctx, projectID, bookstoreID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PublicationStatistics), args.Error(1)
}
