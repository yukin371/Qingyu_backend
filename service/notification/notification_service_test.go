package notification

import (
	"context"
	"testing"
	"time"

	notifModel "Qingyu_backend/models/notification"
	notifRepo "Qingyu_backend/repository/interfaces/notification"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// =========================
// Mock Repository实现
// =========================

// MockNotificationRepository Mock通知仓储
type MockNotificationRepository struct {
	mock.Mock
}

func (m *MockNotificationRepository) Create(ctx context.Context, notif *notifModel.Notification) error {
	args := m.Called(ctx, notif)
	return args.Error(0)
}

func (m *MockNotificationRepository) GetByID(ctx context.Context, id string) (*notifModel.Notification, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notifModel.Notification), args.Error(1)
}

func (m *MockNotificationRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockNotificationRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockNotificationRepository) List(ctx context.Context, filter *notifModel.NotificationFilter) ([]*notifModel.Notification, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*notifModel.Notification), args.Error(1)
}

func (m *MockNotificationRepository) Count(ctx context.Context, filter *notifModel.NotificationFilter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockNotificationRepository) CountUnread(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockNotificationRepository) GetStats(ctx context.Context, userID string) (*notifModel.NotificationStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notifModel.NotificationStats), args.Error(1)
}

func (m *MockNotificationRepository) BatchMarkAsRead(ctx context.Context, ids []string) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockNotificationRepository) MarkAllAsReadForUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockNotificationRepository) BatchDelete(ctx context.Context, ids []string) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockNotificationRepository) DeleteAllForUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockNotificationRepository) DeleteExpired(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockNotificationRepository) DeleteOldNotifications(ctx context.Context, beforeDate time.Time) (int64, error) {
	args := m.Called(ctx, beforeDate)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockNotificationRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockNotificationRepository) GetUnreadByType(ctx context.Context, userID string, notificationType notifModel.NotificationType) ([]*notifModel.Notification, error) {
	args := m.Called(ctx, userID, notificationType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*notifModel.Notification), args.Error(1)
}

// MockNotificationPreferenceRepository Mock通知偏好仓储
type MockNotificationPreferenceRepository struct {
	mock.Mock
}

func (m *MockNotificationPreferenceRepository) Create(ctx context.Context, pref *notifModel.NotificationPreference) error {
	args := m.Called(ctx, pref)
	return args.Error(0)
}

func (m *MockNotificationPreferenceRepository) GetByID(ctx context.Context, id string) (*notifModel.NotificationPreference, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notifModel.NotificationPreference), args.Error(1)
}

func (m *MockNotificationPreferenceRepository) GetByUserID(ctx context.Context, userID string) (*notifModel.NotificationPreference, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notifModel.NotificationPreference), args.Error(1)
}

func (m *MockNotificationPreferenceRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockNotificationPreferenceRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockNotificationPreferenceRepository) Exists(ctx context.Context, userID string) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockNotificationPreferenceRepository) BatchUpdate(ctx context.Context, ids []string, updates map[string]interface{}) error {
	args := m.Called(ctx, ids, updates)
	return args.Error(0)
}

// MockPushDeviceRepository Mock推送设备仓储
type MockPushDeviceRepository struct {
	mock.Mock
}

func (m *MockPushDeviceRepository) Create(ctx context.Context, device *notifModel.PushDevice) error {
	args := m.Called(ctx, device)
	return args.Error(0)
}

func (m *MockPushDeviceRepository) GetByID(ctx context.Context, id string) (*notifModel.PushDevice, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notifModel.PushDevice), args.Error(1)
}

func (m *MockPushDeviceRepository) GetByDeviceID(ctx context.Context, deviceID string) (*notifModel.PushDevice, error) {
	args := m.Called(ctx, deviceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notifModel.PushDevice), args.Error(1)
}

func (m *MockPushDeviceRepository) GetActiveByUserID(ctx context.Context, userID string) ([]*notifModel.PushDevice, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*notifModel.PushDevice), args.Error(1)
}

func (m *MockPushDeviceRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockPushDeviceRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPushDeviceRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockPushDeviceRepository) GetByUserID(ctx context.Context, userID string) ([]*notifModel.PushDevice, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*notifModel.PushDevice), args.Error(1)
}

