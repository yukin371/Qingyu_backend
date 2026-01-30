package mongodb_test

import (
	"context"
	"testing"

	"Qingyu_backend/models/bookstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mongodb "Qingyu_backend/repository/mongodb/bookstore"
	"Qingyu_backend/test/testutil"
)

// TestMongoBookDataStatisticsRepository_CountByStatus 测试统计指定状态的书籍数量
func TestMongoBookDataStatisticsRepository_CountByStatus(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookDataStatisticsRepository(db.Client(), db.Name())
	mutationRepo := mongodb.NewMongoBookDataMutationRepository(db.Client(), db.Name())
	ctx := context.Background()

	// 创建多个测试书籍
	for i := 0; i < 5; i++ {
		book := &bookstore.Book{
			Title:  "测试书籍",
			Status: bookstore.BookStatusOngoing,
		}
		mutationRepo.Create(ctx, book)
	}

	// Act
	count, err := repo.CountByStatus(ctx, bookstore.BookStatusOngoing)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(5))
}

// TestMongoBookDataStatisticsRepository_IncrementViewCount 测试增加浏览次数
func TestMongoBookDataStatisticsRepository_IncrementViewCount(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookDataStatisticsRepository(db.Client(), db.Name())
	mutationRepo := mongodb.NewMongoBookDataMutationRepository(db.Client(), db.Name())
	ctx := context.Background()

	// 创建测试书籍
	book := &bookstore.Book{
		Title:     "热门书籍",
		Author:    "测试作者",
		Status:    bookstore.BookStatusOngoing,
		ViewCount: 100,
	}
	mutationRepo.Create(ctx, book)

	// Act
	err := repo.IncrementViewCount(ctx, book.ID.Hex())

	// Assert
	require.NoError(t, err)
}

// TestMongoBookDataStatisticsRepository_GetStats 测试获取统计信息
func TestMongoBookDataStatisticsRepository_GetStats(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookDataStatisticsRepository(db.Client(), db.Name())
	ctx := context.Background()

	// Act
	stats, err := repo.GetStats(ctx)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, stats)
	assert.GreaterOrEqual(t, stats.TotalBooks, int64(0))
}

// TestMongoBookDataStatisticsRepository_IncrementLikeCount 测试增加点赞数
func TestMongoBookDataStatisticsRepository_IncrementLikeCount(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookDataStatisticsRepository(db.Client(), db.Name())
	mutationRepo := mongodb.NewMongoBookDataMutationRepository(db.Client(), db.Name())
	ctx := context.Background()

	// 创建测试书籍
	book := &bookstore.Book{
		Title:  "测试书籍",
		Author: "测试作者",
		Status: bookstore.BookStatusOngoing,
	}
	mutationRepo.Create(ctx, book)

	// Act
	err := repo.IncrementLikeCount(ctx, book.ID.Hex())

	// Assert
	require.NoError(t, err)
}
