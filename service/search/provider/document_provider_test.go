package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/models/search"
	searchengine "Qingyu_backend/service/search/engine"
)

// MockEngineForDocument 是用于测试 DocumentProvider 的模拟搜索引擎
type MockEngineForDocument struct {
	shouldFail bool
	emptyResult bool
}

func (m *MockEngineForDocument) Search(ctx context.Context, index string, query interface{}, opts *searchengine.SearchOptions) (*searchengine.SearchResult, error) {
	if m.shouldFail {
		return nil, fmt.Errorf("mock search engine error")
	}

	if m.emptyResult {
		return &searchengine.SearchResult{
			Total: 0,
			Hits:  []searchengine.Hit{},
		}, nil
	}

	// 模拟搜索结果
	return &searchengine.SearchResult{
		Total: 2,
		Hits: []searchengine.Hit{
			{
				ID:    "507f1f77bcf86cd799439011",
				Score: 1.0,
				Source: map[string]interface{}{
					"title":       "第一章 开始",
					"content":     "这是第一章的内容",
					"project_id":  "507f1f77bcf86cd799439011",
					"user_id":     "user123",
					"status":      "writing",
					"type":        "chapter",
					"word_count":  1000,
					"level":       1,
					"order":       1,
					"created_at":  "2024-01-01T00:00:00Z",
					"updated_at":  "2024-01-02T00:00:00Z",
				},
			},
			{
				ID:    "507f1f77bcf86cd799439012",
				Score: 0.8,
				Source: map[string]interface{}{
					"title":       "第二章 继续",
					"content":     "这是第二章的内容",
					"project_id":  "507f1f77bcf86cd799439011",
					"user_id":     "user123",
					"status":      "completed",
					"type":        "chapter",
					"word_count":  2000,
					"level":       1,
					"order":       2,
					"created_at":  "2024-01-01T00:00:00Z",
					"updated_at":  "2024-01-03T00:00:00Z",
				},
			},
		},
	}, nil
}

func (m *MockEngineForDocument) Index(ctx context.Context, index string, documents []searchengine.Document) error {
	return nil
}

func (m *MockEngineForDocument) Update(ctx context.Context, index string, id string, document searchengine.Document) error {
	return nil
}

func (m *MockEngineForDocument) Delete(ctx context.Context, index string, id string) error {
	return nil
}

func (m *MockEngineForDocument) CreateIndex(ctx context.Context, index string, mapping interface{}) error {
	return nil
}

func (m *MockEngineForDocument) Health(ctx context.Context) error {
	return nil
}

func TestNewDocumentProvider(t *testing.T) {
	engine := &MockEngineForDocument{}
	config := &DocumentProviderConfig{}

	provider, err := NewDocumentProvider(engine, config)

	assert.NoError(t, err)
	assert.NotNil(t, provider)
	assert.Equal(t, search.SearchTypeDocuments, provider.Type())
	assert.NotNil(t, provider.logger)
}

func TestNewDocumentProvider_WithNilEngine(t *testing.T) {
	config := &DocumentProviderConfig{}

	provider, err := NewDocumentProvider(nil, config)

	assert.Error(t, err)
	assert.Nil(t, provider)
	assert.Contains(t, err.Error(), "engine cannot be nil")
}

func TestNewDocumentProvider_WithNilConfig(t *testing.T) {
	engine := &MockEngineForDocument{}

	provider, err := NewDocumentProvider(engine, nil)

	assert.NoError(t, err)
	assert.NotNil(t, provider)
	assert.NotNil(t, provider.config)
}

func TestDocumentProvider_Type(t *testing.T) {
	engine := &MockEngineForDocument{}
	config := &DocumentProviderConfig{}
	provider, err := NewDocumentProvider(engine, config)

	require.NoError(t, err)
	assert.Equal(t, search.SearchTypeDocuments, provider.Type())
}

