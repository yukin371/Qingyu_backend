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

// TestMongoBookManagementRepository_BatchUpdateStatus 测试批量更新书籍状态
func TestMongoBookManagementRepository_BatchUpdateStatus(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookManagementRepository(db.Client(), db.Name())
	mutationRepo := mongodb.NewMongoBookDataMutationRepository(db.Client(), db.Name())
	ctx := context.Background()

	// 创建测试书籍
	book1 := &bookstore.Book{Title: "书籍1", Status: bookstore.BookStatusDraft}
	book2 := &bookstore.Book{Title: "书籍2", Status: bookstore.BookStatusDraft}
	book3 := &bookstore.Book{Title: "书籍3", Status: bookstore.BookStatusDraft}

	mutationRepo.Create(ctx, book1)
	mutationRepo.Create(ctx, book2)
	mutationRepo.Create(ctx, book3)

	// Act - 批量更新状态
	bookIDs := []string{book1.ID.Hex(), book2.ID.Hex(), book3.ID.Hex()}
	err := repo.BatchUpdateStatus(ctx, bookIDs, bookstore.BookStatusOngoing)

	// Assert
	require.NoError(t, err)
}

// TestMongoBookManagementRepository_BatchUpdateRecommended 测试批量更新推荐状态
func TestMongoBookManagementRepository_BatchUpdateRecommended(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookManagementRepository(db.Client(), db.Name())
	mutationRepo := mongodb.NewMongoBookDataMutationRepository(db.Client(), db.Name())
	ctx := context.Background()

	// 创建测试书籍
	book1 := &bookstore.Book{
		Title:         "书籍1",
		Status:        bookstore.BookStatusOngoing,
		IsRecommended: false,
	}
	book2 := &bookstore.Book{
		Title:         "书籍2",
		Status:        bookstore.BookStatusOngoing,
		IsRecommended: false,
	}
	mutationRepo.Create(ctx, book1)
	mutationRepo.Create(ctx, book2)

	// Act - 批量设置为推荐
	bookIDs := []string{book1.ID.Hex(), book2.ID.Hex()}
	err := repo.BatchUpdateRecommended(ctx, bookIDs, true)

	// Assert
	require.NoError(t, err)
}

// TestMongoBookManagementRepository_GetYears 测试获取年份列表
func TestMongoBookManagementRepository_GetYears(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookManagementRepository(db.Client(), db.Name())
	ctx := context.Background()

	// Act
	years, err := repo.GetYears(ctx)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, years)
}

// TestMongoBookManagementRepository_GetTags 测试获取标签列表
func TestMongoBookManagementRepository_GetTags(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookManagementRepository(db.Client(), db.Name())
	ctx := context.Background()

	// Act
	tags, err := repo.GetTags(ctx, nil)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, tags)
}
