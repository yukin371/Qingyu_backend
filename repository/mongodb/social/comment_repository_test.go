package social_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/social"
	socialRepo "Qingyu_backend/repository/mongodb/social"
	"Qingyu_backend/test/testutil"
)

// setupCommentRepo 测试辅助函数
func setupCommentRepo(t *testing.T) (*socialRepo.MongoCommentRepository, context.Context, func()) {
	db, cleanup := testutil.SetupTestDB(t)
	repo := socialRepo.NewMongoCommentRepository(db)
	ctx := context.Background()
	return repo, ctx, cleanup
}

// TestCommentRepository_Create 测试创建评论
func TestCommentRepository_Create(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCommentRepo(t)
	defer cleanup()

	comment := &social.Comment{
		AuthorID:   "user123",
		TargetType: social.CommentTargetTypeBook,
		TargetID:   "book123",
		Content:    "这是一条测试评论",
		Rating:     5,
		State:      social.CommentStateNormal,
	}

	// Act
	err := repo.Create(ctx, comment)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, comment.ID)
	assert.NotZero(t, comment.CreatedAt)
	assert.NotZero(t, comment.UpdatedAt)
}

// TestCommentRepository_GetByID 测试根据ID获取评论
func TestCommentRepository_GetByID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCommentRepo(t)
	defer cleanup()

	comment := &social.Comment{
		AuthorID:   "user123",
		TargetType: social.CommentTargetTypeBook,
		TargetID:   "book123",
		Content:    "测试评论",
		Rating:     4,
		State:      social.CommentStateNormal,
	}
	err := repo.Create(ctx, comment)
	require.NoError(t, err)

	// Act
	found, err := repo.GetByID(ctx, comment.ID.Hex())

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, comment.Content, found.Content)
	assert.Equal(t, comment.Rating, found.Rating)
}

// TestCommentRepository_GetByID_NotFound 测试获取不存在的评论
func TestCommentRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCommentRepo(t)
	defer cleanup()

	// Act - 使用有效的ObjectID格式但不存在
	nonExistentID := primitive.NewObjectID().Hex()
	found, err := repo.GetByID(ctx, nonExistentID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, found)
}

// TestCommentRepository_Update 测试更新评论
func TestCommentRepository_Update(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCommentRepo(t)
	defer cleanup()

	comment := &social.Comment{
		AuthorID:   "user123",
		TargetType: social.CommentTargetTypeBook,
		TargetID:   "book123",
		Content:    "原始内容",
		Rating:     3,
		State:      social.CommentStateNormal,
	}
	err := repo.Create(ctx, comment)
	require.NoError(t, err)

	// Act - 更新评论
	updates := map[string]interface{}{
		"content": "更新后的内容",
		"rating":  5,
	}
	err = repo.Update(ctx, comment.ID.Hex(), updates)

	// Assert
	require.NoError(t, err)

	// 验证更新
	found, err := repo.GetByID(ctx, comment.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, "更新后的内容", found.Content)
	assert.Equal(t, 5, found.Rating)
}

// TestCommentRepository_Delete 测试删除评论
func TestCommentRepository_Delete(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCommentRepo(t)
	defer cleanup()

	comment := &social.Comment{
		AuthorID:   "user123",
		TargetType: social.CommentTargetTypeBook,
		TargetID:   "book123",
		Content:    "待删除的评论",
		State:      social.CommentStateNormal,
	}
	err := repo.Create(ctx, comment)
	require.NoError(t, err)

	// Act - 删除评论
	err = repo.Delete(ctx, comment.ID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证状态已更新为deleted
	found, err := repo.GetByID(ctx, comment.ID.Hex())
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, social.CommentStateDeleted, found.State)
}

// TestCommentRepository_GetCommentsByBookID 测试获取书籍评论列表
func TestCommentRepository_GetCommentsByBookID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCommentRepo(t)
	defer cleanup()

	bookID := "test_book_123"

	// 创建多条评论
	for i := 0; i < 5; i++ {
		comment := &social.Comment{
			AuthorID:   "user123",
			TargetType: social.CommentTargetTypeBook,
			TargetID:   bookID,
			Content:    "测试评论",
			Rating:     5,
			State:      social.CommentStateNormal,
		}
		err := repo.Create(ctx, comment)
		require.NoError(t, err)
	}

	// Act
	comments, total, err := repo.GetCommentsByBookID(ctx, bookID, 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, comments)
	assert.GreaterOrEqual(t, total, int64(5))
	assert.GreaterOrEqual(t, len(comments), 5)
}

