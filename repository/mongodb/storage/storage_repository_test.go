package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	storageModel "Qingyu_backend/models/storage"
	"Qingyu_backend/repository/interfaces/shared"
	storageRepo "Qingyu_backend/repository/mongodb/storage"
	"Qingyu_backend/test/testutil"
)

// 测试辅助函数
func setupStorageRepo(t *testing.T) (shared.StorageRepository, context.Context, func()) {
	db, cleanup := testutil.SetupTestDB(t)
	repo := storageRepo.NewMongoStorageRepository(db)
	ctx := context.Background()
	return repo, ctx, cleanup
}

// TestStorageRepository_CreateFile 测试创建文件记录
func TestStorageRepository_CreateFile(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupStorageRepo(t)
	defer cleanup()

	file := &storageModel.FileInfo{
		Filename:     "test.jpg",
		OriginalName: "original_test.jpg",
		ContentType:  "image/jpeg",
		Size:         1024,
		Path:         "/uploads/test.jpg",
		UserID:       "user_123",
		IsPublic:     false,
		Category:     "avatar",
		MD5:          primitive.NewObjectID().Hex(),
	}

	// Act
	err := repo.CreateFile(ctx, file)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, file.ID)
	assert.NotZero(t, file.CreatedAt)
	assert.Equal(t, "test.jpg", file.Filename)
}

// TestStorageRepository_CreateFile_NilFile 测试创建空文件
func TestStorageRepository_CreateFile_NilFile(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupStorageRepo(t)
	defer cleanup()

	// Act
	err := repo.CreateFile(ctx, nil)

	// Assert
	assert.Error(t, err)
}

// TestStorageRepository_GetFile 测试根据ID获取文件
func TestStorageRepository_GetFile(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupStorageRepo(t)
	defer cleanup()

	file := &storageModel.FileInfo{
		Filename:     "test.jpg",
		OriginalName: "test.jpg",
		ContentType:  "image/jpeg",
		Size:         1024,
		Path:         "/uploads/test.jpg",
		UserID:       "user_123",
		IsPublic:     true,
		Category:     "avatar",
		MD5:          primitive.NewObjectID().Hex(),
	}
	err := repo.CreateFile(ctx, file)
	require.NoError(t, err)

	// Act
	retrieved, err := repo.GetFile(ctx, file.ID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, file.ID, retrieved.ID)
	assert.Equal(t, file.Filename, retrieved.Filename)
	assert.Equal(t, file.MD5, retrieved.MD5)
}

// TestStorageRepository_GetFile_NotFound 测试获取不存在的文件
func TestStorageRepository_GetFile_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupStorageRepo(t)
	defer cleanup()

	// Act
	retrieved, err := repo.GetFile(ctx, "nonexistent_id")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, retrieved)
}

// TestStorageRepository_GetFileByMD5 测试根据MD5获取文件
func TestStorageRepository_GetFileByMD5(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupStorageRepo(t)
	defer cleanup()

	md5Hash := primitive.NewObjectID().Hex()
	file := &storageModel.FileInfo{
		Filename:     "test.jpg",
		OriginalName: "test.jpg",
		ContentType:  "image/jpeg",
		Size:         1024,
		Path:         "/uploads/test.jpg",
		UserID:       "user_123",
		MD5:          md5Hash,
	}
	err := repo.CreateFile(ctx, file)
	require.NoError(t, err)

	// Act
	byMD5, err := repo.GetFileByMD5(ctx, md5Hash)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, file.ID, byMD5.ID)
	assert.Equal(t, md5Hash, byMD5.MD5)
}

// TestStorageRepository_GetFileByMD5_NotFound 测试根据不存在的MD5获取文件
func TestStorageRepository_GetFileByMD5_NotFound(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupStorageRepo(t)
	defer cleanup()

	// Act
	byMD5, err := repo.GetFileByMD5(ctx, "nonexistent_md5")

	// Assert
	require.NoError(t, err)
	assert.Nil(t, byMD5)
}

