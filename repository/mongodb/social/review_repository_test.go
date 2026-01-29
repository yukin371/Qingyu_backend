package social_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/models/social"
	socialIntf "Qingyu_backend/repository/interfaces/social"
	mongoSocial "Qingyu_backend/repository/mongodb/social"
	"Qingyu_backend/test/testutil"
)

// setupReviewRepo 测试辅助函数
func setupReviewRepo(t *testing.T) (socialIntf.ReviewRepository, context.Context, func()) {
	db, cleanup := testutil.SetupTestDB(t)
	repo := mongoSocial.NewMongoReviewRepository(db)
	ctx := context.Background()
	return repo, ctx, cleanup
}

// TestReviewRepository_CreateReview 测试创建书评
func TestReviewRepository_CreateReview(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupReviewRepo(t)
	defer cleanup()

	review := &social.Review{
		BookID:     "book123",
		UserID:     "user123",
		UserName:   "测试用户",
		UserAvatar: "avatar.jpg",
		Title:      "非常棒的书",
		Content:    "这本书真的很精彩，强烈推荐！",
		Rating:     5,
		IsSpoiler:  false,
		IsPublic:   true,
	}

	// Act
	err := repo.CreateReview(ctx, review)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, review.ID)
	assert.NotZero(t, review.CreatedAt)
	assert.NotZero(t, review.UpdatedAt)
}

// TestReviewRepository_GetReviewByID 测试根据ID获取书评
func TestReviewRepository_GetReviewByID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupReviewRepo(t)
	defer cleanup()

	review := &social.Review{
		BookID:     "book123",
		UserID:     "user123",
		UserName:   "测试用户",
		Title:      "测试书评标题",
		Content:    "测试书评内容",
		Rating:     4,
		IsPublic:   true,
	}
	err := repo.CreateReview(ctx, review)
	require.NoError(t, err)

	// Act
	found, err := repo.GetReviewByID(ctx, review.ID.Hex())

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, review.Title, found.Title)
	assert.Equal(t, review.Content, found.Content)
	assert.Equal(t, review.Rating, found.Rating)
	assert.Equal(t, review.BookID, found.BookID)
	assert.Equal(t, review.UserID, found.UserID)
}

// TestReviewRepository_GetReviewByID_NotFound 测试获取不存在的书评
func TestReviewRepository_GetReviewByID_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupReviewRepo(t)
	defer cleanup()

	// Act
	found, err := repo.GetReviewByID(ctx, "nonexistent_id")

	// Assert
	require.Error(t, err)
	assert.Nil(t, found)
	assert.Contains(t, err.Error(), "review not found")
}

// TestReviewRepository_GetReviewsByBook 测试获取书籍的书评列表
func TestReviewRepository_GetReviewsByBook(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupReviewRepo(t)
	defer cleanup()

	bookID := "test_book_123"

	// 创建多条书评
	for i := 0; i < 5; i++ {
		review := &social.Review{
			BookID:   bookID,
			UserID:   "user123",
			UserName: "测试用户",
			Title:    "书评标题",
			Content:  "书评内容",
			Rating:   5,
			IsPublic: true,
		}
		err := repo.CreateReview(ctx, review)
		require.NoError(t, err)
	}

	// Act
	reviews, total, err := repo.GetReviewsByBook(ctx, bookID, 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, reviews)
	assert.GreaterOrEqual(t, total, int64(5))
	assert.GreaterOrEqual(t, len(reviews), 5)
}

// TestReviewRepository_GetReviewsByUser 测试获取用户的书评列表
func TestReviewRepository_GetReviewsByUser(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupReviewRepo(t)
	defer cleanup()

	userID := "test_user_123"

	// 创建多条书评
	for i := 0; i < 3; i++ {
		review := &social.Review{
			BookID:   "book123",
			UserID:   userID,
			UserName: "测试用户",
			Title:    "书评标题",
			Content:  "书评内容",
			Rating:   4,
			IsPublic: true,
		}
		err := repo.CreateReview(ctx, review)
		require.NoError(t, err)
	}

	// Act
	reviews, total, err := repo.GetReviewsByUser(ctx, userID, 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, reviews)
	assert.GreaterOrEqual(t, total, int64(3))
	assert.GreaterOrEqual(t, len(reviews), 3)
}

