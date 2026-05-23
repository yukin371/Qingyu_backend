package events

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/service/search/engine"
)

type stubSearchEngine struct {
	indexFunc func(ctx context.Context, index string, documents []engine.Document) error
}

func (s *stubSearchEngine) Search(ctx context.Context, index string, query interface{}, opts *engine.SearchOptions) (*engine.SearchResult, error) {
	return nil, nil
}

func (s *stubSearchEngine) Index(ctx context.Context, index string, documents []engine.Document) error {
	if s.indexFunc != nil {
		return s.indexFunc(ctx, index, documents)
	}
	return nil
}

func (s *stubSearchEngine) Update(ctx context.Context, index string, id string, document engine.Document) error {
	return nil
}

func (s *stubSearchEngine) Delete(ctx context.Context, index string, id string) error {
	return nil
}

func (s *stubSearchEngine) CreateIndex(ctx context.Context, index string, mapping interface{}) error {
	return nil
}

func (s *stubSearchEngine) Health(ctx context.Context) error {
	return nil
}

func TestSearchIndexHandler_Handle_SkipsDuplicateKeyError(t *testing.T) {
	t.Parallel()

	handler := NewSearchIndexHandler(&stubSearchEngine{
		indexFunc: func(ctx context.Context, index string, documents []engine.Document) error {
			require.Equal(t, booksIndexName, index)
			require.Len(t, documents, 1)
			assert.Equal(t, "book-1", documents[0].ID)
			return errors.New("E11000 duplicate key error")
		},
	})

	err := handler.Handle(context.Background(), newTestEvent("project.published", map[string]interface{}{
		"bookstoreId": "book-1",
		"projectId":   "project-1",
	}))

	assert.NoError(t, err)
}

func TestSearchIndexHandler_Handle_BuildsSearchDocument(t *testing.T) {
	t.Parallel()

	publishedAt := time.Now().UTC().Truncate(time.Second)
	var captured engine.Document

	handler := NewSearchIndexHandler(&stubSearchEngine{
		indexFunc: func(ctx context.Context, index string, documents []engine.Document) error {
			require.Equal(t, booksIndexName, index)
			require.Len(t, documents, 1)
			captured = documents[0]
			return nil
		},
	})

	err := handler.Handle(context.Background(), newTestEvent("project.published", map[string]interface{}{
		"bookstoreId": "book-2",
		"projectId":   "project-2",
		"publishedAt": publishedAt.Format(time.RFC3339),
	}))

	require.NoError(t, err)
	assert.Equal(t, "book-2", captured.ID)
	assert.Equal(t, "project-2", captured.Source["project_id"])
	assert.Equal(t, "book-2", captured.Source["bookstore_id"])
	assert.Equal(t, "ongoing", captured.Source["status"])
	assert.Equal(t, publishedAt, captured.Source["published_at"])
}