func (m *MockPushDeviceRepository) BatchDelete(ctx context.Context, ids []string) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockPushDeviceRepository) DeactivateAllForUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockPushDeviceRepository) UpdateLastUsed(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPushDeviceRepository) DeleteInactiveDevices(ctx context.Context, beforeDate time.Time) (int64, error) {
	args := m.Called(ctx, beforeDate)
	return args.Get(0).(int64), args.Error(1)
}

// MockNotificationTemplateRepository Mock通知模板仓储
type MockNotificationTemplateRepository struct {
	mock.Mock
}

func (m *MockNotificationTemplateRepository) GetActiveTemplate(ctx context.Context, notifType notifModel.NotificationType, action, language string) (*notifModel.NotificationTemplate, error) {
	args := m.Called(ctx, notifType, action, language)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notifModel.NotificationTemplate), args.Error(1)
}

func (m *MockNotificationTemplateRepository) Create(ctx context.Context, template *notifModel.NotificationTemplate) error {
	args := m.Called(ctx, template)
	return args.Error(0)
}

func (m *MockNotificationTemplateRepository) GetByID(ctx context.Context, id string) (*notifModel.NotificationTemplate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notifModel.NotificationTemplate), args.Error(1)
}

func (m *MockNotificationTemplateRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockNotificationTemplateRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockNotificationTemplateRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockNotificationTemplateRepository) List(ctx context.Context, filter *notifRepo.TemplateFilter) ([]*notifModel.NotificationTemplate, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*notifModel.NotificationTemplate), args.Error(1)
}

func (m *MockNotificationTemplateRepository) GetByTypeAndAction(ctx context.Context, templateType notifModel.NotificationType, action string) ([]*notifModel.NotificationTemplate, error) {
	args := m.Called(ctx, templateType, action)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*notifModel.NotificationTemplate), args.Error(1)
}

// =========================
// 测试辅助函数
// =========================

// setupNotificationService 创建测试用的NotificationService实例
func setupNotificationService() (*notificationServiceImpl, *MockNotificationRepository, *MockNotificationPreferenceRepository, *MockPushDeviceRepository, *MockNotificationTemplateRepository) {
	mockNotifRepo := new(MockNotificationRepository)
	mockPrefRepo := new(MockNotificationPreferenceRepository)
	mockDeviceRepo := new(MockPushDeviceRepository)
	mockTemplateRepo := new(MockNotificationTemplateRepository)

	service := NewNotificationService(mockNotifRepo, mockPrefRepo, mockDeviceRepo, mockTemplateRepo)

	return service.(*notificationServiceImpl), mockNotifRepo, mockPrefRepo, mockDeviceRepo, mockTemplateRepo
}

// =========================
// 创建通知相关测试
// =========================

// TestNotificationService_CreateNotification_Success 测试创建通知成功
func TestNotificationService_CreateNotification_Success(t *testing.T) {
	// Arrange
	service, mockNotifRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	req := &CreateNotificationRequest{
		UserID:   "user123",
		Type:     notifModel.NotificationTypeSystem,
		Title:    "测试通知",
		Content:  "这是一条测试通知",
		Priority: notifModel.NotificationPriorityNormal,
	}

	mockNotifRepo.On("Create", ctx, mock.Anything).Return(nil)

	// Act
	notif, err := service.CreateNotification(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, notif)
	assert.Equal(t, req.UserID, notif.UserID)
	assert.Equal(t, req.Type, notif.Type)
	assert.Equal(t, req.Title, notif.Title)
	assert.Equal(t, req.Content, notif.Content)
	assert.Equal(t, req.Priority, notif.Priority)
	assert.False(t, notif.Read)

	mockNotifRepo.AssertExpectations(t)
}

// TestNotificationService_CreateNotification_WithExpiration 测试创建带过期时间的通知
func TestNotificationService_CreateNotification_WithExpiration(t *testing.T) {
	// Arrange
	service, mockNotifRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	expiresIn := 3600 // 1小时
	req := &CreateNotificationRequest{
		UserID:    "user123",
		Type:      notifModel.NotificationTypeSystem,
		Title:     "测试通知",
		Content:   "这是一条测试通知",
		ExpiresIn: &expiresIn,
	}

	mockNotifRepo.On("Create", ctx, mock.MatchedBy(func(n *notifModel.Notification) bool {
		return n.ExpiresAt != nil && n.ExpiresAt.After(time.Now())
	})).Return(nil)

	// Act
	notif, err := service.CreateNotification(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, notif)
	assert.NotNil(t, notif.ExpiresAt)

	mockNotifRepo.AssertExpectations(t)
}

