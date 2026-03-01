package shared_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/models/users"
	"Qingyu_backend/repository/mongodb/shared"
	"Qingyu_backend/test/testutil"
)

// TestMongoUserRepository_Create 测试创建用户
func TestMongoUserRepository_Create(t *testing.T) {
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewMongoUserRepository(db)
	ctx := context.Background()

	user := &users.User{
		Username:      "testuser",
		Email:         "test@example.com",
		Password:      "hashed_password",
		Roles:         []string{"user"},
		Status:        users.UserStatusActive,
		EmailVerified: false,
	}

	err := repo.Create(ctx, user)

	require.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)
}

// TestMongoUserRepository_GetByID 测试根据ID获取用户
func TestMongoUserRepository_GetByID_Success(t *testing.T) {
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewMongoUserRepository(db)
	ctx := context.Background()

	// 先创建用户
	user := &users.User{
		Username: "getbyid",
		Email:    "getbyid@example.com",
		Password: "hashed",
		Status:   users.UserStatusActive,
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// 通过ID获取
	found, err := repo.GetByID(ctx, user.ID.Hex())

	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, "getbyid", found.Username)
}

// TestMongoUserRepository_GetByID_NotFound 测试获取不存在的用户
func TestMongoUserRepository_GetByID_NotFound(t *testing.T) {
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewMongoUserRepository(db)
	ctx := context.Background()

	// 使用不存在的ID
	found, err := repo.GetByID(ctx, "507f1f77bcf86cd799439011")

	assert.Error(t, err)
	assert.Nil(t, found)
}

// TestMongoUserRepository_GetByID_InvalidID 测试无效ID格式
func TestMongoUserRepository_GetByID_InvalidID(t *testing.T) {
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewMongoUserRepository(db)
	ctx := context.Background()

	found, err := repo.GetByID(ctx, "invalid-id-format")

	assert.Error(t, err)
	assert.Nil(t, found)
}

// TestMongoUserRepository_GetByEmail 测试根据邮箱获取用户
func TestMongoUserRepository_GetByEmail_Success(t *testing.T) {
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewMongoUserRepository(db)
	ctx := context.Background()

	user := &users.User{
		Username: "emailtest",
		Email:    "email@example.com",
		Password: "hashed",
		Status:   users.UserStatusActive,
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	found, err := repo.GetByEmail(ctx, "email@example.com")

	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, "emailtest", found.Username)
}

// TestMongoUserRepository_UpdateStatus 测试更新用户状态
func TestMongoUserRepository_UpdateStatus_Success(t *testing.T) {
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewMongoUserRepository(db)
	ctx := context.Background()

	user := &users.User{
		Username: "statusupdatetest",
		Email:    "statusupdate@example.com",
		Password: "hashed",
		Status:   users.UserStatusActive,
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// 更新状态为封禁
	err = repo.UpdateStatus(ctx, user.ID.Hex(), users.UserStatusBanned)

	require.NoError(t, err)

	// 验证状态已更新
	updated, _ := repo.GetByID(ctx, user.ID.Hex())
	assert.Equal(t, users.UserStatusBanned, updated.Status)
}

// TestMongoUserRepository_BatchUpdateStatus 测试批量更新状态
func TestMongoUserRepository_BatchUpdateStatus_Success(t *testing.T) {
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewMongoUserRepository(db)
	ctx := context.Background()

	// 创建3个测试用户
	var ids []string
	for i := 0; i < 3; i++ {
		user := &users.User{
			Username: "batchuser",
			Email:    "batch@example.com",
			Password: "hashed",
			Status:   users.UserStatusActive,
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)
		ids = append(ids, user.ID.Hex())
	}

	// 批量更新为封禁状态
	err := repo.BatchUpdateStatus(ctx, ids, users.UserStatusBanned)

	require.NoError(t, err)

	// 验证所有用户状态已更新
	for _, id := range ids {
		user, _ := repo.GetByID(ctx, id)
		assert.Equal(t, users.UserStatusBanned, user.Status)
	}
}

// TestMongoUserRepository_BatchUpdateStatus_EmptyList 测试空列表批量更新
func TestMongoUserRepository_BatchUpdateStatus_EmptyList(t *testing.T) {
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewMongoUserRepository(db)
	ctx := context.Background()

	err := repo.BatchUpdateStatus(ctx, []string{}, users.UserStatusBanned)

	assert.Error(t, err)
}

// TestMongoUserRepository_CountByStatus 测试按状态统计用户
func TestMongoUserRepository_CountByStatus(t *testing.T) {
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := shared.NewMongoUserRepository(db)
	ctx := context.Background()

	// 创建不同状态的用户
	activeUser := &users.User{Username: "active", Email: "active@test.com", Password: "hash", Status: users.UserStatusActive}
	bannedUser := &users.User{Username: "banned", Email: "banned@test.com", Password: "hash", Status: users.UserStatusBanned}

	require.NoError(t, repo.Create(ctx, activeUser))
	require.NoError(t, repo.Create(ctx, bannedUser))

	// 统计活跃用户数量
	activeCount, err := repo.CountByStatus(ctx, users.UserStatusActive)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, activeCount, int64(1))
}
