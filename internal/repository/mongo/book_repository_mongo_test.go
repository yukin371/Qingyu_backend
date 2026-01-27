package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"qingteng-qa/internal/domain"
	"qingteng-qa/internal/repository"
	"qingteng-qa/pkg/errors"
)

// TestMongoBookRepository_Create_Success 测试成功创建书籍
func TestMongoBookRepository_Create_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := setupTestRepository(t)
	book := createTestBook(t)

	// Act
	err := repo.Create(ctx, book)

	// Assert
	require.NoError(t, err)
	assert.NotEqual(t, primitive.NilObjectID, book.ID, "ID应该被设置")
	assert.WithinDuration(t, time.Now(), book.CreatedAt, 2*time.Second, "创建时间应该被设置")
	assert.WithinDuration(t, time.Now(), book.UpdatedAt, 2*time.Second, "更新时间应该被设置")
}

// TestMongoBookRepository_Create_DuplicateISBN 测试ISBN重复
func TestMongoBookRepository_Create_DuplicateISBN(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := setupTestRepository(t)
	book1 := createTestBook(t)
	book1.ISBN = "9787111111111"
	_ = repo.Create(ctx, book1)

	book2 := createTestBook(t)
	book2.ISBN = "9787111111111" // 相同的ISBN
	book2.Title = "另一本书"

	// Act
	err := repo.Create(ctx, book2)

	// Assert
	require.Error(t, err)
	var bizErr *errors.BusinessError
	assert.True(t, errors.As(err, &bizErr), "应该返回业务错误")
	assert.Equal(t, errors.ErrCodeDuplicateKey, bizErr.Code, "错误码应该是重复键")
}

// TestMongoBookRepository_GetByID_Success 测试成功获取书籍
func TestMongoBookRepository_GetByID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := setupTestRepository(t)
	book := createTestBook(t)
	_ = repo.Create(ctx, book)
	bookID := book.ID.Hex()

	// Act
	foundBook, err := repo.GetByID(ctx, bookID)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, foundBook)
	assert.Equal(t, book.ID, foundBook.ID)
	assert.Equal(t, book.Title, foundBook.Title)
	assert.Equal(t, book.Author, foundBook.Author)
}

// TestMongoBookRepository_GetByID_NotFound 测试书籍不存在
func TestMongoBookRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := setupTestRepository(t)
	nonExistentID := primitive.NewObjectID().Hex()

	// Act
	book, err := repo.GetByID(ctx, nonExistentID)

	// Assert
	require.Error(t, err)
	require.Nil(t, book)
	var bizErr *errors.BusinessError
	assert.True(t, errors.As(err, &bizErr), "应该返回业务错误")
	assert.Equal(t, errors.ErrCodeBookNotFound, bizErr.Code, "错误码应该是书籍不存在")
}

// TestMongoBookRepository_GetByID_InvalidID 测试无效ID
func TestMongoBookRepository_GetByID_InvalidID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := setupTestRepository(t)
	invalidID := "invalid-id"

	// Act
	book, err := repo.GetByID(ctx, invalidID)

	// Assert
	require.Error(t, err)
	require.Nil(t, book)
	var bizErr *errors.BusinessError
	assert.True(t, errors.As(err, &bizErr), "应该返回业务错误")
	assert.Equal(t, errors.ErrCodeInvalidID, bizErr.Code, "错误码应该是无效ID")
}

// TestMongoBookRepository_Update_Success 测试成功更新书籍
func TestMongoBookRepository_Update_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := setupTestRepository(t)
	book := createTestBook(t)
	_ = repo.Create(ctx, book)

	// 修改书籍信息
	book.Title = "更新后的标题"
	book.Price = 9999
	book.Stock = 50

	// Act
	err := repo.Update(ctx, book)

	// Assert
	require.NoError(t, err)

	// 验证更新
	updatedBook, err := repo.GetByID(ctx, book.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, "更新后的标题", updatedBook.Title)
	assert.Equal(t, int64(9999), updatedBook.Price)
	assert.Equal(t, 50, updatedBook.Stock)
	assert.WithinDuration(t, time.Now(), updatedBook.UpdatedAt, 2*time.Second, "更新时间应该被更新")
}

// TestMongoBookRepository_Update_NotFound 测试更新不存在的书籍
func TestMongoBookRepository_Update_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := setupTestRepository(t)
	book := createTestBook(t)
	book.ID = primitive.NewObjectID() // 不存在的ID

	// Act
	err := repo.Update(ctx, book)

	// Assert
	require.Error(t, err)
	var bizErr *errors.BusinessError
	assert.True(t, errors.As(err, &bizErr), "应该返回业务错误")
	assert.Equal(t, errors.ErrCodeBookNotFound, bizErr.Code)
}

