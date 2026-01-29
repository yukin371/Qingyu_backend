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

// setupNotificationRepo 测试辅助函数
func setupNotificationRepo(t *testing.T) (*notificationRepo.NotificationRepositoryImpl, context.Context, func()) {
	db, cleanup := testutil.SetupTestDB(t)
	repo := notificationRepo.NewNotificationRepository(db).(*notificationRepo.NotificationRepositoryImpl)
	ctx := context.Background()
	return repo, ctx, cleanup
}

// TestNotificationRepository_Create 测试创建通知
func TestNotificationRepository_Create(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	notif := &notification.Notification{
		UserID:    "user123",
		Type:      notification.NotificationTypeSystem,
		Priority:  notification.NotificationPriorityNormal,
		Title:     "测试通知",
		Content:   "这是一条测试通知",
		Data:      map[string]interface{}{"key": "value"},
		Read:      false,
		CreatedAt: time.Now(),
	}

	// Act
	err := repo.Create(ctx, notif)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, notif.ID)
}

// TestNotificationRepository_Create_WithInvalidID 测试使用无效ID创建通知
func TestNotificationRepository_Create_WithInvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	notif := &notification.Notification{
		ID:        "invalid-id",
		UserID:    "user123",
		Type:      notification.NotificationTypeSystem,
		Priority:  notification.NotificationPriorityNormal,
		Title:     "测试通知",
		Content:   "这是一条测试通知",
		Read:      false,
		CreatedAt: time.Now(),
	}

	// Act
	err := repo.Create(ctx, notif)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, notif.ID)
	assert.NotEqual(t, "invalid-id", notif.ID)
}

// TestNotificationRepository_GetByID 测试根据ID获取通知
func TestNotificationRepository_GetByID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	notif := &notification.Notification{
		UserID:    "user123",
		Type:      notification.NotificationTypeSocial,
		Priority:  notification.NotificationPriorityHigh,
		Title:     "社交通知",
		Content:   "有人关注了你",
		Read:      false,
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, notif)
	require.NoError(t, err)

	// Act
	found, err := repo.GetByID(ctx, notif.ID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, notif.Title, found.Title)
	assert.Equal(t, notif.Content, found.Content)
	assert.Equal(t, notif.Type, found.Type)
	assert.Equal(t, notif.Priority, found.Priority)
}

// TestNotificationRepository_GetByID_NotFound 测试获取不存在的通知
func TestNotificationRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	// Act
	found, err := repo.GetByID(ctx, "507f1f77bcf86cd799439011")

	// Assert
	require.NoError(t, err)
	assert.Nil(t, found)
}

// TestNotificationRepository_GetByID_InvalidID 测试使用无效ID获取通知
func TestNotificationRepository_GetByID_InvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	// Act
	found, err := repo.GetByID(ctx, "invalid-id")

	// Assert
	require.Error(t, err)
	assert.Nil(t, found)
	assert.Contains(t, err.Error(), "无效的通知ID")
}

// TestNotificationRepository_Update 测试更新通知
func TestNotificationRepository_Update(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	notif := &notification.Notification{
		UserID:    "user123",
		Type:      notification.NotificationTypeContent,
		Priority:  notification.NotificationPriorityNormal,
		Title:     "原始标题",
		Content:   "原始内容",
		Read:      false,
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, notif)
	require.NoError(t, err)

	// Act - 更新通知
	updates := map[string]interface{}{
		"title": "更新后的标题",
		"read":  true,
	}
	err = repo.Update(ctx, notif.ID, updates)

	// Assert
	require.NoError(t, err)

	// 验证更新
	found, err := repo.GetByID(ctx, notif.ID)
	require.NoError(t, err)
	assert.Equal(t, "更新后的标题", found.Title)
	assert.True(t, found.Read)
	assert.NotNil(t, found.ReadAt)
}

// TestNotificationRepository_Update_NotFound 测试更新不存在的通知
func TestNotificationRepository_Update_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	// Act
	updates := map[string]interface{}{
		"title": "更新后的标题",
	}
	err := repo.Update(ctx, "507f1f77bcf86cd799439011", updates)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "通知不存在")
}

// TestNotificationRepository_Delete 测试删除通知
func TestNotificationRepository_Delete(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	notif := &notification.Notification{
		UserID:    "user123",
		Type:      notification.NotificationTypeReward,
		Priority:  notification.NotificationPriorityNormal,
		Title:     "打赏通知",
		Content:   "收到打赏",
		Read:      false,
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, notif)
	require.NoError(t, err)

	// Act - 删除通知
	err = repo.Delete(ctx, notif.ID)

	// Assert
	require.NoError(t, err)

	// 验证已删除
	found, err := repo.GetByID(ctx, notif.ID)
	require.NoError(t, err)
	assert.Nil(t, found)
}

