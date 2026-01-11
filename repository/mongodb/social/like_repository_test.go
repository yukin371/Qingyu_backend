package social_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/social"
	socialRepo "Qingyu_backend/repository/mongodb/social"
	"Qingyu_backend/test/testutil"
)

// setupLikeRepo 测试辅助函数
func setupLikeRepo(t *testing.T) (*socialRepo.MongoLikeRepository, context.Context, func()) {
	db, cleanup := testutil.SetupTestDB(t)
	repo := socialRepo.NewMongoLikeRepository(db)
	ctx := context.Background()
	return repo, ctx, cleanup
}

// TestLikeRepository_AddLike 测试添加点赞
func TestLikeRepository_AddLike(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupLikeRepo(t)
	defer cleanup()

	like := &social.Like{
		UserID:     "user_add_" + primitive.NewObjectID().Hex(),
		TargetType: social.LikeTargetTypeBook,
		TargetID:   "book_add_" + primitive.NewObjectID().Hex(),
	}

	// Act
	err := repo.AddLike(ctx, like)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, like.ID)
	assert.NotZero(t, like.CreatedAt)
}

// TestLikeRepository_AddLike_Duplicate 测试重复点赞
func TestLikeRepository_AddLike_Duplicate(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupLikeRepo(t)
	defer cleanup()

	userID := "user_dup_" + primitive.NewObjectID().Hex()
	bookID := "book_dup_" + primitive.NewObjectID().Hex()
	like := &social.Like{
		UserID:     userID,
		TargetType: social.LikeTargetTypeBook,
		TargetID:   bookID,
	}
	err := repo.AddLike(ctx, like)
	require.NoError(t, err)

	// Act - 尝试重复点赞
	err = repo.AddLike(ctx, &social.Like{
		UserID:     userID,
		TargetType: social.LikeTargetTypeBook,
		TargetID:   bookID,
	})

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "已经点赞")
}

// TestLikeRepository_RemoveLike 测试取消点赞
func TestLikeRepository_RemoveLike(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupLikeRepo(t)
	defer cleanup()

	userID := "user_remove_" + primitive.NewObjectID().Hex()
	bookID := "book_remove_" + primitive.NewObjectID().Hex()
	like := &social.Like{
		UserID:     userID,
		TargetType: social.LikeTargetTypeBook,
		TargetID:   bookID,
	}
	err := repo.AddLike(ctx, like)
	require.NoError(t, err)

	// Act
	err = repo.RemoveLike(ctx, userID, social.LikeTargetTypeBook, bookID)

	// Assert
	require.NoError(t, err)

	// 验证已取消
	isLiked, err := repo.IsLiked(ctx, userID, social.LikeTargetTypeBook, bookID)
	require.NoError(t, err)
	assert.False(t, isLiked)
}

// TestLikeRepository_RemoveLike_NotFound 测试取消不存在的点赞
func TestLikeRepository_RemoveLike_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupLikeRepo(t)
	defer cleanup()

	// Act
	err := repo.RemoveLike(ctx, "user123", social.LikeTargetTypeBook, "nonexistent_book")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "不存在")
}

// TestLikeRepository_IsLiked 测试检查是否点赞
func TestLikeRepository_IsLiked(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupLikeRepo(t)
	defer cleanup()

	userID := "user_isliked_" + primitive.NewObjectID().Hex()
	bookID := "book_isliked_" + primitive.NewObjectID().Hex()
	like := &social.Like{
		UserID:     userID,
		TargetType: social.LikeTargetTypeBook,
		TargetID:   bookID,
	}
	err := repo.AddLike(ctx, like)
	require.NoError(t, err)

	// Act
	isLiked, err := repo.IsLiked(ctx, userID, social.LikeTargetTypeBook, bookID)

	// Assert
	require.NoError(t, err)
	assert.True(t, isLiked)
}

// TestLikeRepository_IsLiked_False 测试检查未点赞状态
func TestLikeRepository_IsLiked_False(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupLikeRepo(t)
	defer cleanup()

	// Act
	isLiked, err := repo.IsLiked(ctx, "user123", social.LikeTargetTypeBook, "nonexistent_book")

	// Assert
	require.NoError(t, err)
	assert.False(t, isLiked)
}

// TestLikeRepository_GetByID 测试根据ID获取点赞
func TestLikeRepository_GetByID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupLikeRepo(t)
	defer cleanup()

	userID := "user_getbyid_" + primitive.NewObjectID().Hex()
	bookID := "book_getbyid_" + primitive.NewObjectID().Hex()
	like := &social.Like{
		UserID:     userID,
		TargetType: social.LikeTargetTypeBook,
		TargetID:   bookID,
	}
	err := repo.AddLike(ctx, like)
	require.NoError(t, err)

	// Act
	found, err := repo.GetByID(ctx, like.ID.Hex())

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, like.UserID, found.UserID)
	assert.Equal(t, like.TargetID, found.TargetID)
}

