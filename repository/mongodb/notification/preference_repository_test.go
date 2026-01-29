package notification_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/models/notification"
	notificationRepo "Qingyu_backend/repository/mongodb/notification"
	"Qingyu_backend/test/testutil"
)

// setupPreferenceRepo 测试辅助函数
func setupPreferenceRepo(t *testing.T) (*notificationRepo.NotificationPreferenceRepositoryImpl, context.Context, func()) {
	db, cleanup := testutil.SetupTestDB(t)
	repo := notificationRepo.NewNotificationPreferenceRepository(db).(*notificationRepo.NotificationPreferenceRepositoryImpl)
	ctx := context.Background()
	return repo, ctx, cleanup
}

// TestNotificationPreferenceRepository_Create 测试创建通知偏好设置
func TestNotificationPreferenceRepository_Create(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPreferenceRepo(t)
	defer cleanup()

	preference := &notification.NotificationPreference{
		UserID:            "user123",
		EnableSystem:      true,
		EnableSocial:      true,
		EnableContent:     true,
		EnableReward:      true,
		EnableMessage:     true,
		EnableUpdate:      true,
		EnableMembership:  true,
		EmailNotification: notification.EmailNotificationSettings{
			Enabled:   true,
			Types:     []string{"system", "social"},
			Frequency: "daily",
		},
		SMSNotification: notification.SMSNotificationSettings{
			Enabled: false,
			Types:   []string{},
		},
		PushNotification: true,
		QuietHoursStart:  func() *string { s := "22:00"; return &s }(),
		QuietHoursEnd:    func() *string { s := "08:00"; return &s }(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Act
	err := repo.Create(ctx, preference)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, preference.ID)
}

// TestNotificationPreferenceRepository_Create_WithInvalidID 测试使用无效ID创建偏好设置
func TestNotificationPreferenceRepository_Create_WithInvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPreferenceRepo(t)
	defer cleanup()

	preference := &notification.NotificationPreference{
		ID:               "invalid-id",
		UserID:           "user123",
		EnableSystem:     true,
		EnableSocial:     true,
		EnableContent:    true,
		EnableReward:     true,
		EnableMessage:    true,
		EnableUpdate:     true,
		EnableMembership: true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Act
	err := repo.Create(ctx, preference)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, preference.ID)
	assert.NotEqual(t, "invalid-id", preference.ID)
}

// TestNotificationPreferenceRepository_GetByID 测试根据ID获取通知偏好设置
func TestNotificationPreferenceRepository_GetByID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPreferenceRepo(t)
	defer cleanup()

	preference := &notification.NotificationPreference{
		UserID:           "user456",
		EnableSystem:     true,
		EnableSocial:     false,
		EnableContent:    true,
		EnableReward:     true,
		EnableMessage:    true,
		EnableUpdate:     true,
		EnableMembership: true,
		EmailNotification: notification.EmailNotificationSettings{
			Enabled:   true,
			Types:     []string{"system"},
			Frequency: "immediate",
		},
		SMSNotification: notification.SMSNotificationSettings{
			Enabled: false,
			Types:   []string{},
		},
		PushNotification: false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	err := repo.Create(ctx, preference)
	require.NoError(t, err)

	// Act
	found, err := repo.GetByID(ctx, preference.ID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, preference.UserID, found.UserID)
	assert.Equal(t, preference.EnableSystem, found.EnableSystem)
	assert.Equal(t, preference.EnableSocial, found.EnableSocial)
	assert.Equal(t, preference.EmailNotification.Enabled, found.EmailNotification.Enabled)
}

// TestNotificationPreferenceRepository_GetByID_NotFound 测试获取不存在的偏好设置
func TestNotificationPreferenceRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPreferenceRepo(t)
	defer cleanup()

	// Act
	found, err := repo.GetByID(ctx, "507f1f77bcf86cd799439011")

	// Assert
	require.NoError(t, err)
	assert.Nil(t, found)
}