// TestNotificationService_CreateNotification_RepositoryError 测试创建通知-仓储错误
func TestNotificationService_CreateNotification_RepositoryError(t *testing.T) {
	// Arrange
	service, mockNotifRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	req := &CreateNotificationRequest{
		UserID:  "user123",
		Type:    notifModel.NotificationTypeSystem,
		Title:   "测试通知",
		Content: "这是一条测试通知",
	}

	mockNotifRepo.On("Create", ctx, mock.Anything).Return(assert.AnError)

	// Act
	notif, err := service.CreateNotification(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, notif)
	assert.Contains(t, err.Error(), "创建通知失败")

	mockNotifRepo.AssertExpectations(t)
}

// =========================
// 获取通知相关测试
// =========================

// TestNotificationService_GetNotification_Success 测试获取通知成功
func TestNotificationService_GetNotification_Success(t *testing.T) {
	// Arrange
	service, mockNotifRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	notifID := primitive.NewObjectID().Hex()
	expectedNotif := &notifModel.Notification{
		ID:      notifID,
		UserID:  "user123",
		Title:   "测试通知",
		Content: "测试内容",
		Read:    false,
	}

	mockNotifRepo.On("GetByID", ctx, notifID).Return(expectedNotif, nil)

	// Act
	notif, err := service.GetNotification(ctx, notifID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedNotif, notif)

	mockNotifRepo.AssertExpectations(t)
}

// TestNotificationService_GetNotification_NotFound 测试获取通知-不存在
func TestNotificationService_GetNotification_NotFound(t *testing.T) {
	// Arrange
	service, mockNotifRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	notifID := primitive.NewObjectID().Hex()
	mockNotifRepo.On("GetByID", ctx, notifID).Return(nil, nil)

	// Act
	notif, err := service.GetNotification(ctx, notifID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, notif)
	assert.Contains(t, err.Error(), "not found")

	mockNotifRepo.AssertExpectations(t)
}

