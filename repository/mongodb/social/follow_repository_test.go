package social_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	socialModel "Qingyu_backend/models/social"
	mongoSocial "Qingyu_backend/repository/mongodb/social"
	"Qingyu_backend/test/testutil"
)

// setupFollowRepo 测试辅助函数
func setupFollowRepo(t *testing.T) (*mongoSocial.MongoFollowRepository, context.Context, func()) {
	db, cleanup := testutil.SetupTestDB(t)
	repoInterface := mongoSocial.NewMongoFollowRepository(db)
	repo := repoInterface.(*mongoSocial.MongoFollowRepository)
	ctx := context.Background()

	// 在cleanup中添加follows和author_follows集合的清理
	originalCleanup := cleanup
	newCleanup := func() {
		ctx := context.Background()
		_ = db.Collection("follows").Drop(ctx)
		_ = db.Collection("author_follows").Drop(ctx)
		originalCleanup()
	}

	return repo, ctx, newCleanup
}

// ========== 用户关注测试 ==========

// TestFollowRepository_CreateFollow 测试创建关注关系
func TestFollowRepository_CreateFollow(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	follow := &socialModel.Follow{
		FollowerID:  "user123",
		FollowingID: "user456",
		FollowType:  "user",
	}

	// Act
	err := repo.CreateFollow(ctx, follow)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, follow.ID)
	assert.NotZero(t, follow.CreatedAt)
	assert.NotZero(t, follow.UpdatedAt)
}

// TestFollowRepository_CreateFollow_Duplicate 测试创建重复关注
func TestFollowRepository_CreateFollow_Duplicate(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	follow := &socialModel.Follow{
		FollowerID:  "user123",
		FollowingID: "user456",
		FollowType:  "user",
	}
	err := repo.CreateFollow(ctx, follow)
	require.NoError(t, err)

	// Act - 尝试创建相同的关注关系
	duplicateFollow := &socialModel.Follow{
		FollowerID:  "user123",
		FollowingID: "user456",
		FollowType:  "user",
	}
	err = repo.CreateFollow(ctx, duplicateFollow)

	// Assert - 应该返回错误（唯一索引约束）
	assert.Error(t, err)
}

// TestFollowRepository_DeleteFollow 测试删除关注关系
func TestFollowRepository_DeleteFollow(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	follow := &socialModel.Follow{
		FollowerID:  "user123",
		FollowingID: "user456",
		FollowType:  "user",
	}
	err := repo.CreateFollow(ctx, follow)
	require.NoError(t, err)

	// Act
	err = repo.DeleteFollow(ctx, "user123", "user456", "user")

	// Assert
	require.NoError(t, err)

	// 验证已删除
	found, err := repo.GetFollow(ctx, "user123", "user456", "user")
	assert.NoError(t, err)
	assert.Nil(t, found)
}

// TestFollowRepository_GetFollow 测试获取关注关系
func TestFollowRepository_GetFollow(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	follow := &socialModel.Follow{
		FollowerID:  "user123",
		FollowingID: "user456",
		FollowType:  "user",
	}
	err := repo.CreateFollow(ctx, follow)
	require.NoError(t, err)

	// Act
	found, err := repo.GetFollow(ctx, "user123", "user456", "user")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, "user123", found.FollowerID)
	assert.Equal(t, "user456", found.FollowingID)
	assert.Equal(t, "user", found.FollowType)
}

// TestFollowRepository_GetFollow_NotFound 测试获取不存在的关注关系
func TestFollowRepository_GetFollow_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	// Act
	found, err := repo.GetFollow(ctx, "nonexistent1", "nonexistent2", "user")

	// Assert
	require.NoError(t, err)
	assert.Nil(t, found)
}

// TestFollowRepository_IsFollowing 测试检查是否已关注
func TestFollowRepository_IsFollowing(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	follow := &socialModel.Follow{
		FollowerID:  "user123",
		FollowingID: "user456",
		FollowType:  "user",
	}
	err := repo.CreateFollow(ctx, follow)
	require.NoError(t, err)

	// Act
	isFollowing, err := repo.IsFollowing(ctx, "user123", "user456", "user")

	// Assert
	require.NoError(t, err)
	assert.True(t, isFollowing)
}

// TestFollowRepository_IsFollowing_NotFollowing 测试检查未关注状态
func TestFollowRepository_IsFollowing_NotFollowing(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	// Act
	isFollowing, err := repo.IsFollowing(ctx, "user123", "user456", "user")

	// Assert
	require.NoError(t, err)
	assert.False(t, isFollowing)
}

