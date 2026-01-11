package user_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	usersModel "Qingyu_backend/models/users"
	userInterface "Qingyu_backend/repository/interfaces/user"
	userMongo "Qingyu_backend/repository/mongodb/user"
	"Qingyu_backend/test/testutil"
)

// TestUserRepository_Create 测试创建用户
func TestUserRepository_Create(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userMongo.NewMongoUserRepository(db)
	ctx := context.Background()

	user := &usersModel.User{
		Username:      "testuser",
		Email:         "test@example.com",
		Password:      "hashed_password_123",
		Roles:         []string{"user"},
		Status:        usersModel.UserStatusActive,
		EmailVerified: false,
		PhoneVerified: false,
	}

	// Act
	err := repo.Create(ctx, user)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)
}

// TestUserRepository_GetByID 测试根据ID获取用户
func TestUserRepository_GetByID(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userMongo.NewMongoUserRepository(db)
	ctx := context.Background()

	testUser := &usersModel.User{
		Username: "gettest",
		Email:    "get@example.com",
		Password: "hashed_password_123",
		Status:   usersModel.UserStatusActive,
	}
	err := repo.Create(ctx, testUser)
	require.NoError(t, err)

	// Act
	found, err := repo.GetByID(ctx, testUser.ID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, testUser.Username, found.Username)
	assert.Equal(t, testUser.Email, found.Email)
}

// TestUserRepository_GetByID_NotFound 测试获取不存在的用户
func TestUserRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userMongo.NewMongoUserRepository(db)
	ctx := context.Background()

	// Act
	found, err := repo.GetByID(ctx, "nonexistent_id")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, found)
	assert.True(t, userInterface.IsNotFoundError(err))
}

// TestUserRepository_GetByEmail 测试根据邮箱获取用户
func TestUserRepository_GetByEmail(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userMongo.NewMongoUserRepository(db)
	ctx := context.Background()

	testUser := &usersModel.User{
		Username: "emailtest",
		Email:    "emailtest@example.com",
		Password: "hashed_password_123",
		Status:   usersModel.UserStatusActive,
	}
	err := repo.Create(ctx, testUser)
	require.NoError(t, err)

	// Act
	found, err := repo.GetByEmail(ctx, "emailtest@example.com")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, testUser.Username, found.Username)
	assert.Equal(t, testUser.Email, found.Email)
}

// TestUserRepository_GetByEmail_NotFound 测试根据邮箱获取不存在的用户
func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userMongo.NewMongoUserRepository(db)
	ctx := context.Background()

	// Act
	found, err := repo.GetByEmail(ctx, "nonexistent@example.com")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, found)
}

// TestUserRepository_GetByUsername 测试根据用户名获取用户
func TestUserRepository_GetByUsername(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userMongo.NewMongoUserRepository(db)
	ctx := context.Background()

	testUser := &usersModel.User{
		Username: "usernametest",
		Email:    "usernametest@example.com",
		Password: "hashed_password_123",
		Status:   usersModel.UserStatusActive,
	}
	err := repo.Create(ctx, testUser)
	require.NoError(t, err)

	// Act
	found, err := repo.GetByUsername(ctx, "usernametest")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, "usernametest", found.Username)
}

// TestUserRepository_Update 测试更新用户
func TestUserRepository_Update(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userMongo.NewMongoUserRepository(db)
	ctx := context.Background()

	testUser := &usersModel.User{
		Username: "updatetest",
		Email:    "update@example.com",
		Password: "hashed_password_123",
		Status:   usersModel.UserStatusActive,
	}
	err := repo.Create(ctx, testUser)
	require.NoError(t, err)

	// Act - 更新用户
	updates := map[string]interface{}{
		"username": "updated_username",
		"status":   usersModel.UserStatusInactive,
	}
	err = repo.Update(ctx, testUser.ID, updates)

	// Assert
	require.NoError(t, err)

	// 验证更新
	found, err := repo.GetByID(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, "updated_username", found.Username)
	assert.Equal(t, usersModel.UserStatusInactive, found.Status)
}

// TestUserRepository_Update_NotFound 测试更新不存在的用户
func TestUserRepository_Update_NotFound(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userMongo.NewMongoUserRepository(db)
	ctx := context.Background()

	updates := map[string]interface{}{
		"username": "updated",
	}

	// Act
	err := repo.Update(ctx, "nonexistent_id", updates)

	// Assert
	assert.Error(t, err)
	assert.True(t, userInterface.IsNotFoundError(err))
}