// TestNotificationService_GetNotifications_Success 测试获取通知列表成功
func TestNotificationService_GetNotifications_Success(t *testing.T) {
	// Arrange
	service, mockNotifRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	req := &GetNotificationsRequest{
		UserID: "user123",
		Limit:  20,
		Offset: 0,
	}

	expectedNotifs := []*notifModel.Notification{
		{ID: "notif1", UserID: "user123", Title: "通知1"},
		{ID: "notif2", UserID: "user123", Title: "通知2"},
	}

	mockNotifRepo.On("List", ctx, mock.Anything).Return(expectedNotifs, nil)
	mockNotifRepo.On("Count", ctx, mock.Anything).Return(int64(2), nil)
	mockNotifRepo.On("CountUnread", ctx, "user123").Return(int64(1), nil)

	// Act
	resp, err := service.GetNotifications(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.Len(t, resp.Notifications, 2)
	assert.Equal(t, int64(2), resp.Total)
	assert.Equal(t, int64(1), resp.UnreadCount)

	mockNotifRepo.AssertExpectations(t)
}

// TestNotificationService_GetNotifications_DefaultLimit 测试获取通知列表-使用默认限制
func TestNotificationService_GetNotifications_DefaultLimit(t *testing.T) {
	// Arrange
	service, mockNotifRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	req := &GetNotificationsRequest{
		UserID: "user123",
		Limit:  0, // 使用默认值
	}

	mockNotifRepo.On("List", ctx, mock.MatchedBy(func(f *notifModel.NotificationFilter) bool {
		return f.Limit == 20
	})).Return([]*notifModel.Notification{}, nil)
	mockNotifRepo.On("Count", ctx, mock.Anything).Return(int64(0), nil)
	mockNotifRepo.On("CountUnread", ctx, "user123").Return(int64(0), nil)

	// Act
	resp, err := service.GetNotifications(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)

	mockNotifRepo.AssertExpectations(t)
}

// TestNotificationService_GetNotifications_MaxLimit 测试获取通知列表-超过最大限制
func TestNotificationService_GetNotifications_MaxLimit(t *testing.T) {
	// Arrange
	service, mockNotifRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	req := &GetNotificationsRequest{
		UserID: "user123",
		Limit:  200, // 超过最大值100
	}

	mockNotifRepo.On("List", ctx, mock.MatchedBy(func(f *notifModel.NotificationFilter) bool {
		return f.Limit == 100
	})).Return([]*notifModel.Notification{}, nil)
	mockNotifRepo.On("Count", ctx, mock.Anything).Return(int64(0), nil)
	mockNotifRepo.On("CountUnread", ctx, "user123").Return(int64(0), nil)

	// Act
	resp, err := service.GetNotifications(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)

	mockNotifRepo.AssertExpectations(t)
}

// =========================
// 标记通知为已读相关测试
// =========================

// TestNotificationService_MarkAsRead_Success 测试标记通知为已读成功
func TestNotificationService_MarkAsRead_Success(t *testing.T) {
	// Arrange
	service, mockNotifRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	notifID := primitive.NewObjectID().Hex()
	userID := "user123"

	notif := &notifModel.Notification{
		ID:     notifID,
		UserID: userID,
		Read:   false,
	}

	mockNotifRepo.On("GetByID", ctx, notifID).Return(notif, nil)
	mockNotifRepo.On("Update", ctx, notifID, mock.MatchedBy(func(u map[string]interface{}) bool {
		return u["read"] == true && u["read_at"] != nil
	})).Return(nil)

	// Act
	err := service.MarkAsRead(ctx, notifID, userID)

	// Assert
	require.NoError(t, err)

	mockNotifRepo.AssertExpectations(t)
}

// TestNotificationService_MarkAsRead_AlreadyRead 测试标记已读通知-已经已读
func TestNotificationService_MarkAsRead_AlreadyRead(t *testing.T) {
	// Arrange
	service, mockNotifRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	notifID := primitive.NewObjectID().Hex()
	userID := "user123"

	notif := &notifModel.Notification{
		ID:     notifID,
		UserID: userID,
		Read:   true, // 已经已读
	}

	mockNotifRepo.On("GetByID", ctx, notifID).Return(notif, nil)

	// Act
	err := service.MarkAsRead(ctx, notifID, userID)

	// Assert
	require.NoError(t, err)

	// 验证没有调用Update
	mockNotifRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)

	mockNotifRepo.AssertExpectations(t)
}

// TestNotificationService_MarkAsRead_NotFound 测试标记通知为已读-通知不存在
func TestNotificationService_MarkAsRead_NotFound(t *testing.T) {
	// Arrange
	service, mockNotifRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	notifID := primitive.NewObjectID().Hex()
	userID := "user123"

	mockNotifRepo.On("GetByID", ctx, notifID).Return(nil, nil)

	// Act
	err := service.MarkAsRead(ctx, notifID, userID)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	mockNotifRepo.AssertExpectations(t)
}

// TestNotificationService_MarkAsRead_NoPermission 测试标记通知为已读-无权限
func TestNotificationService_MarkAsRead_NoPermission(t *testing.T) {
	// Arrange
	service, mockNotifRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	notifID := primitive.NewObjectID().Hex()
	userID := "user123"
	ownerID := "user456"

	notif := &notifModel.Notification{
		ID:     notifID,
		UserID: ownerID, // 不是当前用户的通知
		Read:   false,
	}

	mockNotifRepo.On("GetByID", ctx, notifID).Return(notif, nil)

	// Act
	err := service.MarkAsRead(ctx, notifID, userID)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "无权操作")

	mockNotifRepo.AssertExpectations(t)
}

// TestNotificationService_MarkAllAsRead_Success 测试标记所有通知为已读成功
func TestNotificationService_MarkAllAsRead_Success(t *testing.T) {
	// Arrange
	service, mockNotifRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	userID := "user123"

	mockNotifRepo.On("MarkAllAsReadForUser", ctx, userID).Return(nil)

	// Act
	err := service.MarkAllAsRead(ctx, userID)

	// Assert
	require.NoError(t, err)

	mockNotifRepo.AssertExpectations(t)
}

