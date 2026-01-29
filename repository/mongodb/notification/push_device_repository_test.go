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

// setupPushDeviceRepo 测试辅助函数
func setupPushDeviceRepo(t *testing.T) (*notificationRepo.PushDeviceRepositoryImpl, context.Context, func()) {
	db, cleanup := testutil.SetupTestDB(t)
	repo := notificationRepo.NewPushDeviceRepository(db).(*notificationRepo.PushDeviceRepositoryImpl)
	ctx := context.Background()
	return repo, ctx, cleanup
}

// TestPushDeviceRepository_Create 测试创建推送设备
func TestPushDeviceRepository_Create(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	device := &notification.PushDevice{
		UserID:      "user123",
		DeviceType:  "ios",
		DeviceToken: "test_token_123",
		DeviceID:    "device_abc",
		IsActive:    true,
		LastUsedAt:  time.Now(),
		CreatedAt:   time.Now(),
	}

	// Act
	err := repo.Create(ctx, device)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, device.ID)
}

// TestPushDeviceRepository_Create_WithInvalidID 测试使用无效ID创建推送设备
func TestPushDeviceRepository_Create_WithInvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	device := &notification.PushDevice{
		ID:          "invalid-id",
		UserID:      "user123",
		DeviceType:  "android",
		DeviceToken: "test_token_456",
		DeviceID:    "device_def",
		IsActive:    true,
		LastUsedAt:  time.Now(),
		CreatedAt:   time.Now(),
	}

	// Act
	err := repo.Create(ctx, device)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, device.ID)
	assert.NotEqual(t, "invalid-id", device.ID)
}

// TestPushDeviceRepository_GetByID 测试根据ID获取推送设备
func TestPushDeviceRepository_GetByID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	device := &notification.PushDevice{
		UserID:      "user456",
		DeviceType:  "ios",
		DeviceToken: "test_token_789",
		DeviceID:    "device_ghi",
		IsActive:    true,
		LastUsedAt:  time.Now(),
		CreatedAt:   time.Now(),
	}
	err := repo.Create(ctx, device)
	require.NoError(t, err)

	// Act
	found, err := repo.GetByID(ctx, device.ID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, device.UserID, found.UserID)
	assert.Equal(t, device.DeviceType, found.DeviceType)
	assert.Equal(t, device.DeviceToken, found.DeviceToken)
	assert.Equal(t, device.DeviceID, found.DeviceID)
	assert.Equal(t, device.IsActive, found.IsActive)
}

// TestPushDeviceRepository_GetByID_NotFound 测试获取不存在的推送设备
func TestPushDeviceRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	// Act
	found, err := repo.GetByID(ctx, "507f1f77bcf86cd799439011")

	// Assert
	require.NoError(t, err)
	assert.Nil(t, found)
}

// TestPushDeviceRepository_GetByID_InvalidID 测试使用无效ID获取推送设备
func TestPushDeviceRepository_GetByID_InvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	// Act
	found, err := repo.GetByID(ctx, "invalid-id")

	// Assert
	require.Error(t, err)
	assert.Nil(t, found)
	assert.Contains(t, err.Error(), "无效的设备ID")
}

// TestPushDeviceRepository_Update 测试更新推送设备
func TestPushDeviceRepository_Update(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	device := &notification.PushDevice{
		UserID:      "user789",
		DeviceType:  "android",
		DeviceToken: "old_token",
		DeviceID:    "device_jkl",
		IsActive:    true,
		LastUsedAt:  time.Now(),
		CreatedAt:   time.Now(),
	}
	err := repo.Create(ctx, device)
	require.NoError(t, err)

	// Act - 更新设备
	updates := map[string]interface{}{
		"device_token": "new_token",
		"is_active":    false,
	}
	err = repo.Update(ctx, device.ID, updates)

	// Assert
	require.NoError(t, err)

	// 验证更新
	found, err := repo.GetByID(ctx, device.ID)
	require.NoError(t, err)
	assert.Equal(t, "new_token", found.DeviceToken)
	assert.False(t, found.IsActive)
}

// TestPushDeviceRepository_Update_NotFound 测试更新不存在的推送设备
func TestPushDeviceRepository_Update_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	// Act
	updates := map[string]interface{}{
		"device_token": "new_token",
	}
	err := repo.Update(ctx, "507f1f77bcf86cd799439011", updates)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "推送设备不存在")
}