// TestNotificationPreferenceRepository_GetByID_InvalidID 测试使用无效ID获取偏好设置
func TestNotificationPreferenceRepository_GetByID_InvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPreferenceRepo(t)
	defer cleanup()

	// Act
	found, err := repo.GetByID(ctx, "invalid-id")

	// Assert
	require.Error(t, err)
	assert.Nil(t, found)
	assert.Contains(t, err.Error(), "无效的通知偏好设置ID")
}

// TestNotificationPreferenceRepository_GetByUserID 测试根据用户ID获取通知偏好设置
func TestNotificationPreferenceRepository_GetByUserID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPreferenceRepo(t)
	defer cleanup()

	preference := &notification.NotificationPreference{
		UserID:           "user789",
		EnableSystem:     true,
		EnableSocial:     true,
		EnableContent:    true,
		EnableReward:     true,
		EnableMessage:    true,
		EnableUpdate:     true,
		EnableMembership: true,
		EmailNotification: notification.EmailNotificationSettings{
			Enabled:   false,
			Types:     []string{},
			Frequency: "immediate",
		},
		SMSNotification: notification.SMSNotificationSettings{
			Enabled: false,
			Types:   []string{},
		},
		PushNotification: true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	err := repo.Create(ctx, preference)
	require.NoError(t, err)

	// Act
	found, err := repo.GetByUserID(ctx, preference.UserID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, preference.UserID, found.UserID)
	assert.Equal(t, preference.EnableSystem, found.EnableSystem)
	assert.Equal(t, preference.PushNotification, found.PushNotification)
}

// TestNotificationPreferenceRepository_GetByUserID_NotFound 测试获取不存在的用户偏好设置
func TestNotificationPreferenceRepository_GetByUserID_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPreferenceRepo(t)
	defer cleanup()

	// Act
	found, err := repo.GetByUserID(ctx, "nonexistent_user")

	// Assert
	require.NoError(t, err)
	assert.Nil(t, found)
}

// TestNotificationPreferenceRepository_Update 测试更新通知偏好设置
func TestNotificationPreferenceRepository_Update(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPreferenceRepo(t)
	defer cleanup()

	preference := &notification.NotificationPreference{
		UserID:           "user101",
		EnableSystem:     true,
		EnableSocial:     true,
		EnableContent:    true,
		EnableReward:     true,
		EnableMessage:    true,
		EnableUpdate:     true,
		EnableMembership: true,
		EmailNotification: notification.EmailNotificationSettings{
			Enabled:   false,
			Types:     []string{},
			Frequency: "immediate",
		},
		SMSNotification: notification.SMSNotificationSettings{
			Enabled: false,
			Types:   []string{},
		},
		PushNotification: true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	err := repo.Create(ctx, preference)
	require.NoError(t, err)

	// Act - 更新偏好设置
	updates := map[string]interface{}{
		"enable_social":       false,
		"push_notification":   false,
		"email_notification": notification.EmailNotificationSettings{
			Enabled:   true,
			Types:     []string{"system", "content"},
			Frequency: "daily",
		},
		"updated_at": time.Now(),
	}
	err = repo.Update(ctx, preference.ID, updates)

	// Assert
	require.NoError(t, err)

	// 验证更新
	found, err := repo.GetByID(ctx, preference.ID)
	require.NoError(t, err)
	assert.False(t, found.EnableSocial)
	assert.False(t, found.PushNotification)
	assert.True(t, found.EmailNotification.Enabled)
}

// TestNotificationPreferenceRepository_Update_NotFound 测试更新不存在的偏好设置
func TestNotificationPreferenceRepository_Update_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPreferenceRepo(t)
	defer cleanup()

	// Act
	updates := map[string]interface{}{
		"enable_social": false,
	}
	err := repo.Update(ctx, "507f1f77bcf86cd799439011", updates)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "通知偏好设置不存在")
}