// TestNotificationService_MarkMultipleAsRead_Success 测试批量标记通知为已读成功
func TestNotificationService_MarkMultipleAsRead_Success(t *testing.T) {
	// Arrange
	service, mockNotifRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	ids := []string{"notif1", "notif2", "notif3"}
	userID := "user123"

	for _, id := range ids {
		mockNotifRepo.On("GetByID", ctx, id).Return(&notifModel.Notification{
			ID:     id,
			UserID: userID,
		}, nil)
	}

	mockNotifRepo.On("BatchMarkAsRead", ctx, ids).Return(nil)

	// Act
	err := service.MarkMultipleAsRead(ctx, ids, userID)

	// Assert
	require.NoError(t, err)

	mockNotifRepo.AssertExpectations(t)
}

// TestNotificationService_MarkMultipleAsRead_EmptyIds 测试批量标记-空ID列表
func TestNotificationService_MarkMultipleAsRead_EmptyIds(t *testing.T) {
	// Arrange
	service, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	// Act
	err := service.MarkMultipleAsRead(ctx, []string{}, "user123")

	// Assert
	require.NoError(t, err)
}

// =========================
// 删除通知相关测试
// =========================

// TestNotificationService_DeleteNotification_Success 测试删除通知成功
func TestNotificationService_DeleteNotification_Success(t *testing.T) {
	// Arrange
	service, mockNotifRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	notifID := primitive.NewObjectID().Hex()
	userID := "user123"

	notif := &notifModel.Notification{
		ID:     notifID,
		UserID: userID,
	}

	mockNotifRepo.On("GetByID", ctx, notifID).Return(notif, nil)
	mockNotifRepo.On("Delete", ctx, notifID).Return(nil)

	// Act
	err := service.DeleteNotification(ctx, notifID, userID)

	// Assert
	require.NoError(t, err)

	mockNotifRepo.AssertExpectations(t)
}

// TestNotificationService_DeleteNotification_NoPermission 测试删除通知-无权限
func TestNotificationService_DeleteNotification_NoPermission(t *testing.T) {
	// Arrange
	service, mockNotifRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	notifID := primitive.NewObjectID().Hex()
	userID := "user123"
	ownerID := "user456"

	notif := &notifModel.Notification{
		ID:     notifID,
		UserID: ownerID,
	}

	mockNotifRepo.On("GetByID", ctx, notifID).Return(notif, nil)

	// Act
	err := service.DeleteNotification(ctx, notifID, userID)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "无权操作")

	mockNotifRepo.AssertExpectations(t)
}

// TestNotificationService_BatchDeleteNotifications_Success 测试批量删除通知成功
func TestNotificationService_BatchDeleteNotifications_Success(t *testing.T) {
	// Arrange
	service, mockNotifRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	ids := []string{"notif1", "notif2"}
	userID := "user123"

	for _, id := range ids {
		mockNotifRepo.On("GetByID", ctx, id).Return(&notifModel.Notification{
			ID:     id,
			UserID: userID,
		}, nil)
	}

	mockNotifRepo.On("BatchDelete", ctx, ids).Return(nil)

	// Act
	err := service.BatchDeleteNotifications(ctx, ids, userID)

	// Assert
	require.NoError(t, err)

	mockNotifRepo.AssertExpectations(t)
}

// TestNotificationService_DeleteAllNotifications_Success 测试删除所有通知成功
func TestNotificationService_DeleteAllNotifications_Success(t *testing.T) {
	// Arrange
	service, mockNotifRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	userID := "user123"

	mockNotifRepo.On("DeleteAllForUser", ctx, userID).Return(nil)

	// Act
	err := service.DeleteAllNotifications(ctx, userID)

	// Assert
	require.NoError(t, err)

	mockNotifRepo.AssertExpectations(t)
}

// =========================
// 发送通知相关测试
// =========================

// TestNotificationService_SendNotification_Success 测试发送通知成功
func TestNotificationService_SendNotification_Success(t *testing.T) {
	// Arrange
	service, mockNotifRepo, mockPrefRepo, _, _ := setupNotificationService()
	ctx := context.Background()

	userID := "user123"

	// Mock偏好设置
	pref := &notifModel.NotificationPreference{
		ID:           "pref123",
		UserID:       userID,
		EnableSystem: true,
	}

	mockPrefRepo.On("GetByUserID", ctx, userID).Return(pref, nil)
	mockNotifRepo.On("Create", ctx, mock.Anything).Return(nil)

	// Act
	err := service.SendNotification(ctx, userID, notifModel.NotificationTypeSystem, "测试标题", "测试内容", nil)

	// Assert
	require.NoError(t, err)

	mockPrefRepo.AssertExpectations(t)
	mockNotifRepo.AssertExpectations(t)
}