// TestReviewRepository_GetPublicReviews 测试获取公开书评列表
func TestReviewRepository_GetPublicReviews(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupReviewRepo(t)
	defer cleanup()

	// 创建公开书评
	publicReview := &social.Review{
		BookID:   "book123",
		UserID:   "user123",
		UserName: "测试用户",
		Title:    "公开书评",
		Content:  "公开内容",
		Rating:   5,
		IsPublic: true,
	}
	err := repo.CreateReview(ctx, publicReview)
	require.NoError(t, err)

	// 创建私密书评
	privateReview := &social.Review{
		BookID:   "book456",
		UserID:   "user456",
		UserName: "测试用户2",
		Title:    "私密书评",
		Content:  "私密内容",
		Rating:   3,
		IsPublic: false,
	}
	err = repo.CreateReview(ctx, privateReview)
	require.NoError(t, err)

	// Act
	reviews, total, err := repo.GetPublicReviews(ctx, 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, reviews)
	// 应该只返回公开的书评
	assert.GreaterOrEqual(t, total, int64(1))
	for _, review := range reviews {
		assert.True(t, review.IsPublic)
	}
}

// TestReviewRepository_GetReviewsByRating 测试根据评分获取书评
func TestReviewRepository_GetReviewsByRating(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupReviewRepo(t)
	defer cleanup()

	bookID := "test_book_rating"

	// 创建不同评分的书评
	ratings := []int{5, 5, 4, 4, 3}
	for _, rating := range ratings {
		review := &social.Review{
			BookID:   bookID,
			UserID:   "user123",
			UserName: "测试用户",
			Title:    "书评标题",
			Content:  "书评内容",
			Rating:   rating,
			IsPublic: true,
		}
		err := repo.CreateReview(ctx, review)
		require.NoError(t, err)
	}

	// Act - 获取5星书评
	reviews, total, err := repo.GetReviewsByRating(ctx, bookID, 5, 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, reviews)
	assert.GreaterOrEqual(t, total, int64(2))
	for _, review := range reviews {
		assert.Equal(t, 5, review.Rating)
	}
}

// TestReviewRepository_UpdateReview 测试更新书评
func TestReviewRepository_UpdateReview(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupReviewRepo(t)
	defer cleanup()

	review := &social.Review{
		BookID:   "book123",
		UserID:   "user123",
		UserName: "测试用户",
		Title:    "原始标题",
		Content:  "原始内容",
		Rating:   3,
		IsPublic: true,
	}
	err := repo.CreateReview(ctx, review)
	require.NoError(t, err)

	// Act - 更新书评
	updates := map[string]interface{}{
		"title":   "更新后的标题",
		"content": "更新后的内容",
		"rating":  5,
	}
	err = repo.UpdateReview(ctx, review.ID.Hex(), updates)

	// Assert
	require.NoError(t, err)

	// 验证更新
	found, err := repo.GetReviewByID(ctx, review.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, "更新后的标题", found.Title)
	assert.Equal(t, "更新后的内容", found.Content)
	assert.Equal(t, 5, found.Rating)
}

// TestReviewRepository_DeleteReview 测试删除书评
func TestReviewRepository_DeleteReview(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupReviewRepo(t)
	defer cleanup()

	review := &social.Review{
		BookID:   "book123",
		UserID:   "user123",
		UserName: "测试用户",
		Title:    "待删除的书评",
		Content:  "待删除的内容",
		Rating:   4,
		IsPublic: true,
	}
	err := repo.CreateReview(ctx, review)
	require.NoError(t, err)

	// Act - 删除书评
	err = repo.DeleteReview(ctx, review.ID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证已删除
	found, err := repo.GetReviewByID(ctx, review.ID.Hex())
	require.Error(t, err)
	assert.Nil(t, found)
}

// TestReviewRepository_CreateReviewLike 测试创建书评点赞
func TestReviewRepository_CreateReviewLike(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupReviewRepo(t)
	defer cleanup()

	// 先创建书评
	review := &social.Review{
		BookID:   "book123",
		UserID:   "user123",
		UserName: "测试用户",
		Title:    "书评标题",
		Content:  "书评内容",
		Rating:   5,
		IsPublic: true,
	}
	err := repo.CreateReview(ctx, review)
	require.NoError(t, err)

	// 创建点赞
	reviewLike := &social.ReviewLike{
		ReviewID: review.ID.Hex(),
		UserID:   "user456",
	}

	// Act
	err = repo.CreateReviewLike(ctx, reviewLike)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, reviewLike.ID)
	assert.NotZero(t, reviewLike.CreatedAt)
}

