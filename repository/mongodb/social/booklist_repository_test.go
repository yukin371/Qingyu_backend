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

// setupBookListRepo 测试辅助函数
func setupBookListRepo(t *testing.T) (*socialRepo.MongoBookListRepository, context.Context, func()) {
	db, cleanup := testutil.SetupTestDB(t)
	repo := socialRepo.NewMongoBookListRepository(db)
	ctx := context.Background()
	return repo, ctx, cleanup
}

// createTestBookList 创建测试书单的辅助函数
func createTestBookList(t *testing.T, repo *socialRepo.MongoBookListRepository, ctx context.Context, userID string) *social.BookList {
	bookList := &social.BookList{
		UserID:      userID,
		UserName:    "测试用户",
		UserAvatar:  "avatar.jpg",
		Title:       "我的测试书单",
		Description: "这是一个测试书单",
		Cover:       "cover.jpg",
		Books:       []social.BookListItem{},
		BookCount:   0,
		LikeCount:   0,
		ForkCount:   0,
		ViewCount:   0,
		IsPublic:    true,
		Tags:        []string{"小说", "科幻"},
		Category:    "科幻小说",
	}
	err := repo.CreateBookList(ctx, bookList)
	require.NoError(t, err)
	return bookList
}

// TestBookListRepository_CreateBookList 测试创建书单
func TestBookListRepository_CreateBookList(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()
	bookList := &social.BookList{
		UserID:      userID,
		UserName:    "测试用户",
		Title:       "我的第一本书单",
		Description: "这是我创建的书单",
		IsPublic:    true,
		Tags:        []string{"文学", "经典"},
		Category:    "文学",
		Books:       []social.BookListItem{},
	}

	// Act
	err := repo.CreateBookList(ctx, bookList)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, bookList.ID)
	assert.NotZero(t, bookList.CreatedAt)
	assert.NotZero(t, bookList.UpdatedAt)
	assert.Equal(t, userID, bookList.UserID)
}

// TestBookListRepository_GetBookListByID 测试根据ID获取书单
func TestBookListRepository_GetBookListByID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()
	created := createTestBookList(t, repo, ctx, userID)

	// Act
	found, err := repo.GetBookListByID(ctx, created.ID.Hex())

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, created.ID, found.ID)
	assert.Equal(t, created.Title, found.Title)
	assert.Equal(t, created.Description, found.Description)
	assert.Equal(t, created.IsPublic, found.IsPublic)
}

// TestBookListRepository_GetBookListByID_NotFound 测试获取不存在的书单
func TestBookListRepository_GetBookListByID_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	// Act
	found, err := repo.GetBookListByID(ctx, "nonexistent_id")

	// Assert
	require.Error(t, err)
	assert.Nil(t, found)
	assert.Contains(t, err.Error(), "无效的书单ID")
}

// TestBookListRepository_GetBookListByID_InvalidID 测试无效ID格式
func TestBookListRepository_GetBookListByID_InvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	// Act
	found, err := repo.GetBookListByID(ctx, "invalid_objectid")

	// Assert
	require.Error(t, err)
	assert.Nil(t, found)
	assert.Contains(t, err.Error(), "无效的书单ID")
}

// TestBookListRepository_GetBookListByID_NotExists 测试获取已存在但无效的ID
func TestBookListRepository_GetBookListByID_NotExists(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	validID := primitive.NewObjectID().Hex()

	// Act
	found, err := repo.GetBookListByID(ctx, validID)

	// Assert
	require.NoError(t, err)
	assert.Nil(t, found)
}

// TestBookListRepository_GetBookListsByUser 测试获取用户书单列表
func TestBookListRepository_GetBookListsByUser(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()

	// 创建多个书单
	for i := 0; i < 5; i++ {
		bookList := &social.BookList{
			UserID:      userID,
			UserName:    "测试用户",
			Title:       "书单标题",
			Description: "书单描述",
			IsPublic:    true,
		}
		err := repo.CreateBookList(ctx, bookList)
		require.NoError(t, err)
	}

	// Act
	bookLists, total, err := repo.GetBookListsByUser(ctx, userID, 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, bookLists)
	assert.GreaterOrEqual(t, total, int64(5))
	assert.GreaterOrEqual(t, len(bookLists), 5)
}

