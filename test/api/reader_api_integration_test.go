//go:build integration
// +build integration

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	readerAPI "Qingyu_backend/api/v1/reader"
	readerModel "Qingyu_backend/models/reading/reader"
	"Qingyu_backend/service/reading"
	"Qingyu_backend/test/testutil"
)

// TestReaderAPIIntegration 阅读器API集成测试
type TestReaderAPIIntegration struct {
	readerService *reading.ReaderService
	chapterAPI    *readerAPI.ChaptersAPI
	settingAPI    *readerAPI.SettingAPI
	router        *gin.Engine
	testUserID    string
	testBookID    primitive.ObjectID
	testChapterID primitive.ObjectID
}

// setupReaderAPITest 设置阅读器API测试环境
func setupReaderAPITest(t *testing.T) *TestReaderAPIIntegration {
	gin.SetMode(gin.TestMode)

	// 使用测试工具创建数据库连接
	ctx := context.Background()
	cleanup := testutil.SetupTestDB(t, ctx)
	t.Cleanup(cleanup)

	// TODO: 创建 Repository 的 Mock 或测试实例
	// 目前使用nil，需要补充完整的Repository Mock
	// 参考 test/service/reader_service_test.go 中的Mock实现

	// 创建ReaderService实例
	// readerService := reading.NewReaderService(chapterRepo, progressRepo, settingsRepo)

	// 暂时返回nil，直到Repository Mock准备好
	t.Skip("等待 Repository Mock 实现")
	return nil
}

// TestGetChapterByID_Success 测试获取章节成功
func TestGetChapterByID_Success(t *testing.T) {
	suite := setupReaderAPITest(t)
	if suite == nil {
		return
	}

	// 准备测试数据
	chapterID := suite.testChapterID.Hex()

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/reader/chapters/"+chapterID, nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
	assert.NotNil(t, response["data"])
}

// TestGetChapterContent_Success 测试获取章节内容成功
func TestGetChapterContent_Success(t *testing.T) {
	suite := setupReaderAPITest(t)
	if suite == nil {
		return
	}

	chapterID := suite.testChapterID.Hex()

	// 创建带认证的请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/reader/chapters/"+chapterID+"/content", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()

	// 模拟认证中间件设置用户ID
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", suite.testUserID)
		c.Next()
	})
	router.GET("/api/v1/reader/chapters/:id/content", suite.chapterAPI.GetChapterContent)

	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
	assert.NotNil(t, response["data"])
}

// TestGetChapterContent_Unauthorized 测试未认证获取章节内容
func TestGetChapterContent_Unauthorized(t *testing.T) {
	suite := setupReaderAPITest(t)
	if suite == nil {
		return
	}

	chapterID := suite.testChapterID.Hex()

	// 创建不带认证的请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/reader/chapters/"+chapterID+"/content", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	// 验证响应 - 应该返回401
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestGetBookChapters_Success 测试获取书籍章节列表
func TestGetBookChapters_Success(t *testing.T) {
	suite := setupReaderAPITest(t)
	if suite == nil {
		return
	}

	bookID := suite.testBookID.Hex()

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/reader/chapters?bookId="+bookID+"&page=1&size=20", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
	data := response["data"].(map[string]interface{})
	assert.NotNil(t, data["chapters"])
	assert.NotNil(t, data["total"])
}

// TestGetReadingSettings_Success 测试获取阅读设置
func TestGetReadingSettings_Success(t *testing.T) {
	suite := setupReaderAPITest(t)
	if suite == nil {
		return
	}

	// 创建带认证的请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/reader/settings", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", suite.testUserID)
		c.Next()
	})
	router.GET("/api/v1/reader/settings", suite.settingAPI.GetReadingSettings)

	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
}

// TestSaveReadingSettings_Success 测试保存阅读设置
func TestSaveReadingSettings_Success(t *testing.T) {
	suite := setupReaderAPITest(t)
	if suite == nil {
		return
	}

	// 准备请求数据
	settings := readerModel.ReadingSettings{
		UserID:          suite.testUserID,
		FontSize:        16,
		FontFamily:      "宋体",
		LineHeight:      1.5,
		BackgroundColor: "#FFFFFF",
		TextColor:       "#000000",
		PageMode:        "scroll",
		AutoSave:        true,
		ShowProgress:    true,
		UpdatedAt:       time.Now(),
	}

	body, _ := json.Marshal(settings)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reader/settings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", suite.testUserID)
		c.Next()
	})
	router.POST("/api/v1/reader/settings", suite.settingAPI.SaveReadingSettings)

	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
}

// TestAPI_ValidationErrors 测试参数验证错误
func TestAPI_ValidationErrors(t *testing.T) {
	suite := setupReaderAPITest(t)
	if suite == nil {
		return
	}

	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
	}{
		{
			name:           "获取章节列表 - 缺少bookId参数",
			method:         http.MethodGet,
			url:            "/api/v1/reader/chapters",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "获取章节 - 无效的章节ID",
			method:         http.MethodGet,
			url:            "/api/v1/reader/chapters/invalid-id",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()

			suite.router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

/*
// 以下是需要Repository Mock完成后才能实现的测试

// TestGetChapterNavigation_Success 测试章节导航
func TestGetChapterNavigation_Success(t *testing.T) {
	// 测试上一章、下一章功能
}

// TestSaveReadingProgress_Success 测试保存阅读进度
func TestSaveReadingProgress_Success(t *testing.T) {
	// 测试保存阅读进度
}

// TestGetRecentReading_Success 测试获取最近阅读
func TestGetRecentReading_Success(t *testing.T) {
	// 测试获取最近阅读记录
}

// TestAnnotations_CRUD 测试笔记CRUD
func TestAnnotations_CRUD(t *testing.T) {
	// 测试创建、获取、更新、删除笔记
}
*/
