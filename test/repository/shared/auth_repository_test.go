package shared

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/global"
	authModel "Qingyu_backend/models/shared/auth"
	"Qingyu_backend/repository/mongodb/shared"
	"Qingyu_backend/test/testutil"
)

var repo *shared.AuthRepositoryImpl

func setupTest(t *testing.T) {
	// 设置测试数据库
	testutil.SetupTestDB(t)

	// 创建Repository实例
	repo = shared.NewAuthRepository(global.DB).(*shared.AuthRepositoryImpl)

	// 清理测试数据
	ctx := context.Background()
	_ = global.DB.Collection("roles").Drop(ctx)
	_ = global.DB.Collection("users").Drop(ctx)
}

func createTestRole(name string) *authModel.Role {
	return &authModel.Role{
		Name:        name,
		Description: "Test role: " + name,
		Permissions: []string{"test.read", "test.write"},
		IsSystem:    false,
	}
}

// createTestUserDoc 创建测试用户文档（使用bson.M以支持roles数组）
func createTestUserDoc(username string) bson.M {
	return bson.M{
		"username":   username,
		"email":      username + "@test.com",
		"password":   "hashed_password",
		"status":     "active",
		"roles":      []string{}, // roles数组字段
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}
}

// ============ 角色管理测试 ============

func TestAuthRepository_CreateRole(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	role := createTestRole("test_role")

	err := repo.CreateRole(ctx, role)
	require.NoError(t, err)
	assert.NotEmpty(t, role.ID)
	assert.NotZero(t, role.CreatedAt)
	assert.NotZero(t, role.UpdatedAt)
}

func TestAuthRepository_GetRole(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 先创建角色
	role := createTestRole("test_role")
	err := repo.CreateRole(ctx, role)
	require.NoError(t, err)

	// 查询角色
	result, err := repo.GetRole(ctx, role.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, role.Name, result.Name)
	assert.Equal(t, role.Description, result.Description)
}

