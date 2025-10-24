package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
	"Qingyu_backend/global"
	authModel "Qingyu_backend/models/shared/auth"
	"Qingyu_backend/repository/mongodb/user"
)

// TestRoleRepository_Integration 角色Repository集成测试
func TestRoleRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试（使用 -short 标志）")
	}

	// 加载配置
	cfg, err := config.LoadConfig("../../../config/config.yaml")
	require.NoError(t, err, "加载配置失败")

	// 初始化全局配置
	config.GlobalConfig = cfg

	// 初始化数据库
	err = core.InitDB()
	require.NoError(t, err, "初始化数据库失败")

	// 创建Repository
	roleRepo := user.NewMongoRoleRepository(global.DB)
	ctx := context.Background()

	// 健康检查
	t.Run("Health", func(t *testing.T) {
		err := roleRepo.Health(ctx)
		assert.NoError(t, err, "健康检查应该成功")
	})

	// 创建测试角色
	testRole := &authModel.Role{
		Name:        "test_role_" + time.Now().Format("20060102150405"),
		Description: "测试角色",
		IsSystem:    false,
		Permissions: []string{authModel.PermUserRead, authModel.PermUserWrite},
	}

	// 测试创建角色
	t.Run("Create", func(t *testing.T) {
		err := roleRepo.Create(ctx, testRole)
		assert.NoError(t, err, "创建角色应该成功")
		assert.NotEmpty(t, testRole.ID, "角色ID应该被设置")
	})

	// 测试根据ID获取角色
	t.Run("GetByID", func(t *testing.T) {
		role, err := roleRepo.GetByID(ctx, testRole.ID)
		assert.NoError(t, err, "获取角色应该成功")
		assert.Equal(t, testRole.Name, role.Name, "角色名应该匹配")
		assert.Equal(t, testRole.Description, role.Description, "描述应该匹配")
	})

	// 测试根据名称获取角色
	t.Run("GetByName", func(t *testing.T) {
		role, err := roleRepo.GetByName(ctx, testRole.Name)
		assert.NoError(t, err, "根据名称获取角色应该成功")
		assert.Equal(t, testRole.ID, role.ID, "角色ID应该匹配")
	})

	// 测试角色名是否存在
	t.Run("ExistsByName", func(t *testing.T) {
		exists, err := roleRepo.ExistsByName(ctx, testRole.Name)
		assert.NoError(t, err, "检查角色名存在应该成功")
		assert.True(t, exists, "角色名应该存在")

		exists, err = roleRepo.ExistsByName(ctx, "nonexistent_role")
		assert.NoError(t, err, "检查不存在的角色名应该成功")
		assert.False(t, exists, "不存在的角色名应该返回false")
	})

	// 测试更新角色
	t.Run("Update", func(t *testing.T) {
		updates := map[string]interface{}{
			"description": "更新后的描述",
			"is_system":   true,
		}
		err := roleRepo.Update(ctx, testRole.ID, updates)
		assert.NoError(t, err, "更新角色应该成功")

		// 验证更新
		role, err := roleRepo.GetByID(ctx, testRole.ID)
		assert.NoError(t, err, "获取更新后的角色应该成功")
		assert.Equal(t, "更新后的描述", role.Description, "描述应该已更新")
		assert.True(t, role.IsSystem, "IsSystem应该已更新")
	})

	// 测试获取角色权限
	t.Run("GetRolePermissions", func(t *testing.T) {
		permissions, err := roleRepo.GetRolePermissions(ctx, testRole.ID)
		assert.NoError(t, err, "获取角色权限应该成功")
		assert.Equal(t, testRole.Permissions, permissions, "权限列表应该匹配")
	})

	// 测试添加权限
	t.Run("AddPermission", func(t *testing.T) {
		err := roleRepo.AddPermission(ctx, testRole.ID, authModel.PermDocumentRead)
		assert.NoError(t, err, "添加权限应该成功")

		permissions, err := roleRepo.GetRolePermissions(ctx, testRole.ID)
		assert.NoError(t, err, "获取权限应该成功")
		assert.Contains(t, permissions, authModel.PermDocumentRead, "权限列表应该包含新权限")
	})

	// 测试移除权限
	t.Run("RemovePermission", func(t *testing.T) {
		err := roleRepo.RemovePermission(ctx, testRole.ID, authModel.PermDocumentRead)
		assert.NoError(t, err, "移除权限应该成功")

		permissions, err := roleRepo.GetRolePermissions(ctx, testRole.ID)
		assert.NoError(t, err, "获取权限应该成功")
		assert.NotContains(t, permissions, authModel.PermDocumentRead, "权限列表不应该包含已移除的权限")
	})

	// 测试更新角色权限
	t.Run("UpdateRolePermissions", func(t *testing.T) {
		newPermissions := []string{
			authModel.PermUserRead,
			authModel.PermDocumentWrite,
			authModel.PermAdminAccess,
		}
		err := roleRepo.UpdateRolePermissions(ctx, testRole.ID, newPermissions)
		assert.NoError(t, err, "更新角色权限应该成功")

		permissions, err := roleRepo.GetRolePermissions(ctx, testRole.ID)
		assert.NoError(t, err, "获取权限应该成功")
		assert.ElementsMatch(t, newPermissions, permissions, "权限列表应该匹配")
	})

	// 测试列出所有角色
	t.Run("ListAllRoles", func(t *testing.T) {
		roles, err := roleRepo.ListAllRoles(ctx)
		assert.NoError(t, err, "列出所有角色应该成功")
		assert.NotEmpty(t, roles, "角色列表不应为空")

		// 验证测试角色在列表中
		found := false
		for _, role := range roles {
			if role.ID == testRole.ID {
				found = true
				break
			}
		}
		assert.True(t, found, "应该找到测试角色")
	})

	// 测试列出默认角色
	t.Run("ListDefaultRoles", func(t *testing.T) {
		roles, err := roleRepo.ListDefaultRoles(ctx)
		assert.NoError(t, err, "列出默认角色应该成功")

		// 验证所有角色都是系统角色
		for _, role := range roles {
			assert.True(t, role.IsSystem, "所有角色都应该是系统角色")
		}
	})

	// 测试获取默认角色
	t.Run("GetDefaultRole", func(t *testing.T) {
		// 首先创建一个默认角色
		defaultRole := &authModel.Role{
			Name:        "default_user_" + time.Now().Format("20060102150405"),
			Description: "默认用户角色",
			IsSystem:    true,
			IsDefault:   true,
			Permissions: []string{authModel.PermUserRead},
		}
		err := roleRepo.Create(ctx, defaultRole)
		require.NoError(t, err, "创建默认角色应该成功")

		// 测试获取默认角色
		role, err := roleRepo.GetDefaultRole(ctx)
		assert.NoError(t, err, "获取默认角色应该成功")
		if role != nil {
			assert.True(t, role.IsSystem, "应该是系统角色")
			assert.True(t, role.IsDefault, "应该是默认角色")
		}

		// 清理默认角色
		if defaultRole.ID != "" {
			roleRepo.Delete(ctx, defaultRole.ID)
		}
	})

	// 测试统计
	t.Run("CountByName", func(t *testing.T) {
		count, err := roleRepo.CountByName(ctx, testRole.Name)
		assert.NoError(t, err, "按名称统计应该成功")
		assert.Equal(t, int64(1), count, "应该只有一个匹配的角色")
	})

	// 清理测试数据
	t.Run("Delete", func(t *testing.T) {
		err := roleRepo.Delete(ctx, testRole.ID)
		assert.NoError(t, err, "删除角色应该成功")

		// 验证删除
		_, err = roleRepo.GetByID(ctx, testRole.ID)
		assert.Error(t, err, "获取已删除的角色应该失败")
	})
}

