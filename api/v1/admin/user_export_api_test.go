package admin

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xuri/excelize/v2"

	"Qingyu_backend/models/users"
)

// setupUserExportAPITestRouter 设置用户导出测试路由
func setupUserExportAPITestRouter(userService *MockUserAdminService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	api := NewUserExportAPI(userService)

	v1 := r.Group("/api/v1/admin")
	{
		v1.GET("/users/export", api.ExportUsers)
		v1.GET("/users/export/template", api.GetUserExportTemplate)
	}

	return r
}

// ==================== ExportUsers Tests ====================

func TestUserExportAPI_ExportUsersToCSV_Success(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserExportAPITestRouter(mockService)

	userList := []*users.User{
		{Username: "user1", Email: "user1@example.com", Roles: []string{"reader"}, Status: users.UserStatusActive},
		{Username: "user2", Email: "user2@example.com", Roles: []string{"author"}, Status: users.UserStatusActive},
	}

	mockService.On("GetUserList", mock.Anything, mock.Anything, 1, 10000).Return(userList, int64(2), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/users/export?format=csv", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/csv")

	// 验证CSV内容
	r := csv.NewReader(strings.NewReader(w.Body.String()))
	records, err := r.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(records)) // Header + 2 rows
	assert.Equal(t, "用户名", records[0][1])
	assert.Equal(t, "user1", records[1][1])

	mockService.AssertExpectations(t)
}

func TestUserExportAPI_ExportUsersToExcel_Success(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserExportAPITestRouter(mockService)

	userList := []*users.User{
		{Username: "user1", Email: "user1@example.com", Roles: []string{"reader"}, Status: users.UserStatusActive},
		{Username: "user2", Email: "user2@example.com", Roles: []string{"author"}, Status: users.UserStatusActive},
	}

	mockService.On("GetUserList", mock.Anything, mock.Anything, 1, 10000).Return(userList, int64(2), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/users/export?format=xlsx", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", w.Header().Get("Content-Type"))

	// 验证Excel内容
	_, err := excelize.OpenReader(bytes.NewReader(w.Body.Bytes()))
	assert.NoError(t, err)

	mockService.AssertExpectations(t)
}

func TestUserExportAPI_ExportUsersWithFilter_Success(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserExportAPITestRouter(mockService)

	userList := []*users.User{
		{Username: "author1", Email: "author1@example.com", Roles: []string{"author"}, Status: users.UserStatusActive},
	}

	mockService.On("GetUserList", mock.Anything, mock.Anything, 1, 10000).Return(userList, int64(1), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/users/export?format=csv&role=author&status=active", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/csv")

	mockService.AssertExpectations(t)
}

func TestUserExportAPI_ExportUsersEmpty_Success(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserExportAPITestRouter(mockService)

	userList := []*users.User{}
	mockService.On("GetUserList", mock.Anything, mock.Anything, 1, 10000).Return(userList, int64(0), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/users/export?format=csv", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/csv")

	// 验证CSV只有header
	r := csv.NewReader(strings.NewReader(w.Body.String()))
	records, err := r.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(records)) // Only header

	mockService.AssertExpectations(t)
}

func TestUserExportAPI_ExportUsersInvalidFormat_Error(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserExportAPITestRouter(mockService)

	// When - 无效格式应该在验证阶段就返回，不调用service
	req, _ := http.NewRequest("GET", "/api/v1/admin/users/export?format=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEqual(t, float64(0), response["code"])
}

func TestUserExportAPI_GetUserExportTemplate_Success(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserExportAPITestRouter(mockService)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/users/export/template", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].([]interface{})
	assert.True(t, len(data) > 0)
}
