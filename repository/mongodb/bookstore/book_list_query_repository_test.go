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

// TestMongoBookListQueryRepository_GetByID 测试根据ID获取书籍
func TestMongoBookListQueryRepository_GetByID(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookListQueryRepository(db.Client(), db.Name())
	mutationRepo := mongodb.NewMongoBookDataMutationRepository(db.Client(), db.Name())
	ctx := context.Background()

	// 创建测试书籍
	book := &bookstore.Book{
		Title:        "测试书籍",
		Author:       "测试作者",
		Introduction: "测试简介",
		Status:       bookstore.BookStatusOngoing,
	}
	err := mutationRepo.Create(ctx, book)
	require.NoError(t, err)

	// Act
	result, err := repo.GetByID(ctx, book.ID.Hex())

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, book.Title, result.Title)
	assert.Equal(t, book.Author, result.Author)
}

// TestMongoBookListQueryRepository_List 测试列出所有书籍
func TestMongoBookListQueryRepository_List(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookListQueryRepository(db.Client(), db.Name())
	mutationRepo := mongodb.NewMongoBookDataMutationRepository(db.Client(), db.Name())
	ctx := context.Background()

	// 创建测试书籍
	book := &bookstore.Book{
		Title:  "测试书籍",
		Status: bookstore.BookStatusOngoing,
	}
	mutationRepo.Create(ctx, book)

	// Act
	result, err := repo.List(ctx, nil)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result), 1)
}

// TestMongoBookListQueryRepository_Count 测试统计书籍数量
func TestMongoBookListQueryRepository_Count(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookListQueryRepository(db.Client(), db.Name())
	mutationRepo := mongodb.NewMongoBookDataMutationRepository(db.Client(), db.Name())
	ctx := context.Background()

	// 创建测试书籍
	book := &bookstore.Book{
		Title:  "测试书籍",
		Status: bookstore.BookStatusOngoing,
	}
	mutationRepo.Create(ctx, book)

	// Act
	count, err := repo.Count(ctx, nil)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(1))
}

// TestMongoBookListQueryRepository_Health 测试健康检查
func TestMongoBookListQueryRepository_Health(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookListQueryRepository(db.Client(), db.Name())
	ctx := context.Background()

	// Act
	err := repo.Health(ctx)

	// Assert
	require.NoError(t, err)
}
