package reading_test

import (
	"Qingyu_backend/models/reading/reader"
	readingInterfaces "Qingyu_backend/repository/interfaces/reading"
	"Qingyu_backend/repository/mongodb/reading"
	"Qingyu_backend/test/testutil"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 测试辅助函数
func setupReadingSettingsRepo(t *testing.T) (readingInterfaces.ReadingSettingsRepository, context.Context, func()) {
	db, cleanup := testutil.SetupTestDB(t)
	repo := reading.NewMongoReadingSettingsRepository(db)
	ctx := context.Background()
	return repo, ctx, cleanup
}

func createTestReadingSettings(userID string) *reader.ReadingSettings {
	return &reader.ReadingSettings{
		ID:          userID + "_settings",
		UserID:      userID,
		Theme:       "light",
		FontSize:    16,
		FontFamily:  "sans-serif",
		LineHeight:  1.5,
		Background:  "#FFFFFF",
		PageMode:    1, // 滑动
		AutoScroll:  false,
		ScrollSpeed: 50,
	}
}

// 1. 测试创建阅读设置
func TestReadingSettingsRepository_Create(t *testing.T) {
	repo, ctx, cleanup := setupReadingSettingsRepo(t)
	defer cleanup()

	settings := createTestReadingSettings("user123")
	err := repo.Create(ctx, settings)

	require.NoError(t, err)
	assert.NotZero(t, settings.CreatedAt)
	assert.NotZero(t, settings.UpdatedAt)
}

// 2. 测试创建空设置
func TestReadingSettingsRepository_Create_Nil(t *testing.T) {
	repo, ctx, cleanup := setupReadingSettingsRepo(t)
	defer cleanup()

	err := repo.Create(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "阅读设置对象不能为空")
}

// 3. 测试根据ID获取
func TestReadingSettingsRepository_GetByID(t *testing.T) {
	repo, ctx, cleanup := setupReadingSettingsRepo(t)
	defer cleanup()

	// 创建设置
	settings := createTestReadingSettings("user123")
	err := repo.Create(ctx, settings)
	require.NoError(t, err)

	// 获取设置
	retrieved, err := repo.GetByID(ctx, settings.ID)
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, settings.ID, retrieved.ID)
	assert.Equal(t, settings.UserID, retrieved.UserID)
	assert.Equal(t, "light", retrieved.Theme)
}

// 4. 测试获取不存在的设置
func TestReadingSettingsRepository_GetByID_NotFound(t *testing.T) {
	repo, ctx, cleanup := setupReadingSettingsRepo(t)
	defer cleanup()

	retrieved, err := repo.GetByID(ctx, "non_existent_id")
	require.NoError(t, err)
	assert.Nil(t, retrieved)
}

// 5. 测试根据UserID获取
func TestReadingSettingsRepository_GetByUserID(t *testing.T) {
	repo, ctx, cleanup := setupReadingSettingsRepo(t)
	defer cleanup()

	userID := "user123"
	settings := createTestReadingSettings(userID)
	err := repo.Create(ctx, settings)
	require.NoError(t, err)

	// 根据UserID获取
	retrieved, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, userID, retrieved.UserID)
	assert.Equal(t, "light", retrieved.Theme)
}

// 6. 测试根据UserID获取不存在的设置
func TestReadingSettingsRepository_GetByUserID_NotFound(t *testing.T) {
	repo, ctx, cleanup := setupReadingSettingsRepo(t)
	defer cleanup()

	retrieved, err := repo.GetByUserID(ctx, "non_existent_user")
	require.NoError(t, err)
	assert.Nil(t, retrieved)
}

// 7. 测试更新设置
func TestReadingSettingsRepository_Update(t *testing.T) {
	repo, ctx, cleanup := setupReadingSettingsRepo(t)
	defer cleanup()

	// 创建设置
	settings := createTestReadingSettings("user123")
	err := repo.Create(ctx, settings)
	require.NoError(t, err)

	// 更新设置
	updates := map[string]interface{}{
		"theme":       "dark",
		"font_size":   18,
		"line_height": 2.0,
	}
	err = repo.Update(ctx, settings.ID, updates)
	require.NoError(t, err)

	// 验证更新
	retrieved, err := repo.GetByID(ctx, settings.ID)
	require.NoError(t, err)
	assert.Equal(t, "dark", retrieved.Theme)
	assert.Equal(t, 18, retrieved.FontSize)
	assert.Equal(t, 2.0, retrieved.LineHeight)
}

