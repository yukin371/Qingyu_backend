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

// setupCollectionRepo 测试辅助函数
func setupCollectionRepo(t *testing.T) (*socialRepo.MongoCollectionRepository, context.Context, func()) {
	db, cleanup := testutil.SetupTestDB(t)
	repo := socialRepo.NewMongoCollectionRepository(db)
	ctx := context.Background()
	return repo, ctx, cleanup
}

// =========================
// 收藏管理测试
// =========================

// TestCollectionRepository_Create 测试创建收藏
func TestCollectionRepository_Create(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	collection := &social.Collection{
		UserID:   "user_create_" + primitive.NewObjectID().Hex(),
		BookID:   "book_create_" + primitive.NewObjectID().Hex(),
		FolderID: "",
		Tags:     []string{"玄幻", "推荐"},
		Note:     "这是一本好书",
		IsPublic: false,
	}

	// Act
	err := repo.Create(ctx, collection)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, collection.ID)
	assert.NotZero(t, collection.CreatedAt)
	assert.NotZero(t, collection.UpdatedAt)
}

// TestCollectionRepository_Create_Duplicate 测试重复收藏
func TestCollectionRepository_Create_Duplicate(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	userID := "user_dup_" + primitive.NewObjectID().Hex()
	bookID := "book_dup_" + primitive.NewObjectID().Hex()
	collection := &social.Collection{
		UserID:   userID,
		BookID:   bookID,
		IsPublic: false,
	}
	err := repo.Create(ctx, collection)
	require.NoError(t, err)

	// Act - 尝试重复收藏
	err = repo.Create(ctx, &social.Collection{
		UserID:   userID,
		BookID:   bookID,
		IsPublic: false,
	})

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "已经收藏")
}

// TestCollectionRepository_Create_MissingParams 测试缺少必需参数
func TestCollectionRepository_Create_MissingParams(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	tests := []struct {
		name       string
		collection *social.Collection
		wantErr    string
	}{
		{
			name: "缺少UserID",
			collection: &social.Collection{
				BookID: "book123",
			},
			wantErr: "用户ID和书籍ID不能为空",
		},
		{
			name: "缺少BookID",
			collection: &social.Collection{
				UserID: "user123",
			},
			wantErr: "用户ID和书籍ID不能为空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := repo.Create(ctx, tt.collection)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

// TestCollectionRepository_GetByID 测试根据ID获取收藏
func TestCollectionRepository_GetByID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	collection := &social.Collection{
		UserID:   "user_getbyid_" + primitive.NewObjectID().Hex(),
		BookID:   "book_getbyid_" + primitive.NewObjectID().Hex(),
		IsPublic: false,
	}
	err := repo.Create(ctx, collection)
	require.NoError(t, err)

	// Act
	found, err := repo.GetByID(ctx, collection.ID.Hex())

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, collection.UserID, found.UserID)
	assert.Equal(t, collection.BookID, found.BookID)
}

// TestCollectionRepository_GetByID_NotFound 测试获取不存在的收藏
func TestCollectionRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	fakeID := primitive.NewObjectID().Hex()

	// Act
	found, err := repo.GetByID(ctx, fakeID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, found)
	assert.Contains(t, err.Error(), "not found")
}

// TestCollectionRepository_GetByID_InvalidID 测试无效ID
func TestCollectionRepository_GetByID_InvalidID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	// Act
	found, err := repo.GetByID(ctx, "invalid_id")

	// Assert
	require.Error(t, err)
	assert.Nil(t, found)
	assert.Contains(t, err.Error(), "invalid collection ID")
}

// TestCollectionRepository_GetByUserAndBook 测试根据用户和书籍获取收藏
func TestCollectionRepository_GetByUserAndBook(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	userID := "user_getuserandbook_" + primitive.NewObjectID().Hex()
	bookID := "book_getuserandbook_" + primitive.NewObjectID().Hex()
	collection := &social.Collection{
		UserID:   userID,
		BookID:   bookID,
		IsPublic: false,
	}
	err := repo.Create(ctx, collection)
	require.NoError(t, err)

	// Act
	found, err := repo.GetByUserAndBook(ctx, userID, bookID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, userID, found.UserID)
	assert.Equal(t, bookID, found.BookID)
}

// TestCollectionRepository_GetByUserAndBook_NotCollected 测试获取未收藏的书籍
func TestCollectionRepository_GetByUserAndBook_NotCollected(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	// Act
	found, err := repo.GetByUserAndBook(ctx, "user123", "book123")

	// Assert
	require.NoError(t, err)
	assert.Nil(t, found)
}