// TestReviewRepository_DeleteReviewLike 测试删除书评点赞
func TestReviewRepository_DeleteReviewLike(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupReviewRepo(t)
	defer cleanup()

	// 先创建书评
	review := &social.Review{
		BookID:   "book123",
		UserID:   "user123",
		UserName: "测试用户",
		Title:    "书评标题",
		Content:  "书评内容",
		Rating:   5,
		IsPublic: true,
	}
	err := repo.CreateReview(ctx, review)
	require.NoError(t, err)

	// 创建点赞
	reviewLike := &social.ReviewLike{
		ReviewID: review.ID.Hex(),
		UserID:   "user456",
	}
	err = repo.CreateReviewLike(ctx, reviewLike)
	require.NoError(t, err)

	// Act - 删除点赞
	err = repo.DeleteReviewLike(ctx, review.ID.Hex(), "user456")

	// Assert
	require.NoError(t, err)

	// 验证已删除
	found, err := repo.GetReviewLike(ctx, review.ID.Hex(), "user456")
	require.Error(t, err)
	assert.Nil(t, found)
}

// TestReviewRepository_GetReviewLike 测试获取书评点赞记录
func TestReviewRepository_GetReviewLike(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupReviewRepo(t)
	defer cleanup()

	// 先创建书评
	review := &social.Review{
		BookID:   "book123",
		UserID:   "user123",
		UserName: "测试用户",
		Title:    "书评标题",
		Content:  "书评内容",
		Rating:   5,
		IsPublic: true,
	}
	err := repo.CreateReview(ctx, review)
	require.NoError(t, err)

	// 创建点赞
	reviewLike := &social.ReviewLike{
		ReviewID: review.ID.Hex(),
		UserID:   "user456",
	}
	err = repo.CreateReviewLike(ctx, reviewLike)
	require.NoError(t, err)

	// Act
	found, err := repo.GetReviewLike(ctx, review.ID.Hex(), "user456")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, review.ID.Hex(), found.ReviewID)
	assert.Equal(t, "user456", found.UserID)
}

// TestReviewRepository_IsReviewLiked 测试检查是否已点赞
func TestReviewRepository_IsReviewLiked(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupReviewRepo(t)
	defer cleanup()

	// 先创建书评
	review := &social.Review{
		BookID:   "book123",
		UserID:   "user123",
		UserName: "测试用户",
		Title:    "书评标题",
		Content:  "书评内容",
		Rating:   5,
		IsPublic: true,
	}
	err := repo.CreateReview(ctx, review)
	require.NoError(t, err)

	// 创建点赞
	reviewLike := &social.ReviewLike{
		ReviewID: review.ID.Hex(),
		UserID:   "user456",
	}
	err = repo.CreateReviewLike(ctx, reviewLike)
	require.NoError(t, err)

	// Act - 检查已点赞的用户
	liked, err := repo.IsReviewLiked(ctx, review.ID.Hex(), "user456")

	// Assert
	require.NoError(t, err)
	assert.True(t, liked)

	// Act - 检查未点赞的用户
	liked, err = repo.IsReviewLiked(ctx, review.ID.Hex(), "user789")

	// Assert
	require.NoError(t, err)
	assert.False(t, liked)
}

// TestReviewRepository_GetReviewLikes 测试获取书评点赞列表
func TestReviewRepository_GetReviewLikes(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupReviewRepo(t)
	defer cleanup()

	// 先创建书评
	review := &social.Review{
		BookID:   "book123",
		UserID:   "user123",
		UserName: "测试用户",
		Title:    "书评标题",
		Content:  "书评内容",
		Rating:   5,
		IsPublic: true,
	}
	err := repo.CreateReview(ctx, review)
	require.NoError(t, err)

	// 创建多个点赞
	for i := 0; i < 3; i++ {
		reviewLike := &social.ReviewLike{
			ReviewID: review.ID.Hex(),
			UserID:   "user" + string(rune(i)),
		}
		err = repo.CreateReviewLike(ctx, reviewLike)
		require.NoError(t, err)
	}

	// Act
	likes, total, err := repo.GetReviewLikes(ctx, review.ID.Hex(), 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, likes)
	assert.GreaterOrEqual(t, total, int64(3))
	assert.GreaterOrEqual(t, len(likes), 3)
}

