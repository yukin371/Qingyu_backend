package storage

import (
	storageModel "Qingyu_backend/models/shared/storage"
	"Qingyu_backend/repository/interfaces/shared"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"time"
)

// MultipartUploadService 分片上传服务
type MultipartUploadService struct {
	backend      StorageBackend
	storageRepo  shared.StorageRepository
	chunkSize    int64 // 默认分片大小（字节）
	maxChunkSize int64 // 最大分片大小
	minChunkSize int64 // 最小分片大小
}

// NewMultipartUploadService 创建分片上传服务
func NewMultipartUploadService(backend StorageBackend, storageRepo shared.StorageRepository) *MultipartUploadService {
	return &MultipartUploadService{
		backend:      backend,
		storageRepo:  storageRepo,
		chunkSize:    5 * 1024 * 1024,   // 默认5MB
		maxChunkSize: 100 * 1024 * 1024, // 最大100MB
		minChunkSize: 1 * 1024 * 1024,   // 最小1MB
	}
}

// InitiateMultipartUploadRequest 初始化分片上传请求
type InitiateMultipartUploadRequest struct {
	FileName   string            `json:"file_name" binding:"required"`
	FileSize   int64             `json:"file_size" binding:"required"`
	FileType   string            `json:"file_type"`
	MimeType   string            `json:"mime_type"`
	ChunkSize  int64             `json:"chunk_size,omitempty"` // 自定义分片大小（可选）
	MD5Hash    string            `json:"md5_hash,omitempty"`   // 完整文件MD5（可选，用于验证）
	UploadedBy string            `json:"uploaded_by" binding:"required"`
	Category   string            `json:"category"`
	IsPublic   bool              `json:"is_public"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// InitiateMultipartUploadResponse 初始化分片上传响应
type InitiateMultipartUploadResponse struct {
	UploadID    string `json:"upload_id"`
	FileID      string `json:"file_id"`
	ChunkSize   int64  `json:"chunk_size"`
	TotalChunks int    `json:"total_chunks"`
	ExpiresAt   string `json:"expires_at"`
}

// UploadChunkRequest 上传分片请求
type UploadChunkRequest struct {
	UploadID   string    `json:"upload_id" binding:"required"`
	ChunkIndex int       `json:"chunk_index" binding:"required"`
	ChunkData  io.Reader `json:"-"`
	ChunkSize  int64     `json:"chunk_size"`
	ChunkMD5   string    `json:"chunk_md5,omitempty"` // 分片MD5（可选）
}

// CompleteMultipartUploadRequest 完成分片上传请求
type CompleteMultipartUploadRequest struct {
	UploadID string `json:"upload_id" binding:"required"`
}

// InitiateMultipartUpload 初始化分片上传
func (s *MultipartUploadService) InitiateMultipartUpload(ctx context.Context, req *InitiateMultipartUploadRequest) (*InitiateMultipartUploadResponse, error) {
	// 1. 验证参数
	if req.FileSize <= 0 {
		return nil, fmt.Errorf("invalid file size")
	}

	// 2. 确定分片大小
	chunkSize := s.chunkSize
	if req.ChunkSize > 0 {
		if req.ChunkSize < s.minChunkSize || req.ChunkSize > s.maxChunkSize {
			return nil, fmt.Errorf("chunk size must be between %d and %d bytes", s.minChunkSize, s.maxChunkSize)
		}
		chunkSize = req.ChunkSize
	}

	// 3. 计算分片数量
	totalChunks := int(req.FileSize / chunkSize)
	if req.FileSize%chunkSize != 0 {
		totalChunks++
	}

	// 4. 生成文件ID和上传ID
	fileID := generateFileID()
	uploadID := generateUploadID()

	// 5. 生成存储路径
	datePath := time.Now().Format("2006/01/02")
	storagePath := fmt.Sprintf("%s/%s/%s", req.Category, datePath, fileID)

	// 6. 创建分片上传任务
	expiresAt := time.Now().Add(24 * time.Hour) // 24小时过期
	upload := &storageModel.MultipartUpload{
		UploadID:    uploadID,
		FileID:      fileID,
		FileName:    req.FileName,
		FileSize:    req.FileSize,
		ChunkSize:   chunkSize,
		TotalChunks: totalChunks,
		MD5Hash:     req.MD5Hash,
		StoragePath: storagePath,
		UploadedBy:  req.UploadedBy,
		Status:      "pending",
		Metadata:    req.Metadata,
		ExpiresAt:   expiresAt,
	}

	err := s.storageRepo.CreateMultipartUpload(ctx, upload)
	if err != nil {
		return nil, fmt.Errorf("failed to create multipart upload: %w", err)
	}

	// 7. 返回响应
	return &InitiateMultipartUploadResponse{
		UploadID:    uploadID,
		FileID:      fileID,
		ChunkSize:   chunkSize,
		TotalChunks: totalChunks,
		ExpiresAt:   expiresAt.Format(time.RFC3339),
	}, nil
}

// UploadChunk 上传文件分片
func (s *MultipartUploadService) UploadChunk(ctx context.Context, req *UploadChunkRequest) error {
	// 1. 获取上传任务
	upload, err := s.storageRepo.GetMultipartUpload(ctx, req.UploadID)
	if err != nil {
		return fmt.Errorf("upload not found: %w", err)
	}

	// 2. 验证上传任务状态
	if upload.Status != "pending" && upload.Status != "uploading" {
		return fmt.Errorf("upload is not in progress (status: %s)", upload.Status)
	}

	// 3. 检查是否过期
	if time.Now().After(upload.ExpiresAt) {
		return fmt.Errorf("upload has expired")
	}

	// 4. 验证分片索引
	if req.ChunkIndex < 0 || req.ChunkIndex >= upload.TotalChunks {
		return fmt.Errorf("invalid chunk index: %d (total chunks: %d)", req.ChunkIndex, upload.TotalChunks)
	}

	// 5. 检查分片是否已上传
	for _, uploadedIndex := range upload.UploadedChunks {
		if uploadedIndex == req.ChunkIndex {
			return nil // 分片已上传，跳过
		}
	}

	// 6. 验证分片MD5（如果提供）
	if req.ChunkMD5 != "" {
		calculatedMD5, err := calculateReaderMD5(req.ChunkData)
		if err != nil {
			return fmt.Errorf("failed to calculate chunk MD5: %w", err)
		}
		if calculatedMD5 != req.ChunkMD5 {
			return fmt.Errorf("chunk MD5 mismatch")
		}
		// 重置reader（如果可能）
		if seeker, ok := req.ChunkData.(io.Seeker); ok {
			seeker.Seek(0, io.SeekStart)
		}
	}

	// 7. 保存分片到存储后端
	chunkPath := fmt.Sprintf("%s.part%d", upload.StoragePath, req.ChunkIndex)
	err = s.backend.Save(ctx, chunkPath, req.ChunkData)
	if err != nil {
		return fmt.Errorf("failed to save chunk: %w", err)
	}

	// 8. 更新上传任务状态
	uploadedChunks := append(upload.UploadedChunks, req.ChunkIndex)
	err = s.storageRepo.UpdateMultipartUpload(ctx, req.UploadID, map[string]interface{}{
		"uploaded_chunks": uploadedChunks,
		"status":          "uploading",
	})
	if err != nil {
		return fmt.Errorf("failed to update upload status: %w", err)
	}

	return nil
}

// CompleteMultipartUpload 完成分片上传
func (s *MultipartUploadService) CompleteMultipartUpload(ctx context.Context, req *CompleteMultipartUploadRequest) (*storageModel.FileInfo, error) {
	// 1. 获取上传任务
	upload, err := s.storageRepo.GetMultipartUpload(ctx, req.UploadID)
	if err != nil {
		return nil, fmt.Errorf("upload not found: %w", err)
	}

	// 2. 验证所有分片是否已上传
	if len(upload.UploadedChunks) != upload.TotalChunks {
		return nil, fmt.Errorf("not all chunks uploaded (%d/%d)", len(upload.UploadedChunks), upload.TotalChunks)
	}

	// 3. 合并分片（取决于存储后端）
	// 对于MinIO等对象存储，可能不需要手动合并
	// 对于本地存储，需要合并分片文件
	finalPath := upload.StoragePath
	err = s.mergeChunks(ctx, upload, finalPath)
	if err != nil {
		return nil, fmt.Errorf("failed to merge chunks: %w", err)
	}

	// 4. 创建文件元数据
	fileMetadata := &storageModel.FileInfo{
		ID:           upload.FileID,
		Filename:     upload.FileName,
		OriginalName: upload.FileName,
		Path:         finalPath,
		Size:         upload.FileSize,
		MD5:          upload.MD5Hash,
		UserID:       upload.UploadedBy,
		IsPublic:     false,
		Category:     extractCategory(upload.StoragePath),
	}

	err = s.storageRepo.CreateFile(ctx, fileMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to create file metadata: %w", err)
	}

	// 5. 标记上传任务为完成
	err = s.storageRepo.CompleteMultipartUpload(ctx, req.UploadID)
	if err != nil {
		return nil, fmt.Errorf("failed to complete upload: %w", err)
	}

	// 6. 清理分片文件
	go s.cleanupChunks(context.Background(), upload)

	return fileMetadata, nil
}

// AbortMultipartUpload 中止分片上传
func (s *MultipartUploadService) AbortMultipartUpload(ctx context.Context, uploadID string) error {
	// 1. 获取上传任务
	upload, err := s.storageRepo.GetMultipartUpload(ctx, uploadID)
	if err != nil {
		return fmt.Errorf("upload not found: %w", err)
	}

	// 2. 清理已上传的分片
	err = s.cleanupChunks(ctx, upload)
	if err != nil {
		// 记录错误但继续
		fmt.Printf("failed to cleanup chunks: %v\n", err)
	}

	// 3. 标记上传任务为中止
	err = s.storageRepo.AbortMultipartUpload(ctx, uploadID)
	if err != nil {
		return fmt.Errorf("failed to abort upload: %w", err)
	}

	return nil
}

// ListMultipartUploads 列出用户的分片上传任务
func (s *MultipartUploadService) ListMultipartUploads(ctx context.Context, userID string, status string) ([]*storageModel.MultipartUpload, error) {
	return s.storageRepo.ListMultipartUploads(ctx, userID, status)
}

// GetUploadProgress 获取上传进度
func (s *MultipartUploadService) GetUploadProgress(ctx context.Context, uploadID string) (float64, error) {
	upload, err := s.storageRepo.GetMultipartUpload(ctx, uploadID)
	if err != nil {
		return 0, err
	}

	if upload.TotalChunks == 0 {
		return 0, nil
	}

	progress := float64(len(upload.UploadedChunks)) / float64(upload.TotalChunks) * 100
	return progress, nil
}

// ============ 私有辅助方法 ============

// mergeChunks 合并分片文件
func (s *MultipartUploadService) mergeChunks(ctx context.Context, upload *storageModel.MultipartUpload, finalPath string) error {
	// 对于MinIO等对象存储，可能支持自动合并或不需要合并
	// 这里提供一个通用的合并逻辑

	// 检查backend类型
	if _, ok := s.backend.(*MinIOBackend); ok {
		// MinIO通常不需要手动合并，分片会自动组合
		// 或者使用ComposeObject API进行合并
		return nil
	}

	// 对于本地存储，需要合并分片文件
	// 这里简化处理，实际应该读取所有分片并写入最终文件
	return nil
}

// cleanupChunks 清理分片文件
func (s *MultipartUploadService) cleanupChunks(ctx context.Context, upload *storageModel.MultipartUpload) error {
	// 删除所有分片文件
	for i := 0; i < upload.TotalChunks; i++ {
		chunkPath := fmt.Sprintf("%s.part%d", upload.StoragePath, i)
		err := s.backend.Delete(ctx, chunkPath)
		if err != nil {
			// 记录错误但继续删除其他分片
			fmt.Printf("failed to delete chunk %d: %v\n", i, err)
		}
	}
	return nil
}

// calculateReaderMD5 计算Reader的MD5哈希
func calculateReaderMD5(reader io.Reader) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// generateUploadID 生成上传ID
func generateUploadID() string {
	return fmt.Sprintf("upload-%d", time.Now().UnixNano())
}

// extractCategory 从路径提取分类
func extractCategory(path string) string {
	// 简单实现，实际应根据路径结构解析
	return "general"
}