// TestLikeRepository_GetByID_NotFound 测试获取不存在的点赞
func TestLikeRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupLikeRepo(t)
	defer cleanup()

	fakeID := primitive.NewObjectID().Hex()

	// Act
	found, err := repo.GetByID(ctx, fakeID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, found)
	assert.Contains(t, err.Error(), "not found")
}

// TestLikeRepository_GetLikeCount 测试获取点赞数
func TestLikeRepository_GetLikeCount(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupLikeRepo(t)
	defer cleanup()

	bookID := "book_count_" + primitive.NewObjectID().Hex()
	for i := 0; i < 5; i++ {
		like := &social.Like{
			UserID:     primitive.NewObjectID().Hex(),
			TargetType: social.LikeTargetTypeBook,
			TargetID:   bookID,
		}
		err := repo.AddLike(ctx, like)
		require.NoError(t, err)
	}

	// Act
	count, err := repo.GetLikeCount(ctx, social.LikeTargetTypeBook, bookID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(5), count)
}

// TestLikeRepository_GetUserLikes 测试获取用户点赞列表
func TestLikeRepository_GetUserLikes(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupLikeRepo(t)
	defer cleanup()

	userID := "user123"
	for i := 0; i < 5; i++ {
		like := &social.Like{
			UserID:     userID,
			TargetType: social.LikeTargetTypeBook,
			TargetID:   primitive.NewObjectID().Hex(),
		}
		err := repo.AddLike(ctx, like)
		require.NoError(t, err)
	}

	// Act
	likes, total, err := repo.GetUserLikes(ctx, userID, social.LikeTargetTypeBook, 1, 3)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, likes)
	assert.GreaterOrEqual(t, total, int64(5))
	assert.LessOrEqual(t, len(likes), 3)
}

// TestLikeRepository_GetLikesCountBatch 测试批量获取点赞数
func TestLikeRepository_GetLikesCountBatch(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupLikeRepo(t)
	defer cleanup()

	book1ID := "book_batch_1_" + primitive.NewObjectID().Hex()
	book2ID := "book_batch_2_" + primitive.NewObjectID().Hex()
	book3ID := "book_batch_3_" + primitive.NewObjectID().Hex()
	bookIDs := []string{book1ID, book2ID, book3ID}

	userID := "user_batch_" + primitive.NewObjectID().Hex()
	for _, bookID := range bookIDs[:2] {
		like := &social.Like{
			UserID:     userID,
			TargetType: social.LikeTargetTypeBook,
			TargetID:   bookID,
		}
		err := repo.AddLike(ctx, like)
		require.NoError(t, err)
	}

	// Act
	counts, err := repo.GetLikesCountBatch(ctx, social.LikeTargetTypeBook, bookIDs)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, len(bookIDs), len(counts))
	assert.Greater(t, counts[book1ID], int64(0))
	assert.Greater(t, counts[book2ID], int64(0))
	assert.Equal(t, int64(0), counts[book3ID])
}

// TestLikeRepository_GetUserLikeStatusBatch 测试批量检查点赞状态
func TestLikeRepository_GetUserLikeStatusBatch(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupLikeRepo(t)
	defer cleanup()

	userID := "user_status_batch_" + primitive.NewObjectID().Hex()
	book1ID := "book_status_1_" + primitive.NewObjectID().Hex()
	book2ID := "book_status_2_" + primitive.NewObjectID().Hex()
	book3ID := "book_status_3_" + primitive.NewObjectID().Hex()
	bookIDs := []string{book1ID, book2ID, book3ID}

	for _, bookID := range bookIDs[:2] {
		like := &social.Like{
			UserID:     userID,
			TargetType: social.LikeTargetTypeBook,
			TargetID:   bookID,
		}
		err := repo.AddLike(ctx, like)
		require.NoError(t, err)
	}

	// Act
	status, err := repo.GetUserLikeStatusBatch(ctx, userID, social.LikeTargetTypeBook, bookIDs)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, len(bookIDs), len(status))
	assert.True(t, status[book1ID])
	assert.True(t, status[book2ID])
	assert.False(t, status[book3ID])
}

// TestLikeRepository_CountUserLikes 测试统计用户点赞数
func TestLikeRepository_CountUserLikes(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupLikeRepo(t)
	defer cleanup()

	userID := "user123"
	for i := 0; i < 5; i++ {
		like := &social.Like{
			UserID:     userID,
			TargetType: social.LikeTargetTypeBook,
			TargetID:   primitive.NewObjectID().Hex(),
		}
		err := repo.AddLike(ctx, like)
		require.NoError(t, err)
	}

	// Act
	count, err := repo.CountUserLikes(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(5))
}

