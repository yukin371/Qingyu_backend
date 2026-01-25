package bookstore

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	bookstoreModel "Qingyu_backend/models/bookstore"
	"Qingyu_backend/models/shared"
)

// TestChineseSearch 测试中文搜索功能
func TestChineseSearch(t *testing.T) {
	// 测试用例结构
	testCases := []struct {
		name           string
		keyword        string
		author         string
		tags           []string
		expectedStatus int
		description    string
		skipMock       bool // 是否跳过 mock 设置
	}{
		{
			name:           "中文关键词_玄幻",
			keyword:        "玄幻",
			expectedStatus: http.StatusOK,
			description:    "搜索包含'玄幻'的书籍",
		},
		{
			name:           "中文关键词_小说",
			keyword:        "小说",
			expectedStatus: http.StatusOK,
			description:    "搜索包含'小说'的书籍",
		},
		{
			name:           "中文关键词_斗罗",
			keyword:        "斗罗",
			expectedStatus: http.StatusOK,
			description:    "搜索包含'斗罗'的书籍",
		},
		{
			name:           "中文作者搜索",
			keyword:        "书",
			author:         "唐家三少",
			expectedStatus: http.StatusOK,
			description:    "按作者唐家三少搜索",
		},
		{
			name:           "中文标签搜索",
			keyword:        "小说",
			tags:           []string{"热血"},
			expectedStatus: http.StatusOK,
			description:    "按标签'热血'搜索",
		},
		{
			name:           "组合搜索_关键词加作者",
			keyword:        "斗罗",
			author:         "唐家三少",
			expectedStatus: http.StatusOK,
			description:    "关键词+作者组合搜索",
		},
		{
			name:           "无关键词无过滤器",
			keyword:        "",
			expectedStatus: http.StatusBadRequest,
			description:    "没有任何搜索条件应该返回400错误",
			skipMock:       true, // 这个测试用例不会调用 service
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			mockService := new(MockBookstoreService)
			router := setupBookstoreTestRouter(mockService)

			// 模拟返回数据
			books := []*bookstoreModel.Book{
				{
					IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
					Title:            "测试中文书名",
					Author:           "测试作者",
				},
			}

			// 设置 mock 期望
			if !tc.skipMock && tc.expectedStatus == http.StatusOK {
				mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(books, int64(1), nil)
			}

			// When - 构建请求 URL
			requestURL := "/api/v1/bookstore/books/search"

			// 添加参数
			params := url.Values{}
			if tc.keyword != "" {
				params.Add("keyword", tc.keyword)
			}
			if tc.author != "" {
				params.Add("author", tc.author)
			}
			for _, tag := range tc.tags {
				params.Add("tags", tag)
			}

			if len(params) > 0 {
				requestURL += "?" + params.Encode()
			}

			req, _ := http.NewRequest("GET", requestURL, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Then
			assert.Equal(t, tc.expectedStatus, w.Code, tc.description)

			if tc.expectedStatus == http.StatusOK {
				// 验证响应体包含中文内容
				responseBody := w.Body.String()
				assert.Contains(t, responseBody, "code", "响应应该包含code字段")
				assert.Contains(t, responseBody, "data", "响应应该包含data字段")
			}
		})
	}
}

// TestChineseSearchJSONResponse 测试中文 JSON 响应不转义
func TestChineseSearchJSONResponse(t *testing.T) {
	// Given
	mockService := new(MockBookstoreService)
	router := setupBookstoreTestRouter(mockService)

	// 创建包含中文的测试数据
	books := []*bookstoreModel.Book{
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			Title:            "斗罗大陆",
			Author:           "唐家三少",
			Introduction:     "这是一个关于玄幻世界的故事",
			Categories:       []string{"玄幻", "魔法"},
			Tags:             []string{"热血", "冒险"},
		},
	}

	mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(books, int64(1), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/bookstore/books/search?keyword=斗罗", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	responseBody := w.Body.String()

	// 验证中文字符没有被转义为 Unicode
	assert.Contains(t, responseBody, "斗罗大陆", "响应应该包含中文字符'斗罗大陆'")
	assert.Contains(t, responseBody, "唐家三少", "响应应该包含中文字符'唐家三少'")
	assert.Contains(t, responseBody, "玄幻", "响应应该包含中文字符'玄幻'")
	assert.NotContains(t, responseBody, "\\u6597", "响应不应该包含Unicode转义字符")
}

