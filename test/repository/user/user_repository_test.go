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
	usersModel "Qingyu_backend/models/users"
	"Qingyu_backend/repository/mongodb/user"
)

// TestUserRepository_Integration 用户Repository集成测试
func TestUserRepository_Integration(t *testing.T) {
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
	userRepo := user.NewMongoUserRepository(global.DB)
	ctx := context.Background()

	// 健康检查
	t.Run("Health", func(t *testing.T) {
		err := userRepo.Health(ctx)
		assert.NoError(t, err, "健康检查应该成功")
	})

	// 创建测试用户
	testUser := &usersModel.User{
		Username: "testuser_" + time.Now().Format("20060102150405"),
		Email:    "test_" + time.Now().Format("20060102150405") + "@example.com",
		Phone:    "13800138000",
		Password: "hashed_password_123",
		Role:     usersModel.RoleUser,
		Status:   usersModel.UserStatusActive,
		Nickname: "测试用户",
		Bio:      "这是一个测试用户",
	}

	// 测试创建用户
	t.Run("Create", func(t *testing.T) {
		err := userRepo.Create(ctx, testUser)
		assert.NoError(t, err, "创建用户应该成功")
		assert.NotEmpty(t, testUser.ID, "用户ID应该被设置")
	})

	// 测试根据ID获取用户
	t.Run("GetByID", func(t *testing.T) {
		user, err := userRepo.GetByID(ctx, testUser.ID)
		assert.NoError(t, err, "获取用户应该成功")
		assert.Equal(t, testUser.Username, user.Username, "用户名应该匹配")
		assert.Equal(t, testUser.Email, user.Email, "邮箱应该匹配")
	})

	// 测试根据Email获取用户
	t.Run("GetByEmail", func(t *testing.T) {
		user, err := userRepo.GetByEmail(ctx, testUser.Email)
		assert.NoError(t, err, "根据邮箱获取用户应该成功")
		assert.Equal(t, testUser.ID, user.ID, "用户ID应该匹配")
	})

	// 测试根据Phone获取用户
	t.Run("GetByPhone", func(t *testing.T) {
		user, err := userRepo.GetByPhone(ctx, testUser.Phone)
		assert.NoError(t, err, "根据手机号获取用户应该成功")
		assert.Equal(t, testUser.ID, user.ID, "用户ID应该匹配")
	})

	// 测试邮箱是否存在
	t.Run("ExistsByEmail", func(t *testing.T) {
		exists, err := userRepo.ExistsByEmail(ctx, testUser.Email)
		assert.NoError(t, err, "检查邮箱存在应该成功")
		assert.True(t, exists, "邮箱应该存在")

		exists, err = userRepo.ExistsByEmail(ctx, "nonexistent@example.com")
		assert.NoError(t, err, "检查不存在的邮箱应该成功")
		assert.False(t, exists, "不存在的邮箱应该返回false")
	})

	// 测试更新用户
	t.Run("Update", func(t *testing.T) {
		updates := map[string]interface{}{
			"nickname": "更新后的昵称",
			"bio":      "更新后的个人简介",
		}
		err := userRepo.Update(ctx, testUser.ID, updates)
		assert.NoError(t, err, "更新用户应该成功")

		// 验证更新
		user, err := userRepo.GetByID(ctx, testUser.ID)
		assert.NoError(t, err, "获取更新后的用户应该成功")
		assert.Equal(t, "更新后的昵称", user.Nickname, "昵称应该已更新")
		assert.Equal(t, "更新后的个人简介", user.Bio, "简介应该已更新")
	})

	// 测试更新最后登录时间
	t.Run("UpdateLastLogin", func(t *testing.T) {
		err := userRepo.UpdateLastLogin(ctx, testUser.ID, "192.168.1.1")
		assert.NoError(t, err, "更新最后登录时间应该成功")

		user, err := userRepo.GetByID(ctx, testUser.ID)
		assert.NoError(t, err, "获取用户应该成功")
		assert.NotNil(t, user.LastLoginAt, "最后登录时间应该被设置")
		assert.Equal(t, "192.168.1.1", user.LastLoginIP, "登录IP应该匹配")
	})

	// 测试更新用户状态
	t.Run("UpdateStatus", func(t *testing.T) {
		err := userRepo.UpdateStatus(ctx, testUser.ID, usersModel.UserStatusInactive)
		assert.NoError(t, err, "更新状态应该成功")

		user, err := userRepo.GetByID(ctx, testUser.ID)
		assert.NoError(t, err, "获取用户应该成功")
		assert.Equal(t, usersModel.UserStatusInactive, user.Status, "状态应该已更新")
	})

	// 测试设置邮箱验证状态
	t.Run("SetEmailVerified", func(t *testing.T) {
		err := userRepo.SetEmailVerified(ctx, testUser.ID, true)
		assert.NoError(t, err, "设置邮箱验证状态应该成功")

		user, err := userRepo.GetByID(ctx, testUser.ID)
		assert.NoError(t, err, "获取用户应该成功")
		assert.True(t, user.EmailVerified, "邮箱应该已验证")
	})

	// 测试高级查询
	t.Run("FindWithFilter", func(t *testing.T) {
		filter := &usersModel.UserFilter{
			Role:     usersModel.RoleUser,
			Status:   usersModel.UserStatusInactive,
			Page:     1,
			PageSize: 10,
		}

		users, total, err := userRepo.FindWithFilter(ctx, filter)
		assert.NoError(t, err, "高级查询应该成功")
		assert.GreaterOrEqual(t, total, int64(1), "应该至少有一个用户")
		assert.NotEmpty(t, users, "用户列表不应为空")
	})

	// 测试搜索用户
	t.Run("SearchUsers", func(t *testing.T) {
		users, err := userRepo.SearchUsers(ctx, "测试", 10)
		assert.NoError(t, err, "搜索用户应该成功")
		assert.NotEmpty(t, users, "搜索结果不应为空")
	})

	// 测试统计方法
	t.Run("CountByRole", func(t *testing.T) {
		count, err := userRepo.CountByRole(ctx, usersModel.RoleUser)
		assert.NoError(t, err, "按角色统计应该成功")
		assert.GreaterOrEqual(t, count, int64(1), "应该至少有一个用户")
	})

	// 清理测试数据
	t.Run("Delete", func(t *testing.T) {
		err := userRepo.Delete(ctx, testUser.ID)
		assert.NoError(t, err, "删除用户应该成功")

		// 验证删除
		_, err = userRepo.GetByID(ctx, testUser.ID)
		assert.Error(t, err, "获取已删除的用户应该失败")
	})
}

