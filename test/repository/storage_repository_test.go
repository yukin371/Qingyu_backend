package repository

import (
	"Qingyu_backend/repository/interfaces/shared"
	"Qingyu_backend/repository/mongodb"
	"Qingyu_backend/test/testutil"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStorageRepository_CRUD 测试基础CRUD操作
func TestStorageRepository_CRUD(t *testing.T) {
	// 1. 连接测试数据库
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	// 2. 创建Repository
	repo := mongodb.NewMongoStorageRepository(db)

	// 3. 测试创建文件
	file := &shared.FileMetadata{
		FileName:     "test.jpg",
		OriginalName: "original_test.jpg",
		FilePath:     "/uploads/test.jpg",
		FileSize:     1024,
		FileType:     "image",
		MimeType:     "image/jpeg",
		Extension:    ".jpg",
		MD5Hash:      "abc123",
		StorageType:  "local",
		StoragePath:  "/uploads/test.jpg",
		Status:       "active",
		IsPublic:     false,
		Category:     "avatar",
		UploadedBy:   "user_123",
	}

	err := repo.Create(context.Background(), file)
	require.NoError(t, err)
	assert.NotEmpty(t, file.ID)
	assert.NotZero(t, file.CreatedAt)

	// 4. 测试根据ID查询
	retrieved, err := repo.GetByID(context.Background(), file.ID)
	require.NoError(t, err)
	assert.Equal(t, file.FileName, retrieved.FileName)
	assert.Equal(t, file.MD5Hash, retrieved.MD5Hash)

	// 5. 测试根据MD5查询
	byMD5, err := repo.GetByMD5(context.Background(), file.MD5Hash)
	require.NoError(t, err)
	assert.Equal(t, file.ID, byMD5.ID)

	// 6. 测试更新
	updates := map[string]interface{}{
		"status": "archived",
	}
	err = repo.Update(context.Background(), file.ID, updates)
	require.NoError(t, err)

	updated, err := repo.GetByID(context.Background(), file.ID)
	require.NoError(t, err)
	assert.Equal(t, "archived", updated.Status)

	// 7. 测试删除（软删除）
	err = repo.Delete(context.Background(), file.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(context.Background(), file.ID)
	assert.Error(t, err) // 软删除后应该找不到
}

// TestStorageRepository_List 测试列表查询
func TestStorageRepository_List(t *testing.T) {
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoStorageRepository(db)

	// 创建多个测试文件
	userID := "user_123"
	for i := 0; i < 5; i++ {
		file := &shared.FileMetadata{
			FileName:     "test.jpg",
			OriginalName: "test.jpg",
			FilePath:     "/uploads/test.jpg",
			FileSize:     int64(1024 * (i + 1)),
			FileType:     "image",
			MimeType:     "image/jpeg",
			Extension:    ".jpg",
			StorageType:  "local",
			StoragePath:  "/uploads/test.jpg",
			Status:       "active",
			IsPublic:     true,
			Category:     "avatar",
			UploadedBy:   userID,
		}
		err := repo.Create(context.Background(), file)
		require.NoError(t, err)
	}

	// 测试列表查询
	filter := &shared.FileFilter{
		UserID:   userID,
		Status:   "active",
		Page:     1,
		PageSize: 10,
	}

	files, total, err := repo.List(context.Background(), filter)
	require.NoError(t, err)
	assert.Len(t, files, 5)
	assert.Equal(t, int64(5), total)
}

// TestStorageRepository_Permissions 测试权限管理
func TestStorageRepository_Permissions(t *testing.T) {
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoStorageRepository(db)

	// 创建测试文件
	file := &shared.FileMetadata{
		FileName:     "test.jpg",
		OriginalName: "test.jpg",
		FilePath:     "/uploads/test.jpg",
		FileSize:     1024,
		FileType:     "image",
		MimeType:     "image/jpeg",
		Extension:    ".jpg",
		StorageType:  "local",
		StoragePath:  "/uploads/test.jpg",
		Status:       "active",
		IsPublic:     false,
		Category:     "avatar",
		UploadedBy:   "user_123",
	}
	err := repo.Create(context.Background(), file)
	require.NoError(t, err)

	// 测试授予权限
	err = repo.GrantAccess(context.Background(), file.ID, "user_456", "read")
	require.NoError(t, err)

	// 测试检查权限
	hasAccess, err := repo.CheckAccess(context.Background(), file.ID, "user_456")
	require.NoError(t, err)
	assert.True(t, hasAccess)

	// 测试撤销权限
	err = repo.RevokeAccess(context.Background(), file.ID, "user_456")
	require.NoError(t, err)

	hasAccess, err = repo.CheckAccess(context.Background(), file.ID, "user_456")
	require.NoError(t, err)
	assert.False(t, hasAccess)
}

// TestStorageRepository_MultipartUpload 测试分片上传
func TestStorageRepository_MultipartUpload(t *testing.T) {
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoStorageRepository(db)

	// 创建分片上传任务
	upload := &shared.MultipartUpload{
		FileName:    "large_file.mp4",
		FileSize:    100 * 1024 * 1024, // 100MB
		ChunkSize:   5 * 1024 * 1024,   // 5MB
		TotalChunks: 20,
		StoragePath: "/uploads/large_file.mp4",
		UploadedBy:  "user_123",
		Metadata: map[string]string{
			"content_type": "video/mp4",
		},
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err := repo.CreateMultipartUpload(context.Background(), upload)
	require.NoError(t, err)
	assert.NotEmpty(t, upload.ID)
	assert.NotEmpty(t, upload.UploadID)
	assert.Equal(t, "pending", upload.Status)

	// 测试获取上传任务
	retrieved, err := repo.GetMultipartUpload(context.Background(), upload.UploadID)
	require.NoError(t, err)
	assert.Equal(t, upload.UploadID, retrieved.UploadID)
	assert.Equal(t, 20, retrieved.TotalChunks)

	// 测试更新上传任务（模拟上传分片）
	err = repo.UpdateMultipartUpload(context.Background(), upload.UploadID, map[string]interface{}{
		"uploaded_chunks": []int{0, 1, 2},
		"status":          "uploading",
	})
	require.NoError(t, err)

	updated, err := repo.GetMultipartUpload(context.Background(), upload.UploadID)
	require.NoError(t, err)
	assert.Equal(t, "uploading", updated.Status)
	assert.Len(t, updated.UploadedChunks, 3)

	// 测试完成上传
	err = repo.CompleteMultipartUpload(context.Background(), upload.UploadID)
	require.NoError(t, err)

	completed, err := repo.GetMultipartUpload(context.Background(), upload.UploadID)
	require.NoError(t, err)
	assert.Equal(t, "completed", completed.Status)
}

// TestStorageRepository_Stats 测试统计功能
func TestStorageRepository_Stats(t *testing.T) {
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoStorageRepository(db)

	userID := "user_123"

	// 创建不同类型的文件
	fileTypes := []string{"image", "video", "audio", "document"}
	for i, fileType := range fileTypes {
		file := &shared.FileMetadata{
			FileName:     "test.jpg",
			OriginalName: "test.jpg",
			FilePath:     "/uploads/test.jpg",
			FileSize:     int64(1024 * (i + 1)),
			FileType:     fileType,
			MimeType:     "application/octet-stream",
			Extension:    ".bin",
			StorageType:  "local",
			StoragePath:  "/uploads/test.jpg",
			Status:       "active",
			IsPublic:     false,
			Category:     "attachment",
			UploadedBy:   userID,
		}
		err := repo.Create(context.Background(), file)
		require.NoError(t, err)
	}

	// 获取统计信息
	stats, err := repo.GetStorageStats(context.Background(), userID)
	require.NoError(t, err)
	assert.Equal(t, int64(4), stats.TotalFiles)
	assert.Greater(t, stats.TotalSize, int64(0))
	assert.Equal(t, int64(1), stats.ImageCount)
	assert.Equal(t, int64(1), stats.VideoCount)
	assert.Equal(t, int64(1), stats.AudioCount)
	assert.Equal(t, int64(1), stats.DocumentCount)
}

// TestStorageRepository_Health 测试健康检查
func TestStorageRepository_Health(t *testing.T) {
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoStorageRepository(db)

	err := repo.Health(context.Background())
	assert.NoError(t, err)
}