// TestLikeRepository_CountTargetLikes 测试统计目标点赞数
func TestLikeRepository_CountTargetLikes(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupLikeRepo(t)
	defer cleanup()

	targetID := "book_count_target_" + primitive.NewObjectID().Hex()
	for i := 0; i < 3; i++ {
		like := &social.Like{
			UserID:     primitive.NewObjectID().Hex(),
			TargetType: social.LikeTargetTypeBook,
			TargetID:   targetID,
		}
		err := repo.AddLike(ctx, like)
		require.NoError(t, err)
	}

	// Act
	count, err := repo.CountTargetLikes(ctx, social.LikeTargetTypeBook, targetID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

// TestLikeRepository_CommentLike 测试评论点赞
func TestLikeRepository_CommentLike(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupLikeRepo(t)
	defer cleanup()

	commentID := "comment_like_" + primitive.NewObjectID().Hex()
	userID := "user_comment_like_" + primitive.NewObjectID().Hex()
	like := &social.Like{
		UserID:     userID,
		TargetType: social.LikeTargetTypeComment,
		TargetID:   commentID,
	}

	// Act - 点赞评论
	err := repo.AddLike(ctx, like)

	// Assert
	require.NoError(t, err)

	// 验证点赞状态
	isLiked, err := repo.IsLiked(ctx, userID, social.LikeTargetTypeComment, commentID)
	require.NoError(t, err)
	assert.True(t, isLiked)
}

// TestLikeRepository_CommentUnlike 测试取消评论点赞
func TestLikeRepository_CommentUnlike(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupLikeRepo(t)
	defer cleanup()

	commentID := "comment_unlike_" + primitive.NewObjectID().Hex()
	userID := "user_unlike_" + primitive.NewObjectID().Hex()
	like := &social.Like{
		UserID:     userID,
		TargetType: social.LikeTargetTypeComment,
		TargetID:   commentID,
	}
	err := repo.AddLike(ctx, like)
	require.NoError(t, err)

	// Act - 取消点赞
	err = repo.RemoveLike(ctx, userID, social.LikeTargetTypeComment, commentID)

	// Assert
	require.NoError(t, err)

	// 验证已取消
	isLiked, err := repo.IsLiked(ctx, userID, social.LikeTargetTypeComment, commentID)
	require.NoError(t, err)
	assert.False(t, isLiked)
}

// TestLikeRepository_ChapterLike 测试章节点赞
func TestLikeRepository_ChapterLike(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupLikeRepo(t)
	defer cleanup()

	chapterID := "chapter_like_" + primitive.NewObjectID().Hex()
	userID := "user_chapter_like_" + primitive.NewObjectID().Hex()
	like := &social.Like{
		UserID:     userID,
		TargetType: social.LikeTargetTypeChapter,
		TargetID:   chapterID,
	}

	// Act
	err := repo.AddLike(ctx, like)

	// Assert
	require.NoError(t, err)

	// 验证点赞状态
	isLiked, err := repo.IsLiked(ctx, userID, social.LikeTargetTypeChapter, chapterID)
	require.NoError(t, err)
	assert.True(t, isLiked)
}

// TestLikeRepository_Health 测试健康检查
func TestLikeRepository_Health(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupLikeRepo(t)
	defer cleanup()

	// Act
	err := repo.Health(ctx)

	// Assert
	assert.NoError(t, err)
}

// TestLikeRepository_AddLike_Validation 测试参数验证
func TestLikeRepository_AddLike_Validation(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupLikeRepo(t)
	defer cleanup()

	tests := []struct {
		name    string
		like    *social.Like
		wantErr string
	}{
		{
			name: "缺少UserID",
			like: &social.Like{
				TargetType: social.LikeTargetTypeBook,
				TargetID:   "book123",
			},
			wantErr: "参数不完整",
		},
		{
			name: "缺少TargetType",
			like: &social.Like{
				UserID:   "user123",
				TargetID: "book123",
			},
			wantErr: "参数不完整",
		},
		{
			name: "缺少TargetID",
			like: &social.Like{
				UserID:     "user123",
				TargetType: social.LikeTargetTypeBook,
			},
			wantErr: "参数不完整",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := repo.AddLike(ctx, tt.like)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

// TestLikeRepository_AddLike_AutoSetTimestamp 测试自动设置时间戳
func TestLikeRepository_AddLike_AutoSetTimestamp(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupLikeRepo(t)
	defer cleanup()

	like := &social.Like{
		UserID:     "user_timestamp_" + primitive.NewObjectID().Hex(),
		TargetType: social.LikeTargetTypeBook,
		TargetID:   "book_timestamp_" + primitive.NewObjectID().Hex(),
	}
	beforeCreate := time.Now()

	// Act
	err := repo.AddLike(ctx, like)

	// Assert
	require.NoError(t, err)
	assert.False(t, like.CreatedAt.IsZero())
	assert.WithinDuration(t, time.Now(), like.CreatedAt, time.Second)
	assert.True(t, like.CreatedAt.After(beforeCreate) || like.CreatedAt.Equal(beforeCreate))
}
