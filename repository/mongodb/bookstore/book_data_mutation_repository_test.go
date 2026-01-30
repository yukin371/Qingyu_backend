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

// TestMongoBookDataMutationRepository_Create 测试创建书籍
func TestMongoBookDataMutationRepository_Create(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookDataMutationRepository(db.Client(), db.Name())
	ctx := context.Background()

	book := &bookstore.Book{
		Title:         "测试书籍",
		Author:        "测试作者",
		Introduction:  "这是一本测试书籍",
		Status:        bookstore.BookStatusOngoing,
		WordCount:     100000,
		ChapterCount:  100,
		Price:         19.99,
		IsFree:        false,
		IsRecommended: false,
		IsFeatured:    false,
	}

	// Act
	err := repo.Create(ctx, book)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, book.ID)
}

// TestMongoBookDataMutationRepository_Update 测试更新书籍
func TestMongoBookDataMutationRepository_Update(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookDataMutationRepository(db.Client(), db.Name())
	ctx := context.Background()

	// 创建测试书籍
	book := &bookstore.Book{
		Title:        "原始标题",
		Author:       "测试作者",
		Introduction: "原始简介",
		Status:       bookstore.BookStatusDraft,
	}
	err := repo.Create(ctx, book)
	require.NoError(t, err)

	// Act - 更新书籍
	updates := map[string]interface{}{
		"title":        "更新后的标题",
		"introduction": "更新后的简介",
		"status":       bookstore.BookStatusOngoing,
	}
	err = repo.Update(ctx, book.ID.Hex(), updates)

	// Assert
	require.NoError(t, err)
}

// TestMongoBookDataMutationRepository_Delete 测试删除书籍
func TestMongoBookDataMutationRepository_Delete(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookDataMutationRepository(db.Client(), db.Name())
	ctx := context.Background()

	// 创建测试书籍
	book := &bookstore.Book{
		Title:  "待删除书籍",
		Author: "测试作者",
		Status: bookstore.BookStatusDraft,
	}
	err := repo.Create(ctx, book)
	require.NoError(t, err)

	// Act - 删除书籍
	err = repo.Delete(ctx, book.ID.Hex())

	// Assert
	require.NoError(t, err)
}