// TestPushDeviceRepository_Delete 测试删除推送设备
func TestPushDeviceRepository_Delete(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	device := &notification.PushDevice{
		UserID:      "user101",
		DeviceType:  "ios",
		DeviceToken: "test_token_delete",
		DeviceID:    "device_mno",
		IsActive:    true,
		LastUsedAt:  time.Now(),
		CreatedAt:   time.Now(),
	}
	err := repo.Create(ctx, device)
	require.NoError(t, err)

	// Act - 删除设备
	err = repo.Delete(ctx, device.ID)

	// Assert
	require.NoError(t, err)

	// 验证已删除
	found, err := repo.GetByID(ctx, device.ID)
	require.NoError(t, err)
	assert.Nil(t, found)
}

// TestPushDeviceRepository_Delete_NotFound 测试删除不存在的推送设备
func TestPushDeviceRepository_Delete_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	// Act
	err := repo.Delete(ctx, "507f1f77bcf86cd799439011")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "推送设备不存在")
}

// TestPushDeviceRepository_Exists 测试检查推送设备是否存在
func TestPushDeviceRepository_Exists(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	device := &notification.PushDevice{
		UserID:      "user202",
		DeviceType:  "android",
		DeviceToken: "test_token_exists",
		DeviceID:    "device_pqr",
		IsActive:    true,
		LastUsedAt:  time.Now(),
		CreatedAt:   time.Now(),
	}
	err := repo.Create(ctx, device)
	require.NoError(t, err)

	// Act
	exists, err := repo.Exists(ctx, device.ID)

	// Assert
	require.NoError(t, err)
	assert.True(t, exists)
}

// TestPushDeviceRepository_Exists_NotFound 测试检查不存在的推送设备
func TestPushDeviceRepository_Exists_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	// Act
	exists, err := repo.Exists(ctx, "507f1f77bcf86cd799439011")

	// Assert
	require.NoError(t, err)
	assert.False(t, exists)
}

// TestPushDeviceRepository_GetByUserID 测试根据用户ID获取推送设备列表
func TestPushDeviceRepository_GetByUserID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	userID := "user303"

	// 创建多个设备
	deviceTypes := []string{"ios", "android", "web"}
	for i, deviceType := range deviceTypes {
		device := &notification.PushDevice{
			UserID:      userID,
			DeviceType:  deviceType,
			DeviceToken: func() string { return "token_" + string(rune('0'+i)) }(),
			DeviceID:    func() string { return "device_" + string(rune('0'+i)) }(),
			IsActive:    true,
			LastUsedAt:  time.Now().Add(time.Duration(i) * time.Hour),
			CreatedAt:   time.Now(),
		}
		err := repo.Create(ctx, device)
		require.NoError(t, err)
	}

	// Act
	devices, err := repo.GetByUserID(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, devices)
	assert.GreaterOrEqual(t, len(devices), 3)
}

// TestPushDeviceRepository_GetByDeviceID 测试根据设备ID获取推送设备
func TestPushDeviceRepository_GetByDeviceID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	device := &notification.PushDevice{
		UserID:      "user404",
		DeviceType:  "ios",
		DeviceToken: "test_token_device_id",
		DeviceID:    "unique_device_id",
		IsActive:    true,
		LastUsedAt:  time.Now(),
		CreatedAt:   time.Now(),
	}
	err := repo.Create(ctx, device)
	require.NoError(t, err)

	// Act
	found, err := repo.GetByDeviceID(ctx, device.DeviceID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, device.DeviceID, found.DeviceID)
	assert.Equal(t, device.UserID, found.UserID)
}

// TestPushDeviceRepository_GetByDeviceID_NotFound 测试根据设备ID获取不存在的设备
func TestPushDeviceRepository_GetByDeviceID_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	// Act
	found, err := repo.GetByDeviceID(ctx, "nonexistent_device_id")

	// Assert
	require.NoError(t, err)
	assert.Nil(t, found)
}

// TestPushDeviceRepository_GetActiveByUserID 测试获取用户的活跃推送设备列表
func TestPushDeviceRepository_GetActiveByUserID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	userID := "user505"

	// 创建活跃设备
	activeDevice := &notification.PushDevice{
		UserID:      userID,
		DeviceType:  "ios",
		DeviceToken: "active_token",
		DeviceID:    "active_device",
		IsActive:    true,
		LastUsedAt:  time.Now(),
		CreatedAt:   time.Now(),
	}
	err := repo.Create(ctx, activeDevice)
	require.NoError(t, err)

	// 创建不活跃设备
	inactiveDevice := &notification.PushDevice{
		UserID:      userID,
		DeviceType:  "android",
		DeviceToken: "inactive_token",
		DeviceID:    "inactive_device",
		IsActive:    false,
		LastUsedAt:  time.Now(),
		CreatedAt:   time.Now(),
	}
	err = repo.Create(ctx, inactiveDevice)
	require.NoError(t, err)

	// Act
	devices, err := repo.GetActiveByUserID(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, devices)
	assert.GreaterOrEqual(t, len(devices), 1)

	// 验证所有返回的设备都是活跃的
	for _, device := range devices {
		assert.True(t, device.IsActive)
	}
}