// TestNotificationService_SendNotification_TypeDisabled 测试发送通知-类型已禁用
func TestNotificationService_SendNotification_TypeDisabled(t *testing.T) {
	// Arrange
	service, mockNotifRepo, mockPrefRepo, _, _ := setupNotificationService()
	ctx := context.Background()

	userID := "user123"

	// Mock偏好设置-系统通知已禁用
	pref := &notifModel.NotificationPreference{
		ID:           "pref123",
		UserID:       userID,
		EnableSystem: false,
	}

	mockPrefRepo.On("GetByUserID", ctx, userID).Return(pref, nil)

	// Act
	err := service.SendNotification(ctx, userID, notifModel.NotificationTypeSystem, "测试标题", "测试内容", nil)

	// Assert
	require.NoError(t, err) // 不应该报错，只是不发送

	// 验证没有调用Create
	mockNotifRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)

	mockPrefRepo.AssertExpectations(t)
}

// TestNotificationService_SendNotification_CreateDefaultPreference 测试发送通知-创建默认偏好
func TestNotificationService_SendNotification_CreateDefaultPreference(t *testing.T) {
	// Arrange
	service, mockNotifRepo, mockPrefRepo, _, _ := setupNotificationService()
	ctx := context.Background()

	userID := "user123"

	mockPrefRepo.On("GetByUserID", ctx, userID).Return(nil, nil)
	mockPrefRepo.On("Create", ctx, mock.Anything).Return(nil)
	mockNotifRepo.On("Create", ctx, mock.Anything).Return(nil)

	// Act
	err := service.SendNotification(ctx, userID, notifModel.NotificationTypeSystem, "测试标题", "测试内容", nil)

	// Assert
	require.NoError(t, err)

	mockPrefRepo.AssertExpectations(t)
	mockNotifRepo.AssertExpectations(t)
}

// =========================
// 通知偏好设置相关测试
// =========================

// TestNotificationService_GetNotificationPreference_Success 测试获取通知偏好设置成功
func TestNotificationService_GetNotificationPreference_Success(t *testing.T) {
	// Arrange
	service, _, mockPrefRepo, _, _ := setupNotificationService()
	ctx := context.Background()

	userID := "user123"
	expectedPref := &notifModel.NotificationPreference{
		ID:           "pref123",
		UserID:       userID,
		EnableSystem: true,
	}

	mockPrefRepo.On("GetByUserID", ctx, userID).Return(expectedPref, nil)

	// Act
	pref, err := service.GetNotificationPreference(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPref, pref)

	mockPrefRepo.AssertExpectations(t)
}

// TestNotificationService_GetNotificationPreference_Default 测试获取通知偏好设置-返回默认设置
func TestNotificationService_GetNotificationPreference_Default(t *testing.T) {
	// Arrange
	service, _, mockPrefRepo, _, _ := setupNotificationService()
	ctx := context.Background()

	userID := "user123"

	mockPrefRepo.On("GetByUserID", ctx, userID).Return(nil, nil)

	// Act
	pref, err := service.GetNotificationPreference(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, pref)
	assert.Equal(t, userID, pref.UserID)

	mockPrefRepo.AssertExpectations(t)
}

// TestNotificationService_UpdateNotificationPreference_Success 测试更新通知偏好设置成功
func TestNotificationService_UpdateNotificationPreference_Success(t *testing.T) {
	// Arrange
	service, _, mockPrefRepo, _, _ := setupNotificationService()
	ctx := context.Background()

	userID := "user123"
	existingPref := &notifModel.NotificationPreference{
		ID:     "pref123",
		UserID: userID,
	}

	req := &UpdateNotificationPreferenceRequest{
		EnableSystem: boolPtr(false),
		EnableSocial: boolPtr(true),
	}

	mockPrefRepo.On("GetByUserID", ctx, userID).Return(existingPref, nil)
	mockPrefRepo.On("Update", ctx, "pref123", mock.MatchedBy(func(u map[string]interface{}) bool {
		return u["enable_system"] == false && u["enable_social"] == true
	})).Return(nil)

	// Act
	err := service.UpdateNotificationPreference(ctx, userID, req)

	// Assert
	require.NoError(t, err)

	mockPrefRepo.AssertExpectations(t)
}

