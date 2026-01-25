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

// MockEngineForProject 是用于测试 ProjectProvider 的模拟搜索引擎
type MockEngineForProject struct {
	shouldFail bool
}

func (m *MockEngineForProject) Search(ctx context.Context, index string, query interface{}, opts *searchengine.SearchOptions) (*searchengine.SearchResult, error) {
	if m.shouldFail {
		return nil, fmt.Errorf("mock search engine error")
	}

	// 模拟搜索结果
	return &searchengine.SearchResult{
		Total: 2,
		Hits: []searchengine.Hit{
			{
				ID:    "507f1f77bcf86cd799439011",
				Score: 1.0,
				Source: map[string]interface{}{
					"title":        "测试项目1",
					"summary":      "这是测试项目1的描述",
					"writing_type": "novel",
					"status":       "draft",
					"visibility":   "public",
					"category":     "玄幻",
					"author_id":    "user123",
					"statistics": map[string]interface{}{
						"total_words":    10000,
						"chapter_count":  5,
						"document_count": 3,
					},
					"created_at": "2024-01-01T00:00:00Z",
					"updated_at": "2024-01-02T00:00:00Z",
				},
			},
			{
				ID:    "507f1f77bcf86cd799439012",
				Score: 0.8,
				Source: map[string]interface{}{
					"title":        "测试项目2",
					"summary":      "这是测试项目2的描述",
					"writing_type": "article",
					"status":       "serializing",
					"visibility":   "private",
					"category":     "科幻",
					"author_id":    "user123",
					"statistics": map[string]interface{}{
						"total_words":    20000,
						"chapter_count":  10,
						"document_count": 5,
					},
					"created_at": "2024-01-01T00:00:00Z",
					"updated_at": "2024-01-03T00:00:00Z",
				},
			},
		},
	}, nil
}

func (m *MockEngineForProject) Index(ctx context.Context, index string, documents []searchengine.Document) error {
	return nil
}

func (m *MockEngineForProject) Update(ctx context.Context, index string, id string, document searchengine.Document) error {
	return nil
}

func (m *MockEngineForProject) Delete(ctx context.Context, index string, id string) error {
	return nil
}

func (m *MockEngineForProject) CreateIndex(ctx context.Context, index string, mapping interface{}) error {
	return nil
}

func (m *MockEngineForProject) Health(ctx context.Context) error {
	return nil
}

func TestNewProjectProvider(t *testing.T) {
	engine := &MockEngineForProject{}
	config := &ProjectProviderConfig{}

	provider, err := NewProjectProvider(engine, config)

	assert.NoError(t, err)
	assert.NotNil(t, provider)
	assert.Equal(t, search.SearchTypeProjects, provider.Type())
	assert.NotNil(t, provider.logger)
}

func TestNewProjectProvider_WithNilEngine(t *testing.T) {
	config := &ProjectProviderConfig{}

	provider, err := NewProjectProvider(nil, config)

	assert.Error(t, err)
	assert.Nil(t, provider)
	assert.Contains(t, err.Error(), "engine cannot be nil")
}

func TestNewProjectProvider_WithNilConfig(t *testing.T) {
	engine := &MockEngineForProject{}

	provider, err := NewProjectProvider(engine, nil)

	assert.NoError(t, err)
	assert.NotNil(t, provider)
	assert.NotNil(t, provider.config)
}

func TestProjectProvider_Type(t *testing.T) {
	engine := &MockEngineForProject{}
	config := &ProjectProviderConfig{}
	provider, err := NewProjectProvider(engine, config)

	require.NoError(t, err)
	assert.Equal(t, search.SearchTypeProjects, provider.Type())
}