// TestBookListRepository_GetBookListsByUser_Pagination 测试分页功能
func TestBookListRepository_GetBookListsByUser_Pagination(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()

	// 创建15个书单
	for i := 0; i < 15; i++ {
		bookList := &social.BookList{
			UserID:      userID,
			UserName:    "测试用户",
			Title:       "书单标题",
			IsPublic:    true,
		}
		err := repo.CreateBookList(ctx, bookList)
		require.NoError(t, err)
	}

	// Act - 第一页
	page1, total1, err := repo.GetBookListsByUser(ctx, userID, 1, 10)
	require.NoError(t, err)

	// Act - 第二页
	page2, total2, err := repo.GetBookListsByUser(ctx, userID, 2, 10)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, total1, total2)
	assert.GreaterOrEqual(t, total1, int64(15))
	assert.GreaterOrEqual(t, len(page1), 10)
	assert.LessOrEqual(t, len(page2), 5)
}

// TestBookListRepository_GetBookListsByUser_NotExists 测试无效用户ID
func TestBookListRepository_GetBookListsByUser_NotExists(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	// Act
	bookLists, total, err := repo.GetBookListsByUser(ctx, "invalid_userid", 1, 10)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Equal(t, 0, len(bookLists))
}

// TestBookListRepository_GetPublicBookLists 测试获取公开书单列表
func TestBookListRepository_GetPublicBookLists(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()

	// 创建公开书单
	publicList := &social.BookList{
		UserID:      userID,
		UserName:    "测试用户",
		Title:       "公开书单",
		Description: "这是公开书单",
		IsPublic:    true,
	}
	err := repo.CreateBookList(ctx, publicList)
	require.NoError(t, err)

	// 创建私有书单
	privateList := &social.BookList{
		UserID:      userID,
		UserName:    "测试用户",
		Title:       "私有书单",
		Description: "这是私有书单",
		IsPublic:    false,
	}
	err = repo.CreateBookList(ctx, privateList)
	require.NoError(t, err)

	// Act
	bookLists, total, err := repo.GetPublicBookLists(ctx, 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, bookLists)

	// 验证只有公开书单
	for _, list := range bookLists {
		assert.True(t, list.IsPublic)
	}
	assert.GreaterOrEqual(t, total, int64(1))
}

// TestBookListRepository_GetBookListsByCategory 测试按分类获取书单
func TestBookListRepository_GetBookListsByCategory(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()
	category := "科幻小说"

	// 创建指定分类的书单
	bookList := &social.BookList{
		UserID:      userID,
		UserName:    "测试用户",
		Title:       "科幻书单",
		Description: "科幻小说推荐",
		IsPublic:    true,
		Category:    category,
	}
	err := repo.CreateBookList(ctx, bookList)
	require.NoError(t, err)

	// Act
	bookLists, total, err := repo.GetBookListsByCategory(ctx, category, 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, bookLists)
	assert.GreaterOrEqual(t, total, int64(1))

	// 验证分类匹配
	for _, list := range bookLists {
		assert.Equal(t, category, list.Category)
		assert.True(t, list.IsPublic)
	}
}

// TestBookListRepository_GetBookListsByTag 测试按标签获取书单
func TestBookListRepository_GetBookListsByTag(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()
	tag := "推荐"

	// 创建带标签的书单
	bookList := &social.BookList{
		UserID:      userID,
		UserName:    "测试用户",
		Title:       "推荐书单",
		Description: "推荐的好书",
		IsPublic:    true,
		Tags:        []string{tag, "经典"},
	}
	err := repo.CreateBookList(ctx, bookList)
	require.NoError(t, err)

	// Act
	bookLists, total, err := repo.GetBookListsByTag(ctx, tag, 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, bookLists)
	assert.GreaterOrEqual(t, total, int64(1))

	// 验证标签匹配且为公开
	for _, list := range bookLists {
		assert.Contains(t, list.Tags, tag)
		assert.True(t, list.IsPublic)
	}
}

