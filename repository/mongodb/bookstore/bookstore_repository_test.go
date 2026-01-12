package mongodb_test

import (
	"Qingyu_backend/models/bookstore"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/repository/mongodb/bookstore"
	"Qingyu_backend/test/testutil"
)

// TestMongoBookRepository_Create 测试创建书籍
func TestMongoBookRepository_Create(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookRepository(db.Client(), db.Name())
	ctx := context.Background()

	book := &bookstore.Book{
		Title:         "测试书籍",
		Author:        "测试作者",
		Introduction:  "这是一本测试书籍",
		Cover:         "https://example.com/cover.jpg",
		Status:        bookstore.BookStatusPublished,
		WordCount:     100000,
		ChapterCount:  100,
		Price:         19.99,
		IsFree:        false,
		IsRecommended: false,
		IsFeatured:    false,
		IsHot:         false,
	}

	// Act
	err := repo.Create(ctx, book)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, book.ID)
	assert.NotZero(t, book.CreatedAt)
	assert.NotZero(t, book.UpdatedAt)
}

// TestMongoBookRepository_GetByID 测试根据ID获取书籍
func TestMongoBookRepository_GetByID(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookRepository(db.Client(), db.Name())
	ctx := context.Background()

	// 创建测试书籍
	book := &bookstore.Book{
		Title:        "测试书籍",
		Author:       "测试作者",
		Introduction: "测试简介",
		Status:       bookstore.BookStatusPublished,
	}
	err := repo.Create(ctx, book)
	require.NoError(t, err)

	// Act
	result, err := repo.GetByID(ctx, book.ID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, book.Title, result.Title)
	assert.Equal(t, book.Author, result.Author)
}

// TestMongoBookRepository_GetByID_NotFound 测试获取不存在的书籍
func TestMongoBookRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookRepository(db.Client(), db.Name())
	ctx := context.Background()

	fakeID := primitive.NewObjectID()

	// Act
	result, err := repo.GetByID(ctx, fakeID)

	// Assert
	require.NoError(t, err)
	assert.Nil(t, result)
}

// TestMongoBookRepository_Update 测试更新书籍
func TestMongoBookRepository_Update(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookRepository(db.Client(), db.Name())
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
		"status":       bookstore.BookStatusPublished,
	}
	err = repo.Update(ctx, book.ID, updates)

	// Assert
	require.NoError(t, err)

	// 验证更新
	result, err := repo.GetByID(ctx, book.ID)
	require.NoError(t, err)
	assert.Equal(t, "更新后的标题", result.Title)
	assert.Equal(t, "更新后的简介", result.Introduction)
}

// TestMongoBookRepository_Delete 测试删除书籍
func TestMongoBookRepository_Delete(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookRepository(db.Client(), db.Name())
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
	err = repo.Delete(ctx, book.ID)

	// Assert
	require.NoError(t, err)

	// 验证已删除
	result, err := repo.GetByID(ctx, book.ID)
	require.NoError(t, err)
	assert.Nil(t, result)
}

// TestMongoBookRepository_List 测试列出书籍
func TestMongoBookRepository_List(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookRepository(db.Client(), db.Name())
	ctx := context.Background()

	// 创建多个测试书籍
	books := []*bookstore.Book{
		{Title: "书籍1", Author: "作者1", Status: bookstore.BookStatusPublished},
		{Title: "书籍2", Author: "作者2", Status: bookstore.BookStatusPublished},
		{Title: "书籍3", Author: "作者3", Status: bookstore.BookStatusDraft},
	}

	for _, book := range books {
		err := repo.Create(ctx, book)
		require.NoError(t, err)
	}

	// Act - 获取所有已发布书籍
	result, err := repo.GetByStatus(ctx, bookstore.BookStatusPublished, 10, 0)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result), 2)
}

// TestMongoBookRepository_GetByStatus 测试根据状态获取书籍
func TestMongoBookRepository_GetByStatus(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookRepository(db.Client(), db.Name())
	ctx := context.Background()

	// 创建测试书籍
	book1 := &bookstore.Book{Title: "已发布书籍", Status: bookstore.BookStatusPublished}
	book2 := &bookstore.Book{Title: "草稿书籍", Status: bookstore.BookStatusDraft}
	repo.Create(ctx, book1)
	repo.Create(ctx, book2)

	// Act
	result, err := repo.GetByStatus(ctx, bookstore.BookStatusPublished, 10, 0)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	// 验证返回的都是已发布状态
	for _, book := range result {
		assert.Equal(t, bookstore.BookStatusPublished, book.Status)
	}
}

