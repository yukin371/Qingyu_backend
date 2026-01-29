package social_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	socialModel "Qingyu_backend/models/social"
	socialRepo "Qingyu_backend/repository/interfaces/social"
	impl "Qingyu_backend/repository/mongodb/social"
	"Qingyu_backend/test/testutil"
)

// setupUserRelationRepo 测试辅助函数
func setupUserRelationRepo(t *testing.T) (socialRepo.UserRelationRepository, context.Context, func()) {
	db, cleanup := testutil.SetupTestDB(t)
	repo := impl.NewMongoUserRelationRepository(db)
	ctx := context.Background()
	return repo, ctx, cleanup
}

// TestUserRelationRepository_Create 测试创建用户关系
func TestUserRelationRepository_Create(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	relation := &socialModel.UserRelation{
		FollowerID: "user123",
		FolloweeID: "user456",
		Status:     socialModel.RelationStatusActive,
	}

	// Act
	err := repo.Create(ctx, relation)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, relation.ID)
	assert.NotZero(t, relation.CreatedAt)
	assert.NotZero(t, relation.UpdatedAt)
}

// TestUserRelationRepository_Create_Multiple 测试创建多个用户关系
func TestUserRelationRepository_Create_Multiple(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	// 创建多个不同的关系
	relations := []*socialModel.UserRelation{
		{
			FollowerID: "user1",
			FolloweeID: "user2",
			Status:     socialModel.RelationStatusActive,
		},
		{
			FollowerID: "user1",
			FolloweeID: "user3",
			Status:     socialModel.RelationStatusActive,
		},
		{
			FollowerID: "user2",
			FolloweeID: "user3",
			Status:     socialModel.RelationStatusActive,
		},
	}

	// Act
	for _, relation := range relations {
		err := repo.Create(ctx, relation)
		require.NoError(t, err)
	}

	// Assert - 验证所有关系都已创建
	count, err := repo.CountFollowing(ctx, "user1")
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

// TestUserRelationRepository_GetByID 测试根据ID获取用户关系
func TestUserRelationRepository_GetByID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	relation := &socialModel.UserRelation{
		FollowerID: "user123",
		FolloweeID: "user456",
		Status:     socialModel.RelationStatusActive,
	}
	err := repo.Create(ctx, relation)
	require.NoError(t, err)

	// Act
	found, err := repo.GetByID(ctx, relation.ID.Hex())

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, relation.FollowerID, found.FollowerID)
	assert.Equal(t, relation.FolloweeID, found.FolloweeID)
	assert.Equal(t, relation.Status, found.Status)
}

// TestUserRelationRepository_GetByID_NotFound 测试获取不存在的关系
func TestUserRelationRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	// Act
	found, err := repo.GetByID(ctx, "nonexistent_id")

	// Assert
	require.Error(t, err)
	assert.Nil(t, found)
}

// TestUserRelationRepository_GetByID_InvalidID 测试使用无效ID获取关系
func TestUserRelationRepository_GetByID_InvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	// Act
	found, err := repo.GetByID(ctx, "invalid_object_id")

	// Assert
	require.Error(t, err)
	assert.Nil(t, found)
}

// TestUserRelationRepository_Update 测试更新用户关系
func TestUserRelationRepository_Update(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	relation := &socialModel.UserRelation{
		FollowerID: "user123",
		FolloweeID: "user456",
		Status:     socialModel.RelationStatusActive,
	}
	err := repo.Create(ctx, relation)
	require.NoError(t, err)

	// Act - 更新关系状态
	updates := map[string]interface{}{
		"status": socialModel.RelationStatusInactive,
	}
	err = repo.Update(ctx, relation.ID.Hex(), updates)

	// Assert
	require.NoError(t, err)

	// 验证更新
	found, err := repo.GetByID(ctx, relation.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, socialModel.RelationStatusInactive, found.Status)
}