// TestNotificationRepository_Delete_NotFound 测试删除不存在的通知
func TestNotificationRepository_Delete_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	// Act
	err := repo.Delete(ctx, "507f1f77bcf86cd799439011")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "通知不存在")
}

// TestNotificationRepository_Exists 测试检查通知是否存在
func TestNotificationRepository_Exists(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	notif := &notification.Notification{
		UserID:    "user123",
		Type:      notification.NotificationTypeMessage,
		Priority:  notification.NotificationPriorityNormal,
		Title:     "私信通知",
		Content:   "收到新私信",
		Read:      false,
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, notif)
	require.NoError(t, err)

	// Act
	exists, err := repo.Exists(ctx, notif.ID)

	// Assert
	require.NoError(t, err)
	assert.True(t, exists)
}

// TestNotificationRepository_Exists_NotFound 测试检查不存在的通知
func TestNotificationRepository_Exists_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	// Act
	exists, err := repo.Exists(ctx, "507f1f77bcf86cd799439011")

	// Assert
	require.NoError(t, err)
	assert.False(t, exists)
}

// TestNotificationRepository_List 测试获取通知列表
func TestNotificationRepository_List(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	userID := "user123"

	// 创建多条通知
	for i := 0; i < 5; i++ {
		notif := &notification.Notification{
			UserID:    userID,
			Type:      notification.NotificationTypeSystem,
			Priority:  notification.NotificationPriorityNormal,
			Title:     "测试通知",
			Content:   "这是一条测试通知",
			Read:      false,
			CreatedAt: time.Now(),
		}
		err := repo.Create(ctx, notif)
		require.NoError(t, err)
	}

	// Act
	filter := &notification.NotificationFilter{
		UserID: &userID,
		Limit:  10,
		Offset: 0,
	}
	notifications, err := repo.List(ctx, filter)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, notifications)
	assert.GreaterOrEqual(t, len(notifications), 5)
}

// TestNotificationRepository_List_WithSort 测试带排序的通知列表
func TestNotificationRepository_List_WithSort(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	userID := "user456"

	// 创建多条通知
	for i := 0; i < 3; i++ {
		notif := &notification.Notification{
			UserID:    userID,
			Type:      notification.NotificationTypeSystem,
			Priority:  notification.NotificationPriorityNormal,
			Title:     "测试通知",
			Content:   "这是一条测试通知",
			Read:      false,
			CreatedAt: time.Now().Add(time.Duration(i) * time.Hour),
		}
		err := repo.Create(ctx, notif)
		require.NoError(t, err)
	}

	// Act - 降序排序
	filter := &notification.NotificationFilter{
		UserID:    &userID,
		SortBy:    "created_at",
		SortOrder: "desc",
		Limit:     10,
		Offset:    0,
	}
	notifications, err := repo.List(ctx, filter)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, notifications)
	assert.GreaterOrEqual(t, len(notifications), 3)
}

// TestNotificationRepository_Count 测试统计通知数量
func TestNotificationRepository_Count(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	userID := "user789"

	// 创建多条通知
	for i := 0; i < 5; i++ {
		notif := &notification.Notification{
			UserID:    userID,
			Type:      notification.NotificationTypeSystem,
			Priority:  notification.NotificationPriorityNormal,
			Title:     "测试通知",
			Content:   "这是一条测试通知",
			Read:      false,
			CreatedAt: time.Now(),
		}
		err := repo.Create(ctx, notif)
		require.NoError(t, err)
	}

	// Act
	filter := &notification.NotificationFilter{
		UserID: &userID,
	}
	count, err := repo.Count(ctx, filter)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(5))
}

// TestNotificationRepository_BatchMarkAsRead 测试批量标记为已读
func TestNotificationRepository_BatchMarkAsRead(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	var ids []string
	userID := "user999"

	// 创建多条通知
	for i := 0; i < 3; i++ {
		notif := &notification.Notification{
			UserID:    userID,
			Type:      notification.NotificationTypeSystem,
			Priority:  notification.NotificationPriorityNormal,
			Title:     "测试通知",
			Content:   "这是一条测试通知",
			Read:      false,
			CreatedAt: time.Now(),
		}
		err := repo.Create(ctx, notif)
		require.NoError(t, err)
		ids = append(ids, notif.ID)
	}

	// Act
	err := repo.BatchMarkAsRead(ctx, ids)

	// Assert
	require.NoError(t, err)

	// 验证已标记为已读
	for _, id := range ids {
		found, err := repo.GetByID(ctx, id)
		require.NoError(t, err)
		assert.True(t, found.Read)
		assert.NotNil(t, found.ReadAt)
	}
}

