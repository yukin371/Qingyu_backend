package storage

import (
	storageModel "Qingyu_backend/models/storage"
	"Qingyu_backend/repository/interfaces/shared"
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	defaultChunkSize    int64 = 5 * 1024 * 1024
	maxChunkSize        int64 = 100 * 1024 * 1024
	minChunkSize        int64 = 1 * 1024 * 1024
	uploadExpirationTTL       = 24 * time.Hour

	statusPending   = "pending"
	statusUploading = "uploading"
	statusCompleted = "completed"
	statusAborted   = "aborted"

	defaultCategory = "general"
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
		chunkSize:    defaultChunkSize,
		maxChunkSize: maxChunkSize,
		minChunkSize: minChunkSize,
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
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}
	if req.FileName == "" {
		return nil, fmt.Errorf("file name is required")
	}
	if req.UploadedBy == "" {
		return nil, fmt.Errorf("uploaded_by is required")
	}

	// 1. 验证参数
	if req.FileSize <= 0 {
		return nil, fmt.Errorf("invalid file size")
	}

	// 2. 确定分片大小
	chunkSize, err := s.resolveChunkSize(req.ChunkSize)
	if err != nil {
		return nil, err
	}

	// 3. 计算分片数量
	totalChunks := calculateTotalChunks(req.FileSize, chunkSize)

	// 4. 生成文件ID和上传ID
	fileID := generateFileID()
	uploadID := generateUploadID()

	// 5. 生成存储路径
	datePath := time.Now().Format("2006/01/02")
	storagePath := fmt.Sprintf("%s/%s/%s", normalizeCategory(req.Category), datePath, fileID)

	// 6. 创建分片上传任务
	expiresAt := time.Now().Add(uploadExpirationTTL)
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
		Status:      statusPending,
		Metadata:    req.Metadata,
		ExpiresAt:   expiresAt,
	}

	if err := s.storageRepo.CreateMultipartUpload(ctx, upload); err != nil {
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
	if err := ctx.Err(); err != nil {
		return err
	}
	if req == nil {
		return fmt.Errorf("request is required")
	}
	if req.UploadID == "" {
		return fmt.Errorf("upload_id is required")
	}
	if req.ChunkData == nil {
		return fmt.Errorf("chunk data is required")
	}

	// 1. 获取上传任务
	upload, err := s.getUploadByID(ctx, req.UploadID)
	if err != nil {
		return err
	}

	// 2. 验证上传任务状态
	if !isUploadInProgress(upload.Status) {
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
	if containsChunk(upload.UploadedChunks, req.ChunkIndex) {
		return nil // 分片已上传，跳过
	}

	chunkReader, err := s.buildChunkReader(req)
	if err != nil {
		return err
	}

	// 7. 保存分片到存储后端
	chunkPath := buildChunkPath(upload.StoragePath, req.ChunkIndex)
	if err = s.backend.Save(ctx, chunkPath, chunkReader); err != nil {
		return fmt.Errorf("failed to save chunk: %w", err)
	}

	// 8. 更新上传任务状态
	uploadedChunks := append(upload.UploadedChunks, req.ChunkIndex)
	err = s.storageRepo.UpdateMultipartUpload(ctx, req.UploadID, map[string]interface{}{
		"uploaded_chunks": uploadedChunks,
		"status":          statusUploading,
	})
	if err != nil {
		return fmt.Errorf("failed to update upload status: %w", err)
	}

	return nil
}

// CompleteMultipartUpload 完成分片上传
func (s *MultipartUploadService) CompleteMultipartUpload(ctx context.Context, req *CompleteMultipartUploadRequest) (*storageModel.FileInfo, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}
	if req.UploadID == "" {
		return nil, fmt.Errorf("upload_id is required")
	}

	// 1. 获取上传任务
	upload, err := s.getUploadByID(ctx, req.UploadID)
	if err != nil {
		return nil, err
	}
	if upload.Status == statusCompleted {
		fileInfo, getErr := s.storageRepo.GetFile(ctx, upload.FileID)
		if getErr != nil {
			return nil, fmt.Errorf("upload already completed and file metadata missing: %w", getErr)
		}
		return fileInfo, nil
	}
	if upload.Status == statusAborted {
		return nil, fmt.Errorf("upload is aborted")
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
	if err := ctx.Err(); err != nil {
		return err
	}
	if uploadID == "" {
		return fmt.Errorf("upload_id is required")
	}

	// 1. 获取上传任务
	upload, err := s.getUploadByID(ctx, uploadID)
	if err != nil {
		return err
	}
	if upload.Status == statusCompleted {
		return fmt.Errorf("cannot abort completed upload")
	}

	// 2. 清理已上传的分片
	cleanupErr := s.cleanupChunks(ctx, upload)

	// 3. 标记上传任务为中止
	err = s.storageRepo.AbortMultipartUpload(ctx, uploadID)
	if err != nil {
		return fmt.Errorf("failed to abort upload: %w", err)
	}
	if cleanupErr != nil {
		return fmt.Errorf("upload aborted with cleanup errors: %w", cleanupErr)
	}

	return nil
}

// ListMultipartUploads 列出用户的分片上传任务
func (s *MultipartUploadService) ListMultipartUploads(ctx context.Context, userID string, status string) ([]*storageModel.MultipartUpload, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return s.storageRepo.ListMultipartUploads(ctx, userID, status)
}

// GetUploadProgress 获取上传进度
func (s *MultipartUploadService) GetUploadProgress(ctx context.Context, uploadID string) (float64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	if uploadID == "" {
		return 0, fmt.Errorf("upload_id is required")
	}
	upload, err := s.getUploadByID(ctx, uploadID)
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
	var cleanupErr error

	// 删除所有分片文件
	for i := 0; i < upload.TotalChunks; i++ {
		if ctx != nil {
			if err := ctx.Err(); err != nil {
				return errors.Join(cleanupErr, err)
			}
		}
		chunkPath := buildChunkPath(upload.StoragePath, i)
		err := s.backend.Delete(ctx, chunkPath)
		if err != nil {
			cleanupErr = errors.Join(cleanupErr, fmt.Errorf("delete chunk %d: %w", i, err))
		}
	}
	return cleanupErr
}

func (s *MultipartUploadService) getUploadByID(ctx context.Context, uploadID string) (*storageModel.MultipartUpload, error) {
	upload, err := s.storageRepo.GetMultipartUpload(ctx, uploadID)
	if err != nil {
		return nil, fmt.Errorf("upload not found: %w", err)
	}
	return upload, nil
}

// generateUploadID 生成上传ID
func generateUploadID() string {
	return fmt.Sprintf("upload-%d", time.Now().UnixNano())
}

// extractCategory 从路径提取分类
func extractCategory(path string) string {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return defaultCategory
	}
	parts := strings.Split(trimmed, "/")
	if len(parts) == 0 || parts[0] == "" {
		return defaultCategory
	}
	return parts[0]
}

