package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/users"
	adminrepo "Qingyu_backend/repository/interfaces/admin"
	adminservice "Qingyu_backend/service/admin"
)

// MockUserAdminService 模拟UserAdminService
type MockUserAdminService struct {
	mock.Mock
}

func (m *MockUserAdminService) GetUserList(ctx context.Context, filter *adminrepo.UserFilter, page, size int) ([]*users.User, int64, error) {
	args := m.Called(ctx, filter, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*users.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserAdminService) GetUserDetail(ctx context.Context, userID string) (*users.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

func (m *MockUserAdminService) UpdateUserStatus(ctx context.Context, userID string, status users.UserStatus) error {
	args := m.Called(ctx, userID, status)
	return args.Error(0)
}

func (m *MockUserAdminService) UpdateUserRole(ctx context.Context, userID, role string) error {
	args := m.Called(ctx, userID, role)
	return args.Error(0)
}

func (m *MockUserAdminService) DeleteUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserAdminService) GetUserActivities(ctx context.Context, userID string, page, size int) ([]*users.UserActivity, int64, error) {
	args := m.Called(ctx, userID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*users.UserActivity), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserAdminService) GetUserStatistics(ctx context.Context, userID string) (*users.UserStatistics, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.UserStatistics), args.Error(1)
}

func (m *MockUserAdminService) ResetUserPassword(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockUserAdminService) BatchUpdateStatus(ctx context.Context, userIds []string, status users.UserStatus) error {
	args := m.Called(ctx, userIds, status)
	return args.Error(0)
}

func (m *MockUserAdminService) BatchDeleteUsers(ctx context.Context, userIds []string) error {
	args := m.Called(ctx, userIds)
	return args.Error(0)
}