// TestNotificationService_UpdateNotificationPreference_CreateNew 测试更新通知偏好设置-创建新的
func TestNotificationService_UpdateNotificationPreference_CreateNew(t *testing.T) {
	// Arrange
	service, _, mockPrefRepo, _, _ := setupNotificationService()
	ctx := context.Background()

	userID := "user123"

	req := &UpdateNotificationPreferenceRequest{
		EnableSystem: boolPtr(true),
	}

	mockPrefRepo.On("GetByUserID", ctx, userID).Return(nil, nil)
	mockPrefRepo.On("Create", ctx, mock.Anything).Return(nil)
	mockPrefRepo.On("Update", ctx, mock.Anything, mock.Anything).Return(nil)

	// Act
	err := service.UpdateNotificationPreference(ctx, userID, req)

	// Assert
	require.NoError(t, err)

	mockPrefRepo.AssertExpectations(t)
}

// TestNotificationService_ResetNotificationPreference_Success 测试重置通知偏好设置成功
func TestNotificationService_ResetNotificationPreference_Success(t *testing.T) {
	// Arrange
	service, _, mockPrefRepo, _, _ := setupNotificationService()
	ctx := context.Background()

	userID := "user123"
	existingPref := &notifModel.NotificationPreference{
		ID:     "pref123",
		UserID: userID,
	}

	mockPrefRepo.On("GetByUserID", ctx, userID).Return(existingPref, nil)
	mockPrefRepo.On("Delete", ctx, "pref123").Return(nil)
	mockPrefRepo.On("Create", ctx, mock.Anything).Return(nil)

	// Act
	err := service.ResetNotificationPreference(ctx, userID)

	// Assert
	require.NoError(t, err)

	mockPrefRepo.AssertExpectations(t)
}

// TestNotificationService_ResetNotificationPreference_NoExisting 测试重置通知偏好设置-无现有设置
func TestNotificationService_ResetNotificationPreference_NoExisting(t *testing.T) {
	// Arrange
	service, _, mockPrefRepo, _, _ := setupNotificationService()
	ctx := context.Background()

	userID := "user123"

	mockPrefRepo.On("GetByUserID", ctx, userID).Return(nil, nil)
	mockPrefRepo.On("Create", ctx, mock.Anything).Return(nil)

	// Act
	err := service.ResetNotificationPreference(ctx, userID)

	// Assert
	require.NoError(t, err)

	// 验证没有调用Delete
	mockPrefRepo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)

	mockPrefRepo.AssertExpectations(t)
}

// =========================
// 推送设备管理相关测试
// =========================

// TestNotificationService_RegisterPushDevice_NewDevice 测试注册推送设备-新设备
func TestNotificationService_RegisterPushDevice_NewDevice(t *testing.T) {
	// Arrange
	service, _, _, mockDeviceRepo, _ := setupNotificationService()
	ctx := context.Background()

	req := &RegisterPushDeviceRequest{
		UserID:      "user123",
		DeviceID:    "device123",
		DeviceType:  "ios",
		DeviceToken: "token123",
	}

	mockDeviceRepo.On("GetByDeviceID", ctx, req.DeviceID).Return(nil, nil)
	mockDeviceRepo.On("Create", ctx, mock.Anything).Return(nil)

	// Act
	device, err := service.RegisterPushDevice(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, device)

	mockDeviceRepo.AssertExpectations(t)
}