// TestNotificationPreferenceRepository_Delete 测试删除通知偏好设置
func TestNotificationPreferenceRepository_Delete(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPreferenceRepo(t)
	defer cleanup()

	preference := &notification.NotificationPreference{
		UserID:           "user202",
		EnableSystem:     true,
		EnableSocial:     true,
		EnableContent:    true,
		EnableReward:     true,
		EnableMessage:    true,
		EnableUpdate:     true,
		EnableMembership: true,
		EmailNotification: notification.EmailNotificationSettings{
			Enabled:   false,
			Types:     []string{},
			Frequency: "immediate",
		},
		SMSNotification: notification.SMSNotificationSettings{
			Enabled: false,
			Types:   []string{},
		},
		PushNotification: true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	err := repo.Create(ctx, preference)
	require.NoError(t, err)

	// Act - 删除偏好设置
	err = repo.Delete(ctx, preference.ID)

	// Assert
	require.NoError(t, err)

	// 验证已删除
	found, err := repo.GetByID(ctx, preference.ID)
	require.NoError(t, err)
	assert.Nil(t, found)
}

// TestNotificationPreferenceRepository_Delete_NotFound 测试删除不存在的偏好设置
func TestNotificationPreferenceRepository_Delete_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPreferenceRepo(t)
	defer cleanup()

	// Act
	err := repo.Delete(ctx, "507f1f77bcf86cd799439011")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "通知偏好设置不存在")
}

// TestNotificationPreferenceRepository_Exists 测试检查偏好设置是否存在
func TestNotificationPreferenceRepository_Exists(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPreferenceRepo(t)
	defer cleanup()

	preference := &notification.NotificationPreference{
		UserID:           "user303",
		EnableSystem:     true,
		EnableSocial:     true,
		EnableContent:    true,
		EnableReward:     true,
		EnableMessage:    true,
		EnableUpdate:     true,
		EnableMembership: true,
		EmailNotification: notification.EmailNotificationSettings{
			Enabled:   false,
			Types:     []string{},
			Frequency: "immediate",
		},
		SMSNotification: notification.SMSNotificationSettings{
			Enabled: false,
			Types:   []string{},
		},
		PushNotification: true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	err := repo.Create(ctx, preference)
	require.NoError(t, err)

	// Act
	exists, err := repo.Exists(ctx, preference.UserID)

	// Assert
	require.NoError(t, err)
	assert.True(t, exists)
}

// TestNotificationPreferenceRepository_Exists_NotFound 测试检查不存在的偏好设置
func TestNotificationPreferenceRepository_Exists_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPreferenceRepo(t)
	defer cleanup()

	// Act
	exists, err := repo.Exists(ctx, "nonexistent_user")

	// Assert
	require.NoError(t, err)
	assert.False(t, exists)
}