// TestCommentRepository_GetRepliesByCommentID 测试获取回复列表
func TestCommentRepository_GetRepliesByCommentID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCommentRepo(t)
	defer cleanup()

	// 创建父评论
	parentComment := &social.Comment{
		AuthorID:   "user123",
		TargetType: social.CommentTargetTypeBook,
		TargetID:   "book123",
		Content:    "父评论",
		State:      social.CommentStateNormal,
	}
	err := repo.Create(ctx, parentComment)
	require.NoError(t, err)

	// 创建回复
	reply := &social.Comment{
		AuthorID:   "user456",
		TargetType: social.CommentTargetTypeBook,
		TargetID:   "book123",
		Content:    "回复内容",
		State:      social.CommentStateNormal,
	}
	// 设置ParentID以建立回复关系
	parentIDStr := parentComment.ID.Hex()
	reply.ParentID = &parentIDStr
	err = repo.Create(ctx, reply)
	require.NoError(t, err)

	// Act
	replies, err := repo.GetRepliesByCommentID(ctx, parentComment.ID.Hex())

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, replies)
	assert.Greater(t, len(replies), 0)
}

// TestCommentRepository_UpdateCommentStatus 测试更新评论状态
func TestCommentRepository_UpdateCommentStatus(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCommentRepo(t)
	defer cleanup()

	comment := &social.Comment{
		AuthorID:   "user123",
		TargetType: social.CommentTargetTypeBook,
		TargetID:   "book123",
		Content:    "待审核评论",
		State:      social.CommentStateNormal,
	}
	err := repo.Create(ctx, comment)
	require.NoError(t, err)

	// Act - 更新审核状态
	err = repo.UpdateCommentStatus(ctx, comment.ID.Hex(), string(social.CommentStateRejected), "包含敏感词")

	// Assert
	require.NoError(t, err)

	// 验证状态
	found, err := repo.GetByID(ctx, comment.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, social.CommentStateRejected, found.State)
	assert.Equal(t, "包含敏感词", found.RejectReason)
}

// TestCommentRepository_IncrementLikeCount 测试增加点赞数
func TestCommentRepository_IncrementLikeCount(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCommentRepo(t)
	defer cleanup()

	comment := &social.Comment{
		AuthorID:   "user123",
		TargetType: social.CommentTargetTypeBook,
		TargetID:   "book123",
		Content:    "点赞测试",
		State:      social.CommentStateNormal,
	}
	err := repo.Create(ctx, comment)
	require.NoError(t, err)

	// Act - 增加点赞数
	err = repo.IncrementLikeCount(ctx, comment.ID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证点赞数
	found, err := repo.GetByID(ctx, comment.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, int64(1), found.LikeCount)
}

// TestCommentRepository_DecrementLikeCount 测试减少点赞数
func TestCommentRepository_DecrementLikeCount(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCommentRepo(t)
	defer cleanup()

	comment := &social.Comment{
		AuthorID:   "user123",
		TargetType: social.CommentTargetTypeBook,
		TargetID:   "book123",
		Content:    "取消点赞测试",
		State:      social.CommentStateNormal,
	}
	err := repo.Create(ctx, comment)
	require.NoError(t, err)

	// Act - 先增加几次点赞
	for i := 0; i < 5; i++ {
		err = repo.IncrementLikeCount(ctx, comment.ID.Hex())
		require.NoError(t, err)
	}

	// Act - 减少点赞数
	err = repo.DecrementLikeCount(ctx, comment.ID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证点赞数
	found, err := repo.GetByID(ctx, comment.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, int64(4), found.LikeCount)
}

// TestCommentRepository_Health 测试健康检查
func TestCommentRepository_Health(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCommentRepo(t)
	defer cleanup()

	// Act
	err := repo.Health(ctx)

	// Assert
	assert.NoError(t, err)
}
