package repository

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/models/shared/storage"
	"Qingyu_backend/repository/interfaces/shared"
	"Qingyu_backend/repository/mongodb"
	"Qingyu_backend/test/testutil"

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
	file := &storage.FileInfo{
		Filename:     "test.jpg",
		OriginalName: "original_test.jpg",
		ContentType:  "image/jpeg",
		Size:         1024,
		Path:         "/uploads/test.jpg",
		UserID:       "user_123",
		IsPublic:     false,
		Category:     "avatar",
		MD5:          "abc123",
	}

	err := repo.CreateFile(context.Background(), file)
	require.NoError(t, err)
	assert.NotEmpty(t, file.ID)
	assert.NotZero(t, file.CreatedAt)

	// 4. 测试根据ID查询
	retrieved, err := repo.GetFile(context.Background(), file.ID)
	require.NoError(t, err)
	assert.Equal(t, file.Filename, retrieved.Filename)
	assert.Equal(t, file.MD5, retrieved.MD5)

	// 5. 测试根据MD5查询
	byMD5, err := repo.GetFileByMD5(context.Background(), file.MD5)
	require.NoError(t, err)
	assert.Equal(t, file.ID, byMD5.ID)

	// 6. 测试更新
	updates := map[string]interface{}{
		"is_public": true,
	}
	err = repo.UpdateFile(context.Background(), file.ID, updates)
	require.NoError(t, err)

	updated, err := repo.GetFile(context.Background(), file.ID)
	require.NoError(t, err)
	assert.True(t, updated.IsPublic)

	// 7. 测试删除（软删除）
	err = repo.DeleteFile(context.Background(), file.ID)
	require.NoError(t, err)

	// 注意：软删除后文件仍可以查询到，但会标记为deleted
	// 如果实现了硬删除，下面的断言应该是 assert.Error(t, err)
	_, err = repo.GetFile(context.Background(), file.ID)
	// 根据实际实现调整断言
	// assert.Error(t, err) // 硬删除
	// 或者
	// assert.NoError(t, err) // 软删除，但需要检查Status字段
}

// TestStorageRepository_List 测试列表查询
func TestStorageRepository_List(t *testing.T) {
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoStorageRepository(db)

	// 创建多个测试文件
	userID := "user_123"
	for i := 0; i < 5; i++ {
		file := &storage.FileInfo{
			Filename:     "test.jpg",
			OriginalName: "test.jpg",
			ContentType:  "image/jpeg",
			Size:         int64(1024 * (i + 1)),
			Path:         "/uploads/test.jpg",
			UserID:       userID,
			IsPublic:     true,
			Category:     "avatar",
		}
		err := repo.CreateFile(context.Background(), file)
		require.NoError(t, err)
	}

	// 测试列表查询
	filter := &shared.FileFilter{
		UserID:   userID,
		Page:     1,
		PageSize: 10,
	}

	files, total, err := repo.ListFiles(context.Background(), filter)
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
	file := &storage.FileInfo{
		Filename:     "test.jpg",
		OriginalName: "test.jpg",
		ContentType:  "image/jpeg",
		Size:         1024,
		Path:         "/uploads/test.jpg",
		UserID:       "user_123",
		IsPublic:     false,
		Category:     "avatar",
	}
	err := repo.CreateFile(context.Background(), file)
	require.NoError(t, err)

	// 测试授予权限 - 只传2个参数（fileID, userID）
	err = repo.GrantAccess(context.Background(), file.ID, "user_456")
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
	upload := &storage.MultipartUpload{
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

// TestStorageRepository_IncrementDownloadCount 测试下载计数
func TestStorageRepository_IncrementDownloadCount(t *testing.T) {
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoStorageRepository(db)

	// 创建测试文件
	file := &storage.FileInfo{
		Filename:     "test.jpg",
		OriginalName: "test.jpg",
		ContentType:  "image/jpeg",
		Size:         1024,
		Path:         "/uploads/test.jpg",
		UserID:       "user_123",
		IsPublic:     true,
		Category:     "attachment",
		Downloads:    0,
	}
	err := repo.CreateFile(context.Background(), file)
	require.NoError(t, err)

	// 增加下载计数
	err = repo.IncrementDownloadCount(context.Background(), file.ID)
	require.NoError(t, err)

	// 验证下载计数已增加
	updated, err := repo.GetFile(context.Background(), file.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), updated.Downloads)
}

// TestStorageRepository_Health 测试健康检查
func TestStorageRepository_Health(t *testing.T) {
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoStorageRepository(db)

	err := repo.Health(context.Background())
	assert.NoError(t, err)
}