func (s *MultipartUploadService) resolveChunkSize(requestChunkSize int64) (int64, error) {
	if requestChunkSize == 0 {
		return s.chunkSize, nil
	}
	if requestChunkSize < s.minChunkSize || requestChunkSize > s.maxChunkSize {
		return 0, fmt.Errorf("chunk size must be between %d and %d bytes", s.minChunkSize, s.maxChunkSize)
	}
	return requestChunkSize, nil
}

func (s *MultipartUploadService) buildChunkReader(req *UploadChunkRequest) (io.Reader, error) {
	if req.ChunkMD5 == "" {
		return req.ChunkData, nil
	}

	data, err := io.ReadAll(req.ChunkData)
	if err != nil {
		return nil, fmt.Errorf("failed to read chunk data: %w", err)
	}
	calculatedMD5 := calculateBytesMD5(data)
	if calculatedMD5 != req.ChunkMD5 {
		return nil, fmt.Errorf("chunk MD5 mismatch")
	}

	return bytes.NewReader(data), nil
}

func isUploadInProgress(status string) bool {
	return status == statusPending || status == statusUploading
}

func containsChunk(chunks []int, chunkIndex int) bool {
	for _, uploadedIndex := range chunks {
		if uploadedIndex == chunkIndex {
			return true
		}
	}
	return false
}

func calculateTotalChunks(fileSize int64, chunkSize int64) int {
	totalChunks := int(fileSize / chunkSize)
	if fileSize%chunkSize != 0 {
		totalChunks++
	}
	return totalChunks
}

func normalizeCategory(category string) string {
	if strings.TrimSpace(category) == "" {
		return defaultCategory
	}
	return category
}

func buildChunkPath(basePath string, chunkIndex int) string {
	return fmt.Sprintf("%s.part%d", basePath, chunkIndex)
}

func calculateBytesMD5(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}