// TestReviewRepository_IncrementReviewLikeCount 测试增加书评点赞数
func TestReviewRepository_IncrementReviewLikeCount(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupReviewRepo(t)
	defer cleanup()

	review := &social.Review{
		BookID:     "book123",
		UserID:     "user123",
		UserName:   "测试用户",
		Title:      "书评标题",
		Content:    "书评内容",
		Rating:     5,
		LikeCount:  0,
		IsPublic:   true,
	}
	err := repo.CreateReview(ctx, review)
	require.NoError(t, err)

	// Act - 增加点赞数
	err = repo.IncrementReviewLikeCount(ctx, review.ID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证点赞数
	found, err := repo.GetReviewByID(ctx, review.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, 1, found.LikeCount)
}

// TestReviewRepository_DecrementReviewLikeCount 测试减少书评点赞数
func TestReviewRepository_DecrementReviewLikeCount(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupReviewRepo(t)
	defer cleanup()

	review := &social.Review{
		BookID:     "book123",
		UserID:     "user123",
		UserName:   "测试用户",
		Title:      "书评标题",
		Content:    "书评内容",
		Rating:     5,
		LikeCount:  5,
		IsPublic:   true,
	}
	err := repo.CreateReview(ctx, review)
	require.NoError(t, err)

	// Act - 减少点赞数
	err = repo.DecrementReviewLikeCount(ctx, review.ID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证点赞数
	found, err := repo.GetReviewByID(ctx, review.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, 4, found.LikeCount)
}

// TestReviewRepository_GetAverageRating 测试获取书籍平均评分
func TestReviewRepository_GetAverageRating(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupReviewRepo(t)
	defer cleanup()

	bookID := "test_book_avg_rating"

	// 创建不同评分的书评
	ratings := []int{5, 4, 3, 5, 4} // 平均4.2
	for _, rating := range ratings {
		review := &social.Review{
			BookID:   bookID,
			UserID:   "user123",
			UserName: "测试用户",
			Title:    "书评标题",
			Content:  "书评内容",
			Rating:   rating,
			IsPublic: true,
		}
		err := repo.CreateReview(ctx, review)
		require.NoError(t, err)
	}

	// Act
	avgRating, err := repo.GetAverageRating(ctx, bookID)

	// Assert
	require.NoError(t, err)
	assert.InDelta(t, 4.2, avgRating, 0.1)
}

// TestReviewRepository_GetRatingDistribution 测试获取评分分布
func TestReviewRepository_GetRatingDistribution(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupReviewRepo(t)
	defer cleanup()

	bookID := "test_book_distribution"

	// 创建不同评分的书评
	ratings := []int{5, 5, 4, 4, 3}
	for _, rating := range ratings {
		review := &social.Review{
			BookID:   bookID,
			UserID:   "user123",
			UserName: "测试用户",
			Title:    "书评标题",
			Content:  "书评内容",
			Rating:   rating,
			IsPublic: true,
		}
		err := repo.CreateReview(ctx, review)
		require.NoError(t, err)
	}

	// Act
	distribution, err := repo.GetRatingDistribution(ctx, bookID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, distribution)
	assert.Equal(t, int64(2), distribution[5])
	assert.Equal(t, int64(2), distribution[4])
	assert.Equal(t, int64(1), distribution[3])
}

// TestReviewRepository_CountReviews 测试统计书评数
func TestReviewRepository_CountReviews(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupReviewRepo(t)
	defer cleanup()

	bookID := "test_book_count"

	// 创建多条书评
	for i := 0; i < 7; i++ {
		review := &social.Review{
			BookID:   bookID,
			UserID:   "user123",
			UserName: "测试用户",
			Title:    "书评标题",
			Content:  "书评内容",
			Rating:   5,
			IsPublic: true,
		}
		err := repo.CreateReview(ctx, review)
		require.NoError(t, err)
	}

	// Act
	count, err := repo.CountReviews(ctx, bookID)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(7))
}

// TestReviewRepository_CountUserReviews 测试统计用户书评数
func TestReviewRepository_CountUserReviews(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupReviewRepo(t)
	defer cleanup()

	userID := "test_user_count"

	// 创建多条书评
	for i := 0; i < 4; i++ {
		review := &social.Review{
			BookID:   "book123",
			UserID:   userID,
			UserName: "测试用户",
			Title:    "书评标题",
			Content:  "书评内容",
			Rating:   5,
			IsPublic: true,
		}
		err := repo.CreateReview(ctx, review)
		require.NoError(t, err)
	}

	// Act
	count, err := repo.CountUserReviews(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(4))
}

// TestReviewRepository_Health 测试健康检查
func TestReviewRepository_Health(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupReviewRepo(t)
	defer cleanup()

	// Act
	err := repo.Health(ctx)

	// Assert
	assert.NoError(t, err)
}