// TestCollectionRepository_GetCollectionsByUser 测试获取用户收藏列表
func TestCollectionRepository_GetCollectionsByUser(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	userID := "user_getcollections_" + primitive.NewObjectID().Hex()
	for i := 0; i < 5; i++ {
		collection := &social.Collection{
			UserID:   userID,
			BookID:   primitive.NewObjectID().Hex(),
			IsPublic: false,
		}
		err := repo.Create(ctx, collection)
		require.NoError(t, err)
	}

	// Act
	collections, total, err := repo.GetCollectionsByUser(ctx, userID, "", 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, collections)
	assert.GreaterOrEqual(t, total, int64(5))
	assert.LessOrEqual(t, len(collections), 10)
}

// TestCollectionRepository_GetCollectionsByUser_WithFolder 测试获取指定收藏夹的收藏
func TestCollectionRepository_GetCollectionsByUser_WithFolder(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	userID := "user_withfolder_" + primitive.NewObjectID().Hex()
	folder := &social.CollectionFolder{
		UserID:      userID,
		Name:        "测试收藏夹",
		Description: "测试用",
		IsPublic:    false,
	}
	err := repo.CreateFolder(ctx, folder)
	require.NoError(t, err)

	// 创建收藏
	collection := &social.Collection{
		UserID:   userID,
		BookID:   primitive.NewObjectID().Hex(),
		FolderID: folder.ID.Hex(),
		IsPublic: false,
	}
	err = repo.Create(ctx, collection)
	require.NoError(t, err)

	// Act
	collections, total, err := repo.GetCollectionsByUser(ctx, userID, folder.ID.Hex(), 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, collections)
	assert.GreaterOrEqual(t, total, int64(1))
	// 验证返回的收藏都在指定收藏夹中
	for _, col := range collections {
		assert.Equal(t, folder.ID.Hex(), col.FolderID)
	}
}

// TestCollectionRepository_GetCollectionsByTag 测试根据标签获取收藏
func TestCollectionRepository_GetCollectionsByTag(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	userID := "user_getbytag_" + primitive.NewObjectID().Hex()
	tag := "玄幻"
	for i := 0; i < 3; i++ {
		collection := &social.Collection{
			UserID:   userID,
			BookID:   primitive.NewObjectID().Hex(),
			Tags:     []string{tag, "推荐"},
			IsPublic: false,
		}
		err := repo.Create(ctx, collection)
		require.NoError(t, err)
	}

	// Act
	collections, total, err := repo.GetCollectionsByTag(ctx, userID, tag, 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, collections)
	assert.GreaterOrEqual(t, total, int64(3))
	// 验证返回的收藏都包含指定标签
	for _, col := range collections {
		assert.Contains(t, col.Tags, tag)
	}
}

// TestCollectionRepository_Update 测试更新收藏
func TestCollectionRepository_Update(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	collection := &social.Collection{
		UserID:   "user_update_" + primitive.NewObjectID().Hex(),
		BookID:   "book_update_" + primitive.NewObjectID().Hex(),
		IsPublic: false,
	}
	err := repo.Create(ctx, collection)
	require.NoError(t, err)

	// Act - 更新收藏
	updates := map[string]interface{}{
		"note":      "更新后的笔记",
		"is_public": true,
		"tags":      []string{"更新", "标签"},
	}
	err = repo.Update(ctx, collection.ID.Hex(), updates)

	// Assert
	require.NoError(t, err)

	// 验证更新
	found, err := repo.GetByID(ctx, collection.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, "更新后的笔记", found.Note)
	assert.True(t, found.IsPublic)
	assert.Equal(t, []string{"更新", "标签"}, found.Tags)
}

// TestCollectionRepository_Update_NotFound 测试更新不存在的收藏
func TestCollectionRepository_Update_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	fakeID := primitive.NewObjectID().Hex()
	updates := map[string]interface{}{
		"note": "更新后的笔记",
	}

	// Act
	err := repo.Update(ctx, fakeID, updates)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestCollectionRepository_Delete 测试删除收藏
