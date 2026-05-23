package search

import (
	"context"
	"log"
	"testing"

	searchmodel "Qingyu_backend/models/search"
	"Qingyu_backend/service/search/provider"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type failingProvider struct{}

func (f *failingProvider) Search(context.Context, *searchmodel.SearchRequest) (*searchmodel.SearchResponse, error) {
	return &searchmodel.SearchResponse{
		Success: false,
		Error: &searchmodel.ErrorInfo{
			Code:    searchmodel.ErrCodeEngineFailure,
			Message: "provider failed",
		},
	}, nil
}

func (f *failingProvider) Type() searchmodel.SearchType { return searchmodel.SearchTypeBooks }
func (f *failingProvider) Validate(*searchmodel.SearchRequest) error {
	return nil
}
func (f *failingProvider) GetByID(context.Context, string) (*searchmodel.SearchItem, error) {
	return nil, nil
}
func (f *failingProvider) GetBatch(context.Context, []string) ([]searchmodel.SearchItem, error) {
	return nil, nil
}

var _ provider.Provider = (*failingProvider)(nil)

func TestSearchService_Search_ProviderFailureResponseDoesNotPanic(t *testing.T) {
	svc := NewSearchService(log.Default(), &Config{EnableCache: false, MaxConcurrentSearches: 1}, nil)
	svc.RegisterProvider(&failingProvider{})

	req := &searchmodel.SearchRequest{
		Type:     searchmodel.SearchTypeBooks,
		Query:    "*",
		Page:     1,
		PageSize: 20,
	}

	resp, err := svc.Search(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.False(t, resp.Success)
	require.NotNil(t, resp.Error)
	assert.Equal(t, searchmodel.ErrCodeEngineFailure, resp.Error.Code)
	assert.NotNil(t, resp.Meta)
}