// TestFollowRepository_GetFollowers 测试获取粉丝列表
func TestFollowRepository_GetFollowers(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	userID := "target_user"

	// 创建多个粉丝
	for i := 0; i < 5; i++ {
		follow := &socialModel.Follow{
			FollowerID:  primitive.NewObjectID().Hex(),
			FollowingID: userID,
			FollowType:  "user",
		}
		err := repo.CreateFollow(ctx, follow)
		require.NoError(t, err)
	}

	// Act
	followers, total, err := repo.GetFollowers(ctx, userID, "user", 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, followers)
	assert.GreaterOrEqual(t, total, int64(5))
	assert.GreaterOrEqual(t, len(followers), 5)
}

// TestFollowRepository_GetFollowers_Empty 测试获取空粉丝列表
func TestFollowRepository_GetFollowers_Empty(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	// Act
	followers, total, err := repo.GetFollowers(ctx, "nonexistent_user", "user", 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, followers)
	assert.Equal(t, int64(0), total)
	assert.Equal(t, 0, len(followers))
}

// TestFollowRepository_GetFollowers_Pagination 测试粉丝列表分页
func TestFollowRepository_GetFollowers_Pagination(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	userID := "target_user"

	// 创建15个粉丝
	for i := 0; i < 15; i++ {
		follow := &socialModel.Follow{
			FollowerID:  primitive.NewObjectID().Hex(),
			FollowingID: userID,
			FollowType:  "user",
		}
		err := repo.CreateFollow(ctx, follow)
		require.NoError(t, err)
	}

	// Act - 获取第一页
	page1, total1, err := repo.GetFollowers(ctx, userID, "user", 1, 10)
	require.NoError(t, err)

	// Act - 获取第二页
	page2, total2, err := repo.GetFollowers(ctx, userID, "user", 2, 10)

	// Assert
	assert.Equal(t, int64(15), total1)
	assert.Equal(t, int64(15), total2)
	assert.LessOrEqual(t, len(page1), 10)
	assert.LessOrEqual(t, len(page2), 10)
}

// TestFollowRepository_GetFollowing 测试获取关注列表
func TestFollowRepository_GetFollowing(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	userID := "follower_user"

	// 创建多个关注
	for i := 0; i < 5; i++ {
		follow := &socialModel.Follow{
			FollowerID:  userID,
			FollowingID: primitive.NewObjectID().Hex(),
			FollowType:  "user",
		}
		err := repo.CreateFollow(ctx, follow)
		require.NoError(t, err)
	}

	// Act
	following, total, err := repo.GetFollowing(ctx, userID, "user", 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, following)
	assert.GreaterOrEqual(t, total, int64(5))
	assert.GreaterOrEqual(t, len(following), 5)
}

// TestFollowRepository_GetFollowing_Empty 测试获取空关注列表
func TestFollowRepository_GetFollowing_Empty(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	// Act
	following, total, err := repo.GetFollowing(ctx, "nonexistent_user", "user", 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, following)
	assert.Equal(t, int64(0), total)
	assert.Equal(t, 0, len(following))
}

// TestFollowRepository_UpdateMutualStatus 测试更新互相关注状态
func TestFollowRepository_UpdateMutualStatus(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	follow := &socialModel.Follow{
		FollowerID:  "user123",
		FollowingID: "user456",
		FollowType:  "user",
		IsMutual:    false,
	}
	err := repo.CreateFollow(ctx, follow)
	require.NoError(t, err)

	// Act - 更新为互相关注
	err = repo.UpdateMutualStatus(ctx, "user123", "user456", "user", true)

	// Assert
	require.NoError(t, err)

	// 验证状态已更新
	found, err := repo.GetFollow(ctx, "user123", "user456", "user")
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.True(t, found.IsMutual)
}

// ========== 作者关注测试 ==========

// TestFollowRepository_CreateAuthorFollow 测试创建作者关注
func TestFollowRepository_CreateAuthorFollow(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	authorFollow := &socialModel.AuthorFollow{
		UserID:        "user123",
		AuthorID:      "author456",
		AuthorName:    "Test Author",
		AuthorAvatar:  "avatar.jpg",
		NotifyNewBook: true,
	}

	// Act
	err := repo.CreateAuthorFollow(ctx, authorFollow)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, authorFollow.ID)
	assert.NotZero(t, authorFollow.CreatedAt)
	assert.NotZero(t, authorFollow.UpdatedAt)
}

// TestFollowRepository_CreateAuthorFollow_Duplicate 测试创建重复作者关注
func TestFollowRepository_CreateAuthorFollow_Duplicate(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	authorFollow := &socialModel.AuthorFollow{
		UserID:   "user123",
		AuthorID: "author456",
	}
	err := repo.CreateAuthorFollow(ctx, authorFollow)
	require.NoError(t, err)

	// Act - 尝试创建相同的作者关注
	duplicateFollow := &socialModel.AuthorFollow{
		UserID:   "user123",
		AuthorID: "author456",
	}
	err = repo.CreateAuthorFollow(ctx, duplicateFollow)

	// Assert - 应该返回错误（唯一索引约束）
	assert.Error(t, err)
}