// TestUserRelationRepository_Update_InvalidID 测试使用无效ID更新
func TestUserRelationRepository_Update_InvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	// Act
	updates := map[string]interface{}{
		"status": socialModel.RelationStatusInactive,
	}
	err := repo.Update(ctx, "invalid_id", updates)

	// Assert
	require.Error(t, err)
}

// TestUserRelationRepository_Delete 测试删除用户关系
func TestUserRelationRepository_Delete(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	relation := &socialModel.UserRelation{
		FollowerID: "user123",
		FolloweeID: "user456",
		Status:     socialModel.RelationStatusActive,
	}
	err := repo.Create(ctx, relation)
	require.NoError(t, err)

	// Act - 删除关系
	err = repo.Delete(ctx, relation.ID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证已删除
	found, err := repo.GetByID(ctx, relation.ID.Hex())
	require.Error(t, err)
	assert.Nil(t, found)
}

// TestUserRelationRepository_Delete_InvalidID 测试使用无效ID删除
func TestUserRelationRepository_Delete_InvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	// Act
	err := repo.Delete(ctx, "invalid_id")

	// Assert
	require.Error(t, err)
}

// TestUserRelationRepository_GetRelation 测试获取用户关系
func TestUserRelationRepository_GetRelation(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	relation := &socialModel.UserRelation{
		FollowerID: "user123",
		FolloweeID: "user456",
		Status:     socialModel.RelationStatusActive,
	}
	err := repo.Create(ctx, relation)
	require.NoError(t, err)

	// Act
	found, err := repo.GetRelation(ctx, "user123", "user456")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, "user123", found.FollowerID)
	assert.Equal(t, "user456", found.FolloweeID)
}

// TestUserRelationRepository_GetRelation_NotFound 测试获取不存在的关系
func TestUserRelationRepository_GetRelation_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	// Act
	found, err := repo.GetRelation(ctx, "user123", "user456")

	// Assert
	require.Error(t, err)
	assert.Nil(t, found)
}

// TestUserRelationRepository_IsFollowing 测试检查是否关注
func TestUserRelationRepository_IsFollowing(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	relation := &socialModel.UserRelation{
		FollowerID: "user123",
		FolloweeID: "user456",
		Status:     socialModel.RelationStatusActive,
	}
	err := repo.Create(ctx, relation)
	require.NoError(t, err)

	// Act
	isFollowing, err := repo.IsFollowing(ctx, "user123", "user456")

	// Assert
	require.NoError(t, err)
	assert.True(t, isFollowing)
}

// TestUserRelationRepository_IsFollowing_NotFollowing 测试检查未关注的情况
func TestUserRelationRepository_IsFollowing_NotFollowing(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	// Act
	isFollowing, err := repo.IsFollowing(ctx, "user123", "user456")

	// Assert
	require.NoError(t, err)
	assert.False(t, isFollowing)
}

// TestUserRelationRepository_IsFollowing_InactiveStatus 测试非活跃状态不计数
func TestUserRelationRepository_IsFollowing_InactiveStatus(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	relation := &socialModel.UserRelation{
		FollowerID: "user123",
		FolloweeID: "user456",
		Status:     socialModel.RelationStatusInactive,
	}
	err := repo.Create(ctx, relation)
	require.NoError(t, err)

	// Act
	isFollowing, err := repo.IsFollowing(ctx, "user123", "user456")

	// Assert
	require.NoError(t, err)
	assert.False(t, isFollowing, "非活跃状态不应该被认为是关注")
}

// TestUserRelationRepository_GetFollowers 测试获取粉丝列表
func TestUserRelationRepository_GetFollowers(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	followeeID := "user_target"

	// 创建多个粉丝
	followerIDs := []string{"user1", "user2", "user3"}
	for _, followerID := range followerIDs {
		relation := &socialModel.UserRelation{
			FollowerID: followerID,
			FolloweeID: followeeID,
			Status:     socialModel.RelationStatusActive,
		}
		err := repo.Create(ctx, relation)
		require.NoError(t, err)
	}

	// Act
	followers, total, err := repo.GetFollowers(ctx, followeeID, 10, 0)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, followers)
	assert.Equal(t, int64(3), total)
	assert.Len(t, followers, 3)
}

