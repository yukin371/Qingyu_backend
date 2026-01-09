package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/reader"
	readerRepo "Qingyu_backend/repository/mongodb/reader"
)

// TestCollectionRepository 收藏Repository测试
func TestCollectionRepository(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := readerRepo.NewMongoCollectionRepository(db)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("Create_Success", func(t *testing.T) {
		collection := &reader.Collection{
			UserID:    testUserID,
			BookID:    testBookID,
			Note:      "这本书太好了！",
			Tags:      []string{"玄幻", "推荐"},
			IsPublic:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, collection)
		assert.NoError(t, err)
		assert.False(t, collection.ID.IsZero(), "收藏ID应该被设置")

		t.Logf("✓ 创建收藏成功，ID: %s", collection.ID.Hex())
	})

	t.Run("Create_Duplicate", func(t *testing.T) {
		collection := &reader.Collection{
			UserID:    testUserID,
			BookID:    testBookID,
			Note:      "重复收藏",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, collection)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "已经收藏")

		t.Logf("✓ 重复收藏检测通过")
	})

	t.Run("GetByID_Success", func(t *testing.T) {
		// 创建收藏
		collection := &reader.Collection{
			UserID:    testUserID,
			BookID:    primitive.NewObjectID().Hex(),
			Note:      "测试获取",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, collection)
		assert.NoError(t, err)

		// 获取收藏
		found, err := repo.GetByID(ctx, collection.ID.Hex())
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, collection.Note, found.Note)

		t.Logf("✓ 获取收藏成功")
	})

	t.Run("GetByUserAndBook_Success", func(t *testing.T) {
		found, err := repo.GetByUserAndBook(ctx, testUserID, testBookID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, testUserID, found.UserID)
		assert.Equal(t, testBookID, found.BookID)

		t.Logf("✓ 根据用户和书籍获取收藏成功")
	})

	t.Run("GetByUserAndBook_NotFound", func(t *testing.T) {
		fakeBookID := primitive.NewObjectID().Hex()
		found, err := repo.GetByUserAndBook(ctx, testUserID, fakeBookID)
		assert.NoError(t, err)
		assert.Nil(t, found, "未收藏应返回nil")

		t.Logf("✓ 未收藏正确返回nil")
	})

	t.Run("GetCollectionsByUser_WithPagination", func(t *testing.T) {
		// 创建多个收藏
		for i := 0; i < 5; i++ {
			collection := &reader.Collection{
				UserID:    testUserID,
				BookID:    primitive.NewObjectID().Hex(),
				Note:      "测试收藏",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err := repo.Create(ctx, collection)
			assert.NoError(t, err)
		}

		// 查询
		collections, total, err := repo.GetCollectionsByUser(ctx, testUserID, "", 1, 3)
		assert.NoError(t, err)
		assert.Greater(t, total, int64(0))
		assert.LessOrEqual(t, len(collections), 3)

		t.Logf("✓ 分页查询收藏成功，总数: %d, 本页: %d", total, len(collections))
	})

	t.Run("Update_Success", func(t *testing.T) {
		// 创建收藏
		collection := &reader.Collection{
			UserID:    testUserID,
			BookID:    primitive.NewObjectID().Hex(),
			Note:      "原始笔记",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, collection)
		assert.NoError(t, err)

		// 更新收藏
		updates := map[string]interface{}{
			"note": "更新后的笔记",
			"tags": []string{"新标签"},
		}
		err = repo.Update(ctx, collection.ID.Hex(), updates)
		assert.NoError(t, err)

		// 验证更新
		found, err := repo.GetByID(ctx, collection.ID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, "更新后的笔记", found.Note)

		t.Logf("✓ 更新收藏成功")
	})

	t.Run("Delete_Success", func(t *testing.T) {
		// 创建收藏
		collection := &reader.Collection{
			UserID:    testUserID,
			BookID:    primitive.NewObjectID().Hex(),
			Note:      "待删除",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, collection)
		assert.NoError(t, err)

		// 删除收藏
		err = repo.Delete(ctx, collection.ID.Hex())
		assert.NoError(t, err)

		// 验证已删除
		found, err := repo.GetByID(ctx, collection.ID.Hex())
		assert.Error(t, err)
		assert.Nil(t, found)

		t.Logf("✓ 删除收藏成功")
	})

	t.Run("CountUserCollections_Success", func(t *testing.T) {
		count, err := repo.CountUserCollections(ctx, testUserID)
		assert.NoError(t, err)
		assert.Greater(t, count, int64(0))

		t.Logf("✓ 统计用户收藏数成功: %d", count)
	})

	t.Run("Health_Success", func(t *testing.T) {
		err := repo.Health(ctx)
		assert.NoError(t, err)

		t.Logf("✓ 健康检查通过")
	})
}