// TestFollowRepository_DeleteAuthorFollow 测试删除作者关注
func TestFollowRepository_DeleteAuthorFollow(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	authorFollow := &socialModel.AuthorFollow{
		UserID:   "user123",
		AuthorID: "author456",
	}
	err := repo.CreateAuthorFollow(ctx, authorFollow)
	require.NoError(t, err)

	// Act
	err = repo.DeleteAuthorFollow(ctx, "user123", "author456")

	// Assert
	require.NoError(t, err)

	// 验证已删除
	found, err := repo.GetAuthorFollow(ctx, "user123", "author456")
	assert.NoError(t, err)
	assert.Nil(t, found)
}

// TestFollowRepository_GetAuthorFollow 测试获取作者关注
func TestFollowRepository_GetAuthorFollow(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	authorFollow := &socialModel.AuthorFollow{
		UserID:     "user123",
		AuthorID:   "author456",
		AuthorName: "Test Author",
	}
	err := repo.CreateAuthorFollow(ctx, authorFollow)
	require.NoError(t, err)

	// Act
	found, err := repo.GetAuthorFollow(ctx, "user123", "author456")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, "user123", found.UserID)
	assert.Equal(t, "author456", found.AuthorID)
	assert.Equal(t, "Test Author", found.AuthorName)
}

// TestFollowRepository_GetAuthorFollow_NotFound 测试获取不存在的作者关注
func TestFollowRepository_GetAuthorFollow_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	// Act
	found, err := repo.GetAuthorFollow(ctx, "nonexistent_user", "nonexistent_author")

	// Assert
	require.NoError(t, err)
	assert.Nil(t, found)
}

// TestFollowRepository_GetAuthorFollowers 测试获取作者粉丝列表
func TestFollowRepository_GetAuthorFollowers(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	authorID := "author123"

	// 创建多个粉丝
	for i := 0; i < 5; i++ {
		authorFollow := &socialModel.AuthorFollow{
			UserID:     primitive.NewObjectID().Hex(),
			AuthorID:   authorID,
			AuthorName: "Test Author",
		}
		err := repo.CreateAuthorFollow(ctx, authorFollow)
		require.NoError(t, err)
	}

	// Act
	followers, total, err := repo.GetAuthorFollowers(ctx, authorID, 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, followers)
	assert.GreaterOrEqual(t, total, int64(5))
	assert.GreaterOrEqual(t, len(followers), 5)
}

// TestFollowRepository_GetAuthorFollowers_Empty 测试获取空的作者粉丝列表
func TestFollowRepository_GetAuthorFollowers_Empty(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	// Act
	followers, total, err := repo.GetAuthorFollowers(ctx, "nonexistent_author", 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, followers)
	assert.Equal(t, int64(0), total)
	assert.Equal(t, 0, len(followers))
}

// TestFollowRepository_GetUserFollowingAuthors 测试获取用户关注的作者列表
func TestFollowRepository_GetUserFollowingAuthors(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	userID := "user123"

	// 创建多个作者关注
	for i := 0; i < 5; i++ {
		authorFollow := &socialModel.AuthorFollow{
			UserID:     userID,
			AuthorID:   primitive.NewObjectID().Hex(),
			AuthorName: "Test Author",
		}
		err := repo.CreateAuthorFollow(ctx, authorFollow)
		require.NoError(t, err)
	}

	// Act
	authors, total, err := repo.GetUserFollowingAuthors(ctx, userID, 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, authors)
	assert.GreaterOrEqual(t, total, int64(5))
	assert.GreaterOrEqual(t, len(authors), 5)
}

// TestFollowRepository_GetUserFollowingAuthors_Empty 测试获取空的关注作者列表
func TestFollowRepository_GetUserFollowingAuthors_Empty(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	// Act
	authors, total, err := repo.GetUserFollowingAuthors(ctx, "nonexistent_user", 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, authors)
	assert.Equal(t, int64(0), total)
	assert.Equal(t, 0, len(authors))
}

// ========== 统计测试 ==========

// TestFollowRepository_GetFollowStats 测试获取关注统计
func TestFollowRepository_GetFollowStats(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	userID := "user123"

	// 创建粉丝
	for i := 0; i < 5; i++ {
		follow := &socialModel.Follow{
			FollowerID:  primitive.NewObjectID().Hex(),
			FollowingID: userID,
			FollowType:  "user",
		}
		err := repo.CreateFollow(ctx, follow)
		require.NoError(t, err)
	}

	// 创建关注
	for i := 0; i < 3; i++ {
		follow := &socialModel.Follow{
			FollowerID:  userID,
			FollowingID: primitive.NewObjectID().Hex(),
			FollowType:  "user",
		}
		err := repo.CreateFollow(ctx, follow)
		require.NoError(t, err)
	}

	// Act
	stats, err := repo.GetFollowStats(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, userID, stats.UserID)
	assert.Equal(t, 5, stats.FollowerCount)
	assert.Equal(t, 3, stats.FollowingCount)
	assert.NotZero(t, stats.UpdatedAt)
}