func TestProjectProvider_Validate(t *testing.T) {
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

	engine := &MockEngineForProject{}
	config := &ProjectProviderConfig{}
	provider, err := NewProjectProvider(engine, config)
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

func TestProjectProvider_Search(t *testing.T) {
	tests := []struct {
		name    string
		req     *search.SearchRequest
		wantErr bool
		check   func(t *testing.T, resp *search.SearchResponse)
	}{
		{
			name: "successful search with user_id",
			req: &search.SearchRequest{
				Query: "测试项目",
				Filter: map[string]interface{}{
					"user_id": "user123",
				},
				Page:     1,
				PageSize: 20,
			},
			wantErr: false,
			check: func(t *testing.T, resp *search.SearchResponse) {
				assert.True(t, resp.Success)
				assert.NotNil(t, resp.Data)
				assert.Equal(t, search.SearchTypeProjects, resp.Data.Type)
				assert.Equal(t, int64(2), resp.Data.Total)
				assert.Len(t, resp.Data.Results, 2)
				assert.Equal(t, "测试项目1", resp.Data.Results[0].Data["title"])
			},
		},
		{
			name: "search without user_id should fail",
			req: &search.SearchRequest{
				Query:    "测试项目",
				Page:     1,
				PageSize: 20,
			},
			wantErr: false,
			check: func(t *testing.T, resp *search.SearchResponse) {
				assert.False(t, resp.Success)
				assert.NotNil(t, resp.Error)
				assert.Equal(t, search.ErrCodeUnauthorized, resp.Error.Code)
				assert.Contains(t, resp.Error.Message, "user_id")
			},
		},
		{
			name: "search with status filter",
			req: &search.SearchRequest{
				Query: "测试项目",
				Filter: map[string]interface{}{
					"user_id": "user123",
					"status":  "draft",
				},
				Page:     1,
				PageSize: 20,
			},
			wantErr: false,
			check: func(t *testing.T, resp *search.SearchResponse) {
				assert.True(t, resp.Success)
				assert.NotNil(t, resp.Data)
			},
		},
		{
			name: "search with visibility filter",
			req: &search.SearchRequest{
				Query: "测试项目",
				Filter: map[string]interface{}{
					"user_id":   "user123",
					"is_public": true,
				},
				Page:     1,
				PageSize: 20,
			},
			wantErr: false,
			check: func(t *testing.T, resp *search.SearchResponse) {
				assert.True(t, resp.Success)
				assert.NotNil(t, resp.Data)
			},
		},
		{
			name: "search with category filter",
			req: &search.SearchRequest{
				Query: "测试项目",
				Filter: map[string]interface{}{
					"user_id":  "user123",
					"category": "玄幻",
				},
				Page:     1,
				PageSize: 20,
			},
			wantErr: false,
			check: func(t *testing.T, resp *search.SearchResponse) {
				assert.True(t, resp.Success)
				assert.NotNil(t, resp.Data)
			},
		},
		{
			name: "invalid request",
			req: &search.SearchRequest{
				Query: "",
			},
			wantErr: false,
			check: func(t *testing.T, resp *search.SearchResponse) {
				assert.False(t, resp.Success)
				assert.NotNil(t, resp.Error)
				assert.Equal(t, search.ErrCodeInvalidRequest, resp.Error.Code)
			},
		},
	}

	engine := &MockEngineForProject{}
	config := &ProjectProviderConfig{}
	provider, err := NewProjectProvider(engine, config)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestProjectProvider_buildUserFilters(t *testing.T) {
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
			name: "status filter",
			filter: map[string]interface{}{
				"status": "draft",
			},
			want: 1,
		},
		{
			name: "is_public filter - true",
			filter: map[string]interface{}{
				"is_public": true,
			},
			want: 1,
		},
		{
			name: "is_public filter - false",
			filter: map[string]interface{}{
				"is_public": false,
			},
			want: 1,
		},
		{
			name: "category filter",
			filter: map[string]interface{}{
				"category": "玄幻",
			},
			want: 1,
		},
		{
			name: "multiple filters",
			filter: map[string]interface{}{
				"status":    "serializing",
				"is_public": true,
				"category":  "科幻",
			},
			want: 3,
		},
	}

	engine := &MockEngineForProject{}
	config := &ProjectProviderConfig{}
	provider, err := NewProjectProvider(engine, config)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// buildUserFilters 会处理除了 author_id/user_id 之外的过滤条件
			// author_id 过滤在 Search 方法中单独添加
			filters := provider.buildUserFilters(tt.filter)
			assert.Equal(t, tt.want, len(filters))
		})
	}
}

func TestProjectProvider_validateSortField(t *testing.T) {
	engine := &MockEngineForProject{}
	config := &ProjectProviderConfig{}
	provider, err := NewProjectProvider(engine, config)
	require.NoError(t, err)

	tests := []struct {
		name     string
		field    string
		expected string
	}{
		{"valid created_at", "created_at", "created_at"},
		{"valid updated_at", "updated_at", "updated_at"},
		{"valid published_at", "published_at", "published_at"},
		{"valid title", "title", "title"},
		{"valid total_words", "total_words", "statistics.total_words"},
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

func TestProjectProvider_getPageSize(t *testing.T) {
	engine := &MockEngineForProject{}
	config := &ProjectProviderConfig{}
	provider, err := NewProjectProvider(engine, config)
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

func TestProjectProvider_getPage(t *testing.T) {
	engine := &MockEngineForProject{}
	config := &ProjectProviderConfig{}
	provider, err := NewProjectProvider(engine, config)
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

func TestProjectProvider_calculateOffset(t *testing.T) {
	engine := &MockEngineForProject{}
	config := &ProjectProviderConfig{}
	provider, err := NewProjectProvider(engine, config)
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

func TestProjectProvider_RequiresAuthentication(t *testing.T) {
	engine := &MockEngineForProject{}
	config := &ProjectProviderConfig{}
	provider, err := NewProjectProvider(engine, config)
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