func TestDocumentProvider_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *search.SearchRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
			errMsg:  "request cannot be nil",
		},
		{
			name: "empty query",
			req: &search.SearchRequest{
				Query: "",
			},
			wantErr: true,
			errMsg:  "query cannot be empty",
		},
		{
			name: "valid request",
			req: &search.SearchRequest{
				Query:    "test",
				Page:     1,
				PageSize: 20,
			},
			wantErr: false,
		},
		{
			name: "invalid page",
			req: &search.SearchRequest{
				Query: "test",
				Page:  0,
			},
			wantErr: true,
			errMsg:  "page must be greater than 0",
		},
		{
			name: "invalid page size",
			req: &search.SearchRequest{
				Query:    "test",
				Page:     1,
				PageSize: 0,
			},
			wantErr: true,
			errMsg:  "page_size must be greater than 0",
		},
		{
			name: "page size exceeds maximum",
			req: &search.SearchRequest{
				Query:    "test",
				Page:     1,
				PageSize: 101,
			},
			wantErr: true,
			errMsg:  "page_size cannot exceed",
		},
	}

	engine := &MockEngineForDocument{}
	config := &DocumentProviderConfig{}
	provider, err := NewDocumentProvider(engine, config)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := provider.Validate(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDocumentProvider_Search(t *testing.T) {
	tests := []struct {
		name    string
		req     *search.SearchRequest
		engine  *MockEngineForDocument
		wantErr bool
		check   func(t *testing.T, resp *search.SearchResponse)
	}{
		{
			name: "successful search with user_id",
			req: &search.SearchRequest{
				Query: "第一章",
				Filter: map[string]interface{}{
					"user_id": "user123",
				},
				Page:     1,
				PageSize: 20,
			},
			engine:  &MockEngineForDocument{},
			wantErr: false,
			check: func(t *testing.T, resp *search.SearchResponse) {
				assert.True(t, resp.Success)
				assert.NotNil(t, resp.Data)
				assert.Equal(t, search.SearchTypeDocuments, resp.Data.Type)
				assert.Equal(t, int64(2), resp.Data.Total)
				assert.Len(t, resp.Data.Results, 2)
				assert.Equal(t, "第一章 开始", resp.Data.Results[0].Data["title"])
			},
		},
		{
			name: "search without user_id should fail",
			req: &search.SearchRequest{
				Query:    "第一章",
				Page:     1,
				PageSize: 20,
			},
			engine:  &MockEngineForDocument{},
			wantErr: false,
			check: func(t *testing.T, resp *search.SearchResponse) {
				assert.False(t, resp.Success)
				assert.NotNil(t, resp.Error)
				assert.Equal(t, search.ErrCodeUnauthorized, resp.Error.Code)
				assert.Contains(t, resp.Error.Message, "user_id")
			},
		},
		{
			name: "search with project_id filter",
			req: &search.SearchRequest{
				Query: "第一章",
				Filter: map[string]interface{}{
					"user_id":    "user123",
					"project_id": "507f1f77bcf86cd799439011",
				},
				Page:     1,
				PageSize: 20,
			},
			engine:  &MockEngineForDocument{},
			wantErr: false,
			check: func(t *testing.T, resp *search.SearchResponse) {
				assert.True(t, resp.Success)
				assert.NotNil(t, resp.Data)
			},
		},
		{
			name: "search with status filter",
			req: &search.SearchRequest{
				Query: "第一章",
				Filter: map[string]interface{}{
					"user_id": "user123",
					"status":  "writing",
				},
				Page:     1,
				PageSize: 20,
			},
			engine:  &MockEngineForDocument{},
			wantErr: false,
			check: func(t *testing.T, resp *search.SearchResponse) {
				assert.True(t, resp.Success)
				assert.NotNil(t, resp.Data)
			},
		},
		{
			name: "search with type filter",
			req: &search.SearchRequest{
				Query: "第一章",
				Filter: map[string]interface{}{
					"user_id": "user123",
					"type":    "chapter",
				},
				Page:     1,
				PageSize: 20,
			},
			engine:  &MockEngineForDocument{},
			wantErr: false,
			check: func(t *testing.T, resp *search.SearchResponse) {
				assert.True(t, resp.Success)
				assert.NotNil(t, resp.Data)
			},
		},
		{
			name: "search with empty results",
			req: &search.SearchRequest{
				Query: "不存在的文档",
				Filter: map[string]interface{}{
					"user_id": "user123",
				},
				Page:     1,
				PageSize: 20,
			},
			engine:  &MockEngineForDocument{emptyResult: true},
			wantErr: false,
			check: func(t *testing.T, resp *search.SearchResponse) {
				assert.True(t, resp.Success)
				assert.NotNil(t, resp.Data)
				assert.Equal(t, int64(0), resp.Data.Total)
				assert.Len(t, resp.Data.Results, 0)
			},
		},
		{
			name: "search engine failure",
			req: &search.SearchRequest{
				Query: "第一章",
				Filter: map[string]interface{}{
					"user_id": "user123",
				},
				Page:     1,
				PageSize: 20,
			},
			engine:  &MockEngineForDocument{shouldFail: true},
			wantErr: false,
			check: func(t *testing.T, resp *search.SearchResponse) {
				assert.False(t, resp.Success)
				assert.NotNil(t, resp.Error)
				assert.Equal(t, search.ErrCodeEngineFailure, resp.Error.Code)
			},
		},
		{
			name: "invalid request",
			req: &search.SearchRequest{
				Query: "",
			},
			engine:  &MockEngineForDocument{},
			wantErr: false,
			check: func(t *testing.T, resp *search.SearchResponse) {
				assert.False(t, resp.Success)
				assert.NotNil(t, resp.Error)
				assert.Equal(t, search.ErrCodeInvalidRequest, resp.Error.Code)
			},
		},
	}

	config := &DocumentProviderConfig{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := NewDocumentProvider(tt.engine, config)
			require.NoError(t, err)

			ctx := context.Background()
			resp, err := provider.Search(ctx, tt.req)

			assert.NoError(t, err)
			assert.NotNil(t, resp)

			if tt.check != nil {
				tt.check(t, resp)
			}
		})
	}
}