// TestRoleRepository_DefaultRole 默认角色测试
func TestRoleRepository_DefaultRole(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试（使用 -short 标志）")
	}

	// 加载配置
	cfg, err := config.LoadConfig("../../../config/config.yaml")
	require.NoError(t, err, "加载配置失败")

	// 初始化全局配置
	config.GlobalConfig = cfg

	// 初始化数据库
	err = core.InitDB()
	require.NoError(t, err, "初始化数据库失败")

	// 创建Repository
	roleRepo := user.NewMongoRoleRepository(global.DB)
	ctx := context.Background()

	// 创建多个默认角色
	timestamp := time.Now().Format("20060102150405")
	defaultRoles := []*authModel.Role{
		{
			Name:        "default_role1_" + timestamp,
			Description: "默认角色1",
			IsSystem:    true,
			IsDefault:   true, // 标记为默认角色
			Permissions: []string{authModel.PermUserRead, authModel.PermBookRead},
		},
		{
			Name:        "default_role2_" + timestamp,
			Description: "默认角色2",
			IsSystem:    true,
			IsDefault:   false, // 非默认角色
			Permissions: []string{authModel.PermUserRead, authModel.PermBookRead},
		},
	}

	// 创建角色
	for _, role := range defaultRoles {
		err := roleRepo.Create(ctx, role)
		require.NoError(t, err, "创建默认角色应该成功")
	}

	// 测试获取默认角色
	t.Run("GetDefaultRole", func(t *testing.T) {
		role, err := roleRepo.GetDefaultRole(ctx)
		assert.NoError(t, err, "获取默认角色应该成功")
		assert.NotNil(t, role, "应该返回一个默认角色")
		if role != nil {
			assert.True(t, role.IsSystem, "应该是系统角色")
			assert.True(t, role.IsDefault, "应该是默认角色")
		}
	})

	// 测试列出所有默认角色
	t.Run("ListDefaultRoles", func(t *testing.T) {
		roles, err := roleRepo.ListDefaultRoles(ctx)
		assert.NoError(t, err, "列出默认角色应该成功")
		assert.GreaterOrEqual(t, len(roles), 2, "应该至少有两个默认角色")
	})

	// 清理测试数据
	for _, role := range defaultRoles {
		_ = roleRepo.Delete(ctx, role.ID)
	}
}