// TestMongoBookRepository_Delete_Success 测试成功删除书籍
func TestMongoBookRepository_Delete_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := setupTestRepository(t)
	book := createTestBook(t)
	_ = repo.Create(ctx, book)
	bookID := book.ID.Hex()

	// Act
	err := repo.Delete(ctx, bookID)

	// Assert
	require.NoError(t, err)

	// 验证软删除: 状态应该变为deleted
	deletedBook, err := repo.GetByID(ctx, bookID)
	require.NoError(t, err) // 软删除后仍能获取
	assert.Equal(t, domain.BookStatusDeleted, deletedBook.Status)
}

// TestMongoBookRepository_Delete_NotFound 测试删除不存在的书籍
func TestMongoBookRepository_Delete_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := setupTestRepository(t)
	nonExistentID := primitive.NewObjectID().Hex()

	// Act
	err := repo.Delete(ctx, nonExistentID)

	// Assert
	require.Error(t, err)
	var bizErr *errors.BusinessError
	assert.True(t, errors.As(err, &bizErr), "应该返回业务错误")
	assert.Equal(t, errors.ErrCodeBookNotFound, bizErr.Code)
}

// TestMongoBookRepository_List_WithoutFilter 测试无条件查询
func TestMongoBookRepository_List_WithoutFilter(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := setupTestRepository(t)

	// 创建测试数据
	book1 := createTestBook(t)
	book1.Status = domain.BookStatusOngoing
	_ = repo.Create(ctx, book1)

	book2 := createTestBook(t)
	book2.Status = domain.BookStatusCompleted
	_ = repo.Create(ctx, book2)

	book3 := createTestBook(t)
	book3.Status = domain.BookStatusDeleted
	_ = repo.Create(ctx, book3)

	// Act
	books, total, err := repo.List(ctx, nil, 1, 10)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(2), total, "总数应该为2 (不包含已删除)")
	assert.Len(t, books, 2, "应该返回2本书")
}

// TestMongoBookRepository_List_WithKeyword 测试关键词搜索
func TestMongoBookRepository_List_WithKeyword(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := setupTestRepository(t)

	book1 := createTestBook(t)
	book1.Title = "三体"
	_ = repo.Create(ctx, book1)

	book2 := createTestBook(t)
	book2.Title = "三体II"
	_ = repo.Create(ctx, book2)

	book3 := createTestBook(t)
	book3.Title = "平凡的世界"
	_ = repo.Create(ctx, book3)

	keyword := "三体"
	filter := &repository.BookFilter{
		Keyword: &keyword,
	}

	// Act
	books, total, err := repo.List(ctx, filter, 1, 10)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, books, 2)
	for _, book := range books {
		assert.Contains(t, book.Title, "三体")
	}
}

// TestMongoBookRepository_List_WithTags 测试标签过滤
func TestMongoBookRepository_List_WithTags(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := setupTestRepository(t)

	book1 := createTestBook(t)
	book1.Tags = []string{"玄幻", "修仙"}
	_ = repo.Create(ctx, book1)

	book2 := createTestBook(t)
	book2.Tags = []string{"仙侠", "修仙"}
	_ = repo.Create(ctx, book2)

	book3 := createTestBook(t)
	book3.Tags = []string{"科幻", "未来"}
	_ = repo.Create(ctx, book3)

	filter := &repository.BookFilter{
		Tags: []string{"玄幻", "仙侠"}, // ANY语义
	}

	// Act
	books, total, err := repo.List(ctx, filter, 1, 10)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(2), total, "应该匹配包含任一标签的书籍")
	assert.Len(t, books, 2)
}

// TestMongoBookRepository_UpdateStock_Success 测试成功更新库存
func TestMongoBookRepository_UpdateStock_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := setupTestRepository(t)
	book := createTestBook(t)
	book.Stock = 100
	_ = repo.Create(ctx, book)

	// Act
	err := repo.UpdateStock(ctx, book.ID.Hex(), -10)

	// Assert
	require.NoError(t, err)

	// 验证库存
	updatedBook, _ := repo.GetByID(ctx, book.ID.Hex())
	assert.Equal(t, 90, updatedBook.Stock, "库存应该减少10")
}