// TestBookListRepository_SearchBookLists 测试搜索书单
func TestBookListRepository_SearchBookLists(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()

	// 创建可搜索的书单
	bookList := &social.BookList{
		UserID:      userID,
		UserName:    "测试用户",
		Title:       "三体系列书单",
		Description: "刘慈欣科幻作品推荐",
		IsPublic:    true,
	}
	err := repo.CreateBookList(ctx, bookList)
	require.NoError(t, err)

	// Act - 搜索标题
	bookLists1, total1, err := repo.SearchBookLists(ctx, "三体", 1, 10)

	// Act - 搜索描述
	bookLists2, total2, err := repo.SearchBookLists(ctx, "刘慈欣", 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, bookLists1)
	assert.GreaterOrEqual(t, total1, int64(1))

	require.NoError(t, err)
	assert.NotNil(t, bookLists2)
	assert.GreaterOrEqual(t, total2, int64(1))
}

// TestBookListRepository_SearchBookLists_CaseInsensitive 测试搜索不区分大小写
func TestBookListRepository_SearchBookLists_CaseInsensitive(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()

	bookList := &social.BookList{
		UserID:      userID,
		UserName:    "测试用户",
		Title:       "Python编程书籍",
		Description: "Python学习资源",
		IsPublic:    true,
	}
	err := repo.CreateBookList(ctx, bookList)
	require.NoError(t, err)

	// Act - 小写搜索
	bookLists, total, err := repo.SearchBookLists(ctx, "python", 1, 10)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	assert.GreaterOrEqual(t, len(bookLists), 1)
}

// TestBookListRepository_UpdateBookList 测试更新书单
func TestBookListRepository_UpdateBookList(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()
	bookList := createTestBookList(t, repo, ctx, userID)

	// Act - 更新书单
	updates := map[string]interface{}{
		"title":       "更新后的标题",
		"description": "更新后的描述",
		"is_public":   false,
	}
	err := repo.UpdateBookList(ctx, bookList.ID.Hex(), updates)

	// Assert
	require.NoError(t, err)

	// 验证更新
	found, err := repo.GetBookListByID(ctx, bookList.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, "更新后的标题", found.Title)
	assert.Equal(t, "更新后的描述", found.Description)
	assert.Equal(t, false, found.IsPublic)
}

// TestBookListRepository_UpdateBookList_NotFound 测试更新不存在的书单
func TestBookListRepository_UpdateBookList_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	validID := primitive.NewObjectID().Hex()

	// Act
	updates := map[string]interface{}{
		"title": "新标题",
	}
	err := repo.UpdateBookList(ctx, validID, updates)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "书单不存在")
}

// TestBookListRepository_UpdateBookList_InvalidID 测试更新无效ID
func TestBookListRepository_UpdateBookList_InvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	// Act
	updates := map[string]interface{}{
		"title": "新标题",
	}
	err := repo.UpdateBookList(ctx, "invalid_id", updates)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "无效的书单ID")
}

// TestBookListRepository_DeleteBookList 测试删除书单
func TestBookListRepository_DeleteBookList(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()
	bookList := createTestBookList(t, repo, ctx, userID)

	// Act - 删除书单
	err := repo.DeleteBookList(ctx, bookList.ID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证已删除
	found, err := repo.GetBookListByID(ctx, bookList.ID.Hex())
	require.NoError(t, err)
	assert.Nil(t, found)
}

// TestBookListRepository_DeleteBookList_NotFound 测试删除不存在的书单
func TestBookListRepository_DeleteBookList_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	validID := primitive.NewObjectID().Hex()

	// Act
	err := repo.DeleteBookList(ctx, validID)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "书单不存在")
}

// TestBookListRepository_DeleteBookList_InvalidID 测试删除无效ID
func TestBookListRepository_DeleteBookList_InvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	// Act
	err := repo.DeleteBookList(ctx, "invalid_id")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "无效的书单ID")
}

// TestBookListRepository_AddBookToList 测试添加书籍到书单
func TestBookListRepository_AddBookToList(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()
	bookList := createTestBookList(t, repo, ctx, userID)

	bookItem := &social.BookListItem{
		BookID:      "book123",
		BookTitle:   "测试书籍",
		BookCover:   "cover.jpg",
		AuthorName:  "测试作者",
		Description: "书籍描述",
		Comment:     "推荐语",
		Order:       1,
	}

	// Act
	err := repo.AddBookToList(ctx, bookList.ID.Hex(), bookItem)

	// Assert
	require.NoError(t, err)

	// 验证书籍已添加
	found, err := repo.GetBookListByID(ctx, bookList.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, 1, len(found.Books))
	assert.Equal(t, "book123", found.Books[0].BookID)
}