// TestMongoBookRepository_Search 测试搜索书籍
func TestMongoBookRepository_Search(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookRepository(db.Client(), db.Name())
	ctx := context.Background()

	// 创建测试书籍
	book := &bookstore.Book{
		Title:        "玄幻小说测试",
		Author:       "知名作家",
		Introduction: "这是一部精彩的玄幻小说",
		Status:       bookstore.BookStatusPublished,
	}
	repo.Create(ctx, book)

	// Act
	result, err := repo.Search(ctx, "玄幻", 10, 0)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	// 验证搜索结果包含关键词
	found := false
	for _, b := range result {
		if b.ID == book.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "搜索结果应包含创建的测试书籍")
}

// TestMongoBookRepository_IncrementViewCount 测试增加浏览次数
func TestMongoBookRepository_IncrementViewCount(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookRepository(db.Client(), db.Name())
	ctx := context.Background()

	// 创建测试书籍
	book := &bookstore.Book{
		Title:     "热门书籍",
		Author:    "测试作者",
		Status:    bookstore.BookStatusPublished,
		ViewCount: 100,
	}
	repo.Create(ctx, book)

	// Act
	err := repo.IncrementViewCount(ctx, book.ID)

	// Assert
	require.NoError(t, err)

	// 验证浏览次数增加
	result, err := repo.GetByID(ctx, book.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, result.ViewCount, int64(100))
}

// TestMongoBookRepository_CountByStatus 测试统计指定状态的书籍数量
func TestMongoBookRepository_CountByStatus(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookRepository(db.Client(), db.Name())
	ctx := context.Background()

	// 创建多个测试书籍
	for i := 0; i < 5; i++ {
		book := &bookstore.Book{
			Title:  "测试书籍",
			Status: bookstore.BookStatusPublished,
		}
		repo.Create(ctx, book)
	}

	// Act
	count, err := repo.CountByStatus(ctx, bookstore.BookStatusPublished)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(5))
}

// TestMongoBookRepository_BatchUpdateStatus 测试批量更新书籍状态
func TestMongoBookRepository_BatchUpdateStatus(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookRepository(db.Client(), db.Name())
	ctx := context.Background()

	// 创建测试书籍
	book1 := &bookstore.Book{Title: "书籍1", Status: bookstore.BookStatusDraft}
	book2 := &bookstore.Book{Title: "书籍2", Status: bookstore.BookStatusDraft}
	book3 := &bookstore.Book{Title: "书籍3", Status: bookstore.BookStatusDraft}

	repo.Create(ctx, book1)
	repo.Create(ctx, book2)
	repo.Create(ctx, book3)

	// Act - 批量更新状态
	bookIDs := []primitive.ObjectID{book1.ID, book2.ID, book3.ID}
	err := repo.BatchUpdateStatus(ctx, bookIDs, bookstore.BookStatusPublished)

	// Assert
	require.NoError(t, err)

	// 验证所有书籍都已更新
	for _, bookID := range bookIDs {
		book, _ := repo.GetByID(ctx, bookID)
		assert.Equal(t, bookstore.BookStatusPublished, book.Status)
	}
}

// TestMongoBookRepository_BatchUpdateRecommended 测试批量更新推荐状态
func TestMongoBookRepository_BatchUpdateRecommended(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookRepository(db.Client(), db.Name())
	ctx := context.Background()

	// 创建测试书籍
	book1 := &bookstore.Book{
		Title:         "书籍1",
		Status:        bookstore.BookStatusPublished,
		IsRecommended: false,
	}
	book2 := &bookstore.Book{
		Title:         "书籍2",
		Status:        bookstore.BookStatusPublished,
		IsRecommended: false,
	}
	repo.Create(ctx, book1)
	repo.Create(ctx, book2)

	// Act - 批量设置为推荐
	bookIDs := []primitive.ObjectID{book1.ID, book2.ID}
	err := repo.BatchUpdateRecommended(ctx, bookIDs, true)

	// Assert
	require.NoError(t, err)

	// 验证推荐设置
	result1, _ := repo.GetByID(ctx, book1.ID)
	result2, _ := repo.GetByID(ctx, book2.ID)
	assert.True(t, result1.IsRecommended)
	assert.True(t, result2.IsRecommended)
}

// TestMongoBookRepository_Health 测试健康检查
func TestMongoBookRepository_Health(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoBookRepository(db.Client(), db.Name())
	ctx := context.Background()

	// Act
	err := repo.Health(ctx)

	// Assert
	require.NoError(t, err)
}