// TestPushDeviceRepository_BatchDelete 测试批量删除推送设备
func TestPushDeviceRepository_BatchDelete(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	var ids []string
	userID := "user606"

	// 创建多个设备
	for i := 0; i < 3; i++ {
		device := &notification.PushDevice{
			UserID:      userID,
			DeviceType:  "ios",
			DeviceToken: func() string { return "token_" + string(rune('0'+i)) }(),
			DeviceID:    func() string { return "device_" + string(rune('0'+i)) }(),
			IsActive:    true,
			LastUsedAt:  time.Now(),
			CreatedAt:   time.Now(),
		}
		err := repo.Create(ctx, device)
		require.NoError(t, err)
		ids = append(ids, device.ID)
	}

	// Act
	err := repo.BatchDelete(ctx, ids)

	// Assert
	require.NoError(t, err)

	// 验证已删除
	for _, id := range ids {
		found, err := repo.GetByID(ctx, id)
		require.NoError(t, err)
		assert.Nil(t, found)
	}
}

// TestPushDeviceRepository_BatchDelete_WithInvalidID 测试批量删除包含无效ID
func TestPushDeviceRepository_BatchDelete_WithInvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	device := &notification.PushDevice{
		UserID:      "user707",
		DeviceType:  "ios",
		DeviceToken: "test_token_batch",
		DeviceID:    "device_batch",
		IsActive:    true,
		LastUsedAt:  time.Now(),
		CreatedAt:   time.Now(),
	}
	err := repo.Create(ctx, device)
	require.NoError(t, err)

	// Act - 包含无效ID
	ids := []string{device.ID, "invalid-id"}
	err = repo.BatchDelete(ctx, ids)

	// Assert
	require.NoError(t, err)

	// 验证有效ID的设备已删除
	found, err := repo.GetByID(ctx, device.ID)
	require.NoError(t, err)
	assert.Nil(t, found)
}

// TestPushDeviceRepository_DeactivateAllForUser 测试停用用户所有推送设备
func TestPushDeviceRepository_DeactivateAllForUser(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	userID := "user808"

	// 创建多个设备
	for i := 0; i < 3; i++ {
		device := &notification.PushDevice{
			UserID:      userID,
			DeviceType:  "ios",
			DeviceToken: func() string { return "token_" + string(rune('0'+i)) }(),
			DeviceID:    func() string { return "device_" + string(rune('0'+i)) }(),
			IsActive:    true,
			LastUsedAt:  time.Now(),
			CreatedAt:   time.Now(),
		}
		err := repo.Create(ctx, device)
		require.NoError(t, err)
	}

	// Act
	err := repo.DeactivateAllForUser(ctx, userID)

	// Assert
	require.NoError(t, err)

	// 验证所有设备已停用
	devices, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	for _, device := range devices {
		assert.False(t, device.IsActive)
	}
}

// TestPushDeviceRepository_UpdateLastUsed 测试更新最后使用时间
func TestPushDeviceRepository_UpdateLastUsed(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	oldTime := time.Now().Add(-24 * time.Hour)
	device := &notification.PushDevice{
		UserID:      "user909",
		DeviceType:  "ios",
		DeviceToken: "test_token_last_used",
		DeviceID:    "device_last_used",
		IsActive:    true,
		LastUsedAt:  oldTime,
		CreatedAt:   time.Now(),
	}
	err := repo.Create(ctx, device)
	require.NoError(t, err)

	// 等待一小段时间确保时间不同
	time.Sleep(10 * time.Millisecond)

	// Act
	err = repo.UpdateLastUsed(ctx, device.ID)

	// Assert
	require.NoError(t, err)

	// 验证最后使用时间已更新
	found, err := repo.GetByID(ctx, device.ID)
	require.NoError(t, err)
	assert.True(t, found.LastUsedAt.After(oldTime))
}

// TestPushDeviceRepository_UpdateLastUsed_NotFound 测试更新不存在设备的最后使用时间
func TestPushDeviceRepository_UpdateLastUsed_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	// Act
	err := repo.UpdateLastUsed(ctx, "507f1f77bcf86cd799439011")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "推送设备不存在")
}