func TestCollectionRepository_Delete(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	collection := &social.Collection{
		UserID:   "user_delete_" + primitive.NewObjectID().Hex(),
		BookID:   "book_delete_" + primitive.NewObjectID().Hex(),
		IsPublic: false,
	}
	err := repo.Create(ctx, collection)
	require.NoError(t, err)

	// Act - 删除收藏
	err = repo.Delete(ctx, collection.ID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证已删除
	found, err := repo.GetByID(ctx, collection.ID.Hex())
	assert.Error(t, err)
	assert.Nil(t, found)
}

// TestCollectionRepository_Delete_NotFound 测试删除不存在的收藏
func TestCollectionRepository_Delete_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	fakeID := primitive.NewObjectID().Hex()

	// Act
	err := repo.Delete(ctx, fakeID)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// =========================
// 收藏夹管理测试
// =========================

// TestCollectionRepository_CreateFolder 测试创建收藏夹
func TestCollectionRepository_CreateFolder(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	folder := &social.CollectionFolder{
		UserID:      "user_createfolder_" + primitive.NewObjectID().Hex(),
		Name:        "我的收藏夹",
		Description: "收藏我喜欢的书籍",
		IsPublic:    false,
	}

	// Act
	err := repo.CreateFolder(ctx, folder)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, folder.ID)
	assert.NotZero(t, folder.CreatedAt)
	assert.NotZero(t, folder.UpdatedAt)
	assert.Equal(t, 0, folder.BookCount)
}

// TestCollectionRepository_CreateFolder_MissingParams 测试缺少必需参数
func TestCollectionRepository_CreateFolder_MissingParams(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	tests := []struct {
		name    string
		folder  *social.CollectionFolder
		wantErr string
	}{
		{
			name: "缺少UserID",
			folder: &social.CollectionFolder{
				Name: "测试收藏夹",
			},
			wantErr: "用户ID和收藏夹名称不能为空",
		},
		{
			name: "缺少Name",
			folder: &social.CollectionFolder{
				UserID: "user123",
			},
			wantErr: "用户ID和收藏夹名称不能为空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := repo.CreateFolder(ctx, tt.folder)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

// TestCollectionRepository_GetFolderByID 测试根据ID获取收藏夹
func TestCollectionRepository_GetFolderByID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	folder := &social.CollectionFolder{
		UserID:      "user_getfolderbyid_" + primitive.NewObjectID().Hex(),
		Name:        "测试收藏夹",
		Description: "测试用",
		IsPublic:    false,
	}
	err := repo.CreateFolder(ctx, folder)
	require.NoError(t, err)

	// Act
	found, err := repo.GetFolderByID(ctx, folder.ID.Hex())

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, folder.Name, found.Name)
	assert.Equal(t, folder.Description, found.Description)
}

// TestCollectionRepository_GetFolderByID_NotFound 测试获取不存在的收藏夹
func TestCollectionRepository_GetFolderByID_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	fakeID := primitive.NewObjectID().Hex()

	// Act
	found, err := repo.GetFolderByID(ctx, fakeID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, found)
	assert.Contains(t, err.Error(), "not found")
}

// TestCollectionRepository_GetFoldersByUser 测试获取用户收藏夹列表
func TestCollectionRepository_GetFoldersByUser(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	userID := "user_getfolders_" + primitive.NewObjectID().Hex()
	for i := 0; i < 3; i++ {
		folder := &social.CollectionFolder{
			UserID:      userID,
			Name:        "收藏夹" + string(rune('1'+i)),
			Description: "测试收藏夹",
			IsPublic:    false,
		}
		err := repo.CreateFolder(ctx, folder)
		require.NoError(t, err)
	}

	// Act
	folders, err := repo.GetFoldersByUser(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, folders)
	assert.GreaterOrEqual(t, len(folders), 3)
}

// TestCollectionRepository_UpdateFolder 测试更新收藏夹
func TestCollectionRepository_UpdateFolder(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	folder := &social.CollectionFolder{
		UserID:      "user_updatefolder_" + primitive.NewObjectID().Hex(),
		Name:        "原始名称",
		Description: "原始描述",
		IsPublic:    false,
	}
	err := repo.CreateFolder(ctx, folder)
	require.NoError(t, err)

	// Act - 更新收藏夹
	updates := map[string]interface{}{
		"name":        "更新后的名称",
		"description": "更新后的描述",
		"is_public":   true,
	}
	err = repo.UpdateFolder(ctx, folder.ID.Hex(), updates)

	// Assert
	require.NoError(t, err)

	// 验证更新
	found, err := repo.GetFolderByID(ctx, folder.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, "更新后的名称", found.Name)
	assert.Equal(t, "更新后的描述", found.Description)
	assert.True(t, found.IsPublic)
}