// TestStorageRepository_UpdateFile 测试更新文件信息
func TestStorageRepository_UpdateFile(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupStorageRepo(t)
	defer cleanup()

	file := &storageModel.FileInfo{
		Filename:     "test.jpg",
		OriginalName: "test.jpg",
		ContentType:  "image/jpeg",
		Size:         1024,
		Path:         "/uploads/test.jpg",
		UserID:       "user_123",
		IsPublic:     false,
		Category:     "avatar",
		MD5:          primitive.NewObjectID().Hex(),
	}
	err := repo.CreateFile(ctx, file)
	require.NoError(t, err)

	// Act
	updates := map[string]interface{}{
		"is_public": true,
		"filename":  "updated_test.jpg",
	}
	err = repo.UpdateFile(ctx, file.ID, updates)

	// Assert
	require.NoError(t, err)

	updated, err := repo.GetFile(ctx, file.ID)
	require.NoError(t, err)
	assert.True(t, updated.IsPublic)
	assert.Equal(t, "updated_test.jpg", updated.Filename)
}

// TestStorageRepository_DeleteFile 测试删除文件
func TestStorageRepository_DeleteFile(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupStorageRepo(t)
	defer cleanup()

	file := &storageModel.FileInfo{
		Filename:     "test.jpg",
		OriginalName: "test.jpg",
		ContentType:  "image/jpeg",
		Size:         1024,
		Path:         "/uploads/test.jpg",
		UserID:       "user_123",
		MD5:          primitive.NewObjectID().Hex(),
	}
	err := repo.CreateFile(ctx, file)
	require.NoError(t, err)

	// Act
	err = repo.DeleteFile(ctx, file.ID)

	// Assert
	require.NoError(t, err)
}

// TestStorageRepository_ListFiles 测试列表查询
func TestStorageRepository_ListFiles(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupStorageRepo(t)
	defer cleanup()

	userID := "user_list_" + primitive.NewObjectID().Hex()
	for i := 0; i < 5; i++ {
		file := &storageModel.FileInfo{
			Filename:     primitive.NewObjectID().Hex() + ".jpg",
			OriginalName: "test.jpg",
			ContentType:  "image/jpeg",
			Size:         int64(1024 * (i + 1)),
			Path:         "/uploads/test.jpg",
			UserID:       userID,
			IsPublic:     true,
			Category:     "avatar",
		}
		err := repo.CreateFile(ctx, file)
		require.NoError(t, err)
	}

	// Act
	filter := &shared.FileFilter{
		UserID:   userID,
		Page:     1,
		PageSize: 10,
	}

	files, total, err := repo.ListFiles(ctx, filter)

	// Assert
	require.NoError(t, err)
	assert.Len(t, files, 5)
	assert.Equal(t, int64(5), total)
}

// TestStorageRepository_ListFiles_Pagination 测试分页查询
func TestStorageRepository_ListFiles_Pagination(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupStorageRepo(t)
	defer cleanup()

	userID := "user_page_" + primitive.NewObjectID().Hex()
	for i := 0; i < 15; i++ {
		file := &storageModel.FileInfo{
			Filename:     primitive.NewObjectID().Hex() + ".jpg",
			OriginalName: "test.jpg",
			ContentType:  "image/jpeg",
			Size:         1024,
			Path:         "/uploads/test.jpg",
			UserID:       userID,
			IsPublic:     true,
			Category:     "avatar",
		}
		err := repo.CreateFile(ctx, file)
		require.NoError(t, err)
	}

	// Act - 第一页
	filter1 := &shared.FileFilter{
		UserID:   userID,
		Page:     1,
		PageSize: 10,
	}
	files1, total1, err := repo.ListFiles(ctx, filter1)

	// Assert
	require.NoError(t, err)
	assert.Len(t, files1, 10)
	assert.Equal(t, int64(15), total1)

	// Act - 第二页
	filter2 := &shared.FileFilter{
		UserID:   userID,
		Page:     2,
		PageSize: 10,
	}
	files2, total2, err := repo.ListFiles(ctx, filter2)

	// Assert
	require.NoError(t, err)
	assert.Len(t, files2, 5)
	assert.Equal(t, int64(15), total2)
}