// TestChineseSearchURLEncoding 测试中文 URL 编码
func TestChineseSearchURLEncoding(t *testing.T) {
	testCases := []struct {
		name        string
		keyword     string
		author      string
		description string
	}{
		{
			name:        "URL编码_关键词",
			keyword:     "修真世界",
			description: "测试中文关键词的URL编码",
		},
		{
			name:        "URL编码_作者",
			keyword:     "小说",
			author:      "我吃西红柿",
			description: "测试中文作者的URL编码",
		},
		{
			name:        "URL编码_特殊字符",
			keyword:     "三体（地球往事）",
			description: "测试包含特殊字符的中文URL编码",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			mockService := new(MockBookstoreService)
			router := setupBookstoreTestRouter(mockService)

			books := []*bookstoreModel.Book{}
			mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(books, int64(0), nil)

			// When - 使用 URL 编码
			requestURL := "/api/v1/bookstore/books/search"
			params := url.Values{}
			params.Add("keyword", tc.keyword)
			if tc.author != "" {
				params.Add("author", tc.author)
			}
			requestURL += "?" + params.Encode()

			req, _ := http.NewRequest("GET", requestURL, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Then
			assert.Equal(t, http.StatusOK, w.Code, tc.description)
		})
	}
}

// TestChineseSearchPagination 测试中文搜索的分页功能
func TestChineseSearchPagination(t *testing.T) {
	testCases := []struct {
		name        string
		keyword     string
		page        int
		size        int
		description string
	}{
		{
			name:        "分页_第一页",
			keyword:     "玄幻",
			page:        1,
			size:        10,
			description: "测试中文搜索的第一页",
		},
		{
			name:        "分页_第二页",
			keyword:     "小说",
			page:        2,
			size:        20,
			description: "测试中文搜索的第二页",
		},
		{
			name:        "分页_大页码",
			keyword:     "武侠",
			page:        5,
			size:        50,
			description: "测试中文搜索的大页码",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			mockService := new(MockBookstoreService)
			router := setupBookstoreTestRouter(mockService)

			books := []*bookstoreModel.Book{}
			mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(books, int64(0), nil)

			// When
			requestURL := "/api/v1/bookstore/books/search"
			params := url.Values{}
			params.Add("keyword", tc.keyword)
			params.Add("page", string(rune(tc.page)))
			params.Add("size", string(rune(tc.size)))
			requestURL += "?" + params.Encode()

			req, _ := http.NewRequest("GET", requestURL, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Then
			assert.Equal(t, http.StatusOK, w.Code, tc.description)
		})
	}
}

// TestChineseSearchSortOptions 测试中文搜索的排序功能
func TestChineseSearchSortOptions(t *testing.T) {
	testCases := []struct {
		name        string
		keyword     string
		sortBy      string
		sortOrder   string
		description string
	}{
		{
			name:        "排序_最新创建",
			keyword:     "玄幻",
			sortBy:      "created_at",
			sortOrder:   "desc",
			description: "按创建时间降序",
		},
		{
			name:        "排序_最早创建",
			keyword:     "小说",
			sortBy:      "created_at",
			sortOrder:   "asc",
			description: "按创建时间升序",
		},
		{
			name:        "排序_字数降序",
			keyword:     "武侠",
			sortBy:      "word_count",
			sortOrder:   "desc",
			description: "按字数降序",
		},
		{
			name:        "排序_字数升序",
			keyword:     "仙侠",
			sortBy:      "word_count",
			sortOrder:   "asc",
			description: "按字数升序",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			mockService := new(MockBookstoreService)
			router := setupBookstoreTestRouter(mockService)

			books := []*bookstoreModel.Book{}
			mockService.On("SearchBooksWithFilter", mock.Anything, mock.Anything).Return(books, int64(0), nil)

			// When
			requestURL := "/api/v1/bookstore/books/search"
			params := url.Values{}
			params.Add("keyword", tc.keyword)
			params.Add("sortBy", tc.sortBy)
			params.Add("sortOrder", tc.sortOrder)
			requestURL += "?" + params.Encode()

			req, _ := http.NewRequest("GET", requestURL, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Then
			assert.Equal(t, http.StatusOK, w.Code, tc.description)
		})
	}
}