// TestNotificationService_RegisterPushDevice_ExistingDevice 测试注册推送设备-已存在
func TestNotificationService_RegisterPushDevice_ExistingDevice(t *testing.T) {
	// Arrange
	service, _, _, mockDeviceRepo, _ := setupNotificationService()
	ctx := context.Background()

	req := &RegisterPushDeviceRequest{
		UserID:      "user123",
		DeviceID:    "device123",
		DeviceType:  "ios",
		DeviceToken: "new_token",
	}

	existingDevice := &notifModel.PushDevice{
		ID:          "push123",
		UserID:      "user123",
		DeviceID:    "device123",
		DeviceToken: "old_token",
	}

	mockDeviceRepo.On("GetByDeviceID", ctx, req.DeviceID).Return(existingDevice, nil)
	mockDeviceRepo.On("Update", ctx, "push123", mock.MatchedBy(func(u map[string]interface{}) bool {
		return u["device_token"] == "new_token"
	})).Return(nil)
	mockDeviceRepo.On("GetByID", ctx, "push123").Return(existingDevice, nil)

	// Act
	device, err := service.RegisterPushDevice(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, existingDevice, device)

	mockDeviceRepo.AssertExpectations(t)
}

// TestNotificationService_UnregisterPushDevice_Success 测试取消注册推送设备成功
func TestNotificationService_UnregisterPushDevice_Success(t *testing.T) {
	// Arrange
	service, _, _, mockDeviceRepo, _ := setupNotificationService()
	ctx := context.Background()

	deviceID := "device123"
	userID := "user123"

	device := &notifModel.PushDevice{
		ID:     "push123",
		UserID: userID,
	}

	mockDeviceRepo.On("GetByDeviceID", ctx, deviceID).Return(device, nil)
	mockDeviceRepo.On("Delete", ctx, "push123").Return(nil)

	// Act
	err := service.UnregisterPushDevice(ctx, deviceID, userID)

	// Assert
	require.NoError(t, err)

	mockDeviceRepo.AssertExpectations(t)
}

// TestNotificationService_UnregisterPushDevice_NotFound 测试取消注册推送设备-设备不存在
func TestNotificationService_UnregisterPushDevice_NotFound(t *testing.T) {
	// Arrange
	service, _, _, mockDeviceRepo, _ := setupNotificationService()
	ctx := context.Background()

	deviceID := "device123"
	userID := "user123"

	mockDeviceRepo.On("GetByDeviceID", ctx, deviceID).Return(nil, nil)

	// Act
	err := service.UnregisterPushDevice(ctx, deviceID, userID)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	mockDeviceRepo.AssertExpectations(t)
}

// TestNotificationService_GetUserPushDevices_Success 测试获取用户推送设备列表成功
func TestNotificationService_GetUserPushDevices_Success(t *testing.T) {
	// Arrange
	service, _, _, mockDeviceRepo, _ := setupNotificationService()
	ctx := context.Background()

	userID := "user123"
	expectedDevices := []*notifModel.PushDevice{
		{ID: "device1", UserID: userID},
		{ID: "device2", UserID: userID},
	}

	mockDeviceRepo.On("GetActiveByUserID", ctx, userID).Return(expectedDevices, nil)

	// Act
	devices, err := service.GetUserPushDevices(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.Len(t, devices, 2)

	mockDeviceRepo.AssertExpectations(t)
}

// =========================
// 清理相关测试
// =========================

// TestNotificationService_CleanupExpiredNotifications_Success 测试清理过期通知成功
func TestNotificationService_CleanupExpiredNotifications_Success(t *testing.T) {
	// Arrange
	service, mockNotifRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	mockNotifRepo.On("DeleteExpired", ctx).Return(int64(100), nil)

	// Act
	count, err := service.CleanupExpiredNotifications(ctx)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(100), count)

	mockNotifRepo.AssertExpectations(t)
}

// TestNotificationService_CleanupOldNotifications_Success 测试清理旧通知成功
func TestNotificationService_CleanupOldNotifications_Success(t *testing.T) {
	// Arrange
	service, mockNotifRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	days := 30
	mockNotifRepo.On("DeleteOldNotifications", ctx, mock.Anything).Return(int64(50), nil)

	// Act
	count, err := service.CleanupOldNotifications(ctx, days)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(50), count)

	mockNotifRepo.AssertExpectations(t)
}

// TestNotificationService_CleanupOldNotifications_InvalidDays 测试清理旧通知-无效天数
func TestNotificationService_CleanupOldNotifications_InvalidDays(t *testing.T) {
	// Arrange
	service, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	// Act
	count, err := service.CleanupOldNotifications(ctx, 0)

	// Assert
	require.Error(t, err)
	assert.Equal(t, int64(0), count)
}

// =========================
// 辅助函数
// =========================

// boolPtr 返回bool指针
func boolPtr(b bool) *bool {
	return &b
}