// TestStorageRepository_GrantAccess 测试授予访问权限
func TestStorageRepository_GrantAccess(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupStorageRepo(t)
	defer cleanup()

	file := &storageModel.FileInfo{
		Filename:     "test.jpg",
		OriginalName: "test.jpg",
		ContentType:  "image/jpeg",
		Size:         1024,
		Path:         "/uploads/test.jpg",
		UserID:       "user_123",
		IsPublic:     false,
		Category:     "avatar",
	}
	err := repo.CreateFile(ctx, file)
	require.NoError(t, err)

	// Act - 授予访问权限
	err = repo.GrantAccess(ctx, file.ID, "user_456")

	// Assert
	require.NoError(t, err)
}

// TestStorageRepository_CheckAccess 测试检查访问权限
func TestStorageRepository_CheckAccess(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupStorageRepo(t)
	defer cleanup()

	file := &storageModel.FileInfo{
		Filename:     "test.jpg",
		OriginalName: "test.jpg",
		ContentType:  "image/jpeg",
		Size:         1024,
		Path:         "/uploads/test.jpg",
		UserID:       "user_123",
		IsPublic:     false,
		Category:     "avatar",
	}
	err := repo.CreateFile(ctx, file)
	require.NoError(t, err)

	err = repo.GrantAccess(ctx, file.ID, "user_456")
	require.NoError(t, err)

	// Act
	hasAccess, err := repo.CheckAccess(ctx, file.ID, "user_456")

	// Assert
	require.NoError(t, err)
	assert.True(t, hasAccess)
}

// TestStorageRepository_CheckAccess_NoAccess 测试无访问权限
func TestStorageRepository_CheckAccess_NoAccess(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupStorageRepo(t)
	defer cleanup()

	file := &storageModel.FileInfo{
		Filename:     "test.jpg",
		OriginalName: "test.jpg",
		ContentType:  "image/jpeg",
		Size:         1024,
		Path:         "/uploads/test.jpg",
		UserID:       "user_123",
		IsPublic:     false,
		Category:     "avatar",
	}
	err := repo.CreateFile(ctx, file)
	require.NoError(t, err)

	// Act
	hasAccess, err := repo.CheckAccess(ctx, file.ID, "user_789")

	// Assert
	require.NoError(t, err)
	assert.False(t, hasAccess)
}

// TestStorageRepository_RevokeAccess 测试撤销访问权限
func TestStorageRepository_RevokeAccess(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupStorageRepo(t)
	defer cleanup()

	file := &storageModel.FileInfo{
		Filename:     "test.jpg",
		OriginalName: "test.jpg",
		ContentType:  "image/jpeg",
		Size:         1024,
		Path:         "/uploads/test.jpg",
		UserID:       "user_123",
		IsPublic:     false,
		Category:     "avatar",
	}
	err := repo.CreateFile(ctx, file)
	require.NoError(t, err)

	err = repo.GrantAccess(ctx, file.ID, "user_456")
	require.NoError(t, err)

	// Act
	err = repo.RevokeAccess(ctx, file.ID, "user_456")

	// Assert
	require.NoError(t, err)

	hasAccess, err := repo.CheckAccess(ctx, file.ID, "user_456")
	require.NoError(t, err)
	assert.False(t, hasAccess)
}

// TestStorageRepository_CreateMultipartUpload 测试创建分片上传任务
func TestStorageRepository_CreateMultipartUpload(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupStorageRepo(t)
	defer cleanup()

	upload := &storageModel.MultipartUpload{
		FileName:    "large_file.mp4",
		FileSize:    100 * 1024 * 1024,
		ChunkSize:   5 * 1024 * 1024,
		TotalChunks: 20,
		StoragePath: "/uploads/large_file.mp4",
		UploadedBy:  "user_123",
		Metadata: map[string]string{
			"content_type": "video/mp4",
		},
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// Act
	err := repo.CreateMultipartUpload(ctx, upload)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, upload.ID)
	assert.NotEmpty(t, upload.UploadID)
	assert.Equal(t, "pending", upload.Status)
}