// TestUserRelationRepository_GetFollowers_WithPagination 测试分页获取粉丝
func TestUserRelationRepository_GetFollowers_WithPagination(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	followeeID := "user_target"

	// 创建5个粉丝
	for i := 1; i <= 5; i++ {
		relation := &socialModel.UserRelation{
			FollowerID: "user" + string(rune('0'+i)),
			FolloweeID: followeeID,
			Status:     socialModel.RelationStatusActive,
		}
		err := repo.Create(ctx, relation)
		require.NoError(t, err)
	}

	// Act - 获取第一页
	followers1, total1, err := repo.GetFollowers(ctx, followeeID, 2, 0)
	require.NoError(t, err)

	// Act - 获取第二页
	followers2, total2, err := repo.GetFollowers(ctx, followeeID, 2, 2)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, int64(5), total1)
	assert.Equal(t, int64(5), total2)
	assert.Len(t, followers1, 2)
	assert.Len(t, followers2, 2)

	// 确保两页的数据不重复
	followerIDs1 := make([]string, len(followers1))
	for i, f := range followers1 {
		followerIDs1[i] = f.FollowerID
	}
	for _, f := range followers2 {
		assert.NotContains(t, followerIDs1, f.FollowerID)
	}
}

// TestUserRelationRepository_GetFollowers_OnlyActive 测试只返回活跃状态的粉丝
func TestUserRelationRepository_GetFollowers_OnlyActive(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	followeeID := "user_target"

	// 创建活跃粉丝
	activeRelation := &socialModel.UserRelation{
		FollowerID: "user_active",
		FolloweeID: followeeID,
		Status:     socialModel.RelationStatusActive,
	}
	err := repo.Create(ctx, activeRelation)
	require.NoError(t, err)

	// 创建非活跃粉丝
	inactiveRelation := &socialModel.UserRelation{
		FollowerID: "user_inactive",
		FolloweeID: followeeID,
		Status:     socialModel.RelationStatusInactive,
	}
	err = repo.Create(ctx, inactiveRelation)
	require.NoError(t, err)

	// Act
	followers, total, err := repo.GetFollowers(ctx, followeeID, 10, 0)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, followers, 1)
	assert.Equal(t, "user_active", followers[0].FollowerID)
}

// TestUserRelationRepository_GetFollowing 测试获取关注列表
func TestUserRelationRepository_GetFollowing(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	followerID := "user_follower"

	// 创建多个关注
	followeeIDs := []string{"user1", "user2", "user3"}
	for _, followeeID := range followeeIDs {
		relation := &socialModel.UserRelation{
			FollowerID: followerID,
			FolloweeID: followeeID,
			Status:     socialModel.RelationStatusActive,
		}
		err := repo.Create(ctx, relation)
		require.NoError(t, err)
	}

	// Act
	following, total, err := repo.GetFollowing(ctx, followerID, 10, 0)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, following)
	assert.Equal(t, int64(3), total)
	assert.Len(t, following, 3)
}

// TestUserRelationRepository_GetFollowing_WithPagination 测试分页获取关注
func TestUserRelationRepository_GetFollowing_WithPagination(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	followerID := "user_follower"

	// 创建5个关注
	for i := 1; i <= 5; i++ {
		relation := &socialModel.UserRelation{
			FollowerID: followerID,
			FolloweeID: "user" + string(rune('0'+i)),
			Status:     socialModel.RelationStatusActive,
		}
		err := repo.Create(ctx, relation)
		require.NoError(t, err)
	}

	// Act - 获取第一页
	following1, total1, err := repo.GetFollowing(ctx, followerID, 2, 0)
	require.NoError(t, err)

	// Act - 获取第二页
	following2, total2, err := repo.GetFollowing(ctx, followerID, 2, 2)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, int64(5), total1)
	assert.Equal(t, int64(5), total2)
	assert.Len(t, following1, 2)
	assert.Len(t, following2, 2)

	// 确保两页的数据不重复
	followeeIDs1 := make([]string, len(following1))
	for i, f := range following1 {
		followeeIDs1[i] = f.FolloweeID
	}
	for _, f := range following2 {
		assert.NotContains(t, followeeIDs1, f.FolloweeID)
	}
}

