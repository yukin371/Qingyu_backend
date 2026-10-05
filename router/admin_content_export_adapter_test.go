package router

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"Qingyu_backend/models/bookstore"
	infrastructure "Qingyu_backend/repository/interfaces/infrastructure"
	serviceinterfaces "Qingyu_backend/service/interfaces"
)

type adminContentBookRepositoryStub struct {
	filter infrastructure.Filter
}

func (s *adminContentBookRepositoryStub) List(_ context.Context, filter infrastructure.Filter) ([]*bookstore.Book, error) {
	s.filter = filter
	return nil, nil
}

func TestAdminContentInfrastructureFilterMapsPersistentDateFields(t *testing.T) {
	start := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 10, 2, 23, 59, 59, 0, time.UTC)
	repo := &adminContentBookRepositoryStub{}
	port := newAdminContentBookPort(repo)

	_, err := port.List(context.Background(), serviceinterfaces.AdminContentFilter{
		Status: "ongoing",
		CreatedAt: &serviceinterfaces.AdminContentTimeRange{
			Start: &start,
			End:   &end,
		},
	})

	require.NoError(t, err)
	require.Equal(t, "ongoing", repo.filter.GetConditions()["status"])
	require.Equal(t, start, repo.filter.GetConditions()["created_at.$gte"])
	require.Equal(t, end, repo.filter.GetConditions()["created_at.$lte"])
}

type adminContentChapterRepositoryStub struct {
	filter infrastructure.Filter
}

func (s *adminContentChapterRepositoryStub) List(_ context.Context, filter infrastructure.Filter) ([]*bookstore.Chapter, error) {
	s.filter = filter
	return nil, nil
}

func TestAdminContentInfrastructureFilterMapsChapterPublishFields(t *testing.T) {
	start := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 10, 2, 23, 59, 59, 0, time.UTC)
	repo := &adminContentChapterRepositoryStub{}
	port := newAdminContentChapterPort(repo)

	_, err := port.List(context.Background(), serviceinterfaces.AdminContentFilter{
		BookID: "book-1",
		PublishedAt: &serviceinterfaces.AdminContentTimeRange{
			Start: &start,
			End:   &end,
		},
	})

	require.NoError(t, err)
	require.Equal(t, "book-1", repo.filter.GetConditions()["book_id"])
	require.Equal(t, start, repo.filter.GetConditions()["publish_time.$gte"])
	require.Equal(t, end, repo.filter.GetConditions()["publish_time.$lte"])
}