func (m *MockUserAdminService) SearchUsers(ctx context.Context, keyword string, page, size int) ([]*users.User, int64, error) {
	args := m.Called(ctx, keyword, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*users.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserAdminService) CountByStatus(ctx context.Context) (map[string]int64, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *MockUserAdminService) GetUsersByRole(ctx context.Context, role string, page, size int) ([]*users.User, int64, error) {
	args := m.Called(ctx, role, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*users.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserAdminService) GetRecentUsers(ctx context.Context, limit int) ([]*users.User, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*users.User), args.Error(1)
}

func (m *MockUserAdminService) GetActiveUsers(ctx context.Context, days, limit int) ([]*users.User, error) {
	args := m.Called(ctx, days, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*users.User), args.Error(1)
}

func (m *MockUserAdminService) CreateUser(ctx context.Context, req *adminservice.CreateUserRequest) (*users.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

func (m *MockUserAdminService) BatchCreateUsers(ctx context.Context, req *adminservice.BatchCreateUserRequest) ([]*users.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*users.User), args.Error(1)
}

// setupUserAdminTestRouter 设置用户管理测试路由
func setupUserAdminTestRouter(userAdminService *MockUserAdminService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	api := NewUserAdminAPI(userAdminService)

	v1 := r.Group("/api/v1/admin/users")
	{
		v1.GET("", api.ListUsers)
		v1.GET(":id", api.GetUserDetail)
		v1.PUT(":id/status", api.UpdateUserStatus)
		v1.PUT(":id/role", api.UpdateUserRole)
		v1.DELETE(":id", api.DeleteUser)
		v1.GET(":id/activities", api.GetUserActivities)
		v1.GET(":id/statistics", api.GetUserStatistics)
		v1.POST(":id/reset-password", api.ResetUserPassword)
		v1.POST("batch-update-status", api.BatchUpdateStatus)
		v1.POST("batch-delete", api.BatchDeleteUsers)
		v1.GET("search", api.SearchUsers)
		v1.GET("count-by-status", api.CountByStatus)
	}

	return r
}

// ==================== ListUsers Tests ====================

func TestUserAdminAPI_ListUsers_Success(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	expectedUsers := []*users.User{
		func() *users.User {
			u := &users.User{Username: "user1", Email: "user1@example.com"}
			u.ID = primitive.NewObjectID()
			return u
		}(),
		func() *users.User {
			u := &users.User{Username: "user2", Email: "user2@example.com"}
			u.ID = primitive.NewObjectID()
			return u
		}(),
	}

	mockService.On("GetUserList", mock.Anything, mock.Anything, 1, 20).
		Return(expectedUsers, int64(2), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/users?page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "操作成功", response["message"])

	mockService.AssertExpectations(t)
}

func TestUserAdminAPI_ListUsers_WithFilter(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	mockService.On("GetUserList", mock.Anything, mock.MatchedBy(func(filter *adminrepo.UserFilter) bool {
		return filter.Keyword == "test" && filter.Role == "author"
	}), 1, 20).Return([]*users.User{}, int64(0), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/users?keyword=test&role=author", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestUserAdminAPI_ListUsers_ServiceError(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	mockService.On("GetUserList", mock.Anything, mock.Anything, 1, 20).
		Return(nil, int64(0), assert.AnError)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== GetUserDetail Tests ====================

func TestUserAdminAPI_GetUserDetail_Success(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	expectedUser := &users.User{
		Username: "testuser",
		Email:    "test@example.com",
	}
	expectedUser.ID = primitive.NewObjectID()
	userID := expectedUser.ID.Hex()

	mockService.On("GetUserDetail", mock.Anything, userID).Return(expectedUser, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/users/"+userID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])

	mockService.AssertExpectations(t)
}

func TestUserAdminAPI_GetUserDetail_EmptyID(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/users/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	// Gin redirects trailing slash, so we expect 301
	assert.Equal(t, http.StatusMovedPermanently, w.Code)
}

func TestUserAdminAPI_GetUserDetail_NotFound(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	userID := primitive.NewObjectID().Hex()
	mockService.On("GetUserDetail", mock.Anything, userID).Return(nil, adminservice.ErrUserNotFound)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/users/"+userID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== UpdateUserStatus Tests ====================

func TestUserAdminAPI_UpdateUserStatus_Success(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	userID := primitive.NewObjectID().Hex()
	reqBody := UpdateUserStatusRequest{Status: "active"}

	mockService.On("UpdateUserStatus", mock.Anything, userID, users.UserStatus("active")).Return(nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/admin/users/"+userID+"/status", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])

	mockService.AssertExpectations(t)
}

func TestUserAdminAPI_UpdateUserStatus_EmptyID(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	reqBody := UpdateUserStatusRequest{Status: "active"}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/admin/users//status", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserAdminAPI_UpdateUserStatus_InvalidJSON(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	userID := primitive.NewObjectID().Hex()

	// When
	req, _ := http.NewRequest("PUT", "/api/v1/admin/users/"+userID+"/status", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserAdminAPI_UpdateUserStatus_UserNotFound(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	userID := primitive.NewObjectID().Hex()
	reqBody := UpdateUserStatusRequest{Status: "active"}

	mockService.On("UpdateUserStatus", mock.Anything, userID, users.UserStatus("active")).Return(adminservice.ErrUserNotFound)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/admin/users/"+userID+"/status", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.AssertExpectations(t)
}

func TestUserAdminAPI_UpdateUserStatus_CannotModifySuperAdmin(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	userID := primitive.NewObjectID().Hex()
	reqBody := UpdateUserStatusRequest{Status: "active"}

	mockService.On("UpdateUserStatus", mock.Anything, userID, users.UserStatus("active")).Return(adminservice.ErrCannotModifySuperAdmin)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/admin/users/"+userID+"/status", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusForbidden, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== UpdateUserRole Tests ====================

func TestUserAdminAPI_UpdateUserRole_Success(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	userID := primitive.NewObjectID().Hex()
	reqBody := UpdateUserRoleRequest{Role: "admin"}

	mockService.On("UpdateUserRole", mock.Anything, userID, "admin").Return(nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/admin/users/"+userID+"/role", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestUserAdminAPI_UpdateUserRole_EmptyID(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	reqBody := UpdateUserRoleRequest{Role: "admin"}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/admin/users//role", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserAdminAPI_UpdateUserRole_InvalidRole(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	userID := primitive.NewObjectID().Hex()
	reqBody := UpdateUserRoleRequest{Role: "admin"}

	mockService.On("UpdateUserRole", mock.Anything, userID, "admin").Return(adminservice.ErrInvalidRole)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/admin/users/"+userID+"/role", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== DeleteUser Tests ====================

func TestUserAdminAPI_DeleteUser_Success(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	userID := primitive.NewObjectID().Hex()
	mockService.On("DeleteUser", mock.Anything, userID).Return(nil)

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/admin/users/"+userID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "操作成功", response["message"])

	mockService.AssertExpectations(t)
}

func TestUserAdminAPI_DeleteUser_CannotModifySuperAdmin(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	userID := primitive.NewObjectID().Hex()
	mockService.On("DeleteUser", mock.Anything, userID).Return(adminservice.ErrCannotModifySuperAdmin)

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/admin/users/"+userID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusForbidden, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== GetUserActivities Tests ====================

func TestUserAdminAPI_GetUserActivities_Success(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	userID := primitive.NewObjectID().Hex()
	expectedActivities := []*users.UserActivity{
		{Action: "login", Description: "用户登录"},
	}

	mockService.On("GetUserActivities", mock.Anything, userID, 1, 20).
		Return(expectedActivities, int64(1), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/users/"+userID+"/activities?page=1&size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestUserAdminAPI_GetUserActivities_InvalidUserID(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	userID := primitive.NewObjectID().Hex()
	mockService.On("GetUserActivities", mock.Anything, userID, 1, 20).Return(nil, int64(0), adminservice.ErrInvalidUserID)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/users/"+userID+"/activities", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== GetUserStatistics Tests ====================

func TestUserAdminAPI_GetUserStatistics_Success(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	userID := primitive.NewObjectID().Hex()
	expectedStats := &users.UserStatistics{
		TotalBooks: 10,
		TotalWords: 100000,
	}

	mockService.On("GetUserStatistics", mock.Anything, userID).Return(expectedStats, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/users/"+userID+"/statistics", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestUserAdminAPI_GetUserStatistics_EmptyID(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/users//statistics", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==================== ResetUserPassword Tests ====================

func TestUserAdminAPI_ResetUserPassword_Success(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	userID := primitive.NewObjectID().Hex()
	newPassword := "newPassword123"
	mockService.On("ResetUserPassword", mock.Anything, userID).Return(newPassword, nil)

	// When
	req, _ := http.NewRequest("POST", "/api/v1/admin/users/"+userID+"/reset-password", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, newPassword, data["newPassword"])

	mockService.AssertExpectations(t)
}

func TestUserAdminAPI_ResetUserPassword_UserNotFound(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	userID := primitive.NewObjectID().Hex()
	mockService.On("ResetUserPassword", mock.Anything, userID).Return("", adminservice.ErrUserNotFound)

	// When
	req, _ := http.NewRequest("POST", "/api/v1/admin/users/"+userID+"/reset-password", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== BatchUpdateStatus Tests ====================

func TestUserAdminAPI_BatchUpdateStatus_Success(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	reqBody := BatchUpdateStatusRequest{
		UserIds: []string{"user1", "user2"},
		Status:  "active",
	}

	mockService.On("BatchUpdateStatus", mock.Anything, reqBody.UserIds, users.UserStatus("active")).Return(nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/users/batch-update-status", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestUserAdminAPI_BatchUpdateStatus_EmptyUserIds(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	reqBody := BatchUpdateStatusRequest{
		UserIds: []string{},
		Status:  "active",
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/users/batch-update-status", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserAdminAPI_BatchUpdateStatus_MissingRequiredFields(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	reqBody := map[string]interface{}{
		"userIds": []string{"user1"},
		// Missing status
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/users/batch-update-status", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==================== BatchDeleteUsers Tests ====================

func TestUserAdminAPI_BatchDeleteUsers_Success(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	reqBody := BatchDeleteUsersRequest{
		UserIds: []string{"user1", "user2"},
	}

	mockService.On("BatchDeleteUsers", mock.Anything, reqBody.UserIds).Return(nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/users/batch-delete", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestUserAdminAPI_BatchDeleteUsers_EmptyUserIds(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	reqBody := BatchDeleteUsersRequest{
		UserIds: []string{},
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/admin/users/batch-delete", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==================== SearchUsers Tests ====================

func TestUserAdminAPI_SearchUsers_Success(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	keyword := "testuser"
	expectedUsers := []*users.User{
		func() *users.User {
			u := &users.User{Username: "testuser1", Email: "test1@example.com"}
			u.ID = primitive.NewObjectID()
			return u
		}(),
	}

	mockService.On("SearchUsers", mock.Anything, keyword, 1, 20).
		Return(expectedUsers, int64(1), nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/users/search?keyword="+keyword, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestUserAdminAPI_SearchUsers_EmptyKeyword(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/users/search", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==================== CountByStatus Tests ====================

func TestUserAdminAPI_CountByStatus_Success(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	expectedCounts := map[string]int64{
		"active":   100,
		"inactive": 10,
		"banned":   5,
	}

	mockService.On("CountByStatus", mock.Anything).Return(expectedCounts, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/users/count-by-status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])

	mockService.AssertExpectations(t)
}

func TestUserAdminAPI_CountByStatus_ServiceError(t *testing.T) {
	// Given
	mockService := new(MockUserAdminService)
	router := setupUserAdminTestRouter(mockService)

	mockService.On("CountByStatus", mock.Anything).Return(nil, assert.AnError)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/admin/users/count-by-status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}