// TestPushDeviceRepository_DeleteInactiveDevices 测试删除不活跃的设备
func TestPushDeviceRepository_DeleteInactiveDevices(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	// 创建不活跃设备（30天前使用）
	oldTime := time.Now().Add(-30 * 24 * time.Hour)
	inactiveDevice := &notification.PushDevice{
		UserID:      "user1010",
		DeviceType:  "ios",
		DeviceToken: "inactive_token",
		DeviceID:    "inactive_device",
		IsActive:    true,
		LastUsedAt:  oldTime,
		CreatedAt:   time.Now(),
	}
	err := repo.Create(ctx, inactiveDevice)
	require.NoError(t, err)

	// 创建活跃设备（刚刚使用）
	activeDevice := &notification.PushDevice{
		UserID:      "user1010",
		DeviceType:  "android",
		DeviceToken: "active_token",
		DeviceID:    "active_device",
		IsActive:    true,
		LastUsedAt:  time.Now(),
		CreatedAt:   time.Now(),
	}
	err = repo.Create(ctx, activeDevice)
	require.NoError(t, err)

	// Act - 删除7天前使用的设备
	beforeDate := time.Now().Add(-7 * 24 * time.Hour)
	count, err := repo.DeleteInactiveDevices(ctx, beforeDate)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(1))

	// 验证不活跃设备已删除
	found, err := repo.GetByID(ctx, inactiveDevice.ID)
	require.NoError(t, err)
	assert.Nil(t, found)

	// 验证活跃设备仍然存在
	found, err = repo.GetByID(ctx, activeDevice.ID)
	require.NoError(t, err)
	assert.NotNil(t, found)
}

// TestPushDeviceRepository_Create_Default 测试创建默认推送设备
func TestPushDeviceRepository_Create_Default(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	device := notification.NewPushDevice("user1111", "ios", "test_token", "device_id")

	// Act
	err := repo.Create(ctx, device)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, device.ID)
	assert.True(t, device.IsActive)
}

// TestPushDeviceRepository_BatchDelete_EmptyIDs 测试批量删除空ID列表
func TestPushDeviceRepository_BatchDelete_EmptyIDs(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	// Act
	err := repo.BatchDelete(ctx, []string{})

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "没有有效的设备ID")
}

// TestPushDeviceRepository_BatchDelete_AllInvalidIDs 测试批量删除全部为无效ID
func TestPushDeviceRepository_BatchDelete_AllInvalidIDs(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	// Act
	err := repo.BatchDelete(ctx, []string{"invalid-id-1", "invalid-id-2"})

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "没有有效的设备ID")
}

// TestPushDeviceRepository_GetByUserID_MultipleUsers 测试多个用户的设备列表
func TestPushDeviceRepository_GetByUserID_MultipleUsers(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	// 为不同用户创建设备
	user1ID := "user1212"
	user2ID := "user1313"

	// 用户1的设备
	device1 := &notification.PushDevice{
		UserID:      user1ID,
		DeviceType:  "ios",
		DeviceToken: "user1_token",
		DeviceID:    "user1_device",
		IsActive:    true,
		LastUsedAt:  time.Now(),
		CreatedAt:   time.Now(),
	}
	err := repo.Create(ctx, device1)
	require.NoError(t, err)

	// 用户2的设备
	device2 := &notification.PushDevice{
		UserID:      user2ID,
		DeviceType:  "android",
		DeviceToken: "user2_token",
		DeviceID:    "user2_device",
		IsActive:    true,
		LastUsedAt:  time.Now(),
		CreatedAt:   time.Now(),
	}
	err = repo.Create(ctx, device2)
	require.NoError(t, err)

	// Act - 获取用户1的设备
	user1Devices, err := repo.GetByUserID(ctx, user1ID)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(user1Devices), 1)

	// 验证只返回用户1的设备
	for _, device := range user1Devices {
		assert.Equal(t, user1ID, device.UserID)
	}
}

// TestPushDeviceRepository_Update_Token 测试更新设备Token
func TestPushDeviceRepository_Update_Token(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupPushDeviceRepo(t)
	defer cleanup()

	device := &notification.PushDevice{
		UserID:      "user1414",
		DeviceType:  "ios",
		DeviceToken: "old_token",
		DeviceID:    "device_token_test",
		IsActive:    true,
		LastUsedAt:  time.Now(),
		CreatedAt:   time.Now(),
	}
	err := repo.Create(ctx, device)
	require.NoError(t, err)

	// Act - 更新token
	updates := map[string]interface{}{
		"device_token": "new_token",
	}
	err = repo.Update(ctx, device.ID, updates)

	// Assert
	require.NoError(t, err)

	// 验证token已更新
	found, err := repo.GetByID(ctx, device.ID)
	require.NoError(t, err)
	assert.Equal(t, "new_token", found.DeviceToken)
}