// TestCollectionRepository_UpdateFolder_NotFound 测试更新不存在的收藏夹
func TestCollectionRepository_UpdateFolder_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	fakeID := primitive.NewObjectID().Hex()
	updates := map[string]interface{}{
		"name": "更新后的名称",
	}

	// Act
	err := repo.UpdateFolder(ctx, fakeID, updates)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestCollectionRepository_DeleteFolder 测试删除收藏夹
func TestCollectionRepository_DeleteFolder(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	folder := &social.CollectionFolder{
		UserID:      "user_deletefolder_" + primitive.NewObjectID().Hex(),
		Name:        "待删除收藏夹",
		Description: "测试用",
		IsPublic:    false,
	}
	err := repo.CreateFolder(ctx, folder)
	require.NoError(t, err)

	// Act - 删除收藏夹
	err = repo.DeleteFolder(ctx, folder.ID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证已删除
	found, err := repo.GetFolderByID(ctx, folder.ID.Hex())
	assert.Error(t, err)
	assert.Nil(t, found)
}

// TestCollectionRepository_DeleteFolder_NotFound 测试删除不存在的收藏夹
func TestCollectionRepository_DeleteFolder_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	fakeID := primitive.NewObjectID().Hex()

	// Act
	err := repo.DeleteFolder(ctx, fakeID)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestCollectionRepository_IncrementFolderBookCount 测试增加收藏夹书籍数量
func TestCollectionRepository_IncrementFolderBookCount(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	folder := &social.CollectionFolder{
		UserID:      "user_increment_" + primitive.NewObjectID().Hex(),
		Name:        "测试收藏夹",
		Description: "测试用",
		IsPublic:    false,
	}
	err := repo.CreateFolder(ctx, folder)
	require.NoError(t, err)

	// Act - 增加书籍数量
	err = repo.IncrementFolderBookCount(ctx, folder.ID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证数量已增加
	found, err := repo.GetFolderByID(ctx, folder.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, 1, found.BookCount)
}

// TestCollectionRepository_DecrementFolderBookCount 测试减少收藏夹书籍数量
func TestCollectionRepository_DecrementFolderBookCount(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	folder := &social.CollectionFolder{
		UserID:      "user_decrement_" + primitive.NewObjectID().Hex(),
		Name:        "测试收藏夹",
		Description: "测试用",
		IsPublic:    false,
	}
	err := repo.CreateFolder(ctx, folder)
	require.NoError(t, err)

	// 先增加书籍数量
	err = repo.IncrementFolderBookCount(ctx, folder.ID.Hex())
	require.NoError(t, err)

	// Act - 减少书籍数量
	err = repo.DecrementFolderBookCount(ctx, folder.ID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证数量已减少
	found, err := repo.GetFolderByID(ctx, folder.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, 0, found.BookCount)
}

// TestCollectionRepository_IncrementFolderBookCount_EmptyFolderID 测试空收藏夹ID
func TestCollectionRepository_IncrementFolderBookCount_EmptyFolderID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	// Act - 空收藏夹ID不应该报错
	err := repo.IncrementFolderBookCount(ctx, "")

	// Assert
	require.NoError(t, err)
}

// TestCollectionRepository_DecrementFolderBookCount_EmptyFolderID 测试空收藏夹ID
func TestCollectionRepository_DecrementFolderBookCount_EmptyFolderID(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	// Act - 空收藏夹ID不应该报错
	err := repo.DecrementFolderBookCount(ctx, "")

	// Assert
	require.NoError(t, err)
}

// =========================
// 公开收藏测试
// =========================

// TestCollectionRepository_GetPublicCollections 测试获取公开收藏列表
func TestCollectionRepository_GetPublicCollections(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	// 创建公开收藏
	for i := 0; i < 3; i++ {
		collection := &social.Collection{
			UserID:   primitive.NewObjectID().Hex(),
			BookID:   primitive.NewObjectID().Hex(),
			IsPublic: true,
		}
		err := repo.Create(ctx, collection)
		require.NoError(t, err)
	}

	// Act
	collections, total, err := repo.GetPublicCollections(ctx, 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, collections)
	assert.GreaterOrEqual(t, total, int64(3))
	// 验证都是公开收藏
	for _, col := range collections {
		assert.True(t, col.IsPublic)
	}
}

// TestCollectionRepository_GetPublicFolders 测试获取公开收藏夹列表
func TestCollectionRepository_GetPublicFolders(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	// 创建公开收藏夹
	for i := 0; i < 3; i++ {
		folder := &social.CollectionFolder{
			UserID:      primitive.NewObjectID().Hex(),
			Name:        "公开收藏夹" + string(rune('1'+i)),
			Description: "测试用",
			IsPublic:    true,
		}
		err := repo.CreateFolder(ctx, folder)
		require.NoError(t, err)
	}

	// Act
	folders, total, err := repo.GetPublicFolders(ctx, 1, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, folders)
	assert.GreaterOrEqual(t, total, int64(3))
	// 验证都是公开收藏夹
	for _, folder := range folders {
		assert.True(t, folder.IsPublic)
	}
}

// =========================
// 统计测试
// =========================

// TestCollectionRepository_CountUserCollections 测试统计用户收藏数
func TestCollectionRepository_CountUserCollections(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	userID := "user_count_" + primitive.NewObjectID().Hex()
	for i := 0; i < 5; i++ {
		collection := &social.Collection{
			UserID:   userID,
			BookID:   primitive.NewObjectID().Hex(),
			IsPublic: false,
		}
		err := repo.Create(ctx, collection)
		require.NoError(t, err)
	}

	// Act
	count, err := repo.CountUserCollections(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(5))
}

// TestCollectionRepository_CountUserCollections_Zero 测试统计无收藏的用户
func TestCollectionRepository_CountUserCollections_Zero(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	// Act
	count, err := repo.CountUserCollections(ctx, "nonexistent_user")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

// =========================
// 健康检查测试
// =========================

// TestCollectionRepository_Health 测试健康检查
func TestCollectionRepository_Health(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	// Act
	err := repo.Health(ctx)

	// Assert
	assert.NoError(t, err)
}

// =========================
// 综合场景测试
// =========================

// TestCollectionRepository_CompleteWorkflow 测试完整的收藏工作流
func TestCollectionRepository_CompleteWorkflow(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	userID := "user_workflow_" + primitive.NewObjectID().Hex()

	// 1. 创建收藏夹
	folder := &social.CollectionFolder{
		UserID:      userID,
		Name:        "我的玄幻收藏",
		Description: "收藏玄幻小说",
		IsPublic:    false,
	}
	err := repo.CreateFolder(ctx, folder)
	require.NoError(t, err)

	// 2. 添加收藏到收藏夹
	bookID := primitive.NewObjectID().Hex()
	collection := &social.Collection{
		UserID:   userID,
		BookID:   bookID,
		FolderID: folder.ID.Hex(),
		Tags:     []string{"玄幻", "推荐"},
		Note:     "很好看",
		IsPublic: false,
	}
	err = repo.Create(ctx, collection)
	require.NoError(t, err)

	// 3. 增加收藏夹书籍数量
	err = repo.IncrementFolderBookCount(ctx, folder.ID.Hex())
	require.NoError(t, err)

	// 4. 验证收藏夹书籍数量
	folderAfter, err := repo.GetFolderByID(ctx, folder.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, 1, folderAfter.BookCount)

	// 5. 获取用户收藏列表
	collections, total, err := repo.GetCollectionsByUser(ctx, userID, "", 1, 10)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	assert.Len(t, collections, 1)

	// 6. 统计用户收藏数
	count, err := repo.CountUserCollections(ctx, userID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(1))

	// 7. 取消收藏
	err = repo.Delete(ctx, collection.ID.Hex())
	require.NoError(t, err)

	// 8. 减少收藏夹书籍数量
	err = repo.DecrementFolderBookCount(ctx, folder.ID.Hex())
	require.NoError(t, err)

	// 9. 验证收藏已删除
	_, err = repo.GetByID(ctx, collection.ID.Hex())
	assert.Error(t, err)

	// 10. 验证收藏夹书籍数量已减少
	folderFinal, err := repo.GetFolderByID(ctx, folder.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, 0, folderFinal.BookCount)
}

// TestCollectionRepository_AutoSetTimestamp 测试自动设置时间戳
func TestCollectionRepository_AutoSetTimestamp(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupCollectionRepo(t)
	defer cleanup()

	collection := &social.Collection{
		UserID:   "user_timestamp_" + primitive.NewObjectID().Hex(),
		BookID:   "book_timestamp_" + primitive.NewObjectID().Hex(),
		IsPublic: false,
	}
	beforeCreate := time.Now()

	// Act
	err := repo.Create(ctx, collection)

	// Assert
	require.NoError(t, err)
	assert.False(t, collection.CreatedAt.IsZero())
	assert.False(t, collection.UpdatedAt.IsZero())
	assert.WithinDuration(t, time.Now(), collection.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), collection.UpdatedAt, time.Second)
	assert.True(t, collection.CreatedAt.After(beforeCreate) || collection.CreatedAt.Equal(beforeCreate))
}
