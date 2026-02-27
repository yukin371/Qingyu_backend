package shared

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetUserID(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func(*gin.Context)
		expectedUserID string
		expectedOk     bool
	}{
		{
			name: "用户ID存在且类型正确",
			setupContext: func(c *gin.Context) {
				c.Set("user_id", "test-user-123")
			},
			expectedUserID: "test-user-123",
			expectedOk:     true,
		},
		{
			name:           "用户ID不存在",
			setupContext:   func(c *gin.Context) {},
			expectedUserID: "",
			expectedOk:     false,
		},
		{
			name: "用户ID类型错误",
			setupContext: func(c *gin.Context) {
				c.Set("user_id", 123) // 错误类型
			},
			expectedUserID: "",
			expectedOk:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.setupContext(c)

			userID, ok := GetUserID(c)

			assert.Equal(t, tt.expectedUserID, userID)
			assert.Equal(t, tt.expectedOk, ok)

			// 如果期望失败，检查是否发送了响应
			if !tt.expectedOk {
				assert.True(t, w.Header().Get("Content-Type") == "application/json; charset=utf-8" || w.Code != 0)
			}
		})
	}
}

func TestGetUserIDOptional(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func(*gin.Context)
		expectedUserID string
	}{
		{
			name: "用户ID存在",
			setupContext: func(c *gin.Context) {
				c.Set("user_id", "test-user-123")
			},
			expectedUserID: "test-user-123",
		},
		{
			name:           "用户ID不存在",
			setupContext:   func(c *gin.Context) {},
			expectedUserID: "",
		},
		{
			name: "用户ID类型错误",
			setupContext: func(c *gin.Context) {
				c.Set("user_id", 123)
			},
			expectedUserID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.setupContext(c)

			userID := GetUserIDOptional(c)
			assert.Equal(t, tt.expectedUserID, userID)
		})
	}
}

func TestGetRequiredParam(t *testing.T) {
	tests := []struct {
		name         string
		paramValue   string
		displayName  string
		expectedOk   bool
		expectedCode int
	}{
		{
			name:         "参数存在",
			paramValue:   "test-id-123",
			displayName:  "书籍ID",
			expectedOk:   true,
			expectedCode: 0,
		},
		{
			name:         "参数为空",
			paramValue:   "",
			displayName:  "书籍ID",
			expectedOk:   false,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = []gin.Param{
				{Key: "id", Value: tt.paramValue},
			}

			param, ok := GetRequiredParam(c, "id", tt.displayName)

			assert.Equal(t, tt.expectedOk, ok)

			if tt.expectedOk {
				assert.Equal(t, tt.paramValue, param)
			} else {
				assert.Equal(t, tt.expectedCode, w.Code)
			}
		})
	}
}

func TestGetRequiredQuery(t *testing.T) {
	tests := []struct {
		name         string
		queryValue   string
		displayName  string
		expectedOk   bool
		expectedCode int
	}{
		{
			name:         "查询参数存在",
			queryValue:   "search-keyword",
			displayName:  "搜索关键词",
			expectedOk:   true,
			expectedCode: 0,
		},
		{
			name:         "查询参数为空",
			queryValue:   "",
			displayName:  "搜索关键词",
			expectedOk:   false,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.queryValue != "" {
				c.Request = httptest.NewRequest("GET", "/?query="+tt.queryValue, nil)
			} else {
				c.Request = httptest.NewRequest("GET", "/", nil)
			}

			param, ok := GetRequiredQuery(c, "query", tt.displayName)

			assert.Equal(t, tt.expectedOk, ok)

			if tt.expectedOk {
				assert.Equal(t, tt.queryValue, param)
			} else {
				assert.Equal(t, tt.expectedCode, w.Code)
			}
		})
	}
}

func TestGetPaginationParamsStandard(t *testing.T) {
	tests := []struct {
		name              string
		queryPage         string
		querySize         string
		expectedPage      int
		expectedPageSize  int
		expectedLimit     int
		expectedOffset    int
	}{
		{
			name:           "默认值",
			queryPage:      "",
			querySize:      "",
			expectedPage:   1,
			expectedPageSize: 20,
			expectedLimit:  20,
			expectedOffset: 0,
		},
		{
			name:           "正常值",
			queryPage:      "2",
			querySize:      "30",
			expectedPage:   2,
			expectedPageSize: 30,
			expectedLimit:  30,
			expectedOffset: 30,
		},
		{
			name:           "页码小于1",
			queryPage:      "0",
			querySize:      "20",
			expectedPage:   1,
			expectedPageSize: 20,
			expectedLimit:  20,
			expectedOffset: 0,
		},
		{
			name:           "每页数量超过最大值",
			queryPage:      "1",
			querySize:      "150",
			expectedPage:   1,
			expectedPageSize: 100,
			expectedLimit:  100,
			expectedOffset: 0,
		},
		{
			name:           "无效数字使用默认值",
			queryPage:      "invalid",
			querySize:      "invalid",
			expectedPage:   1,
			expectedPageSize: 20,
			expectedLimit:  20,
			expectedOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			url := "/?"
			if tt.queryPage != "" {
				url += "page=" + tt.queryPage
			}
			if tt.querySize != "" {
				if tt.queryPage != "" {
					url += "&"
				}
				url += "size=" + tt.querySize
			}
			c.Request = httptest.NewRequest("GET", url, nil)

			params := GetPaginationParamsStandard(c)

			assert.Equal(t, tt.expectedPage, params.Page)
			assert.Equal(t, tt.expectedPageSize, params.PageSize)
			assert.Equal(t, tt.expectedLimit, params.Limit)
			assert.Equal(t, tt.expectedOffset, params.Offset)
		})
	}
}

func TestGetIntParam(t *testing.T) {
	tests := []struct {
		name         string
		paramValue   string
		defaultValue int
		min          int
		max          int
		expected     int
	}{
		{
			name:         "正常值",
			paramValue:   "50",
			defaultValue: 10,
			min:          1,
			max:          100,
			expected:     50,
		},
		{
			name:         "空值使用默认值",
			paramValue:   "",
			defaultValue: 10,
			min:          0,
			max:          0,
			expected:     10,
		},
		{
			name:         "小于最小值",
			paramValue:   "0",
			defaultValue: 10,
			min:          1,
			max:          100,
			expected:     1,
		},
		{
			name:         "大于最大值",
			paramValue:   "150",
			defaultValue: 10,
			min:          1,
			max:          100,
			expected:     100,
		},
		{
			name:         "无效值使用默认值",
			paramValue:   "invalid",
			defaultValue: 10,
			min:          0,
			max:          0,
			expected:     10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = []gin.Param{
				{Key: "limit", Value: tt.paramValue},
			}

			result := GetIntParam(c, "limit", false, tt.defaultValue, tt.min, tt.max)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateBatchIDs(t *testing.T) {
	tests := []struct {
		name         string
		ids          []string
		displayName  string
		expectedOk   bool
		expectedCode int
	}{
		{
			name:         "有效ID列表",
			ids:          []string{"id1", "id2", "id3"},
			displayName:  "用户ID",
			expectedOk:   true,
			expectedCode: 0,
		},
		{
			name:         "空ID列表",
			ids:          []string{},
			displayName:  "用户ID",
			expectedOk:   false,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "ID列表超过限制",
			ids:          make([]string, 1001),
			displayName:  "用户ID",
			expectedOk:   false,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			ok := ValidateBatchIDs(c, tt.ids, tt.displayName)

			assert.Equal(t, tt.expectedOk, ok)

			if !tt.expectedOk {
				assert.Equal(t, tt.expectedCode, w.Code)
			}
		})
	}
}
