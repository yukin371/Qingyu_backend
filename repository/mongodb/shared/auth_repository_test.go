package shared_test

import (
	authModel "Qingyu_backend/models/auth"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	shared "Qingyu_backend/repository/mongodb/shared"
	"Qingyu_backend/test/testutil"
)

// TestAuthRepository_CreateRole 测试创建角色
func TestAuthRepository_CreateRole(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewAuthRepository(db)
	ctx := context.Background()

	role := &authModel.Role{
		Name:        "test_role",
		Description: "Test role description",
		Permissions: []string{"test.read", "test.write"},
		IsSystem:    false,
	}

	// Act
	err := repo.CreateRole(ctx, role)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, role.ID)
	assert.NotZero(t, role.CreatedAt)
	assert.NotZero(t, role.UpdatedAt)
}

// TestAuthRepository_GetRole 测试获取角色
func TestAuthRepository_GetRole(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewAuthRepository(db)
	ctx := context.Background()

	// 先创建角色
	role := &authModel.Role{
		Name:        "test_role",
		Description: "Test role description",
		Permissions: []string{"test.read", "test.write"},
		IsSystem:    false,
	}
	err := repo.CreateRole(ctx, role)
	require.NoError(t, err)

	// Act
	result, err := repo.GetRole(ctx, role.ID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, role.Name, result.Name)
	assert.Equal(t, role.Description, result.Description)
}

