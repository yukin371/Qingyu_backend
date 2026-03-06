package auth_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"

	authModel "Qingyu_backend/models/auth"
	repoAuth "Qingyu_backend/repository/interfaces/auth"
	authMongo "Qingyu_backend/repository/mongodb/auth"
	"Qingyu_backend/test/testutil"
)

// TestPermissionTemplateRepository_CreateTemplate 测试创建模板
func TestPermissionTemplateRepository_CreateTemplate(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := authMongo.NewPermissionTemplateRepositoryMongo(db.Client(), db.Name(), nil)
	ctx := context.Background()

	template := &authModel.PermissionTemplate{
		Name:        "测试模板",
		Code:        "test_template",
		Description: "测试用模板",
		Permissions: []string{"user.read", "book.read"},
		IsSystem:    false,
		Category:    authModel.CategoryCustom,
	}

	// 执行
	err := repo.CreateTemplate(ctx, template)

	// 验证
	require.NoError(t, err)
	assert.NotEmpty(t, template.ID)

	// 从数据库验证
	stored, err := repo.GetTemplateByCode(ctx, "test_template")
	require.NoError(t, err)
	assert.Equal(t, template.Name, stored.Name)
	assert.Equal(t, template.Code, stored.Code)
	assert.Equal(t, template.Permissions, stored.Permissions)

	// 清理
	_ = repo.DeleteTemplate(ctx, template.ID)
}

// TestPermissionTemplateRepository_CreateSystemTemplate 测试创建系统模板
func TestPermissionTemplateRepository_CreateSystemTemplate(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := authMongo.NewPermissionTemplateRepositoryMongo(db.Client(), db.Name(), nil)
	ctx := context.Background()

	template := &authModel.PermissionTemplate{
		Name:        "系统测试模板",
		Code:        "system_test_template",
		Description: "系统测试用模板",
		Permissions: []string{"user.read", "user.write"},
		IsSystem:    true,
		Category:    authModel.CategoryAdmin,
	}

	err := repo.CreateTemplate(ctx, template)
	require.NoError(t, err)
	assert.NotEmpty(t, template.ID)

	// 验证是系统模板
	stored, err := repo.GetTemplateByID(ctx, template.ID)
	require.NoError(t, err)
	assert.True(t, stored.IsSystem)

	// 清理
	_ = repo.DeleteTemplate(ctx, template.ID)
}

// TestPermissionTemplateRepository_CodeUnique 测试模板代码唯一性
func TestPermissionTemplateRepository_CodeUnique(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := authMongo.NewPermissionTemplateRepositoryMongo(db.Client(), db.Name(), nil)
	ctx := context.Background()

	// 创建第一个模板
	template1 := &authModel.PermissionTemplate{
		Name:        "模板1",
		Code:        "unique_code",
		Description: "第一个模板",
		Permissions: []string{"user.read"},
		IsSystem:    false,
		Category:    authModel.CategoryCustom,
	}

	err := repo.CreateTemplate(ctx, template1)
	require.NoError(t, err)

	// 尝试创建相同代码的模板
	template2 := &authModel.PermissionTemplate{
		Name:        "模板2",
		Code:        "unique_code",
		Description: "第二个模板",
		Permissions: []string{"book.read"},
		IsSystem:    false,
		Category:    authModel.CategoryCustom,
	}

	err = repo.CreateTemplate(ctx, template2)
	assert.Error(t, err)

	// 清理
	_ = repo.DeleteTemplate(ctx, template1.ID)
}

// TestPermissionTemplateRepository_GetTemplateByID 测试根据ID获取模板
func TestPermissionTemplateRepository_GetTemplateByID(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := authMongo.NewPermissionTemplateRepositoryMongo(db.Client(), db.Name(), nil)
	ctx := context.Background()

	// 创建模板
	template := &authModel.PermissionTemplate{
		Name:        "测试模板",
		Code:        "get_by_id_test",
		Description: "测试用",
		Permissions: []string{"user.read"},
		IsSystem:    false,
		Category:    authModel.CategoryCustom,
	}

	err := repo.CreateTemplate(ctx, template)
	require.NoError(t, err)

	// 获取模板
	stored, err := repo.GetTemplateByID(ctx, template.ID)
	require.NoError(t, err)
	assert.Equal(t, template.Name, stored.Name)
	assert.Equal(t, template.Code, stored.Code)

	// 测试不存在的ID
	_, err = repo.GetTemplateByID(ctx, "nonexistent")
	assert.Error(t, err)

	// 清理
	_ = repo.DeleteTemplate(ctx, template.ID)
}