// TestBookListRepository_AddBookToList_InvalidBookListID 测试添加书籍到无效书单
func TestBookListRepository_AddBookToList_InvalidBookListID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	bookItem := &social.BookListItem{
		BookID:    "book123",
		BookTitle: "测试书籍",
		Order:     1,
	}

	// Act
	err := repo.AddBookToList(ctx, "invalid_id", bookItem)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "无效的书单ID")
}

// TestBookListRepository_AddBookToList_BookListNotExists 测试添加书籍到不存在的书单
func TestBookListRepository_AddBookToList_BookListNotExists(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	validID := primitive.NewObjectID().Hex()
	bookItem := &social.BookListItem{
		BookID:    "book123",
		BookTitle: "测试书籍",
		Order:     1,
	}

	// Act
	err := repo.AddBookToList(ctx, validID, bookItem)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "书单不存在")
}

// TestBookListRepository_RemoveBookFromList 测试从书单移除书籍
func TestBookListRepository_RemoveBookFromList(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()
	bookList := createTestBookList(t, repo, ctx, userID)

	// 添加书籍
	bookItem := &social.BookListItem{
		BookID:    "book123",
		BookTitle: "测试书籍",
		Order:     1,
	}
	err := repo.AddBookToList(ctx, bookList.ID.Hex(), bookItem)
	require.NoError(t, err)

	// Act - 移除书籍
	err = repo.RemoveBookFromList(ctx, bookList.ID.Hex(), "book123")

	// Assert
	require.NoError(t, err)

	// 验证书籍已移除
	found, err := repo.GetBookListByID(ctx, bookList.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, 0, len(found.Books))
}

// TestBookListRepository_RemoveBookFromList_InvalidBookListID 测试从无效书单移除书籍
func TestBookListRepository_RemoveBookFromList_InvalidBookListID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	// Act
	err := repo.RemoveBookFromList(ctx, "invalid_id", "book123")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "无效的书单ID")
}

// TestBookListRepository_RemoveBookFromList_BookListNotExists 测试从不存在的书单移除书籍
func TestBookListRepository_RemoveBookFromList_BookListNotExists(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	validID := primitive.NewObjectID().Hex()

	// Act
	err := repo.RemoveBookFromList(ctx, validID, "book123")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "书单不存在")
}

// TestBookListRepository_UpdateBookInList 测试更新书单中的书籍
func TestBookListRepository_UpdateBookInList(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()
	bookList := createTestBookList(t, repo, ctx, userID)

	// 添加书籍
	bookItem := &social.BookListItem{
		BookID:    "book123",
		BookTitle: "原始标题",
		Comment:   "原始推荐语",
		Order:     1,
	}
	err := repo.AddBookToList(ctx, bookList.ID.Hex(), bookItem)
	require.NoError(t, err)

	// Act - 更新书籍信息
	updates := map[string]interface{}{
		"comment": "更新后的推荐语",
		"order":   2,
	}
	err = repo.UpdateBookInList(ctx, bookList.ID.Hex(), "book123", updates)

	// Assert
	require.NoError(t, err)

	// 验证更新
	books, err := repo.GetBooksInList(ctx, bookList.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, 1, len(books))
	assert.Equal(t, "更新后的推荐语", books[0].Comment)
	assert.Equal(t, 2, books[0].Order)
}

// TestBookListRepository_UpdateBookInList_BookNotFound 测试更新不存在的书籍
func TestBookListRepository_UpdateBookInList_BookNotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()
	bookList := createTestBookList(t, repo, ctx, userID)

	// Act
	updates := map[string]interface{}{
		"comment": "新推荐语",
	}
	err := repo.UpdateBookInList(ctx, bookList.ID.Hex(), "nonexistent_book", updates)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "书单或书籍不存在")
}