// TestNotificationRepository_BatchMarkAsRead_WithInvalidID 测试批量标记包含无效ID
func TestNotificationRepository_BatchMarkAsRead_WithInvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	notif := &notification.Notification{
		UserID:    "user123",
		Type:      notification.NotificationTypeSystem,
		Priority:  notification.NotificationPriorityNormal,
		Title:     "测试通知",
		Content:   "这是一条测试通知",
		Read:      false,
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, notif)
	require.NoError(t, err)

	// Act - 包含无效ID
	ids := []string{notif.ID, "invalid-id"}
	err = repo.BatchMarkAsRead(ctx, ids)

	// Assert
	require.NoError(t, err)

	// 验证有效ID的通知已被标记为已读
	found, err := repo.GetByID(ctx, notif.ID)
	require.NoError(t, err)
	assert.True(t, found.Read)
}

// TestNotificationRepository_MarkAllAsReadForUser 测试标记用户所有通知为已读
func TestNotificationRepository_MarkAllAsReadForUser(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	userID := "user111"

	// 创建多条未读通知
	for i := 0; i < 5; i++ {
		notif := &notification.Notification{
			UserID:    userID,
			Type:      notification.NotificationTypeSystem,
			Priority:  notification.NotificationPriorityNormal,
			Title:     "测试通知",
			Content:   "这是一条测试通知",
			Read:      false,
			CreatedAt: time.Now(),
		}
		err := repo.Create(ctx, notif)
		require.NoError(t, err)
	}

	// Act
	err := repo.MarkAllAsReadForUser(ctx, userID)

	// Assert
	require.NoError(t, err)

	// 验证所有通知已标记为已读
	filter := &notification.NotificationFilter{
		UserID: &userID,
		Read:   func() *bool { b := false; return &b }(),
	}
	notifications, err := repo.List(ctx, filter)
	require.NoError(t, err)
	assert.Equal(t, 0, len(notifications))
}

