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

// TestMongoBookSearchRepository_Search 测试搜索书籍
func TestMongoBookSearchRepository_Search(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookSearchRepository(db.Client(), db.Name())
	mutationRepo := mongodb.NewMongoBookDataMutationRepository(db.Client(), db.Name())
	ctx := context.Background()

	// 创建测试书籍
	book := &bookstore.Book{
		Title:        "玄幻小说测试",
		Author:       "知名作家",
		Introduction: "这是一部精彩的玄幻小说",
		Status:       bookstore.BookStatusOngoing,
	}
	mutationRepo.Create(ctx, book)

	// Act
	result, err := repo.Search(ctx, "玄幻", 10, 0)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
}

// TestMongoBookSearchRepository_SearchWithFilter 测试使用过滤器搜索书籍
func TestMongoBookSearchRepository_SearchWithFilter(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookSearchRepository(db.Client(), db.Name())
	mutationRepo := mongodb.NewMongoBookDataMutationRepository(db.Client(), db.Name())
	ctx := context.Background()

	// 创建测试书籍
	book := &bookstore.Book{
		Title:  "玄幻小说",
		Status: bookstore.BookStatusOngoing,
	}
	mutationRepo.Create(ctx, book)

	// Act
	filter := &bookstore.BookFilter{
		Limit:  10,
		Offset: 0,
	}
	result, err := repo.SearchWithFilter(ctx, filter)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
}

// TestMongoBookSearchRepository_GetByPriceRange 测试按价格区间获取书籍
func TestMongoBookSearchRepository_GetByPriceRange(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookSearchRepository(db.Client(), db.Name())
	mutationRepo := mongodb.NewMongoBookDataMutationRepository(db.Client(), db.Name())
	ctx := context.Background()

	// 创建测试书籍
	book := &bookstore.Book{
		Title:  "测试书籍",
		Price:  19.99,
		Status: bookstore.BookStatusOngoing,
	}
	mutationRepo.Create(ctx, book)

	// Act
	result, err := repo.GetByPriceRange(ctx, 0, 100, 10, 0)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
}