// TestNotificationPreferenceRepository_BatchUpdate 测试批量更新通知偏好设置
func TestNotificationPreferenceRepository_BatchUpdate(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPreferenceRepo(t)
	defer cleanup()

	var ids []string

	// 创建多个偏好设置
	for i := 0; i < 3; i++ {
		preference := &notification.NotificationPreference{
			UserID:           func() string { return "user404_" + string(rune('0'+i)) }(),
			EnableSystem:     true,
			EnableSocial:     true,
			EnableContent:    true,
			EnableReward:     true,
			EnableMessage:    true,
			EnableUpdate:     true,
			EnableMembership: true,
			EmailNotification: notification.EmailNotificationSettings{
				Enabled:   false,
				Types:     []string{},
				Frequency: "immediate",
			},
			SMSNotification: notification.SMSNotificationSettings{
				Enabled: false,
				Types:   []string{},
			},
			PushNotification: true,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		err := repo.Create(ctx, preference)
		require.NoError(t, err)
		ids = append(ids, preference.ID)
	}

	// Act - 批量更新
	updates := map[string]interface{}{
		"push_notification": false,
		"updated_at":        time.Now(),
	}
	err := repo.BatchUpdate(ctx, ids, updates)

	// Assert
	require.NoError(t, err)

	// 验证所有偏好设置都已更新
	for _, id := range ids {
		found, err := repo.GetByID(ctx, id)
		require.NoError(t, err)
		assert.False(t, found.PushNotification)
	}
}

// TestNotificationPreferenceRepository_BatchUpdate_WithInvalidID 测试批量更新包含无效ID
func TestNotificationPreferenceRepository_BatchUpdate_WithInvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPreferenceRepo(t)
	defer cleanup()

	preference := &notification.NotificationPreference{
		UserID:           "user505",
		EnableSystem:     true,
		EnableSocial:     true,
		EnableContent:    true,
		EnableReward:     true,
		EnableMessage:    true,
		EnableUpdate:     true,
		EnableMembership: true,
		EmailNotification: notification.EmailNotificationSettings{
			Enabled:   false,
			Types:     []string{},
			Frequency: "immediate",
		},
		SMSNotification: notification.SMSNotificationSettings{
			Enabled: false,
			Types:   []string{},
		},
		PushNotification: true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	err := repo.Create(ctx, preference)
	require.NoError(t, err)

	// Act - 包含无效ID
	ids := []string{preference.ID, "invalid-id"}
	updates := map[string]interface{}{
		"push_notification": false,
		"updated_at":        time.Now(),
	}
	err = repo.BatchUpdate(ctx, ids, updates)

	// Assert
	require.NoError(t, err)

	// 验证有效ID的偏好设置已更新
	found, err := repo.GetByID(ctx, preference.ID)
	require.NoError(t, err)
	assert.False(t, found.PushNotification)
}

// TestNotificationPreferenceRepository_Create_Default 测试创建默认偏好设置
func TestNotificationPreferenceRepository_Create_Default(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPreferenceRepo(t)
	defer cleanup()

	preference := notification.NewNotificationPreference("user606")

	// Act
	err := repo.Create(ctx, preference)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, preference.ID)
	assert.True(t, preference.EnableSystem)
	assert.True(t, preference.EnableSocial)
	assert.True(t, preference.EnableContent)
	assert.True(t, preference.EnableReward)
	assert.True(t, preference.EnableMessage)
	assert.True(t, preference.EnableUpdate)
	assert.True(t, preference.EnableMembership)
	assert.True(t, preference.PushNotification)
}