// TestPermissionTemplateRepository_UpdateTemplate 测试更新模板
func TestPermissionTemplateRepository_UpdateTemplate(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := authMongo.NewPermissionTemplateRepositoryMongo(db.Client(), db.Name(), nil)
	ctx := context.Background()

	// 创建模板
	template := &authModel.PermissionTemplate{
		Name:        "原始名称",
		Code:        "update_test",
		Description: "原始描述",
		Permissions: []string{"user.read"},
		IsSystem:    false,
		Category:    authModel.CategoryCustom,
	}

	err := repo.CreateTemplate(ctx, template)
	require.NoError(t, err)

	// 更新模板
	updates := map[string]interface{}{
		"name":        "新名称",
		"description": "新描述",
		"permissions": []string{"user.read", "user.write"},
	}

	err = repo.UpdateTemplate(ctx, template.ID, updates)
	require.NoError(t, err)

	// 验证更新
	stored, err := repo.GetTemplateByID(ctx, template.ID)
	require.NoError(t, err)
	assert.Equal(t, "新名称", stored.Name)
	assert.Equal(t, "新描述", stored.Description)
	assert.Equal(t, []string{"user.read", "user.write"}, stored.Permissions)

	// 清理
	_ = repo.DeleteTemplate(ctx, template.ID)
}

// TestPermissionTemplateRepository_UpdateSystemTemplate 测试不能更新系统模板
func TestPermissionTemplateRepository_UpdateSystemTemplate(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := authMongo.NewPermissionTemplateRepositoryMongo(db.Client(), db.Name(), nil)
	ctx := context.Background()

	// 创建系统模板
	template := &authModel.PermissionTemplate{
		Name:        "系统模板",
		Code:        "system_template",
		Description: "系统模板描述",
		Permissions: []string{"user.read"},
		IsSystem:    true,
		Category:    authModel.CategoryAdmin,
	}

	err := repo.CreateTemplate(ctx, template)
	require.NoError(t, err)

	// 尝试更新系统模板
	updates := map[string]interface{}{
		"name": "新名称",
	}

	err = repo.UpdateTemplate(ctx, template.ID, updates)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "系统模板")

	// 清理
	_ = repo.DeleteTemplate(ctx, template.ID)
}

// TestPermissionTemplateRepository_DeleteTemplate 测试删除模板
func TestPermissionTemplateRepository_DeleteTemplate(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := authMongo.NewPermissionTemplateRepositoryMongo(db.Client(), db.Name(), nil)
	ctx := context.Background()

	// 创建模板
	template := &authModel.PermissionTemplate{
		Name:        "待删除模板",
		Code:        "delete_test",
		Description: "将被删除",
		Permissions: []string{"user.read"},
		IsSystem:    false,
		Category:    authModel.CategoryCustom,
	}

	err := repo.CreateTemplate(ctx, template)
	require.NoError(t, err)

	// 删除模板
	err = repo.DeleteTemplate(ctx, template.ID)
	require.NoError(t, err)

	// 验证已删除
	_, err = repo.GetTemplateByID(ctx, template.ID)
	assert.Error(t, err)
}

// TestPermissionTemplateRepository_DeleteSystemTemplate 测试不能删除系统模板
func TestPermissionTemplateRepository_DeleteSystemTemplate(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := authMongo.NewPermissionTemplateRepositoryMongo(db.Client(), db.Name(), nil)
	ctx := context.Background()

	// 创建系统模板
	template := &authModel.PermissionTemplate{
		Name:        "系统模板",
		Code:        "system_delete_test",
		Description: "系统模板",
		Permissions: []string{"user.read"},
		IsSystem:    true,
		Category:    authModel.CategoryAdmin,
	}

	err := repo.CreateTemplate(ctx, template)
	require.NoError(t, err)

	// 尝试删除系统模板
	err = repo.DeleteTemplate(ctx, template.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "系统模板")

	// 清理
	_ = repo.DeleteTemplate(ctx, template.ID)
}

// TestPermissionTemplateRepository_ListTemplates 测试列出所有模板
func TestPermissionTemplateRepository_ListTemplates(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := authMongo.NewPermissionTemplateRepositoryMongo(db.Client(), db.Name(), nil)
	ctx := context.Background()

	// 创建多个模板
	templates := []*authModel.PermissionTemplate{
		{
			Name:        "模板1",
			Code:        "list_test_1",
			Description: "第一个",
			Permissions: []string{"user.read"},
			IsSystem:    false,
			Category:    authModel.CategoryReader,
		},
		{
			Name:        "模板2",
			Code:        "list_test_2",
			Description: "第二个",
			Permissions: []string{"book.read"},
			IsSystem:    false,
			Category:    authModel.CategoryAuthor,
		},
	}

	for _, tmpl := range templates {
		err := repo.CreateTemplate(ctx, tmpl)
		require.NoError(t, err)
	}

	// 列出所有模板
	result, err := repo.ListTemplates(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(result), 2)

	// 清理
	for _, tmpl := range templates {
		_ = repo.DeleteTemplate(ctx, tmpl.ID)
	}
}

