package user

import (
	"context"
	"testing"

	authModel "Qingyu_backend/models/auth"
	usersModel "Qingyu_backend/models/users"
	repoInterfaces "Qingyu_backend/repository/interfaces/user"
	user2 "Qingyu_backend/service/interfaces/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// =========================
// 移除角色相关测试
// =========================

// TestUserService_RemoveRole_Success 测试移除角色成功
func TestUserService_RemoveRole_Success(t *testing.T) {
	// Arrange
	service, mockUserRepo, mockAuthRepo := setupUserService()
	ctx := context.Background()

	req := &user2.RemoveRoleRequest{
		UserID: "user123",
		RoleID: "role123",
	}

	user := &usersModel.User{}
	user.ID, _ = primitive.ObjectIDFromHex(req.UserID)

	mockUserRepo.On("GetByID", ctx, req.UserID).Return(user, nil)
	mockAuthRepo.On("RemoveUserRole", ctx, req.UserID, req.RoleID).Return(nil)

	// Act
	resp, err := service.RemoveRole(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Removed)

	mockUserRepo.AssertExpectations(t)
	mockAuthRepo.AssertExpectations(t)
}

// TestUserService_RemoveRole_UserNotFound 测试移除角色-用户不存在
func TestUserService_RemoveRole_UserNotFound(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	req := &user2.RemoveRoleRequest{
		UserID: "nonexistent",
		RoleID: "role123",
	}

	mockUserRepo.On("GetByID", ctx, req.UserID).Return(nil, repoInterfaces.NewUserRepositoryError(repoInterfaces.ErrorTypeNotFound, "user not found", nil))

	// Act
	resp, err := service.RemoveRole(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户不存在")

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_RemoveRole_EmptyUserID 测试移除角色-用户ID为空
func TestUserService_RemoveRole_EmptyUserID(t *testing.T) {
	// Arrange
	service, _, _ := setupUserService()
	ctx := context.Background()

	req := &user2.RemoveRoleRequest{
		UserID: "",
		RoleID: "role123",
	}

	// Act
	resp, err := service.RemoveRole(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户ID不能为空")
}

// TestUserService_RemoveRole_EmptyRoleID 测试移除角色-角色ID为空
func TestUserService_RemoveRole_EmptyRoleID(t *testing.T) {
	// Arrange
	service, _, _ := setupUserService()
	ctx := context.Background()

	req := &user2.RemoveRoleRequest{
		UserID: "user123",
		RoleID: "",
	}

	// Act
	resp, err := service.RemoveRole(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "角色ID不能为空")
}

// TestUserService_RemoveRole_RepositoryError 测试移除角色-仓库错误
func TestUserService_RemoveRole_RepositoryError(t *testing.T) {
	// Arrange
	service, mockUserRepo, mockAuthRepo := setupUserService()
	ctx := context.Background()

	req := &user2.RemoveRoleRequest{
		UserID: "user123",
		RoleID: "role123",
	}

	user := &usersModel.User{}
	user.ID, _ = primitive.ObjectIDFromHex(req.UserID)

	mockUserRepo.On("GetByID", ctx, req.UserID).Return(user, nil)
	mockAuthRepo.On("RemoveUserRole", ctx, req.UserID, req.RoleID).Return(assert.AnError)

	// Act
	resp, err := service.RemoveRole(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "移除角色失败")

	mockUserRepo.AssertExpectations(t)
	mockAuthRepo.AssertExpectations(t)
}

// =========================
// 获取用户角色相关测试
// =========================

// TestUserService_GetUserRoles_Success 测试获取用户角色成功
func TestUserService_GetUserRoles_Success(t *testing.T) {
	// Arrange
	service, mockUserRepo, mockAuthRepo := setupUserService()
	ctx := context.Background()

	req := &user2.GetUserRolesRequest{
		UserID: "user123",
	}

	user := &usersModel.User{}
	user.ID, _ = primitive.ObjectIDFromHex(req.UserID)

	roles := []*authModel.Role{
		{ID: "role1", Name: "reader"},
		{ID: "role2", Name: "author"},
	}

	mockUserRepo.On("GetByID", ctx, req.UserID).Return(user, nil)
	mockAuthRepo.On("GetUserRoles", ctx, req.UserID).Return(roles, nil)

	// Act
	resp, err := service.GetUserRoles(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Roles, 2)
	assert.Contains(t, resp.Roles, "reader")
	assert.Contains(t, resp.Roles, "author")

	mockUserRepo.AssertExpectations(t)
	mockAuthRepo.AssertExpectations(t)
}

// TestUserService_GetUserRoles_UserNotFound 测试获取用户角色-用户不存在
func TestUserService_GetUserRoles_UserNotFound(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	req := &user2.GetUserRolesRequest{
		UserID: "nonexistent",
	}

	mockUserRepo.On("GetByID", ctx, req.UserID).Return(nil, repoInterfaces.NewUserRepositoryError(repoInterfaces.ErrorTypeNotFound, "user not found", nil))

	// Act
	resp, err := service.GetUserRoles(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户不存在")

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_GetUserRoles_EmptyUserID 测试获取用户角色-用户ID为空
func TestUserService_GetUserRoles_EmptyUserID(t *testing.T) {
	// Arrange
	service, _, _ := setupUserService()
	ctx := context.Background()

	req := &user2.GetUserRolesRequest{
		UserID: "",
	}

	// Act
	resp, err := service.GetUserRoles(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户ID不能为空")
}

// TestUserService_GetUserRoles_NoRoles 测试获取用户角色-用户没有角色
func TestUserService_GetUserRoles_NoRoles(t *testing.T) {
	// Arrange
	service, mockUserRepo, mockAuthRepo := setupUserService()
	ctx := context.Background()

	req := &user2.GetUserRolesRequest{
		UserID: "user123",
	}

	user := &usersModel.User{}
	user.ID, _ = primitive.ObjectIDFromHex(req.UserID)

	roles := []*authModel.Role{}

	mockUserRepo.On("GetByID", ctx, req.UserID).Return(user, nil)
	mockAuthRepo.On("GetUserRoles", ctx, req.UserID).Return(roles, nil)

	// Act
	resp, err := service.GetUserRoles(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Roles, 0)

	mockUserRepo.AssertExpectations(t)
	mockAuthRepo.AssertExpectations(t)
}

// TestUserService_GetUserRoles_RepositoryError 测试获取用户角色-仓库错误
func TestUserService_GetUserRoles_RepositoryError(t *testing.T) {
	// Arrange
	service, mockUserRepo, mockAuthRepo := setupUserService()
	ctx := context.Background()

	req := &user2.GetUserRolesRequest{
		UserID: "user123",
	}

	user := &usersModel.User{}
	user.ID, _ = primitive.ObjectIDFromHex(req.UserID)

	mockUserRepo.On("GetByID", ctx, req.UserID).Return(user, nil)
	mockAuthRepo.On("GetUserRoles", ctx, req.UserID).Return(nil, assert.AnError)

	// Act
	resp, err := service.GetUserRoles(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "获取用户角色失败")

	mockUserRepo.AssertExpectations(t)
	mockAuthRepo.AssertExpectations(t)
}

// =========================
// 获取用户权限相关测试
// =========================

// TestUserService_GetUserPermissions_Success 测试获取用户权限成功
func TestUserService_GetUserPermissions_Success(t *testing.T) {
	// Arrange
	service, mockUserRepo, mockAuthRepo := setupUserService()
	ctx := context.Background()

	req := &user2.GetUserPermissionsRequest{
		UserID: "user123",
	}

	user := &usersModel.User{}
	user.ID, _ = primitive.ObjectIDFromHex(req.UserID)

	permissions := []string{"read:books", "write:books", "delete:books"}

	mockUserRepo.On("GetByID", ctx, req.UserID).Return(user, nil)
	mockAuthRepo.On("GetUserPermissions", ctx, req.UserID).Return(permissions, nil)

	// Act
	resp, err := service.GetUserPermissions(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Permissions, 3)
	assert.Contains(t, resp.Permissions, "read:books")

	mockUserRepo.AssertExpectations(t)
	mockAuthRepo.AssertExpectations(t)
}

// TestUserService_GetUserPermissions_UserNotFound 测试获取用户权限-用户不存在
func TestUserService_GetUserPermissions_UserNotFound(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	req := &user2.GetUserPermissionsRequest{
		UserID: "nonexistent",
	}

	mockUserRepo.On("GetByID", ctx, req.UserID).Return(nil, repoInterfaces.NewUserRepositoryError(repoInterfaces.ErrorTypeNotFound, "user not found", nil))

	// Act
	resp, err := service.GetUserPermissions(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户不存在")

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_GetUserPermissions_EmptyUserID 测试获取用户权限-用户ID为空
func TestUserService_GetUserPermissions_EmptyUserID(t *testing.T) {
	// Arrange
	service, _, _ := setupUserService()
	ctx := context.Background()

	req := &user2.GetUserPermissionsRequest{
		UserID: "",
	}

	// Act
	resp, err := service.GetUserPermissions(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户ID不能为空")
}

// TestUserService_GetUserPermissions_NoPermissions 测试获取用户权限-用户没有权限
func TestUserService_GetUserPermissions_NoPermissions(t *testing.T) {
	// Arrange
	service, mockUserRepo, mockAuthRepo := setupUserService()
	ctx := context.Background()

	req := &user2.GetUserPermissionsRequest{
		UserID: "user123",
	}

	user := &usersModel.User{}
	user.ID, _ = primitive.ObjectIDFromHex(req.UserID)

	permissions := []string{}

	mockUserRepo.On("GetByID", ctx, req.UserID).Return(user, nil)
	mockAuthRepo.On("GetUserPermissions", ctx, req.UserID).Return(permissions, nil)

	// Act
	resp, err := service.GetUserPermissions(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Permissions, 0)

	mockUserRepo.AssertExpectations(t)
	mockAuthRepo.AssertExpectations(t)
}