func TestDocumentProvider_buildUserFilters(t *testing.T) {
	tests := []struct {
		name   string
		filter map[string]interface{}
		want   int // 期望的过滤条件数量
	}{
		{
			name:   "no filters",
			filter: map[string]interface{}{},
			want:   0,
		},
		{
			name: "project_id filter",
			filter: map[string]interface{}{
				"project_id": "507f1f77bcf86cd799439011",
			},
			want: 1,
		},
		{
			name: "status filter",
			filter: map[string]interface{}{
				"status": "writing",
			},
			want: 1,
		},
		{
			name: "type filter",
			filter: map[string]interface{}{
				"type": "chapter",
			},
			want: 1,
		},
		{
			name: "level filter",
			filter: map[string]interface{}{
				"level": 1,
			},
			want: 1,
		},
		{
			name: "multiple filters",
			filter: map[string]interface{}{
				"project_id": "507f1f77bcf86cd799439011",
				"status":     "writing",
				"type":       "chapter",
			},
			want: 3,
		},
	}

	engine := &MockEngineForDocument{}
	config := &DocumentProviderConfig{}
	provider, err := NewDocumentProvider(engine, config)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filters := provider.buildUserFilters(tt.filter)
			assert.Equal(t, tt.want, len(filters))
		})
	}
}