// TestFollowRepository_GetFollowStats_NoFollows 测试获取无关注的用户统计
func TestFollowRepository_GetFollowStats_NoFollows(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	userID := "nonexistent_user"

	// Act
	stats, err := repo.GetFollowStats(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, userID, stats.UserID)
	assert.Equal(t, 0, stats.FollowerCount)
	assert.Equal(t, 0, stats.FollowingCount)
}

// TestFollowRepository_CountFollowers 测试统计粉丝数
func TestFollowRepository_CountFollowers(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	userID := "user123"

	// 创建粉丝
	for i := 0; i < 7; i++ {
		follow := &socialModel.Follow{
			FollowerID:  primitive.NewObjectID().Hex(),
			FollowingID: userID,
			FollowType:  "user",
		}
		err := repo.CreateFollow(ctx, follow)
		require.NoError(t, err)
	}

	// Act
	count, err := repo.CountFollowers(ctx, userID, "user")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(7), count)
}

// TestFollowRepository_CountFollowing 测试统计关注数
func TestFollowRepository_CountFollowing(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	userID := "user123"

	// 创建关注
	for i := 0; i < 4; i++ {
		follow := &socialModel.Follow{
			FollowerID:  userID,
			FollowingID: primitive.NewObjectID().Hex(),
			FollowType:  "user",
		}
		err := repo.CreateFollow(ctx, follow)
		require.NoError(t, err)
	}

	// Act
	count, err := repo.CountFollowing(ctx, userID, "user")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(4), count)
}

// TestFollowRepository_UpdateFollowStats 测试更新关注统计
func TestFollowRepository_UpdateFollowStats(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	userID := "user123"

	// Act
	err := repo.UpdateFollowStats(ctx, userID, 5, -3)

	// Assert - 根据实现，这个方法总是返回nil（简化实现）
	require.NoError(t, err)
}

// ========== 健康检查测试 ==========

// TestFollowRepository_Health 测试健康检查
func TestFollowRepository_Health(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	// Act
	err := repo.Health(ctx)

	// Assert
	assert.NoError(t, err)
}

// TestFollowRepository_MixedFollowTypes 测试混合关注类型
func TestFollowRepository_MixedFollowTypes(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	userID := "user123"
	otherUserID := "user456"

	// 创建user类型关注
	userFollow := &socialModel.Follow{
		FollowerID:  userID,
		FollowingID: otherUserID,
		FollowType:  "user",
	}
	err := repo.CreateFollow(ctx, userFollow)
	require.NoError(t, err)

	// 创建author类型关注（同一对用户）
	authorFollow := &socialModel.Follow{
		FollowerID:  userID,
		FollowingID: otherUserID,
		FollowType:  "author",
	}
	err = repo.CreateFollow(ctx, authorFollow)
	require.NoError(t, err)

	// Act & Assert - 验证两种类型都存在
	isUserFollowing, err := repo.IsFollowing(ctx, userID, otherUserID, "user")
	require.NoError(t, err)
	assert.True(t, isUserFollowing)

	isAuthorFollowing, err := repo.IsFollowing(ctx, userID, otherUserID, "author")
	require.NoError(t, err)
	assert.True(t, isAuthorFollowing)

	// 验证统计
	userCount, err := repo.CountFollowing(ctx, userID, "user")
	require.NoError(t, err)
	assert.Equal(t, int64(1), userCount)

	authorCount, err := repo.CountFollowing(ctx, userID, "author")
	require.NoError(t, err)
	assert.Equal(t, int64(1), authorCount)
}

// TestFollowRepository_DeleteNonExistentFollow 测试删除不存在的关注关系
func TestFollowRepository_DeleteNonExistentFollow(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	// Act - 尝试删除不存在的关注关系
	err := repo.DeleteFollow(ctx, "nonexistent1", "nonexistent2", "user")

	// Assert - 删除操作应该成功（即使不存在）
	assert.NoError(t, err)
}

// TestFollowRepository_DeleteNonExistentAuthorFollow 测试删除不存在的作者关注
func TestFollowRepository_DeleteNonExistentAuthorFollow(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupFollowRepo(t)
	defer cleanup()

	// Act - 尝试删除不存在的作者关注
	err := repo.DeleteAuthorFollow(ctx, "nonexistent_user", "nonexistent_author")

	// Assert - 删除操作应该成功（即使不存在）
	assert.NoError(t, err)
}