// TestUserRelationRepository_GetFollowing_OnlyActive 测试只返回活跃状态的关注
func TestUserRelationRepository_GetFollowing_OnlyActive(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	followerID := "user_follower"

	// 创建活跃关注
	activeRelation := &socialModel.UserRelation{
		FollowerID: followerID,
		FolloweeID: "user_active",
		Status:     socialModel.RelationStatusActive,
	}
	err := repo.Create(ctx, activeRelation)
	require.NoError(t, err)

	// 创建非活跃关注
	inactiveRelation := &socialModel.UserRelation{
		FollowerID: followerID,
		FolloweeID: "user_inactive",
		Status:     socialModel.RelationStatusInactive,
	}
	err = repo.Create(ctx, inactiveRelation)
	require.NoError(t, err)

	// Act
	following, total, err := repo.GetFollowing(ctx, followerID, 10, 0)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, following, 1)
	assert.Equal(t, "user_active", following[0].FolloweeID)
}

// TestUserRelationRepository_CountFollowers 测试统计粉丝数
func TestUserRelationRepository_CountFollowers(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	userID := "user_target"

	// 创建多个粉丝
	for i := 1; i <= 3; i++ {
		relation := &socialModel.UserRelation{
			FollowerID: "follower" + string(rune('0'+i)),
			FolloweeID: userID,
			Status:     socialModel.RelationStatusActive,
		}
		err := repo.Create(ctx, relation)
		require.NoError(t, err)
	}

	// Act
	count, err := repo.CountFollowers(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

// TestUserRelationRepository_CountFollowers_OnlyActive 测试只统计活跃状态的粉丝
func TestUserRelationRepository_CountFollowers_OnlyActive(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	userID := "user_target"

	// 创建活跃粉丝
	activeRelation := &socialModel.UserRelation{
		FollowerID: "follower_active",
		FolloweeID: userID,
		Status:     socialModel.RelationStatusActive,
	}
	err := repo.Create(ctx, activeRelation)
	require.NoError(t, err)

	// 创建非活跃粉丝
	inactiveRelation := &socialModel.UserRelation{
		FollowerID: "follower_inactive",
		FolloweeID: userID,
		Status:     socialModel.RelationStatusInactive,
	}
	err = repo.Create(ctx, inactiveRelation)
	require.NoError(t, err)

	// Act
	count, err := repo.CountFollowers(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

// TestUserRelationRepository_CountFollowing 测试统计关注数
func TestUserRelationRepository_CountFollowing(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	userID := "user_follower"

	// 创建多个关注
	for i := 1; i <= 3; i++ {
		relation := &socialModel.UserRelation{
			FollowerID: userID,
			FolloweeID: "followee" + string(rune('0'+i)),
			Status:     socialModel.RelationStatusActive,
		}
		err := repo.Create(ctx, relation)
		require.NoError(t, err)
	}

	// Act
	count, err := repo.CountFollowing(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

// TestUserRelationRepository_CountFollowing_OnlyActive 测试只统计活跃状态的关注
func TestUserRelationRepository_CountFollowing_OnlyActive(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	userID := "user_follower"

	// 创建活跃关注
	activeRelation := &socialModel.UserRelation{
		FollowerID: userID,
		FolloweeID: "followee_active",
		Status:     socialModel.RelationStatusActive,
	}
	err := repo.Create(ctx, activeRelation)
	require.NoError(t, err)

	// 创建非活跃关注
	inactiveRelation := &socialModel.UserRelation{
		FollowerID: userID,
		FolloweeID: "followee_inactive",
		Status:     socialModel.RelationStatusInactive,
	}
	err = repo.Create(ctx, inactiveRelation)
	require.NoError(t, err)

	// Act
	count, err := repo.CountFollowing(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

// TestUserRelationRepository_BatchCreate 测试批量创建用户关系
func TestUserRelationRepository_BatchCreate(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	relations := []*socialModel.UserRelation{
		{
			FollowerID: "user1",
			FolloweeID: "user2",
			Status:     socialModel.RelationStatusActive,
		},
		{
			FollowerID: "user3",
			FolloweeID: "user4",
			Status:     socialModel.RelationStatusActive,
		},
		{
			FollowerID: "user5",
			FolloweeID: "user6",
			Status:     socialModel.RelationStatusActive,
		},
	}

	// Act
	err := repo.BatchCreate(ctx, relations)

	// Assert
	require.NoError(t, err)
	for _, relation := range relations {
		assert.NotEmpty(t, relation.ID)
		assert.NotZero(t, relation.CreatedAt)
		assert.NotZero(t, relation.UpdatedAt)
	}

	// 验证数据已保存
	count, err := repo.CountFollowing(ctx, "user1")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(1))
}

// TestUserRelationRepository_BatchCreate_Empty 测试批量创建空列表
func TestUserRelationRepository_BatchCreate_Empty(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	// Act
	err := repo.BatchCreate(ctx, []*socialModel.UserRelation{})

	// Assert
	require.NoError(t, err)
}

// TestUserRelationRepository_BatchCreate_WithNilID 测试批量创建时自动生成ID
func TestUserRelationRepository_BatchCreate_WithNilID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	relations := []*socialModel.UserRelation{
		{
			FollowerID: "user1",
			FolloweeID: "user2",
			Status:     socialModel.RelationStatusActive,
		},
	}

	// Act
	err := repo.BatchCreate(ctx, relations)

	// Assert
	require.NoError(t, err)
	assert.False(t, relations[0].ID.IsZero())
}

// TestUserRelationRepository_BatchCreate_WithExistingID 测试批量创建时保留已有ID
func TestUserRelationRepository_BatchCreate_WithExistingID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	existingID := primitive.NewObjectID()
	relations := []*socialModel.UserRelation{
		{
			IdentifiedEntity: socialModel.IdentifiedEntity{ID: existingID},
			FollowerID:       "user1",
			FolloweeID:       "user2",
			Status:           socialModel.RelationStatusActive,
		},
	}

	// Act
	err := repo.BatchCreate(ctx, relations)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, existingID, relations[0].ID)
}

// TestUserRelationRepository_MutualFollow 测试双向关注场景
func TestUserRelationRepository_MutualFollow(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupUserRelationRepo(t)
	defer cleanup()

	// 用户A关注用户B
	relation1 := &socialModel.UserRelation{
		FollowerID: "userA",
		FolloweeID: "userB",
		Status:     socialModel.RelationStatusActive,
	}
	err := repo.Create(ctx, relation1)
	require.NoError(t, err)

	// 用户B关注用户A
	relation2 := &socialModel.UserRelation{
		FollowerID: "userB",
		FolloweeID: "userA",
		Status:     socialModel.RelationStatusActive,
	}
	err = repo.Create(ctx, relation2)
	require.NoError(t, err)

	// Act & Assert - 验证双向关注
	isAFollowingB, err := repo.IsFollowing(ctx, "userA", "userB")
	require.NoError(t, err)
	assert.True(t, isAFollowingB)

	isBFollowingA, err := repo.IsFollowing(ctx, "userB", "userA")
	require.NoError(t, err)
	assert.True(t, isBFollowingA)

	// 验证粉丝数
	aFollowers, _ := repo.CountFollowers(ctx, "userA")
	bFollowers, _ := repo.CountFollowers(ctx, "userB")
	assert.Equal(t, int64(1), aFollowers)
	assert.Equal(t, int64(1), bFollowers)

	// 验证关注数
	aFollowing, _ := repo.CountFollowing(ctx, "userA")
	bFollowing, _ := repo.CountFollowing(ctx, "userB")
	assert.Equal(t, int64(1), aFollowing)
	assert.Equal(t, int64(1), bFollowing)
}