func TestDocumentProvider_validateSortField(t *testing.T) {
	engine := &MockEngineForDocument{}
	config := &DocumentProviderConfig{}
	provider, err := NewDocumentProvider(engine, config)
	require.NoError(t, err)

	tests := []struct {
		name     string
		field    string
		expected string
	}{
		{"valid created_at", "created_at", "created_at"},
		{"valid updated_at", "updated_at", "updated_at"},
		{"valid title", "title", "title"},
		{"valid word_count", "word_count", "word_count"},
		{"valid level", "level", "level"},
		{"valid order", "order", "order"},
		{"invalid field", "invalid_field", ""},
		{"empty field", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.validateSortField(tt.field)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDocumentProvider_getPageSize(t *testing.T) {
	engine := &MockEngineForDocument{}
	config := &DocumentProviderConfig{}
	provider, err := NewDocumentProvider(engine, config)
	require.NoError(t, err)

	tests := []struct {
		name     string
		pageSize int
		expected int
	}{
		{"zero size", 0, 20},
		{"negative size", -1, 20},
		{"normal size", 10, 10},
		{"max size", 100, 100},
		{"exceed max", 101, 100},
		{"large size", 1000, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.getPageSize(tt.pageSize)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDocumentProvider_getPage(t *testing.T) {
	engine := &MockEngineForDocument{}
	config := &DocumentProviderConfig{}
	provider, err := NewDocumentProvider(engine, config)
	require.NoError(t, err)

	tests := []struct {
		name     string
		page     int
		expected int
	}{
		{"zero page", 0, 1},
		{"negative page", -1, 1},
		{"normal page", 2, 2},
		{"large page", 100, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.getPage(tt.page)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDocumentProvider_calculateOffset(t *testing.T) {
	engine := &MockEngineForDocument{}
	config := &DocumentProviderConfig{}
	provider, err := NewDocumentProvider(engine, config)
	require.NoError(t, err)

	tests := []struct {
		name     string
		page     int
		pageSize int
		expected int
	}{
		{"first page", 1, 20, 0},
		{"second page", 2, 20, 20},
		{"third page", 3, 10, 20},
		{"default page size", 2, 0, 20}, // page size 0 使用默认值 20
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.calculateOffset(tt.page, tt.pageSize)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDocumentProvider_RequiresAuthentication(t *testing.T) {
	engine := &MockEngineForDocument{}
	config := &DocumentProviderConfig{}
	provider, err := NewDocumentProvider(engine, config)
	require.NoError(t, err)

	// 测试：没有 user_id 的请求应该返回认证错误
	req := &search.SearchRequest{
		Query:    "test",
		Page:     1,
		PageSize: 20,
	}

	ctx := context.Background()
	resp, err := provider.Search(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Success)
	assert.NotNil(t, resp.Error)
	assert.Equal(t, search.ErrCodeUnauthorized, resp.Error.Code)
	assert.Contains(t, resp.Error.Message, "user_id")
}

func TestDocumentProvider_ExtractUserID(t *testing.T) {
	engine := &MockEngineForDocument{}
	config := &DocumentProviderConfig{}
	provider, err := NewDocumentProvider(engine, config)
	require.NoError(t, err)

	tests := []struct {
		name     string
		filter   map[string]interface{}
		wantOK   bool
		wantID   string
	}{
		{
			name:     "no filter",
			filter:   nil,
			wantOK:   false,
			wantID:   "",
		},
		{
			name: "empty filter",
			filter:   map[string]interface{}{},
			wantOK:   false,
			wantID:   "",
		},
		{
			name: "valid user_id",
			filter: map[string]interface{}{
				"user_id": "user123",
			},
			wantOK: true,
			wantID: "user123",
		},
		{
			name: "user_id with other filters",
			filter: map[string]interface{}{
				"user_id":    "user456",
				"project_id": "proj123",
			},
			wantOK: true,
			wantID: "user456",
		},
		{
			name: "user_id not string",
			filter: map[string]interface{}{
				"user_id": 123,
			},
			wantOK: false,
			wantID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &search.SearchRequest{
				Query:  "test",
				Filter: tt.filter,
			}
			userID, ok := provider.extractUserID(req)
			assert.Equal(t, tt.wantOK, ok)
			assert.Equal(t, tt.wantID, userID)
		})
	}
}