// TestMongoBookRepository_UpdateStock_OutOfStock 测试库存不足
func TestMongoBookRepository_UpdateStock_OutOfStock(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := setupTestRepository(t)
	book := createTestBook(t)
	book.Stock = 5
	_ = repo.Create(ctx, book)

	// Act
	err := repo.UpdateStock(ctx, book.ID.Hex(), -10)

	// Assert
	require.Error(t, err)
	var bizErr *errors.BusinessError
	assert.True(t, errors.As(err, &bizErr), "应该返回业务错误")
	assert.Equal(t, errors.ErrCodeBookOutOfStock, bizErr.Code, "错误码应该是库存不足")

	// 验证库存未被修改
	originalBook, _ := repo.GetByID(ctx, book.ID.Hex())
	assert.Equal(t, 5, originalBook.Stock, "库存应该保持不变")
}

// TestMongoBookRepository_Transaction_Success 测试事务成功
func TestMongoBookRepository_Transaction_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := setupTestRepository(t)
	book := createTestBook(t)
	_ = repo.Create(ctx, book)

	// 注意: 这里需要完整的事务管理器实现才能测试
	// 当前为占位测试,验证WithTransaction方法存在

	// Act & Assert
	txAware, ok := interface{}(repo).(interface{ WithTransaction(tx repository.Transaction) interface{} })
	assert.True(t, ok, "应该实现WithTransaction方法")
	assert.NotNil(t, txAware, "WithTransaction应该返回非nil值")
}

// TestMongoBookRepository_IncrementView 测试增加浏览次数
func TestMongoBookRepository_IncrementView(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := setupTestRepository(t)
	book := createTestBook(t)
	book.ViewCount = 100
	_ = repo.Create(ctx, book)

	// Act
	err := repo.IncrementView(ctx, book.ID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证浏览次数
	updatedBook, _ := repo.GetByID(ctx, book.ID.Hex())
	assert.Equal(t, int64(101), updatedBook.ViewCount, "浏览次数应该增加1")
}

// TestMongoBookRepository_IncrementLike 测试增加点赞数
func TestMongoBookRepository_IncrementLike(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := setupTestRepository(t)
	book := createTestBook(t)
	book.LikeCount = 50
	_ = repo.Create(ctx, book)

	// Act
	err := repo.IncrementLike(ctx, book.ID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证点赞数
	updatedBook, _ := repo.GetByID(ctx, book.ID.Hex())
	assert.Equal(t, int64(51), updatedBook.LikeCount, "点赞数应该增加1")
}

// TestMongoBookRepository_DecrementLike 测试减少点赞数
func TestMongoBookRepository_DecrementLike(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := setupTestRepository(t)
	book := createTestBook(t)
	book.LikeCount = 50
	_ = repo.Create(ctx, book)

	// Act
	err := repo.DecrementLike(ctx, book.ID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证点赞数
	updatedBook, _ := repo.GetByID(ctx, book.ID.Hex())
	assert.Equal(t, int64(49), updatedBook.LikeCount, "点赞数应该减少1")
}

// setupTestRepository 创建测试Repository
// 注意: 这里使用内存数据库或mock会更好,当前实现需要真实MongoDB连接
func setupTestRepository(t *testing.T) *MongoBookRepository {
	// TODO: 使用testcontainers或mock MongoDB
	// 当前使用占位实现
	t.Skip("需要MongoDB测试环境,跳过集成测试")
	return nil
}

// createTestBook 创建测试书籍数据
func createTestBook(t *testing.T) *domain.Book {
	now := time.Now()
	return &domain.Book{
		ID:          primitive.NewObjectID(),
		Title:       "测试书籍",
		Author:      "测试作者",
		ISBN:        generateTestISBN(t),
		Description: "这是一本测试书籍",
		CategoryID:  primitive.NewObjectIDPtr(),
		Tags:        []string{"测试", "样例"},
		CoverImage:  "https://example.com/cover.jpg",
		Price:       5999, // 59.99元
		Stock:       100,
		Status:      domain.BookStatusOngoing,
		ViewCount:   0,
		LikeCount:   0,
		CollectCount: 0,
		WordCount:   100000,
		PublishDate: &now,
		Publisher:   "测试出版社",
		PublishedAt: now,
		UpdatedAt:   now,
		CreatedAt:   now,
	}
}

// generateTestISBN 生成测试ISBN
func generateTestISBN(t *testing.T) string {
	// 生成唯一的测试ISBN
	return fmt.Sprintf("9787111%06d", time.Now().UnixNano()%1000000)
}