// TestNotificationRepository_BatchDelete 测试批量删除通知
func TestNotificationRepository_BatchDelete(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	var ids []string
	userID := "user222"

	// 创建多条通知
	for i := 0; i < 3; i++ {
		notif := &notification.Notification{
			UserID:    userID,
			Type:      notification.NotificationTypeSystem,
			Priority:  notification.NotificationPriorityNormal,
			Title:     "测试通知",
			Content:   "这是一条测试通知",
			Read:      false,
			CreatedAt: time.Now(),
		}
		err := repo.Create(ctx, notif)
		require.NoError(t, err)
		ids = append(ids, notif.ID)
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

// TestNotificationRepository_DeleteAllForUser 测试删除用户所有通知
func TestNotificationRepository_DeleteAllForUser(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	userID := "user333"

	// 创建多条通知
	for i := 0; i < 5; i++ {
		notif := &notification.Notification{
			UserID:    userID,
			Type:      notification.NotificationTypeSystem,
			Priority:  notification.NotificationPriorityNormal,
			Title:     "测试通知",
			Content:   "这是一条测试通知",
			Read:      false,
			CreatedAt: time.Now(),
		}
		err := repo.Create(ctx, notif)
		require.NoError(t, err)
	}

	// Act
	err := repo.DeleteAllForUser(ctx, userID)

	// Assert
	require.NoError(t, err)

	// 验证所有通知已删除
	filter := &notification.NotificationFilter{
		UserID: &userID,
	}
	notifications, err := repo.List(ctx, filter)
	require.NoError(t, err)
	assert.Equal(t, 0, len(notifications))
}

// TestNotificationRepository_CountUnread 测试统计未读通知数量
func TestNotificationRepository_CountUnread(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	userID := "user444"

	// 创建未读通知
	for i := 0; i < 3; i++ {
		notif := &notification.Notification{
			UserID:    userID,
			Type:      notification.NotificationTypeSystem,
			Priority:  notification.NotificationPriorityNormal,
			Title:     "未读通知",
			Content:   "这是一条未读通知",
			Read:      false,
			CreatedAt: time.Now(),
		}
		err := repo.Create(ctx, notif)
		require.NoError(t, err)
	}

	// 创建已读通知
	readNotif := &notification.Notification{
		UserID:    userID,
		Type:      notification.NotificationTypeSystem,
		Priority:  notification.NotificationPriorityNormal,
		Title:     "已读通知",
		Content:   "这是一条已读通知",
		Read:      true,
		ReadAt:    func() *time.Time { t := time.Now(); return &t }(),
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, readNotif)
	require.NoError(t, err)

	// Act
	count, err := repo.CountUnread(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(3))
}

// TestNotificationRepository_GetStats 测试获取通知统计
func TestNotificationRepository_GetStats(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	userID := "user555"

	// 创建不同类型的通知
	notifTypes := []notification.NotificationType{
		notification.NotificationTypeSystem,
		notification.NotificationTypeSocial,
		notification.NotificationTypeContent,
	}

	for _, notifType := range notifTypes {
		for i := 0; i < 2; i++ {
			notif := &notification.Notification{
				UserID:    userID,
				Type:      notifType,
				Priority:  notification.NotificationPriorityNormal,
				Title:     "测试通知",
				Content:   "这是一条测试通知",
				Read:      i%2 == 0,
				CreatedAt: time.Now(),
			}
			if notif.Read {
				notif.ReadAt = func() *time.Time { t := time.Now(); return &t }()
			}
			err := repo.Create(ctx, notif)
			require.NoError(t, err)
		}
	}

	// Act
	stats, err := repo.GetStats(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, stats)
	assert.GreaterOrEqual(t, stats.TotalCount, int64(6))
	assert.GreaterOrEqual(t, stats.UnreadCount, int64(3))
	assert.NotNil(t, stats.TypeCounts)
	assert.NotNil(t, stats.PriorityCounts)
}

// TestNotificationRepository_GetUnreadByType 测试获取指定类型的未读通知
func TestNotificationRepository_GetUnreadByType(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	userID := "user666"

	// 创建指定类型的未读通知
	for i := 0; i < 3; i++ {
		notif := &notification.Notification{
			UserID:    userID,
			Type:      notification.NotificationTypeSocial,
			Priority:  notification.NotificationPriorityNormal,
			Title:     "社交通知",
			Content:   "这是一条社交通知",
			Read:      false,
			CreatedAt: time.Now(),
		}
		err := repo.Create(ctx, notif)
		require.NoError(t, err)
	}

	// Act
	notifications, err := repo.GetUnreadByType(ctx, userID, notification.NotificationTypeSocial)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, notifications)
	assert.GreaterOrEqual(t, len(notifications), 3)
}

// TestNotificationRepository_DeleteExpired 测试删除过期通知
func TestNotificationRepository_DeleteExpired(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	// 创建过期通知
	past := time.Now().Add(-24 * time.Hour)
	expiredNotif := &notification.Notification{
		UserID:    "user777",
		Type:      notification.NotificationTypeSystem,
		Priority:  notification.NotificationPriorityNormal,
		Title:     "过期通知",
		Content:   "这是一条过期通知",
		Read:      false,
		CreatedAt: time.Now(),
		ExpiresAt: &past,
	}
	err := repo.Create(ctx, expiredNotif)
	require.NoError(t, err)

	// 创建未过期通知
	future := time.Now().Add(24 * time.Hour)
	validNotif := &notification.Notification{
		UserID:    "user777",
		Type:      notification.NotificationTypeSystem,
		Priority:  notification.NotificationPriorityNormal,
		Title:     "有效通知",
		Content:   "这是一条有效通知",
		Read:      false,
		CreatedAt: time.Now(),
		ExpiresAt: &future,
	}
	err = repo.Create(ctx, validNotif)
	require.NoError(t, err)

	// Act
	count, err := repo.DeleteExpired(ctx)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(1))

	// 验证过期通知已删除
	found, err := repo.GetByID(ctx, expiredNotif.ID)
	require.NoError(t, err)
	assert.Nil(t, found)

	// 验证有效通知仍然存在
	found, err = repo.GetByID(ctx, validNotif.ID)
	require.NoError(t, err)
	assert.NotNil(t, found)
}