// TestUserRepository_BatchOperations 批量操作测试
func TestUserRepository_BatchOperations(t *testing.T) {
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
	userRepo := user.NewMongoUserRepository(global.DB)
	ctx := context.Background()

	// 创建测试用户
	timestamp := time.Now().Format("20060102150405")
	testUsers := []*usersModel.User{
		{
			Username: "batch_user1_" + timestamp,
			Email:    "batch1_" + timestamp + "@example.com",
			Password: "hashed_password_123",
			Role:     usersModel.RoleUser,
			Status:   usersModel.UserStatusActive,
		},
		{
			Username: "batch_user2_" + timestamp,
			Email:    "batch2_" + timestamp + "@example.com",
			Password: "hashed_password_123",
			Role:     usersModel.RoleUser,
			Status:   usersModel.UserStatusActive,
		},
	}

	// 测试创建多个用户（使用单独创建）
	t.Run("CreateMultiple", func(t *testing.T) {
		for _, u := range testUsers {
			err := userRepo.Create(ctx, u)
			assert.NoError(t, err, "创建用户应该成功")
			assert.NotEmpty(t, u.ID, "用户ID应该被设置")
		}
	})

	// 测试批量更新状态
	t.Run("BatchUpdateStatus", func(t *testing.T) {
		ids := []string{testUsers[0].ID, testUsers[1].ID}
		err := userRepo.BatchUpdateStatus(ctx, ids, usersModel.UserStatusInactive)
		assert.NoError(t, err, "批量更新状态应该成功")

		// 验证更新
		for _, id := range ids {
			user, err := userRepo.GetByID(ctx, id)
			assert.NoError(t, err, "获取用户应该成功")
			assert.Equal(t, usersModel.UserStatusInactive, user.Status, "状态应该已更新")
		}
	})

	// 测试批量删除
	t.Run("BatchDelete", func(t *testing.T) {
		ids := []string{testUsers[0].ID, testUsers[1].ID}
		err := userRepo.BatchDelete(ctx, ids)
		assert.NoError(t, err, "批量删除应该成功")

		// 验证删除（软删除，状态应该变为deleted）
		for _, id := range ids {
			user, err := userRepo.GetByID(ctx, id)
			if err == nil {
				assert.Equal(t, usersModel.UserStatusDeleted, user.Status, "状态应该为deleted")
			}
		}
	})
}
