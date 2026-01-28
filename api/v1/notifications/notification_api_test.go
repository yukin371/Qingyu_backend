package notifications_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/api/v1/notifications"
	"Qingyu_backend/api/v1/notifications/dto"
	notifModels "Qingyu_backend/models/notification"
	notifService "Qingyu_backend/service/notification"
)

// MockNotificationService 模拟通知服务
type MockNotificationService struct {
	mock.Mock
}

// GetNotifications 模拟获取通知列表
func (m *MockNotificationService) GetNotifications(ctx context.Context, req *notifService.GetNotificationsRequest) (*notifService.GetNotificationsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notifService.GetNotificationsResponse), args.Error(1)
}

// GetNotification 模拟获取单个通知
func (m *MockNotificationService) GetNotification(ctx context.Context, id string) (*notifModels.Notification, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notifModels.Notification), args.Error(1)
}

// MarkAsRead 模拟标记为已读
func (m *MockNotificationService) MarkAsRead(ctx context.Context, id, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

// MarkMultipleAsReadWithResult 模拟批量标记为已读
func (m *MockNotificationService) MarkMultipleAsReadWithResult(ctx context.Context, ids []string, userID string) (int, int, error) {
	args := m.Called(ctx, ids, userID)
	return args.Int(0), args.Int(1), args.Error(2)
}

// MarkAllAsRead 模拟全部标记为已读
func (m *MockNotificationService) MarkAllAsRead(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// DeleteNotification 模拟删除通知
func (m *MockNotificationService) DeleteNotification(ctx context.Context, id, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

// BatchDeleteNotificationsWithResult 模拟批量删除通知
func (m *MockNotificationService) BatchDeleteNotificationsWithResult(ctx context.Context, ids []string, userID string) (int, int, error) {
	args := m.Called(ctx, ids, userID)
	return args.Int(0), args.Int(1), args.Error(2)
}

// BatchDeleteNotifications 模拟批量删除通知
func (m *MockNotificationService) BatchDeleteNotifications(ctx context.Context, ids []string, userID string) error {
	args := m.Called(ctx, ids, userID)
	return args.Error(0)
}

// DeleteAllNotifications 模拟删除所有通知
func (m *MockNotificationService) DeleteAllNotifications(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// GetUnreadCount 模拟获取未读数量
func (m *MockNotificationService) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// GetNotificationStats 模拟获取通知统计
func (m *MockNotificationService) GetNotificationStats(ctx context.Context, userID string) (*notifModels.NotificationStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notifModels.NotificationStats), args.Error(1)
}

// GetNotificationPreference 模拟获取通知偏好设置
func (m *MockNotificationService) GetNotificationPreference(ctx context.Context, userID string) (*notifModels.NotificationPreference, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notifModels.NotificationPreference), args.Error(1)
}

// UpdateNotificationPreference 模拟更新通知偏好设置
func (m *MockNotificationService) UpdateNotificationPreference(ctx context.Context, userID string, req *notifService.UpdateNotificationPreferenceRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

// ResetNotificationPreference 模拟重置通知偏好设置
func (m *MockNotificationService) ResetNotificationPreference(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// GetEmailNotificationSettings 模拟获取邮件通知设置
func (m *MockNotificationService) GetEmailNotificationSettings(ctx context.Context, userID string) (*notifModels.EmailNotificationSettings, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notifModels.EmailNotificationSettings), args.Error(1)
}

// UpdateEmailNotificationSettings 模拟更新邮件通知设置
func (m *MockNotificationService) UpdateEmailNotificationSettings(ctx context.Context, userID string, settings *notifModels.EmailNotificationSettings) error {
	args := m.Called(ctx, userID, settings)
	return args.Error(0)
}

// GetSMSNotificationSettings 模拟获取短信通知设置
func (m *MockNotificationService) GetSMSNotificationSettings(ctx context.Context, userID string) (*notifModels.SMSNotificationSettings, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notifModels.SMSNotificationSettings), args.Error(1)
}

// UpdateSMSNotificationSettings 模拟更新短信通知设置
func (m *MockNotificationService) UpdateSMSNotificationSettings(ctx context.Context, userID string, settings *notifModels.SMSNotificationSettings) error {
	args := m.Called(ctx, userID, settings)
	return args.Error(0)
}

// RegisterPushDevice 模拟注册推送设备
func (m *MockNotificationService) RegisterPushDevice(ctx context.Context, req *notifService.RegisterPushDeviceRequest) (*notifModels.PushDevice, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notifModels.PushDevice), args.Error(1)
}

// UnregisterPushDevice 模拟取消注册推送设备
func (m *MockNotificationService) UnregisterPushDevice(ctx context.Context, deviceID, userID string) error {
	args := m.Called(ctx, deviceID, userID)
	return args.Error(0)
}

// GetUserPushDevices 模拟获取用户推送设备
func (m *MockNotificationService) GetUserPushDevices(ctx context.Context, userID string) ([]*notifModels.PushDevice, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*notifModels.PushDevice), args.Error(1)
}