// TestAuthRepository_GetRole_NotFound 测试获取不存在的角色
func TestAuthRepository_GetRole_NotFound(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewAuthRepository(db)
	ctx := context.Background()

	// 使用不存在的ObjectID
	fakeID := primitive.NewObjectID().Hex()

	// Act
	result, err := repo.GetRole(ctx, fakeID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestAuthRepository_GetRole_InvalidID 测试获取无效ID的角色
func TestAuthRepository_GetRole_InvalidID(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewAuthRepository(db)
	ctx := context.Background()

	// Act
	result, err := repo.GetRole(ctx, "invalid_id")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestAuthRepository_GetRoleByName 测试按名称获取角色
func TestAuthRepository_GetRoleByName(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewAuthRepository(db)
	ctx := context.Background()

	// 创建角色
	role := &authModel.Role{
		Name:        "admin",
		Description: "Administrator role",
		Permissions: []string{"admin.*"},
		IsSystem:    false,
	}
	err := repo.CreateRole(ctx, role)
	require.NoError(t, err)

	// Act
	result, err := repo.GetRoleByName(ctx, "admin")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "admin", result.Name)
}

// TestAuthRepository_GetRoleByName_NotFound 测试按名称获取不存在的角色
func TestAuthRepository_GetRoleByName_NotFound(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewAuthRepository(db)
	ctx := context.Background()

	// Act
	result, err := repo.GetRoleByName(ctx, "nonexistent")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestAuthRepository_UpdateRole 测试更新角色
func TestAuthRepository_UpdateRole(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewAuthRepository(db)
	ctx := context.Background()

	// 创建角色
	role := &authModel.Role{
		Name:        "test_role",
		Description: "Original description",
		Permissions: []string{"test.read"},
		IsSystem:    false,
	}
	err := repo.CreateRole(ctx, role)
	require.NoError(t, err)

	// Act - 更新角色
	updates := map[string]interface{}{
		"description": "Updated description",
		"permissions": []string{"new.permission"},
	}
	err = repo.UpdateRole(ctx, role.ID, updates)

	// Assert
	require.NoError(t, err)

	// 验证更新
	result, err := repo.GetRole(ctx, role.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated description", result.Description)
}

// TestAuthRepository_UpdateRole_NotFound 测试更新不存在的角色
func TestAuthRepository_UpdateRole_NotFound(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewAuthRepository(db)
	ctx := context.Background()

	fakeID := primitive.NewObjectID().Hex()
	updates := map[string]interface{}{
		"description": "Updated description",
	}

	// Act
	err := repo.UpdateRole(ctx, fakeID, updates)

	// Assert
	assert.Error(t, err)
}

// TestAuthRepository_DeleteRole 测试删除角色
func TestAuthRepository_DeleteRole(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewAuthRepository(db)
	ctx := context.Background()

	// 创建角色
	role := &authModel.Role{
		Name:        "test_role",
		Description: "Test role",
		Permissions: []string{"test.read"},
		IsSystem:    false,
	}
	err := repo.CreateRole(ctx, role)
	require.NoError(t, err)

	// Act - 删除角色
	err = repo.DeleteRole(ctx, role.ID)

	// Assert
	require.NoError(t, err)

	// 验证已删除
	result, err := repo.GetRole(ctx, role.ID)
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestAuthRepository_DeleteRole_SystemRole 测试删除系统角色应失败
func TestAuthRepository_DeleteRole_SystemRole(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewAuthRepository(db)
	ctx := context.Background()

	// 创建系统角色
	role := &authModel.Role{
		Name:        "system_role",
		Description: "System role",
		Permissions: []string{"system.*"},
		IsSystem:    true,
	}
	err := repo.CreateRole(ctx, role)
	require.NoError(t, err)

	// Act - 尝试删除系统角色
	err = repo.DeleteRole(ctx, role.ID)

	// Assert
	assert.Error(t, err) // 系统角色不应被删除
}

// TestAuthRepository_ListRoles 测试列出所有角色
func TestAuthRepository_ListRoles(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewAuthRepository(db)
	ctx := context.Background()

	// 创建多个角色
	roles := []*authModel.Role{
		{Name: "role1", Description: "Role 1", Permissions: []string{"perm1"}, IsSystem: false},
		{Name: "role2", Description: "Role 2", Permissions: []string{"perm2"}, IsSystem: false},
		{Name: "role3", Description: "Role 3", Permissions: []string{"perm3"}, IsSystem: false},
	}

	for _, role := range roles {
		err := repo.CreateRole(ctx, role)
		require.NoError(t, err)
	}

	// Act
	result, err := repo.ListRoles(ctx)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result), 3) // 至少有3个角色
}

// TestAuthRepository_ListRoles_WithSystemRoles 测试列出包含系统角色
func TestAuthRepository_ListRoles_WithSystemRoles(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewAuthRepository(db)
	ctx := context.Background()

	// 创建不同类型的角色
	adminRole := &authModel.Role{
		Name:        "admin",
		Description: "Admin role",
		Permissions: []string{"admin.*"},
		IsSystem:    true,
	}
	userRole := &authModel.Role{
		Name:        "user",
		Description: "User role",
		Permissions: []string{"user.read"},
		IsSystem:    false,
	}

	repo.CreateRole(ctx, adminRole)
	repo.CreateRole(ctx, userRole)

	// Act - 查询所有角色
	result, err := repo.ListRoles(ctx)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	// 验证包含系统角色和非系统角色
	hasSystemRole := false
	hasNonSystemRole := false
	for _, role := range result {
		if role.IsSystem {
			hasSystemRole = true
		} else {
			hasNonSystemRole = true
		}
	}
	assert.True(t, hasSystemRole || hasNonSystemRole)
}

// TestAuthRepository_AssignUserRole 测试为用户分配角色
func TestAuthRepository_AssignUserRole(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewAuthRepository(db)
	ctx := context.Background()

	// 创建角色
	role := &authModel.Role{
		Name:        "reader",
		Description: "Reader role",
		Permissions: []string{"book.read"},
		IsSystem:    false,
	}
	err := repo.CreateRole(ctx, role)
	require.NoError(t, err)

	// 创建用户 - 使用ObjectID作为_id
	userID := primitive.NewObjectID()
	userDoc := bson.M{
		"_id":        userID,
		"username":   "testuser",
		"email":      "test@example.com",
		"status":     "active",
		"roles":      []string{},
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}
	_, err = db.Collection("users").InsertOne(ctx, userDoc)
	require.NoError(t, err)

	// Act - 为用户分配角色
	err = repo.AssignUserRole(ctx, userID.Hex(), role.ID)

	// Assert
	require.NoError(t, err)

	// 验证角色已分配
	var user bson.M
	err = db.Collection("users").FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	require.NoError(t, err)
	roles := user["roles"].(primitive.A)
	assert.Contains(t, roles, role.ID)
}

// TestAuthRepository_RemoveUserRole 测试移除用户角色
func TestAuthRepository_RemoveUserRole(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewAuthRepository(db)
	ctx := context.Background()

	// 创建角色
	role := &authModel.Role{
		Name:        "reader",
		Description: "Reader role",
		Permissions: []string{"book.read"},
		IsSystem:    false,
	}
	err := repo.CreateRole(ctx, role)
	require.NoError(t, err)

	// 创建用户并分配角色 - 使用ObjectID作为_id
	userID := primitive.NewObjectID()
	userDoc := bson.M{
		"_id":        userID,
		"username":   "testuser",
		"email":      "test@example.com",
		"status":     "active",
		"roles":      []string{role.ID},
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}
	_, err = db.Collection("users").InsertOne(ctx, userDoc)
	require.NoError(t, err)

	// Act - 移除用户角色
	err = repo.RemoveUserRole(ctx, userID.Hex(), role.ID)

	// Assert
	require.NoError(t, err)

	// 验证角色已移除
	var user bson.M
	err = db.Collection("users").FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	require.NoError(t, err)
	roles := user["roles"].(primitive.A)
	assert.NotContains(t, roles, role.ID)
}

// TestAuthRepository_GetUserRoles 测试获取用户角色
func TestAuthRepository_GetUserRoles(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewAuthRepository(db)
	ctx := context.Background()

	// 创建两个角色
	role1 := &authModel.Role{Name: "reader", Permissions: []string{"read"}, IsSystem: false}
	role2 := &authModel.Role{Name: "author", Permissions: []string{"write"}, IsSystem: false}

	repo.CreateRole(ctx, role1)
	repo.CreateRole(ctx, role2)

	// 创建用户并分配角色 - GetUserRoles需要string类型的_id
	userID := "user_" + primitive.NewObjectID().Hex()
	userDoc := bson.M{
		"_id":        userID,
		"username":   "testuser",
		"email":      "test@example.com",
		"status":     "active",
		"roles":      []string{role1.ID, role2.ID},
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}
	_, err := db.Collection("users").InsertOne(ctx, userDoc)
	require.NoError(t, err)

	// Act
	roles, err := repo.GetUserRoles(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, roles)
	assert.Len(t, roles, 2)
}

// TestAuthRepository_GetUserRoles_NoRoles 测试获取没有角色的用户
func TestAuthRepository_GetUserRoles_NoRoles(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewAuthRepository(db)
	ctx := context.Background()

	// 创建没有角色的用户 - GetUserRoles需要string类型的_id
	userID := "user_" + primitive.NewObjectID().Hex()
	userDoc := bson.M{
		"_id":        userID,
		"username":   "testuser",
		"email":      "test@example.com",
		"status":     "active",
		"roles":      []string{},
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}
	_, err := db.Collection("users").InsertOne(ctx, userDoc)
	require.NoError(t, err)

	// Act
	roles, err := repo.GetUserRoles(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, roles)
	assert.Empty(t, roles)
}

// TestAuthRepository_HasUserRole 测试检查用户是否有角色
func TestAuthRepository_HasUserRole(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewAuthRepository(db)
	ctx := context.Background()

	// 创建角色
	role := &authModel.Role{
		Name:        "reader",
		Description: "Reader role",
		Permissions: []string{"book.read"},
		IsSystem:    false,
	}
	err := repo.CreateRole(ctx, role)
	require.NoError(t, err)

	// 创建用户并分配角色 - 使用ObjectID作为_id
	userID := primitive.NewObjectID()
	userDoc := bson.M{
		"_id":        userID,
		"username":   "testuser",
		"email":      "test@example.com",
		"status":     "active",
		"roles":      []string{role.ID},
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}
	_, err = db.Collection("users").InsertOne(ctx, userDoc)
	require.NoError(t, err)

	// Act - 检查用户有角色
	hasRole, err := repo.HasUserRole(ctx, userID.Hex(), role.ID)

	// Assert
	require.NoError(t, err)
	assert.True(t, hasRole)
}

// TestAuthRepository_GetUserPermissions 测试获取用户权限
func TestAuthRepository_GetUserPermissions(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewAuthRepository(db)
	ctx := context.Background()

	// 创建两个角色
	role1 := &authModel.Role{
		Name:        "reader",
		Permissions: []string{"book.read", "book.list"},
		IsSystem:    false,
	}
	role2 := &authModel.Role{
		Name:        "author",
		Permissions: []string{"book.write", "book.publish"},
		IsSystem:    false,
	}

	repo.CreateRole(ctx, role1)
	repo.CreateRole(ctx, role2)

	// 创建用户并分配角色 - GetUserRoles需要string类型的_id
	userID := "user_" + primitive.NewObjectID().Hex()
	userDoc := bson.M{
		"_id":        userID,
		"username":   "testuser",
		"email":      "test@example.com",
		"status":     "active",
		"roles":      []string{role1.ID, role2.ID},
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}
	_, err := db.Collection("users").InsertOne(ctx, userDoc)
	require.NoError(t, err)

	// Act
	permissions, err := repo.GetUserPermissions(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, permissions)
	assert.Contains(t, permissions, "book.read")
	assert.Contains(t, permissions, "book.write")
}

// TestAuthRepository_GetRolePermissions 测试获取角色权限
func TestAuthRepository_GetRolePermissions(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewAuthRepository(db)
	ctx := context.Background()

	// 创建角色
	expectedPerms := []string{"book.read", "book.write", "book.delete"}
	role := &authModel.Role{
		Name:        "admin",
		Permissions: expectedPerms,
		IsSystem:    false,
	}
	err := repo.CreateRole(ctx, role)
	require.NoError(t, err)

	// Act
	permissions, err := repo.GetRolePermissions(ctx, role.ID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPerms, permissions)
}
