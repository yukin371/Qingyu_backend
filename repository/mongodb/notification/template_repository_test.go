package notification_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/models/notification"
	notificationRepo "Qingyu_backend/repository/mongodb/notification"
	repoInterfaces "Qingyu_backend/repository/interfaces/notification"
	"Qingyu_backend/test/testutil"
)

// setupTemplateRepo 测试辅助函数
func setupTemplateRepo(t *testing.T) (*notificationRepo.NotificationTemplateRepositoryImpl, context.Context, func()) {
	db, cleanup := testutil.SetupTestDB(t)
	repo := notificationRepo.NewNotificationTemplateRepository(db).(*notificationRepo.NotificationTemplateRepositoryImpl)
	ctx := context.Background()
	return repo, ctx, cleanup
}

// TestNotificationTemplateRepository_Create 测试创建通知模板
func TestNotificationTemplateRepository_Create(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	template := &notification.NotificationTemplate{
		Type:      notification.NotificationTypeSystem,
		Action:    "announcement",
		Title:     "系统公告",
		Content:   "这是一个系统公告模板",
		Variables: []string{"title", "content"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	err := repo.Create(ctx, template)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, template.ID)
}

// TestNotificationTemplateRepository_Create_WithInvalidID 测试使用无效ID创建通知模板
func TestNotificationTemplateRepository_Create_WithInvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	template := &notification.NotificationTemplate{
		ID:        "invalid-id",
		Type:      notification.NotificationTypeSocial,
		Action:    "follow",
		Title:     "关注通知",
		Content:   "有人关注了你",
		Variables: []string{"username"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	err := repo.Create(ctx, template)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, template.ID)
	assert.NotEqual(t, "invalid-id", template.ID)
}

// TestNotificationTemplateRepository_GetByID 测试根据ID获取通知模板
func TestNotificationTemplateRepository_GetByID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	template := &notification.NotificationTemplate{
		Type:      notification.NotificationTypeContent,
		Action:    "review",
		Title:     "审核通知",
		Content:   "您的作品已通过审核",
		Variables: []string{"book_title"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(ctx, template)
	require.NoError(t, err)

	// Act
	found, err := repo.GetByID(ctx, template.ID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, template.Type, found.Type)
	assert.Equal(t, template.Action, found.Action)
	assert.Equal(t, template.Title, found.Title)
	assert.Equal(t, template.Content, found.Content)
	assert.Equal(t, template.Language, found.Language)
}

// TestNotificationTemplateRepository_GetByID_NotFound 测试获取不存在的通知模板
func TestNotificationTemplateRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	// Act
	found, err := repo.GetByID(ctx, "507f1f77bcf86cd799439011")

	// Assert
	require.NoError(t, err)
	assert.Nil(t, found)
}

// TestNotificationTemplateRepository_GetByID_InvalidID 测试使用无效ID获取通知模板
func TestNotificationTemplateRepository_GetByID_InvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	// Act
	found, err := repo.GetByID(ctx, "invalid-id")

	// Assert
	require.Error(t, err)
	assert.Nil(t, found)
	assert.Contains(t, err.Error(), "无效的模板ID")
}

// TestNotificationTemplateRepository_Update 测试更新通知模板
func TestNotificationTemplateRepository_Update(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	template := &notification.NotificationTemplate{
		Type:      notification.NotificationTypeReward,
		Action:    "received",
		Title:     "收到打赏",
		Content:   "您收到了一笔打赏",
		Variables: []string{"amount", "sender"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(ctx, template)
	require.NoError(t, err)

	// Act - 更新模板
	updates := map[string]interface{}{
		"title":      "更新后的打赏标题",
		"is_active":  false,
		"updated_at": time.Now(),
	}
	err = repo.Update(ctx, template.ID, updates)

	// Assert
	require.NoError(t, err)

	// 验证更新
	found, err := repo.GetByID(ctx, template.ID)
	require.NoError(t, err)
	assert.Equal(t, "更新后的打赏标题", found.Title)
	assert.False(t, found.IsActive)
}

// TestNotificationTemplateRepository_Update_NotFound 测试更新不存在的通知模板
func TestNotificationTemplateRepository_Update_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	// Act
	updates := map[string]interface{}{
		"title": "更新后的标题",
	}
	err := repo.Update(ctx, "507f1f77bcf86cd799439011", updates)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "通知模板不存在")
}