// TestStorageRepository_GetMultipartUpload 测试获取分片上传任务
func TestStorageRepository_GetMultipartUpload(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupStorageRepo(t)
	defer cleanup()

	upload := &storageModel.MultipartUpload{
		FileName:    "large_file.mp4",
		FileSize:    100 * 1024 * 1024,
		ChunkSize:   5 * 1024 * 1024,
		TotalChunks: 20,
		StoragePath: "/uploads/large_file.mp4",
		UploadedBy:  "user_123",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}
	err := repo.CreateMultipartUpload(ctx, upload)
	require.NoError(t, err)

	// Act
	retrieved, err := repo.GetMultipartUpload(ctx, upload.UploadID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, upload.UploadID, retrieved.UploadID)
	assert.Equal(t, 20, retrieved.TotalChunks)
}

// TestStorageRepository_UpdateMultipartUpload 测试更新分片上传任务
func TestStorageRepository_UpdateMultipartUpload(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupStorageRepo(t)
	defer cleanup()

	upload := &storageModel.MultipartUpload{
		FileName:    "large_file.mp4",
		FileSize:    100 * 1024 * 1024,
		ChunkSize:   5 * 1024 * 1024,
		TotalChunks: 20,
		StoragePath: "/uploads/large_file.mp4",
		UploadedBy:  "user_123",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}
	err := repo.CreateMultipartUpload(ctx, upload)
	require.NoError(t, err)

	// Act
	err = repo.UpdateMultipartUpload(ctx, upload.UploadID, map[string]interface{}{
		"uploaded_chunks": []int{0, 1, 2},
		"status":          "uploading",
	})

	// Assert
	require.NoError(t, err)

	updated, err := repo.GetMultipartUpload(ctx, upload.UploadID)
	require.NoError(t, err)
	assert.Equal(t, "uploading", updated.Status)
	assert.Len(t, updated.UploadedChunks, 3)
}

// TestStorageRepository_CompleteMultipartUpload 测试完成分片上传
func TestStorageRepository_CompleteMultipartUpload(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupStorageRepo(t)
	defer cleanup()

	upload := &storageModel.MultipartUpload{
		FileName:    "large_file.mp4",
		FileSize:    100 * 1024 * 1024,
		ChunkSize:   5 * 1024 * 1024,
		TotalChunks: 20,
		StoragePath: "/uploads/large_file.mp4",
		UploadedBy:  "user_123",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}
	err := repo.CreateMultipartUpload(ctx, upload)
	require.NoError(t, err)

	// Act
	err = repo.CompleteMultipartUpload(ctx, upload.UploadID)

	// Assert
	require.NoError(t, err)

	completed, err := repo.GetMultipartUpload(ctx, upload.UploadID)
	require.NoError(t, err)
	assert.Equal(t, "completed", completed.Status)
}

// TestStorageRepository_IncrementDownloadCount 测试增加下载计数
func TestStorageRepository_IncrementDownloadCount(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupStorageRepo(t)
	defer cleanup()

	file := &storageModel.FileInfo{
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
	err := repo.CreateFile(ctx, file)
	require.NoError(t, err)

	// Act
	err = repo.IncrementDownloadCount(ctx, file.ID)

	// Assert
	require.NoError(t, err)

	updated, err := repo.GetFile(ctx, file.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), updated.Downloads)
}

// TestStorageRepository_Health 测试健康检查
func TestStorageRepository_Health(t *testing.T) {
	// Arrange
	repo, ctx, cleanup := setupStorageRepo(t)
	defer cleanup()

	// Act
	err := repo.Health(ctx)

	// Assert
	assert.NoError(t, err)
}