// TestBookListRepository_ReorderBooks 测试重新排序书籍
func TestBookListRepository_ReorderBooks(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()
	bookList := createTestBookList(t, repo, ctx, userID)

	// 添加多本书
	book1 := &social.BookListItem{BookID: "book1", BookTitle: "书籍1", Order: 1}
	book2 := &social.BookListItem{BookID: "book2", BookTitle: "书籍2", Order: 2}
	book3 := &social.BookListItem{BookID: "book3", BookTitle: "书籍3", Order: 3}

	err := repo.AddBookToList(ctx, bookList.ID.Hex(), book1)
	require.NoError(t, err)
	err = repo.AddBookToList(ctx, bookList.ID.Hex(), book2)
	require.NoError(t, err)
	err = repo.AddBookToList(ctx, bookList.ID.Hex(), book3)
	require.NoError(t, err)

	// Act - 重新排序
	newOrders := map[string]int{
		"book1": 3,
		"book2": 1,
		"book3": 2,
	}
	err = repo.ReorderBooks(ctx, bookList.ID.Hex(), newOrders)

	// Assert
	require.NoError(t, err)

	// 验证排序
	books, err := repo.GetBooksInList(ctx, bookList.ID.Hex())
	require.NoError(t, err)

	// 创建bookID到order的映射
	bookOrderMap := make(map[string]int)
	for _, book := range books {
		bookOrderMap[book.BookID] = book.Order
	}

	assert.Equal(t, 3, bookOrderMap["book1"])
	assert.Equal(t, 1, bookOrderMap["book2"])
	assert.Equal(t, 2, bookOrderMap["book3"])
}

// TestBookListRepository_ReorderBooks_InvalidID 测试重新排序无效书单
func TestBookListRepository_ReorderBooks_InvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	// Act
	bookOrders := map[string]int{
		"book1": 1,
	}
	err := repo.ReorderBooks(ctx, "invalid_id", bookOrders)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "无效的书单ID")
}

// TestBookListRepository_GetBooksInList 测试获取书单中的书籍
func TestBookListRepository_GetBooksInList(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()
	bookList := createTestBookList(t, repo, ctx, userID)

	// 添加书籍
	book1 := &social.BookListItem{BookID: "book1", BookTitle: "书籍1", Order: 1}
	book2 := &social.BookListItem{BookID: "book2", BookTitle: "书籍2", Order: 2}

	err := repo.AddBookToList(ctx, bookList.ID.Hex(), book1)
	require.NoError(t, err)
	err = repo.AddBookToList(ctx, bookList.ID.Hex(), book2)
	require.NoError(t, err)

	// Act
	books, err := repo.GetBooksInList(ctx, bookList.ID.Hex())

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 2, len(books))
}

// TestBookListRepository_GetBooksInList_InvalidID 测试获取无效书单的书籍
func TestBookListRepository_GetBooksInList_InvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	// Act
	books, err := repo.GetBooksInList(ctx, "invalid_id")

	// Assert
	require.Error(t, err)
	assert.Nil(t, books)
	assert.Contains(t, err.Error(), "无效的书单ID")
}

// TestBookListRepository_GetBooksInList_BookListNotExists 测试获取不存在书单的书籍
func TestBookListRepository_GetBooksInList_BookListNotExists(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	validID := primitive.NewObjectID().Hex()

	// Act
	books, err := repo.GetBooksInList(ctx, validID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, books)
	assert.Contains(t, err.Error(), "书单不存在")
}

// TestBookListRepository_CreateBookListLike 测试创建书单点赞
func TestBookListRepository_CreateBookListLike(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()
	bookList := createTestBookList(t, repo, ctx, userID)

	likerID := primitive.NewObjectID().Hex()
	bookListLike := &social.BookListLike{
		BookListID: bookList.ID.Hex(),
		UserID:     likerID,
	}

	// Act
	err := repo.CreateBookListLike(ctx, bookListLike)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, bookListLike.ID)
	assert.NotZero(t, bookListLike.CreatedAt)

	// 验证书单点赞数增加
	found, err := repo.GetBookListByID(ctx, bookList.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, 1, found.LikeCount)
}

// TestBookListRepository_DeleteBookListLike 测试删除书单点赞
func TestBookListRepository_DeleteBookListLike(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()
	bookList := createTestBookList(t, repo, ctx, userID)

	likerID := primitive.NewObjectID().Hex()
	bookListLike := &social.BookListLike{
		BookListID: bookList.ID.Hex(),
		UserID:     likerID,
	}

	// 先点赞
	err := repo.CreateBookListLike(ctx, bookListLike)
	require.NoError(t, err)

	// Act - 取消点赞
	err = repo.DeleteBookListLike(ctx, bookList.ID.Hex(), likerID)

	// Assert
	require.NoError(t, err)

	// 验证书单点赞数减少
	found, err := repo.GetBookListByID(ctx, bookList.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, 0, found.LikeCount)
}

