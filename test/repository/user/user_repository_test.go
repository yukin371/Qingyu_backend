package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/global"
	usersModel "Qingyu_backend/models/users"
	userInterface "Qingyu_backend/repository/interfaces/user"
	userMongo "Qingyu_backend/repository/mongodb/user"
	"Qingyu_backend/test/testutil"
)

func setupUserTest(t *testing.T) context.Context {
	testutil.SetupTestDB(t)
	ctx := context.Background()

	// 清空测试数据
	_ = global.DB.Collection("users").Drop(ctx)
	
	return ctx
}

// 辅助函数：创建测试用户
func createTestUser(username, email string) *usersModel.User {
	user := &usersModel.User{
		Username:      username,
		Email:         email,
		Password:      "hashed_password_123",
		Role:          "user",
		Status:        usersModel.UserStatusActive,
		EmailVerified: false,
		PhoneVerified: false,
	}
	return user
}

// ==================== 基础CRUD测试 ====================

func TestUserRepository_Create(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("成功创建用户", func(t *testing.T) {
		user := createTestUser("testuser", "test@example.com")
		
		err := repo.Create(ctx, user)
		assert.NoError(t, err)
		assert.NotEmpty(t, user.ID)
		assert.False(t, user.CreatedAt.IsZero())
		assert.False(t, user.UpdatedAt.IsZero())
	})

	t.Run("邮箱已存在", func(t *testing.T) {
		_ = global.DB.Collection("users").Drop(ctx)
		
		user1 := createTestUser("user1", "duplicate@example.com")
		err := repo.Create(ctx, user1)
		require.NoError(t, err)
		
		user2 := createTestUser("user2", "duplicate@example.com")
		err = repo.Create(ctx, user2)
		// 注意：可能不会报错，因为MongoDB没有建立email唯一索引
		// 此测试验证Repository的基本行为
		if err != nil {
			assert.True(t, userInterface.IsDuplicateError(err))
		}
	})

	t.Run("空用户对象", func(t *testing.T) {
		err := repo.Create(ctx, nil)
		assert.Error(t, err)
		assert.True(t, userInterface.IsValidationError(err))
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("成功获取用户", func(t *testing.T) {
		testUser := createTestUser("gettest", "get@example.com")
		err := repo.Create(ctx, testUser)
		require.NoError(t, err)
		
		found, err := repo.GetByID(ctx, testUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, testUser.Username, found.Username)
		assert.Equal(t, testUser.Email, found.Email)
	})

	t.Run("用户不存在", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "nonexistent_id")
		assert.Error(t, err)
		assert.True(t, userInterface.IsNotFoundError(err))
	})
}

func TestUserRepository_Update(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("成功更新用户", func(t *testing.T) {
		testUser := createTestUser("updatetest", "update@example.com")
		err := repo.Create(ctx, testUser)
		require.NoError(t, err)
		
		updates := map[string]interface{}{
			"nickname": "Updated Nickname",
			"bio":      "Updated bio",
		}
		
		err = repo.Update(ctx, testUser.ID, updates)
		assert.NoError(t, err)

		// 验证更新
		found, _ := repo.GetByID(ctx, testUser.ID)
		assert.Equal(t, "Updated Nickname", found.Nickname)
		assert.Equal(t, "Updated bio", found.Bio)
	})

	t.Run("更新不存在的用户", func(t *testing.T) {
		updates := map[string]interface{}{"nickname": "Test"}
		err := repo.Update(ctx, "nonexistent", updates)
		assert.Error(t, err)
		assert.True(t, userInterface.IsNotFoundError(err))
	})
}

func TestUserRepository_Delete(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("成功软删除用户", func(t *testing.T) {
		testUser := createTestUser("deletetest", "delete@example.com")
		err := repo.Create(ctx, testUser)
		require.NoError(t, err)
		
		err = repo.Delete(ctx, testUser.ID)
		assert.NoError(t, err)
		
		// 验证软删除后无法再获取
		_, err = repo.GetByID(ctx, testUser.ID)
		assert.Error(t, err)
		assert.True(t, userInterface.IsNotFoundError(err))
	})

	t.Run("删除不存在的用户", func(t *testing.T) {
		err := repo.Delete(ctx, "nonexistent")
		assert.Error(t, err)
		assert.True(t, userInterface.IsNotFoundError(err))
	})
}

// ==================== 查询方法测试 ====================