// TestNotificationRepository_DeleteOldNotifications 测试删除旧通知
func TestNotificationRepository_DeleteOldNotifications(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	// 创建旧通知
	oldTime := time.Now().Add(-30 * 24 * time.Hour)
	oldNotif := &notification.Notification{
		UserID:    "user888",
		Type:      notification.NotificationTypeSystem,
		Priority:  notification.NotificationPriorityNormal,
		Title:     "旧通知",
		Content:   "这是一条旧通知",
		Read:      false,
		CreatedAt: oldTime,
	}
	err := repo.Create(ctx, oldNotif)
	require.NoError(t, err)

	// 创建新通知
	newNotif := &notification.Notification{
		UserID:    "user888",
		Type:      notification.NotificationTypeSystem,
		Priority:  notification.NotificationPriorityNormal,
		Title:     "新通知",
		Content:   "这是一条新通知",
		Read:      false,
		CreatedAt: time.Now(),
	}
	err = repo.Create(ctx, newNotif)
	require.NoError(t, err)

	// Act - 删除7天前的通知
	beforeDate := time.Now().Add(-7 * 24 * time.Hour)
	count, err := repo.DeleteOldNotifications(ctx, beforeDate)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(1))

	// 验证旧通知已删除
	found, err := repo.GetByID(ctx, oldNotif.ID)
	require.NoError(t, err)
	assert.Nil(t, found)

	// 验证新通知仍然存在
	found, err = repo.GetByID(ctx, newNotif.ID)
	require.NoError(t, err)
	assert.NotNil(t, found)
}

// TestNotificationRepository_DeleteReadForUser 测试删除用户所有已读通知
func TestNotificationRepository_DeleteReadForUser(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	userID := "user999"

	// 创建已读通知
	for i := 0; i < 3; i++ {
		notif := &notification.Notification{
			UserID:    userID,
			Type:      notification.NotificationTypeSystem,
			Priority:  notification.NotificationPriorityNormal,
			Title:     "已读通知",
			Content:   "这是一条已读通知",
			Read:      true,
			ReadAt:    func() *time.Time { t := time.Now(); return &t }(),
			CreatedAt: time.Now(),
		}
		err := repo.Create(ctx, notif)
		require.NoError(t, err)
	}

	// 创建未读通知
	unreadNotif := &notification.Notification{
		UserID:    userID,
		Type:      notification.NotificationTypeSystem,
		Priority:  notification.NotificationPriorityNormal,
		Title:     "未读通知",
		Content:   "这是一条未读通知",
		Read:      false,
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, unreadNotif)
	require.NoError(t, err)

	// Act
	count, err := repo.DeleteReadForUser(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(3))

	// 验证未读通知仍然存在
	found, err := repo.GetByID(ctx, unreadNotif.ID)
	require.NoError(t, err)
	assert.NotNil(t, found)
}

// TestNotificationRepository_List_WithFilter 测试带过滤条件的通知列表
func TestNotificationRepository_List_WithFilter(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	userID := "user101"
	searchKeyword := "重要"

	// 创建包含关键词的通知
	importantNotif := &notification.Notification{
		UserID:    userID,
		Type:      notification.NotificationTypeSystem,
		Priority:  notification.NotificationPriorityHigh,
		Title:     "重要通知",
		Content:   "这是一条重要通知",
		Read:      false,
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, importantNotif)
	require.NoError(t, err)

	// 创建不包含关键词的通知
	normalNotif := &notification.Notification{
		UserID:    userID,
		Type:      notification.NotificationTypeSystem,
		Priority:  notification.NotificationPriorityNormal,
		Title:     "普通通知",
		Content:   "这是一条普通通知",
		Read:      false,
		CreatedAt: time.Now(),
	}
	err = repo.Create(ctx, normalNotif)
	require.NoError(t, err)

	// Act - 使用关键词过滤
	filter := &notification.NotificationFilter{
		UserID:  &userID,
		Keyword: &searchKeyword,
		Limit:   10,
		Offset:  0,
	}
	notifications, err := repo.List(ctx, filter)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, notifications)
	assert.GreaterOrEqual(t, len(notifications), 1)
}

// TestNotificationRepository_List_WithDateRange 测试带日期范围的通知列表
func TestNotificationRepository_List_WithDateRange(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupNotificationRepo(t)
	defer cleanup()

	userID := "user102"
	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now().Add(24 * time.Hour)

	// 创建在日期范围内的通知
	notif := &notification.Notification{
		UserID:    userID,
		Type:      notification.NotificationTypeSystem,
		Priority:  notification.NotificationPriorityNormal,
		Title:     "测试通知",
		Content:   "这是一条测试通知",
		Read:      false,
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, notif)
	require.NoError(t, err)

	// Act - 使用日期范围过滤
	filter := &notification.NotificationFilter{
		UserID:    &userID,
		StartDate: &startDate,
		EndDate:   &endDate,
		Limit:     10,
		Offset:    0,
	}
	notifications, err := repo.List(ctx, filter)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, notifications)
	assert.GreaterOrEqual(t, len(notifications), 1)
}