// TestNotificationTemplateRepository_Delete 测试删除通知模板
func TestNotificationTemplateRepository_Delete(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	template := &notification.NotificationTemplate{
		Type:      notification.NotificationTypeMessage,
		Action:    "new_message",
		Title:     "新消息",
		Content:   "您有一条新消息",
		Variables: []string{"sender", "content"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(ctx, template)
	require.NoError(t, err)

	// Act - 删除模板
	err = repo.Delete(ctx, template.ID)

	// Assert
	require.NoError(t, err)

	// 验证已删除
	found, err := repo.GetByID(ctx, template.ID)
	require.NoError(t, err)
	assert.Nil(t, found)
}

// TestNotificationTemplateRepository_Delete_NotFound 测试删除不存在的通知模板
func TestNotificationTemplateRepository_Delete_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	// Act
	err := repo.Delete(ctx, "507f1f77bcf86cd799439011")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "通知模板不存在")
}

// TestNotificationTemplateRepository_Exists 测试检查通知模板是否存在
func TestNotificationTemplateRepository_Exists(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	template := &notification.NotificationTemplate{
		Type:      notification.NotificationTypeUpdate,
		Action:    "new_chapter",
		Title:     "新章节",
		Content:   "关注的作品更新了",
		Variables: []string{"book_title", "chapter_title"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(ctx, template)
	require.NoError(t, err)

	// Act
	exists, err := repo.Exists(ctx, template.ID)

	// Assert
	require.NoError(t, err)
	assert.True(t, exists)
}

// TestNotificationTemplateRepository_Exists_NotFound 测试检查不存在的通知模板
func TestNotificationTemplateRepository_Exists_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	// Act
	exists, err := repo.Exists(ctx, "507f1f77bcf86cd799439011")

	// Assert
	require.NoError(t, err)
	assert.False(t, exists)
}

// TestNotificationTemplateRepository_List 测试获取通知模板列表
func TestNotificationTemplateRepository_List(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	// 创建多个模板
	actions := []string{"follow", "like", "comment"}
	for i, action := range actions {
		template := &notification.NotificationTemplate{
			Type:      notification.NotificationTypeSocial,
			Action:    action,
			Title:     func() string { return "社交动作" + string(rune('0'+i)) }(),
			Content:   "这是一个社交通知模板",
			Variables: []string{"username"},
			Language:  "zh-CN",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, template)
		require.NoError(t, err)
	}

	// Act
	filter := &repoInterfaces.TemplateFilter{
		Limit:  10,
		Offset: 0,
	}
	templates, err := repo.List(ctx, filter)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, templates)
	assert.GreaterOrEqual(t, len(templates), 3)
}