// TestUserRepository_Delete 测试删除用户
func TestUserRepository_Delete(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userMongo.NewMongoUserRepository(db)
	ctx := context.Background()

	testUser := &usersModel.User{
		Username: "deletetest",
		Email:    "delete@example.com",
		Password: "hashed_password_123",
		Status:   usersModel.UserStatusActive,
	}
	err := repo.Create(ctx, testUser)
	require.NoError(t, err)

	// Act - 删除用户
	err = repo.Delete(ctx, testUser.ID)

	// Assert
	require.NoError(t, err)

	// 验证已删除
	found, err := repo.GetByID(ctx, testUser.ID)
	assert.Error(t, err)
	assert.Nil(t, found)
}

// TestUserRepository_Delete_NotFound 测试删除不存在的用户
func TestUserRepository_Delete_NotFound(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userMongo.NewMongoUserRepository(db)
	ctx := context.Background()

	// Act
	err := repo.Delete(ctx, "nonexistent_id")

	// Assert
	assert.Error(t, err)
}

// TestUserRepository_List 测试列出用户
func TestUserRepository_List(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userMongo.NewMongoUserRepository(db)
	ctx := context.Background()

	// 创建多个测试用户
	users := []*usersModel.User{
		{Username: "user1", Email: "user1@example.com", Password: "hash", Status: usersModel.UserStatusActive},
		{Username: "user2", Email: "user2@example.com", Password: "hash", Status: usersModel.UserStatusActive},
		{Username: "user3", Email: "user3@example.com", Password: "hash", Status: usersModel.UserStatusInactive},
	}

	for _, user := range users {
		err := repo.Create(ctx, user)
		require.NoError(t, err)
	}

	// Act - 使用基础filter
	filter := &userInterface.UserFilter{
		Limit: 10,
	}
	result, err := repo.List(ctx, filter)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result), 3)
}

// TestUserRepository_Count 测试统计用户数量
func TestUserRepository_Count(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userMongo.NewMongoUserRepository(db)
	ctx := context.Background()

	// 创建测试用户
	for i := 0; i < 5; i++ {
		user := &usersModel.User{
			Username: fmt.Sprintf("countuser_%d", i),
			Email:    fmt.Sprintf("count%d@example.com", i),
			Password: "hash",
			Status:   usersModel.UserStatusActive,
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)
	}

	// Act
	filter := &userInterface.UserFilter{}
	count, err := repo.Count(ctx, filter)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(5))
}

// TestUserRepository_UpdateStatus 测试更新用户状态
func TestUserRepository_UpdateStatus(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userMongo.NewMongoUserRepository(db)
	ctx := context.Background()

	testUser := &usersModel.User{
		Username: "statustest",
		Email:    "status@example.com",
		Password: "hash",
		Status:   usersModel.UserStatusActive,
	}
	err := repo.Create(ctx, testUser)
	require.NoError(t, err)

	// Act - 更新状态
	err = repo.UpdateStatus(ctx, testUser.ID, usersModel.UserStatusBanned)

	// Assert
	require.NoError(t, err)

	// 验证状态更新
	found, err := repo.GetByID(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, usersModel.UserStatusBanned, found.Status)
}

// TestUserRepository_UpdatePassword 测试更新密码
func TestUserRepository_UpdatePassword(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userMongo.NewMongoUserRepository(db)
	ctx := context.Background()

	testUser := &usersModel.User{
		Username: "passwordtest",
		Email:    "password@example.com",
		Password: "old_hashed_password",
		Status:   usersModel.UserStatusActive,
	}
	err := repo.Create(ctx, testUser)
	require.NoError(t, err)

	// Act - 更新密码
	newPassword := "new_hashed_password"
	err = repo.UpdatePassword(ctx, testUser.ID, newPassword)

	// Assert
	require.NoError(t, err)

	// 验证密码已更新
	found, err := repo.GetByID(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, newPassword, found.Password)
}

// TestUserRepository_UpdateLastLogin 测试更新最后登录时间
func TestUserRepository_UpdateLastLogin(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userMongo.NewMongoUserRepository(db)
	ctx := context.Background()

	testUser := &usersModel.User{
		Username: "logintest",
		Email:    "login@example.com",
		Password: "hash",
		Status:   usersModel.UserStatusActive,
	}
	err := repo.Create(ctx, testUser)
	require.NoError(t, err)

	// Act - 更新最后登录IP
	err = repo.UpdateLastLogin(ctx, testUser.ID, "192.168.1.1")

	// Assert
	require.NoError(t, err)
}