// TestBookListRepository_DeleteBookListLike_NotExists 测试删除不存在的点赞记录
func TestBookListRepository_DeleteBookListLike_NotExists(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	// Act - 删除不存在的点赞记录应该返回成功（幂等操作）
	err := repo.DeleteBookListLike(ctx, "booklist123", "user123")

	// Assert
	require.NoError(t, err)
}

// TestBookListRepository_GetBookListLike 测试获取书单点赞记录
func TestBookListRepository_GetBookListLike(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()
	bookList := createTestBookList(t, repo, ctx, userID)

	likerID := primitive.NewObjectID().Hex()
	bookListLike := &social.BookListLike{
		BookListID: bookList.ID.Hex(),
		UserID:     likerID,
	}
	err := repo.CreateBookListLike(ctx, bookListLike)
	require.NoError(t, err)

	// Act
	found, err := repo.GetBookListLike(ctx, bookList.ID.Hex(), likerID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, bookList.ID.Hex(), found.BookListID)
	assert.Equal(t, likerID, found.UserID)
}

// TestBookListRepository_GetBookListLike_NotFound 测试获取不存在的点赞记录
func TestBookListRepository_GetBookListLike_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	// Act
	found, err := repo.GetBookListLike(ctx, "booklist123", "user123")

	// Assert
	require.NoError(t, err)
	assert.Nil(t, found)
}

// TestBookListRepository_IsBookListLiked 测试检查是否已点赞
func TestBookListRepository_IsBookListLiked(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()
	bookList := createTestBookList(t, repo, ctx, userID)

	likerID := primitive.NewObjectID().Hex()

	// Act - 点赞前
	liked1, err := repo.IsBookListLiked(ctx, bookList.ID.Hex(), likerID)

	// 点赞
	bookListLike := &social.BookListLike{
		BookListID: bookList.ID.Hex(),
		UserID:     likerID,
	}
	err = repo.CreateBookListLike(ctx, bookListLike)
	require.NoError(t, err)

	// Act - 点赞后
	liked2, err := repo.IsBookListLiked(ctx, bookList.ID.Hex(), likerID)

	// Assert
	require.NoError(t, err)
	assert.False(t, liked1)

	require.NoError(t, err)
	assert.True(t, liked2)
}

// TestBookListRepository_GetBookListLikes 测试获取书单点赞列表
func TestBookListRepository_GetBookListLikes(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()
	bookList := createTestBookList(t, repo, ctx, userID)

	// 创建多个点赞
	for i := 0; i < 5; i++ {
		likerID := primitive.NewObjectID().Hex()
		bookListLike := &social.BookListLike{
			BookListID: bookList.ID.Hex(),
			UserID:     likerID,
		}
		err := repo.CreateBookListLike(ctx, bookListLike)
		require.NoError(t, err)
	}

	// Act
	likes, total, err := repo.GetBookListLikes(ctx, bookList.ID.Hex(), 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, likes)
	assert.GreaterOrEqual(t, total, int64(5))
	assert.GreaterOrEqual(t, len(likes), 5)
}

// TestBookListRepository_IncrementBookListLikeCount 测试增加书单点赞数
func TestBookListRepository_IncrementBookListLikeCount(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()
	bookList := createTestBookList(t, repo, ctx, userID)

	// Act
	err := repo.IncrementBookListLikeCount(ctx, bookList.ID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证点赞数
	found, err := repo.GetBookListByID(ctx, bookList.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, 1, found.LikeCount)
}

// TestBookListRepository_IncrementBookListLikeCount_InvalidID 测试增加无效书单点赞数
func TestBookListRepository_IncrementBookListLikeCount_InvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	// Act
	err := repo.IncrementBookListLikeCount(ctx, "invalid_id")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "无效的书单ID")
}