// TestCollectionRepositoryTags 标签查询测试
func TestCollectionRepositoryTags(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := readerRepo.NewMongoCollectionRepository(db)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("GetCollectionsByTag_Success", func(t *testing.T) {
		// 创建带标签的收藏
		tags := [][]string{
			{"玄幻", "热血"},
			{"玄幻", "推荐"},
			{"都市", "言情"},
		}

		for _, tagSet := range tags {
			collection := &reader.Collection{
				UserID:    testUserID,
				BookID:    primitive.NewObjectID().Hex(),
				Tags:      tagSet,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err := repo.Create(ctx, collection)
			assert.NoError(t, err)
		}

		// 按标签查询
		_, total, err := repo.GetCollectionsByTag(ctx, testUserID, "玄幻", 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total, "应该有2条玄幻标签的收藏")

		t.Logf("✓ 按标签查询成功，找到%d条", total)
	})
}

// TestCollectionFolderRepository 收藏夹Repository测试
func TestCollectionFolderRepository(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := readerRepo.NewMongoCollectionRepository(db)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("CreateFolder_Success", func(t *testing.T) {
		folder := &reader.CollectionFolder{
			UserID:      testUserID,
			Name:        "我的最爱",
			Description: "收藏的经典作品",
			IsPublic:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err := repo.CreateFolder(ctx, folder)
		assert.NoError(t, err)
		assert.False(t, folder.ID.IsZero())
		assert.Equal(t, 0, folder.BookCount, "初始书籍数应该为0")

		t.Logf("✓ 创建收藏夹成功，ID: %s", folder.ID.Hex())
	})

	t.Run("GetFolderByID_Success", func(t *testing.T) {
		// 创建收藏夹
		folder := &reader.CollectionFolder{
			UserID:    testUserID,
			Name:      "测试收藏夹",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.CreateFolder(ctx, folder)
		assert.NoError(t, err)

		// 获取收藏夹
		found, err := repo.GetFolderByID(ctx, folder.ID.Hex())
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, folder.Name, found.Name)

		t.Logf("✓ 获取收藏夹成功")
	})

	t.Run("GetFoldersByUser_Success", func(t *testing.T) {
		folders, err := repo.GetFoldersByUser(ctx, testUserID)
		assert.NoError(t, err)
		assert.Greater(t, len(folders), 0)

		t.Logf("✓ 获取用户收藏夹列表成功，数量: %d", len(folders))
	})

	t.Run("UpdateFolder_Success", func(t *testing.T) {
		// 创建收藏夹
		folder := &reader.CollectionFolder{
			UserID:    testUserID,
			Name:      "原始名称",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.CreateFolder(ctx, folder)
		assert.NoError(t, err)

		// 更新收藏夹
		updates := map[string]interface{}{
			"name":        "新名称",
			"description": "新描述",
		}
		err = repo.UpdateFolder(ctx, folder.ID.Hex(), updates)
		assert.NoError(t, err)

		// 验证更新
		found, err := repo.GetFolderByID(ctx, folder.ID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, "新名称", found.Name)

		t.Logf("✓ 更新收藏夹成功")
	})

	t.Run("IncrementFolderBookCount_Success", func(t *testing.T) {
		// 创建收藏夹
		folder := &reader.CollectionFolder{
			UserID:    testUserID,
			Name:      "计数测试",
			BookCount: 0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.CreateFolder(ctx, folder)
		assert.NoError(t, err)

		// 增加计数
		err = repo.IncrementFolderBookCount(ctx, folder.ID.Hex())
		assert.NoError(t, err)

		// 验证计数
		found, err := repo.GetFolderByID(ctx, folder.ID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, 1, found.BookCount)

		t.Logf("✓ 增加收藏夹书籍计数成功")
	})

	t.Run("DecrementFolderBookCount_Success", func(t *testing.T) {
		// 创建收藏夹
		folder := &reader.CollectionFolder{
			UserID:    testUserID,
			Name:      "减计数测试",
			BookCount: 5,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.CreateFolder(ctx, folder)
		assert.NoError(t, err)

		// 减少计数
		err = repo.DecrementFolderBookCount(ctx, folder.ID.Hex())
		assert.NoError(t, err)

		// 验证计数
		found, err := repo.GetFolderByID(ctx, folder.ID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, 4, found.BookCount)

		t.Logf("✓ 减少收藏夹书籍计数成功")
	})

	t.Run("DeleteFolder_Success", func(t *testing.T) {
		// 创建收藏夹
		folder := &reader.CollectionFolder{
			UserID:    testUserID,
			Name:      "待删除",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.CreateFolder(ctx, folder)
		assert.NoError(t, err)

		// 删除收藏夹
		err = repo.DeleteFolder(ctx, folder.ID.Hex())
		assert.NoError(t, err)

		// 验证已删除
		found, err := repo.GetFolderByID(ctx, folder.ID.Hex())
		assert.Error(t, err)
		assert.Nil(t, found)

		t.Logf("✓ 删除收藏夹成功")
	})
}

// TestCollectionRepositoryPublic 公开收藏测试
func TestCollectionRepositoryPublic(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := readerRepo.NewMongoCollectionRepository(db)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("GetPublicCollections_Success", func(t *testing.T) {
		// 创建公开收藏
		for i := 0; i < 3; i++ {
			collection := &reader.Collection{
				UserID:    testUserID,
				BookID:    primitive.NewObjectID().Hex(),
				IsPublic:  true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err := repo.Create(ctx, collection)
			assert.NoError(t, err)
		}

		// 查询公开收藏
		_, total, err := repo.GetPublicCollections(ctx, 1, 10)
		assert.NoError(t, err)
		assert.Greater(t, total, int64(0))

		t.Logf("✓ 获取公开收藏成功，数量: %d", total)
	})

	t.Run("GetPublicFolders_Success", func(t *testing.T) {
		// 创建公开收藏夹
		folder := &reader.CollectionFolder{
			UserID:    testUserID,
			Name:      "公开收藏夹",
			IsPublic:  true,
			BookCount: 10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.CreateFolder(ctx, folder)
		assert.NoError(t, err)

		// 查询公开收藏夹
		_, total, err := repo.GetPublicFolders(ctx, 1, 10)
		assert.NoError(t, err)
		assert.Greater(t, total, int64(0))

		t.Logf("✓ 获取公开收藏夹成功，数量: %d", total)
	})
}

// TestCollectionRepositoryFolderQuery 收藏夹过滤查询测试
func TestCollectionRepositoryFolderQuery(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := readerRepo.NewMongoCollectionRepository(db)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	// 创建收藏夹
	folder := &reader.CollectionFolder{
		UserID:    testUserID,
		Name:      "测试收藏夹",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.CreateFolder(ctx, folder)
	assert.NoError(t, err)

	t.Run("GetCollectionsByFolder_Success", func(t *testing.T) {
		// 在收藏夹中添加收藏
		for i := 0; i < 3; i++ {
			collection := &reader.Collection{
				UserID:    testUserID,
				BookID:    primitive.NewObjectID().Hex(),
				FolderID:  folder.ID.Hex(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err := repo.Create(ctx, collection)
			assert.NoError(t, err)
		}

		// 按收藏夹查询
		_, total, err := repo.GetCollectionsByUser(ctx, testUserID, folder.ID.Hex(), 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total, "应该有3条收藏在该收藏夹中")

		t.Logf("✓ 按收藏夹查询成功，找到%d条", total)
	})
}