// TestPermissionTemplateRepository_ListTemplatesByCategory 测试按分类列出模板
func TestPermissionTemplateRepository_ListTemplatesByCategory(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := authMongo.NewPermissionTemplateRepositoryMongo(db.Client(), db.Name(), nil)
	ctx := context.Background()

	// 创建不同分类的模板
	templates := []*authModel.PermissionTemplate{
		{
			Name:        "读者模板1",
			Code:        "reader_test_1",
			Description: "读者模板",
			Permissions: []string{"book.read"},
			IsSystem:    false,
			Category:    authModel.CategoryReader,
		},
		{
			Name:        "读者模板2",
			Code:        "reader_test_2",
			Description: "读者模板2",
			Permissions: []string{"document.read"},
			IsSystem:    false,
			Category:    authModel.CategoryReader,
		},
		{
			Name:        "作者模板",
			Code:        "author_test",
			Description: "作者模板",
			Permissions: []string{"book.write"},
			IsSystem:    false,
			Category:    authModel.CategoryAuthor,
		},
	}

	for _, tmpl := range templates {
		err := repo.CreateTemplate(ctx, tmpl)
		require.NoError(t, err)
	}

	// 按分类列出
	readerTemplates, err := repo.ListTemplatesByCategory(ctx, authModel.CategoryReader)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(readerTemplates), 2)

	authorTemplates, err := repo.ListTemplatesByCategory(ctx, authModel.CategoryAuthor)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(authorTemplates), 1)

	// 清理
	for _, tmpl := range templates {
		_ = repo.DeleteTemplate(ctx, tmpl.ID)
	}
}

// TestPermissionTemplateRepository_InitializeSystemTemplates 测试初始化系统模板
func TestPermissionTemplateRepository_InitializeSystemTemplates(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := authMongo.NewPermissionTemplateRepositoryMongo(db.Client(), db.Name(), nil)
	ctx := context.Background()

	// 初始化系统模板
	err := repo.InitializeSystemTemplates(ctx)
	require.NoError(t, err)

	// 验证系统模板已创建
	readerTemplate, err := repo.GetTemplateByCode(ctx, authModel.TemplateReader)
	require.NoError(t, err)
	assert.NotNil(t, readerTemplate)
	assert.True(t, readerTemplate.IsSystem)
	assert.Equal(t, authModel.CategoryReader, readerTemplate.Category)
	assert.NotEmpty(t, readerTemplate.Permissions)

	authorTemplate, err := repo.GetTemplateByCode(ctx, authModel.TemplateAuthor)
	require.NoError(t, err)
	assert.NotNil(t, authorTemplate)
	assert.True(t, authorTemplate.IsSystem)
	assert.Equal(t, authModel.CategoryAuthor, authorTemplate.Category)

	adminTemplate, err := repo.GetTemplateByCode(ctx, authModel.TemplateAdmin)
	require.NoError(t, err)
	assert.NotNil(t, adminTemplate)
	assert.True(t, adminTemplate.IsSystem)
	assert.Equal(t, authModel.CategoryAdmin, adminTemplate.Category)
}

// ============ 辅助函数 ============

// setupPermissionTemplateTestRepo 设置权限模板测试仓储
func setupPermissionTemplateTestRepo(t *testing.T, redisClient *redis.Client) (repoAuth.PermissionTemplateRepository, *mongo.Database, func()) {
	db, cleanup := testutil.SetupTestDB(t)
	repo := authMongo.NewPermissionTemplateRepositoryMongo(db.Client(), db.Name(), redisClient)

	// 在清理函数中添加清除权限模板集合
	enhancedCleanup := func() {
		ctx := context.Background()
		_ = db.Collection("permission_templates").Drop(ctx)
		cleanup()
	}

	return repo, db, enhancedCleanup
}

// createTestRoleForTemplate 创建测试角色
func createTestRoleForTemplate(t *testing.T, db *mongo.Database, name string) string {
	ctx := context.Background()
	collection := db.Collection("roles")

	// 简单创建一个测试角色
	role := map[string]interface{}{
		"name":        name,
		"permissions": []string{},
		"is_system":   false,
		"created_at":  "2024-01-01T00:00:00Z",
		"updated_at":  "2024-01-01T00:00:00Z",
	}

	result, err := collection.InsertOne(ctx, role)
	if err != nil {
		t.Fatalf("创建测试角色失败: %v", err)
	}

	if oid, ok := result.InsertedID.(string); ok {
		return oid
	}
	return result.InsertedID.(string)
}

// cleanupTestRoleForTemplate 清理测试角色
func cleanupTestRoleForTemplate(t *testing.T, db *mongo.Database, roleID string) {
	ctx := context.Background()
	_, _ = db.Collection("roles").DeleteOne(ctx, map[string]interface{}{"_id": roleID})
}