// TestBookListRepository_DecrementBookListLikeCount 测试减少书单点赞数
func TestBookListRepository_DecrementBookListLikeCount(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()
	bookList := createTestBookList(t, repo, ctx, userID)

	// 先增加几次点赞
	err := repo.IncrementBookListLikeCount(ctx, bookList.ID.Hex())
	require.NoError(t, err)
	err = repo.IncrementBookListLikeCount(ctx, bookList.ID.Hex())
	require.NoError(t, err)
	err = repo.IncrementBookListLikeCount(ctx, bookList.ID.Hex())
	require.NoError(t, err)

	// Act - 减少点赞数
	err = repo.DecrementBookListLikeCount(ctx, bookList.ID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证点赞数
	found, err := repo.GetBookListByID(ctx, bookList.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, 2, found.LikeCount)
}

// TestBookListRepository_ForkBookList 测试复制书单
func TestBookListRepository_ForkBookList(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()
	originalList := createTestBookList(t, repo, ctx, userID)

	// 添加书籍到原始书单
	bookItem := &social.BookListItem{
		BookID:    "book123",
		BookTitle: "测试书籍",
		Order:     1,
	}
	err := repo.AddBookToList(ctx, originalList.ID.Hex(), bookItem)
	require.NoError(t, err)

	forkUserID := primitive.NewObjectID().Hex()

	// Act
	forkedList, err := repo.ForkBookList(ctx, originalList.ID.Hex(), forkUserID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, forkedList)
	assert.NotEqual(t, originalList.ID, forkedList.ID)
	assert.Equal(t, forkUserID, forkedList.UserID)
	assert.Equal(t, originalList.Title, forkedList.Title)
	assert.Equal(t, originalList.Description, forkedList.Description)
	assert.Equal(t, originalList.IsPublic, forkedList.IsPublic)
	assert.Equal(t, 1, len(forkedList.Books))
	assert.NotNil(t, forkedList.OriginalID)
	assert.Equal(t, originalList.ID, *forkedList.OriginalID)

	// 验证原始书单的fork_count增加
	original, err := repo.GetBookListByID(ctx, originalList.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, 1, original.ForkCount)
}

// TestBookListRepository_ForkBookList_OriginalNotExists 测试复制不存在的书单
func TestBookListRepository_ForkBookList_OriginalNotExists(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	validID := primitive.NewObjectID().Hex()
	userID := primitive.NewObjectID().Hex()

	// Act
	forkedList, err := repo.ForkBookList(ctx, validID, userID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, forkedList)
	assert.Contains(t, err.Error(), "原始书单不存在")
}

// TestBookListRepository_IncrementForkCount 测试增加被复制次数
func TestBookListRepository_IncrementForkCount(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()
	bookList := createTestBookList(t, repo, ctx, userID)

	// Act
	err := repo.IncrementForkCount(ctx, bookList.ID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证fork_count
	found, err := repo.GetBookListByID(ctx, bookList.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, 1, found.ForkCount)
}

// TestBookListRepository_IncrementForkCount_InvalidID 测试增加无效书单的fork_count
func TestBookListRepository_IncrementForkCount_InvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	// Act
	err := repo.IncrementForkCount(ctx, "invalid_id")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "无效的书单ID")
}

// TestBookListRepository_GetForkedBookLists 测试获取复制的书单列表
func TestBookListRepository_GetForkedBookLists(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()
	originalList := createTestBookList(t, repo, ctx, userID)

	// 创建多个复制书单
	for i := 0; i < 3; i++ {
		forkUserID := primitive.NewObjectID().Hex()
		_, err := repo.ForkBookList(ctx, originalList.ID.Hex(), forkUserID)
		require.NoError(t, err)
	}

	// Act
	forkedLists, total, err := repo.GetForkedBookLists(ctx, originalList.ID.Hex(), 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, forkedLists)
	assert.GreaterOrEqual(t, total, int64(3))
	assert.GreaterOrEqual(t, len(forkedLists), 3)

	// 验证所有书单都指向原始书单
	for _, list := range forkedLists {
		assert.NotNil(t, list.OriginalID)
		assert.Equal(t, originalList.ID, *list.OriginalID)
	}
}

// TestBookListRepository_GetForkedBookLists_InvalidID 测试获取无效书单的复制列表
func TestBookListRepository_GetForkedBookLists_InvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	// Act
	forkedLists, total, err := repo.GetForkedBookLists(ctx, "invalid_id", 1, 10)

	// Assert
	require.Error(t, err)
	assert.Nil(t, forkedLists)
	assert.Equal(t, int64(0), total)
	assert.Contains(t, err.Error(), "无效的原始书单ID")
}