// TestUserRepository_SearchUsers 测试搜索用户
func TestUserRepository_SearchUsers(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userMongo.NewMongoUserRepository(db)
	ctx := context.Background()

	// 创建测试用户
	testUser := &usersModel.User{
		Username: "searchtest",
		Email:    "search@example.com",
		Password: "hash",
		Status:   usersModel.UserStatusActive,
	}
	err := repo.Create(ctx, testUser)
	require.NoError(t, err)

	// Act - 搜索用户
	result, err := repo.SearchUsers(ctx, "searchtest", 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	// 验证搜索结果包含关键词
	found := false
	for _, user := range result {
		if user.ID == testUser.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "搜索结果应包含创建的测试用户")
}

// TestUserRepository_Health 测试健康检查
func TestUserRepository_Health(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userMongo.NewMongoUserRepository(db)
	ctx := context.Background()

	// Act
	err := repo.Health(ctx)

	// Assert
	assert.NoError(t, err)
}

// TestUserRepository_ExistsByEmail 测试检查邮箱是否存在
func TestUserRepository_ExistsByEmail(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userMongo.NewMongoUserRepository(db)
	ctx := context.Background()

	testUser := &usersModel.User{
		Username: "existsTest",
		Email:    "exists@example.com",
		Password: "hash",
		Status:   usersModel.UserStatusActive,
	}
	err := repo.Create(ctx, testUser)
	require.NoError(t, err)

	// Act
	exists, err := repo.ExistsByEmail(ctx, "exists@example.com")

	// Assert
	require.NoError(t, err)
	assert.True(t, exists)
}

// TestUserRepository_ExistsByUsername 测试检查用户名是否存在
func TestUserRepository_ExistsByUsername(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userMongo.NewMongoUserRepository(db)
	ctx := context.Background()

	testUser := &usersModel.User{
		Username: "existsuser",
		Email:    "existsuser@example.com",
		Password: "hash",
		Status:   usersModel.UserStatusActive,
	}
	err := repo.Create(ctx, testUser)
	require.NoError(t, err)

	// Act
	exists, err := repo.ExistsByUsername(ctx, "existsuser")

	// Assert
	require.NoError(t, err)
	assert.True(t, exists)
}

// TestUserRepository_SetEmailVerified 测试设置邮箱验证状态
func TestUserRepository_SetEmailVerified(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userMongo.NewMongoUserRepository(db)
	ctx := context.Background()

	testUser := &usersModel.User{
		Username:      "verifytest",
		Email:         "verify@example.com",
		Password:      "hash",
		Status:        usersModel.UserStatusActive,
		EmailVerified: false,
	}
	err := repo.Create(ctx, testUser)
	require.NoError(t, err)

	// Act - 设置为已验证
	err = repo.SetEmailVerified(ctx, testUser.ID, true)

	// Assert
	require.NoError(t, err)

	// 验证状态
	found, err := repo.GetByID(ctx, testUser.ID)
	require.NoError(t, err)
	assert.True(t, found.EmailVerified)
}

// TestUserRepository_BatchUpdateStatus 测试批量更新用户状态
func TestUserRepository_BatchUpdateStatus(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userMongo.NewMongoUserRepository(db)
	ctx := context.Background()

	// 创建多个测试用户
	userIDs := []string{}
	for i := 0; i < 3; i++ {
		user := &usersModel.User{
			Username: fmt.Sprintf("batchuser_%d", i),
			Email:    fmt.Sprintf("batch%d@example.com", i),
			Password: "hash",
			Status:   usersModel.UserStatusActive,
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)
		userIDs = append(userIDs, user.ID)
	}

	// Act - 批量更新状态
	err := repo.BatchUpdateStatus(ctx, userIDs, usersModel.UserStatusBanned)

	// Assert
	require.NoError(t, err)

	// 验证所有用户都已更新
	for _, userID := range userIDs {
		found, err := repo.GetByID(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, usersModel.UserStatusBanned, found.Status)
	}
}

// TestUserRepository_GetActiveUsers 测试获取活跃用户列表
func TestUserRepository_GetActiveUsers(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userMongo.NewMongoUserRepository(db)
	ctx := context.Background()

	// 创建活跃用户
	for i := 0; i < 3; i++ {
		user := &usersModel.User{
			Username: fmt.Sprintf("activeuser_%d", i),
			Email:    fmt.Sprintf("active%d@example.com", i),
			Password: "hash",
			Status:   usersModel.UserStatusActive,
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)
	}

	// Act
	result, err := repo.GetActiveUsers(ctx, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result), 3)

	// 验证都是活跃用户
	for _, user := range result {
		assert.Equal(t, usersModel.UserStatusActive, user.Status)
	}
}

// TestUserRepository_CountByStatus 测试根据状态统计用户
func TestUserRepository_CountByStatus(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := userMongo.NewMongoUserRepository(db)
	ctx := context.Background()

	// 创建活跃用户
	for i := 0; i < 5; i++ {
		user := &usersModel.User{
			Username: fmt.Sprintf("active_%d", i),
			Email:    fmt.Sprintf("active_%d@example.com", i),
			Password: "hash",
			Status:   usersModel.UserStatusActive,
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)
	}

	// Act
	count, err := repo.CountByStatus(ctx, usersModel.UserStatusActive)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(5))
}