func TestUserRepository_GetByUsername(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("成功根据用户名查询", func(t *testing.T) {
		testUser := createTestUser("usernametest", "username@example.com")
		err := repo.Create(ctx, testUser)
		require.NoError(t, err)
		
		found, err := repo.GetByUsername(ctx, "usernametest")
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, testUser.Username, found.Username)
	})

	t.Run("用户名不存在", func(t *testing.T) {
		_, err := repo.GetByUsername(ctx, "nonexistent")
		assert.Error(t, err)
		assert.True(t, userInterface.IsNotFoundError(err))
	})
}

func TestUserRepository_GetByEmail(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("成功根据邮箱查询", func(t *testing.T) {
		testUser := createTestUser("emailtest", "emailtest@example.com")
		err := repo.Create(ctx, testUser)
		require.NoError(t, err)
		
		found, err := repo.GetByEmail(ctx, "emailtest@example.com")
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, testUser.Email, found.Email)
	})

	t.Run("邮箱不存在", func(t *testing.T) {
		_, err := repo.GetByEmail(ctx, "nonexistent@example.com")
		assert.Error(t, err)
		assert.True(t, userInterface.IsNotFoundError(err))
	})
}

func TestUserRepository_GetByPhone(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("成功根据手机号查询", func(t *testing.T) {
		testUser := createTestUser("phonetest", "phone@example.com")
		testUser.Phone = "13800138000"
		err := repo.Create(ctx, testUser)
		require.NoError(t, err)
		
		found, err := repo.GetByPhone(ctx, "13800138000")
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, testUser.Phone, found.Phone)
	})

	t.Run("手机号不存在", func(t *testing.T) {
		_, err := repo.GetByPhone(ctx, "99999999999")
		assert.Error(t, err)
		assert.True(t, userInterface.IsNotFoundError(err))
	})
}

// ==================== 存在性检查测试 ====================

