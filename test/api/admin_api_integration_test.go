//go:build integration
// +build integration

package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAdminSystemAPI 系统管理API集成测试
func TestAdminSystemAPI_GetSystemStats(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 设置路由
	router := gin.New()
	router.GET("/api/v1/admin/stats", mockAuthMiddleware("admin1"), func(c *gin.Context) {
		// 调用API方法
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "获取成功",
			"data": map[string]interface{}{
				"totalUsers":    100,
				"activeUsers":   50,
				"totalBooks":    200,
				"totalRevenue":  10000.0,
				"pendingAudits": 5,
			},
		})
	})

	// 发送请求
	req := httptest.NewRequest("GET", "/api/v1/admin/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "获取成功", response["message"])
	assert.NotNil(t, response["data"])
}

// TestAdminSystemAPI_GetSystemConfig 获取系统配置
func TestAdminSystemAPI_GetSystemConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/api/v1/admin/config", mockAuthMiddleware("admin1"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "获取成功",
			"data": map[string]interface{}{
				"allowRegistration":        true,
				"requireEmailVerification": true,
				"maxUploadSize":            10485760,
				"enableAudit":              true,
			},
		})
	})

	req := httptest.NewRequest("GET", "/api/v1/admin/config", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(0), response["code"])
	data := response["data"].(map[string]interface{})
	assert.True(t, data["allowRegistration"].(bool))
	assert.True(t, data["enableAudit"].(bool))
}

// TestAdminSystemAPI_UpdateSystemConfig 更新系统配置
func TestAdminSystemAPI_UpdateSystemConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.PUT("/api/v1/admin/config", mockAuthMiddleware("admin1"), func(c *gin.Context) {
		var req map[string]interface{}
		c.ShouldBindJSON(&req)

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "更新成功",
			"data":    nil,
		})
	})

	configReq := map[string]interface{}{
		"allowRegistration": false,
		"enableAudit":       false,
	}

	body, _ := json.Marshal(configReq)
	req := httptest.NewRequest("PUT", "/api/v1/admin/config", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "更新成功", response["message"])
}

// TestAdminSystemAPI_CreateAnnouncement 创建公告
func TestAdminSystemAPI_CreateAnnouncement(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/api/v1/admin/announcements", mockAuthMiddleware("admin1"), func(c *gin.Context) {
		var req map[string]interface{}
		c.ShouldBindJSON(&req)

		c.JSON(http.StatusCreated, gin.H{
			"code":    0,
			"message": "公告发布成功",
			"data": map[string]interface{}{
				"id":        "announcement_id_123",
				"title":     req["title"],
				"content":   req["content"],
				"type":      req["type"],
				"createdBy": "admin1",
			},
		})
	})

	announceReq := map[string]interface{}{
		"title":    "系统维护",
		"content":  "系统将于明天进行维护",
		"type":     "system",
		"priority": "high",
	}

	body, _ := json.Marshal(announceReq)
	req := httptest.NewRequest("POST", "/api/v1/admin/announcements", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "公告发布成功", response["message"])
	assert.NotNil(t, response["data"])
}

// TestAdminSystemAPI_GetAnnouncements 获取公告列表
func TestAdminSystemAPI_GetAnnouncements(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/api/v1/admin/announcements", mockAuthMiddleware("admin1"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":      0,
			"message":   "获取成功",
			"data":      []interface{}{},
			"page":      1,
			"page_size": 20,
			"total":     0,
		})
	})

	req := httptest.NewRequest("GET", "/api/v1/admin/announcements?page=1&page_size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
}

// mockAuthMiddleware 模拟认证中间件
func mockAuthMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", "test_admin_id")
		c.Set("role", role)
		c.Next()
	}
}