func TestAuthRepository_GetRole_NotFound(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 使用不存在的ObjectID
	fakeID := primitive.NewObjectID().Hex()
	result, err := repo.GetRole(ctx, fakeID)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestAuthRepository_GetRole_InvalidID(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	result, err := repo.GetRole(ctx, "invalid_id")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestAuthRepository_GetRoleByName(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建角色
	role := createTestRole("admin")
	err := repo.CreateRole(ctx, role)
	require.NoError(t, err)

	// 按名称查询
	result, err := repo.GetRoleByName(ctx, "admin")
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "admin", result.Name)
}

func TestAuthRepository_GetRoleByName_NotFound(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	result, err := repo.GetRoleByName(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestAuthRepository_UpdateRole(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建角色
	role := createTestRole("test_role")
	err := repo.CreateRole(ctx, role)
	require.NoError(t, err)

	// 更新角色
	updates := map[string]interface{}{
		"description": "Updated description",
		"permissions": []string{"new.permission"},
	}
	err = repo.UpdateRole(ctx, role.ID, updates)
	require.NoError(t, err)

	// 验证更新
	result, err := repo.GetRole(ctx, role.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated description", result.Description)
}

func TestAuthRepository_UpdateRole_NotFound(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	fakeID := primitive.NewObjectID().Hex()
	updates := map[string]interface{}{
		"description": "Updated description",
	}
	err := repo.UpdateRole(ctx, fakeID, updates)
	assert.Error(t, err)
}

func TestAuthRepository_DeleteRole(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建角色
	role := createTestRole("test_role")
	err := repo.CreateRole(ctx, role)
	require.NoError(t, err)

	// 删除角色
	err = repo.DeleteRole(ctx, role.ID)
	require.NoError(t, err)

	// 验证已删除
	result, err := repo.GetRole(ctx, role.ID)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestAuthRepository_DeleteRole_SystemRole(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建系统角色
	role := createTestRole("system_role")
	role.IsSystem = true
	err := repo.CreateRole(ctx, role)
	require.NoError(t, err)

	// 尝试删除系统角色
	err = repo.DeleteRole(ctx, role.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "不能删除系统角色")
}

func TestAuthRepository_ListRoles(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建多个角色
	role1 := createTestRole("role1")
	role2 := createTestRole("role2")
	role3 := createTestRole("role3")

	err := repo.CreateRole(ctx, role1)
	require.NoError(t, err)
	time.Sleep(10 * time.Millisecond)

	err = repo.CreateRole(ctx, role2)
	require.NoError(t, err)
	time.Sleep(10 * time.Millisecond)

	err = repo.CreateRole(ctx, role3)
	require.NoError(t, err)

	// 查询所有角色
	results, err := repo.ListRoles(ctx)
	require.NoError(t, err)
	assert.Len(t, results, 3)
}

// ============ 用户角色关联测试 ============

func TestAuthRepository_AssignUserRole(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建测试角色
	role := createTestRole("reader")
	err := repo.CreateRole(ctx, role)
	require.NoError(t, err)

	// 创建测试用户
	userDoc := createTestUserDoc("testuser")
	result, err := global.DB.Collection("users").InsertOne(ctx, userDoc)
	require.NoError(t, err)
	userID := result.InsertedID.(primitive.ObjectID).Hex()

	// 分配角色
	err = repo.AssignUserRole(ctx, userID, role.ID)
	require.NoError(t, err)

	// 验证角色已分配
	hasRole, err := repo.HasUserRole(ctx, userID, role.ID)
	require.NoError(t, err)
	assert.True(t, hasRole)
}

func TestAuthRepository_AssignUserRole_InvalidRole(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建测试用户
	userDoc := createTestUserDoc("testuser")
	result, err := global.DB.Collection("users").InsertOne(ctx, userDoc)
	require.NoError(t, err)
	userID := result.InsertedID.(primitive.ObjectID).Hex()

	// 尝试分配不存在的角色
	fakeRoleID := primitive.NewObjectID().Hex()
	err = repo.AssignUserRole(ctx, userID, fakeRoleID)
	assert.Error(t, err)
}

func TestAuthRepository_RemoveUserRole(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建角色和用户
	role := createTestRole("reader")
	err := repo.CreateRole(ctx, role)
	require.NoError(t, err)

	userDoc := createTestUserDoc("testuser")
	result, err := global.DB.Collection("users").InsertOne(ctx, userDoc)
	require.NoError(t, err)
	userID := result.InsertedID.(primitive.ObjectID).Hex()

	// 分配角色
	err = repo.AssignUserRole(ctx, userID, role.ID)
	require.NoError(t, err)

	// 移除角色
	err = repo.RemoveUserRole(ctx, userID, role.ID)
	require.NoError(t, err)

	// 验证角色已移除
	hasRole, err := repo.HasUserRole(ctx, userID, role.ID)
	require.NoError(t, err)
	assert.False(t, hasRole)
}

func TestAuthRepository_GetUserRoles(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建多个角色
	role1 := createTestRole("reader")
	role2 := createTestRole("author")
	err := repo.CreateRole(ctx, role1)
	require.NoError(t, err)
	err = repo.CreateRole(ctx, role2)
	require.NoError(t, err)

	// 创建用户
	userDoc := createTestUserDoc("testuser")
	result, err := global.DB.Collection("users").InsertOne(ctx, userDoc)
	require.NoError(t, err)
	userID := result.InsertedID.(primitive.ObjectID).Hex()

	// 分配多个角色
	err = repo.AssignUserRole(ctx, userID, role1.ID)
	require.NoError(t, err)
	err = repo.AssignUserRole(ctx, userID, role2.ID)
	require.NoError(t, err)

	// 查询用户角色
	roles, err := repo.GetUserRoles(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, roles, 2)
}

func TestAuthRepository_GetUserRoles_NoRoles(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建没有角色的用户
	userDoc := createTestUserDoc("testuser")
	result, err := global.DB.Collection("users").InsertOne(ctx, userDoc)
	require.NoError(t, err)
	userID := result.InsertedID.(primitive.ObjectID).Hex()

	// 查询用户角色
	roles, err := repo.GetUserRoles(ctx, userID)
	require.NoError(t, err)
	assert.Empty(t, roles)
}

func TestAuthRepository_HasUserRole(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建角色和用户
	role := createTestRole("reader")
	err := repo.CreateRole(ctx, role)
	require.NoError(t, err)

	userDoc := createTestUserDoc("testuser")
	result, err := global.DB.Collection("users").InsertOne(ctx, userDoc)
	require.NoError(t, err)
	userID := result.InsertedID.(primitive.ObjectID).Hex()

	// 分配角色
	err = repo.AssignUserRole(ctx, userID, role.ID)
	require.NoError(t, err)

	// 检查用户是否有角色
	hasRole, err := repo.HasUserRole(ctx, userID, role.ID)
	require.NoError(t, err)
	assert.True(t, hasRole)

	// 检查用户是否没有其他角色
	fakeRoleID := primitive.NewObjectID().Hex()
	hasRole, err = repo.HasUserRole(ctx, userID, fakeRoleID)
	require.NoError(t, err)
	assert.False(t, hasRole)
}

// ============ 权限查询测试 ============

func TestAuthRepository_GetRolePermissions(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建带权限的角色
	role := createTestRole("admin")
	role.Permissions = []string{"user.read", "user.write", "user.delete"}
	err := repo.CreateRole(ctx, role)
	require.NoError(t, err)

	// 查询角色权限
	permissions, err := repo.GetRolePermissions(ctx, role.ID)
	require.NoError(t, err)
	assert.Len(t, permissions, 3)
	assert.Contains(t, permissions, "user.read")
	assert.Contains(t, permissions, "user.write")
	assert.Contains(t, permissions, "user.delete")
}

func TestAuthRepository_GetUserPermissions(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建多个带权限的角色
	role1 := createTestRole("reader")
	role1.Permissions = []string{"book.read", "chapter.read"}
	err := repo.CreateRole(ctx, role1)
	require.NoError(t, err)

	role2 := createTestRole("author")
	role2.Permissions = []string{"book.write", "book.read"} // book.read重复
	err = repo.CreateRole(ctx, role2)
	require.NoError(t, err)

	// 创建用户并分配角色
	userDoc := createTestUserDoc("testuser")
	result, err := global.DB.Collection("users").InsertOne(ctx, userDoc)
	require.NoError(t, err)
	userID := result.InsertedID.(primitive.ObjectID).Hex()

	err = repo.AssignUserRole(ctx, userID, role1.ID)
	require.NoError(t, err)
	err = repo.AssignUserRole(ctx, userID, role2.ID)
	require.NoError(t, err)

	// 查询用户权限（应该去重）
	permissions, err := repo.GetUserPermissions(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, permissions, 3) // 去重后：book.read, chapter.read, book.write
	assert.Contains(t, permissions, "book.read")
	assert.Contains(t, permissions, "chapter.read")
	assert.Contains(t, permissions, "book.write")
}

func TestAuthRepository_GetUserPermissions_NoRoles(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建没有角色的用户
	userDoc := createTestUserDoc("testuser")
	result, err := global.DB.Collection("users").InsertOne(ctx, userDoc)
	require.NoError(t, err)
	userID := result.InsertedID.(primitive.ObjectID).Hex()

	// 查询用户权限
	permissions, err := repo.GetUserPermissions(ctx, userID)
	require.NoError(t, err)
	assert.Empty(t, permissions)
}

// ============ 健康检查测试 ============

func TestAuthRepository_Health(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	err := repo.Health(ctx)
	assert.NoError(t, err)
}