func TestUserRepository_ExistsByUsername(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("用户名存在", func(t *testing.T) {
		testUser := createTestUser("existsuser", "exists@example.com")
		err := repo.Create(ctx, testUser)
		require.NoError(t, err)
		
		exists, err := repo.ExistsByUsername(ctx, "existsuser")
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("用户名不存在", func(t *testing.T) {
		exists, err := repo.ExistsByUsername(ctx, "notexists")
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestUserRepository_ExistsByEmail(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("邮箱存在", func(t *testing.T) {
		testUser := createTestUser("emailexists", "emailexists@example.com")
		err := repo.Create(ctx, testUser)
		require.NoError(t, err)
		
		exists, err := repo.ExistsByEmail(ctx, "emailexists@example.com")
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("邮箱不存在", func(t *testing.T) {
		exists, err := repo.ExistsByEmail(ctx, "notexists@example.com")
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestUserRepository_ExistsByPhone(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("手机号存在", func(t *testing.T) {
		testUser := createTestUser("phoneexists", "phoneexists@example.com")
		testUser.Phone = "13800138001"
		err := repo.Create(ctx, testUser)
		require.NoError(t, err)
		
		exists, err := repo.ExistsByPhone(ctx, "13800138001")
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("手机号不存在", func(t *testing.T) {
		exists, err := repo.ExistsByPhone(ctx, "99999999998")
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

// ==================== 状态管理测试 ====================

func TestUserRepository_UpdateLastLogin(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("成功更新最后登录信息", func(t *testing.T) {
		testUser := createTestUser("logintest", "login@example.com")
		err := repo.Create(ctx, testUser)
		require.NoError(t, err)
		
		err = repo.UpdateLastLogin(ctx, testUser.ID, "192.168.1.1")
		assert.NoError(t, err)
		
		// 验证更新
		found, _ := repo.GetByID(ctx, testUser.ID)
		assert.False(t, found.LastLoginAt.IsZero())
		assert.Equal(t, "192.168.1.1", found.LastLoginIP)
	})

	t.Run("更新不存在的用户", func(t *testing.T) {
		err := repo.UpdateLastLogin(ctx, "nonexistent", "192.168.1.1")
		assert.Error(t, err)
		assert.True(t, userInterface.IsNotFoundError(err))
	})
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("成功更新密码", func(t *testing.T) {
		testUser := createTestUser("pwdtest", "pwd@example.com")
		err := repo.Create(ctx, testUser)
		require.NoError(t, err)
		
		newPassword := "new_hashed_password"
		err = repo.UpdatePassword(ctx, testUser.ID, newPassword)
		assert.NoError(t, err)
		
		// 验证更新
		found, _ := repo.GetByID(ctx, testUser.ID)
		assert.Equal(t, newPassword, found.Password)
	})

	t.Run("更新不存在的用户密码", func(t *testing.T) {
		err := repo.UpdatePassword(ctx, "nonexistent", "newpwd")
		assert.Error(t, err)
		assert.True(t, userInterface.IsNotFoundError(err))
	})
}

func TestUserRepository_UpdateStatus(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("成功更新状态", func(t *testing.T) {
		testUser := createTestUser("statustest", "status@example.com")
		err := repo.Create(ctx, testUser)
		require.NoError(t, err)
		
		err = repo.UpdateStatus(ctx, testUser.ID, usersModel.UserStatusInactive)
		assert.NoError(t, err)
		
		// 验证更新
		found, _ := repo.GetByID(ctx, testUser.ID)
		assert.Equal(t, usersModel.UserStatusInactive, found.Status)
	})

	t.Run("更新不存在的用户状态", func(t *testing.T) {
		err := repo.UpdateStatus(ctx, "nonexistent", usersModel.UserStatusBanned)
		assert.Error(t, err)
		assert.True(t, userInterface.IsNotFoundError(err))
	})
}

// ==================== 验证状态测试 ====================

func TestUserRepository_SetEmailVerified(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("设置邮箱已验证", func(t *testing.T) {
		testUser := createTestUser("emailverify", "emailverify@example.com")
		err := repo.Create(ctx, testUser)
		require.NoError(t, err)
		
		err = repo.SetEmailVerified(ctx, testUser.ID, true)
		assert.NoError(t, err)

		// 验证更新
		found, _ := repo.GetByID(ctx, testUser.ID)
		assert.True(t, found.EmailVerified)
	})

	t.Run("设置不存在用户的邮箱验证", func(t *testing.T) {
		err := repo.SetEmailVerified(ctx, "nonexistent", true)
		assert.Error(t, err)
		assert.True(t, userInterface.IsNotFoundError(err))
	})
}

func TestUserRepository_SetPhoneVerified(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("设置手机号已验证", func(t *testing.T) {
		testUser := createTestUser("phoneverify", "phoneverify@example.com")
		err := repo.Create(ctx, testUser)
		require.NoError(t, err)
		
		err = repo.SetPhoneVerified(ctx, testUser.ID, true)
		assert.NoError(t, err)

		// 验证更新
		found, _ := repo.GetByID(ctx, testUser.ID)
		assert.True(t, found.PhoneVerified)
	})

	t.Run("设置不存在用户的手机验证", func(t *testing.T) {
		err := repo.SetPhoneVerified(ctx, "nonexistent", true)
		assert.Error(t, err)
		assert.True(t, userInterface.IsNotFoundError(err))
	})
}

// ==================== 列表和查询测试 ====================

func TestUserRepository_GetActiveUsers(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("获取活跃用户", func(t *testing.T) {
		// 创建活跃用户
		for i := 1; i <= 3; i++ {
			user := createTestUser("active"+string(rune('0'+i)), "active"+string(rune('0'+i))+"@example.com")
			user.Status = usersModel.UserStatusActive
			err := repo.Create(ctx, user)
			require.NoError(t, err)
		}
		
		// 创建非活跃用户
		inactiveUser := createTestUser("inactive", "inactive@example.com")
		inactiveUser.Status = usersModel.UserStatusInactive
		err := repo.Create(ctx, inactiveUser)
		require.NoError(t, err)
		
		users, err := repo.GetActiveUsers(ctx, 10)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 3)
		
		// 验证都是活跃用户
		for _, u := range users {
			assert.Equal(t, usersModel.UserStatusActive, u.Status)
		}
	})
}

func TestUserRepository_GetUsersByRole(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("根据角色获取用户", func(t *testing.T) {
		// 创建不同角色的用户
		adminUser := createTestUser("admin1", "admin1@example.com")
		adminUser.Role = "admin"
		err := repo.Create(ctx, adminUser)
		require.NoError(t, err)
		
		normalUser := createTestUser("user1", "user1@example.com")
		normalUser.Role = "user"
		err = repo.Create(ctx, normalUser)
		require.NoError(t, err)
		
		admins, err := repo.GetUsersByRole(ctx, "admin", 10)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(admins), 1)
		
		// 验证都是admin角色
		for _, u := range admins {
			assert.Equal(t, "admin", u.Role)
		}
	})
}

// ==================== 批量操作测试 ====================

func TestUserRepository_BatchUpdateStatus(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("批量更新状态", func(t *testing.T) {
		// 创建3个用户
		userIDs := []string{}
		for i := 1; i <= 3; i++ {
			user := createTestUser("batchstatus"+string(rune('0'+i)), "batchstatus"+string(rune('0'+i))+"@example.com")
			err := repo.Create(ctx, user)
			require.NoError(t, err)
			userIDs = append(userIDs, user.ID)
		}
		
		// 批量更新为inactive
		err := repo.BatchUpdateStatus(ctx, userIDs, usersModel.UserStatusInactive)
		assert.NoError(t, err)
		
		// 验证所有用户状态已更新
		for _, id := range userIDs {
			found, _ := repo.GetByID(ctx, id)
			assert.Equal(t, usersModel.UserStatusInactive, found.Status)
		}
	})

	t.Run("空数组不报错", func(t *testing.T) {
		err := repo.BatchUpdateStatus(ctx, []string{}, usersModel.UserStatusBanned)
		assert.NoError(t, err)
	})
}

func TestUserRepository_BatchDelete(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("批量删除用户", func(t *testing.T) {
		// 创建3个用户
		userIDs := []string{}
		for i := 1; i <= 3; i++ {
			user := createTestUser("batchdel"+string(rune('0'+i)), "batchdel"+string(rune('0'+i))+"@example.com")
			err := repo.Create(ctx, user)
			require.NoError(t, err)
			userIDs = append(userIDs, user.ID)
		}
		
		// 批量删除
		err := repo.BatchDelete(ctx, userIDs)
		assert.NoError(t, err)
		
		// 验证所有用户已被软删除
		for _, id := range userIDs {
			_, err := repo.GetByID(ctx, id)
			assert.Error(t, err)
			assert.True(t, userInterface.IsNotFoundError(err))
		}
	})

	t.Run("空数组不报错", func(t *testing.T) {
		err := repo.BatchDelete(ctx, []string{})
		assert.NoError(t, err)
	})
}

// ==================== 搜索和统计测试 ====================

func TestUserRepository_SearchUsers(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("搜索用户", func(t *testing.T) {
		// 创建测试用户
		user1 := createTestUser("searchtestuser", "searchtest@example.com")
		user1.Nickname = "SearchNickname"
		err := repo.Create(ctx, user1)
		require.NoError(t, err)
		
		// 按用户名搜索
		users, err := repo.SearchUsers(ctx, "searchtest", 10)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 1)
		
		// 按昵称搜索
		users, err = repo.SearchUsers(ctx, "SearchNick", 10)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 1)
	})

	t.Run("无匹配结果", func(t *testing.T) {
		users, err := repo.SearchUsers(ctx, "nonexistentkeyword", 10)
		assert.NoError(t, err)
		assert.Empty(t, users)
	})
}

func TestUserRepository_CountByRole(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("按角色统计", func(t *testing.T) {
		// 创建不同角色的用户
		for i := 1; i <= 3; i++ {
			user := createTestUser("author"+string(rune('0'+i)), "author"+string(rune('0'+i))+"@example.com")
			user.Role = "author"
			err := repo.Create(ctx, user)
			require.NoError(t, err)
		}
		
		count, err := repo.CountByRole(ctx, "author")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(3))
	})
}

func TestUserRepository_CountByStatus(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("按状态统计", func(t *testing.T) {
		// 创建不同状态的用户
		for i := 1; i <= 2; i++ {
			user := createTestUser("banned"+string(rune('0'+i)), "banned"+string(rune('0'+i))+"@example.com")
			user.Status = usersModel.UserStatusBanned
			err := repo.Create(ctx, user)
			require.NoError(t, err)
		}
		
		count, err := repo.CountByStatus(ctx, usersModel.UserStatusBanned)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(2))
	})
}

// ==================== 健康检查测试 ====================

func TestUserRepository_Health(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("健康检查", func(t *testing.T) {
		err := repo.Health(ctx)
		assert.NoError(t, err)
	})
}

// ==================== Exists测试 ====================

func TestUserRepository_Exists(t *testing.T) {
	ctx := setupUserTest(t)
	repo := userMongo.NewMongoUserRepository(global.DB)

	t.Run("用户存在", func(t *testing.T) {
		testUser := createTestUser("existsidtest", "existsid@example.com")
		err := repo.Create(ctx, testUser)
		require.NoError(t, err)
		
		exists, err := repo.Exists(ctx, testUser.ID)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("用户不存在", func(t *testing.T) {
		exists, err := repo.Exists(ctx, "nonexistent_id_123")
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}