// TestNotificationTemplateRepository_List_WithTypeFilter 测试带类型过滤的模板列表
func TestNotificationTemplateRepository_List_WithTypeFilter(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	notifType := notification.NotificationTypeSystem

	// 创建指定类型的模板
	for i := 0; i < 3; i++ {
		template := &notification.NotificationTemplate{
			Type:      notifType,
			Action:    func() string { return "action_" + string(rune('0'+i)) }(),
			Title:     "系统通知",
			Content:   "这是一个系统通知模板",
			Variables: []string{"content"},
			Language:  "zh-CN",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, template)
		require.NoError(t, err)
	}

	// 创建其他类型的模板
	otherTemplate := &notification.NotificationTemplate{
		Type:      notification.NotificationTypeSocial,
		Action:    "follow",
		Title:     "关注通知",
		Content:   "这是一个社交通知模板",
		Variables: []string{"username"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(ctx, otherTemplate)
	require.NoError(t, err)

	// Act - 使用类型过滤
	filter := &repoInterfaces.TemplateFilter{
		Type:   &notifType,
		Limit:  10,
		Offset: 0,
	}
	templates, err := repo.List(ctx, filter)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, templates)
	assert.GreaterOrEqual(t, len(templates), 3)

	// 验证所有返回的模板都是指定类型
	for _, template := range templates {
		assert.Equal(t, notifType, template.Type)
	}
}

// TestNotificationTemplateRepository_List_WithLanguageFilter 测试带语言过滤的模板列表
func TestNotificationTemplateRepository_List_WithLanguageFilter(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	language := "en-US"

	// 创建指定语言的模板
	template := &notification.NotificationTemplate{
		Type:      notification.NotificationTypeSystem,
		Action:    "announcement_en",
		Title:     "Announcement",
		Content:   "This is an announcement template",
		Variables: []string{"title", "content"},
		Language:  language,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(ctx, template)
	require.NoError(t, err)

	// 创建其他语言的模板
	otherTemplate := &notification.NotificationTemplate{
		Type:      notification.NotificationTypeSystem,
		Action:    "announcement_zh",
		Title:     "公告",
		Content:   "这是一个公告模板",
		Variables: []string{"title", "content"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = repo.Create(ctx, otherTemplate)
	require.NoError(t, err)

	// Act - 使用语言过滤
	filter := &repoInterfaces.TemplateFilter{
		Language: &language,
		Limit:    10,
		Offset:   0,
	}
	templates, err := repo.List(ctx, filter)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, templates)
	assert.GreaterOrEqual(t, len(templates), 1)

	// 验证所有返回的模板都是指定语言
	for _, template := range templates {
		assert.Equal(t, language, template.Language)
	}
}

// TestNotificationTemplateRepository_List_WithActiveFilter 测试带活跃状态过滤的模板列表
func TestNotificationTemplateRepository_List_WithActiveFilter(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	isActive := true

	// 创建活跃模板
	activeTemplate := &notification.NotificationTemplate{
		Type:      notification.NotificationTypeContent,
		Action:    "review_approved",
		Title:     "审核通过",
		Content:   "您的作品已通过审核",
		Variables: []string{"book_title"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(ctx, activeTemplate)
	require.NoError(t, err)

	// 创建不活跃模板
	inactiveTemplate := &notification.NotificationTemplate{
		Type:      notification.NotificationTypeContent,
		Action:    "review_rejected",
		Title:     "审核拒绝",
		Content:   "您的作品未通过审核",
		Variables: []string{"book_title", "reason"},
		Language:  "zh-CN",
		IsActive:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = repo.Create(ctx, inactiveTemplate)
	require.NoError(t, err)

	// Act - 使用活跃状态过滤
	filter := &repoInterfaces.TemplateFilter{
		IsActive: &isActive,
		Limit:    10,
		Offset:   0,
	}
	templates, err := repo.List(ctx, filter)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, templates)
	assert.GreaterOrEqual(t, len(templates), 1)

	// 验证所有返回的模板都是活跃的
	for _, template := range templates {
		assert.True(t, template.IsActive)
	}
}

// TestNotificationTemplateRepository_GetByTypeAndAction 测试根据类型和操作获取通知模板
func TestNotificationTemplateRepository_GetByTypeAndAction(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	templateType := notification.NotificationTypeSocial
	action := "like"

	// 创建指定类型和操作的模板
	template := &notification.NotificationTemplate{
		Type:      templateType,
		Action:    action,
		Title:     "点赞通知",
		Content:   "有人赞了你的作品",
		Variables: []string{"username", "book_title"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(ctx, template)
	require.NoError(t, err)

	// Act
	templates, err := repo.GetByTypeAndAction(ctx, templateType, action)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, templates)
	assert.GreaterOrEqual(t, len(templates), 1)
}

// TestNotificationTemplateRepository_GetActiveTemplate 测试获取活跃的通知模板
func TestNotificationTemplateRepository_GetActiveTemplate(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	templateType := notification.NotificationTypeReward
	action := "received"
	language := "zh-CN"

	// 创建活跃模板
	activeTemplate := &notification.NotificationTemplate{
		Type:      templateType,
		Action:    action,
		Title:     "收到打赏",
		Content:   "您收到了一笔打赏",
		Variables: []string{"amount", "sender"},
		Language:  language,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(ctx, activeTemplate)
	require.NoError(t, err)

	// 创建不活跃模板
	inactiveTemplate := &notification.NotificationTemplate{
		Type:      templateType,
		Action:    action,
		Title:     "收到打赏（旧）",
		Content:   "您收到了一笔打赏",
		Variables: []string{"amount", "sender"},
		Language:  language,
		IsActive:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = repo.Create(ctx, inactiveTemplate)
	require.NoError(t, err)

	// Act
	found, err := repo.GetActiveTemplate(ctx, templateType, action, language)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.True(t, found.IsActive)
	assert.Equal(t, templateType, found.Type)
	assert.Equal(t, action, found.Action)
	assert.Equal(t, language, found.Language)
}

// TestNotificationTemplateRepository_GetActiveTemplate_NotFound 测试获取不存在的活跃模板
func TestNotificationTemplateRepository_GetActiveTemplate_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	// Act
	found, err := repo.GetActiveTemplate(ctx, notification.NotificationTypeSystem, "nonexistent_action", "zh-CN")

	// Assert
	require.NoError(t, err)
	assert.Nil(t, found)
}

// TestNotificationTemplateRepository_List_WithPagination 测试分页获取模板列表
func TestNotificationTemplateRepository_List_WithPagination(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	// 创建多个模板
	for i := 0; i < 15; i++ {
		template := &notification.NotificationTemplate{
			Type:      notification.NotificationTypeSystem,
			Action:    func() string { return "action_" + string(rune('0'+i)) }(),
			Title:     "测试模板",
			Content:   "这是一个测试模板",
			Variables: []string{"content"},
			Language:  "zh-CN",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, template)
		require.NoError(t, err)
	}

	// Act - 第一页
	filter1 := &repoInterfaces.TemplateFilter{
		Limit:  10,
		Offset: 0,
	}
	templates1, err := repo.List(ctx, filter1)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, templates1)
	assert.LessOrEqual(t, len(templates1), 10)

	// Act - 第二页
	filter2 := &repoInterfaces.TemplateFilter{
		Limit:  10,
		Offset: 10,
	}
	templates2, err := repo.List(ctx, filter2)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, templates2)
	assert.GreaterOrEqual(t, len(templates2), 5)
}

// TestNotificationTemplateRepository_List_WithActionFilter 测试带操作过滤的模板列表
func TestNotificationTemplateRepository_List_WithActionFilter(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	action := "comment"

	// 创建指定操作的模板
	for i := 0; i < 3; i++ {
		template := &notification.NotificationTemplate{
			Type:      notification.NotificationTypeSocial,
			Action:    action,
			Title:     "评论通知",
			Content:   "有人评论了你的作品",
			Variables: []string{"username", "content"},
			Language:  "zh-CN",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, template)
		require.NoError(t, err)
	}

	// 创建其他操作的模板
	otherTemplate := &notification.NotificationTemplate{
		Type:      notification.NotificationTypeSocial,
		Action:    "like",
		Title:     "点赞通知",
		Content:   "有人赞了你的作品",
		Variables: []string{"username"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(ctx, otherTemplate)
	require.NoError(t, err)

	// Act - 使用操作过滤
	filter := &repoInterfaces.TemplateFilter{
		Action: &action,
		Limit:  10,
		Offset: 0,
	}
	templates, err := repo.List(ctx, filter)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, templates)
	assert.GreaterOrEqual(t, len(templates), 3)

	// 验证所有返回的模板都是指定操作
	for _, template := range templates {
		assert.Equal(t, action, template.Action)
	}
}

// TestNotificationTemplateRepository_Update_MultipleFields 测试更新多个字段
func TestNotificationTemplateRepository_Update_MultipleFields(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	template := &notification.NotificationTemplate{
		Type:      notification.NotificationTypeMembership,
		Action:    "expiring",
		Title:     "会员即将到期",
		Content:   "您的会员即将到期",
		Variables: []string{"expiry_date"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(ctx, template)
	require.NoError(t, err)

	// Act - 更新多个字段
	newVariables := []string{"expiry_date", "renewal_url"}
	updates := map[string]interface{}{
		"title":     "会员到期提醒",
		"content":   "您的会员即将到期，请及时续费",
		"variables": newVariables,
		"language":  "en-US",
		"is_active": false,
		"updated_at": time.Now(),
	}
	err = repo.Update(ctx, template.ID, updates)

	// Assert
	require.NoError(t, err)

	// 验证更新
	found, err := repo.GetByID(ctx, template.ID)
	require.NoError(t, err)
	assert.Equal(t, "会员到期提醒", found.Title)
	assert.Equal(t, "您的会员即将到期，请及时续费", found.Content)
	assert.Equal(t, newVariables, found.Variables)
	assert.Equal(t, "en-US", found.Language)
	assert.False(t, found.IsActive)
}

// TestNotificationTemplateRepository_Create_WithData 测试创建包含数据的模板
func TestNotificationTemplateRepository_Create_WithData(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	data := map[string]interface{}{
		"priority": "high",
		"sound":    "default",
	}

	template := &notification.NotificationTemplate{
		Type:      notification.NotificationTypeSystem,
		Action:    "urgent",
		Title:     "紧急通知",
		Content:   "这是一条紧急通知",
		Variables: []string{"message"},
		Data:      data,
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	err := repo.Create(ctx, template)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, template.ID)

	// 验证数据已保存
	found, err := repo.GetByID(ctx, template.ID)
	require.NoError(t, err)
	assert.NotNil(t, found.Data)
	assert.Equal(t, "high", found.Data["priority"])
	assert.Equal(t, "default", found.Data["sound"])
}

// TestNotificationTemplateRepository_List_MultipleLanguages 测试多语言模板列表
func TestNotificationTemplateRepository_List_MultipleLanguages(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	languages := []string{"zh-CN", "en-US", "ja-JP"}

	// 创建多语言模板
	for _, language := range languages {
		template := &notification.NotificationTemplate{
			Type:      notification.NotificationTypeSystem,
			Action:    "welcome",
			Title:     "欢迎",
			Content:   "欢迎使用青鱼阅读",
			Variables: []string{"username"},
			Language:  language,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, template)
		require.NoError(t, err)
	}

	// Act
	filter := &repoInterfaces.TemplateFilter{
		Limit:  10,
		Offset: 0,
	}
	templates, err := repo.List(ctx, filter)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, templates)
	assert.GreaterOrEqual(t, len(templates), 3)
}