// TestNotificationPreferenceRepository_IsTypeEnabled 测试检查类型是否启用
func TestNotificationPreferenceRepository_IsTypeEnabled(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPreferenceRepo(t)
	defer cleanup()

	preference := &notification.NotificationPreference{
		UserID:           "user707",
		EnableSystem:     true,
		EnableSocial:     false,
		EnableContent:    true,
		EnableReward:     true,
		EnableMessage:    true,
		EnableUpdate:     true,
		EnableMembership: true,
		EmailNotification: notification.EmailNotificationSettings{
			Enabled:   false,
			Types:     []string{},
			Frequency: "immediate",
		},
		SMSNotification: notification.SMSNotificationSettings{
			Enabled: false,
			Types:   []string{},
		},
		PushNotification: true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	err := repo.Create(ctx, preference)
	require.NoError(t, err)

	// Act & Assert
	found, err := repo.GetByUserID(ctx, preference.UserID)
	require.NoError(t, err)

	assert.True(t, found.IsTypeEnabled(notification.NotificationTypeSystem))
	assert.False(t, found.IsTypeEnabled(notification.NotificationTypeSocial))
	assert.True(t, found.IsTypeEnabled(notification.NotificationTypeContent))
}

// TestNotificationPreferenceRepository_IsEmailEnabledForType 测试检查邮件通知是否启用
func TestNotificationPreferenceRepository_IsEmailEnabledForType(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPreferenceRepo(t)
	defer cleanup()

	preference := &notification.NotificationPreference{
		UserID:           "user808",
		EnableSystem:     true,
		EnableSocial:     true,
		EnableContent:    true,
		EnableReward:     true,
		EnableMessage:    true,
		EnableUpdate:     true,
		EnableMembership: true,
		EmailNotification: notification.EmailNotificationSettings{
			Enabled:   true,
			Types:     []string{"system", "content"},
			Frequency: "daily",
		},
		SMSNotification: notification.SMSNotificationSettings{
			Enabled: false,
			Types:   []string{},
		},
		PushNotification: true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	err := repo.Create(ctx, preference)
	require.NoError(t, err)

	// Act & Assert
	found, err := repo.GetByUserID(ctx, preference.UserID)
	require.NoError(t, err)

	assert.True(t, found.IsEmailEnabledForType(notification.NotificationTypeSystem))
	assert.False(t, found.IsEmailEnabledForType(notification.NotificationTypeSocial))
	assert.True(t, found.IsEmailEnabledForType(notification.NotificationTypeContent))
}

// TestNotificationPreferenceRepository_IsSMSEnabledForType 测试检查短信通知是否启用
func TestNotificationPreferenceRepository_IsSMSEnabledForType(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPreferenceRepo(t)
	defer cleanup()

	preference := &notification.NotificationPreference{
		UserID:           "user909",
		EnableSystem:     true,
		EnableSocial:     true,
		EnableContent:    true,
		EnableReward:     true,
		EnableMessage:    true,
		EnableUpdate:     true,
		EnableMembership: true,
		EmailNotification: notification.EmailNotificationSettings{
			Enabled:   false,
			Types:     []string{},
			Frequency: "immediate",
		},
		SMSNotification: notification.SMSNotificationSettings{
			Enabled: true,
			Types:   []string{"system", "reward"},
		},
		PushNotification: true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	err := repo.Create(ctx, preference)
	require.NoError(t, err)

	// Act & Assert
	found, err := repo.GetByUserID(ctx, preference.UserID)
	require.NoError(t, err)

	assert.True(t, found.IsSMSEnabledForType(notification.NotificationTypeSystem))
	assert.False(t, found.IsSMSEnabledForType(notification.NotificationTypeSocial))
	assert.True(t, found.IsSMSEnabledForType(notification.NotificationTypeReward))
}

// TestNotificationPreferenceRepository_Update_QuietHours 测试更新免打扰时间
func TestNotificationPreferenceRepository_Update_QuietHours(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPreferenceRepo(t)
	defer cleanup()

	preference := &notification.NotificationPreference{
		UserID:           "user1010",
		EnableSystem:     true,
		EnableSocial:     true,
		EnableContent:    true,
		EnableReward:     true,
		EnableMessage:    true,
		EnableUpdate:     true,
		EnableMembership: true,
		EmailNotification: notification.EmailNotificationSettings{
			Enabled:   false,
			Types:     []string{},
			Frequency: "immediate",
		},
		SMSNotification: notification.SMSNotificationSettings{
			Enabled: false,
			Types:   []string{},
		},
		PushNotification:  true,
		QuietHoursStart:   func() *string { s := "22:00"; return &s }(),
		QuietHoursEnd:     func() *string { s := "08:00"; return &s }(),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	err := repo.Create(ctx, preference)
	require.NoError(t, err)

	// Act - 更新免打扰时间
	newStart := "23:00"
	newEnd := "07:00"
	updates := map[string]interface{}{
		"quiet_hours_start": &newStart,
		"quiet_hours_end":   &newEnd,
		"updated_at":        time.Now(),
	}
	err = repo.Update(ctx, preference.ID, updates)

	// Assert
	require.NoError(t, err)

	// 验证更新
	found, err := repo.GetByID(ctx, preference.ID)
	require.NoError(t, err)
	assert.Equal(t, "23:00", *found.QuietHoursStart)
	assert.Equal(t, "07:00", *found.QuietHoursEnd)
}

// TestNotificationPreferenceRepository_BatchUpdate_EmptyIDs 测试批量更新空ID列表
func TestNotificationPreferenceRepository_BatchUpdate_EmptyIDs(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPreferenceRepo(t)
	defer cleanup()

	// Act
	updates := map[string]interface{}{
		"push_notification": false,
		"updated_at":        time.Now(),
	}
	err := repo.BatchUpdate(ctx, []string{}, updates)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "没有有效的通知偏好设置ID")
}

// TestNotificationPreferenceRepository_BatchUpdate_AllInvalidIDs 测试批量更新全部为无效ID
func TestNotificationPreferenceRepository_BatchUpdate_AllInvalidIDs(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPreferenceRepo(t)
	defer cleanup()

	// Act
	updates := map[string]interface{}{
		"push_notification": false,
		"updated_at":        time.Now(),
	}
	err := repo.BatchUpdate(ctx, []string{"invalid-id-1", "invalid-id-2"}, updates)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "没有有效的通知偏好设置ID")
}