// ClearReadNotifications 模拟清除已读通知
func (m *MockNotificationService) ClearReadNotifications(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// ResendNotification 模拟重新发送通知
func (m *MockNotificationService) ResendNotification(ctx context.Context, id, userID, method string) error {
	args := m.Called(ctx, id, userID, method)
	return args.Error(0)
}

// CreateNotification 模拟创建通知
func (m *MockNotificationService) CreateNotification(ctx context.Context, req *notifService.CreateNotificationRequest) (*notifModels.Notification, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notifModels.Notification), args.Error(1)
}

// MarkMultipleAsRead 模拟批量标记为已读
func (m *MockNotificationService) MarkMultipleAsRead(ctx context.Context, ids []string, userID string) error {
	args := m.Called(ctx, ids, userID)
	return args.Error(0)
}

// SendNotification 模拟发送通知
func (m *MockNotificationService) SendNotification(ctx context.Context, userID string, notificationType notifModels.NotificationType, title, content string, data map[string]interface{}) error {
	args := m.Called(ctx, userID, notificationType, title, content, data)
	return args.Error(0)
}

// SendNotificationWithTemplate 模拟使用模板发送通知
func (m *MockNotificationService) SendNotificationWithTemplate(ctx context.Context, userID string, notificationType notifModels.NotificationType, action string, variables map[string]interface{}) error {
	args := m.Called(ctx, userID, notificationType, action, variables)
	return args.Error(0)
}

// BatchSendNotification 模拟批量发送通知
func (m *MockNotificationService) BatchSendNotification(ctx context.Context, userIDs []string, notificationType notifModels.NotificationType, title, content string, data map[string]interface{}) error {
	args := m.Called(ctx, userIDs, notificationType, title, content, data)
	return args.Error(0)
}

// UpdatePushDeviceToken 模拟更新推送设备token
func (m *MockNotificationService) UpdatePushDeviceToken(ctx context.Context, deviceID, userID, deviceToken string) error {
	args := m.Called(ctx, deviceID, userID, deviceToken)
	return args.Error(0)
}

// CleanupExpiredNotifications 模拟清理过期通知
func (m *MockNotificationService) CleanupExpiredNotifications(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// CleanupOldNotifications 模拟清理旧通知
func (m *MockNotificationService) CleanupOldNotifications(ctx context.Context, days int) (int64, error) {
	args := m.Called(ctx, days)
	return args.Get(0).(int64), args.Error(1)
}

// setupNotificationTestRouter 设置通知测试路由
func setupNotificationTestRouter(service *MockNotificationService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 添加middleware来设置user_id
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	})

	api := notifications.NewNotificationAPI(service)

	v1 := r.Group("/api/v1/notifications")
	{
		v1.GET("", api.GetNotifications)
		v1.GET("/:id", api.GetNotification)
		v1.POST("/:id/read", api.MarkAsRead)
		v1.POST("/batch-read", api.MarkMultipleAsRead)
		v1.POST("/read-all", api.MarkAllAsRead)
		v1.DELETE("/:id", api.DeleteNotification)
		v1.POST("/batch-delete", api.BatchDeleteNotifications)
		v1.DELETE("/delete-all", api.DeleteAllNotifications)
		v1.GET("/unread-count", api.GetUnreadCount)
		v1.GET("/stats", api.GetNotificationStats)
		v1.GET("/preferences", api.GetNotificationPreference)
		v1.PUT("/preferences", api.UpdateNotificationPreference)
		v1.POST("/preferences/reset", api.ResetNotificationPreference)
		v1.POST("/clear-read", api.ClearReadNotifications)
		v1.POST("/:id/resend", api.ResendNotification)
		v1.GET("/ws-endpoint", api.GetWSEndpoint)
	}

	return r
}