// 8. 测试根据UserID更新
func TestReadingSettingsRepository_UpdateByUserID(t *testing.T) {
	repo, ctx, cleanup := setupReadingSettingsRepo(t)
	defer cleanup()

	userID := "user123"
	settings := createTestReadingSettings(userID)
	err := repo.Create(ctx, settings)
	require.NoError(t, err)

	// 根据UserID更新
	newSettings := createTestReadingSettings(userID)
	newSettings.Theme = "dark"
	newSettings.FontSize = 20
	err = repo.UpdateByUserID(ctx, userID, newSettings)
	require.NoError(t, err)

	// 验证更新
	retrieved, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, "dark", retrieved.Theme)
	assert.Equal(t, 20, retrieved.FontSize)
}

// 9. 测试删除设置
func TestReadingSettingsRepository_Delete(t *testing.T) {
	repo, ctx, cleanup := setupReadingSettingsRepo(t)
	defer cleanup()

	// 创建设置
	settings := createTestReadingSettings("user123")
	err := repo.Create(ctx, settings)
	require.NoError(t, err)

	// 删除设置
	err = repo.Delete(ctx, settings.ID)
	require.NoError(t, err)

	// 验证已删除
	retrieved, err := repo.GetByID(ctx, settings.ID)
	require.NoError(t, err)
	assert.Nil(t, retrieved)
}

// 10. 测试Exists方法
func TestReadingSettingsRepository_Exists(t *testing.T) {
	repo, ctx, cleanup := setupReadingSettingsRepo(t)
	defer cleanup()

	// 创建设置
	settings := createTestReadingSettings("user123")
	err := repo.Create(ctx, settings)
	require.NoError(t, err)

	// 检查存在
	exists, err := repo.Exists(ctx, settings.ID)
	require.NoError(t, err)
	assert.True(t, exists)

	// 检查不存在
	exists, err = repo.Exists(ctx, "non_existent_id")
	require.NoError(t, err)
	assert.False(t, exists)
}

// 11. 测试列表查询
func TestReadingSettingsRepository_List(t *testing.T) {
	repo, ctx, cleanup := setupReadingSettingsRepo(t)
	defer cleanup()

	// 创建多个设置
	for i := 1; i <= 3; i++ {
		userID := "user" + string(rune('0'+i))
		settings := createTestReadingSettings(userID)
		err := repo.Create(ctx, settings)
		require.NoError(t, err)
	}

	// 查询列表
	filter := &testutil.SimpleFilter{
		Conditions: map[string]interface{}{},
	}
	list, err := repo.List(ctx, filter)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 3)
}

// 12. 测试统计数量
func TestReadingSettingsRepository_Count(t *testing.T) {
	repo, ctx, cleanup := setupReadingSettingsRepo(t)
	defer cleanup()

	// 创建多个设置
	for i := 1; i <= 5; i++ {
		userID := "user" + string(rune('0'+i))
		settings := createTestReadingSettings(userID)
		err := repo.Create(ctx, settings)
		require.NoError(t, err)
	}

	// 统计数量
	filter := &testutil.SimpleFilter{
		Conditions: map[string]interface{}{},
	}
	count, err := repo.Count(ctx, filter)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(5))
}

// 13. 测试ExistsByUserID方法
func TestReadingSettingsRepository_ExistsByUserID(t *testing.T) {
	repo, ctx, cleanup := setupReadingSettingsRepo(t)
	defer cleanup()

	userID := "user123"
	settings := createTestReadingSettings(userID)
	err := repo.Create(ctx, settings)
	require.NoError(t, err)

	// 检查用户设置存在
	exists, err := repo.ExistsByUserID(ctx, userID)
	require.NoError(t, err)
	assert.True(t, exists)

	// 检查不存在的用户
	exists, err = repo.ExistsByUserID(ctx, "non_existent_user")
	require.NoError(t, err)
	assert.False(t, exists)
}

// 14. 测试批量创建和查询
func TestReadingSettingsRepository_BatchOperations(t *testing.T) {
	repo, ctx, cleanup := setupReadingSettingsRepo(t)
	defer cleanup()

	// 批量创建
	userIDs := []string{"user1", "user2", "user3"}
	for _, userID := range userIDs {
		settings := createTestReadingSettings(userID)
		settings.Theme = "dark" // 统一主题
		err := repo.Create(ctx, settings)
		require.NoError(t, err)
	}

	// 验证每个都能查到
	for _, userID := range userIDs {
		retrieved, err := repo.GetByUserID(ctx, userID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, "dark", retrieved.Theme)
	}
}

// 15. 测试CreateDefaultSettings方法
func TestReadingSettingsRepository_CreateDefaultSettings(t *testing.T) {
	repo, ctx, cleanup := setupReadingSettingsRepo(t)
	defer cleanup()

	userID := "new_user"

	// 创建默认设置
	settings, err := repo.CreateDefaultSettings(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, settings)
	assert.Equal(t, userID, settings.UserID)
	assert.Equal(t, "light", settings.Theme)
	assert.Equal(t, 16, settings.FontSize)

	// 验证可以通过UserID查询到
	retrieved, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, userID, retrieved.UserID)
}
