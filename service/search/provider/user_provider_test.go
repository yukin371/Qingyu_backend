package provider

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/models/search"
	searchengine "Qingyu_backend/service/search/engine"
)

// MockEngine 是用于测试的模拟搜索引擎
type MockEngine struct{}

func (m *MockEngine) Search(ctx context.Context, index string, query interface{}, opts *searchengine.SearchOptions) (*searchengine.SearchResult, error) {
	// 简单模拟搜索结果
	return &searchengine.SearchResult{
		Total: 2,
		Hits: []searchengine.Hit{
			{
				ID:    "user1",
				Score: 1.0,
				Source: map[string]interface{}{
					"username":       "testuser1",
					"nickname":       "测试用户1",
					"bio":            "这是测试用户1的简介",
					"avatar":         "https://example.com/avatar1.jpg",
					"roles":          []string{"reader"},
					"email_verified": true,
					"status":         "active",
				},
			},
			{
				ID:    "user2",
				Score: 0.8,
				Source: map[string]interface{}{
					"username":       "testuser2",
					"nickname":       "测试用户2",
					"bio":            "这是测试用户2的简介",
					"avatar":         "https://example.com/avatar2.jpg",
					"roles":          []string{"author"},
					"email_verified": false,
					"status":         "active",
				},
			},
		},
	}, nil
}

func (m *MockEngine) Index(ctx context.Context, index string, documents []searchengine.Document) error {
	return nil
}

func (m *MockEngine) Update(ctx context.Context, index string, id string, document searchengine.Document) error {
	return nil
}

func (m *MockEngine) Delete(ctx context.Context, index string, id string) error {
	return nil
}

func (m *MockEngine) CreateIndex(ctx context.Context, index string, mapping interface{}) error {
	return nil
}

func (m *MockEngine) Health(ctx context.Context) error {
	return nil
}

func TestNewUserProvider(t *testing.T) {
	engine := &MockEngine{}
	config := &UserProviderConfig{}

	provider, err := NewUserProvider(engine, config)

	assert.NoError(t, err)
	assert.NotNil(t, provider)
	assert.Equal(t, search.SearchTypeUsers, provider.Type())
}

func TestUserProvider_Type(t *testing.T) {
	engine := &MockEngine{}
	config := &UserProviderConfig{}
	provider, err := NewUserProvider(engine, config)

	require.NoError(t, err)
	assert.Equal(t, search.SearchTypeUsers, provider.Type())
}

func TestUserProvider_Validate(t *testing.T) {
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
				Query:    "testuser",
				Page:     1,
				PageSize: 20,
			},
			wantErr: false,
		},
		{
			name: "invalid page",
			req: &search.SearchRequest{
				Query: "testuser",
				Page:  0,
			},
			wantErr: true,
			errMsg:  "page must be greater than 0",
		},
		{
			name: "invalid page size",
			req: &search.SearchRequest{
				Query:    "testuser",
				Page:     1,
				PageSize: 0,
			},
			wantErr: true,
			errMsg:  "page_size must be greater than 0",
		},
		{
			name: "page size exceeds maximum",
			req: &search.SearchRequest{
				Query:    "testuser",
				Page:     1,
				PageSize: 101,
			},
			wantErr: true,
			errMsg:  "page_size cannot exceed",
		},
	}

	engine := &MockEngine{}
	config := &UserProviderConfig{}
	provider, err := NewUserProvider(engine, config)
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

func TestUserProvider_Search(t *testing.T) {
	tests := []struct {
		name    string
		req     *search.SearchRequest
		wantErr bool
		check   func(t *testing.T, resp *search.SearchResponse)
	}{
		{
			name: "successful search",
			req: &search.SearchRequest{
				Query:    "testuser",
				Page:     1,
				PageSize: 20,
			},
			wantErr: false,
			check: func(t *testing.T, resp *search.SearchResponse) {
				assert.True(t, resp.Success)
				assert.NotNil(t, resp.Data)
				assert.Equal(t, search.SearchTypeUsers, resp.Data.Type)
				assert.Equal(t, int64(2), resp.Data.Total)
				assert.Len(t, resp.Data.Results, 2)
				assert.Equal(t, "testuser1", resp.Data.Results[0].Data["username"])
			},
		},
		{
			name: "search with filter",
			req: &search.SearchRequest{
				Query: "testuser",
				Filter: map[string]interface{}{
					"role": "reader",
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
			wantErr: false, // Validate 不返回 error，Search 返回错误响应
			check: func(t *testing.T, resp *search.SearchResponse) {
				assert.False(t, resp.Success)
				assert.NotNil(t, resp.Error)
				assert.Equal(t, search.ErrCodeInvalidRequest, resp.Error.Code)
			},
		},
	}

	engine := &MockEngine{}
	config := &UserProviderConfig{}
	provider, err := NewUserProvider(engine, config)
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

func TestUserProvider_buildUserFilters(t *testing.T) {
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
			name: "role filter",
			filter: map[string]interface{}{
				"role": "author",
			},
			want: 1,
		},
		{
			name: "status filter",
			filter: map[string]interface{}{
				"status": "active",
			},
			want: 1,
		},
		{
			name: "is_verified filter",
			filter: map[string]interface{}{
				"is_verified": true,
			},
			want: 1,
		},
		{
			name: "multiple filters",
			filter: map[string]interface{}{
				"role":        "reader",
				"status":      "active",
				"is_verified": true,
			},
			want: 3,
		},
	}

	engine := &MockEngine{}
	config := &UserProviderConfig{}
	provider, err := NewUserProvider(engine, config)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filters := provider.buildUserFilters(tt.filter)
			assert.Equal(t, tt.want, len(filters))
		})
	}
}

func TestUserProvider_validateSortField(t *testing.T) {
	engine := &MockEngine{}
	config := &UserProviderConfig{}
	provider, err := NewUserProvider(engine, config)
	require.NoError(t, err)

	tests := []struct {
		name     string
		field    string
		expected string
	}{
		{"valid username", "username", "username"},
		{"valid created_at", "created_at", "created_at"},
		{"valid updated_at", "updated_at", "updated_at"},
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

func TestUserProvider_getPageSize(t *testing.T) {
	engine := &MockEngine{}
	config := &UserProviderConfig{}
	provider, err := NewUserProvider(engine, config)
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

func TestUserProvider_getPage(t *testing.T) {
	engine := &MockEngine{}
	config := &UserProviderConfig{}
	provider, err := NewUserProvider(engine, config)
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

func TestUserProvider_calculateOffset(t *testing.T) {
	engine := &MockEngine{}
	config := &UserProviderConfig{}
	provider, err := NewUserProvider(engine, config)
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