// TestBookListRepository_IncrementViewCount 测试增加浏览次数
func TestBookListRepository_IncrementViewCount(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()
	bookList := createTestBookList(t, repo, ctx, userID)

	// Act
	err := repo.IncrementViewCount(ctx, bookList.ID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证view_count
	found, err := repo.GetBookListByID(ctx, bookList.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, 1, found.ViewCount)
}

// TestBookListRepository_IncrementViewCount_InvalidID 测试增加无效书单浏览次数
func TestBookListRepository_IncrementViewCount_InvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	// Act
	err := repo.IncrementViewCount(ctx, "invalid_id")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "无效的书单ID")
}

// TestBookListRepository_CountUserBookLists 测试统计用户书单数
func TestBookListRepository_CountUserBookLists(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()

	// 创建多个书单
	for i := 0; i < 5; i++ {
		bookList := &social.BookList{
			UserID:   userID,
			UserName: "测试用户",
			Title:    "书单",
			IsPublic: true,
		}
		err := repo.CreateBookList(ctx, bookList)
		require.NoError(t, err)
	}

	// Act
	count, err := repo.CountUserBookLists(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(5))
}

// TestBookListRepository_CountUserBookLists_NotExists 测试统计不存在用户的书单数
func TestBookListRepository_CountUserBookLists_NotExists(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	// Act
	count, err := repo.CountUserBookLists(ctx, "nonexistent_user_id")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

// TestBookListRepository_Health 测试健康检查
func TestBookListRepository_Health(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	// Act
	err := repo.Health(ctx)

	// Assert
	assert.NoError(t, err)
}

// TestBookListRepository_ComprehensiveFlow 测试综合流程
func TestBookListRepository_ComprehensiveFlow(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupBookListRepo(t)
	defer cleanup()

	userID := primitive.NewObjectID().Hex()

	// 1. 创建书单
	bookList := &social.BookList{
		UserID:      userID,
		UserName:    "测试用户",
		Title:       "综合测试书单",
		Description: "这是一个综合测试",
		IsPublic:    true,
		Tags:        []string{"测试", "综合"},
		Category:    "测试分类",
	}
	err := repo.CreateBookList(ctx, bookList)
	require.NoError(t, err)

	// 2. 添加书籍
	book1 := &social.BookListItem{BookID: "book1", BookTitle: "书籍1", Order: 1}
	book2 := &social.BookListItem{BookID: "book2", BookTitle: "书籍2", Order: 2}
	err = repo.AddBookToList(ctx, bookList.ID.Hex(), book1)
	require.NoError(t, err)
	err = repo.AddBookToList(ctx, bookList.ID.Hex(), book2)
	require.NoError(t, err)

	// 3. 获取书单中的书籍
	books, err := repo.GetBooksInList(ctx, bookList.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, 2, len(books))

	// 4. 更新书单
	updates := map[string]interface{}{
		"title": "更新后的综合测试书单",
	}
	err = repo.UpdateBookList(ctx, bookList.ID.Hex(), updates)
	require.NoError(t, err)

	// 5. 点赞书单
	likerID := primitive.NewObjectID().Hex()
	bookListLike := &social.BookListLike{
		BookListID: bookList.ID.Hex(),
		UserID:     likerID,
	}
	err = repo.CreateBookListLike(ctx, bookListLike)
	require.NoError(t, err)

	// 6. 检查点赞
	liked, err := repo.IsBookListLiked(ctx, bookList.ID.Hex(), likerID)
	require.NoError(t, err)
	assert.True(t, liked)

	// 7. 增加浏览次数
	err = repo.IncrementViewCount(ctx, bookList.ID.Hex())
	require.NoError(t, err)

	// 8. 复制书单
	forkUserID := primitive.NewObjectID().Hex()
	forkedList, err := repo.ForkBookList(ctx, bookList.ID.Hex(), forkUserID)
	require.NoError(t, err)
	assert.NotNil(t, forkedList)

	// 9. 验证最终状态
	final, err := repo.GetBookListByID(ctx, bookList.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, "更新后的综合测试书单", final.Title)
	assert.Equal(t, 1, final.LikeCount)
	assert.Equal(t, 1, final.ViewCount)
	assert.Equal(t, 1, final.ForkCount)

	// 10. 统计用户书单数
	count, err := repo.CountUserBookLists(ctx, userID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(1))
}