// =========================
// MarkAsRead 测试
// =========================

// TestNotificationAPI_MarkAsRead_MissingUserID 测试缺少用户ID
func TestNotificationAPI_MarkAsRead_MissingUserID(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	router := setupNotificationTestRouter(mockService, "") // 不设置userID

	notificationID := primitive.NewObjectID().Hex()
	req, _ := http.NewRequest("POST", "/api/v1/notifications/"+notificationID+"/read", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestNotificationAPI_MarkAsRead_MissingNotificationID 测试缺少通知ID
func TestNotificationAPI_MarkAsRead_MissingNotificationID(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	userID := primitive.NewObjectID().Hex()
	router := setupNotificationTestRouter(mockService, userID)

	req, _ := http.NewRequest("POST", "/api/v1/notifications//read", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestNotificationAPI_MarkAsRead_Success 测试标记已读成功
func TestNotificationAPI_MarkAsRead_Success(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	userID := primitive.NewObjectID().Hex()
	router := setupNotificationTestRouter(mockService, userID)

	notificationID := primitive.NewObjectID().Hex()
	mockService.On("MarkAsRead", mock.Anything, notificationID, userID).Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/notifications/"+notificationID+"/read", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// =========================
// BatchMarkRead 测试
// =========================

// TestNotificationAPI_BatchMarkRead_MissingUserID 测试缺少用户ID
func TestNotificationAPI_BatchMarkRead_MissingUserID(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	router := setupNotificationTestRouter(mockService, "")

	body := dto.BatchMarkReadRequest{
		NotificationIDs: []string{primitive.NewObjectID().Hex()},
		ReadAt:          time.Now().Unix(),
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/notifications/batch-read", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestNotificationAPI_BatchMarkRead_EmptyNotificationIDs 测试空的通知ID列表
func TestNotificationAPI_BatchMarkRead_EmptyNotificationIDs(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	userID := primitive.NewObjectID().Hex()
	router := setupNotificationTestRouter(mockService, userID)

	body := dto.BatchMarkReadRequest{
		NotificationIDs: []string{},
		ReadAt:          time.Now().Unix(),
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/notifications/batch-read", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestNotificationAPI_BatchMarkRead_MissingReadAt 测试缺少read_at
func TestNotificationAPI_BatchMarkRead_MissingReadAt(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	userID := primitive.NewObjectID().Hex()
	router := setupNotificationTestRouter(mockService, userID)

	body := map[string]interface{}{
		"notification_ids": []string{primitive.NewObjectID().Hex()},
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/notifications/batch-read", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestNotificationAPI_BatchMarkRead_Success 测试批量标记成功
func TestNotificationAPI_BatchMarkRead_Success(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	userID := primitive.NewObjectID().Hex()
	router := setupNotificationTestRouter(mockService, userID)

	notificationIDs := []string{primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex()}
	body := struct {
		IDs []string `json:"ids"`
	}{
		IDs: notificationIDs,
	}
	bodyBytes, _ := json.Marshal(body)

	mockService.On("MarkMultipleAsRead", mock.Anything, notificationIDs, userID).Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/notifications/batch-read", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// =========================
// MarkAllAsRead 测试
// =========================

// TestNotificationAPI_MarkAllAsRead_MissingUserID 测试缺少用户ID
func TestNotificationAPI_MarkAllAsRead_MissingUserID(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	router := setupNotificationTestRouter(mockService, "")

	req, _ := http.NewRequest("POST", "/api/v1/notifications/read-all", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestNotificationAPI_MarkAllAsRead_Success 测试全部标记成功
func TestNotificationAPI_MarkAllAsRead_Success(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	userID := primitive.NewObjectID().Hex()
	router := setupNotificationTestRouter(mockService, userID)

	mockService.On("MarkAllAsRead", mock.Anything, userID).Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/notifications/read-all", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// =========================
// DeleteNotification 测试
// =========================

// TestNotificationAPI_DeleteNotification_MissingUserID 测试缺少用户ID
func TestNotificationAPI_DeleteNotification_MissingUserID(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	router := setupNotificationTestRouter(mockService, "")

	notificationID := primitive.NewObjectID().Hex()
	req, _ := http.NewRequest("DELETE", "/api/v1/notifications/"+notificationID, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestNotificationAPI_DeleteNotification_MissingNotificationID 测试缺少通知ID
func TestNotificationAPI_DeleteNotification_MissingNotificationID(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	userID := primitive.NewObjectID().Hex()
	router := setupNotificationTestRouter(mockService, userID)

	req, _ := http.NewRequest("DELETE", "/api/v1/notifications/", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestNotificationAPI_DeleteNotification_Success 测试删除成功
func TestNotificationAPI_DeleteNotification_Success(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	userID := primitive.NewObjectID().Hex()
	router := setupNotificationTestRouter(mockService, userID)

	notificationID := primitive.NewObjectID().Hex()
	mockService.On("DeleteNotification", mock.Anything, notificationID, userID).Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/v1/notifications/"+notificationID, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// =========================
// BatchDeleteNotifications 测试
// =========================

// TestNotificationAPI_BatchDeleteNotifications_MissingUserID 测试缺少用户ID
func TestNotificationAPI_BatchDeleteNotifications_MissingUserID(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	router := setupNotificationTestRouter(mockService, "")

	body := dto.BatchDeleteRequest{
		NotificationIDs: []string{primitive.NewObjectID().Hex()},
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/notifications/batch-delete", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestNotificationAPI_BatchDeleteNotifications_EmptyNotificationIDs 测试空的通知ID列表
func TestNotificationAPI_BatchDeleteNotifications_EmptyNotificationIDs(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	userID := primitive.NewObjectID().Hex()
	router := setupNotificationTestRouter(mockService, userID)

	body := dto.BatchDeleteRequest{
		NotificationIDs: []string{},
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/notifications/batch-delete", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestNotificationAPI_BatchDeleteNotifications_Success 测试批量删除成功
func TestNotificationAPI_BatchDeleteNotifications_Success(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	userID := primitive.NewObjectID().Hex()
	router := setupNotificationTestRouter(mockService, userID)

	notificationIDs := []string{primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex()}
	body := struct {
		IDs []string `json:"ids"`
	}{
		IDs: notificationIDs,
	}
	bodyBytes, _ := json.Marshal(body)

	mockService.On("BatchDeleteNotifications", mock.Anything, notificationIDs, userID).Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/notifications/batch-delete", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// =========================
// ClearReadNotifications 测试
// =========================

// TestNotificationAPI_ClearReadNotifications_MissingUserID 测试缺少用户ID
func TestNotificationAPI_ClearReadNotifications_MissingUserID(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	router := setupNotificationTestRouter(mockService, "")

	req, _ := http.NewRequest("POST", "/api/v1/notifications/clear-read", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestNotificationAPI_ClearReadNotifications_Success 测试清除成功
func TestNotificationAPI_ClearReadNotifications_Success(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	userID := primitive.NewObjectID().Hex()
	router := setupNotificationTestRouter(mockService, userID)

	mockService.On("ClearReadNotifications", mock.Anything, userID).Return(int64(10), nil)

	req, _ := http.NewRequest("POST", "/api/v1/notifications/clear-read", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// =========================
// ResendNotification 测试
// =========================

// TestNotificationAPI_ResendNotification_MissingUserID 测试缺少用户ID
func TestNotificationAPI_ResendNotification_MissingUserID(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	router := setupNotificationTestRouter(mockService, "")

	notificationID := primitive.NewObjectID().Hex()
	body := map[string]string{"method": "email"}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/notifications/"+notificationID+"/resend", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestNotificationAPI_ResendNotification_MissingNotificationID 测试缺少通知ID
func TestNotificationAPI_ResendNotification_MissingNotificationID(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	userID := primitive.NewObjectID().Hex()
	router := setupNotificationTestRouter(mockService, userID)

	body := map[string]string{"method": "email"}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/notifications//resend", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestNotificationAPI_ResendNotification_MissingMethod 测试缺少method
func TestNotificationAPI_ResendNotification_MissingMethod(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	userID := primitive.NewObjectID().Hex()
	router := setupNotificationTestRouter(mockService, userID)

	notificationID := primitive.NewObjectID().Hex()
	body := map[string]string{}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/notifications/"+notificationID+"/resend", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestNotificationAPI_ResendNotification_InvalidMethod 测试无效的method
func TestNotificationAPI_ResendNotification_InvalidMethod(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	userID := primitive.NewObjectID().Hex()
	router := setupNotificationTestRouter(mockService, userID)

	notificationID := primitive.NewObjectID().Hex()
	body := map[string]string{"method": "invalid"}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/notifications/"+notificationID+"/resend", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestNotificationAPI_ResendNotification_Success 测试重新发送成功
func TestNotificationAPI_ResendNotification_Success(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	userID := primitive.NewObjectID().Hex()
	router := setupNotificationTestRouter(mockService, userID)

	notificationID := primitive.NewObjectID().Hex()
	body := map[string]string{"method": "email"}
	bodyBytes, _ := json.Marshal(body)

	mockService.On("ResendNotification", mock.Anything, notificationID, userID, "email").Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/notifications/"+notificationID+"/resend", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// TestNotificationAPI_ResendNotification_PushMethod 测试使用push方法
func TestNotificationAPI_ResendNotification_PushMethod(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	userID := primitive.NewObjectID().Hex()
	router := setupNotificationTestRouter(mockService, userID)

	notificationID := primitive.NewObjectID().Hex()
	body := map[string]string{"method": "push"}
	bodyBytes, _ := json.Marshal(body)

	mockService.On("ResendNotification", mock.Anything, notificationID, userID, "push").Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/notifications/"+notificationID+"/resend", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// =========================
// GetWSEndpoint 测试
// =========================

// TestNotificationAPI_GetWSEndpoint_MissingUserID 测试缺少用户ID
func TestNotificationAPI_GetWSEndpoint_MissingUserID(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	router := setupNotificationTestRouter(mockService, "")

	req, _ := http.NewRequest("GET", "/api/v1/notifications/ws-endpoint", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestNotificationAPI_GetWSEndpoint_MissingToken 测试缺少token
func TestNotificationAPI_GetWSEndpoint_MissingToken(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	userID := primitive.NewObjectID().Hex()
	router := setupNotificationTestRouter(mockService, userID)

	req, _ := http.NewRequest("GET", "/api/v1/notifications/ws-endpoint", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestNotificationAPI_GetWSEndpoint_Success 测试获取成功
func TestNotificationAPI_GetWSEndpoint_Success(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	userID := primitive.NewObjectID().Hex()
	router := setupNotificationTestRouter(mockService, userID)

	token := "test-token-123"
	req, _ := http.NewRequest("GET", "/api/v1/notifications/ws-endpoint", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	assert.Contains(t, data["url"], "ws://")
	assert.Contains(t, data["url"], "token="+token)
}

// TestNotificationAPI_GetWSEndpoint_HTTPS 测试HTTPS环境
func TestNotificationAPI_GetWSEndpoint_HTTPS(t *testing.T) {
	// Given
	mockService := new(MockNotificationService)
	userID := primitive.NewObjectID().Hex()
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		// 模拟HTTPS
		c.Request.TLS = &tls.ConnectionState{}
		c.Next()
	})

	api := notifications.NewNotificationAPI(mockService)
	r.GET("/api/v1/notifications/ws-endpoint", api.GetWSEndpoint)

	token := "test-token-123"
	req, _ := http.NewRequest("GET", "/api/v1/notifications/ws-endpoint", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// When
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	assert.Contains(t, data["url"], "wss://")
}
